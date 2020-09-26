package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/czhj/ahfs/modules/utils"
)

type WorkerPoolConfiguration struct {
	BlockTimeout time.Duration
	BoostTimeout time.Duration
	BoostWorker  int
	MaxWorkers   int
	QueueLength  int
	BatchLength  int
}

type WorkerPool struct {
	m                  sync.Mutex
	baseCtx            context.Context
	cancel             context.CancelFunc
	qid                int64
	maxNumberOfWorkers int
	numberOfWorkers    int
	batchLength        int
	blockTimeout       time.Duration
	boostTimeout       time.Duration
	boostWorkers       int
	numInQueue         int64
	handle             HandleFunc
	dataChan           chan Data
}

func NewWorkerPool(handle HandleFunc, config WorkerPoolConfiguration) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	dataChan := make(chan Data, config.QueueLength)
	pool := &WorkerPool{
		baseCtx:            ctx,
		cancel:             cancel,
		maxNumberOfWorkers: config.MaxWorkers,
		batchLength:        config.BatchLength,
		blockTimeout:       config.BlockTimeout,
		boostTimeout:       config.BoostTimeout,
		boostWorkers:       config.BoostWorker,
		handle:             handle,
		dataChan:           dataChan,
	}

	return pool
}

func (p *WorkerPool) Push(data Data) {
	atomic.AddInt64(&p.numInQueue, 1)
	p.m.Lock()
	if p.blockTimeout > 0 && p.boostTimeout > 0 && (p.numberOfWorkers <= p.maxNumberOfWorkers || p.maxNumberOfWorkers < 0) {
		p.m.Unlock()
		p.pushBoost(data)
	} else {
		p.m.Unlock()
		p.dataChan <- data
	}
}

func (p *WorkerPool) pushBoost(data Data) {
	select {
	case p.dataChan <- data:
	default:
		p.m.Lock()
		if p.blockTimeout <= 0 {
			p.m.Unlock()
			p.dataChan <- data
			return
		}

		ourTimeout := p.blockTimeout
		timer := time.NewTimer(ourTimeout)
		p.m.Unlock()

		select {
		case p.dataChan <- data:
			utils.StopTimer(timer)
		case <-timer.C:
			p.m.Lock()
			if p.blockTimeout > ourTimeout || (p.numberOfWorkers >= p.maxNumberOfWorkers && p.maxNumberOfWorkers > 0) {
				p.m.Unlock()
				p.dataChan <- data
				return
			}

			p.blockTimeout *= 2
			ctx, cancel := context.WithCancel(p.baseCtx)
			mq := GetManager().GetManagedQueue(p.qid)
			boost := p.boostWorkers

			if (boost + p.numberOfWorkers) > p.maxNumberOfWorkers {
				boost = p.maxNumberOfWorkers - p.numberOfWorkers
			}

			if mq != nil {
				start := time.Now()
				end := start.Add(p.boostTimeout)
				pid := mq.RegisterWorker(boost, cancel, start, end, false)

				go func() {
					<-ctx.Done()
					mq.RemoveWorker(pid)
					cancel()
				}()
			}

			go func() {
				<-time.After(p.boostTimeout)
				cancel()

				p.m.Lock()
				p.boostTimeout /= 2
				p.m.Unlock()
			}()

			p.m.Unlock()
			p.addWorkers(ctx, boost)
			p.dataChan <- data
		}
	}
}

func (p *WorkerPool) commonRegisterWorker(count int, timeout time.Duration) context.CancelFunc {
	var start time.Time
	var end time.Time
	var ctx context.Context
	var cancel context.CancelFunc
	if timeout > 0 {
		start = time.Now()
		end = start.Add(timeout)
		ctx, cancel = context.WithTimeout(p.baseCtx, timeout)
	} else {
		ctx, cancel = context.WithCancel(p.baseCtx)
	}

	mq := GetManager().GetManagedQueue(p.qid)
	if mq != nil {
		pid := mq.RegisterWorker(count, cancel, start, end, false)

		go func() {
			<-ctx.Done()
			mq.RemoveWorker(pid)
			cancel()
		}()
	}

	p.addWorkers(ctx, count)
	return cancel
}

