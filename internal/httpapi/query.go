package httpapi

import (
	"net/http"
	"strconv"
)

func intQuery(r *http.Request, key string, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil {
		return fallback
	}
	return value
}
func boolQuery(r *http.Request, key string, fallback bool) bool {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback
	}
	return raw == "true" || raw == "1"
}
