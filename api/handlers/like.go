package handlers

import (
	"github.com/woodhds/vk.service/internal/notifier"
	"net/http"
)

func LikeHandler(service *notifier.NotifyService) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		go func() {
			service.Success("like complete")
		}()
		rw.WriteHeader(http.StatusOK)
	})
}
