import { displayDomain } from "@/lib/format";
import type { NewsItem } from "@/lib/types";

interface SourceLinkProps {
  item: NewsItem;
}

/**
 * design-system v2 muc 5.3 - da sua theo phan hoi khach hang.
 *
 * Van chi HIEN THI domain rut gon (giu dung anh tham chieu), nhung bo sung
 * ngu canh "bam vao se ra bai nao" qua `title`, `aria-label` va mot span
 * sr-only ghi tieu de bai goc.
 *
 * `short_link_alive === false` thi khong render thanh <a> nua: hien text mo
 * kem chu thich, tranh dua nguoi dung toi mot lien ket da chet.
 */
export function SourceLink({ item }: SourceLinkProps) {
  const { source, short_link: shortLink, short_link_alive: alive } = item;
  const domain = displayDomain(source.domain);
  const articleTitle = source.article_title?.trim() || item.title;
  const publishedDate = item.published_date_vn;

  // Hai ly do khac nhau, phai noi ro ly do nao:
  //  - is_demo: 5 tin dung lai tu anh tham chieu, URL khong co that.
  //  - alive === false: bai that nhung nguon da go xuong.
  // Ca hai deu KHONG duoc render thanh <a>. Mot link dep ma bam vao ra 404
  // pha huy dung thu san pham nay ban: do tin cay cua trich dan.
  if (item.is_demo || alive === false) {
    const note = item.is_demo
      ? "Dữ liệu mẫu, chưa có bài gốc"
      : "Bài gốc không còn khả dụng";

    return (
      <span className="source-dead">
        {domain}
        <span className="source-dead-note">{note}</span>
        <span className="sr-only">, bài: {articleTitle}</span>
      </span>
    );
  }

  return (
    <a
      className="source-link"
      href={shortLink}
      target="_blank"
      rel="nofollow noopener noreferrer"
      title={`${articleTitle} · ${domain} · ${publishedDate}`}
      aria-label={`Đọc bài gốc: ${articleTitle}. Nguồn ${source.name}, đăng ${publishedDate} (mở tab mới)`}
    >
      {domain}
      <span className="sr-only">, bài: {articleTitle}</span>
    </a>
  );
}
