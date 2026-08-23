package route_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/route"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouteTable(t *testing.T) {
	table := &route.Table{}
	table.Add("GET", "/v1/work-orders/:id", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	request := httptest.NewRequest("GET", "/v1/work-orders/wo-1", nil)
	response := httptest.NewRecorder()
	table.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("code=%d", response.Code)
	}
	request = httptest.NewRequest("POST", "/v1/work-orders/wo-1", nil)
	response = httptest.NewRecorder()
	table.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("code=%d", response.Code)
	}
}
