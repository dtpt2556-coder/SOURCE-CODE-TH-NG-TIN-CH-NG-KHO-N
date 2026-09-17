import { isFilterableNewsType } from "./news-types";
import type { ApiNewsType, RelevanceFilter } from "./types";

/**
 * Toàn bộ state bộ lọc nằm trong query string - BA1 US-2.1 AC3:
 * copy URL rồi mở ở máy khác thì filter phải được giữ nguyên.
 *
 *   /?ma=FPT,HPG&loai=nganh&q=data%20center&tu=2026-08-20&den=2026-08-27
 */
export interface NewsFilters {
  ma: string[];
  loai: ApiNewsType[];
  q: string;
  tu: string;
  den: string;
  relevance: RelevanceFilter;
}

export type RawSearchParams = Record<string, string | string[] | undefined>;

export const PAGE_SIZE = 20;

export const EMPTY_FILTERS: NewsFilters = {
  ma: [],
  loai: [],
  q: "",
  tu: "",
  den: "",
  relevance: "primary",
};

function first(value: string | string[] | undefined): string {
  if (Array.isArray(value)) return value[0] ?? "";
  return value ?? "";
}

function splitCsv(value: string): string[] {
  return value
    .split(",")
    .map((part) => part.trim())
    .filter(Boolean);
}

const SYMBOL_PATTERN = /^[A-Z0-9]{3,10}$/;
const DATE_PATTERN = /^\d{4}-\d{2}-\d{2}$/;

export function normalizeSymbol(value: string): string {
  return value.trim().toUpperCase();
}

export function parseFilters(searchParams: RawSearchParams): NewsFilters {
  const ma = Array.from(
    new Set(splitCsv(first(searchParams.ma)).map(normalizeSymbol).filter((s) => SYMBOL_PATTERN.test(s))),
  );

  const loai = Array.from(
    new Set(
      splitCsv(first(searchParams.loai))
        .map((value) => value.toLowerCase())
        .filter(isFilterableNewsType),
    ),
  );

  const tu = first(searchParams.tu);
  const den = first(searchParams.den);
  const relevanceRaw = first(searchParams.relevance);

  return {
    ma,
    loai,
    q: first(searchParams.q).trim().slice(0, 120),
    tu: DATE_PATTERN.test(tu) ? tu : "",
    den: DATE_PATTERN.test(den) ? den : "",
    relevance: relevanceRaw === "all" ? "all" : "primary",
  };
}

export function hasActiveFilters(filters: NewsFilters): boolean {
  return (
    filters.ma.length > 0 ||
    filters.loai.length > 0 ||
    filters.q.length > 0 ||
    filters.tu.length > 0 ||
    filters.den.length > 0 ||
    filters.relevance !== "primary"
  );
}

/** Query string cho URL trình duyệt (bỏ trường rỗng để URL luôn gọn). */
export function buildSearchParams(filters: NewsFilters): string {
  const params = new URLSearchParams();
  if (filters.ma.length) params.set("ma", filters.ma.join(","));
  if (filters.loai.length) params.set("loai", filters.loai.join(","));
  if (filters.q) params.set("q", filters.q);
  if (filters.tu) params.set("tu", filters.tu);
  if (filters.den) params.set("den", filters.den);
  if (filters.relevance !== "primary") params.set("relevance", filters.relevance);
  return params.toString();
}

export function buildHref(filters: NewsFilters, pathname = "/"): string {
  const qs = buildSearchParams(filters);
  return qs ? `${pathname}?${qs}` : pathname;
}

/** Query string cho `GET /api/v1/news` - architecture.md §6. */
export function buildNewsApiQuery(
  filters: NewsFilters,
  offset = 0,
  limit = PAGE_SIZE,
): string {
  const params = new URLSearchParams();
  if (filters.ma.length) params.set("ma", filters.ma.join(","));
  if (filters.loai.length) params.set("loai", filters.loai.join(","));
  if (filters.q) params.set("q", filters.q);
  if (filters.tu) params.set("tu", filters.tu);
  if (filters.den) params.set("den", filters.den);
  params.set("relevance", filters.relevance);
  params.set("limit", String(limit));
  params.set("offset", String(offset));
  return params.toString();
}

export function toggleSymbol(filters: NewsFilters, symbol: string): NewsFilters {
  const value = normalizeSymbol(symbol);
  const exists = filters.ma.includes(value);
  return {
    ...filters,
    ma: exists ? filters.ma.filter((s) => s !== value) : [...filters.ma, value],
  };
}
