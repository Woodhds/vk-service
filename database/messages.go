package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/woodhds/vk.service/message"
)

type MessagesQueryService interface {
	GetMessages(search string, ctx context.Context) ([]*message.VkMessageModel, error)
	GetMessageById(ownerId int, id int) *message.SimpleMessageModel
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

	var data []*message.VkMessageModel

	for res.Next() {
		m := &message.VkMessageModel{}
		e := res.Scan(&m.ID, &m.FromID, &m.Date, &m.Images, &m.LikesCount, &m.Owner, &m.OwnerID, &m.RepostsCount, &m.Text, &m.UserReposted)
		if e == nil {
			data = append(data, m)
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
