package cache_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/cache"
	"testing"
	"time"
)

func TestCacheExpiryAndPurge(t *testing.T) {
	now := time.Now()
	current := now
	c := cache.New[string](func() time.Time { return current })
	c.Set("a", "value", time.Minute)
	if value, ok := c.Get("a"); !ok || value != "value" {
		t.Fatalf("value=%s ok=%v", value, ok)
	}
	current = current.Add(2 * time.Minute)
	if _, ok := c.Get("a"); ok {
		t.Fatal("expired value returned")
	}
	c.Set("b", "other", time.Minute)
	current = current.Add(2 * time.Minute)
	c.Purge()
	if _, ok := c.Get("b"); ok {
		t.Fatal("purge failed")
	}
}
