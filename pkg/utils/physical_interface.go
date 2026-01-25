package utils

import (
	"net"
	"runtime"
	"strings"
	"sync"
)

var (
	// Windows 物理适配器缓存（MAC 地址 -> 适配器信息）
	windowsPhysicalAdapters map[string]*WindowsNetworkAdapter
	windowsAdapterOnce      sync.Once
	windowsAdapterErr       error
)

// WindowsNetworkAdapter Windows 网络适配器信息（通过 WMI 获取）
type WindowsNetworkAdapter struct {
	Name            string
	Description     string
	MACAddress      *string
	PhysicalAdapter bool
	NetEnabled      *bool
	Manufacturer    *string
}

// getWindowsPhysicalAdapters 通过 WMI 获取物理网络适配器（Windows 专用）
func getWindowsPhysicalAdapters() (map[string]*WindowsNetworkAdapter, error) {
	if runtime.GOOS != "windows" {
		return make(map[string]*WindowsNetworkAdapter), nil
	}
	return getWindowsPhysicalAdaptersWMI()
}

// IsWindowsPhysicalInterface Windows 专用的物理接口判断
func IsWindowsPhysicalInterface(iface *net.Interface) bool {
	if runtime.GOOS != "windows" {
		return false
	}

	// 1. 排除 loopback
	if iface.Flags&net.FlagLoopback != 0 {
		return false
	}

	// 2. 必须有 MAC
	if len(iface.HardwareAddr) == 0 {
		return false
	}

	mac := iface.HardwareAddr.String()
	name := iface.Name

	// 3. 确保 WMI 查询已执行
	windowsAdapterOnce.Do(func() {
		windowsPhysicalAdapters, windowsAdapterErr = getWindowsPhysicalAdapters()
		if windowsAdapterErr != nil {
			// 如果 WMI 查询失败，回退到名称过滤
			windowsPhysicalAdapters = make(map[string]*WindowsNetworkAdapter)
		}
	})

	// 如果 WMI 查询失败，回退到名称过滤
	if windowsAdapterErr != nil {
		nameLower := strings.ToLower(name)
		virtualKeywords := []string{
			"virtual", "vmware", "hyper-v",
			"vbox", "tap", "tun", "vpn", "wsl",
			"loopback", "isatap", "teredo",
		}
		for _, k := range virtualKeywords {
			if strings.Contains(nameLower, k) {
				return false
			}
		}
		return true
	}

	// 标准化 MAC 地址
	macNormalized := NormalizeMAC(mac)

	// 检查是否在物理适配器映射中
	_, exists := windowsPhysicalAdapters[macNormalized]
	return exists
}

// IsLinuxPhysicalInterface Linux 专用的物理接口判断
// 只做基本过滤，不使用 sysfs 检查（适用于容器环境）
func IsLinuxPhysicalInterface(iface *net.Interface) bool {
	if runtime.GOOS != "linux" {
		return false
	}

	// 1. 排除 loopback
	if iface.Flags&net.FlagLoopback != 0 {
		return false
	}

	// 2. 接口名快速排除
	name := iface.Name
	virtualPrefixes := []string{
		"lo", "docker", "veth", "cni",
		"flannel", "virbr", "br-",
		"tun", "tap",
	}
	for _, p := range virtualPrefixes {
		if strings.HasPrefix(name, p) {
			return false
		}
	}

	// 3. 必须有 MAC
	if len(iface.HardwareAddr) == 0 {
		return false
	}

	return true
}

// IsMacPhysicalInterface macOS 专用的物理接口判断
func IsMacPhysicalInterface(iface *net.Interface) bool {
	if runtime.GOOS != "darwin" {
		return false
	}

	// 1. 排除 loopback
	if iface.Flags&net.FlagLoopback != 0 {
		return false
	}

	// 2. 必须有 MAC
	if len(iface.HardwareAddr) == 0 {
		return false
	}

	// 3. 接口名称必须以 "en" 开头
	if !strings.HasPrefix(iface.Name, "en") {
		return false
	}

	// 4. MAC 地址不能是本地管理的
	mac := iface.HardwareAddr.String()
	if IsLocallyAdministeredMAC(mac) {
		return false
	}

	return true
}

// isVirtualInterface 判断是否是虚拟网络接口
// 返回 true 表示是虚拟接口（应该跳过），返回 false 表示是物理接口（应该使用）
// 根据不同的系统使用不同的判断机制
func isVirtualInterface(iface *net.Interface) bool {
	// 1. 排除 loopback 接口
	if iface.Flags&net.FlagLoopback != 0 {
		return true
	}

	// 2. 排除没有硬件地址的接口
	if len(iface.HardwareAddr) == 0 {
		return true
	}

	// 3. 只使用有 UP 标志且不是点对点接口
	if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagPointToPoint != 0 {
		return true
	}

	// 4. 根据平台判断是否是物理接口
	if runtime.GOOS == "darwin" {
		return !IsMacPhysicalInterface(iface)
	}

	if runtime.GOOS == "windows" {
		return !IsWindowsPhysicalInterface(iface)
	}

	if runtime.GOOS == "linux" {
		return !IsLinuxPhysicalInterface(iface)
	}

	return false
}
