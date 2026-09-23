package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDoRetriesRetryableErrorAndThenSucceeds(t *testing.T) {
	attempts := 0

	err := Do(context.Background(), testConfig(), func(ctx context.Context) error {
		attempts++
		if attempts == 1 {
			return status.Error(codes.Unavailable, "service unavailable")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Do() error = %v, want nil", err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

func TestDoDoesNotRetryNonRetryableError(t *testing.T) {
	businessErr := errors.New("validation failed")
	attempts := 0

	err := Do(context.Background(), testConfig(), func(ctx context.Context) error {
		attempts++
		return businessErr
	})

	if !errors.Is(err, businessErr) {
		t.Fatalf("Do() error = %v, want %v", err, businessErr)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}

func TestDoStopsWhenContextIsCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0

	err := Do(ctx, testConfig(), func(ctx context.Context) error {
		attempts++
		cancel()
		return status.Error(codes.Unavailable, "service unavailable")
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Do() error = %v, want context.Canceled", err)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}

func TestDoValueReturnsResult(t *testing.T) {
	result, err := DoValue(context.Background(), testConfig(), func(ctx context.Context) (string, error) {
		return "ok", nil
	})

	if err != nil {
		t.Fatalf("DoValue() error = %v, want nil", err)
	}
	if result != "ok" {
		t.Fatalf("DoValue() result = %q, want %q", result, "ok")
	}
}

func TestIsRetryableHTTPStatus(t *testing.T) {
	for _, statusCode := range []int{429, 502, 503, 504} {
		if !IsRetryableHTTPStatus(statusCode) {
			t.Fatalf("IsRetryableHTTPStatus(%d) = false, want true", statusCode)
		}
	}

	for _, statusCode := range []int{400, 401, 403, 404, 409} {
		if IsRetryableHTTPStatus(statusCode) {
			t.Fatalf("IsRetryableHTTPStatus(%d) = true, want false", statusCode)
		}
	}
}

func testConfig() Config {
	return Config{
		Attemps: 3,
		Timeout: time.Second,
		Delays:  []time.Duration{time.Millisecond},
		Jitter:  0,
	}
}
