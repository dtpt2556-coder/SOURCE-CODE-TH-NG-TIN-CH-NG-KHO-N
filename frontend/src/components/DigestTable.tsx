import { SourceLink } from "@/components/SourceLink";
import { SummaryText } from "@/components/SummaryText";
import { TickerList } from "@/components/TickerList";
import { formatDateTimeVN } from "@/lib/format";
import { newsTypeLabel } from "@/lib/news-types";
import type { NewsItem } from "@/lib/types";

interface DigestTableProps {
  items: NewsItem[];
  /** `<caption>` sr-only - a11y checklist §8. */
  caption: string;
}

/**
 * Bảng tin tổng hợp - design-system.md §5.1.
 * Hairline ngăn hàng, kẻ đậm dưới header, KHÔNG zebra, KHÔNG border dọc,
 * KHÔNG bo góc, KHÔNG shadow. Responsive do CSS đảm nhiệm (globals.css):
 * <768px → card dọc · 768–1023px → 4 cột · ≥1024px → 5 cột.
 *
 * `role` khai báo tường minh để giữ ngữ nghĩa bảng trong cây accessibility
 * kể cả khi CSS đổi `display` ở breakpoint mobile.
 */
export function DigestTable({ items, caption }: DigestTableProps) {
  return (
    <table className="digest-table" role="table">
      <caption className="sr-only">{caption}</caption>
      <thead role="rowgroup">
        <tr role="row">
          <th role="columnheader" scope="col" data-col="date">
            Ngày
          </th>
          <th role="columnheader" scope="col" data-col="tickers">
            Mã CK
          </th>
          <th role="columnheader" scope="col" data-col="summary">
            Tóm tắt thông tin
          </th>
          <th role="columnheader" scope="col" data-col="source">
            Source
          </th>
          <th role="columnheader" scope="col" data-col="type">
            Loại tin
          </th>
        </tr>
      </thead>
      <tbody role="rowgroup">
        {items.map((item) => {
          const label = newsTypeLabel(item.news_type, item.news_type_label);
          return (
            <tr role="row" key={item.id}>
              <td role="cell" data-col="date">
                <time className="digest-date" dateTime={item.published_at}>
                  {item.published_at_estimated ? "~" : ""}
                  {item.published_date_vn}
                </time>
                {item.published_at_estimated ? (
                  <span className="sr-only"> (thời gian đăng là ước tính)</span>
                ) : null}
                <span className="digest-type-inline"> · {label}</span>
                {item.updated_at ? (
                  <span className="digest-updated">
                    Đã cập nhật {formatDateTimeVN(item.updated_at)}
                  </span>
                ) : null}
              </td>
              <td role="cell" data-col="tickers">
                <TickerList tickers={item.tickers} />
              </td>
              <td role="cell" data-col="summary">
                <SummaryText className="digest-summary" value={item.summary_md} />
              </td>
              <td role="cell" data-col="source">
                <SourceLink item={item} />
              </td>
              <td role="cell" data-col="type">
                <span className="digest-type">{label}</span>
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}
