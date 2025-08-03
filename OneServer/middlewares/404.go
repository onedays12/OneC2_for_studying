package middlewares

import (
	"OneServer/profile"
	"github.com/gin-gonic/gin"
)

func Default404Middleware(config profile.ServerResponseConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Errors) > 0 && !c.Writer.Written() {
			for header, value := range config.Headers {
				c.Header(header, value)
			}
			c.String(config.Status, config.Page)
			c.Abort()
			return
		}

		c.Next()

		if len(c.Errors) > 0 && !c.Writer.Written() {
			for header, value := range config.Headers {
				c.Header(header, value)
			}
			c.String(config.Status, config.Page)
		}
	}
}
