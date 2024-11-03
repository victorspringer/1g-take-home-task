package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Run starts the HTTP server on specified port.
func Run(port int, logger *zap.Logger) {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		logger.Info("test", zap.String("requestUrl", c.Request.URL.Path))
		c.JSON(http.StatusOK, gin.H{
			"requested": c.Request.URL.Path,
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
		<-sigChan

		// Received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.With(zap.Error(err)).Error("error while attempting graceful shutdown")
		}

		close(idleConnsClosed)
	}()

	logger.With(zap.Int("port", port)).Info("starting http server")

	if err := srv.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		logger.Info("server closed")
	} else {
		logger.With(zap.Error(err)).Error("http server start/close error")
	}

	<-idleConnsClosed
}
