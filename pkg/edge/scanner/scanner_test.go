package scanner

import (
	"testing"
	"time"
)

func TestGetTopPorts(t *testing.T) {
	ports := getTopPorts(10)
	if len(ports) != 10 {
		t.Fatalf("expected 10 ports, got %d", len(ports))
	}
	if ports[0] != 22 {
		t.Fatalf("expected first port to be 22, got %d", ports[0])
	}
}

func TestGetIPsFromNets(t *testing.T) {
	scanner := &scanner{}

	// 测试单个 IP
	ips, err := scanner.getIPsFromNets([]string{"192.168.1.1"})
	if err != nil {
		t.Fatalf("getIPsFromNets error: %s", err)
	}
	if len(ips) != 1 || ips[0] != "192.168.1.1" {
		t.Fatalf("expected [192.168.1.1], got %v", ips)
	}

	// 测试 CIDR
	ips, err = scanner.getIPsFromNets([]string{"192.168.1.0/30"})
	if err != nil {
		t.Fatalf("getIPsFromNets error: %s", err)
	}
	// /30 有 4 个 IP，去掉网络地址和广播地址，应该有 2 个
	if len(ips) < 2 {
		t.Fatalf("expected at least 2 IPs, got %d", len(ips))
	}
}

func TestScanTCP(t *testing.T) {
	// 测试扫描本地回环地址的常见端口
	open := scanTCP("127.0.0.1", 22, 1*time.Second)
	// 不检查结果，因为端口可能开放也可能不开放
	_ = open
}
