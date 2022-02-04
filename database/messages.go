package database

import (
	"database/sql"
	"github.com/woodhds/vk.service/message"
)

type MessagesQueryService interface {
	GetMessages(search string) ([]*message.VkCategorizedMessageModel, error)
	GetMessageById(ownerId int, id int) *message.SimpleMessageModel
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
			       ts_headline(messages.text, phraseto_tsquery($1), 'HighlightAll = true') as Text,
			       UserReposted
			FROM messages inner join messages_search as s on messages.Id = s.Id AND messages.OwnerId = s.OwnerId 
				where s.Text @@ phraseto_tsquery($1)
				order by ts_rank(to_tsvector(s.text), phraseto_tsquery($1)) desc
				`, search)

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

func (m messageQueryService) GetMessageById(ownerId int, id int) *message.SimpleMessageModel {
	res := m.conn.QueryRow(`SELECT messages.Id,
			       OwnerId,
			       messages.text as Text
				FROM messages where OwnerId = $1 AND Id = $2`, ownerId, id)

	if res == nil {
		return nil
	}
	var data message.SimpleMessageModel
	if e:= res.Scan(&data.ID, &data.OwnerID, &data.Text); e != nil {
		return nil
	}
	return &data
}

func NewMessageQueryService(conn *sql.DB) MessagesQueryService {
	return &messageQueryService{
		conn: conn,
	}
}
