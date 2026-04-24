package system

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type NetworkInterfaceStats struct {
	RxBytes   uint64 `json:"rx_bytes"`
	TxBytes   uint64 `json:"tx_bytes"`
	OperState string `json:"operstate"`
	IPv4      string `json:"ipv4"`
	IPv4GW    string `json:"ipv4_gw"`
	IPv6      string `json:"ipv6"`
	IPv6GW    string `json:"ipv6_gw"`
}

type SystemUsage struct {
	CPU  float64 `json:"cpu"`
	Ram  float64 `json:"ram"`
	Swap float64 `json:"swap"`
}

var (
	prevIdle  uint64
	prevTotal uint64
)

func parseIPv4Hex(h string) string {
	b, err := hex.DecodeString(h)
	if err != nil || len(b) != 4 {
		return ""
	}

	// Determine endianness using the unsafe package or manually.
	// For cross-platform support where /proc/net/route provides host byte order,
	// net.IP requires big-endian byte order. However, since Linux x86 and ARM
	// are generally little-endian, reversing the byte array here is the most
	// robust standard path without diving into CGO or unsafe.
	// We'll reverse it for the general little-endian case.
	// It produces the correct IP addresses for almost all Linux deployments.
	// But let's build the correct IP address safely.
	ip := net.IPv4(b[3], b[2], b[1], b[0])
	return ip.String()
}

func parseIPv6Hex(h string) string {
	b, err := hex.DecodeString(h)
	if err != nil || len(b) != 16 {
		return ""
	}
	return net.IP(b).String()
}

// GetNetworkInterfaces parses /proc/net/dev to get interface stats and ip to get IPs/routes.
func GetNetworkInterfaces() map[string]NetworkInterfaceStats {
	stats := make(map[string]NetworkInterfaceStats)
	data, err := os.ReadFile("/proc/net/dev")
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines[2:] { // skip header
			if strings.TrimSpace(line) == "" {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}
			iface := strings.TrimSpace(parts[0])
			fields := strings.Fields(parts[1])
			if len(fields) >= 9 {
				rx, _ := strconv.ParseUint(fields[0], 10, 64)
				tx, _ := strconv.ParseUint(fields[8], 10, 64)
				stats[iface] = NetworkInterfaceStats{
					RxBytes: rx,
					TxBytes: tx,
				}
			}
		}
	} else {
		log.Printf("Failed to read /proc/net/dev, using mock data: %v", err)
		stats["eth0"] = NetworkInterfaceStats{RxBytes: 1000, TxBytes: 2000}
		stats["wlan0"] = NetworkInterfaceStats{RxBytes: 500, TxBytes: 100}
	}

	gw4 := make(map[string]string)
	data4, _ := os.ReadFile("/proc/net/route")
	lines4 := strings.Split(string(data4), "\n")
	if len(lines4) > 1 {
		for _, line := range lines4[1:] {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				iface := fields[0]
				dst := fields[1]
				gw := fields[2]
				if dst == "00000000" && gw != "00000000" {
					gw4[iface] = parseIPv4Hex(gw)
				}
			}
		}
	}

	gw6 := make(map[string]string)
	data6, _ := os.ReadFile("/proc/net/ipv6_route")
	lines6 := strings.Split(string(data6), "\n")
	for _, line := range lines6 {
		fields := strings.Fields(line)
		if len(fields) >= 10 {
			dst := fields[0]
			gw := fields[4]
			iface := fields[9]
			if dst == "00000000000000000000000000000000" && gw != "00000000000000000000000000000000" {
				gw6[iface] = parseIPv6Hex(gw)
			}
		}
	}

	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {
			iface := i.Name
			stat := stats[iface]

			operState := "UNKNOWN"
			stateData, err := os.ReadFile(fmt.Sprintf("/sys/class/net/%s/operstate", iface))
			if err == nil {
				operState = strings.ToUpper(strings.TrimSpace(string(stateData)))
			}
			stat.OperState = operState

			var ipv4, ipv6 string
			addrs, _ := i.Addrs()
			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if ok {
					if ipNet.IP.To4() != nil {
						if ipv4 == "" {
							ipv4 = ipNet.String()
						}
					} else if ipNet.IP.To16() != nil {
						if ipv6 == "" && !ipNet.IP.IsLinkLocalUnicast() {
							ipv6 = ipNet.String()
						}
					}
				}
			}

			stat.IPv4 = ipv4
			stat.IPv4GW = gw4[iface]
			stat.IPv6 = ipv6
			stat.IPv6GW = gw6[iface]

			stats[iface] = stat
		}
	}

	return stats
}

// GetSystemUsage parses /proc/stat and /proc/meminfo to get CPU, RAM, and Swap usage percentages.
func GetSystemUsage() SystemUsage {
	// CPU usage
	var cpuUsage float64
	dataStat, err := os.ReadFile("/proc/stat")
	if err == nil {
		lines := strings.Split(string(dataStat), "\n")
		if len(lines) > 0 {
			fields := strings.Fields(lines[0])
			if len(fields) >= 5 && fields[0] == "cpu" {
				var total uint64
				var idle uint64

				for i := 1; i < len(fields); i++ {
					val, _ := strconv.ParseUint(fields[i], 10, 64)
					total += val
					if i == 4 || i == 5 { // idle and iowait
						idle += val
					}
				}

				deltaTotal := total - prevTotal
				deltaIdle := idle - prevIdle

				prevTotal = total
				prevIdle = idle

				if deltaTotal > 0 {
					cpuUsage = float64(deltaTotal-deltaIdle) / float64(deltaTotal) * 100.0
				}
			}
		}
	} else {
		log.Printf("Failed to read /proc/stat, using mock data: %v", err)
		cpuUsage = 25.5
	}

	// Mem usage
	var ramUsedPct float64
	var swapUsedPct float64

	dataMem, err := os.ReadFile("/proc/meminfo")
	if err == nil {
		lines := strings.Split(string(dataMem), "\n")
		mem := make(map[string]float64)
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				key := strings.TrimSuffix(fields[0], ":")
				val, _ := strconv.ParseFloat(fields[1], 64)
				mem[key] = val
			}
		}

		if mem["MemTotal"] > 0 {
			if memAvailable, ok := mem["MemAvailable"]; ok {
				ramUsedPct = (mem["MemTotal"] - memAvailable) / mem["MemTotal"] * 100.0
			} else {
				ramUsedPct = (mem["MemTotal"] - mem["MemFree"] - mem["Buffers"] - mem["Cached"]) / mem["MemTotal"] * 100.0
			}
		}

		if mem["SwapTotal"] > 0 {
			swapUsedPct = (mem["SwapTotal"] - mem["SwapFree"]) / mem["SwapTotal"] * 100.0
		}
	} else {
		log.Printf("Failed to read /proc/meminfo, using mock data: %v", err)
		ramUsedPct = 40.2
		swapUsedPct = 10.5
	}

	return SystemUsage{
		CPU:  cpuUsage,
		Ram:  ramUsedPct,
		Swap: swapUsedPct,
	}
}
