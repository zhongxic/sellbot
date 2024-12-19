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
	interactiveRespond, err := c.botService.Prologue(prologueDTO)
	if err != nil {
		slog.Error("process prologue request failed", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, result.FailedWithErrorCode(errorcode.SystemError, http.StatusText(http.StatusInternalServerError)))
		return
	}
	interactiveResponse := convertInteractiveRespondToInteractiveResponse(interactiveRespond)
	ctx.JSON(http.StatusOK, result.SuccessWithData(interactiveResponse))
}
