package routes

import (
	"sync"

	"github.com/gin-gonic/gin"
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
	return r
}
