// Package extractor bóc tách nội dung sạch từ HTML bằng goquery
// (readability-lite).
package extractor

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// MinBodyLength — dưới ngưỡng này coi là paywall hoặc bài không đủ text để tóm
// tắt (F-07, F-11).
const MinBodyLength = 400

// Article là kết quả extract.
type Article struct {
	Title        string
	Body         string
	CanonicalURL string
	PublishedAt  time.Time
	HasDate      bool
	IsLive       bool
}

// Extractor bóc tách title, body, canonical URL và thời gian đăng.
type Extractor struct {
	// Location là múi giờ mặc định khi nguồn không ghi offset.
	Location *time.Location
	// Now cho phép test bơm thời gian cố định khi parse "2 giờ trước".
	Now func() time.Time
}

// New tạo extractor. loc=nil thì dùng UTC.
func New(loc *time.Location) *Extractor {
	if loc == nil {
		loc = time.UTC
	}
	return &Extractor{Location: loc, Now: time.Now}
}

// noiseSelectors là các khối không thuộc nội dung bài. Lọt vào thì tóm tắt lẫn
// nội dung bài khác và validator so sai số (E-05).
var noiseSelectors = []string{
	"script", "style", "noscript", "iframe", "svg", "form", "nav", "header",
	"footer", "aside", "figure", "figcaption", ".related", ".relate",
	".tin-lien-quan", ".box-tin-lien-quan", ".comment", ".comments",
	".advertisement", ".ads", ".ad", ".banner", ".social", ".share",
	".breadcrumb", ".tag", ".tags", ".author", ".source", ".copyright",
	"[class*=related]", "[id*=related]", "[class*=comment]",
	// Layout dạng thẻ tin: mỗi thẻ là một khối text độc lập, không thuộc thân
	// bài. Nếu để lọt, thẻ quảng bá dùng chung cho cả site sẽ trở thành "thân
	// bài" của mọi bài trên nguồn đó.
	"[class*=news-feed]", "[class*=new-item]", "[class*=general-item]",
	"[class*=ai-news]", "[class*=box-adv]", "[class*=advertisement]",
	"[data-ad]", "[class*=sidebar]", "#sidebar", "[class*=widget]",
	"[class*=box-tin]", "[class*=doc-them]", "[class*=xem-them]",
	"[class*=newsletter]", "[class*=subscribe]", "[class*=most-read]",
	"[class*=popular]", "[class*=trending]",
}

// bodySelectors theo thứ tự ưu tiên: selector riêng của các CMS báo VN phổ biến
// trước, heuristic chung sau.
var bodySelectors = []string{
	"[itemprop=articleBody]",
	"[data-field=body]",
	"article .fck_detail",
	".fck_detail",
	".detail-content",
	".detail__content",
	".article-editor",
	".article-content",
	".article__body",
	".contentdetail",
	"#mainContent",
	".singular-content",
	".entry-content",
	"main",
	"article",
}

// Extract phân tích HTML và trả về nội dung sạch.
func (e *Extractor) Extract(html, pageURL string) (Article, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return Article{}, fmt.Errorf("parse HTML of %s: %w", pageURL, err)
	}

	article := Article{
		Title:        extractTitle(doc),
		CanonicalURL: extractCanonical(doc),
		IsLive:       looksLikeLiveBlog(doc),
	}
	if t, ok := e.extractPublishedAt(doc); ok {
		article.PublishedAt = t
		article.HasDate = true
	}

	removeNoise(doc)
	article.Body = extractBody(doc)

	if article.Title == "" {
		return article, fmt.Errorf("no title found in %s", pageURL)
	}
	if length := len([]rune(article.Body)); length < MinBodyLength {
		return article, fmt.Errorf("body is only %d chars (< %d), likely paywalled or not an article: %s",
			length, MinBodyLength, pageURL)
	}
	return article, nil
}

// newspaperSuffixRe tách phần đuôi tên báo trong <title> ("... - VnExpress").
var titleSeparators = []string{" - ", " | ", " — ", " – ", " · "}

func extractTitle(doc *goquery.Document) string {
	for _, sel := range []string{
		`meta[property="og:title"]`,
		`meta[name="twitter:title"]`,
		`meta[itemprop="headline"]`,
	} {
		if v, ok := doc.Find(sel).First().Attr("content"); ok {
			if s := clean(v); s != "" {
				return s
			}
		}
	}
	if s := clean(doc.Find("h1").First().Text()); s != "" {
		return s
	}
	return stripNewspaperSuffix(clean(doc.Find("title").First().Text()))
}

// stripNewspaperSuffix cắt đuôi tên báo khỏi <title> (E-07). Chỉ cắt khi phần
// đuôi ngắn, để không xén nhầm tiêu đề có chứa dấu gạch.
func stripNewspaperSuffix(title string) string {
	for _, sep := range titleSeparators {
		if idx := strings.LastIndex(title, sep); idx > 0 {
			head, tail := title[:idx], title[idx+len(sep):]
			if len(tail) <= 30 && len(head) >= 20 {
				return strings.TrimSpace(head)
			}
		}
	}
	return title
}

// extractCanonical đọc <link rel="canonical"> — nguồn sự thật đáng tin nhất cho
// URL bài, đứng trên cả URL đã fetch (E-01).
func extractCanonical(doc *goquery.Document) string {
	if v, ok := doc.Find(`link[rel="canonical"]`).First().Attr("href"); ok {
		return clean(v)
	}
	if v, ok := doc.Find(`meta[property="og:url"]`).First().Attr("content"); ok {
		return clean(v)
	}
	return ""
}

