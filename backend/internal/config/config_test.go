package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/config"
)

// isolate xoá mọi biến môi trường của TenPoint để test không phụ thuộc môi
// trường máy chạy (CI, máy dev có sẵn ANTHROPIC_MODEL...).
func isolate(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"PORT", "DATABASE_URL", "TZ", "CRON_SPEC", "CRON_ENABLED",
		"SUMMARIZER_PROVIDER", "ANTHROPIC_API_KEY", "ANTHROPIC_MODEL",
		"ADMIN_TOKEN", "PUBLIC_BASE_URL", "FETCH_TIMEOUT_SEC",
		"MAX_ARTICLES_PER_SOURCE", "USER_AGENT", "READY_MAX_AGE_HOURS",
	} {
		t.Setenv(key, "")
	}
}

func TestLoadDefaults(t *testing.T) {
	isolate(t)
	t.Setenv("DATABASE_URL", "postgres://tenpoint:secret@db:5432/tenpoint?sslmode=disable")

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "Asia/Ho_Chi_Minh", cfg.TZ)
	assert.Equal(t, "0 6,14,22 * * *", cfg.CronSpec)
	assert.True(t, cfg.CronEnabled)
	assert.Equal(t, config.ProviderExtractive, cfg.SummarizerProvider,
		"mặc định phải là extractive — hệ thống chạy được không cần API key (ADR-007)")
	assert.Equal(t, "claude-sonnet-5", cfg.AnthropicModel)
	assert.Equal(t, 20*time.Second, cfg.FetchTimeout)
	assert.Equal(t, 40, cfg.MaxArticlesPerSource)
	assert.Equal(t, "TenPointBot/1.0 (+https://tenpoint.vn/bot)", cfg.UserAgent)
	assert.Equal(t, "http://localhost", cfg.PublicBaseURL)
	require.NotNil(t, cfg.Location)
}

func TestLoadRequiresDatabaseURL(t *testing.T) {
	isolate(t)
	t.Setenv("DATABASE_URL", "")
	_, err := config.Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DATABASE_URL")
}

func TestCronCanBeDisabled(t *testing.T) {
	isolate(t)
	t.Setenv("DATABASE_URL", "postgres://x/y")
	t.Setenv("CRON_ENABLED", "false")

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.False(t, cfg.CronEnabled)
}

func TestAnthropicProviderRequiresAPIKey(t *testing.T) {
	isolate(t)
	t.Setenv("DATABASE_URL", "postgres://x/y")
	t.Setenv("SUMMARIZER_PROVIDER", "anthropic")
	t.Setenv("ANTHROPIC_API_KEY", "")

	_, err := config.Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ANTHROPIC_API_KEY")
}

func TestUnknownProviderIsRejected(t *testing.T) {
	isolate(t)
	t.Setenv("DATABASE_URL", "postgres://x/y")
	t.Setenv("SUMMARIZER_PROVIDER", "openai")

	_, err := config.Load()
	assert.Error(t, err)
}

func TestPublicBaseURLTrailingSlashStripped(t *testing.T) {
	isolate(t)
	t.Setenv("DATABASE_URL", "postgres://x/y")
	t.Setenv("PUBLIC_BASE_URL", "https://tenpoint.vn/")

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "https://tenpoint.vn", cfg.PublicBaseURL)
}

func TestNoSecretsHardcoded(t *testing.T) {
	isolate(t)
	t.Setenv("DATABASE_URL", "postgres://x/y")
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Empty(t, cfg.AnthropicAPIKey, "không được hardcode secret")
	assert.Empty(t, cfg.AdminToken, "ADMIN_TOKEN rỗng => /admin/ingest bị khoá")
}
