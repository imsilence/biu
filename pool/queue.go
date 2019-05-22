package pool

import (
	"container/list"
	"sync"
)

type Queue struct {
	l     *list.List
	mutex sync.RWMutex
}

func NewQueue() *Queue {
	return &Queue{
		l: list.New(),
	}
}

func (q *Queue) Push(v interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.l.PushBack(v)
}

func (q *Queue) Len() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return q.l.Len()
}

func (q *Queue) Pop() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	element := q.l.Front()
	if element != nil {
		q.l.Remove(element)
		return element.Value
	}
	return nil
}
