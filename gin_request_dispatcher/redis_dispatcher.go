package gin_request_dispatcher

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/redis/go-redis/v9"
	"net/http"
	"sync/atomic"
)

type RedisDispatcher struct {
	client *redis.Client
	reqMap cmap.ConcurrentMap[string, interface{}] // store real gin context
	key    atomic.Uint64                           // format redis list val, cause redis is not appropriate to save big val
}

const (
	RedisDispatcherListKey = "dispatcher_router" // key of redis queue dispatcher
)

func NewRedisDispatcher() *RedisDispatcher {
	return &RedisDispatcher{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "root",
			DB:       0,
		}),
		reqMap: cmap.New[interface{}](),
	}
}

func (r *RedisDispatcher) StoreReq(g *gin.Context) {
	key := fmt.Sprintf("req_seq_%d", r.key.Load())
	r.reqMap.Set(key, g)
	r.key.Add(1)
	r.client.LPush(g, RedisDispatcherListKey, key) // push request to list
}

func (r *RedisDispatcher) DoReq() {
	strCmd := r.client.RPop(context.Background(), RedisDispatcherListKey) // pop request to list
	val, flag := r.reqMap.Get(strCmd.Val())
	if flag {
		if gCtx, ok := val.(*gin.Context); ok {
			_, err := http.DefaultClient.Do(gCtx.Request)
			if err != nil {
				return
			}
		}
	}
}

// Process example to use this dispatcher
func Process(g *gin.Context) {
	r := NewRedisDispatcher()
	for i := 0; i < 10000; i++ {
		go func() {
			r.StoreReq(g)
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			r.DoReq()
		}()
	}

	d := NewDefaultDispatcher()
	for i := 0; i < 10000; i++ {
		go func() {
			d.StoreReq(g)
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			d.DoReq()
		}()
	}
}
