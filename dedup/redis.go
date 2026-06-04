// Package dedup provides Redis-based event deduplication for at-least-once
// delivery semantics. Processed event IDs are tracked with a configurable TTL
// window so that duplicate events are safely ignored within that window.
package dedup

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Deduplicator prevents duplicate event processing using Redis.
// Stores processed event IDs with a TTL-based retention window.
type Deduplicator struct {
	client *redis.Client
	ttl    time.Duration
	prefix string
}

// New creates a Deduplicator backed by the given Redis client.
// The ttl controls how long processed event IDs are retained; events arriving
// after the TTL has expired will be treated as new.
func New(client *redis.Client, ttl time.Duration) *Deduplicator {
	return &Deduplicator{
		client: client,
		ttl:    ttl,
		prefix: "event_dedup:",
	}
}

// IsProcessed checks whether an event ID has already been recorded.
// Returns true when the event was previously processed, false when it is new.
func (d *Deduplicator) IsProcessed(ctx context.Context, eventID string) (bool, error) {
	key := d.prefix + eventID
	exists, err := d.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

// MarkProcessed records an event ID so that future IsProcessed calls return
// true. The entry expires after the configured TTL.
func (d *Deduplicator) MarkProcessed(ctx context.Context, eventID string) error {
	key := d.prefix + eventID
	return d.client.Set(ctx, key, "1", d.ttl).Err()
}

// CheckAndMark atomically checks and records an event ID via SETNX.
// Returns true when the event is new (first time being seen).
// Returns false when the event was already processed.
func (d *Deduplicator) CheckAndMark(ctx context.Context, eventID string) (bool, error) {
	key := d.prefix + eventID
	result, err := d.client.SetNX(ctx, key, "1", d.ttl).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}
