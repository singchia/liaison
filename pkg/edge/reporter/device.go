package reporter

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/jumboframes/armorigo/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/singchia/liaison/pkg/proto"
	"github.com/singchia/liaison/pkg/utils"
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
	fingerprint, err := utils.GetFingerprint()
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
		// 在 macOS、Windows 和 Linux 上，只使用物理接口
		if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
			if len(iface.HardwareAddr) == 0 {
				continue
			}
			mac := iface.HardwareAddr.String()
			// 只使用物理接口
			if !utils.IsPhysicalInterface(iface.Name, mac) {
				continue
			}
		} else if runtime.GOOS == "linux" {
			// Linux: 使用 sysfs 判断物理接口
			if !utils.IsPhysicalInterfaceLinux(iface) {
				continue
			}
		} else {
			// 其他系统，过滤掉常见的虚拟接口
			name := strings.ToLower(iface.Name)
			if strings.Contains(name, "utun") || // VPN 接口
				strings.Contains(name, "bridge") || // 桥接接口
				strings.Contains(name, "vmnet") || // VMware 虚拟接口
				strings.Contains(name, "vboxnet") || // VirtualBox 虚拟接口
				strings.Contains(name, "awdl") || // Apple Wireless Direct Link (可能不稳定)
				strings.Contains(name, "anpi") { // Apple Network Packet Injection
				continue
			}
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
	fingerprint, err := utils.GetFingerprint()
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
	return utils.GetFingerprint()
}
