package transport

import (
	"context"
	"net"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// proxy 端口监听器
type proxy struct {
	port     int
	listener net.Listener
	wg       sync.WaitGroup
}

// acceptConnections 接受连接并处理
func (proxy *proxy) accept(ctx context.Context) {
	defer proxy.wg.Done()

	logrus.Infof("listening on port %d", proxy.port)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := proxy.listener.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					return
				}
				select {
				case <-ctx.Done():
					return
				default:
					logrus.Errorf("failed to accept connection on port %d: %s", proxy.port, err)
					continue
				}
			}

			// 处理连接
			proxy.wg.Add(1)
			go proxy.handleConnection(conn)
		}
	}
}

// handleConnection 处理单个连接
func (p *proxy) handleConnection(conn net.Conn) {
	defer p.wg.Done()
	defer conn.Close()

	logrus.Debugf("new connection from %s to port %d", conn.RemoteAddr(), p.port)

	logrus.Infof("connection from %s forwarded successfully", conn.RemoteAddr())
}

// Close 关闭监听器
func (p *proxy) close() {
	p.listener.Close()
	p.wg.Wait()
}
