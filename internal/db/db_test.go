package db_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zhanglei10281852-gif/windsea/internal/db"
	"testing"
)

func TestMigrationsCreateRelationalSchema(t *testing.T) {
	database, err := db.Open(context.Background(), "file:migrations?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	enabled, err := db.ForeignKeysEnabled(context.Background(), database.DB)
	if err != nil || !enabled {
		t.Fatalf("foreign keys enabled=%v err=%v", enabled, err)
	}
	count, err := db.TableCount(context.Background(), database.DB)
	if err != nil || count < 15 {
		t.Fatalf("tables=%d err=%v", count, err)
	}
	if err := db.Migrate(context.Background(), database); err != nil {
		t.Fatal(err)
	}
}
func TestTransactionRollback(t *testing.T) {
	database, err := db.Open(context.Background(), "file:rollback?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	err = database.WithTx(context.Background(), func(tx *sql.Tx) error {
		if _, err := tx.Exec("INSERT INTO farms(id,name,timezone,created_at) VALUES('f','F','UTC',CURRENT_TIMESTAMP)"); err != nil {
			return err
		}
		return errors.New("stop")
	})
	if err == nil {
		t.Fatal("expected rollback error")
	}
	var count int
	if err := database.QueryRow("SELECT COUNT(*) FROM farms").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("rollback left %d rows", count)
	}
}
