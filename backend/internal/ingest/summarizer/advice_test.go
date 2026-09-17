package summarizer_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/ingest/summarizer"
)

// Bản cũ chỉ coi khoảng trắng và dấu phẩy là biên từ nên mọi cụm cấm đứng ở
// CUỐI CÂU đều lọt. Đây là 16 mẫu QA đo được, 14 trong số đó từng lọt.
func TestBannedPhraseAtEveryPunctuationBoundary(t *testing.T) {
	leaked := []string{
		"Chuyên gia khuyến nghị nắm giữ cổ phiếu.",
		"Chuyên gia khuyến nghị nắm giữ, theo báo cáo.",
		"Chuyên gia khuyến nghị nắm giữ.",
		"SSI khuyến nghị mua.",
		"VCSC khuyến nghị bán.",
		"Đây là cơ hội đầu tư.",
		"Nhà đầu tư có thể chốt lời.",
		"Cổ phiếu đang bắt đáy.",
		"Cổ phiếu có tiềm năng tăng giá.",
		"Nhà đầu tư có thể canh mua.",
		"Quan điểm: khuyến nghị nắm giữ.",
		"**Khuyến nghị nắm giữ.**",
		"Chuyên gia khuyến nghị nắm giữ;",
		"Chuyên gia khuyến nghị nắm giữ!",
		"Chuyên gia khuyến nghị nắm giữ?",
		"Chuyên gia khuyến nghị nắm giữ)",
	}
	for _, s := range leaked {
		phrase, found := summarizer.FindBannedPhrase(s)
		assert.True(t, found, "phải chặn %q", s)
		assert.NotEmpty(t, phrase)
	}
}

func TestBannedPhraseBoundaryVariants(t *testing.T) {
	// Cùng một cụm, đặt ở đủ mọi vị trí và kiểu dấu câu.
	for _, s := range []string{
		"Khuyến nghị mua cổ phiếu ngành thép.",
		"Báo cáo kết luận: khuyến nghị mua.",
		"Báo cáo kết luận (khuyến nghị mua).",
		"Báo cáo kết luận — khuyến nghị mua…",
		"KHUYẾN NGHỊ MUA.",
		"Khuyến nghị mua\ntrong phiên tới.",
		"\"Khuyến nghị mua\" là quan điểm của môi giới.",
		"Khuyến nghị mua/bán đều có rủi ro.",
	} {
		_, found := summarizer.FindBannedPhrase(s)
		assert.True(t, found, "phải chặn %q", s)
	}
}

// Câu thật đang publish mà QA tìm thấy ở tin id 8.
func TestBlocksRealLeakedAdviceSentence(t *testing.T) {
	leaked := "Trong bối cảnh này, nhà đầu tư ngắn hạn nên duy trì tỷ trọng cổ phiếu " +
		"ở mức trung bình, tránh mua đuổi và tập trung giao dịch theo vùng hỗ trợ " +
		"1.810-1.820 điểm và kháng cự 1.830-1.840 điểm."

	phrase, found := summarizer.FindBannedPhrase(leaked)
	require.True(t, found, "câu tư vấn phân bổ danh mục phải bị chặn")
	t.Logf("chặn bởi: %q", phrase)

	require.NotNil(t, summarizer.CheckQuality(leaked))
	assert.Equal(t, "banned_phrase", summarizer.CheckQuality(leaked).Code)
}

func TestBlocksAdvicePatternsNotInFixedList(t *testing.T) {
	for _, s := range []string{
		"Nhà đầu tư ngắn hạn nên duy trì tỷ trọng ở mức trung bình.",
		"Nhà đầu tư dài hạn có thể giải ngân dần.",
		"Chuyên gia cho rằng nên tăng tỷ trọng cổ phiếu ngân hàng.",
		"Báo cáo đề xuất giảm tỷ trọng danh mục xuống mức thấp.",
		"Vùng mua tốt quanh 1.800 điểm.",
		"Tránh mua đuổi bằng mọi giá.",
		"Khuyến nghị khả quan với ngành thép.",
	} {
		_, found := summarizer.FindBannedPhrase(s)
		assert.True(t, found, "phải chặn %q", s)
	}
}

