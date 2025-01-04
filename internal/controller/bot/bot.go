package bot

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/internal/service/bot"
	"github.com/zhongxic/sellbot/internal/traceid"
	"github.com/zhongxic/sellbot/pkg/errorcode"
	"github.com/zhongxic/sellbot/pkg/middleware"
	"github.com/zhongxic/sellbot/pkg/result"
)

type Controller struct {
	botService bot.Service
}

func (c *Controller) Prologue(ctx *gin.Context) {
	traceId := ctx.GetString(middleware.ContextKeyTraceId)
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
	traceContext := context.WithValue(context.Background(), traceid.TraceId{}, traceId)
	interactiveRespond, err := c.botService.Prologue(traceContext, prologueDTO)
	if err != nil {
		slog.Error("process prologue request failed", "traceId", traceId, "error", err)
		ctx.JSON(http.StatusInternalServerError, result.FailedWithErrorCode(errorcode.SystemError, http.StatusText(http.StatusInternalServerError)))
		return
	}
	interactiveResponse := convertInteractiveRespondToInteractiveResponse(interactiveRespond)
	ctx.JSON(http.StatusOK, result.SuccessWithData(interactiveResponse))
}

func (c *Controller) Connect(ctx *gin.Context) {
	traceId := ctx.GetString(middleware.ContextKeyTraceId)
	request := &ConnectRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		slog.Error("bind connect request failed", "traceId", traceId, "error", err)
		ctx.JSON(http.StatusBadRequest, result.FailedWithErrorCode(errorcode.ParamsError, errorcode.MessageRequestBodyNotBindable))
		return
	}
	slog.Info("connect request received", "traceId", traceId, "body", request)
	if request.SessionId == "" {
		ctx.JSON(http.StatusBadRequest, result.FailedWithErrorCode(errorcode.ParamsError, "sessionId is required"))
		return
	}
	connectDTO := convertConnectRequestToConnectDTO(request)
	traceContext := context.WithValue(context.Background(), traceid.TraceId{}, traceId)
	connectRespond, err := c.botService.Connect(traceContext, connectDTO)
	if err != nil {
		slog.Error("process connect request failed", "traceId", traceId, "error", err)
		ctx.JSON(http.StatusInternalServerError, result.FailedWithErrorCode(errorcode.SystemError, http.StatusText(http.StatusInternalServerError)))
		return
	}
	connectResponse := convertConnectRespondToResponse(connectRespond)
	ctx.JSON(http.StatusOK, result.SuccessWithData(connectResponse))
}

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

func NewController(botService bot.Service) *Controller {
	return &Controller{botService: botService}
}
