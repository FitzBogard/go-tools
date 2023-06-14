package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type (
	Worker struct {
		processors []chan Task // 管道切片长度对应 worker 数量
		workerNum  int
	}
	Task struct {
		id    int
		fn    ExecFunc
		param interface{}
	}
	ExecFunc func(ctx context.Context, param interface{}) error
	Option   func(*Worker)
)

func New(options ...Option) *Worker {
	w := &Worker{workerNum: 1, processors: make([]chan Task, 1)}
	for _, opt := range options {
		opt(w)
	}
	return w
}

func SetWorkerNum(workerNum int) Option {
	if workerNum <= 1 {
		return func(w *Worker) {
			// default num
		}
	}
	return func(w *Worker) {
		w.workerNum = workerNum
		w.processors = make([]chan Task, workerNum)
	}
}

func (t Task) Exec(ctx context.Context) error {
	return t.fn(ctx, t.param)
}

func (w *Worker) Run(ctx context.Context) {
	for i := range w.processors {
		w.processors[i] = make(chan Task, 1)
		go process(ctx, w.processors[i], func(task Task, err error) {
			fmt.Println(fmt.Sprintf("task:%v, err:%v", task.param, err.Error()))
		})
	}
}

func process(ctx context.Context, t chan Task, errFunc func(task Task, err error)) {
	interrupt := make(chan os.Signal, 1)
	select {
	case task := <-t:
		if err := task.Exec(ctx); err != nil {
			errFunc(task, err)
		}
	case <-interrupt:
		time.Sleep(time.Second * 3)
		return
	}
}

func main() {
	w := New(SetWorkerNum(5))
	w.Run(context.Background())
	go func() {
		for {
			t := Task{
				id:    rand.Int()>>2 ^ 31,
				param: make(map[string]interface{}),
				fn: func(ctx context.Context, param interface{}) error {
					fmt.Println(param)
					return nil
				},
			}
			w.processors[t.id%w.workerNum] <- t // 任务根据 id 进入对应取模的管道分片去处理
		}
	}()
	time.Sleep(time.Second * 3)
}
