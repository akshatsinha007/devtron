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

package router

import (
	"github.com/devtron-labs/devtron/api/restHandler"
	"github.com/gorilla/mux"
)

type WebhookListenerRouter interface {
	InitWebhookListenerRouter(gocdRouter *mux.Router)
}

type WebhookListenerRouterImpl struct {
	webhookEventHandler restHandler.WebhookEventHandler
}

func NewWebhookListenerRouterImpl(webhookEventHandler restHandler.WebhookEventHandler) *WebhookListenerRouterImpl {
	return &WebhookListenerRouterImpl{
		webhookEventHandler: webhookEventHandler,
	}
}

func (impl WebhookListenerRouterImpl) InitWebhookListenerRouter(configRouter *mux.Router) {
	configRouter.Path("/{gitHostId}").
		HandlerFunc(impl.webhookEventHandler.OnWebhookEvent).
		Methods("POST")
	configRouter.Path("/{gitHostId}/{secret}").
		HandlerFunc(impl.webhookEventHandler.OnWebhookEvent).
		Methods("POST")
}
