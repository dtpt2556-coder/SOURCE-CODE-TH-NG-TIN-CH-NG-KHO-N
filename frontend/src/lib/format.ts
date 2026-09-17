/**
 * Định dạng số / ngày kiểu Việt Nam, múi giờ Asia/Ho_Chi_Minh.
 * BA1 §7: `408.000`, `1,67%`, ngày `DD/MM/YYYY`.
 */

export const TIMEZONE = "Asia/Ho_Chi_Minh";
export const LOCALE = "vi-VN";

const dateFormatter = new Intl.DateTimeFormat(LOCALE, {
  timeZone: TIMEZONE,
  day: "2-digit",
  month: "2-digit",
  year: "numeric",
});

const timeFormatter = new Intl.DateTimeFormat(LOCALE, {
  timeZone: TIMEZONE,
  hour: "2-digit",
  minute: "2-digit",
  hour12: false,
});

function toDate(value: string | Date): Date | null {
  const date = value instanceof Date ? value : new Date(value);
  return Number.isNaN(date.getTime()) ? null : date;
}

/** `26/08/2026` */
export function formatDateVN(value: string | Date): string {
  const date = toDate(value);
  return date ? dateFormatter.format(date) : "-";
}

/** `06:00 16/09/2026` */
export function formatDateTimeVN(value: string | Date): string {
  const date = toDate(value);
  if (!date) return "-";
  return `${timeFormatter.format(date)} ${dateFormatter.format(date)}`;
}

/** `2026-08-26` - dùng cho `<input type="date">` và thuộc tính `dateTime`. */
export function toISODateVN(value: string | Date): string {
  const date = toDate(value);
  if (!date) return "";
  const parts = dateFormatter.formatToParts(date);
  const get = (type: Intl.DateTimeFormatPartTypes) =>
    parts.find((p) => p.type === type)?.value ?? "";
  return `${get("year")}-${get("month")}-${get("day")}`;
}

/** `408.000` - dấu `.` phân nhóm nghìn, dấu `,` thập phân. */
export function formatNumberVN(value: number, fractionDigits = 0): string {
  if (!Number.isFinite(value)) return "-";
  return new Intl.NumberFormat(LOCALE, {
    minimumFractionDigits: fractionDigits,
    maximumFractionDigits: fractionDigits,
  }).format(value);
}

/**
 * `+1,67%` - luôn kèm dấu `+`/`−` để không truyền đạt thông tin chỉ bằng màu
 * (a11y checklist §8).
 */
export function formatPercentVN(value: number, fractionDigits = 2): string {
  if (!Number.isFinite(value)) return "-";
  const sign = value > 0 ? "+" : value < 0 ? "−" : "";
  return `${sign}${formatNumberVN(Math.abs(value), fractionDigits)}%`;
}

const TWELVE_HOURS_MS = 12 * 60 * 60 * 1000;

/** BA1 US-1.3 AC2: lần crawl gần nhất > 12 giờ trước → cảnh báo. */
export function isStale(lastCrawlAt: string | Date, now: Date = new Date()): boolean {
  const date = toDate(lastCrawlAt);
  if (!date) return true;
  return now.getTime() - date.getTime() > TWELVE_HOURS_MS;
}

/** Bỏ `www.` khỏi domain - BA1 US-3.1 AC1. */
export function displayDomain(domain: string): string {
  return domain.replace(/^www\./i, "");
}
