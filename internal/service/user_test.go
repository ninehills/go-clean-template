package service_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ninehills/go-webapp-template/internal/dao"
	"github.com/ninehills/go-webapp-template/internal/entity"
	"github.com/ninehills/go-webapp-template/internal/infra/dependency"
	"github.com/ninehills/go-webapp-template/internal/infra/exception"
	"github.com/ninehills/go-webapp-template/internal/service"
	"github.com/ninehills/go-webapp-template/mocks"
	"github.com/ninehills/go-webapp-template/pkg/cache"
	"github.com/ninehills/go-webapp-template/pkg/logger"
)

var errInternal = errors.New("internal error")

type test struct {
	name string
	id   string
	mock func(string)
	res  interface{}
	err  error
}

func bootstrap(t *testing.T) (*service.UserService, *mocks.MockQuerier, *mocks.MockCacher) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	querier := mocks.NewMockQuerier(mockCtl)
	cacher := mocks.NewMockCacher(mockCtl)

	userService := service.NewUserService(&dependency.Dependency{
		DAO: querier,
		Logger: logger.New(logger.Config{
			Format:  "text",
			Level:   "debug",
			NoColor: false,
		}),
		Cache: cacher,
	}, &service.Services{})

	return userService, querier, cacher
}

// 自定义 UserMatcher，只比较 Username.
type userMatcher struct {
	Username string
}

func (e userMatcher) Matches(x interface{}) bool {
	arg, ok := x.(dao.CreateUserParams)
	if !ok {
		return false
	}

	return e.Username == arg.Username
}

func (e userMatcher) String() string {
	return fmt.Sprintf("is equal to %v", e.Username)
}

func TestUserCreate(t *testing.T) {
	t.Parallel()
	userService, querier, _ := bootstrap(t)

	tests := []test{
		{
			name: "Create() - Duplicate entry",
			id:   uuid.NewString(),
			mock: func(id string) {
				querier.EXPECT().CreateUser(
					context.Background(), userMatcher{id},
				).Return(errors.New(":Duplicate entry"))
			},
			res: entity.User{},
			err: exception.Conflict(nil),
		},
		{
			name: "Create() - create failed",
			id:   uuid.NewString(),
			mock: func(id string) {
				querier.EXPECT().CreateUser(
					context.Background(), userMatcher{id},
				).Return(errInternal)
			},
			res: entity.User{},
			err: errInternal,
		},
		{
			name: "Create() - get failed",
			id:   uuid.NewString(),
			mock: func(id string) {
				querier.EXPECT().CreateUser(
					context.Background(), userMatcher{id},
				).Return(nil)
				querier.EXPECT().GetUser(context.Background(), id).Return(dao.User{}, errInternal)
			},
			res: entity.User{},
			err: errInternal,
		},
		{
			name: "Create() - success",
			id:   uuid.NewString(),
			mock: func(id string) {
				querier.EXPECT().CreateUser(
					context.Background(), userMatcher{id},
				).Return(nil)
				querier.EXPECT().GetUser(context.Background(), id).Return(dao.User{}, nil)
			},
			res: entity.User{},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock(tc.id)
			res, err := userService.Create(context.Background(), entity.User{Username: tc.id})
			require.EqualValues(t, res, tc.res)
			require.True(t, exception.Is(err, tc.err)) //nolint:testifylint
		})
	}
}

func TestUserGet(t *testing.T) {
	t.Parallel()

	userService, querier, _ := bootstrap(t)

	tests := []test{
		{
			name: "Get() - empty",
			id:   uuid.NewString(),
			mock: func(id string) {
				querier.EXPECT().GetUser(context.Background(), id).Return(dao.User{}, sql.ErrNoRows)
			},
			res: entity.User{},
			err: exception.NotFound(nil),
		},
		{
			name: "Get() - failed",
			id:   uuid.NewString(),
			mock: func(id string) {
				querier.EXPECT().GetUser(context.Background(), id).Return(dao.User{}, errInternal)
			},
			res: entity.User{},
			err: errInternal,
		},
		{
			name: "Get() - success",
			id:   uuid.NewString(),
			mock: func(id string) {
				querier.EXPECT().GetUser(context.Background(), id).Return(dao.User{}, nil)
			},
			res: entity.User{},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock(tc.id)
			res, err := userService.Get(context.Background(), tc.id)
			require.EqualValues(t, res, tc.res)
			require.True(t, exception.Is(err, tc.err)) //nolint:testifylint
		})
	}
}

func TestUserCacheGet(t *testing.T) {
	t.Parallel()

	userService, querier, cacher := bootstrap(t)

	userCacheKey := func(username string) string {
		return fmt.Sprintf("%s%s", service.UserCacheKeyPrefix, username)
	}

	tests := []test{
		{
			name: "CacheGet() - cache hit",
			id:   uuid.NewString(),
			mock: func(id string) {
				cacher.EXPECT().Get(context.Background(), userCacheKey(id), gomock.Any()).Return(nil)
			},
			res: entity.User{},
			err: nil,
		},
		{
			name: "CacheGet() - cache failed but db success",
			id:   uuid.NewString(),
			mock: func(id string) {
				cacher.EXPECT().Get(context.Background(), userCacheKey(id), gomock.Any()).Return(errInternal)
				querier.EXPECT().GetUser(context.Background(), id).Return(dao.User{}, nil)
				cacher.EXPECT().Set(context.Background(), userCacheKey(id), gomock.Any()).Return(nil)
			},
			res: entity.User{},
			err: nil,
		},
		{
			name: "CacheGet() - cache miss and db fail",
			id:   uuid.NewString(),
			mock: func(id string) {
				cacher.EXPECT().Get(context.Background(), userCacheKey(id), gomock.Any()).Return(cache.ErrMiss)
				querier.EXPECT().GetUser(context.Background(), id).Return(dao.User{}, errInternal)
			},
			res: entity.User{},
			err: errInternal,
		},
		{
			name: "CacheGet() - 缓存设定失败，但是整个 CacheGet 会成功",
			id:   uuid.NewString(),
			mock: func(id string) {
				cacher.EXPECT().Get(context.Background(), userCacheKey(id), gomock.Any()).Return(cache.ErrMiss)
				querier.EXPECT().GetUser(context.Background(), id).Return(dao.User{}, nil)
				cacher.EXPECT().Set(context.Background(), userCacheKey(id), gomock.Any()).Return(errInternal)
			},
			res: entity.User{},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock(tc.id)
			res, err := userService.CacheGet(context.Background(), tc.id)
			require.EqualValues(t, res, tc.res)
			require.True(t, exception.Is(err, tc.err)) //nolint:testifylint
		})
	}
}
