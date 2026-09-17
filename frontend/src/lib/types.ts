/**
 * Kiểu dữ liệu khớp hợp đồng API REST v1 - docs/03-sa/architecture.md §6.
 * Thời gian luôn là ISO-8601 UTC; trường `*_vn` đã được BE format theo giờ VN.
 */

/**
 * Taxonomy 9 nhan - PRD muc 5. BE da bo sung `phan_tich` va
 * `trai_phieu_tin_dung` vao enum o migration 0005, nen ca 9 deu loc duoc.
 */
export const API_NEWS_TYPES = [
  "vi_mo",
  "nganh",
  "doanh_nghiep",
  "thi_truong",
  "khoi_ngoai",
  "co_tuc_phat_hanh",
  "phap_ly",
  "phan_tich",
  "trai_phieu_tin_dung",
] as const;

export const NEWS_TYPES = API_NEWS_TYPES;

export type ApiNewsType = (typeof API_NEWS_TYPES)[number];
export type NewsType = (typeof NEWS_TYPES)[number];

export type TickerRelevance = "primary" | "mentioned";

/**
 * Giá trị hợp lệ của query param `relevance` trên `GET /api/v1/news`:
 * `primary` (mặc định) hoặc `all` (gồm cả mã chỉ được nhắc tới).
 */
export type RelevanceFilter = "primary" | "all";

export interface TickerRef {
  symbol: string;
  relevance: TickerRelevance;
}

export interface NewsSource {
  name: string;
  domain: string;
  tier: number;
  /**
   * Tieu de bai goc tren bao nguon - design-system v2 muc 5.3.
   * BE dang bo sung; khi chua co thi FE fallback ve `NewsItem.title`.
   */
  article_title?: string;
}

export interface NewsItem {
  id: number;
  published_at: string;
  published_date_vn: string;
  title: string;
  /** Markdown thô - chỉ chứa `**...**`. FE bắt buộc sanitize (BA1 US-1.2 AC2). */
  summary_md: string;
  news_type: NewsType;
  news_type_label: string;
  tickers: TickerRef[];
  source: NewsSource;
  short_link: string;
  /**
   * `short_links.target_alive` - BA1 Q6. `false` thi khong render thanh <a>
   * ma hien text mo kem chu thich (design-system v2 muc 5.3).
   * Thieu truong nay thi coi nhu con song.
   */
  short_link_alive?: boolean;
  /**
   * 5 tin dung lai tu anh tham chieu co URL khong co that (ingest-edge-cases
   * muc 0). Tuyet doi khong cho bam.
   */
  is_demo?: boolean;
  /** ISO-8601 UTC. Khac null nghia la nguon da sua bai sau khi dang (U-03). */
  updated_at?: string | null;
  /** T-06: khong tim duoc thoi gian dang, dang lay tam `fetched_at`. */
  published_at_estimated?: boolean;
  /** Tang moi lan bai duoc tom tat lai (U-03). */
  revision?: number;
}

export interface ListMeta {
  total: number;
  limit: number;
  offset: number;
  has_more: boolean;
}

export interface NewsListResponse {
  data: NewsItem[];
  meta: ListMeta;
}

export interface Ticker {
  symbol: string;
  company_name: string;
  short_name: string;
  exchange: string;
  sector: string;
  in_vn30: boolean;
}

export interface TickerListResponse {
  data: Ticker[];
}

export interface ResearchNoteSummary {
  /**
   * BE dinh danh research note bang `slug`, khong tra `id`
   * (backend/internal/http/dto.go NoteItem). Chi fixture moi co.
   */
  id?: number;
  slug: string;
  symbol: string;
  title: string;
  published_at: string;
  published_date_vn: string;
  excerpt: string;
}

export interface ResearchNotePoint {
  ordinal: number;
  /** Dòng dẫn in đậm */
  lead: string;
  /** Đoạn phân tích */
  body: string;
}

export interface ResearchNote extends ResearchNoteSummary {
  section_heading: string;
  disclaimer: string;
  points: ResearchNotePoint[];
  related_symbols: string[];
}

export interface NoteListResponse {
  data: ResearchNoteSummary[];
  meta: ListMeta;
}

export interface NoteResponse {
  data: ResearchNote;
}

export interface TickerDetail {
  ticker: Ticker;
  news: NewsItem[];
  notes: ResearchNoteSummary[];
}

export interface TickerDetailResponse {
  data: TickerDetail;
}

export interface NewsTypeCount {
  news_type: NewsType;
  label: string;
  count: number;
}

export interface SiteMeta {
  last_crawl_at: string;
  is_stale: boolean;
  total_articles: number;
  news_types: NewsTypeCount[];
}

export interface SiteMetaResponse {
  data: SiteMeta;
}

/** Định dạng lỗi thống nhất - architecture.md §6 */
export interface ApiError {
  error: {
    code: string;
    message: string;
  };
}
