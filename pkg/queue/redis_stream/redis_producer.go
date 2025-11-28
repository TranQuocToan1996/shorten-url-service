package redis_stream

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"shorten/pkg/queue"
)

const defaultValueField = "payload"

type ProducerOption func(*redisStreamProducer)

func WithMaxLen(maxLen int64, approximate bool) ProducerOption {
	return func(p *redisStreamProducer) {
		p.maxLen = maxLen
		p.approximate = approximate
	}
}

func WithValueField(field string) ProducerOption {
	return func(p *redisStreamProducer) {
		if field != "" {
			p.valueField = field
		}
	}
}

func WithProducerOwnsClient() ProducerOption {
	return func(p *redisStreamProducer) {
		p.closeClient = true
	}
}

type redisStreamProducer struct {
	client      *redis.Client
	maxLen      int64
	approximate bool
	valueField  string
	closeClient bool
}

var _ queue.Producer = (*redisStreamProducer)(nil)

func NewRedisStreamProducer(client *redis.Client, opts ...ProducerOption) queue.Producer {
	p := &redisStreamProducer{
		client:      client,
		valueField:  defaultValueField,
		approximate: true,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *redisStreamProducer) Publish(ctx context.Context, queueName string, payload []byte) error {
	if p.client == nil {
		return errors.New("redis client is nil")
	}
	if queueName == "" {
		return errors.New("queue name is required")
	}
	if len(payload) == 0 {
		return errors.New("payload is required")
	}

	args := &redis.XAddArgs{
		Stream: queueName,
		Values: map[string]any{
			p.valueField: payload,
		},
	}

	if p.maxLen > 0 {
		args.MaxLen = p.maxLen
		args.Approx = p.approximate
	}

	return p.client.XAdd(ctx, args).Err()
}

func (p *redisStreamProducer) Close() error {
	if p.closeClient && p.client != nil {
		return p.client.Close()
	}
	return nil
}
