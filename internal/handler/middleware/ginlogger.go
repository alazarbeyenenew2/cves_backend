// internal/handler/middleware/ginlogger.go
package middleware

import (
	"context"
	"time"

	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GinLogger(log logger.Logger) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.Query()
		requestID := ctx.GetHeader("X-Request-Id")
		channel := "mobile"
		ctx.Set("X-Request-Id", requestID)
		ctx.Set("x-start-time", start)
		ctx.Set("user-agent", ctx.Request.UserAgent())
		ctx.Set("ip", ctx.ClientIP())
		ctx.Set("X-Channel-Id", channel)

		ctx.Set("channel", ctx.GetHeader("X-Channel-Id"))

		ctx.Next()
		end := time.Now()
		latency := end.Sub(start)
		fields := []zapcore.Field{
			zap.Int("status", ctx.Writer.Status()),
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.String("ip", ctx.ClientIP()),
			zap.Any("query", query),
			zap.String("user-agent", ctx.Request.UserAgent()),
			zap.Float64("latency", latency.Minutes()),
		}
		log.Info(context.Background(), "Gin Request", fields...)
	}
}
