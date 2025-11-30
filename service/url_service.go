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
	"shorten/pkg/webhook"
	"shorten/repo"
)

type URLService interface {
	SubmitURL(ctx context.Context, longURL string, callbackURL string) error
	HandleShortenURL(ctx context.Context, queueName string, payload []byte) error
	GetDecode(ctx context.Context, shortenURL string) (*model.ShortenURL, error)
	GetByLongURL(ctx context.Context, longURL string) (*model.ShortenURL, error)
}

type UrlService struct {
	queueName string
	urlRepo   repo.URLRepository
	producer  queue.Producer
	encoder   URLEncoder
	config    config.Config
	cache     cache.Cache
	webhook   webhook.Client
}

func NewURLService(
	config config.Config,
	urlRepo repo.URLRepository,
	producer queue.Producer,
	cache cache.Cache,
	encoder URLEncoder,
	webhook webhook.Client,
) URLService {
	return &UrlService{
		queueName: config.QUEUE_NAME,
		urlRepo:   urlRepo,
		producer:  producer,
		encoder:   encoder,
		config:    config,
		cache:     cache,
		webhook:   webhook,
	}
}

func (s *UrlService) SubmitURL(ctx context.Context, longURL string, callbackURL string) error {
	if !url_utils.IsValidURL(longURL) {
		return fmt.Errorf("invalid URL: %s", longURL)
	}
	msg := dto.URLMessage{
		URL:         longURL,
		CallbackURL: callbackURL,
	}
	return s.producer.Publish(ctx, s.queueName, msg.Bytes())
}

func (s *UrlService) HandleShortenURL(ctx context.Context, queueName string, payload []byte) error {
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
		// URL already processed, notify success
		shortURL := fmt.Sprintf("%s/%s", s.config.REDIRECT_HOST, existing.Code)
		s.notifyWebhook(msg.CallbackURL, dto.WebhookPayload{
			Status:   existing.Status,
			LongURL:  msg.URL,
			ShortURL: shortURL,
			Code:     existing.Code,
		})
		return nil
	}

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

	// Notify success
	shortURL := fmt.Sprintf("%s/%s", s.config.REDIRECT_HOST, code)
	s.notifyWebhook(msg.CallbackURL, dto.WebhookPayload{
		Status:   shortenURL.Status,
		LongURL:  msg.URL,
		ShortURL: shortURL,
		Code:     code,
	})

	return nil
}

func (s *UrlService) notifyWebhook(callbackURL string, payload dto.WebhookPayload) {
	if s.webhook == nil || callbackURL == "" {
		return
	}
	go func() {
		if err := s.webhook.Notify(context.Background(), callbackURL, payload); err != nil {
			// Log error but don't fail the main process
			// In production, you might want to use a logger here
			_ = err
		}
	}()
}

func (s *UrlService) GetDecode(ctx context.Context, shortenURL string) (*model.ShortenURL, error) {
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

func (s *UrlService) GetByLongURL(ctx context.Context, longURL string) (*model.ShortenURL, error) {
	return s.urlRepo.GetByLongURL(longURL)
}

func (s *UrlService) UrlCacheKey(code string) string {
	return fmt.Sprintf("url:code:%s", code)
}
