package middleware

import (
	"bytes"
	"io"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/ninehills/go-webapp-template/pkg/logger"
)

type AuditMiddleware struct {
	l logger.Logger
}

func NewAuditMiddleware(l logger.Logger) *AuditMiddleware {
	return &AuditMiddleware{
		l: l,
	}
}

// 返回审计中间件.
func (a *AuditMiddleware) Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()

			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		c.Next()
		a.l.Info(
			"AUDIT_LOG",
			map[string]interface{}{
				"log_type":    "audit",
				"method":      c.Request.Method,
				"remote_ip":   c.ClientIP(),
				"host":        c.Request.Host,
				"path":        c.Request.URL.Path,
				"body":        string(body),
				"status_code": c.Writer.Status(),
				"user_agent":  c.Request.UserAgent(),
				"raw_query":   c.Request.URL.RawQuery,
				"request_id":  requestid.Get(c),
			},
		)
	}
}
