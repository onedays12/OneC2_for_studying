package middlewares

import "github.com/gin-gonic/gin"

func ResponseHeaderMiddleware(headers map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for k, v := range headers {
			c.Header(k, v)
		}
		c.Next()
	}
}
