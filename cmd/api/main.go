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
	"shorten/pkg/cache/redis_cache"
	"shorten/pkg/config"
	"shorten/pkg/db"
	"shorten/pkg/db/migrations"
	"shorten/pkg/queue/redis_stream"
	"shorten/pkg/webhook"
	"shorten/repo"
	"shorten/service"

	"github.com/gin-contrib/cors"
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

	// Configure CORS - flexible for dev and production
	corsOrigins := cfg.CORSOrigins()
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// If origins contains "*", allow all origins (development mode)
	if len(corsOrigins) == 1 && corsOrigins[0] == "*" {
		corsConfig.AllowAllOrigins = true
		log.Println("CORS: Allowing all origins (development mode)")
	} else {
		corsConfig.AllowOrigins = corsOrigins
		log.Printf("CORS: Allowing origins: %v", corsOrigins)
	}

	r.Use(cors.New(corsConfig))

	urlRepo := repo.NewURLRepository(database)
	cache := redis_cache.NewRedisCache(redisClient)
	webhookClient := webhook.NewHTTPWebhookClient()
	urlService := service.NewURLService(*cfg, urlRepo, redis_stream.NewRedisStreamProducer(redisClient), cache, service.NewBase62Encoder(*cfg), webhookClient)
	urlHandler := handler.NewShortenURLHandler(urlService, *cfg)

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
		v1.GET("/urls/long", urlHandler.GetURLEncodeByLongURL)
		v1.GET("/:code", urlHandler.RedirectLongURL)
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
