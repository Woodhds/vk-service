package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/woodhds/vk.service/protos"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

type saveRequest struct {
	Category string `json:"category"`
	Text     string `json:"text"`
}

func MessageSaveHandler(host string) http.Handler {
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

		var data saveRequest
		json.NewDecoder(r.Body).Decode(&data)

		if owner == 0 || messageId == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if conn, e := grpc.Dial(host, grpc.WithInsecure()); e != nil {
			fmt.Println(e)
			rw.WriteHeader(http.StatusBadRequest)
		} else {
			client := protos.NewMessageSaveServiceClient(conn)
			if _, e := client.SaveMessage(context.Background(), &protos.MessageSaveRequest{
				OwnerId:  int32(owner),
				Id:       int32(messageId),
				Category: data.Category,
				Text:     data.Text,
			}); e != nil {
				rw.WriteHeader(http.StatusBadRequest)
			}
			Json(rw, true)
		}

	})
}
