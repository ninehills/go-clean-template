// Package usecase implements application business logic. Each logic group in own file.
package service

import (
	"context"

	"github.com/ninehills/go-webapp-template/internal/entity"
	"github.com/ninehills/go-webapp-template/internal/infra/dependency"
)

// 定义 Service 聚合结构
type Services struct {
	User User
}

// 创建所有 Service，另外将srvs 注入到各个 Service 中，方便相互之间的引用
func NewServices(deps *dependency.Dependency) *Services {
	svcs := &Services{}
	svcs.User = NewUserService(deps, svcs)
	return svcs
}

type (
	// User Interface
	User interface {
		// 创建用户
		Create(ctx context.Context, in entity.User) (entity.User, error)
		// 根据用户名称获取用户
		Get(ctx context.Context, username string) (entity.User, error)
		// 根据用户名称获取用户（带 Cache）
		CacheGet(ctx context.Context, username string) (entity.User, error)
		// 更新用户
		Update(ctx context.Context, in entity.User) (entity.User, error)
		// 删除用户
		Delete(ctx context.Context, username string) error
		// 分页查询用户信息，支持丰富的查询条件
		Query(ctx context.Context, p entity.PageQuery, o entity.OrderQuery, u entity.UserQuery) (entity.PageResult, []entity.User, error)
		// 验证密码是否正确
		AuthenticationPassword(ctx context.Context, username, password string) (ok bool, reason string, err error)
	}
)
