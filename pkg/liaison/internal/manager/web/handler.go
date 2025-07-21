package web

import (
	"context"

	v1 "github.com/singchia/liaison/api/v1"
)

//-- Edge --//

// @Summary CreateEdge
// @Tags 1.0
// @Param params query v1.CreateEdgeRequest true "queries"
// @Success 200 {object} v1.CreateEdgeResponse
// @Router /api/v1/edge [post]
func (web *web) CreateEdge(ctx context.Context, req *v1.CreateEdgeRequest) (*v1.CreateEdgeResponse, error) {
	return nil, nil
}

// @Summary GetEdge
// @Tags 1.0
// @Param id path int true "edge id"
// @Param params query v1.GetEdgeRequest true "queries"
// @Success 200 {object} v1.GetEdgeResponse
// @Router /api/v1/edge/{id} [get]
func (web *web) GetEdge(ctx context.Context, req *v1.GetEdgeRequest) (*v1.GetEdgeResponse, error) {
	return nil, nil
}

// @Summary ListEdges
// @Tags 1.0
// @Param params query v1.ListEdgesRequest true "queries"
// @Success 200 {object} v1.ListEdgesResponse
// @Router /api/v1/edges [get]
func (web *web) ListEdges(ctx context.Context, req *v1.ListEdgesRequest) (*v1.ListEdgesResponse, error) {
	return nil, nil
}

// @Summary UpdateEdge
// @Tags 1.0
// @Param id path int true "edge id"
// @Param params query v1.UpdateEdgeRequest true "queries"
// @Success 200 {object} v1.UpdateEdgeResponse
// @Router /api/v1/edge/{id} [put]
func (web *web) UpdateEdge(ctx context.Context, req *v1.UpdateEdgeRequest) (*v1.UpdateEdgeResponse, error) {
	return nil, nil
}

// @Summary DeleteEdge
// @Tags 1.0
// @Param id path int true "edge id"
// @Success 200 {object} v1.DeleteEdgeResponse
// @Router /api/v1/edge/{id} [delete]
func (web *web) DeleteEdge(ctx context.Context, req *v1.DeleteEdgeRequest) (*v1.DeleteEdgeResponse, error) {
	return nil, nil
}

//-- Device --//

// @Summary ListDevices
// @Tags 1.0
// @Param params query v1.ListDevicesRequest true "queries"
// @Success 200 {object} v1.ListDevicesResponse
// @Router /api/v1/devices [get]
func (web *web) ListDevices(ctx context.Context, req *v1.ListDevicesRequest) (*v1.ListDevicesResponse, error) {
	return nil, nil
}

// @Summary GetDevice
// @Tags 1.0
// @Param id path int true "device id"
// @Param params query v1.GetDeviceRequest true "queries"
// @Success 200 {object} v1.GetDeviceResponse
// @Router /api/v1/device/{id} [get]
func (web *web) GetDevice(ctx context.Context, req *v1.GetDeviceRequest) (*v1.GetDeviceResponse, error) {
	return nil, nil
}

// @Summary UpdateDevice
// @Tags 1.0
// @Param id path int true "device id"
// @Param params query v1.UpdateDeviceRequest true "queries"
// @Success 200 {object} v1.UpdateDeviceResponse
// @Router /api/v1/device/{id} [put]
func (web *web) UpdateDevice(ctx context.Context, req *v1.UpdateDeviceRequest) (*v1.UpdateDeviceResponse, error) {
	return nil, nil
}

//-- Application --//

// @Summary ListApplications
// @Tags 1.0
// @Param params query v1.ListApplicationsRequest true "queries"
// @Success 200 {object} v1.ListApplicationsResponse
// @Router /api/v1/applications [get]
func (web *web) ListApplications(ctx context.Context, req *v1.ListApplicationsRequest) (*v1.ListApplicationsResponse, error) {
	return nil, nil
}

func (web *web) UpdateApplication(ctx context.Context, req *v1.UpdateApplicationRequest) (*v1.UpdateApplicationResponse, error) {
	return nil, nil
}

func (web *web) DeleteApplication(ctx context.Context, req *v1.DeleteApplicationRequest) (*v1.DeleteApplicationResponse, error) {
	return nil, nil
}

func (web *web) ListProxies(ctx context.Context, req *v1.ListProxiesRequest) (*v1.ListProxiesResponse, error) {
	return nil, nil
}

func (web *web) CreateProxy(ctx context.Context, req *v1.CreateProxyRequest) (*v1.CreateProxyResponse, error) {
	return nil, nil
}

func (web *web) UpdateProxy(ctx context.Context, req *v1.UpdateProxyRequest) (*v1.UpdateProxyResponse, error) {
	return nil, nil
}

func (web *web) DeleteProxy(ctx context.Context, req *v1.DeleteProxyRequest) (*v1.DeleteProxyResponse, error) {
	return nil, nil
}
