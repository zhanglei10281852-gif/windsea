package validation_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/validation"
	"testing"
)

func TestValidationRules(t *testing.T) {
	if err := validation.Required("a", "b"); err != nil {
		t.Fatal(err)
	}
	if err := validation.Required(" "); err == nil {
		t.Fatal("blank accepted")
	}
	if err := validation.Positive(1); err != nil {
		t.Fatal(err)
	}
	if err := validation.Positive(0); err == nil {
		t.Fatal("zero accepted")
	}
	if err := validation.Range(3, 1, 5); err != nil {
		t.Fatal(err)
	}
	if err := validation.Range(7, 1, 5); err == nil {
		t.Fatal("out of range accepted")
	}
}
