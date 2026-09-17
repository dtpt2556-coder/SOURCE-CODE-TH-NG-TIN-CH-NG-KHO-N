import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";

import { DigestTable } from "@/components/DigestTable";
import { SampleDataNotice } from "@/components/Notices";
import { fetchNote, fetchTickerDetail } from "@/lib/api";
import { truncate } from "@/lib/markdown";
import { absoluteUrl, SITE_NAME } from "@/lib/site";
import type { ResearchNote } from "@/lib/types";

interface NotePageProps {
  params: Promise<{ slug: string }>;
}

export async function generateMetadata({ params }: NotePageProps): Promise<Metadata> {
  const { slug } = await params;
  const note = await fetchNote(slug);

  if (!note) {
    return { title: "Không tìm thấy bài nhận định", robots: { index: false, follow: false } };
  }

  const description = truncate(note.data.excerpt || note.data.points[0]?.lead || note.data.title);

  return {
    title: note.data.title,
    description,
    alternates: { canonical: `/nhan-dinh/${note.data.slug}` },
    openGraph: {
      type: "article",
      title: note.data.title,
      description,
      url: `/nhan-dinh/${note.data.slug}`,
      publishedTime: note.data.published_at,
      siteName: SITE_NAME,
      locale: "vi_VN",
    },
  };
}

/** JSON-LD NewsArticle - BA1 §7 (SEO). */
function buildJsonLd(note: ResearchNote): string {
  const payload = {
    "@context": "https://schema.org",
    "@type": "NewsArticle",
    headline: note.title,
    datePublished: note.published_at,
    dateModified: note.published_at,
    inLanguage: "vi-VN",
    description: note.excerpt,
    mainEntityOfPage: {
      "@type": "WebPage",
      "@id": absoluteUrl(`/nhan-dinh/${note.slug}`),
    },
    author: { "@type": "Organization", name: SITE_NAME },
    publisher: { "@type": "Organization", name: SITE_NAME },
    about: note.related_symbols.map((symbol) => ({
      "@type": "Corporation",
      tickerSymbol: symbol,
    })),
  };

  // Escape `<` để chuỗi JSON không bao giờ đóng sớm thẻ <script>.
  return JSON.stringify(payload).replace(/</g, "\\u003c");
}

export default async function NotePage({ params }: NotePageProps) {
  const { slug } = await params;
  const note = await fetchNote(slug);

  if (!note) notFound();

  const data = note.data;
  const related = await fetchTickerDetail(data.symbol);
  const relatedNews = related?.data.news.slice(0, 5) ?? [];

  return (
    <article className="mx-auto max-w-content px-4 pb-16 pt-8 md:px-6 md:pt-12">
      <script
        type="application/ld+json"
        // Nội dung là JSON do chính FE sinh ra, đã escape `<`.
        dangerouslySetInnerHTML={{ __html: buildJsonLd(data) }}
      />

      <nav aria-label="Đường dẫn" className="text-[13px] text-ink-muted">
        <Link href="/nhan-dinh" className="underline underline-offset-2 hover:text-accent">
          Nhận định
        </Link>
        <span aria-hidden="true"> / </span>
        <Link
          href={`/ma/${data.symbol}`}
          className="font-mono tracking-[0.04em] underline underline-offset-2 hover:text-accent"
        >
          {data.symbol}
        </Link>
      </nav>

      <header className="mt-6">
        <h1 className="note-title">{data.title}</h1>
        <p className="mt-5 text-[13.5px] text-ink-muted">
          <time dateTime={data.published_at}>{data.published_date_vn}</time>
          <span aria-hidden="true"> · </span>
          <span>Mã liên quan: </span>
          {data.related_symbols.map((symbol, index) => (
            <span key={symbol}>
              {index > 0 ? ", " : null}
              <Link
                href={`/ma/${symbol}`}
                className="ticker-link font-mono tracking-[0.04em] text-ink"
              >
                {symbol}
              </Link>
            </span>
          ))}
        </p>
      </header>

      {note.fromFixture ? (
        <div className="mt-6 max-w-prose">
          <SampleDataNotice />
        </div>
      ) : null}

      <section className="mt-10" aria-labelledby="luan-diem">
        <h2 id="luan-diem" className="note-section">
          {data.section_heading || "LUẬN ĐIỂM ĐẦU TƯ"}
        </h2>

        <ol className="note-points">
          {data.points.map((point) => (
            <li className="note-point" key={point.ordinal}>
              {/* Số đánh inline ngay trong dòng dẫn - đúng ảnh tham chiếu §5.5 */}
              <h3 className="note-lead">
                {point.ordinal}. {point.lead}
              </h3>
              <p className="note-body">{point.body}</p>
            </li>
          ))}
        </ol>
      </section>

      <hr className="mt-12 border-0 border-t border-rule" />
      <p className="mt-6 max-w-prose text-[13.5px] leading-[1.5] text-ink-muted">
        {data.disclaimer ||
          "Nội dung chỉ mang tính thông tin, không phải khuyến nghị đầu tư."}
      </p>

      {relatedNews.length ? (
        <section className="mt-16" aria-labelledby="tin-lien-quan">
          <h2
            id="tin-lien-quan"
            className="font-sans text-[21px] font-bold uppercase leading-[1.3] tracking-[0.06em] text-ink"
          >
            Tin gần đây của {data.symbol}
          </h2>
          <div className="mt-6">
            <DigestTable
              items={relatedNews}
              caption={`Tin gần đây của mã ${data.symbol}`}
            />
          </div>
        </section>
      ) : null}
    </article>
  );
}
