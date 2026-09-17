# TenPoint — Đặc tả Edge Case cho Pipeline Cập nhật Tin

> Phát sinh từ phản hồi khách hàng: *"tin tức cập nhật BE phải suy nghĩ ra cover được các case"* và *"phần trích dẫn link chưa đưa đến thông tin trang cụ thể"*.
> Mỗi case có mã. BE implement theo mã, QA test theo mã.
> Bổ sung cho `docs/03-sa/architecture.md` mục 5 và 8.

---

## 0. Vấn đề gốc phải sửa trước: link không dẫn tới bài cụ thể

Dữ liệu seed demo hiện tại chứa URL **bịa**:

```
https://vnexpress.net/12-ngan-hang-cam-ket-408000-ty-dong-tin-dung-cho-dnnvv-4790001.html
https://cafef.vn/vn-index-dong-cua-1821-diem-lan-dau-vuot-moc-1800-188260826154500123.chn
```

Trông đúng định dạng của từng báo, nhưng **không tồn tại**. Bấm vào ra 404. Đây chính là lỗi khách hàng nhìn thấy.

### Cách sửa — 3 tầng

| Tầng | Nội dung |
|---|---|
| **T1. Bootstrap dữ liệu thật** | Khi khởi động, nếu chưa có bài nào `is_demo = false`, chạy ingest thật ngay lập tức (không chờ cron). Sau vài giây bảng digest hiển thị **tin thật, link thật**. Đây là đường đi chính của sản phẩm. |
| **T2. Đánh dấu dữ liệu demo** | Thêm cột `articles.is_demo BOOLEAN DEFAULT false`. 5 bài seed từ ảnh tham chiếu đặt `is_demo = true`, `short_links.target_alive = false`. FE render chúng dưới dạng text mờ kèm chú thích "Dữ liệu mẫu", **không** phải link bấm được. Không bao giờ để người dùng bấm vào một URL bịa. |
| **T3. Ẩn demo khi có dữ liệu thật** | `GET /api/v1/news` mặc định lọc `is_demo = false` **nếu** tồn tại ≥ 1 bài thật; chỉ fallback sang demo khi CSDL rỗng (không có mạng, chưa crawl lần nào). |

> **Nguyên tắc:** thà hiển thị "Dữ liệu mẫu" trung thực còn hơn một link đẹp mà bấm vào ra 404. Sản phẩm này bán bằng **độ tin cậy của trích dẫn**.

---

## 1. Khám phá bài (Discovery)

| ID | Tình huống | Hành vi đúng |
|---|---|---|
| **D-01** | Nguồn có RSS hợp lệ | Ưu tiên RSS. Gửi kèm `If-Modified-Since` / `If-None-Match` từ lần chạy trước; `304` thì bỏ qua nguồn, không tính là lỗi |
| **D-02** | Nguồn không có RSS | Fallback parse trang danh sách bằng CSS selector lưu trong `sources.list_selector` |
| **D-03** | RSS trả XML hỏng, thiếu đóng thẻ | Parse lỗi không được panic. Ghi `crawl_run_sources.error_message`, thử fallback HTML |
| **D-04** | RSS trả 0 item | Ghi nhận `found = 0`. **2 run liên tiếp bằng 0 thì cảnh báo** (PRD Q7): nhiều khả năng nguồn đổi cấu trúc |
| **D-05** | RSS trả 500 item (backfill lịch sử) | Cắt theo `MAX_ARTICLES_PER_SOURCE` (mặc định 40) **sau khi đã sắp xếp theo thời gian giảm dần**, để lấy tin mới nhất chứ không phải 40 bài đầu tiên gặp |
| **D-06** | `robots.txt` cấm đường dẫn | **Không fetch.** Ghi log `skipped_by_robots`. Cache `robots.txt` 24h |
| **D-07** | `robots.txt` trả 404 hoặc timeout | Coi như cho phép (chuẩn RFC 9309), nhưng vẫn giữ rate limit |
| **D-08** | Item RSS trỏ tới trang chuyên mục, không phải bài | Nhận diện bằng heuristic (URL không có slug hoặc ID bài, thiếu thẻ `article`) rồi bỏ qua |
| **D-09** | Nguồn phân trang, tin mới nằm ở trang 2 | v1 chỉ lấy trang 1. Khoảng trống này được che bởi cron 8h và `MAX_ARTICLES_PER_SOURCE` đủ lớn. Ghi vào phần hạn chế đã biết |

