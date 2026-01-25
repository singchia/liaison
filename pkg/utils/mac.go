package utils

import (
	"net"
	"strings"
)

// NormalizeMAC 标准化 MAC 地址（去除分隔符，转为小写）
func NormalizeMAC(mac string) string {
	mac = strings.ToLower(strings.ReplaceAll(mac, ":", ""))
	mac = strings.ReplaceAll(mac, "-", "")
	return mac
}

// IsLocallyAdministeredMAC 判断 MAC 地址是否是本地管理的
// 本地管理的 MAC 地址的第 2 位为 1
func IsLocallyAdministeredMAC(mac string) bool {
	hw, err := net.ParseMAC(mac)
	if err != nil {
		return false
	}
	if len(hw) == 0 {
		return false
	}
	// 检查是否是本地管理的 MAC 地址（第 2 位为 1 表示本地管理）
	return hw[0]&2 != 0
}

// IsValidMAC 判断是否是有效的 MAC 地址
func IsValidMAC(mac string) bool {
	hw, err := net.ParseMAC(mac)
	if err != nil {
		return false
	}
	return len(hw) > 0
}
