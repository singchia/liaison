package frontierbound

import (
	"context"
	"net"

	"github.com/singchia/frontier/api/dataplane/v1/service"
	"github.com/singchia/liaison/pkg/liaison/internal/config"
	"github.com/singchia/liaison/pkg/liaison/internal/repo"
)

type FrontierBound interface {
}

type frontierBound struct {
	repo repo.Repo
	svc  service.Service
}

func NewFrontierBound(conf *config.Configuration, repo repo.Repo) (FrontierBound, error) {
	fb := &frontierBound{
		repo: repo,
	}

	dialer := func() (net.Conn, error) {
		return net.Dial("tcp", conf.Frontier.Addr)
	}
	svc, err := service.NewService(dialer)
	if err != nil {
		return nil, err
	}
	// 注册frontier回调函数
	err = svc.RegisterGetEdgeID(context.Background(), fb.getID)
	if err != nil {
		return nil, err
	}
	err = svc.RegisterEdgeOnline(context.Background(), fb.online)
	if err != nil {
		return nil, err
	}
	err = svc.RegisterEdgeOffline(context.Background(), fb.offline)
	if err != nil {
		return nil, err
	}
	// 注册liaison函数
	err = svc.Register(context.Background(), "report_device_usage", fb.reportDeviceUsage)
	if err != nil {
		return nil, err
	}

	fb.svc = svc
	return fb, nil
}
