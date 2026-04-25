package payment

import (
	"context"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"time"
)

// RetryConfig holds the retry configuration
type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

// ExecutePayment executes a payment with automatic retry on temporary failures
func ExecutePayment(ctx context.Context, cfg RetryConfig, baseURL string, payload []byte) (int, []byte, error) {
	var lastErr error

	for attempt := 0; attempt < cfg.MaxRetries; attempt++ {
		if ctx.Err() != nil {
			return 0, nil, fmt.Errorf("context cancelled/expired: %w", ctx.Err())
		}

		resp, err := http.Post(baseURL+"/payment", "application/json", nil)

		if err != nil {
			lastErr = err
			if attempt < cfg.MaxRetries-1 {
				delay := calculateBackoff(attempt, cfg.BaseDelay, cfg.MaxDelay)
				fmt.Printf("Attempt %d failed: %v, waiting %v before next retry...\n", attempt+1, err, delay)
				time.Sleep(delay)
			}
			continue
		}

		if isRetryable(resp) {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("status %d: %s", resp.StatusCode, string(body))

			if attempt < cfg.MaxRetries-1 {
				delay := calculateBackoff(attempt, cfg.BaseDelay, cfg.MaxDelay)
				fmt.Printf("Attempt %d failed: status %d, waiting %v before next retry...\n", attempt+1, resp.StatusCode, delay)
				time.Sleep(delay)
			}
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return resp.StatusCode, nil, err
		}
		return resp.StatusCode, body, nil
	}

	if lastErr != nil {
		return 0, nil, lastErr
	}
	return 0, nil, fmt.Errorf("max retries exceeded")
}

func isRetryable(resp *http.Response) bool {
	retryableCodes := map[int]bool{
		429: true,
		500: true,
		502: true,
		503: true,
		504: true,
	}

	nonRetryableCodes := map[int]bool{
		400: true,
		401: true,
		403: true,
		404: true,
	}

	if nonRetryableCodes[resp.StatusCode] {
		return false
	}

	return retryableCodes[resp.StatusCode] || (resp.StatusCode >= 500)
}

func calculateBackoff(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	backoffTime := float64(baseDelay) * math.Pow(2, float64(attempt))

	if backoffTime > float64(maxDelay) {
		backoffTime = float64(maxDelay)
	}

	jitter := time.Duration(rand.Int63n(int64(backoffTime)))
	return jitter
}

// IsRetryable checks if a response code indicates a temporary error
func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}

	retryableCodes := map[int]bool{
		429: true,
		500: true,
		502: true,
		503: true,
		504: true,
	}

	return retryableCodes[resp.StatusCode] || (resp.StatusCode >= 500)
}

// CalculateBackoff calculates exponential backoff with full jitter
func CalculateBackoff(attempt int) time.Duration {
	baseDelay := 100 * time.Millisecond
	maxDelay := 5 * time.Second

	backoffTime := float64(baseDelay) * math.Pow(2, float64(attempt))
	if backoffTime > float64(maxDelay) {
		backoffTime = float64(maxDelay)
	}

	jitter := time.Duration(rand.Int63n(int64(backoffTime)))
	return jitter
}
