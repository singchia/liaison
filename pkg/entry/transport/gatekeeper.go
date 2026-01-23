package transport

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jumboframes/armorigo/log"
	"github.com/jumboframes/armorigo/rproxy"
	"github.com/singchia/liaison/pkg/entry/frontierbound"
	"github.com/singchia/liaison/pkg/lerrors"
	"github.com/singchia/liaison/pkg/proto"
)

// Gatekeeper 端口管理器，负责动态管理TCP端口监听
type Gatekeeper struct {
	mu             sync.RWMutex
	proxies        map[int]*proxy // id -> listener
	proxiesIdxPort map[int]int    // port -> id
	// proxy ID -> application ID 映射（用于流量统计）
	proxyAppMap map[int]uint

	// frontier
	frontierBound frontierbound.FrontierBound
	// 流量统计器（可选，如果设置了则统计流量）
	trafficCollector interface {
		RecordTraffic(proxyID, applicationID uint, bytesIn, bytesOut int64)
	}
	// 流量统计数据（每分钟上报一次）
	trafficStats map[string]*trafficStats // key: "proxyID:applicationID"
	stop         chan struct{}
}

type trafficStats struct {
	ProxyID       uint
	ApplicationID uint
	BytesIn       int64
	BytesOut      int64
}

func NewGatekeeper(frontierBound frontierbound.FrontierBound) *Gatekeeper {
	gk := &Gatekeeper{
		proxies:        make(map[int]*proxy),
		proxiesIdxPort: make(map[int]int),
		proxyAppMap:    make(map[int]uint),
		frontierBound:  frontierBound,
		trafficStats:   make(map[string]*trafficStats),
		stop:           make(chan struct{}),
	}
	// 启动定时上报任务（每分钟上报一次）
	go gk.reportLoop()
	return gk
}

// SetTrafficCollector 设置流量统计器
func (m *Gatekeeper) SetTrafficCollector(collector interface {
	RecordTraffic(proxyID, applicationID uint, bytesIn, bytesOut int64)
}) {
	m.trafficCollector = collector
}

func (m *Gatekeeper) CreateProxy(ctx context.Context, protoproxy *proto.Proxy) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已存在该代理
	_, exists := m.proxies[protoproxy.ID]
	if exists {
		log.Warnf("port %d is already in use", protoproxy.ProxyPort)
		return nil
	}
	// 如果端口为0，系统会自动分配端口
	requestedPort := protoproxy.ProxyPort
	if requestedPort == 0 {
		// 端口为0时，不检查冲突，因为系统会自动分配
	} else {
		// 检查端口是否和其他代理冲突
		id, exists := m.proxiesIdxPort[requestedPort]
		if exists && id != protoproxy.ID {
			log.Errorf("port %d conflict with proxy %d", requestedPort, id)
			return lerrors.ErrPortConflict
		}
	}

	// 监听
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", requestedPort))
	if err != nil {
		log.Errorf("failed to listen on port %d: %s", requestedPort, err)
		return err
	}
	
	// 获取实际监听的端口（如果端口为0，系统会分配一个端口）
	actualPort := listener.Addr().(*net.TCPAddr).Port
	if requestedPort == 0 {
		log.Infof("proxy %d: system allocated port %d", protoproxy.ID, actualPort)
		// 更新 protoproxy 的端口，以便后续使用
		protoproxy.ProxyPort = actualPort
	}
	// hook 函数
	postAccept := func(clientAddr net.Addr, _ net.Addr) (custom interface{}, err error) {
		pc := &proxyContext{
			edgeID:        protoproxy.EdgeID,
			dst:           protoproxy.Dst,
			applicationID: protoproxy.ApplicationID,
			proxyID:       uint(protoproxy.ID),
			gatekeeper:    m,
		}
		return pc, nil
	}
	proxyDial := func(dst net.Addr, custom interface{}) (target net.Conn, err error) {
		pc := custom.(*proxyContext)
		stream, err := m.frontierBound.OpenStream(context.TODO(), pc.edgeID)
		if err != nil {
			return nil, err
		}
		// 包装stream连接以统计流量
		return newCountingConn(stream, pc), nil
	}
	preWrite := func(writer io.Writer, custom interface{}) error {
		pc := custom.(*proxyContext)
		dst := proto.Dst{
			Addr:          pc.dst,
			ApplicationID: pc.applicationID,
			ProxyID:       pc.proxyID,
		}
		data, err := json.Marshal(dst)
		if err != nil {
			log.Errorf("failed to marshal dst: %s", err)
			return err
		}
		lengthBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(lengthBuf, uint32(len(data)))
		_, err = writer.Write(lengthBuf)
		if err != nil {
			log.Errorf("failed to write dst length: %s", err)
			return err
		}
		_, err = writer.Write(data)
		if err != nil {
			log.Errorf("failed to write dst: %s", err)
			return err
		}
		return nil
	}

	rp, err := rproxy.NewRProxy(listener,
		rproxy.OptionRProxyPostAccept(postAccept),
		rproxy.OptionRProxyDial(proxyDial),
		rproxy.OptionRProxyPreWrite(preWrite))
	if err != nil {
		log.Errorf("failed to create rproxy: %s", err)
		return err
	}

	go rp.Proxy(context.Background())

	p := &proxy{
		port: actualPort, // 使用实际端口
		rp:   rp,
	}
	m.proxies[protoproxy.ID] = p
	m.proxiesIdxPort[actualPort] = protoproxy.ID
	// 保存 proxy ID 和 application ID 的映射
	if protoproxy.ApplicationID > 0 {
		m.proxyAppMap[protoproxy.ID] = protoproxy.ApplicationID
	}

	return nil
}

