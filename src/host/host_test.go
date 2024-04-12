package host

import "testing"

func TestDetermineHostTypeValidIPAddress(t *testing.T) {
	hostInfo := HostInfoCmd{}
	host := "127.0.0.1"

	hostType, err := hostInfo.DetermineHostType(host)

	if err != nil || hostType != IP_ADDR {
		t.Fatalf("Failed parsing %s as %s, got %s or err %v", host, IP_ADDR, hostType, err)
	}
}

func TestDetermineHostTypeValidHostname(t *testing.T) {
	hostInfo := HostInfoCmd{}
	host := "www.google.com"

	hostType, err := hostInfo.DetermineHostType(host)

	if err != nil || hostType != HOSTNAME {
		t.Fatalf("Failed parsing %s as %s, got %s or err %v", host, HOSTNAME, hostType, err)
	}
}
