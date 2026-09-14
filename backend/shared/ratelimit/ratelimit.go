package ratelimit

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	apperrors "shared/app_errors"
)

type Config struct {
	Requests int
	Window   time.Duration
	Burst    int
	KeyFunc  func(*http.Request) string
}

type bucket struct {
	tokens float64
	last   time.Time
	seen   time.Time
}

type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     float64
	capacity float64
	window   time.Duration
	keyFunc  func(*http.Request) string
	now      func() time.Time
}

func New(config Config) *Limiter {
	if config.Requests <= 0 {
		config.Requests = 60
	}
	if config.Window <= 0 {
		config.Window = time.Minute
	}
	if config.Burst <= 0 {
		config.Burst = config.Requests
	}
	if config.KeyFunc == nil {
		config.KeyFunc = IPKey
	}

	return &Limiter{
		buckets:  make(map[string]*bucket),
		rate:     float64(config.Requests) / config.Window.Seconds(),
		capacity: float64(config.Burst),
		window:   config.Window,
		keyFunc:  config.KeyFunc,
		now:      time.Now,
	}
}

func Middleware(config Config) func(http.Handler) http.Handler {
	limiter := New(config)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowed, retryAfter := limiter.Allow(r)
			if !allowed {
				WriteTooManyRequests(w, retryAfter)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (l *Limiter) Allow(r *http.Request) (bool, time.Duration) {
	key := l.keyFunc(r)
	if key == "" {
		key = "unknown"
	}

	now := l.now()

	l.mu.Lock()
	defer l.mu.Unlock()

	b := l.buckets[key]
	if b == nil {
		b = &bucket{tokens: l.capacity, last: now, seen: now}
		l.buckets[key] = b
	}

	elapsed := now.Sub(b.last).Seconds()
	b.tokens += elapsed * l.rate
	if b.tokens > l.capacity {
		b.tokens = l.capacity
	}
	b.last = now
	b.seen = now

	if b.tokens >= 1 {
		b.tokens--
		l.prune(now)
		return true, 0
	}

	missing := 1 - b.tokens
	retryAfter := time.Duration(missing/l.rate*float64(time.Second) + float64(time.Millisecond-1))
	l.prune(now)
	return false, retryAfter
}

func (l *Limiter) prune(now time.Time) {
	maxIdle := 2 * l.window
	if maxIdle < time.Minute {
		maxIdle = time.Minute
	}

	for key, b := range l.buckets {
		if now.Sub(b.seen) > maxIdle {
			delete(l.buckets, key)
		}
	}
}

func IPKey(r *http.Request) string {
	if ip := strings.TrimSpace(r.Header.Get("X-Real-IP")); ip != "" {
		return ip
	}

	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return strings.TrimSpace(r.RemoteAddr)
	}

	return host
}

func UserIDKey(contextKey any) func(*http.Request) string {
	return func(r *http.Request) string {
		if id, ok := r.Context().Value(contextKey).(interface{ String() string }); ok {
			return id.String()
		}

		return IPKey(r)
	}
}

func PathKey(base func(*http.Request) string, names ...string) func(*http.Request) string {
	return func(r *http.Request) string {
		key := base(r)
		for _, name := range names {
			value := strings.TrimSpace(r.PathValue(name))
			if value != "" {
				key += ":" + value
			}
		}
		return key
	}
}

func WriteTooManyRequests(w http.ResponseWriter, retryAfter time.Duration) {
	seconds := int(retryAfter.Round(time.Second).Seconds())
	if seconds < 1 {
		seconds = 1
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", strconv.Itoa(seconds))
	w.WriteHeader(http.StatusTooManyRequests)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"code":     apperrors.CodeRateLimited,
			"category": apperrors.CategoryRateLimited,
			"message":  "rate limit exceeded",
			"details": map[string]any{
				"retry_after_seconds": seconds,
			},
		},
	})
}
