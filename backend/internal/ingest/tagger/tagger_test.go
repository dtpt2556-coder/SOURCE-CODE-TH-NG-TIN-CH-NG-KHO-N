package tagger_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/domain"
	"github.com/tenpoint/tenpoint-api/internal/ingest/tagger"
)

// testTickers là rổ mã rút gọn dùng cho toàn bộ test — offline, không cần DB.
func testTickers() []tagger.Ticker {
	return []tagger.Ticker{
		{Symbol: "FPT", Sector: "Công nghệ", Aliases: []string{"Tập đoàn FPT"}},
		{Symbol: "CMG", Sector: "Công nghệ", Aliases: []string{"CMC", "Tập đoàn CMC"}},
		{Symbol: "HPG", Sector: "Thép", Aliases: []string{"Hòa Phát", "Hoà Phát", "Tập đoàn Hòa Phát"}},
		{Symbol: "HSG", Sector: "Thép", Aliases: []string{"Hoa Sen", "Tôn Hoa Sen"}},
		{Symbol: "NKG", Sector: "Thép", Aliases: []string{"Nam Kim", "Thép Nam Kim"}},
		{Symbol: "TIS", Sector: "Thép", Aliases: []string{"Gang thép Thái Nguyên", "Tisco"}},
		{Symbol: "VCB", Sector: "Ngân hàng", Aliases: []string{"Vietcombank", "Ngân hàng Ngoại Thương"}},
		{Symbol: "TCB", Sector: "Ngân hàng", Aliases: []string{"Techcombank"}},
		{Symbol: "ACB", Sector: "Ngân hàng", Aliases: []string{"Ngân hàng Á Châu"}},
		{Symbol: "BID", Sector: "Ngân hàng", Aliases: []string{"BIDV"}},
		{Symbol: "CTG", Sector: "Ngân hàng", Aliases: []string{"VietinBank"}},
		{Symbol: "STB", Sector: "Ngân hàng", Aliases: []string{"Sacombank"}},
		{Symbol: "SHB", Sector: "Ngân hàng", Aliases: []string{"Ngân hàng Sài Gòn Hà Nội"}},
		{Symbol: "MSB", Sector: "Ngân hàng", Aliases: []string{"Maritime Bank"}},
		{Symbol: "TPB", Sector: "Ngân hàng", Aliases: []string{"TPBank"}},
		{Symbol: "VIC", Sector: "Bất động sản", Aliases: []string{"Vingroup", "Tập đoàn Vingroup"}},
		{Symbol: "VHM", Sector: "Bất động sản", Aliases: []string{"Vinhomes"}},
		{Symbol: "VND", Sector: "Dịch vụ tài chính", Aliases: []string{"VNDIRECT"}},
		{Symbol: "SSI", Sector: "Dịch vụ tài chính"},
		{Symbol: "MBS", Sector: "Dịch vụ tài chính"},
		{Symbol: "DBC", Sector: "Chăn nuôi", Aliases: []string{"Dabaco"}},
		{Symbol: "BAF", Sector: "Chăn nuôi", Aliases: []string{"Nông nghiệp BAF"}},
		{Symbol: "MML", Sector: "Chăn nuôi", Aliases: []string{"Masan MEATLife"}},
		{Symbol: "HAG", Sector: "Chăn nuôi", Aliases: []string{"Hoàng Anh Gia Lai", "HAGL"}},
	}
}

func symbols(ms []tagger.Match) []string {
	out := make([]string, 0, len(ms))
	for _, m := range ms {
		out = append(out, m.Symbol)
	}
	return out
}

func find(ms []tagger.Match, symbol string) (tagger.Match, bool) {
	for _, m := range ms {
		if m.Symbol == symbol {
			return m, true
		}
	}
	return tagger.Match{}, false
}

// --------------------------------------------------------------------- T1

