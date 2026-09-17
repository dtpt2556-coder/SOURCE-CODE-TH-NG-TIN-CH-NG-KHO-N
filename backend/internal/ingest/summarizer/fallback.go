package summarizer

import (
	"context"
	"fmt"
	"log/slog"
)

// SummarizeWithFallback chạy provider chính rồi tự rơi về `extractive` khi LLM
// lỗi (timeout, 429, hết quota) hoặc trả nội dung vi phạm rule hình thức.
// Trả về tên provider đã thực sự sinh ra kết quả để ghi vào summary_provider.
//
// Pipeline không bao giờ được dừng vì LLM: mất một bản tóm tắt chất lượng cao
// còn hơn mất cả run (S-05, S-06).
func SummarizeWithFallback(ctx context.Context, primary, backup Provider, in Input, log *slog.Logger) (Output, string, error) {
	if log == nil {
		log = slog.Default()
	}

	out, err := primary.Summarize(ctx, in)
	if err != nil {
		log.Warn("primary summarizer failed, falling back to extractive",
			"provider", primary.Name(), "err", err)
		return runBackup(ctx, backup, in)
	}

	out = polish(out)
	issue := CheckQuality(out.SummaryMD)
	if issue == nil {
		return out, primary.Name(), nil
	}

	if primary.Name() == backup.Name() {
		return out, primary.Name(), fmt.Errorf("summary rejected: %s (%s)", issue.Code, issue.Detail)
	}

	log.Warn("primary summary rejected, retrying once",
		"provider", primary.Name(), "issue", issue.Code, "detail", issue.Detail)

	retried, retryErr := primary.Summarize(ctx, in)
	if retryErr == nil {
		retried = polish(retried)
		if CheckQuality(retried.SummaryMD) == nil {
			return retried, primary.Name(), nil
		}
	}
	return runBackup(ctx, backup, in)
}

func runBackup(ctx context.Context, backup Provider, in Input) (Output, string, error) {
	out, err := backup.Summarize(ctx, in)
	if err != nil {
		return Output{}, backup.Name(), fmt.Errorf("fallback summarizer %s: %w", backup.Name(), err)
	}
	out = polish(out)
	if issue := CheckQuality(out.SummaryMD); issue != nil {
		// Trả kèm nội dung để caller còn lưu lại chờ duyệt thay vì mất bài.
		return out, backup.Name(), fmt.Errorf("fallback summary rejected: %s (%s)", issue.Code, issue.Detail)
	}
	return out, backup.Name(), nil
}

func polish(out Output) Output {
	out.SummaryMD = SanitizeMarkdown(out.SummaryMD)
	out.BoldSpans = BoldSpans(out.SummaryMD)
	return out
}
