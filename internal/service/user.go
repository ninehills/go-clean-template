package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ninehills/go-webapp-template/internal/dao"
	"github.com/ninehills/go-webapp-template/internal/entity"
	"github.com/ninehills/go-webapp-template/internal/infra/dependency"
	"github.com/ninehills/go-webapp-template/internal/infra/exception"
	"github.com/ninehills/go-webapp-template/pkg/cache"
	"github.com/ninehills/go-webapp-template/pkg/logger"
	"github.com/ninehills/go-webapp-template/pkg/password"
)

const userCacheKeyPrefix = "cache:user:"

// UserService 实现了 User 接口.
type UserService struct {
	db    dao.Querier
	l     logger.Logger
	cache cache.Cacher
	svcs  *Services
}

// New -.
func NewUserService(deps *dependency.Dependency, svcs *Services) *UserService {
	return &UserService{
		db:    deps.DAO,
		l:     deps.Logger,
		cache: deps.Cache,
		svcs:  svcs,
	}
}

// Create - 创建 User.
func (s *UserService) Create(ctx context.Context, in entity.User) (entity.User, error) {
	encryptPassword, err := password.EncryptPassword(in.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("- UserService - Create - encrypt password failed: %w", err)
	}

	arg := dao.CreateUserParams{
		Username:    in.Username,
		Status:      in.Status,
		Email:       in.Email,
		Description: in.Description,
		Password:    encryptPassword,
	}

	err = s.db.CreateUser(ctx, arg)
	if err != nil {
		// Ugly hack to handle duplicate key error, only for mysql.
		if strings.Contains(err.Error(), "Duplicate entry") {
			return entity.User{}, exception.Conflict(fmt.Errorf("duplicate name or user id"))
		}

		return entity.User{}, fmt.Errorf("- UserService - Create - create failed: %w", err)
	}

	u, err := s.db.GetUser(ctx, in.Username)
	if err != nil {
		return entity.User{}, fmt.Errorf("- UserService - Create - get failed: %w", err)
	}

	return entity.ToUser(u), nil
}

// Get - 根据 User ID 获取 User.
func (s *UserService) Get(ctx context.Context, username string) (entity.User, error) {
	u, err := s.db.GetUser(ctx, username)
	s.l.Ctx(ctx).Debugf("UserService - Get - username: %s, error[%v]", username, err)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, exception.Conflict(fmt.Errorf("user %s not found", username))
		}

		return entity.User{}, fmt.Errorf("- UserService - Get - get failed: %w", err)
	}

	return entity.ToUser(u), nil
}

// CacheGet -.
func (s *UserService) CacheGet(ctx context.Context, username string) (entity.User, error) {
	var user entity.User

	key := userCacheKeyPrefix + username

	err := s.cache.Get(ctx, key, &user)
	if err == nil {
		s.l.Ctx(ctx).Debugf("UserService - GetCache - cache hit %s: %+v", key, user)

		return user, nil
	}

	if errors.Is(err, cache.ErrMiss) {
		s.l.Ctx(ctx).Debugf("UserService - GetCache - cache miss %s", key)
	} else {
		s.l.Ctx(ctx).Warnf("UserService - GetCache - get from cache %s failed: %v", key, err)
	}

	u, err := s.Get(ctx, username)
	if err != nil {
		return entity.User{}, err
	}

	err = s.cache.Set(ctx, key, u)
	if err != nil {
		s.l.Ctx(ctx).Warnf("UserService - GetCache - set to cache %s/%+v failed: %v", key, u, err)
	}

	s.l.Ctx(ctx).Debugf("UserService - GetCache - set cache %s/%+v success", key, u)

	return u, nil
}

// QueryUser - 分页查询 User 信息.
func (s *UserService) Query(
	ctx context.Context, p entity.PageQuery, o entity.OrderQuery, u entity.UserQuery) (
	entity.PageResult, []entity.User, error,
) {
	us, count, err := s.db.QueryUser(ctx, dao.QueryUserParams{
		Offset:   (p.PageNo - 1) * p.PageSize,
		Limit:    p.PageSize,
		Order:    o.Order,
		OrderBy:  o.OrderBy,
		Username: u.Username,
		Status:   u.Status,
		Email:    u.Email,
	})
	if err != nil {
		return entity.PageResult{}, nil, fmt.Errorf("- UserService - ListWithPages - list failed: %w", err)
	}

	users := make([]entity.User, len(us))
	for i, u := range us {
		users[i] = entity.ToUser(u)
	}

	return entity.PageResult{
		PageNo:     p.PageNo,
		PageSize:   p.PageSize,
		TotalCount: count,
	}, users, nil
}

// Update - 更新 User.
func (s *UserService) Update(ctx context.Context, in entity.User) (entity.User, error) {
	// check if User exists
	_, err := s.db.GetUser(ctx, in.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, exception.NotFound(fmt.Errorf("user %s not found", in.Username))
		}

		return entity.User{}, fmt.Errorf("- UserService - Update - get failed: %w", err)
	}

	arg := dao.UpdateUserParams{
		Username:    in.Username,
		Status:      sql.NullInt32{Int32: in.Status, Valid: in.Status != 0},
		Email:       sql.NullString{String: in.Email, Valid: in.Email != ""},
		Description: sql.NullString{String: in.Description, Valid: in.Description != ""},
	}

	if in.Password != "" {
		encryptPassword, err := password.EncryptPassword(in.Password)
		if err != nil {
			return entity.User{}, fmt.Errorf("- UserService - Update - encrypt password failed: %w", err)
		}

		arg.Password = sql.NullString{String: encryptPassword, Valid: true}
	}

	err = s.db.UpdateUser(ctx, arg)
	if err != nil {
		return entity.User{}, fmt.Errorf("- UserService - Update - update failed: %w", err)
	}

	// 先更新后查询
	u, err := s.db.GetUser(ctx, in.Username)
	if err != nil {
		return entity.User{}, fmt.Errorf("- UserService - Update - get failed: %w", err)
	}

	// 删除缓存
	key := userCacheKeyPrefix + in.Username

	err = s.cache.Del(ctx, key)
	if err != nil {
		s.l.Ctx(ctx).Warnf("UserService - Update - del cache %s failed: %v", key, err)
	}

	return entity.ToUser(u), nil
}

// Delete - 删除 User，操作是幂等的，也就是如果 User 不存在时返回成功.
func (s *UserService) Delete(ctx context.Context, username string) error {
	// check if User exists
	_, err := s.db.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("- UserService - Delete - get failed: %w", err)
	}

	err = s.db.DeleteUser(ctx, username)
	if err != nil {
		return fmt.Errorf("- UserService - Delete - delete failed: %w", err)
	}

	// 删除缓存
	key := userCacheKeyPrefix + username

	err = s.cache.Del(ctx, key)
	if err != nil {
		s.l.Ctx(ctx).Warnf("UserService - Update - del cache %s failed: %v", key, err)
	}

	return nil
}

// AuthenticationPassword - 验证用户密码.
func (s *UserService) AuthenticationPassword(ctx context.Context, username, pass string) (bool, string, error) {
	u, err := s.db.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Sprintf("User %s not found", username), nil
		}

		return false, "", fmt.Errorf("- UserService - AuthenticationPassword - get failed: %w", err)
	}

	err = password.CompareHashAndPassword(u.Password, pass)
	if err != nil {
		return false, err.Error(), err
	}

	return true, "", nil
}
