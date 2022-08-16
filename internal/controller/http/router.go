package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ninehills/go-webapp-template/internal/entity"
	"github.com/ninehills/go-webapp-template/internal/infra/dependency"
	"github.com/ninehills/go-webapp-template/internal/infra/middleware"
	"github.com/ninehills/go-webapp-template/internal/service"
)

// NewRouter -.
// Swagger spec:
// @title       GO WEBAPP TEMPLATE API
// @description GO WEBAPP TEMPLATE API
// @version     1.0
// @host        localhost:8080
// @BasePath    /.
func NewRouter(handler *gin.Engine, deps *dependency.Dependency) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 创建所有 Service
	svcs := service.NewServices(deps)

	// 创建非全局的 middleware
	middlewares := middleware.NewMiddlewares(deps)

	// Init default user
	InitDefaultUser(svcs.User, deps.Config.App.SuperUser, deps.Config.App.SuperPassword)

	l := deps.Logger

	// Routers

	// v1 API
	v1 := handler.Group("/v1")
	{
		newUserRoutes(v1, l, svcs, middlewares)
	}
}

func InitDefaultUser(user service.User, username, password string) {
	_, err := user.Get(context.Background(), username)
	if err != nil {
		log.Printf("Init default user %s", username)

		_, err := user.Create(context.Background(), entity.User{
			Username:    username,
			Password:    password,
			Email:       fmt.Sprintf("%s@example.com", username),
			Status:      entity.UserStatusActive,
			Description: "Default created super user",
		})
		if err != nil {
			log.Printf("Init default user failed: %s", err)
			panic(err)
		}
	}
}
