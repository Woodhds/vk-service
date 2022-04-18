package database

import (
	"context"
	"github.com/woodhds/vk.service/message"
)

type MessagesQueryService interface {
	GetMessages(search string, ctx context.Context) ([]*message.VkCategorizedMessageModel, error)
	GetMessageById(ownerId int, id int) *message.SimpleMessageModel
}

type messageQueryService struct {
	factory ConnectionFactory
}

func (m messageQueryService) GetMessages(search string, ctx context.Context) ([]*message.VkCategorizedMessageModel, error) {
	conn, _ := m.factory.GetConnection(ctx)
	defer conn.Close()

	res, e := conn.QueryContext(ctx, `SELECT messages.Id,
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
				order by ts_rank(to_tsvector(s.text), phraseto_tsquery($1)) desc`)

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
	conn, _ := m.factory.GetConnection(context.Background())
	defer conn.Close()
	res := conn.QueryRowContext(context.Background(), `SELECT messages.Id,
			       OwnerId,
			       messages.text as Text
				FROM messages where OwnerId = $1 AND Id = $2`, ownerId, id)

	if res == nil {
		return nil
	}
	var data message.SimpleMessageModel
	if e := res.Scan(&data.ID, &data.OwnerID, &data.Text); e != nil {
		return nil
	}
	return &data
}

func NewMessageQueryService(conn ConnectionFactory) MessagesQueryService {
	return &messageQueryService{
		factory: conn,
	}
}
