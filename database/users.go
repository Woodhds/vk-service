package database

import (
	"database/sql"
	"errors"
)

type UsersQueryService interface {
	GetAll() ([]int, error)
}

type userQueryService struct {
	conn *sql.DB
}

func (m userQueryService) GetAll() ([]int, error) {
	rows, _ := m.conn.Query(`select Id from VkUserModel`)
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

func NewUserQueryService(conn *sql.DB) (UsersQueryService, error) {
	if conn == nil {
		return nil, errors.New("factory empty")
	}

	return &userQueryService{
		conn: conn,
	}, nil
}
