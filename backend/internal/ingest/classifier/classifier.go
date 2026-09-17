// Package classifier gán "Loại tin" bằng cây quyết định rule-based
// (§3.3 BA2 + thứ tự ưu tiên PRD mục 5). Tất định, không gọi LLM — QA chạy lại
// cho ra đúng kết quả.
package classifier

import (
	"github.com/tenpoint/tenpoint-api/internal/domain"
	"github.com/tenpoint/tenpoint-api/internal/textutil"
)

// Input là dữ liệu cần để phân loại.
type Input struct {
	Title          string
	Body           string
	PrimarySymbols []string
	SourceTier     int
}

// Từ khoá cho từng nhãn, viết CÓ DẤU và so khớp trên văn bản có dấu
// (textutil.Fold). Bỏ dấu ở bước này gây va chạm nghĩa — xem P1-NEW-1.
var (
	regulatory = []string{
		"xử phạt vi phạm", "quyết định xử phạt", "xử phạt hành chính",
		"đình chỉ giao dịch", "khởi tố", "thanh tra", "uỷ ban chứng khoán nhà nước",
		"ủy ban chứng khoán nhà nước", "vi phạm công bố thông tin",
		"huỷ niêm yết bắt buộc", "nghị định", "thông tư",
	}
	corpAction = []string{
		"cổ tức", "gdkhq", "ngày đăng ký cuối cùng", "chào bán",
		"phát hành riêng lẻ", "cổ phiếu quỹ", "cổ phiếu thưởng", "esop",
		"chia tách", "niêm yết bổ sung", "chuyển sàn", "trả cổ tức",
	}
	foreign = []string{
		"khối ngoại", "nhà đầu tư nước ngoài", "mua ròng", "bán ròng", "etf",
		"ftse", "msci", "nâng hạng", "room ngoại", "quỹ ngoại", "cơ cấu danh mục",
	}
	credit = []string{
		"trái phiếu doanh nghiệp", "tpdn", "trái phiếu", "chậm trả gốc lãi",
		"xếp hạng tín nhiệm", "nợ xấu", "lãi suất huy động", "room tín dụng",
		"tín dụng", "lãi suất cho vay", "gói tín dụng", "dư nợ",
	}
	analysis = []string{
		"research", "báo cáo phân tích", "giá mục tiêu", "dự phóng",
		"nhóm phân tích", "báo cáo chiến lược", "triển vọng ngành",
	}
	market = []string{
		"vn index", "vnindex", "vn 30", "vn30", "hnx index", "upcom index",
		"thanh khoản thị trường", "độ rộng thị trường",
	}
	macro = []string{
		"cpi", "gdp", "pmi", "iip", "lạm phát", "lãi suất điều hành", "tỷ giá",
		"nhnn", "ngân hàng nhà nước", "chính phủ", "bộ tài chính", "fed", "ecb",
		"xuất nhập khẩu", "fdi",
	}
	sector = []string{
		"ngành thép", "ngành ngân hàng", "ngành điện", "ngành phân bón",
		"ngành xi măng", "giá hàng hoá", "chuỗi giá trị", "toàn ngành",
		"các doanh nghiệp trong ngành", "nhóm ngành", "chăn nuôi", "cbam",
	}
)

// Classifier áp cây quyết định §3.3 BA2 theo thứ tự ưu tiên của PRD mục 5.
type Classifier struct{}

// New tạo classifier (stateless).
func New() *Classifier { return &Classifier{} }

// Classify trả về loại tin và độ tin cậy 0..1.
//
// Thứ tự PRD: Pháp lý > Cổ tức/Phát hành > Khối ngoại > Trái phiếu/Tín dụng >
// Phân tích > Doanh nghiệp > Ngành > Thị trường > Vĩ mô.
//
// Một sai lệch có chủ đích so với danh sách đó: khi TIÊU ĐỀ xoay quanh tên chỉ
// số, "Thị trường" được quyết ngay, trước các luật đếm mã. Bài tổng kết phiên
// luôn liệt kê vài mã kèm số liệu nên luật đếm sẽ gán nhầm thành "Ngành" —
// trong khi chủ thể thật của tin là chỉ số. Thứ tự PRD vẫn được giữ nguyên giữa
// các nhãn nhận diện bằng từ khoá.
func (c *Classifier) Classify(in Input) (domain.NewsType, float64) {
	title := textutil.FoldPad(in.Title)
	full := title + " " + textutil.FoldPad(in.Body)

	switch {
	case containsAny(full, regulatory):
		return domain.NewsTypePhapLy, 0.9
	case containsAny(full, corpAction):
		return domain.NewsTypeCoTucPhatHanh, 0.85
	case containsAny(title, foreign):
		return domain.NewsTypeKhoiNgoai, 0.85
	case containsAny(full, credit):
		return domain.NewsTypeTraiPhieuTinDung, 0.8
	case containsAny(full, analysis):
		return domain.NewsTypePhanTich, 0.75
	case containsAny(title, market):
		return domain.NewsTypeThiTruong, 0.85
	}

	switch n := len(in.PrimarySymbols); {
	case n == 1:
		return domain.NewsTypeDoanhNghiep, 0.8
	case n >= 3:
		return domain.NewsTypeNganh, 0.75
	}

	switch {
	case containsAny(full, sector):
		return domain.NewsTypeNganh, 0.7
	case containsAny(full, macro):
		return domain.NewsTypeViMo, 0.7
	case len(in.PrimarySymbols) == 2:
		return domain.NewsTypeNganh, 0.6
	}
	return domain.NewsTypeThiTruong, 0.5
}

func containsAny(foldedPadded string, keywords []string) bool {
	for _, kw := range keywords {
		if textutil.CountPhrase(foldedPadded, kw) > 0 {
			return true
		}
	}
	return false
}
