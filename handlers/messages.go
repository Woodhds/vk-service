package handlers

import (
	"database/sql"
	"fmt"
	"github.com/woodhds/vk.service/message"
	"github.com/woodhds/vk.service/predictor"
	"net/http"
)

type VkCategorizedMessageModel struct {
	*message.VkMessageModel
	Category string `json:"category"`
	IsAccept bool   `json:"isAccept"`
}

func MessagesHandler(conn *sql.DB, predictorClient predictor.Predictor) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")

		res, e := conn.Query(`
			SELECT messages.Id, FromId, Date, Images, LikesCount, Owner, messages.OwnerId, RepostedFrom, RepostsCount, highlight(messages_search, 2, '<b><i><big>', '</big></i></b>') as Text, UserReposted
			FROM messages inner join messages_search as search  on messages.Id = search.Id AND  messages.OwnerId = search.OwnerId 
				where search.Text MATCH @search
				order by rank desc
				`, sql.Named("search", fmt.Sprintf(`"%s"`, search)))

		if e != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		var data []*VkCategorizedMessageModel
		var predictions []*predictor.PredictMessage

		for res.Next() {
			m := VkCategorizedMessageModel{
				VkMessageModel: &message.VkMessageModel{},
			}
			e := res.Scan(&m.ID, &m.FromID, &m.Date, &m.Images, &m.LikesCount, &m.Owner, &m.OwnerID, &m.RepostedFrom, &m.RepostsCount, &m.Text, &m.UserReposted)
			if e == nil {
				data = append(data, &m)
				predictions = append(predictions, &predictor.PredictMessage{
					OwnerId:  m.OwnerID,
					Id:       m.ID,
					Category: "",
					Text:     m.Text,
				})
			}
		}
		res.Close()

		if len(data) > 0 {
			if respPredictions, e := predictorClient.Predict(predictions); e == nil {
				MapCategoriesToMessages(data, respPredictions)
			}
		}

		Json(rw, data)
	})
}

func MapCategoriesToMessages(data []*VkCategorizedMessageModel, predictions []*predictor.PredictMessage) {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(predictions); j++ {
			if predictions[j].Id == data[i].ID && data[i].OwnerID == predictions[j].OwnerId {
				data[i].Category = predictions[j].Category
				data[i].IsAccept = predictions[j].IsAccept
				break
			}
		}
	}
}
