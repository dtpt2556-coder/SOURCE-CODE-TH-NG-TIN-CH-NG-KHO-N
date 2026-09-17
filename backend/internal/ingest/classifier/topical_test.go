package classifier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tenpoint/tenpoint-api/internal/ingest/classifier"
)

// Ba tin đầu trang chủ mà QA chụp được ở vòng 2 — đều gắn nhãn "Thị trường",
// không mã nào, và không liên quan gì tới chứng khoán.
func TestRejectsOffTopicArticlesFoundByQA(t *testing.T) {
	cases := []classifier.TopicalInput{
		{
			Title: "Trung Quốc mở kênh đào dài 134km ở Quảng Tây, nối sông với Biển Đông",
			Body: "Dự án kênh đào Bình Lục dài 134,2 km, tổng vốn 72,7 tỷ nhân dân tệ, " +
				"rút ngắn quãng đường ra biển khoảng 560 km cho tàu hàng nội địa.",
		},
		{
			Title: "Lên men cà phê: Ở đâu, khi nào và vì sao?",
			Body: "Quá trình lên men giúp loại bỏ lớp nhầy quanh hạt cà phê. Nhiệt độ lên men " +
				"lý tưởng dao động 18-24 độ C trong 12 đến 36 giờ tuỳ vùng trồng.",
		},
		{
			Title: "Nhà xuất bản Trẻ và những trang sách tuổi thơ",
			Body: "Suốt 45 năm, nhà xuất bản đã in hàng nghìn đầu sách cho thiếu nhi, " +
				"trở thành một phần ký ức của nhiều thế hệ bạn đọc phương Nam.",
		},
	}
	for _, in := range cases {
		got := classifier.IsOnTopic(in)
		assert.False(t, got.OnTopic, "%q phải bị loại (tín hiệu=%d)", in.Title, got.Signals)
		assert.Equal(t, "off_topic", got.Reason)
	}
}

func TestKeepsArticlesWithTickers(t *testing.T) {
	got := classifier.IsOnTopic(classifier.TopicalInput{
		Title:       "Chuyện gì đang xảy ra với doanh nghiệp này?",
		Body:        "Nội dung không có từ khoá tài chính nào cả, chỉ kể chuyện nội bộ công ty.",
		TickerCount: 2,
	})
	assert.True(t, got.OnTopic, "bài đã gắn được mã thì luôn đúng chủ đề")
}

func TestKeepsMacroFinanceArticlesWithoutTickers(t *testing.T) {
	cases := []classifier.TopicalInput{
		{
			Title: "Ngân hàng Nhà nước nâng tỷ giá trung tâm lên kỷ lục",
			Body:  "Tỷ giá tham chiếu USD/VND tăng lên 25.617 đồng trước cuộc họp của Fed về lãi suất.",
		},
		{
			Title: "Khối ngoại mua ròng phiên thứ năm liên tiếp",
			Body:  "Giá trị mua ròng đạt 260 tỷ đồng trên HOSE, tập trung ở nhóm cổ phiếu ngân hàng.",
		},
		{
			Title: "Doanh nghiệp bất động sản phát hành 5.000 tỷ đồng trái phiếu",
			Body:  "Lô trái phiếu doanh nghiệp kỳ hạn 3 năm, lãi suất 11%/năm, dư nợ tăng mạnh.",
		},
	}
	for _, in := range cases {
		got := classifier.IsOnTopic(in)
		assert.True(t, got.OnTopic, "%q phải được giữ (tín hiệu=%d)", in.Title, got.Signals)
	}
}

// Nhắc thoáng qua một từ tài chính không đủ để lên digest.
func TestSingleFinanceMentionIsNotEnough(t *testing.T) {
	got := classifier.IsOnTopic(classifier.TopicalInput{
		Title: "Làng nghề gốm Bát Tràng thích ứng với du lịch",
		Body: "Nhiều hộ sản xuất vay tín dụng để mở rộng xưởng, nhưng khó khăn lớn nhất " +
			"vẫn là tìm đầu ra cho sản phẩm thủ công truyền thống.",
	})
	assert.False(t, got.OnTopic, "tín hiệu=%d", got.Signals)
}
