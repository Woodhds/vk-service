package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/woodhds/vk.service/protos"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

func MessageSaveHandler(conn grpc.ClientConnInterface) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		vars:=mux.Vars(r)
		var owner int
		var messageId int
		if ownerId, e := strconv.Atoi(vars["ownerId"]); e == nil {
			owner = ownerId
		}

		if id, e := strconv.Atoi(vars["id"]); e == nil {
			messageId = id
		}

		category := r.URL.Query().Get("category")

		if owner == 0 || messageId == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}


		client := protos.NewMessageSaveServiceClient(conn)
		if _, e := client.SaveMessage(context.Background(), &protos.MessageSaveRequest{
			OwnerId: int32(owner),
			Id:       int32(messageId),
			Category: category,
		}); e != nil {
			rw.WriteHeader(http.StatusBadRequest)
		}
		Json(rw, true)
	})
}