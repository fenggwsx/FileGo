package middleware

import (
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RecoveryMiddleware struct {
	Logger *zap.Logger
}

func NewRecoveryMiddleware(logger *zap.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		Logger: logger,
	}
}

func (m *RecoveryMiddleware) GetMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it's not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") ||
							strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// Create the request header.
				httpRequest, _ := httputil.DumpRequest(ctx.Request, false)

				// Split by lines for further processing.
				headers := strings.Split(string(httpRequest), "\r\n")

				// Mask the sensitive data.
				for idx, header := range headers {
					if idx > 0 {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" || current[0] == "Cookie" {
							headers[idx] = current[0] + ": *"
						}
					}
				}

				// Join each line together to restore the header string.
				headersToStr := strings.Join(headers, "\r\n")

				if brokenPipe {
					// Log a warning since it's not really a fatal error.
					m.Logger.Warn("broken connection",
						zap.String("headers", headersToStr),
						zap.Any("error", err),
					)

					// If the connection is dead, we can't write a status to it.
					ctx.Error(err.(error)) //nolint: errcheck
					ctx.Abort()
				} else {
					// Log an error for bug fix in the future.
					m.Logger.Error("recovery from panic",
						zap.String("headers", headersToStr),
						zap.Any("error", err),
					)

					// Set the response status to `Internal Server Error`.
					ctx.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		ctx.Next()
	}
}
