package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/tenpoint/tenpoint-api/internal/config"
	"github.com/tenpoint/tenpoint-api/internal/ingest"
	"github.com/tenpoint/tenpoint-api/internal/shortlink"
	"github.com/tenpoint/tenpoint-api/internal/store"
)

// Runner là hợp đồng tối thiểu server cần từ pipeline (cho phép mock trong test).
type Runner interface {
	Run(ctx context.Context, trigger string) (ingest.Result, error)
}

// Deps là toàn bộ dependency của tầng HTTP, inject qua struct.
type Deps struct {
	Store  *store.Store
	Config *config.Config
	Log    *slog.Logger
	Runner Runner
	Cache  *shortlink.LRUCache
	Clicks *shortlink.ClickCounter
}

// Server giữ dependency và cài các handler.
type Server struct {
	store     *store.Store
	cfg       *config.Config
	log       *slog.Logger
	runner    Runner
	cache     *shortlink.LRUCache
	clicks    *shortlink.ClickCounter
	redirectR *rateLimiter
}

// NewRouter dựng router chi với đầy đủ middleware và route của REST v1.
func NewRouter(d Deps) http.Handler {
	log := d.Log
	if log == nil {
		log = slog.Default()
	}
	cache := d.Cache
	if cache == nil {
		cache = shortlink.NewLRUCache(shortlink.DefaultCacheSize)
	}

	s := &Server{
		store:     d.Store,
		cfg:       d.Config,
		log:       log,
		runner:    d.Runner,
		cache:     cache,
		clicks:    d.Clicks,
		redirectR: newRateLimiter(60, time.Minute),
	}

	r := chi.NewRouter()
	r.Use(Recoverer(log))
	r.Use(RequestLogger(log))
	r.Use(SecurityHeaders)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, CodeNotFound, "Không tìm thấy đường dẫn này.")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusMethodNotAllowed, CodeInvalidRequest,
			"Phương thức không được hỗ trợ cho đường dẫn này.")
	})

	// Probe phải nhận cả GET lẫn HEAD: `wget --spider` và phần lớn healthcheck
	// của container gửi HEAD, router chỉ đăng ký GET sẽ trả 405 và container
	// không bao giờ chuyển sang healthy.
	for _, path := range []string{"/healthz", "/readyz"} {
		handler := s.handleHealthz
		if path == "/readyz" {
			handler = s.handleReadyz
		}
		r.Get(path, handler)
		r.Head(path, handler)
	}
	r.Get("/r/{code}", s.handleRedirect)

	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/news", s.handleListNews)
		api.Get("/news/{id}", s.handleGetNews)
		api.Get("/tickers", s.handleListTickers)
		api.Get("/tickers/{symbol}", s.handleGetTicker)
		api.Get("/notes", s.handleListNotes)
		api.Get("/notes/{slug}", s.handleGetNote)
		api.Get("/meta", s.handleMeta)
		api.Post("/admin/ingest", s.handleAdminIngest)
	})

	return r
}

// location trả về múi giờ hiển thị (Asia/Ho_Chi_Minh).
func (s *Server) location() *time.Location {
	if s.cfg != nil && s.cfg.Location != nil {
		return s.cfg.Location
	}
	return time.UTC
}
