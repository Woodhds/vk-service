package vkclient

import "testing"

func TestSearch(t *testing.T) {
	client, e := NewUserClient("", "5.130")
	if client == nil {
		t.Error("client is nil", e)
	}

	resp, err := client.Search("Сергей")

	if len(resp) == 0 {
		t.Error("Empty response", err)
	}
}
