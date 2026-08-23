package encoding

import (
	"encoding/base64"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"strconv"
	"strings"
)

type Cursor struct {
	At string
	ID string
}

func EncodeCursor(cursor Cursor) string {
	raw := cursor.At + "|" + cursor.ID
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}
func DecodeCursor(value string) (Cursor, error) {
	if value == "" {
		return Cursor{}, nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return Cursor{}, fmt.Errorf("%w: cursor", domain.ErrValidation)
	}
	parts := strings.Split(string(raw), "|")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Cursor{}, domain.ErrValidation
	}
	return Cursor{At: parts[0], ID: parts[1]}, nil
}
func ParseInt(value string) (int, error) {
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%w: integer", domain.ErrValidation)
	}
	return number, nil
}