func TestTagsTickerInTitleAsPrimary(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "FPT báo lãi quý 2 tăng 20%",
		Body:  "Cổ phiếu FPT tăng 2,69% trong phiên hôm nay. FPT ghi nhận doanh thu 15.000 tỷ đồng.",
	})

	m, ok := find(got, "FPT")
	require.True(t, ok, "phải gắn mã FPT, nhận được %v", symbols(got))
	assert.Equal(t, domain.RelevancePrimary, m.Relevance)
	assert.GreaterOrEqual(t, m.Score, 0.8)
}

func TestTagsTickerWithPriceMoveContext(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Phiên 26/8: VN-Index tăng mạnh",
		Body:  "TCB trần lên 33.450 đ khớp 39,17 triệu cp, VIC +4,31%, FPT +2,69%.",
	})

	assert.Contains(t, symbols(got), "VIC")
	assert.Contains(t, symbols(got), "TCB")
	assert.Contains(t, symbols(got), "FPT")
}

// --------------------------------------------------------------- Stoplist B1

func TestStoplistBlocksGDPAndUSD(t *testing.T) {
	// Cố tình nhét GDP/USD/ETF vào rổ mã để chứng minh stoplist chặn được, kể cả
	// khi bảng tickers bị nhiễm dữ liệu bẩn.
	dirty := append(testTickers(),
		tagger.Ticker{Symbol: "GDP", Sector: "Vĩ mô"},
		tagger.Ticker{Symbol: "USD", Sector: "Tiền tệ"},
		tagger.Ticker{Symbol: "ETF", Sector: "Quỹ"},
		tagger.Ticker{Symbol: "CPI", Sector: "Vĩ mô"},
		tagger.Ticker{Symbol: "ATC", Sector: "Giao dịch"},
	)
	tg := tagger.New(dirty)
	got := tg.Tag(tagger.Input{
		Title: "GDP quý III tăng 7,1%, CPI tăng 3,2%",
		Body:  "Tỷ giá USD/VND vượt 26.000. Quỹ ETF hút vốn trong phiên ATC, cổ phiếu tăng giá.",
	})

	for _, banned := range []string{"GDP", "CPI", "USD", "ETF", "ATC"} {
		assert.NotContains(t, symbols(got), banned, "stoplist phải chặn %s", banned)
	}
}

func TestStoplistContainsRequiredAbbreviations(t *testing.T) {
	set := map[string]bool{}
	for _, s := range tagger.DefaultStoplist() {
		set[s] = true
	}
	for _, required := range []string{"GDP", "CPI", "ETF", "HSX", "HNX", "USD", "CEO", "ATC"} {
		assert.True(t, set[required], "stoplist phải chứa %s", required)
	}
	// Bất biến: stoplist KHÔNG được chứa mã thật.
	for _, real := range []string{"FPT", "HPG", "VIC", "VCB", "TCB", "SSI", "VND", "CTG"} {
		assert.False(t, set[real], "stoplist không được chứa mã thật %s", real)
	}
}

func TestHoseAndHnxNeverTagged(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "HOSE và HNX công bố dữ liệu giao dịch",
		Body:  "Thanh khoản trên HOSE đạt 20.000 tỷ đồng, HNX đạt 1.500 tỷ đồng, cổ phiếu tăng giá.",
	})
	assert.Empty(t, symbols(got), "không được gắn mã nào, nhận được %v", symbols(got))
}

// ------------------------------------------------------------------- Rule B2

func TestVNDAsCurrencyIsNotTagged(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Doanh nghiệp huy động 5.000 tỷ VND trái phiếu",
		Body:  "Tổng giá trị phát hành đạt 5.000 tỷ VND, lãi suất 9,5%/năm.",
	})
	assert.NotContains(t, symbols(got), "VND", "‘tỷ VND’ là đơn vị tiền, không phải mã")
}

func TestVNDAsTickerIsTagged(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Cổ phiếu VND tăng trần phiên 26/8",
		Body:  "Mã VND khớp lệnh 10 triệu cp, thị giá lên 18.500 đồng.",
	})
	m, ok := find(got, "VND")
	require.True(t, ok, "phải gắn mã VND, nhận được %v", symbols(got))
	assert.Equal(t, domain.RelevancePrimary, m.Relevance)
}

