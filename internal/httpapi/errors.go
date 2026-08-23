package httpapi

import (
	"errors"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"net/http"
)

func status(err error) int {
	switch {
	case errors.Is(err, domain.ErrValidation):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrConflict), errors.Is(err, domain.ErrCapacity):
		return http.StatusConflict
	case errors.Is(err, domain.ErrInvalidState):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
func fail(w http.ResponseWriter, err error) {
	writeJSON(w, status(err), map[string]any{"code": code(err), "message": err.Error()})
}
func code(err error) string {
	switch {
	case errors.Is(err, domain.ErrValidation):
		return "validation_error"
	case errors.Is(err, domain.ErrForbidden):
		return "forbidden"
	case errors.Is(err, domain.ErrNotFound):
		return "not_found"
	case errors.Is(err, domain.ErrConflict):
		return "conflict"
	case errors.Is(err, domain.ErrCapacity):
		return "capacity"
	case errors.Is(err, domain.ErrInvalidState):
		return "invalid_state"
	default:
		return "internal_error"
	}
}
