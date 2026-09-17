// Package config nạp và validate cấu hình từ biến môi trường.
// Xem mục 12 của docs/03-sa/architecture.md.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Provider names hợp lệ cho SUMMARIZER_PROVIDER.
const (
	ProviderExtractive = "extractive"
	ProviderAnthropic  = "anthropic"
)

// Config là toàn bộ cấu hình runtime của service. Không có biến global mutable:
// mọi thành phần nhận Config qua dependency injection.
type Config struct {
	Port        string
	DatabaseURL string

	TZ       string
	Location *time.Location

	CronSpec    string
	CronEnabled bool

	SummarizerProvider string
	AnthropicAPIKey    string
	AnthropicModel     string

	AdminToken    string
	PublicBaseURL string

	FetchTimeout         time.Duration
	MaxArticlesPerSource int
	UserAgent            string

	// ReadyMaxAge: dữ liệu cũ hơn ngưỡng này thì /readyz trả 503.
	ReadyMaxAge time.Duration

	// BootstrapIngest: chạy ingest ngay lúc khởi động nếu CSDL chưa có bài thật,
	// để người dùng thấy tin thật kèm link thật mà không phải chờ tới mốc cron
	// kế tiếp (T1 mục 0 của ingest-edge-cases).
	BootstrapIngest bool

	// IndustryTier bật tầng T3 gắn mã theo ngữ cảnh ngành. Mặc định TẮT: đo
	// trên dữ liệu thật, 38/38 tag của tầng này đều sai. Chỉ bật lại khi đạt
	// precision trên corpus do QA tạo.
	IndustryTier bool
}

// Load đọc file .env (nếu có) rồi nạp cấu hình từ môi trường.
// Mọi biến đều có default an toàn trừ DATABASE_URL.
func Load() (*Config, error) {
	// .env là tuỳ chọn — thiếu file không phải lỗi.
	_ = godotenv.Load()

	tz := env("TZ", "Asia/Ho_Chi_Minh")
	loc, err := loadLocation(tz)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Port:                 env("PORT", "8080"),
		DatabaseURL:          env("DATABASE_URL", ""),
		TZ:                   tz,
		Location:             loc,
		CronSpec:             env("CRON_SPEC", "0 6,14,22 * * *"),
		CronEnabled:          envBool("CRON_ENABLED", true),
		SummarizerProvider:   strings.ToLower(env("SUMMARIZER_PROVIDER", ProviderExtractive)),
		AnthropicAPIKey:      env("ANTHROPIC_API_KEY", ""),
		AnthropicModel:       env("ANTHROPIC_MODEL", "claude-sonnet-5"),
		AdminToken:           env("ADMIN_TOKEN", ""),
		PublicBaseURL:        strings.TrimRight(env("PUBLIC_BASE_URL", "http://localhost"), "/"),
		FetchTimeout:         time.Duration(envInt("FETCH_TIMEOUT_SEC", 20)) * time.Second,
		MaxArticlesPerSource: envInt("MAX_ARTICLES_PER_SOURCE", 40),
		UserAgent:            env("USER_AGENT", "TenPointBot/1.0 (+https://tenpoint.vn/bot)"),
		ReadyMaxAge:          time.Duration(envInt("READY_MAX_AGE_HOURS", 26)) * time.Hour,
		BootstrapIngest:      envBool("BOOTSTRAP_INGEST", true),
		IndustryTier:         envBool("TAGGER_INDUSTRY_TIER", false),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate kiểm tra các ràng buộc bắt buộc.
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("missing required environment variable DATABASE_URL")
	}
	switch c.SummarizerProvider {
	case ProviderExtractive, ProviderAnthropic:
	default:
		return fmt.Errorf("invalid SUMMARIZER_PROVIDER=%q (accepted: %q, %q)",
			c.SummarizerProvider, ProviderExtractive, ProviderAnthropic)
	}
	if c.SummarizerProvider == ProviderAnthropic && c.AnthropicAPIKey == "" {
		return fmt.Errorf("SUMMARIZER_PROVIDER=anthropic requires ANTHROPIC_API_KEY")
	}
	if c.MaxArticlesPerSource <= 0 {
		return fmt.Errorf("MAX_ARTICLES_PER_SOURCE must be > 0, got %d", c.MaxArticlesPerSource)
	}
	if c.FetchTimeout <= 0 {
		return fmt.Errorf("FETCH_TIMEOUT_SEC must be > 0")
	}
	return nil
}

// loadLocation nạp múi giờ và dừng ngay nếu image thiếu tzdata.
//
// Im lặng rơi về UTC là lỗi tệ nhất có thể xảy ra ở đây: cron sẽ chạy lệch 7
// tiếng và mọi ngày trên bảng digest sẽ sai, mà không có dấu hiệu nào trong log
// (R-05).
func loadLocation(tz string) (*time.Location, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf(
			"load timezone %q: %w (image is missing tzdata; install the tzdata package "+
				"or import _ \"time/tzdata\" — refusing to start rather than silently running on UTC)", tz, err)
	}
	if tz == "Asia/Ho_Chi_Minh" {
		if _, offset := time.Now().In(loc).Zone(); offset != 7*3600 {
			return nil, fmt.Errorf(
				"timezone %q resolved to UTC offset %+d seconds instead of +25200: zoneinfo database looks broken", tz, offset)
		}
	}
	return loc, nil
}

func env(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}
	return def
}

func envInt(key string, def int) int {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	n, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return def
	}
	return n
}

func envBool(key string, def bool) bool {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	b, err := strconv.ParseBool(strings.TrimSpace(raw))
	if err != nil {
		return def
	}
	return b
}
