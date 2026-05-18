package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEnvelopeRoundTrip(t *testing.T) {
	envelope := &Envelope{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Subject:     "msg.inbound.telegram.user",
		TenantID:    "tenant-123",
		PublishedAt: time.Date(2026, 5, 19, 10, 30, 0, 0, time.UTC),
		TraceID:     "trace-abc",
		Data:        []byte(`{"platform":"telegram"}`),
	}

	data, err := json.Marshal(envelope)
	require.NoError(t, err)

	var decoded Envelope
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, envelope.ID, decoded.ID)
	require.Equal(t, envelope.Subject, decoded.Subject)
	require.Equal(t, envelope.TenantID, decoded.TenantID)
	require.Equal(t, envelope.TraceID, decoded.TraceID)
	require.Equal(t, envelope.Data, decoded.Data)
}

func TestInboundMessageRoundTrip(t *testing.T) {
	msg := &InboundMessage{
		Platform:   "telegram",
		TenantID:   "tenant-123",
		AccountID:  "account-456",
		ChatID:     "chat-789",
		ThreadID:   "thread-101",
		FromUserID: "user-202",
		Username:   "john_doe",
		Text:       "Hello bot!",
		Attachments: []Attachment{
			{
				Type:     "image",
				MimeType: "image/jpeg",
				FileName: "photo.jpg",
				URL:      "https://example.com/photo.jpg",
				FileID:   "fileid-123",
			},
		},
		DeepLink:     "course_abc123",
		CallbackData: "btn_confirm",
		MessageID:    "msg-303",
		ReceivedAt:   time.Date(2026, 5, 19, 10, 30, 0, 0, time.UTC),
		TraceID:      "trace-xyz",
		RawPayload:   []byte(`{"update_id":123}`),
	}

	data, err := json.Marshal(msg)
	require.NoError(t, err)

	var decoded InboundMessage
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, msg.Platform, decoded.Platform)
	require.Equal(t, msg.TenantID, decoded.TenantID)
	require.Equal(t, msg.Text, decoded.Text)
	require.Equal(t, len(msg.Attachments), len(decoded.Attachments))
	require.Equal(t, msg.Attachments[0].Type, decoded.Attachments[0].Type)
	require.Equal(t, msg.DeepLink, decoded.DeepLink)
	require.Equal(t, msg.MessageID, decoded.MessageID)
}

func TestOutboundMessageRoundTrip(t *testing.T) {
	msg := &OutboundMessage{
		Platform:  "whatsapp",
		TenantID:  "tenant-123",
		AccountID: "account-456",
		ChatID:    "chat-789",
		ThreadID:  "thread-101",
		ReplyTo:   "msg-202",
		Text:      "Response text",
		ParseMode: "html",
		Buttons: [][]Button{
			{
				{Text: "Yes", CallbackData: "yes"},
				{Text: "No", CallbackData: "no"},
			},
		},
		DisablePreview:      true,
		DisableNotification: false,
		ChatAction:          "typing",
		TraceID:             "trace-xyz",
	}

	data, err := json.Marshal(msg)
	require.NoError(t, err)

	var decoded OutboundMessage
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, msg.Platform, decoded.Platform)
	require.Equal(t, msg.Text, decoded.Text)
	require.Equal(t, msg.ParseMode, decoded.ParseMode)
	require.Equal(t, len(msg.Buttons), len(decoded.Buttons))
	require.Equal(t, msg.DisablePreview, decoded.DisablePreview)
	require.Equal(t, msg.ChatAction, decoded.ChatAction)
}

