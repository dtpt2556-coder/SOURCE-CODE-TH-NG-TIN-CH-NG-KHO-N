// Command api khởi động HTTP server + cron scheduler + ingest pipeline trong
// một binary duy nhất (ADR-002: monolith cho quy mô MVP).
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tenpoint/tenpoint-api/internal/config"
	"github.com/tenpoint/tenpoint-api/internal/domain"
	tphttp "github.com/tenpoint/tenpoint-api/internal/http"
	"github.com/tenpoint/tenpoint-api/internal/ingest"
	"github.com/tenpoint/tenpoint-api/internal/ingest/summarizer"
	"github.com/tenpoint/tenpoint-api/internal/scheduler"
	"github.com/tenpoint/tenpoint-api/internal/shortlink"
	"github.com/tenpoint/tenpoint-api/internal/store"
	"github.com/tenpoint/tenpoint-api/migrations"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(log)

	if err := run(log); err != nil {
		log.Error("process exited with error", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	// Graceful shutdown: SIGINT/SIGTERM huỷ ctx gốc.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log.Info("configuration loaded",
		"port", cfg.Port, "tz", cfg.TZ, "cron_enabled", cfg.CronEnabled,
		"cron_spec", cfg.CronSpec, "summarizer", cfg.SummarizerProvider)

	st, err := store.New(ctx, cfg.DatabaseURL, log)
	if err != nil {
		return err
	}
	defer st.Close()

	// Migration tự chạy khi khởi động — idempotent.
	migrateCtx, cancelMigrate := context.WithTimeout(ctx, 2*time.Minute)
	err = st.Migrate(migrateCtx, migrations.FS)
	cancelMigrate()
	if err != nil {
		return err
	}
	log.Info("migrations up to date")

	provider := buildSummarizer(cfg, log)
	pipeline := ingest.New(ingest.Config{
		UserAgent:            cfg.UserAgent,
		FetchTimeout:         cfg.FetchTimeout,
		MaxArticlesPerSource: cfg.MaxArticlesPerSource,
		Workers:              ingest.DefaultWorkers,
		Location:             cfg.Location,
		IndustryTier:         cfg.IndustryTier,
	}, st, provider, log)

	clicks := shortlink.NewClickCounter(st.FlushClicks, shortlink.FlushInterval, log)
	go clicks.Run(ctx)

	if cfg.BootstrapIngest {
		go bootstrapIngest(ctx, st, pipeline, log)
	}

	handler := tphttp.NewRouter(tphttp.Deps{
		Store:  st,
		Config: cfg,
		Log:    log,
		Runner: pipeline,
		Cache:  shortlink.NewLRUCache(shortlink.DefaultCacheSize),
		Clicks: clicks,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	var sched *scheduler.Scheduler
	if cfg.CronEnabled {
		sched, err = scheduler.New(cfg.CronSpec, cfg.Location, func(jobCtx context.Context) error {
			_, runErr := pipeline.Run(jobCtx, domain.TriggerCron)
			if errors.Is(runErr, ingest.ErrAlreadyRunning) {
				log.Warn("cron tick skipped, another pipeline run is in progress")
				return nil
			}
			return runErr
		}, log)
		if err != nil {
			return err
		}
		sched.Start()
	} else {
		log.Info("cron disabled (CRON_ENABLED=false)")
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Info("HTTP server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		log.Info("shutdown signal received, draining")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if sched != nil {
		sched.Stop(shutdownCtx)
	}
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	// Flush nốt click đang đệm trước khi đóng pool.
	clicks.Flush(shutdownCtx)
	log.Info("shutdown complete")
	return nil
}

// bootstrapIngest crawl ngay ở lần khởi động đầu tiên khi CSDL chưa có bài thật.
//
// Không có bước này, sau `docker compose up` người dùng chỉ thấy 5 bài seed với
// URL không tồn tại cho tới mốc cron kế tiếp — đúng lỗi khách hàng báo về link
// không dẫn tới bài cụ thể. Chạy nền để không chặn ListenAndServe.
func bootstrapIngest(ctx context.Context, st *store.Store, runner *ingest.Pipeline, log *slog.Logger) {
	hasReal, err := st.HasRealArticles(ctx)
	if err != nil {
		log.Error("bootstrap: cannot check for real articles", "err", err)
		return
	}
	if hasReal {
		return
	}

	log.Info("bootstrap: no real articles yet, running first ingest now")
	runCtx, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()

	res, err := runner.Run(runCtx, domain.TriggerManual)
	if err != nil {
		if errors.Is(err, ingest.ErrAlreadyRunning) {
			return
		}
		log.Error("bootstrap ingest failed", "err", err)
		return
	}
	log.Info("bootstrap ingest finished", "run_id", res.RunID, "new", res.New, "errors", res.Errors)
}

// buildSummarizer chọn provider theo cấu hình. Không có API key -> extractive,
// hệ thống vẫn chạy được end-to-end (ADR-007).
func buildSummarizer(cfg *config.Config, log *slog.Logger) summarizer.Provider {
	if cfg.SummarizerProvider == config.ProviderAnthropic && cfg.AnthropicAPIKey != "" {
		log.Info("using Anthropic summarizer", "model", cfg.AnthropicModel)
		return summarizer.NewAnthropic(cfg.AnthropicAPIKey, cfg.AnthropicModel, 60*time.Second)
	}
	log.Info("using extractive summarizer (no API key required)")
	return summarizer.NewExtractive()
}
