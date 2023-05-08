package fsm

import (
	"context"
	"errors"
	"fmt"
)

// FSM -> finite state machine
type FSM interface {
	// Do make specified business types implementing this function to do their own features while transferring fromState to toState
	Do(ctx context.Context, fromState, toState string, param interface{}) error
}

type Impl struct {
	transition map[string][]tranAction // key -> fromState
}

type tranAction struct {
	toState string
	action  func(ctx context.Context, param interface{}) error
}

const (
	Init       = "init"
	Running    = "running"
	WaitFinish = "wait_finish"
	Success    = "success"
	Failed     = "failed"
)

func initTransition() map[string][]tranAction {
	return map[string][]tranAction{ // transition example
		Init: {
			{
				toState: Running,
				action: func(ctx context.Context, param interface{}) error {
					if param == nil {
						return errors.New("param is nil")
					}
					// todo implement this
					return nil
				},
			},
		},
		Running: {
			{
				toState: WaitFinish,
				action: func(ctx context.Context, param interface{}) error {
					if param == nil {
						return errors.New("param is nil")
					}
					// todo implement this
					return nil
				},
			},
			{
				toState: Failed,
				action: func(ctx context.Context, param interface{}) error {
					if param == nil {
						return errors.New("param is nil")
					}
					// todo implement this
					return nil
				},
			},
		},
		WaitFinish: {
			{
				toState: Failed,
				action: func(ctx context.Context, param interface{}) error {
					if param == nil {
						return errors.New("param is nil")
					}
					// todo implement this
					return nil
				},
			},
			{
				toState: Success,
				action: func(ctx context.Context, param interface{}) error {
					if param == nil {
						return errors.New("param is nil")
					}
					// todo implement this
					return nil
				},
			},
		},
	}
}

func NewFSM() *Impl {
	return &Impl{
		transition: initTransition(),
	}
}

func (f *Impl) Do(ctx context.Context, fromState, toState string, param interface{}) error {
	if _, ok := f.transition[fromState]; !ok {
		return errors.New("nonexistent from_state")
	}
	for _, action := range f.transition[fromState] {
		if action.toState == toState {
			if err := action.action(ctx, param); err != nil {
				return err
			} else {
				return nil
			}
		}
	}
	return errors.New("nonexistent from_state to to_state action")
}

func Example() {
	fsm, ctx := NewFSM(), context.Background()
	err := fsm.Do(ctx, Init, Running, "this is param")
	if err != nil {
		fmt.Println(err)
	}
	err = fsm.Do(ctx, Init, Running, "this is param")
	if err != nil {
		fmt.Println(err)
	}
}