---

## 2. Tải bài (Fetch)

| ID | Tình huống | Hành vi đúng |
|---|---|---|
| **F-01** | Timeout | Timeout cứng `FETCH_TIMEOUT_SEC` (20s). Retry tối đa 2 lần, backoff mũ có jitter (1s, 3s) |
| **F-02** | HTTP `429 Too Many Requests` | **Không retry ngay.** Tôn trọng `Retry-After`; nếu không có thì bỏ nguồn đó ở run này. Nhân đôi `rate_limit_ms` cho nguồn trong 24h |
| **F-03** | HTTP `403`, bị chặn Cloudflare | Không tìm cách vượt qua. Ghi lỗi, tắt tạm nguồn nếu lặp lại 3 run liên tiếp. **Không giả mạo User-Agent trình duyệt** |
| **F-04** | `5xx` phía nguồn | Retry 2 lần rồi bỏ qua bài đó. Không làm hỏng cả nguồn |
| **F-05** | Redirect chuỗi | Đi theo tối đa 5 redirect. **URL cuối cùng mới là `canonical_url`** |
| **F-06** | Redirect vòng lặp | Phát hiện và cắt, ghi lỗi |
| **F-07** | Bài nằm sau tường phí | Nhận diện (text quá ngắn, có class paywall) rồi đặt `status = 'rejected'`, lý do `paywalled`. Không lưu nội dung một phần |
| **F-08** | Body khổng lồ (> 5MB) | Giới hạn bằng `io.LimitReader`. Chống cạn bộ nhớ |
| **F-09** | Encoding không phải UTF-8 (một số CMS cũ dùng windows-1258) | Phát hiện charset từ header hoặc meta, chuyển sang UTF-8. **Sai bước này là hỏng toàn bộ dấu tiếng Việt** |
| **F-10** | Lỗi chứng chỉ TLS | Không bỏ qua verify. Ghi lỗi |
| **F-11** | Bài là video, podcast, infographic, gần như không có text | `status = 'rejected'`, lý do `insufficient_text` (dưới 200 ký tự) |
| **F-12** | Bài trả về bản AMP | Chuẩn hoá về URL non-AMP bằng `<link rel="canonical">` |

---

## 3. Bóc tách nội dung (Extract)

| ID | Tình huống | Hành vi đúng |
|---|---|---|
| **E-01** | Có `<link rel="canonical">` | **Luôn ưu tiên** giá trị này làm `canonical_url` |
| **E-02** | URL có tham số theo dõi (`utm_*`, `fbclid`, `gclid`, `zarsrc`, `src`) | Loại bỏ trước khi hash. `?utm_source=zalo` và URL trần phải cho ra **cùng một** `url_hash` |
| **E-03** | URL khác nhau chỉ ở fragment `#comment` | Bỏ fragment |
| **E-04** | URL khác nhau ở `http` với `https`, có hoặc không `www`, dấu `/` cuối | Chuẩn hoá: ép `https`, bỏ `www.`, bỏ `/` cuối, hạ chữ thường phần host (**không** hạ chữ thường phần path) |
| **E-05** | Bài có box "Tin liên quan", quảng cáo, form đăng ký | Loại bỏ trước khi tóm tắt. Nếu lọt vào sẽ làm hỏng tóm tắt và làm validator so sai số |
| **E-06** | Bài có bảng số liệu | Giữ nội dung bảng dưới dạng text có cấu trúc, vì đây thường là chỗ chứa số quan trọng nhất |
| **E-07** | Tiêu đề trong `<title>` khác tiêu đề trong `<h1>` | Ưu tiên `og:title`, rồi `<h1>`, cuối cùng `<title>` đã cắt đuôi tên báo |
| **E-08** | Bài dạng live-blog cập nhật liên tục | Nhận diện, đặt `is_live = true`, cho phép cập nhật lại (xem U-03) |
| **E-09** | Bài có nhiều trang (`?page=2`) | v1 chỉ lấy trang 1. Ghi nhận hạn chế |

---

## 4. Phân tích thời gian (đặc thù tiếng Việt)

