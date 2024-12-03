package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/pkg/errorcode"
	"github.com/zhongxic/sellbot/pkg/result"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("recovered", "error", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					result.FailedWithErrorCode(errorcode.SystemError, "internal server error"))
			}
		}()
		c.Next()
	}
}
