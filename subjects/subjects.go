// Package subjects defines the canonical NATS subject namespace for the
// platform. All inter-service communication uses these helpers so that subject
// strings never drift across repos.
//
// Subject hierarchy:
//
//	msg.inbound.<platform>.<tenant_id>.<bot_id>    — messenger-gateway → products
//	msg.outbound.<platform>.<tenant_id>.<bot_id>   — products → messenger-gateway
//	ingest.completed.<tenant_id>                   — ingestion-service → products
//	ingest.failed.<tenant_id>                      — ingestion-service → products
//	crawl.site.completed.<tenant_id>               — crawler-service → products
//	crawl.site.failed.<tenant_id>                  — crawler-service → products
//	billing.subscription.updated.<tenant_id>       — billing-service → products
//	schedule.fire.<job_type>.<tenant_id>           — scheduler-service → products
//	identity.bot.<action>.<bot_id>                 — identity-service → gateway
//	identity.tenant.<action>.<tenant_id>           — identity-service → all
package subjects

import "strings"

const (
	// Stream names — each stream covers one subject prefix family.
	StreamMessages     = "MESSAGES"
	StreamIngest       = "INGEST"
	StreamCrawl        = "CRAWL"
	StreamBilling      = "BILLING"
	StreamSchedule     = "SCHEDULE"
	StreamIdentity     = "IDENTITY"
	StreamConversation = "CONVERSATION"

	// Wildcard filter subjects per stream (used when creating JetStream streams).
	FilterMessages     = "msg.>"
	FilterIngest       = "ingest.>"
	FilterCrawl        = "crawl.>"
	FilterBilling      = "billing.>"
	FilterSchedule     = "schedule.>"
	FilterIdentity     = "identity.>"
	FilterAgent        = "agent.>"
	FilterConversation = "conversation.>"
	FilterRealtime     = "tenant.>"
)

// MsgInbound returns the subject for an inbound message arriving from a
// messenger platform. Products subscribe with platform=">" to receive all, or
// with a specific platform to filter.
//
//	msg.inbound.telegram.tenant123.bot456
func MsgInbound(platform, tenantID, botID string) string {
	return join("msg", "inbound", platform, tenantID, botID)
}

// MsgInboundFilter returns a wildcard filter for inbound messages.
// Pass "*" for any segment you want to wildcard at one level.
//
//	msg.inbound.telegram.>.  (all bots for telegram)
//	msg.inbound.>.          (all platforms, all tenants)
func MsgInboundFilter(platform, tenantID, botID string) string {
	return join("msg", "inbound", platform, tenantID, botID)
}

// MsgOutbound returns the subject a product publishes to when it wants the
// messenger-gateway to send a message.
//
//	msg.outbound.telegram.tenant123.bot456
func MsgOutbound(platform, tenantID, botID string) string {
	return join("msg", "outbound", platform, tenantID, botID)
}

// IngestCompleted is published by the ingestion-service when a document has
// been fully processed into chunks and stored.
//
//	ingest.completed.tenant123
func IngestCompleted(tenantID string) string {
	return join("ingest", "completed", tenantID)
}

// IngestFailed is published when ingestion fails after all retries.
func IngestFailed(tenantID string) string {
	return join("ingest", "failed", tenantID)
}

// BillingSubscriptionUpdated is published by the billing-service whenever a
// subscription state changes (created, renewed, cancelled, expired).
//
//	billing.subscription.updated.tenant123
func BillingSubscriptionUpdated(tenantID string) string {
	return join("billing", "subscription", "updated", tenantID)
}

// BillingUsageRecorded is published by agent-runtime after cost is captured
// for a message execution. Analytics workers consume this for aggregation.
//
//	billing.usage.recorded.tenant123
func BillingUsageRecorded(tenantID string) string {
	return join("billing", "usage", "recorded", tenantID)
}

// ScheduleFire is published by the scheduler-service when a scheduled job
// fires. Products subscribe to the job types they care about.
//
//	schedule.fire.daily_lesson.tenant123
func ScheduleFire(jobType, tenantID string) string {
	return join("schedule", "fire", jobType, tenantID)
}

// ScheduleFireFilter returns a filter for all jobs of a given type.
//
//	schedule.fire.daily_lesson.>
func ScheduleFireFilter(jobType string) string {
	return join("schedule", "fire", jobType, ">")
}

// IdentityBotUpdated is published when a bot's config changes in
// identity-service. messenger-gateway uses this to hot-reload its registry.
//
//	identity.bot.updated.bot456
func IdentityBotUpdated(botID string) string {
	return join("identity", "bot", "updated", botID)
}

// IdentityBotDeleted is published when a bot is removed.
func IdentityBotDeleted(botID string) string {
	return join("identity", "bot", "deleted", botID)
}

// IdentityTenantUpdated is published when a tenant's profile changes.
func IdentityTenantUpdated(tenantID string) string {
	return join("identity", "tenant", "updated", tenantID)
}

// CrawlSiteCompleted is published by crawler-service when a full site crawl
// finishes successfully.
//
//	crawl.site.completed.tenant123
func CrawlSiteCompleted(tenantID string) string {
	return join("crawl", "site", "completed", tenantID)
}

// CrawlSiteFailed is published when a site crawl fails unrecoverably.
func CrawlSiteFailed(tenantID string) string {
	return join("crawl", "site", "failed", tenantID)
}

// CrawlAgentPersonaUpdated is published by crawler-service when it synthesises
// a new system-prompt persona for a tenant agent after a crawl finalises.
// identity-service subscribes and applies the overrides.
//
//	crawl.agent.persona_updated.tenant123
func CrawlAgentPersonaUpdated(tenantID string) string {
	return join("crawl", "agent", "persona_updated", tenantID)
}

// ConversationCompactRequested is published by agent-runtime after a turn that
// pushes a conversation over the compaction threshold. The compactor consumer
// picks it up and runs the LLM summariser asynchronously.
//
//	conversation.compact.requested.<tenant_id>.<conversation_id>
func ConversationCompactRequested(tenantID, convID string) string {
	return join("conversation", "compact", "requested", tenantID, convID)
}

// ConversationCompactRequestedFilter returns a wildcard for all compaction
// requests (used by the compactor NATS consumer).
func ConversationCompactRequestedFilter() string {
	return join("conversation", "compact", "requested", ">")
}

// ConversationCompacted is published after a successful compaction run.
//
//	conversation.compacted.<tenant_id>
func ConversationCompacted(tenantID string) string {
	return join("conversation", "compacted", tenantID)
}

// Realtime returns a browser-facing tenant-scoped subject.
//
//	tenant.<tenant_id>.<domain>.<resource>.<event>
func Realtime(tenantID, domain, resource, event string) string {
	return join("tenant", tenantID, domain, resource, event)
}

// RealtimeFilter returns a wildcard tenant-scoped filter for dashboard clients.
func RealtimeFilter(tenantID, domain, resource string) string {
	return join("tenant", tenantID, domain, resource, ">")
}

func join(parts ...string) string {
	return strings.Join(parts, ".")
}
