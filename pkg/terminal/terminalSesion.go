/*
 * Copyright (c) 2024. Devtron Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package terminal

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/devtron-labs/common-lib/async"

	"github.com/caarlos0/env"
	"github.com/devtron-labs/common-lib/utils/k8s"
	"github.com/devtron-labs/devtron/internal/middleware"
	bean3 "github.com/devtron-labs/devtron/pkg/argoApplication/bean"
	"github.com/devtron-labs/devtron/pkg/argoApplication/read/config"
	"github.com/devtron-labs/devtron/pkg/cluster"
	"github.com/devtron-labs/devtron/pkg/cluster/bean"
	"github.com/devtron-labs/devtron/pkg/cluster/environment"
	bean2 "github.com/devtron-labs/devtron/pkg/cluster/environment/bean"
	"github.com/devtron-labs/devtron/pkg/cluster/read"
	"github.com/devtron-labs/devtron/pkg/cluster/repository"
	errors1 "github.com/juju/errors"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"

	"gopkg.in/igm/sockjs-go.v3/sockjs"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

const END_OF_TRANSMISSION = "\u0004"
const ProcessExitedMsg = "Process exited"
const ProcessTimedOut = "Process timedOut"

// PtyHandler is what remotecommand expects from a pty
type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

// TerminalSession implements PtyHandler (using a SockJS connection)
type TerminalSession struct {
	id                string
	bound             chan error
	sockJSSession     sockjs.Session
	sizeChan          chan remotecommand.TerminalSize
	doneChan          chan struct{}
	context           context.Context
	contextCancelFunc context.CancelFunc
	podName           string
	namespace         string
	clusterId         string
	startedOn         time.Time
}

// TerminalMessage is the messaging protocol between ShellController and TerminalSession.
//
// OP      DIRECTION  FIELD(S) USED  DESCRIPTION
// ---------------------------------------------------------------------
// bind    fe->be     SessionID      Id sent back from TerminalResponse
// stdin   fe->be     Data           Keystrokes/paste buffer
// resize  fe->be     Rows, Cols     New terminal size
// stdout  be->fe     Data           Output from the process
// toast   be->fe     Data           OOB message to be shown to the user
type TerminalMessage struct {
	Op, Data, SessionID string
	Rows, Cols          uint16
}

// TerminalSize handles pty->process resize events
// Called in a loop from remotecommand as long as the process is running
func (t TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

// Read handles pty->process messages (stdin, resize)
// Called in a loop from remotecommand as long as the process is running
func (t TerminalSession) Read(p []byte) (int, error) {
	m, err := t.sockJSSession.Recv()
	if err != nil {
		// Send terminated signal to process to avoid resource leak
		return copy(p, END_OF_TRANSMISSION), err
	}

	var msg TerminalMessage
	if err := json.Unmarshal([]byte(m), &msg); err != nil {
		return copy(p, END_OF_TRANSMISSION), err
	}

	switch msg.Op {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	default:
		return copy(p, END_OF_TRANSMISSION), fmt.Errorf("unknown message type '%s'", msg.Op)
	}
}

// Write handles process->pty stdout
// Called from remotecommand whenever there is any output
func (t TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Op:   "stdout",
		Data: string(p),
	})
	if err != nil {
		return 0, err
	}

	if err = t.sockJSSession.Send(string(msg)); err != nil {
		return 0, err
	}
	return len(p), nil
}

// Toast can be used to send the user any OOB messages
// hterm puts these in the center of the terminal
func (t TerminalSession) Toast(p string) error {
	msg, err := json.Marshal(TerminalMessage{
		Op:   "toast",
		Data: p,
	})
	if err != nil {
		return err
	}

	if err = t.sockJSSession.Send(string(msg)); err != nil {
		return err
	}
	return nil
}

// SessionMap stores a map of all TerminalSession objects and a lock to avoid concurrent conflict
type SessionMap struct {
	Sessions map[string]TerminalSession
	Lock     sync.RWMutex
}

// Get return a given terminalSession by sessionId
func (sm *SessionMap) Get(sessionId string) TerminalSession {
	sm.Lock.RLock()
	defer sm.Lock.RUnlock()
	session := sm.Sessions[sessionId]
	return session
}

// Set store a TerminalSession to SessionMap
func (sm *SessionMap) Set(sessionId string, session TerminalSession) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	sm.Sessions[sessionId] = session
}

func (sm *SessionMap) SetTerminalSessionStartTime(sessionId string) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	if session, ok := sm.Sessions[sessionId]; ok {
		session.startedOn = time.Now()
		sm.Sessions[sessionId] = session
	}
}

func (sm *SessionMap) setAndSendSignal(sessionId string, session sockjs.Session) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	terminalSession, ok := sm.Sessions[sessionId]
	if ok && terminalSession.id == "" {
		log.Printf("handleTerminalSession: can't find session '%s'", sessionId)
		session.Close(http.StatusGone, fmt.Sprintf("handleTerminalSession: can't find session '%s'", sessionId))
		return
	} else if ok {
		terminalSession.sockJSSession = session
		sm.Sessions[sessionId] = terminalSession

		select {
		case terminalSession.bound <- nil:
			log.Printf("message sent on bound channel for sessionId : %s", sessionId)
		default:
			// if a request from the front end is not received within a particular time frame, and no one is reading from the bound channel, we will ignore sending on the bound channel.
			log.Printf("skipping send on bound, channel receiver possibly timed out. sessionId: %s", sessionId)
		}

	}
}

// Close shuts down the SockJS connection and sends the status code and reason to the client
// Can happen if the process exits or if there is an error starting up the process
// For now the status code is unused and reason is shown to the user (unless "")
func (sm *SessionMap) Close(sessionId string, status uint32, reason string) {

	sm.Lock.Lock()
	defer sm.Lock.Unlock()

	terminalSession := sm.Sessions[sessionId]

	if terminalSession.sockJSSession != nil {

		err := terminalSession.sockJSSession.Close(status, reason)
		if err != nil {
			log.Println(err)
		}

		close(terminalSession.doneChan)

		isErroredConnectionTermination := isConnectionClosedByError(status)
		middleware.IncTerminalSessionRequestCounter(SessionTerminated, strconv.FormatBool(isErroredConnectionTermination))
		middleware.RecordTerminalSessionDurationMetrics(terminalSession.podName, terminalSession.namespace, terminalSession.clusterId, time.Since(terminalSession.startedOn).Seconds())
		terminalSession.contextCancelFunc()
		close(terminalSession.bound)
		delete(sm.Sessions, sessionId)
	}

}

func isConnectionClosedByError(status uint32) bool {
	if status == 2 {
		return true
	}
	return false
}

var terminalSessions = SessionMap{Sessions: make(map[string]TerminalSession)}

// handleTerminalSession is Called by net/http for any new /api/sockjs connections
func handleTerminalSession(session sockjs.Session) {
	var (
		buf string
		err error
		msg TerminalMessage
	)

	if buf, err = session.Recv(); err != nil {
		log.Printf("handleTerminalSession: can't Recv: %v", err)
		return
	}

	if err = json.Unmarshal([]byte(buf), &msg); err != nil {
		log.Printf("handleTerminalSession: can't UnMarshal (%v): %s", err, buf)
		return
	}

	if msg.Op != "bind" {
		log.Printf("handleTerminalSession: expected 'bind' message, got: %s", buf)
		session.Close(http.StatusBadRequest, fmt.Sprintf("expected 'bind' message, got '%s'", buf))
		return
	}

	terminalSessions.setAndSendSignal(msg.SessionID, session)

}

type SocketConfig struct {
	SocketHeartbeatSeconds int `env:"SOCKET_HEARTBEAT_SECONDS" envDefault:"25" description:"In order to keep proxies and load balancers from closing long running http requests we need to pretend that the connection is active and send a heartbeat packet once in a while. This setting controls how often this is done. By default a heartbeat packet is sent every 25 seconds."`
	SocketDisconnectDelay  int `env:"SOCKET_DISCONNECT_DELAY_SECONDS" envDefault:"5" description:"The server closes a session when a client receiving connection have not been seen for a while.This delay is configured by this setting. By default the session is closed when a receiving connection wasn't seen for 5 seconds."`
}

var cfg *SocketConfig

// CreateAttachHandler is called from main for /api/sockjs
func CreateAttachHandler(path string) http.Handler {
	if cfg == nil {
		cfg = &SocketConfig{}
		env.Parse(cfg)
	}

	opts := sockjs.DefaultOptions
	opts.HeartbeatDelay = time.Duration(cfg.SocketHeartbeatSeconds) * time.Second
	opts.DisconnectDelay = time.Duration(cfg.SocketDisconnectDelay) * time.Second
	return sockjs.NewHandler(path, opts, handleTerminalSession)
}

// startProcess is called by handleAttach
// Executed cmd in the container specified in request and connects it up with the ptyHandler (a session)
func startProcess(ctx context.Context, k8sClient kubernetes.Interface, cfg *rest.Config,
	cmd []string, ptyHandler PtyHandler, sessionRequest *TerminalSessionRequest) error {
	namespace := sessionRequest.Namespace
	podName := sessionRequest.PodName
	containerName := sessionRequest.ContainerName

	exec, err := getExecutor(k8sClient, cfg, podName, namespace, containerName, cmd, true, true)

	if err != nil {
		log.Println("error in getting terminal executor ", "err: ", err)
		return err
	}

	streamOptions := remotecommand.StreamOptions{
		Stdin:             ptyHandler,
		Stdout:            ptyHandler,
		Stderr:            ptyHandler,
		TerminalSizeQueue: ptyHandler,
		Tty:               true,
	}
	isErroredConnectionTermination := false
	middleware.IncTerminalSessionRequestCounter(SessionInitiating, strconv.FormatBool(isErroredConnectionTermination))
	terminalSessions.SetTerminalSessionStartTime(sessionRequest.SessionId)
	err = execWithStreamOptions(ctx, exec, streamOptions)
	if err != nil {
		log.Println("error in terminal exec with stream opts: ", "err: ", err)
		return err
	}
	return nil
}

func execWithStreamOptions(ctx context.Context, exec remotecommand.Executor, streamOptions remotecommand.StreamOptions) error {
	return exec.StreamWithContext(ctx, streamOptions)
}

func getExecutor(k8sClient kubernetes.Interface, cfg *rest.Config, podName, namespace, containerName string, cmd []string, stdin bool, tty bool) (remotecommand.Executor, error) {
	req := k8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	req.VersionedParams(&v1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdin:     stdin,
		Stdout:    true,
		Stderr:    true,
		TTY:       tty,
	}, scheme.ParameterCodec)

	// Use the new fallback executor instead of SPDY directly
	exec, err := NewFallbackExecutor(cfg, "POST", req.URL())
	return exec, err
}

// genTerminalSessionId generates a random session ID string. The format is not really interesting.
// This ID is used to identify the session when the client opens the SockJS connection.
// Not the same as the SockJS session id! We can't use that as that is generated
// on the client side and we don't have it yet at this point.
func genTerminalSessionId() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	id := make([]byte, hex.EncodedLen(len(bytes)))
	hex.Encode(id, bytes)
	return string(id), nil
}

// isValidShell checks if the Shell is an allowed one
func isValidShell(validShells []string, shell string) bool {
	for _, validShell := range validShells {
		if validShell == shell {
			return true
		}
	}
	return false
}

type TerminalSessionRequest struct {
	Shell         string
	SessionId     string
	Namespace     string
	PodName       string
	ContainerName string
	//ApplicationId is helm app Id
	ApplicationId string
	EnvironmentId int
	AppId         int
	//ClusterId is optional
	ClusterId                        int
	UserId                           int32
	ExternalArgoApplicationName      string
	ExternalArgoApplicationNamespace string
	ExternalArgoAppIdentifier        *bean3.ArgoAppIdentifier
}

const CommandExecutionFailed = "Failed to Execute Command"
const PodNotFound = "Pod NotFound"

var validShells = []string{"bash", "sh", "powershell", "cmd"}

// WaitForTerminal is called from apihandler.handleAttach as a goroutine
// Waits for the SockJS connection to be opened by the client the session to be bound in handleTerminalSession
func WaitForTerminal(k8sClient kubernetes.Interface, cfg *rest.Config, request *TerminalSessionRequest) {

	session := terminalSessions.Get(request.SessionId)
	sessionCtx := session.context
	timedCtx, _ := context.WithTimeout(sessionCtx, 60*time.Second)
	select {
	case <-session.bound:

		var err error
		if isValidShell(validShells, request.Shell) {
			cmd := []string{request.Shell}

			err = startProcess(sessionCtx, k8sClient, cfg, cmd, terminalSessions.Get(request.SessionId), request)
		} else {
			// No Shell given or it was not valid: try some shells until one succeeds or all fail
			// FIXME: if the first Shell fails then the first keyboard event is lost
			for _, testShell := range validShells {
				cmd := []string{testShell}
				if err = startProcess(sessionCtx, k8sClient, cfg, cmd, terminalSessions.Get(request.SessionId), request); err == nil || errors2.Is(err, context.Canceled) {
					break
				}
			}
		}

		if err != nil && !errors2.Is(err, context.Canceled) {
			terminalSessions.Close(request.SessionId, 2, err.Error())
			return
		}

		terminalSessions.Close(request.SessionId, 1, ProcessExitedMsg)
	case <-timedCtx.Done():
		// handle case when connection has not been initiated from FE side within particular time
		terminalSessions.Close(request.SessionId, 1, ProcessTimedOut)
	}
}

type TerminalSessionHandler interface {
	GetTerminalSession(req *TerminalSessionRequest) (statusCode int, message *TerminalMessage, err error)
	Close(sessionId string, statusCode uint32, msg string)
	ValidateSession(sessionId string) bool
	ValidateShell(req *TerminalSessionRequest) (bool, error)
	AutoSelectShell(req *TerminalSessionRequest) (string, error)
	RunCmdInRemotePod(req *TerminalSessionRequest, cmds []string) (*bytes.Buffer, *bytes.Buffer, error)
}

type TerminalSessionHandlerImpl struct {
	environmentService           environment.EnvironmentService
	logger                       *zap.SugaredLogger
	k8sUtil                      *k8s.K8sServiceImpl
	ephemeralContainerService    cluster.EphemeralContainerService
	argoApplicationConfigService config.ArgoApplicationConfigService
	ClusterReadService           read.ClusterReadService
	asyncRunnable                *async.Runnable
}

func NewTerminalSessionHandlerImpl(environmentService environment.EnvironmentService,
	logger *zap.SugaredLogger, k8sUtil *k8s.K8sServiceImpl, ephemeralContainerService cluster.EphemeralContainerService,
	argoApplicationConfigService config.ArgoApplicationConfigService,
	ClusterReadService read.ClusterReadService, asyncRunnable *async.Runnable) *TerminalSessionHandlerImpl {
	return &TerminalSessionHandlerImpl{
		environmentService:           environmentService,
		logger:                       logger,
		k8sUtil:                      k8sUtil,
		ephemeralContainerService:    ephemeralContainerService,
		argoApplicationConfigService: argoApplicationConfigService,
		ClusterReadService:           ClusterReadService,
		asyncRunnable:                asyncRunnable,
	}
}

func (impl *TerminalSessionHandlerImpl) Close(sessionId string, statusCode uint32, msg string) {
	terminalSessions.Close(sessionId, statusCode, msg)
}

func (impl *TerminalSessionHandlerImpl) ValidateSession(sessionId string) bool {
	if sessionId == "" {
		return false
	}
	terminalSession := terminalSessions.Get(sessionId)
	sockJSSession := terminalSession.sockJSSession
	if sockJSSession != nil {
		sessionState := sockJSSession.GetSessionState()
		return sessionState == sockjs.SessionActive
	}
	return false
}

func (impl *TerminalSessionHandlerImpl) GetTerminalSession(req *TerminalSessionRequest) (statusCode int, message *TerminalMessage, err error) {
	sessionID, err := genTerminalSessionId()
	if err != nil {
		statusCode := http.StatusInternalServerError
		statusError, ok := err.(*errors.StatusError)
		if ok && statusError.Status().Code > 0 {
			statusCode = int(statusError.Status().Code)
		}
		return statusCode, nil, err
	}
	req.SessionId = sessionID
	sessionCtx, cancelFunc := context.WithCancel(context.Background())
	terminalSessions.Set(sessionID, TerminalSession{
		id:                sessionID,
		bound:             make(chan error),
		sizeChan:          make(chan remotecommand.TerminalSize),
		doneChan:          make(chan struct{}),
		context:           sessionCtx,
		contextCancelFunc: cancelFunc,
		podName:           req.PodName,
		namespace:         req.Namespace,
		clusterId:         strconv.Itoa(req.ClusterId),
	})
	config, client, err := impl.getClientSetAndRestConfigForTerminalConn(req)

	impl.asyncRunnable.Execute(func() {
		err := impl.saveEphemeralContainerTerminalAccessAudit(req)
		if err != nil {
			impl.logger.Errorw("error in saving ephemeral container terminal access audit,so skipping auditing", "err", err)
		}
	})

	if err != nil {
		impl.logger.Errorw("error in fetching config", "err", err)
		return http.StatusInternalServerError, nil, err
	}
	impl.asyncRunnable.Execute(func() { WaitForTerminal(client, config, req) })
	return http.StatusOK, &TerminalMessage{SessionID: sessionID}, nil
}

func (impl *TerminalSessionHandlerImpl) getClientSetAndRestConfigForTerminalConn(req *TerminalSessionRequest) (*rest.Config, *kubernetes.Clientset, error) {
	var clusterBean *bean.ClusterBean
	var clusterConfig *k8s.ClusterConfig
	var restConfig *rest.Config
	var err error
	if req.ExternalArgoAppIdentifier != nil {
		restConfig, err = impl.argoApplicationConfigService.GetRestConfigForExternalArgo(context.Background(), req.ExternalArgoAppIdentifier)
		if err != nil {
			impl.logger.Errorw("error in getting rest config", "err", err, "clusterId", req.ClusterId, "externalArgoApplicationName", req.ExternalArgoApplicationName)
			return nil, nil, err
		}

		_, clientSet, err := impl.k8sUtil.GetK8sConfigAndClientsByRestConfig(restConfig)
		if err != nil {
			impl.logger.Errorw("error in clientSet", "err", err)
			return nil, nil, err
		}
		return restConfig, clientSet, nil
	} else {
		if req.ClusterId != 0 {
			clusterBean, err = impl.ClusterReadService.FindById(req.ClusterId)
			if err != nil {
				impl.logger.Errorw("error in fetching cluster detail", "err", err, "clusterId", req.ClusterId)
				return nil, nil, err
			}
		} else if req.EnvironmentId != 0 {
			clusterBean, err = impl.environmentService.FindClusterByEnvId(req.EnvironmentId)
			if err != nil {
				impl.logger.Errorw("error in fetching cluster detail", "envId", req.EnvironmentId, "err", err)
				return nil, nil, err
			}
		} else {
			return nil, nil, fmt.Errorf("not able to find cluster-config")
		}

		clusterConfig = clusterBean.GetClusterConfig()
		restConfig, err = impl.k8sUtil.GetRestConfigByCluster(clusterConfig, k8s.WithDefaultHttpTransport())
		if err != nil {
			impl.logger.Errorw("error in getting rest config by cluster", "err", err, "clusterName", clusterConfig.ClusterName)
			return nil, nil, err
		}

		_, clientSet, err := impl.k8sUtil.GetK8sConfigAndClientsByRestConfig(restConfig, k8s.WithDefaultHttpTransport())
		if err != nil {
			impl.logger.Errorw("error in clientSet", "err", err)
			return nil, nil, err
		}
		return restConfig, clientSet, nil
	}
}

func (impl *TerminalSessionHandlerImpl) AutoSelectShell(req *TerminalSessionRequest) (string, error) {
	var err1 error
	for _, testShell := range validShells {
		req.Shell = testShell
		isValid, err := impl.ValidateShell(req)
		if isValid {
			return testShell, nil
		}
		if err != nil {
			err1 = err
		}
	}
	if err1 != nil && err1.Error() != CommandExecutionFailed {
		return "", err1
	}
	return "", errors1.New("no shell is supported")
}
func (impl *TerminalSessionHandlerImpl) ValidateShell(req *TerminalSessionRequest) (bool, error) {
	impl.logger.Infow("Inside ValidateShell method in TerminalSessionHandlerImpl", "shellName", req.Shell, "podName", req.PodName, "nameSpace", req.Namespace)

	cmd := fmt.Sprintf("/bin/%s", req.Shell)
	cmdArray := []string{cmd}

	buf, errBuf, err := impl.RunCmdInRemotePod(req, cmdArray)
	if err != nil {
		impl.logger.Errorw("failed to execute commands ", "err", err, "commands", cmdArray, "podName", req.PodName, "namespace", req.Namespace)
		return false, getErrorMsg(err.Error())
	}
	errBufString := errBuf.String()
	if errBufString != "" {
		impl.logger.Errorw("error response on executing commands ", "err", errBufString, "commands", cmdArray, "podName", req.PodName, "namespace", req.Namespace)
		return false, getErrorMsg(errBufString)
	}
	impl.logger.Infow("validated Shell,returning from validateShell method", "StdOut", buf.String())
	return true, nil
}

func getErrorMsg(err string) error {
	if strings.Contains(err, "pods") && strings.Contains(err, "not found") {
		return errors1.New(PodNotFound)
	}
	return errors1.New(CommandExecutionFailed)
}

func (impl *TerminalSessionHandlerImpl) RunCmdInRemotePod(req *TerminalSessionRequest, cmds []string) (*bytes.Buffer, *bytes.Buffer, error) {
	config, client, err := impl.getClientSetAndRestConfigForTerminalConn(req)
	if err != nil {
		impl.logger.Errorw("error in fetching config", "err", err)
		return nil, nil, err
	}
	impl.logger.Debug("reached getExecutor method call")
	exec, err := getExecutor(client, config, req.PodName, req.Namespace, req.ContainerName, cmds, false, false)
	if err != nil {
		impl.logger.Errorw("error occurred in getting remoteCommand executor", "err", err)
		return nil, nil, err
	}
	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	impl.logger.Debug("reached execWithStreamOptions method call")
	err = execWithStreamOptions(context.Background(), exec, remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: errBuf,
	})
	return buf, errBuf, err
}

func (impl *TerminalSessionHandlerImpl) saveEphemeralContainerTerminalAccessAudit(req *TerminalSessionRequest) error {
	var restConfig *rest.Config
	var err error
	if req.ExternalArgoAppIdentifier != nil {
		restConfig, err = impl.argoApplicationConfigService.GetRestConfigForExternalArgo(context.Background(), req.ExternalArgoAppIdentifier)
		if err != nil {
			impl.logger.Errorw("error in getting rest config", "err", err, "clusterId", req.ClusterId, "externalArgoApplicationName", req.ExternalArgoApplicationName)
			return err
		}
	} else {
		clusterBean, err := impl.ClusterReadService.FindById(req.ClusterId)
		if err != nil {
			impl.logger.Errorw("error occurred in finding clusterBean by Id", "clusterId", req.ClusterId, "err", err)
			return err
		}
		clusterConfig := clusterBean.GetClusterConfig()
		restConfig, err = impl.k8sUtil.GetRestConfigByCluster(clusterConfig)
		if err != nil {
			impl.logger.Errorw("error in getting rest config", "err", err, "clusterId", req.ClusterId, "externalArgoApplicationName", req.ExternalArgoApplicationName)
			return err
		}
	}

	v1Client, err := impl.k8sUtil.GetCoreV1ClientByRestConfig(restConfig)
	if err != nil {
		impl.logger.Errorw("error, GetCoreV1ClientByRestConfig", "err", err)
		return err
	}
	pod, err := impl.k8sUtil.GetPodByName(req.Namespace, req.PodName, v1Client)
	if err != nil {
		impl.logger.Errorw("error in getting pod", "clusterId", req.ClusterId, "namespace", req.Namespace, "podName", req.PodName, "err", err)
		return err
	}
	var ephemeralContainer *v1.EphemeralContainer
	for _, ec := range pod.Spec.EphemeralContainers {
		if ec.Name == req.ContainerName {
			ephemeralContainer = &ec
			break
		}
	}
	if ephemeralContainer == nil {
		impl.logger.Infow("terminal session requested for non ephemeral container,so not auditing the terminal access", "clusterId", req.ClusterId, "namespace", req.Namespace, "podName", req.PodName)
		return nil
	}
	ephemeralContainerJson, err := json.Marshal(ephemeralContainer)
	if err != nil {
		impl.logger.Errorw("error occurred while marshaling ephemeralContainer object", "err", err, "ephemeralContainer", ephemeralContainer)
		return err
	}
	ephemeralReq := bean2.EphemeralContainerRequest{
		PodName:   req.PodName,
		Namespace: req.Namespace,
		ClusterId: req.ClusterId,
		BasicData: &bean2.EphemeralContainerBasicData{
			ContainerName:       req.ContainerName,
			TargetContainerName: ephemeralContainer.TargetContainerName,
			Image:               ephemeralContainer.Image,
		},
		AdvancedData: &bean2.EphemeralContainerAdvancedData{
			Manifest: string(ephemeralContainerJson),
		},
		UserId: req.UserId,
	}
	err = impl.ephemeralContainerService.AuditEphemeralContainerAction(ephemeralReq, repository.ActionAccessed)
	if err != nil {
		impl.logger.Errorw("error occurred while requesting ephemeral container terminal access audit", "err", err, "clusterId", req.ClusterId, "namespace", req.Namespace, "podName", req.PodName)
	}
	return err
}
