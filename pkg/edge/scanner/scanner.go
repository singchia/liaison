package scanner

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jumboframes/armorigo/log"
	"github.com/singchia/geminio"
	"github.com/singchia/liaison/pkg/edge/frontierbound"
	"github.com/singchia/liaison/pkg/proto"
)

type Scanner interface{}

type scanner struct {
	frontierBound frontierbound.FrontierBound

	// tasks
	mu    sync.Mutex
	tasks map[uint]*proto.ScanApplicationTaskRequest
}

func NewScanner(frontierBound frontierbound.FrontierBound) (Scanner, error) {

	s := &scanner{
		frontierBound: frontierBound,
		tasks:         make(map[uint]*proto.ScanApplicationTaskRequest),
	}

	// 注册函数
	err := frontierBound.RegisterRPCHandler("scan_application", s.scanApplication)
	if err != nil {
		log.Errorf("scanner register func err: %s", err)
		return nil, err
	}
	// 拉取扫描任务
	go s.pullTaskScanApplication()

	return s, nil
}

func (s *scanner) pullTaskScanApplication() {
	edgeID, _ := s.frontierBound.EdgeID()
	request := proto.PullTaskScanApplicationRequest{
		EdgeID: edgeID,
	}
	data, err := json.Marshal(request)
	if err != nil {
		log.Errorf("marshal pull task scan application request error: %s", err)
		return
	}
	req := s.frontierBound.NewRequest(data)
	rsp, err := s.frontierBound.Call(context.Background(), "pull_task_scan_application", req)
	if err != nil {
		log.Errorf("pull task scan application call error: %s", err)
		return
	}
	if rsp.Error() != nil {
		log.Errorf("pull task scan application return error: %s", rsp.Error())
		return
	}
	var response proto.PullTaskScanApplicationResponse
	err = json.Unmarshal(rsp.Data(), &response)
	if err != nil {
		log.Errorf("unmarshal pull task scan application response error: %s", err)
		return
	}
	log.Infof("pull task scan application response: %v", response)

	// 扫描
	err = s.scan(context.Background(), &proto.ScanApplicationTaskRequest{
		TaskID:   response.TaskID,
		Nets:     response.Nets,
		Port:     response.Port,
		Protocol: response.Protocol,
	})
	if err != nil {
		rsp.SetError(err)
		return
	}
}

func (s *scanner) scanApplication(ctx context.Context, req geminio.Request, rsp geminio.Response) {
	log.Infof("scanner scan application: %s", req.Data())
	var task proto.ScanApplicationTaskRequest
	if err := json.Unmarshal(req.Data(), &task); err != nil {
		rsp.SetError(err)
		log.Errorf("unmarshal scan application task request error: %s", err)
		return
	}

	// 扫描
	err := s.scan(ctx, &task)
	if err != nil {
		rsp.SetError(err)
		log.Errorf("scan application error: %s", err)
		return
	}
}

func (s *scanner) scan(ctx context.Context, task *proto.ScanApplicationTaskRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.tasks[task.TaskID]
	if ok {
		log.Errorf("task already exists")
		return errors.New("task already exists")
	}
	s.tasks[task.TaskID] = task

	go func() {
		defer func() {
			s.mu.Lock()
			delete(s.tasks, task.TaskID)
			s.mu.Unlock()
		}()

		running := proto.ScanApplicationTaskResult{
			TaskID:              task.TaskID,
			ScannedApplications: make([]proto.ScannedApplication, 0),
			Status:              "running",
		}
		// 上报扫描中
		data, err := json.Marshal(running)
		if err != nil {
			log.Errorf("marshal report task scan application result error: %s", err)
			return
		}
		reportReq := s.frontierBound.NewRequest(data)
		_, err = s.frontierBound.Call(ctx, "report_task_scan_application", reportReq)
		if err != nil {
			log.Errorf("call report task scan application error: %s", err)
			return
		}
		log.Infof("scan with go start: %v", task)

		// 使用纯 Go 实现的扫描（不依赖 libpcap）
		err = s.scanWithGo(ctx, task, &running)
		if err != nil {
			log.Errorf("scan with go error: %s", err)
			return
		}

		completed := proto.ScanApplicationTaskResult{
			TaskID:              task.TaskID,
			ScannedApplications: running.ScannedApplications,
			Status:              "completed",
		}
		data, err = json.Marshal(completed)
		if err != nil {
			log.Errorf("marshal report task scan application result error: %s", err)
			return
		}
		reportReq = s.frontierBound.NewRequest(data)
		_, err = s.frontierBound.Call(ctx, "report_task_scan_application", reportReq)
		if err != nil {
			log.Errorf("call report task scan application error: %s", err)
			return
		}
		log.Infof("scan with go completed: %v", string(data))
	}()
	return nil
}

func (s *scanner) Close() error {
	return nil
}

