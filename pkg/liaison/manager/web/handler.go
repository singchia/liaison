package web

import (
	"context"
	"errors"

	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/liaison/manager/iam"
)

//-- Edge --//

// @Summary CreateEdge
// @Tags 1.0
// @Param params query v1.CreateEdgeRequest true "queries"
// @Success 200 {object} v1.CreateEdgeResponse
// @Router /api/v1/edges [post]
func (web *web) CreateEdge(ctx context.Context, req *v1.CreateEdgeRequest) (*v1.CreateEdgeResponse, error) {
	return web.controlPlane.CreateEdge(ctx, req)
}

// @Summary GetEdge
// @Tags 1.0
// @Param id path int true "edge id"
// @Param params query v1.GetEdgeRequest true "queries"
// @Success 200 {object} v1.GetEdgeResponse
// @Router /api/v1/edges/{id} [get]
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
// @Router /api/v1/edges/{id} [put]
func (web *web) UpdateEdge(ctx context.Context, req *v1.UpdateEdgeRequest) (*v1.UpdateEdgeResponse, error) {
	return web.controlPlane.UpdateEdge(ctx, req)
}

// @Summary DeleteEdge
// @Tags 1.0
// @Param id path int true "edge id"
// @Success 200 {object} v1.DeleteEdgeResponse
// @Router /api/v1/edges/{id} [delete]
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
// @Router /api/v1/devices/{id} [get]
func (web *web) GetDevice(ctx context.Context, req *v1.GetDeviceRequest) (*v1.GetDeviceResponse, error) {
	return web.controlPlane.GetDevice(ctx, req)
}

// @Summary UpdateDevice
// @Tags 1.0
// @Param id path int true "device id"
// @Param params query v1.UpdateDeviceRequest true "queries"
// @Success 200 {object} v1.UpdateDeviceResponse
// @Router /api/v1/devices/{id} [put]
func (web *web) UpdateDevice(ctx context.Context, req *v1.UpdateDeviceRequest) (*v1.UpdateDeviceResponse, error) {
	return web.controlPlane.UpdateDevice(ctx, req)
}

//-- Application --//

// @Summary CreateApplication
// @Tags 1.0
// @Param params query v1.CreateApplicationRequest true "queries"
// @Success 200 {object} v1.CreateApplicationResponse
// @Router /api/v1/applications [post]
func (web *web) CreateApplication(ctx context.Context, req *v1.CreateApplicationRequest) (*v1.CreateApplicationResponse, error) {
	return web.controlPlane.CreateApplication(ctx, req)
}

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
// @Router /api/v1/applications/{id} [put]
func (web *web) UpdateApplication(ctx context.Context, req *v1.UpdateApplicationRequest) (*v1.UpdateApplicationResponse, error) {
	return web.controlPlane.UpdateApplication(ctx, req)
}

// @Summary DeleteApplication
// @Tags 1.0
// @Param id path int true "application id"
// @Success 200 {object} v1.DeleteApplicationResponse
// @Router /api/v1/applications/{id} [delete]
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
// @Router /api/v1/proxies [post]
func (web *web) CreateProxy(ctx context.Context, req *v1.CreateProxyRequest) (*v1.CreateProxyResponse, error) {
	return web.controlPlane.CreateProxy(ctx, req)
}

// @Summary UpdateProxy
// @Tags 1.0
// @Param id path int true "proxy id"
// @Param params query v1.UpdateProxyRequest true "queries"
// @Success 200 {object} v1.UpdateProxyResponse
// @Router /api/v1/proxies/{id} [put]
func (web *web) UpdateProxy(ctx context.Context, req *v1.UpdateProxyRequest) (*v1.UpdateProxyResponse, error) {
	return web.controlPlane.UpdateProxy(ctx, req)
}

// @Summary DeleteProxy
// @Tags 1.0
// @Param id path int true "proxy id"
// @Success 200 {object} v1.DeleteProxyResponse
// @Router /api/v1/proxies/{id} [delete]
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

//-- Auth --//

// Login 用户登录
func (web *web) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	// 转换请求类型
	iamReq := &iam.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := web.iamService.Login(iamReq)
	if err != nil {
		return nil, err
	}

	// 转换响应类型
	return &v1.LoginResponse{
		Code:    200,
		Message: "success",
		Data: &v1.LoginData{
			Token: resp.Token,
			User: &v1.User{
				Id:    uint64(resp.User.ID),
				Email: resp.User.Email,
			},
		},
	}, nil
}

// GetProfile 获取用户信息
func (web *web) GetProfile(ctx context.Context, req *v1.GetProfileRequest) (*v1.GetProfileResponse, error) {
	// 从context中获取用户信息（需要中间件设置）
	userID := ctx.Value("user_id")
	userEmail := ctx.Value("user_email")

	if userID == nil || userEmail == nil {
		return nil, errors.New("未认证")
	}

	user := &v1.User{
		Id:    uint64(userID.(uint)),
		Email: userEmail.(string),
	}

	return &v1.GetProfileResponse{
		Code:    200,
		Message: "success",
		Data:    user,
	}, nil
}

// Logout 用户登出
func (web *web) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	// JWT是无状态的，登出只需要客户端删除token
	return &v1.LogoutResponse{
		Code:    200,
		Message: "登出成功",
	}, nil
}

// Health 健康检查
func (web *web) Health(ctx context.Context, req *v1.HealthRequest) (*v1.HealthResponse, error) {
	return &v1.HealthResponse{
		Status: "ok",
	}, nil
}
