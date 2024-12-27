package bot

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/internal/traceid"
	"github.com/zhongxic/sellbot/pkg/errorcode"
	"github.com/zhongxic/sellbot/pkg/middleware"
	"github.com/zhongxic/sellbot/pkg/result"
)

func (c *Controller) Chat(ctx *gin.Context) {
	traceId := ctx.GetString(middleware.ContextKeyTraceId)
	request := &ChatRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		slog.Error("bind chat request failed", "traceId", traceId, "error", err)
		ctx.JSON(http.StatusBadRequest, result.FailedWithErrorCode(errorcode.ParamsError, errorcode.MessageRequestBodyNotBindable))
		return
	}
	slog.Info("chat request received", "traceId", traceId, "body", request)
	if request.SessionId == "" {
		ctx.JSON(http.StatusBadRequest, result.FailedWithErrorCode(errorcode.ParamsError, "sessionId is required"))
		return
	}
	chatDTO := convertChatRequestToChatDTO(request)
	traceContext := context.WithValue(context.Background(), traceid.TraceId{}, traceId)
	interactiveRespond, err := c.botService.Chat(traceContext, chatDTO)
	if err != nil {
		slog.Error("process chat request failed", "traceId", traceId, "error", err)
		ctx.JSON(http.StatusInternalServerError, result.FailedWithErrorCode(errorcode.SystemError, http.StatusText(http.StatusInternalServerError)))
		return
	}
	interactiveResponse := convertInteractiveRespondToInteractiveResponse(interactiveRespond)
	ctx.JSON(http.StatusOK, result.SuccessWithData(interactiveResponse))
}
