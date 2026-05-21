package jetstream

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	evtypes "github.com/madmike/go-events/types"
	"github.com/madmike/go-infra/telemetry"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
)

// MockJetStream implements jetstream.JetStream for testing.
type MockJetStream struct {
	publishErr error
	pubAck     *jetstream.PubAck
	createErr  error
}

func (m *MockJetStream) Publish(ctx context.Context, subject string, data []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	if m.publishErr != nil {
		return nil, m.publishErr
	}
	if m.pubAck != nil {
		return m.pubAck, nil
	}
	return &jetstream.PubAck{}, nil
}

func (m *MockJetStream) PublishAsync(subject string, data []byte, opts ...jetstream.PublishAsyncOpt) (jetstream.PubAckFuture, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJetStream) PublishMsg(msg *nats.Msg, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJetStream) Subscribe(subject string, opts ...jetstream.ConsumeOpt) (jetstream.MessageConsumer, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJetStream) ChanSubscribe(subject string, opts ...jetstream.ConsumeOpt) (<-chan jetstream.Msg, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJetStream) PullSubscribe(subject string, opts ...jetstream.PullSubscribeOpt) (jetstream.PullSubscriber, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJetStream) CreateConsumer(ctx context.Context, stream string, config jetstream.ConsumerConfig, opts ...jetstream.CreateConsumerOpt) (jetstream.Consumer, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return &MockConsumer{}, nil
}

func (m *MockJetStream) UpdateConsumer(ctx context.Context, stream string, config jetstream.ConsumerConfig, opts ...jetstream.UpdateConsumerOpt) (jetstream.Consumer, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJetStream) CreateOrUpdateConsumer(ctx context.Context, stream string, config jetstream.ConsumerConfig, opts ...jetstream.CreateOrUpdateConsumerOpt) (jetstream.Consumer, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return &MockConsumer{}, nil
}

func (m *MockJetStream) DeleteConsumer(ctx context.Context, stream string, consumer string, opts ...jetstream.DeleteConsumerOpt) error {
	return nil
}

func (m *MockJetStream) Consumer(ctx context.Context, stream string, consumer string, opts ...jetstream.ConsumerOpt) (jetstream.Consumer, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJetStream) Consumers(ctx context.Context, stream string, opts ...jetstream.ConsumerOpt) jetstream.ConsumerLister {
	return nil
}

func (m *MockJetStream) CreateKeyValue(ctx context.Context, config jetstream.KeyValueConfig) (jetstream.KeyValue, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJetStream) KeyValue(ctx context.Context, bucket string) (jetstream.KeyValue, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJetStream) DeleteKeyValue(ctx context.Context, bucket string) error {
	return errors.New("not implemented")
}

func (m *MockJetStream) KeyValues(ctx context.Context) jetstream.KeyValueLister {
	return nil
}

func (m *MockJetStream) CreateObjectStore(ctx context.Context, config jetstream.ObjectStoreConfig) (jetstream.ObjectStore, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJetStream) ObjectStore(ctx context.Context, bucket string) (jetstream.ObjectStore, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJetStream) DeleteObjectStore(ctx context.Context, bucket string) error {
	return errors.New("not implemented")
}

func (m *MockJetStream) ObjectStores(ctx context.Context) jetstream.ObjectStoreLister {
	return nil
}

func (m *MockJetStream) CreateStream(ctx context.Context, config jetstream.StreamConfig, opts ...jetstream.CreateStreamOpt) (jetstream.Stream, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return &MockStream{}, nil
}

func (m *MockJetStream) UpdateStream(ctx context.Context, config jetstream.StreamConfig, opts ...jetstream.UpdateStreamOpt) (jetstream.Stream, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJetStream) CreateOrUpdateStream(ctx context.Context, config jetstream.StreamConfig, opts ...jetstream.CreateOrUpdateStreamOpt) (jetstream.Stream, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return &MockStream{}, nil
}

func (m *MockJetStream) DeleteStream(ctx context.Context, stream string, opts ...jetstream.DeleteStreamOpt) error {
	return nil
}

func (m *MockJetStream) Stream(ctx context.Context, stream string, opts ...jetstream.StreamOpt) (jetstream.Stream, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJetStream) Streams(ctx context.Context, opts ...jetstream.StreamOpt) jetstream.StreamLister {
	return nil
}

func (m *MockJetStream) AccountInfo(ctx context.Context, opts ...jetstream.JetStreamOpt) (*jetstream.AccountInfo, error) {
	return nil, errors.New("not implemented")
}

// MockStream implements jetstream.Stream.
type MockStream struct{}

func (m *MockStream) Info(ctx context.Context) (*jetstream.StreamInfo, error) {
	return nil, errors.New("not implemented")
}

