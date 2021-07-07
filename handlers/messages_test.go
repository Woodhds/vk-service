package handlers

import (
	"github.com/woodhds/vk.service/message"
	"github.com/woodhds/vk.service/predictor"
	"testing"
)

func BenchmarkMapCategoriesToMessages(t *testing.B) {
	t.ReportAllocs()
	data := make([]*VkCategorizedMessageModel, 110000)
	for i := 0; i < 110000; i++ {
		data[i] = &VkCategorizedMessageModel{
			VkMessageModel: &message.VkMessageModel{
				ID:      1,
				OwnerID: 1,
			},
			Category: "test",
		}
	}


	predictions := []*predictor.PredictMessage{{OwnerId: 1, Id: 1, Category: "Test", Text: "tt"}}
	MapCategoriesToMessages(data, predictions)

	for i := 0; i < len(data); i++ {
		if data[i].Category == "" {
			t.Error("Category empty")
		}
	}

}
