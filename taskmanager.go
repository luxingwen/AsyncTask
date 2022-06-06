package task

import (
	"context"

	"github.com/panjf2000/ants"
)

const (
	defaultSize = 100
)

type taskManger struct {
	taskQueue chan Task
	size      int
	worker    *ants.PoolWithFunc
}

func (self *taskManger) Size() int {
	return self.size
}

func NewTaskManager(size int, ctx context.Context) *taskManger {
	if size <= defaultSize {
		size = defaultSize
	}
	t := &taskManger{
		taskQueue: make(chan Task, size),
		size:      size,
	}

	var err error

	t.worker, err = ants.NewPoolWithFunc(size, t.Exec)
	if err != nil {
		panic(err)
	}
	go t.run(ctx)
	return t
}

func (self *taskManger) Enqueue(t Task) bool {
	select {
	case self.taskQueue <- t:
		return true
	default:
		return false
	}
}

func (self *taskManger) run(ctx context.Context) {
	for {
		select {
		case t := <-self.taskQueue:
			self.worker.Invoke(t)
		case <-ctx.Done():
			return
		}
	}
}

func (self *taskManger) Exec(req interface{}) {
	t, ok := req.(Task)
	if !ok {
		return
	}
	TaskRun(t)

	state := t.State()
	if state != TaskStateDone && state != TaskStateFailed {
		self.Enqueue(t)
	}
}
