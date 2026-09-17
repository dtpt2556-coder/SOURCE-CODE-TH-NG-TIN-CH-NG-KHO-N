// Package tagger gắn mã chứng khoán vào bài viết bằng từ điển + luật
// (ADR-005: tất định, kiểm chứng được, chạy offline — KHÔNG dùng LLM).
//
// Thuật toán 3 tầng theo mục 8 architecture.md và §4 BA2:
//
//	T1 — Mã trực tiếp   : regex \b[A-Z]{3}\b đối chiếu bảng tickers          (0.90)
//	T2 — Tên doanh nghiệp: khớp ticker_aliases.alias_normalized               (0.80)
//	T3 — Ngữ cảnh ngành  : từ khoá ngành -> nhóm mã, CHỈ gán `mentioned`      (0.40)
package tagger

import (
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/tenpoint/tenpoint-api/internal/domain"
	"github.com/tenpoint/tenpoint-api/internal/textutil"
)

// MaxTickersPerArticle giới hạn số mã hiển thị (rule B8 của BA2).
const MaxTickersPerArticle = 8

// MaxSectorTickers giới hạn số mã mà tầng ngành được phép tự thêm.
const MaxSectorTickers = 4

// SectorMentionsRequired là số lần từ khoá ngành phải xuất hiện trong thân bài
// để coi ngành đó là chủ đề, khi tiêu đề không nêu.
const SectorMentionsRequired = 2

// TickerListNeighbours là số mã khác phải cùng xuất hiện trong cửa sổ để coi đó
// là một danh sách mã chứng khoán.
const TickerListNeighbours = 2

// Điểm khởi điểm của từng tầng.
const (
	scoreDirect      = 0.90
	scoreAlias       = 0.80
	scoreConfirmed   = 0.95
	scoreSectorOnly  = 0.40
	contextWindowLen = 80
)

// codeRe bắt mã 3 chữ cái in hoa. \b của RE2 là ASCII nên các từ tiếng Việt có
// dấu không tạo ra false positive.
var codeRe = regexp.MustCompile(`\b[A-Z]{3}\b`)

// Ticker là dữ liệu tối thiểu tagger cần về một mã.
type Ticker struct {
	Symbol  string
	Sector  string
	Aliases []string
	// NegativeAliases là các chuỗi KHÔNG được map vào mã này: "Vinhomes" không
	// phải VIC, "Masan MEATLife" không phải MSN (rule G-04).
	NegativeAliases []string
	// InVN30 dùng để chọn mã tiêu biểu khi tầng ngành phải cắt bớt.
	InVN30 bool
}

// Match là kết quả gắn mã cho một bài.
type Match struct {
	Symbol    string
	Relevance string // domain.RelevancePrimary | domain.RelevanceMentioned
	Score     float64
}

// Input là văn bản đã extract (KHÔNG phải raw HTML — tránh bắt mã trong URL).
type Input struct {
	Title string
	Body  string
}

// Tagger là bộ gắn mã, an toàn khi dùng đồng thời (chỉ đọc sau khi New).
type Tagger struct {
	symbols      map[string]Ticker // symbol -> ticker
	aliasToSym   map[string]string // alias đã normalize -> symbol
	sectorToSyms map[string][]string
	sectorKw     map[string]string // từ khoá ngành đã normalize -> tên ngành
	stoplist     map[string]bool
	negatives    map[string][]string // symbol -> cụm phủ định đã normalize
	vn30         map[string]bool
	industryTier bool
}

// New dựng tagger từ danh sách mã, với tầng ngành TẮT.
func New(tickers []Ticker) *Tagger {
	return NewWithOptions(tickers, DefaultStoplist(), DefaultSectorKeywords(), false)
}

// NewWithIndustryTier bật tầng T3 ngữ cảnh ngành.
//
// Mặc định tắt (env TAGGER_INDUSTRY_TIER). QA đếm TOÀN BỘ tầng này trên dữ liệu
// thật: 38 tag trên 12 bài, 38/38 đều sai — "thương hiệu bò H'Mông" gắn
// BAF/DBC/HAG/MML, "Canada và thế khó của nền kinh tế hướng Mỹ" gắn
// HPG/HSG/NKG/TIS. Riêng nó tạo trần cứng 86,7% precision cho toàn quần thể.
//
// Cách nhận diện hiện tại chỉ dò từ khoá ngành nên mọi bài chứa chữ "thép" hay
// "bò" đều dính. Làm đúng thì bài phải thực sự nói về một sự kiện kinh tế tác
// động lên nhà sản xuất Việt Nam (thuế, giá, chính sách, xuất nhập khẩu) — đó
// là bài toán riêng của v1.1.
//
// ĐIỀU KIỆN BẬT LẠI: đạt precision đo được trên corpus do QA tạo
// (backend/testdata/tagger_corpus_qa.json), không phải trên corpus tự tạo.
func NewWithIndustryTier(tickers []Ticker) *Tagger {
	return NewWithOptions(tickers, DefaultStoplist(), DefaultSectorKeywords(), true)
}

