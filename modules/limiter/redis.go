package limiter

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisLimiter struct {
	c      *redis.Client
	script *redis.Script
}

func NewRedisLimiter(opt *redis.Options) (Limiter, error) {
	limiter := &RedisLimiter{
		c:      redis.NewClient(opt),
		script: newRedisScript(),
	}

	_, err := limiter.c.Ping().Result()
	if err != nil {
		return nil, err
	}
	return limiter, nil
}

func newRedisScript() *redis.Script {
	return redis.NewScript(`
	if redis.call("GET", KEYS[1]) == false then
		if redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2]) == false then
			return -1
		end
		return ARGV[1]
	end
	return redis.call("INCRBY", KEYS[1], ARGV[1])`)
}

func (l *RedisLimiter) Request(key string, count int, d time.Duration) (int, error) {
	r := l.c
	script := l.script
	result := script.Eval(r, []string{key}, count, d)
	return result.Int()
}
