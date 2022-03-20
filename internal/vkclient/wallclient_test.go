package vkclient

import (
	"testing"

	"github.com/woodhds/vk.service/message"
)

func TestRepost(t *testing.T) {
	token := ""
	version := ""
	wallClient, e := NewWallClient(&token, &version)

	if e != nil {
		t.Error(e)
	}

	if e := wallClient.Repost(&message.VkRepostMessage{OwnerID: -174563218, ID: 415478}); e != nil {
		t.Error(e)
	}
}
