package routes

import (
	"sync"

	"github.com/gin-gonic/gin"
	botctl "github.com/zhongxic/sellbot/internal/controller/bot"
	"github.com/zhongxic/sellbot/internal/controller/ping"
	botserve "github.com/zhongxic/sellbot/internal/service/bot"
	"github.com/zhongxic/sellbot/pkg/middleware"
)

var (
	initOnce sync.Once
	engine   *gin.Engine
)

func Init() *gin.Engine {
	initOnce.Do(func() {
		engine = initRoutes()
	})
	return engine
}

func initRoutes() *gin.Engine {
	r := gin.New()
	registerMiddleware(r)
	registerRoutes(r)
	return r
}

func registerMiddleware(r *gin.Engine) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
}

func registerRoutes(r *gin.Engine) {
	pingController := ping.NewController()
	botOptions := botserve.Options{
		TestProcessDir:    "",
		ReleaseProcessDir: "",
	}
	botController := botctl.NewController(botserve.NewService(botOptions))
	r.GET("/ping", pingController.Ping)
	r.POST("/prologue", botController.Prologue)
}
