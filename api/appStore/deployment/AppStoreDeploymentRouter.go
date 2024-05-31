/*
 * Copyright (c) 2020-2024. Devtron Inc.
 */

package appStoreDeployment

import (
	"github.com/gorilla/mux"
)

type AppStoreDeploymentRouter interface {
	Init(configRouter *mux.Router)
}

type AppStoreDeploymentRouterImpl struct {
	appStoreDeploymentRestHandler AppStoreDeploymentRestHandler
}

func NewAppStoreDeploymentRouterImpl(appStoreDeploymentRestHandler AppStoreDeploymentRestHandler) *AppStoreDeploymentRouterImpl {
	return &AppStoreDeploymentRouterImpl{
		appStoreDeploymentRestHandler: appStoreDeploymentRestHandler,
	}
}

func (router AppStoreDeploymentRouterImpl) Init(configRouter *mux.Router) {
	configRouter.Path("/application/install").
		HandlerFunc(router.appStoreDeploymentRestHandler.InstallApp).Methods("POST")

	configRouter.Path("/application/update").
		HandlerFunc(router.appStoreDeploymentRestHandler.UpdateInstalledApp).Methods("PUT")

	configRouter.Path("/installed-app/{appStoreId}").
		HandlerFunc(router.appStoreDeploymentRestHandler.GetInstalledAppsByAppStoreId).Methods("GET")

	configRouter.Path("/application/delete/{id}").
		HandlerFunc(router.appStoreDeploymentRestHandler.DeleteInstalledApp).Methods("DELETE")

	configRouter.Path("/application/helm/link-to-chart-store").
		HandlerFunc(router.appStoreDeploymentRestHandler.LinkHelmApplicationToChartStore).Methods("PUT")

	configRouter.Path("/application/version/{installedAppVersionId}").
		HandlerFunc(router.appStoreDeploymentRestHandler.GetInstalledAppVersion).Methods("GET")

	configRouter.Path("/application/update/project").
		HandlerFunc(router.appStoreDeploymentRestHandler.UpdateProjectHelmApp).Methods("PUT")

}