Đây là chỗ hay sai nhất, và hậu quả là tin hiện sai ngày trên bảng digest.

| ID | Định dạng gặp thực tế | Hành vi đúng |
|---|---|---|
| **T-01** | `"Thứ Tư, 26/8/2026, 08:15 (GMT+7)"` | Parse được, đổi sang UTC |
| **T-02** | `"26/08/2026 08:15"` | Mặc định `Asia/Ho_Chi_Minh` khi không ghi múi giờ |
| **T-03** | `"2 giờ trước"`, `"30 phút trước"`, `"hôm qua"` | Tính ngược từ thời điểm fetch |
| **T-04** | RSS `pubDate` chuẩn RFC 1123 | Parse chuẩn |
| **T-05** | `<meta property="article:published_time">` ISO-8601 | **Ưu tiên cao nhất**, đáng tin hơn text hiển thị |
| **T-06** | Không tìm thấy thời gian nào | Dùng `fetched_at`, đặt cờ `published_at_estimated = true`. FE hiển thị dấu `~` trước ngày |
| **T-07** | Thời gian ở **tương lai** do lỗi CMS nguồn | Nếu lệch quá 2h so với hiện tại thì kẹp về `fetched_at`, ghi cảnh báo |
| **T-08** | Thời gian quá cũ (hơn 90 ngày) trong feed "mới nhất" | Bỏ qua, nhiều khả năng là bài cũ được đẩy lại |
| **T-09** | Bài đăng lúc 23:50 ngày 26/08 giờ VN | Hiển thị `26/08/2026`, **không** phải `27/08`. Mọi hiển thị ngày phải đổi sang `Asia/Ho_Chi_Minh` trước khi cắt ngày |

---

## 5. Trùng lặp và gom cụm

| ID | Tình huống | Hành vi đúng |
|---|---|---|
| **C-01** | Đúng cùng URL, chạy lại lần 2 | `url_hash` UNIQUE chặn. Không tạo bản ghi mới, không tính là lỗi |
| **C-02** | Cùng bài, khác URL do tham số theo dõi | E-02 chuẩn hoá xong ra cùng hash |
| **C-03** | 5 báo cùng đưa một sự kiện, tiêu đề gần giống | Simhash tiêu đề, khoảng cách Hamming tối đa 3 trong cửa sổ 72h thì cùng `cluster_id`. Chọn đại diện theo `sources.tier` thấp nhất; hoà thì chọn bài đăng sớm nhất |
| **C-04** | Cùng tiêu đề nhưng cách nhau 6 tháng (tin định kỳ kiểu "VN-Index đóng cửa...") | Cửa sổ 72h ngăn gom nhầm. **Không được bỏ cửa sổ thời gian** |
| **C-05** | Một báo đăng lại chính bài của mình với URL mới | Cùng nguồn và simhash gần thì đánh dấu `duplicate_of`, không hiện trùng |
| **C-06** | Bài dịch lại từ nguồn nước ngoài, nội dung khác nhưng cùng sự kiện | v1 không bắt được vì simhash chỉ so tiêu đề. Ghi nhận hạn chế |
| **C-07** | Đại diện cụm bị xoá khỏi nguồn | Bầu lại đại diện từ các bài còn lại trong cụm |

---

## 6. Tin được cập nhật sau khi đăng (phần bị thiếu ở thiết kế v1)

Báo Việt Nam **thường xuyên sửa bài sau khi đăng**: sửa số liệu, sửa tiêu đề, bổ sung diễn biến. Thiết kế v1 chỉ xử lý "bài mới", không xử lý "bài đã đổi". Đây là thiếu sót thật.

