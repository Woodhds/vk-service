package handlers

import (
	"github.com/gorilla/mux"
	"github.com/woodhds/vk.service/predictor"
	"io"
	"net/http"
	"strconv"
)

func PredictHandler(predictorClient predictor.Predictor) http.Handler {
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

		bodyBytes, _ := io.ReadAll(r.Body)
		text := string(bodyBytes)

		resp, e := predictorClient.PredictMessage(&predictor.PredictMessage{OwnerId: owner, Id: messageId, Text: text})

		if e != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		Json(rw, resp)
	})
}
