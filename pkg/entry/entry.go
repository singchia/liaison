package entry

import (
	"context"
	"sync"
	"time"

	"github.com/singchia/liaison/pkg/entry/transport"
)

type Entry struct {
	gatekeeper *transport.Gatekeeper
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

func NewEntry() (*Entry, error) {

	// 创建端口管理器
	gatekeeper := transport.NewGatekeeper()

	ctx, cancel := context.WithCancel(context.Background())

	entry := &Entry{
		gatekeeper: gatekeeper,
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
