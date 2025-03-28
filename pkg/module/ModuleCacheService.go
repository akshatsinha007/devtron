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

package module

import (
	"context"
	"github.com/devtron-labs/common-lib/utils/k8s"
	"github.com/devtron-labs/devtron/pkg/module/bean"
	moduleRepo "github.com/devtron-labs/devtron/pkg/module/repo"
	serverBean "github.com/devtron-labs/devtron/pkg/server/bean"
	serverEnvConfig "github.com/devtron-labs/devtron/pkg/server/config"
	serverDataStore "github.com/devtron-labs/devtron/pkg/server/store"
	"github.com/devtron-labs/devtron/pkg/team"
	util2 "github.com/devtron-labs/devtron/util"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

type ModuleCacheService interface {
}

type ModuleCacheServiceImpl struct {
	logger           *zap.SugaredLogger
	mutex            sync.Mutex
	K8sUtil          *k8s.K8sServiceImpl
	moduleEnvConfig  *bean.ModuleEnvConfig
	serverEnvConfig  *serverEnvConfig.ServerEnvConfig
	serverDataStore  *serverDataStore.ServerDataStore
	moduleRepository moduleRepo.ModuleRepository
	teamService      team.TeamService
}

func NewModuleCacheServiceImpl(logger *zap.SugaredLogger, K8sUtil *k8s.K8sServiceImpl, moduleEnvConfig *bean.ModuleEnvConfig, serverEnvConfig *serverEnvConfig.ServerEnvConfig,
	serverDataStore *serverDataStore.ServerDataStore, moduleRepository moduleRepo.ModuleRepository, teamService team.TeamService) (*ModuleCacheServiceImpl, error) {
	impl := &ModuleCacheServiceImpl{
		logger:           logger,
		K8sUtil:          K8sUtil,
		moduleEnvConfig:  moduleEnvConfig,
		serverEnvConfig:  serverEnvConfig,
		serverDataStore:  serverDataStore,
		moduleRepository: moduleRepository,
		teamService:      teamService,
	}

	// DB migration - if server mode is not base stack and data in modules table is empty, then insert entries in DB
	if !util2.IsBaseStack() {
		exists, err := impl.moduleRepository.ModuleExists()
		if err != nil {
			log.Println("Error while checking if any module exists in database.", "error", err)
			return nil, err
		}
		if !exists {
			// insert cicd module entry
			err = impl.updateModuleToInstalled(bean.ModuleNameCiCd)
			if err != nil {
				return nil, err
			}

		}
	}

	// if devtron user type is OSS_HELM then only installer object and modules installation is useful
	if serverEnvConfig.DevtronInstallationType == serverBean.DevtronInstallationTypeOssHelm {
		// listen in installer object to save status in-memory
		// build informer to listen on installer object
		err := impl.buildInformerToListenOnInstallerObject()
		if err != nil {
			log.Println("Error building informer:", err)
			return nil, err
		}
	}

	return impl, nil
}

func (impl *ModuleCacheServiceImpl) updateModuleToInstalled(moduleName string) error {
	module := &moduleRepo.Module{
		Name:      moduleName,
		Version:   impl.serverDataStore.CurrentVersion,
		Status:    bean.ModuleStatusInstalled,
		UpdatedOn: time.Now(),
	}
	err := impl.moduleRepository.Save(module)
	if err != nil {
		log.Println("Error while saving module.", "moduleName", moduleName, "error", err)
		return err
	}
	return nil
}

func (impl *ModuleCacheServiceImpl) buildInformerToListenOnInstallerObject() error {
	impl.logger.Debug("building informer cache to listen on installer object")
	_, _, clusterDynamicClient, err := impl.K8sUtil.GetK8sInClusterConfigAndDynamicClients()
	if err != nil {
		log.Println("not able to get k8s cluster rest config.", "error", err)
		return err
	}
	go func() {
		installerResource := schema.GroupVersionResource{
			Group:    impl.serverEnvConfig.InstallerCrdObjectGroupName,
			Version:  impl.serverEnvConfig.InstallerCrdObjectVersion,
			Resource: impl.serverEnvConfig.InstallerCrdObjectResource,
		}
		factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
			clusterDynamicClient, time.Minute, impl.serverEnvConfig.InstallerCrdNamespace, nil)
		informer := factory.ForResource(installerResource).Informer()

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				impl.handleInstallerObjectChange(obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				impl.handleInstallerObjectChange(newObj)
			},
			DeleteFunc: func(obj interface{}) {
				impl.serverDataStore.InstallerCrdObjectStatus = ""
				impl.serverDataStore.InstallerCrdObjectExists = false
			},
		})

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		go informer.Run(ctx.Done())
		<-ctx.Done()
	}()
	return nil
}

func (impl *ModuleCacheServiceImpl) handleInstallerObjectChange(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	val, _, _ := unstructured.NestedString(u.Object, "status", "sync", "status")
	impl.serverDataStore.InstallerCrdObjectStatus = val
	impl.serverDataStore.InstallerCrdObjectExists = true
}
