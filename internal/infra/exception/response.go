package exception

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/ninehills/go-webapp-template/apis/httpv1"
)

func CodeResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(
		code,
		httpv1.ErrorResponse{
			Message: msg, Code: http.StatusText(code), RequestID: requestid.Get(c),
		},
	)
}

func ResponseWithError(c *gin.Context, err error) {
	var e *Error
	if errors.As(err, &e) {
		CodeResponse(c, e.Code(), e.Error())
	} else {
		CodeResponse(c, http.StatusInternalServerError, err.Error())
	}
}
