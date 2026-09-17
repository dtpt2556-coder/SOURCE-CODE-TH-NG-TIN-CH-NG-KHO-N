package summarizer

import (
	"regexp"
	"strings"

	"github.com/tenpoint/tenpoint-api/internal/textutil"
)

// Tư vấn đầu tư là nghiệp vụ có điều kiện — chỉ công ty chứng khoán được cấp
// phép mới được đưa khuyến nghị. Một bản tóm tắt tự động nói "nên mua" đặt cả
// sản phẩm vào rủi ro mức tồn vong (PRD mục 4 & 11).
//
// Mọi cụm dưới đây viết CÓ DẤU và được so khớp trên văn bản CÓ DẤU. Bản trước
// bỏ dấu hai vế trước khi so, nên "nền tảng" bị bắt vì trùng "nên tăng" và
// "căn bản" bị bắt vì trùng "cần bán" — 6/10 bài bị cách ly oan.

var advicePhrases = []string{
	// Khuyến nghị trực tiếp
	"nên mua", "nên bán", "nên nắm giữ", "nên duy trì", "nên gom",
	"nên giải ngân", "nên chốt lời", "nên cắt lỗ", "nên mua vào", "nên bán ra",
	"nên hạn chế", "nên ưu tiên", "nên tránh", "nên xem xét mua",
	"khuyến nghị", "khuyến cáo nhà đầu tư",
	// Ngôn ngữ môi giới
	"giá mục tiêu", "mục tiêu giá", "tránh mua đuổi", "mua đuổi",
	"canh mua", "canh bán", "bắt đáy", "chốt lời", "cắt lỗ", "lướt sóng",
	"khả quan", "kém khả quan", "outperform", "underperform",
	"overweight", "underweight",
	// Mời chào
	"tiềm năng tăng giá", "cơ hội đầu tư", "đáng để đầu tư", "nhà đầu tư nên",
	"có thể mua vào", "có thể bán ra", "có thể giải ngân",
	"giải ngân dần", "giải ngân thêm", "chủ động giải ngân", "vùng giải ngân",
	// Câu dẫn thừa và giọng không trung lập
	"theo thông tin từ", "bài viết cho biết", "như vậy có thể thấy",
	"tóm lại", "nhìn chung", "gây sốc", "lao dốc không phanh", "bốc hơi",
	"thổi bay", "chúng tôi", "chúng ta",
}

// advicePatterns bắt các mẫu sinh sôi mà liệt kê cứng không phủ hết.
//
// Hai từ phụ thuộc ngữ cảnh được xử lý riêng thay vì chặn trần:
//
//   - "tỷ trọng": "tỷ trọng ngành ngân hàng trong VN-Index là 30%" là dữ kiện.
//     Chỉ chặn khi đi với động từ hành động, lúc nó thành lời khuyên phân bổ.
//   - "giải ngân": "vốn FDI giải ngân 12 tỷ USD" là dữ kiện. Chỉ chặn ở dạng
//     hướng dẫn hành động.
var advicePatterns = []*regexp.Regexp{
	regexp.MustCompile(`nhà đầu tư (ngắn hạn |dài hạn |cá nhân )?(nên|có thể|được khuyến)`),
	regexp.MustCompile(`(nên|cần|hãy) (mua|bán|nắm giữ|duy trì|gom|giải ngân|chốt lời|cắt lỗ)`),
	regexp.MustCompile(`(duy trì|tăng|giảm|nâng|hạ|phân bổ|gia tăng|cắt giảm) tỷ trọng`),
	regexp.MustCompile(`tỷ trọng (cổ phiếu|danh mục) (ở mức|về mức|lên mức|xuống mức)`),
	regexp.MustCompile(`khuyến nghị (mua|bán|nắm giữ|khả quan|trung lập|theo dõi)`),
	regexp.MustCompile(`(vùng|điểm) (mua|bán|giải ngân) (tốt|hợp lý|phù hợp)`),
	regexp.MustCompile(`(tránh|hạn chế) (mua|bán) (đuổi|tháo|bằng mọi giá)`),
	regexp.MustCompile(`giải ngân (dần|thêm|mạnh|từng phần|ở vùng)`),
}

// FindBannedPhrase trả về cụm vi phạm đầu tiên tìm thấy.
//
// Biên từ dựng bằng textutil.Fold: mọi ký tự không phải chữ/số thành dấu cách,
// nên cụm cấm ở cuối câu, trong ngoặc hay sau dấu hai chấm đều bị bắt. Không
// dùng \b của regexp: RE2 coi ký tự tiếng Việt có dấu là ký tự không phải từ
// nên biên từ sẽ sai ngay giữa chữ.
func FindBannedPhrase(summaryMD string) (string, bool) {
	flat := textutil.FoldPad(strings.ReplaceAll(summaryMD, "*", ""))

	for _, phrase := range advicePhrases {
		if textutil.CountPhrase(flat, phrase) > 0 {
			return phrase, true
		}
	}
	for _, pattern := range advicePatterns {
		if m := pattern.FindString(flat); m != "" {
			return m, true
		}
	}
	return "", false
}
