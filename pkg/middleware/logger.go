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
	"github.com/zhongxic/sellbot/pkg/errorcode"
	"github.com/zhongxic/sellbot/pkg/model"
)

const (
	contentType     = "Content-Type"
	applicationJson = "application/json"
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
		body, err := dumpRequest(c)
		if err != nil {
			result := model.FailedWithCode(errorcode.SystemError, "dump request failed")
			c.AbortWithStatusJSON(http.StatusInternalServerError, result)
		}

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

func dumpRequest(c *gin.Context) (body *map[string]any, err error) {
	if c.Request.Header.Get(contentType) == applicationJson {
		if data, err := io.ReadAll(c.Request.Body); err != nil {
			return nil, err
		} else {
			m := &map[string]any{}
			if err := json.Unmarshal(data, m); err != nil {
				return nil, err
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
			return m, nil
		}
	}
	return &map[string]any{}, nil
}

func dumpResponse(writer *ResponseWriterWrapper, c *gin.Context) (body *map[string]any) {
	m := &map[string]any{"addition": "not in json format or dump failed"}
	if strings.Contains(c.Writer.Header().Get(contentType), applicationJson) {
		if data, err := io.ReadAll(writer.Body); err != nil {
			return m
		} else {
			res := &map[string]any{}
			if err := json.Unmarshal(data, res); err != nil {
				return m
			}
			return res
		}
	}
	return m
}
