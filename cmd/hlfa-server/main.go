package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/makarivmri12/bug/internal/api"
	"github.com/makarivmri12/bug/internal/engine"
	"github.com/makarivmri12/bug/internal/repository"
	"github.com/makarivmri12/bug/internal/scanner"
	"github.com/makarivmri12/bug/pkg/config"
)

func main() {
	// Parse command-line flags
	cfgPath := flag.String("config", "configs/config.yaml", "Path to configuration file")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	// Initialize logger
	logger, err := initializeLogger(*logLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("HLFA Server starting", zap.String("version", "1.0.0"))

	// Load configuration
	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	logger.Info("Configuration loaded", zap.String("config_path", *cfgPath))

	// Initialize database connection
	db, err := repository.InitPostgres(cfg.Database.DSN, logger)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer repository.Close(db, logger)

	// Initialize Redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	})
	defer redisClient.Close()

	// Test Redis connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	logger.Info("Redis connected successfully")

	// Initialize scan engine
	scanEngine := engine.NewScanEngine(cfg.Engine.Workers, redisClient, logger)

	// Register scanners
	scanEngine.RegisterScanner(scanner.NewWebScanner(logger))
	scanEngine.RegisterScanner(scanner.NewNetworkScanner(logger))

	logger.Info("Scanners registered", zap.Int("count", 2))

	// Start scan engine
	if err := scanEngine.Start(context.Background()); err != nil {
		logger.Fatal("Failed to start scan engine", zap.Error(err))
	}

	// Initialize API server
	apiServer := api.NewServer(logger, scanEngine)

	// Start API server in goroutine
	go func() {
		logger.Info("API server starting",
			zap.String("host", cfg.API.Host),
			zap.Int("port", cfg.API.Port),
		)

		addr := fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port)
		if err := startAPIServer(addr, apiServer, logger); err != nil {
			logger.Error("API server error", zap.Error(err))
		}
	}()

	// Setup graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	<-shutdownChan
	logger.Info("Shutdown signal received")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	scanEngine.Stop()
	logger.Info("HLFA Server stopped")
}

// initializeLogger initializes the zap logger
func initializeLogger(level string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	return cfg.Build()
}

// startAPIServer starts the HTTP API server
func startAPIServer(addr string, server *api.Server, logger *zap.Logger) error {
	import "net/http"

	http := &http.Server{
		Addr:    addr,
		Handler: server,
	}

	return http.ListenAndServe()
}
