package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zhongxic/sellbot/pkg/errorcode"
	"github.com/zhongxic/sellbot/pkg/result"
)

const TraceId = "X-Trace-Id"

type ResponseWriterWrapper struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w ResponseWriterWrapper) Write(b []byte) (int, error) {
	if n, err := w.Body.Write(b); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriterWrapper) WriteString(s string) (int, error) {
	if n, err := w.Body.WriteString(s); err != nil {
		return n, err
	}
	return w.ResponseWriter.WriteString(s)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.Request.Header.Get(TraceId)
		if requestId == "" {
			requestId = uuid.New().String()
		}
		c.Set(TraceId, requestId)
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		body, err := c.GetRawData()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				result.FailedWithErrorCode(errorcode.SystemError, "dump request failed"))
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		writer := &ResponseWriterWrapper{c.Writer, &bytes.Buffer{}}
		c.Writer = writer

		c.Next()

		status := c.Writer.Status()
		elapsed := time.Since(start).Milliseconds()
		response, _ := io.ReadAll(writer.Body)

		slog.Info("completed",
			slog.String("traceId", requestId),
			slog.Group("request",
				slog.String("path", path),
				slog.String("query", query),
				slog.String("body", string(body)),
				slog.Int64("elapsed", elapsed)),
			slog.Group("response",
				slog.Int64("status", int64(status)),
				slog.String("body", string(response))))
	}
}
