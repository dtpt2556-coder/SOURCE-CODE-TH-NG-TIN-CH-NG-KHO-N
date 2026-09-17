# TenPoint — Test Plan (QA)

> Kế hoạch kiểm thử cho MVP v1. Nghiệm thu theo bảng chỉ tiêu mục 10 của `02-pm/PRD.md`.
> Truy vết yêu cầu: `01-ba/BA1-users-and-usecases.md` (user story + AC), `03-sa/architecture.md` (hành vi kỹ thuật), `04-design/design-system.md` (thị giác + a11y).

---

## 1. Phạm vi

**Trong phạm vi:** Backend Go (API, pipeline, tagger, summarizer, validator, short link, cron) · Frontend Next.js (3 trang, filter, responsive, a11y, SEO) · Docker Compose (khởi động, healthcheck, biến môi trường) · Đối chiếu thị giác với 2 ảnh tham chiếu.

**Ngoài phạm vi v1:** đăng nhập · thanh toán · realtime · app native · kiểm thử tải quy mô lớn (chỉ smoke performance).

---

## 2. Chiến lược & mức độ ưu tiên

| Mức | Nghĩa | Xử lý khi fail |
|---|---|---|
| **P0 — Blocker** | Sai số liệu tài chính, lỗ hổng bảo mật, sản phẩm không chạy | **Chặn phát hành**, sửa ngay |
| **P1 — Major** | Sai lệch rõ so với ảnh/AC, tính năng chính hỏng | Sửa trước khi phát hành |
| **P2 — Minor** | Lệch thị giác nhỏ, văn bản, edge case hiếm | Ghi nhận, xử lý ở v1.1 |

**Nguyên tắc QA của dự án này:** đây là sản phẩm tài chính. **Một con số sai nghiêm trọng hơn mười lỗi giao diện.** Mọi case nhóm TC-VAL và TC-SEC đều là P0 tuyệt đối.

---

## 3. Lệnh kiểm thử tự động

```bash
# Backend
cd backend
go build ./...            # phải sạch
go vet ./...              # phải sạch
gofmt -l .                # phải không in ra file nào
go test ./... -count=1    # tất cả PASS
go test ./... -race -count=1   # không có data race

# Frontend
cd frontend
npx tsc --noEmit
npm run lint
npm run build

# Hạ tầng
docker compose -f deploy/docker-compose.yml config
docker compose -f deploy/docker-compose.yml up -d
```

---

## 4. Test case

### 4.1 TC-VAL — Kỷ luật số liệu (P0, quan trọng nhất)

| ID | Mô tả | Dữ liệu vào | Kỳ vọng |
|---|---|---|---|
| TC-VAL-01 | Tóm tắt chứa số có trong bài gốc | body chứa `2.323 tỷ đồng`, summary chứa `2.323 tỷ đồng` | Validator PASS → `status='published'` |
| TC-VAL-02 | **Tóm tắt bịa số** | body chứa `2.323`, summary chứa `2.523` | Validator FAIL → `status='pending_review'`, **KHÔNG publish** |
| TC-VAL-03 | Định dạng nghìn kiểu VN | body `408.000 tỷ`, summary `408.000 tỷ` | PASS (dấu `.` là phân nhóm nghìn, không phải thập phân) |
| TC-VAL-04 | Thập phân kiểu VN | body `+1,67%`, summary `+1,67%` | PASS |
| TC-VAL-05 | Số bị làm tròn | body `29,91 điểm`, summary `30 điểm` | FAIL — không cho phép tự làm tròn |
| TC-VAL-06 | Tóm tắt không có số | bài định tính | PASS (không có số thì không có gì để đối chiếu) |
| TC-VAL-07 | Số trong đơn vị ghép | body `hơn 100 USD/tấn thép` | PASS, giữ nguyên đơn vị |
| TC-VAL-08 | Độ dài tóm tắt | summary 220 từ | Reject — vượt ngưỡng 180 từ (NUM-1) |
| TC-VAL-09 | Số cụm bold | summary có 5 cụm `**...**` | Cảnh báo — vượt ngưỡng 1–3 (NUM-3) |
| TC-VAL-10 | **Từ ngữ khuyến nghị đầu tư** | summary chứa "nên mua", "khuyến nghị nắm giữ" | **Reject** — vi phạm NUM-7, rủi ro pháp lý (PRD mục 11) |

