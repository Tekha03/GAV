package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLimiterAllowsRequestsWithinBurst(t *testing.T) {
	limiter := New(Config{
		Requests: 2,
		Window:   time.Minute,
		Burst:    2,
		KeyFunc: func(r *http.Request) string {
			return "user-1"
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	if allowed, _ := limiter.Allow(req); !allowed {
		t.Fatal("first request should be allowed")
	}
	if allowed, _ := limiter.Allow(req); !allowed {
		t.Fatal("second request should be allowed")
	}
	if allowed, retryAfter := limiter.Allow(req); allowed || retryAfter <= 0 {
		t.Fatalf("third request should be limited with retryAfter, allowed=%v retryAfter=%s", allowed, retryAfter)
	}
}
