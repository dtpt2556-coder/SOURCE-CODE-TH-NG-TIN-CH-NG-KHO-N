import Link from "next/link";
import { Fragment } from "react";

import type { TickerRef } from "@/lib/types";

const MAX_VISIBLE = 8;

interface TickerLinkProps {
  symbol: string;
}

/** §5.2 - mã CK luôn là link tới `/ma/{mã}`. */
export function TickerLink({ symbol }: TickerLinkProps) {
  return (
    <Link className="ticker-link" href={`/ma/${encodeURIComponent(symbol)}`}>
      {symbol}
    </Link>
  );
}

interface TickerListProps {
  tickers: TickerRef[];
}

/**
 * Danh sách mã, nối bằng `, `, cho phép wrap.
 * > 8 mã: hiện 8 mã đầu + `+N` mở rộng bằng `<details>` (không cần JS).
 */
export function TickerList({ tickers }: TickerListProps) {
  if (!tickers.length) {
    return <span className="ticker-cell text-ink-faint">-</span>;
  }

  const visible = tickers.slice(0, MAX_VISIBLE);
  const hidden = tickers.slice(MAX_VISIBLE);

  return (
    <span className="ticker-cell">
      {visible.map((ticker, index) => (
        <Fragment key={ticker.symbol}>
          {index > 0 ? ", " : null}
          <TickerLink symbol={ticker.symbol} />
        </Fragment>
      ))}
      {hidden.length > 0 ? (
        <>
          {", "}
          <details className="ticker-more">
            <summary aria-label={`Hiện thêm ${hidden.length} mã chứng khoán`}>
              +{hidden.length}
            </summary>
            {hidden.map((ticker, index) => (
              <Fragment key={ticker.symbol}>
                {index > 0 ? ", " : " "}
                <TickerLink symbol={ticker.symbol} />
              </Fragment>
            ))}
          </details>
        </>
      ) : null}
    </span>
  );
}