### 4.2 TC-TAG — Gắn mã chứng khoán (P0/P1)

Chỉ tiêu: **precision ≥ 96%**, recall ≥ 85% (PRD mục 10).

| ID | Câu vào | Kỳ vọng |
|---|---|---|
| TC-TAG-01 | "FPT +2,69% trong phiên hôm nay" | `FPT` = `primary` |
| TC-TAG-02 | "GDP quý II tăng 6,8%" | **Không** gắn mã nào (`GDP` trong stoplist) |
| TC-TAG-03 | "giá trị 1,5 tỷ USD" | **Không** gắn `USD` |
| TC-TAG-04 | "kết phiên trên HOSE" | **Không** gắn `HOSE` |
| TC-TAG-05 | "công suất 30 MW" | **Không** gắn `MW` |
| TC-TAG-06 | "cổ phiếu VND tăng trần" | `VND` = mã VNDIRECT (có ngữ cảnh "cổ phiếu") |
| TC-TAG-07 | "doanh thu 2.323 tỷ VND" | **Không** gắn `VND` (là đơn vị tiền) |
| TC-TAG-08 | "Theo MBS Research/CBRE, thị trường..." | **Không** gắn `MBS` (nằm trong cụm trích dẫn nguồn) |
| TC-TAG-09 | "ACBS ước tính 5.588 tỷ đồng chảy vào..." | **Không** gắn `ACBS`; chỉ gắn mã được nhắc trong nội dung |
| TC-TAG-10 | "Hoà Phát báo lãi quý II" | `HPG` = `primary` (alias có dấu) |
| TC-TAG-11 | "hoa phat bao lai" (không dấu) | `HPG` — alias normalized |
| TC-TAG-12 | "Vinhomes mở bán dự án" | `VHM`, **không phải** `VIC` (phân giải mẹ–con) |
| TC-TAG-13 | "Vietcombank giảm lãi suất" | `VCB` |
| TC-TAG-14 | Tiêu đề có `CMG`, thân bài nhắc 1 lần | `CMG` = `primary` (xuất hiện ở tiêu đề) |
| TC-TAG-15 | Bài về VN-Index liệt kê 20 mã cuối bài | Các mã đó = `mentioned`, **không** `primary` |
| TC-TAG-16 | Bài CBAM không nhắc tên mã nào, chỉ nói "ngành thép" | `HPG, HSG, NKG, TIS` = `mentioned` (tầng T3 ngữ cảnh ngành) |
| TC-TAG-17 | Mã nằm trong URL `/tin-tuc/FPT-abc.html` | **Không** gắn (chỉ quét text đã extract) |
| TC-TAG-18 | Bài dùng "CEO", "CTO", "CFO" | **Không** gắn mã |

**Cách đo precision/recall:** lấy 50 bài mẫu đã gán nhãn thủ công, so với đầu ra của tagger.

### 4.3 TC-API — Hợp đồng API (P1)

| ID | Request | Kỳ vọng |
|---|---|---|
| TC-API-01 | `GET /api/v1/news` | `200`, có `data[]` + `meta{total,limit,offset,has_more}`, mặc định 7 ngày gần nhất, sắp xếp `published_at` giảm dần |
| TC-API-02 | `GET /api/v1/news?ma=FPT` | Chỉ trả tin có `FPT` |
| TC-API-03 | `GET /api/v1/news?ma=FPT,HPG` | Logic **OR** — tin khớp bất kỳ mã nào |
| TC-API-04 | `GET /api/v1/news?relevance=all` | Bao gồm cả mã `mentioned` |
| TC-API-05 | `GET /api/v1/news?loai=thi_truong` | Chỉ trả đúng loại |
| TC-API-06 | `GET /api/v1/news?q=chung khoan` | Tìm được bài có "chứng khoán" (**không dấu vẫn ra** — US-2.4) |
| TC-API-07 | `GET /api/v1/news?limit=999` | Bị chặn về `100` |
| TC-API-08 | `GET /api/v1/news?tu=2026-08-20&den=2026-08-27` | Lọc đúng khoảng, bao gồm 2 đầu mút theo giờ VN |
| TC-API-09 | `GET /api/v1/news/999999` | `404` + `{"error":{"code":"not_found","message":"..."}}` tiếng Việt |
| TC-API-10 | `GET /api/v1/tickers/FPT` | `200` với tin gần đây + research notes |
| TC-API-11 | `GET /api/v1/notes/{slug}` | `200` với đủ các luận điểm theo `ordinal` |
| TC-API-12 | `GET /api/v1/meta` | Có `last_crawl_at`, `is_stale` |
| TC-API-13 | `GET /healthz` | `200` |
| TC-API-14 | `GET /readyz` khi DB chết | `503` |
| TC-API-15 | `POST /api/v1/admin/ingest` không token | `401` |
| TC-API-16 | `POST /api/v1/admin/ingest` token sai | `401` |
| TC-API-17 | `POST /api/v1/admin/ingest` token đúng | `202` |
| TC-API-18 | Gọi ingest 2 lần liên tiếp | Lần 2 trả `409` (advisory lock chặn chạy chồng) |
| TC-API-19 | Kiểm tra `published_date_vn` | Đúng định dạng `DD/MM/YYYY` theo `Asia/Ho_Chi_Minh` |
| TC-API-20 | Kiểm tra `summary_md` | Là markdown thô chứa `**...**`, **không phải HTML** |

