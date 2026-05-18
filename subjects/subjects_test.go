package subjects

import (
	"strings"
	"testing"
)

func TestMsgInbound(t *testing.T) {
	tests := []struct {
		platform, tenantID, botID string
		expected                   string
	}{
		{"telegram", "tenant-123", "bot-456", "msg.inbound.telegram.tenant-123.bot-456"},
		{"whatsapp", "t1", "b1", "msg.inbound.whatsapp.t1.b1"},
		{"", "", "", "msg.inbound..."},
	}

	for _, tt := range tests {
		t.Run(tt.platform+"/"+tt.tenantID, func(t *testing.T) {
			result := MsgInbound(tt.platform, tt.tenantID, tt.botID)
			if result != tt.expected {
				t.Errorf("MsgInbound(%q, %q, %q) = %q, want %q", tt.platform, tt.tenantID, tt.botID, result, tt.expected)
			}
		})
	}
}

func TestMsgInboundFilter(t *testing.T) {
	tests := []struct {
		platform, tenantID, botID string
		expected                   string
		description                string
	}{
		{"telegram", ">", ">", "msg.inbound.telegram.>.>", "all tenants/bots for telegram"},
		{"telegram", "tenant-123", ">", "msg.inbound.telegram.tenant-123.>", "all bots for tenant"},
		{">", ">", ">", "msg.inbound.>.>.>", "all platforms/tenants/bots"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := MsgInboundFilter(tt.platform, tt.tenantID, tt.botID)
			if result != tt.expected {
				t.Errorf("MsgInboundFilter(%q, %q, %q) = %q, want %q", tt.platform, tt.tenantID, tt.botID, result, tt.expected)
			}
		})
	}
}

func TestMsgOutbound(t *testing.T) {
	tests := []struct {
		platform, tenantID, botID string
		expected                   string
	}{
		{"telegram", "tenant-123", "bot-456", "msg.outbound.telegram.tenant-123.bot-456"},
		{"whatsapp", "t2", "b2", "msg.outbound.whatsapp.t2.b2"},
	}

	for _, tt := range tests {
		t.Run(tt.platform+"/"+tt.tenantID, func(t *testing.T) {
			result := MsgOutbound(tt.platform, tt.tenantID, tt.botID)
			if result != tt.expected {
				t.Errorf("MsgOutbound(%q, %q, %q) = %q, want %q", tt.platform, tt.tenantID, tt.botID, result, tt.expected)
			}
		})
	}
}

func TestIngestCompleted(t *testing.T) {
	tests := []struct {
		tenantID string
		expected string
	}{
		{"tenant-123", "ingest.completed.tenant-123"},
		{"t1", "ingest.completed.t1"},
		{"", "ingest.completed."},
	}

	for _, tt := range tests {
		t.Run(tt.tenantID, func(t *testing.T) {
			result := IngestCompleted(tt.tenantID)
			if result != tt.expected {
				t.Errorf("IngestCompleted(%q) = %q, want %q", tt.tenantID, result, tt.expected)
			}
		})
	}
}

func TestIngestFailed(t *testing.T) {
	result := IngestFailed("tenant-123")
	expected := "ingest.failed.tenant-123"
	if result != expected {
		t.Errorf("IngestFailed(tenant-123) = %q, want %q", result, expected)
	}
}

func TestBillingSubscriptionUpdated(t *testing.T) {
	result := BillingSubscriptionUpdated("tenant-abc")
	expected := "billing.subscription.updated.tenant-abc"
	if result != expected {
		t.Errorf("BillingSubscriptionUpdated(tenant-abc) = %q, want %q", result, expected)
	}
}

func TestBillingUsageRecorded(t *testing.T) {
	result := BillingUsageRecorded("tenant-xyz")
	expected := "billing.usage.recorded.tenant-xyz"
	if result != expected {
		t.Errorf("BillingUsageRecorded(tenant-xyz) = %q, want %q", result, expected)
	}
}

