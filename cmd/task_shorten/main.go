package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shorten/pkg/cache/redis_cache"
	"shorten/pkg/config"
	"shorten/pkg/db"
	"shorten/pkg/queue/redis_stream"
	"shorten/repo"
	"shorten/service"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Initialize database
	database, err := db.NewGormDB(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.REDIS_PASSWORD,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("Connected to Redis successfully")
	}
	defer redisClient.Close()

	// Initialize dependencies
	urlRepo := repo.NewURLRepository(database)
	cache := redis_cache.NewRedisCache(redisClient)
	urlService := service.NewURLService(*cfg, urlRepo, nil, cache, service.NewBase62Encoder(*cfg))

	consumerName := fmt.Sprintf("shorten-url-worker-%s", uuid.NewString())
	log.Printf("Consumer shorten_URL start: %v", consumerName)

	// Create consumer with handler
	consumer, err := redis_stream.NewRedisStreamConsumer(
		redisClient,
		urlService.HandleShortenURL,
		redis_stream.WithConsumerGroup("shorten-url-group"),
		redis_stream.WithConsumerName(consumerName),
		redis_stream.WithEnsureGroup(),
	)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Start consumer in goroutine
	consumerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Printf("Starting consumer for queue: %s", cfg.QUEUE_NAME)
		if err := consumer.Consume(consumerCtx, cfg.QUEUE_NAME, nil); err != nil {
			log.Printf("Consume err: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Wait for interrupt signal
	<-quit
	log.Println("Received shutdown signal, starting graceful shutdown...")

	// Wait for graceful shutdown or timeout
	<-time.After(5 * time.Second)

	log.Println("Worker stopped")
}