func TestVNDTaggedWhenVNDirectAppears(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "VNDIRECT công bố kết quả kinh doanh",
		Body:  "Công ty VNDIRECT ghi nhận doanh thu 1.200 tỷ đồng. Mã VND tăng 3%.",
	})
	assert.Contains(t, symbols(got), "VND")
}

// ------------------------------------------------------------------- Rule B3

func TestAttributionMBSResearchIsNotTagged(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Thị trường Data Center Việt Nam còn nhiều dư địa",
		Body:  "Theo MBS Research, quy mô thị trường đạt 12% tổng công suất khu vực. MBS Research dự phóng tăng trưởng 25% mỗi năm.",
	})
	assert.NotContains(t, symbols(got), "MBS", "mã trong cụm trích dẫn nguồn không được gắn")
}

func TestAttributionChungKhoanSSIIsNotTagged(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Dòng tiền vào nhóm ngân hàng",
		Body:  "Theo Chứng khoán SSI, dòng tiền vào nhóm ngân hàng đạt 2.000 tỷ đồng trong phiên. SSI ước tính xu hướng còn tiếp diễn.",
	})
	assert.NotContains(t, symbols(got), "SSI")
}

func TestACBSDoesNotLeakIntoACB(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Dòng vốn ngoại vào cổ phiếu Việt Nam",
		Body:  "ACBS ước tính 5.588 tỷ đồng chảy vào 117 cổ phiếu Việt Nam kỳ cơ cấu FTSE GEIS.",
	})
	assert.NotContains(t, symbols(got), "ACB", "‘ACBS’ không được hiểu thành mã ACB")
}

// ------------------------------------------------------------------- Rule B4

func TestTickerWithoutStockContextIsDropped(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Tỉnh công bố kế hoạch VIC-2030 giai đoạn mới",
		Body:  "Kế hoạch VIC-2030 của tỉnh tập trung vào giao thông nông thôn.",
	})
	assert.NotContains(t, symbols(got), "VIC", "thiếu ngữ cảnh chứng khoán thì không gắn mã")
}

// ------------------------------------------------------------- T2 — tên công ty

func TestCompanyNameWithDiacriticsMapsToSymbol(t *testing.T) {
	tg := tagger.New(testTickers())

	// Hai cách gõ dấu khác nhau của cùng một tên.
	for _, title := range []string{
		"Hoà Phát khởi công dự án Dung Quất 3",
		"Hòa Phát khởi công dự án Dung Quất 3",
	} {
		got := tg.Tag(tagger.Input{
			Title: title,
			Body:  "Tập đoàn Hòa Phát cho biết tổng vốn đầu tư đạt 85.000 tỷ đồng.",
		})
		m, ok := find(got, "HPG")
		require.True(t, ok, "tiêu đề %q phải gắn HPG, nhận được %v", title, symbols(got))
		assert.Equal(t, domain.RelevancePrimary, m.Relevance)
	}
}

func TestBrandNameMapsToSymbol(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Vietcombank giảm lãi suất cho vay",
		Body:  "Ngân hàng Vietcombank đăng ký 50.000 tỷ đồng tín dụng ưu đãi cho doanh nghiệp.",
	})
	m, ok := find(got, "VCB")
	require.True(t, ok, "phải gắn VCB, nhận được %v", symbols(got))
	assert.Equal(t, domain.RelevancePrimary, m.Relevance)
}

// ---------------------------------------------------- primary vs mentioned

func TestPrimaryVsMentioned(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "VN-Index vượt 1.800 điểm",
		Body: "Chỉ số VN-Index đóng cửa tại 1.821 điểm, tăng 1,67%.\n\n" +
			"Cổ phiếu TCB tăng trần lên 33.450 đồng, khớp 39,17 triệu cp. TCB dẫn dắt nhóm ngân hàng.\n\n" +
			"Khối ngoại bán ròng ACB 78 tỷ đồng.",
	})

	tcb, ok := find(got, "TCB")
	require.True(t, ok, "phải gắn TCB, nhận được %v", symbols(got))
	assert.Equal(t, domain.RelevancePrimary, tcb.Relevance, "TCB xuất hiện 2 lần -> primary")

	acb, ok := find(got, "ACB")
	require.True(t, ok, "phải gắn ACB, nhận được %v", symbols(got))
	assert.Equal(t, domain.RelevanceMentioned, acb.Relevance, "ACB chỉ nhắc 1 lần cuối bài -> mentioned")
}

