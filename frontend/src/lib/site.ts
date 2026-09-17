export const SITE_NAME = "TenPoint";

/**
 * URL công khai - dùng cho canonical, Open Graph, sitemap, robots và short link
 * tuyệt đối khi copy bản tin. Khớp biến `PUBLIC_BASE_URL` của backend
 * (architecture.md §12).
 */
export const SITE_URL = (
  process.env.NEXT_PUBLIC_SITE_URL ?? "http://localhost:3000"
).replace(/\/$/, "");

export function absoluteUrl(path: string): string {
  return `${SITE_URL}${path.startsWith("/") ? path : `/${path}`}`;
}
