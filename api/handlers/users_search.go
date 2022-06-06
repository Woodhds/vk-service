package handlers

import (
	"encoding/json"
	"github.com/woodhds/vk.service/internal/vkclient"
	"net/http"
)

func UsersSearchHandler(token string, version string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")

		if q == "" {
			return
		}

		client, _ := vkclient.NewUserClient(token, version)
		response, _ := client.Search(q)

		json.NewEncoder(rw).Encode(response)
	})
}
