package handlers

import (
	"encoding/json"
	"github.com/woodhds/vk.service/message"
	"github.com/woodhds/vk.service/vkclient"
	"net/http"
)

func RepostHandler(token string, version string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var d []message.VkRepostMessage
		json.NewDecoder(r.Body).Decode(&d)

		go func() {

			wallClient, _ := vkclient.NewWallClient(token, version)
			groupClient, _ := vkclient.NewGroupClient(token, version)

			data, _ := wallClient.GetById(&d, "is_member")

			for _, d := range data.Response.Groups {
				if d.IsMember == 0 {
					groupClient.Join(d.ID)
				}
			}

			for _, i := range data.Response.Items {
				wallClient.Repost(&message.VkRepostMessage{OwnerID: i.OwnerID, ID: i.ID})
			}

		}()

		rw.WriteHeader(http.StatusOK)
	})
}
