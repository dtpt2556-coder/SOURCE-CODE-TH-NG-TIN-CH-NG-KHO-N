"use server";

import { fetchNews } from "@/lib/api";
import { PAGE_SIZE, type NewsFilters } from "@/lib/query";
import type { NewsItem } from "@/lib/types";

export interface LoadMoreResult {
  items: NewsItem[];
  hasMore: boolean;
  total: number;
}

/**
 * BA1 US-1.4: nút "Xem thêm" tải tiếp 20 tin, không reload trang.
 * Chạy qua Server Action để `API_BASE_URL` không bao giờ rời khỏi server.
 */
export async function loadMoreNews(
  filters: NewsFilters,
  offset: number,
): Promise<LoadMoreResult> {
  const safeOffset = Number.isFinite(offset) && offset > 0 ? Math.floor(offset) : 0;
  const { data } = await fetchNews(filters, safeOffset, PAGE_SIZE);

  return {
    items: data.data,
    hasMore: data.meta.has_more,
    total: data.meta.total,
  };
}
