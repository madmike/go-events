package jetstream

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	evtypes "github.com/madmike/go-events/types"
	"github.com/nats-io/nats.go/jetstream"
)

// Publisher wraps the JetStream client with typed publish methods.
// Construct one per service; it is safe for concurrent use.
type Publisher struct {
	c *Client
}

// NewPublisher returns a Publisher backed by the given Client.
func NewPublisher(c *Client) *Publisher { return &Publisher{c: c} }

// Close drains the underlying NATS connection.
func (p *Publisher) Close() { p.c.Close() }

// Publish serialises payload as JSON, wraps it in an Envelope, and publishes
// to the given subject. Returns the server ack or an error.
func (p *Publisher) Publish(ctx context.Context, subject, tenantID string, payload any) (*jetstream.PubAck, error) {
	return p.publishWithTraceID(ctx, subject, tenantID, "", payload)
}

// PublishRealtime publishes a tenant-scoped dashboard event over plain NATS.
// These subjects intentionally bypass JetStream streams so browser subscribers
// receive low-latency ephemeral updates without requiring stream retention.
func (p *Publisher) PublishRealtime(ctx context.Context, subject, tenantID string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal realtime payload: %w", err)
	}
	env := evtypes.Envelope{
		ID:          uuid.Must(uuid.NewV7()).String(),
		Subject:     subject,
		TenantID:    tenantID,
		PublishedAt: time.Now().UTC(),
		Data:        data,
	}
	envBytes, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("marshal realtime envelope: %w", err)
	}
	if err := p.c.NC().Publish(subject, envBytes); err != nil {
		return fmt.Errorf("publish realtime %s: %w", subject, err)
	}
	return nil
}

// PublishWithTraceID publishes with a custom trace ID for distributed tracing.
func (p *Publisher) PublishWithTraceID(ctx context.Context, subject, tenantID, traceID string, payload any) (*jetstream.PubAck, error) {
	return p.publishWithTraceID(ctx, subject, tenantID, traceID, payload)
}

// publishWithTraceID is the internal implementation.
func (p *Publisher) publishWithTraceID(ctx context.Context, subject, tenantID, traceID string, payload any) (*jetstream.PubAck, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	env := evtypes.Envelope{
		ID:          uuid.Must(uuid.NewV7()).String(),
		Subject:     subject,
		TenantID:    tenantID,
		TraceID:     traceID,
		PublishedAt: time.Now().UTC(),
		Data:        data,
	}
	envBytes, err := json.Marshal(env)
	if err != nil {
		return nil, fmt.Errorf("marshal envelope: %w", err)
	}

	ack, err := p.c.js.Publish(ctx, subject, envBytes)
	if err != nil {
		return nil, fmt.Errorf("publish %s: %w", subject, err)
	}
	return ack, nil
}

// PublishInbound publishes an InboundMessage event.
func (p *Publisher) PublishInbound(ctx context.Context, msg evtypes.InboundMessage) error {
	subj := fmt.Sprintf("msg.inbound.%s.%s.%s", msg.Platform, msg.TenantID, msg.AccountID)
	_, err := p.PublishWithTraceID(ctx, subj, msg.TenantID, msg.TraceID, msg)
	return err
}

// PublishOutbound publishes an OutboundMessage event.
func (p *Publisher) PublishOutbound(ctx context.Context, msg evtypes.OutboundMessage) error {
	subj := fmt.Sprintf("msg.outbound.%s.%s.%s", msg.Platform, msg.TenantID, msg.AccountID)
	_, err := p.PublishWithTraceID(ctx, subj, msg.TenantID, msg.TraceID, msg)
	return err
}

// PublishIngestCompleted publishes an IngestCompleted event.
func (p *Publisher) PublishIngestCompleted(ctx context.Context, evt evtypes.IngestCompleted) error {
	subj := fmt.Sprintf("ingest.completed.%s", evt.TenantID)
	_, err := p.Publish(ctx, subj, evt.TenantID, evt)
	return err
}

// PublishIngestFailed publishes an IngestFailed event.
func (p *Publisher) PublishIngestFailed(ctx context.Context, evt evtypes.IngestFailed) error {
	subj := fmt.Sprintf("ingest.failed.%s", evt.TenantID)
	_, err := p.Publish(ctx, subj, evt.TenantID, evt)
	return err
}

// PublishBillingSubscriptionUpdated publishes a subscription state-change event.
func (p *Publisher) PublishBillingSubscriptionUpdated(ctx context.Context, evt evtypes.BillingSubscriptionUpdated) error {
	subj := fmt.Sprintf("billing.subscription.updated.%s", evt.TenantID)
	_, err := p.Publish(ctx, subj, evt.TenantID, evt)
	return err
}

// PublishScheduleFire publishes a schedule.fire event from the scheduler.
func (p *Publisher) PublishScheduleFire(ctx context.Context, evt evtypes.ScheduleFireEvent) error {
	subj := fmt.Sprintf("schedule.fire.%s.%s", evt.JobType, evt.TenantID)
	_, err := p.Publish(ctx, subj, evt.TenantID, evt)
	return err
}

// PublishIdentityAccount publishes an account identity-change event.
func (p *Publisher) PublishIdentityAccount(ctx context.Context, evt evtypes.IdentityAccountEvent) error {
	subj := fmt.Sprintf("identity.account.%s.%s", evt.Action, evt.AccountID)
	_, err := p.Publish(ctx, subj, evt.TenantID, evt)
	return err
}

// PublishIdentityTenant publishes a tenant identity-change event.
func (p *Publisher) PublishIdentityTenant(ctx context.Context, evt evtypes.IdentityTenantEvent) error {
	subj := fmt.Sprintf("identity.tenant.%s.%s", evt.Action, evt.TenantID)
	_, err := p.Publish(ctx, subj, evt.TenantID, evt)
	return err
}

// PublishIdentityBot publishes a bot identity-change event.
func (p *Publisher) PublishIdentityBot(ctx context.Context, evt evtypes.IdentityBotEvent) error {
	subj := fmt.Sprintf("identity.bot.%s.%s", evt.Action, evt.TenantAgentID)
	_, err := p.Publish(ctx, subj, evt.TenantID, evt)
	return err
}

// PublishCrawlSiteCompleted is published by crawler-service when a site crawl finishes.
func (p *Publisher) PublishCrawlSiteCompleted(ctx context.Context, evt evtypes.CrawlSiteCompleted) error {
	subj := fmt.Sprintf("crawl.site.completed.%s", evt.TenantID)
	_, err := p.Publish(ctx, subj, evt.TenantID, evt)
	return err
}

// PublishCrawlSiteFailed is published by crawler-service when a site crawl fails unrecoverably.
func (p *Publisher) PublishCrawlSiteFailed(ctx context.Context, evt evtypes.CrawlSiteFailed) error {
	subj := fmt.Sprintf("crawl.site.failed.%s", evt.TenantID)
	_, err := p.Publish(ctx, subj, evt.TenantID, evt)
	return err
}

// PublishCrawlAgentPersonaUpdated is published by crawler-service when it
// regenerates a tenant agent's system-prompt persona after a crawl finalises.
// identity-service subscribes and applies the overrides to platform.tenant_agents.
func (p *Publisher) PublishCrawlAgentPersonaUpdated(ctx context.Context, evt evtypes.CrawlAgentPersonaUpdated) error {
	subj := fmt.Sprintf("crawl.agent.persona_updated.%s", evt.TenantID)
	_, err := p.Publish(ctx, subj, evt.TenantID, evt)
	return err
}
