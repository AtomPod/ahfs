package locker

import (
	"context"
	"errors"
	"reflect"

	"github.com/czhj/ahfs/modules/log"
	"go.uber.org/zap"
)

var (
	ErrKeyLocked      = errors.New("key is locked")
	ErrLockerIsExpire = errors.New("locker is expire")
	ErrLockerNotFound = errors.New("locker does not exists")
	ErrMaxAttempts    = errors.New("maximum number of attempts reached")
)

var (
	defaultLocker Locker = &OSLocker{
		lockers: make(map[string]*osLocker),
	}
)

func SetDefaultLocker(l Locker) {
	defaultLocker = l
}

func init() {
	if err := defaultLocker.Init(""); err != nil {
		log.Error("Cannot initialize locker", zap.String("name", reflect.TypeOf(defaultLocker).String()), zap.Error(err))
	}
}

func Lock(ctx context.Context, key string, opts ...Option) (string, error) {
	return defaultLocker.Lock(ctx, key, opts...)
}

func Unlock(ctx context.Context, key string, id string) error {
	return defaultLocker.Unlock(ctx, key, id)
}
