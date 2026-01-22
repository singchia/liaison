package proxy

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/jumboframes/armorigo/log"
	"github.com/singchia/geminio"
	"github.com/singchia/liaison/pkg/edge/frontierbound"
	"github.com/singchia/liaison/pkg/proto"
)

type Proxy interface{}

type trafficStats struct {
	ProxyID       uint
	ApplicationID uint
	BytesIn       int64
	BytesOut      int64
}

type proxy struct {
	frontierBound frontierbound.FrontierBound
	mu            sync.RWMutex
	// key: "proxyID:applicationID", value: traffic stats
	stats map[string]*trafficStats
	stop  chan struct{}
}

func NewProxy(frontierBound frontierbound.FrontierBound) (Proxy, error) {
	proxy := &proxy{
		frontierBound: frontierBound,
		stats:         make(map[string]*trafficStats),
		stop:          make(chan struct{}),
	}

	proxy.frontierBound.RegisterStreamHandler(proxy.proxy)

	// 启动定时上报任务（每分钟上报一次）
	go proxy.reportLoop()

	return proxy, nil
}

func (p *proxy) proxy(ctx context.Context, stream geminio.Stream) {
	// 读取前4个字节获取meta长度
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(stream, lengthBuf)
	if err != nil {
		log.Errorf("proxy stream read meta length err: %s", err)
		return
	}
	length := binary.BigEndian.Uint32(lengthBuf)
	dataBuf := make([]byte, length)
	_, err = io.ReadFull(stream, dataBuf)
	if err != nil {
		log.Errorf("proxy stream read meta data err: %s", err)
		return
	}

	var dst proto.Dst
	if err := json.Unmarshal(dataBuf, &dst); err != nil {
		log.Errorf("proxy stream meta unmarshal err: %s", err)
		return
	}

	conn, err := net.Dial("tcp", dst.Addr)
	if err != nil {
		log.Errorf("proxy stream dial err: %s", err)
		return
	}

	// 流量统计
	var bytesIn, bytesOut int64

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()

		// 使用自定义的 Copy 来统计流量
		// 从 stream 读取到 conn，这是入站流量（从客户端到服务器）
		n, err := copyWithStats(conn, stream, &bytesIn)
		if err != nil && !IsErrClosed(err) {
			log.Errorf("read stream, src: %s, dst: %s; to conn, src: %s, dst: %s; err: %s",
				stream.RemoteAddr(), stream.LocalAddr(), conn.LocalAddr(), conn.RemoteAddr(), err)
		}
		_ = stream.Close()
		_ = conn.Close()
		_ = n // 避免未使用变量警告
	}()

	go func() {
		defer wg.Done()

		// 从 conn 读取到 stream，这是出站流量（从服务器到客户端）
		n, err := copyWithStats(stream, conn, &bytesOut)
		if err != nil && !IsErrClosed(err) {
			log.Errorf("read conn, src: %s, dst: %s; to stream, src: %s, dst: %s; err: %s",
				conn.LocalAddr(), conn.RemoteAddr(), stream.RemoteAddr(), stream.LocalAddr(), err)
		}
		_ = stream.Close()
		_ = conn.Close()
		_ = n // 避免未使用变量警告
	}()

	wg.Wait()

	// 累积流量统计（如果有 ApplicationID 和 ProxyID）
	// 不再立即上报，而是累积到stats中，由定时器每分钟上报一次
	if dst.ApplicationID > 0 && dst.ProxyID > 0 {
		p.recordTraffic(dst.ProxyID, dst.ApplicationID, bytesIn, bytesOut)
	}
}

func IsErrClosed(err error) bool {
	if strings.Contains(err.Error(), net.ErrClosed.Error()) {
		return true
	}
	return false
}

// copyWithStats 带流量统计的 io.Copy
func copyWithStats(dst io.Writer, src io.Reader, bytes *int64) (written int64, err error) {
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				*bytes += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		written += int64(nr)
	}
	return written, err
}

// recordTraffic 记录流量（线程安全，累积到stats中）
func (p *proxy) recordTraffic(proxyID, applicationID uint, bytesIn, bytesOut int64) {
	if bytesIn == 0 && bytesOut == 0 {
		return // 没有流量，不需要记录
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	key := trafficKey(proxyID, applicationID)
	stats, exists := p.stats[key]
	if !exists {
		stats = &trafficStats{
			ProxyID:       proxyID,
			ApplicationID: applicationID,
		}
		p.stats[key] = stats
	}

	stats.BytesIn += bytesIn
	stats.BytesOut += bytesOut
}

// reportLoop 每分钟上报一次流量统计
func (p *proxy) reportLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.flushAndReport()
		case <-p.stop:
			// 退出前最后一次上报
			p.flushAndReport()
			return
		}
	}
}

// flushAndReport 上报并清零流量统计
func (p *proxy) flushAndReport() {
	p.mu.Lock()
	if len(p.stats) == 0 {
		p.mu.Unlock()
		return
	}

	// 复制统计数据
	statsToReport := make([]*trafficStats, 0, len(p.stats))
	for _, stats := range p.stats {
		statsToReport = append(statsToReport, &trafficStats{
			ProxyID:       stats.ProxyID,
			ApplicationID: stats.ApplicationID,
			BytesIn:       stats.BytesIn,
			BytesOut:      stats.BytesOut,
		})
	}

	// 清空统计数据
	p.stats = make(map[string]*trafficStats)
	p.mu.Unlock()

	// 上报（在锁外执行，避免阻塞）
	ctx := context.Background()
	for _, stats := range statsToReport {
		p.reportTraffic(ctx, stats.ProxyID, stats.ApplicationID, stats.BytesIn, stats.BytesOut)
	}

	log.Debugf("reported %d traffic metrics", len(statsToReport))
}

// reportTraffic 上报流量统计到 manager
func (p *proxy) reportTraffic(ctx context.Context, proxyID, applicationID uint, bytesIn, bytesOut int64) {
	if bytesIn == 0 && bytesOut == 0 {
		return // 没有流量，不需要上报
	}

	// 通过 RPC 上报流量
	log.Debugf("traffic stats: proxyID=%d, applicationID=%d, bytesIn=%d, bytesOut=%d",
		proxyID, applicationID, bytesIn, bytesOut)

	// 通过 RPC 发送到 manager 的流量统计器
	data, err := json.Marshal(proto.ReportTrafficMetricRequest{
		ProxyID:       proxyID,
		ApplicationID: applicationID,
		BytesIn:       bytesIn,
		BytesOut:      bytesOut,
	})
	if err != nil {
		log.Errorf("marshal traffic metric request error: %s", err)
		return
	}

	req := p.frontierBound.NewRequest(data)
	_, err = p.frontierBound.Call(ctx, "report_traffic_metric", req)
	if err != nil {
		log.Errorf("call report traffic metric error: %s", err)
		return
	}
}

func trafficKey(proxyID, applicationID uint) string {
	return fmt.Sprintf("%d:%d", proxyID, applicationID)
}
