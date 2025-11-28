Analyzing the codebase for delivery guarantees, ordering, and potential conflicts with multiple producers and consumers.


[4 tools called]


## Analysis: 3 Producers + 3 Consumers on Same Queue

### 1. Delivery mode: At-least-once delivery

Current behavior:
- `autoAck = true` (default)
- ACK happens after handler succeeds (lines 182-185)
- If handler fails, message is not ACKed (line 178-180)
- Pending messages are reclaimed after `pendingIdle` (30s default)

```169:189:pkg/queue/redis_stream/redis_consumer.go
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
	}

	return nil
}
```

Result: At-least-once delivery
- Messages are redelivered if processing fails
- Messages can be processed multiple times if a consumer crashes mid-processing

### 2. Ordering: FIFO within stream

How it works:
- Redis Streams maintains FIFO ordering
- Multiple producers: messages appended in order they arrive at Redis
- Multiple consumers in the same group: each message delivered to one consumer, maintaining order per consumer

```84:91:pkg/queue/redis_stream/redis_consumer.go
res, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
	Group:    c.group,
	Consumer: c.name,
	Streams:  streams,
	Count:    1,
	Block:    c.block,
	NoAck:    !c.autoAck,
}).Result()
```

Result: FIFO ordering guaranteed
- Messages processed in the order they were published
- Each consumer processes messages sequentially

### 3. At-least-once send: Yes

Producer behavior:
```59:83:pkg/queue/redis_stream/redis_producer.go
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
```

- `XAdd` is atomic
- If it succeeds, the message is in the stream
- If it fails, an error is returned
- Result: At-least-once send (if you retry on failure)

### 4. At-least-once delivery: Yes

Consumer behavior:
- Messages ACKed only after successful processing
- Failed messages remain in the pending list
- Pending messages are reclaimed after `pendingIdle` (30s)

```121:167:pkg/queue/redis_stream/redis_consumer.go
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
		if err := c.processMessage(ctx, queueName, msg); err != nil {
			return err
		}
	}

	return nil
}
```

Result: At-least-once delivery
- Messages are redelivered if processing fails or a consumer crashes

### 5. Consumer conflicts: No conflicts (load balanced)

How it works:
- All consumers in the same group share work
- Each message is delivered to one consumer only
- Redis Streams handles distribution atomically

Scenario: 3 consumers, same group
```
Consumer-1: Processes messages 1, 4, 7, 10...
Consumer-2: Processes messages 2, 5, 8, 11...
Consumer-3: Processes messages 3, 6, 9, 12...
```

Potential issue: duplicate processing
- If Consumer-1 crashes after receiving message but before ACK, the message becomes pending
- After `pendingIdle` (30s), another consumer can reclaim and process it
- Result: the same message may be processed twice (at-least-once)

### 6. Producer conflicts: No conflicts (atomic append)

How it works:
- `XAdd` is atomic
- Multiple producers can append concurrently
- Redis handles ordering

```70:82:pkg/queue/redis_stream/redis_producer.go
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
```

Result: No conflicts
- All producers can publish concurrently
- Messages are appended in order they arrive at Redis

---

## Summary table

| Aspect | Status | Details |
|--------|--------|---------|
| Delivery mode | At-least-once | Messages can be redelivered if processing fails |
| Ordering | FIFO | Messages processed in publish order |
| At-least-once send | Yes | `XAdd` is atomic, retry on failure |
| At-least-once delivery | Yes | Failed messages are redelivered |
| Consumer conflicts | No | Load balanced, each message to one consumer |
| Producer conflicts | No | Atomic append, concurrent writes safe |

---

## Potential issues and recommendations

### Issue 1: Duplicate processing
Problem: If a consumer crashes after processing but before ACK, the message may be reprocessed.

Current code:
```go
if err := c.handler(ctx, queueName, payload); err != nil {
	return err  // Message not ACKed, will be redelivered
}
```

Recommendation: Make handlers idempotent or add deduplication.

### Issue 2: Error handling stops consumer
Problem: If `processMessage` returns an error, the consumer loop exits (line 102-104).

```100:107:pkg/queue/redis_stream/redis_consumer.go
for _, str := range res {
	for _, msg := range str.Messages {
		if err := c.processMessage(ctx, queueName, msg); err != nil {
			return err  // Consumer stops!
		}
	}
}
```

Recommendation: Log errors and continue processing instead of stopping.

### Issue 3: No retry mechanism
Problem: Failed messages are redelivered after 30s, but there's no exponential backoff or max retries.

Recommendation: Add retry logic with exponential backoff and max retry limits.

---

## Recommended improvements

1. Make handlers idempotent (check if URL already processed)
2. Continue processing on errors instead of stopping the consumer
3. Add retry logic with exponential backoff
4. Consider dead-letter queue for messages that fail repeatedly
5. Add monitoring for pending messages and consumer lag

Current implementation is safe for concurrent producers and consumers, with at-least-once delivery guarantees.