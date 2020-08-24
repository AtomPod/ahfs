package locker

import (
	"context"
	"time"
)

type Options struct {
	expiration  time.Duration
	maxAttempts int
}

type Option func(*Options)

type Locker interface {
	Init(config string) error
	Close() error

	Lock(ctx context.Context, key string, opts ...Option) (string, error)
	Unlock(ctx context.Context, key string, id string) error
}

func WithExpiration(d time.Duration) Option {
	return func(o *Options) {
		o.expiration = d
	}
}

func WithMaxAttempts(count int) Option {
	return func(o *Options) {
		o.maxAttempts = count
	}
}
