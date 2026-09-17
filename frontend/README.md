# TenPoint — Frontend (Next.js)

Giao diện web cho nền tảng tổng hợp & tóm tắt tin chứng khoán Việt Nam.

Nguồn sự thật thiết kế: `docs/04-design/design-system.md`.
Hợp đồng API: `docs/03-sa/architecture.md` §6.

## Chạy

```bash
npm run dev        # http://localhost:3000
npm run build && npm run start
```

```bash
npm run typecheck  # tsc --noEmit
npm run lint       # eslint (Next 16 đã bỏ lệnh `next lint`)
```

## Biến môi trường

Xem `.env.example`.

| Biến | Mặc định | Ghi chú |
|---|---|---|
| `API_BASE_URL` | `http://localhost:8080` | Chỉ đọc phía server. Trong Docker Compose dùng `http://api:8080`. |
| `API_TIMEOUT_MS` | `4000` | Quá hạn → trang render bằng dữ liệu mẫu. |
| `NEXT_PUBLIC_SITE_URL` | `http://localhost:3000` | Canonical / Open Graph / sitemap / robots. |

## Trang

| Route | Nội dung | Render |
|---|---|---|
| `/` | News Digest — bảng 5 cột, FilterBar, "Xem thêm", "Copy bản tin" | SSR |
| `/ma/[symbol]` | Thông tin mã, dòng thời gian tin, nhận định liên quan | SSR |
| `/nhan-dinh` | Danh sách research note | SSR |
| `/nhan-dinh/[slug]` | Research Note + JSON-LD `NewsArticle` | SSR |
| `/sitemap.xml`, `/robots.txt` | SEO | Metadata route |

## Ghi chú kỹ thuật

- **Fallback fixture** — khi không gọi được API, `src/lib/api.ts` trả dữ liệu mẫu
  từ `src/lib/fixtures.ts` và trang hiện banner "Đang dùng dữ liệu mẫu". Không
  bao giờ crash, nhờ vậy `npm run build` và demo luôn chạy được khi backend chưa
  sẵn sàng.
- **Sanitize `summary_md`** — `src/lib/markdown.ts` parse markdown inline thành
  token rồi React render ra `<strong>`/`<em>`. Không đi qua chuỗi HTML nào, nên
  không cần `dangerouslySetInnerHTML` (BA1 US-1.2 AC2).
- **Bộ lọc nằm trong query string** (`/?ma=FPT,HPG&loai=nganh&q=…`) để copy URL
  là giữ nguyên trạng thái (BA1 US-2.1 AC3).
- **Tailwind v4** — không có `tailwind.config.ts`; token khai báo bằng `@theme`
  trong `src/app/globals.css`, cùng các biến `:root` theo design-system §2.
