package jetstream

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConfigSetDefaults(t *testing.T) {
	cfg := Config{}
	cfg.setDefaults()

	require.Equal(t, "nats://127.0.0.1:4222", cfg.URL)
	require.Equal(t, -1, cfg.MaxReconnects)
	require.Equal(t, 2*time.Second, cfg.ReconnectWait)
	require.Equal(t, 10*time.Second, cfg.ConnectTimeout)
}

func TestConfigSetDefaultsPreservesProvidedValues(t *testing.T) {
	cfg := Config{
		URL:            "nats://nats:4222",
		MaxReconnects:  5,
		ReconnectWait:  3 * time.Second,
		ConnectTimeout: 7 * time.Second,
	}
	cfg.setDefaults()

	require.Equal(t, "nats://nats:4222", cfg.URL)
	require.Equal(t, 5, cfg.MaxReconnects)
	require.Equal(t, 3*time.Second, cfg.ReconnectWait)
	require.Equal(t, 7*time.Second, cfg.ConnectTimeout)
}

func TestParseDuration(t *testing.T) {
	d, err := parseDuration("45s")
	require.NoError(t, err)
	require.Equal(t, 45*time.Second, d)
}

func TestParseDurationInvalid(t *testing.T) {
	_, err := parseDuration("not-a-duration")
	require.Error(t, err)
}

func TestPointerHelper(t *testing.T) {
	v := 42
	p := Pointer(v)
	require.NotNil(t, p)
	require.Equal(t, 42, *p)
}
