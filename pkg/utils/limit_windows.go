//go:build windows
// +build windows

package utils

// SetRLimit 在 Windows 平台上不支持的实现
// Windows 不使用 Unix 风格的 rlimit，所以这里返回 nil
func SetRLimit(fileLimit uint64) error {
	// Windows 不需要设置 rlimit
	return nil
}
