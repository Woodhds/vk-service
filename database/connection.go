package database

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

const (
	maxConnection     = 20
	connectionTimeout = 60
)

type ConnectionFactory interface {
	GetConnection(ctx context.Context) (*sql.Conn, error)
}

type connectionFactory struct {
	db *sql.DB
}

func (factory connectionFactory) GetConnection(ctx context.Context) (*sql.Conn, error) {
	return factory.db.Conn(ctx)
}

func NewConnectionFactory(connectionString *string) (ConnectionFactory, error) {
	if len(*connectionString) == 0 {
		return nil, errors.New("connection string is empty")
	}
	if conn, err := sql.Open("pgx", *connectionString); err != nil {
		return nil, err
	} else {
		if e := conn.Ping(); e != nil {
			return nil, e
		}

		conn.SetMaxOpenConns(maxConnection)
		conn.SetConnMaxLifetime(connectionTimeout * time.Second)

		return &connectionFactory{db: conn}, nil
	}
}
