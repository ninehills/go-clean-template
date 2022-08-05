// Package app configures and runs application.
package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/ninehills/go-webapp-template/config"
	"github.com/ninehills/go-webapp-template/internal/entity/validation"
	"github.com/ninehills/go-webapp-template/internal/infra/dependency"
	"github.com/ninehills/go-webapp-template/internal/infra/middleware"
	"github.com/ninehills/go-webapp-template/pkg/httpserver"

	http "github.com/ninehills/go-webapp-template/internal/controller/http"
)

// Run creates objects via constructors.
func Run(cfgFile string) {
	var err error
	// 加载日志
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	log.Printf("Config: %+v", cfg)

	// 初始化全部依赖
	dep := dependency.NewDependency(cfg)
	defer dep.Close()

	// 启动日志动态加载逻辑
	go config.ConfigWatcher(cfgFile, dep.ReloadLogger)

	l := dep.Logger

	// 初始化 Gin
	if !cfg.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	handler := gin.New()

	// 绑定自定义的 Validator 参数校验器
	validation.BindValidator()

	// 初始化中间件
	middleware.RegisterGlobalMiddleware(handler, dep.Logger)

	// 初始化 router
	l.Info().Msg("Controller router init...")
	http.NewRouter(handler, dep)
	l.Info().Msgf("Start http server at %s", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info().Msg("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error().Err(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}
	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error().Err(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
