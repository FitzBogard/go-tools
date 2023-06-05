package task_scheduler

import (
	"context"
	"time"
)

type Worker struct {
	ch chan Task
}

func New() *Worker {
	return &Worker{ch: make(chan Task, 1)} // 每个worker只处理1个任务
}

type Task struct {
	fn    func(ctx context.Context, param interface{}) error
	param interface{}
}

func (t Task) Exec(ctx context.Context) error {
	return t.fn(ctx, t.param)
}

func (w *Worker) Work(ctx context.Context) error {
	select {
	case t := <-w.ch:
		return t.Exec(ctx)
	default:
		time.Sleep(time.Millisecond * 300) // 等待300ms
	}
	return nil
}
