package database_test

import (
	"context"
	"github.com/woodhds/vk.service/database"
	"testing"
)

func TestGetConnection(t *testing.T) {
	connString := "host=localhost user=postgres port=5432 dbname=postgres"
	factory, _ := database.NewConnectionFactory(&connString)
	for i := 0; i < 50; i++ {
		con, _ := factory.GetConnection(context.Background())
		con.Close()
	}

	if stat := factory.Info(); stat.OpenConnections > 20 {
		t.Error("Max opened connections > 20")
	}
}
