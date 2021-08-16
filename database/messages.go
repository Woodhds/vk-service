package database

import (
	"database/sql"
	"fmt"
	"github.com/woodhds/vk.service/message"
)

type MessagesQueryService interface {
	GetMessages(search string) ([]*message.VkCategorizedMessageModel, error)
}

type messageQueryService struct {
	conn *sql.DB
}

func (m messageQueryService) GetMessages(search string) ([]*message.VkCategorizedMessageModel, error) {
	s := fmt.Sprintf(`"%s"`, search)
	res, e := m.conn.Query(`
			SELECT messages.Id,
			   FromId,
			   Date,
			   Images,
			   LikesCount,
			   Owner,
			   messages.OwnerId,
			   RepostsCount,
			   messages.Text,
			   UserReposted,
			   ts_rank(to_tsvector(s.text), plainto_tsquery($1)) as rank
		FROM messages inner join messages_search as s on messages.Id = s.Id AND messages.OwnerId = s.OwnerId
		where s.Text @@ plainto_tsquery($2)
		order by rank desc
				`,  s, s)

	if e != nil {
		return nil, e
	}

	var data []*message.VkCategorizedMessageModel

	for res.Next() {
		m := message.VkCategorizedMessageModel{
			VkMessageModel: &message.VkMessageModel{},
		}
		e := res.Scan(&m.ID, &m.FromID, &m.Date, &m.Images, &m.LikesCount, &m.Owner, &m.OwnerID, &m.RepostsCount, &m.Text, &m.UserReposted)
		if e == nil {
			data = append(data, &m)
		}
	}
	defer res.Close()

	return data, nil
}

func NewMessageQueryService(conn *sql.DB) MessagesQueryService {
	return &messageQueryService{
		conn: conn,
	}
}
