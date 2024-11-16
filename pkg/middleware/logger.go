package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/pkg/model"
)

const (
	contentType     = "Content-Type"
	applicationJson = "application/json"

	msgErrDumpRequest  = "dump request failed"
	msgErrDumpResponse = "dump response failed"
	msgNotInJsonFormat = "not in json format"
)

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
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		body := dumpRequest(c)

		writer := &ResponseWriterWrapper{c.Writer, &bytes.Buffer{}}
		c.Writer = writer

		c.Next()

		status := c.Writer.Status()
		elapsed := time.Since(start).Milliseconds()
		response := dumpResponse(writer, c)

		slog.Info("completed",
			slog.Group("request",
				slog.String("path", path),
				slog.String("query", query),
				slog.Any("body", body),
				slog.Int64("elapsed", elapsed)),
			slog.Group("response",
				slog.Int64("status", int64(status)),
				slog.Any("body", response)))
	}
}

func dumpRequest(c *gin.Context) any {
	if c.Request.Header.Get(contentType) == applicationJson {
		if data, err := io.ReadAll(c.Request.Body); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				model.FailedWithMessage(model.SYSTEM_ERROR, msgErrDumpRequest))
			return msgErrDumpRequest
		} else {
			body := map[string]any{}
			if err = json.Unmarshal(data, &body); err != nil {
				return msgErrDumpRequest
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
			return body
		}
	}
	return msgNotInJsonFormat
}

func dumpResponse(writer *ResponseWriterWrapper, c *gin.Context) any {
	if strings.Contains(c.Writer.Header().Get(contentType), applicationJson) {
		if data, err := io.ReadAll(writer.Body); err != nil {
			return msgErrDumpResponse
		} else {
			response := map[string]any{}
			if err := json.Unmarshal(data, &response); err != nil {
				return msgErrDumpResponse
			}
			return response
		}
	}
	return msgNotInJsonFormat
}
