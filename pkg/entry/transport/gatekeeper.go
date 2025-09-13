package transport

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/jumboframes/armorigo/rproxy"
	"github.com/singchia/liaison/pkg/entry/frontierbound"
	"github.com/singchia/liaison/pkg/lerrors"
	"github.com/singchia/liaison/pkg/proto"
	"github.com/sirupsen/logrus"
)

// Gatekeeper 端口管理器，负责动态管理TCP端口监听
type Gatekeeper struct {
	mu             sync.RWMutex
	proxies        map[int]*proxy // id -> listener
	proxiesIdxPort map[int]int    // port -> id

	rp *rproxy.RProxy

	// frontier
	frontier frontierbound.FrontierBound
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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", protoproxy.ProxyPort))
	if err != nil {
		logrus.Errorf("failed to listen on port %d: %s", protoproxy.ProxyPort, err)
		return err
	}
	// hook 函数
	postAccept := func(_ net.Addr, _ net.Addr) (custom interface{}, err error) {
		pc := proxyContext{
			edgeID: protoproxy.EdgeID,
			dst:    protoproxy.Dst,
		}
		return &pc, nil
	}
	proxyDial := func(dst net.Addr, custom interface{}) (target net.Conn, err error) {
		pc := custom.(*proxyContext)
		return m.frontier.OpenStream(context.TODO(), pc.edgeID)
	}
	preWrite := func(writer io.Writer, custom interface{}) error {
		pc := custom.(*proxyContext)
		dst := proto.Dst{
			Addr: pc.dst,
		}
		data, err := json.Marshal(dst)
		if err != nil {
			logrus.Errorf("failed to marshal dst: %s", err)
			return err
		}
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(len(data)))
		_, err = writer.Write(buf)
		if err != nil {
			logrus.Errorf("failed to write dst length: %s", err)
			return err
		}
		_, err = writer.Write(data)
		if err != nil {
			logrus.Errorf("failed to write dst: %s", err)
			return err
		}
		return nil
	}

	rp, err := rproxy.NewRProxy(listener,
		rproxy.OptionRProxyPostAccept(postAccept),
		rproxy.OptionRProxyDial(proxyDial),
		rproxy.OptionRProxyPreWrite(preWrite))
	if err != nil {
		logrus.Errorf("failed to create rproxy: %s", err)
		return err
	}

	go rp.Proxy(context.Background())
	return nil
}

type proxyContext struct {
	edgeID uint64
	dst    string
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
