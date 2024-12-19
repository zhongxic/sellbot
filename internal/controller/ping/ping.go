package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/pkg/result"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, result.SuccessWithData("pong"))
}
