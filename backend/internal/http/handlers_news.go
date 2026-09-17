package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tenpoint/tenpoint-api/internal/domain"
	"github.com/tenpoint/tenpoint-api/internal/store"
)

// Giới hạn phân trang theo mục 6 architecture.md.
const (
	defaultLimit = 20
	maxLimit     = 100
	defaultDays  = 7
)

// handleListNews phục vụ GET /api/v1/news.
func (s *Server) handleListNews(w http.ResponseWriter, r *http.Request) {
	loc := s.location()
	q := r.URL.Query()

	filter := domain.ArticleFilter{
		Relevance: "primary",
		Limit:     defaultLimit,
	}

	if v := strings.TrimSpace(q.Get("relevance")); v != "" {
		if v != "primary" && v != "all" {
			writeError(w, http.StatusBadRequest, CodeInvalidRequest,
				"Tham số relevance chỉ nhận giá trị primary hoặc all.")
			return
		}
		filter.Relevance = v
	}

	if v := strings.TrimSpace(q.Get("ma")); v != "" {
		for _, sym := range strings.Split(v, ",") {
			sym = strings.ToUpper(strings.TrimSpace(sym))
			if sym != "" {
				filter.Symbols = append(filter.Symbols, sym)
			}
		}
	}

	if v := strings.TrimSpace(q.Get("loai")); v != "" {
		for _, raw := range strings.Split(v, ",") {
			raw = strings.ToLower(strings.TrimSpace(raw))
			if raw == "" {
				continue
			}
			nt, ok := domain.ParseNewsType(raw)
			if !ok {
				writeError(w, http.StatusBadRequest, CodeInvalidRequest,
					"Giá trị loai không hợp lệ: "+raw+". Các giá trị hợp lệ: "+domain.NewsTypeValues()+".")
				return
			}
			filter.NewsTypes = append(filter.NewsTypes, nt)
		}
	}

	from, to, err := parseDateRange(q.Get("tu"), q.Get("den"), loc)
	if err != nil {
		writeError(w, http.StatusBadRequest, CodeInvalidRequest, err.Error())
		return
	}
	filter.From, filter.To = from, to
	filter.Query = strings.TrimSpace(q.Get("q"))

	if v := strings.TrimSpace(q.Get("limit")); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			writeError(w, http.StatusBadRequest, CodeInvalidRequest,
				"Tham số limit phải là số nguyên dương.")
			return
		}
		if n > maxLimit {
			n = maxLimit
		}
		filter.Limit = n
	}
	if v := strings.TrimSpace(q.Get("offset")); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			writeError(w, http.StatusBadRequest, CodeInvalidRequest,
				"Tham số offset phải là số nguyên không âm.")
			return
		}
		filter.Offset = n
	}

	// Chỉ hiện dữ liệu mẫu khi chưa crawl được bài thật nào, để bảng digest ở
	// lần chạy đầu không trắng trơn. Có tin thật là demo biến mất (T3 mục 0).
	hasReal, err := s.store.HasRealArticles(r.Context())
	if err != nil {
		s.log.Error("check for real articles", "err", err)
		writeError(w, http.StatusInternalServerError, CodeInternal,
			"Không tải được danh sách tin, vui lòng thử lại sau.")
		return
	}
	filter.IncludeDemo = !hasReal

	articles, total, err := s.store.ListArticles(r.Context(), filter)
	if err != nil {
		s.log.Error("list articles", "err", err)
		writeError(w, http.StatusInternalServerError, CodeInternal,
			"Không tải được danh sách tin, vui lòng thử lại sau.")
		return
	}

	writeJSON(w, http.StatusOK, listEnvelope{
		Data: toNewsItems(articles, loc),
		Meta: &ListMeta{
			Total:   total,
			Limit:   filter.Limit,
			Offset:  filter.Offset,
			HasMore: filter.Offset+len(articles) < total,
		},
	})
}

// handleGetNews phục vụ GET /api/v1/news/{id}.
func (s *Server) handleGetNews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, CodeInvalidRequest, "Mã tin không hợp lệ.")
		return
	}

	loc := s.location()
	article, err := s.store.GetArticle(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, CodeNotFound, "Không tìm thấy tin này.")
			return
		}
		s.log.Error("get article", "id", id, "err", err)
		writeError(w, http.StatusInternalServerError, CodeInternal,
			"Không tải được tin, vui lòng thử lại sau.")
		return
	}

	siblings, err := s.store.ListClusterSiblings(r.Context(), article.ID, article.ClusterID)
	if err != nil {
		s.log.Error("list cluster siblings", "id", id, "err", err)
		siblings = nil
	}

	writeJSON(w, http.StatusOK, itemEnvelope{Data: NewsDetail{
		NewsItem:       toNewsItem(article, loc),
		AlsoReportedBy: toNewsItems(siblings, loc),
	}})
}

// parseDateRange đọc `tu`/`den` theo giờ VN. Mặc định 7 ngày gần nhất.
func parseDateRange(tu, den string, loc *time.Location) (*time.Time, *time.Time, error) {
	const layout = "2006-01-02"
	var from, to *time.Time

	tu = strings.TrimSpace(tu)
	den = strings.TrimSpace(den)

	if tu != "" {
		t, err := time.ParseInLocation(layout, tu, loc)
		if err != nil {
			return nil, nil, errInvalidDate("tu")
		}
		u := t.UTC()
		from = &u
	}
	if den != "" {
		t, err := time.ParseInLocation(layout, den, loc)
		if err != nil {
			return nil, nil, errInvalidDate("den")
		}
		// `den` là ngày bao gồm cả ngày đó -> chặn trên là 00:00 hôm sau.
		u := t.AddDate(0, 0, 1).UTC()
		to = &u
	}
	if from == nil && to == nil {
		d := time.Now().In(loc).AddDate(0, 0, -defaultDays)
		start := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, loc).UTC()
		from = &start
	}
	if from != nil && to != nil && !from.Before(*to) {
		return nil, nil, errDateRange()
	}
	return from, to, nil
}

type apiError struct{ msg string }

func (e apiError) Error() string { return e.msg }

func errInvalidDate(field string) error {
	return apiError{msg: "Tham số " + field + " phải có định dạng YYYY-MM-DD."}
}

func errDateRange() error {
	return apiError{msg: "Khoảng thời gian không hợp lệ: tu phải trước den."}
}
