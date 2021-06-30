package handlers

import (
	"database/sql"
	"github.com/woodhds/vk.service/protos"
	"net/http"
)

type handler struct {
	protos.
}

func MessageSaveHandler(conn *sql.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		protos.
	})
}
