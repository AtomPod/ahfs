package locker

import (
	"context"
	"sync"

	"github.com/rs/xid"
)

type osLocker struct {
	id     string
	closed chan bool
}

type OSLocker struct {
	m       sync.Mutex
	lockers map[string]*osLocker
}

func (o *OSLocker) Init(config string) error {
	return nil
}

func (o *OSLocker) Close() error {
	return nil
}

func (o *OSLocker) Lock(ctx context.Context, key string, opts ...Option) (string, error) {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}

	var attempts = 0
	for {
		o.m.Lock()
		locker, ok := o.lockers[key]
		if ok {
			o.m.Unlock()

			attempts++
			if options.maxAttempts > 0 && options.maxAttempts <= attempts {
				return "", ErrMaxAttempts
			}

			select {
			case <-locker.closed:
			case <-ctx.Done():
				return "", ctx.Err()
			}
		} else {
			id := xid.New().String()
			o.lockers[key] = &osLocker{
				id:     id,
				closed: make(chan bool),
			}
			o.m.Unlock()
			return id, nil
		}
	}
}

func (o *OSLocker) Unlock(ctx context.Context, key string, id string) error {
	o.m.Lock()
	defer o.m.Unlock()

	locker, ok := o.lockers[key]
	if !ok || locker.id != id {
		return ErrLockerNotFound
	}

	delete(o.lockers, key)
	close(locker.closed)

	return nil
}
