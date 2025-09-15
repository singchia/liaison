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

	// frontier
	frontier frontierbound.FrontierBound
}

func NewGatekeeper() *Gatekeeper {
	return &Gatekeeper{
		proxies: make(map[int]*proxy),
	}
}

func (m *Gatekeeper) CreateProxy(ctx context.Context, protoproxy *proto.Proxy) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已存在该代理
	_, exists := m.proxies[protoproxy.ID]
	if exists {
		logrus.Warnf("port %d is already in use", protoproxy.ProxyPort)
		return nil
	}
	// 检查端口是否和其他代理冲突
	id, exists := m.proxiesIdxPort[protoproxy.ProxyPort]
	if exists && id != protoproxy.ID {
		logrus.Errorf("port %d conflict with proxy %d", protoproxy.ProxyPort, id)
		return lerrors.ErrPortConflict
	}

	// 监听
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

	p := &proxy{
		port: protoproxy.ProxyPort,
		rp:   rp,
	}
	m.proxies[protoproxy.ID] = p
	m.proxiesIdxPort[protoproxy.ProxyPort] = protoproxy.ID

	return nil
}

func (m *Gatekeeper) DeleteProxy(ctx context.Context, id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查端口是否存在
	p, exists := m.proxies[id]
	if !exists {
		logrus.Warnf("proxy %d not found", id)
		return nil
	}

	// 关闭监听器
	p.rp.Close()

	// 删除映射
	delete(m.proxies, id)
	delete(m.proxiesIdxPort, p.port)

	return nil
}

// Close 关闭端口管理器
func (m *Gatekeeper) Close() {

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, p := range m.proxies {
		p.rp.Close()
	}
	m.proxies = make(map[int]*proxy)
}

type proxyContext struct {
	edgeID uint64
	dst    string
}

type proxy struct {
	port int
	rp   *rproxy.RProxy
}
