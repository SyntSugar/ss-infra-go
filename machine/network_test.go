package machine

import (
	"net"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	ip, err := GetLocalIP()

	if err != nil {
		t.Fatalf("Failed to get local IP: %v", err)
	}

	if net.ParseIP(ip) == nil {
		t.Errorf("Invalid local IP address: %s", ip)
	}
}
