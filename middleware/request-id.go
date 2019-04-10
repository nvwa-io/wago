package middleware

import (
	"github.com/nvwa-io/wago"
	"github.com/nvwa-io/wago/util/uuid"
)

// register request id to  context
func RequestId() wago.MiddleWareHandler {
	return func(c *wago.Context) {
		c.Set(wago.REQUEST_ID, uuid.New())
		c.Next()
	}
}
