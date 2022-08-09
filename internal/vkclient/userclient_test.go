package vkclient

import "testing"

func TestSearch(t *testing.T) {
	token := ""
	version := ""
	client, e := NewUserClient(token, version)
	if client == nil {
		t.Error("client is nil", e)
	}

	resp, err := client.Search("Сергей")

	if len(resp) == 0 {
		t.Error("Empty response", err)
	}
}
