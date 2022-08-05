package exception

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// Error response for http
type Response struct {
	Message   string `json:"message" example:"message"`
	Code      string `json:"code" example:"Conflict"`
	RequestID string `json:"requestId" example:"b5953bf0-9f15-4c42-afb4-1c125b40d7ce"`
}

func CodeResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, Response{msg, http.StatusText(code), requestid.Get(c)})
}

func ErrorResponse(c *gin.Context, err error) {
	var e *Error
	if errors.As(err, &e) {
		CodeResponse(c, e.Code(), e.Error())
	} else {
		CodeResponse(c, http.StatusInternalServerError, err.Error())
	}
}