### 4.4 TC-SEC — Bảo mật (P0)

| ID | Kịch bản | Kỳ vọng |
|---|---|---|
| TC-SEC-01 | Tạo short link trỏ tới `https://evil.example.com` | **Bị từ chối** — ngoài allowlist domain, có log cảnh báo |
| TC-SEC-02 | Tạo short link trỏ tới `https://sub.cafef.vn/...` | Chấp nhận (subdomain của nguồn hợp lệ) |
| TC-SEC-03 | Tạo short link trỏ tới `https://cafef.vn.evil.com/...` | **Bị từ chối** — không được khớp lỏng theo chuỗi |
| TC-SEC-04 | `GET /r/khongtontai` | `404` + trang lỗi tiếng Việt, **KHÔNG redirect về trang chủ** (US-3.2 AC3) |
| TC-SEC-05 | SQL injection qua `?q=' OR 1=1--` | Không lỗi, không lộ dữ liệu (pgx dùng prepared statement) |
| TC-SEC-06 | XSS qua `summary_md` chứa `<script>alert(1)</script>` | FE render ra text thuần, **không thực thi** |
| TC-SEC-07 | `summary_md` chứa `<img src=x onerror=alert(1)>` | Bị sanitize, chỉ `<strong>`/`<em>` được phép |
| TC-SEC-08 | Kiểm tra link nguồn trong HTML | Có đủ `rel="nofollow noopener noreferrer"` và `target="_blank"` |
| TC-SEC-09 | Quét secret trong repo | Không có API key / mật khẩu hardcode |
| TC-SEC-10 | Container chạy bằng user nào | **Non-root** cho cả `api` và `web` |

### 4.5 TC-PIPE — Pipeline & Cron (P1)

| ID | Kịch bản | Kỳ vọng |
|---|---|---|
| TC-PIPE-01 | Chạy pipeline với 1 nguồn trả lỗi HTTP 500 | Các nguồn khác **vẫn chạy xong**; lỗi ghi vào `crawl_run_sources.error_message` (US-6.1 AC2) |
| TC-PIPE-02 | Chạy pipeline 2 lần, cùng 1 bài | Lần 2 không tạo bản ghi trùng (`url_hash` UNIQUE) |
| TC-PIPE-03 | 2 bài khác nguồn, tiêu đề gần giống | Được gom vào cùng `cluster_id`; chỉ 1 bài `is_canonical_in_cluster=true`, chọn nguồn tier cao hơn |
| TC-PIPE-04 | 2 bài tiêu đề khác hẳn | Khác cluster |
| TC-PIPE-05 | Kiểm tra `CRON_SPEC` mặc định | `0 6,14,22 * * *` theo `Asia/Ho_Chi_Minh` |
| TC-PIPE-06 | Container thiếu `tzdata` | Phát hiện được — cron sẽ chạy sai giờ. **Phải có tzdata trong image** |
| TC-PIPE-07 | `CRON_ENABLED=false` | Cron không chạy |
| TC-PIPE-08 | Nguồn trả 0 bài 2 run liên tiếp | Có cảnh báo (PRD Q7) |
| TC-PIPE-09 | robots.txt cấm đường dẫn | Pipeline **không** fetch đường dẫn đó |
| TC-PIPE-10 | Thời lượng 1 run | < 10 phút (PRD mục 10) |

