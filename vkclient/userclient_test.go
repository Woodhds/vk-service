package vkclient

import "testing"

func TestSearch(t *testing.T) {
	client, e := NewUserClient("6b066d614f742ff2850d568b8676e4e0240c2768088a5c3c58b2306047544a650219d8b8ecb079f5c1e72", "5.130")
	if client == nil {
		t.Error("client is nil", e)
	}

	resp, err := client.Search("Сергей")

	if len(resp) == 0 {
		t.Error("Empty response", err)
	}
}
