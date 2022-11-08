package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ninehills/go-webapp-template/apis/httpv1"
	"github.com/ninehills/go-webapp-template/internal/entity"
	"github.com/ninehills/go-webapp-template/internal/infra/exception"
	"github.com/ninehills/go-webapp-template/internal/infra/middleware"
	"github.com/ninehills/go-webapp-template/internal/service"
	"github.com/ninehills/go-webapp-template/pkg/logger"
	"github.com/ninehills/go-webapp-template/pkg/password"
)

type userRoutes struct {
	s service.User
	l logger.Logger
}

func newUserRoutes(handler *gin.RouterGroup, l logger.Logger, serv *service.Services, midd *middleware.Middlewares) {
	r := &userRoutes{
		l: l,
		s: serv.User,
	}
	handler.POST("/users",
		midd.Audit.Audit(),
		r.createUser)
	handler.GET("/users",
		r.ListUsers)
	handler.GET("/users/:username",
		r.getUser)
	handler.PUT("/users/:username",
		midd.Audit.Audit(),
		r.updateUser)
	handler.DELETE("/users/:username",
		midd.Audit.Audit(),
		r.deleteUser)
}

// @Summary     Get user
// @Description Get user by username
// @ID          get-user
// @Tags  	    user
// @Param 		username path string true "Username"
// @Produce     json
// @Success     200 {object} httpv1.GetUserResponse
// @Failure     500 {object} httpv1.ErrorResponse
// @Router      /v1/users/:username [get].
func (r *userRoutes) getUser(c *gin.Context) {
	user, err := r.s.Get(c, c.Param("username"))
	if err != nil {
		r.l.Ctx(c).Err(err).Error("http - v1 - getUser failed")
		exception.ResponseWithError(c, err)

		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary     Update user
// @Description Update user by username
// @ID          update-user
// @Tags  	    user
// @Param 		username path string true "username"
// @Produce     json
// @Success     200 {object} httpv1.UpdateUserResponse
// @Failure     500 {object} httpv1.ErrorResponse
// @Router      /v1/users/:username [PUT].
func (r *userRoutes) updateUser(c *gin.Context) {
	username := c.Param("username")

	var request httpv1.UpdateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Ctx(c).Err(err).Warn("http - v1 - updateUser invalid request body")
		exception.CodeResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	// 当请求中密码不为空时，才会更新密码
	if request.Password != "" || request.ConfirmPassword != "" {
		err := password.ValidatePassword(request.Password, request.ConfirmPassword)
		if err != nil {
			r.l.Ctx(c).Err(err).Warn("http - v1 - updateUser password is invalid")
			exception.CodeResponse(c, http.StatusBadRequest, err.Error())

			return
		}
	}

	user, err := r.s.Update(
		c, entity.User{
			Username:    username,
			Status:      request.Status,
			Email:       request.Email,
			Description: request.Description,
			Password:    request.Password,
		},
	)
	if err != nil {
		r.l.Ctx(c).Err(err).Error("http - v1 - updateUser failed")
		exception.ResponseWithError(c, err)

		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary     List users
// @Description List user with pages
// @ID          list-users
// @Tags  	    user
// @Param		pageNo		query	int64	true	"Page number"
// @Param		pageSize	query	int64	true	"Page size"
// @Param		order		query	string	true	"Order asc/desc"
// @Param		orderBy		query	string	true	"Order by create_time"
// @Param		username	query	string	true	"Username"
// @Param		status		query	int32	true	"Status 1/2"
// @Produce     json
// @Success     200 {object} httpv1.ListUserResponse
// @Failure     500 {object} httpv1.ErrorResponse
// @Router      /v1/users [get].
func (r *userRoutes) ListUsers(c *gin.Context) {
	var request httpv1.ListUserRequest

	err := c.ShouldBindQuery(&request)
	if err != nil {
		r.l.Ctx(c).Err(err).Warn("http - v1 - ListUsers invalid request body")
		exception.CodeResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	pageresult, users, sErr := r.s.Query(
		c, entity.PageQuery{
			PageNo:   request.PageNo,
			PageSize: request.PageSize,
		}, entity.OrderQuery{
			Order:   request.Order,
			OrderBy: request.OrderBy,
		}, entity.UserQuery{
			Username: request.Username,
			Status:   request.Status,
			Email:    request.Email,
		},
	)
	if sErr != nil {
		r.l.Ctx(c).Err(err).Error("http - v1 - ListUsers failed")
		exception.ResponseWithError(c, sErr)

		return
	}

	c.JSON(http.StatusOK, httpv1.ListUserResponse{
		PageNo:     pageresult.PageNo,
		PageSize:   pageresult.PageSize,
		TotalCount: pageresult.TotalCount,
		Result:     users,
	})
}

// @Summary     Create user
// @Description Create user, user_id is generated random
// @ID          create-user
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Param       request body httpv1.CreateUserRequest true "Set up user"
// @Success     200 {object} httpv1.CreateUserResponse
// @Failure     400 {object} httpv1.ErrorResponse
// @Failure     500 {object} httpv1.ErrorResponse
// @Router      /v1/users [post].
func (r *userRoutes) createUser(c *gin.Context) {
	var request httpv1.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Ctx(c).Err(err).Warn("http - v1 - createUser invalid request body")
		exception.CodeResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	err := password.ValidatePassword(request.Password, request.ConfirmPassword)
	if err != nil {
		r.l.Ctx(c).Err(err).Warn("http - v1 - createUser password is invalid")
		exception.CodeResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	user := entity.User{
		Username:    request.Username,
		Email:       request.Email,
		Description: request.Description,
		Password:    request.Password,
		Status:      entity.UserStatusActive,
	}

	ret, err := r.s.Create(c, user)
	if err != nil {
		r.l.Ctx(c).Err(err).Error("http - v1 - createUser user failed")
		exception.ResponseWithError(c, err)

		return
	}

	c.JSON(http.StatusOK, ret)
}

// @Summary     Delete user
// @Description Delete user by username
// @ID          delete-user
// @Tags  	    user
// @Param 		username path string true "Username"
// @Produce     json
// @Success     200 {object} httpv1.DeleteUserResponse
// @Failure     500 {object} httpv1.ErrorResponse
// @Router      /v1/users/:username [DELETE].
func (r *userRoutes) deleteUser(c *gin.Context) {
	username := c.Param("username")

	if err := r.s.Delete(c, username); err != nil {
		r.l.Ctx(c).Err(err).Error("http - v1 - deleteUser failed")
		exception.ResponseWithError(c, err)

		return
	}

	c.Status(http.StatusOK)
}
