package controlplane

import (
	"context"

	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/liaison/internal/repo"
)

type ControlPlane interface {
	CreateEdge(ctx context.Context, req *v1.CreateEdgeRequest) (*v1.CreateEdgeResponse, error)
	GetEdge(ctx context.Context, req *v1.GetEdgeRequest) (*v1.GetEdgeResponse, error)
	ListEdges(ctx context.Context, req *v1.ListEdgesRequest) (*v1.ListEdgesResponse, error)
	UpdateEdge(ctx context.Context, req *v1.UpdateEdgeRequest) (*v1.UpdateEdgeResponse, error)
	DeleteEdge(ctx context.Context, req *v1.DeleteEdgeRequest) (*v1.DeleteEdgeResponse, error)

	ListDevices(ctx context.Context, req *v1.ListDevicesRequest) (*v1.ListDevicesResponse, error)
	GetDevice(ctx context.Context, req *v1.GetDeviceRequest) (*v1.GetDeviceResponse, error)
	UpdateDevice(ctx context.Context, req *v1.UpdateDeviceRequest) (*v1.UpdateDeviceResponse, error)

	ListApplications(ctx context.Context, req *v1.ListApplicationsRequest) (*v1.ListApplicationsResponse, error)
	UpdateApplication(ctx context.Context, req *v1.UpdateApplicationRequest) (*v1.UpdateApplicationResponse, error)
	DeleteApplication(ctx context.Context, req *v1.DeleteApplicationRequest) (*v1.DeleteApplicationResponse, error)

	ListProxies(ctx context.Context, req *v1.ListProxiesRequest) (*v1.ListProxiesResponse, error)
	CreateProxy(ctx context.Context, req *v1.CreateProxyRequest) (*v1.CreateProxyResponse, error)
	UpdateProxy(ctx context.Context, req *v1.UpdateProxyRequest) (*v1.UpdateProxyResponse, error)
	DeleteProxy(ctx context.Context, req *v1.DeleteProxyRequest) (*v1.DeleteProxyResponse, error)
}

type controlPlane struct {
	repo repo.Repo
}