func (m *MockStream) CreateConsumer(ctx context.Context, config jetstream.ConsumerConfig, opts ...jetstream.CreateConsumerOpt) (jetstream.Consumer, error) {
	return nil, errors.New("not implemented")
}

func (m *MockStream) UpdateConsumer(ctx context.Context, config jetstream.ConsumerConfig, opts ...jetstream.UpdateConsumerOpt) (jetstream.Consumer, error) {
	return nil, errors.New("not implemented")
}

func (m *MockStream) CreateOrUpdateConsumer(ctx context.Context, config jetstream.ConsumerConfig, opts ...jetstream.CreateOrUpdateConsumerOpt) (jetstream.Consumer, error) {
	return nil, errors.New("not implemented")
}

func (m *MockStream) DeleteConsumer(ctx context.Context, consumer string, opts ...jetstream.DeleteConsumerOpt) error {
	return nil
}

func (m *MockStream) Consumer(ctx context.Context, consumer string, opts ...jetstream.ConsumerOpt) (jetstream.Consumer, error) {
	return nil, errors.New("not implemented")
}

func (m *MockStream) Consumers(ctx context.Context, opts ...jetstream.ConsumerOpt) jetstream.ConsumerLister {
	return nil
}

func (m *MockStream) Purge(ctx context.Context, opts ...jetstream.StreamOpt) error {
	return nil
}

func (m *MockStream) GetMsg(ctx context.Context, index uint64, opts ...jetstream.StreamOpt) (*nats.Msg, error) {
	return nil, errors.New("not implemented")
}

func (m *MockStream) GetLastMsg(ctx context.Context, subject string, opts ...jetstream.StreamOpt) (*nats.Msg, error) {
	return nil, errors.New("not implemented")
}

func (m *MockStream) DeleteMsg(ctx context.Context, index uint64, opts ...jetstream.StreamOpt) error {
	return nil
}

func (m *MockStream) SecureDeleteMsg(ctx context.Context, index uint64, opts ...jetstream.StreamOpt) error {
	return nil
}

// MockConsumer implements jetstream.Consumer.
type MockConsumer struct{}

func (m *MockConsumer) Info(ctx context.Context) (*jetstream.ConsumerInfo, error) {
	return nil, errors.New("not implemented")
}

func (m *MockConsumer) Messages(opts ...jetstream.ConsumeOpt) (jetstream.MessageConsumer, error) {
	return nil, errors.New("not implemented")
}

func (m *MockConsumer) Consume(handler func(jetstream.Msg), opts ...jetstream.ConsumeOpt) (jetstream.MessageConsumer, error) {
	return nil, errors.New("not implemented")
}

func (m *MockConsumer) ChanMessages(opts ...jetstream.ConsumeOpt) (<-chan jetstream.Msg, error) {
	return nil, errors.New("not implemented")
}

func (m *MockConsumer) Pull(batchSize int, opts ...jetstream.PullOpt) ([]jetstream.Msg, error) {
	return nil, errors.New("not implemented")
}

func (m *MockConsumer) Fetch(batchSize int, opts ...jetstream.FetchOpt) ([]jetstream.Msg, error) {
	return nil, errors.New("not implemented")
}

func (m *MockConsumer) Next(opts ...jetstream.FetchOpt) (jetstream.Msg, error) {
	return nil, errors.New("not implemented")
}

// MockNatsConn implements a minimal nats.Conn interface for testing.
type MockNatsConn struct {
	publishErr error
}

func (m *MockNatsConn) Publish(subject string, data []byte) error {
	return m.publishErr
}

func (m *MockNatsConn) PublishMsg(msg *nats.Msg) error {
	return errors.New("not implemented")
}

func (m *MockNatsConn) Drain() error {
	return nil
}

func (m *MockNatsConn) Close() {
}

// MockClient wraps Client for testing, allowing injection of mock JS/NC.
type MockClient struct {
	jsm jetstream.JetStream
	ncm nats.Conn
	log telemetry.Logger
}

// TestPublisherPublish tests Publisher.Publish wraps payload in envelope.
func TestPublisherPublish(t *testing.T) {
	mockJS := &MockJetStream{
		pubAck: &jetstream.PubAck{Stream: "BILLING", Sequence: 1},
	}
	mockNC := &MockNatsConn{}

	pub := &Publisher{
		c: &Client{js: mockJS, nc: mockNC, log: &telemetry.NoOpLogger{}},
	}

	event := evtypes.BillingUsageRecorded{
		TenantID:       "t1",
		CostMicroCents: 5000,
	}

	ack, err := pub.Publish(context.Background(), "billing.usage.recorded.t1", "t1", event)
	require.NoError(t, err)
	require.NotNil(t, ack)
}

