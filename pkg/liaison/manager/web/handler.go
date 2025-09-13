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
	return web.controlPlane.CreateEdge(ctx, req)
}

// @Summary GetEdge
// @Tags 1.0
// @Param id path int true "edge id"
// @Param params query v1.GetEdgeRequest true "queries"
// @Success 200 {object} v1.GetEdgeResponse
// @Router /api/v1/edge/{id} [get]
func (web *web) GetEdge(ctx context.Context, req *v1.GetEdgeRequest) (*v1.GetEdgeResponse, error) {
	return web.controlPlane.GetEdge(ctx, req)
}

// @Summary ListEdges
// @Tags 1.0
// @Param params query v1.ListEdgesRequest true "queries"
// @Success 200 {object} v1.ListEdgesResponse
// @Router /api/v1/edges [get]
func (web *web) ListEdges(ctx context.Context, req *v1.ListEdgesRequest) (*v1.ListEdgesResponse, error) {
	return web.controlPlane.ListEdges(ctx, req)
}

// @Summary UpdateEdge
// @Tags 1.0
// @Param id path int true "edge id"
// @Param params query v1.UpdateEdgeRequest true "queries"
// @Success 200 {object} v1.UpdateEdgeResponse
// @Router /api/v1/edge/{id} [put]
func (web *web) UpdateEdge(ctx context.Context, req *v1.UpdateEdgeRequest) (*v1.UpdateEdgeResponse, error) {
	return web.controlPlane.UpdateEdge(ctx, req)
}

// @Summary DeleteEdge
// @Tags 1.0
// @Param id path int true "edge id"
// @Success 200 {object} v1.DeleteEdgeResponse
// @Router /api/v1/edge/{id} [delete]
func (web *web) DeleteEdge(ctx context.Context, req *v1.DeleteEdgeRequest) (*v1.DeleteEdgeResponse, error) {
	return web.controlPlane.DeleteEdge(ctx, req)
}

//-- Device --//

// @Summary ListDevices
// @Tags 1.0
// @Param params query v1.ListDevicesRequest true "queries"
// @Success 200 {object} v1.ListDevicesResponse
// @Router /api/v1/devices [get]
func (web *web) ListDevices(ctx context.Context, req *v1.ListDevicesRequest) (*v1.ListDevicesResponse, error) {
	return web.controlPlane.ListDevices(ctx, req)
}

// @Summary GetDevice
// @Tags 1.0
// @Param id path int true "device id"
// @Param params query v1.GetDeviceRequest true "queries"
// @Success 200 {object} v1.GetDeviceResponse
// @Router /api/v1/device/{id} [get]
func (web *web) GetDevice(ctx context.Context, req *v1.GetDeviceRequest) (*v1.GetDeviceResponse, error) {
	return web.controlPlane.GetDevice(ctx, req)
}

// @Summary UpdateDevice
// @Tags 1.0
// @Param id path int true "device id"
// @Param params query v1.UpdateDeviceRequest true "queries"
// @Success 200 {object} v1.UpdateDeviceResponse
// @Router /api/v1/device/{id} [put]
func (web *web) UpdateDevice(ctx context.Context, req *v1.UpdateDeviceRequest) (*v1.UpdateDeviceResponse, error) {
	return web.controlPlane.UpdateDevice(ctx, req)
}

//-- Application --//

// @Summary ListApplications
// @Tags 1.0
// @Param params query v1.ListApplicationsRequest true "queries"
// @Success 200 {object} v1.ListApplicationsResponse
// @Router /api/v1/applications [get]
func (web *web) ListApplications(ctx context.Context, req *v1.ListApplicationsRequest) (*v1.ListApplicationsResponse, error) {
	return web.controlPlane.ListApplications(ctx, req)
}

// @Summary UpdateApplication
// @Tags 1.0
// @Param id path int true "application id"
// @Param params query v1.UpdateApplicationRequest true "queries"
// @Success 200 {object} v1.UpdateApplicationResponse
// @Router /api/v1/application/{id} [put]
func (web *web) UpdateApplication(ctx context.Context, req *v1.UpdateApplicationRequest) (*v1.UpdateApplicationResponse, error) {
	return web.controlPlane.UpdateApplication(ctx, req)
}

// @Summary DeleteApplication
// @Tags 1.0
// @Param id path int true "application id"
// @Success 200 {object} v1.DeleteApplicationResponse
// @Router /api/v1/application/{id} [delete]
func (web *web) DeleteApplication(ctx context.Context, req *v1.DeleteApplicationRequest) (*v1.DeleteApplicationResponse, error) {
	return web.controlPlane.DeleteApplication(ctx, req)
}

//-- Proxy --//

// @Summary ListProxies
// @Tags 1.0
// @Param params query v1.ListProxiesRequest true "queries"
// @Success 200 {object} v1.ListProxiesResponse
// @Router /api/v1/proxies [get]
func (web *web) ListProxies(ctx context.Context, req *v1.ListProxiesRequest) (*v1.ListProxiesResponse, error) {
	return web.controlPlane.ListProxies(ctx, req)
}

// @Summary CreateProxy
// @Tags 1.0
// @Param params query v1.CreateProxyRequest true "queries"
// @Success 200 {object} v1.CreateProxyResponse
// @Router /api/v1/proxy [post]
func (web *web) CreateProxy(ctx context.Context, req *v1.CreateProxyRequest) (*v1.CreateProxyResponse, error) {
	return web.controlPlane.CreateProxy(ctx, req)
}

// @Summary UpdateProxy
// @Tags 1.0
// @Param id path int true "proxy id"
// @Param params query v1.UpdateProxyRequest true "queries"
// @Success 200 {object} v1.UpdateProxyResponse
// @Router /api/v1/proxy/{id} [put]
func (web *web) UpdateProxy(ctx context.Context, req *v1.UpdateProxyRequest) (*v1.UpdateProxyResponse, error) {
	return web.controlPlane.UpdateProxy(ctx, req)
}

// @Summary DeleteProxy
// @Tags 1.0
// @Param id path int true "proxy id"
// @Success 200 {object} v1.DeleteProxyResponse
// @Router /api/v1/proxy/{id} [delete]
func (web *web) DeleteProxy(ctx context.Context, req *v1.DeleteProxyRequest) (*v1.DeleteProxyResponse, error) {
	return web.controlPlane.DeleteProxy(ctx, req)
}

//-- Task --//

// @Summary CreateEdgeScanApplicationTask
// @Tags 1.0
// @Param params query v1.CreateEdgeScanApplicationTaskRequest true "queries"
// @Success 200 {object} v1.CreateEdgeScanApplicationTaskResponse
// @Router /api/v1/edges/{edge_id}/scan_application_tasks [post]
func (web *web) CreateEdgeScanApplicationTask(ctx context.Context, req *v1.CreateEdgeScanApplicationTaskRequest) (*v1.CreateEdgeScanApplicationTaskResponse, error) {
	return web.controlPlane.CreateEdgeScanApplicationTask(ctx, req)
}

// @Summary GetEdgeScanApplicationTask
// @Tags 1.0
// @Param params query v1.GetEdgeScanApplicationTaskRequest true "queries"
// @Success 200 {object} v1.GetEdgeScanApplicationTaskResponse
// @Router /api/v1/edges/{edge_id}/scan_application_tasks [get]
func (web *web) GetEdgeScanApplicationTask(ctx context.Context, req *v1.GetEdgeScanApplicationTaskRequest) (*v1.GetEdgeScanApplicationTaskResponse, error) {
	return web.controlPlane.GetEdgeScanApplicationTask(ctx, req)
}
