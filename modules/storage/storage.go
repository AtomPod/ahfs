package storage

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/czhj/ahfs/modules/log"
	"github.com/czhj/ahfs/modules/setting"
	"go.uber.org/zap"
)

var (
	ErrNotFound = errors.New("not found")
)

type ErrInvalidConfiguration struct {
	config interface{}
	err    error
}

func (e ErrInvalidConfiguration) Error() string {
	if e.err != nil {
		return fmt.Sprintf("invalid configuration: %v, error: %v", e.config, e.err)
	}

	return fmt.Sprintf("invalid configuration: %v", e.config)
}

func IsErrInvalidConfiguration(err error) bool {
	if _, ok := err.(ErrInvalidConfiguration); ok {
		return true
	}
	return false
}

type Type string
type Generator func(context.Context, interface{}) (Storage, error)

var (
	storageMap map[Type]Generator = make(map[Type]Generator)
)

type ID string

type WriteOptions struct {
	ID uint
}

type WriteOption func(*WriteOptions)

type ReadOptions struct {
}

type ReadOption func(*ReadOptions)

type Object struct {
	Size   int64
	Reader io.ReadCloser
}

type Storage interface {
	Write(f *Object, opts ...WriteOption) (ID, error)
	Read(id ID, opts ...ReadOption) (*Object, error)
	Delete(id ID) error
}

func WithID(id uint) WriteOption {
	return func(wo *WriteOptions) {
		wo.ID = id
	}
}

var (
	LFS Storage
)

func NewStorage(typ string, cfg interface{}) (Storage, error) {
	if len(typ) == 0 {
		typ = "local"
	}

	generator, ok := storageMap[Type(typ)]
	if !ok {
		return nil, fmt.Errorf("Unsupported storage type: %v", typ)
	}

	return generator(context.Background(), cfg)
}

func Init() error {
	return initLFS()
}

func initLFS() (err error) {
	log.Info("Initialising LFS storage", zap.String("type", setting.LFS.Storage.Type))
	LFS, err = NewStorage(setting.LFS.Storage.Type, &setting.LFS.Storage)
	return err
}
