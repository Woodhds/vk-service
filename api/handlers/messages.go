package handlers

import (
	"context"
	"github.com/woodhds/vk.service/database"
	vkMessages "github.com/woodhds/vk.service/gen/messages"
	"github.com/woodhds/vk.service/internal/predictor"
	"github.com/woodhds/vk.service/message"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type messagesImplementation struct {
	vkMessages.UnimplementedMessagesServiceServer
	messagesQueryService database.MessagesQueryService
	predictorClient      predictor.Predictor
}

func (m messagesImplementation) GetMessages(ctx context.Context, r *vkMessages.GetMessagesRequest) (*vkMessages.GetMessagesResponse, error) {
	search := r.GetSearch()

	data, e := m.messagesQueryService.GetMessages(search.GetValue(), ctx)

	if e != nil {
		return nil, e
	}

	response := mapToResponse(data)

	if len(data) > 0 {

		predictions := make([]*predictor.PredictMessage, len(data))

		for i := 0; i < len(predictions); i++ {
			predictions[i] = &predictor.PredictMessage{
				OwnerId:  data[i].OwnerID,
				Id:       data[i].ID,
				Category: "",
				Text:     data[i].Text,
			}
		}

		if respPredictions, e := m.predictorClient.Predict(predictions); e == nil {
			mapCategoriesToMessages(response, respPredictions)
		}
	}

	return &vkMessages.GetMessagesResponse{
		Messages: response,
	}, nil
}

func mapCategoriesToMessages(data []*vkMessages.VkMessageExt, predictions []*predictor.PredictMessage) {

	for i := 0; i < len(data); i++ {

		for j := 0; j < len(predictions); j++ {
			if int32(predictions[j].Id) == data[i].Id && int32(predictions[j].OwnerId) == data[i].OwnerId {
				data[i].Category = predictions[j].Category
				data[i].IsAccept = predictions[j].IsAccept
				data[i].Scores = predictions[j].Scores
				break
			}
		}
	}
}

func mapToResponse(data []*message.VkCategorizedMessageModel) []*vkMessages.VkMessageExt {
	n := len(data)
	res := make([]*vkMessages.VkMessageExt, n, n)

	for i := 0; i < n; i++ {
		res[i] = &vkMessages.VkMessageExt{
			Id:           int32(data[i].ID),
			FromId:       int32(data[i].FromID),
			Date:         timestamppb.New(time.Time(*data[i].Date)),
			Images:       data[i].Images,
			LikesCount:   int32(data[i].LikesCount),
			Owner:        data[i].Owner,
			OwnerId:      int32(data[i].OwnerID),
			RepostsCount: int32(data[i].RepostsCount),
			Text:         data[i].Text,
			UserReposted: data[i].UserReposted,
			Category:     "",
			IsAccept:     false,
			Scores:       nil,
		}
	}

	return res
}

func NewMessageHandler(messagesQueryService database.MessagesQueryService, predictorClient predictor.Predictor) vkMessages.MessagesServiceServer {
	return messagesImplementation{messagesQueryService: messagesQueryService, predictorClient: predictorClient}
}
