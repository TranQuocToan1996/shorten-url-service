package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shorten/handler"
	"shorten/pkg/config"
	"shorten/pkg/db"
	"shorten/pkg/db/migrations"
	"shorten/pkg/queue/redis_stream"
	"shorten/repo"
	"shorten/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

func main() {
	cfg := config.LoadConfig()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	dsn := cfg.DSN()
	// Run migrations
	if err := migrations.RunMigrations(dsn); err != nil {
		log.Printf("Migration failed: %v", err)
	}

	// Initialize database
	database, err := db.NewGormDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.REDIS_PASSWORD,
	})
	defer redisClient.Close()

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("Connected to Redis successfully")
	}

	// Setup routes
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	urlRepo := repo.NewURLRepository(database)
	urlService := service.NewURLService(*cfg, urlRepo, redis_stream.NewRedisStreamProducer(redisClient))
	urlHandler := handler.NewShortenURLHandler(urlService)

	// Swagger
	if cfg.HOST != "" {
		docs.SwaggerInfo.Host = cfg.HOST
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/health", healthCheck)

	// API routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/encode", urlHandler.SubmitEncode)
		v1.GET("/decode", urlHandler.GetDecode)
		// TODO: Implement
		// v1.POST("/decode/:code", urlHandler.SubmitShortURL)
	}

	srv := &http.Server{
		Addr:    cfg.SERVER_PORT,
		Handler: r,
	}

	// Run server in goroutine
	go func() {
		log.Printf("API server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("API server stopped")
}

// healthCheck godoc
// @Summary      Health check
// @Description  Check if the service is running
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
