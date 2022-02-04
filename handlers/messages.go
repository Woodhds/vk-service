package handlers

import (
	"github.com/woodhds/vk.service/database"
	"github.com/woodhds/vk.service/message"
	"github.com/woodhds/vk.service/predictor"
	"net/http"
)

func MessagesHandler(messageQueryService database.MessagesQueryService, predictorClient predictor.Predictor) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")

		data, e := messageQueryService.GetMessages(search)

		if e != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

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

			if respPredictions, e := predictorClient.Predict(predictions); e == nil {
				MapCategoriesToMessages(data, respPredictions)
			}
		}

		Json(rw, data)
	})
}

func MapCategoriesToMessages(data []*message.VkCategorizedMessageModel, predictions []*predictor.PredictMessage) {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(predictions); j++ {
			if predictions[j].Id == data[i].ID && data[i].OwnerID == predictions[j].OwnerId {
				data[i].Category = predictions[j].Category
				data[i].IsAccept = predictions[j].IsAccept
				data[i].Scores = predictions[j].Scores
				break
			}
		}
	}
}
