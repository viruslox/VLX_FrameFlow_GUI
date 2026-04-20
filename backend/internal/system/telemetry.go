package system

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type NetworkInterfaceStats struct {
	RxBytes uint64 `json:"rx_bytes"`
	TxBytes uint64 `json:"tx_bytes"`
}

type GPSData struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Speed float64 `json:"speed"`
}

// GetNetworkInterfaces parses /proc/net/dev to get interface stats.
// Falls back to mock data if it fails.
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
	return stats
}

// GetSystemLoad parses /proc/loadavg.
// Falls back to mock data if it fails.
func GetSystemLoad() []float64 {
	data, err := os.ReadFile("/proc/loadavg")
	if err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			load1, _ := strconv.ParseFloat(fields[0], 64)
			load5, _ := strconv.ParseFloat(fields[1], 64)
			load15, _ := strconv.ParseFloat(fields[2], 64)
			return []float64{load1, load5, load15}
		}
	}
	log.Printf("Failed to read /proc/loadavg, using mock data: %v", err)
	return []float64{0.5, 0.4, 0.3}
}

// GetGPSData returns mock GPS data for now.
func GetGPSData() GPSData {
	return GPSData{
		Lat:   37.7749,
		Lon:   -122.4194,
		Speed: 15.2,
	}
}

// GetFFmpegLogs returns mock latest ffmpeg logs.
func GetFFmpegLogs() []string {
	return []string{
		"ffmpeg log 1: frame=  100 fps= 30 q=28.0 size=    1024kB time=00:00:03.33 bitrate=2516.5kbits/s speed=   1x",
		"ffmpeg log 2: frame=  130 fps= 30 q=28.0 size=    1280kB time=00:00:04.33 bitrate=2419.8kbits/s speed=   1x",
	}
}
