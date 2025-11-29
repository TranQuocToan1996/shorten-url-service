package redis_stream

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"shorten/pkg/queue"
	"shorten/pkg/utils/patterns"

	"github.com/redis/go-redis/v9"
)

// TODO: Retry logic redis stream
// TODO: delete msg after ack, config maxlen
type redisStreamConsumer struct {
	client       *redis.Client
	handler      queue.MessageHandler
	group        string
	name         string
	block        time.Duration
	startID      string
	autoAck      bool
	valueField   string
	closeClient  bool
	ensureGroup  bool
	pendingIdle  time.Duration
	reclaimBatch int64
	pool         *patterns.WorkerPool
}

var _ queue.Consumer = (*redisStreamConsumer)(nil)

func NewRedisStreamConsumer(client *redis.Client, handler queue.MessageHandler, opts ...ConsumerOption) (queue.Consumer, error) {
	if client == nil {
		return nil, errors.New("redis client is required")
	}
	if handler == nil {
		return nil, errors.New("message handler is required")
	}

	pool := patterns.NewWorkerPool(runtime.NumCPU() * 16)
	pool.Run()

	c := &redisStreamConsumer{
		client:       client,
		handler:      handler,
		group:        "default-group",
		name:         fmt.Sprintf("consumer-%d", time.Now().UnixNano()),
		block:        5 * time.Second,
		startID:      ">",
		autoAck:      true,
		valueField:   defaultValueField,
		pendingIdle:  10 * time.Second,
		reclaimBatch: 20,
		pool:         pool,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *redisStreamConsumer) Consume(ctx context.Context, queueName string, _ []byte) error {
	if queueName == "" {
		return errors.New("queue name is required")
	}

	if c.ensureGroup {
		if err := c.ensureConsumerGroup(ctx, queueName); err != nil && !isBusyGroupErr(err) {
			return err
		}
	}

	streams := []string{queueName, c.startID}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := c.reclaimPending(ctx, queueName); err != nil {
			return err
		}

		res, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    c.group,
			Consumer: c.name,
			Streams:  streams,
			Count:    1,
			Block:    c.block,
			NoAck:    !c.autoAck,
		}).Result()

		if errors.Is(err, redis.Nil) {
			continue
		}
		if err != nil {
			return err
		}

		for _, str := range res {
			for _, msg := range str.Messages {
				msgCopy := msg
				c.pool.Submit(func() error {
					if err := c.processMessage(ctx, queueName, msgCopy); err != nil {
						log.Printf("Failed to process message %s: %v", msgCopy.ID, err)
					}
					return nil // Always return nil so worker continues
				})
			}
		}
	}
}

func (c *redisStreamConsumer) Close() error {
	if c.closeClient && c.client != nil {
		return c.client.Close()
	}
	c.pool.Close()
	return nil
}

func (c *redisStreamConsumer) ensureConsumerGroup(ctx context.Context, queueName string) error {
	return c.client.XGroupCreateMkStream(ctx, queueName, c.group, "0").Err()
}

func (c *redisStreamConsumer) reclaimPending(ctx context.Context, queueName string) error {
	if c.pendingIdle <= 0 || c.reclaimBatch <= 0 {
		return nil
	}

	pending, err := c.client.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: queueName,
		Group:  c.group,
		Start:  "-",
		End:    "+",
		Count:  c.reclaimBatch,
		Idle:   c.pendingIdle,
	}).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(pending) == 0 {
		return nil
	}

	ids := make([]string, 0, len(pending))
	for _, p := range pending {
		ids = append(ids, p.ID)
	}

	claimed, err := c.client.XClaim(ctx, &redis.XClaimArgs{
		Stream:   queueName,
		Group:    c.group,
		Consumer: c.name,
		Messages: ids,
		MinIdle:  c.pendingIdle,
	}).Result()
	if err != nil {
		return err
	}

	for _, msg := range claimed {
		msgCopy := msg
		c.pool.Submit(func() error {
			if err := c.processMessage(ctx, queueName, msgCopy); err != nil {
				log.Printf("Failed to process message %s: %v", msgCopy.ID, err)
			}
			return nil // Always return nil so worker continues
		})
	}

	return nil
}

func (c *redisStreamConsumer) processMessage(ctx context.Context, queueName string, msg redis.XMessage) error {
	payload, err := c.extractPayload(msg)
	if err != nil {
		if c.autoAck {
			_ = c.client.XAck(ctx, queueName, c.group, msg.ID).Err()
		}
		return err
	}

	if err := c.handler(ctx, queueName, payload); err != nil {
		return err
	}

	if c.autoAck {
		if err := c.client.XAck(ctx, queueName, c.group, msg.ID).Err(); err != nil {
			return err
		}
		// Delete msg in stream after process
		// c.client.XDel(ctx, queueName, msg.ID).Err()
	}

	return nil
}

func (c *redisStreamConsumer) extractPayload(msg redis.XMessage) ([]byte, error) {
	raw, ok := msg.Values[c.valueField]
	if !ok {
		return nil, fmt.Errorf("field %s not found in message", c.valueField)
	}

	switch v := raw.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported payload type %T", raw)
	}
}

type ConsumerOption func(*redisStreamConsumer)

func WithConsumerGroup(group string) ConsumerOption {
	return func(c *redisStreamConsumer) {
		if group != "" {
			c.group = group
		}
	}
}

func WithConsumerName(name string) ConsumerOption {
	return func(c *redisStreamConsumer) {
		if name != "" {
			c.name = name
		}
	}
}

func WithBlockTimeout(block time.Duration) ConsumerOption {
	return func(c *redisStreamConsumer) {
		if block > 0 {
			c.block = block
		}
	}
}

func WithStartID(startID string) ConsumerOption {
	return func(c *redisStreamConsumer) {
		if startID != "" {
			c.startID = startID
		}
	}
}

func WithManualAck() ConsumerOption {
	return func(c *redisStreamConsumer) {
		c.autoAck = false
	}
}

func WithConsumerValueField(field string) ConsumerOption {
	return func(c *redisStreamConsumer) {
		if field != "" {
			c.valueField = field
		}
	}
}

func WithConsumerOwnsClient() ConsumerOption {
	return func(c *redisStreamConsumer) {
		c.closeClient = true
	}
}

func WithEnsureGroup() ConsumerOption {
	return func(c *redisStreamConsumer) {
		c.ensureGroup = true
	}
}

func WithPendingIdle(idle time.Duration) ConsumerOption {
	return func(c *redisStreamConsumer) {
		if idle > 0 {
			c.pendingIdle = idle
		}
	}
}

func WithReclaimBatch(batch int64) ConsumerOption {
	return func(c *redisStreamConsumer) {
		if batch > 0 {
			c.reclaimBatch = batch
		}
	}
}

func isBusyGroupErr(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "BUSYGROUP")
}
