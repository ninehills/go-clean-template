package http

import (
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message   string `json:"message" example:"message"`
	Code      string `json:"code" example:"Conflict"`
	RequestID string `json:"requestId" example:"b5953bf0-9f15-4c42-afb4-1c125b40d7ce"`
}

func ErrorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, Response{msg, http.StatusText(code), requestid.Get(c)})
}
