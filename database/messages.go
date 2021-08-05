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
	res, e := m.conn.Query(`
			SELECT messages.Id, 
			       FromId, 
			       Date, 
			       Images, 
			       LikesCount, 
			       Owner, 
			       messages.OwnerId,
			       RepostsCount, 
			       highlight(messages_search, 2, '<b><i><big>', '</big></i></b>') as Text, 
			       UserReposted
			FROM messages inner join messages_search as search  on messages.Id = search.Id AND  messages.OwnerId = search.OwnerId 
				where search.Text MATCH @search
				order by rank desc
				`, sql.Named("search", fmt.Sprintf(`"%s"`, search)))

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
