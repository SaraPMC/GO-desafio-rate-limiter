package strategy

import "context"

// StorageStrategy defines the interface for rate limiter storage strategies
type StorageStrategy interface {
	// IncrementCounter increments the request counter for a key
	IncrementCounter(ctx context.Context, key string) (int, error)

	// SetExpiration sets expiration time for a key
	SetExpiration(ctx context.Context, key string, ttlSeconds int) error

	// GetCounter gets the current counter value for a key
	GetCounter(ctx context.Context, key string) (int, error)

	// Exists checks if a key exists
	Exists(ctx context.Context, key string) (bool, error)

	// Delete removes a key
	Delete(ctx context.Context, key string) error

	// Close closes the connection to the storage
	Close() error
}
