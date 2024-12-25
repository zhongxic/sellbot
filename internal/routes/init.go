package routes

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/config"
	botctl "github.com/zhongxic/sellbot/internal/controller/bot"
	"github.com/zhongxic/sellbot/internal/controller/ping"
	botserve "github.com/zhongxic/sellbot/internal/service/bot"
	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/pkg/cache"
	"github.com/zhongxic/sellbot/pkg/middleware"
)

var (
	initOnce sync.Once
	engine   *gin.Engine
)

func Init(cfg *config.Config) *gin.Engine {
	initOnce.Do(func() {
		engine = initRoutes(cfg)
	})
	return engine
}

func initRoutes(cfg *config.Config) *gin.Engine {
	r := gin.New()
	registerMiddleware(r)
	registerRoutes(r, cfg)
	return r
}

func registerMiddleware(r *gin.Engine) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
}

func registerRoutes(r *gin.Engine, cfg *config.Config) {
	pingController := ping.NewController()
	testProcessStorage := cache.NewCache[*process.Process](cache.Options{
		DefaultExpiration: time.Duration(cfg.Process.Cache.Expiration) * time.Second,
		CleanupInterval:   time.Duration(cfg.Process.Cache.CleanupInterval) * time.Second,
	})
	releaseProcessStorage := cache.NewCache[*process.Process](cache.Options{
		DefaultExpiration: time.Duration(cfg.Process.Cache.Expiration) * time.Second,
		CleanupInterval:   time.Duration(cfg.Process.Cache.CleanupInterval) * time.Second,
	})
	testLoader := process.NewCachedLoader(process.NewFileLoader(cfg.Process.Directory.Test), testProcessStorage)
	releaseLoader := process.NewCachedLoader(process.NewFileLoader(cfg.Process.Directory.Release), releaseProcessStorage)
	botOptions := botserve.Options{
		ExtraDict:     cfg.Tokenizer.ExtraDict,
		TestLoader:    testLoader,
		ReleaseLoader: releaseLoader,
	}
	botController := botctl.NewController(botserve.NewService(botOptions))
	r.GET("/ping", pingController.Ping)
	r.POST("/prologue", botController.Prologue)
}
