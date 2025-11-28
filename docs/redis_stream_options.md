Reviewing the producer and consumer implementations to list all options:


[2 tools called]


## Producer options

### 1. `WithMaxLen(maxLen int64, approximate bool)`
- Purpose: Limits the Redis Stream length to prevent unbounded growth
- Parameters:
  - `maxLen`: Maximum number of messages to keep (0 = no limit)
  - `approximate`: `true` for approximate trimming (faster), `false` for exact trimming
- Default: `0` (no limit), `approximate = true`
- Usage:
```go
redis_stream.WithMaxLen(10000, true)  // Keep max 10k, approximate
redis_stream.WithMaxLen(5000, false) // Keep max 5k, exact
```
- When to use: Production to prevent memory issues

### 2. `WithValueField(field string)`
- Purpose: Sets the Redis Stream field name where payload is stored
- Parameters:
  - `field`: Field name (must match consumer's `WithConsumerValueField`)
- Default: `"payload"`
- Usage:
```go
redis_stream.WithValueField("payload")  // Default
redis_stream.WithValueField("data")     // Custom
```
- When to use: If you need a custom field name (must match consumer)

### 3. `WithProducerOwnsClient()`
- Purpose: Makes producer close the Redis client when `Close()` is called
- Parameters: None
- Default: `false` (doesn't close client)
- Usage:
```go
redis_stream.WithProducerOwnsClient()
```
- When to use: Only if the producer is the sole owner of the Redis client
- Warning: Don't use in API server if the client is shared

---

## Consumer options

### 1. `WithConsumerGroup(group string)`
- Purpose: Sets the consumer group name for Redis Streams consumer groups
- Parameters:
  - `group`: Consumer group name
- Default: `"default-group"`
- Usage:
```go
redis_stream.WithConsumerGroup("shorten-url-group")
```
- When to use: Always set a meaningful name for your application
- Note: Multiple consumers can share the same group for load balancing

### 2. `WithConsumerName(name string)`
- Purpose: Sets the unique consumer name within the group
- Parameters:
  - `name`: Consumer identifier
- Default: `"consumer-{timestamp}"` (auto-generated)
- Usage:
```go
redis_stream.WithConsumerName("shorten-url-worker")
redis_stream.WithConsumerName("worker-1")
```
- When to use: For identification and monitoring
- Note: Each consumer instance should have a unique name

### 3. `WithBlockTimeout(block time.Duration)`
- Purpose: Sets how long to block waiting for new messages
- Parameters:
  - `block`: Duration to wait (0 = no blocking)
- Default: `5 * time.Second`
- Usage:
```go
redis_stream.WithBlockTimeout(10 * time.Second)  // Wait 10s
redis_stream.WithBlockTimeout(1 * time.Second)   // Wait 1s
```
- When to use: Adjust based on latency requirements
- Note: Longer timeouts reduce CPU usage but increase latency

### 4. `WithStartID(startID string)`
- Purpose: Sets the starting message ID for reading from the stream
- Parameters:
  - `startID`: Message ID or special value
- Default: `">"` (read only new messages)
- Special values:
  - `">"`: Only new messages (default)
  - `"0"`: From the beginning
  - `"$"`: Latest message
  - Specific ID: `"1234567890-0"`
- Usage:
```go
redis_stream.WithStartID(">")        // New messages only
redis_stream.WithStartID("0")        // From beginning
redis_stream.WithStartID("1234-0")   // From specific ID
```
- When to use: Usually keep `">"` unless replaying

### 5. `WithManualAck()`
- Purpose: Disables automatic acknowledgment (requires manual ACK)
- Parameters: None
- Default: `autoAck = true` (automatic)
- Usage:
```go
redis_stream.WithManualAck()
```
- When to use: When you need explicit control over message acknowledgment
- Note: With manual ACK, you must acknowledge messages yourself (not implemented in current consumer)

### 6. `WithConsumerValueField(field string)`
- Purpose: Sets the field name to extract payload from stream messages
- Parameters:
  - `field`: Field name (must match producer's `WithValueField`)
- Default: `"payload"`
- Usage:
```go
redis_stream.WithConsumerValueField("payload")  // Default
redis_stream.WithConsumerValueField("data")     // Custom
```
- When to use: Must match producer's value field if customized
- Critical: Must match producer's `WithValueField()` option

### 7. `WithConsumerOwnsClient()`
- Purpose: Makes consumer close the Redis client when `Close()` is called
- Parameters: None
- Default: `false` (doesn't close client)
- Usage:
```go
redis_stream.WithConsumerOwnsClient()
```
- When to use: Only if the consumer is the sole owner of the Redis client
- Warning: Don't use if the client is shared

### 8. `WithEnsureGroup()`
- Purpose: Automatically creates the consumer group if it doesn't exist
- Parameters: None
- Default: `false` (group must exist)
- Usage:
```go
redis_stream.WithEnsureGroup()
```
- When to use: Recommended for production to avoid manual setup
- Note: Creates group with `XGroupCreateMkStream` if missing

### 9. `WithPendingIdle(idle time.Duration)`
- Purpose: Sets the minimum idle time before reclaiming pending messages
- Parameters:
  - `idle`: Minimum idle duration (0 = disable reclaim)
- Default: `30 * time.Second`
- Usage:
```go
redis_stream.WithPendingIdle(60 * time.Second)  // Reclaim after 60s
redis_stream.WithPendingIdle(0)                 // Disable reclaim
```
- When to use: Adjust based on processing time
- Note: Messages idle longer than this are reclaimed from other consumers

### 10. `WithReclaimBatch(batch int64)`
- Purpose: Sets how many pending messages to check for reclaiming per iteration
- Parameters:
  - `batch`: Number of messages to check (0 = disable reclaim)
- Default: `20`
- Usage:
```go
redis_stream.WithReclaimBatch(50)  // Check 50 messages
redis_stream.WithReclaimBatch(0)   // Disable reclaim
```
- When to use: Adjust based on expected pending message volume
- Note: Works with `WithPendingIdle()` for message recovery

---

## Compatibility matrix

| Option | Producer | Consumer | Must Match? |
|--------|----------|----------|-------------|
| Value Field | `WithValueField()` | `WithConsumerValueField()` | Yes |
| Queue Name | Passed to `Publish()` | Passed to `Consume()` | Yes |
| Consumer Group | N/A | `WithConsumerGroup()` | N/A |
| Consumer Name | N/A | `WithConsumerName()` | N/A |
| MaxLen | `WithMaxLen()` | N/A | N/A |
| Ensure Group | N/A | `WithEnsureGroup()` | N/A |

---

## Recommended configuration

### Producer (API Server):
```go
redis_stream.NewRedisStreamProducer(redisClient,
    redis_stream.WithMaxLen(10000, true), // Prevent unbounded growth
    // Don't use WithProducerOwnsClient() if client is shared
)
```

### Consumer (Worker):
```go
redis_stream.NewRedisStreamConsumer(redisClient, handler,
    redis_stream.WithConsumerGroup("shorten-url-group"),
    redis_stream.WithConsumerName("shorten-url-worker"),
    redis_stream.WithEnsureGroup(),                    // Auto-create group
    redis_stream.WithBlockTimeout(5 * time.Second),   // Default is fine
    redis_stream.WithPendingIdle(30 * time.Second),   // Default is fine
    redis_stream.WithReclaimBatch(20),                 // Default is fine
    // WithConsumerValueField("payload") - only if producer uses custom field
)
```

This setup ensures compatibility and production readiness.