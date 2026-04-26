package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/viruslox/VLX_FrameFlow_GUI/backend/internal/api"
	"github.com/viruslox/VLX_FrameFlow_GUI/backend/internal/system"
)

type TelemetryData struct {
	Type              string                                  `json:"type"`
	NetworkInterfaces map[string]system.NetworkInterfaceStats `json:"network_interfaces"`
	SystemUsage       system.SystemUsage                      `json:"system_usage"`
	WifiMode          string                                  `json:"wifi_mode"`
}

func StartTelemetryWorker(wsHub *api.WSHub) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			data := TelemetryData{
				Type:              "telemetry",
				NetworkInterfaces: system.GetNetworkInterfaces(),
				SystemUsage:       system.GetSystemUsage(),
				WifiMode:          system.GetWifiMode(),
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("Failed to marshal telemetry data: %v", err)
				continue
			}

			wsHub.Broadcast(jsonData)
		}
	}()
}
