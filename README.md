# TenPoint

**Nền tảng tổng hợp và chưng cất tin tức chứng khoán Việt Nam.**

TenPoint thu thập tin từ các nguồn uy tín 8 tiếng một lần (06:00 · 14:00 · 22:00 giờ VN),
tự động gắn mã chứng khoán, phân loại tin, rồi rút ra **một bảng digest giàu số liệu** —
để nhà đầu tư nắm được diễn biến thị trường trong vài phút thay vì lướt 6 tờ báo.

Sản phẩm **không** lưu và **không** hiển thị công khai nội dung gốc của bài báo. Đầu ra là
bản tóm tắt do hệ thống sinh ra + metadata + short link trả người đọc về bài gốc
(xem [ADR-006](docs/03-sa/architecture.md#10-quyết-định-kiến-trúc-adr)).

---

## Kiến trúc

```
                         ┌──────────────────────┐
     Trình duyệt  ──────►│        caddy         │  :80 / :443
     (mobile-first)      │  TLS · gzip/zstd     │  reverse proxy
                         │  security headers    │
                         └───────┬──────────┬───┘
                 / , /ma/* ,     │          │   /api/* , /r/*
                 /nhan-dinh/*    │          │
                         ┌───────▼──────┐   │
                         │     web      │   │
                         │  Next.js SSR │   │
                         │  Node 22     │   │
                         │  :3000       │   │
                         └───────┬──────┘   │
                                 │ SSR fetch (server-side, trong docker network)
                                 │ http://api:8080
                         ┌───────▼──────────▼───┐          ┌────────────────┐
                         │         api          │─────────►│  Nguồn tin VN  │
                         │  Go 1.26 · chi · pgx │  8h/lần  │  (RSS + HTML)  │
                         │  ┣ HTTP API v1       │          └────────────────┘
                         │  ┣ Cron scheduler    │          ┌────────────────┐
                         │  ┗ Ingest pipeline   │─────────►│   Claude API   │
                         │  :8080               │ tuỳ chọn │  (tóm tắt)     │
                         └───────────┬──────────┘          └────────────────┘
                                     │
                         ┌───────────▼──────────┐
                         │          db          │
                         │  PostgreSQL 16       │
                         │  unaccent · pg_trgm  │
                         │  volume: pgdata      │
                         └──────────────────────┘
```

| Container | Công nghệ | Trách nhiệm |
|---|---|---|
| `caddy` | Caddy 2 | TLS tự động, reverse proxy, nén, security headers |
| `web`   | Next.js App Router (SSR) | Render HTML cho SEO, không giữ state nghiệp vụ |
| `api`   | Go 1.26 + chi + pgx | REST API + cron + ingest + short link redirect |
| `db`    | PostgreSQL 16 | Nguồn sự thật duy nhất |

Chi tiết: [`docs/03-sa/architecture.md`](docs/03-sa/architecture.md).

---

## Yêu cầu hệ thống

| Thành phần | Phiên bản tối thiểu | Ghi chú |
|---|---|---|
| Docker Engine | 24+ | cần BuildKit (mặc định từ 23) |
| Docker Compose | v2.24+ | dùng `docker compose`, không phải `docker-compose` |
| RAM | 2 GB trống | 4 GB nếu chạy cùng lúc cả build |
| Ổ đĩa | 5 GB trống | image + volume Postgres |
| `make` | có sẵn trên macOS/Linux | Windows: dùng WSL2 |

Chỉ cần Docker. **Không** cần cài Go hay Node trên máy để chạy hệ thống
(chỉ cần khi muốn chạy `make test-be` / `make test-fe` bằng toolchain local).

---

## Chạy trong 3 lệnh

```bash
cp deploy/.env.example deploy/.env    # 1. tạo cấu hình (mặc định chạy được ngay)
make build                            # 2. build image api + web
make dev                              # 3. khởi động — mở http://localhost:3000
```

Sau bước 3:

| Địa chỉ | Là gì |
|---|---|
| <http://localhost> | Trang chính qua Caddy (đúng đường đi production) |
| <http://localhost:3000> | Next.js trực tiếp (bỏ qua proxy, tiện debug) |
| <http://localhost:8080/healthz> | Liveness của API |
| <http://localhost:8080/readyz> | Readiness — `503` nếu dữ liệu quá cũ |
| `localhost:5432` | PostgreSQL (chỉ bind `127.0.0.1`) |

> Hệ thống chạy được **ngay, không cần API key nào**: summarizer mặc định là
> `extractive` (heuristic chọn câu giàu số liệu). Muốn tóm tắt chất lượng cao hơn,
> đặt `SUMMARIZER_PROVIDER=anthropic` + `ANTHROPIC_API_KEY` trong `deploy/.env`.

Chạy ở chế độ production (không mở cổng debug, Caddy xin TLS thật):

```bash
make up
```

---

## Bảng biến môi trường

Khai báo trong `deploy/.env` (copy từ `deploy/.env.example`).
Nguồn chuẩn: [architecture.md mục 12](docs/03-sa/architecture.md#12-cấu-hình-biến-môi-trường).

### Bắt buộc đổi trước khi lên production

| Biến | Vì sao |
|---|---|
| `POSTGRES_PASSWORD` | Mặc định là mật khẩu demo công khai |
| `ADMIN_TOKEN` | Bearer token cho `POST /api/v1/admin/ingest`. Để trống = tắt endpoint |
| `SITE_ADDRESS` | Đổi từ `:80` sang domain thật để Caddy cấp TLS |
| `PUBLIC_BASE_URL` | Dùng dựng short link tuyệt đối và metadata SEO |
| `ACME_EMAIL` | Nhận cảnh báo chứng chỉ sắp hết hạn |

### Toàn bộ biến

| Biến | Mặc định | Mô tả |
|---|---|---|
| `POSTGRES_USER` | `tenpoint` | User Postgres (cũng là superuser để tạo extension) |
| `POSTGRES_PASSWORD` | `tenpoint_dev_only_change_me` | 🔴 Mật khẩu DB |
| `POSTGRES_DB` | `tenpoint` | Tên database |
| `DATABASE_URL` | *(tự ghép)* | Compose tự dựng từ 3 biến trên → `postgres://…@db:5432/…` |
| `PORT` | `8080` | Cổng HTTP của API |
| `TZ` | `Asia/Ho_Chi_Minh` | Múi giờ cho cron + hiển thị. **Sai = cron lệch 7 tiếng** |
| `CRON_SPEC` | `0 6,14,22 * * *` | Lịch chạy pipeline (giờ VN) |
| `CRON_ENABLED` | `true` | `false` khi test. `make dev` tự ép `false` |
| `SUMMARIZER_PROVIDER` | `extractive` | `extractive` (không cần key) \| `anthropic` |
| `ANTHROPIC_API_KEY` | *(trống)* | 🔴 Secret. Chỉ cần khi provider = `anthropic` |
| `ANTHROPIC_MODEL` | `claude-sonnet-5` | Model dùng khi provider = `anthropic` |
| `ADMIN_TOKEN` | *(trống)* | 🔴 Bearer token cho `/api/v1/admin/ingest` |
| `PUBLIC_BASE_URL` | `http://localhost` | URL công khai, không có `/` cuối. Cũng cấp cho FE qua `NEXT_PUBLIC_SITE_URL` |
| `API_TIMEOUT_MS` | `4000` | Timeout mỗi lời gọi API từ SSR của Next.js |
| `FETCH_TIMEOUT_SEC` | `20` | Timeout tải mỗi bài |
| `MAX_ARTICLES_PER_SOURCE` | `40` | Trần số bài mỗi nguồn trong 1 run |
| `USER_AGENT` | `TenPointBot/1.0 (+https://tenpoint.vn/bot)` | Định danh khi crawl |
| `READY_MAX_AGE_HOURS` | `26` | `/readyz` trả `503` khi dữ liệu cũ hơn ngưỡng |
| `LOG_LEVEL` | `info` | `debug` \| `info` \| `warn` \| `error` |
| `LOG_FORMAT` | `json` | `json` (production) \| `text` (dev) |
| `SOURCES_PROFILE` | `tier0` | `tier0` (production) \| `all` (dev/QA) — xem Lưu ý pháp lý |
| `SITE_ADDRESS` | `:80` | `:80` = HTTP local; `tenpoint.vn` = TLS tự động |
| `ACME_EMAIL` | `devops@example.com` | Email đăng ký Let's Encrypt |
| `CSP_SCRIPT_SRC` | `'self' 'unsafe-inline'` | Nguồn script trong CSP — xem `deploy/Caddyfile` |
| `HTTP_PORT` / `HTTPS_PORT` | `80` / `443` | Cổng Caddy chiếm trên máy host |
| `DEV_DB_PORT` / `DEV_API_PORT` / `DEV_WEB_PORT` | `5432` / `8080` / `3000` | Chỉ có tác dụng với `make dev` |
| `APP_VERSION` / `IMAGE_PREFIX` | `dev` / `tenpoint` | Tag của image build ra |

---

## Chạy ingest thủ công

Cron chạy 3 lần/ngày. Khi cần lấy tin ngay (dev, backfill, sau khi sửa parser):

```bash
make ingest          # chạy binary `ingest` một lần trong container, rồi thoát
```

Hoặc qua HTTP, nếu API đang chạy và `ADMIN_TOKEN` đã đặt:

```bash
curl -X POST http://localhost:8080/api/v1/admin/ingest \
     -H "Authorization: Bearer $ADMIN_TOKEN"
# 202 Accepted — hoặc 409 nếu đang có pipeline chạy (advisory lock)
```

Kiểm tra kết quả:

```bash
make psql
tenpoint=# SELECT id, trigger, started_at, finished_at, articles_found, articles_new, errors
           FROM crawl_runs ORDER BY id DESC LIMIT 5;
```

---

## Xem log

```bash
make logs              # toàn bộ service, realtime
make logs S=api        # chỉ API (nơi có log pipeline + cron)
make logs S=db
make logs S=caddy
make ps                # trạng thái + healthcheck của từng container
```

Log được giới hạn `10 MB × 5 file` mỗi service (json-file driver) nên không làm đầy ổ đĩa.

---

## Toàn bộ lệnh

```bash
make help              # danh sách đầy đủ, tự sinh
```

| Lệnh | Việc |
|---|---|
| `make up` / `make down` | Khởi động / dừng (giữ dữ liệu) |
| `make dev` | Khởi động chế độ dev (mở cổng, tắt cron, log verbose) |
| `make build` / `make rebuild` | Build image (có cache / bỏ cache) |
| `make logs` / `make ps` | Log realtime / trạng thái |
| `make migrate` | Khởi động lại API để chạy migration |
| `make ingest` | Chạy pipeline thủ công |
| `make psql` | Mở psql vào database |
| `make backup-db` / `make restore-db FILE=…` | Sao lưu / khôi phục |
| `make test-be` / `make test-fe` / `make lint` | Kiểm thử & lint |
| `make clean` | ⚠️ Xoá cả volume — mất sạch dữ liệu |

---

## Cấu trúc thư mục

```
tenpoint/
├── backend/                  # Go 1.26 — module github.com/tenpoint/tenpoint-api
│   ├── cmd/api/              #   HTTP server + cron
│   ├── cmd/ingest/           #   CLI chạy pipeline 1 lần
│   ├── internal/             #   config · domain · store · http · ingest · shortlink
│   └── migrations/           #   SQL, được embed vào binary, tự chạy khi khởi động
├── frontend/                 # Next.js App Router + TypeScript + Tailwind
├── deploy/                   # ★ Toàn bộ hạ tầng
│   ├── Dockerfile.api        #   multi-stage Go → alpine non-root
│   ├── Dockerfile.web        #   multi-stage Next.js → standalone non-root
│   ├── web-entrypoint.sh     #   chọn standalone hoặc fallback `npm start`
│   ├── docker-compose.yml    #   4 service: caddy · web · api · db
│   ├── docker-compose.dev.yml#   override cho dev
│   ├── Caddyfile             #   routing + TLS + security headers
│   ├── .env.example          #   mẫu biến môi trường (có chú thích tiếng Việt)
│   └── README.md             #   ★ sổ tay vận hành
├── docs/                     # Tài liệu bàn giao
├── Makefile                  # Mọi lệnh vận hành
└── backups/                  # Bản dump DB (git-ignored)
```

---

## Tài liệu

| Tài liệu | Nội dung |
|---|---|
| [`docs/01-ba/BA1-users-and-usecases.md`](docs/01-ba/BA1-users-and-usecases.md) | Người dùng, use case, tiêu chí chấp nhận |
| [`docs/01-ba/BA2-data-sources-pipeline.md`](docs/01-ba/BA2-data-sources-pipeline.md) | Nguồn dữ liệu, phân tier, pipeline |
| [`docs/01-ba/BA3-market-differentiation-business.md`](docs/01-ba/BA3-market-differentiation-business.md) | Thị trường, khác biệt, mô hình kinh doanh |
| [`docs/02-pm/PRD.md`](docs/02-pm/PRD.md) | Phạm vi MVP, roadmap, rủi ro, **quyết định về nguồn (mục 4)** |
| [`docs/03-sa/architecture.md`](docs/03-sa/architecture.md) | Kiến trúc C4, ERD, API contract, ADR, **biến môi trường (mục 12)** |
| [`docs/04-design/design-system.md`](docs/04-design/design-system.md) | Design system, a11y |
| [`docs/00-inputs/reference-images.md`](docs/00-inputs/reference-images.md) | Ảnh tham chiếu gốc |
| [`deploy/README.md`](deploy/README.md) | Vận hành: deploy, đổi domain, backup, troubleshoot |

---

## Lưu ý pháp lý

> Tóm tắt [PRD mục 4](docs/02-pm/PRD.md). Đây là **rủi ro mức tồn vong** của dự án và
> **không giải quyết được bằng code**.

Tại Việt Nam, tổng hợp tin từ báo chí có thể cần **giấy phép "trang thông tin điện tử
tổng hợp"** (NĐ 72/2013, NĐ 147/2024) và **thoả thuận bằng văn bản với từng nguồn**.
Quy định yêu cầu **trích dẫn nguyên văn** — mâu thuẫn trực tiếp với việc tóm tắt bằng LLM.

Vì vậy hệ thống được xây **trung lập với nguồn**, bật/tắt bằng cấu hình
(`sources.enabled`) chứ không phải bằng code:

| Môi trường | Nguồn được bật | Biến |
|---|---|---|
| Dev / Demo / QA | Toàn bộ Tier 0 + Tier 1 + Tier 2 | `SOURCES_PROFILE=all` |
| **Production khi chưa có thoả thuận** | **Chỉ Tier 0** — HOSE, HNX, UBCKNN, NHNN, TCTK | `SOURCES_PROFILE=tier0` ← **mặc định** |
| Production sau khi ký thoả thuận | Mở thêm từng nguồn Tier 1 đã ký | bật lẻ trong bảng `sources` |

**`deploy/.env.example` mặc định là `tier0`.** `make dev` mới bật `all`.
Đổi sang `all` trên production là **quyết định của pháp chế, không phải của DevOps**.

**Cổng chặn phát hành — bắt buộc trước khi mở công khai:**

1. Tham vấn luật sư về giấy phép trang TTĐT tổng hợp.
2. Bật rule kỹ thuật chống vi phạm bản quyền: không copy ≥15 từ liên tiếp từ bài gốc,
   trùng 5-gram < 25%, không sao chép tiêu đề nguyên văn, không sao chép ảnh,
   không cache full-text ra public.
3. Banned-phrase list chặn ngôn ngữ tư vấn đầu tư **ở tầng sinh nội dung** — tư vấn đầu tư
   là nghiệp vụ có điều kiện, chỉ công ty chứng khoán được cấp phép mới được làm.

---

## Việc cần FE làm

### 1. Bật `output: 'standalone'` trong `frontend/next.config.ts` ⚠️

Hiện tại file này **chưa bật** standalone. `deploy/Dockerfile.web` đã xử lý cả hai trường
hợp: nếu không có `.next/standalone`, container tự fallback sang `npm start` và in cảnh báo
trong log. Nhưng đó chỉ là chống cháy — image nặng hơn nhiều (phải mang cả `node_modules`)
và khởi động chậm hơn.

```ts
// frontend/next.config.ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",   // ← thêm dòng này
};

export default nextConfig;
```

Sau khi thêm, chạy lại `make rebuild` và kiểm tra log `make logs S=web` — phải thấy
`[web] chế độ standalone`.

### 2. (Tuỳ chọn) Siết CSP bằng nonce

CSP hiện tại cho phép `script-src 'self' 'unsafe-inline'` vì Next.js App Router chèn các
thẻ `<script>` inline để stream RSC payload. Nếu FE thêm middleware sinh nonce
([hướng dẫn Next.js](https://nextjs.org/docs/app/guides/content-security-policy)),
DevOps sẽ đổi `CSP_SCRIPT_SRC="'self'"` trong `deploy/.env` để siết lại.

### 3. Đọc API qua biến môi trường

SSR phải gọi `process.env.API_BASE_URL` (mặc định `http://api:8080` trong Docker network),
**không** hard-code `localhost:8080`. Xem
[architecture.md mục 2](docs/03-sa/architecture.md#2-kiến-trúc-container-c4--level-2).

---

## Việc cần BE làm

`deploy/Dockerfile.api` build hai binary: `./cmd/api` và `./cmd/ingest`. Hai package này
phải tồn tại và biên dịch được thì `make build` mới chạy. Ngoài ra API cần có:

- `GET /healthz` — dùng cho `HEALTHCHECK` của container (không được phụ thuộc DB).
- `GET /readyz` — kiểm tra DB + độ tươi dữ liệu.
- Migrations embed trong binary, tự chạy khi khởi động (compose không có bước migrate riêng).
