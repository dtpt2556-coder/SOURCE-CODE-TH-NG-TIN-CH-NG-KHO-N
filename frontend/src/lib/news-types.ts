import { API_NEWS_TYPES, NEWS_TYPES, type ApiNewsType, type NewsType } from "./types";

/** Nhan hien thi - design-system v2 muc 5.4. Text thuan, khong badge, khong cham mau. */
export const NEWS_TYPE_LABELS: Record<NewsType, string> = {
  vi_mo: "Vĩ mô",
  nganh: "Ngành",
  doanh_nghiep: "Doanh nghiệp",
  thi_truong: "Thị trường",
  khoi_ngoai: "Khối ngoại",
  co_tuc_phat_hanh: "Cổ tức/Phát hành",
  phap_ly: "Pháp lý",
  phan_tich: "Phân tích",
  trai_phieu_tin_dung: "Trái phiếu/Tín dụng",
};

/** Cac loai tin duoc phep gui len `GET /api/v1/news?loai=` (hop dong API muc 6). */
export const FILTERABLE_NEWS_TYPES: readonly ApiNewsType[] = API_NEWS_TYPES;

export function isNewsType(value: string): value is NewsType {
  return (NEWS_TYPES as readonly string[]).includes(value);
}

export function isFilterableNewsType(value: string): value is ApiNewsType {
  return (API_NEWS_TYPES as readonly string[]).includes(value);
}

export function newsTypeLabel(value: string, fallback?: string): string {
  if (isNewsType(value)) return NEWS_TYPE_LABELS[value];
  return fallback ?? value;
}
