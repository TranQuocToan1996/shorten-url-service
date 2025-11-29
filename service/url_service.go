package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"shorten/model"
	"shorten/pkg/cache"
	"shorten/pkg/config"
	"shorten/pkg/dto"
	"shorten/pkg/queue"
	"shorten/pkg/utils/url_utils"
	"shorten/repo"
)

type URLService interface {
	SubmitURL(ctx context.Context, longURL string) error
	HandleShortenURL(ctx context.Context, queueName string, payload []byte) error
	GetDecode(ctx context.Context, shortenURL string) (*model.ShortenURL, error)
}

type urlService struct {
	queueName string
	urlRepo   repo.URLRepository
	producer  queue.Producer
	encoder   URLEncoder
	config    config.Config
	cache     cache.Cache
}

func NewURLService(
	config config.Config,
	urlRepo repo.URLRepository,
	producer queue.Producer,
	cache cache.Cache,
) URLService {
	return &urlService{
		queueName: config.QUEUE_NAME,
		urlRepo:   urlRepo,
		producer:  producer,
		encoder:   NewBase62Encoder(config),
		config:    config,
		cache:     cache,
	}
}

// TODO: Save submit status -> fail/ok
func (s *urlService) SubmitURL(ctx context.Context, longURL string) error {
	if !url_utils.IsValidURL(longURL) {
		return fmt.Errorf("invalid URL: %s", longURL)
	}
	msg := dto.URLMessage{URL: longURL}
	return s.producer.Publish(ctx, s.queueName, msg.Bytes())
}

func (s *urlService) HandleShortenURL(ctx context.Context, queueName string, payload []byte) error {
	msg := dto.URLMessage{}
	err := msg.Unmarshal(payload)
	if err != nil {
		return err
	}
	if !url_utils.IsValidURL(msg.URL) {
		return fmt.Errorf("invalid URL: %s", msg.URL)
	}
	existing, err := s.urlRepo.GetByLongURL(msg.URL)
	if err == nil && existing != nil {
		// URL already processed, skip
		return nil
	}
	// TODO: Notify client ok and fail case
	code, err := s.encoder.Encode(msg.URL)
	if err != nil {
		return err
	}
	shortenURL := &model.ShortenURL{
		LongURL: msg.URL,
		Code:    code,
		Algo:    model.AlgoBase62,
		Status:  model.StatusEncoded,
	}
	err = s.urlRepo.Save(shortenURL)
	if err != nil {
		return err
	}

	return nil
}

func (s *urlService) GetDecode(ctx context.Context, shortenURL string) (*model.ShortenURL, error) {
	code := strings.TrimPrefix(shortenURL, s.config.REDIRECT_HOST)

	// Try to get from cache first
	cacheKey := s.UrlCacheKey(code)
	cachedData, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		// Cache hit - unmarshal and return
		var shortenURL model.ShortenURL
		if err := json.Unmarshal(cachedData, &shortenURL); err == nil {
			return &shortenURL, nil
		}
		// If unmarshal fails, continue to database lookup
	}

	// Cache miss or error - get from database
	result, err := s.urlRepo.GetByCode(code)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests (cache-aside pattern)
	if result != nil {
		if data, err := json.Marshal(result); err == nil {
			go s.cache.Set(context.Background(), cacheKey, data, time.Hour)
		}
	}

	return result, nil
}

func (s *urlService) UrlCacheKey(code string) string {
	return fmt.Sprintf("url:code:%s", code)
}
