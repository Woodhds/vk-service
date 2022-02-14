package handlers

import (
	"encoding/json"
	"github.com/woodhds/vk.service/database"
	"github.com/woodhds/vk.service/message"
	"github.com/woodhds/vk.service/vkclient"
	"net/http"
)

func RepostHandler(factory database.ConnectionFactory, token *string, version *string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var d []*message.VkRepostMessage
		json.NewDecoder(r.Body).Decode(&d)

		wallClient, _ := vkclient.NewWallClient(token, version)
		groupClient, _ := vkclient.NewGroupClient(token, version)

		data, _ := wallClient.GetById(d, "is_member")

		for _, d := range data.Response.Groups {
			if d.IsMember == 0 {
				groupClient.Join(d.ID)
			}
		}

		conn, _ := factory.GetConnection(r.Context())
		for _, i := range data.Response.Items {
			if e := wallClient.Repost(&message.VkRepostMessage{OwnerID: i.OwnerID, ID: i.ID}); e == nil {
				if _, e := conn.ExecContext(r.Context(), "UPDATE messages SET UserReposted = true where Id = $1 and OwnerId = $2", i.ID, i.OwnerID); e != nil {
					rw.WriteHeader(http.StatusBadRequest)
					rw.Write([]byte(e.Error()))
				}
			} else {
				rw.WriteHeader(http.StatusBadRequest)
			}
		}

		rw.WriteHeader(http.StatusOK)
	})
}