func TestScheduleFire(t *testing.T) {
	tests := []struct {
		jobType, tenantID string
		expected          string
	}{
		{"daily_lesson", "tenant-123", "schedule.fire.daily_lesson.tenant-123"},
		{"weekly_recap", "t2", "schedule.fire.weekly_recap.t2"},
		{"reminder", "tenant-456", "schedule.fire.reminder.tenant-456"},
	}

	for _, tt := range tests {
		t.Run(tt.jobType, func(t *testing.T) {
			result := ScheduleFire(tt.jobType, tt.tenantID)
			if result != tt.expected {
				t.Errorf("ScheduleFire(%q, %q) = %q, want %q", tt.jobType, tt.tenantID, result, tt.expected)
			}
		})
	}
}

func TestScheduleFireFilter(t *testing.T) {
	tests := []struct {
		jobType  string
		expected string
	}{
		{"daily_lesson", "schedule.fire.daily_lesson.>"},
		{"weekly_recap", "schedule.fire.weekly_recap.>"},
	}

	for _, tt := range tests {
		t.Run(tt.jobType, func(t *testing.T) {
			result := ScheduleFireFilter(tt.jobType)
			if result != tt.expected {
				t.Errorf("ScheduleFireFilter(%q) = %q, want %q", tt.jobType, result, tt.expected)
			}
		})
	}
}

func TestIdentityBotUpdated(t *testing.T) {
	result := IdentityBotUpdated("bot-123")
	expected := "identity.bot.updated.bot-123"
	if result != expected {
		t.Errorf("IdentityBotUpdated(bot-123) = %q, want %q", result, expected)
	}
}

func TestIdentityBotDeleted(t *testing.T) {
	result := IdentityBotDeleted("bot-456")
	expected := "identity.bot.deleted.bot-456"
	if result != expected {
		t.Errorf("IdentityBotDeleted(bot-456) = %q, want %q", result, expected)
	}
}

func TestIdentityTenantUpdated(t *testing.T) {
	result := IdentityTenantUpdated("tenant-789")
	expected := "identity.tenant.updated.tenant-789"
	if result != expected {
		t.Errorf("IdentityTenantUpdated(tenant-789) = %q, want %q", result, expected)
	}
}

func TestCrawlSiteCompleted(t *testing.T) {
	result := CrawlSiteCompleted("tenant-123")
	expected := "crawl.site.completed.tenant-123"
	if result != expected {
		t.Errorf("CrawlSiteCompleted(tenant-123) = %q, want %q", result, expected)
	}
}

func TestCrawlSiteFailed(t *testing.T) {
	result := CrawlSiteFailed("tenant-456")
	expected := "crawl.site.failed.tenant-456"
	if result != expected {
		t.Errorf("CrawlSiteFailed(tenant-456) = %q, want %q", result, expected)
	}
}

func TestCrawlAgentPersonaUpdated(t *testing.T) {
	result := CrawlAgentPersonaUpdated("tenant-789")
	expected := "crawl.agent.persona_updated.tenant-789"
	if result != expected {
		t.Errorf("CrawlAgentPersonaUpdated(tenant-789) = %q, want %q", result, expected)
	}
}

func TestConversationCompactRequested(t *testing.T) {
	result := ConversationCompactRequested("tenant-123", "conv-456")
	expected := "conversation.compact.requested.tenant-123.conv-456"
	if result != expected {
		t.Errorf("ConversationCompactRequested(tenant-123, conv-456) = %q, want %q", result, expected)
	}
}

func TestConversationCompactRequestedFilter(t *testing.T) {
	result := ConversationCompactRequestedFilter()
	expected := "conversation.compact.requested.>"
	if result != expected {
		t.Errorf("ConversationCompactRequestedFilter() = %q, want %q", result, expected)
	}
}

func TestConversationCompacted(t *testing.T) {
	result := ConversationCompacted("tenant-789")
	expected := "conversation.compacted.tenant-789"
	if result != expected {
		t.Errorf("ConversationCompacted(tenant-789) = %q, want %q", result, expected)
	}
}