| ID | Tình huống | Hành vi đúng |
|---|---|---|
| **U-01** | Đã có bài, lần crawl sau nội dung **giống hệt** | Chỉ cập nhật `last_seen_at`. Không tóm tắt lại, tránh tốn chi phí LLM vô ích |
| **U-02** | Nội dung đổi nhẹ (sửa chính tả, khác dưới 5%) | Cập nhật `content_hash`, **không** tóm tắt lại |
| **U-03** | Nội dung đổi đáng kể (từ 5% trở lên, hoặc **số liệu đổi**, hoặc tiêu đề đổi) | Tóm tắt lại, tăng `revision`, ghi `updated_at`. FE hiện nhãn "Đã cập nhật lúc HH:mm" |
| **U-04** | Bài bị xoá khỏi nguồn (`404` hoặc `410`) | Đặt `short_links.target_alive = false`, `articles.status = 'source_gone'`. **Giữ lại tóm tắt** vì đó đã là nội dung của ta, nhưng FE không render link bấm được |
| **U-05** | Bài đổi URL, nội dung giữ nguyên | Bắt qua redirect (F-05) và `content_hash`, cập nhật `canonical_url`, giữ nguyên `article.id` và `short_link.code` để link đã chia sẻ không chết |
| **U-06** | Cách phát hiện đổi nội dung | `content_hash = sha256(normalized_body)`, trong đó chuẩn hoá khoảng trắng, bỏ timestamp động, bỏ số lượt xem và bình luận. **Không chuẩn hoá kỹ thì mọi bài đều bị coi là "đổi" mỗi lần crawl** |
| **U-07** | Bài cũ được sửa nhưng tóm tắt mới lại **fail validator** | Giữ tóm tắt cũ đang publish, đưa bản mới vào `pending_review`. **Không bao giờ thay nội dung đang publish bằng nội dung chưa qua kiểm chứng** |

---

## 7. Gắn mã chứng khoán

Bổ sung cho `architecture.md` mục 8.

| ID | Tình huống | Hành vi đúng |
|---|---|---|
| **G-01** | Từ viết tắt 3 chữ trùng mã: GDP, CPI, ETF, USD, EUR, JPY, CEO, CFO, CTO, IPO, ROE, ROA, EPS, NHNN, HOSE, HNX, UPC, ATC, ATO, FDI, ODA, VAT, EVN, BOT, PPP, ESG, CSR, MW, KWH, EU, WTO, IMF, ADB | Stoplist cứng, không bao giờ coi là mã |
| **G-02** | `VND` | Là mã VNDIRECT **chỉ khi** trong câu có từ khoá cổ phiếu (`cổ phiếu`, `mã`, `khớp lệnh`, `thị giá`, `phiên`). Đứng sau số hoặc cạnh `tỷ`, `đồng` thì là đơn vị tiền |
| **G-03** | Mã nằm trong cụm trích dẫn nguồn: `"theo MBS Research"`, `"ACBS ước tính"`, `"SSI Research cho rằng"` | **Không gắn.** Đây là tên tổ chức phân tích, không phải chủ thể của tin |
| **G-04** | Tên công ty mẹ và con: `Vinhomes` ra `VHM` chứ không phải `VIC`; `Vincom Retail` ra `VRE`; `Vietcombank` ra `VCB` | Cần `negative_aliases` cho mỗi mã |
| **G-05** | Tên có dấu và không dấu: `Hoà Phát`, `Hòa Phát`, `hoa phat` | Chuẩn hoá NFC, bỏ dấu, hạ chữ thường trước khi khớp. **Lưu ý `oà` và `òa` là hai cách gõ khác nhau của cùng một từ** |
| **G-06** | Mã xuất hiện trong URL hoặc tên ảnh | Chỉ quét text đã extract, không quét HTML thô |
| **G-07** | Bài liệt kê 30 mã ở cuối (thống kê phiên) | Tất cả là `mentioned`. `primary` chỉ khi ở tiêu đề, đoạn đầu, hoặc xuất hiện từ 2 lần trở lên |
| **G-08** | Bài nói về ngành, không nhắc mã nào (ví dụ tin CBAM) | Tầng T3 ngữ cảnh ngành gán `mentioned` cho nhóm mã của ngành đó |
| **G-09** | Mã mới niêm yết chưa có trong CSDL | Không gắn. Ghi log `unknown_ticker_candidate` để bổ sung bảng `tickers` |
| **G-10** | Mã bị huỷ niêm yết | Giữ trong CSDL, đặt `delisted_at`. Vẫn gắn cho tin lịch sử |
| **G-11** | Một tin gắn hơn 15 mã | Cảnh báo, nhiều khả năng tagger sai. QA rà soát |

---

## 8. Tóm tắt và kiểm chứng số liệu

