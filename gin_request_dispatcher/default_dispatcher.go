package gin_request_dispatcher

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync/atomic"
	"unsafe"
)

type Node struct {
	val  interface{}
	next unsafe.Pointer
}

type Queue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func NewQueue() *Queue {
	n := unsafe.Pointer(&Node{})
	return &Queue{head: n, tail: n}
}

func (q *Queue) Push(val interface{}) {
	node := &Node{val: val}
	for {
		tail := q.tail
		next := atomic.LoadPointer(&(*Node)(tail).next)
		if next == nil {
			if atomic.CompareAndSwapPointer(&(*Node)(tail).next, next, unsafe.Pointer(node)) {
				atomic.CompareAndSwapPointer(&q.tail, tail, unsafe.Pointer(node))
				return
			}
		} else {
			atomic.CompareAndSwapPointer(&q.tail, tail, next)
		}
	}
}

func (q *Queue) Pop() interface{} {
	for {
		head := q.head
		tail := q.tail
		next := atomic.LoadPointer(&(*Node)(head).next)
		if head == q.head {
			if head == nil {
				if next == nil {
					return nil
				}
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			}
		} else {
			value := (*Node)(next).val
			if atomic.CompareAndSwapPointer(&q.head, head, next) {
				return value
			}
		}
	}
}

type DefaultDispatcher struct {
	queue *Queue
}

func NewDefaultDispatcher() *DefaultDispatcher {
	return &DefaultDispatcher{
		queue: NewQueue(),
	}
}

func (d *DefaultDispatcher) StoreReq(g *gin.Context) {
	d.queue.Push(g)
}
func (d *DefaultDispatcher) DoReq() {
	val := d.queue.Pop()
	if gCtx, ok := val.(*gin.Context); ok {
		_, err := http.DefaultClient.Do(gCtx.Request)
		if err != nil {
			return
		}
	}
}