func TestIngestCompletedRoundTrip(t *testing.T) {
	evt := &IngestCompleted{
		TenantID:   "tenant-123",
		JobID:      "job-456",
		SourceType: "pdf",
		SourceRef:  "s3://bucket/file.pdf",
		Title:      "Important Document",
		Language:   "en",
		ChunkCount: 42,
		Chunks: []TextChunk{
			{
				Index: 0,
				Text:  "First paragraph...",
			},
		},
		StorageRef: "platform.ingest_chunks/job_id=job-456",
		CallbackMetadata: map[string]any{
			"course_id": "course-789",
			"user_id":   "user-101",
		},
		CompletedAt: time.Date(2026, 5, 19, 10, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(evt)
	require.NoError(t, err)

	var decoded IngestCompleted
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, evt.TenantID, decoded.TenantID)
	require.Equal(t, evt.JobID, decoded.JobID)
	require.Equal(t, evt.SourceType, decoded.SourceType)
	require.Equal(t, evt.ChunkCount, decoded.ChunkCount)
	require.Equal(t, len(evt.Chunks), len(decoded.Chunks))
	require.Equal(t, evt.CallbackMetadata["course_id"], decoded.CallbackMetadata["course_id"])
}

func TestIngestFailedRoundTrip(t *testing.T) {
	evt := &IngestFailed{
		TenantID:  "tenant-123",
		JobID:     "job-456",
		SourceRef: "s3://bucket/file.pdf",
		CallbackMetadata: map[string]any{
			"retry_count": 3,
		},
		Error:    "timeout after 3 retries",
		FailedAt: time.Date(2026, 5, 19, 10, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(evt)
	require.NoError(t, err)

	var decoded IngestFailed
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, evt.TenantID, decoded.TenantID)
	require.Equal(t, evt.Error, decoded.Error)
	// JSON unmarshals integers as float64 in maps
	require.Equal(t, float64(3), decoded.CallbackMetadata["retry_count"])
}

func TestCrawlSiteCompletedRoundTrip(t *testing.T) {
	evt := &CrawlSiteCompleted{
		TenantID:     "tenant-123",
		SourceID:     "source-456",
		RootURL:      "https://example.com",
		EntryCount:   150,
		PagesCrawled: 45,
		CompletedAt:  time.Date(2026, 5, 19, 10, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(evt)
	require.NoError(t, err)

	var decoded CrawlSiteCompleted
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, evt.TenantID, decoded.TenantID)
	require.Equal(t, evt.EntryCount, decoded.EntryCount)
	require.Equal(t, evt.PagesCrawled, decoded.PagesCrawled)
	require.Equal(t, evt.RootURL, decoded.RootURL)
}

func TestCrawlSiteFailedRoundTrip(t *testing.T) {
	evt := &CrawlSiteFailed{
		TenantID: "tenant-123",
		SourceID: "source-456",
		RootURL:  "https://example.com",
		Error:    "robots.txt denied access",
		FailedAt: time.Date(2026, 5, 19, 10, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(evt)
	require.NoError(t, err)

	var decoded CrawlSiteFailed
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, evt.TenantID, decoded.TenantID)
	require.Equal(t, evt.Error, decoded.Error)
}

func TestAttachmentRoundTrip(t *testing.T) {
	tests := []struct {
		name       string
		attachment Attachment
	}{
		{
			"image_with_url",
			Attachment{
				Type:     "image",
				MimeType: "image/png",
				FileName: "screenshot.png",
				URL:      "https://cdn.example.com/img.png",
			},
		},
		{
			"document_with_file_id",
			Attachment{
				Type:   "document",
				FileID: "tg_fileid_abc123",
				Data:   []byte("fake pdf data"),
			},
		},
		{
			"voice_with_caption",
			Attachment{
				Type:    "voice",
				Caption: "User voice message",
				FileID:  "voice_id",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.attachment)
			require.NoError(t, err)

			var decoded Attachment
			require.NoError(t, json.Unmarshal(data, &decoded))

			require.Equal(t, tt.attachment.Type, decoded.Type)
			require.Equal(t, tt.attachment.MimeType, decoded.MimeType)
			require.Equal(t, tt.attachment.FileName, decoded.FileName)
			require.Equal(t, tt.attachment.URL, decoded.URL)
			require.Equal(t, tt.attachment.FileID, decoded.FileID)
			require.Equal(t, tt.attachment.Caption, decoded.Caption)
		})
	}
}

func TestButtonRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		button Button
	}{
		{"callback", Button{Text: "Confirm", CallbackData: "confirm"}},
		{"url", Button{Text: "Visit", URL: "https://example.com"}},
		{"simple", Button{Text: "OK"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.button)
			require.NoError(t, err)

			var decoded Button
			require.NoError(t, json.Unmarshal(data, &decoded))

			require.Equal(t, tt.button.Text, decoded.Text)
			require.Equal(t, tt.button.CallbackData, decoded.CallbackData)
			require.Equal(t, tt.button.URL, decoded.URL)
		})
	}
}

