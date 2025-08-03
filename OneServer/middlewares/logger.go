package middlewares

import (
	"OneServer/logs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// GinLogger 是一个 Gin 中间件，用于记录请求日志。
// 该中间件会在每次请求结束后，使用 Zap 日志记录请求信息。
// 通过此中间件，可以方便地追踪每个请求的情况以及性能。
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 获取请求的路径和查询参数
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 继续执行后续的处理
		c.Next()

		// 计算请求处理的耗时
		cost := time.Since(start)

		// 使用 Zap 记录请求日志
		logs.Logger.Info(path,
			// 记录响应状态码
			zap.Int("status", c.Writer.Status()),
			// 记录请求方法
			zap.String("method", c.Request.Method),
			// 记录请求路径
			zap.String("path", path),
			// 记录查询参数
			zap.String("query", query),
			// 记录客户端 IP          zap.String("ip", c.ClientIP()),
			// 记录 User-Agent 信息
			zap.String("user-agent", c.Request.UserAgent()),
			// 记录错误信息（如果有）
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			// 记录请求耗时
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery 是一个 Gin 中间件，用于捕获和处理请求中的 panic 错误。
// 该中间件的主要作用是确保服务在遇到未处理的异常时不会崩溃，并通过日志系统提供详细的错误追踪。
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 defer 确保 panic 被捕获，并且处理函数会在 panic 后执行
		defer func() {
			// 检查是否发生了 panic 错误
			if err := recover(); err != nil {
				// 检查是否是连接被断开的问题（如 broken pipe），这些错误不需要记录堆栈信息
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取请求信息，包括请求体等
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// 如果是 broken pipe 错误，则只记录错误信息，不记录堆栈信息
				if brokenPipe {
					logs.Logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// 由于连接断开，不能再向客户端写入状态码
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()                // 中止请求处理
					return
				}

				// 如果是其他类型的 panic，根据 `stack` 参数决定是否记录堆栈信息
				if stack {
					// 记录详细的错误和堆栈信息
					logs.Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					// 只记录错误信息，不记录堆栈
					logs.Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				// 返回 500 错误状态码，表示服务器内部错误
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		// 继续执行后续的请求处理
		c.Next()
	}
}
