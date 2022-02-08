package vkclient

import (
	"testing"

	"github.com/woodhds/vk.service/message"
)

func TestRepost(t *testing.T) {
	wallClient, e := NewWallClient("", "5.130")

	if e != nil {
		t.Error(e)
	}

	if e := wallClient.Repost(&message.VkRepostMessage{OwnerID: -174563218, ID: 415478}); e != nil {
		t.Error(e)
	}
}
