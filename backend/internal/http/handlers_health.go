package http

import (
	"context"
	"net/http"
	"time"
)

// handleHealthz là liveness probe — không chạm DB.
func (s *Server) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleReadyz là readiness probe: kiểm tra DB + độ tươi dữ liệu.
//
// Quy ước: khi hệ thống mới cài và CHƯA có crawl run nào, dữ liệu seed vẫn hợp
// lệ nên trả 200 (is_stale=false). Chỉ báo 503 khi đã từng crawl nhưng lần crawl
// gần nhất quá cũ.
func (s *Server) handleReadyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := s.store.Ping(ctx); err != nil {
		s.log.Error("readyz: database unreachable", "err", err)
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": "unavailable",
			"db":     "down",
			"reason": "Không kết nối được cơ sở dữ liệu.",
		})
		return
	}

	last, err := s.store.LastCrawlAt(ctx)
	if err != nil {
		s.log.Error("readyz: cannot read crawl_runs", "err", err)
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": "unavailable",
			"db":     "up",
			"reason": "Không đọc được lịch sử thu thập tin.",
		})
		return
	}

	body := map[string]any{"status": "ok", "db": "up", "is_stale": false}
	if last != nil {
		body["last_crawl_at"] = last.UTC().Format(time.RFC3339)
		if s.isStale(*last) {
			body["status"] = "stale"
			body["is_stale"] = true
			body["reason"] = "Dữ liệu tin tức đã quá cũ so với chu kỳ thu thập."
			writeJSON(w, http.StatusServiceUnavailable, body)
			return
		}
	}
	writeJSON(w, http.StatusOK, body)
}
