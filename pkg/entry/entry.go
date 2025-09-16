package entry

import (
	"context"
	"errors"
	"fmt"

	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/entry/transport"
	"github.com/singchia/liaison/pkg/liaison/manager/controlplane"
	"github.com/singchia/liaison/pkg/proto"
)

type Entry struct {
	gatekeeper *transport.Gatekeeper
	// liaison manager
	manager controlplane.ControlPlane
}

func NewEntry(manager controlplane.ControlPlane) (*Entry, error) {

	// 创建端口管理器
	gatekeeper := transport.NewGatekeeper()
	manager.RegisterProxyManager(gatekeeper)

	entry := &Entry{
		gatekeeper: gatekeeper,
		manager:    manager,
	}

	err := entry.pullProxyConfigs()
	if err != nil {
		return nil, err
	}

	return entry, nil
}

// pullProxyConfigs 定期从manager同步Proxy配置
func (e *Entry) pullProxyConfigs() error {
	rsp, err := e.manager.ListProxies(context.Background(), &v1.ListProxiesRequest{
		Page:     -1,
		PageSize: -1,
	})
	if err != nil {
		return err
	}
	if rsp.Code != 200 {
		return errors.New(rsp.Message)
	}

	for _, proxy := range rsp.Data.GetProxies() {
		application := proxy.Application
		dst := fmt.Sprintf("%s:%d", application.Ip, application.Port)
		e.gatekeeper.CreateProxy(context.Background(), &proto.Proxy{
			ID:        int(proxy.Id),
			Name:      proxy.Name,
			ProxyPort: int(proxy.Port),
			EdgeID:    application.EdgeId,
			Dst:       dst,
		})
	}

	return nil
}

func (e *Entry) Close() error {
	return nil
}