// NewWithOptions cho phép nạp stoplist / từ khoá ngành riêng (dùng trong test).
func NewWithOptions(tickers []Ticker, stoplist []string, sectorKeywords map[string]string, industryTier bool) *Tagger {
	t := &Tagger{
		symbols:      make(map[string]Ticker, len(tickers)),
		aliasToSym:   make(map[string]string),
		sectorToSyms: make(map[string][]string),
		sectorKw:     make(map[string]string, len(sectorKeywords)),
		stoplist:     make(map[string]bool, len(stoplist)),
		negatives:    make(map[string][]string),
		vn30:         make(map[string]bool),
		industryTier: industryTier,
	}
	for _, s := range stoplist {
		t.stoplist[strings.ToUpper(strings.TrimSpace(s))] = true
	}
	for kw, sector := range sectorKeywords {
		t.sectorKw[normPhrase(kw)] = normPhrase(sector)
	}
	for _, tk := range tickers {
		sym := strings.ToUpper(strings.TrimSpace(tk.Symbol))
		if sym == "" {
			continue
		}
		tk.Symbol = sym
		t.symbols[sym] = tk
		t.vn30[sym] = tk.InVN30
		for _, neg := range tk.NegativeAliases {
			if key := normPhrase(neg); len([]rune(key)) >= 3 {
				t.negatives[sym] = append(t.negatives[sym], key)
			}
		}
		for _, a := range tk.Aliases {
			key := normPhrase(a)
			// Alias quá ngắn (< 3 ký tự) gây nhiễu nặng -> bỏ.
			if len([]rune(key)) < 3 {
				continue
			}
			// Alias trùng chính mã sẽ do T1 xử lý (có đủ luật khử nhập nhằng).
			if strings.EqualFold(key, sym) {
				continue
			}
			if _, exists := t.aliasToSym[key]; !exists {
				t.aliasToSym[key] = sym
			}
		}
		if sec := normPhrase(tk.Sector); sec != "" {
			t.sectorToSyms[sec] = append(t.sectorToSyms[sec], sym)
		}
	}
	for sec := range t.sectorToSyms {
		sort.Strings(t.sectorToSyms[sec])
	}
	return t
}

// hasPrimaryEvidence cho biết T1/T2 đã tìm được mã nào đủ mạnh để làm chủ thể
// của bài chưa.
func hasPrimaryEvidence(ev map[string]*evidence) bool {
	for _, e := range ev {
		if e.sectorOnly {
			continue
		}
		if e.titleHit || e.firstParaHit || e.bodyCount >= 2 {
			return true
		}
	}
	return false
}

// evidence gom bằng chứng của một mã trong bài.
type evidence struct {
	titleHit     bool
	firstParaHit bool
	bodyCount    int

	codeOccur      int
	codeAttributed int
	codeContext    int

	aliasHit   bool
	sectorOnly bool
}

