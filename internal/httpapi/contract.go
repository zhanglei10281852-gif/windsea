package httpapi

import (
	"encoding/json"
	"net/http"
)

type ErrorEnvelope struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func writeError(w http.ResponseWriter, status int, code, message, requestID string) {
	writeJSON(w, status, ErrorEnvelope{Code: code, Message: message, RequestID: requestID})
}
func decodeStrict(r *http.Request, value any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(value)
}
