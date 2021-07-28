package handlers

import (
	"github.com/woodhds/vk.service/notifier"
	"net/http"
)

func LikeHandler(service notifier.Notifier) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		service.Success("like complete")
	})
}