// ------------------------------------------------------------------- T3 ngành

func TestSectorContextTagsSteelGroup(t *testing.T) {
	tg := tagger.NewWithIndustryTier(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Ngành thép Việt Nam đối mặt CBAM của EU",
		Body:  "Cơ chế điều chỉnh carbon của EU áp lên ngành thép từ đầu 2026. Chi phí bù đắp có thể hơn 100 USD/tấn.",
	})

	for _, sym := range []string{"HPG", "HSG", "NKG", "TIS"} {
		m, ok := find(got, sym)
		require.True(t, ok, "tin ngành thép phải gắn %s, nhận được %v", sym, symbols(got))
		assert.Equal(t, domain.RelevanceMentioned, m.Relevance, "T3 chỉ được gán mentioned")
	}
}

// Thân bài lấy từ tin nhập khẩu thịt trong ảnh tham chiếu: ngành chăn nuôi được
// nhắc nhiều lần nên T3 được phép suy luận.
func TestSectorContextTagsLivestockGroup(t *testing.T) {
	tg := tagger.NewWithIndustryTier(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Sáu tháng đầu 2026, Việt Nam nhập 494.000 tấn thịt trị giá 1,5 tỷ USD",
		Body: "Sáu tháng đầu 2026 nhập 494.000 tấn thịt, khoảng 1,5 tỷ USD. Gia cầm 186.400 tấn, " +
			"thịt trâu 108.000 tấn, thịt heo khoảng 80.000 tấn.\n\n" +
			"Quý II nhập 42.600 tấn thịt heo, giá bình quân 2.152 USD/tấn giảm 18,5%.\n\n" +
			"Giá heo hơi trong nước 56.000 đồng/kg chịu sức ép từ hàng đông lạnh giá rẻ, " +
			"ảnh hưởng trực tiếp tới biên lợi nhuận ngành chăn nuôi trong nước.",
	})
	for _, sym := range []string{"DBC", "BAF", "MML", "HAG"} {
		assert.Contains(t, symbols(got), sym)
	}
	for _, m := range got {
		assert.Equal(t, domain.RelevanceMentioned, m.Relevance, "T3 chỉ được gán mentioned")
	}
}

// G-08: từ khoá ngành nhắc thoáng qua KHÔNG đủ để suy luận cả rổ mã.
// Đây là nguồn của 29/31 false positive QA tìm thấy.
func TestSectorContextIgnoresPassingMention(t *testing.T) {
	tg := tagger.NewWithIndustryTier(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Trung Quốc mở kênh đào dài 134km ở Quảng Tây, nối sông với Biển Đông",
		Body: "Dự án có tổng vốn 72,7 tỷ nhân dân tệ, được tài trợ bởi một nhóm ngân hàng " +
			"quốc doanh. Tuyến đường thuỷ rút ngắn quãng đường ra biển khoảng 560 km.",
	})
	assert.Empty(t, symbols(got),
		"một lần nhắc 'ngân hàng' không được kéo theo cả rổ mã, nhận được %v", symbols(got))
}

// G-08: bài đã nêu doanh nghiệp cụ thể thì không suy luận thêm theo ngành.
func TestSectorContextSkippedWhenCompanyIsSubject(t *testing.T) {
	tg := tagger.NewWithIndustryTier(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Hòa Phát khởi công dự án Dung Quất 3",
		Body: "Tập đoàn Hòa Phát cho biết tổng vốn đầu tư đạt 85.000 tỷ đồng. " +
			"Dự án nằm trong chiến lược mở rộng công suất ngành thép của doanh nghiệp, " +
			"hướng tới sản lượng thép chất lượng cao phục vụ ngành thép trong nước.",
	})
	assert.Equal(t, []string{"HPG"}, symbols(got),
		"đã có chủ thể HPG thì không kéo thêm HSG/NKG/TIS theo ngành")
}

