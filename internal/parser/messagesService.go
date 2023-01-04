package parser

import (
	"fmt"
	"github.com/woodhds/vk.service/internal/vkclient"
	"github.com/woodhds/vk.service/message"
)

type VkMessagesService interface {
	GetMessages(id int, page int, count int) []*message.VkRepostMessage
	GetById(message []*message.VkRepostMessage) []*message.VkMessageModel
}

type vkMessageService struct {
	wallClient vkclient.WallClient
}

func (m *vkMessageService) GetMessages(id int, page int, count int) []*message.VkRepostMessage {
	data, e := m.wallClient.Get(&vkclient.WallGetRequest{OwnerId: id, Offset: (page - 1) * count, Count: count})

	if e != nil {
		fmt.Println(e)
		return nil
	}

	var messages []*message.VkRepostMessage

	for _, v := range data.Items {
		if len(v.CopyHistory) > 0 {
			messages = append(messages, &message.VkRepostMessage{OwnerID: v.CopyHistory[0].OwnerID, ID: v.CopyHistory[0].ID})
		}
	}

	return messages
}

func (m *vkMessageService) GetById(reposts []*message.VkRepostMessage) []*message.VkMessageModel {
	posts, e := m.wallClient.GetById(reposts)
	if e != nil {
		return nil
	}

	result := make([]*message.VkMessageModel, len(posts.Items))

	for i, post := range posts.Items {
		result[i] = message.New(post, posts.Groups)
	}

	return result
}

func NewMessageService(wallClient vkclient.WallClient) VkMessagesService {
	return &vkMessageService{wallClient}
}
