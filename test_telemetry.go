package main

import (
	"encoding/json"
	"fmt"
	"github.com/viruslox/VLX_FrameFlow_GUI/backend/internal/system"
)

func main() {
	ni := system.GetNetworkInterfaces()
	b, _ := json.MarshalIndent(ni, "", "  ")
	fmt.Println(string(b))
}
