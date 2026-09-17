import { formatDateTimeVN, isStale } from "@/lib/format";

interface UpdatedAtBannerProps {
  lastCrawlAt: string;
  /** BE co the tu danh dau; FE van kiem tra lai moc 12h. */
  isStaleFromApi?: boolean;
}

/** muc 5.7 + BA1 US-1.3 */
export function UpdatedAtBanner({ lastCrawlAt, isStaleFromApi }: UpdatedAtBannerProps) {
  const stale = isStaleFromApi || isStale(lastCrawlAt);

  return (
    <div className="text-[13.5px] leading-[1.5] text-ink-muted">
      <p>
        Cập nhật lần cuối:{" "}
        <time className="font-mono tracking-[0.02em]" dateTime={lastCrawlAt}>
          {formatDateTimeVN(lastCrawlAt)}
        </time>
      </p>
      {stale ? (
        // Khong dung mau gia o day: trong san pham chung khoan VN, chu do doc ra
        // "giam gia" chu khong phai "du lieu cu". Phan biet bang weight + ky hieu
        // + ke trai, nen dat luon WCAG 1.4.1.
        <p className="mt-1 border-l-2 border-accent pl-2 font-semibold text-ink">
          <span aria-hidden="true">⚠ </span>
          Cảnh báo: dữ liệu có thể chưa mới nhất.
        </p>
      ) : null}
    </div>
  );
}

/** Banner nho khi API chua san sang va trang dang render bang du lieu mau. */
export function SampleDataNotice() {
  return (
    <p
      role="status"
      className="border border-rule bg-surface-raised px-3 py-2 text-[13px] leading-[1.5] text-ink-muted"
    >
      Đang dùng dữ liệu mẫu: chưa kết nối được tới API nội bộ.
    </p>
  );
}

interface EmptyStateProps {
  title: string;
  hint?: string;
}

/** BA1 US-1.1 AC3: khong hien bang rong. */
export function EmptyState({ title, hint }: EmptyStateProps) {
  return (
    <div className="border-t border-rule-strong py-16 text-center">
      <p className="font-serif text-[17px] leading-[1.62] text-ink">{title}</p>
      {hint ? <p className="mt-2 text-[14px] text-ink-muted">{hint}</p> : null}
    </div>
  );
}
