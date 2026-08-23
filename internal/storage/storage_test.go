package storage_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/storage"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
	"testing"
	"time"
)

func TestTxStoreAndLockTable(t *testing.T) {
	database, _ := testkit.Open(t)
	store := storage.TxStore{DB: database.DB}
	if err := store.EnsureFarm(context.Background(), "missing"); err == nil {
		t.Fatal("missing farm accepted")
	}
	locks := &storage.LockTable{DB: database.DB}
	if err := locks.Acquire(context.Background(), "campaign-1", time.Minute); err != nil {
		t.Fatal(err)
	}
	if err := locks.Acquire(context.Background(), "campaign-1", time.Minute); err == nil {
		t.Fatal("lock acquired twice")
	}
	if err := locks.Release(context.Background(), "campaign-1"); err != nil {
		t.Fatal(err)
	}
}
