package dedupe_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/ingest/dedupe"
)

const originalBody = `VN-Index đóng cửa 1.821 điểm, tăng 29,91 điểm tương đương 1,67%, lần đầu vượt 1.800.

Giá trị giao dịch đạt khoảng 20.000 tỷ đồng, vượt bình quân 20 phiên gần nhất.

TCB tăng trần lên 33.450 đồng, khớp 39,17 triệu cp. VIC tăng 4,31% và FPT tăng 2,69%.

Khối ngoại mua ròng khoảng 14 tỷ đồng trên toàn thị trường nhưng bán ròng 43 tỷ đồng trên HOSE.`

// U-01: crawl lại bài không đổi phải ra đúng hash cũ.
func TestContentHashStableAcrossIdenticalCrawls(t *testing.T) {
	assert.Equal(t, dedupe.ContentHash(originalBody), dedupe.ContentHash(originalBody))
}

// U-06: phần biến động theo từng lần crawl không được ảnh hưởng tới hash.
// Không loại chúng thì mọi bài đều bị coi là "đã cập nhật" ở mỗi run.
func TestContentHashIgnoresVolatileFragments(t *testing.T) {
	base := dedupe.ContentHash(originalBody)

	variants := map[string]string{
		"khác giờ cập nhật":   "Cập nhật lúc 15:45 " + originalBody,
		"khác lượt xem":       originalBody + " 12.480 lượt xem",
		"khác số bình luận":   originalBody + " 37 bình luận",
		"khác thứ trong tuần": "Thứ Tư, " + originalBody,
		"khác khoảng trắng":   strings.ReplaceAll(originalBody, " ", "  "),
		"khác viết hoa":       strings.ToUpper(originalBody),
	}
	for name, variant := range variants {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, base, dedupe.ContentHash(variant))
		})
	}
}

func TestContentHashChangesWhenTextChanges(t *testing.T) {
	changed := strings.Replace(originalBody, "1.821 điểm", "1.845 điểm", 1)
	assert.NotEqual(t, dedupe.ContentHash(originalBody), dedupe.ContentHash(changed))
}

func TestChangeRatio(t *testing.T) {
	assert.Equal(t, 0.0, dedupe.ChangeRatio(originalBody, originalBody))

	typoFixed := strings.Replace(originalBody, "vượt bình quân", "vượt mức bình quân", 1)
	assert.Less(t, dedupe.ChangeRatio(originalBody, typoFixed), dedupe.RewriteThreshold,
		"sửa một từ phải nằm dưới ngưỡng viết lại")

	rewritten := "Thị trường chứng khoán Việt Nam ghi nhận phiên giao dịch ảm đạm với thanh khoản thấp. " +
		"Nhóm dầu khí và bán lẻ chịu áp lực bán mạnh trong suốt cả phiên chiều nay."
	assert.Greater(t, dedupe.ChangeRatio(originalBody, rewritten), 0.5)
}

// U-02: sửa nhẹ thì không gọi lại summarizer.
func TestNeedsResummarizeSkipsMinorEdits(t *testing.T) {
	typoFixed := strings.Replace(originalBody, "vượt bình quân", "vượt mức bình quân", 1)
	assert.False(t, dedupe.NeedsResummarize(
		"VN-Index đóng cửa 1.821 điểm", originalBody,
		"VN-Index đóng cửa 1.821 điểm", typoFixed))
}

// U-03: số liệu đổi là tín hiệu bắt buộc tóm tắt lại, kể cả khi chỉ đổi 1 ký tự.
func TestNeedsResummarizeWhenNumbersChange(t *testing.T) {
	changed := strings.Replace(originalBody, "1.821 điểm", "1.845 điểm", 1)
	assert.True(t, dedupe.NeedsResummarize(
		"VN-Index đóng cửa phiên hôm nay", originalBody,
		"VN-Index đóng cửa phiên hôm nay", changed))
}

// U-03: tiêu đề đổi cũng bắt buộc tóm tắt lại.
func TestNeedsResummarizeWhenTitleChanges(t *testing.T) {
	assert.True(t, dedupe.NeedsResummarize(
		"VN-Index đóng cửa 1.821 điểm", originalBody,
		"VN-Index lần đầu vượt mốc 1.800 điểm", originalBody))
}

func TestNeedsResummarizeWhenBodyRewritten(t *testing.T) {
	rewritten := "Thị trường chứng khoán Việt Nam ghi nhận phiên giao dịch ảm đạm với thanh khoản thấp. " +
		"Nhóm dầu khí và bán lẻ chịu áp lực bán mạnh trong suốt cả phiên chiều nay."
	require.True(t, dedupe.NeedsResummarize("Tiêu đề", originalBody, "Tiêu đề", rewritten))
}
