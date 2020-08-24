package session

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/czhj/ahfs/modules/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memcached"
	"github.com/gin-contrib/sessions/memstore"
	redisSession "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type Options struct {
	Provider       string
	ProviderConfig string
	HttpOnly       bool
	Secure         bool
	CookieName     string
	CookiePath     string
	MaxAge         int64
	Domain         string
	Secret         []byte
}

func NewSession(opts Options) gin.HandlerFunc {
	store := optionsToStore(opts)
	store.Options(sessions.Options{
		Path:     opts.CookiePath,
		Domain:   opts.Domain,
		MaxAge:   int(opts.MaxAge),
		HttpOnly: opts.HttpOnly,
		Secure:   opts.Secure,
	})
	return sessions.Sessions(opts.CookieName, store)
}

func optionsToStore(opts Options) sessions.Store {
	switch opts.Provider {
	case "redis":
		store, err := parseRedisConfig(opts.ProviderConfig, opts.Secret)
		if err != nil {
			log.Fatal("Cannot create redis session store", zap.Error(err))
		}
		return store
	case "memory":
		return memstore.NewStore(opts.Secret)
	case "memcached":
		return parseMemcacheConfig(opts.ProviderConfig, opts.Secret)
	}

	log.Fatal("Session provider is not supported", zap.String("provider", opts.Provider))
	return nil
}

func parseRedisConfig(config string, secret []byte) (redisSession.Store, error) {
	opts, err := redis.ParseURL(config)
	if err != nil {
		return nil, fmt.Errorf("ParseURL: %v", err)
	}

	store, err := redisSession.NewStoreWithDB(opts.PoolSize, opts.Network, opts.Addr, opts.Password, strconv.Itoa(opts.DB), []byte(secret))
	if err != nil {
		return nil, fmt.Errorf("NewStoreWithDB: %v", err)
	}
	return store, nil
}

func parseMemcacheConfig(config string, secret []byte) memcached.Store {
	servers := strings.Split(config, ";")
	return memcached.NewStore(memcache.New(servers...), "", secret)
}
