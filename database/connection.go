package database

import (
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type ConnectionFactory interface {
	GetConnection() (*sql.DB, error)
}

type connectionFactory struct {
	connectionString string
}

func (factory connectionFactory) GetConnection() (*sql.DB, error) {
	if conn, err := sql.Open("pgx", factory.connectionString); err != nil {
		return nil, err
	} else {
		return conn, nil
	}
}

func NewConnectionFactory(connectionString *string) (ConnectionFactory, error) {
	if len(*connectionString) == 0 {
		return nil, errors.New("connection string is empty")
	}
	return &connectionFactory{connectionString: *connectionString}, nil
}