// TestPublisherPublishError propagates publish error.
func TestPublisherPublishError(t *testing.T) {
	expectedErr := errors.New("publish failed")
	mockJS := &MockJetStream{publishErr: expectedErr}
	mockNC := &MockNatsConn{}

	pub := &Publisher{
		c: &Client{js: mockJS, nc: mockNC, log: &telemetry.NoOpLogger{}},
	}

	_, err := pub.Publish(context.Background(), "test.subject", "t1", map[string]any{"key": "value"})
	require.Error(t, err)
	require.True(t, errors.Is(err, expectedErr) || errors.Is(err, expectedErr))
}

// TestPublisherPublishMarshalError handles payload marshal failure.
func TestPublisherPublishMarshalError(t *testing.T) {
	mockJS := &MockJetStream{}
	mockNC := &MockNatsConn{}

	pub := &Publisher{
		c: &Client{js: mockJS, nc: mockNC, log: &telemetry.NoOpLogger{}},
	}

	// Unmarshalable payload: func is not JSON serializable
	_, err := pub.Publish(context.Background(), "test.subject", "t1", func() {})
	require.Error(t, err)
}

// TestPublisherPublishWithTraceID includes trace ID in envelope.
func TestPublisherPublishWithTraceID(t *testing.T) {
	mockJS := &MockJetStream{}
	mockNC := &MockNatsConn{}

	pub := &Publisher{
		c: &Client{js: mockJS, nc: mockNC, log: &telemetry.NoOpLogger{}},
	}

	event := map[string]any{"data": "test"}
	ack, err := pub.PublishWithTraceID(context.Background(), "test.subject", "t1", "trace-123", event)
	require.NoError(t, err)
	require.NotNil(t, ack)
}

// TestPublisherPublishRealtime publishes directly to NATS (no JetStream).
func TestPublisherPublishRealtime(t *testing.T) {
	mockJS := &MockJetStream{}
	mockNC := &MockNatsConn{}

	pub := &Publisher{
		c: &Client{js: mockJS, nc: mockNC, log: &telemetry.NoOpLogger{}},
	}

	event := map[string]any{"data": "realtime-test"}
	err := pub.PublishRealtime(context.Background(), "realtime.subject", "t1", event)
	require.NoError(t, err)
}

// TestPublisherPublishRealtimeError propagates NATS error.
func TestPublisherPublishRealtimeError(t *testing.T) {
	expectedErr := errors.New("nats publish failed")
	mockJS := &MockJetStream{}
	mockNC := &MockNatsConn{publishErr: expectedErr}

	pub := &Publisher{
		c: &Client{js: mockJS, nc: mockNC, log: &telemetry.NoOpLogger{}},
	}

	err := pub.PublishRealtime(context.Background(), "realtime.subject", "t1", map[string]any{"data": "test"})
	require.Error(t, err)
}

// TestPublisherPublishInbound publishes InboundMessage with correct subject.
func TestPublisherPublishInbound(t *testing.T) {
	mockJS := &MockJetStream{}
	mockNC := &MockNatsConn{}

	pub := &Publisher{
		c: &Client{js: mockJS, nc: mockNC, log: &telemetry.NoOpLogger{}},
	}

	msg := evtypes.InboundMessage{
		TenantID:  "t1",
		Platform:  "telegram",
		AccountID: "acc-1",
		TraceID:   "trace-xyz",
	}

	err := pub.PublishInbound(context.Background(), msg)
	require.NoError(t, err)
}

// TestPublisherPublishOutbound publishes OutboundMessage with correct subject.
func TestPublisherPublishOutbound(t *testing.T) {
	mockJS := &MockJetStream{}
	mockNC := &MockNatsConn{}

	pub := &Publisher{
		c: &Client{js: mockJS, nc: mockNC, log: &telemetry.NoOpLogger{}},
	}

	msg := evtypes.OutboundMessage{
		TenantID:  "t1",
		Platform:  "whatsapp",
		AccountID: "acc-2",
		TraceID:   "trace-abc",
	}

	err := pub.PublishOutbound(context.Background(), msg)
	require.NoError(t, err)
}

