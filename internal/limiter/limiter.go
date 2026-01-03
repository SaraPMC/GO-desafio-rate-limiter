package limiter

import (
	"context"
	"fmt"

	"github.com/SaraPMC/GO-desafio-rate-limiter/internal/strategy"
)

// RateLimiter handles the rate limiting logic separated from middleware
type RateLimiter struct {
	storage              strategy.StorageStrategy
	defaultLimit         int
	defaultBlockDuration int
	tokenLimits          map[string]int
	tokenBlockDurations  map[string]int
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(
	storage strategy.StorageStrategy,
	defaultLimit int,
	defaultBlockDuration int,
) *RateLimiter {
	return &RateLimiter{
		storage:              storage,
		defaultLimit:         defaultLimit,
		defaultBlockDuration: defaultBlockDuration,
		tokenLimits:          make(map[string]int),
		tokenBlockDurations:  make(map[string]int),
	}
}

// ConfigureToken sets a custom limit and block duration for a specific token
func (rl *RateLimiter) ConfigureToken(token string, limit int, blockDuration int) {
	rl.tokenLimits[token] = limit
	rl.tokenBlockDurations[token] = blockDuration
}

// Allow checks if a request should be allowed based on IP or token
func (rl *RateLimiter) Allow(ctx context.Context, ip string, token string) (bool, error) {
	// Token takes precedence over IP
	if token != "" {
		return rl.checkToken(ctx, token)
	}

	return rl.checkIP(ctx, ip)
}

// checkIP checks if the request from the IP should be allowed
func (rl *RateLimiter) checkIP(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("limiter:ip:%s", ip)
	return rl.check(ctx, key, rl.defaultLimit, rl.defaultBlockDuration)
}

// checkToken checks if the request with the token should be allowed
func (rl *RateLimiter) checkToken(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("limiter:token:%s", token)

	limit := rl.defaultLimit
	blockDuration := rl.defaultBlockDuration

	if customLimit, exists := rl.tokenLimits[token]; exists {
		limit = customLimit
	}

	if customDuration, exists := rl.tokenBlockDurations[token]; exists {
		blockDuration = customDuration
	}

	return rl.check(ctx, key, limit, blockDuration)
}

// check performs the actual rate limit check
func (rl *RateLimiter) check(ctx context.Context, key string, limit int, blockDuration int) (bool, error) {
	// Get current counter
	counter, err := rl.storage.GetCounter(ctx, key)
	if err != nil {
		return false, err
	}

	// If blocked, deny the request
	if counter >= limit {
		return false, nil
	}

	// Increment counter
	newCounter, err := rl.storage.IncrementCounter(ctx, key)
	if err != nil {
		return false, err
	}

	// Set expiration only on first request in the window
	if newCounter == 1 {
		err = rl.storage.SetExpiration(ctx, key, blockDuration)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// Reset resets the counter for a specific key (useful for testing)
func (rl *RateLimiter) Reset(ctx context.Context, ip string, token string) error {
	if token != "" {
		key := fmt.Sprintf("limiter:token:%s", token)
		return rl.storage.Delete(ctx, key)
	}

	if ip != "" {
		key := fmt.Sprintf("limiter:ip:%s", ip)
		return rl.storage.Delete(ctx, key)
	}

	return nil
}

// Close closes the underlying storage connection
func (rl *RateLimiter) Close() error {
	return rl.storage.Close()
}
