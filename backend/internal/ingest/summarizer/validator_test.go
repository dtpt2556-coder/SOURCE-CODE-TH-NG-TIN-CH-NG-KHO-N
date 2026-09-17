package summarizer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/ingest/summarizer"
	"github.com/tenpoint/tenpoint-api/internal/textutil"
)

// body mẫu lấy từ tin "12 ngân hàng cam kết 408.000 tỷ đồng" của ảnh tham chiếu.
const creditBody = `12 ngân hàng cam kết 408.000 tỷ đồng tín dụng cho DNNVV theo chương trình ` +
	`do Ngân hàng Nhà nước chủ trì. Nhóm Big4 đăng ký 220.000 tỷ đồng, trong đó Agribank ` +
	`70.000 tỷ đồng; BIDV, Vietcombank và VietinBank mỗi bên 50.000 tỷ đồng. Lãi suất thấp ` +
	`hơn ít nhất 1 điểm % so với bình quân cùng kỳ hạn. 8 ngân hàng tư nhân đăng ký 188.000 ` +
	`tỷ đồng, giảm lãi 0,5-2 điểm % tuỳ lĩnh vực. Cuối tháng 7 chỉ 8,8% DNNVV tiếp cận được ` +
	`vốn so với trên 47% ở doanh nghiệp lớn.`

const marketBody = `VN-Index đóng cửa 1.821 điểm, tăng 29,91 điểm tương đương 1,67%, ` +
	`lần đầu vượt 1.800. Giá trị giao dịch khoảng 20.000 tỷ đồng với 182 mã tăng và 125 mã giảm.`

func TestValidatorPassesWhenAllNumbersGrounded(t *testing.T) {
	v := summarizer.NewValidator()
	summary := "12 ngân hàng cam kết **408.000 tỷ đồng** tín dụng cho DNNVV. " +
		"Big4 đăng ký 220.000 tỷ đồng, trong đó Agribank 70.000 tỷ đồng."

	res := v.Validate(summary, creditBody)
	assert.True(t, res.Grounded, "mọi số đều có trong bài gốc, vi phạm: %v", res.Violations)
	assert.Empty(t, res.Violations)
	assert.Empty(t, res.Error())
}

func TestValidatorFailsOnFabricatedNumber(t *testing.T) {
	v := summarizer.NewValidator()
	// 999.000 tỷ đồng KHÔNG có trong bài gốc — đây chính là ảo giác số liệu.
	summary := "12 ngân hàng cam kết **999.000 tỷ đồng** tín dụng cho DNNVV."

	res := v.Validate(summary, creditBody)
	require.False(t, res.Grounded)
	assert.Contains(t, res.Violations, "999.000")
	assert.Contains(t, res.Error(), "999.000")
}

func TestValidatorHandlesVietnameseThousandSeparator(t *testing.T) {
	v := summarizer.NewValidator()

	// "408.000" phải được hiểu là 408000, không phải 408,0.
	res := v.Validate("Cam kết **408.000 tỷ đồng** tín dụng.", creditBody)
	assert.True(t, res.Grounded, "vi phạm: %v", res.Violations)

	// Trong khi 408 (không có phần nghìn) KHÔNG xuất hiện trong bài -> fail.
	res = v.Validate("Cam kết 408 tỷ đồng tín dụng.", creditBody)
	assert.False(t, res.Grounded)
	assert.Contains(t, res.Violations, "408")
}

func TestValidatorHandlesVietnameseDecimalComma(t *testing.T) {
	v := summarizer.NewValidator()

	res := v.Validate("VN-Index tăng **1,67%** lên 1.821 điểm.", marketBody)
	assert.True(t, res.Grounded, "1,67 có trong bài gốc, vi phạm: %v", res.Violations)

	res = v.Validate("VN-Index tăng 1,68% lên 1.821 điểm.", marketBody)
	assert.False(t, res.Grounded, "1,68 không có trong bài gốc")
	assert.Contains(t, res.Violations, "1,68")
}

func TestNormalizeVietnameseNumberFormat(t *testing.T) {
	cases := []struct {
		raw  string
		want string
		ok   bool
	}{
		{"408.000", "408000", true},
		{"1,67", "1.67", true},
		{"29,91", "29.91", true},
		{"1.821", "1821", true},
		{"20.000", "20000", true},
		{"0,5", "0.5", true},
		{"12", "12", true},
		{"5.588", "5588", true},
		{"2026.", "2026", true},
		{",,", "", false},
	}
	for _, tc := range cases {
		got, ok := textutil.NormalizeVNNumber(tc.raw)
		assert.Equal(t, tc.ok, ok, "raw=%q", tc.raw)
		if tc.ok {
			assert.Equal(t, tc.want, got, "raw=%q", tc.raw)
		}
	}
}

func TestValidatorIgnoresDuplicateViolations(t *testing.T) {
	v := summarizer.NewValidator()
	res := v.Validate("Con số 777.777 và lại 777.777 lần nữa.", creditBody)
	require.False(t, res.Grounded)
	assert.Len(t, res.Violations, 1)
}

func TestBoldSpans(t *testing.T) {
	spans := summarizer.BoldSpans("Cam kết **408.000 tỷ đồng** và **220.000 tỷ**.")
	assert.Equal(t, []string{"408.000 tỷ đồng", "220.000 tỷ"}, spans)
	assert.Empty(t, summarizer.BoldSpans("Không có cụm in đậm."))
}

func TestFindBannedPhrase(t *testing.T) {
	banned := []string{
		"Nhà đầu tư nên mua cổ phiếu này.",
		"Cổ phiếu có tiềm năng tăng giá mạnh.",
		"Nhìn chung, thị trường tích cực.",
		"Chúng tôi cho rằng xu hướng còn tiếp diễn.",
	}
	for _, s := range banned {
		_, found := summarizer.FindBannedPhrase(s)
		assert.True(t, found, "phải chặn: %q", s)
	}

	clean := "VN-Index đóng cửa 1.821 điểm, tăng 1,67%, thanh khoản 20.000 tỷ đồng."
	_, found := summarizer.FindBannedPhrase(clean)
	assert.False(t, found, "tóm tắt trung lập không được bị chặn")
}
