package httpv1

import "github.com/ninehills/go-webapp-template/internal/entity"

type GetUserResponse entity.User

type UpdateUserRequest struct {
	Status          int32  `json:"status" binding:"omitempty,oneof=1 2"`
	Email           string `json:"email" binding:"omitempty,min=1,max=64,email"`
	Description     string `json:"description" binding:"omitempty,min=0,max=140"`
	Password        string `json:"password" binding:"omitempty"` // 密码有单独的方法进行校验
	ConfirmPassword string `json:"confirmPassword" binding:"omitempty"`
}

type UpdateUserResponse entity.User

type ListUserRequest struct {
	PageNo   int64  `form:"pageNo,default=1" binding:"gte=1"`
	PageSize int64  `form:"pageSize,default=100" binding:"gte=1"`
	Order    string `form:"order,default=" binding:"omitempty,oneof='asc' 'desc'"`
	OrderBy  string `form:"orderBy,default=" binding:"omitempty,oneof='created_at' 'updated_at' 'username'"`
	Username string `form:"username,default=" binding:"omitempty,min=1,max=64,username"`
	Status   int32  `form:"status,default=0" binding:"omitempty,oneof=1 2"`
	Email    string `form:"email,default=" binding:"omitempty,min=1,max=64,email"`
}

type ListUserResponse struct {
	PageNo     int64         `json:"pageNo"`
	PageSize   int64         `json:"pageSize"`
	TotalCount int64         `json:"totalCount"`
	Result     []entity.User `json:"result"`
}

type CreateUserRequest struct {
	// username 应该是长度1-64位的字母、数字或"_"
	Username    string `json:"username"  binding:"username,min=1,max=64,required"`
	Email       string `json:"email" binding:"required,email,min=0,max=64"`
	Description string `json:"description" binding:"min=0,max=140"`
	// 密码的验证比较复杂，有单独的方法进行验证
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

type CreateUserResponse entity.User
