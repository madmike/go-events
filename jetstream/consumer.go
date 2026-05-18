package jetstream

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	evtypes "github.com/madmike/go-events/types"
	"github.com/madmike/go-infra/telemetry"
	"github.com/nats-io/nats.go/jetstream"
)

// Handler is a typed callback for a consumed envelope. Call msg.Ack() when
// processing succeeds, or msg.Nak() / msg.Term() on failure.
type Handler func(ctx context.Context, env evtypes.Envelope, msg jetstream.Msg) error

// ConsumerConfig defines a durable push consumer. Services create one per
// subscription they need. DurableName must be unique across the cluster.
type ConsumerConfig struct {
	// Stream is the JetStream stream name (e.g. subjects.StreamMessages).
	Stream string
	// DurableName uniquely identifies this consumer in the stream.
	// Convention: "<service-name>-<purpose>", e.g. "learnbot-inbound".
	DurableName string
	// FilterSubject narrows delivery (wildcards OK), e.g. "msg.inbound.telegram.>".
	FilterSubject string
	// MaxDeliver is the maximum number of delivery attempts before Term. 0 = unlimited.
	MaxDeliver int
	// AckWait is how long NATS waits for an Ack before redelivering.
	AckWait string // duration string, e.g. "30s"
	// DeliverPolicy defines where the consumer starts. Defaults to DeliverNewPolicy.
	// For WorkQueue streams, this MUST be DeliverAllPolicy.
	DeliverPolicy *jetstream.DeliverPolicy
}

// Consumer wraps a JetStream consumer and delivers typed Envelope messages.
type Consumer struct {
	c   *Client
	log telemetry.Logger
}

// NewConsumer creates a Consumer backed by the given Client.
func NewConsumer(c *Client, log telemetry.Logger) *Consumer {
	return &Consumer{c: c, log: log}
}

// Subscribe creates (or resumes) a durable consumer and starts delivering
// messages to handler. Blocks until ctx is cancelled. Safe to call in a
// goroutine per subscription.
func (cs *Consumer) Subscribe(ctx context.Context, cfg ConsumerConfig, handler Handler) error {
	ackWait := defaultAckWait
	if cfg.AckWait != "" {
		d, err := parseDuration(cfg.AckWait)
		if err == nil {
			ackWait = d
		}
	}

	maxDeliver := 5
	if cfg.MaxDeliver > 0 {
		maxDeliver = cfg.MaxDeliver
	}

	deliverPolicy := jetstream.DeliverNewPolicy
	if cfg.DeliverPolicy != nil {
		deliverPolicy = *cfg.DeliverPolicy
	}

	jsCfg := jetstream.ConsumerConfig{
		Durable:       cfg.DurableName,
		FilterSubject: cfg.FilterSubject,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxDeliver:    maxDeliver,
		AckWait:       ackWait,
		DeliverPolicy: deliverPolicy,
	}

	consumer, err := cs.c.js.CreateOrUpdateConsumer(ctx, cfg.Stream, jsCfg)
	if err != nil {
		return fmt.Errorf("create consumer %s/%s: %w", cfg.Stream, cfg.DurableName, err)
	}

	cs.log.Info("NATS consumer ready",
		telemetry.String("stream", cfg.Stream),
		telemetry.String("durable", cfg.DurableName),
		telemetry.String("filter", cfg.FilterSubject),
	)

	it, err := consumer.Messages()
	if err != nil {
		return fmt.Errorf("consumer messages iterator: %w", err)
	}
	defer it.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		msg, err := it.Next()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			cs.log.Warn("consumer next error", telemetry.Err(err))
			continue
		}
		cs.log.Info("consumer received message", telemetry.String("subject", msg.Subject()))

		var env evtypes.Envelope
		if err := json.Unmarshal(msg.Data(), &env); err != nil {
			cs.log.Error("unmarshal envelope", telemetry.Err(err),
				telemetry.String("subject", msg.Subject()))
			_ = msg.Term()
			continue
		}

		if err := handler(ctx, env, msg); err != nil {
			md, metaErr := msg.Metadata()
			if metaErr == nil && maxDeliver > 0 && md.NumDelivered >= uint64(maxDeliver) {
				cs.log.Error("max deliver reached, moving to DLQ",
					telemetry.String("subject", msg.Subject()),
					telemetry.String("id", env.ID),
					telemetry.Err(err),
				)
				_, pubErr := cs.c.js.Publish(ctx, "DLQ."+msg.Subject(), msg.Data())
				if pubErr != nil {
					cs.log.Error("failed to publish to DLQ", telemetry.Err(pubErr))
					_ = msg.Nak()
				} else {
					_ = msg.Term()
				}
			} else {
				cs.log.Warn("handler error — naking",
					telemetry.Err(err),
					telemetry.String("subject", msg.Subject()),
					telemetry.String("id", env.ID),
				)
				_ = msg.Nak()
			}
		} else {
			_ = msg.Ack()
		}
	}
}

// Unmarshal is a helper that decodes the envelope's Data field into target.
func Unmarshal[T any](env evtypes.Envelope) (T, error) {
	var v T
	if err := json.Unmarshal(env.Data, &v); err != nil {
		return v, fmt.Errorf("unmarshal %T: %w", v, err)
	}
	return v, nil
}

// ─────────────────────────────────────────────────────────────
// helpers
// ─────────────────────────────────────────────────────────────

func Pointer[T any](v T) *T {
	return &v
}

const defaultAckWait = 30 * time.Second

func parseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}
