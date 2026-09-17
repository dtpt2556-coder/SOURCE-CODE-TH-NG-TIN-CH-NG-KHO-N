import {
  FIXTURE_NEWS,
  FIXTURE_NOTES,
  FIXTURE_TICKERS,
  fixtureNoteSummaries,
  fixtureSiteMeta,
} from "./fixtures";
import { toPlainText } from "./markdown";
import { buildNewsApiQuery, PAGE_SIZE, type NewsFilters } from "./query";
import type {
  ListMeta,
  NewsItem,
  NewsListResponse,
  NoteListResponse,
  ResearchNote,
  ResearchNoteSummary,
  SiteMeta,
  Ticker,
  TickerDetail,
  TickerListResponse,
} from "./types";

/**
 * API client - chỉ chạy phía server (SSR). Base URL không bao giờ lộ ra client.
 * architecture.md §2: `web` gọi `api` qua tên service nội bộ trong Docker network.
 */
export const API_BASE_URL = process.env.API_BASE_URL ?? "http://localhost:8080";

const TIMEOUT_MS = Number(process.env.API_TIMEOUT_MS ?? 4000);

/** Kết quả kèm cờ cho biết dữ liệu đến từ fixture hay từ API thật. */
export interface Fetched<T> {
  data: T;
  fromFixture: boolean;
}

async function getJson<T>(path: string): Promise<T | null> {
  try {
    const response = await fetch(new URL(path, API_BASE_URL), {
      // SSR cho mọi trang công khai - luôn lấy dữ liệu mới nhất.
      cache: "no-store",
      headers: { Accept: "application/json" },
      signal: AbortSignal.timeout(TIMEOUT_MS),
    });

    if (!response.ok) return null;
    return (await response.json()) as T;
  } catch {
    // Backend chưa chạy / timeout / JSON hỏng → rơi về fixture, không bao giờ crash.
    return null;
  }
}

/* ------------------------------------------------------------------ */
/* Lọc fixture phía FE - mô phỏng hành vi của `GET /api/v1/news`        */
/* ------------------------------------------------------------------ */

function normalizeText(value: string): string {
  return value
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/đ/g, "d")
    .replace(/Đ/g, "D")
    .toLowerCase();
}

function isoDateOf(item: NewsItem): string {
  const [day, month, year] = item.published_date_vn.split("/");
  return `${year}-${month}-${day}`;
}

function filterFixtureNews(filters: NewsFilters): NewsItem[] {
  const needle = filters.q ? normalizeText(filters.q) : "";

  return FIXTURE_NEWS.filter((item) => {
    if (filters.ma.length) {
      const symbols = item.tickers
        .filter((t) => filters.relevance === "all" || t.relevance === "primary")
        .map((t) => t.symbol);
      if (!filters.ma.some((symbol) => symbols.includes(symbol))) return false;
    }

    if (
      filters.loai.length &&
      !(filters.loai as readonly string[]).includes(item.news_type)
    ) {
      return false;
    }

    const date = isoDateOf(item);
    if (filters.tu && date < filters.tu) return false;
    if (filters.den && date > filters.den) return false;

    if (needle) {
      const haystack = normalizeText(`${item.title} ${toPlainText(item.summary_md)}`);
      if (!haystack.includes(needle)) return false;
    }

    return true;
  }).sort((a, b) => b.published_at.localeCompare(a.published_at));
}

function paginate<T>(items: T[], offset: number, limit: number): { page: T[]; meta: ListMeta } {
  const page = items.slice(offset, offset + limit);
  return {
    page,
    meta: {
      total: items.length,
      limit,
      offset,
      has_more: offset + page.length < items.length,
    },
  };
}

/* ------------------------------------------------------------------ */
/* Public API                                                          */
/* ------------------------------------------------------------------ */

export async function fetchNews(
  filters: NewsFilters,
  offset = 0,
  limit = PAGE_SIZE,
): Promise<Fetched<NewsListResponse>> {
  const query = buildNewsApiQuery(filters, offset, limit);
  const response = await getJson<NewsListResponse>(`/api/v1/news?${query}`);

  if (response?.data) {
    return { data: response, fromFixture: false };
  }

  const { page, meta } = paginate(filterFixtureNews(filters), offset, limit);
  return { data: { data: page, meta }, fromFixture: true };
}

/* ------------------------------------------------------------------ */
/* Ánh xạ hình dạng BE sang hình dạng FE                                */
/* ------------------------------------------------------------------ */

/**
 * BE trả `TickerDetail` phẳng kèm `recent_news`/`research_notes`, và `NoteDetail`
 * không có `related_symbols` (backend/internal/http/dto.go:94-140). FE lại dùng
 * `{ticker, news, notes}` và gọi `.map()` thẳng trên `related_symbols`.
 *
 * Trước đây `getJson` ép kiểu `as T` mà không kiểm chứng, nên chỗ lệch này
 * không lộ ra lúc `tsc` mà nổ thành HTTP 500 lúc chạy thật (QA P0-1).
 * Ánh xạ gom về đây để các trang không phải biết tới khác biệt đó.
 */

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

function str(value: unknown, fallback = ""): string {
  return typeof value === "string" ? value : fallback;
}

function adaptNoteSummary(raw: unknown): ResearchNoteSummary | null {
  if (!isRecord(raw)) return null;
  const slug = str(raw.slug);
  if (!slug) return null;

  return {
    slug,
    symbol: str(raw.symbol),
    title: str(raw.title),
    published_at: str(raw.published_at),
    published_date_vn: str(raw.published_date_vn),
    excerpt: str(raw.excerpt),
  };
}

