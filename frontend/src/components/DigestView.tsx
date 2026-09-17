"use client";

import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState, useSyncExternalStore } from "react";

import { loadMoreNews } from "@/app/actions";
import { DigestTable } from "@/components/DigestTable";
import { FilterBar } from "@/components/FilterBar";
import { EmptyState } from "@/components/Notices";
import { formatNumberVN } from "@/lib/format";
import { toPlainText } from "@/lib/markdown";
import { buildHref, EMPTY_FILTERS, type NewsFilters } from "@/lib/query";
import type { ListMeta, NewsItem, NewsTypeCount, Ticker } from "@/lib/types";
import {
  clearWatchlist,
  getServerWatchlistSnapshot,
  getWatchlistSnapshot,
  markWatchlistApplied,
  saveWatchlist,
  subscribeWatchlist,
} from "@/lib/watchlist";

const TOAST_MS = 2000;

interface DigestViewProps {
  initialItems: NewsItem[];
  initialMeta: ListMeta;
  filters: NewsFilters;
  tickers: Ticker[];
  typeCounts: NewsTypeCount[];
  caption: string;
}

export function DigestView({
  initialItems,
  initialMeta,
  filters,
  tickers,
  typeCounts,
  caption,
}: DigestViewProps) {
  const router = useRouter();

  const [items, setItems] = useState<NewsItem[]>(initialItems);
  const [hasMore, setHasMore] = useState(initialMeta.has_more);
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState<string | null>(null);

  const watchlist = useSyncExternalStore(
    subscribeWatchlist,
    getWatchlistSnapshot,
    getServerWatchlistSnapshot,
  );

  /**
   * BA1 US-5.2 AC1: quay lại sau thì bộ lọc tự áp dụng từ localStorage.
   * Chỉ chạy 1 lần mỗi phiên và chỉ khi URL chưa có bộ lọc mã nào.
   */
  useEffect(() => {
    const stored = getWatchlistSnapshot();
    if (!stored.length || filters.ma.length > 0) return;
    if (!markWatchlistApplied()) return;

    router.replace(buildHref({ ...EMPTY_FILTERS, ...filters, ma: [...stored] }), {
      scroll: false,
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (!toast) return;
    const timer = window.setTimeout(() => setToast(null), TOAST_MS);
    return () => window.clearTimeout(timer);
  }, [toast]);

  const handleLoadMore = useCallback(async () => {
    setLoading(true);
    try {
      const result = await loadMoreNews(filters, items.length);
      setItems((previous) => {
        const seen = new Set(previous.map((item) => item.id));
        return [...previous, ...result.items.filter((item) => !seen.has(item.id))];
      });
      setHasMore(result.hasMore);
    } catch {
      setToast("Không tải thêm được tin. Vui lòng thử lại.");
    } finally {
      setLoading(false);
    }
  }, [filters, items.length]);

  /**
   * BA1 US-5.1 AC1: mỗi tin là một khối 3 dòng.
   * Ba dòng thay cho một dòng dài vì môi giới dán thẳng sang Zalo.
   */
  const handleCopy = useCallback(async () => {
    const origin = window.location.origin;
    const text = items
      .map((item) => {
        const symbols = item.tickers.map((ticker) => ticker.symbol).join(", ");
        return [
          `[${item.published_date_vn}] [${symbols}]`,
          toPlainText(item.summary_md),
          `Nguồn: ${origin}${item.short_link}`,
        ].join("\n");
      })
      .join("\n\n");

    try {
      await navigator.clipboard.writeText(text);
      setToast(`Đã copy ${formatNumberVN(items.length)} tin vào clipboard.`);
    } catch {
      setToast("Trình duyệt chặn quyền truy cập clipboard.");
    }
  }, [items]);

  const handleSaveWatchlist = useCallback((symbols: string[]) => {
    saveWatchlist(symbols);
    setToast("Đã lưu danh sách theo dõi trên trình duyệt này.");
  }, []);

  const handleClearWatchlist = useCallback(() => {
    clearWatchlist();
    setToast("Đã xoá danh sách theo dõi.");
  }, []);

  return (
    <>
      <FilterBar
        filters={filters}
        tickers={tickers}
        typeCounts={typeCounts}
        watchlist={watchlist}
        onCopy={handleCopy}
        copyDisabled={items.length === 0}
        onSaveWatchlist={handleSaveWatchlist}
        onClearWatchlist={handleClearWatchlist}
      />

      <div className="mx-auto max-w-content px-4 md:px-6">
        {items.length === 0 ? (
          <EmptyState
            title="Chưa có tin trong khoảng thời gian này."
            hint="Thử bỏ bớt mã CK, mở rộng khoảng ngày hoặc xoá từ khoá tìm kiếm."
          />
        ) : (
          <>
            <div className="fade-in">
              <DigestTable items={items} caption={caption} />
            </div>

            <p className="sr-only" role="status">
              Đang hiển thị {formatNumberVN(items.length)} trên tổng{" "}
              {formatNumberVN(initialMeta.total)} tin.
            </p>

            <div className="mt-8 flex items-center gap-4 border-t border-rule pt-6">
              {hasMore ? (
                <button
                  className="btn btn-strong"
                  type="button"
                  onClick={handleLoadMore}
                  disabled={loading}
                >
                  {loading ? "Đang tải…" : "Xem thêm"}
                </button>
              ) : (
                <p className="text-[13.5px] text-ink-muted">Đã hiển thị hết tin.</p>
              )}
              <p className="text-[13.5px] text-ink-muted">
                {formatNumberVN(items.length)}/{formatNumberVN(initialMeta.total)} tin
              </p>
            </div>
          </>
        )}
      </div>

      {toast ? (
        <p className="toast" role="status" aria-live="polite">
          {toast}
        </p>
      ) : null}
    </>
  );
}
