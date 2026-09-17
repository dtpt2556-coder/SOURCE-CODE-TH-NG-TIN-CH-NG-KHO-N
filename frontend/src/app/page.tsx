import type { Metadata } from "next";

import { DigestView } from "@/components/DigestView";
import { SampleDataNotice, UpdatedAtBanner } from "@/components/Notices";
import { fetchNews, fetchSiteMeta, fetchTickers } from "@/lib/api";
import { buildHeadline } from "@/lib/digest";
import { formatDateTimeVN } from "@/lib/format";
import { NEWS_TYPE_LABELS } from "@/lib/news-types";
import {
  buildSearchParams,
  hasActiveFilters,
  parseFilters,
  type RawSearchParams,
} from "@/lib/query";
import { SITE_NAME } from "@/lib/site";

interface HomeProps {
  searchParams: Promise<RawSearchParams>;
}

export async function generateMetadata({ searchParams }: HomeProps): Promise<Metadata> {
  const filters = parseFilters(await searchParams);
  const canonicalQuery = buildSearchParams(filters);
  const canonical = canonicalQuery ? `/?${canonicalQuery}` : "/";

  const scope: string[] = [];
  if (filters.ma.length) scope.push(filters.ma.join(", "));
  if (filters.loai.length) {
    scope.push(filters.loai.map((type) => NEWS_TYPE_LABELS[type]).join(", "));
  }

  const title = scope.length
    ? `Tin chứng khoán ${scope.join(" · ")}`
    : "Bản tin chứng khoán Việt Nam hôm nay";

  const description = scope.length
    ? `Tin chứng khoán ${scope.join(" · ")} được tóm tắt cô đọng, giàu số liệu, kèm liên kết bài gốc.`
    : "Bảng tin chứng khoán Việt Nam tóm tắt cô đọng, giàu số liệu, gắn mã CK và loại tin, kèm liên kết về bài gốc.";

  return {
    title,
    description,
    alternates: { canonical },
    openGraph: { title: `${title} · ${SITE_NAME}`, description, url: canonical },
    // Trang đã lọc không cần index riêng để tránh loãng tín hiệu SEO.
    robots: hasActiveFilters(filters) ? { index: false, follow: true } : undefined,
  };
}

export default async function Home({ searchParams }: HomeProps) {
  const filters = parseFilters(await searchParams);

  const [news, tickers, meta] = await Promise.all([
    fetchNews(filters),
    fetchTickers(),
    fetchSiteMeta(),
  ]);

  const items = news.data.data;
  const headline = buildHeadline(news.data.meta.total, items);
  const usingFixture = news.fromFixture || meta.fromFixture;
  const caption = `Bảng tin chứng khoán tổng hợp, cập nhật ${formatDateTimeVN(meta.data.last_crawl_at)}`;

  return (
    <div className="pb-16">
      <DigestHeader
        headline={headline}
        lastCrawlAt={meta.data.last_crawl_at}
        isStaleFromApi={meta.data.is_stale}
        usingFixture={usingFixture}
      />

      <DigestView
        key={buildSearchParams(filters)}
        initialItems={items}
        initialMeta={news.data.meta}
        filters={filters}
        tickers={tickers.data}
        typeCounts={meta.data.news_types}
        caption={caption}
      />
    </div>
  );
}

interface DigestHeaderProps {
  headline: string;
  lastCrawlAt: string;
  isStaleFromApi: boolean;
  usingFixture: boolean;
}

function DigestHeader({
  headline,
  lastCrawlAt,
  isStaleFromApi,
  usingFixture,
}: DigestHeaderProps) {
  return (
    <div className="mx-auto max-w-content px-4 pb-8 pt-8 md:px-6 md:pt-12">
      {/* Màn hình B dùng serif cho toàn bộ nội dung - design-system.md §3 */}
      <h1 className="max-w-title font-serif text-[30px] font-bold leading-[1.2] text-ink">
        {headline}
      </h1>
      <div className="mt-3">
        <UpdatedAtBanner lastCrawlAt={lastCrawlAt} isStaleFromApi={isStaleFromApi} />
      </div>
      {usingFixture ? (
        <div className="mt-4 max-w-prose">
          <SampleDataNotice />
        </div>
      ) : null}
    </div>
  );
}