// Tag chạy T1 -> T2 -> T3 và hợp nhất kết quả.
func (t *Tagger) Tag(in Input) []Match {
	ev := map[string]*evidence{}
	get := func(sym string) *evidence {
		e, ok := ev[sym]
		if !ok {
			e = &evidence{}
			ev[sym] = e
		}
		return e
	}

	firstPara := textutil.FirstParagraph(in.Body)

	// ---- T2 trước: tên doanh nghiệp xác nhận mã (rule B5 của BA2) ----
	for sym, n := range t.scanAliases(in.Title) {
		e := get(sym)
		e.aliasHit = true
		e.titleHit = true
		e.bodyCount += n
	}
	for sym := range t.scanAliases(firstPara) {
		e := get(sym)
		e.aliasHit = true
		e.firstParaHit = true
	}
	for sym, n := range t.scanAliases(in.Body) {
		e := get(sym)
		e.aliasHit = true
		e.bodyCount += n
	}

	// ---- T1: mã trực tiếp ----
	for sym, st := range t.scanCodes(in.Title) {
		e := get(sym)
		e.codeOccur += st.total
		e.codeAttributed += st.attributed
		e.codeContext += st.withContext
		if st.total > st.attributed {
			e.titleHit = true
			e.bodyCount += st.total - st.attributed
		}
	}
	for sym, st := range t.scanCodes(firstPara) {
		e := get(sym)
		if st.total > st.attributed {
			e.firstParaHit = true
		}
	}
	for sym, st := range t.scanCodes(in.Body) {
		e := get(sym)
		e.codeOccur += st.total
		e.codeAttributed += st.attributed
		e.codeContext += st.withContext
		e.bodyCount += st.total - st.attributed
	}

	// ---- T3: ngữ cảnh ngành (chỉ `mentioned`) ----
	// Chỉ chạy khi T1/T2 không tìm được mã primary nào. Bài đã nêu doanh nghiệp
	// cụ thể thì suy luận theo ngành chỉ thêm nhiễu (G-08).
	if t.industryTier && !hasPrimaryEvidence(ev) {
		for _, sym := range t.scanSectors(in.Title, in.Body) {
			if _, ok := ev[sym]; !ok {
				ev[sym] = &evidence{sectorOnly: true}
			}
		}
	}

	out := make([]Match, 0, len(ev))
	for sym, e := range ev {
		m, ok := t.decide(sym, e)
		if !ok {
			continue
		}
		out = append(out, m)
	}

	sort.Slice(out, func(i, j int) bool {
		if (out[i].Relevance == domain.RelevancePrimary) != (out[j].Relevance == domain.RelevancePrimary) {
			return out[i].Relevance == domain.RelevancePrimary
		}
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		return out[i].Symbol < out[j].Symbol
	})
	if len(out) > MaxTickersPerArticle {
		out = out[:MaxTickersPerArticle]
	}
	return out
}

// decide áp các rule khử nhập nhằng còn lại và chấm relevance.
func (t *Tagger) decide(sym string, e *evidence) (Match, bool) {
	if e.sectorOnly {
		return Match{Symbol: sym, Relevance: domain.RelevanceMentioned, Score: scoreSectorOnly}, true
	}

	if !e.aliasHit {
		// Chỉ do T1 phát hiện.
		valid := e.codeOccur - e.codeAttributed
		if valid <= 0 {
			// B3: mã chỉ nằm trong cụm trích dẫn nguồn -> bỏ.
			return Match{}, false
		}
		// B4: bắt buộc có ngữ cảnh chứng khoán quanh mã.
		if e.codeContext <= 0 {
			return Match{}, false
		}
	}

	score := scoreAlias
	switch {
	case e.aliasHit && e.codeOccur > e.codeAttributed:
		score = scoreConfirmed
	case !e.aliasHit:
		score = scoreDirect
	}

	rel := domain.RelevanceMentioned
	if e.titleHit || e.firstParaHit || e.bodyCount >= 2 {
		rel = domain.RelevancePrimary
	}
	return Match{Symbol: sym, Relevance: rel, Score: score}, true
}

// ---------------------------------------------------------------- T1 scanning

type codeStat struct {
	total       int
	attributed  int
	withContext int
}

