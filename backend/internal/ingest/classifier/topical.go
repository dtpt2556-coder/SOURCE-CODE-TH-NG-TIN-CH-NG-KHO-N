package classifier

import "github.com/tenpoint/tenpoint-api/internal/textutil"

// TenPoint là sản phẩm tin chứng khoán. Feed "kinh doanh" chung của tuoitre và
// vnexpress kéo theo tin lên men cà phê, kênh đào Trung Quốc, tuỳ bút về nhà
// xuất bản — QA đo 3/3 tin đầu trang chủ là tin ngoài chủ đề, 112/179 bài
// (62,6%) không có mã nào.
//
// Cổng này chạy SAU tagger: bài được giữ khi có mã CK, hoặc khi nói về tài
// chính vĩ mô đủ đậm. Không đạt thì `rejected` với lý do `off_topic` thay vì
// publish.

// financeKeywords là các chủ đề tài chính đáng lên digest dù không có mã nào.
var financeKeywords = []string{
	// Chính sách tiền tệ và vĩ mô
	"lãi suất", "tỷ giá", "lạm phát", "cpi", "gdp", "ngân hàng nhà nước", "nhnn",
	"chính sách tiền tệ", "room tín dụng", "tín dụng", "dư nợ", "nợ xấu",
	"trái phiếu", "tpdn", "xếp hạng tín nhiệm", "fed", "ecb", "cung tiền",
	// Thị trường chứng khoán
	"vn index", "vnindex", "vn30", "hnx index", "upcom", "hose", "hnx",
	"cổ phiếu", "chứng khoán", "thanh khoản", "vốn hoá", "vốn hóa",
	"khối ngoại", "nhà đầu tư nước ngoài", "mua ròng", "bán ròng",
	"nâng hạng", "ftse", "msci", "etf", "quỹ đầu tư", "danh mục đầu tư",
	"ipo", "niêm yết", "cổ tức", "cổ đông", "phát hành riêng lẻ",
	"uỷ ban chứng khoán", "ủy ban chứng khoán", "công ty chứng khoán",
	// Kết quả kinh doanh doanh nghiệp niêm yết
	"lợi nhuận sau thuế", "lnst", "doanh thu thuần", "biên lợi nhuận",
	"báo cáo tài chính", "kết quả kinh doanh", "vốn điều lệ", "đại hội cổ đông",
	"công bố thông tin", "cbtt", "m&a", "thoái vốn", "vốn fdi",
}

// MinFinanceSignals là số từ khoá tài chính riêng biệt tối thiểu để giữ một bài
// không gắn được mã nào. Một lần nhắc thoáng qua là chưa đủ: bài về kênh đào
// vẫn có thể nhắc "tín dụng" một lần.
const MinFinanceSignals = 3

// TopicalInput là dữ liệu cần để quyết định bài có thuộc chủ đề hay không.
type TopicalInput struct {
	Title       string
	Body        string
	TickerCount int
}

// TopicalVerdict là kết quả của cổng lọc chủ đề.
type TopicalVerdict struct {
	OnTopic bool
	Signals int
	Reason  string
}

// IsOnTopic quyết định bài có được lên digest hay không.
func IsOnTopic(in TopicalInput) TopicalVerdict {
	if in.TickerCount > 0 {
		return TopicalVerdict{OnTopic: true, Reason: "có mã chứng khoán"}
	}

	title := textutil.FoldPad(in.Title)
	body := textutil.FoldPad(in.Body)

	signals := 0
	titleHit := false
	for _, kw := range financeKeywords {
		inTitle := textutil.CountPhrase(title, kw) > 0
		if inTitle {
			titleHit = true
		}
		if inTitle || textutil.CountPhrase(body, kw) > 0 {
			signals++
		}
	}

	// Từ khoá tài chính ngay trên tiêu đề là tín hiệu mạnh; giữa thân bài thì
	// cần nhiều tín hiệu độc lập mới đủ.
	if titleHit && signals >= 2 {
		return TopicalVerdict{OnTopic: true, Signals: signals, Reason: "chủ đề tài chính nêu ở tiêu đề"}
	}
	if signals >= MinFinanceSignals {
		return TopicalVerdict{OnTopic: true, Signals: signals, Reason: "đủ tín hiệu tài chính trong bài"}
	}
	return TopicalVerdict{OnTopic: false, Signals: signals, Reason: "off_topic"}
}
