package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/woodhds/vk.service/database"
	"github.com/woodhds/vk.service/internal/predictor"
	"net/http"
	"strconv"
)

func MessageSaveHandler(predict predictor.Predictor, factory database.ConnectionFactory) http.Handler {
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

		var data predictor.SavePredictRequest
		json.NewDecoder(r.Body).Decode(&data)

		if owner == 0 || messageId == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		var message struct {
			text  string
			owner string
		}

		conn, _ := factory.GetConnection(r.Context())
		defer conn.Close()

		if d, e := conn.QueryContext(r.Context(), "SELECT Text, Owner from messages where OwnerId = $1 and Id = $2", owner, messageId); e != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		} else {
			for d.Next() {
				d.Scan(&message.text, &message.owner)
			}

			d.Close()
		}

		if e := predict.SaveMessage(owner, messageId, message.text, message.owner, data.Category); e != nil {
			rw.WriteHeader(http.StatusBadRequest)
		}
		Json(rw, true)
	})
}
