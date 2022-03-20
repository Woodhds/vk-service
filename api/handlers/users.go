package handlers

import (
	"encoding/json"
	"github.com/woodhds/vk.service/database"
	"github.com/woodhds/vk.service/internal/vkclient"
	"net/http"
)

func UsersHandler(usersService database.UsersQueryService) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if rows, e := usersService.GetFullUsers(r.Context()); e != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			} else {

				Json(rw, rows)
			}
		}

		if r.Method == http.MethodPost {
			u := &vkclient.VkUserMdodel{}
			json.NewDecoder(r.Body).Decode(u)

			if u.Id == 0 {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			if e := usersService.InsertNew(u.Id, u.Name, u.Avatar, r.Context()); e != nil {
				rw.WriteHeader(http.StatusBadRequest)
			}
		}
	})
}

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
