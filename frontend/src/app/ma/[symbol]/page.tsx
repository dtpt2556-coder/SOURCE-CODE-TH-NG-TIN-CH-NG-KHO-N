import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";

import { DigestTable } from "@/components/DigestTable";
import { EmptyState, SampleDataNotice } from "@/components/Notices";
import { fetchTickerDetail } from "@/lib/api";
import { formatNumberVN } from "@/lib/format";
import { normalizeSymbol } from "@/lib/query";
import { SITE_NAME } from "@/lib/site";

interface TickerPageProps {
  params: Promise<{ symbol: string }>;
}

export async function generateMetadata({ params }: TickerPageProps): Promise<Metadata> {
  const { symbol } = await params;
  const upper = normalizeSymbol(decodeURIComponent(symbol));
  const detail = await fetchTickerDetail(upper);

  if (!detail) {
    return { title: `Không tìm thấy mã ${upper}`, robots: { index: false, follow: false } };
  }

  const { ticker } = detail.data;
  const title = `${ticker.symbol} · ${ticker.company_name}`;
  const description = `Tin tức và nhận định mới nhất về ${ticker.symbol} (${ticker.company_name}), sàn ${ticker.exchange}, ngành ${ticker.sector}. Tóm tắt giàu số liệu, kèm liên kết bài gốc.`;

  return {
    title,
    description,
    alternates: { canonical: `/ma/${ticker.symbol}` },
    openGraph: {
      type: "website",
      title: `${title} · ${SITE_NAME}`,
      description,
      url: `/ma/${ticker.symbol}`,
    },
  };
}

export default async function TickerPage({ params }: TickerPageProps) {
  const { symbol } = await params;
  const upper = normalizeSymbol(decodeURIComponent(symbol));
  const detail = await fetchTickerDetail(upper);

  if (!detail) notFound();

  const { ticker, news, notes } = detail.data;

  return (
    <div className="mx-auto max-w-content px-4 pb-16 pt-8 md:px-6 md:pt-12">
      <nav aria-label="Đường dẫn" className="text-[13px] text-ink-muted">
        <Link href="/" className="underline underline-offset-2 hover:text-accent">
          Bản tin
        </Link>
        <span aria-hidden="true"> / </span>
        <span>{ticker.symbol}</span>
      </nav>

      <header className="mt-4">
        <h1 className="max-w-title font-serif text-[30px] font-bold leading-[1.2] text-ink">
          <span className="font-mono tracking-[0.04em]">{ticker.symbol}</span>{" · "}
          {ticker.company_name}
        </h1>
        <dl className="mt-3 flex flex-wrap gap-x-8 gap-y-1 text-[13.5px] text-ink-muted">
          <div className="flex gap-2">
            <dt>Sàn:</dt>
            <dd className="text-ink">{ticker.exchange}</dd>
          </div>
          <div className="flex gap-2">
            <dt>Ngành:</dt>
            <dd className="text-ink">{ticker.sector}</dd>
          </div>
          <div className="flex gap-2">
            <dt>Rổ VN30:</dt>
            <dd className="text-ink">{ticker.in_vn30 ? "Có" : "Không"}</dd>
          </div>
        </dl>
      </header>

      {detail.fromFixture ? (
        <div className="mt-4 max-w-prose">
          <SampleDataNotice />
        </div>
      ) : null}

      <section className="mt-10" aria-labelledby="dong-thoi-gian">
        <h2
          id="dong-thoi-gian"
          className="font-sans text-[21px] font-bold uppercase leading-[1.3] tracking-[0.06em] text-ink"
        >
          Dòng thời gian tin
        </h2>
        <p className="mt-2 text-[13.5px] text-ink-muted">
          {formatNumberVN(news.length)} tin gần nhất có gắn mã {ticker.symbol}.
        </p>
        <div className="mt-6">
          {news.length ? (
            <DigestTable
              items={news}
              caption={`Dòng thời gian tin của mã ${ticker.symbol}`}
            />
          ) : (
            <EmptyState
              title={`Chưa có tin nào gắn mã ${ticker.symbol}.`}
              hint="Dữ liệu được cập nhật 3 lần mỗi ngày."
            />
          )}
        </div>
      </section>

      <section className="mt-16" aria-labelledby="nhan-dinh-lien-quan">
        <h2
          id="nhan-dinh-lien-quan"
          className="font-sans text-[21px] font-bold uppercase leading-[1.3] tracking-[0.06em] text-ink"
        >
          Nhận định liên quan
        </h2>
        {notes.length ? (
          <ul className="mt-6 border-t border-rule-strong">
            {notes.map((note) => (
              <li key={note.slug} className="border-b border-rule py-5">
                <p className="text-[13.5px] text-ink-muted">
                  <time dateTime={note.published_at}>{note.published_date_vn}</time>
                  <span aria-hidden="true"> · </span>
                  <span className="font-mono tracking-[0.04em]">{note.symbol}</span>
                </p>
                <h3 className="mt-1 max-w-prose font-serif text-[18px] font-bold leading-[1.4]">
                  <Link
                    href={`/nhan-dinh/${note.slug}`}
                    className="no-underline hover:underline"
                  >
                    {note.title}
                  </Link>
                </h3>
                {note.excerpt ? (
                  <p className="mt-2 max-w-prose font-serif text-[17px] leading-[1.62] text-ink-2">
                    {note.excerpt}
                  </p>
                ) : null}
              </li>
            ))}
          </ul>
        ) : (
          <p className="mt-4 text-[14px] text-ink-muted">
            Chưa có bài nhận định nào cho mã này.
          </p>
        )}
      </section>
    </div>
  );
}
