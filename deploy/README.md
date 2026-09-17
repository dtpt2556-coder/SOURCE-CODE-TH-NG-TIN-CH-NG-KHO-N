# TenPoint — Sổ tay vận hành (`deploy/`)

Tài liệu dành cho DevOps/SRE. Phần giới thiệu sản phẩm và hướng dẫn chạy nhanh nằm ở
[`../README.md`](../README.md).

---

## Nội dung thư mục

| File | Vai trò |
|---|---|
| `docker-compose.yml` | Cấu hình **production**. 4 service: `caddy` · `web` · `api` · `db` |
| `docker-compose.dev.yml` | **Override** cho dev: mở cổng debug, tắt cron, log verbose |
| `Dockerfile.api` | Multi-stage Go 1.26 → `alpine:3.20`, non-root `tenpoint` (UID 10001) |
| `Dockerfile.web` | Multi-stage Next.js → `node:22-alpine`, non-root `nextjs` (UID 1001) |
| `web-entrypoint.sh` | Chọn `node server.js` (standalone) hoặc fallback `npm start` |
| `Caddyfile` | Routing, TLS tự động, nén, security headers |
| `.env.example` | Mẫu biến môi trường, chú thích tiếng Việt |

**Build context của cả hai Dockerfile là thư mục GỐC dự án** (`context: ..` trong compose),
vì vậy mọi `COPY` đều bắt đầu bằng `backend/…` hoặc `frontend/…`.

---

## 1. Deploy production

### 1.1. Chuẩn bị máy chủ

```bash
# Ubuntu 22.04/24.04
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker "$USER" && newgrp docker
docker compose version      # phải ≥ v2.24
```

Mở firewall cho 80/tcp, 443/tcp và 443/udp (HTTP/3). **Không** mở 5432/8080/3000 —
cấu hình production không map các cổng đó ra host.

### 1.2. Cấu hình

```bash
git clone <repo> /opt/tenpoint && cd /opt/tenpoint
cp deploy/.env.example deploy/.env
```

Sửa `deploy/.env`, tối thiểu 5 dòng sau:

```dotenv
POSTGRES_PASSWORD=<openssl rand -hex 32>
ADMIN_TOKEN=<openssl rand -hex 32>
SITE_ADDRESS=tenpoint.vn
PUBLIC_BASE_URL=https://tenpoint.vn
ACME_EMAIL=devops@tenpoint.vn
```

Giữ nguyên `SOURCES_PROFILE=tier0` cho tới khi pháp chế phê duyệt
(xem "Lưu ý pháp lý" trong README gốc).

```bash
chmod 600 deploy/.env
```

### 1.3. Khởi động

```bash
make build
make up
make ps        # cả 4 service phải ở trạng thái healthy sau ~1 phút
```

Thứ tự khởi động do compose đảm bảo: `db` (healthy) → `api` (healthy) → `web` → `caddy`.
API tự chạy migration khi khởi động (migrations được embed trong binary).

### 1.4. Nghiệm thu sau deploy

```bash
curl -fsS https://tenpoint.vn/api/v1/meta | jq .        # có last_crawl_at, is_stale
curl -fsSI https://tenpoint.vn/ | grep -iE 'content-security-policy|x-frame|strict-transport'
docker compose -f deploy/docker-compose.yml exec api date   # phải ra giờ VN (+07)
```

Dòng thứ ba là kiểm tra quan trọng nhất: nếu `date` ra giờ UTC thì cron sẽ chạy sai
(xem mục 5.2).

---

## 2. Đổi domain

1. Trỏ bản ghi DNS `A` (và `AAAA` nếu có IPv6) của domain về IP máy chủ. **Đợi DNS
   propagate xong** — kiểm tra bằng `dig +short tenpoint.vn`. Caddy xin chứng chỉ bằng
   HTTP-01 challenge, DNS chưa trỏ đúng thì chắc chắn thất bại.
2. Sửa `deploy/.env`:
   ```dotenv
   SITE_ADDRESS=tenpoint.vn, www.tenpoint.vn
   PUBLIC_BASE_URL=https://tenpoint.vn
   ```
3. Nạp lại chỉ service caddy:
   ```bash
   docker compose -f deploy/docker-compose.yml up -d --force-recreate caddy
   make logs S=caddy
   ```
4. `api` và `web` cũng cần restart vì `PUBLIC_BASE_URL` / `NEXT_PUBLIC_SITE_URL` đổi:
   ```bash
   docker compose -f deploy/docker-compose.yml up -d --force-recreate api web
   ```

