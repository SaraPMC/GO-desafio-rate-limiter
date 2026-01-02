package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/seu-usuario/desafio-rate-limiter/internal/strategy"
)

// MockStorage is a mock implementation of StorageStrategy for testing
type MockStorage struct {
	counters map[string]int
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		counters: make(map[string]int),
	}
}

func (m *MockStorage) IncrementCounter(ctx context.Context, key string) (int, error) {
	m.counters[key]++
	return m.counters[key], nil
}

func (m *MockStorage) SetExpiration(ctx context.Context, key string, ttlSeconds int) error {
	return nil
}

func (m *MockStorage) GetCounter(ctx context.Context, key string) (int, error) {
	return m.counters[key], nil
}

func (m *MockStorage) Exists(ctx context.Context, key string) (bool, error) {
	_, exists := m.counters[key]
	return exists, nil
}

func (m *MockStorage) Delete(ctx context.Context, key string) error {
	delete(m.counters, key)
	return nil
}

func (m *MockStorage) Close() error {
	return nil
}

// TestIPRateLimit tests rate limiting by IP address
func TestIPRateLimit(t *testing.T) {
	storage := NewMockStorage()
	limiter := NewRateLimiter(storage, 5, 300)
	defer limiter.Close()

	ctx := context.Background()

	// Allow 5 requests
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(ctx, "192.168.1.1", "")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !allowed {
			t.Fatalf("Request %d should have been allowed", i+1)
		}
	}

	// 6th request should be denied
	allowed, err := limiter.Allow(ctx, "192.168.1.1", "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("6th request should have been denied")
	}
}

// TestTokenRateLimit tests rate limiting by token
func TestTokenRateLimit(t *testing.T) {
	storage := NewMockStorage()
	limiter := NewRateLimiter(storage, 5, 300)
	defer limiter.Close()

	// Configure token with limit of 10
	limiter.ConfigureToken("token123", 10, 600)

	ctx := context.Background()

	// Allow 10 requests with token
	for i := 0; i < 10; i++ {
		allowed, err := limiter.Allow(ctx, "192.168.1.1", "token123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !allowed {
			t.Fatalf("Request %d with token should have been allowed", i+1)
		}
	}

	// 11th request should be denied
	allowed, err := limiter.Allow(ctx, "192.168.1.1", "token123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("11th request with token should have been denied")
	}
}

// TestTokenPrecedence tests that token limits take precedence over IP limits
func TestTokenPrecedence(t *testing.T) {
	storage := NewMockStorage()
	limiter := NewRateLimiter(storage, 5, 300) // IP limit: 5
	defer limiter.Close()

	// Configure token with higher limit
	limiter.ConfigureToken("token123", 100, 600)

	ctx := context.Background()

	// Allow 10 requests with token (higher than IP limit of 5)
	for i := 0; i < 10; i++ {
		allowed, err := limiter.Allow(ctx, "192.168.1.1", "token123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !allowed {
			t.Fatalf("Request %d with high-limit token should have been allowed", i+1)
		}
	}
}

// TestReset tests the reset functionality
func TestReset(t *testing.T) {
	storage := NewMockStorage()
	limiter := NewRateLimiter(storage, 5, 300)
	defer limiter.Close()

	ctx := context.Background()

	// Fill up the limit
	for i := 0; i < 5; i++ {
		limiter.Allow(ctx, "192.168.1.1", "")
	}

	// Verify blocked
	allowed, _ := limiter.Allow(ctx, "192.168.1.1", "")
	if allowed {
		t.Fatal("Should be blocked at limit")
	}

	// Reset
	limiter.Reset(ctx, "192.168.1.1", "")

	// Should be allowed again
	allowed, err := limiter.Allow(ctx, "192.168.1.1", "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("Should be allowed after reset")
	}
}

// BenchmarkAllow benchmarks the Allow function
func BenchmarkAllow(b *testing.B) {
	storage := NewMockStorage()
	limiter := NewRateLimiter(storage, 5, 300)
	defer limiter.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow(ctx, "192.168.1.1", "")
	}
}
