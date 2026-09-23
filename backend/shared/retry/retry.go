package retry

import (
	"context"
	"errors"
	"io"
	"math/rand/v2"
	"net"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	Attemps    int
	Jitter     float64
	Timeout    time.Duration
	Delays     []time.Duration
	RetryCheck func(error) bool
}

func IsRetryableGRPC(err error) bool {
	code := status.Code(err)

	switch code {
	case codes.Unavailable,
		codes.DeadlineExceeded,
		codes.ResourceExhausted:
		return true
	default:
		return false
	}
}

func IsRetryableHTTPStatus(status int) bool {
	switch status {
	case 429, 502, 503, 504:
		return true
	default:
		return false
	}
}

func IsRetryableNetwork(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, io.EOF) {
		return true
	}

	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func IsRetryable(err error) bool {
	return IsRetryableNetwork(err) || IsRetryableGRPC(err)
}

func DefaultConfig() Config {
	return Config{
		Attemps: 3,
		Jitter:  0.25,
		Timeout: 2 * time.Second,
		Delays: []time.Duration{
			100 * time.Millisecond,
			300 * time.Millisecond,
			700 * time.Millisecond,
		},
		RetryCheck: IsRetryable,
	}
}

func Do(ctx context.Context, cfg Config, operation func(context.Context) error) error {
	_, err := DoValue(ctx, cfg, func(ctx context.Context) (struct{}, error) {
		return struct{}{}, operation(ctx)
	})
	return err
}

func DoValue[T any](ctx context.Context, cfg Config, operation func(context.Context) (T, error)) (T, error) {
	var zero T

	if ctx == nil {
		ctx = context.Background()
	}
	if operation == nil {
		return zero, nil
	}

	cfg = withDefaults(cfg)

	var lastErr error
	for attempt := 0; attempt < cfg.Attemps; attempt++ {
		attemptCtx := ctx
		cancel := func() {}
		if cfg.Timeout > 0 {
			attemptCtx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		}

		result, err := operation(attemptCtx)
		cancel()

		if err == nil {
			return result, nil
		}

		lastErr = err
		if attempt == cfg.Attemps-1 || !cfg.RetryCheck(err) {
			return zero, err
		}

		if err := wait(ctx, delayForAttempt(cfg, attempt)); err != nil {
			return zero, err
		}
	}

	return zero, lastErr
}

func withDefaults(cfg Config) Config {
	defaults := DefaultConfig()

	if cfg.Attemps <= 0 {
		cfg.Attemps = defaults.Attemps
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = defaults.Timeout
	}
	if len(cfg.Delays) == 0 {
		cfg.Delays = defaults.Delays
	}
	if cfg.RetryCheck == nil {
		cfg.RetryCheck = defaults.RetryCheck
	}
	if cfg.Jitter < 0 {
		cfg.Jitter = 0
	}

	return cfg
}

func delayForAttempt(cfg Config, attempt int) time.Duration {
	delay := cfg.Delays[len(cfg.Delays)-1]
	if attempt < len(cfg.Delays) {
		delay = cfg.Delays[attempt]
	}

	if cfg.Jitter <= 0 || delay <= 0 {
		return delay
	}

	delta := int64(float64(delay) * cfg.Jitter)
	if delta <= 0 {
		return delay
	}

	minDelay := int64(delay) - delta
	maxDelay := int64(delay) + delta
	return time.Duration(rand.Int64N(maxDelay-minDelay+1) + minDelay)
}

func wait(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return nil
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
