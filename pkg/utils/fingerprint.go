package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
)

// isVirtualInterface 判断是否是虚拟网络接口
// 通过接口名称和 MAC 地址特征来判断
func isVirtualInterface(iface *net.Interface) bool {
	name := strings.ToLower(iface.Name)

	// 在 macOS 上，使用更严格的物理接口判断
	if runtime.GOOS == "darwin" {
		if len(iface.HardwareAddr) > 0 {
			mac := iface.HardwareAddr.String()
			// 只使用物理接口（enX 开头且非本地管理的 MAC）
			if !IsPhysicalInterface(iface.Name, mac) {
				return true
			}
		}
		// 即使没有 MAC 地址，如果不是 enX 开头，也认为是虚拟接口
		if !strings.HasPrefix(iface.Name, "en") {
			return true
		}
		return false
	}

	// 在 Windows 上，使用 WMI 查询物理接口
	if runtime.GOOS == "windows" {
		if len(iface.HardwareAddr) == 0 {
			return true // 没有 MAC 地址，认为是虚拟接口
		}
		mac := iface.HardwareAddr.String()
		// 使用 WMI 查询结果判断
		if !IsPhysicalInterface(iface.Name, mac) {
			return true
		}
		return false
	}

	// 在 Linux 上，使用 sysfs 判断物理接口
	if runtime.GOOS == "linux" {
		// 使用 IsPhysicalInterfaceLinux 判断
		if !IsPhysicalInterfaceLinux(*iface) {
			return true // 不是物理接口，返回 true 表示是虚拟接口
		}
		return false // 是物理接口，返回 false 表示不是虚拟接口
	}

	// 其他系统，通过接口名称判断
	if strings.Contains(name, "utun") || // VPN 接口
		strings.Contains(name, "bridge") || // 桥接接口
		strings.Contains(name, "vmnet") || // VMware 虚拟接口
		strings.Contains(name, "vboxnet") || // VirtualBox 虚拟接口
		strings.Contains(name, "awdl") || // Apple Wireless Direct Link (可能不稳定)
		strings.Contains(name, "anpi") { // Apple Network Packet Injection
		return true
	}

	return false
}

// GetFingerprint 获取设备指纹
// 基于 MAC、CPU、磁盘序列号生成指纹
// 所有组件都进行排序，确保指纹稳定性
func GetFingerprint() (string, error) {
	// 获取所有非 loopback 网卡的 MAC 地址，并排序
	// 在 Mac 上，需要过滤掉不稳定的虚拟接口（如 VPN、虚拟网卡等）
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %w", err)
	}
	macs := make([]string, 0)
	for _, iface := range interfaces {
		// 跳过 loopback 接口
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 跳过没有硬件地址的接口
		if len(iface.HardwareAddr) == 0 {
			continue
		}
		// 在 Mac 上，过滤掉常见的虚拟接口
		// 这些接口可能会动态变化，导致指纹不稳定
		if isVirtualInterface(&iface) {
			continue
		}
		// 只使用物理接口（有 UP 标志且不是点对点接口）
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagPointToPoint == 0 {
			macs = append(macs, iface.HardwareAddr.String())
		}
	}
	if len(macs) == 0 {
		// 如果过滤后没有 MAC 地址，回退到使用所有非 loopback 的 MAC
		for _, iface := range interfaces {
			if len(iface.HardwareAddr) > 0 && iface.Flags&net.FlagLoopback == 0 {
				// 在 macOS、Windows 和 Linux 上，回退时也使用物理接口判断
				if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
					mac := iface.HardwareAddr.String()
					if !IsPhysicalInterface(iface.Name, mac) {
						continue
					}
				} else if runtime.GOOS == "linux" {
					if !IsPhysicalInterfaceLinux(iface) {
						continue
					}
				}
				macs = append(macs, iface.HardwareAddr.String())
			}
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
	// 在 Mac 上，只使用物理磁盘的序列号，过滤掉虚拟磁盘和挂载点
	diskIDs := make([]string, 0)
	counts, err := disk.IOCounters()
	if err == nil && len(counts) > 0 {
		for _, stats := range counts {
			// 只使用有序列号的磁盘
			if stats.SerialNumber != "" {
				// 在 Mac 上，过滤掉可能的虚拟磁盘（名称包含特定关键词）
				name := strings.ToLower(stats.Name)
				// 跳过明显的虚拟磁盘（可以根据实际情况调整）
				if !strings.Contains(name, "disk image") &&
					!strings.Contains(name, "virtual") {
					diskIDs = append(diskIDs, stats.SerialNumber)
				}
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
	fingerprint := hex.EncodeToString(sum[:])

	return fingerprint, nil
}
