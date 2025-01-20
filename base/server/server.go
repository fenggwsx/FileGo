package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/fenggwsx/filego/base/conf"
	"github.com/fenggwsx/filego/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HttpServer struct {
	Addr   string
	Config conf.ServerConfig
	Logger *zap.Logger
	Engine *gin.Engine
}

func NewHttpServer(
	config *conf.Config,
	logger *zap.Logger,
	loggerMiddleware *middleware.LoggerMiddleware,
	recoveryMiddleware *middleware.RecoveryMiddleware,
) *HttpServer {
	engine := gin.New()

	engine.Use(loggerMiddleware.GetMiddleware())
	engine.Use(recoveryMiddleware.GetMiddleware())

	return &HttpServer{
		Addr:   fmt.Sprintf(":%d", config.Port),
		Config: config.Server,
		Logger: logger,
		Engine: engine,
	}
}

func (m *HttpServer) Run() {
	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:              m.Addr,
		Handler:           m.Engine.Handler(),
		ReadTimeout:       m.Config.ReadTimeout,
		ReadHeaderTimeout: m.Config.ReadHeaderTimeout,
		WriteTimeout:      m.Config.WriteTimeout,
		IdleTimeout:       m.Config.IdleTimeout,
		MaxHeaderBytes:    m.Config.MaxHeaderBytes,
	}

	// Create a pending channel that synchronizes the goroutines.
	pending := make(chan struct{})

	// Initialize the server in a goroutine so that
	// it won't block the graceful shutdown handling below.
	go func() {
		ln, err := net.Listen("tcp", server.Addr)
		if err != nil {
			m.Logger.Panic("failed to start server", zap.Error(err))
		}

		m.Logger.Debug(fmt.Sprintf("begin to listen on %v", m.Addr))

		if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			m.Logger.Error("an error occured during serving", zap.Error(err))
		}

		// Close the channel to cancel the pending state of the main goroutine.
		close(pending)
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	m.Logger.Debug("begin to shutdown server")

	// The context is used to inform the server it has a specified period of time
	// to finish the request it is currently handling.
	ctx, cancel := context.WithTimeout(context.Background(), m.Config.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); errors.Is(err, context.DeadlineExceeded) {
		m.Logger.Error("shutdown server timeout", zap.Error(err))
	} else {
		if err != nil {
			m.Logger.Error("shutdown server with an error", zap.Error(err))
		}
		m.Logger.Debug("shutdown server successfully")

		// Keep pending until the goroutine exits.
		<-pending
	}
}
