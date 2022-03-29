package queue

import (
	"context"
	"fmt"

	"github.com/czhj/ahfs/modules/log"
	"go.uber.org/zap"
)

const ChannelQueueType Type = "channel"

type ChannelQueueConfiguration struct {
	WorkerPoolConfiguration
	Workers int
	Name    string
}

type ChannelQueue struct {
	*WorkerPool
	exemplar interface{}
	workers  int
	name     string
}

func NewChannelQueue(handle HandleFunc, cfg, exemplar interface{}) (Queue, error) {
	configInterface, err := toConfig(ChannelQueueConfiguration{}, cfg)
	if err != nil {
		return nil, err
	}

	config := configInterface.(ChannelQueueConfiguration)
	if config.BatchLength == 0 {
		config.BatchLength = 1
	}

	queue := &ChannelQueue{
		WorkerPool: NewWorkerPool(handle, config.WorkerPoolConfiguration),
		workers:    config.Workers,
		name:       config.Name,
		exemplar:   exemplar,
	}

	queue.qid = GetManager().Add(queue.WorkerPool, ChannelQueueType, config, exemplar)
	return queue, nil
}

func (q *ChannelQueue) Run(atShutdown, atTerminate func(context.Context, func())) {
	atShutdown(context.Background(), func() {
		log.Warn("ChannelQueue is not shutdownable!", zap.String("name", q.name))
	})

	atTerminate(context.Background(), func() {
		log.Warn("ChannelQueue is not terminate!", zap.String("name", q.name))
	})

	go func() {
		_ = q.AddWorker(q.workers, 0)
	}()
}

func (q *ChannelQueue) Push(d Data) error {
	if !assignableTo(d, q.exemplar) {
		return fmt.Errorf("Unable to assign data %s to type exemplar: %v in queue: %s", d, q.exemplar, q.name)
	}
	q.WorkerPool.Push(d)
	return nil
}

func (q *ChannelQueue) Flush(ctx context.Context) error {
	return nil
}

func init() {
	RegisterAsType(ChannelQueueType, NewChannelQueue)
}
