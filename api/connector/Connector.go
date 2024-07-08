/*
 * Copyright (c) 2020-2024. Devtron Inc.
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

package connector

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/devtron-labs/devtron/api/bean"
	"github.com/gogo/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/juju/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var delimiter = []byte("\n\n")

type Pump interface {
	StartStreamWithTransformer(w http.ResponseWriter, recv func() (proto.Message, error), err error, transformer func(interface{}) interface{})
	StartK8sStreamWithHeartBeat(w http.ResponseWriter, isReconnect bool, stream io.ReadCloser, err error)
}

type PumpImpl struct {
	logger *zap.SugaredLogger
}

func NewPumpImpl(logger *zap.SugaredLogger) *PumpImpl {
	return &PumpImpl{
		logger: logger,
	}
}

func (impl PumpImpl) StartK8sStreamWithHeartBeat(w http.ResponseWriter, isReconnect bool, stream io.ReadCloser, err error) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "unexpected server doesnt support streaming", http.StatusInternalServerError)
	}

	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-cache, no-transform")

	if err != nil {
		err := impl.sendEvent(nil, []byte("CUSTOM_ERR_STREAM"), []byte(err.Error()), w)
		if err != nil {
			impl.logger.Errorw("error in writing data over sse", "err", err)
		}
		return
	}

	if isReconnect {
		err := impl.sendEvent(nil, []byte("RECONNECT_STREAM"), []byte("RECONNECT_STREAM"), w)
		if err != nil {
			impl.logger.Errorw("error in writing data over sse", "err", err)
			return
		}
	}
	// heartbeat start
	ticker := time.NewTicker(30 * time.Second)
	done := make(chan bool)
	var mux sync.Mutex
	go func() error {
		for {
			select {
			case <-done:
				return nil
			case t := <-ticker.C:
				mux.Lock()
				err := impl.sendEvent(nil, []byte("PING"), []byte(t.String()), w)
				mux.Unlock()
				if err != nil {
					impl.logger.Errorw("error in writing PING over sse", "err", err)
					return err
				}
				f.Flush()
			}
		}
	}()
	defer func() {
		ticker.Stop()
		done <- true
	}()

	bufReader := bufio.NewReader(stream)
	eof := false
	for !eof {
		log, err := bufReader.ReadString('\n')
		if err == io.EOF {
			eof = true
			// stop if we reached end of stream and the next line is empty
			if log == "" {
				return
			}
		} else if err != nil && err != io.EOF {
			impl.logger.Errorw("error in reading buffer string, StartK8sStreamWithHeartBeat", "err", err)
			return
		}
		log = strings.TrimSpace(log) // Remove trailing line ending
		a := regexp.MustCompile(" ")
		splitLog := a.Split(log, 2)
		parsedTime, err := time.Parse(time.RFC3339, splitLog[0])
		if err != nil {
			impl.logger.Errorw("error in writing data over sse", "err", err)
			return
		}
		eventId := strconv.FormatInt(parsedTime.UnixNano(), 10)
		mux.Lock()
		if len(splitLog) == 2 {
			err = impl.sendEvent([]byte(eventId), nil, []byte(splitLog[1]), w)
		}
		mux.Unlock()
		if err != nil {
			impl.logger.Errorw("error in writing data over sse", "err", err)
			return
		}
		f.Flush()
	}
	// heartbeat end
}

func (impl PumpImpl) StartStreamWithTransformer(w http.ResponseWriter, recv func() (proto.Message, error), err error, transformer func(interface{}) interface{}) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "unexpected server doesnt support streaming", http.StatusInternalServerError)
		impl.logger.Debugw("unexpected server doesnt support streaming", "internal Server error", http.StatusInternalServerError)
	}
	if err != nil {
		http.Error(w, errors.Details(err), http.StatusInternalServerError)
		impl.logger.Debugw("error occurred during the setting up the flusher in wrote header  ", "error", err)
	}
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	var wroteHeader bool
	for {
		resp, err := recv()
		if err == io.EOF {
			impl.logger.Debugw("error in receiving stream proto message ", "error", err)
			return
		}
		impl.logger.Debugw("receiving stream proto message ", "response", resp, "error", err)

		if err != nil {
			impl.logger.Errorf("Error occurred while reading data from argocd %+v\n", err)
			impl.handleForwardResponseStreamError(wroteHeader, w, err)
			impl.logger.Debugw("Failed to send response chunk ", "error", err)

			return
		}
		response := bean.Response{}
		response.Result = transformer(resp)
		buf, err := json.Marshal(response)
		data := "data: " + string(buf)

		if _, err = w.Write([]byte(data)); err != nil {
			impl.logger.Errorf("Failed to send response chunk: %v", err)
			impl.logger.Debugw("Failed to send response chunk ", "error", err)
			return
		}
		wroteHeader = true
		if _, err = w.Write(delimiter); err != nil {
			impl.logger.Errorf("Failed to send delimiter chunk: %v", err)
			impl.logger.Debugw("failed to send delimiter chunk ", "error", err)
			return
		}
		impl.logger.Debugw("sending chunk data ", "data", data, "wroteheader", w)
		f.Flush()
	}
}

func (impl *PumpImpl) sendEvent(eventId []byte, eventName []byte, payload []byte, w http.ResponseWriter) error {
	var res []byte
	if len(eventId) > 0 {
		res = append(res, "id: "...)
		res = append(res, eventId...)
		res = append(res, '\n')
	}
	if len(eventName) > 0 {
		res = append(res, "event:"...)
		res = append(res, eventName...)
		res = append(res, '\n')
	}
	if len(payload) > 0 {
		res = append(res, "data:"...)
		res = append(res, payload...)
	}
	res = append(res, '\n', '\n')
	if _, err := w.Write(res); err != nil {
		impl.logger.Errorf("Failed to send response chunk: %v", err)
		return err
	}

	return nil
}

func (impl PumpImpl) handleForwardResponseStreamError(wroteHeader bool, w http.ResponseWriter, err error) {
	code := "000"
	if !wroteHeader {
		s, ok := status.FromError(err)
		if !ok {
			s = status.New(codes.Unknown, err.Error())
		}
		w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
		code = fmt.Sprint(s.Code())
	}
	response := bean.Response{}
	apiErr := bean.ApiError{}
	apiErr.Code = code // 000=unknown
	apiErr.InternalMessage = errors.Details(err)
	response.Errors = []bean.ApiError{apiErr}
	buf, merr := json.Marshal(response)
	if merr != nil {
		impl.logger.Errorf("Failed to marshal response %+v\n", merr)
	}
	if _, werr := w.Write(buf); werr != nil {
		impl.logger.Errorf("Failed to notify error to client: %v", werr)
		return
	}
}