> ⚠️ Let's Encrypt giới hạn **5 chứng chỉ trùng nhau / 7 ngày**. Khi thử nghiệm domain
> mới, dùng CA staging trước bằng cách thêm vào block global của `Caddyfile`:
> `acme_ca https://acme-staging-v02.api.letsencrypt.org/directory`

---

## 3. Sao lưu & khôi phục

### 3.1. Sao lưu thủ công

```bash
make backup-db        # → backups/tenpoint-<YYYYmmdd-HHMMSS>.dump (format custom)
```

### 3.2. Sao lưu tự động (cron trên host)

```cron
# /etc/cron.d/tenpoint-backup — 03:15 mỗi ngày, giữ 14 bản gần nhất
15 3 * * *  root  cd /opt/tenpoint && make backup-db >> /var/log/tenpoint-backup.log 2>&1
30 3 * * *  root  find /opt/tenpoint/backups -name '*.dump' -mtime +14 -delete
```

Đẩy bản dump ra ngoài máy chủ (S3/rclone) — volume Docker nằm cùng ổ đĩa với DB, mất
ổ là mất cả hai.

### 3.3. Khôi phục

```bash
make restore-db FILE=backups/tenpoint-20260916-031500.dump
```

Lệnh dùng `pg_restore --clean --if-exists`, tức **ghi đè** dữ liệu hiện tại. Có 5 giây
để Ctrl-C huỷ. Nên dừng `api` trước để pipeline không ghi xen vào:

```bash
docker compose -f deploy/docker-compose.yml stop api
make restore-db FILE=backups/....dump
docker compose -f deploy/docker-compose.yml start api
```

### 3.4. Sao lưu chứng chỉ TLS

Volume `tenpoint_caddy_data` chứa private key + chứng chỉ. Mất volume này không mất
dữ liệu nghiệp vụ, nhưng Caddy phải xin cert lại (dễ dính rate limit).

```bash
docker run --rm -v tenpoint_caddy_data:/data -v "$PWD/backups":/out alpine \
    tar czf /out/caddy-data-$(date +%F).tgz -C /data .
```

---

## 4. Xem log & quan sát

```bash
make logs                 # tất cả service
make logs S=api           # pipeline, cron, HTTP request
make logs S=caddy         # access log JSON + ACME
make ps                   # health status từng container

# Lọc theo mốc thời gian
docker compose -f deploy/docker-compose.yml logs --since 30m api

# Tình trạng pipeline từ DB
make psql
tenpoint=# SELECT id, trigger, started_at, finished_at,
                  articles_found, articles_new, articles_rejected, errors
           FROM crawl_runs ORDER BY id DESC LIMIT 10;
tenpoint=# SELECT s.code, crs.found, crs.new_count, crs.error_message
           FROM crawl_run_sources crs JOIN sources s ON s.id = crs.source_id
           WHERE crs.run_id = (SELECT max(id) FROM crawl_runs);
```

Log driver là `json-file`, xoay vòng `10 MB × 5 file` mỗi service (tối đa ~200 MB tổng).

Endpoint quan sát (chỉ truy cập được **trong** docker network — Caddy không proxy ra ngoài,
đây là chủ ý để không lộ thông tin nội bộ):

```bash
docker compose -f deploy/docker-compose.yml exec api wget -qO- http://127.0.0.1:8080/readyz
```

---

## 5. Troubleshoot

### 5.1. `api` khởi động lại liên tục — DB chưa sẵn sàng

**Triệu chứng:** `make ps` thấy `tenpoint-api` ở trạng thái `Restarting`; log có
`connection refused` hoặc `the database system is starting up`.

**Kiểm tra:**
```bash
docker compose -f deploy/docker-compose.yml ps db      # db đã `healthy` chưa?
make logs S=db
```

**Nguyên nhân & xử lý:**
- Lần khởi động đầu, `initdb` mất 15–30s. `depends_on: db: condition: service_healthy`
  đã xử lý việc này — nếu vẫn lỗi, đọc tiếp.
- Sai mật khẩu: volume `tenpoint_pgdata` **chỉ nhận `POSTGRES_PASSWORD` ở lần initdb đầu
  tiên**. Đổi mật khẩu trong `.env` sau đó **không** đổi mật khẩu trong DB. Sửa bằng cách
  đổi mật khẩu trực tiếp:
  ```bash
  make psql
  tenpoint=# ALTER USER tenpoint WITH PASSWORD 'mật_khẩu_mới';
  ```
  (Hoặc `make clean` để xoá volume — **mất toàn bộ dữ liệu**.)
- Thiếu `DATABASE_URL`: config của BE fail-fast khi biến này trống. Kiểm tra bằng
  `docker compose -f deploy/docker-compose.yml config | grep DATABASE_URL`.

