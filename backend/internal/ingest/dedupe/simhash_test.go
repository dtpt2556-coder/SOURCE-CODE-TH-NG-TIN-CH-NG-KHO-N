package dedupe_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/ingest/dedupe"
)

func TestNearIdenticalTitlesHaveSmallHammingDistance(t *testing.T) {
	cases := []struct {
		name string
		a    string
		b    string
	}{
		{
			name: "khác dấu câu và viết hoa",
			a:    "VN-Index đóng cửa 1.821 điểm, lần đầu vượt 1.800",
			b:    "VN-INDEX ĐÓNG CỬA 1.821 ĐIỂM — LẦN ĐẦU VƯỢT 1.800!",
		},
		{
			name: "khác cách gõ dấu",
			a:    "Hoà Phát khởi công dự án Dung Quất 3",
			b:    "Hòa Phát khởi công dự án Dung Quất 3",
		},
		{
			name: "khác stopword",
			a:    "12 ngân hàng cam kết 408.000 tỷ đồng tín dụng cho DNNVV",
			b:    "12 ngân hàng cam kết 408.000 tỷ đồng tín dụng của DNNVV",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := dedupe.HammingDistance(dedupe.TitleSimhash(tc.a), dedupe.TitleSimhash(tc.b))
			assert.LessOrEqual(t, d, dedupe.MaxTitleHammingDistance,
				"hai tiêu đề gần giống phải có khoảng cách Hamming nhỏ, thực tế %d", d)
			assert.True(t, dedupe.IsNearDuplicate(dedupe.TitleSimhash(tc.a), dedupe.TitleSimhash(tc.b)))
		})
	}
}

// Tiêu đề bị "xào" thêm một từ vẫn phải gần hơn hẳn so với tiêu đề khác chủ đề.
func TestRewrittenTitleIsMuchCloserThanUnrelatedTitle(t *testing.T) {
	base := dedupe.TitleSimhash("VN-Index đóng cửa 1.821 điểm, lần đầu vượt 1.800")
	rewritten := dedupe.TitleSimhash("VN-Index đóng cửa 1.821 điểm, lần đầu tiên vượt 1.800")
	unrelated := dedupe.TitleSimhash("Sáu tháng đầu 2026, Việt Nam nhập 494.000 tấn thịt trị giá 1,5 tỷ USD")

	near := dedupe.HammingDistance(base, rewritten)
	far := dedupe.HammingDistance(base, unrelated)

	assert.LessOrEqual(t, near, 8, "tiêu đề xào lại phải có khoảng cách nhỏ, thực tế %d", near)
	assert.Less(t, near*2, far, "khoảng cách gần (%d) phải nhỏ hơn hẳn khoảng cách xa (%d)", near, far)
}

func TestDifferentTitlesHaveLargeHammingDistance(t *testing.T) {
	a := dedupe.TitleSimhash("VN-Index đóng cửa 1.821 điểm, lần đầu vượt 1.800")
	b := dedupe.TitleSimhash("Sáu tháng đầu 2026, Việt Nam nhập 494.000 tấn thịt trị giá 1,5 tỷ USD")

	d := dedupe.HammingDistance(a, b)
	assert.Greater(t, d, 12, "hai tiêu đề khác hẳn phải có khoảng cách lớn, thực tế %d", d)
	assert.False(t, dedupe.IsNearDuplicate(a, b))
}

func TestIdenticalTitleHasZeroDistance(t *testing.T) {
	title := "CBAM: Doanh nghiệp thép, nhôm Việt đối mặt chi phí carbon lớn khi xuất sang EU"
	assert.Equal(t, 0, dedupe.HammingDistance(dedupe.TitleSimhash(title), dedupe.TitleSimhash(title)))
}

func TestEncodeDecodeSimhashRoundTrip(t *testing.T) {
	h := dedupe.TitleSimhash("ACBS ước 5.588 tỷ đồng chảy vào 117 cổ phiếu Việt Nam")
	encoded := dedupe.EncodeSimhash(h)
	require.Len(t, encoded, 16)

	decoded, err := dedupe.DecodeSimhash(encoded)
	require.NoError(t, err)
	assert.Equal(t, h, decoded)
}

func TestEmptyTitleIsZero(t *testing.T) {
	assert.Equal(t, uint64(0), dedupe.TitleSimhash(""))
}
