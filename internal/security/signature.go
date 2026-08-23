package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

func Sign(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
func Verify(secret, payload, signature string) bool {
	expected := Sign(secret, payload)
	return hmac.Equal([]byte(strings.ToLower(expected)), []byte(strings.ToLower(signature)))
}
func Fresh(now, issued time.Time, ttl time.Duration) error {
	if issued.After(now) || now.Sub(issued) > ttl {
		return errors.New("signature expired")
	}
	return nil
}
