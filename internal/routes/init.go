package routes

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/internal/controller/ping"
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
	pingController := &ping.Controller{}
	r.GET("/ping", pingController.Ping)
}
