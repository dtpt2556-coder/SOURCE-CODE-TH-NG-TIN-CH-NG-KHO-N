/**
 * Sanitize + parse markdown inline cho cột "Tóm tắt thông tin".
 *
 * BA1 US-1.2 AC2 / design-system.md §3.3 rule B-1:
 *   chỉ `<strong>` và `<em>` được phép - mọi HTML khác bị loại bỏ.
 *
 * Cách tiếp cận: KHÔNG sinh chuỗi HTML nào cả. Hàm này trả về danh sách token
 * thuần dữ liệu; React render token thành phần tử `<strong>`/`<em>`/text node.
 * Vì không có bước "chuỗi HTML", `dangerouslySetInnerHTML` là không cần thiết
 * và XSS qua summary_md là bất khả thi về mặt cấu trúc.
 */

export type InlineToken =
  | { type: "text"; value: string }
  | { type: "strong"; value: string }
  | { type: "em"; value: string };

/** Loại bỏ mọi thẻ giống HTML trước khi parse markdown. */
const HTML_TAG = /<[^>]*>/g;

/** `**đậm**` · `__đậm__` · `*nghiêng*` · `_nghiêng_` */
const INLINE_PATTERN =
  /\*\*([^\n*](?:[^\n]*?[^\n*])?)\*\*|__([^\n_](?:[^\n]*?[^\n_])?)__|\*([^\n*]+?)\*|_([^\n_]+?)_/g;

function stripHtml(input: string): string {
  return input.replace(HTML_TAG, "");
}

export function parseInlineMarkdown(input: string | null | undefined): InlineToken[] {
  if (!input) return [];

  const source = stripHtml(input);
  const tokens: InlineToken[] = [];
  let cursor = 0;

  for (const match of source.matchAll(INLINE_PATTERN)) {
    const index = match.index ?? 0;
    if (index > cursor) {
      tokens.push({ type: "text", value: source.slice(cursor, index) });
    }

    const strong = match[1] ?? match[2];
    if (strong !== undefined) {
      tokens.push({ type: "strong", value: strong });
    } else {
      tokens.push({ type: "em", value: (match[3] ?? match[4]) as string });
    }

    cursor = index + match[0].length;
  }

  if (cursor < source.length) {
    tokens.push({ type: "text", value: source.slice(cursor) });
  }

  return tokens;
}

/** Số cụm được bold - design-system.md B-2: đúng 1–3 cụm mỗi tóm tắt. */
export function countStrongSpans(input: string | null | undefined): number {
  return parseInlineMarkdown(input).filter((token) => token.type === "strong").length;
}

/** Text thuần (bỏ dấu markdown) - dùng cho "Copy bản tin" và thẻ meta description. */
export function toPlainText(input: string | null | undefined): string {
  return parseInlineMarkdown(input)
    .map((token) => token.value)
    .join("")
    .replace(/\s+/g, " ")
    .trim();
}

/** Cắt gọn text cho `<meta name="description">`. */
export function truncate(input: string, max = 160): string {
  const text = input.trim();
  if (text.length <= max) return text;
  return `${text.slice(0, max - 1).trimEnd()}…`;
}
