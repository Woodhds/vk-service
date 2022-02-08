package handlers

import (
	"github.com/gorilla/mux"
	"github.com/woodhds/vk.service/database"
	"github.com/woodhds/vk.service/message"
	"github.com/woodhds/vk.service/predictor"
	"github.com/woodhds/vk.service/vkclient"
	"net/http"
	"strconv"
)

type PredictResponse struct {
	Message *message.SimpleMessageModel `json:"message"`
	Predict *map[string]float32         `json:"predict"`
}

func PredictHandler(predictorClient predictor.Predictor, db database.MessagesQueryService, client vkclient.WallClient) http.Handler {
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
		mes := db.GetMessageById(owner, messageId)

		if mes == nil {
			if res, e := client.GetById([]*message.VkRepostMessage{{owner, messageId}}); e != nil || len(res.Response.Items) == 0 {
				return
			} else {
				mes = &message.SimpleMessageModel{OwnerID: res.Response.Items[0].OwnerID, ID: res.Response.Items[0].ID, Text: res.Response.Items[0].Text}
			}
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