### 4.6 TC-UI — Giao diện, đối chiếu ảnh tham chiếu (P1)

Dùng dữ liệu seed demo (5 tin + 1 research note lấy nguyên từ ảnh).

| ID | Kiểm tra | Kỳ vọng |
|---|---|---|
| TC-UI-01 | Bảng digest desktop ≥1024px | Đúng 5 cột theo thứ tự `Ngày · Mã CK · Tóm tắt thông tin · Source · Loại tin` |
| TC-UI-02 | Đường kẻ bảng | Hairline 1px `--rule` (light `#D5D8D2`, dark `#262C34`) ngăn hàng; kẻ đậm `--rule-strong` dưới header; **không** border dọc, **không** zebra |
| TC-UI-03 | Nền trang | `--surface`: light `#F4F5F2`, dark `#101317`. **Không** trắng thuần, **không** gradient, và **không** thuộc nhóm kem/beige `#F7F5F0`/`#f5f1ea`/`#fbf8f1` mà design-system v2 mục 10 cấm |
| TC-UI-02b | Token màu | Đối chiếu **toàn bộ** token mục 2.1 và 2.2 design-system v2 bằng computed style, ở **cả hai** chế độ sáng/tối |
| TC-UI-04 | Typography bảng | Serif; tiêu đề Research Note dùng sans-serif đậm UPPERCASE |
| TC-UI-05 | Nhấn số liệu | Chỉ bằng `font-weight: 700`, **không** đổi màu, **không** highlight nền |
| TC-UI-06 | Cột Source | Hiện domain rút gọn có gạch chân (`vnexpress.net`), không phải URL đầy đủ |
| TC-UI-07 | Tiêu đề động | Đúng dạng `{N} tin mới đáng chú ý — {mã} · {mã} · {mã}.` |
| TC-UI-08 | `Cập nhật lần cuối` | Hiển thị đúng giờ VN |
| TC-UI-09 | Dữ liệu cũ > 12h | Hiện cảnh báo |
| TC-UI-10 | Research Note | Tiêu đề UPPERCASE, section `LUẬN ĐIỂM ĐẦU TƯ`, 4 luận điểm đánh số với dòng dẫn bold |
| TC-UI-11 | Disclaimer | **Có mặt** ở cuối mọi research note |
| TC-UI-12 | Bo góc / shadow | **Không có** (trừ chip mã CK radius 2px) |
| TC-UI-13 | Empty state | Thông báo tiếng Việt, không phải bảng rỗng |
| TC-UI-14 | Mã CK > 8 mã | Hiện 8 + `+N`, không để ô cao quá |

### 4.7 TC-RWD — Responsive (P1)

| ID | Viewport | Kỳ vọng |
|---|---|---|
| TC-RWD-01 | 375px (iPhone SE) | Bảng → card dọc; **KHÔNG scroll ngang** (yêu cầu cứng) |
| TC-RWD-02 | 375px | Thứ tự trong card: `Ngày · Loại tin` → `Mã CK` → `Tóm tắt` → `Source` |
| TC-RWD-03 | 768px | 4 cột, `Loại tin` gộp vào dòng meta |
| TC-RWD-04 | 1024px | Đủ 5 cột |
| TC-RWD-05 | 1920px | Container khoá `1180px`, không giãn tràn |
| TC-RWD-06 | Mọi viewport | `document.body.scrollWidth <= window.innerWidth` |

### 4.8 TC-A11Y — Khả năng tiếp cận (P1)

| ID | Kiểm tra | Kỳ vọng |
|---|---|---|
| TC-A11Y-01 | `<html lang>` | `vi` |
| TC-A11Y-02 | Cấu trúc bảng | Có `<caption>` (sr-only) + `<th scope="col">` |
| TC-A11Y-03 | Contrast | Mọi cặp text/nền ≥ 4.5:1 |
| TC-A11Y-04 | Điều hướng bàn phím | Tab qua được toàn bộ filter + link; focus ring luôn nhìn thấy |
| TC-A11Y-05 | Filter chips | Là `<button aria-pressed>`, không phải `<div onClick>` |
| TC-A11Y-06 | Link mở tab mới | Có `aria-label` ghi rõ "(mở tab mới)" |
| TC-A11Y-07 | Skip link | Có "Bỏ qua tới nội dung chính" |
| TC-A11Y-08 | Target size | ≥ 24×24px (WCAG 2.2 — 2.5.8) |
| TC-A11Y-09 | Font tiếng Việt | Hiển thị đúng `Ừ Ữ Ỡ Ợ Ặ Ẫ Ỹ ọ ự ẳ` — **không ô vuông, không mất dấu** |
| TC-A11Y-10 | `prefers-reduced-motion` | Được tôn trọng |

