package entry

import (
	"context"
	"errors"
	"fmt"

	"github.com/jumboframes/armorigo/log"
	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/entry/frontierbound"
	"github.com/singchia/liaison/pkg/entry/transport"
	"github.com/singchia/liaison/pkg/liaison/config"
	"github.com/singchia/liaison/pkg/liaison/manager/controlplane"
	"github.com/singchia/liaison/pkg/proto"
)

type Entry struct {
	gatekeeper *transport.Gatekeeper
	// liaison manager
	manager controlplane.ControlPlane
}

func NewEntry(conf *config.Configuration, manager controlplane.ControlPlane) (*Entry, error) {

	frontierBound, err := frontierbound.NewFrontierBound(conf)
	if err != nil {
		return nil, err
	}

	// 创建端口管理器
	gatekeeper := transport.NewGatekeeper(frontierBound)
	manager.RegisterProxyManager(gatekeeper)

	entry := &Entry{
		gatekeeper: gatekeeper,
		manager:    manager,
	}

	err = entry.pullProxyConfigs()
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
		log.Errorf("failed to list proxies: %s", err)
		return err
	}
	if rsp.Code != 200 {
		log.Errorf("failed to list proxies: %s", rsp.Message)
		return errors.New(rsp.Message)
	}
	log.Infof("list proxies: %s", rsp.Data.GetProxies())

	for _, proxy := range rsp.Data.GetProxies() {
		application := proxy.Application
		if application == nil {
			log.Warnf("proxy %d (name: %s) has no associated application, skipping", proxy.Id, proxy.Name)
			continue
		}
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