// scanWithGo 使用纯 Go 实现的端口扫描（不依赖 libpcap）
func (s *scanner) scanWithGo(ctx context.Context, task *proto.ScanApplicationTaskRequest, result *proto.ScanApplicationTaskResult) error {
	// 获取要扫描的端口列表
	ports, err := s.getPorts(task.Port)
	if err != nil {
		return fmt.Errorf("get ports error: %s", err)
	}

	// 获取要扫描的 IP 列表
	ips, err := s.getIPsFromNets(task.Nets)
	if err != nil {
		return fmt.Errorf("get IPs from nets error: %s", err)
	}

	// 并发扫描
	var wg sync.WaitGroup
	var mu sync.Mutex
	timeout := 2 * time.Second
	maxConcurrency := 100
	sem := make(chan struct{}, maxConcurrency)

	for _, ip := range ips {
		for _, port := range ports {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			wg.Add(1)
			sem <- struct{}{} // 获取信号量

			go func(ip string, port int) {
				defer wg.Done()
				defer func() { <-sem }() // 释放信号量

				var open bool
				protocol := strings.ToLower(task.Protocol)
				switch protocol {
				case "udp":
					open = scanUDP(ip, port, timeout)
				default: // tcp 或其他
					open = scanTCP(ip, port, timeout)
				}

				if open {
					mu.Lock()
					result.ScannedApplications = append(result.ScannedApplications, proto.ScannedApplication{
						IP:       ip,
						Port:     port,
						Protocol: protocol,
					})
					mu.Unlock()
				}
			}(ip, port)
		}
	}

	wg.Wait()
	return nil
}

// getPorts 获取要扫描的端口列表
func (s *scanner) getPorts(port int) ([]int, error) {
	if port > 0 {
		return []int{port}, nil
	}
	// 返回常见端口列表（前100个）
	return getTopPorts(100), nil
}

// getIPsFromNets 从网段列表中获取所有 IP 地址
func (s *scanner) getIPsFromNets(nets []string) ([]string, error) {
	var ips []string
	for _, netStr := range nets {
		// 尝试解析为 CIDR
		ip, ipNet, err := net.ParseCIDR(netStr)
		if err != nil {
			// 如果不是 CIDR，尝试作为单个 IP
			if ip := net.ParseIP(netStr); ip != nil {
				ips = append(ips, netStr)
				continue
			}
			return nil, fmt.Errorf("invalid network format: %s", netStr)
		}

		// 遍历 CIDR 中的所有 IP
		for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
			// 跳过网络地址和广播地址
			if isNetworkOrBroadcast(ip, ipNet) {
				continue
			}
			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}

// inc 增加 IP 地址
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// isNetworkOrBroadcast 检查是否是网络地址或广播地址
func isNetworkOrBroadcast(ip net.IP, ipNet *net.IPNet) bool {
	// 简化处理：对于小网段（/24 及以上），跳过第一个和最后一个
	ones, bits := ipNet.Mask.Size()
	if ones >= 24 && bits == 32 {
		ip4 := ip.To4()
		if ip4 != nil {
			lastOctet := ip4[3]
			return lastOctet == 0 || lastOctet == 255
		}
	}
	return false
}

// getTopPorts 返回常见端口列表
func getTopPorts(count int) []int {
	// 常见端口列表（基于实际使用频率）
	commonPorts := []int{
		22, 23, 25, 53, 80, 110, 111, 135, 139, 143,
		443, 445, 993, 995, 1723, 3306, 3389, 5900, 8080, 8443,
		// 添加更多常见端口
		21, 69, 79, 88, 102, 110, 143, 389, 443, 445,
		636, 993, 995, 1433, 1521, 1723, 3306, 3389, 5432, 5900,
		6000, 8000, 8080, 8443, 8888, 9000, 9090, 9200, 9300, 10000,
		// 继续添加
		7, 9, 13, 17, 19, 20, 21, 22, 23, 25,
		26, 37, 53, 79, 80, 81, 88, 106, 110, 111,
		113, 119, 135, 139, 143, 144, 179, 199, 389, 427,
		443, 444, 445, 465, 514, 515, 543, 544, 548, 554,
		587, 631, 646, 873, 990, 993, 995, 1025, 1026, 1027,
		1028, 1029, 1110, 1433, 1720, 1723, 1755, 1900, 2000, 2049,
		2121, 2717, 3000, 3128, 3306, 3389, 3986, 4899, 5000, 5009,
		5051, 5060, 5101, 5190, 5357, 5432, 5631, 5666, 5800, 5900,
		6000, 6001, 6646, 7070, 8000, 8008, 8009, 8080, 8081, 8443,
		8888, 9100, 9999, 10000, 32768, 49152, 49153, 49154, 49155, 49156,
	}

	if count > len(commonPorts) {
		count = len(commonPorts)
	}
	return commonPorts[:count]
}

func scanTCP(ip string, port int, timeout time.Duration) bool {
	addr := net.JoinHostPort(ip, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	// 尝试读取一点数据（非必须）
	_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	reader := bufio.NewReader(conn)
	_, _ = reader.Peek(1)

	return true
}

func scanUDP(ip string, port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("udp", net.JoinHostPort(ip, strconv.Itoa(port)), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	// 发个探测包
	_ = conn.SetWriteDeadline(time.Now().Add(timeout))
	_, _ = conn.Write([]byte("\n"))

	// 看能不能收到回应
	buf := make([]byte, 1)
	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	_, err = conn.Read(buf)
	if err != nil {
		return false
	}
	return true
}
