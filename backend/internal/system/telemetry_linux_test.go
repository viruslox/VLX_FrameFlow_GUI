//go:build linux

package system

import (
	"testing"
)

func TestGetNetworkInterfacesLinux(t *testing.T) {
	GetNetworkInterfaces()
}
