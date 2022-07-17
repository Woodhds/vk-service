package database

import (
	"context"
	"database/sql"
)

type DistributedLock interface {
	Lock(id int64, ctx context.Context) bool
	Unlock(id int64, ctx context.Context)
}

type distributedLockImpl struct {
	connection *sql.Conn
}

func (d distributedLockImpl) Unlock(id int64, ctx context.Context) {
	d.connection.ExecContext(ctx, "select pg_advisory_unlock($1)", id)
}

func (d distributedLockImpl) Lock(id int64, context context.Context) bool {
	row := d.connection.QueryRowContext(context, "select pg_try_advisory_lock($1)", id)
	isLock := false

	row.Scan(&isLock)

	return isLock
}

func NewDistributedLock(connection *sql.Conn) DistributedLock {
	return &distributedLockImpl{connection}
}
