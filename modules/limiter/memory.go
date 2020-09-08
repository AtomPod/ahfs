package limiter

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type MemoryLimiter struct {
	cache *cache.Cache
}

func NewMemoryLimiter() Limiter {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &MemoryLimiter{
		cache: c,
	}
}

func (l *MemoryLimiter) Request(key string, count int, d time.Duration) (int, error) {
	c := l.cache

	result, err := c.IncrementInt(key, count)
	if err != nil {
		c.Set(key, count, d)
		return count, nil
	}
	return result, nil
}
