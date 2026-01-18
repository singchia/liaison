package reporter

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jumboframes/armorigo/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/singchia/liaison/pkg/proto"
)

func (r *reporter) loopReportDevice(ctx context.Context) {
	for {
		device, err := getDevice()
		if err != nil {
			log.Errorf("get device error: %v", err)
			// 失败后等待 5 分钟再重试
			time.Sleep(5 * time.Minute)
			continue
		}
		device.EdgeID, _ = r.frontierBound.EdgeID()
		err = r.reportDevice(ctx, device)
		if err != nil {
			log.Errorf("report device error: %v", err)
			// 失败后等待 5 分钟再重试
			time.Sleep(5 * time.Minute)
			continue
		}
		// 成功后等待 1 小时
		time.Sleep(time.Hour)
	}
}

func (r *reporter) loopReportDeviceUsage(ctx context.Context) {
	for {
		deviceUsage, err := getDeviceUsage()
		if err != nil {
			log.Errorf("get device usage error: %v", err)
			// 失败后等待 1 分钟再重试
			time.Sleep(time.Minute)
			continue
		}
		err = r.reportDeviceUsage(ctx, deviceUsage)
		if err != nil {
			log.Errorf("report device usage error: %v", err)
			// 失败后等待 1 分钟再重试
			time.Sleep(time.Minute)
			continue
		}
		// 成功后等待 1 分钟
		time.Sleep(time.Minute)
	}
}

func (r *reporter) reportDevice(ctx context.Context, device *proto.Device) error {
	data, err := json.Marshal(device)
	if err != nil {
		return err
	}

	req := r.frontierBound.NewRequest(data)
	rsp, err := r.frontierBound.Call(ctx, "report_device", req)
	if err != nil {
		return err
	}
	if rsp.Error() != nil {
		return rsp.Error()
	}
	log.Infof("report device success: %s", string(data))
	return nil
}

func (r *reporter) reportDeviceUsage(ctx context.Context, deviceUsage *proto.DeviceUsage) error {
	data, err := json.Marshal(deviceUsage)
	if err != nil {
		return err
	}
	req := r.frontierBound.NewRequest(data)
	rsp, err := r.frontierBound.Call(ctx, "report_device_usage", req)
	if err != nil {
		return err
	}
	if rsp.Error() != nil {
		return rsp.Error()
	}
	log.Infof("report device usage success: %s", string(data))
	return nil
}

// 注意：指纹才是唯一标识一台主机，而不是边缘和AKSK
// 如果当前edge扫描的指纹发生改变，则会重新建立edge和device的关联
func getDevice() (*proto.Device, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	// 基础信息
	cpuCount, err := cpu.Counts(true)
	if err != nil {
		// 如果无法获取 CPU 数量，使用默认值 1
		cpuCount = 1
		log.Warnf("failed to get CPU count, using default: 1, error: %v", err)
	}

	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}
	memMB := memory.Total / 1024 / 1024

	// OS
	info, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %w", err)
	}
	os := info.OS
	osVersion := info.KernelVersion

	// 根磁盘大小
	diskUsage, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("failed to get disk usage: %w", err)
	}
	diskMB := diskUsage.Total / 1024 / 1024

	// 网络接口
	interfaces, err := getDeviceEthernetInterface()
	if err != nil {
		// 如果无法获取网络接口，使用空列表
		interfaces = []*proto.DeviceEthernetInterface{}
		log.Warnf("failed to get network interfaces, using empty list, error: %v", err)
	}
	// 指纹
	fingerprint, err := getFingerprint()
	if err != nil {
		return nil, fmt.Errorf("failed to get fingerprint: %w", err)
	}

	return &proto.Device{
		Fingerprint: fingerprint,
		HostName:    hostname,
		CPU:         cpuCount,
		Memory:      int(memMB),
		Disk:        int(diskMB),
		OS:          os,
		OSVersion:   osVersion,
		Interfaces:  interfaces,
	}, nil
}

