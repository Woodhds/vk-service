package database

import (
	"context"
	"fmt"
	"github.com/woodhds/vk.service/internal/vkclient"
)

type UsersQueryService interface {
	GetAll() ([]int, error)
	GetFullUsers(ctx context.Context) ([]*vkclient.VkUserModel, error)
	InsertNew(id int, name string, avatar string, ctx context.Context) error
	Delete(id int, ctx context.Context) error
}

type userQueryService struct {
	factory ConnectionFactory
}

func (m *userQueryService) GetAll() ([]int, error) {
	conn, _ := m.factory.GetConnection(context.Background())
	defer conn.Close()

	rows, _ := conn.QueryContext(context.Background(), `select Id from VkUserModel`)
	var err error

	var ids []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err == nil {
			ids = append(ids, id)
		}
	}

	rows.Close()

	return ids, nil
}

func (m *userQueryService) GetFullUsers(ctx context.Context) ([]*vkclient.VkUserModel, error) {
	conn, _ := m.factory.GetConnection(ctx)
	defer conn.Close()
	if rows, e := conn.QueryContext(ctx, `SELECT Id, coalesce(Name, '') as Name, coalesce(Avatar,'') as Avatar from VkUserModel`); e != nil {
		return nil, e
	} else {
		defer rows.Close()

		var res []*vkclient.VkUserModel

		for rows.Next() {
			u := vkclient.VkUserModel{}
			if e := rows.Scan(&u.Id, &u.Name, &u.Avatar); e == nil {
				res = append(res, &u)
			} else {
				fmt.Println(e)
			}

		}

		return res, nil
	}
}

func (m *userQueryService) InsertNew(id int, name string, avatar string, ctx context.Context) error {
	conn, _ := m.factory.GetConnection(ctx)
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, "INSERT INTO VkUserModel (Id, Avatar, Name) VALUES ($1, $2, $3)", id, avatar, name); err != nil {
		return err
	}

	return nil
}

func (m *userQueryService) Delete(id int, ctx context.Context) error {
	conn, _ := m.factory.GetConnection(ctx)

	defer conn.Close()

	if _, err := conn.ExecContext(ctx, "DELETE FROM VkUserModel where id = $1", id); err != nil {
		return err
	}

	return nil
}

func NewUserQueryService(conn ConnectionFactory) UsersQueryService {
	return &userQueryService{
		factory: conn,
	}
}
