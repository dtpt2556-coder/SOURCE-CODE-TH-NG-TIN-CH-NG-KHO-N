package http

import (
	"context"
	"crypto/subtle"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/tenpoint/tenpoint-api/internal/domain"
	"github.com/tenpoint/tenpoint-api/internal/ingest"
)

// manualIngestTimeout giới hạn thời gian một run thủ công.
const manualIngestTimeout = 15 * time.Minute

// handleAdminIngest phục vụ POST /api/v1/admin/ingest (Bearer token).
// Trả 202 ngay và chạy pipeline nền; 409 nếu đang có run khác.
func (s *Server) handleAdminIngest(w http.ResponseWriter, r *http.Request) {
	if !s.authorizeAdmin(r) {
		writeError(w, http.StatusUnauthorized, CodeUnauthorized,
			"Token quản trị không hợp lệ.")
		return
	}
	if s.runner == nil {
		writeError(w, http.StatusInternalServerError, CodeInternal,
			"Pipeline chưa được cấu hình trên tiến trình này.")
		return
	}

	started := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), manualIngestTimeout)
		defer cancel()
		_, runErr := s.runner.Run(ctx, domain.TriggerManual)
		started <- runErr
		if runErr != nil && !errors.Is(runErr, ingest.ErrAlreadyRunning) {
			s.log.Error("manual run failed", "err", runErr)
		}
	}()

	// Chờ ngắn để phân biệt 409 (đang chạy) với 202 (đã nhận).
	select {
	case runErr := <-started:
		if errors.Is(runErr, ingest.ErrAlreadyRunning) {
			writeError(w, http.StatusConflict, CodeConflict,
				"Đang có tiến trình thu thập tin chạy, vui lòng thử lại sau.")
			return
		}
	case <-time.After(300 * time.Millisecond):
	}

	writeJSON(w, http.StatusAccepted, itemEnvelope{Data: map[string]string{
		"status":  "accepted",
		"trigger": domain.TriggerManual,
		"message": "Đã nhận yêu cầu thu thập tin, pipeline đang chạy nền.",
	}})
}

// authorizeAdmin so sánh Bearer token bằng thời gian hằng số.
func (s *Server) authorizeAdmin(r *http.Request) bool {
	if s.cfg == nil || s.cfg.AdminToken == "" {
		// Không cấu hình token => endpoint bị khoá hoàn toàn (mặc định an toàn).
		return false
	}
	header := r.Header.Get("Authorization")
	token, ok := strings.CutPrefix(header, "Bearer ")
	if !ok {
		return false
	}
	token = strings.TrimSpace(token)
	return subtle.ConstantTimeCompare([]byte(token), []byte(s.cfg.AdminToken)) == 1
}
