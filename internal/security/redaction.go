package security

import "strings"

func RedactToken(token string) string {
	if len(token) <= 8 {
		return "[redacted]"
	}
	return token[:4] + "…" + token[len(token)-4:]
}
func SafeEmail(email string) string {
	at := strings.Index(email, "@")
	if at <= 1 {
		return "[redacted]"
	}
	return email[:1] + "***" + email[at:]
}
func IsSensitiveKey(key string) bool {
	switch strings.ToLower(key) {
	case "token", "password", "secret", "authorization", "cookie":
		return true
	default:
		return false
	}
}
