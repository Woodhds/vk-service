package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/woodhds/vk.service/predictor"
	"net/http"
	"strconv"
)

func MessageSaveHandler(predict predictor.Predictor, conn *sql.DB) http.Handler {
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

		var text string

		if d, e := conn.Query("SELECT Text from messages where OwnerId = $1 and Id = $2", owner, messageId); e != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		} else {
			for d.Next() {
				d.Scan(&text)
			}

			d.Close()
		}

		if e := predict.SaveMessage(owner, messageId, text, data.Category); e != nil {
			rw.WriteHeader(http.StatusBadRequest)
		}
		Json(rw, true)
	})
}
