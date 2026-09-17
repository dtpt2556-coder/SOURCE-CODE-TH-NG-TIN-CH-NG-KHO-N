package tagger_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/ingest/tagger"
)

// MinPrecision là ngưỡng PRD mục 10 cho việc gắn mã. Gắn sai mã tệ hơn bỏ sót:
// người dùng mất niềm tin ngay khi thấy một mã không liên quan trong cột Mã CK.
const MinPrecision = 0.96

type corpusArticle struct {
	ID       int      `json:"id"`
	Source   string   `json:"source"`
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Expected []string `json:"expected"`
	Note     string   `json:"note"`
}

type corpusFile struct {
	Description string          `json:"description"`
	Source      string          `json:"source"`
	Rules       []string        `json:"labelling_rules"`
	Articles    []corpusArticle `json:"articles"`
}

func loadCorpus(t *testing.T) corpusFile {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("../../../testdata", "tagger_corpus.json"))
	require.NoError(t, err, "đọc corpus gán nhãn")

	var c corpusFile
	require.NoError(t, json.Unmarshal(raw, &c))
	require.Len(t, c.Articles, 50, "test-plan mục 4.2 yêu cầu 50 bài gán nhãn")
	return c
}

// corpusTickers nạp ĐÚNG từ điển production, xuất từ các file migration seed
// (0003 + 0005 + 0006 + 0007) sang testdata/ticker_dictionary.json.
//
// Trước đây hàm này dựng từ điển viết tay và thiếu negative_aliases, nên cổng
// precision chạy trên một từ điển khác với thứ đang chạy thật — đó là lý do
// "FPT Retail (FRT)" không bị chặn trong test dù DB đã khai báo phủ định.
func corpusTickers(t *testing.T) []tagger.Ticker {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("../../../testdata", "ticker_dictionary.json"))
	require.NoError(t, err, "đọc từ điển mã")

	var rows []struct {
		Symbol          string   `json:"symbol"`
		Sector          string   `json:"sector"`
		InVN30          bool     `json:"in_vn30"`
		Aliases         []string `json:"aliases"`
		NegativeAliases []string `json:"negative_aliases"`
	}
	require.NoError(t, json.Unmarshal(raw, &rows))
	require.GreaterOrEqual(t, len(rows), 50, "từ điển seed phải có ít nhất 50 mã")

	out := make([]tagger.Ticker, 0, len(rows))
	for _, r := range rows {
		out = append(out, tagger.Ticker{
			Symbol:          r.Symbol,
			Sector:          r.Sector,
			InVN30:          r.InVN30,
			Aliases:         r.Aliases,
			NegativeAliases: r.NegativeAliases,
		})
	}
	return out
}

// TestTaggerPrecisionOnLabelledCorpus là cổng chất lượng của tagger.
//
// Không có tập gán nhãn thì không ai đo được precision — QA vòng trước đã phải
// ước lượng bằng tay trên mẫu tiện lợi và đúng khi từ chối bịa con số recall.
func TestTaggerPrecisionOnLabelledCorpus(t *testing.T) {
	corpus := loadCorpus(t)
	tg := tagger.New(corpusTickers(t))

	var truePositive, falsePositive, falseNegative int
	var worst []string

	for _, article := range corpus.Articles {
		expected := map[string]bool{}
		for _, sym := range article.Expected {
			expected[sym] = true
		}

		got := map[string]bool{}
		for _, m := range tg.Tag(tagger.Input{Title: article.Title, Body: article.Body}) {
			got[m.Symbol] = true
		}

		var extra, missing []string
		for sym := range got {
			if expected[sym] {
				truePositive++
			} else {
				falsePositive++
				extra = append(extra, sym)
			}
		}
		for sym := range expected {
			if !got[sym] {
				falseNegative++
				missing = append(missing, sym)
			}
		}
		sort.Strings(extra)
		sort.Strings(missing)
		if len(extra) > 0 || len(missing) > 0 {
			worst = append(worst, fmt.Sprintf("  #%d [%s] thừa=%v thiếu=%v — %s",
				article.ID, article.Source, extra, missing, truncate(article.Title, 58)))
		}
	}

	precision := ratio(truePositive, truePositive+falsePositive)
	recall := ratio(truePositive, truePositive+falseNegative)

	t.Logf("corpus %d bài · TP=%d FP=%d FN=%d", len(corpus.Articles), truePositive, falsePositive, falseNegative)
	t.Logf("PRECISION = %.1f%% (ngưỡng PRD %.0f%%)", precision*100, MinPrecision*100)
	t.Logf("RECALL    = %.1f%%", recall*100)
	for _, line := range worst {
		t.Log(line)
	}

	assert.GreaterOrEqual(t, precision, MinPrecision,
		"precision %.1f%% dưới ngưỡng PRD %.0f%% — gắn sai mã tệ hơn bỏ sót",
		precision*100, MinPrecision*100)
}

// Bài không nêu doanh nghiệp nào thì cột Mã CK phải để trống. Đây là nhóm sinh
// ra 29/31 false positive mà QA tìm thấy.
func TestNoTickerInventedOnArticlesWithoutCompanies(t *testing.T) {
	corpus := loadCorpus(t)
	tg := tagger.New(corpusTickers(t))

	for _, article := range corpus.Articles {
		if len(article.Expected) > 0 {
			continue
		}
		got := tg.Tag(tagger.Input{Title: article.Title, Body: article.Body})
		assert.Empty(t, symbols(got),
			"#%d %q không nêu doanh nghiệp nào nhưng bị gắn %v",
			article.ID, truncate(article.Title, 60), symbols(got))
	}
}

func ratio(num, den int) float64 {
	if den == 0 {
		return 1
	}
	return float64(num) / float64(den)
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}
