package setting

import (
	"time"

	"github.com/czhj/ahfs/modules/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Cache struct {
	Enabled  bool
	Adapter  string
	Url      string
	Interval time.Duration
	TTL      time.Duration
}

var (
	CacheService = struct {
		Cache
	}{
		Cache: Cache{
			Enabled:  true,
			Adapter:  "memory",
			Interval: 60 * time.Second,
			TTL:      24 * time.Hour,
		},
	}
)

func newCacheService() {
	viper.SetDefault("cache", map[string]interface{}{
		"enabled":  true,
		"adapter":  "memory",
		"interval": 60 * time.Second,
		"ttl":      24 * time.Hour,
	})

	cacheCfg := viper.Sub("cache")
	if err := cacheCfg.Unmarshal(&CacheService.Cache); err != nil {
		log.Fatal("Cannot unmarshal cache config", zap.Error(err))
	}

	adapter := CacheService.Adapter
	switch adapter {
	case "memory":
	case "redis", "memcache":
	case "":
		CacheService.Enabled = false
	default:
		log.Fatal("Unknow cache adapter", zap.String("adapter", adapter))
	}

	if CacheService.Enabled {
		log.Info("Cache Service Enabled")
	}
}