// Tin tức trung lập KHÔNG được bị chặn nhầm — cổng này chạy trên mọi bài.
func TestNeutralFinancialReportingIsNotBlocked(t *testing.T) {
	for _, s := range []string{
		"VN-Index đóng cửa **1.821 điểm**, tăng 1,67%, thanh khoản 20.000 tỷ đồng.",
		"12 ngân hàng cam kết **408.000 tỷ đồng** tín dụng cho DNNVV theo chương trình NHNN chủ trì.",
		"Tỷ trọng ngành ngân hàng trong VN-Index đạt 30,2% vốn hoá toàn thị trường.",
		"Khối ngoại mua ròng 203 tỷ đồng cổ phiếu FPT trong phiên 26/8.",
		"CBAM của EU áp dụng từ đầu 2026, chi phí bù đắp **hơn 100 USD/tấn thép**.",
		"Vốn FDI giải ngân 6 tháng đạt 12,3 tỷ USD, tăng 8,4% so với cùng kỳ.",
	} {
		phrase, found := summarizer.FindBannedPhrase(s)
		assert.False(t, found, "không được chặn tin trung lập %q (khớp %q)", s, phrase)
	}
}

// "tỷ trọng" trần là dữ kiện hợp lệ; chỉ thành khuyến nghị khi đi với động từ
// hành động.
func TestPortfolioWeightOnlyBlockedWithActionVerb(t *testing.T) {
	_, found := summarizer.FindBannedPhrase("Tỷ trọng cổ phiếu ngân hàng chiếm 30% danh mục quỹ.")
	assert.False(t, found)

	_, found = summarizer.FindBannedPhrase("Quỹ khuyến nghị duy trì tỷ trọng cổ phiếu ở mức 70%.")
	assert.True(t, found)
}

func TestAdviceFilterIgnoresMarkdownEmphasis(t *testing.T) {
	_, found := summarizer.FindBannedPhrase("Báo cáo nêu **khuyến nghị mua**.")
	assert.True(t, found, "dấu ** không được che cụm cấm")
}

// Danh sách tối thiểu coordinator yêu cầu phủ, kiểm trong câu tư vấn thật.
// "giải ngân" và "tỷ trọng" cố ý chỉ chặn ở nghĩa khuyến nghị (xem advice.go),
// nên câu kiểm của chúng là câu hướng dẫn hành động chứ không phải từ trần.
func TestRequiredPhraseCoverage(t *testing.T) {
	sentences := map[string]string{
		"nên duy trì":    "Nhà đầu tư nên duy trì danh mục hiện tại.",
		"nên mua":        "Chuyên gia cho rằng nên mua cổ phiếu thép.",
		"nên bán":        "Báo cáo nêu nên bán khi chạm kháng cự.",
		"nên nắm giữ":    "Quan điểm: nên nắm giữ.",
		"nên gom":        "Nhà đầu tư nên gom thêm ở vùng giá thấp.",
		"tránh mua đuổi": "Tránh mua đuổi bằng mọi giá.",
		"khuyến nghị":    "Báo cáo đưa khuyến nghị.",
		"mục tiêu giá":   "Mục tiêu giá 45.000 đồng/cp.",
		"giá mục tiêu":   "Giá mục tiêu được nâng lên 45.000 đồng/cp.",
		"chốt lời":       "Nhà đầu tư có thể chốt lời.",
		"cắt lỗ":         "Cần cắt lỗ nếu thủng hỗ trợ.",
		"bắt đáy":        "Cổ phiếu đang bắt đáy.",
		"giải ngân":      "Nhà đầu tư nên giải ngân dần trong phiên tới.",
		"tỷ trọng":       "Nên tăng tỷ trọng cổ phiếu ngân hàng.",
		"canh mua":       "Nhà đầu tư có thể canh mua.",
		"canh bán":       "Nhà đầu tư có thể canh bán.",
		"khả quan":       "Báo cáo nêu khả quan với ngành thép.",
		"kém khả quan":   "Báo cáo nêu kém khả quan với ngành thép.",
		"outperform":     "Báo cáo nêu outperform.",
		"underperform":   "Báo cáo nêu underperform.",
	}
	for phrase, sentence := range sentences {
		_, found := summarizer.FindBannedPhrase(sentence)
		assert.True(t, found, "danh sách cấm phải phủ %q qua câu %q", phrase, sentence)
	}
}

