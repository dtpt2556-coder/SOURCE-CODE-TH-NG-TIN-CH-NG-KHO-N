import type { Metadata } from "next";
import Link from "next/link";

import { EmptyState, SampleDataNotice } from "@/components/Notices";
import { fetchNotes } from "@/lib/api";
import { SITE_NAME } from "@/lib/site";

export const metadata: Metadata = {
  title: "Nhận định & luận điểm đầu tư",
  description:
    "Các bài nhận định dạng luận điểm đầu tư theo từng mã chứng khoán: bức tranh dài hạn, số liệu cụ thể, nguồn tham khảo rõ ràng.",
  alternates: { canonical: "/nhan-dinh" },
  openGraph: {
    type: "website",
    title: `Nhận định & luận điểm đầu tư · ${SITE_NAME}`,
    description:
      "Các bài nhận định dạng luận điểm đầu tư theo từng mã chứng khoán, có số liệu và nguồn tham khảo.",
    url: "/nhan-dinh",
  },
};

export default async function NotesPage() {
  const notes = await fetchNotes();
  const items = notes.data.data;

  return (
    <div className="mx-auto max-w-content px-4 pb-16 pt-8 md:px-6 md:pt-12">
      <h1 className="max-w-title font-sans text-[30px] font-bold leading-[1.2] text-ink">
        Nhận định &amp; luận điểm đầu tư
      </h1>
      <p className="mt-3 max-w-prose font-serif text-[17px] leading-[1.62] text-ink-2">
        Mỗi bài là một tập luận điểm có đánh số, bám vào số liệu công bố và luôn kèm
        khuyến cáo ở cuối bài.
      </p>

      {notes.fromFixture ? (
        <div className="mt-4 max-w-prose">
          <SampleDataNotice />
        </div>
      ) : null}

      <div className="mt-10">
        {items.length ? (
          <ul className="border-t border-rule-strong">
            {items.map((note) => (
              <li key={note.slug} className="border-b border-rule py-6">
                <p className="text-[13.5px] text-ink-muted">
                  <time dateTime={note.published_at}>{note.published_date_vn}</time>
                  <span aria-hidden="true"> · </span>
                  <Link
                    href={`/ma/${note.symbol}`}
                    className="ticker-link font-mono tracking-[0.04em] text-ink"
                  >
                    {note.symbol}
                  </Link>
                </p>
                <h2 className="mt-2 max-w-title font-sans text-[21px] font-extrabold uppercase leading-[1.25] tracking-[-0.01em]">
                  <Link
                    href={`/nhan-dinh/${note.slug}`}
                    className="no-underline hover:underline"
                  >
                    {note.title}
                  </Link>
                </h2>
                {note.excerpt ? (
                  <p className="mt-3 max-w-prose font-serif text-[17px] leading-[1.62] text-ink-2">
                    {note.excerpt}
                  </p>
                ) : null}
              </li>
            ))}
          </ul>
        ) : (
          <EmptyState
            title="Chưa có bài nhận định nào."
            hint="Các bài nhận định được đăng theo từng mã sau mỗi kỳ công bố kết quả kinh doanh."
          />
        )}
      </div>
    </div>
  );
}
