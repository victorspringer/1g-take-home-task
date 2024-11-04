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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/victorspringer/1g-take-home-task/docs"
	"github.com/victorspringer/1g-take-home-task/internal/pkg/device"
	"go.uber.org/zap"
)

// @title Devices Service
// @version 1.0
// @description Devices Service for technical challenge.
// @contact.name Victor Springer
// @license.name MIT License
// @host localhost:8080
// @BasePath /

// Run starts the HTTP server on specified port.
func Run(port int, logger *zap.Logger, deviceRepository device.Repository) {
	handler := &handler{
		logger:           logger,
		deviceRepository: deviceRepository,
	}

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/", handler.healthCheck)
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	devices := router.Group("devices")

	devices.GET("/", handler.listAllDevices)
	devices.GET("/:id", handler.getDeviceByID)
	devices.GET("/search", handler.searchDevices)

	devices.POST("/", handler.addDevice)

	devices.PATCH("/:id", handler.updateDevice)

	devices.DELETE("/:id", handler.deleteDevice)

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