func (t *Tagger) scanCodes(text string) map[string]*codeStat {
	res := map[string]*codeStat{}
	if text == "" {
		return res
	}
	normFull := normPhrase(text)
	hasVNDirect := strings.Contains(normFull, "vndirect")

	for _, loc := range codeRe.FindAllStringIndex(text, -1) {
		code := text[loc[0]:loc[1]]
		if t.stoplist[code] {
			continue // B1: stoplist tuyệt đối
		}
		if _, ok := t.symbols[code]; !ok {
			continue
		}
		windowStart, windowEnd := max(0, loc[0]-contextWindowLen), min(len(text), loc[1]+contextWindowLen)
		before := normPhrase(text[windowStart:loc[0]])
		after := normPhrase(text[loc[1]:windowEnd])
		rawAfter := text[loc[1]:min(len(text), loc[1]+8)]
		neighbours := t.knownCodesAround(text[windowStart:windowEnd], code)

		if code == "VND" && !isVNDirectOccurrence(before, after, hasVNDirect) {
			continue // B2: "5.000 tỷ VND" là đơn vị tiền, không phải mã
		}
		if t.insideNegativeAlias(code, text, loc[0]) {
			// "FPT Retail (FRT)" không được gắn FPT: negative_aliases trước đây
			// chỉ áp cho tầng alias, nên mã trần vẫn lọt ở điểm 0.90.
			continue
		}
		if isFundTicker(before, after) {
			// "VanEck Vectors Vietnam ETF (VNM ETF)" là mã quỹ nước ngoài trùng
			// ký hiệu với mã sàn HOSE. Đây là lỗi primary QA tìm thấy ở tin 26.
			continue
		}

		st, ok := res[code]
		if !ok {
			st = &codeStat{}
			res[code] = st
		}
		st.total++
		if isAttribution(before, after) {
			st.attributed++
			continue
		}
		// Mã nằm giữa một danh sách mã khác là ngữ cảnh chứng khoán rõ ràng, kể
		// cả khi không có từ khoá nào trong cửa sổ: "GAS, BSR, PVS" chỉ có thể
		// là liệt kê cổ phiếu.
		if hasStockContext(before, after, rawAfter) || neighbours >= TickerListNeighbours {
			st.withContext++
		}
	}
	return res
}

// knownCodesAround đếm số mã KHÁC trong bảng tickers xuất hiện trong cửa sổ.
func (t *Tagger) knownCodesAround(window, self string) int {
	seen := map[string]bool{}
	for _, code := range codeRe.FindAllString(window, -1) {
		if code == self || t.stoplist[code] || seen[code] {
			continue
		}
		if _, ok := t.symbols[code]; ok {
			seen[code] = true
		}
	}
	return len(seen)
}

// attributionPrefixes: cụm đứng NGAY TRƯỚC mã trong vai trò nguồn trích dẫn.
var attributionPrefixes = []string{
	"theo", "chứng khoán", "báo cáo của", "báo cáo từ", "nhóm phân tích",
	"nhóm nghiên cứu", "công ty chứng khoán", "số liệu của", "dự báo của",
	"ước tính của", "đánh giá của", "khuyến nghị của", "ctck",
}

// attributionSuffixes: cụm đứng NGAY SAU mã trong vai trò nguồn trích dẫn.
var attributionSuffixes = []string{
	"research", "securities", "ước tính", "dự báo", "dự phóng", "nhận định",
	"cho rằng", "đánh giá rằng", "khuyến nghị", "kỳ vọng", "công bố báo cáo",
}

// isAttribution nhận diện rule B3: "theo MBS Research", "ACBS ước tính"...
func isAttribution(before, after string) bool {
	b := strings.TrimRight(before, " ")
	for _, p := range attributionPrefixes {
		if b == p || strings.HasSuffix(b, " "+p) {
			return true
		}
	}
	a := strings.TrimLeft(after, " ")
	for _, s := range attributionSuffixes {
		if a == s || strings.HasPrefix(a, s+" ") || strings.HasPrefix(a, s+",") {
			return true
		}
	}
	return false
}

// stockContextKeywords — rule B4 của BA2 (đã bỏ dấu).
var stockContextKeywords = []string{
	"cổ phiếu", "mã ck", "mã chứng khoán", "cp", "doanh nghiệp", "công ty",
	"tập đoàn", "ngân hàng", "tăng", "giảm", "trần", "sàn", "khớp lệnh",
	"thị giá", "vốn hoá", "vốn hóa", "lợi nhuận", "doanh thu", "cổ đông",
	"lnst", "cbtt", "hose", "hnx", "upcom", "phiên", "đồng/cp", "triệu cp",
	"điểm", "mua ròng", "bán ròng", "vốn điều lệ", "niêm yết", "cổ tức",
}

// hasStockContext kiểm tra có ≥1 từ khoá ngữ cảnh chứng khoán trong cửa sổ ±80
// ký tự, hoặc ngay sau mã là biến động giá dạng "+4,31%".
func hasStockContext(before, after, rawAfter string) bool {
	if priceMoveRe.MatchString(rawAfter) {
		return true
	}
	win := before + " " + after
	for _, kw := range stockContextKeywords {
		if containsWord(win, kw) {
			return true
		}
	}
	return false
}

// priceMoveRe bắt biến động giá viết liền sau mã: "VIC +4,31%".
// BẮT BUỘC có khoảng trắng trước dấu, nếu không "VIC-2030" (tên kế hoạch) sẽ bị
// hiểu nhầm là "VIC giảm 2030".
var priceMoveRe = regexp.MustCompile(`^\s+[+\-]\s*\d`)

