package queue

import "context"

type Producer interface {
	Publish(ctx context.Context, queueName string, payload []byte) error
	Close() error
}

type Consumer interface {
	Consume(ctx context.Context, queueName string, payload []byte) error
	Close() error
}

type MessageHandler func(ctx context.Context, queueName string, payload []byte) error
