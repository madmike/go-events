// Package jetstream provides a thin, opinionated wrapper around the NATS
// JetStream client. It manages stream creation/upsert at startup and exposes
// typed publish and consume helpers used across all platform services.
package jetstream

import (
	"context"
	"fmt"
	"time"

	"github.com/madmike/go-events/subjects"
	"github.com/madmike/go-infra/telemetry"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Config holds NATS connection parameters, read from service config files.
type Config struct {
	URL            string        `yaml:"url"`             // nats://nats:4222
	MaxReconnects  int           `yaml:"max_reconnects"`  // -1 = unlimited
	ReconnectWait  time.Duration `yaml:"reconnect_wait"`  // default 2s
	ConnectTimeout time.Duration `yaml:"connect_timeout"` // default 10s
}

func (c *Config) setDefaults() {
	if c.URL == "" {
		c.URL = nats.DefaultURL
	}
	if c.MaxReconnects == 0 {
		c.MaxReconnects = -1
	}
	if c.ReconnectWait == 0 {
		c.ReconnectWait = 2 * time.Second
	}
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = 10 * time.Second
	}
}

// Client wraps a NATS connection and a JetStream context.
// Obtain one via Connect(); all services share this shape.
type Client struct {
	nc  *nats.Conn
	js  jetstream.JetStream
	log telemetry.Logger
}

// Connect dials NATS, ensures all platform streams exist, and returns a ready
// Client. Call Close() when the service shuts down.
func Connect(cfg Config, log telemetry.Logger) (*Client, error) {
	cfg.setDefaults()

	opts := []nats.Option{
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.Timeout(cfg.ConnectTimeout),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				log.Warn("NATS disconnected", telemetry.Err(err))
			}
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Info("NATS reconnected")
		}),
	}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("nats connect %s: %w", cfg.URL, err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("jetstream context: %w", err)
	}

	c := &Client{nc: nc, js: js, log: log}
	if err := c.ensureStreams(context.Background()); err != nil {
		nc.Close()
		return nil, fmt.Errorf("ensure streams: %w", err)
	}

	log.Info("NATS JetStream connected", telemetry.String("url", cfg.URL))
	return c, nil
}

// JS returns the underlying JetStream context for advanced usage.
func (c *Client) JS() jetstream.JetStream { return c.js }

// NC returns the underlying NATS connection.
func (c *Client) NC() *nats.Conn { return c.nc }

// Close drains and closes the NATS connection gracefully.
func (c *Client) Close() {
	if err := c.nc.Drain(); err != nil {
		c.log.Warn("NATS drain error", telemetry.Err(err))
	}
}

// platformStreams defines every JetStream stream used by the platform.
// Streams are created with AddOrUpdateStream so they're idempotent on restart.
var platformStreams = []jetstream.StreamConfig{
	{
		Name:        subjects.StreamMessages,
		Description: "Inbound, outbound, and agent execution events",
		Subjects:    []string{subjects.FilterMessages, subjects.FilterAgent},
		Retention:   jetstream.LimitsPolicy,
		MaxAge:      24 * time.Hour, // messages are transient; 24h for replay
		Storage:     jetstream.FileStorage,
		Replicas:    1,
		Discard:     jetstream.DiscardOld,
	},
	{
		Name:        subjects.StreamIngest,
		Description: "Document ingestion lifecycle events (ingestion-service)",
		Subjects:    []string{subjects.FilterIngest},
		Retention:   jetstream.LimitsPolicy,
		MaxAge:      7 * 24 * time.Hour,
		Storage:     jetstream.FileStorage,
		Replicas:    1,
	},
	{
		Name:        subjects.StreamCrawl,
		Description: "Site crawl lifecycle events (crawler-service)",
		Subjects:    []string{subjects.FilterCrawl},
		Retention:   jetstream.LimitsPolicy,
		MaxAge:      7 * 24 * time.Hour,
		Storage:     jetstream.FileStorage,
		Replicas:    1,
	},
	{
		Name:        subjects.StreamBilling,
		Description: "Subscription lifecycle events from billing-service",
		Subjects:    []string{subjects.FilterBilling},
		Retention:   jetstream.LimitsPolicy,
		MaxAge:      30 * 24 * time.Hour,
		Storage:     jetstream.FileStorage,
		Replicas:    1,
	},
	{
		Name:        subjects.StreamSchedule,
		Description: "Scheduled job fire events from scheduler-service",
		Subjects:    []string{subjects.FilterSchedule},
		Retention:   jetstream.WorkQueuePolicy, // each job fires exactly once
		MaxAge:      24 * time.Hour,
		Storage:     jetstream.FileStorage,
		Replicas:    1,
	},
	{
		Name:        subjects.StreamIdentity,
		Description: "Identity change events (bots, tenants)",
		Subjects:    []string{subjects.FilterIdentity},
		Retention:   jetstream.LimitsPolicy,
		MaxAge:      7 * 24 * time.Hour,
		Storage:     jetstream.FileStorage,
		Replicas:    1,
	},
	{
		Name:        "MESSAGES_DLQ",
		Description: "Dead letter queue for message processing failures",
		Subjects:    []string{"DLQ.msg.>"},
		Retention:   jetstream.LimitsPolicy,
		MaxAge:      7 * 24 * time.Hour,
		Storage:     jetstream.FileStorage,
		Replicas:    1,
	},
	{
		Name:        subjects.StreamConversation,
		Description: "Conversation lifecycle events: compaction requests and completions",
		Subjects:    []string{subjects.FilterConversation},
		Retention:   jetstream.LimitsPolicy,
		MaxAge:      24 * time.Hour,
		Storage:     jetstream.FileStorage,
		Replicas:    1,
	},
}

func (c *Client) ensureStreams(ctx context.Context) error {
	for _, cfg := range platformStreams {
		cfg := cfg
		_, err := c.js.CreateOrUpdateStream(ctx, cfg)
		if err != nil {
			return fmt.Errorf("stream %s: %w", cfg.Name, err)
		}
		c.log.Info("NATS stream ready", telemetry.String("stream", cfg.Name))
	}
	return nil
}
