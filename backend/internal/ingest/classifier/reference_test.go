package classifier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tenpoint/tenpoint-api/internal/domain"
	"github.com/tenpoint/tenpoint-api/internal/ingest/classifier"
)

// Ảnh tham chiếu gán cả 5/5 dòng là "Ngành" — BA2 §3.4 chỉ ra 3 trong số đó bị
// gán quá thô. Seed demo vẫn giữ nhãn "Ngành" để đối chiếu thị giác 1:1 với
// ảnh; test này chạy classifier trên chính nội dung đó để chứng minh sản phẩm
// phân loại đúng hơn bản mẫu.
func TestClassifierBeatsReferenceMockup(t *testing.T) {
	c := classifier.New()

	cases := []struct {
		name    string
		in      classifier.Input
		mockup  domain.NewsType
		want    domain.NewsType
		because string
	}{
		{
			name: "12 ngân hàng cam kết 408.000 tỷ tín dụng",
			in: classifier.Input{
				Title: "12 ngân hàng cam kết 408.000 tỷ đồng tín dụng cho DNNVV",
				Body: "12 ngân hàng cam kết 408.000 tỷ đồng tín dụng cho DNNVV theo chương trình NHNN chủ trì. " +
					"Big4 đăng ký 220.000 tỷ (Agribank 70.000; BIDV, Vietcombank, VietinBank mỗi bên 50.000 tỷ), " +
					"lãi suất thấp hơn ít nhất 1 điểm % so bình quân cùng kỳ hạn. 8 ngân hàng tư nhân đăng ký " +
					"188.000 tỷ, giảm lãi 0,5-2 điểm % tuỳ lĩnh vực kèm miễn giảm phí. Cuối tháng 7 chỉ 8,8% " +
					"DNNVV tiếp cận được vốn so với trên 47% ở doanh nghiệp lớn.",
				PrimarySymbols: []string{"BID", "VCB", "CTG", "SHB", "MSB", "STB", "BVB", "NAB", "NVB", "SGB", "TPB"},
			},
			mockup:  domain.NewsTypeNganh,
			want:    domain.NewsTypeTraiPhieuTinDung,
			because: "nội dung chính là gói tín dụng và lãi suất, không phải đặc thù ngành",
		},
		{
			name: "VN-Index đóng cửa 1.821 điểm",
			in: classifier.Input{
				Title: "VN-Index đóng cửa 1.821 điểm, lần đầu vượt mốc 1.800",
				Body: "VN-Index đóng cửa 1.821 điểm (+29,91 điểm, +1,67%), lần đầu vượt 1.800, phiên tăng thứ 5 " +
					"liên tiếp. Giá trị giao dịch ~20.000 tỷ, vượt bình quân 20 phiên; 182 mã tăng / 125 giảm. " +
					"TCB trần lên 33.450 đ khớp 39,17 triệu cp, BCM trần 44.450 đ, VIC +4,31%, FPT +2,69%.",
				PrimarySymbols: []string{"TCB", "BCM", "VIC", "FPT"},
			},
			mockup:  domain.NewsTypeNganh,
			want:    domain.NewsTypeThiTruong,
			because: "chủ thể của tin là chỉ số, các mã chỉ là minh hoạ",
		},
		{
			name: "ACBS ước 5.588 tỷ vào kỳ cơ cấu FTSE GEIS",
			in: classifier.Input{
				Title: "ACBS: 5.588 tỷ đồng có thể chảy vào 117 cổ phiếu Việt Nam kỳ cơ cấu FTSE GEIS",
				Body: "ACBS ước 5.588 tỷ đồng (216 triệu USD) chảy vào 117 cổ phiếu Việt Nam kỳ cơ cấu FTSE GEIS " +
					"hiệu lực 21/9/2026, từ 17 quỹ thụ động; phân bổ 3 large-cap, 3 mid-cap, 21 small-cap, 90 micro-cap. " +
					"Ba mã hút mạnh nhất: VIC 80,67 triệu USD, VHM 25,43 triệu USD, HPG 13,11 triệu USD.",
				PrimarySymbols: []string{"VIC", "VHM", "HPG"},
			},
			mockup:  domain.NewsTypeNganh,
			want:    domain.NewsTypeKhoiNgoai,
			because: "chủ thể là quỹ ngoại và tổ chức xếp hạng thị trường",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, confidence := c.Classify(tc.in)
			assert.Equal(t, tc.want, got, "%s", tc.because)
			assert.NotEqual(t, tc.mockup, got, "phải khác nhãn thô của bản mẫu")
			assert.GreaterOrEqual(t, confidence, 0.7)
		})
	}
}

