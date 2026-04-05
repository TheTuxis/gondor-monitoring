package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/TheTuxis/gondor-monitoring/internal/config"
	"github.com/TheTuxis/gondor-monitoring/internal/handler"
	"github.com/TheTuxis/gondor-monitoring/internal/middleware"
	"github.com/TheTuxis/gondor-monitoring/internal/model"
	jwtpkg "github.com/TheTuxis/gondor-monitoring/internal/pkg/jwt"
	"github.com/TheTuxis/gondor-monitoring/internal/repository"
	"github.com/TheTuxis/gondor-monitoring/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Init logger
	var logger *zap.Logger
	var err error
	if cfg.Environment == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = logger.Sync() }()

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("failed to get underlying sql.DB", zap.Error(err))
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Auto-migrate models
	if err := db.AutoMigrate(
		&model.AlertRule{},
		&model.Alert{},
		&model.AuditLog{},
		&model.ServiceStatus{},
	); err != nil {
		logger.Fatal("failed to auto-migrate", zap.Error(err))
	}
	logger.Info("database migration completed")

	// Init Redis client
	var redisClient *redis.Client
	if cfg.RedisURL != "" {
		opts, err := redis.ParseURL("redis://" + cfg.RedisURL)
		if err != nil {
			// Fallback: treat as host:port
			opts = &redis.Options{Addr: cfg.RedisURL}
		}
		redisClient = redis.NewClient(opts)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := redisClient.Ping(ctx).Err(); err != nil {
			logger.Warn("redis connection failed, continuing without redis", zap.Error(err))
			redisClient = nil
		} else {
			logger.Info("connected to Redis")
		}
	}

	// Init JWT manager (validate-only — tokens are issued by gondor-users-security)
	jwtManager := jwtpkg.NewManager(cfg.JWTSecret)

	// Init repositories
	alertRuleRepo := repository.NewAlertRuleRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)
	serviceStatusRepo := repository.NewServiceStatusRepository(db)

	// Init services
	alertRuleService := service.NewAlertRuleService(alertRuleRepo, logger)
	alertService := service.NewAlertService(alertRepo, logger)
	auditLogService := service.NewAuditLogService(auditLogRepo, logger)
	serviceStatusService := service.NewServiceStatusService(serviceStatusRepo, logger)

	// Init handlers
	healthHandler := handler.NewHealthHandler(db, redisClient)
	alertRuleHandler := handler.NewAlertRuleHandler(alertRuleService)
	alertHandler := handler.NewAlertHandler(alertService)
	auditLogHandler := handler.NewAuditLogHandler(auditLogService)
	serviceStatusHandler := handler.NewServiceStatusHandler(serviceStatusService)

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.AuthMiddleware(jwtManager))

	// Health & metrics (no auth required — handled by skip list)
	router.GET("/health", healthHandler.Health)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Monitoring routes
	v1 := router.Group("/v1/monitoring")
	{
		// Alert rules
		v1.GET("/alerts/rules", alertRuleHandler.List)
		v1.POST("/alerts/rules", alertRuleHandler.Create)
		v1.GET("/alerts/rules/:id", alertRuleHandler.GetByID)
		v1.PUT("/alerts/rules/:id", alertRuleHandler.Update)
		v1.DELETE("/alerts/rules/:id", alertRuleHandler.Delete)

		// Alerts
		v1.GET("/alerts", alertHandler.List)
		v1.POST("/alerts/:id/acknowledge", alertHandler.Acknowledge)

		// Audit logs
		v1.GET("/audit-logs", auditLogHandler.List)
		v1.POST("/audit-logs", auditLogHandler.Create)

		// Service status
		v1.GET("/services/status", serviceStatusHandler.List)
		v1.POST("/services/status", serviceStatusHandler.UpdateStatus)
	}

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("starting server", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	if redisClient != nil {
		_ = redisClient.Close()
	}
	_ = sqlDB.Close()

	logger.Info("server stopped")
}
