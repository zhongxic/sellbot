package bot

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/pkg/errorcode"
	"github.com/zhongxic/sellbot/pkg/middleware"
	"github.com/zhongxic/sellbot/pkg/result"
)

func (c *Controller) Prologue(ctx *gin.Context) {
	traceId := ctx.GetString(middleware.TraceId)
	request := &PrologueRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		slog.Error("bind prologue request failed", "traceId", traceId, "error", err)
		ctx.JSON(http.StatusBadRequest, result.FailedWithErrorCode(errorcode.ParamsError, errorcode.MessageRequestBodyNotBindable))
		return
	}
	slog.Info("prologue request received", "traceId", traceId, "body", request)
	if request.ProcessId == "" {
		ctx.JSON(http.StatusBadRequest, result.FailedWithErrorCode(errorcode.ParamsError, "processId is required"))
		return
	}
	prologueDTO := convertPrologueRequestToPrologueDTO(request)
	traceContext := context.WithValue(context.Background(), "traceId", traceId)
	interactiveRespond, err := c.botService.Prologue(traceContext, prologueDTO)
	if err != nil {
		slog.Error("process prologue request failed", "traceId", traceId, "error", err)
		ctx.JSON(http.StatusInternalServerError, result.FailedWithErrorCode(errorcode.SystemError, http.StatusText(http.StatusInternalServerError)))
		return
	}
	interactiveResponse := convertInteractiveRespondToInteractiveResponse(interactiveRespond)
	ctx.JSON(http.StatusOK, result.SuccessWithData(interactiveResponse))
}
