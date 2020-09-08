package limiter

import (
	"fmt"
	"log"
	"time"

	"github.com/czhj/ahfs/modules/setting"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var (
	limiter Limiter
)

func NewContext() {
	if limiter == nil && setting.CacheService.Enabled {
		if err := newLimiter(setting.CacheService.Cache); err != nil {
			log.Fatal("Failed to initialize limiter", zap.Error(err))
		}
	}
}

func newLimiter(config setting.Cache) error {
	switch config.Adapter {
	case "redis":
		options, err := redis.ParseURL(config.Url)
		if err != nil {
			return err
		}
		limiter, err = NewRedisLimiter(options)
		if err != nil {
			return err
		}
	case "memory":
		limiter = NewMemoryLimiter()
	default:
		return fmt.Errorf("Unknow adapter [%s]", config.Adapter)
	}
	return nil
}

func Request(key string, count int, d time.Duration) (int, error) {
	return limiter.Request(key, count, d)
}
