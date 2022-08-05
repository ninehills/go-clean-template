package dependency

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/ninehills/go-webapp-template/config"
	"github.com/ninehills/go-webapp-template/internal/dao"
	"github.com/ninehills/go-webapp-template/pkg/cache"
	"github.com/ninehills/go-webapp-template/pkg/logger"
	"github.com/ninehills/go-webapp-template/pkg/mysql"
)

// 全局依赖
type Dependency struct {
	Config *config.Config
	Logger *logger.Logger
	MySQL  *mysql.MySQL
	DAO    dao.Querier
	Redis  *redis.Client
	Cache  cache.Cacher
}

// 动态加载日志级别
func (d *Dependency) ReloadLogger(cfg *config.Config) {
	l, err := logger.New(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		d.Logger.Fatal().Err(err).Msgf("base - ReloadLogger - logger.New: level[%s] format[%s]", cfg.Log.Level, cfg.Log.Format)
	}
	d.Logger.Info().Msgf("base - ReloadLogger - logger.New: level[%s] format[%s]", cfg.Log.Level, cfg.Log.Format)
	d.Logger = l
}

// 初始化全局依赖
func NewDependency(cfg *config.Config) *Dependency {
	// 初始化日志 logger
	l, err := logger.New(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		panic(err)
	}

	// 初始化 MySQL 数据库
	ms, err := mysql.New(
		cfg.MySQL.DSN,
		mysql.ConnMaxLifetime(cfg.MySQL.ConnMaxLifetime),
		mysql.MaxOpenConns(cfg.MySQL.MaxOpenConns),
		mysql.MaxIdleConns(cfg.MySQL.MaxIdleConns),
	)
	if err != nil {
		l.Fatal().Err(fmt.Errorf("base - NewDependency - mysql.New: %w", err))
	}
	// Dao 初始化，生成 queries 对象
	queries := dao.New(ms.DB)

	// 初始化 Redis 数据库
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		l.Fatal().Err(fmt.Errorf("base - NewDependency - redis parse url failed: %w", err))
	}
	rdb := redis.NewClient(opt)

	// 初始化 Cache，默认过期时间是5分钟
	cache := cache.NewCache(rdb, cache.DefaultCacheExpires)

	deps := Dependency{
		Config: cfg,
		Logger: l,
		MySQL:  ms,
		DAO:    queries,
		Redis:  rdb,
		Cache:  cache,
	}

	return &deps
}

func (d *Dependency) Close() {
	d.MySQL.DB.Close()
	d.Redis.Close()
}
