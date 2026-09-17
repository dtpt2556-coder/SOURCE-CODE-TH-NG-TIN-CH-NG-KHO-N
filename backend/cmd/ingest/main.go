// Command ingest chạy pipeline đúng MỘT lần rồi thoát — dùng cho debug,
// backfill và kiểm thử thủ công của QA/DevOps.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tenpoint/tenpoint-api/internal/config"
	"github.com/tenpoint/tenpoint-api/internal/ingest"
	"github.com/tenpoint/tenpoint-api/internal/ingest/summarizer"
	"github.com/tenpoint/tenpoint-api/internal/store"
	"github.com/tenpoint/tenpoint-api/migrations"
)

func main() {
	trigger := flag.String("trigger", "manual", "nguồn kích hoạt: cron | manual")
	timeout := flag.Duration("timeout", 20*time.Minute, "thời gian tối đa cho một lần chạy")
	skipMigrate := flag.Bool("skip-migrate", false, "bỏ qua bước chạy migrations")
	flag.Parse()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(log)

	if err := run(log, *trigger, *timeout, *skipMigrate); err != nil {
		log.Error("ingest failed", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger, trigger string, timeout time.Duration, skipMigrate bool) error {
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithTimeout(rootCtx, timeout)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	st, err := store.New(ctx, cfg.DatabaseURL, log)
	if err != nil {
		return err
	}
	defer st.Close()

	if !skipMigrate {
		if err := st.Migrate(ctx, migrations.FS); err != nil {
			return err
		}
	}

	var provider summarizer.Provider = summarizer.NewExtractive()
	if cfg.SummarizerProvider == config.ProviderAnthropic && cfg.AnthropicAPIKey != "" {
		provider = summarizer.NewAnthropic(cfg.AnthropicAPIKey, cfg.AnthropicModel, 60*time.Second)
	}

	pipeline := ingest.New(ingest.Config{
		UserAgent:            cfg.UserAgent,
		FetchTimeout:         cfg.FetchTimeout,
		MaxArticlesPerSource: cfg.MaxArticlesPerSource,
		Workers:              ingest.DefaultWorkers,
		Location:             cfg.Location,
		IndustryTier:         cfg.IndustryTier,
	}, st, provider, log)

	res, err := pipeline.Run(ctx, trigger)
	if err != nil {
		if errors.Is(err, ingest.ErrAlreadyRunning) {
			return fmt.Errorf("cannot start: %w", err)
		}
		return err
	}

	log.Info("ingest finished",
		"run_id", res.RunID, "found", res.Found, "new", res.New,
		"rejected", res.Rejected, "errors", res.Errors, "duration", res.Duration.String())
	return nil
}