func (m *Gatekeeper) DeleteProxy(ctx context.Context, id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查端口是否存在
	p, exists := m.proxies[id]
	if !exists {
		log.Warnf("proxy %d not found", id)
		return nil
	}

	// 关闭监听器
	p.rp.Close()

	// 删除映射
	delete(m.proxies, id)
	delete(m.proxiesIdxPort, p.port)
	delete(m.proxyAppMap, id)

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
	edgeID        uint64
	dst           string
	applicationID uint
	proxyID       uint
	// 流量统计
	bytesIn  int64 // 入站流量（从客户端到服务器）
	bytesOut int64 // 出站流量（从服务器到客户端）
	// 用于记录流量到collector
	gatekeeper *Gatekeeper
	// 连接关闭标记
	closed int32
}

type proxy struct {
	port int
	rp   *rproxy.RProxy
}

// recordTraffic 记录流量（累积到stats中，由定时器每分钟上报一次）
func (m *Gatekeeper) recordTraffic(proxyID, applicationID uint, bytesIn, bytesOut int64) {
	if bytesIn == 0 && bytesOut == 0 {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%d:%d", proxyID, applicationID)
	stats, exists := m.trafficStats[key]
	if !exists {
		stats = &trafficStats{
			ProxyID:       proxyID,
			ApplicationID: applicationID,
		}
		m.trafficStats[key] = stats
	}

	stats.BytesIn += bytesIn
	stats.BytesOut += bytesOut
}

// reportLoop 每分钟上报一次流量统计
func (m *Gatekeeper) reportLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.flushAndReport()
		case <-m.stop:
			// 退出前最后一次上报
			m.flushAndReport()
			return
		}
	}
}

// flushAndReport 上报并清零流量统计
func (m *Gatekeeper) flushAndReport() {
	m.mu.Lock()
	if len(m.trafficStats) == 0 {
		m.mu.Unlock()
		return
	}

	// 复制统计数据
	statsToReport := make([]*trafficStats, 0, len(m.trafficStats))
	for _, stats := range m.trafficStats {
		statsToReport = append(statsToReport, &trafficStats{
			ProxyID:       stats.ProxyID,
			ApplicationID: stats.ApplicationID,
			BytesIn:       stats.BytesIn,
			BytesOut:      stats.BytesOut,
		})
	}

	// 清空统计数据
	m.trafficStats = make(map[string]*trafficStats)
	m.mu.Unlock()

	// 上报（在锁外执行，避免阻塞）
	if m.trafficCollector != nil {
		for _, stats := range statsToReport {
			m.trafficCollector.RecordTraffic(stats.ProxyID, stats.ApplicationID, stats.BytesIn, stats.BytesOut)
		}
		log.Debugf("reported %d traffic metrics", len(statsToReport))
	}
}

// countingConn 包装net.Conn以统计流量
// rproxy内部会进行双向数据复制：
// 1. 从客户端读取 -> 写入stream（入站流量，通过stream.Write统计）
// 2. 从stream读取 -> 写入客户端（出站流量，通过stream.Read统计）
type countingConn struct {
	net.Conn
	pc *proxyContext
}

func newCountingConn(conn net.Conn, pc *proxyContext) *countingConn {
	return &countingConn{
		Conn: conn,
		pc:   pc,
	}
}

func (c *countingConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	if n > 0 {
		// 从stream读取，是出站流量（从服务器到客户端）
		atomic.AddInt64(&c.pc.bytesOut, int64(n))
		// 实时累积到gatekeeper的stats中（不等待连接关闭）
		c.pc.gatekeeper.recordTraffic(c.pc.proxyID, c.pc.applicationID, 0, int64(n))
	}
	return n, err
}

func (c *countingConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if n > 0 {
		// 向stream写入，是入站流量（从客户端到服务器）
		atomic.AddInt64(&c.pc.bytesIn, int64(n))
		// 实时累积到gatekeeper的stats中（不等待连接关闭）
		c.pc.gatekeeper.recordTraffic(c.pc.proxyID, c.pc.applicationID, int64(n), 0)
	}
	return n, err
}

func (c *countingConn) Close() error {
	return c.Conn.Close()
}
