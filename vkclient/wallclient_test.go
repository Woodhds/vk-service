package vkclient

import (
	"testing"

	"github.com/woodhds/vk.service/message"
)

func TestRepost(t *testing.T) {
	wallClient, e := NewWallClient("6b066d614f742ff2850d568b8676e4e0240c2768088a5c3c58b2306047544a650219d8b8ecb079f5c1e72", "5.130")

	if e != nil {
		t.Error(e)
	}

	if e := wallClient.Repost(&message.VkRepostMessage{OwnerID: -174563218, ID: 415478}); e != nil {
		t.Error(e)
	}
}
