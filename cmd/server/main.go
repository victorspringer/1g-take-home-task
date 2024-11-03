package main

import (
	"os"
	"strconv"

	"github.com/victorspringer/1g-take-home-task/internal/app"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	port, err := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	if err != nil {
		logger.With(zap.Error(err)).Fatal("unable to parse HTTP_PORT env var value")
	}

	app.Run(port, logger)
}

func getEnv(name string, defaultValue string) string {
	if value, found := os.LookupEnv(name); found {
		return value
	}
	return defaultValue
}