function adaptNoteDetail(raw: unknown): ResearchNote | null {
  const summary = adaptNoteSummary(raw);
  if (!summary || !isRecord(raw)) return null;

  const points = Array.isArray(raw.points)
    ? raw.points.filter(isRecord).map((point) => ({
        ordinal: typeof point.ordinal === "number" ? point.ordinal : 0,
        lead: str(point.lead),
        body: str(point.body),
      }))
    : [];

  // BE chưa trả `related_symbols`; mã của chính bài luôn là mã liên quan.
  const related = Array.isArray(raw.related_symbols)
    ? raw.related_symbols.filter((item): item is string => typeof item === "string")
    : summary.symbol
      ? [summary.symbol]
      : [];

  return {
    ...summary,
    excerpt: summary.excerpt || points[0]?.lead || summary.title,
    section_heading: str(raw.section_heading, "LUẬN ĐIỂM ĐẦU TƯ"),
    disclaimer: str(raw.disclaimer),
    points,
    related_symbols: related,
  };
}

function adaptTickerDetail(raw: unknown): TickerDetail | null {
  if (!isRecord(raw)) return null;
  const symbol = str(raw.symbol);
  if (!symbol) return null;

  return {
    ticker: {
      symbol,
      company_name: str(raw.company_name),
      short_name: str(raw.short_name),
      exchange: str(raw.exchange),
      sector: str(raw.sector),
      in_vn30: raw.in_vn30 === true,
    },
    news: Array.isArray(raw.recent_news) ? (raw.recent_news as NewsItem[]) : [],
    notes: Array.isArray(raw.research_notes)
      ? raw.research_notes
          .map(adaptNoteSummary)
          .filter((note): note is ResearchNoteSummary => note !== null)
      : [],
  };
}

function adaptSiteMeta(raw: unknown): SiteMeta | null {
  if (!isRecord(raw)) return null;

  // BE dùng khoá `value`, FE dùng `news_type`.
  const types = Array.isArray(raw.news_types)
    ? raw.news_types.filter(isRecord).map((entry) => ({
        news_type: str(entry.news_type) || str(entry.value),
        label: str(entry.label),
        count: typeof entry.count === "number" ? entry.count : 0,
      }))
    : [];

  return {
    last_crawl_at: str(raw.last_crawl_at),
    is_stale: raw.is_stale === true,
    total_articles: typeof raw.total_articles === "number" ? raw.total_articles : 0,
    news_types: types as SiteMeta["news_types"],
  };
}

export async function fetchTickers(): Promise<Fetched<Ticker[]>> {
  const response = await getJson<TickerListResponse>("/api/v1/tickers");
  if (response?.data) return { data: response.data, fromFixture: false };
  return { data: FIXTURE_TICKERS, fromFixture: true };
}

export async function fetchTickerDetail(symbol: string): Promise<Fetched<TickerDetail> | null> {
  const upper = symbol.toUpperCase();
  const response = await getJson<{ data?: unknown }>(
    `/api/v1/tickers/${encodeURIComponent(upper)}`,
  );
  const adapted = adaptTickerDetail(response?.data);
  if (adapted) return { data: adapted, fromFixture: false };

  const ticker = FIXTURE_TICKERS.find((t) => t.symbol === upper);
  if (!ticker) return null;

  const news = FIXTURE_NEWS.filter((item) =>
    item.tickers.some((t) => t.symbol === upper),
  ).sort((a, b) => b.published_at.localeCompare(a.published_at));

  const notes = fixtureNoteSummaries().filter(
    (note) =>
      note.symbol === upper ||
      FIXTURE_NOTES.find((n) => n.slug === note.slug)?.related_symbols.includes(upper),
  );

  return { data: { ticker, news, notes }, fromFixture: true };
}

export async function fetchNotes(limit = 50): Promise<Fetched<NoteListResponse>> {
  const response = await getJson<{ data?: unknown; meta?: ListMeta }>(
    `/api/v1/notes?limit=${limit}`,
  );
  if (Array.isArray(response?.data)) {
    const items = response.data
      .map(adaptNoteSummary)
      .filter((note): note is ResearchNoteSummary => note !== null);
    const meta = response.meta ?? {
      total: items.length,
      limit,
      offset: 0,
      has_more: false,
    };
    return { data: { data: items, meta }, fromFixture: false };
  }

  const summaries: ResearchNoteSummary[] = fixtureNoteSummaries();
  const { page, meta } = paginate(summaries, 0, limit);
  return { data: { data: page, meta }, fromFixture: true };
}

export async function fetchNote(slug: string): Promise<Fetched<ResearchNote> | null> {
  const response = await getJson<{ data?: unknown }>(
    `/api/v1/notes/${encodeURIComponent(slug)}`,
  );
  const adapted = adaptNoteDetail(response?.data);
  if (adapted) return { data: adapted, fromFixture: false };

  const note = FIXTURE_NOTES.find((n) => n.slug === slug);
  return note ? { data: note, fromFixture: true } : null;
}

export async function fetchSiteMeta(): Promise<Fetched<SiteMeta>> {
  const response = await getJson<{ data?: unknown }>("/api/v1/meta");
  const adapted = adaptSiteMeta(response?.data);
  if (adapted) return { data: adapted, fromFixture: false };
  return { data: fixtureSiteMeta(), fromFixture: true };
}
