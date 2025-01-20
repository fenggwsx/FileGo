package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoggerMiddleware struct {
	Logger *zap.Logger
}

func NewLoggerMiddleware(logger *zap.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		Logger: logger,
	}
}

func (m *LoggerMiddleware) GetMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Start the timer for latency calculation.
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		// Join the non-empty query string to the path.
		if raw != "" {
			path = path + "?" + raw
		}

		// Process the incoming request.
		ctx.Next()

		// Log the request and processing information.
		m.Logger.Info("incoming request",
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.Int("status", ctx.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", ctx.ClientIP()),
			zap.String("user-agent", ctx.Request.UserAgent()),
			zap.Strings("errors", ctx.Errors.ByType(gin.ErrorTypePrivate).Errors()),
		)
	}
}
