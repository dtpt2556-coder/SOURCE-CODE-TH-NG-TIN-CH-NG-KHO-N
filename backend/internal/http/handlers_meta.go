package http

import (
	"net/http"
	"time"

	"github.com/tenpoint/tenpoint-api/internal/domain"
)

// handleMeta phục vụ GET /api/v1/meta: độ tươi dữ liệu + thống kê loại tin.
func (s *Server) handleMeta(w http.ResponseWriter, r *http.Request) {
	stats, err := s.store.Meta(r.Context())
	if err != nil {
		s.log.Error("load meta stats", "err", err)
		writeError(w, http.StatusInternalServerError, CodeInternal,
			"Không tải được thông tin hệ thống, vui lòng thử lại sau.")
		return
	}

	hasReal, err := s.store.HasRealArticles(r.Context())
	if err != nil {
		s.log.Error("check for real articles", "err", err)
		writeError(w, http.StatusInternalServerError, CodeInternal,
			"Không tải được thông tin hệ thống, vui lòng thử lại sau.")
		return
	}

	resp := MetaResponse{
		TotalArticles:   stats.TotalArticles,
		TickersCount:    stats.TickersCount,
		HasRealArticles: hasReal,
		NewsTypes:       make([]NewsTypeCount, 0, len(domain.AllNewsTypes())),
	}
	if stats.LastCrawlAt != nil {
		formatted := stats.LastCrawlAt.UTC().Format(time.RFC3339)
		resp.LastCrawlAt = &formatted
		resp.IsStale = s.isStale(*stats.LastCrawlAt)
	}
	for _, nt := range domain.AllNewsTypes() {
		resp.NewsTypes = append(resp.NewsTypes, NewsTypeCount{
			Value: string(nt),
			Label: nt.Label(),
			Count: stats.NewsTypes[nt],
		})
	}
	writeJSON(w, http.StatusOK, itemEnvelope{Data: resp})
}

// isStale: dữ liệu quá cũ so với ngưỡng cấu hình (mặc định 26 giờ ~ 3 chu kỳ cron).
func (s *Server) isStale(last time.Time) bool {
	maxAge := 26 * time.Hour
	if s.cfg != nil && s.cfg.ReadyMaxAge > 0 {
		maxAge = s.cfg.ReadyMaxAge
	}
	return time.Since(last) > maxAge
}
