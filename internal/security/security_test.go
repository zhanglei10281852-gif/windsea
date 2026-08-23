package security_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/security"
	"testing"
	"time"
)

func TestSignature(t *testing.T) {
	signature := security.Sign("secret", "payload")
	if !security.Verify("secret", "payload", signature) {
		t.Fatal("signature failed")
	}
	if security.Verify("secret", "other", signature) {
		t.Fatal("wrong payload accepted")
	}
	if security.RedactToken("abcdefgh1234") != "abcd…1234" {
		t.Fatal("redaction")
	}
	if !security.IsSensitiveKey("Authorization") {
		t.Fatal("sensitive key")
	}
}
func TestSignatureFreshness(t *testing.T) {
	now := time.Now()
	if err := security.Fresh(now, now.Add(-time.Minute), time.Hour); err != nil {
		t.Fatal(err)
	}
	if err := security.Fresh(now, now.Add(-2*time.Hour), time.Hour); err == nil {
		t.Fatal("expired signature")
	}
}
