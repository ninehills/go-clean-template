package httpv1

import "github.com/ninehills/go-webapp-template/internal/entity"

type GetUserResponse entity.User

type UpdateUserRequest struct {
	Status          int32  `binding:"omitempty,oneof=1 2"          json:"status"`
	Email           string `binding:"omitempty,min=1,max=64,email" json:"email"`
	Description     string `binding:"omitempty,min=0,max=140"      json:"description"`
	Password        string `binding:"omitempty"                    json:"password"` // 密码有单独的方法进行校验
	ConfirmPassword string `binding:"omitempty"                    json:"confirmPassword"`
}

type UpdateUserResponse entity.User

type ListUserRequest struct {
	PageNo   int64  `binding:"gte=1"                                                form:"pageNo,default=1"`
	PageSize int64  `binding:"gte=1"                                                form:"pageSize,default=100"`
	Order    string `binding:"omitempty,oneof='asc' 'desc'"                         form:"order,default="`
	OrderBy  string `binding:"omitempty,oneof='created_at' 'updated_at' 'username'" form:"orderBy,default="`
	Username string `binding:"omitempty,min=1,max=64,username"                      form:"username,default="`
	Status   int32  `binding:"omitempty,oneof=1 2"                                  form:"status,default=0"`
	Email    string `binding:"omitempty,min=1,max=64,email"                         form:"email,default="`
}

type ListUserResponse struct {
	PageNo     int64         `json:"pageNo"`
	PageSize   int64         `json:"pageSize"`
	TotalCount int64         `json:"totalCount"`
	Result     []entity.User `json:"result"`
}

type CreateUserRequest struct {
	// username 应该是长度1-64位的字母、数字或"_"
	Username    string `binding:"username,min=1,max=64,required" json:"username"`
	Email       string `binding:"required,email,min=0,max=64"    json:"email"`
	Description string `binding:"min=0,max=140"                  json:"description"`
	// 密码的验证比较复杂，有单独的方法进行验证
	Password        string `binding:"required" json:"password"`
	ConfirmPassword string `binding:"required" json:"confirmPassword"`
}

type CreateUserResponse entity.User

type DeleteUserResponse struct{}
