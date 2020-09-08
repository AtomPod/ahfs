package limiter

import "time"

type Limiter interface {
	Request(key string, count int, d time.Duration) (int, error)
}
