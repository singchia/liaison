package entry

import (
	"context"
	"sync"
	"time"

	"github.com/singchia/liaison/pkg/entry/config"
	"github.com/singchia/liaison/pkg/entry/transport"
)

type Entry struct {
	config *config.Configuration

	gatekeeper *transport.Gatekeeper
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

func NewEntry() (*Entry, error) {
	// 初始化配置
	cfg, err := config.Init()
	if err != nil {
		return nil, err
	}

	// 创建端口管理器
	portManager := transport.NewGatekeeper()

	ctx, cancel := context.WithCancel(context.Background())

	entry := &Entry{
		config: cfg,

		gatekeeper: portManager,
		ctx:        ctx,
		cancel:     cancel,
	}

	// 启动配置同步
	entry.wg.Add(1)
	go entry.syncProxyConfigs()

	return entry, nil
}

// syncProxyConfigs 定期从manager同步Proxy配置
func (e *Entry) syncProxyConfigs() {
	defer e.wg.Done()

	ticker := time.NewTicker(5 * time.Second) // 每5秒同步一次
	defer ticker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			return
		case <-ticker.C:
			e.updateProxyConfigs()
		}
	}
}

// updateProxyConfigs 更新代理配置
func (e *Entry) updateProxyConfigs() {

}

func (e *Entry) Close() error {
	e.cancel()
	e.wg.Wait()

	return nil
}
