/**
 * Watchlist lưu trên trình duyệt - BA1 US-5.2 / ADR-010 (không đăng nhập ở v1).
 *
 * Expose dưới dạng external store để component đọc bằng `useSyncExternalStore`:
 * tránh `setState` trong `useEffect` và tránh lệch hydration (server luôn trả
 * snapshot rỗng).
 */

export const WATCHLIST_KEY = "tenpoint.watchlist";
export const WATCHLIST_APPLIED_KEY = "tenpoint.watchlist-applied";

const EMPTY: readonly string[] = Object.freeze([]);

const listeners = new Set<() => void>();

let cachedRaw: string | null = null;
let cachedValue: readonly string[] = EMPTY;

function parse(raw: string | null): readonly string[] {
  if (!raw) return EMPTY;
  try {
    const parsed: unknown = JSON.parse(raw);
    if (!Array.isArray(parsed)) return EMPTY;
    const symbols = parsed.filter((value): value is string => typeof value === "string");
    return symbols.length ? Object.freeze(symbols) : EMPTY;
  } catch {
    return EMPTY;
  }
}

function emit() {
  for (const listener of listeners) listener();
}

export function subscribeWatchlist(listener: () => void): () => void {
  listeners.add(listener);
  window.addEventListener("storage", listener);
  return () => {
    listeners.delete(listener);
    window.removeEventListener("storage", listener);
  };
}

/** Memo hoá theo chuỗi thô để `useSyncExternalStore` nhận đúng tham chiếu ổn định. */
export function getWatchlistSnapshot(): readonly string[] {
  let raw: string | null = null;
  try {
    raw = window.localStorage.getItem(WATCHLIST_KEY);
  } catch {
    raw = null;
  }

  if (raw !== cachedRaw) {
    cachedRaw = raw;
    cachedValue = parse(raw);
  }
  return cachedValue;
}

export function getServerWatchlistSnapshot(): readonly string[] {
  return EMPTY;
}

export function saveWatchlist(symbols: string[]): void {
  try {
    window.localStorage.setItem(WATCHLIST_KEY, JSON.stringify(symbols));
    window.sessionStorage.setItem(WATCHLIST_APPLIED_KEY, "1");
  } catch {
    // Trình duyệt chặn storage (private mode) - bỏ qua, không làm hỏng UI.
  }
  emit();
}

export function clearWatchlist(): void {
  try {
    window.localStorage.removeItem(WATCHLIST_KEY);
  } catch {
    // bỏ qua
  }
  emit();
}

/** Đã áp dụng watchlist trong phiên này chưa (tránh ghi đè lựa chọn của người dùng). */
export function markWatchlistApplied(): boolean {
  try {
    if (window.sessionStorage.getItem(WATCHLIST_APPLIED_KEY) === "1") return false;
    window.sessionStorage.setItem(WATCHLIST_APPLIED_KEY, "1");
    return true;
  } catch {
    return false;
  }
}
