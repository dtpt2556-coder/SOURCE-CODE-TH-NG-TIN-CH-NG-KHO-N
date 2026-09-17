/**
 * Dữ liệu mẫu — lấy nguyên văn từ `docs/00-inputs/reference-images.md`.
 *
 * Mục đích: khi `api` chưa chạy (dev, `npm run build`, demo), toàn bộ trang vẫn
 * render được thay vì crash. Mọi trang dùng fixture đều hiện banner
 * "Đang dùng dữ liệu mẫu" để không ai nhầm đây là dữ liệu thật.
 */

import type {
  NewsItem,
  ResearchNote,
  ResearchNoteSummary,
  SiteMeta,
  Ticker,
} from "./types";

export const FIXTURE_NEWS: NewsItem[] = [
  {
    id: 1001,
    published_at: "2026-08-26T01:30:00Z",
    published_date_vn: "26/08/2026",
    title: "12 ngân hàng cam kết 408.000 tỷ đồng tín dụng cho DNNVV",
    summary_md:
      "12 ngân hàng cam kết **408.000 tỷ đồng** tín dụng cho DNNVV theo chương trình NHNN chủ trì. Big4 đăng ký 220.000 tỷ (Agribank 70.000; BIDV, Vietcombank, VietinBank mỗi bên 50.000 tỷ), lãi suất thấp hơn ít nhất 1 điểm % so bình quân cùng kỳ hạn. 8 ngân hàng tư nhân đăng ký 188.000 tỷ, giảm lãi 0,5-2 điểm % tuỳ lĩnh vực kèm miễn giảm phí. Cuối tháng 7 chỉ 8,8% DNNVV tiếp cận được vốn so với trên 47% ở doanh nghiệp lớn.",
    news_type: "nganh",
    news_type_label: "Ngành",
    tickers: [
      { symbol: "BID", relevance: "primary" },
      { symbol: "VCB", relevance: "primary" },
      { symbol: "CTG", relevance: "primary" },
      { symbol: "SHB", relevance: "mentioned" },
      { symbol: "MSB", relevance: "mentioned" },
      { symbol: "STB", relevance: "mentioned" },
      { symbol: "BVB", relevance: "mentioned" },
      { symbol: "NAB", relevance: "mentioned" },
      { symbol: "NVB", relevance: "mentioned" },
      { symbol: "SGB", relevance: "mentioned" },
      { symbol: "TPB", relevance: "mentioned" },
    ],
    source: { name: "VnExpress", domain: "vnexpress.net", tier: 1 },
    short_link: "/r/a7Kx2p",
  },
  {
    id: 1002,
    published_at: "2026-08-27T02:15:00Z",
    published_date_vn: "27/08/2026",
    title: "CBAM của EU áp dụng từ đầu 2026, thép Việt chịu chi phí bù đắp lớn",
    summary_md:
      "CBAM của EU áp dụng từ đầu 2026, doanh nghiệp có nguy cơ bù đắp **hơn 100 USD/tấn thép** xuất sang EU, nhôm tới cả nghìn USD/tấn. Giá chứng chỉ carbon EU trên 75 USD/tấn CO2. Lô 100.000 tấn thép đơn giá 500-550 USD/tấn phải nộp khoảng 300 tỷ đồng, tương đương 20% giá trị đơn hàng. Cường độ phát thải thép Việt 2,5 tấn CO2/tấn, nhôm 14 tấn so chuẩn EU 1,4 tấn. Sáu nhóm chịu CBAM: sắt thép, nhôm, xi măng, phân bón, điện, hydrogen.",
    news_type: "nganh",
    news_type_label: "Ngành",
    tickers: [
      { symbol: "HPG", relevance: "primary" },
      { symbol: "HSG", relevance: "primary" },
      { symbol: "NKG", relevance: "mentioned" },
      { symbol: "TIS", relevance: "mentioned" },
    ],
    source: { name: "VnExpress", domain: "vnexpress.net", tier: 1 },
    short_link: "/r/b3Qm9d",
  },
  {
    id: 1003,
    published_at: "2026-08-26T09:05:00Z",
    published_date_vn: "26/08/2026",
    title: "VN-Index lần đầu vượt 1.800 điểm, phiên tăng thứ 5 liên tiếp",
    summary_md:
      "VN-Index đóng cửa **1.821 điểm (+29,91 điểm, +1,67%)**, lần đầu vượt 1.800, phiên tăng thứ 5 liên tiếp. Giá trị giao dịch ~20.000 tỷ, vượt bình quân 20 phiên; 182 mã tăng / 125 giảm. TCB trần lên 33.450 đ khớp 39,17 triệu cp, BCM trần 44.450 đ, VIC +4,31%, FPT +2,69%. Khối ngoại mua ròng ~14 tỷ toàn thị trường nhưng bán ròng 43 tỷ trên HOSE: mua FPT 203 tỷ, TCB 136 tỷ, PNJ 100 tỷ; bán CTG 78 tỷ, ACB 78 tỷ, VHM 61 tỷ.",
    news_type: "nganh",
    news_type_label: "Ngành",
    tickers: [
      { symbol: "TCB", relevance: "primary" },
      { symbol: "BCM", relevance: "primary" },
      { symbol: "VIC", relevance: "primary" },
      { symbol: "FPT", relevance: "primary" },
    ],
    source: { name: "CafeF", domain: "cafef.vn", tier: 1 },
    short_link: "/r/c8Tz1k",
  },
  {
    id: 1004,
    published_at: "2026-08-27T03:40:00Z",
    published_date_vn: "27/08/2026",
    title: "Sáu tháng đầu 2026 nhập 494.000 tấn thịt, khoảng 1,5 tỷ USD",
    summary_md:
      "Sáu tháng đầu 2026 nhập **494.000 tấn thịt, ~1,5 tỷ USD (+10% lượng, +36,4% giá trị)**; cả năm 2025 là 1,95 tỷ USD. Gia cầm 186.400 tấn (đùi gà đông lạnh 130.000 tấn), thịt trâu 108.000 tấn, thịt heo ~80.000 tấn. Quý II nhập 42.600 tấn thịt heo, giá bình quân 2.152 USD/tấn giảm 18,5%. Nguồn chính: Mỹ 63.300 tấn, Thổ Nhĩ Kỳ 29.600 tấn (+261%), Tây Ban Nha 13.300 tấn. Giá heo hơi trong nước 56.000 đ/kg chịu sức ép từ hàng đông lạnh giá rẻ.",
    news_type: "nganh",
    news_type_label: "Ngành",
    tickers: [
      { symbol: "DBC", relevance: "primary" },
      { symbol: "BAF", relevance: "primary" },
      { symbol: "MML", relevance: "mentioned" },
      { symbol: "HAG", relevance: "mentioned" },
    ],
    source: { name: "Thanh Niên", domain: "thanhnien.vn", tier: 1 },
    short_link: "/r/d2Wr7v",
  },
  {
    id: 1005,
    published_at: "2026-08-27T06:20:00Z",
    published_date_vn: "27/08/2026",
    title: "ACBS ước 5.588 tỷ đồng chảy vào cổ phiếu Việt kỳ cơ cấu FTSE GEIS",
    summary_md:
      "ACBS ước **5.588 tỷ đồng (216 triệu USD)** chảy vào 117 cổ phiếu Việt Nam kỳ cơ cấu FTSE GEIS hiệu lực 21/9/2026, từ 17 quỹ thụ động; phân bổ 3 large-cap, 3 mid-cap, 21 small-cap, 90 micro-cap. Ba mã hút mạnh nhất: VIC 80,67 triệu USD (2.092 tỷ), VHM 25,43 triệu USD (659 tỷ), HPG 13,11 triệu USD (340 tỷ). Hai quỹ tham chiếu lớn nhất là Vanguard Total International Stock Index Fund 646,2 tỷ USD và Vanguard FTSE Emerging Markets ETF 162,3 tỷ USD.",
    news_type: "nganh",
    news_type_label: "Ngành",
    tickers: [
      { symbol: "VIC", relevance: "primary" },
      { symbol: "VHM", relevance: "primary" },
      { symbol: "HPG", relevance: "primary" },
    ],
    source: { name: "CafeF", domain: "cafef.vn", tier: 1 },
    short_link: "/r/e6Yn4s",
  },
];

