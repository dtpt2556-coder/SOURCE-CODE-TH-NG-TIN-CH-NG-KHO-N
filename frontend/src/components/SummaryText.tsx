import { Fragment } from "react";

import { parseInlineMarkdown } from "@/lib/markdown";

interface SummaryTextProps {
  value: string;
  className?: string;
}

/**
 * Render `summary_md` an toàn: token → phần tử React, không đi qua chuỗi HTML.
 * Chỉ `<strong>` và `<em>` được sinh ra (BA1 US-1.2 AC2).
 */
export function SummaryText({ value, className }: SummaryTextProps) {
  const tokens = parseInlineMarkdown(value);

  if (process.env.NODE_ENV !== "production") {
    const strongCount = tokens.filter((token) => token.type === "strong").length;
    if (strongCount > 3) {
      // design-system.md B-2: đúng 1–3 cụm bold mỗi tóm tắt.
      console.warn(
        `[TenPoint] Tóm tắt có ${strongCount} cụm bold (> 3), điểm nhấn bị loãng: "${value.slice(0, 60)}"`,
      );
    }
  }

  return (
    <p className={className}>
      {tokens.map((token, index) => {
        if (token.type === "strong") {
          return <strong key={index}>{token.value}</strong>;
        }
        if (token.type === "em") {
          return <em key={index}>{token.value}</em>;
        }
        return <Fragment key={index}>{token.value}</Fragment>;
      })}
    </p>
  );
}
