package service

import (
	"context"
	"fmt"
	"strings"

	"shorten/model"
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
}

func NewURLService(
	config config.Config,
	urlRepo repo.URLRepository,
	producer queue.Producer,
) URLService {
	return &urlService{
		queueName: config.QUEUE_NAME,
		urlRepo:   urlRepo,
		producer:  producer,
		encoder:   NewBase62Encoder(config),
		config:    config,
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
	// TODO: Notify client ok and fail case
	code, err := s.encoder.Encode(msg.URL)
	if err != nil {
		return err
	}
	shortenURL := &model.ShortenURL{
		CleanURL: msg.URL,
		Code:     code,
		Algo:     model.AlgoBase62,
		Status:   model.StatusEncoded,
	}
	err = s.urlRepo.Save(shortenURL)
	if err != nil {
		return err
	}
	return nil
}

func (s *urlService) GetDecode(ctx context.Context, shortenURL string) (*model.ShortenURL, error) {
	code := strings.TrimPrefix(shortenURL, s.config.DB_HOST)
	return s.urlRepo.GetByCode(code)
}