export const FIXTURE_TICKERS: Ticker[] = [
  { symbol: "BID", company_name: "Ngân hàng TMCP Đầu tư và Phát triển Việt Nam", short_name: "BIDV", exchange: "HOSE", sector: "Ngân hàng", in_vn30: true },
  { symbol: "VCB", company_name: "Ngân hàng TMCP Ngoại thương Việt Nam", short_name: "Vietcombank", exchange: "HOSE", sector: "Ngân hàng", in_vn30: true },
  { symbol: "CTG", company_name: "Ngân hàng TMCP Công Thương Việt Nam", short_name: "VietinBank", exchange: "HOSE", sector: "Ngân hàng", in_vn30: true },
  { symbol: "SHB", company_name: "Ngân hàng TMCP Sài Gòn – Hà Nội", short_name: "SHB", exchange: "HOSE", sector: "Ngân hàng", in_vn30: true },
  { symbol: "MSB", company_name: "Ngân hàng TMCP Hàng Hải Việt Nam", short_name: "MSB", exchange: "HOSE", sector: "Ngân hàng", in_vn30: false },
  { symbol: "STB", company_name: "Ngân hàng TMCP Sài Gòn Thương Tín", short_name: "Sacombank", exchange: "HOSE", sector: "Ngân hàng", in_vn30: true },
  { symbol: "BVB", company_name: "Ngân hàng TMCP Bản Việt", short_name: "BVBank", exchange: "UPCOM", sector: "Ngân hàng", in_vn30: false },
  { symbol: "NAB", company_name: "Ngân hàng TMCP Nam Á", short_name: "Nam A Bank", exchange: "HOSE", sector: "Ngân hàng", in_vn30: false },
  { symbol: "NVB", company_name: "Ngân hàng TMCP Quốc Dân", short_name: "NCB", exchange: "HNX", sector: "Ngân hàng", in_vn30: false },
  { symbol: "SGB", company_name: "Ngân hàng TMCP Sài Gòn Công Thương", short_name: "Saigonbank", exchange: "UPCOM", sector: "Ngân hàng", in_vn30: false },
  { symbol: "TPB", company_name: "Ngân hàng TMCP Tiên Phong", short_name: "TPBank", exchange: "HOSE", sector: "Ngân hàng", in_vn30: true },
  { symbol: "TCB", company_name: "Ngân hàng TMCP Kỹ Thương Việt Nam", short_name: "Techcombank", exchange: "HOSE", sector: "Ngân hàng", in_vn30: true },
  { symbol: "HPG", company_name: "CTCP Tập đoàn Hòa Phát", short_name: "Hòa Phát", exchange: "HOSE", sector: "Thép", in_vn30: true },
  { symbol: "HSG", company_name: "CTCP Tập đoàn Hoa Sen", short_name: "Hoa Sen", exchange: "HOSE", sector: "Thép", in_vn30: false },
  { symbol: "NKG", company_name: "CTCP Thép Nam Kim", short_name: "Nam Kim", exchange: "HOSE", sector: "Thép", in_vn30: false },
  { symbol: "TIS", company_name: "CTCP Gang thép Thái Nguyên", short_name: "Tisco", exchange: "UPCOM", sector: "Thép", in_vn30: false },
  { symbol: "BCM", company_name: "Tổng Công ty Đầu tư và Phát triển Công nghiệp – CTCP", short_name: "Becamex IDC", exchange: "HOSE", sector: "Bất động sản khu công nghiệp", in_vn30: true },
  { symbol: "VIC", company_name: "Tập đoàn Vingroup – CTCP", short_name: "Vingroup", exchange: "HOSE", sector: "Bất động sản", in_vn30: true },
  { symbol: "VHM", company_name: "CTCP Vinhomes", short_name: "Vinhomes", exchange: "HOSE", sector: "Bất động sản", in_vn30: true },
  { symbol: "FPT", company_name: "CTCP FPT", short_name: "FPT", exchange: "HOSE", sector: "Công nghệ", in_vn30: true },
  { symbol: "CMG", company_name: "CTCP Tập đoàn Công nghệ CMC", short_name: "CMC", exchange: "HOSE", sector: "Công nghệ", in_vn30: false },
  { symbol: "DBC", company_name: "CTCP Tập đoàn Dabaco Việt Nam", short_name: "Dabaco", exchange: "HOSE", sector: "Chăn nuôi", in_vn30: false },
  { symbol: "BAF", company_name: "CTCP Nông nghiệp BAF Việt Nam", short_name: "BAF", exchange: "HOSE", sector: "Chăn nuôi", in_vn30: false },
  { symbol: "MML", company_name: "CTCP Masan MEATLife", short_name: "Masan MEATLife", exchange: "UPCOM", sector: "Chăn nuôi", in_vn30: false },
  { symbol: "HAG", company_name: "CTCP Hoàng Anh Gia Lai", short_name: "HAGL", exchange: "HOSE", sector: "Nông nghiệp", in_vn30: false },
];

