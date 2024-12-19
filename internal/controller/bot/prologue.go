package bot

import (
	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/pkg/errorcode"
	"github.com/zhongxic/sellbot/pkg/result"
	"log/slog"
	"net/http"
)

func (c *Controller) Prologue(ctx *gin.Context) {
	request := &PrologueRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		slog.Error("process prologue request failed", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, result.FailedWithErrorCode(errorcode.ParamsError, errorcode.MessageRequestBodyNotBindable))
		return
	}
	slog.Info("prologue request received", "body", request)
	prologueDTO := convertPrologueRequestToPrologueDTO(request)
	interactiveRespond := c.botService.Prologue(prologueDTO)
	interactiveResponse := convertInteractiveRespondToInteractiveResponse(interactiveRespond)
	ctx.JSON(http.StatusOK, result.SuccessWithData(interactiveResponse))
}
