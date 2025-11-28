package service

import (
	"context"
	"fmt"

	"shorten/model"
	"shorten/pkg/config"
	"shorten/pkg/dto"
	"shorten/pkg/queue"
	"shorten/pkg/utils/url_utils"
	"shorten/repo"
)

type URLService interface {
	SubmitURL(ctx context.Context, url string) error
	HandleShortenURL(ctx context.Context, queueName string, payload []byte) error
}

type urlService struct {
	queueName string
	urlRepo   repo.URLRepository
	queue     queue.Producer
	encoder   URLEncoder
}

func NewURLService(
	config config.Config,
	urlRepo repo.URLRepository,
	queue queue.Producer,
) URLService {
	return &urlService{
		queueName: config.QUEUE_NAME,
		urlRepo:   urlRepo,
		queue:     queue,
		encoder:   NewBase62Encoder(config),
	}
}

// TODO: Save submit status -> fail/ok
func (s *urlService) SubmitURL(ctx context.Context, url string) error {
	if !url_utils.IsValidURL(url) {
		return fmt.Errorf("invalid URL: %s", url)
	}
	msg := dto.URLMessage{URL: url}
	return s.queue.Publish(ctx, s.queueName, msg.Bytes())
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