| ID | Tình huống | Hành vi đúng |
|---|---|---|
| **S-01** | LLM trả về số không có trong bài gốc | Validator chặn, chuyển `pending_review`. **Không bao giờ publish** |
| **S-02** | LLM tự làm tròn `29,91` thành `30` | Fail, vì làm tròn cũng là sai số |
| **S-03** | Số viết bằng chữ trong bài nhưng tóm tắt viết số | Validator phải chuẩn hoá được cả hai chiều, hoặc coi là fail an toàn |
| **S-04** | `408.000` theo quy ước VN và theo quy ước Anh là khác nhau | Chuẩn hoá theo quy ước VN: `.` phân nhóm nghìn, `,` thập phân |
| **S-05** | API LLM timeout, `429`, hết quota | Fallback ngay sang provider `extractive`. **Pipeline không bao giờ dừng vì LLM** |
| **S-06** | LLM trả về ngôn ngữ khuyến nghị đầu tư | Banned-phrase list chặn ở tầng sinh nội dung, sinh lại 1 lần, vẫn sai thì dùng `extractive` |
| **S-07** | LLM trả markdown ngoài `**bold**` (heading, list, link) | Strip, chỉ giữ `**` và `*` |
| **S-08** | Bài quá dài (hơn 15.000 ký tự) | Cắt phần đầu cộng đoạn có mật độ số cao nhất trước khi đưa vào LLM, để không vượt context và không tốn tiền |
| **S-09** | Tóm tắt dài hơn 180 từ | Reject, sinh lại |
| **S-10** | Tóm tắt trùng gần nguyên văn câu trong bài gốc (từ 15 từ liên tiếp) | Reject, vì vi phạm rule bản quyền ở PRD mục 4 |

---

## 9. Short link

| ID | Tình huống | Hành vi đúng |
|---|---|---|
| **L-01** | Domain đích không thuộc allowlist | Từ chối tạo, ghi cảnh báo. Chống open redirect |
| **L-02** | `cafef.vn.evil.com` | Từ chối. So khớp theo nhãn domain, **không** so khớp chuỗi con |
| **L-03** | `m.cafef.vn`, `amp.cafef.vn` | Chấp nhận vì là subdomain hợp lệ, nhưng chuẩn hoá về domain chính nếu có canonical |
| **L-04** | Đụng mã base62 | Retry tối đa 5 lần rồi báo lỗi |
| **L-05** | Bài đổi URL (U-05) | **Giữ nguyên `code`**, chỉ cập nhật `target_url`. Link đã chia sẻ trên Zalo không được chết |
| **L-06** | Bài gốc chết (U-04) | `/r/{code}` trả `410 Gone` kèm trang tiếng Việt giải thích và link tới trang mã liên quan. **Không** redirect về trang chủ |
| **L-07** | Mã không tồn tại | `404`, không redirect |
| **L-08** | Kiểm tra định kỳ | Job hàng tuần gọi `HEAD` từng `target_url`, cập nhật `target_alive` và `last_checked_at` |

---

## 10. Lập lịch và khoảng trống dữ liệu

| ID | Tình huống | Hành vi đúng |
|---|---|---|
| **R-01** | Run trước chưa xong, đã tới giờ run kế | `pg_try_advisory_lock` chặn. Bỏ qua run mới, ghi log `skipped_overlapping` |
| **R-02** | Server tắt lúc 06:00, bật lại lúc 09:00 | Khi khởi động, nếu `now - last_successful_run` vượt chu kỳ cron thì chạy bù ngay một lần |
| **R-03** | Cách nhau 8h giữa 2 run, nguồn đăng 100 bài trong khoảng đó, giới hạn 40 | Sắp theo thời gian giảm dần rồi mới cắt (D-05). Ghi `truncated = true` vào `crawl_run_sources` để nhìn thấy được khoảng trống |
| **R-04** | Chuyển mùa và đổi múi giờ | Việt Nam không có DST. Nhưng cron **phải** dùng `cron.WithLocation(Asia/Ho_Chi_Minh)`, không dựa vào giờ hệ thống của container |
| **R-05** | Thiếu `tzdata` trong image | `time.LoadLocation` lỗi thì **dừng ngay lúc khởi động** kèm thông báo rõ ràng. Không im lặng rơi về UTC, vì sẽ lệch 7 tiếng |
| **R-06** | Ingest thủ công gọi đúng lúc cron chạy | Trả `409 Conflict`, không xếp hàng |
| **R-07** | Pipeline chạy quá 30 phút | Timeout toàn cục, huỷ context, đóng `crawl_run` với trạng thái `timeout` |

