package request_dispatcher

import "github.com/gin-gonic/gin"

type Dispatcher interface {
	StoreReq(ctx *gin.Context)
	DoReq()
}
