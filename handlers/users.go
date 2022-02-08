package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/woodhds/vk.service/database"
	"github.com/woodhds/vk.service/vkclient"
	"net/http"
)

func UsersHandler(factory database.ConnectionFactory) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		conn, _ := factory.GetConnection()

		if r.Method == http.MethodGet {
			if rows, e := conn.Query(`SELECT Id, coalesce(Name, '') as Name, coalesce(Avatar,'') as Avatar from VkUserModel`); e != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			} else {
				defer rows.Close()

				var res []vkclient.VkUserMdodel

				for rows.Next() {
					u := vkclient.VkUserMdodel{}
					if e := rows.Scan(&u.Id, &u.Name, &u.Avatar); e == nil {
						res = append(res, u)
					} else {
						fmt.Println(e)
					}

				}

				Json(rw, res)
			}
		}

		if r.Method == http.MethodPost {
			u := &vkclient.VkUserMdodel{}
			json.NewDecoder(r.Body).Decode(u)

			if u.Id == 0 {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			_, e := conn.Exec(`INSERT INTO VkUserModel (Id, Avatar, Name) VALUES ($1, $2, $3)`, u.Id, u.Avatar, u.Name)

			if e != nil {
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

		json.NewEncoder(rw).Encode(&response)
	})
}
