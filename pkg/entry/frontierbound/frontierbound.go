package frontierbound

import (
	"context"
	"math/rand"
	"net"

	"github.com/singchia/frontier/api/dataplane/v1/service"
	"github.com/singchia/geminio"
	"github.com/singchia/liaison/pkg/liaison/config"
	"github.com/singchia/liaison/pkg/utils"
)

type FrontierBound interface {
	OpenStream(ctx context.Context, edgeID uint64) (geminio.Stream, error)
	Close() error
}

// 这是edge向frontier注册的连接
type frontierBound struct {
	svc service.Service
}

func NewFrontierBound(conf *config.Configuration) *frontierBound {
	dial := conf.Frontier.Dial
	if len(dial.Addrs) == 0 {
		return nil
	}
	dialer := func() (net.Conn, error) {
		return utils.Dial(&dial, rand.Intn(len(dial.Addrs)))
	}
	svc, err := service.NewService(dialer)
	if err != nil {
		return nil
	}
	return &frontierBound{
		svc: svc,
	}
}

func (fb *frontierBound) Close() error {
	return fb.svc.Close()
}
