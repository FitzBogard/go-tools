package err_group

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
)

type Group struct {
	wg   sync.WaitGroup
	err  error
	once sync.Once
}

func NewErrorGroup() *Group {
	return &Group{
		wg:  sync.WaitGroup{},
		err: nil,
	}
}

func (e *Group) Go(ctx context.Context, fn func() error) {
	e.wg.Add(1)
	go func(ctx context.Context) {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
				e.err = fmt.Errorf("panic: %+v", r)
			}
		}()
		defer e.wg.Done()
		if err := fn(); err != nil {
			e.once.Do(func() {
				e.err = err
			})
		}
	}(ctx)
}

func (e *Group) Wait() error {
	e.wg.Wait()
	return e.err
}
