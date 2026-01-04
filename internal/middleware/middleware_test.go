package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SaraPMC/GO-desafio-rate-limiter/internal/limiter"
)

// MockStorage for testing
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

// TestRateLimiterMiddlewareAcceptsRequest tests that valid requests are accepted
func TestRateLimiterMiddlewareAcceptsRequest(t *testing.T) {
	storage := NewMockStorage()
	rateLimiter := limiter.NewRateLimiter(storage, 5, 300)
	defer rateLimiter.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	middleware := RateLimiterMiddleware(rateLimiter)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")

	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

// TestRateLimiterMiddlewareBlocksExceededRequests tests that requests exceeding limit are blocked
func TestRateLimiterMiddlewareBlocksExceededRequests(t *testing.T) {
	storage := NewMockStorage()
	rateLimiter := limiter.NewRateLimiter(storage, 3, 300) // Limit to 3 requests
	defer rateLimiter.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RateLimiterMiddleware(rateLimiter)
	wrappedHandler := middleware(handler)

	ip := "192.168.1.100"

	// Make 3 allowed requests
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("X-Forwarded-For", ip)
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Request %d should be allowed, got status %d", i+1, w.Code)
		}
	}

	// 4th request should be blocked
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("X-Forwarded-For", ip)
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("Expected status 429, got %d", w.Code)
	}

	// Check response message
	var respBody map[string]interface{}
	json.NewDecoder(w.Body).Decode(&respBody)

	expectedMsg := "you have reached the maximum number of requests or actions allowed within a certain time frame"
	if respBody["error"] != expectedMsg {
		t.Fatalf("Expected error message: %s, got: %v", expectedMsg, respBody["error"])
	}
}

// TestRateLimiterMiddlewareTokenPrecedence tests that token limits take precedence
func TestRateLimiterMiddlewareTokenPrecedence(t *testing.T) {
	storage := NewMockStorage()
	rateLimiter := limiter.NewRateLimiter(storage, 2, 300) // IP limit: 2
	rateLimiter.ConfigureToken("premium-token", 5, 600)    // Token limit: 5

	defer rateLimiter.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RateLimiterMiddleware(rateLimiter)
	wrappedHandler := middleware(handler)

	ip := "192.168.1.50"

	// Make 5 requests with token (should succeed despite IP limit of 2)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("X-Forwarded-For", ip)
		req.Header.Set("API_KEY", "premium-token")
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Request %d with token should be allowed, got status %d", i+1, w.Code)
		}
	}

	// 6th request should be blocked
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("X-Forwarded-For", ip)
	req.Header.Set("API_KEY", "premium-token")
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("Expected status 429 after token limit, got %d", w.Code)
	}
}

// TestRateLimiterMiddlewareExtractsIPFromXForwardedFor tests IP extraction
func TestRateLimiterMiddlewareExtractsIPFromXForwardedFor(t *testing.T) {
	storage := NewMockStorage()
	rateLimiter := limiter.NewRateLimiter(storage, 1, 300) // Limit 1 per IP
	defer rateLimiter.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RateLimiterMiddleware(rateLimiter)
	wrappedHandler := middleware(handler)

	// First request from IP1
	req1 := httptest.NewRequest("GET", "/api/test", nil)
	req1.Header.Set("X-Forwarded-For", "192.168.1.1")
	w1 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("First request should be allowed, got status %d", w1.Code)
	}

	// Second request from SAME IP should be blocked
	req2 := httptest.NewRequest("GET", "/api/test", nil)
	req2.Header.Set("X-Forwarded-For", "192.168.1.1")
	w2 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("Second request from same IP should be blocked, got status %d", w2.Code)
	}

	// First request from DIFFERENT IP should succeed
	req3 := httptest.NewRequest("GET", "/api/test", nil)
	req3.Header.Set("X-Forwarded-For", "192.168.1.2")
	w3 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Fatalf("First request from different IP should be allowed, got status %d", w3.Code)
	}
}
