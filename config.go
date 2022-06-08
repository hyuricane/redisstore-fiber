package redisstorage

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// Config defines the config for storage.
type Config struct {
	Client     *redis.Client
	Expiration time.Duration
	Prefix     string
	Secret     string
	Reset      bool
}
