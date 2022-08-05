package middleware

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/ninehills/go-webapp-template/pkg/logger"
)

type AuditMiddleware struct {
	l *logger.Logger
}

func NewAuditMiddleware(l *logger.Logger) *AuditMiddleware {
	return &AuditMiddleware{
		l: l,
	}
}

// 返回审计中间件
func (a *AuditMiddleware) Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		c.Next()
		a.l.WithContext(c).Info().
			Str("log_type", "audit").
			Str("method", c.Request.Method).
			Str("remote_ip", c.ClientIP()).
			Str("path", c.Request.URL.Path).
			Str("user_agent", c.Request.UserAgent()).
			Int("status_code", c.Writer.Status()).
			Str("raw_query", c.Request.URL.RawQuery).
			Bytes("body", body).
			Msg("AUDIT_LOG")
	}
}
