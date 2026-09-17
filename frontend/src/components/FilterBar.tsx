"use client";

import { useRouter } from "next/navigation";
import { useId, useState, useTransition } from "react";

import { ThemeToggle } from "@/components/ThemeToggle";
import { FILTERABLE_NEWS_TYPES, NEWS_TYPE_LABELS } from "@/lib/news-types";
import {
  buildHref,
  EMPTY_FILTERS,
  hasActiveFilters,
  normalizeSymbol,
  toggleSymbol,
  type NewsFilters,
} from "@/lib/query";
import type { ApiNewsType, NewsTypeCount, Ticker } from "@/lib/types";

const MAX_CHIPS = 14;

interface FilterBarProps {
  filters: NewsFilters;
  tickers: Ticker[];
  typeCounts: NewsTypeCount[];
  watchlist: readonly string[];
  onCopy: () => void;
  copyDisabled: boolean;
  onSaveWatchlist: (symbols: string[]) => void;
  onClearWatchlist: () => void;
}

/**
 * §5.6 - sticky, nền `--surface-raised`, hairline dưới, control viền 1px radius 0.
 * Toàn bộ state bộ lọc được đẩy vào query string (BA1 US-2.1 AC3).
 */
export function FilterBar({
  filters,
  tickers,
  typeCounts,
  watchlist,
  onCopy,
  copyDisabled,
  onSaveWatchlist,
  onClearWatchlist,
}: FilterBarProps) {
  const router = useRouter();
  const [isPending, startTransition] = useTransition();

  const searchId = useId();
  const typeId = useId();
  const fromId = useId();
  const toId = useId();
  const addSymbolId = useId();
  const symbolListId = useId();

  const [q, setQ] = useState(filters.q);
  const [tu, setTu] = useState(filters.tu);
  const [den, setDen] = useState(filters.den);
  const [newSymbol, setNewSymbol] = useState("");

  const apply = (next: NewsFilters) => {
    startTransition(() => {
      router.push(buildHref(next), { scroll: false });
    });
  };

  const countOf = (type: ApiNewsType) =>
    typeCounts.find((entry) => entry.news_type === type)?.count;

  const popular = Array.from(
    new Set([
      ...filters.ma,
      ...watchlist,
      ...tickers.filter((t) => t.in_vn30).map((t) => t.symbol),
      ...tickers.map((t) => t.symbol),
    ]),
  ).slice(0, MAX_CHIPS);

  return (
    <div className="filter-bar">
      <div className="mx-auto max-w-content px-4 py-3 md:px-6">
        <form
          className="flex flex-wrap items-end gap-x-3 gap-y-3"
          onSubmit={(event) => {
            event.preventDefault();
            apply({ ...filters, q: q.trim(), tu, den });
          }}
        >
          <div className="flex min-w-[200px] flex-1 flex-col gap-1">
            <label className="text-[12px] font-semibold text-ink-muted" htmlFor={searchId}>
              Tìm kiếm
            </label>
            <input
              id={searchId}
              className="control w-full"
              type="search"
              name="q"
              value={q}
              placeholder="VD: data center, khoi ngoai"
              onChange={(event) => setQ(event.target.value)}
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-[12px] font-semibold text-ink-muted" htmlFor={typeId}>
              Loại tin
            </label>
            <select
              id={typeId}
              className="control"
              value={filters.loai[0] ?? ""}
              onChange={(event) => {
                const value = event.target.value;
                apply({
                  ...filters,
                  q: q.trim(),
                  tu,
                  den,
                  loai: value ? [value as ApiNewsType] : [],
                });
              }}
            >
              <option value="">Tất cả</option>
              {FILTERABLE_NEWS_TYPES.map((type) => {
                const count = countOf(type);
                return (
                  <option key={type} value={type}>
                    {NEWS_TYPE_LABELS[type]}
                    {count !== undefined ? ` (${count})` : ""}
                  </option>
                );
              })}
            </select>
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-[12px] font-semibold text-ink-muted" htmlFor={fromId}>
              Từ ngày
            </label>
            <input
              id={fromId}
              className="control"
              type="date"
              value={tu}
              onChange={(event) => setTu(event.target.value)}
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-[12px] font-semibold text-ink-muted" htmlFor={toId}>
              Đến ngày
            </label>
            <input
              id={toId}
              className="control"
              type="date"
              value={den}
              onChange={(event) => setDen(event.target.value)}
            />
          </div>

          <button className="btn btn-strong" type="submit">
            Áp dụng
          </button>

          <button
            className="btn"
            type="button"
            disabled={!hasActiveFilters(filters)}
            onClick={() => {
              setQ("");
              setTu("");
              setDen("");
              apply(EMPTY_FILTERS);
            }}
          >
            Xoá lọc
          </button>

          <button className="btn" type="button" onClick={onCopy} disabled={copyDisabled}>
            Copy bản tin
          </button>

          <ThemeToggle />
        </form>

        <div className="mt-3 flex flex-wrap items-center gap-2">
          <span className="text-[12px] font-semibold text-ink-muted">Mã CK:</span>
          {popular.map((symbol) => {
            const selected = filters.ma.includes(symbol);
            return (
              <button
                key={symbol}
                type="button"
                className="chip"
                aria-pressed={selected}
                onClick={() => apply({ ...toggleSymbol(filters, symbol), q: q.trim(), tu, den })}
              >
                {symbol}
              </button>
            );
          })}

          <label className="sr-only" htmlFor={addSymbolId}>
            Thêm mã chứng khoán
          </label>
          <input
            id={addSymbolId}
            className="control w-[104px]"
            list={symbolListId}
            placeholder="Thêm mã"
            value={newSymbol}
            onChange={(event) => setNewSymbol(event.target.value.toUpperCase())}
            onKeyDown={(event) => {
              if (event.key !== "Enter") return;
              event.preventDefault();
              const symbol = normalizeSymbol(newSymbol);
              if (!symbol) return;
              setNewSymbol("");
              apply({ ...toggleSymbol(filters, symbol), q: q.trim(), tu, den });
            }}
          />
          <datalist id={symbolListId}>
            {tickers.map((ticker) => (
              <option key={ticker.symbol} value={ticker.symbol}>
                {ticker.short_name}
              </option>
            ))}
          </datalist>
        </div>

        <div className="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1 text-[12.5px] text-ink-muted">
          <button
            className="underline underline-offset-2 hover:text-accent disabled:no-underline disabled:opacity-60"
            type="button"
            disabled={filters.ma.length === 0}
            onClick={() => onSaveWatchlist(filters.ma)}
          >
            Lưu danh sách theo dõi
          </button>
          {watchlist.length > 0 ? (
            <>
              <span>
                Đang lưu: <span className="font-mono">{watchlist.join(", ")}</span>
              </span>
              <button
                className="underline underline-offset-2 hover:text-accent"
                type="button"
                onClick={onClearWatchlist}
              >
                Xoá danh sách theo dõi
              </button>
            </>
          ) : (
            <span>Chưa lưu danh sách theo dõi nào.</span>
          )}
          {isPending ? <span aria-live="polite">Đang áp dụng bộ lọc…</span> : null}
        </div>
      </div>
    </div>
  );
}
