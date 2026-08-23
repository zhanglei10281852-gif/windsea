package pagination_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/pagination"
	"net/url"
	"testing"
)

func TestPageQueryBounds(t *testing.T) {
	page := pagination.FromQuery(url.Values{"limit": []string{"999"}, "offset": []string{"-1"}})
	if page.Limit != 200 || page.Offset != 0 {
		t.Fatalf("page=%+v", page)
	}
	page = pagination.FromQuery(url.Values{"limit": []string{"bad"}})
	if page.Limit != 50 {
		t.Fatalf("fallback=%+v", page)
	}
}
