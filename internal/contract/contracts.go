package contract

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"net/http"
	"time"
)

type APIError struct {
	Code, Message, RequestID string
	Status                   int
	At                       time.Time
}

func (e APIError) Error() string { return e.Code + ": " + e.Message }
func FromError(err error, requestID string, at time.Time) APIError {
	status := http.StatusInternalServerError
	code := "internal_error"
	switch {
	case err == nil:
		status = http.StatusOK
		code = "ok"
	case errors.Is(err, domain.ErrValidation):
		status = http.StatusBadRequest
		code = "validation_error"
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
		code = "not_found"
	case errors.Is(err, domain.ErrConflict):
		status = http.StatusConflict
		code = "conflict"
	case errors.Is(err, domain.ErrForbidden):
		status = http.StatusForbidden
		code = "forbidden"
	}
	return APIError{Code: code, Message: message(err), RequestID: requestID, Status: status, At: at}
}
func message(err error) string {
	if err == nil {
		return "ok"
	}
	return err.Error()
}
func EncodeError(err error) []byte { raw, _ := json.Marshal(err); return raw }
func ValidateContentType(r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("%w: content type", domain.ErrValidation)
	}
	return nil
}
