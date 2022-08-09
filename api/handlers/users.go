package handlers

import (
	"encoding/json"
	"github.com/woodhds/vk.service/database"
	"github.com/woodhds/vk.service/internal/notifier"
	"github.com/woodhds/vk.service/internal/vkclient"
	"net/http"
	"strconv"
)

func UsersHandler(usersService database.UsersQueryService, notifier *notifier.NotifyService) http.Handler {
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
			u := &vkclient.VkUserModel{}
			json.NewDecoder(r.Body).Decode(u)

			if u.Id == 0 {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			if e := usersService.InsertNew(u.Id, u.Name, u.Avatar, r.Context()); e != nil {
				rw.WriteHeader(http.StatusBadRequest)
			}
		}

		if r.Method == http.MethodDelete {
			deleteUser(rw, r, usersService, notifier)
		}
	})
}

func deleteUser(rw http.ResponseWriter, r *http.Request, usersService database.UsersQueryService, notifier *notifier.NotifyService) {
	strId := r.URL.Query().Get("id")

	id, e := strconv.Atoi(strId)
	if e != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	e = usersService.Delete(id, r.Context())

	if e != nil {
		notifier.Danger(e.Error())
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	notifier.Success("User successful deleted")
	rw.WriteHeader(http.StatusOK)
}
