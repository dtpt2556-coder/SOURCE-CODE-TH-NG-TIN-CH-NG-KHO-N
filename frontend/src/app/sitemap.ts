import type { MetadataRoute } from "next";

import { fetchNotes, fetchTickers } from "@/lib/api";
import { absoluteUrl } from "@/lib/site";

/** Luôn dựng lại khi build/deploy; danh mục mã và note thay đổi chậm. */
export const revalidate = 3600;

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const [tickers, notes] = await Promise.all([fetchTickers(), fetchNotes(1000)]);
  const now = new Date();

  const staticRoutes: MetadataRoute.Sitemap = [
    { url: absoluteUrl("/"), lastModified: now, changeFrequency: "hourly", priority: 1 },
    {
      url: absoluteUrl("/nhan-dinh"),
      lastModified: now,
      changeFrequency: "daily",
      priority: 0.8,
    },
  ];

  const tickerRoutes: MetadataRoute.Sitemap = tickers.data.map((ticker) => ({
    url: absoluteUrl(`/ma/${ticker.symbol}`),
    lastModified: now,
    changeFrequency: "daily",
    priority: 0.7,
  }));

  const noteRoutes: MetadataRoute.Sitemap = notes.data.data.map((note) => ({
    url: absoluteUrl(`/nhan-dinh/${note.slug}`),
    lastModified: new Date(note.published_at),
    changeFrequency: "monthly",
    priority: 0.6,
  }));

  return [...staticRoutes, ...tickerRoutes, ...noteRoutes];
}
