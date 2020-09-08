package cache

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/czhj/ahfs/modules/log"
	"github.com/czhj/ahfs/modules/setting"
	"github.com/gin-contrib/cache/persistence"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var (
	ErrNotFound = errors.New("not found")
)

var (
	cacheStore persistence.CacheStore
)

func NewContext() {
	if cacheStore == nil && setting.CacheService.Enabled {
		if err := newCache(setting.CacheService.Cache); err != nil {
			log.Fatal("Failed to initialize cache", zap.Error(err))
		}
	}
}

func newCache(config setting.Cache) error {
	switch config.Adapter {
	case "redis":
		options, err := redis.ParseURL(config.Url)
		if err != nil {
			return err
		}
		cacheStore = persistence.NewRedisCache(options.Addr, options.Password, config.TTL)
	case "memcache":
		memUrl, err := url.Parse(config.Url)
		if err != nil {
			return err
		}
		username := memUrl.User.Username()
		password, _ := memUrl.User.Password()
		cacheStore = persistence.NewMemcachedBinaryStore(memUrl.Host, username, password, config.TTL)
	case "memory":
		cacheStore = persistence.NewInMemoryStore(config.TTL)
	default:
		return fmt.Errorf("Unknow adapter [%s]", config.Adapter)
	}
	return nil
}

func Set(key string, val interface{}, d time.Duration) error {
	return cacheStore.Set(key, val, d)
}

func Get(key string, val interface{}) error {
	err := cacheStore.Get(key, val)
	if err == persistence.ErrCacheMiss {
		return ErrNotFound
	}
	return err
}

func Delete(key string) error {
	err := cacheStore.Delete(key)
	if err == persistence.ErrCacheMiss {
		return ErrNotFound
	}
	return err
}

func Increment(key string, count uint64) (uint64, error) {
	result, err := cacheStore.Increment(key, count)
	if err != nil {
		if err == persistence.ErrCacheMiss {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func Decrement(key string, count uint64) (uint64, error) {
	result, err := cacheStore.Decrement(key, count)
	if err != nil {
		if err == persistence.ErrCacheMiss {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return result, nil
}
