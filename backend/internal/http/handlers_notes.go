package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/tenpoint/tenpoint-api/internal/store"
)

// handleListNotes phục vụ GET /api/v1/notes.
func (s *Server) handleListNotes(w http.ResponseWriter, r *http.Request) {
	loc := s.location()
	notes, err := s.store.ListNotes(r.Context(), 50)
	if err != nil {
		s.log.Error("list research notes", "err", err)
		writeError(w, http.StatusInternalServerError, CodeInternal,
			"Không tải được danh sách nhận định, vui lòng thử lại sau.")
		return
	}

	items := make([]NoteItem, 0, len(notes))
	for _, n := range notes {
		items = append(items, toNoteItem(n, loc))
	}
	writeJSON(w, http.StatusOK, listEnvelope{
		Data: items,
		Meta: &ListMeta{Total: len(items), Limit: len(items), Offset: 0, HasMore: false},
	})
}

// handleGetNote phục vụ GET /api/v1/notes/{slug}.
func (s *Server) handleGetNote(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimSpace(chi.URLParam(r, "slug"))
	if slug == "" || len(slug) > 200 {
		writeError(w, http.StatusBadRequest, CodeInvalidRequest, "Đường dẫn nhận định không hợp lệ.")
		return
	}

	loc := s.location()
	note, err := s.store.GetNote(r.Context(), slug)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, CodeNotFound, "Không tìm thấy nhận định này.")
			return
		}
		s.log.Error("get research note", "slug", slug, "err", err)
		writeError(w, http.StatusInternalServerError, CodeInternal,
			"Không tải được nhận định, vui lòng thử lại sau.")
		return
	}

	points := make([]NotePoint, 0, len(note.Points))
	for _, p := range note.Points {
		points = append(points, NotePoint{Ordinal: p.Ordinal, Lead: p.Lead, Body: p.Body})
	}

	writeJSON(w, http.StatusOK, itemEnvelope{Data: NoteDetail{
		NoteItem:   toNoteItem(note, loc),
		Disclaimer: note.Disclaimer,
		Points:     points,
	}})
}
