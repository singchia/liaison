package controlplane

import (
	"context"
	"time"

	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/liaison/internal/repo/model"
)

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
	proxies, err := cp.repo.ListProxies(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	return &v1.ListProxiesResponse{
		Code:    200,
		Message: "success",
		Data: &v1.Proxies{
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
	return &v1.Proxy{
		Id:        uint64(proxy.ID),
		Name:      proxy.Name,
		CreatedAt: proxy.CreatedAt.Format(time.DateTime),
		UpdatedAt: proxy.UpdatedAt.Format(time.DateTime),
	}
}