### 5.2. Cron chạy sai giờ (lệch 7 tiếng) — thiếu `tzdata`

**Triệu chứng:** pipeline chạy lúc 13:00 / 21:00 / 05:00 thay vì 06:00 / 14:00 / 22:00.
Cột `crawl_runs.started_at` lệch đúng 7 giờ.

**Kiểm tra:**
```bash
docker compose -f deploy/docker-compose.yml exec api date
# ĐÚNG:  Tue Sep 16 06:00:00 +07 2026
# SAI :  Mon Sep 15 23:00:00 UTC 2026

docker compose -f deploy/docker-compose.yml exec api ls /usr/share/zoneinfo/Asia/Ho_Chi_Minh
# Phải tồn tại. Nếu "No such file" → image thiếu gói tzdata.
```

**Nguyên nhân:** Go gọi `time.LoadLocation("Asia/Ho_Chi_Minh")`, hàm này đọc
`/usr/share/zoneinfo`. Image `alpine`/`distroless` **không có sẵn** thư mục đó. Thiếu
tzdata thì `LoadLocation` trả lỗi và scheduler rơi về UTC.

**Xử lý:** `deploy/Dockerfile.api` đã có `apk add --no-cache tzdata`. Nếu ai đó sửa
Dockerfile và bỏ dòng này, thêm lại rồi `make rebuild`. Kiểm tra thêm `TZ` trong
`deploy/.env` đúng chính tả `Asia/Ho_Chi_Minh` (không phải `Asia/Saigon`).

> Đây là lỗi im lặng: hệ thống vẫn chạy, API vẫn trả 200, chỉ có tin ra sai giờ.
> Luôn chạy `exec api date` sau mỗi lần deploy.

### 5.3. Caddy không lấy được chứng chỉ TLS

**Triệu chứng:** trình duyệt báo lỗi chứng chỉ; `make logs S=caddy` có
`could not get certificate from issuer` / `challenge failed`.

**Chẩn đoán theo thứ tự:**

| # | Kiểm tra | Lệnh |
|---|---|---|
| 1 | DNS đã trỏ đúng IP? | `dig +short tenpoint.vn` so với `curl -s ifconfig.me` |
| 2 | Cổng 80 có tới được từ Internet? | `curl -I http://tenpoint.vn/.well-known/acme-challenge/test` |
| 3 | Có dịch vụ khác chiếm 80/443? | `sudo ss -tlnp \| grep -E ':(80\|443)'` |
| 4 | `SITE_ADDRESS` còn là `:80`? | `docker compose -f deploy/docker-compose.yml config \| grep SITE_ADDRESS` |
| 5 | Dính rate limit Let's Encrypt? | log caddy có `too many certificates already issued` |

**Xử lý thường gặp:**
- ACME HTTP-01 **bắt buộc** đi qua cổng 80. Firewall/Cloud Security Group chặn 80 là
  nguyên nhân phổ biến nhất. Mở 80 kể cả khi chỉ dùng HTTPS.
- Cloudflare bật proxy (mây cam) → challenge không tới được Caddy. Tạm chuyển sang
  DNS-only (mây xám) cho tới khi cấp cert xong, hoặc dùng DNS-01 challenge.
- Dính rate limit: chuyển tạm sang CA staging (mục 2), sửa cho chạy được rồi mới đổi lại.
- Reset hoàn toàn trạng thái ACME:
  ```bash
  docker compose -f deploy/docker-compose.yml down
  docker volume rm tenpoint_caddy_data      # ⚠️ xoá hết cert đã cấp
  make up
  ```

### 5.4. `web` in cảnh báo "KHÔNG tìm thấy /app/server.js"

Không phải lỗi — là fallback. `frontend/next.config.ts` chưa bật `output: 'standalone'`,
container tự chuyển sang `npm start`. Trang vẫn chạy nhưng image nặng hơn và khởi động
chậm hơn. Xem mục "Việc cần FE làm" trong README gốc.

### 5.5. Trang render ra nhưng không bấm được gì (CSP chặn script)

**Triệu chứng:** HTML hiện đầy đủ, console trình duyệt báo
`Refused to execute inline script because it violates ... Content Security Policy`.

**Nguyên nhân:** `CSP_SCRIPT_SRC` bị siết thành `'self'` trong khi Next.js App Router
vẫn chèn inline script (`self.__next_f.push(...)`) để stream RSC payload mà chưa có nonce.

**Xử lý:** đặt lại trong `deploy/.env`:
```dotenv
CSP_SCRIPT_SRC='self' 'unsafe-inline'
```
rồi `docker compose -f deploy/docker-compose.yml up -d --force-recreate caddy`.
Chỉ siết về `'self'` sau khi FE thêm middleware sinh nonce.