func TestAdviceFilterHandlesUppercaseAndDiacritics(t *testing.T) {
	for _, s := range []string{
		"NHÀ ĐẦU TƯ NÊN MUA VÀO.",
		"Nhà Đầu Tư Nên Bán Ra.",
		"nhà đầu tư nên nắm giữ",
	} {
		_, found := summarizer.FindBannedPhrase(s)
		assert.True(t, found, "phải chặn %q", s)
	}
}

func TestFallbackRejectsAdviceAtEndOfSentence(t *testing.T) {
	// Cổng end-to-end: LLM trả câu tư vấn kết thúc bằng dấu chấm.
	llm := &stubProvider{name: "anthropic", outputs: []string{
		"VN-Index đóng cửa 1.821 điểm. Quan điểm: khuyến nghị nắm giữ.",
		"Nhà đầu tư nên duy trì tỷ trọng cổ phiếu ở mức trung bình.",
	}}

	out, provider, err := summarizer.SummarizeWithFallback(
		newContext(), llm, summarizer.NewExtractive(), input(), nil)

	require.NoError(t, err)
	assert.Equal(t, 2, llm.calls, "phải thử lại 1 lần rồi mới đổi provider")
	assert.Equal(t, "extractive", provider)
	_, banned := summarizer.FindBannedPhrase(out.SummaryMD)
	assert.False(t, banned)
	assert.False(t, strings.Contains(out.SummaryMD, "khuyến nghị"))
}

// P1-NEW-1: bỏ dấu trước khi so khớp làm "nền tảng" trùng "nên tăng" và
// "căn bản" trùng "cần bán". QA đo được 6/10 bài bị cách ly oan vì lớp lỗi này.
func TestDiacriticCollisionsDoNotTriggerAdviceFilter(t *testing.T) {
	innocent := []string{
		"Nền tảng thương mại điện tử ghi nhận 12 triệu đơn hàng trong quý III.",
		"Đây là nền tảng cho tăng trưởng dài hạn của doanh nghiệp.",
		"Về căn bản, cơ cấu nguồn vốn của ngân hàng không thay đổi.",
		"Yếu tố căn bản của thị trường vẫn được giữ vững.",
		"Nền tảng công nghệ mới giúp giảm 18% chi phí vận hành.",
		"Chính sách căn bản đã được Quốc hội thông qua từ năm 2025.",
	}
	for _, s := range innocent {
		phrase, found := summarizer.FindBannedPhrase(s)
		assert.False(t, found, "câu trung lập %q bị chặn nhầm bởi %q", s, phrase)
	}
}

// Cùng lúc đó, cụm cấm THẬT (có dấu đúng) vẫn phải bị bắt.
func TestRealAdvicePhrasesStillBlockedAfterDiacriticFix(t *testing.T) {
	guilty := []string{
		"Chuyên gia cho rằng nên tăng tỷ trọng cổ phiếu ngân hàng.",
		"Nhà đầu tư cần bán bớt khi chỉ số chạm kháng cự.",
		"Nhà đầu tư nên duy trì tỷ trọng ở mức trung bình.",
		"Báo cáo khuyến nghị mua.",
	}
	for _, s := range guilty {
		_, found := summarizer.FindBannedPhrase(s)
		assert.True(t, found, "câu tư vấn %q phải bị chặn", s)
	}
}
