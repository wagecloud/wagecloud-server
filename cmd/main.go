package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wagecloud/wagecloud-server/config"
	"github.com/wagecloud/wagecloud-server/gen/pb/os/v1/osv1connect"
	"github.com/wagecloud/wagecloud-server/internal/client/libvirt"
	"github.com/wagecloud/wagecloud-server/internal/client/nats"
	"github.com/wagecloud/wagecloud-server/internal/client/pgxpool"
	"github.com/wagecloud/wagecloud-server/internal/client/redis"
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
	paymentsvc "github.com/wagecloud/wagecloud-server/internal/modules/payment/service"
	paymentstorage "github.com/wagecloud/wagecloud-server/internal/modules/payment/storage"
	paymentecho "github.com/wagecloud/wagecloud-server/internal/modules/payment/transport/echo"
	echovalidator "github.com/wagecloud/wagecloud-server/internal/shared/transport/http/validator"
	"github.com/wagecloud/wagecloud-server/internal/utils/net"
	"github.com/wagecloud/wagecloud-server/internal/utils/ptr"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var targetService = flag.String("service", "", "Which service to run")

const defaultConfigFile = "config/config.dev.yml"
const productionConfigFile = "config/config.production.yml"

var configFile string

func main() {
	flag.Parse()

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

	// pgpoolread, err := pgxpool.NewPgxpool(pgxpool.PgxpoolOptions{
	// 	Url:             config.GetConfig().Postgres.Url,
	// 	Host:            config.GetConfig().Postgres.Host,
	// 	Port:            config.GetConfig().Postgres.Port,
	// 	Username:        config.GetConfig().Postgres.Username,
	// 	Password:        config.GetConfig().Postgres.Password,
	// 	Database:        config.GetConfig().Postgres.Database,
	// 	MaxConnections:  config.GetConfig().Postgres.MaxConnections,
	// 	MaxConnIdleTime: config.GetConfig().Postgres.MaxConnIdleTime,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to get pgx pool: %v", err)
	// }

	// pgpoolwrite, err := pgxpool.NewPgxpool(pgxpool.PgxpoolOptions{
	// 	Url:             config.GetConfig().PostgresWrite.Url,
	// 	Host:            config.GetConfig().PostgresWrite.Host,
	// 	Port:            config.GetConfig().PostgresWrite.Port,
	// 	Username:        config.GetConfig().PostgresWrite.Username,
	// 	Password:        config.GetConfig().PostgresWrite.Password,
	// 	Database:        config.GetConfig().PostgresWrite.Database,
	// 	MaxConnections:  config.GetConfig().PostgresWrite.MaxConnections,
	// 	MaxConnIdleTime: config.GetConfig().PostgresWrite.MaxConnIdleTime,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to get pgx pool for write: %v", err)
	// }

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
		log.Fatalf("Failed to get pgx pool for write: %v", err)
	}

	e := echo.New()

	e.Pre(middleware.AddTrailingSlash())
	// e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))
	// e.Use(middleware.Logger())
	e.Validator = echovalidator.NewCustomValidator()

	api := e.Group("/api")
	v1 := api.Group("/v1")

	natsClient, err := nats.NewClient(nats.NATSConfig{
		URL:     config.GetConfig().Nats.Url,
		Timeout: config.GetConfig().Nats.Timeout,
	})
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	redisClient, err := redis.NewClient(redis.RedisConfig{
		Addr:     config.GetConfig().Redis.Addr,
		Password: config.GetConfig().Redis.Password,
		DB:       config.GetConfig().Redis.DB,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	svcCtx := serviceContext{
		db:            pgpool,
		e:             v1,
		targetService: ptr.DerefDefault(targetService, ""),
		httpClient:    &http.Client{},
		mux:           &http.ServeMux{},
		nats:          natsClient,
		redis:         redisClient,
	}

	setupServiceAccount(svcCtx)
	osSvc := setupServiceOS(svcCtx)
	paymentSvc := setupServicePayment(svcCtx)
	setupServiceInstance(svcCtx, osSvc.svc, paymentSvc.svc)

	// Print the api routes
	for _, route := range e.Routes() {
		logger.Log.Info("", zap.String("method", route.Method), zap.String("path", route.Path))
	}

	// Setup 404 handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			c.JSON(he.Code, map[string]string{"error": he.Message.(string)})
		} else {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
	}

	if *targetService != "" {
		// Start the connect/gRPC server
		go func() {
			port, err := net.FindNextAvailablePort(50051)
			if err != nil {
				log.Fatalf("Failed to find available port: %v", err)
			}

			if err = http.ListenAndServe(
				fmt.Sprintf(":%d", port),
				h2c.NewHandler(svcCtx.mux, &http2.Server{}),
			); err != nil {
				log.Fatalf("failed to start server: %v", err)
			}
		}()
	}

	// Start the HTTP server
	go func() {
		port, err := net.FindNextAvailablePort(3000)
		if err != nil {
			log.Fatalf("Failed to find available port: %v", err)
		}

		port = 3000

		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	select {}

}

type serviceContext struct {
	db            pgxpool.DBTX
	e             *echo.Group
	targetService string
	httpClient    *http.Client
	mux           *http.ServeMux
	nats          nats.Client
	redis         redis.Client
}

type service[T any] struct {
	svc   T
	isRPC bool
}

func setupServiceAccount(svcCtx serviceContext) service[accountsvc.Service] {
	var accountSvc accountsvc.Service

	isRPC := svcCtx.targetService != "" && svcCtx.targetService != "account"

	// if isRPC {
	// 	connectClient := accountv1connect.NewAccountServiceClient(
	// 		svcCtx.httpClient,
	// 		"localhost:50051",
	// 		connect.WithGRPC(),
	// 	)
	// 	accountSvc = accountsvc.NewServiceRpc(connectClient)
	// } else {
	accountSvc = accountsvc.NewService(accountstorage.NewStorage(svcCtx.db))
	accountHandler := accountecho.NewEchoHandler(accountSvc)
	// path, handler := accountconnect.NewAccountServiceHandler(accountSvc)
	// svcCtx.mux.Handle(path, handler)

	account := svcCtx.e.Group("/account")
	account.PATCH("/", accountHandler.UpdateAccount)

	user := account.Group("/user")
	user.PATCH("/", accountHandler.UpdateUser)
	user.GET("/", accountHandler.GetUser)
	user.POST("/login/", accountHandler.LoginUser)
	user.POST("/register/", accountHandler.RegisterUser)
	// }

	return service[accountsvc.Service]{
		svc:   accountSvc,
		isRPC: isRPC,
	}
}

func setupServiceOS(svcCtx serviceContext) service[ossvc.Service] {
	var osSvc ossvc.Service

	isRPC := svcCtx.targetService != "" && svcCtx.targetService != "os"

	if isRPC {
		connectClient := osv1connect.NewOSServiceClient(
			svcCtx.httpClient,
			"localhost:50051",
			connect.WithGRPC(),
		)
		osSvc = ossvc.NewServiceRpc(connectClient)
	} else {
		osSvc = ossvc.NewService(osstorage.NewStorage(svcCtx.db))
		osHandler := osecho.NewEchoHandler(osSvc)

		os := svcCtx.e.Group("/os")
		os.GET("/", osHandler.ListOSs)
		os.GET("/:id", osHandler.GetOS)
		os.POST("/", osHandler.CreateOS)
		os.PATCH("/:id", osHandler.UpdateOS)
		os.DELETE("/:id", osHandler.DeleteOS)

		arch := os.Group("/arch")
		arch.GET("/", osHandler.ListArchs)
		arch.GET("/:id", osHandler.GetArch)
		arch.POST("/", osHandler.CreateArch)
		arch.PATCH("/:id", osHandler.UpdateArch)
		arch.DELETE("/:id", osHandler.DeleteArch)
	}

	return service[ossvc.Service]{
		svc:   osSvc,
		isRPC: isRPC,
	}
}

func setupServiceInstance(svcCtx serviceContext, osSvc ossvc.Service, paymentSvc paymentsvc.Service) service[instancesvc.Service] {
	var instanceSvc instancesvc.Service

	isRPC := svcCtx.targetService != "" && svcCtx.targetService != "instance"

	// if isRPC {
	// 	connectClient := instancev1connect.NewInstanceServiceClient(
	// 		svcCtx.httpClient,
	// 		"localhost:50051",
	// 		connect.WithGRPC(),
	// 	)
	// 	instanceSvc = instancesvc.NewServiceRpc(connectClient)
	// } else {
	libvirt := libvirt.NewClient()
	instanceSvc = instancesvc.NewService(
		libvirt,
		svcCtx.nats,
		svcCtx.redis,
		instancestorage.NewStorage(svcCtx.db),
		osSvc,
		paymentSvc,
	)
	instanceHandler := instanceecho.NewEchoHandler(instanceSvc)

	region := svcCtx.e.Group("/region")
	region.GET("/", instanceHandler.ListRegions)
	region.GET("/:id", instanceHandler.GetRegion)
	region.POST("/", instanceHandler.CreateRegion)
	region.PATCH("/:id", instanceHandler.UpdateRegion)
	region.DELETE("/:id", instanceHandler.DeleteRegion)

	instance := svcCtx.e.Group("/instance")
	instance.GET("/", instanceHandler.ListInstances)
	instance.GET("/:id/", instanceHandler.GetInstance)
	instance.GET("/:id/monitor/", instanceHandler.GetInstanceMonitor)
	instance.POST("/", instanceHandler.CreateInstance)
	instance.POST("/start/:id/", instanceHandler.StartInstance)
	instance.POST("/stop/:id/", instanceHandler.StopInstance)
	instance.PATCH("/:id", instanceHandler.UpdateInstance)
	instance.DELETE("/:id", instanceHandler.DeleteInstance)

	log := instance.Group("/log")
	log.GET("/:id", instanceHandler.GetInstanceLog)
	log.GET("/", instanceHandler.ListInstanceLogs)
	log.POST("/", instanceHandler.CreateInstanceLog)
	log.PATCH("/:id", instanceHandler.UpdateInstanceLog)
	log.DELETE("/:id", instanceHandler.DeleteInstanceLog)

	network := svcCtx.e.Group("/network")
	network.GET("/list/", instanceHandler.ListNetworks)
	network.GET("/", instanceHandler.GetNetwork)
	network.POST("/map/", instanceHandler.MapPortNginx)
	network.POST("/unmap/", instanceHandler.UnmapPortNginx)

	network.GET("/domain/", instanceHandler.ListDomains)
	network.GET("/domain/:id", instanceHandler.GetDomain)
	network.POST("/domain/", instanceHandler.CreateDomain)
	network.PATCH("/domain/:id", instanceHandler.UpdateDomain)
	network.DELETE("/domain/:id", instanceHandler.DeleteDomain)

	// }

	return service[instancesvc.Service]{
		svc:   instanceSvc,
		isRPC: isRPC,
	}
}

func setupServicePayment(svcCtx serviceContext) service[paymentsvc.Service] {
	var paymentSvc paymentsvc.Service

	isRPC := svcCtx.targetService != "" && svcCtx.targetService != "payment"

	if isRPC {
		// connectClient := paymentv1connect.NewPaymentServiceClient(
		// 	svcCtx.httpClient,
		// 	"localhost:50051",
		// 	connect.WithGRPC(),
		// )
		// paymentSvc = paymentsvc.NewServiceRpc(connectClient)
	} else {
		paymentSvc = paymentsvc.NewService(paymentstorage.NewStorage(svcCtx.db), svcCtx.nats)
		paymentHandler := paymentecho.NewEchoHandler(paymentSvc)

		paymentHandler.RegisterRoutes(svcCtx.e)

	}

	return service[paymentsvc.Service]{
		svc:   paymentSvc,
		isRPC: isRPC,
	}
}
