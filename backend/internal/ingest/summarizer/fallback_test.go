package summarizer_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/ingest/summarizer"
)

// stubProvider giả lập một LLM: trả lần lượt các kết quả đã nạp sẵn.
type stubProvider struct {
	name    string
	outputs []string
	err     error
	calls   int
}

func (s *stubProvider) Name() string { return s.name }

func (s *stubProvider) Summarize(context.Context, summarizer.Input) (summarizer.Output, error) {
	s.calls++
	if s.err != nil {
		return summarizer.Output{}, s.err
	}
	idx := s.calls - 1
	if idx >= len(s.outputs) {
		idx = len(s.outputs) - 1
	}
	return summarizer.Output{SummaryMD: s.outputs[idx]}, nil
}

func newContext() context.Context { return context.Background() }

func input() summarizer.Input {
	return summarizer.Input{Title: "CBAM áp lên thép Việt", Body: cbamBody, SourceName: "VnExpress"}
}

// S-05: LLM timeout / 429 / hết quota thì rơi về extractive, pipeline không dừng.
func TestFallbackWhenProviderFails(t *testing.T) {
	llm := &stubProvider{name: "anthropic", err: errors.New("429 rate_limit_error")}

	out, provider, err := summarizer.SummarizeWithFallback(
		context.Background(), llm, summarizer.NewExtractive(), input(), nil)

	require.NoError(t, err)
	assert.Equal(t, "extractive", provider, "phải ghi nhận provider thực sự sinh ra tóm tắt")
	assert.NotEmpty(t, out.SummaryMD)
	assert.Equal(t, 1, llm.calls)
}

// S-06: nội dung có từ khuyến nghị đầu tư thì sinh lại 1 lần, vẫn sai thì dùng
// extractive.
func TestFallbackRetriesOnceThenSwitchesProvider(t *testing.T) {
	llm := &stubProvider{name: "anthropic", outputs: []string{
		"Nhà đầu tư nên mua cổ phiếu này vì **408.000 tỷ đồng** tín dụng sắp được giải ngân.",
		"Cổ phiếu có tiềm năng tăng giá mạnh trong quý tới.",
	}}

	out, provider, err := summarizer.SummarizeWithFallback(
		context.Background(), llm, summarizer.NewExtractive(), input(), nil)

	require.NoError(t, err)
	assert.Equal(t, 2, llm.calls, "phải thử lại đúng 1 lần")
	assert.Equal(t, "extractive", provider)
	_, banned := summarizer.FindBannedPhrase(out.SummaryMD)
	assert.False(t, banned)
}

func TestFallbackAcceptsGoodRetry(t *testing.T) {
	llm := &stubProvider{name: "anthropic", outputs: []string{
		"Nhà đầu tư nên mua ngay cổ phiếu thép.",
		"CBAM của EU áp dụng từ đầu 2026, doanh nghiệp bù đắp **hơn 100 USD/tấn thép** xuất sang EU.",
	}}

	out, provider, err := summarizer.SummarizeWithFallback(
		context.Background(), llm, summarizer.NewExtractive(), input(), nil)

	require.NoError(t, err)
	assert.Equal(t, "anthropic", provider)
	assert.Equal(t, 2, llm.calls)
	assert.Contains(t, out.SummaryMD, "100 USD/tấn thép")
}

func TestFallbackKeepsGoodFirstAnswer(t *testing.T) {
	llm := &stubProvider{name: "anthropic", outputs: []string{
		"CBAM của EU áp dụng từ đầu 2026, chi phí bù đắp **hơn 100 USD/tấn thép** khi xuất sang EU.",
	}}

	_, provider, err := summarizer.SummarizeWithFallback(
		context.Background(), llm, summarizer.NewExtractive(), input(), nil)

	require.NoError(t, err)
	assert.Equal(t, "anthropic", provider)
	assert.Equal(t, 1, llm.calls, "kết quả tốt ngay lần đầu thì không gọi lại")
}

// S-07: mọi markdown ngoài **bold** phải bị strip trước khi lưu.
func TestSanitizeMarkdown(t *testing.T) {
	cases := map[string]string{
		"## Tiêu đề\nNội dung **quan trọng** ở đây.": "Tiêu đề Nội dung **quan trọng** ở đây.",
		"- Mục một\n- Mục hai":                       "Mục một Mục hai",
		"Xem [bài gốc](https://cafef.vn/x.chn) ngay": "Xem bài gốc ngay",
		"Giá `33.450` đồng":                          "Giá 33.450 đồng",
	}
	for in, want := range cases {
		assert.Equal(t, want, summarizer.SanitizeMarkdown(in))
	}
	assert.Equal(t, 2, strings.Count(summarizer.SanitizeMarkdown("**408.000 tỷ đồng**"), "**"))
}

// S-09: tóm tắt quá dài bị loại.
func TestCheckQualityRejectsOverlongSummary(t *testing.T) {
	long := strings.Repeat("một số liệu quan trọng ", 60)
	issue := summarizer.CheckQuality(long)
	require.NotNil(t, issue)
	assert.Equal(t, "too_long", issue.Code)

	assert.Nil(t, summarizer.CheckQuality("VN-Index đóng cửa **1.821 điểm**, tăng 1,67%."))
}

func TestCheckQualityRejectsBannedPhrase(t *testing.T) {
	issue := summarizer.CheckQuality("Nhà đầu tư nên mua cổ phiếu này.")
	require.NotNil(t, issue)
	assert.Equal(t, "banned_phrase", issue.Code)
}

// S-10: sao chép nguyên văn quá dài là vi phạm ranh giới bản quyền.
func TestLongestVerbatimRun(t *testing.T) {
	body := "Cơ chế điều chỉnh biên giới carbon của Liên minh châu Âu sẽ được áp dụng chính thức từ đầu năm 2026 cho sáu nhóm hàng."

	copied := "Cơ chế điều chỉnh biên giới carbon của Liên minh châu Âu sẽ được áp dụng chính thức từ đầu năm 2026."
	assert.GreaterOrEqual(t, summarizer.LongestVerbatimRun(copied, body), summarizer.VerbatimRunLimit)

	paraphrased := "EU áp CBAM từ 2026 lên sáu nhóm hàng."
	assert.Less(t, summarizer.LongestVerbatimRun(paraphrased, body), summarizer.VerbatimRunLimit)
}
