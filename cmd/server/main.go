package main

import (
	"log"
	"os"
	"strconv"

	"github.com/victorspringer/1g-take-home-task/internal/app"
	"github.com/victorspringer/1g-take-home-task/internal/pkg/repository"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	port, err := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	if err != nil {
		logger.With(zap.Error(err)).Fatal("unable to parse HTTP_PORT env var value")
	}

	connString := "postgres://user:pass@localhost:5432/challenge?sslmode=disable"
	repo, err := repository.New(connString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer repo.Close()

	app.Run(port, logger, repo)
}

func getEnv(name string, defaultValue string) string {
	if value, found := os.LookupEnv(name); found {
		return value
	}
	return defaultValue
}
