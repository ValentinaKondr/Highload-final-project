package cache

import "time"

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Close() error
}
