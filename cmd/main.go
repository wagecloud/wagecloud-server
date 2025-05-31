package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/wagecloud/wagecloud-server/config"
	"github.com/wagecloud/wagecloud-server/internal/client/libvirt"
	"github.com/wagecloud/wagecloud-server/internal/client/pgxpool"
	"github.com/wagecloud/wagecloud-server/internal/logger"
	accountsvc "github.com/wagecloud/wagecloud-server/internal/modules/account/service"
	accountstorage "github.com/wagecloud/wagecloud-server/internal/modules/account/storage"
	accountecho "github.com/wagecloud/wagecloud-server/internal/modules/account/transport/echo"
	instancesvc "github.com/wagecloud/wagecloud-server/internal/modules/instance/service"
	instancestorage "github.com/wagecloud/wagecloud-server/internal/modules/instance/storage"
	instanceecho "github.com/wagecloud/wagecloud-server/internal/modules/instance/transport/echo"
	ossvc "github.com/wagecloud/wagecloud-server/internal/modules/os/service"
	osstorage "github.com/wagecloud/wagecloud-server/internal/modules/os/storage"
	osecho "github.com/wagecloud/wagecloud-server/internal/modules/os/transport/echo"
)

const defaultConfigFile = "config/config.dev.yml"
const productionConfigFile = "config/config.production.yml"

var configFile string

func main() {
	setUpConfig()
	setupLogger()
	setupSentry()
	setupModules()
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
		Dsn:              config.GetConfig().Sentry.Dsn,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works!")
}

func setupModules() {
	logger.Log.Info("Setting up modules...")

	pgpool, err := pgxpool.NewPgxpool(pgxpool.PgxpoolOptions{
		Url:             config.GetConfig().Postgres.Url,
		Host:            config.GetConfig().Postgres.Host,
		Port:            config.GetConfig().Postgres.Port,
		Username:        config.GetConfig().Postgres.Username,
		Password:        config.GetConfig().Postgres.Password,
		Database:        config.GetConfig().Postgres.Database,
		MaxConnections:  config.GetConfig().Postgres.MaxConnections,
		MaxConnIdleTime: config.GetConfig().Postgres.MaxConnIdleTime,
	})
	if err != nil {
		log.Fatalf("Failed to get pgx pool: %v", err)
	}

	libvirt := libvirt.NewClient()

	accountHandler := accountecho.NewEchoHandler(accountsvc.NewService(accountstorage.NewStorage(pgpool)))
	instanceHandler := instanceecho.NewEchoHandler(instancesvc.NewService(libvirt, instancestorage.NewStorage(pgpool)))
	osHandler := osecho.NewEchoHandler(ossvc.NewService(osstorage.NewStorage(pgpool)))

	e := echo.New()

	api := e.Group("/api")
	v1 := api.Group("/v1")

	// Account
	account := v1.Group("/account")
	account.GET("/", accountHandler.GetAccount)
	account.POST("/login", accountHandler.LoginUser)
	account.POST("/register", accountHandler.RegisterUser)

	// Instance
	instance := v1.Group("/instance")
	instance.GET("/", instanceHandler.ListInstances)
	instance.GET("/:id", instanceHandler.GetInstance)
	instance.POST("/", instanceHandler.CreateInstance)
	instance.PUT("/:id", instanceHandler.UpdateInstance)
	instance.DELETE("/:id", instanceHandler.DeleteInstance)

	// OS
	os := v1.Group("/os")
	os.GET("/", osHandler.ListOSs)
	os.GET("/:id", osHandler.GetOS)
	os.POST("/", osHandler.CreateOS)
	os.PUT("/:id", osHandler.UpdateOS)
	os.DELETE("/:id", osHandler.DeleteOS)
	// Start the server
	if err := e.Start(fmt.Sprintf(":%d", 3000)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Setup CORS middleware
	e.Use(echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}
			return next(c)
		}
	}))

	// Setup 404 handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			c.JSON(he.Code, map[string]string{"error": he.Message.(string)})
		} else {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
	}
}