func getDeviceEthernetInterface() ([]*proto.DeviceEthernetInterface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ethernetInterfaces := make([]*proto.DeviceEthernetInterface, 0)

	for _, iface := range interfaces {
		// 跳过 lo 网卡
		if iface.Name == "lo" || iface.Name == "lo0" {
			continue
		}
		// 跳过没有MAC或未启用的接口
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		ipMasks := make([]*proto.IPMask, 0)
		hasIPv4 := false

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet: // 带掩码的情况
				ip := v.IP.String()
				// 跳过 loopback 和 link-local 地址
				if v.IP.IsLoopback() || v.IP.IsLinkLocalUnicast() {
					continue
				}
				// 检查是否为 IPv4
				if v.IP.To4() != nil {
					hasIPv4 = true
				}
				mask := ""
				if v.Mask != nil {
					// 统一使用前缀长度（CIDR notation）
					ones, _ := v.Mask.Size()
					mask = fmt.Sprintf("%d", ones)
				}
				ipMasks = append(ipMasks, &proto.IPMask{
					IP:      ip,
					Netmask: mask,
				})

			case *net.IPAddr: // 只有IP没有掩码
				// 跳过 loopback 和 link-local 地址
				if v.IP.IsLoopback() || v.IP.IsLinkLocalUnicast() {
					continue
				}
				// 检查是否为 IPv4
				if v.IP.To4() != nil {
					hasIPv4 = true
				}
				ipMasks = append(ipMasks, &proto.IPMask{
					IP:      v.IP.String(),
					Netmask: "",
				})
			}
		}

		// 只添加有IP地址且至少有一个IPv4地址的网卡（跳过只有IPv6的网卡）
		if len(ipMasks) > 0 && hasIPv4 {
			ethernetInterfaces = append(ethernetInterfaces, &proto.DeviceEthernetInterface{
				Name:    iface.Name,
				MAC:     iface.HardwareAddr.String(),
				IPMasks: ipMasks,
			})
		}
	}
	return ethernetInterfaces, nil
}

func getDeviceUsage() (*proto.DeviceUsage, error) {
	// 获取指纹
	fingerprint, err := getFingerprint()
	if err != nil {
		return nil, fmt.Errorf("failed to get fingerprint: %w", err)
	}

	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	// 资源使用
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}
	diskUsage, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}
	deviceUsage := &proto.DeviceUsage{
		Fingerprint: fingerprint,
		CPUUsage:    float32(cpuUsage[0]),
		MemoryUsage: float32(memory.UsedPercent),
		DiskUsage:   float32(diskUsage.UsedPercent),
	}
	return deviceUsage, nil
}

// GetFingerprint 获取设备指纹（导出函数，供命令行使用）
func GetFingerprint() (string, error) {
	return getFingerprint()
}

// 基于 MAC、CPU、磁盘序列号生成指纹
// 所有组件都进行排序，确保指纹稳定性
func getFingerprint() (string, error) {
	// 获取所有非 loopback 网卡的 MAC 地址，并排序
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %w", err)
	}
	macs := make([]string, 0)
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) > 0 && iface.Flags&net.FlagLoopback == 0 {
			macs = append(macs, iface.HardwareAddr.String())
		}
	}
	if len(macs) == 0 {
		return "", fmt.Errorf("failed to get MAC address")
	}
	// 对 MAC 地址进行排序，确保顺序稳定
	sort.Strings(macs)
	mac := strings.Join(macs, ",") // 使用所有 MAC 地址，用逗号分隔

	// CPU 信息
	cpuID := ""
	cpuInfo, err := cpu.Info()
	if err != nil {
		// 如果无法获取 CPU 信息，使用 hostname 作为备用
		hostname, _ := os.Hostname()
		cpuID = hostname
	} else if len(cpuInfo) > 0 {
		// 如果有多个 CPU，对 CPU 信息进行排序
		cpuStrings := make([]string, 0, len(cpuInfo))
		for _, info := range cpuInfo {
			cpuStr := info.ModelName + info.VendorID + info.Family
			cpuStrings = append(cpuStrings, cpuStr)
		}
		// 去重并排序
		cpuMap := make(map[string]bool)
		uniqueCpus := make([]string, 0)
		for _, cpuStr := range cpuStrings {
			if !cpuMap[cpuStr] {
				cpuMap[cpuStr] = true
				uniqueCpus = append(uniqueCpus, cpuStr)
			}
		}
		sort.Strings(uniqueCpus)
		cpuID = strings.Join(uniqueCpus, ",")
	}

	// 磁盘序列号（收集所有磁盘的序列号，并排序）
	diskIDs := make([]string, 0)
	counts, err := disk.IOCounters()
	if err == nil && len(counts) > 0 {
		for _, stats := range counts {
			if stats.SerialNumber != "" {
				diskIDs = append(diskIDs, stats.SerialNumber)
			}
		}
	}
	// 如果无法获取磁盘序列号，使用 hostname 作为备用
	if len(diskIDs) == 0 {
		hostname, _ := os.Hostname()
		diskIDs = []string{hostname}
	}
	// 对磁盘序列号进行排序，确保顺序稳定
	sort.Strings(diskIDs)
	diskID := strings.Join(diskIDs, ",") // 使用所有磁盘序列号，用逗号分隔

	raw := fmt.Sprintf("%s|%s|%s", mac, cpuID, diskID)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:]), nil
}
