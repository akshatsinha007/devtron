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

package cluster

import (
	"context"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/cluster"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/devtron-labs/devtron/client/argocdServer/bean"
	"github.com/devtron-labs/devtron/client/argocdServer/connection"
	"go.uber.org/zap"
	"time"
)

type ServiceClient interface {
	// List returns list of clusters
	List(ctx context.Context, connectionConfig *bean.ArgoGRPCConfig, query *cluster.ClusterQuery) (*v1alpha1.ClusterList, error)
	// Create creates a cluster
	Create(ctx context.Context, connectionConfig *bean.ArgoGRPCConfig, query *cluster.ClusterCreateRequest) (*v1alpha1.Cluster, error)
	// CreateFromKubeConfig installs the argocd-manager service account into the cluster specified in the given kubeconfig and context
	CreateFromKubeConfig(ctx context.Context, connectionConfig *bean.ArgoGRPCConfig, query *cluster.ClusterCreateRequest) (*v1alpha1.Cluster, error)
	// Get returns a cluster by server address
	Get(ctx context.Context, connectionConfig *bean.ArgoGRPCConfig, query *cluster.ClusterQuery) (*v1alpha1.Cluster, error)
	// Update updates a cluster
	Update(ctx context.Context, connectionConfig *bean.ArgoGRPCConfig, query *cluster.ClusterUpdateRequest) (*v1alpha1.Cluster, error)
	// Delete deletes a cluster
	Delete(ctx context.Context, connectionConfig *bean.ArgoGRPCConfig, query *cluster.ClusterQuery) (*cluster.ClusterResponse, error)
}

type ServiceClientImpl struct {
	logger           *zap.SugaredLogger
	argoCdConnection connection.ArgoCDConnectionManager
}

func NewServiceClientImpl(logger *zap.SugaredLogger, argoCdConnection connection.ArgoCDConnectionManager) *ServiceClientImpl {
	return &ServiceClientImpl{
		logger:           logger,
		argoCdConnection: argoCdConnection,
	}
}

func (c ServiceClientImpl) getService(connectionConfig *bean.ArgoGRPCConfig) (cluster.ClusterServiceClient, error) {
	conn := c.argoCdConnection.GetGrpcClientConnection(connectionConfig)
	//defer conn.Close()
	return cluster.NewClusterServiceClient(conn), nil
}

func (c ServiceClientImpl) List(ctx context.Context, connectionConfig *bean.ArgoGRPCConfig, query *cluster.ClusterQuery) (*v1alpha1.ClusterList, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	client, err := c.getService(connectionConfig)
	if err != nil {
		return nil, err
	}
	return client.List(ctx, query)
}

func (c ServiceClientImpl) Create(ctx context.Context, connectionConfig *bean.ArgoGRPCConfig, query *cluster.ClusterCreateRequest) (*v1alpha1.Cluster, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	client, err := c.getService(connectionConfig)
	if err != nil {
		return nil, err
	}
	return client.Create(ctx, query)
}

func (c ServiceClientImpl) CreateFromKubeConfig(ctx context.Context, connectionConfig *bean.ArgoGRPCConfig, query *cluster.ClusterCreateRequest) (*v1alpha1.Cluster, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	client, err := c.getService(connectionConfig)
	if err != nil {
		return nil, err
	}
	return client.Create(ctx, query)
}

func (c ServiceClientImpl) Get(ctx context.Context, connectionConfig *bean.ArgoGRPCConfig, query *cluster.ClusterQuery) (*v1alpha1.Cluster, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	client, err := c.getService(connectionConfig)
	if err != nil {
		return nil, err
	}
	return client.Get(ctx, query)
}

func (c ServiceClientImpl) Update(ctx context.Context, connectionConfig *bean.ArgoGRPCConfig, query *cluster.ClusterUpdateRequest) (*v1alpha1.Cluster, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	client, err := c.getService(connectionConfig)
	if err != nil {
		return nil, err
	}
	return client.Update(ctx, query)
}

func (c ServiceClientImpl) Delete(ctx context.Context, connectionConfig *bean.ArgoGRPCConfig, query *cluster.ClusterQuery) (*cluster.ClusterResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	client, err := c.getService(connectionConfig)
	if err != nil {
		return nil, err
	}
	return client.Delete(ctx, query)
}
