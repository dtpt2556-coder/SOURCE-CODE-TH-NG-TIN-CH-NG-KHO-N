package summarizer_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/ingest/summarizer"
)

const cbamBody = `Cơ chế điều chỉnh biên giới carbon (CBAM) của EU áp dụng từ đầu 2026.

Doanh nghiệp có nguy cơ bù đắp hơn 100 USD/tấn thép xuất sang EU, nhôm tới cả nghìn USD/tấn.

Giá chứng chỉ carbon EU hiện trên 75 USD/tấn CO2.

Một lô 100.000 tấn thép đơn giá 500-550 USD/tấn phải nộp khoảng 300 tỷ đồng, tương đương 20% giá trị đơn hàng.

Cường độ phát thải thép Việt Nam là 2,5 tấn CO2/tấn, nhôm 14 tấn so với chuẩn EU 1,4 tấn.

Sáu nhóm hàng chịu CBAM gồm sắt thép, nhôm, xi măng, phân bón, điện và hydrogen.`

func TestExtractiveProducesGroundedSummary(t *testing.T) {
	p := summarizer.NewExtractive()
	out, err := p.Summarize(context.Background(), summarizer.Input{
		Title:      "CBAM: Doanh nghiệp thép, nhôm Việt đối mặt chi phí carbon lớn khi xuất sang EU",
		Body:       cbamBody,
		SourceName: "VnExpress",
		Tickers:    []string{"HPG", "HSG"},
	})
	require.NoError(t, err)
	require.NotEmpty(t, out.SummaryMD)

	// Bất biến quan trọng nhất: tóm tắt extractive giữ nguyên văn nên MỌI con số
	// phải truy vết được về bài gốc.
	res := summarizer.NewValidator().Validate(out.SummaryMD, cbamBody)
	assert.True(t, res.Grounded, "tóm tắt extractive phải luôn qua validator, vi phạm: %v", res.Violations)
}

func TestExtractiveNameIsStable(t *testing.T) {
	assert.Equal(t, "extractive", summarizer.NewExtractive().Name())
}

func TestExtractiveBoldsExactlyOneSpan(t *testing.T) {
	p := summarizer.NewExtractive()
	out, err := p.Summarize(context.Background(), summarizer.Input{
		Title: "CBAM áp lên thép Việt",
		Body:  cbamBody,
	})
	require.NoError(t, err)

	spans := summarizer.BoldSpans(out.SummaryMD)
	require.Len(t, spans, 1, "rubric C5: đúng 1 cụm bold, nhận được %q", out.SummaryMD)
	assert.Equal(t, spans, out.BoldSpans)
	assert.Equal(t, 2, strings.Count(out.SummaryMD, "**"))
}

func TestExtractiveRespectsWordCap(t *testing.T) {
	p := summarizer.NewExtractive()
	out, err := p.Summarize(context.Background(), summarizer.Input{Title: "CBAM", Body: cbamBody})
	require.NoError(t, err)

	words := len(strings.Fields(out.SummaryMD))
	assert.LessOrEqual(t, words, summarizer.TargetMaxWords, "tóm tắt %d từ, vượt trần", words)
	assert.Greater(t, words, 10)
}

func TestExtractivePrefersNumberRichSentences(t *testing.T) {
	body := `Đây là một câu mở đầu hoàn toàn không chứa số liệu nào cả và khá dài dòng để lọt ngưỡng.
Một câu khác cũng không có dữ kiện định lượng nào đáng chú ý cho nhà đầu tư chuyên nghiệp.
VN-Index đóng cửa 1.821 điểm, tăng 29,91 điểm tương đương 1,67%, thanh khoản 20.000 tỷ đồng.`

	out, err := summarizer.NewExtractive().Summarize(context.Background(),
		summarizer.Input{Title: "Thị trường", Body: body})
	require.NoError(t, err)
	assert.Contains(t, out.SummaryMD, "1.821", "phải ưu tiên câu giàu số liệu")
}

func TestExtractiveRejectsEmptyBody(t *testing.T) {
	_, err := summarizer.NewExtractive().Summarize(context.Background(),
		summarizer.Input{Title: "Tiêu đề", Body: "   "})
	assert.Error(t, err)
}

func TestSplitSentencesKeepsVietnameseNumbers(t *testing.T) {
	got := summarizer.SplitSentences("Cam kết 408.000 tỷ đồng tín dụng. Big4 đăng ký 220.000 tỷ.")
	require.Len(t, got, 2, "dấu chấm trong 408.000 không được cắt câu, nhận được %v", got)
	assert.Contains(t, got[0], "408.000")
	assert.Contains(t, got[1], "220.000")
}

func TestExtractiveSummaryHasNoBannedPhrase(t *testing.T) {
	out, err := summarizer.NewExtractive().Summarize(context.Background(),
		summarizer.Input{Title: "CBAM", Body: cbamBody})
	require.NoError(t, err)
	_, found := summarizer.FindBannedPhrase(out.SummaryMD)
	assert.False(t, found)
}

// TC-VAL-09: rubric NUM-3 yêu cầu 1–3 cụm bold. QA tìm thấy 3 tin có 0 cụm vì
// bài không chứa con số nào đáng nhấn.
func TestExtractiveAlwaysProducesBoldSpan(t *testing.T) {
	body := `Ngân hàng Nhà nước công bố định hướng điều hành chính sách tiền tệ cho giai đoạn tới.

Cơ quan này nhấn mạnh mục tiêu ổn định mặt bằng lãi suất và hỗ trợ tăng trưởng tín dụng bền vững.

Các tổ chức tín dụng được yêu cầu tiết giảm chi phí hoạt động để có dư địa hỗ trợ khách hàng.

Định hướng cũng nêu rõ yêu cầu kiểm soát chặt chẽ chất lượng tín dụng trong toàn hệ thống.`

	out, err := summarizer.NewExtractive().Summarize(context.Background(),
		summarizer.Input{Title: "Định hướng điều hành chính sách tiền tệ", Body: body})
	require.NoError(t, err)

	spans := summarizer.BoldSpans(out.SummaryMD)
	require.Len(t, spans, 1, "bài không có số vẫn phải có đúng 1 cụm bold: %q", out.SummaryMD)
	assert.Equal(t, spans, out.BoldSpans)
	assert.NotEmpty(t, strings.TrimSpace(spans[0]))
}
