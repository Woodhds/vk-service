package vkclient

import "testing"

func TestJoin(t *testing.T) {
	groupClient, e := NewGroupClient("", "5.131")
	if e != nil {
		t.Error(e)
	}

	if e := groupClient.Join(203189081); e != nil {
		t.Error(e)
	}
}
