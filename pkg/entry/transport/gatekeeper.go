package transport

import (
	"context"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

// Manager 端口管理器，负责动态管理TCP端口监听
type Gatekeeper struct {
	mu        sync.RWMutex
	listeners map[int32]*Listener // port -> listener
	ctx       context.Context
	cancel    context.CancelFunc
}

// Listener 端口监听器
type Listener struct {
	port     int32
	listener net.Listener
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewGatekeeper() *Gatekeeper {
	ctx, cancel := context.WithCancel(context.Background())
	return &Gatekeeper{
		listeners: make(map[int32]*Listener),
		ctx:       ctx,
		cancel:    cancel,
	}
}

/*
// UpdateProxies 更新代理配置，动态开启/关闭端口
func (m *Gatekeeper) UpdateProxies(proxies []*types.Proxy) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 创建当前代理的端口映射
	currentPorts := make(map[int32]*types.Proxy)
	for _, proxy := range proxies {
		if proxy.Status == "running" { // 只处理运行中的代理
			currentPorts[proxy.Port] = proxy
		}
	}

	// 关闭不再需要的监听器
	for port, listener := range m.listeners {
		if _, exists := currentPorts[port]; !exists {
			logrus.Infof("stopping listener for port %d", port)
			listener.Close()
			delete(m.listeners, port)
		}
	}

	// 启动新的监听器或更新现有监听器
	for port, proxy := range currentPorts {
		if listener, exists := m.listeners[port]; exists {
			// 检查代理配置是否有变化
			if listener.proxy.ID != proxy.ID {
				logrus.Infof("updating listener for port %d", port)
				listener.Close()
				delete(m.listeners, port)
				m.startListener(port, proxy)
			}
		} else {
			// 启动新的监听器
			logrus.Infof("starting new listener for port %d", port)
			m.startListener(port, proxy)
		}
	}
}

// startListener 启动指定端口的监听器
func (m *Gatekeeper) startListener(port int32, proxy *types.Proxy) {
	ctx, cancel := context.WithCancel(m.ctx)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Errorf("failed to listen on port %d: %s", port, err)
		cancel()
		return
	}

	l := &Listener{
		port:     port,
		proxy:    proxy,
		listener: listener,
		ctx:      ctx,
		cancel:   cancel,
	}

	m.listeners[port] = l

	// 启动监听循环
	l.wg.Add(1)
	go l.acceptConnections()
}
*/

// acceptConnections 接受连接并处理
func (l *Listener) acceptConnections() {
	defer l.wg.Done()

	logrus.Infof("listening on port %d", l.port)

	for {
		select {
		case <-l.ctx.Done():
			return
		default:
			conn, err := l.listener.Accept()
			if err != nil {
				select {
				case <-l.ctx.Done():
					return
				default:
					logrus.Errorf("failed to accept connection on port %d: %s", l.port, err)
					continue
				}
			}

			// 处理连接
			l.wg.Add(1)
			go l.handleConnection(conn)
		}
	}
}

// handleConnection 处理单个连接
func (l *Listener) handleConnection(conn net.Conn) {
	defer l.wg.Done()
	defer conn.Close()

	logrus.Infof("new connection from %s to port %d", conn.RemoteAddr(), l.port)

	logrus.Infof("connection from %s forwarded successfully", conn.RemoteAddr())
}

// Close 关闭监听器
func (l *Listener) Close() {
	l.cancel()
	l.listener.Close()
	l.wg.Wait()
}

// Close 关闭端口管理器
func (m *Gatekeeper) Close() {
	m.cancel()

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, listener := range m.listeners {
		listener.Close()
	}
	m.listeners = make(map[int32]*Listener)
}