func looksLikeLiveBlog(doc *goquery.Document) bool {
	if doc.Find(".live-blog, .liveblog, [class*=live-update]").Length() > 0 {
		return true
	}
	title := strings.ToLower(extractTitle(doc))
	return strings.Contains(title, "trực tiếp:") || strings.Contains(title, "tường thuật trực tiếp")
}

// dateSelectors là các chỗ báo VN thường đặt thời gian, xếp theo độ tin cậy.
// Meta ISO-8601 đáng tin hơn text hiển thị (T-05).
var dateMetaSelectors = []struct{ sel, attr string }{
	{`meta[property="article:published_time"]`, "content"},
	{`meta[itemprop="datePublished"]`, "content"},
	{`meta[name="pubdate"]`, "content"},
	{`meta[name="publish_date"]`, "content"},
	{`meta[name="article:published_time"]`, "content"},
	{`time[datetime]`, "datetime"},
}

var dateTextSelectors = []string{
	".date", ".time", ".publish-date", ".article-date", ".detail-time",
	"[class*=publish]", "[class*=pubtime]", "span.time", "time",
}

func (e *Extractor) extractPublishedAt(doc *goquery.Document) (time.Time, bool) {
	now := time.Now
	if e.Now != nil {
		now = e.Now
	}

	for _, c := range dateMetaSelectors {
		v, ok := doc.Find(c.sel).First().Attr(c.attr)
		if !ok {
			continue
		}
		if t, ok := ParsePublishedAt(clean(v), now(), e.Location); ok {
			return t, true
		}
	}
	for _, sel := range dateTextSelectors {
		text := clean(doc.Find(sel).First().Text())
		if text == "" {
			continue
		}
		if t, ok := ParsePublishedAt(text, now(), e.Location); ok {
			return t, true
		}
	}
	return time.Time{}, false
}

// removeNoise xoá các khối không thuộc bài, nhưng KHÔNG BAO GIỜ xoá một node
// đang chứa thân bài.
//
// Selector nhiễu viết theo wildcard rất dễ trúng nhầm: VnExpress đặt cột nội
// dung chính trong `class="sidebar-1"`, nên "[class*=sidebar]" từng xoá sạch
// toàn bộ 28 bài của nguồn này. Hàng rào cấu trúc này làm cả danh sách nhiễu an
// toàn theo thiết kế, thay vì phải soi từng selector một.
func removeNoise(doc *goquery.Document) {
	for _, sel := range noiseSelectors {
		doc.Find(sel).Each(func(_ int, node *goquery.Selection) {
			if containsArticleBody(node) {
				return
			}
			node.Remove()
		})
	}
}

func containsArticleBody(node *goquery.Selection) bool {
	for _, sel := range bodySelectors {
		if sel == "main" || sel == "article" {
			continue
		}
		if node.Find(sel).Length() > 0 {
			return true
		}
	}
	return false
}

func extractBody(doc *goquery.Document) string {
	best := ""
	for _, sel := range bodySelectors {
		doc.Find(sel).Each(func(_ int, s *goquery.Selection) {
			if txt := paragraphs(s); len(txt) > len(best) {
				best = txt
			}
		})
		if len([]rune(best)) >= MinBodyLength {
			return best
		}
	}
	if best == "" {
		best = paragraphs(doc.Find("body"))
	}
	return best
}

// paragraphs gom sapo, các thẻ <p> và ô bảng thành text sạch. Bảng số liệu
// thường là chỗ chứa những con số quan trọng nhất của bài (E-06).
//
// Trả chuỗi rỗng khi khối không có đoạn văn nào. Bản cũ rơi về s.Text() ở đây,
// nên một thẻ tin quảng bá không có <p> nào vẫn đủ tư cách làm "thân bài" —
// đúng đường đã khiến 36 bài VnEconomy dùng chung một tóm tắt về hội thảo
// tầng ô-dôn.
func paragraphs(s *goquery.Selection) string {
	var parts []string
	s.Find("[data-field=sapo], .sapo, .article-content__lead").Each(func(_ int, lead *goquery.Selection) {
		if t := clean(lead.Text()); len([]rune(t)) >= 20 {
			parts = append(parts, t)
		}
	})
	s.Find("p").Each(func(_ int, p *goquery.Selection) {
		if t := clean(p.Text()); len([]rune(t)) >= 20 {
			parts = append(parts, t)
		}
	})
	s.Find("table tr").Each(func(_ int, tr *goquery.Selection) {
		var cells []string
		tr.Find("th, td").Each(func(_ int, td *goquery.Selection) {
			if t := clean(td.Text()); t != "" {
				cells = append(cells, t)
			}
		})
		if len(cells) >= 2 {
			parts = append(parts, strings.Join(cells, " | "))
		}
	})
	if len(parts) == 0 {
		return ""
	}
	return dedupeParts(parts)
}

// dedupeParts bỏ đoạn lặp (sapo thường được lặp lại làm đoạn mở đầu).
func dedupeParts(parts []string) string {
	seen := make(map[string]bool, len(parts))
	kept := parts[:0]
	for _, part := range parts {
		if seen[part] {
			continue
		}
		seen[part] = true
		kept = append(kept, part)
	}
	return strings.Join(kept, "\n\n")
}

func clean(s string) string {
	s = strings.ReplaceAll(s, " ", " ")
	return strings.Join(strings.Fields(s), " ")
}
