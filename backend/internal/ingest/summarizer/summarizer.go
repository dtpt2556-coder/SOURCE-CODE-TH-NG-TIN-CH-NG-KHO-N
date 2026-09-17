// Package summarizer sinh bản tóm tắt tiếng Việt giàu số liệu cho digest.
//
// ADR-007: hệ thống PHẢI chạy được end-to-end mà không cần API key nào.
// `extractive` là provider mặc định; `anthropic` chỉ bật khi có ANTHROPIC_API_KEY.
package summarizer

import "context"

// Provider là hợp đồng của mọi bộ tóm tắt (mục 9 architecture.md).
type Provider interface {
	Name() string
	Summarize(ctx context.Context, in Input) (Output, error)
}

// Input là dữ liệu đầu vào cho một lần tóm tắt.
type Input struct {
	Title      string
	Body       string
	SourceName string
	Tickers    []string
}

// Output là kết quả tóm tắt.
type Output struct {
	SummaryMD string   // chỉ chứa **bold**, không HTML
	BoldSpans []string // các cụm được bold, để validator kiểm tra
}

// Giới hạn độ dài theo rubric C4 của BA2 §5.2.
const (
	TargetMinWords = 60
	TargetMaxWords = 150
)