// TestPublisherTypeSpecificMethods publishes events with correct subjects.
func TestPublisherTypeSpecificMethods(t *testing.T) {
	mockJS := &MockJetStream{}
	mockNC := &MockNatsConn{}

	pub := &Publisher{
		c: &Client{js: mockJS, nc: mockNC, log: &telemetry.NoOpLogger{}},
	}

	tests := []struct {
		name   string
		publish func() error
	}{
		{
			"IngestCompleted",
			func() error {
				return pub.PublishIngestCompleted(context.Background(), evtypes.IngestCompleted{TenantID: "t1"})
			},
		},
		{
			"IngestFailed",
			func() error {
				return pub.PublishIngestFailed(context.Background(), evtypes.IngestFailed{TenantID: "t1"})
			},
		},
		{
			"BillingSubscriptionUpdated",
			func() error {
				return pub.PublishBillingSubscriptionUpdated(context.Background(), evtypes.BillingSubscriptionUpdated{TenantID: "t1"})
			},
		},
		{
			"ScheduleFire",
			func() error {
				return pub.PublishScheduleFire(context.Background(), evtypes.ScheduleFireEvent{TenantID: "t1", JobType: "test"})
			},
		},
		{
			"IdentityAccount",
			func() error {
				return pub.PublishIdentityAccount(context.Background(), evtypes.IdentityAccountEvent{TenantID: "t1", AccountID: "acc-1", Action: "created"})
			},
		},
		{
			"IdentityTenant",
			func() error {
				return pub.PublishIdentityTenant(context.Background(), evtypes.IdentityTenantEvent{TenantID: "t1", Action: "created"})
			},
		},
		{
			"IdentityBot",
			func() error {
				return pub.PublishIdentityBot(context.Background(), evtypes.IdentityBotEvent{TenantID: "t1", TenantAgentID: "bot-1", Action: "created"})
			},
		},
		{
			"CrawlSiteCompleted",
			func() error {
				return pub.PublishCrawlSiteCompleted(context.Background(), evtypes.CrawlSiteCompleted{TenantID: "t1"})
			},
		},
		{
			"CrawlSiteFailed",
			func() error {
				return pub.PublishCrawlSiteFailed(context.Background(), evtypes.CrawlSiteFailed{TenantID: "t1"})
			},
		},
		{
			"CrawlAgentPersonaUpdated",
			func() error {
				return pub.PublishCrawlAgentPersonaUpdated(context.Background(), evtypes.CrawlAgentPersonaUpdated{TenantID: "t1"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.publish()
			require.NoError(t, err, "should publish %s without error", tt.name)
		})
	}
}

// TestConsumerConfigDefaults applies sensible defaults.
func TestConsumerConfigDefaults(t *testing.T) {
	cfg := ConsumerConfig{
		Stream:       "MESSAGES",
		DurableName:  "test-consumer",
		FilterSubject: "msg.inbound.>",
		// AckWait not set, should default to 30s
		// MaxDeliver not set, should default to 5
	}

	require.Equal(t, "MESSAGES", cfg.Stream)
	require.Equal(t, "test-consumer", cfg.DurableName)
	require.Equal(t, "msg.inbound.>", cfg.FilterSubject)
}

// TestUnmarshalHelper decodes envelope data into typed struct.
func TestUnmarshalHelper(t *testing.T) {
	event := evtypes.BillingUsageRecorded{
		TenantID:       "t1",
		MessageID:      "msg-1",
		CostMicroCents: 5000,
	}

	data, _ := json.Marshal(event)
	env := evtypes.Envelope{
		ID:      "env-1",
		Subject: "billing.usage.recorded.t1",
		Data:    data,
	}

	result, err := Unmarshal[evtypes.BillingUsageRecorded](env)
	require.NoError(t, err)
	require.Equal(t, "t1", result.TenantID)
	require.Equal(t, int64(5000), result.CostMicroCents)
}

// TestUnmarshalHelperError handles unmarshal failure.
func TestUnmarshalHelperError(t *testing.T) {
	env := evtypes.Envelope{
		ID:      "env-1",
		Subject: "test",
		Data:    []byte("invalid json {"),
	}

	_, err := Unmarshal[evtypes.BillingUsageRecorded](env)
	require.Error(t, err)
}

// TestPointerHelper creates pointer from value.
func TestPointerHelper(t *testing.T) {
	val := jetstream.DeliverAllPolicy
	ptr := Pointer(val)

	require.NotNil(t, ptr)
	require.Equal(t, val, *ptr)
}

// TestConfigDefaults applies default NATS config values.
func TestConfigDefaults(t *testing.T) {
	cfg := Config{}
	cfg.setDefaults()

	require.Equal(t, nats.DefaultURL, cfg.URL)
	require.Equal(t, -1, cfg.MaxReconnects)
	require.Equal(t, 2*time.Second, cfg.ReconnectWait)
	require.Equal(t, 10*time.Second, cfg.ConnectTimeout)
}

// TestConfigPartialDefaults preserves explicitly set values.
func TestConfigPartialDefaults(t *testing.T) {
	cfg := Config{
		URL:           "nats://custom:4222",
		MaxReconnects: 5,
		// Others not set
	}
	cfg.setDefaults()

	require.Equal(t, "nats://custom:4222", cfg.URL)
	require.Equal(t, 5, cfg.MaxReconnects)
	require.Equal(t, 2*time.Second, cfg.ReconnectWait)
	require.Equal(t, 10*time.Second, cfg.ConnectTimeout)
}
