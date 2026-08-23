package domain

import "errors"

var (
	ErrNotFound     = errors.New("windsea: not found")
	ErrConflict     = errors.New("windsea: conflict")
	ErrInvalidState = errors.New("windsea: invalid state transition")
	ErrForbidden    = errors.New("windsea: forbidden")
	ErrValidation   = errors.New("windsea: validation failed")
	ErrCancelled    = errors.New("windsea: operation cancelled")
	ErrCapacity     = errors.New("windsea: capacity exceeded")
)
