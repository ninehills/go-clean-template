package validation

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func BindValidator() {
	var err error

	// Register binding
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err = v.RegisterValidation("username", UsernameValidator())
		if err != nil {
			panic(err)
		}
	}
}
