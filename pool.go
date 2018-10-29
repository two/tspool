package tspool

import (
	"errors"
	"runtime"
	"sync"
)

var (
	errWorker = errors.New("all workers is busying")
	errCap    = errors.New("initCap must <= maxCap")
	errArgs   = errors.New("invalied args num, need most two args")
)

var defaultWorkerPoolCap = uint(runtime.NumCPU() * 100)

// WorkerPool provide Get and Put functions of pool
type WorkerPool interface {
	Get() (Worker, error)
	Put(Worker) error
}

type defaultWorkerPool struct {
	initCap uint // initial cap of pool
	maxCap  uint // max cap of pool
	mu      sync.Mutex
	workers []Worker // worker pool
	used    uint     // used worker of pool
}

// Get one idle worker from pool
func (wp *defaultWorkerPool) Get() (Worker, error) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if wp.used >= wp.maxCap {
		wp.used = wp.maxCap
		return nil, errWorker
	}
	if wp.workers[wp.used] == nil {
		defaultWorker := newWorker(wp.used)
		defaultWorker.run()
		wp.workers[wp.used] = defaultWorker
	}
	worker := wp.workers[wp.used]
	wp.used++
	return worker, nil
}

// Put idle worker to pool
func (wp *defaultWorkerPool) Put(worker Worker) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	w := worker.(*defaultWorker)
	wp.used--
	//change used work with idle worker
	wp.workers[w.pos] = wp.workers[wp.used]
	wp.workers[w.pos].(*defaultWorker).pos = w.pos
	w.pos = wp.used
	wp.workers[wp.used] = w
	return nil
}

func (wp *defaultWorkerPool) initWorkerPool() (WorkerPool, error) {
	wp.workers = make([]Worker, wp.maxCap)
	for i := uint(0); i < wp.initCap; i++ {
		wp.workers[i] = newWorker(i)
	}
	return wp, nil
}

func (wp *defaultWorkerPool) initWorker() {
	for i := uint(0); i < wp.initCap; i++ {
		wp.workers[i].(*defaultWorker).run()
	}
}

// DefaultWorkerPool return the default config of pool
func DefaultWorkerPool(cap ...uint) (WorkerPool, error) {
	wp := &defaultWorkerPool{}
	switch len(cap) {
	case 0:
		wp.initCap = defaultWorkerPoolCap
		wp.maxCap = wp.initCap
	case 1:
		wp.initCap = uint(cap[0])
		wp.maxCap = wp.initCap
	case 2:
		wp.initCap = uint(cap[0])
		wp.maxCap = uint(cap[1])
		if wp.maxCap < wp.initCap {
			return nil, errCap
		}
	default:
		return nil, errArgs
	}
	wp.initWorkerPool()
	wp.initWorker()
	return wp, nil
}
