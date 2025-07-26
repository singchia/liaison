package frontierbound

import (
	"context"
	"encoding/json"
	"errors"
	"net"

	"github.com/singchia/frontier/api/dataplane/v1/service"
	"github.com/singchia/liaison/pkg/liaison/internal/config"
	"github.com/singchia/liaison/pkg/liaison/internal/repo"
	"github.com/singchia/liaison/pkg/liaison/internal/repo/model"
	"github.com/singchia/liaison/pkg/proto"
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

// 获取设备ID
func (fb *frontierBound) getID(meta []byte) (uint64, error) {
	var m proto.Meta
	if err := json.Unmarshal(meta, &m); err != nil {
		return 0, err
	}
	ak, edge, err := fb.repo.GetEdgeByAccessKey(m.AccessKey)
	if err != nil {
		return 0, err
	}
	if ak.SecretKey != m.SecretKey {
		return 0, errors.New("invalid secret key")
	}

	return uint64(edge.ID), nil
}

func (fb *frontierBound) online(edgeID uint64, meta []byte, addr net.Addr) error {
	err := fb.repo.UpdateEdgeOnlineStatus(edgeID, model.EdgeOnlineStatusOnline)
	if err != nil {
		return err
	}
	return nil
}

func (fb *frontierBound) offline(edgeID uint64, meta []byte, addr net.Addr) error {
	err := fb.repo.UpdateEdgeOnlineStatus(edgeID, model.EdgeOnlineStatusOffline)
	if err != nil {
		return err
	}
	return nil
}