// isVNDirectOccurrence xử lý rule B2 cho mã VND.
func isVNDirectOccurrence(before, after string, hasVNDirect bool) bool {
	if hasVNDirect {
		return true
	}
	b := strings.TrimRight(before, " ")
	// Cặp tỷ giá "USD/VND", "EUR/VND": luôn là đơn vị tiền tệ.
	for _, cur := range []string{"usd", "eur", "jpy", "cny", "krw", "gbp"} {
		if strings.HasSuffix(b, cur+"/") || strings.HasSuffix(b, " "+cur) {
			return false
		}
	}
	// Đơn vị tiền tệ: đứng sau số hoặc sau "tỷ/triệu/nghìn/đồng".
	for _, cur := range []string{"tỷ", "triệu", "nghìn", "tỉ", "đồng", "ngàn"} {
		if b == cur || strings.HasSuffix(b, " "+cur) {
			return false
		}
	}
	if endsWithDigit(b) {
		return false
	}
	// Là mã khi đi kèm từ khoá cổ phiếu ngay trước.
	for _, p := range []string{"cổ phiếu", "mã", "cp", "chứng khoán"} {
		if b == p || strings.HasSuffix(b, " "+p) {
			return true
		}
	}
	a := strings.TrimLeft(after, " ")
	for _, s := range []string{"tăng", "giảm", "trần", "sàn", "khớp"} {
		if strings.HasPrefix(a, s+" ") || a == s {
			return true
		}
	}
	return false
}

// insideNegativeAlias cho biết lần xuất hiện của mã tại vị trí `at` có nằm
// trong một cụm phủ định của chính mã đó không.
func (t *Tagger) insideNegativeAlias(code, text string, at int) bool {
	negs := t.negatives[code]
	if len(negs) == 0 {
		return false
	}
	for _, neg := range negs {
		// Cụm phủ định luôn bắt đầu bằng chính mã (ví dụ "FPT Retail"), nên chỉ
		// cần soi đoạn ngay sau vị trí mã, có dư ra chút cho khác biệt khoảng
		// trắng và dấu câu.
		end := min(len(text), at+len(neg)*2+8)
		if strings.HasPrefix(normPhrase(text[at:end]), neg) {
			return true
		}
	}
	return false
}

// fundMarkers là các từ cho biết ký hiệu đứng cạnh là mã quỹ, không phải mã
// niêm yết trên sàn Việt Nam.
var fundMarkers = []string{"etf", "fund", "index", "vaneck", "vectors", "ishares", "fubon", "xtrackers"}

func isFundTicker(before, after string) bool {
	a := strings.TrimLeft(after, " ")
	for _, marker := range fundMarkers {
		if strings.HasPrefix(a, marker+" ") || a == marker ||
			strings.HasPrefix(a, marker+")") || strings.HasPrefix(a, marker+",") {
			return true
		}
	}
	b := strings.TrimRight(before, " ")
	for _, marker := range []string{"vaneck", "vectors", "ishares", "fubon", "xtrackers", "quỹ"} {
		if strings.HasSuffix(b, " "+marker) || b == marker {
			return true
		}
	}
	return false
}

func endsWithDigit(s string) bool {
	if s == "" {
		return false
	}
	r := []rune(s)
	return unicode.IsDigit(r[len(r)-1])
}

// ---------------------------------------------------------------- T2 scanning

func (t *Tagger) scanAliases(text string) map[string]int {
	res := map[string]int{}
	if text == "" {
		return res
	}
	padded := " " + normPhrase(text) + " "

	// Cụm phủ định phải được loại TRƯỚC khi đếm alias, nếu không "Điện Máy Xanh"
	// trong "CTCP Đầu tư Điện Máy Xanh (DMX)" vẫn kéo theo MWG (rule G-04).
	masked := make(map[string]string, len(t.negatives))
	for sym, negs := range t.negatives {
		blanked := padded
		for _, neg := range negs {
			blanked = blankPhrase(blanked, neg)
		}
		if blanked != padded {
			masked[sym] = blanked
		}
	}

	for alias, sym := range t.aliasToSym {
		haystack := padded
		if m, ok := masked[sym]; ok {
			haystack = m
		}
		if n := countNonAttributed(haystack, alias); n > 0 {
			res[sym] += n
		}
	}
	return res
}

