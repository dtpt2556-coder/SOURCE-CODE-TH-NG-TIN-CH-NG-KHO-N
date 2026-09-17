package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/tenpoint/tenpoint-api/internal/store"
)

// handleListTickers phục vụ GET /api/v1/tickers.
func (s *Server) handleListTickers(w http.ResponseWriter, r *http.Request) {
	tickers, err := s.store.ListTickers(r.Context())
	if err != nil {
		s.log.Error("list tickers", "err", err)
		writeError(w, http.StatusInternalServerError, CodeInternal,
			"Không tải được danh sách mã, vui lòng thử lại sau.")
		return
	}

	items := make([]TickerItem, 0, len(tickers))
	for _, t := range tickers {
		items = append(items, TickerItem{
			Symbol:       t.Symbol,
			CompanyName:  t.CompanyName,
			ShortName:    t.ShortName,
			Exchange:     t.Exchange,
			Sector:       t.Sector,
			InVN30:       t.InVN30,
			ArticleCount: t.ArticleCount,
		})
	}
	writeJSON(w, http.StatusOK, listEnvelope{
		Data: items,
		Meta: &ListMeta{Total: len(items), Limit: len(items), Offset: 0, HasMore: false},
	})
}

// handleGetTicker phục vụ GET /api/v1/tickers/{symbol}.
func (s *Server) handleGetTicker(w http.ResponseWriter, r *http.Request) {
	symbol := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "symbol")))
	if symbol == "" || len(symbol) > 10 {
		writeError(w, http.StatusBadRequest, CodeInvalidRequest, "Mã chứng khoán không hợp lệ.")
		return
	}

	ctx := r.Context()
	loc := s.location()

	ticker, err := s.store.GetTicker(ctx, symbol)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, CodeNotFound, "Không tìm thấy mã chứng khoán này.")
			return
		}
		s.log.Error("get ticker", "symbol", symbol, "err", err)
		writeError(w, http.StatusInternalServerError, CodeInternal,
			"Không tải được thông tin mã, vui lòng thử lại sau.")
		return
	}

	news, err := s.store.ListArticlesBySymbol(ctx, symbol, 20)
	if err != nil {
		s.log.Error("list articles by ticker", "symbol", symbol, "err", err)
		writeError(w, http.StatusInternalServerError, CodeInternal,
			"Không tải được tin của mã, vui lòng thử lại sau.")
		return
	}

	notes, err := s.store.ListNotesBySymbol(ctx, symbol)
	if err != nil {
		s.log.Error("list notes by ticker", "symbol", symbol, "err", err)
		notes = nil
	}
	noteItems := make([]NoteItem, 0, len(notes))
	for _, n := range notes {
		noteItems = append(noteItems, toNoteItem(n, loc))
	}

	writeJSON(w, http.StatusOK, itemEnvelope{Data: TickerDetail{
		TickerItem: TickerItem{
			Symbol:      ticker.Symbol,
			CompanyName: ticker.CompanyName,
			ShortName:   ticker.ShortName,
			Exchange:    ticker.Exchange,
			Sector:      ticker.Sector,
			InVN30:      ticker.InVN30,
		},
		RecentNews:    toNewsItems(news, loc),
		ResearchNotes: noteItems,
	}})
}
