package system

import (
	"os"
	"strings"
)

// GetWifiMode returns the current mode of the first wireless interface.
func GetWifiMode() string {
	iface := getFirstWifiInterface()
	if iface == "" {
		return "Not found"
	}

	// Try iw first
	cmd := execCommand("iw", "dev", iface, "info")
	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "type ") {
				t := strings.TrimPrefix(line, "type ")
				switch strings.ToLower(t) {
				case "managed":
					return "Managed"
				case "ap":
					return "Master"
				case "monitor":
					return "Monitor"
				case "ibss":
					return "Ad-Hoc"
				case "mesh point":
					return "Mesh"
				case "wds":
					return "Repeater"
				default:
					return t
				}
			}
		}
	}

	// Fallback to iwconfig
	cmd = execCommand("iwconfig", iface)
	out, err = cmd.Output()
	if err == nil {
		s := string(out)
		idx := strings.Index(s, "Mode:")
		if idx != -1 {
			modeStr := s[idx+5:]
			spaceIdx := strings.IndexAny(modeStr, " \n")
			if spaceIdx != -1 {
				modeStr = modeStr[:spaceIdx]
			}
			return modeStr
		}
	}

	return "Not found"
}

func getFirstWifiInterface() string {
	entries, err := os.ReadDir("/sys/class/net/")
	if err != nil {
		return ""
	}
	for _, e := range entries {
		iface := e.Name()
		if iface == "lo" {
			continue
		}
		if _, err := os.Stat("/sys/class/net/" + iface + "/wireless"); err == nil {
			return iface
		}
		if _, err := os.Stat("/sys/class/net/" + iface + "/phy80211"); err == nil {
			return iface
		}
	}
	return ""
}
