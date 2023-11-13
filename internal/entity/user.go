package entity

import (
	"time"

	"github.com/ninehills/go-webapp-template/internal/dao"
)

const (
	// UserStatusActive 正常.
	UserStatusActive = 1
	// UserStatusInactive 禁用.
	UserStatusInactive = 2
)

// User Entity.
type User struct {
	// DB id.
	ID int64 `example:"1" json:"-"`
	// 用户的名称
	Username string `example:"twfbmbsr" json:"username"`
	// 用户状态，1代表启用，2代表禁用
	Status int32 `example:"1" json:"status"`
	// 邮箱
	Email string `example:"xxx@example.com" json:"email"`
	// 加密后的密码，并不输出到前端
	Password string `json:"-"`
	// 备注
	Description string `example:"twfbmbsr" json:"description"`
	// 创建时间
	CreatedAt time.Time `example:"2020-01-01T00:00:00Z" json:"createdAt"`
	// 更新时间
	UpdatedAt time.Time `example:"2020-01-01T00:00:00Z" json:"updatedAt"`
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// 将 dao.models.User 转为 entity.User, 忽略Password 字段.
func ToUser(user dao.User) User {
	return User{
		ID:          user.ID,
		Username:    user.Username,
		Status:      user.Status,
		Email:       user.Email,
		Description: user.Description,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

// 用户查询条件.
type UserQuery struct {
	Username string `json:"username"`
	Status   int32  `json:"status"`
	Email    string `json:"email"`
}
