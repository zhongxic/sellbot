package routes

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/config"
	botctl "github.com/zhongxic/sellbot/internal/controller/bot"
	"github.com/zhongxic/sellbot/internal/controller/ping"
	botserve "github.com/zhongxic/sellbot/internal/service/bot"
	"github.com/zhongxic/sellbot/internal/service/process"
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
	botOptions := botserve.Options{
		DictFile:          cfg.Tokenizer.DictFile,
		TestProcessDir:    cfg.Process.Directory.Test,
		ReleaseProcessDir: cfg.Process.Directory.Release,
	}
	testLoader := process.NewFileLoader(cfg.Process.Directory.Test)
	releaseLoader := process.NewFileLoader(cfg.Process.Directory.Release)
	botController := botctl.NewController(botserve.NewService(botOptions, testLoader, releaseLoader))
	r.GET("/ping", pingController.Ping)
	r.POST("/prologue", botController.Prologue)
}
