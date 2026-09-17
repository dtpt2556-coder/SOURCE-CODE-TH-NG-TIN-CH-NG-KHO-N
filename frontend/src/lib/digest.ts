import { formatNumberVN } from "./format";
import type { NewsItem } from "./types";

/**
 * Mã xuất hiện nhiều nhất trong tập tin hiện tại.
 * Tie-break theo thứ tự xuất hiện đầu tiên để kết quả tất định.
 */
export function topSymbols(items: NewsItem[], limit = 3): string[] {
  const counts = new Map<string, number>();
  const firstSeen = new Map<string, number>();

  items.forEach((item, index) => {
    for (const ticker of item.tickers) {
      counts.set(ticker.symbol, (counts.get(ticker.symbol) ?? 0) + 1);
      if (!firstSeen.has(ticker.symbol)) firstSeen.set(ticker.symbol, index);
    }
  });

  return Array.from(counts.entries())
    .sort((a, b) => {
      if (b[1] !== a[1]) return b[1] - a[1];
      return (firstSeen.get(a[0]) ?? 0) - (firstSeen.get(b[0]) ?? 0);
    })
    .slice(0, limit)
    .map(([symbol]) => symbol);
}

/**
 * BA1 US-1.5: "8 tin mới đáng chú ý — PVS · FPT · VIC."
 * N = số tin của bộ lọc hiện tại, tối đa 3 mã nổi bật nối bằng " · ".
 */
export function buildHeadline(total: number, items: NewsItem[]): string {
  const symbols = topSymbols(items, 3);
  const count = formatNumberVN(total);
  if (!symbols.length) return `${count} tin mới đáng chú ý.`;
  return `${count} tin mới đáng chú ý — ${symbols.join(" · ")}.`;
}
