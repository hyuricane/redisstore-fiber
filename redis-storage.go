package redisstorage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Storage interface that is implemented by storage providers
type Storage struct {
	db         *redis.Client
	expiration time.Duration
	prefix     string
}

// New creates a new redis storage
func New(config ...Config) *Storage {
	// Set default config
	cfg := config[0]
	if cfg.Client == nil {
		cfg.Client = redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
		})
	}
	if cfg.Expiration == 0 {
		cfg.Expiration = time.Hour * 24 * 7
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "sess:"
	}

	db := cfg.Client

	// Test connection
	if err := db.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	// Empty collection if Clear is true
	if cfg.Reset {
		if err := db.FlushDB(context.Background()).Err(); err != nil {
			panic(err)
		}
	}

	// Create new store
	return &Storage{
		db:         db,
		prefix:     cfg.Prefix,
		expiration: cfg.Expiration,
	}
}

// Get value by key
func (s *Storage) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}
	val, err := s.db.Get(context.Background(), s.prefix+key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

// Set key with value
func (s *Storage) Set(key string, val []byte, exp time.Duration) error {
	// Ain't Nobody Got Time For That
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}
	return s.db.Set(context.Background(), s.prefix+key, val, exp).Err()
}

func (s *Storage) SetDef(key string, val []byte) error {
	return s.Set(s.prefix+key, val, s.expiration)
}

// Delete key by key
func (s *Storage) Delete(key string) error {
	// Ain't Nobody Got Time For That
	if len(key) <= 0 {
		return nil
	}
	return s.db.Del(context.Background(), s.prefix+key).Err()
}

// Reset all keys
func (s *Storage) Reset() error {
	return s.db.FlushDB(context.Background()).Err()
}

// Close the database
func (s *Storage) Close() error {
	return s.db.Close()
}
