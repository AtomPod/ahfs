package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type ManagedPool interface {
	BoostWorkers() int
	BoostTimeout() time.Duration
	BlockTimeout() time.Duration
	MaxNumberOfWorkers() int
	NumberOfWorkers() int
	AddWorker(count int, timeout time.Duration) context.CancelFunc
	SetMaxNumberOfWorkers(c int)
	SetManagedSettings(boostNumber int, maxWorkerNumber int, timeout time.Duration)
}

type PoolWorker struct {
	pid       int64
	cancel    context.CancelFunc
	count     int
	start     time.Time
	end       time.Time
	isFlusher bool
}

type ManagedQueue struct {
	m             sync.Mutex
	qid           int64
	counter       int64
	managed       interface{}
	workers       map[int64]*PoolWorker
	exemplar      string
	configuration interface{}
	t             Type
	name          string
}

func NewManagedQueue(qid int64, managed interface{}) *ManagedQueue {
	return &ManagedQueue{
		qid:     qid,
		managed: managed,
		workers: make(map[int64]*PoolWorker),
	}
}

func (q *ManagedQueue) RegisterWorker(count int, cancel context.CancelFunc,
	start time.Time, end time.Time, isFlusher bool) int64 {
	q.m.Lock()
	defer q.m.Unlock()

	q.counter++
	pid := q.counter

	q.workers[pid] = &PoolWorker{
		pid:       pid,
		count:     count,
		cancel:    cancel,
		start:     start,
		end:       end,
		isFlusher: isFlusher,
	}

	return pid
}

func (q *ManagedQueue) CancelWorker(qid int64) bool {
	q.m.Lock()
	worker, ok := q.workers[qid]
	q.m.Unlock()

	if !ok {
		return false
	}
	worker.cancel()
	return true
}

func (q *ManagedQueue) RemoveWorker(qid int64) {
	q.m.Lock()
	worker, ok := q.workers[qid]
	delete(q.workers, qid)
	q.m.Unlock()

	if ok && worker.cancel != nil {
		worker.cancel()
	}
}

func (q *ManagedQueue) WorkerNumber() int {
	q.m.Lock()
	defer q.m.Unlock()

	var count int
	for _, worker := range q.workers {
		count += worker.count
	}
	return count
}

func (q *ManagedQueue) BoostWorkers() int {
	if pool, ok := q.managed.(ManagedPool); ok {
		return pool.BoostWorkers()
	}
	return 0
}

func (q *ManagedQueue) MaxNumberOfWorkers() int {
	if pool, ok := q.managed.(ManagedPool); ok {
		return pool.MaxNumberOfWorkers()
	}
	return 0
}

func (q *ManagedQueue) NumberOfWorkers() int {
	if pool, ok := q.managed.(ManagedPool); ok {
		return pool.NumberOfWorkers()
	}
	return 0
}

func (q *ManagedQueue) BoostTimeout() time.Duration {
	if pool, ok := q.managed.(ManagedPool); ok {
		return pool.BoostTimeout()
	}
	return 0
}

func (q *ManagedQueue) BlockTimeout() time.Duration {
	if pool, ok := q.managed.(ManagedPool); ok {
		return pool.BlockTimeout()
	}
	return 0
}

type Manager struct {
	m       sync.Mutex
	counter int64
	queues  map[int64]*ManagedQueue
}

var (
	manager *Manager
)

func GetManager() *Manager {
	if manager == nil {
		manager = &Manager{
			queues: make(map[int64]*ManagedQueue),
		}
	}
	return manager
}

func (m *Manager) Add(managed interface{},
	t Type,
	configuration interface{},
	exemplar interface{}) int64 {

	cfg, _ := json.Marshal(configuration)

	mq := &ManagedQueue{
		managed:       managed,
		t:             t,
		configuration: string(cfg),
		exemplar:      reflect.TypeOf(exemplar).String(),
		workers:       make(map[int64]*PoolWorker),
	}

	m.m.Lock()
	m.counter++
	mq.qid = m.counter
	mq.name = fmt.Sprintf("queue-%d", mq.qid)
	if named, ok := managed.(Named); ok {
		mq.name = named.Name()
	}
	m.queues[mq.qid] = mq
	m.m.Unlock()
	return mq.qid
}

func (m *Manager) Remove(qid int64) {
	m.m.Lock()
	delete(m.queues, qid)
	m.m.Unlock()
}

func (m *Manager) GetManagedQueue(qid int64) *ManagedQueue {
	m.m.Lock()
	defer m.m.Unlock()
	return m.queues[qid]
}
