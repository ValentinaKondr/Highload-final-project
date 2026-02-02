package cache

import "time"

type testCache struct{}

func newTestCache() *testCache {
	return &testCache{}
}

func (c *testCache) Set(key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (c *testCache) Close() error {
	return nil
}
