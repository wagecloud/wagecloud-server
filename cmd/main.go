package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/wagecloud/wagecloud-server/config"
	pgxutil "github.com/wagecloud/wagecloud-server/internal/db/pgx"
	"github.com/wagecloud/wagecloud-server/internal/logger"
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/service"
	"github.com/wagecloud/wagecloud-server/internal/transport/http"
)

const defaultConfigFile = "config/config.dev.yml"
const productionConfigFile = "config/config.production.yml"

var configFile string

func main() {
	setUpConfig()
	setupLogger()
	setupSentry()
	setupHttp()

	// focal-server-cloudimg-amd64.img
}

func setUpConfig() {
	fmt.Println("APP_STAGE", os.Getenv("APP_STAGE"))

	if os.Getenv("APP_STAGE") == "production" {
		configFile = productionConfigFile
	} else {
		configFile = defaultConfigFile
	}

	log.Default().Printf("Using config file: %s", configFile)
	config.SetConfig(configFile)
}

func setupLogger() {
	log.Default().Printf("Using log level: %s", config.GetConfig().Log.Level)
	logger.InitLogger("zap")
}

func setupSentry() {
	logger.Log.Info("Setting up Sentry...")
	err := sentry.Init(sentry.ClientOptions{
		Dsn: config.GetConfig().Sentry.Dsn,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works!")
}

func setupHttp() {
	log.Default().Printf("Setting up HTTP server...")

	pool, err := pgxutil.GetPgxPool()
	if err != nil {
		log.Fatalf("Failed to get pgx pool: %v", err)
	}

	repo := repository.NewRepository(pool)
	service := service.New(repo)
	server := http.NewServer(service, ":8080")
	server.Start()
}
