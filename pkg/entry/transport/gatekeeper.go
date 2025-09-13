package transport

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/singchia/liaison/pkg/lerrors"
	"github.com/singchia/liaison/pkg/proto"
	"github.com/sirupsen/logrus"
)

// Gatekeeper 端口管理器，负责动态管理TCP端口监听
type Gatekeeper struct {
	mu             sync.RWMutex
	proxies        map[int]*proxy // id -> listener
	proxiesIdxPort map[int]int    // port -> id
}

func NewGatekeeper() *Gatekeeper {
	return &Gatekeeper{
		proxies: make(map[int]*proxy),
	}
}

func (m *Gatekeeper) CreateProxy(protoproxy *proto.Proxy) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查端口是否已存在
	p, exists := m.proxies[protoproxy.ID]
	if exists {

		if protoproxy.ProxyPort == p.port {
			logrus.Warnf("port %d is already in use", protoproxy.ProxyPort)
			return nil
		}
	}
	// 检查端口是否和其他代理冲突
	id, exists := m.proxiesIdxPort[protoproxy.ProxyPort]
	if exists && id != protoproxy.ID {
		logrus.Errorf("port %d conflict with proxy %d", protoproxy.ProxyPort, id)
		return lerrors.ErrPortConflict
	}

	// 启动新的监听器
	m.startListener(protoproxy)
	return nil
}

func (m *Gatekeeper) DeleteProxy(protoproxy *proto.Proxy) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查端口是否存在
	p, exists := m.proxies[protoproxy.ID]
	if !exists {
		logrus.Warnf("proxy %d not found", protoproxy.ID)
		return nil
	}

	// 关闭监听器
	p.close()

	// 删除映射
	delete(m.proxies, protoproxy.ID)
	delete(m.proxiesIdxPort, protoproxy.ProxyPort)

	return nil
}

// startListener 启动指定端口的监听器
func (m *Gatekeeper) startListener(protoproxy *proto.Proxy) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", protoproxy.ProxyPort))
	if err != nil {
		logrus.Errorf("failed to listen on port %d: %s", protoproxy.ProxyPort, err)
		return
	}

	p := &proxy{
		port:     protoproxy.ProxyPort,
		listener: listener,
	}

	m.proxies[protoproxy.ID] = p

	go p.accept(context.Background())
}

// Close 关闭端口管理器
func (m *Gatekeeper) Close() {

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, p := range m.proxies {
		p.close()
	}
	m.proxies = make(map[int]*proxy)
}