const CMG_NOTE: ResearchNote = {
  id: 2001,
  slug: "cmg-quy-2-2026",
  symbol: "CMG",
  title:
    "CMG – Quý 2/2026: Doanh thu tiếp tục tăng nhưng lợi nhuận CĐ mẹ giảm – chu kỳ đầu tư Data Center chưa tạo ra lợi nhuận tương xứng",
  published_at: "2026-09-02T02:00:00Z",
  published_date_vn: "02/09/2026",
  excerpt:
    "Doanh thu đạt 2.323 tỷ đồng (+5,1% YoY) nhưng LNST CĐ mẹ giảm 20,3% còn 75 tỷ đồng khi chi phí lãi vay tăng 48%.",
  section_heading: "LUẬN ĐIỂM ĐẦU TƯ",
  disclaimer:
    "Nội dung chỉ mang tính thông tin, không phải khuyến nghị đầu tư. Nguồn tham khảo: báo cáo tài chính hợp nhất quý 2/2026 của CMC, MBS Research/CBRE và công bố thông tin của doanh nghiệp.",
  related_symbols: ["CMG", "FPT"],
  points: [
    {
      ordinal: 1,
      lead: "Hoạt động kinh doanh cốt lõi chưa xấu đi, nhưng Q2/2026 cho thấy doanh thu và lợi nhuận bắt đầu đi lệch nhau.",
      body: "Q2/2026, tương ứng quý đầu tiên của niên độ tài chính 2026–2027 của CMC, doanh thu đạt 2.323 tỷ đồng, tăng 5,1% YoY; biên gộp tăng nhẹ từ 17,7% lên 17,9%. Tuy nhiên, LNST hợp nhất giảm 13,1% còn khoảng 101 tỷ đồng và LNST CĐ mẹ giảm 20,3% còn 75 tỷ đồng. Nguyên nhân chính không nằm ở biên gộp mà ở chi phí vận hành và đặc biệt chi phí lãi vay tăng 48%. Điều này phản ánh giai đoạn CMC phải bỏ vốn trước cho hạ tầng số trong khi doanh thu mới chưa theo kịp.",
    },
    {
      ordinal: 2,
      lead: "Chất lượng tài sản công nghệ của CMC đáng chú ý hơn tốc độ lợi nhuận hiện tại.",
      body: "CMC có vị thế tương đối mạnh trong ba thị trường: tích hợp hệ thống và chuyển đổi số cho doanh nghiệp lớn, hạ tầng Cloud/Data Center trong nước và xuất khẩu dịch vụ CNTT. Khách hàng có thể xác định gồm Samsung SDS, IBM, Honda, Panasonic, VPBank, BIDV, VIB, ABBank cùng nhiều cơ quan Chính phủ. Theo MBS Research/CBRE, CMC Telecom chiếm khoảng 12% thị trường Data Center Việt Nam; CMC cũng công bố CMC Cloud chiếm trên 25% thị trường cloud nội địa.",
    },
    {
      ordinal: 3,
      lead: "Hyperscale Data Center là tài sản có khả năng thay đổi quy mô CMG, nhưng chưa thể coi là lợi nhuận đã chắc chắn.",
      body: "Dự án tại Khu Công nghệ cao TP.HCM có tổng vốn giai đoạn đầu trên 250 triệu USD, công suất thiết kế ban đầu 30 MW và có khả năng mở rộng trên 100 MW. Nhu cầu AI, Cloud và lưu trữ dữ liệu tạo thị trường thuận lợi, nhưng hiệu quả đầu tư cuối cùng còn phụ thuộc vào tỷ lệ lấp đầy, giá thuê, chi phí điện, cấu trúc vốn và thời gian đưa dự án vào khai thác.",
    },
    {
      ordinal: 4,
      lead: 'Ở 24.150 đồng/cp, định giá đã giảm đáng kể nhưng chưa phải trường hợp "rẻ rõ ràng".',
      body: "P/E TTM hiện khoảng 14 lần, thấp hơn đáng kể vùng định giá trước đây của CMG. Tuy nhiên, ROE chỉ quanh 12% và lợi nhuận đang chịu áp lực trong khi CAPEX tăng mạnh. Biên an toàn chỉ thực sự xuất hiện nếu CMC chứng minh được rằng giai đoạn đầu tư hiện nay sẽ chuyển thành tăng trưởng LNST CĐ mẹ từ 2027 trở đi.",
    },
  ],
};

