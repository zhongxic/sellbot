package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/pkg/model"
)

type Controller struct {
}

func (p *Controller) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, model.SuccessWithData("pong"))
}
