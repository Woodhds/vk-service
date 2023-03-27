package database

import (
	"context"
	"github.com/woodhds/vk.service/message"
)

type GroupsQueryService interface {
	Add(id int, name string, avatar string, ctx context.Context) error
	Remove(id int, ctx context.Context) error
	Get(page int, count int, ctx context.Context) ([]*message.GroupModel, error)
}
type groupsImplementation struct {
	connectionFactory ConnectionFactory
}

func (g *groupsImplementation) Get(page int, count int, ctx context.Context) ([]*message.GroupModel, error) {
	conn, _ := g.connectionFactory.GetConnection(ctx)
	defer conn.Close()

	rows, e := conn.QueryContext(ctx, "SELECT id, name, avatar FROM favorite_groups limit $1 offset $2", count, (page-1)*count)

	if e != nil {
		return nil, e
	}

	res := make([]*message.GroupModel, 0, 8)

	for rows.Next() {
		mes := message.GroupModel{}
		if e := rows.Scan(&mes.Id, &mes.Name, &mes.Avatar); e == nil {
			res = append(res, &mes)
		}
	}

	return res, nil

}

func (g *groupsImplementation) Add(id int, name string, avatar string, ctx context.Context) error {
	conn, _ := g.connectionFactory.GetConnection(ctx)
	defer conn.Close()

	if _, e := conn.ExecContext(ctx, `
        INSERT INTO favorite_groups (id, name, avatar) VALUES ($1, $2, $3) 
        ON CONFLICT (id) DO UPDATE 
            SET name = excluded.name, 
                avatar = excluded.avatar`,
		id, name, avatar); e != nil {
		return e
	}

	return nil
}

func (g *groupsImplementation) Remove(id int, ctx context.Context) error {
	conn, _ := g.connectionFactory.GetConnection(ctx)
	defer conn.Close()

	if _, e := conn.ExecContext(ctx, "DELETE FROM favorite_groups where id = $1", id); e != nil {
		return e
	}

	return nil
}

func NewGroupsQueryService(factory ConnectionFactory) GroupsQueryService {
	return &groupsImplementation{connectionFactory: factory}
}
