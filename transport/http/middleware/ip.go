package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"
)

type realIPKey struct{}

func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realIP(r)
		ctx := context.WithValue(r.Context(), realIPKey{}, ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRealIP(ctx context.Context) (string, bool) {
	ip, ok := ctx.Value(realIPKey{}).(string)
	return ip, ok
}

func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}

func Timeout(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, 30*time.Second, `{"error":"request timeout}`)
}
