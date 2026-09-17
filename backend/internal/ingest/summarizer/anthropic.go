package summarizer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/tenpoint/tenpoint-api/internal/textutil"
)

// DefaultAnthropicEndpoint là endpoint Messages API.
const DefaultAnthropicEndpoint = "https://api.anthropic.com/v1/messages"

// anthropicVersion là header bắt buộc của Anthropic API.
const anthropicVersion = "2023-06-01"

// Anthropic gọi Claude API để tóm tắt. Chỉ dùng net/http chuẩn — không cần SDK.
type Anthropic struct {
	APIKey   string
	Model    string
	Endpoint string
	HTTP     *http.Client
}

// NewAnthropic dựng provider. Endpoint rỗng => dùng endpoint mặc định.
func NewAnthropic(apiKey, model string, timeout time.Duration) *Anthropic {
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	return &Anthropic{
		APIKey:   apiKey,
		Model:    model,
		Endpoint: DefaultAnthropicEndpoint,
		HTTP:     &http.Client{Timeout: timeout},
	}
}

// Name thoả Provider.
func (a *Anthropic) Name() string { return "anthropic" }

type anthropicRequest struct {
	Model       string             `json:"model"`
	MaxTokens   int                `json:"max_tokens"`
	Temperature float64            `json:"temperature"`
	System      string             `json:"system"`
	Messages    []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

type modelOutput struct {
	SummaryMD       string   `json:"summary_md"`
	HeadlineNumber  string   `json:"headline_number"`
	Confidence      float64  `json:"confidence"`
	OmittedBecauseU []string `json:"omitted_because_uncertain"`
}

// Summarize gọi Claude với prompt ép rubric NUM-1..NUM-8.
func (a *Anthropic) Summarize(ctx context.Context, in Input) (Output, error) {
	if a.APIKey == "" {
		return Output{}, fmt.Errorf("ANTHROPIC_API_KEY is empty")
	}
	endpoint := a.Endpoint
	if endpoint == "" {
		endpoint = DefaultAnthropicEndpoint
	}
	client := a.HTTP
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second}
	}

	payload := anthropicRequest{
		Model:       a.Model,
		MaxTokens:   1024,
		Temperature: 0,
		System:      systemPrompt,
		Messages: []anthropicMessage{
			{Role: "user", Content: buildUserPrompt(in)},
		},
	}
	buf, err := json.Marshal(payload)
	if err != nil {
		return Output{}, fmt.Errorf("encode Anthropic request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(buf))
	if err != nil {
		return Output{}, fmt.Errorf("build Anthropic request: %w", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", a.APIKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	resp, err := client.Do(req)
	if err != nil {
		return Output{}, fmt.Errorf("call Anthropic API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Output{}, fmt.Errorf("read Anthropic response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return Output{}, fmt.Errorf("Anthropic API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed anthropicResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return Output{}, fmt.Errorf("decode Anthropic response: %w", err)
	}
	if parsed.Error != nil {
		return Output{}, fmt.Errorf("Anthropic API error %s: %s", parsed.Error.Type, parsed.Error.Message)
	}

	var text strings.Builder
	for _, c := range parsed.Content {
		if c.Type == "text" {
			text.WriteString(c.Text)
		}
	}
	raw := strings.TrimSpace(text.String())
	if raw == "" {
		return Output{}, fmt.Errorf("Anthropic API returned empty content")
	}

	summary := extractSummaryMD(raw)
	if summary == "" {
		return Output{}, fmt.Errorf("no summary_md in Anthropic response")
	}
	return Output{SummaryMD: summary, BoldSpans: BoldSpans(summary)}, nil
}

// extractSummaryMD đọc JSON theo hợp đồng prompt; nếu model trả văn bản thuần
// thì dùng luôn văn bản đó (fail-soft, validator vẫn chặn số bịa).
func extractSummaryMD(raw string) string {
	jsonPart := raw
	if i := strings.Index(raw, "{"); i >= 0 {
		if j := strings.LastIndex(raw, "}"); j > i {
			jsonPart = raw[i : j+1]
		}
	}
	var out modelOutput
	if err := json.Unmarshal([]byte(jsonPart), &out); err == nil && strings.TrimSpace(out.SummaryMD) != "" {
		return strings.TrimSpace(out.SummaryMD)
	}
	if strings.HasPrefix(strings.TrimSpace(raw), "{") {
		return ""
	}
	return raw
}

// MaxPromptBodyChars — bài dài hơn ngưỡng này bị cắt trước khi gửi LLM: giữ
// phần đầu (chứa số liệu headline) cộng đoạn có mật độ số cao nhất. Gửi nguyên
// bài 40.000 ký tự vừa tốn tiền vừa dễ vượt context (S-08).
const MaxPromptBodyChars = 15000

func trimBodyForPrompt(body string) string {
	if len(body) <= MaxPromptBodyChars {
		return body
	}
	head := MaxPromptBodyChars / 2
	paragraphs := strings.Split(body[head:], "\n\n")
	sort.SliceStable(paragraphs, func(i, j int) bool {
		return len(textutil.RawNumbers(paragraphs[i])) > len(textutil.RawNumbers(paragraphs[j]))
	})

	var tail strings.Builder
	for _, para := range paragraphs {
		if tail.Len()+len(para) > MaxPromptBodyChars-head {
			break
		}
		tail.WriteString(para)
		tail.WriteString("\n\n")
	}
	return body[:head] + "\n\n" + tail.String()
}

func buildUserPrompt(in Input) string {
	var b strings.Builder
	b.WriteString("=== METADATA ===\n")
	fmt.Fprintf(&b, "Nguồn: %s\n", in.SourceName)
	fmt.Fprintf(&b, "Mã CK liên quan (đã xác định): %s\n\n", strings.Join(in.Tickers, ", "))
	b.WriteString("=== TIÊU ĐỀ ===\n")
	b.WriteString(in.Title)
	b.WriteString("\n\n=== NỘI DUNG BÀI GỐC ===\n")
	b.WriteString(trimBodyForPrompt(in.Body))
	b.WriteString("\n\n=== YÊU CẦU ===\nViết tóm tắt theo đúng quy tắc trong system prompt.\n")
	if len(in.Tickers) > 0 {
		fmt.Fprintf(&b, "Ưu tiên các số liệu liên quan trực tiếp tới %s.\n", strings.Join(in.Tickers, ", "))
	}
	return b.String()
}

// systemPrompt ép rubric NUM-1..NUM-8 của BA1/BA2 §5.5.
const systemPrompt = `Bạn là biên tập viên tài chính của TenPoint, chuyên tóm tắt tin tức chứng khoán
Việt Nam thành các bản tin cô đọng, giàu số liệu, dành cho nhà đầu tư chuyên nghiệp.

NHIỆM VỤ
Viết một đoạn tóm tắt tiếng Việt từ bài báo được cung cấp.

QUY TẮC BẤT BIẾN — vi phạm bất kỳ quy tắc nào dưới đây là lỗi nghiêm trọng:

NUM-1. TRUNG THỰC SỐ LIỆU
   - Chỉ được dùng những con số XUẤT HIỆN NGUYÊN VĂN trong bài gốc.
   - TUYỆT ĐỐI KHÔNG tự tính toán, không tự quy đổi đơn vị, không làm tròn,
     không suy ra tỷ lệ phần trăm mới.
   - Nếu không chắc một con số có trong bài, hãy BỎ con số đó.

NUM-2. ĐỊNH DẠNG SỐ KIỂU VIỆT NAM
   - Dấu chấm phân cách hàng nghìn: 408.000 tỷ đồng
   - Dấu phẩy phân cách thập phân: 1,67% ; 29,91 điểm
   - Giữ nguyên cách viết của bài gốc.

NUM-3. GIỮ NGUYÊN ĐƠN VỊ GỐC
   - tỷ đồng, nghìn tỷ đồng, triệu USD, %, điểm %, USD/tấn, đồng/kg,
     triệu cp, điểm, tấn CO2/tấn. Không tự quy đổi.

NUM-4. ĐỘ DÀI
   - 60 đến 150 từ, mục tiêu 80-120 từ. Viết 4-6 câu liền mạch, KHÔNG bullet.

NUM-5. CẤU TRÚC
   - Câu 1: con số/sự kiện quan trọng nhất, đủ ngữ cảnh để hiểu độc lập.
   - Câu 2-5: bóc tách chi tiết, so sánh cùng kỳ, các bên liên quan.
   - Ưu tiên chi tiết CÓ SỐ hơn chi tiết định tính.

NUM-6. IN ĐẬM
   - Dùng từ 1 đến 3 cụm in đậm bằng markdown **...**, ưu tiên ĐÚNG MỘT cụm.
   - Cụm bold đầu tiên phải nằm trong câu đầu tiên và là con số quan trọng nhất.
   - Không dùng HTML, chỉ dùng **...**.

NUM-7. GIỌNG VĂN TRUNG LẬP
   - Mô tả như thông tấn xã. Không tính từ cảm thán, không emoji, không câu hỏi.
   - KHÔNG câu dẫn nhập ("Theo bài viết...", "Nhìn chung...", "Tóm lại...").
   - KHÔNG dùng ngôi thứ nhất.

NUM-8. CẤM KHUYẾN NGHỊ ĐẦU TƯ — QUY ĐỊNH PHÁP LÝ
   - TUYỆT ĐỐI KHÔNG viết: nên mua, nên bán, nên nắm giữ, khuyến nghị,
     tiềm năng tăng giá, cơ hội đầu tư, canh mua, chốt lời, bắt đáy.
   - Nếu bài gốc có khuyến nghị của một công ty chứng khoán, chỉ được TƯỜNG
     THUẬT kèm tên tổ chức đó, không trình bày như quan điểm của người viết.

ĐẦU RA
Trả về DUY NHẤT một JSON hợp lệ, không kèm giải thích:
{
  "summary_md": "<đoạn tóm tắt, có cụm **bold**>",
  "headline_number": "<con số đã in đậm>",
  "confidence": <0.0-1.0>,
  "omitted_because_uncertain": ["<số đã bỏ vì không chắc>"]
}`