// Trùng ký hiệu: VNM trong "VNM ETF" là quỹ VanEck, không phải Vinamilk.
func TestETFSuffixIsNotATicker(t *testing.T) {
	tg := tagger.New(append(testTickers(),
		tagger.Ticker{Symbol: "VNM", Sector: "Thực phẩm & đồ uống", Aliases: []string{"Vinamilk"}}))
	got := tg.Tag(tagger.Input{
		Title: "SeABank vào rổ chỉ số của VanEck",
		Body: "Quỹ VanEck Vectors Vietnam ETF (VNM ETF) công bố danh mục kỳ cơ cấu mới, " +
			"cổ phiếu SSB được thêm vào rổ với tỷ trọng 1,2%.",
	})
	assert.NotContains(t, symbols(got), "VNM", "VNM ETF là quỹ, không phải mã Vinamilk")
}

// G-04: alias mơ hồ không được kéo nhầm mã mẹ.
func TestNegativeAliasBlocksWrongParentCompany(t *testing.T) {
	tg := tagger.New([]tagger.Ticker{
		{Symbol: "MWG", Sector: "Bán lẻ",
			Aliases:         []string{"Thế Giới Di Động", "Bách Hoá Xanh"},
			NegativeAliases: []string{"Điện Máy Xanh"}},
		{Symbol: "VIC", Sector: "Bất động sản",
			Aliases: []string{"Vingroup"}, NegativeAliases: []string{"Vinhomes", "Vincom Retail"}},
		{Symbol: "VHM", Sector: "Bất động sản", Aliases: []string{"Vinhomes"}},
	})

	got := tg.Tag(tagger.Input{
		Title: "CTCP Đầu tư Điện Máy Xanh chào sàn HOSE",
		Body:  "CTCP Đầu tư Điện Máy Xanh (DMX – HOSE) niêm yết 120 triệu cổ phiếu trong quý IV.",
	})
	assert.NotContains(t, symbols(got), "MWG", "‘Điện Máy Xanh’ ở đây là pháp nhân DMX")

	got = tg.Tag(tagger.Input{
		Title: "Vinhomes ghi nhận doanh thu quý III tăng 18%",
		Body:  "Vinhomes công bố doanh thu 32.000 tỷ đồng, lợi nhuận sau thuế 9.100 tỷ đồng.",
	})
	assert.Contains(t, symbols(got), "VHM")
	assert.NotContains(t, symbols(got), "VIC", "công ty con không được map về công ty mẹ")
}

// ----------------------------------------------------------------- các rule khác

func TestMaxEightTickersPerArticle(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Nhóm ngân hàng đồng loạt tăng giá",
		Body: "Cổ phiếu VCB tăng 2%, BID tăng 1,5%, CTG tăng 1,2%, TCB tăng 3%, " +
			"ACB tăng 1%, STB tăng 2,5%, SHB tăng 1,8%, MSB tăng 0,9%, TPB tăng 1,1%, " +
			"VHM tăng 0,5%, VIC tăng 4,31%.",
	})
	assert.LessOrEqual(t, len(got), 8, "tối đa 8 mã hiển thị (rule B8)")
	assert.NotEmpty(t, got)
}

func TestDeterministicOutput(t *testing.T) {
	tg := tagger.New(testTickers())
	in := tagger.Input{
		Title: "Hòa Phát và Hoa Sen cùng tăng giá",
		Body:  "Cổ phiếu HPG tăng 2%, HSG tăng 3%, khối lượng khớp lệnh tăng mạnh.",
	}
	first := tg.Tag(in)
	for i := 0; i < 5; i++ {
		assert.Equal(t, first, tg.Tag(in), "tagger phải tất định (ADR-005)")
	}
}

