package httpapi_test

import (
	"bytes"
	"encoding/json"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/httpapi"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthAndReady(t *testing.T) {
	_, services := testkit.Open(t)
	server := httptest.NewServer(httpapi.NewRouter(services, slog.Default()))
	defer server.Close()
	for _, path := range []string{"/healthz", "/readyz"} {
		response, err := http.Get(server.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != http.StatusOK {
			t.Fatalf("%s status=%d", path, response.StatusCode)
		}
		response.Body.Close()
	}
}
func TestCreateCampaignEndpoint(t *testing.T) {
	database, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)
	server := httptest.NewServer(httpapi.NewRouter(services, slog.Default()))
	defer server.Close()
	payload, _ := json.Marshal(map[string]any{"id": "http-c", "farmID": farm.ID, "name": "HTTP campaign", "startAt": time.Now().UTC(), "endAt": time.Now().Add(time.Hour).UTC()})
	response, err := http.Post(server.URL+"/v1/campaigns", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("status=%d", response.StatusCode)
	}
}
func TestInvalidWorkOrderReturnsJSONError(t *testing.T) {
	_, services := testkit.Open(t)
	server := httptest.NewServer(httpapi.NewRouter(services, slog.Default()))
	defer server.Close()
	payload, _ := json.Marshal(domain.WorkOrder{ID: "", Title: ""})
	response, err := http.Post(server.URL+"/v1/work-orders", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("status=%d", response.StatusCode)
	}
	var body map[string]any
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["code"] != "validation_error" {
		t.Fatalf("body=%v", body)
	}
}
