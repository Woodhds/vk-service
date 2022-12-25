package database

import (
	"context"
)

type GroupsQueryService interface {
	Add(id int, name string, avatar string, ctx context.Context) error
	Remove(id int, ctx context.Context) error
	Get(page int, count int, ctx context.Context)
}
type groupsImplementation struct {
	connectionFactory ConnectionFactory
}

func (g *groupsImplementation) Get(page int, count int, ctx context.Context) {
	conn, _ := g.connectionFactory.GetConnection(ctx)
	defer conn.Close()

	conn.QueryContext(ctx, "SELECT id, name, avatar FROM favorite_groups offset $1 limit $2", (page-1)*count, count)

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