### 4.9 TC-SEO (P2)

| ID | Kiểm tra | Kỳ vọng |
|---|---|---|
| TC-SEO-01 | Trang được render server-side | View-source thấy nội dung tin, không phải shell rỗng |
| TC-SEO-02 | `title`/`description` | Động theo từng trang |
| TC-SEO-03 | `/sitemap.xml` | Có, gồm các trang `/ma/{mã}` |
| TC-SEO-04 | `/robots.txt` | Có |
| TC-SEO-05 | JSON-LD | `NewsArticle` trên research note |
| TC-SEO-06 | Canonical URL | Có trên mọi trang |

### 4.10 TC-DEPLOY (P1)

| ID | Kiểm tra | Kỳ vọng |
|---|---|---|
| TC-DEPLOY-01 | `docker compose config` | Parse sạch |
| TC-DEPLOY-02 | `docker compose up` từ máy sạch | 4 service `healthy` |
| TC-DEPLOY-03 | Migration | Tự chạy lúc khởi động, idempotent khi restart |
| TC-DEPLOY-04 | Seed demo | Trang chủ hiện đúng 5 tin từ ảnh tham chiếu |
| TC-DEPLOY-05 | Khởi động khi **không có** `ANTHROPIC_API_KEY` | Hệ thống vẫn chạy với provider `extractive` (ADR-007) |
| TC-DEPLOY-06 | `TZ` trong container `api` | `Asia/Ho_Chi_Minh` |
| TC-DEPLOY-07 | Restart container | Dữ liệu còn nguyên (volume) |
| TC-DEPLOY-08 | `.env.example` | Không chứa secret thật |

---

## 5. Ma trận truy vết yêu cầu

| Yêu cầu | Nguồn | Test case |
|---|---|---|
| Bảng digest 5 cột như ảnh | Ảnh B | TC-UI-01..07, TC-DEPLOY-04 |
| Tóm tắt giàu số liệu, bold số | BA1 NUM-1..8 | TC-VAL-01..10, TC-UI-05 |
| Không bịa số | PRD Q5 (tồn vong) | **TC-VAL-02, TC-VAL-05** |
| Gắn mã CK | BA2 | TC-TAG-01..18 |
| Short link + trích dẫn nguồn | Yêu cầu gốc | TC-SEC-01..04, TC-UI-06 |
| Chống open redirect | BA1 US-3.3 | **TC-SEC-01, 03** |
| Cron 8h/lần | Yêu cầu gốc | TC-PIPE-05, 06, 07 |
| Nguồn uy tín | Yêu cầu gốc | TC-PIPE-01, 09 |
| Lọc theo mã / loại tin | BA1 US-2.x | TC-API-02..05 |
| Research Note như ảnh | Ảnh A | TC-UI-10, 11 |
| Deploy bằng Docker | Yêu cầu gốc | TC-DEPLOY-01..08 |

---

## 6. Tiêu chí phát hành

**Được phát hành khi:**
- [ ] 100% case P0 PASS (TC-VAL, TC-SEC)
- [ ] ≥ 95% case P1 PASS
- [ ] `go build`, `go vet`, `gofmt -l`, `go test ./...` sạch
- [ ] `tsc --noEmit`, `npm run lint`, `npm run build` sạch
- [ ] `docker compose up` cho 4 service healthy
- [ ] Đối chiếu thị giác với 2 ảnh tham chiếu đạt
- [ ] Precision tag mã ≥ 96%
- [ ] Không có lỗ hổng bảo mật mức P0

**Chặn phát hành ngay nếu:** có bất kỳ tin nào publish với số liệu bịa · short link redirect được tới domain ngoài allowlist · bảng scroll ngang trên mobile · thiếu disclaimer trên research note.

---

## 7. Sản phẩm bàn giao của QA

`docs/05-qa/qa-report.md` — kết quả thực tế từng case (PASS/FAIL/BLOCKED), output lệnh thật, danh sách defect kèm mức độ, và kết luận nghiệm thu.
