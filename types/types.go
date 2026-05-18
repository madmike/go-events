// Package types defines the canonical JSON payload structs for every platform
// event. Services marshal/unmarshal these types when publishing or consuming
// NATS messages. All structs are append-only; add fields, never remove them.
package types

import "time"

// ─────────────────────────────────────────────────────────
// Envelope
// ─────────────────────────────────────────────────────────

// Envelope wraps every event with routing and tracing metadata.
// The Data field holds the JSON-encoded domain payload.
type Envelope struct {
	ID          string    `json:"id"`      // UUID v7, for deduplication
	Subject     string    `json:"subject"` // NATS subject this was published to
	TenantID    string    `json:"tenant_id"`
	PublishedAt time.Time `json:"published_at"`
	TraceID     string    `json:"trace_id,omitempty"` // OpenTelemetry trace propagation
	Data        []byte    `json:"data"`               // JSON-encoded domain payload
}

// ─────────────────────────────────────────────────────────
// Messenger events  (msg.inbound.* / msg.outbound.*)
// ─────────────────────────────────────────────────────────

// InboundMessage is published by messenger-gateway when a message arrives
// from any messenger platform. Products subscribe to these.
type InboundMessage struct {
	Platform   string `json:"platform"` // telegram | whatsapp | …
	TenantID   string `json:"tenant_id"`
	AccountID  string `json:"account_id"`
	ChatID     string `json:"chat_id"`
	ThreadID   string `json:"thread_id,omitempty"`
	FromUserID string `json:"from_user_id"`
	Username   string `json:"username,omitempty"`

	// Content
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`

	// Deep-link parameter from /start command (e.g. "course_abc123")
	DeepLink string `json:"deep_link,omitempty"`

	// Inline button callback data
	CallbackData string `json:"callback_data,omitempty"`

	MessageID  string    `json:"message_id"`
	ReceivedAt time.Time `json:"received_at"`
	TraceID    string    `json:"trace_id,omitempty"` // UUID v7 for distributed tracing

	// Full raw JSON payload from the platform for edge cases
	RawPayload []byte `json:"raw_payload,omitempty"`
}

// OutboundMessage is published by product services to ask the
// messenger-gateway to send a message on their behalf.
type OutboundMessage struct {
	Platform  string `json:"platform"`
	TenantID  string `json:"tenant_id"`
	AccountID string `json:"account_id"`
	ChatID    string `json:"chat_id"`
	ThreadID  string `json:"thread_id,omitempty"`
	ReplyTo   string `json:"reply_to,omitempty"`

	Text        string       `json:"text,omitempty"`
	ParseMode   string       `json:"parse_mode,omitempty"` // html | markdown | markdown_v2
	Attachments []Attachment `json:"attachments,omitempty"`

	// Inline keyboard rows
	Buttons [][]Button `json:"buttons,omitempty"`

	DisablePreview      bool `json:"disable_preview,omitempty"`
	DisableNotification bool `json:"disable_notification,omitempty"`
	ChatAction          string `json:"chat_action,omitempty"` // typing | upload_photo | …

	TraceID string `json:"trace_id,omitempty"` // UUID v7 for distributed tracing
}

type Attachment struct {
	Type     string `json:"type"` // image | audio | video | document | voice
	MimeType string `json:"mime_type,omitempty"`
	FileName string `json:"file_name,omitempty"`
	URL      string `json:"url,omitempty"`     // signed download URL
	FileID   string `json:"file_id,omitempty"` // platform-native file id
	Data     []byte `json:"data,omitempty"`    // inline bytes (for small files)
	Caption  string `json:"caption,omitempty"`
}

type Button struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	URL          string `json:"url,omitempty"`
}

// ─────────────────────────────────────────────────────────
// Ingestion events  (ingest.*)
// ─────────────────────────────────────────────────────────

// IngestCompleted is published by ingestion-service when a document has been
// processed into chunks and stored. Products use this to trigger their own
// domain logic (lesson generation, RAG indexing, etc.).
type IngestCompleted struct {
	TenantID   string      `json:"tenant_id"`
	JobID      string      `json:"job_id"`
	SourceType string      `json:"source_type"` // pdf | epub | youtube | url | text
	SourceRef  string      `json:"source_ref"`  // URL, storage path, or youtube video id
	Title      string      `json:"title,omitempty"`
	Language   string      `json:"language,omitempty"`
	ChunkCount int         `json:"chunk_count"`
	Chunks     []TextChunk `json:"chunks,omitempty"`      // inline for small docs (<150 chunks)
	StorageRef string      `json:"storage_ref,omitempty"` // "platform.ingest_chunks/job_id=<id>"

	// CallbackMetadata is opaque to ingestion-service — it is set by the caller
	// on job submission and echoed verbatim here. Products use it to carry
	// correlation IDs (course_id, bot_id, etc.) without a DB lookup.
	CallbackMetadata map[string]any `json:"callback_metadata,omitempty"`

	CompletedAt time.Time `json:"completed_at"`
}

// IngestFailed is published when ingestion fails after all retries.
type IngestFailed struct {
	TenantID         string         `json:"tenant_id"`
	JobID            string         `json:"job_id"`
	SourceRef        string         `json:"source_ref"`
	CallbackMetadata map[string]any `json:"callback_metadata,omitempty"`
	Error            string         `json:"error"`
	FailedAt         time.Time      `json:"failed_at"`
}

// ─────────────────────────────────────────────────────────
// Crawler events  (crawl.*)
// ─────────────────────────────────────────────────────────

// CrawlSiteCompleted is published by crawler-service when a full site crawl
// finishes. The entries are stored in the public.entries table; the event
// carries summary counts rather than the full data.
type CrawlSiteCompleted struct {
	TenantID     string    `json:"tenant_id"`
	SourceID     string    `json:"source_id"` // public.sources.id
	RootURL      string    `json:"root_url"`
	EntryCount   int       `json:"entry_count"`
	PagesCrawled int       `json:"pages_crawled"`
	CompletedAt  time.Time `json:"completed_at"`
}

// CrawlSiteFailed is published when a site crawl fails unrecoverably.
type CrawlSiteFailed struct {
	TenantID string    `json:"tenant_id"`
	SourceID string    `json:"source_id"`
	RootURL  string    `json:"root_url"`
	Error    string    `json:"error"`
	FailedAt time.Time `json:"failed_at"`
}

// CrawlAgentPersonaUpdated is published by crawler-service when it regenerates
// a tenant agent's system-prompt persona from the finalized site crawl.
// identity-service subscribes and applies the overrides to platform.tenant_agents.
type CrawlAgentPersonaUpdated struct {
	TenantID      string         `json:"tenant_id"`
	TenantAgentID string         `json:"tenant_agent_id"`
	SourceID      string         `json:"source_id"`
	Overrides     map[string]any `json:"overrides"` // partial overrides to merge
	UpdatedAt     time.Time      `json:"updated_at"`
}

type TextChunk struct {
	Index    int            `json:"index"`
	Text     string         `json:"text"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// ─────────────────────────────────────────────────────────
// Billing events  (billing.*)
// ─────────────────────────────────────────────────────────

// SubscriptionStatus mirrors the billing-service's lifecycle states.
type SubscriptionStatus string

const (
	SubscriptionActive   SubscriptionStatus = "active"
	SubscriptionTrialing SubscriptionStatus = "trialing"
	SubscriptionPastDue  SubscriptionStatus = "past_due"
	SubscriptionCanceled SubscriptionStatus = "canceled"
	SubscriptionExpired  SubscriptionStatus = "expired"
)

// BillingSubscriptionUpdated is published whenever a subscription state
// changes. Products gate feature access based on this event.
type BillingSubscriptionUpdated struct {
	TenantID       string             `json:"tenant_id"`
	UserID         string             `json:"user_id"`
	SubscriptionID string             `json:"subscription_id"`
	ProductID      string             `json:"product_id,omitempty"` // e.g. course ID
	Plan           string             `json:"plan"`
	Status         SubscriptionStatus `json:"status"`
	Provider       string             `json:"provider"` // stripe | telegram_stars
	PeriodStart    time.Time          `json:"period_start"`
	PeriodEnd      time.Time          `json:"period_end"`
}

// BillingUsageRecorded is published by agent-runtime after message execution
// completes and cost is captured. It feeds usage analytics and billing aggregation.
type BillingUsageRecorded struct {
	TenantID          string    `json:"tenant_id"`
	AssistantID       string    `json:"assistant_id"`
	MessageID         string    `json:"message_id"`
	Tier              string    `json:"tier"`               // simple | complex | default
	Provider          string    `json:"provider"`           // openai | gemini | …
	Model             string    `json:"model"`              // gpt-4o | claude-3-opus | …
	InputTokens       float64        `json:"input_tokens"`
	OutputTokens      float64        `json:"output_tokens"`
	CachedInputTokens float64        `json:"cached_input_tokens"`
	// Cost is stored in micro-cents so sub-cent per-message costs aren't truncated.
	CostMicroCents    int64     `json:"cost_micro_cents"`
	HoldID            string    `json:"hold_id"` // credit hold UUID

	// Breakdown contains detailed cost/token decomposition (e.g. classifier, embedding)
	Breakdown map[string]any `json:"breakdown,omitempty"`

	OccurredAt time.Time  `json:"occurred_at"`
	CanceledAt *time.Time `json:"canceled_at,omitempty"`
	TrialEnd   *time.Time `json:"trial_end,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// ─────────────────────────────────────────────────────────
// Schedule events  (schedule.*)
// ─────────────────────────────────────────────────────────

// ScheduleFireEvent is published by the scheduler-service when a scheduled
// job fires. Products consume job types they own.
type ScheduleFireEvent struct {
	JobID    string         `json:"job_id"`
	JobType  string         `json:"job_type"` // daily_lesson | weekly_recap | reminder | …
	TenantID string         `json:"tenant_id"`
	Payload  map[string]any `json:"payload,omitempty"` // job-specific parameters
	FiredAt  time.Time      `json:"fired_at"`
}

// ─────────────────────────────────────────────────────────
// Identity events  (identity.*)
// ─────────────────────────────────────────────────────────

// IdentityAccountEvent is published when a messenger account is created, updated, or deleted.
// messenger-gateway subscribes to hot-reload its provider registry.
type IdentityAccountEvent struct {
	Action    string `json:"action"` // created | updated | deleted
	AccountID string `json:"account_id"`
	TenantID  string `json:"tenant_id"`
	Platform  string `json:"platform"`
	// Full config is included on created/updated so the gateway can apply
	// it without a round-trip to identity-service.
	Config    *MessengerAccountConfig `json:"config,omitempty"`
	UpdatedAt time.Time               `json:"updated_at"`
}

// MessengerAccountConfig is the subset of account configuration the gateway needs.
type MessengerAccountConfig struct {
	Preset           string         `json:"preset"` // telegram_bot | whatsapp_cloud | …
	APIKey           string         `json:"api_key,omitempty"`           // plaintext for immediate use (legacy)
	APIKeyRef        string         `json:"api_key_ref,omitempty"`       // vault ref (new)
	APISecret        string         `json:"api_secret,omitempty"`        // plaintext for immediate use (legacy)
	APISecretRef     string         `json:"api_secret_ref,omitempty"`    // vault ref (new)
	ExternalID       string         `json:"external_id,omitempty"`
	PublicWebhookURL string         `json:"public_webhook_url"`
	WebhookSecret    string         `json:"webhook_secret,omitempty"`    // plaintext for immediate use (legacy)
	WebhookSecretRef string         `json:"webhook_secret_ref,omitempty"` // vault ref (new)
	WebhookStatus    string         `json:"webhook_status,omitempty"`    // pending | active | failed | password_needed
	Options          map[string]any `json:"options,omitempty"`
}

// IdentityTenantEvent is published when a tenant's profile changes.
type IdentityTenantEvent struct {
	Action    string    `json:"action"` // created | updated | deleted
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IdentityBotEvent is published when a tenant agent (bot) is created, updated, or deleted.
type IdentityBotEvent struct {
	Action        string    `json:"action"` // updated | deleted
	TenantAgentID string    `json:"tenant_agent_id"`
	TenantID      string    `json:"tenant_id"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ─────────────────────────────────────────────────────────
// Agent events (agent.stream.*)
// ─────────────────────────────────────────────────────────

// AgentEventType defines the kind of information in an agent stream event.
type AgentEventType string

const (
	AgentEventThought  AgentEventType = "thought"   // internal reasoning/planning
	AgentEventDelta    AgentEventType = "delta"     // token delta for user
	AgentEventDone     AgentEventType = "done"      // final chunk with usage/metrics
	AgentEventToolCall AgentEventType = "tool_call" // tool execution start
	AgentEventToolRes  AgentEventType = "tool_res"  // tool execution result
	AgentEventTyping   AgentEventType = "typing"    // typing indicator
	AgentEventError    AgentEventType = "error"
)

// AgentStreamEvent is published by agent-runtime during execution to support
// low-latency streaming. The stream is typically consumed via SSE or WebSockets.
type AgentStreamEvent struct {
	Type          AgentEventType `json:"type"`
	TenantID      string         `json:"tenant_id"`
	TenantAgentID string         `json:"tenant_agent_id"`
	RunID         string         `json:"run_id"` // unique ID for this execution turn

	// Messenger metadata for routing back to platforms
	Platform  string `json:"platform,omitempty"`
	AccountID string `json:"account_id,omitempty"`
	ChatID    string `json:"chat_id,omitempty"`

	// Content is the payload corresponding to the type:
	// - thought: reasoning string
	// - delta: token string
	// - done: final text or JSON metrics
	// - tool_call: tool name and call ID
	// - tool_res: tool output JSON
	Content string `json:"content,omitempty"`

	// Metadata contains type-specific structured data
	Metadata map[string]any `json:"metadata,omitempty"`

	Timestamp time.Time `json:"timestamp"`
}
