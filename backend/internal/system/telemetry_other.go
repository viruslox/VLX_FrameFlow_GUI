//go:build !linux

package system

func GetNetworkInterfaces() map[string]NetworkInterfaceStats {
	return make(map[string]NetworkInterfaceStats)
}

func GetSystemUsage() SystemUsage {
	return SystemUsage{}
}
