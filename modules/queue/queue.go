package queue

import (
	"context"
	"fmt"
)

type ErrInvalidConfiguration struct {
	cfg interface{}
}

func (e ErrInvalidConfiguration) Error() string {
	return fmt.Sprintf("Invalid configuration: %v", e.cfg)
}

func IsErrInvalidConfiguration(e error) bool {
	_, ok := e.(ErrInvalidConfiguration)
	return ok
}

type Type string

type Data interface{}

type HandleFunc func(data ...Data)

type NewQueueFunc func(handleFunc HandleFunc, opts, exemplar interface{}) (Queue, error)

type Named interface {
	Name() string
}

type Flushable interface {
	Flush(context.Context) error
}

type Queue interface {
	Flushable
	Run(atShutdown, atTerminate func(context.Context, func()))
	Push(Data) error
}

var queueMap = map[Type]NewQueueFunc{}

func NewQueue(queueType Type, handleFunc HandleFunc, opts, exemplar interface{}) (Queue, error) {
	newFn, ok := queueMap[queueType]
	if !ok {
		return nil, fmt.Errorf("Unsupported queue type: %v", queueType)
	}
	return newFn(handleFunc, opts, exemplar)
}

func RegisterAsType(t Type, fn NewQueueFunc) {
	if _, ok := queueMap[t]; ok {
		panic("RegisterAsType: cannot resgiter twice for queue: " + t)
	}

	queueMap[t] = fn
}

func ResgiteredTypes() []Type {
	types := make([]Type, 0, len(queueMap))

	for t := range queueMap {
		types = append(types, t)
	}
	return types
}
