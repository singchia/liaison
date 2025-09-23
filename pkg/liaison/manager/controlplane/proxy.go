package controlplane

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/liaison/repo/dao"
	"github.com/singchia/liaison/pkg/liaison/repo/model"
	"github.com/singchia/liaison/pkg/proto"
)

func (cp *controlPlane) RegisterProxyManager(proxyManager proto.ProxyManager) {
	cp.proxyManager = proxyManager
}

func (cp *controlPlane) CreateProxy(_ context.Context, req *v1.CreateProxyRequest) (*v1.CreateProxyResponse, error) {
	// 查看是否有冲突
	// 获取application
	application, err := cp.repo.GetApplicationByID(uint(req.ApplicationId))
	if err != nil {
		log.Warnf("application %d not found", req.ApplicationId)
		return nil, err
	}

	// 创建Proxy持久化
	proxy := &model.Proxy{
		Name:          req.Name,
		Status:        model.ProxyStatusRunning,
		Description:   req.Description,
		Port:          int(req.Port),
		ApplicationID: uint(req.ApplicationId),
	}
	err = cp.repo.CreateProxy(proxy)
	if err != nil {
		log.Warnf("failed to create proxy: %s", err)
		return nil, err
	}

	// 创建Proxy
	cp.proxyManager.CreateProxy(context.Background(), &proto.Proxy{
		ID:        int(proxy.ID),
		Name:      proxy.Name,
		ProxyPort: int(proxy.Port),
		EdgeID:    uint64(application.EdgeIDs[0]),
		Dst:       fmt.Sprintf("%s:%d", application.IP, application.Port),
	})
	return &v1.CreateProxyResponse{
		Code:    200,
		Message: "success",
	}, nil
}

func (cp *controlPlane) ListProxies(_ context.Context, req *v1.ListProxiesRequest) (*v1.ListProxiesResponse, error) {
	// list proxies
	proxies, err := cp.repo.ListProxies(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	count, err := cp.repo.CountProxies()
	if err != nil {
		return nil, err
	}
	ids := make([]uint, len(proxies))
	for i, proxy := range proxies {
		ids[i] = proxy.ApplicationID
	}
	// list applications
	applications, err := cp.repo.ListApplications(&dao.ListApplicationsQuery{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
		IDs:      ids,
	})
	if err != nil {
		return nil, err
	}
	// add applications to proxies
	// 创建一个 map 来快速查找 application
	appMap := make(map[uint]*model.Application)
	for _, app := range applications {
		appMap[app.ID] = app
	}

	for i := range proxies {
		if app, exists := appMap[proxies[i].ApplicationID]; exists {
			proxies[i].Application = app
		} else {
			log.Warnf("application %d not found", proxies[i].ApplicationID)
		}
	}

	return &v1.ListProxiesResponse{
		Code:    200,
		Message: "success",
		Data: &v1.Proxies{
			Total:   int32(count),
			Proxies: transformProxies(proxies),
		},
	}, nil
}

// 更新代理，不允许更新代理端口
func (cp *controlPlane) UpdateProxy(_ context.Context, req *v1.UpdateProxyRequest) (*v1.UpdateProxyResponse, error) {
	proxy, err := cp.repo.GetProxyByID(uint(req.Id))
	if err != nil {
		return nil, err
	}
	proxy.Name = req.Name
	err = cp.repo.UpdateProxy(proxy)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateProxyResponse{
		Code:    200,
		Message: "success",
		Data: &v1.Proxy{
			Id:        uint64(proxy.ID),
			Name:      proxy.Name,
			CreatedAt: proxy.CreatedAt.Format(time.DateTime),
			UpdatedAt: proxy.UpdatedAt.Format(time.DateTime),
		},
	}, nil
}

func (cp *controlPlane) DeleteProxy(_ context.Context, req *v1.DeleteProxyRequest) (*v1.DeleteProxyResponse, error) {
	err := cp.repo.DeleteProxy(uint(req.Id))
	if err != nil {
		return nil, err
	}
	// 删除正在工作的代理
	err = cp.proxyManager.DeleteProxy(context.Background(), int(req.Id))
	if err != nil {
		return nil, err
	}
	return &v1.DeleteProxyResponse{
		Code:    200,
		Message: "success",
	}, nil
}

func transformProxies(proxies []*model.Proxy) []*v1.Proxy {
	proxiesV1 := make([]*v1.Proxy, len(proxies))
	for i, proxy := range proxies {
		proxiesV1[i] = transformProxy(proxy)
	}
	return proxiesV1
}

func transformProxy(proxy *model.Proxy) *v1.Proxy {
	var application *v1.Application
	if proxy.Application != nil {
		application = transformApplication(proxy.Application)
	}

	// 将 ProxyStatus 转换为字符串
	var status string
	switch proxy.Status {
	case model.ProxyStatusRunning:
		status = "running"
	case model.ProxyStatusStopped:
		status = "stopped"
	default:
		status = "unknown"
	}

	return &v1.Proxy{
		Id:          uint64(proxy.ID),
		Name:        proxy.Name,
		Port:        int32(proxy.Port),
		Status:      status,
		Application: application,
		Description: proxy.Description,
		CreatedAt:   proxy.CreatedAt.Format(time.DateTime),
		UpdatedAt:   proxy.UpdatedAt.Format(time.DateTime),
	}
}
