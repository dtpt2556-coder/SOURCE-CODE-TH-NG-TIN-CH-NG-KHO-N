package http

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"sync"
	"time"
)

// statusWriter ghi lại status code để log.
type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

// RequestLogger ghi log structured cho mỗi request.
func RequestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := &statusWriter{ResponseWriter: w}
			next.ServeHTTP(sw, r)
			if sw.status == 0 {
				sw.status = http.StatusOK
			}
			log.Info("http",
				"method", r.Method,
				"path", r.URL.Path,
				"status", sw.status,
				"bytes", sw.bytes,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}

// Recoverer biến panic thành 500 JSON thay vì làm sập server.
func Recoverer(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic while handling request",
						"path", r.URL.Path, "panic", rec, "stack", string(debug.Stack()))
					writeError(w, http.StatusInternalServerError, CodeInternal,
						"Có lỗi xảy ra, vui lòng thử lại sau.")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// rateLimiter là bộ đếm cửa sổ trượt đơn giản theo IP, dùng cho /r/{code}
// để chặn enumeration (§7.4 BA2).
type rateLimiter struct {
	mu       sync.Mutex
	hits     map[string][]time.Time
	limit    int
	window   time.Duration
	lastGC   time.Time
	maxHosts int
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		hits: make(map[string][]time.Time), limit: limit, window: window,
		lastGC: time.Now(), maxHosts: 50_000,
	}
}

func (rl *rateLimiter) allow(key string) bool {
	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if now.Sub(rl.lastGC) > rl.window || len(rl.hits) > rl.maxHosts {
		for k, ts := range rl.hits {
			if len(ts) == 0 || now.Sub(ts[len(ts)-1]) > rl.window {
				delete(rl.hits, k)
			}
		}
		rl.lastGC = now
	}

	cutoff := now.Add(-rl.window)
	kept := rl.hits[key][:0]
	for _, t := range rl.hits[key] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= rl.limit {
		rl.hits[key] = kept
		return false
	}
	rl.hits[key] = append(kept, now)
	return true
}

// SecurityHeaders đặt các header an toàn tối thiểu ở tầng ứng dụng
// (CSP đầy đủ do Caddy đảm nhiệm).
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}