func TestRealtime(t *testing.T) {
	tests := []struct {
		tenantID, domain, resource, event string
		expected                           string
	}{
		{"tenant-123", "billing", "usage", "recorded", "tenant.tenant-123.billing.usage.recorded"},
		{"t1", "agent", "message", "new", "tenant.t1.agent.message.new"},
	}

	for _, tt := range tests {
		t.Run(tt.domain+"/"+tt.resource, func(t *testing.T) {
			result := Realtime(tt.tenantID, tt.domain, tt.resource, tt.event)
			if result != tt.expected {
				t.Errorf("Realtime(%q, %q, %q, %q) = %q, want %q", tt.tenantID, tt.domain, tt.resource, tt.event, result, tt.expected)
			}
		})
	}
}

func TestRealtimeFilter(t *testing.T) {
	tests := []struct {
		tenantID, domain, resource string
		expected                    string
	}{
		{"tenant-123", "billing", "usage", "tenant.tenant-123.billing.usage.>"},
		{"t1", "agent", ">", "tenant.t1.agent.>.>"},
	}

	for _, tt := range tests {
		t.Run(tt.domain+"/"+tt.resource, func(t *testing.T) {
			result := RealtimeFilter(tt.tenantID, tt.domain, tt.resource)
			if result != tt.expected {
				t.Errorf("RealtimeFilter(%q, %q, %q) = %q, want %q", tt.tenantID, tt.domain, tt.resource, result, tt.expected)
			}
		})
	}
}

// SubjectPrefixConsistency verifies that all subjects follow the dot-separated namespace convention.
func TestSubjectPrefixConsistency(t *testing.T) {
	prefixes := []string{
		MsgInbound("telegram", "t1", "b1"),
		MsgOutbound("telegram", "t1", "b1"),
		IngestCompleted("t1"),
		IngestFailed("t1"),
		BillingSubscriptionUpdated("t1"),
		BillingUsageRecorded("t1"),
		ScheduleFire("daily_lesson", "t1"),
		IdentityBotUpdated("b1"),
		IdentityBotDeleted("b1"),
		IdentityTenantUpdated("t1"),
		CrawlSiteCompleted("t1"),
		CrawlSiteFailed("t1"),
		CrawlAgentPersonaUpdated("t1"),
		ConversationCompactRequested("t1", "c1"),
		ConversationCompacted("t1"),
		Realtime("t1", "domain", "resource", "event"),
	}

	for _, subject := range prefixes {
		if strings.Contains(subject, " ") {
			t.Errorf("Subject contains spaces: %q", subject)
		}
		if !strings.Contains(subject, ".") {
			t.Errorf("Subject does not contain dots: %q", subject)
		}
		parts := strings.Split(subject, ".")
		if len(parts) < 2 {
			t.Errorf("Subject has fewer than 2 parts: %q", subject)
		}
	}
}

// StreamConstantsExist verifies that expected stream names are defined.
func TestStreamConstants(t *testing.T) {
	streams := map[string]string{
		"StreamMessages":     StreamMessages,
		"StreamIngest":       StreamIngest,
		"StreamCrawl":        StreamCrawl,
		"StreamBilling":      StreamBilling,
		"StreamSchedule":     StreamSchedule,
		"StreamIdentity":     StreamIdentity,
		"StreamConversation": StreamConversation,
	}

	for name, value := range streams {
		if value == "" {
			t.Errorf("Stream constant %s is empty", name)
		}
	}
}

// FilterConstantsExist verifies that filter subjects are defined.
func TestFilterConstants(t *testing.T) {
	filters := map[string]string{
		"FilterMessages":     FilterMessages,
		"FilterIngest":       FilterIngest,
		"FilterCrawl":        FilterCrawl,
		"FilterBilling":      FilterBilling,
		"FilterSchedule":     FilterSchedule,
		"FilterIdentity":     FilterIdentity,
		"FilterAgent":        FilterAgent,
		"FilterConversation": FilterConversation,
		"FilterRealtime":     FilterRealtime,
	}

	for name, value := range filters {
		if value == "" {
			t.Errorf("Filter constant %s is empty", name)
		}
		if !strings.Contains(value, ">") {
			t.Errorf("Filter constant %s does not contain wildcard: %q", name, value)
		}
	}
}
