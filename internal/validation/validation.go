package validation

import (
	"fmt"
	"strings"

	"github.com/zhanglei10281852-gif/windsea/internal/domain"
)

func Required(fields ...string) error {
	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			return fmt.Errorf("%w: required value", domain.ErrValidation)
		}
	}
	return nil
}
func Positive(value int) error {
	if value <= 0 {
		return fmt.Errorf("%w: value must be positive", domain.ErrValidation)
	}
	return nil
}
func Range(value, min, max int) error {
	if value < min || value > max {
		return fmt.Errorf("%w: value outside range", domain.ErrValidation)
	}
	return nil
}
