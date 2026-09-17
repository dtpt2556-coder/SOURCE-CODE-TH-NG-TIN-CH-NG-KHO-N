package domain

// Sàn giao dịch hợp lệ.
const (
	ExchangeHOSE  = "HOSE"
	ExchangeHNX   = "HNX"
	ExchangeUPCOM = "UPCOM"
)

// Ticker là một mã chứng khoán trong `ticker_master`.
type Ticker struct {
	Symbol      string
	CompanyName string
	ShortName   string
	Exchange    string
	Sector      string
	InVN30      bool

	// ArticleCount chỉ được điền ở các truy vấn tổng hợp.
	ArticleCount int
	// Aliases dùng cho tagger (tên có dấu).
	Aliases []TickerAlias
	// NegativeAliases là các chuỗi KHÔNG được map vào mã này (rule G-04).
	NegativeAliases []string
}

// TickerAlias là một cách gọi khác của doanh nghiệp.
// AliasNormalized đã bỏ dấu + lowercase để tagger tra cứu O(1).
type TickerAlias struct {
	ID              int64
	Symbol          string
	Alias           string
	AliasNormalized string
}
