package controlplane

import (
	"context"
	"time"

	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/liaison/repo/dao"
	"github.com/singchia/liaison/pkg/liaison/repo/model"
	"github.com/singchia/liaison/pkg/proto"
)

func (cp *controlPlane) RegisterProxyManager(proxyManager proto.ProxyManager) {
	cp.proxyManager = proxyManager
}

func (cp *controlPlane) CreateProxy(_ context.Context, req *v1.CreateProxyRequest) (*v1.CreateProxyResponse, error) {
	proxy := &model.Proxy{
		Name: req.Name,
	}
	err := cp.repo.CreateProxy(proxy)
	if err != nil {
		return nil, err
	}
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
		}
		// 如果 application 不存在，Application 字段会保持为 nil
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
