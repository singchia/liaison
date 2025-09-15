package entry

import (
	"github.com/singchia/liaison/pkg/entry/transport"
	"github.com/singchia/liaison/pkg/liaison/manager/controlplane"
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

	entry.pullProxyConfigs()

	return entry, nil
}

// pullProxyConfigs 定期从manager同步Proxy配置
func (e *Entry) pullProxyConfigs() {

}

func (e *Entry) Close() error {
	return nil
}
