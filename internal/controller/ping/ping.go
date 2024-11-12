package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingController struct {
}

func (p *PingController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