// Hai dòng còn lại của ảnh thì bản mẫu gán đúng — classifier phải giữ nguyên.
func TestClassifierAgreesWithReferenceWhereMockupWasRight(t *testing.T) {
	c := classifier.New()

	cbam, _ := c.Classify(classifier.Input{
		Title: "CBAM: Doanh nghiệp thép, nhôm Việt đối mặt chi phí carbon lớn khi xuất sang EU",
		Body: "CBAM của EU áp dụng từ đầu 2026, doanh nghiệp có nguy cơ bù đắp hơn 100 USD/tấn thép xuất " +
			"sang EU, nhôm tới cả nghìn USD/tấn. Cường độ phát thải thép Việt 2,5 tấn CO2/tấn.",
		PrimarySymbols: []string{"HPG", "HSG", "NKG", "TIS"},
	})
	assert.Equal(t, domain.NewsTypeNganh, cbam)

	meat, _ := c.Classify(classifier.Input{
		Title: "Sáu tháng đầu 2026, Việt Nam nhập 494.000 tấn thịt trị giá 1,5 tỷ USD",
		Body: "Sáu tháng đầu 2026 nhập 494.000 tấn thịt, ~1,5 tỷ USD. Gia cầm 186.400 tấn, thịt trâu " +
			"108.000 tấn. Giá heo hơi trong nước 56.000 đ/kg chịu sức ép từ hàng đông lạnh giá rẻ.",
		PrimarySymbols: []string{"DBC", "BAF", "MML", "HAG"},
	})
	assert.Equal(t, domain.NewsTypeNganh, meat)
}

func TestClassifyAnalysisReport(t *testing.T) {
	got, _ := classifier.New().Classify(classifier.Input{
		Title: "MBS Research dự phóng thị trường Data Center Việt Nam",
		Body:  "Báo cáo phân tích của MBS Research đưa giá mục tiêu và dự phóng tăng trưởng 25% mỗi năm.",
	})
	assert.Equal(t, domain.NewsTypePhanTich, got)
}

func TestClassifyCorporateBond(t *testing.T) {
	got, _ := classifier.New().Classify(classifier.Input{
		Title: "Doanh nghiệp bất động sản phát hành 5.000 tỷ đồng trái phiếu",
		Body:  "Lô trái phiếu doanh nghiệp kỳ hạn 3 năm, lãi suất 11%/năm.",
	})
	// "phát hành riêng lẻ" không xuất hiện nên không rơi vào Cổ tức/Phát hành.
	assert.Equal(t, domain.NewsTypeTraiPhieuTinDung, got)
}

func TestPriorityOrderMatchesPRD(t *testing.T) {
	assert.Equal(t, []domain.NewsType{
		domain.NewsTypePhapLy,
		domain.NewsTypeCoTucPhatHanh,
		domain.NewsTypeKhoiNgoai,
		domain.NewsTypeTraiPhieuTinDung,
		domain.NewsTypePhanTich,
		domain.NewsTypeDoanhNghiep,
		domain.NewsTypeNganh,
		domain.NewsTypeThiTruong,
		domain.NewsTypeViMo,
	}, domain.NewsTypePriority())
}
