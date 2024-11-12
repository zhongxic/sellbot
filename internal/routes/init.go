package routes

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/internal/controller/ping"
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
	r := gin.Default()
	registerRoutes(r)
	return r
}

func registerRoutes(r *gin.Engine) {
	pingController := &ping.PingController{}
	r.GET("/ping", pingController.Ping)
}
