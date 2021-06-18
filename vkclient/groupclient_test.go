package vkclient

import "testing"

func TestJoin(t *testing.T) {
	groupClient, e := NewGroupClient("6b066d614f742ff2850d568b8676e4e0240c2768088a5c3c58b2306047544a650219d8b8ecb079f5c1e72", "5.130")
	if e != nil {
		t.Error(e)
	}

	if e := groupClient.Join(203189081); e != nil {
		t.Error(e)
	}
}
