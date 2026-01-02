package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/seu-usuario/desafio-rate-limiter/internal/config"
	"github.com/seu-usuario/desafio-rate-limiter/internal/limiter"
	"github.com/seu-usuario/desafio-rate-limiter/internal/middleware"
	"github.com/seu-usuario/desafio-rate-limiter/internal/strategy"
)

type Response struct {
	Message string `json:"message"`
	Time    string `json:"timestamp"`
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Redis storage strategy
	storage, err := strategy.NewRedisStorage(cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to initialize Redis storage: %v", err)
	}
	defer storage.Close()

	// Create rate limiter
	rateLimiter := limiter.NewRateLimiter(
		storage,
		cfg.RateLimitIP,
		cfg.IPBlockDuration,
	)
	defer rateLimiter.Close()

	// Configure tokens (example tokens)
	rateLimiter.ConfigureToken("token123", cfg.RateLimitToken, cfg.TokenBlockDuration)
	rateLimiter.ConfigureToken("premium-token", 100, 60)

	// Create HTTP server
	mux := http.NewServeMux()

	// Health check endpoint (without rate limiting)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Protected endpoints with rate limiting middleware
	rateLimitedMux := http.NewServeMux()
	rateLimitedMux.HandleFunc("/", handleRequest)
	rateLimitedMux.HandleFunc("/api/test", handleTestRequest)

	// Apply middleware to protected endpoints
	rateLimitedHandler := middleware.RateLimiterMiddleware(rateLimiter)(rateLimitedMux)

	// Combine both handlers
	mux.Handle("/api/", rateLimitedHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("✓ Rate Limiter started successfully!")
		log.Printf("✓ Server listening on port :%d", cfg.ServerPort)
		log.Printf("✓ Rate Limit (IP): %d requests/sec", cfg.RateLimitIP)
		log.Printf("✓ Block Duration (IP): %d seconds", cfg.IPBlockDuration)
		log.Printf("✓ Rate Limit (Token): %d requests/sec", cfg.RateLimitToken)
		log.Printf("✓ Block Duration (Token): %d seconds", cfg.TokenBlockDuration)
		log.Printf("✓ Redis: %s:%d", cfg.RedisHost, cfg.RedisPort)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("✓ Server stopped gracefully")
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := Response{
		Message:   "Request accepted! Rate limiter working correctly.",
		Time:      time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

func handleTestRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message": "Test endpoint response",
		"method":  r.Method,
		"path":    r.URL.Path,
		"ip":      r.RemoteAddr,
		"time":    time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}
