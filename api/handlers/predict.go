package handlers

import (
	"github.com/gorilla/mux"
	"github.com/woodhds/vk.service/internal/predictor"
	"github.com/woodhds/vk.service/message"
	"net/http"
	"strconv"
)

type PredictResponse struct {
	Message *message.SimpleMessageModel `json:"message"`
	Predict *map[string]float32         `json:"predict"`
}

func PredictHandler(predictorClient predictor.Predictor, messageService VkMessagesService) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var owner int
		var messageId int
		if ownerId, e := strconv.Atoi(vars["ownerId"]); e == nil {
			owner = ownerId
		}

		if id, e := strconv.Atoi(vars["id"]); e == nil {
			messageId = id
		}

		var mes *message.SimpleMessageModel

		if res := messageService.GetById([]*message.VkRepostMessage{{owner, messageId}}); len(res) == 0 {
			return
		} else {
			mes = &message.SimpleMessageModel{OwnerID: res[0].OwnerID, ID: res[0].ID, Text: res[0].Text}
		}

		resp, e := predictorClient.PredictMessage(&predictor.PredictMessage{OwnerId: owner, Id: messageId, Text: mes.Text})

		if e != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		result := &PredictResponse{Message: mes, Predict: &resp}

		Json(rw, result)
	})
}
