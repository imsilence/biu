package pool

import (
	"math"

	"sync"

	"github.com/sirupsen/logrus"
)

type Pool struct {
	worker  int
	events  chan struct{}
	jobs    *Queue
	Results chan interface{}
	w       sync.WaitGroup
}

func New(worker int) *Pool {
	return &Pool{
		worker:  worker,
		events:  make(chan struct{}, math.MaxInt32),
		jobs:    NewQueue(),
		Results: make(chan interface{}, worker),
	}
}

func (pool *Pool) Add(task func() interface{}) {
	pool.events <- struct{}{}
	pool.jobs.Push(task)
}

func (pool *Pool) Start() {
	for id := 0; id <= pool.worker; id++ {
		pool.w.Add(1)
		go func(id int) {
			for _ = range pool.events {
				job := pool.jobs.Pop()
				if job != nil {
					logrus.WithFields(logrus.Fields{
						"id": id,
					}).Debug("running job")
					pool.Results <- job.(func() interface{})()
				}
			}
			pool.w.Done()
		}(id)
	}
}

func (pool *Pool) CloseAndWait() {
	close(pool.events)
	pool.w.Wait()
	close(pool.Results)
}
