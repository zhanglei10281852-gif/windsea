package encoding_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/encoding"
	"testing"
)

func TestCursorRoundTrip(t *testing.T) {
	value := encoding.Cursor{At: "2026-01-01T00:00:00Z", ID: "wo-1"}
	encoded := encoding.EncodeCursor(value)
	decoded, err := encoding.DecodeCursor(encoded)
	if err != nil || decoded != value {
		t.Fatalf("decoded=%+v err=%v", decoded, err)
	}
	if _, err := encoding.DecodeCursor("!"); err == nil {
		t.Fatal("invalid cursor accepted")
	}
}
