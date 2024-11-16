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
	content_type     = "Content-Type"
	application_json = "application/json"
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
	if c.Request.Header.Get(content_type) == application_json {
		if data, err := io.ReadAll(c.Request.Body); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				model.FailedWithMessage(model.SYSTEM_ERROR, "dump request failed"))
		} else {
			body := map[string]any{}
			if err = json.Unmarshal(data, &body); err != nil {
				return "dump request failed"
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
			return body
		}
	}
	return "not in json format"
}

func dumpResponse(writer *ResponseWriterWrapper, c *gin.Context) any {
	if strings.Contains(c.Writer.Header().Get(content_type), application_json) {
		if data, err := io.ReadAll(writer.Body); err != nil {
			return "dump response failed"
		} else {
			response := map[string]any{}
			if err := json.Unmarshal(data, &response); err != nil {
				return "dump response failed"
			}
			return response
		}
	}
	return "not in json format"
}
