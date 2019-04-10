package middleware

import (
	"github.com/nvwa-io/wago"
	"github.com/nvwa-io/wago/logger"
)

// log before & after request handled
func RequestLogger() wago.MiddleWareHandler {
	return func(c *wago.Context) {
		logger.WithFields(logger.Fields{
			wago.REQUEST_ID: c.GetString(wago.REQUEST_ID),
			"path":          c.Request.URL.Path,
			"host":          c.Request.Host,
			"header":        c.Request.Header,
		}).Debug("before-handle")

		c.Next()

		logger.WithFields(logger.Fields{
			wago.REQUEST_ID: c.GetString(wago.REQUEST_ID),
			"path":          c.Request.URL.Path,
			"host":          c.Request.Host,
			"header":        c.Request.Header,
		}).Debug("after-handle")
	}
}