---

## 11. Chịu lỗi và toàn vẹn dữ liệu

| ID | Tình huống | Hành vi đúng |
|---|---|---|
| **I-01** | 1 nguồn panic | `recover()` trong goroutine của nguồn đó. Các nguồn khác chạy tiếp |
| **I-02** | Mất kết nối DB giữa chừng | Retry với backoff. `crawl_run` được đóng trong `defer` để không treo trạng thái `running` mãi |
| **I-03** | Ghi bài, mã CK và short link | Trong **một transaction**. Không để tồn tại bài không có mã hoặc short link mồ côi |
| **I-04** | Process bị kill giữa run | `crawl_runs` còn trạng thái `running` quá 1h thì job dọn đánh dấu `aborted` |
| **I-05** | Bảng `articles` phình to | Giữ 18 tháng. Job dọn hàng tháng chuyển bài cũ hơn sang bảng lưu trữ |
| **I-06** | Hai instance API cùng chạy khi scale | Advisory lock trên DB đảm bảo chỉ một instance chạy pipeline |

---

## 12. Cột cần bổ sung vào schema

Các case trên yêu cầu thêm những cột sau so với `architecture.md` mục 4:

```sql
ALTER TABLE articles
  ADD COLUMN is_demo                BOOLEAN     NOT NULL DEFAULT false,
  ADD COLUMN content_hash           TEXT,
  ADD COLUMN revision               INT         NOT NULL DEFAULT 1,
  ADD COLUMN updated_at             TIMESTAMPTZ,
  ADD COLUMN last_seen_at           TIMESTAMPTZ,
  ADD COLUMN published_at_estimated BOOLEAN     NOT NULL DEFAULT false,
  ADD COLUMN is_live                BOOLEAN     NOT NULL DEFAULT false,
  ADD COLUMN reject_reason          TEXT,
  ADD COLUMN duplicate_of           BIGINT REFERENCES articles(id);

ALTER TABLE sources
  ADD COLUMN list_selector          TEXT,
  ADD COLUMN etag                   TEXT,
  ADD COLUMN last_modified          TEXT,
  ADD COLUMN consecutive_empty_runs INT         NOT NULL DEFAULT 0,
  ADD COLUMN backoff_until          TIMESTAMPTZ;

ALTER TABLE tickers
  ADD COLUMN delisted_at            DATE,
  ADD COLUMN negative_aliases       TEXT[];

ALTER TABLE crawl_run_sources
  ADD COLUMN truncated              BOOLEAN     NOT NULL DEFAULT false,
  ADD COLUMN skipped_by_robots      INT         NOT NULL DEFAULT 0;
```

Và `GET /api/v1/news` bổ sung cho mỗi phần tử:

```jsonc
{
  "source": {
    "name": "CafeF",
    "domain": "cafef.vn",
    "tier": 1,
    "article_title": "VN-Index đóng cửa 1.821 điểm, lần đầu vượt mốc 1.800"
  },
  "short_link": "/r/a7Kx2p",
  "short_link_alive": true,
  "is_demo": false,
  "updated_at": null,
  "published_at_estimated": false
}
```

`source.article_title` là thứ FE dùng để ghi vào `title` và `aria-label` của link. Đây là câu trả lời trực tiếp cho phản hồi *"link chưa đưa đến thông tin trang cụ thể"*.

---

## 13. Thứ tự ưu tiên khi triển khai

| Ưu tiên | Case |
|---|---|
| **P0, làm ngay** | T1, T2, T3 mục 0 · U-01 đến U-07 · E-01 đến E-04 · T-01 đến T-09 · S-01 đến S-06 · L-01, L-02, L-05, L-06 · R-05 · I-01, I-03 |
| **P1** | D-04, D-05 · F-01 đến F-09 · C-01 đến C-05 · G-01 đến G-08 · R-01 đến R-04 · I-02, I-04 |
| **P2** | D-08, D-09 · F-11, F-12 · E-08, E-09 · C-06, C-07 · G-09 đến G-11 · L-08 · I-05, I-06 |