func (p *WorkerPool) addWorkers(ctx context.Context, number int) {
	for i := 0; i < number; i++ {
		p.m.Lock()
		p.numberOfWorkers++
		p.m.Unlock()

		go func() {
			p.doWork(ctx)

			p.m.Lock()
			p.numberOfWorkers--
			p.m.Unlock()
		}()
	}
}

func (p *WorkerPool) doWork(ctx context.Context) {
	delay := time.Millisecond * 300
	data := make([]Data, 0, p.batchLength)

	for {
		select {
		case <-ctx.Done():
			if len(data) > 0 {
				p.handle(data...)
				atomic.AddInt64(&p.numInQueue, int64(-1*len(data)))
			}
			return
		case dat, ok := <-p.dataChan:
			if !ok {
				if len(data) > 0 {
					p.handle(data...)
					atomic.AddInt64(&p.numInQueue, int64(-1*len(data)))
				}
				return
			}
			data = append(data, dat)
			if len(data) >= p.batchLength {
				p.handle(data...)
				atomic.AddInt64(&p.numInQueue, int64(-1*len(data)))
				data = make([]Data, 0, p.batchLength)
			}
		default:
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				utils.StopTimer(timer)
				if len(data) > 0 {
					p.handle(data...)
					atomic.AddInt64(&p.numInQueue, int64(-1*len(data)))
				}
				return
			case dat, ok := <-p.dataChan:
				utils.StopTimer(timer)
				if !ok {
					if len(data) > 0 {
						p.handle(data...)
						atomic.AddInt64(&p.numInQueue, int64(-1*len(data)))
					}
					return
				}

				data = append(data, dat)
				if len(data) >= p.batchLength {
					p.handle(data...)
					atomic.AddInt64(&p.numInQueue, int64(-1*len(data)))
					data = make([]Data, 0, p.batchLength)
				}
			case <-timer.C:
				delay = time.Millisecond * 100
				if len(data) > 0 {
					p.handle(data...)
					atomic.AddInt64(&p.numInQueue, int64(-1*len(data)))
					data = make([]Data, 0, p.batchLength)
				}
			}
		}
	}
}

func (p *WorkerPool) BoostWorkers() int {
	p.m.Lock()
	defer p.m.Unlock()

	return p.boostWorkers
}

func (p *WorkerPool) BoostTimeout() time.Duration {
	p.m.Lock()
	defer p.m.Unlock()

	return p.boostTimeout
}

func (p *WorkerPool) BlockTimeout() time.Duration {
	p.m.Lock()
	defer p.m.Unlock()

	return p.blockTimeout
}

func (p *WorkerPool) MaxNumberOfWorkers() int {
	p.m.Lock()
	defer p.m.Unlock()

	return p.maxNumberOfWorkers
}

func (p *WorkerPool) NumberOfWorkers() int {
	p.m.Lock()
	defer p.m.Unlock()

	return p.numberOfWorkers
}

func (p *WorkerPool) SetMaxNumberOfWorkers(count int) {
	p.m.Lock()
	defer p.m.Unlock()

	p.maxNumberOfWorkers = count
}

func (p *WorkerPool) SetManagedSettings(boostNumber int, maxWorkerNumber int, timeout time.Duration) {
	p.m.Lock()
	defer p.m.Unlock()

	p.boostTimeout = timeout
	p.maxNumberOfWorkers = maxWorkerNumber
	p.boostWorkers = boostNumber
}

func (p *WorkerPool) AddWorker(count int, timeout time.Duration) context.CancelFunc {
	return p.commonRegisterWorker(count, timeout)
}
