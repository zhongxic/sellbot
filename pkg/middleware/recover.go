package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/pkg/model"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("error recovered", slog.Any("error", err))
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					model.FailedWithMessage(model.SYSTEM_ERROR, "internal server error"))
			}
		}()
		c.Next()
	}
}
