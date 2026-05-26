package middleware

import (
	"log/slog"
	"time"

	"github.com/5gMurilo/helptrix-api/core/logger"
	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)

		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		attrs := []any{
			slog.String("layer", "middleware"),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.String("latency", latency.String()),
			slog.String("ip", c.ClientIP()),
		}

		if ginErr := c.Errors.ByType(gin.ErrorTypePrivate).String(); ginErr != "" {
			attrs = append(attrs, slog.String("gin_error", ginErr))
		}

		logger.Get().Log(c.Request.Context(), level, "http request", attrs...)
	}
}
