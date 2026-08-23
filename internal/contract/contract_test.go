package contract_test

import (
	"errors"
	"github.com/zhanglei10281852-gif/windsea/internal/contract"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"net/http"
	"testing"
	"time"
)

func TestErrorMapping(t *testing.T) {
	apiErr := contract.FromError(domain.ErrConflict, "req", time.Now())
	if apiErr.Status != http.StatusConflict || apiErr.Code != "conflict" {
		t.Fatalf("error=%+v", apiErr)
	}
	if !errors.Is(domain.ErrConflict, domain.ErrConflict) {
		t.Fatal("sentinel")
	}
}
func TestPageParsing(t *testing.T) {
	page, err := contract.ParsePage("20", "40")
	if err != nil || page.Limit != 20 || page.Offset != 40 {
		t.Fatalf("page=%+v err=%v", page, err)
	}
	if _, err := contract.ParsePage("0", "0"); err == nil {
		t.Fatal("invalid page accepted")
	}
	if contract.TotalPages(41, 20) != 3 {
		t.Fatal("pages")
	}
}