export const FIXTURE_NOTES: ResearchNote[] = [CMG_NOTE];

export function fixtureNoteSummaries(): ResearchNoteSummary[] {
  return FIXTURE_NOTES.map((note) => ({
    id: note.id,
    slug: note.slug,
    symbol: note.symbol,
    title: note.title,
    published_at: note.published_at,
    published_date_vn: note.published_date_vn,
    excerpt: note.excerpt,
  }));
}

const VN_OFFSET_MS = 7 * 60 * 60 * 1000;
const CRON_HOURS_VN = [22, 14, 6];

/**
 * Mốc cron gần nhất (06:00 / 14:00 / 22:00 giờ VN — BA1 §8) tính từ thời điểm
 * hiện tại. Giữ cho dữ liệu mẫu luôn "tươi" thay vì luôn báo quá hạn 12h.
 */
export function lastFixtureCrawlAt(now: Date = new Date()): string {
  const vn = new Date(now.getTime() + VN_OFFSET_MS);
  const hour = vn.getUTCHours();
  const slot = CRON_HOURS_VN.find((h) => h <= hour);
  const base = Date.UTC(
    vn.getUTCFullYear(),
    vn.getUTCMonth(),
    vn.getUTCDate(),
    slot ?? 22,
    0,
    0,
  );
  const shifted = slot === undefined ? base - 24 * 60 * 60 * 1000 : base;
  return new Date(shifted - VN_OFFSET_MS).toISOString();
}

export function fixtureSiteMeta(now: Date = new Date()): SiteMeta {
  const counts = new Map<string, number>();
  for (const item of FIXTURE_NEWS) {
    counts.set(item.news_type, (counts.get(item.news_type) ?? 0) + 1);
  }

  return {
    last_crawl_at: lastFixtureCrawlAt(now),
    is_stale: false,
    total_articles: FIXTURE_NEWS.length,
    news_types: Array.from(counts.entries()).map(([news_type, count]) => ({
      news_type: news_type as SiteMeta["news_types"][number]["news_type"],
      label: FIXTURE_NEWS.find((n) => n.news_type === news_type)?.news_type_label ?? news_type,
      count,
    })),
  };
}
