package ledger_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/ledger"
	"testing"
	"time"
)

func TestLedgerBalanceAndOrdering(t *testing.T) {
	now := time.Now()
	entries := []ledger.Entry{{ID: "2", Account: "parts", Delta: -3, At: now.Add(time.Minute)}, {ID: "1", Account: "parts", Delta: 5, At: now}}
	if ledger.Balance(entries, "parts") != 2 {
		t.Fatal("balance")
	}
	if err := ledger.RequireNonNegative(entries, "parts"); err != nil {
		t.Fatal(err)
	}
	if !ledger.Ordered(entries)[0].At.Before(ledger.Ordered(entries)[1].At) {
		t.Fatal("order")
	}
	entries = append(entries, ledger.Entry{ID: "3", Account: "parts", Delta: -4, At: now.Add(2 * time.Minute)})
	if err := ledger.RequireNonNegative(entries, "parts"); err == nil {
		t.Fatal("negative balance accepted")
	}
}
