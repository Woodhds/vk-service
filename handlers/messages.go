package handlers

import (
	"database/sql"
	"fmt"
	"github.com/woodhds/vk.service/message"
	"net/http"
)

func MessagesHandler(conn *sql.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")

		res, e := conn.Query(`
			SELECT messages.Id, FromId, Date, Images, LikesCount, Owner, messages.OwnerId, RepostedFrom, RepostsCount, messages.Text, UserReposted
			FROM messages inner join messages_search as search  on messages.Id = search.Id AND  messages.OwnerId = search.OwnerId 
				where search.Text MATCH @search
				order by rank desc
				`, sql.Named("search", fmt.Sprintf(`"%s"`, search)))

		if e != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		var data []message.VkMessageModel

		for res.Next() {
			m := message.VkMessageModel{}
			e := res.Scan(&m.ID, &m.FromID, &m.Date, &m.Images, &m.LikesCount, &m.Owner, &m.OwnerID, &m.RepostedFrom, &m.RepostsCount, &m.Text, &m.UserReposted)
			if e == nil {
				data = append(data, m)
			}
		}
		defer res.Close()

		Json(rw, data)
	})
}
