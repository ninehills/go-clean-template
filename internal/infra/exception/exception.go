package exception

import (
	"errors"
	"net/http"
)

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

func NewError(code int, err error) *Error {
	return &Error{code, err}
}

// http.StatusConflict.
func Conflict(err error) *Error {
	return NewError(http.StatusConflict, err)
}

// http.StatusBadRequest.
func BadRequest(err error) *Error {
	return NewError(http.StatusBadRequest, err)
}

// http.StatusUnauthorized.
func Unauthorized(err error) *Error {
	return NewError(http.StatusUnauthorized, err)
}

// http.StatusNotFound.
func NotFound(err error) *Error {
	return NewError(http.StatusNotFound, err)
}

// http.StatusInternalServerError.
func InternalServer(err error) *Error {
	return NewError(http.StatusInternalServerError, err)
}

// 错误比较.
func Is(err error, target error) bool {
	if err == nil && target == nil {
		return true
	}
	var e1, e2 *Error
	if errors.As(err, &e1) && errors.As(target, &e2) {
		return e1.Code() == e2.Code()
	}

	return errors.Is(err, target)
}