// SchemaBreakingChange detects if JSON field tags have been removed or renamed
// (regression test for append-only schema contract)
func TestTextChunkRoundTrip(t *testing.T) {
	chunk := TextChunk{
		Index: 0,
		Text:  "This is the first chunk of text.",
		Metadata: map[string]any{
			"page": 1,
			"tags": []string{"intro", "welcome"},
		},
	}

	data, err := json.Marshal(chunk)
	require.NoError(t, err)

	var decoded TextChunk
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, chunk.Index, decoded.Index)
	require.Equal(t, chunk.Text, decoded.Text)
	// JSON unmarshals integers as float64 in maps
	require.Equal(t, float64(1), decoded.Metadata["page"])
}

func TestBillingUsageRecordedRoundTrip(t *testing.T) {
	evt := &BillingUsageRecorded{
		TenantID:          "tenant-123",
		AssistantID:       "asst-456",
		MessageID:         "msg-789",
		Tier:              "complex",
		Provider:          "openai",
		Model:             "gpt-4o",
		InputTokens:       1500,
		OutputTokens:      500,
		CachedInputTokens: 300,
		CostMicroCents:    25000, // 0.25 cents
		HoldID:            "hold-101",
		Breakdown: map[string]any{
			"classifier": 10000,
			"embed":      5000,
			"llm":        10000,
		},
		OccurredAt: time.Date(2026, 5, 19, 10, 30, 0, 0, time.UTC),
		UpdatedAt:  time.Date(2026, 5, 19, 10, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(evt)
	require.NoError(t, err)

	var decoded BillingUsageRecorded
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, evt.TenantID, decoded.TenantID)
	require.Equal(t, evt.Provider, decoded.Provider)
	require.Equal(t, evt.Model, decoded.Model)
	require.Equal(t, evt.CostMicroCents, decoded.CostMicroCents)
	// JSON unmarshals integers as float64 in maps
	require.Equal(t, float64(10000), decoded.Breakdown["classifier"])
}

func TestBillingSubscriptionUpdatedRoundTrip(t *testing.T) {
	evt := &BillingSubscriptionUpdated{
		TenantID:       "tenant-123",
		UserID:         "user-456",
		SubscriptionID: "sub-789",
		ProductID:      "course-101",
		Plan:           "pro",
		Status:         SubscriptionActive,
		Provider:       "stripe",
		PeriodStart:    time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(evt)
	require.NoError(t, err)

	var decoded BillingSubscriptionUpdated
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, evt.TenantID, decoded.TenantID)
	require.Equal(t, evt.Status, decoded.Status)
	require.Equal(t, evt.Plan, decoded.Plan)
}

func TestBillingSubscriptionStatusConstants(t *testing.T) {
	statuses := []SubscriptionStatus{
		SubscriptionActive,
		SubscriptionTrialing,
		SubscriptionPastDue,
		SubscriptionCanceled,
		SubscriptionExpired,
	}

	for _, status := range statuses {
		require.NotEmpty(t, string(status), "status should have a value")
	}
}

func TestScheduleFireEventRoundTrip(t *testing.T) {
	evt := &ScheduleFireEvent{
		JobID:   "job-123",
		JobType: "daily_lesson",
		TenantID: "tenant-456",
		Payload: map[string]any{
			"course_id": "course-789",
			"lesson_index": 5,
		},
		FiredAt: time.Date(2026, 5, 19, 10, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(evt)
	require.NoError(t, err)

	var decoded ScheduleFireEvent
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, evt.JobID, decoded.JobID)
	require.Equal(t, evt.JobType, decoded.JobType)
	require.Equal(t, evt.Payload["course_id"], decoded.Payload["course_id"])
}

func TestCrawlAgentPersonaUpdatedRoundTrip(t *testing.T) {
	evt := &CrawlAgentPersonaUpdated{
		TenantID:      "tenant-123",
		TenantAgentID: "agent-456",
		SourceID:      "source-789",
		Overrides: map[string]any{
			"system_prompt": "You are a helpful assistant...",
			"temperature":   0.7,
		},
		UpdatedAt: time.Date(2026, 5, 19, 10, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(evt)
	require.NoError(t, err)

	var decoded CrawlAgentPersonaUpdated
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, evt.TenantID, decoded.TenantID)
	require.Equal(t, evt.Overrides["system_prompt"], decoded.Overrides["system_prompt"])
}

func TestIdentityAccountEventRoundTrip(t *testing.T) {
	config := &MessengerAccountConfig{
		Preset:           "telegram_bot",
		APIKeyRef:        "vault-key-123",
		PublicWebhookURL: "https://example.com/webhook",
		WebhookStatus:    "active",
		Options: map[string]any{
			"bot_name": "MyBot",
		},
	}

	evt := &IdentityAccountEvent{
		Action:    "created",
		AccountID: "account-456",
		TenantID:  "tenant-123",
		Platform:  "telegram",
		Config:    config,
		UpdatedAt: time.Date(2026, 5, 19, 10, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(evt)
	require.NoError(t, err)

	var decoded IdentityAccountEvent
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, evt.Action, decoded.Action)
	require.Equal(t, evt.Platform, decoded.Platform)
	require.NotNil(t, decoded.Config)
	require.Equal(t, config.Preset, decoded.Config.Preset)
	require.Equal(t, config.Options["bot_name"], decoded.Config.Options["bot_name"])
}

func TestSchemaAppendOnly(t *testing.T) {
	tests := []struct {
		name          string
		marshaler     interface{}
		requiredFields []string
	}{
		{
			"Envelope",
			&Envelope{},
			[]string{"id", "subject", "tenant_id", "published_at", "data"},
		},
		{
			"InboundMessage",
			&InboundMessage{},
			[]string{"platform", "tenant_id", "chat_id", "from_user_id", "message_id", "received_at"},
		},
		{
			"OutboundMessage",
			&OutboundMessage{},
			[]string{"platform", "tenant_id", "chat_id"},
		},
		{
			"IngestCompleted",
			&IngestCompleted{},
			[]string{"tenant_id", "job_id", "source_type", "chunk_count", "completed_at"},
		},
		{
			"BillingUsageRecorded",
			&BillingUsageRecorded{},
			[]string{"tenant_id", "provider", "model", "cost_micro_cents", "occurred_at"},
		},
		{
			"ScheduleFireEvent",
			&ScheduleFireEvent{},
			[]string{"job_id", "job_type", "tenant_id", "fired_at"},
		},
		{
			"IdentityAccountEvent",
			&IdentityAccountEvent{},
			[]string{"action", "account_id", "tenant_id", "platform", "updated_at"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.marshaler)
			require.NoError(t, err)

			var obj map[string]interface{}
			require.NoError(t, json.Unmarshal(data, &obj))

			for _, field := range tt.requiredFields {
				require.Contains(t, obj, field, "required field %s missing from %s", field, tt.name)
			}
		})
	}
}
