package http

import (
	"errors"
	"fmt"
	"html"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/tenpoint/tenpoint-api/internal/shortlink"
	"github.com/tenpoint/tenpoint-api/internal/store"
)

// handleRedirect phục vụ GET /r/{code}.
//
// Yêu cầu hiệu năng < 50ms (BA1 US-3.2 AC4): tra cache LRU trước nên phần lớn
// request không chạm DB, và click được đếm qua buffer trong RAM.
// Dùng 302 chứ không 301 vì 301 bị trình duyệt cache vĩnh viễn, mất tracking từ
// lần click thứ hai và không sửa được đích khi nguồn đổi URL (R7.7).
func (s *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	if !s.redirectR.allow(clientIP(r)) {
		writeError(w, http.StatusTooManyRequests, CodeInvalidRequest,
			"Bạn đang truy cập quá nhanh, vui lòng thử lại sau ít phút.")
		return
	}
	if !shortlink.IsValidCode(code) {
		writeError(w, http.StatusNotFound, CodeNotFound, "Liên kết không tồn tại.")
		return
	}

	target, cached := s.cache.Get(code)
	if !cached {
		link, err := s.store.GetShortLink(r.Context(), code)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				// Không phân biệt "không tồn tại" với "bị vô hiệu" trước scanner.
				writeError(w, http.StatusNotFound, CodeNotFound, "Liên kết không tồn tại.")
				return
			}
			s.log.Error("get short link", "code", code, "err", err)
			writeError(w, http.StatusInternalServerError, CodeInternal,
				"Có lỗi xảy ra, vui lòng thử lại sau.")
			return
		}
		target = shortlink.Target{ArticleID: link.ArticleID, URL: link.TargetURL, Alive: link.TargetAlive}
		s.cache.Put(code, target)
	}

	if !target.Alive {
		s.writeGone(w, r, code)
		return
	}

	if s.clicks != nil {
		s.clicks.Add(code)
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	http.Redirect(w, r, target.URL, http.StatusFound)
}

// writeGone trả 410 kèm trang tiếng Việt và lối đi tiếp sang trang mã liên quan.
// Redirect về trang chủ sẽ khiến người dùng tưởng link hỏng do TenPoint và mất
// hoàn toàn ngữ cảnh họ đang đọc (L-06).
func (s *Server) writeGone(w http.ResponseWriter, r *http.Request, code string) {
	if prefersJSON(r) {
		writeError(w, http.StatusGone, CodeGone, "Bài gốc không còn truy cập được.")
		return
	}

	symbols, err := s.store.TickerSymbolsForShortCode(r.Context(), code)
	if err != nil {
		s.log.Warn("load tickers for gone page", "code", code, "err", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusGone)
	if _, err := w.Write([]byte(gonePage(symbols))); err != nil {
		s.log.Error("write gone page", "code", code, "err", err)
	}
}

func gonePage(symbols []string) string {
	var links strings.Builder
	for _, symbol := range symbols {
		safe := html.EscapeString(symbol)
		fmt.Fprintf(&links, `<li><a href="/ma/%s">Tin về mã %s</a></li>`, safe, safe)
	}
	if links.Len() == 0 {
		links.WriteString(`<li><a href="/">Về trang tin tổng hợp</a></li>`)
	}

	return `<!doctype html>
<html lang="vi">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="robots" content="noindex">
<title>Bài gốc không còn truy cập được — TenPoint</title>
<style>
body{font-family:Georgia,"Times New Roman",serif;background:#F7F5F0;color:#111;
margin:0;padding:48px 20px;line-height:1.6}
main{max-width:640px;margin:0 auto}
h1{font-size:1.5rem;line-height:1.3;margin:0 0 16px}
p{margin:0 0 16px}
ul{padding-left:20px}
a{color:#111}
</style>
</head>
<body>
<main>
<h1>Bài gốc không còn truy cập được</h1>
<p>Nguồn tin đã gỡ hoặc thay đổi bài viết này, nên liên kết không còn dẫn tới
nội dung ban đầu. Bản tóm tắt của TenPoint vẫn được giữ lại trong kho tin.</p>
<p>Bạn có thể xem tiếp:</p>
<ul>` + links.String() + `</ul>
</main>
</body>
</html>`
}

func prefersJSON(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	if accept == "" {
		return false
	}
	if strings.Contains(accept, "text/html") {
		return false
	}
	return strings.Contains(accept, "application/json")
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i > 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
