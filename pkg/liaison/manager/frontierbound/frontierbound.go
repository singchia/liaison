package frontierbound

import (
	"context"
	"errors"
	"math/rand"
	"net"

	"github.com/singchia/frontier/api/dataplane/v1/service"
	"github.com/singchia/liaison/pkg/liaison/config"
	"github.com/singchia/liaison/pkg/liaison/repo"
	"github.com/singchia/liaison/pkg/utils"
)

type FrontierBound interface {
	EmitScanApplications(ctx context.Context, taskID uint, edgeID uint64, net *Net) error
}

type frontierBound struct {
	repo repo.Repo
	svc  service.Service
}

func NewFrontierBound(conf *config.Configuration, repo repo.Repo) (FrontierBound, error) {
	dial := conf.Frontier.Dial
	if len(dial.Addrs) == 0 {
		return nil, errors.New("dial addr is empty")
	}
	fb := &frontierBound{
		repo: repo,
	}

	dialer := func() (net.Conn, error) {
		return utils.Dial(&dial, rand.Intn(len(dial.Addrs)))
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
	err = svc.Register(context.Background(), "report_device", fb.reportDevice)
	if err != nil {
		return nil, err
	}
	err = svc.Register(context.Background(), "report_device_usage", fb.reportDeviceUsage)
	if err != nil {
		return nil, err
	}
	err = svc.Register(context.Background(), "report_edge", fb.reportEdge)
	if err != nil {
		return nil, err
	}
	err = svc.Register(context.Background(), "report_task_scan_application", fb.reportTaskScanApplication)
	if err != nil {
		return nil, err
	}

	fb.svc = svc
	return fb, nil
}

func (fb *frontierBound) Close() error {
	return fb.svc.Close()
}