### 5.6. Pipeline không lấy được bài nào / bị nguồn chặn

```bash
make psql
tenpoint=# SELECT s.code, crs.found, crs.error_message
           FROM crawl_run_sources crs JOIN sources s ON s.id = crs.source_id
           WHERE crs.run_id = (SELECT max(id) FROM crawl_runs);
```

- `found = 0` ở **mọi** nguồn → thường là DNS/egress của container. Thử
  `docker compose -f deploy/docker-compose.yml exec api wget -qO- https://vnexpress.net/rss/kinh-doanh.rss | head`.
- `found = 0` ở **một** nguồn trong 2 run liên tiếp → nguồn đổi HTML, parser gãy.
  Đây là cảnh báo đã được định nghĩa trong PRD (rủi ro "Nguồn đổi HTML → parser gãy").
- Lỗi `403` / `429` → bị rate limit. Tăng `sources.rate_limit_ms` trong DB, giảm
  `MAX_ARTICLES_PER_SOURCE`, và kiểm tra `USER_AGENT` có định danh rõ ràng không.

### 5.7. Giá trị trong container khác với `deploy/.env`

**Triệu chứng:** sửa `deploy/.env` nhưng `docker compose config` vẫn in ra giá trị cũ/lạ.

**Nguyên nhân:** Docker Compose ưu tiên **biến môi trường của shell** cao hơn file `.env`.
Nếu shell của bạn (hoặc `/etc/environment`, hoặc một công cụ dev nào đó) đã export sẵn
`ANTHROPIC_MODEL`, `ANTHROPIC_API_KEY`, `TZ`, `LOG_LEVEL`… thì giá trị đó thắng.

**Kiểm tra:**
```bash
env | grep -E '^(TZ|LOG_|CRON_|ANTHROPIC_|POSTGRES_|SITE_ADDRESS|ADMIN_TOKEN)='
docker compose -f deploy/docker-compose.yml config | grep ANTHROPIC_MODEL
```

**Xử lý:** `unset <BIẾN>` trước khi chạy, hoặc chạy trong shell sạch:
```bash
env -i PATH="$PATH" HOME="$HOME" make up
```

### 5.8. Cổng 80/443 đã bị chiếm trên host

```bash
sudo ss -tlnp | grep -E ':(80|443)'
```
Đổi trong `deploy/.env`:
```dotenv
HTTP_PORT=8081
HTTPS_PORT=8443
```
Lưu ý: ACME HTTP-01 **phải** đi qua cổng 80 công khai. Nếu đổi cổng, cần một reverse
proxy khác ở phía trước forward 80 → 8081, hoặc chuyển sang DNS-01 challenge.

---

## 6. Ghi chú thiết kế

| Quyết định | Lý do |
|---|---|
| Runtime API là `alpine:3.20` chứ không phải distroless | Cần shell + `wget` cho `HEALTHCHECK` gọi `/healthz`. Distroless nhỏ hơn nhưng phải nhúng binary healthcheck riêng — không đáng ở quy mô MVP |
| Cả hai image chạy **non-root** (`tenpoint` 10001, `nextjs` 1001) | Giảm thiệt hại nếu có RCE qua parser HTML hoặc Next.js |
| `tzdata` được cài ở **cả** api và web | Cron phải đúng giờ VN (mục 5.2); web SSR phải hiển thị ngày giờ VN |
| `ca-certificates` trong image api | Bắt buộc để gọi HTTPS tới nguồn tin và Claude API |
| `db` không map cổng ở production | Postgres chỉ tiếp cận được trong docker network. `make dev` mới mở, và chỉ bind `127.0.0.1` |
| Healthcheck ở **cả** Dockerfile và compose | Dockerfile để `docker run` thủ công vẫn có; compose để `depends_on: condition: service_healthy` hoạt động |
| `start_period=40s` cho api | Lần khởi động đầu phải chạy migration trước khi phục vụ `/healthz` |
| Resource limit cho mọi service | Một pipeline lỗi (vòng lặp vô hạn khi parse HTML) không được kéo sập cả máy chủ |
| Caddy chỉ proxy `/api/*` và `/r/*` | `/healthz`, `/readyz`, `/metrics` không lộ ra Internet |
| HSTS chỉ bật khi `scheme == https` | Bật ở `http://localhost` sẽ khiến trình duyệt ghim localhost sang https vĩnh viễn |
| `admin off` trong Caddyfile | Tắt admin API cổng 2019 — bề mặt tấn công thừa |
