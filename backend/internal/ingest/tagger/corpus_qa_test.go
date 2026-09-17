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
	"github.com/tenpoint/tenpoint-api/internal/domain"
	"github.com/tenpoint/tenpoint-api/internal/ingest/tagger"
)

// Cổng nghiệm thu precision chạy trên corpus do QA gán nhãn ĐỘC LẬP.
//
// File testdata/tagger_corpus_qa.json là CHỈ ĐỌC với người viết tagger. Corpus
// cũ (tagger_corpus.json) do chính tôi gán nhãn nên là tập tuning: nó cho 100%
// và con số đó không có giá trị nghiệm thu.

type qaTag struct {
	Symbol    string `json:"symbol"`
	Relevance string `json:"relevance"`
	Reason    string `json:"reason"`
}

type qaSample struct {
	ID         int     `json:"id"`
	Domain     string  `json:"source_domain"`
	Title      string  `json:"title"`
	Body       string  `json:"body"`
	Expected   []qaTag `json:"expected"`
	MustNotTag []qaTag `json:"must_not_tag"`
	Uncertain  []qaTag `json:"uncertain"`
	Group      string  `json:"group"`
	Note       string  `json:"note"`
}

type qaCorpus struct {
	Meta struct {
		Counts struct {
			Samples           int   `json:"samples"`
			ExpectedTags      int   `json:"expected_tags"`
			MustNotTagTags    int   `json:"must_not_tag_tags"`
			UncertainExcluded int   `json:"uncertain_tags_excluded"`
			VerifiedEmptyIDs  []int `json:"verified_empty_expected_ids"`
		} `json:"counts"`
	} `json:"_meta"`
	Samples []qaSample `json:"samples"`
}

func loadQACorpus(t *testing.T) qaCorpus {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("../../../testdata", "tagger_corpus_qa.json"))
	require.NoError(t, err, "đọc corpus nghiệm thu của QA")

	var c qaCorpus
	require.NoError(t, json.Unmarshal(raw, &c))

	// Chốt lại kích thước corpus: nếu ai đó thêm/bớt mẫu để dễ pass thì test này
	// vỡ ngay.
	require.Equal(t, 27, len(c.Samples), "corpus QA phải có đúng 27 mẫu")
	require.Equal(t, 27, c.Meta.Counts.Samples)
	// 96, không phải 95: QA adjudicate thêm `#108 MWG` thành expected/mentioned ở
	// vòng nghiệm thu 3, cùng dạng "Tên công ty (MÃ)" với TCB/MBB ở câu liền trước.
	require.Equal(t, 96, c.Meta.Counts.ExpectedTags)
	require.Equal(t, 42, c.Meta.Counts.MustNotTagTags)
	require.Equal(t, 3, c.Meta.Counts.UncertainExcluded)
	return c
}

func symbolSet(tags []qaTag) map[string]string {
	out := make(map[string]string, len(tags))
	for _, t := range tags {
		out[t.Symbol] = t.Relevance
	}
	return out
}

// TestTaggerPrecisionOnQACorpus là cổng chặn phát hành của tagger.
func TestTaggerPrecisionOnQACorpus(t *testing.T) {
	corpus := loadQACorpus(t)
	tg := tagger.New(corpusTickers(t))

	var tp, fp, fn, skipped int
	byTier := map[string]struct{ tp, fp int }{}
	var misses []string

	for _, sample := range corpus.Samples {
		expected := symbolSet(sample.Expected)
		uncertain := symbolSet(sample.Uncertain)

		got := map[string]tagger.Match{}
		for _, m := range tg.Tag(tagger.Input{Title: sample.Title, Body: sample.Body}) {
			got[m.Symbol] = m
		}

		var extra, missing []string
		for symbol, match := range got {
			// 3 tag `uncertain` phải loại khỏi CẢ tử số lẫn mẫu số.
			if _, ok := uncertain[symbol]; ok {
				skipped++
				continue
			}
			tier := match.Relevance
			stat := byTier[tier]
			if _, ok := expected[symbol]; ok {
				tp++
				stat.tp++
			} else {
				fp++
				stat.fp++
				extra = append(extra, fmt.Sprintf("%s(%.2f)", symbol, match.Score))
			}
			byTier[tier] = stat
		}
		for symbol := range expected {
			if _, ok := got[symbol]; !ok {
				fn++
				missing = append(missing, symbol)
			}
		}

		sort.Strings(extra)
		sort.Strings(missing)
		if len(extra) > 0 {
			misses = append(misses, fmt.Sprintf("  #%d [%s] THỪA=%v — %s",
				sample.ID, sample.Domain, extra, truncate(sample.Title, 52)))
		}
	}

	precision := ratio(tp, tp+fp)
	recall := ratio(tp, tp+fn)

	t.Logf("corpus QA: %d mẫu · TP=%d FP=%d FN=%d · bỏ qua %d tag uncertain",
		len(corpus.Samples), tp, fp, fn, skipped)
	t.Logf("PRECISION = %.2f%%  (ngưỡng PRD > %.0f%%)", precision*100, MinPrecision*100)
	for _, tier := range []string{domain.RelevancePrimary, domain.RelevanceMentioned} {
		s := byTier[tier]
		t.Logf("  %-10s TP=%-3d FP=%-3d precision=%.2f%%", tier, s.tp, s.fp, ratio(s.tp, s.tp+s.fp)*100)
	}
	t.Logf("RECALL    = %.2f%%  (chỉ tham khảo — corpus chưa quét vét cạn mọi mã)", recall*100)
	for _, line := range misses {
		t.Log(line)
	}

	assert.Equal(t, 3, skipped, "phải loại đúng 3 tag uncertain khỏi phép tính")
	assert.Greater(t, precision, MinPrecision,
		"precision %.2f%% chưa vượt ngưỡng PRD %.0f%%", precision*100, MinPrecision*100)
}

