package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/pkg/model"
)

type PingController struct {
}

func (p *PingController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, model.SuccessWithData("pong"))
}
