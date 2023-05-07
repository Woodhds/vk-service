package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/woodhds/vk.service/message"
)

type MessagesQueryService interface {
	GetMessages(search string, ctx context.Context) ([]*message.VkMessageModel, error)
}

type messageQueryService struct {
	factory ConnectionFactory
}

func (m messageQueryService) GetMessages(search string, ctx context.Context) ([]*message.VkMessageModel, error) {
	conn, _ := m.factory.GetConnection(ctx)
	defer conn.Close()

	res, e := conn.QueryContext(ctx, `
			SELECT messages.Id, 
			       FromId,
			       messages.OwnerId
			FROM messages inner join messages_search as search  on messages.Id = search.Id AND  messages.OwnerId = search.OwnerId 
				where search.Text MATCH @search
				order by rank desc
				`, sql.Named("search", fmt.Sprintf(`"%s"`, search)))

	if e != nil {
		return nil, e
	}

	var data []*message.VkMessageModel

	for res.Next() {
		m := &message.VkMessageModel{}
		e := res.Scan(&m.ID, &m.FromID, &m.OwnerID)
		if e == nil {
			data = append(data, m)
		}
	}
	defer res.Close()

	return data, nil
}

func NewMessageQueryService(conn ConnectionFactory) MessagesQueryService {
	return &messageQueryService{
		factory: conn,
	}
}