// Bốn tag sai ngoài tầng T3 mà QA chỉ đích danh, mỗi cái một nguyên nhân gốc.
func TestQACitedFalsePositivesAreFixed(t *testing.T) {
	corpus := loadQACorpus(t)
	tg := tagger.New(corpusTickers(t))

	cases := map[int][]string{
		108: {"FPT", "GAS"}, // negative alias "FPT Retail" + va chạm bỏ dấu "Khí Việt Nam"
		129: {"GAS"},        // cùng lớp lỗi bỏ dấu
		168: {"VND"},        // G-02: USD/VND là tiền tệ, bài không có ngữ cảnh cổ phiếu nào
	}
	for _, sample := range corpus.Samples {
		banned, ok := cases[sample.ID]
		if !ok {
			continue
		}
		got := symbols(tg.Tag(tagger.Input{Title: sample.Title, Body: sample.Body}))
		for _, symbol := range banned {
			assert.NotContains(t, got, symbol,
				"bài %d không được gắn %s — nhận được %v", sample.ID, symbol, got)
		}
	}
}

// Tắt tầng T3 KHÔNG được làm mất các tag điểm cao đúng trên cùng bài.
// QA tự đính chính ở vòng 2: bài 191 nêu đích danh VCB/BID/TCB/MBB kèm số liệu.
func TestArticle191KeepsCorrectHighScoreTags(t *testing.T) {
	corpus := loadQACorpus(t)
	tg := tagger.New(corpusTickers(t))

	for _, sample := range corpus.Samples {
		if sample.ID != 191 {
			continue
		}
		got := symbols(tg.Tag(tagger.Input{Title: sample.Title, Body: sample.Body}))
		for _, symbol := range []string{"BID", "MBB", "TCB", "VCB"} {
			assert.Contains(t, got, symbol,
				"bài 191 phải giữ %s (nêu đích danh kèm số liệu) — nhận được %v", symbol, got)
		}
		for _, symbol := range []string{"GAS", "PLX", "PVD", "PVS"} {
			assert.NotContains(t, got, symbol, "4 mã dầu khí ở bài 191 là tag T3 sai")
		}
		return
	}
	t.Fatal("không tìm thấy bài 191 trong corpus QA")
}

// 10 bài QA đã xác minh là KHÔNG có mã nào. Bài 78 cố tình không nằm trong danh
// sách này vì tag duy nhất của nó bị xếp uncertain.
func TestVerifiedEmptySamplesGetNoTags(t *testing.T) {
	corpus := loadQACorpus(t)
	tg := tagger.New(corpusTickers(t))

	verified := map[int]bool{}
	for _, id := range corpus.Meta.Counts.VerifiedEmptyIDs {
		verified[id] = true
	}
	require.Len(t, verified, 10)

	for _, sample := range corpus.Samples {
		if !verified[sample.ID] {
			continue
		}
		got := symbols(tg.Tag(tagger.Input{Title: sample.Title, Body: sample.Body}))
		assert.Empty(t, got, "bài %d đã xác minh không có mã nhưng bị gắn %v", sample.ID, got)
	}
}
