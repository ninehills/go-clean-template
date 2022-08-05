package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 自定义的Username validator
// 参考 https://github.com/gin-gonic/examples/blob/master/custom-validation/server.go
var usernameValidator validator.Func = func(fl validator.FieldLevel) bool {
	username, ok := fl.Field().Interface().(string)
	r := `^[A-Za-z0-9_]*$`
	if ok {
		if username == "" {
			// 允许空字符串，此处是为了兼容搜索时的空字符串
			// 长度由别的validator保证
			return true
		}
		return regexp.MustCompile(r).MatchString(username)
	}
	return false
}
