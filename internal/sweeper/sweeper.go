package sweeper

import (
	"github.com/woodhds/vk.service/database"
	"golang.org/x/net/context"
)

type MessageSweeper interface {
	Run(ctx context.Context)
}

type sweeperImplementation struct {
	db database.ConnectionFactory
}

func (ms *sweeperImplementation) Run(ctx context.Context) {
	ms.initialize(ctx)
}

func (ms *sweeperImplementation) initialize(ctx context.Context) {
	conn, _ := ms.db.GetConnection(ctx)

	defer conn.Close()

	conn.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS maintenance (id INTEGER UNIQUE, last_execution DATETIME)")
}

func NewSweeper(db database.ConnectionFactory) MessageSweeper {
	return &sweeperImplementation{db: db}
}
