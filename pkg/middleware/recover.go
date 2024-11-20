package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/pkg/errorcode"
	"github.com/zhongxic/sellbot/pkg/model"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("error recovered", slog.Any("error", err))
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					model.FailedWithCode(errorcode.SystemError, "internal server error"))
			}
		}()
		c.Next()
	}
}
