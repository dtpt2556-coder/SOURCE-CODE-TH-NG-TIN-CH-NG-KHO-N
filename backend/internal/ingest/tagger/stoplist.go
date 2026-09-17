package tagger

// DefaultStoplist — rule B1 của BA2: các từ viết tắt 3 ký tự KHÔNG BAO GIỜ được
// coi là mã chứng khoán. Thiếu danh sách này, mọi bài có "GDP", "USD", "ETF"...
// đều bị gắn mã sai hàng loạt.
//
// Bất biến quan trọng: danh sách này KHÔNG được chứa bất kỳ mã nào có trong
// bảng `tickers` (VND là ngoại lệ có chủ đích — xử lý bằng rule B2, không phải
// stoplist, vì VNDIRECT là mã thật).
func DefaultStoplist() []string {
	return []string{
		// Chỉ số & khái niệm kinh tế vĩ mô
		"GDP", "CPI", "PMI", "IIP", "FDI", "ODA", "GNI", "PPI",
		// Sản phẩm tài chính & chỉ tiêu
		"ETF", "IPO", "NAV", "OMO", "NIM", "ROE", "ROA", "EPS", "PEG",
		"YOY", "QOQ", "ESG",
		// Tiền tệ
		"USD", "EUR", "JPY", "CNY", "KRW", "GBP", "AUD", "SGD", "THB", "CHF",
		// Sở GDCK & cơ quan quản lý
		"HSX", "HNX", "SSC", "UBC", "BTC",
		// Chức danh
		"CEO", "CFO", "COO", "CTO", "CIO", "CMO", "CCO",
		// Lệnh & phiên giao dịch
		"ATC", "ATO", "LOC", "MTL", "MOK", "MAK", "PLO",
		// Viết tắt hành chính / kỹ thuật thường gặp trên báo VN
		"KCN", "KDT", "BDS", "CNT", "CNC", "XNK", "GTV",
		"BOT", "BTO", "PPP", "ODM", "OEM", "SME", "DNN",
		// Tổ chức quốc tế & hiệp định
		"IMF", "WTO", "WHO", "ILO", "ADB", "EIB", "ECB", "FED",
		// Đơn vị đo
		"KWH", "MWH", "TCO", "KGS",
		// Công ty chứng khoán hay xuất hiện ở vai trò nguồn trích dẫn nhưng
		// không nằm trong rổ theo dõi (rule B3 xử lý phần còn lại).
		"HCM", "KIS", "VPS",
		// Từ tiếng Anh thông dụng viết hoa trong tiêu đề
		"AND", "FOR", "THE", "NEW", "TOP", "ALL", "OUT",
	}
}

// DefaultSectorKeywords — tầng T3: từ khoá ngành -> giá trị `tickers.sector`.
// Từ khoá viết CÓ DẤU vì được so khớp bằng textutil.Fold.
func DefaultSectorKeywords() map[string]string {
	return map[string]string{
		"ngành thép":          "Thép",
		"thép":                "Thép",
		"tôn mạ":              "Thép",
		"gang thép":           "Thép",
		"ngành ngân hàng":     "Ngân hàng",
		"hệ thống ngân hàng":  "Ngân hàng",
		"nhóm ngân hàng":      "Ngân hàng",
		"ngành chăn nuôi":     "Chăn nuôi",
		"chăn nuôi":           "Chăn nuôi",
		"heo hơi":             "Chăn nuôi",
		"thịt heo":            "Chăn nuôi",
		"ngành dầu khí":       "Dầu khí",
		"nhóm dầu khí":        "Dầu khí",
		"ngành bất động sản":  "Bất động sản",
		"nhóm bất động sản":   "Bất động sản",
		"ngành công nghệ":     "Công nghệ",
		"chuyển đổi số":       "Công nghệ",
		"ngành bán lẻ":        "Bán lẻ",
		"ngành điện":          "Điện",
		"ngành hoá chất":      "Hóa chất",
		"ngành hóa chất":      "Hóa chất",
		"ngành hàng không":    "Hàng không",
		"ngành bảo hiểm":      "Bảo hiểm",
		"ngành thực phẩm":     "Thực phẩm & đồ uống",
		"bđs khu công nghiệp": "BĐS khu công nghiệp",
	}
}
