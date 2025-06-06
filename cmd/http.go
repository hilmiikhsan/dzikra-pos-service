package cmd

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/route"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/rs/zerolog/log"
)

func RunServerHTTP(cmd *flag.FlagSet, args []string) {
	var (
		envs        = config.Envs
		flagAppPort = cmd.String("port", envs.App.Port, "Application port")
		SERVER_PORT string
	)

	// logLevel, err := zerolog.ParseLevel(envs.App.LogLevel)
	// if err != nil {
	// 	logLevel = zerolog.InfoLevel
	// }

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing flags")
	}

	if envs.App.Port != "" {
		SERVER_PORT = envs.App.Port
	} else {
		SERVER_PORT = *flagAppPort
	}

	app := fiber.New()

	// Application Middlewares
	if envs.App.Environtment == constants.EnvProduction {
		app.Use(limiter.New(limiter.Config{
			Max:        50,
			Expiration: 30 * time.Second,
		}))
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS,HEAD",
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
	}))
	// End Application Middlewares

	minioOpt, err := adapter.WithDzikraMinio()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Dzikra Minio")
	}

	adapter.Adapters.Sync(
		adapter.WithRestServer(app),
		adapter.WithDzikraPostgres(),
		adapter.WithDzikraRedis(),
		minioOpt,
		adapter.WithValidator(validator.NewValidator()),
	)

	// infrastructure.InitializeLogger(
	// 	envs.App.Environtment,
	// 	envs.App.LogFile,
	// 	logLevel,
	// )

	app.Get("/metrics", monitor.New(monitor.Config{Title: config.Envs.App.Name + config.Envs.App.Environtment + " Metrics"}))
	route.SetupRoutes(app)

	// print all routes that are registered
	// for _, route := range app.Stack() {
	// 	for _, handler := range route {
	// 		fmt.Printf("Method: %s, Path: %s\n", handler.Method, handler.Path)
	// 	}
	// }

	// Run server in goroutine
	go func() {
		log.Info().Msgf("Server is running on port %s", SERVER_PORT)
		if err := app.Listen(":" + SERVER_PORT); err != nil {
			log.Fatal().Msgf("Error while starting server: %v", err)
		}
	}()
	// End Run server in goroutine

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)

	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	if runtime.GOOS == "windows" {
		shutdownSignals = []os.Signal{os.Interrupt}
	}

	signal.Notify(quit, shutdownSignals...)
	<-quit
	log.Info().Msg("Server is shutting down ...")

	err = adapter.Adapters.Unsync()
	if err != nil {
		log.Error().Msgf("Error while closing adapters: %v", err)
	}

	log.Info().Msg("Server gracefully stopped")
}