// countNonAttributed đếm alias nhưng bỏ những lần nằm trong cụm trích dẫn nguồn.
// "Chứng khoán Vietcap nhận định…" là tên tổ chức phân tích, không phải chủ thể
// của tin — rule B3 trước đây chỉ áp cho mã trần nên alias vẫn lọt.
func countNonAttributed(padded, core string) int {
	n, start := 0, 0
	for start < len(padded) {
		i := strings.Index(padded[start:], core)
		if i < 0 {
			break
		}
		abs := start + i
		end := abs + len(core)
		leftOK := abs == 0 || padded[abs-1] == ' '
		rightOK := end == len(padded) || padded[end] == ' '
		if leftOK && rightOK && !isAttribution(padded[:abs], padded[end:]) {
			n++
		}
		start = abs + 1
	}
	return n
}

// blankPhrase thay mọi lần xuất hiện trọn vẹn của `core` bằng dấu cách, giữ
// nguyên độ dài để các chỉ số khác không lệch.
func blankPhrase(padded, core string) string {
	if core == "" {
		return padded
	}
	out := []byte(padded)
	start := 0
	for start < len(out) {
		i := strings.Index(string(out[start:]), core)
		if i < 0 {
			break
		}
		abs := start + i
		end := abs + len(core)
		if (abs == 0 || out[abs-1] == ' ') && (end == len(out) || out[end] == ' ') {
			for k := abs; k < end; k++ {
				out[k] = ' '
			}
		}
		start = abs + 1
	}
	return string(out)
}

// ---------------------------------------------------------------- T3 scanning

// scanSectors chạy tầng T3.
//
// Rule G-08 chỉ cho phép suy luận theo ngành khi bài THỰC SỰ nói về ngành và
// không nêu doanh nghiệp nào. Bản cũ chạy vô điều kiện cho mọi bài: QA đo được
// 29/31 false positive đến từ đây, kéo precision xuống ~52%. Một từ "ngân hàng"
// thoáng qua trong bài về đường sắt Trung Quốc từng đủ để gắn 8 mã ngân hàng.
func (t *Tagger) scanSectors(title, body string) []string {
	if title == "" && body == "" {
		return nil
	}
	paddedTitle := " " + normPhrase(title) + " "
	paddedBody := " " + normPhrase(body) + " "

	seen := map[string]bool{}
	var out []string
	for kw, sector := range t.sectorKw {
		// Ngành phải là chủ đề thật sự: nêu ở tiêu đề, hoặc được nhắc lại nhiều
		// lần trong bài. Một lần nhắc thoáng qua giữa thân bài thì không.
		if countPhrase(paddedTitle, kw) == 0 &&
			countPhrase(paddedTitle, kw)+countPhrase(paddedBody, kw) < SectorMentionsRequired {
			continue
		}
		for _, sym := range t.sectorToSyms[sector] {
			if !seen[sym] {
				seen[sym] = true
				out = append(out, sym)
			}
		}
	}

	// Ưu tiên mã tiêu biểu (VN30) rồi mới tới thứ tự chữ cái, sau đó cắt bớt:
	// một tin ngành gắn cả rổ mã làm cột "Mã CK" mất hết ý nghĩa.
	sort.Slice(out, func(i, j int) bool {
		if t.vn30[out[i]] != t.vn30[out[j]] {
			return t.vn30[out[i]]
		}
		return out[i] < out[j]
	})
	if len(out) > MaxSectorTickers {
		out = out[:MaxSectorTickers]
	}
	sort.Strings(out)
	return out
}

// ---------------------------------------------------------------- helpers

// normPhrase chuẩn hoá để so khớp mà GIỮ NGUYÊN DẤU.
//
// Bản trước bỏ dấu, nên alias "Khí Việt Nam" của GAS khớp cụm thường
// "Khi Việt Nam" và gắn GAS primary cho bài không liên quan. Dấu tiếng Việt
// phân biệt nghĩa nên không được bỏ ở bước so khớp alias.
func normPhrase(s string) string {
	return textutil.Fold(s)
}

func countPhrase(padded, core string) int {
	return textutil.CountPhrase(padded, core)
}

// containsWord kiểm tra cụm từ xuất hiện trọn vẹn (theo biên từ).
func containsWord(text, phrase string) bool {
	return countPhrase(" "+text+" ", phrase) > 0
}
