package domain

import "strings"

// NewsType là taxonomy "Loại tin" — tập đóng 9 nhãn theo PRD mục 5.
// architecture.md mục 6 chỉ liệt kê 7; hai nhãn phan_tich và trai_phieu_tin_dung
// được bổ sung để khớp PRD (migration 0005).
type NewsType string

const (
	NewsTypeViMo             NewsType = "vi_mo"
	NewsTypeNganh            NewsType = "nganh"
	NewsTypeDoanhNghiep      NewsType = "doanh_nghiep"
	NewsTypeThiTruong        NewsType = "thi_truong"
	NewsTypeKhoiNgoai        NewsType = "khoi_ngoai"
	NewsTypeCoTucPhatHanh    NewsType = "co_tuc_phat_hanh"
	NewsTypePhapLy           NewsType = "phap_ly"
	NewsTypePhanTich         NewsType = "phan_tich"
	NewsTypeTraiPhieuTinDung NewsType = "trai_phieu_tin_dung"
)

// newsTypePriority là thứ tự ưu tiên khi một tin thuộc nhiều nhóm (PRD mục 5):
// nhãn càng cụ thể và hành động được thì càng đứng trước. Nhà đầu tư lọc
// "Cổ tức/Phát hành" để biết ngày chốt quyền — gán nhầm thành "Doanh nghiệp" là
// họ bỏ lỡ.
var newsTypePriority = []NewsType{
	NewsTypePhapLy,
	NewsTypeCoTucPhatHanh,
	NewsTypeKhoiNgoai,
	NewsTypeTraiPhieuTinDung,
	NewsTypePhanTich,
	NewsTypeDoanhNghiep,
	NewsTypeNganh,
	NewsTypeThiTruong,
	NewsTypeViMo,
}

// displayOrder là thứ tự hiển thị trên dropdown lọc của FE — theo tầng tác động
// từ rộng tới hẹp, không phải theo độ ưu tiên gán nhãn.
var displayOrder = []NewsType{
	NewsTypeViMo,
	NewsTypeNganh,
	NewsTypeDoanhNghiep,
	NewsTypeThiTruong,
	NewsTypeKhoiNgoai,
	NewsTypeCoTucPhatHanh,
	NewsTypePhapLy,
	NewsTypePhanTich,
	NewsTypeTraiPhieuTinDung,
}

var newsTypeLabels = map[NewsType]string{
	NewsTypeViMo:             "Vĩ mô",
	NewsTypeNganh:            "Ngành",
	NewsTypeDoanhNghiep:      "Doanh nghiệp",
	NewsTypeThiTruong:        "Thị trường",
	NewsTypeKhoiNgoai:        "Khối ngoại",
	NewsTypeCoTucPhatHanh:    "Cổ tức / Phát hành",
	NewsTypePhapLy:           "Pháp lý",
	NewsTypePhanTich:         "Phân tích",
	NewsTypeTraiPhieuTinDung: "Trái phiếu / Tín dụng",
}

// AllNewsTypes trả về danh sách loại tin theo thứ tự hiển thị trên FE.
func AllNewsTypes() []NewsType {
	out := make([]NewsType, len(displayOrder))
	copy(out, displayOrder)
	return out
}

// NewsTypePriority trả về thứ tự ưu tiên gán nhãn (PRD mục 5).
func NewsTypePriority() []NewsType {
	out := make([]NewsType, len(newsTypePriority))
	copy(out, newsTypePriority)
	return out
}

// Label trả về nhãn tiếng Việt hiển thị trên cột "Loại tin".
func (n NewsType) Label() string {
	if l, ok := newsTypeLabels[n]; ok {
		return l
	}
	return string(n)
}

// Valid cho biết giá trị có nằm trong tập đóng hay không.
func (n NewsType) Valid() bool {
	_, ok := newsTypeLabels[n]
	return ok
}

// ParseNewsType chuẩn hoá chuỗi đầu vào thành NewsType.
func ParseNewsType(s string) (NewsType, bool) {
	nt := NewsType(strings.ToLower(strings.TrimSpace(s)))
	return nt, nt.Valid()
}

// NewsTypeValues liệt kê các giá trị hợp lệ để đưa vào thông điệp lỗi 400 cho
// người dùng, thay vì im lặng trả danh sách rỗng.
func NewsTypeValues() string {
	vals := make([]string, 0, len(displayOrder))
	for _, nt := range displayOrder {
		vals = append(vals, string(nt))
	}
	return strings.Join(vals, ", ")
}
