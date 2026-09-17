package classifier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tenpoint/tenpoint-api/internal/domain"
	"github.com/tenpoint/tenpoint-api/internal/ingest/classifier"
)

func TestClassifyDecisionTree(t *testing.T) {
	c := classifier.New()

	cases := []struct {
		name string
		in   classifier.Input
		want domain.NewsType
	}{
		{
			name: "xử phạt của UBCKNN -> pháp lý",
			in: classifier.Input{
				Title: "UBCKNN xử phạt CTCP X 1,5 tỷ đồng",
				Body:  "Quyết định xử phạt vi phạm hành chính do công bố thông tin sai hạn.",
			},
			want: domain.NewsTypePhapLy,
		},
		{
			name: "chốt quyền cổ tức -> cổ tức/phát hành",
			in: classifier.Input{
				Title: "Doanh nghiệp chốt quyền trả cổ tức 15% tiền mặt",
				Body:  "Ngày đăng ký cuối cùng là 20/9, tỷ lệ 15% bằng tiền mặt.",
			},
			want: domain.NewsTypeCoTucPhatHanh,
		},
		{
			name: "khối ngoại ở tiêu đề -> khối ngoại",
			in: classifier.Input{
				Title: "Khối ngoại mua ròng 1.200 tỷ đồng phiên 26/8",
				Body:  "Dòng vốn ngoại tập trung vào nhóm ngân hàng.",
			},
			want: domain.NewsTypeKhoiNgoai,
		},
		{
			name: "chỉ số ở tiêu đề -> thị trường",
			in: classifier.Input{
				Title: "VN-Index đóng cửa 1.821 điểm, lần đầu vượt 1.800",
				Body:  "Thanh khoản đạt 20.000 tỷ đồng với 182 mã tăng.",
			},
			want: domain.NewsTypeThiTruong,
		},
		{
			name: "đúng 1 mã primary -> doanh nghiệp",
			in: classifier.Input{
				Title:          "FPT ký hợp đồng 100 triệu USD",
				Body:           "Doanh thu mảng xuất khẩu phần mềm tăng 25%.",
				PrimarySymbols: []string{"FPT"},
			},
			want: domain.NewsTypeDoanhNghiep,
		},
		{
			name: "≥3 mã primary -> ngành",
			in: classifier.Input{
				Title:          "Nhóm thép chịu áp lực chi phí carbon",
				Body:           "Chi phí bù đắp hơn 100 USD/tấn thép xuất sang EU.",
				PrimarySymbols: []string{"HPG", "HSG", "NKG", "TIS"},
			},
			want: domain.NewsTypeNganh,
		},
		{
			name: "từ khoá ngành, không có mã -> ngành",
			in: classifier.Input{
				Title: "CBAM của EU áp thuế carbon lên hàng nhập khẩu",
				Body:  "Sáu nhóm hàng chịu CBAM gồm sắt thép, nhôm, xi măng, phân bón, điện và hydrogen.",
			},
			want: domain.NewsTypeNganh,
		},
		{
			name: "chỉ số kinh tế vĩ mô -> vĩ mô",
			in: classifier.Input{
				Title: "CPI tháng 8 tăng 3,2%",
				Body:  "Lạm phát cơ bản duy trì dưới mục tiêu của Chính phủ.",
			},
			want: domain.NewsTypeViMo,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, confidence := c.Classify(tc.in)
			assert.Equal(t, tc.want, got)
			assert.Greater(t, confidence, 0.0)
			assert.LessOrEqual(t, confidence, 1.0)
		})
	}
}

func TestClassifyAlwaysReturnsValidNewsType(t *testing.T) {
	c := classifier.New()
	got, _ := c.Classify(classifier.Input{Title: "Một tiêu đề không rõ chủ đề", Body: "Nội dung chung chung."})
	assert.True(t, got.Valid(), "classifier phải luôn trả nhãn thuộc tập đóng, nhận được %q", got)
}

func TestNewsTypeLabels(t *testing.T) {
	assert.Equal(t, "Ngành", domain.NewsTypeNganh.Label())
	assert.Equal(t, "Vĩ mô", domain.NewsTypeViMo.Label())
	assert.Equal(t, "Doanh nghiệp", domain.NewsTypeDoanhNghiep.Label())
	assert.Equal(t, "Phân tích", domain.NewsTypePhanTich.Label())
	assert.Equal(t, "Trái phiếu / Tín dụng", domain.NewsTypeTraiPhieuTinDung.Label())
	assert.Len(t, domain.AllNewsTypes(), 9, "PRD mục 5 chốt 9 nhãn")

	_, ok := domain.ParseNewsType("khong_ton_tai")
	assert.False(t, ok)
}