func TestEmptyInputReturnsNothing(t *testing.T) {
	tg := tagger.New(testTickers())
	assert.Empty(t, tg.Tag(tagger.Input{}))
}

func TestTickerInsideQuotedSourceOnlyOnceIsDropped(t *testing.T) {
	tg := tagger.New(testTickers())
	got := tg.Tag(tagger.Input{
		Title: "Triển vọng nhóm bất động sản",
		Body:  "Báo cáo của MBS cho thấy nguồn cung căn hộ tăng 12% so với cùng kỳ.",
	})
	assert.NotContains(t, symbols(got), "MBS")
}

// Tầng ngành phải TẮT ở cấu hình mặc định: QA đếm toàn bộ 38 tag của tầng này
// trên dữ liệu thật và cả 38 đều sai.
func TestIndustryTierIsOffByDefault(t *testing.T) {
	got := tagger.New(testTickers()).Tag(tagger.Input{
		Title: "Ngành thép Việt Nam đối mặt CBAM của EU",
		Body:  "Cơ chế điều chỉnh carbon của EU áp lên ngành thép từ đầu 2026. Chi phí ngành thép tăng.",
	})
	assert.Empty(t, symbols(got), "mặc định không được suy luận mã theo ngành, nhận được %v", symbols(got))

	withTier := tagger.NewWithIndustryTier(testTickers()).Tag(tagger.Input{
		Title: "Ngành thép Việt Nam đối mặt CBAM của EU",
		Body:  "Cơ chế điều chỉnh carbon của EU áp lên ngành thép từ đầu 2026. Chi phí ngành thép tăng.",
	})
	assert.NotEmpty(t, withTier, "bật cờ thì tầng ngành vẫn hoạt động như cũ")
}

// P1-NEW-3: alias có dấu không được khớp cụm thường viết khác dấu.
func TestAliasMatchingIsDiacriticSensitive(t *testing.T) {
	tg := tagger.New([]tagger.Ticker{
		{Symbol: "GAS", Sector: "Dầu khí", Aliases: []string{"PV GAS", "Khí Việt Nam"}},
	})

	got := tg.Tag(tagger.Input{
		Title: "Khi Việt Nam bước vào giai đoạn dân số già",
		Body:  "Khi Việt Nam chuyển sang cơ cấu dân số già, chi tiêu y tế được dự báo tăng mạnh.",
	})
	assert.NotContains(t, symbols(got), "GAS", "‘Khi Việt Nam’ không phải alias ‘Khí Việt Nam’")

	got = tg.Tag(tagger.Input{
		Title: "Khí Việt Nam báo lãi quý III",
		Body:  "Tổng Công ty Khí Việt Nam ghi nhận doanh thu 25.000 tỷ đồng, lợi nhuận tăng 12%.",
	})
	assert.Contains(t, symbols(got), "GAS", "alias đúng dấu vẫn phải khớp")
}

// P1-NEW-4: negative_aliases phải chặn cả mã trần ở tầng điểm 0.90.
func TestNegativeAliasBlocksBareCodeTier(t *testing.T) {
	tg := tagger.New([]tagger.Ticker{
		{Symbol: "FPT", Sector: "Công nghệ", Aliases: []string{"Tập đoàn FPT"},
			NegativeAliases: []string{"FPT Retail", "FPT Telecom", "FPT Shop"}},
	})

	got := tg.Tag(tagger.Input{
		Title: "Chuỗi bán lẻ mở rộng quy mô",
		Body:  "FPT Retail (FRT) công bố doanh thu 12.000 tỷ đồng, cổ phiếu tăng 3% trong phiên.",
	})
	assert.NotContains(t, symbols(got), "FPT", "‘FPT Retail’ là pháp nhân FRT, không phải FPT")

	got = tg.Tag(tagger.Input{
		Title: "FPT ký hợp đồng 100 triệu USD",
		Body:  "Cổ phiếu FPT tăng 2,69%. FPT ghi nhận doanh thu tăng 25%.",
	})
	assert.Contains(t, symbols(got), "FPT")
}
