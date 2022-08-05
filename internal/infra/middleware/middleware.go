package middleware

import (
	"io"
	"time"

	ginlog "github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ninehills/go-webapp-template/internal/infra/dependency"
	"github.com/ninehills/go-webapp-template/pkg/logger"
	"github.com/rs/zerolog"
)

const requestIDKey = "X-Request-Id"

// 非全局的中间件集合
type Middlewares struct {
	Audit *AuditMiddleware
}

// 创建非全局的中间件
func NewMiddlewares(deps *dependency.Dependency) *Middlewares {
	auditMiddleware := NewAuditMiddleware(deps.Logger)
	return &Middlewares{
		Audit: auditMiddleware,
	}
}

// 注册全局中间件
func RegisterGlobalMiddleware(handler *gin.Engine, l *logger.Logger) {
	// Register middleware
	handler.Use(
		// request id middleware
		requestid.New(
			requestid.WithGenerator(func() string {
				// Generate a random UUID
				u, _ := uuid.NewRandom()
				return u.String()
			}),
			requestid.WithCustomHeaderStrKey(requestIDKey),
		),
		// logger middleware， 将访问日志也按照规范打到日志中。
		// TODO: 去除掉 Gin 自身打出的 access log.
		ginlog.SetLogger(ginlog.WithLogger(func(c *gin.Context, out io.Writer, latency time.Duration) zerolog.Logger {
			return l.Log().With().
				Str("request_id", requestid.Get(c)).
				Str("remote_ip", c.ClientIP()).
				Str("host", c.Request.Host).
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Str("user_agent", c.Request.UserAgent()).
				Int("status_code", c.Writer.Status()).
				Int("size", c.Writer.Size()).
				Dur("latency_ms", latency).
				Logger()
		})),
	)
}
