package exception

import "net/http"

// 定义自定义错误.
type Error struct {
	// 错误码，列表在 https://pkg.go.dev/net/http#StatusBadRequest
	// 使用时不要直接传入code，而是传入 http.StatusBadRequest
	code int
	// 内部封装的错误
	err error
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Unwarp() error {
	return e.err
}

func New(code int, err error) *Error {
	return &Error{code, err}
}

// http.StatusConflict.
func Conflict(err error) *Error {
	return New(http.StatusConflict, err)
}

// http.StatusBadRequest.
func BadRequest(err error) *Error {
	return New(http.StatusBadRequest, err)
}

// http.StatusUnauthorized.
func Unauthorized(err error) *Error {
	return New(http.StatusUnauthorized, err)
}

// http.StatusNotFound.
func NotFound(err error) *Error {
	return New(http.StatusNotFound, err)
}

// http.StatusInternalServerError.
func InternalServer(err error) *Error {
	return New(http.StatusInternalServerError, err)
}
