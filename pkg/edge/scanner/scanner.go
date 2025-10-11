package scanner

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/singchia/geminio"
	"github.com/singchia/liaison/pkg/edge/frontierbound"
	"github.com/singchia/liaison/pkg/proto"
	"github.com/sirupsen/logrus"
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
	}

	// 注册函数
	err := frontierBound.RegisterRPCHandler("scan_application", s.scanApplication)
	if err != nil {
		logrus.Errorf("scanner register func err: %s", err)
		return nil, err
	}

	return s, nil
}

func (s *scanner) scanApplication(ctx context.Context, req geminio.Request, rsp geminio.Response) {
	var task proto.ScanApplicationTaskRequest
	if err := json.Unmarshal(req.Data(), &task); err != nil {
		rsp.SetError(err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.tasks[task.TaskID]
	if ok {
		rsp.SetError(errors.New("task already exists"))
		return
	}
	s.tasks[task.TaskID] = &task

	go func() {
		defer func() {
			s.mu.Lock()
			delete(s.tasks, task.TaskID)
			s.mu.Unlock()
		}()

		result := proto.ScanApplicationTaskResult{
			TaskID:              task.TaskID,
			ScannedApplications: make([]proto.ScannedApplication, 0),
			Status:              "running",
		}
		// 扫描中
		data, err := json.Marshal(result)
		if err != nil {
			return
		}
		s.frontierBound.NewRequest(data)
		_, err = s.frontierBound.Call(ctx, "report_task_scan_application", req)
		if err != nil {
			return
		}

		for _, ip := range task.Nets {
			open := false
			switch strings.ToLower(task.Protocol) {
			case "udp":
				open = scanUDP(ip, task.Port, 2*time.Second)
			default:
				open = scanTCP(ip, task.Port, 2*time.Second)
			}

			if open {
				result.ScannedApplications = append(result.ScannedApplications, proto.ScannedApplication{
					IP:       ip,
					Port:     task.Port,
					Protocol: task.Protocol,
				})
			}
		}

		result.Status = "completed"

		// 序列化结果
		data, err = json.Marshal(result)
		if err != nil {
			return
		}
		reportReq := s.frontierBound.NewRequest(data)
		_, err = s.frontierBound.Call(ctx, "report_task_scan_application", reportReq)
		if err != nil {
			return
		}
	}()

}

func (s *scanner) Close() error {
	return nil
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
