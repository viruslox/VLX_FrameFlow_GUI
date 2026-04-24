package system

import (
	"testing"
)

func BenchmarkGetNetworkInterfaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetNetworkInterfaces()
	}
}
