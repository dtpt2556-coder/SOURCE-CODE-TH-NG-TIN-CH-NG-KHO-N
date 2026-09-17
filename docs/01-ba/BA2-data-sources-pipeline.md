# BA2 — Dữ liệu, Nguồn tin & Content Pipeline (TenPoint)

> **Tài liệu:** BA2 / Góc nhìn Dữ liệu — Nguồn tin — Content Pipeline
> **Sản phẩm:** TenPoint — nền tảng tổng hợp & tóm tắt tin tức chứng khoán Việt Nam
> **Input tham chiếu:** `docs/00-inputs/reference-images.md` (Màn hình A: Research Note, Màn hình B: News Digest Table)
> **Ý tưởng gốc (verbatim):** *"cập nhật tin tức chứng khoán, tổng hợp nội dung và trích dẫn short link, tin tức có cronjob cập nhật 8h 1 lần, nguồn data là những trang uy tín"*
> **Ngày:** 2026-09-16 · **Phiên bản:** v1.0 (draft để PM review)

---

## 0. Phạm vi & giả định

### 0.1 Tài liệu này trả lời

Bảng digest ở Màn hình B có 5 cột: `Ngày` | `Mã CK` | `Tóm tắt thông tin` | `Source` | `Loại tin`. Mỗi cột là một bài toán dữ liệu riêng:

| Cột UI | Bài toán dữ liệu | Mục trong tài liệu |
|---|---|---|
| `Ngày` | Chuẩn hoá `published_at` từ nhiều nguồn, timezone, tin cập nhật | §2 (Step 3), §9 |
| `Mã CK` | Ticker tagging từ văn bản tiếng Việt | §4 |
| `Tóm tắt thông tin` | Summarization spec + rubric + prompt | §5 |
| `Source` | Short link + hiển thị domain rút gọn | §7 |
| `Loại tin` | Taxonomy chuẩn hoá | §3 |
| (ẩn) 1 row = 1 sự kiện | Dedupe & clustering | §6 |

### 0.2 Giả định đã chốt

- **G1.** Cronjob chạy **mỗi 8 giờ** → 3 lần/ngày. Đề xuất khung giờ ICT (UTC+7): **06:00 / 14:00 / 22:00**.
  - `06:00` — digest trước phiên (gom tin đêm + tin quốc tế + CBTT công bố sau 15:00 hôm trước).
  - `14:00` — digest trong phiên (gom tin sáng + diễn biến phiên sáng).
  - `22:00` — digest sau phiên (gom tin chiều + báo cáo phân tích + CBTT cuối ngày).
  - **Fast-lane ngoại lệ:** CBTT bắt buộc từ HOSE/HNX/UBCKNN (§1.3) nên có job riêng 30 phút/lần — vì tin CBTT có tính vật chất (material) và trễ 8 giờ là không chấp nhận được. Cần PM quyết: **fast-lane có trong MVP hay Phase 2**.
- **G2.** TenPoint **không** lưu và **không** hiển thị full-text bài gốc ra public. Chỉ lưu bản raw phục vụ xử lý nội bộ, public chỉ thấy tóm tắt + link gốc (§8).
- **G3.** Ngôn ngữ đầu ra: tiếng Việt. Nguồn nước ngoài (nếu có ở Phase 2) phải dịch + đánh dấu `translated=true`.
- **G4.** Tài liệu này **không** viết code, không chốt tech stack. Mọi thứ ở mức nghiệp vụ + rule.

### 0.3 Cảnh báo về độ chính xác của tài liệu

> ⚠️ **Các URL RSS trong §1 là phỏng đoán theo pattern CMS phổ biến, CHƯA được verify bằng HTTP request.** Trước khi onboard bất kỳ nguồn nào, bắt buộc chạy checklist §1.5 (probe feed, kiểm robots.txt, đọc ToS). Không được coi bảng §1 là sự thật kỹ thuật.
>
> ⚠️ **Phần pháp lý (§8) là phân tích rủi ro nghiệp vụ, KHÔNG phải tư vấn pháp lý.** Các điều luật/nghị định được dẫn cần luật sư Việt Nam xác nhận hiệu lực và cách áp dụng tại thời điểm triển khai — đặc biệt là quy định về "trang thông tin điện tử tổng hợp" vốn thay đổi nhiều lần.

---

## 1. Danh mục nguồn tin uy tín Việt Nam

### 1.1 Nguyên tắc phân tầng

| Tier | Định nghĩa | Chính sách sử dụng |
|---|---|---|
| **Tier 0 — Nguồn sơ cấp** | Cơ quan quản lý & Sở GDCK. Dữ liệu là **bản gốc**, không qua diễn giải. | Luôn ưu tiên làm canonical source. Không cần dedupe với nhau. Fast-lane. |
| **Tier 1 — Báo/chuyên trang tài chính uy tín** | Có giấy phép báo chí, đội ngũ phóng viên tài chính chuyên trách, độ phủ cao, sai sót số liệu thấp. | Nguồn chính của digest. Được chọn làm đại diện cluster. |
| **Tier 2 — Báo tổng hợp / chuyên trang bổ trợ** | Có giấy phép nhưng mảng chứng khoán không phải thế mạnh chính, hoặc tốc độ/độ sâu kém hơn. | Dùng để bổ sung độ phủ, ít khi làm đại diện cluster. |
| **Tier 3 — Nền tảng dữ liệu/cộng đồng** | Không phải cơ quan báo chí. Dữ liệu tốt nhưng rủi ro pháp lý & ToS cao. | Chỉ dùng để **enrich** (giá, EPS, ngành), **không** trích dẫn làm nguồn tin. |

**Thang độ tin cậy 1–5** (dùng cho `source.trust_score`, ảnh hưởng thứ tự chọn đại diện cluster ở §6.4):

| Điểm | Ý nghĩa |
|---|---|
| 5 | Nguồn sơ cấp/pháp lý. Số liệu là bản gốc. Không cần cross-check. |
| 4 | Báo tài chính chuyên sâu, hiếm sai số liệu, thường trích nguồn gốc rõ ràng. |
| 3 | Báo uy tín nhưng đôi khi đưa lại từ nguồn khác, có sai sót đơn vị/quy đổi. |
| 2 | Tốc độ nhanh nhưng chất lượng biên tập không đồng đều, tít giật. |
| 1 | Không đủ tin cậy để publish. Chỉ dùng làm tín hiệu phát hiện sự kiện. |

---

### 1.2 Tier 1 — Nguồn chính

| # | Tên | Domain | RSS? | URL RSS (⚠️ phỏng đoán, cần verify) | Nội dung mạnh | Độ trễ ước tính | Tin cậy | Ghi chú pháp lý |
|---|---|---|---|---|---|---|---|---|
| 1 | **CafeF** (Kênh thông tin kinh tế – tài chính, thuộc VCCorp) | `cafef.vn` | Có | `https://cafef.vn/thi-truong-chung-khoan.rss`<br>`https://cafef.vn/doanh-nghiep.rss`<br>`https://cafef.vn/tai-chinh-ngan-hang.rss`<br>`https://cafef.vn/vi-mo-dau-tu.rss` | **Doanh nghiệp, Thị trường, Khối ngoại.** Độ phủ CBTT doanh nghiệp tốt nhất thị trường. Nhiều bài "tổng kết phiên" giàu số liệu. | 5–30 phút sau sự kiện; bản tin phiên ~15:30 | **4** | Thuộc VCCorp — doanh nghiệp tư nhân, bảo vệ bản quyền chủ động. Cần thỏa thuận hợp tác nội dung. Rủi ro cao nhất về ToS trong Tier 1. |
| 2 | **Vietstock** | `vietstock.vn`<br>`finance.vietstock.vn` | Có (trang RSS tập trung) | `https://vietstock.vn/rss` (trang index liệt kê feed con) | **Doanh nghiệp, Phân tích, CBTT, Cổ tức/Phát hành.** Dữ liệu tài chính doanh nghiệp + lịch sự kiện (GDKHQ, ĐHCĐ) mạnh nhất. | 10–60 phút | **4** | Là công ty dữ liệu tài chính (Vietstock JSC), có sản phẩm thương mại bán dữ liệu → **rủi ro xung đột lợi ích cao nhất**. Phải xin phép rõ ràng trước khi scrape. |
| 3 | **VnExpress — Kinh doanh** | `vnexpress.net` | Có (feed chính thức, đã biết) | `https://vnexpress.net/rss/kinh-doanh.rss`<br>`https://vnexpress.net/rss/kinh-doanh/chung-khoan.rss` | **Vĩ mô, Ngành, Thị trường.** Văn phong chuẩn, ít giật tít, số liệu được biên tập kỹ → rất hợp cho tóm tắt LLM. | 15–60 phút | **4** | Báo điện tử có giấy phép (thuộc Bộ KH&CN). RSS công khai → tín hiệu ngầm cho phép syndication headline. Vẫn phải attribution đầy đủ. |
| 4 | **Tin nhanh Chứng khoán** (Báo Đầu tư Chứng khoán) | `tinnhanhchungkhoan.vn` | Có | `https://www.tinnhanhchungkhoan.vn/rss/home.rss`<br>`https://www.tinnhanhchungkhoan.vn/rss/chung-khoan.rss` | **Thị trường, Doanh nghiệp, Pháp lý, Phân tích.** Bám sát nghiệp vụ TTCK (quy định UBCKNN, xử phạt, nâng hạng). | 15–60 phút | **4** | Thuộc Báo Đầu tư — **cơ quan báo chí của Bộ Tài chính** (trước là Bộ KH&ĐT). Nguồn nhà nước → attribution nghiêm ngặt, nhưng thái độ thường cởi mở hơn với dẫn link. |
| 5 | **Báo Đầu tư** | `baodautu.vn` | Có | `https://baodautu.vn/rss/home.rss`<br>`https://baodautu.vn/rss/doanh-nghiep.rss`<br>`https://baodautu.vn/rss/tai-chinh-chung-khoan.rss` | **Vĩ mô, Ngành, FDI, Dự án hạ tầng.** Nguồn tốt nhất cho tin đầu tư công, FDI, quy hoạch — tác động gián tiếp lên nhóm BĐS KCN, xây dựng, vật liệu. | 30 phút – 3 giờ | **4** | Cơ quan báo chí nhà nước. Như #4. |
| 6 | **VnEconomy** (Tạp chí Kinh tế Việt Nam) | `vneconomy.vn` | Có | `https://vneconomy.vn/chung-khoan.rss`<br>`https://vneconomy.vn/tai-chinh.rss`<br>`https://vneconomy.vn/kinh-te-vi-mo.rss` | **Vĩ mô, Phân tích chính sách.** Bài chuyên sâu về lãi suất, tỷ giá, tín dụng, chính sách tiền tệ. | 1–4 giờ | **4** | Tạp chí điện tử thuộc Hội Khoa học Kinh tế VN. Có tường phí (paywall) một phần → **không được vượt paywall**, chỉ dùng phần công khai. |

---

### 1.3 Tier 0 — Nguồn sơ cấp (CBTT & cơ quan quản lý)

| # | Tên | Domain | RSS? | Cách lấy | Nội dung mạnh | Độ trễ | Tin cậy | Ghi chú pháp lý |
|---|---|---|---|---|---|---|---|---|
| 7 | **HOSE — Sở GDCK TP.HCM** | `hsx.vn` | Không có RSS công khai đã biết | Trang "Công bố thông tin" + endpoint JSON mà chính website dùng. ⚠️ Không có API public chính thức được document → phải verify & xin phép. | **Doanh nghiệp, Cổ tức/Phát hành, Pháp lý.** CBTT bắt buộc, giao dịch nội bộ, thay đổi tỷ lệ sở hữu, đình chỉ/cảnh báo, tự doanh & khối ngoại. | Gần real-time | **5** | Tổ chức nhà nước. Thông tin CBTT về bản chất là **thông tin phải công bố công khai** → rủi ro bản quyền thấp nhất. Nhưng vẫn phải tôn trọng rate limit & ToS website. |
| 8 | **HNX — Sở GDCK Hà Nội** | `hnx.vn` | Không có RSS công khai đã biết | Trang CBTT (HNX + UPCoM) | Như HOSE, thêm **trái phiếu doanh nghiệp** (chuyên trang TPDN) và **đấu thầu TPCP** — tín hiệu quan trọng cho nhóm ngân hàng/BĐS. | Gần real-time | **5** | Như HOSE. |
| 9 | **UBCKNN — Uỷ ban Chứng khoán Nhà nước** | `ssc.gov.vn` | Không có RSS công khai đã biết | Crawl mục "Tin tức – Sự kiện" + "Xử phạt vi phạm hành chính" | **Pháp lý.** Quyết định xử phạt, quy định mới, cấp phép, thu hồi giấy phép. Đây là nguồn duy nhất cho `Loại tin = Pháp lý` cấp cao. | 1–24 giờ | **5** | Văn bản nhà nước — theo Luật SHTT, văn bản quy phạm pháp luật & bản dịch chính thức **không thuộc đối tượng bảo hộ quyền tác giả**. Rủi ro thấp nhất. (Cần luật sư xác nhận.) |
| 10 | *(Bổ sung đề xuất)* **NHNN** | `sbv.gov.vn` | Không | Crawl thông cáo báo chí | **Vĩ mô.** Lãi suất điều hành, tỷ giá trung tâm, room tín dụng, OMO. | 1–24 giờ | **5** | Như #9. |
| 11 | *(Bổ sung đề xuất)* **Tổng cục Thống kê / Cục Thống kê** | `nso.gov.vn` | Không | Crawl báo cáo định kỳ | **Vĩ mô.** CPI, GDP, IIP, PMI, XNK, FDI giải ngân. | Theo lịch công bố cố định | **5** | Như #9. |

> **Ghi chú quan trọng về Tier 0:** vì HOSE/HNX/UBCKNN không có RSS, chi phí kỹ thuật cao hơn (scrape + parse HTML/JSON, dễ vỡ khi site đổi layout). Nhưng đây là **nguồn có giá trị pháp lý và giá trị thông tin cao nhất**, đồng thời **rủi ro bản quyền thấp nhất**. Khuyến nghị: đưa Tier 0 vào MVP dù tốn công hơn, thay vì phụ thuộc hoàn toàn vào báo chí tư nhân.

---

### 1.4 Tier 2 — Nguồn bổ trợ

| # | Tên | Domain | RSS? | URL RSS (⚠️ phỏng đoán) | Nội dung mạnh | Độ trễ | Tin cậy | Ghi chú pháp lý |
|---|---|---|---|---|---|---|---|---|
| 12 | **Thanh Niên** | `thanhnien.vn` | Có | `https://thanhnien.vn/rss/kinh-te.rss`<br>`https://thanhnien.vn/rss/kinh-te/chung-khoan.rss` | **Vĩ mô, Ngành** (nông nghiệp, tiêu dùng, XNK, năng lượng). Xuất hiện trong ảnh tham chiếu ở tin nhập khẩu thịt → DBC/BAF/MML/HAG. | 1–6 giờ | **3** | Cơ quan báo chí của Hội LHTN VN. Có bảo vệ bản quyền chủ động, từng xử lý vi phạm. |
| 13 | **Tuổi Trẻ** | `tuoitre.vn` | Có | `https://tuoitre.vn/rss/kinh-doanh.rss` | **Vĩ mô, Ngành, Pháp lý** (điều tra, vụ việc doanh nghiệp). | 1–6 giờ | **3** | Cơ quan báo chí của Thành Đoàn TP.HCM. Bảo vệ bản quyền chủ động. |
| 14 | **Nhịp sống Thị trường** (kế thừa NDH / Người Đồng Hành) | `markettimes.vn` | Có (khả năng cao) | `https://markettimes.vn/rss/home.rss` | **Doanh nghiệp, Thị trường.** Cùng hệ sinh thái VCCorp với CafeF → **trùng lặp nội dung rất cao với CafeF**. | 10–60 phút | **3** | ⚠️ **Lưu ý dedupe:** cùng nhà với CafeF, nhiều bài đăng chéo gần như nguyên văn. Phải đặt rule dedupe riêng cho cặp `cafef.vn` ↔ `markettimes.vn` (§6.5). |
| 15 | **Người Quan Sát** | `nguoiquansat.vn` | Có (khả năng cao) | `https://nguoiquansat.vn/rss/home.rss` | **Doanh nghiệp, Thị trường.** Tin nhanh về cổ phiếu riêng lẻ, biến động giá, lãnh đạo doanh nghiệp. | 10–60 phút | **2** | Tốc độ tốt nhưng **tít giật, hay suy diễn nhân-quả giá cổ phiếu**. → Nếu chọn làm đại diện cluster, bắt buộc downgrade. Ưu tiên dùng làm tín hiệu phát hiện sự kiện rồi tìm bản Tier 1. |

---

### 1.5 Tier 3 — Nền tảng dữ liệu (enrichment, KHÔNG trích dẫn làm nguồn tin)

| # | Tên | Domain | RSS? | Vai trò trong TenPoint | Tin cậy | Ghi chú pháp lý — ⚠️ RỦI RO CAO |
|---|---|---|---|---|---|---|
| 16 | **Fireant** | `fireant.vn` | Không | **Enrichment:** giá, khối lượng, dữ liệu tài chính, danh sách mã theo ngành. Có API nội bộ mà web app dùng, **không phải API public được cấp phép**. | 3 | ❌ **Không scrape API nội bộ.** Đây là sản phẩm thương mại; dùng API không cấp phép có thể vi phạm ToS + Điều 289 BLHS (truy cập trái phép). **Chỉ dùng nếu có hợp đồng API thương mại.** |
| 17 | **Simplize** | `simplize.vn` | Không | **Enrichment:** hồ sơ doanh nghiệp, phân ngành ICB, báo cáo phân tích tổng hợp. | 3 | ❌ Như #16. Sản phẩm SaaS trả phí. **Chỉ dùng qua hợp đồng.** |
| 18 | *(Thay thế đề xuất)* **Dữ liệu tự build từ Tier 0** | — | — | Danh sách mã niêm yết + tên doanh nghiệp + ngành lấy trực tiếp từ HOSE/HNX là **hợp pháp và miễn phí**. | 5 | ✅ Khuyến nghị dùng cái này cho seed mapping §4.6 thay vì Fireant/Simplize. |

> **Khuyến nghị PM:** Fireant và Simplize nên được coi là **đối tác tiềm năng hoặc đối thủ**, không phải nguồn scrape. Nếu cần dữ liệu giá/tài chính, ưu tiên: (1) tự build từ HOSE/HNX, (2) mua API thương mại (Vietstock, FiinGroup, Wichart), (3) tuyệt đối không scrape lén.

---

### 1.6 Checklist bắt buộc trước khi onboard 1 nguồn

Mỗi nguồn phải pass đủ 8 mục, ghi kết quả vào `source_registry`:

| # | Mục kiểm tra | Tiêu chí pass | Người chịu trách nhiệm |
|---|---|---|---|
| 1 | Probe RSS/feed | HTTP 200, content-type XML/RSS/Atom hợp lệ, ≥10 item | Eng |
| 2 | `robots.txt` | Path cần crawl **không** bị `Disallow` cho UA của ta | Eng |
| 3 | ToS / Điều khoản sử dụng | Đọc & ghi lại điều khoản về sao chép, trích dẫn, tự động hoá | BA + Legal |
| 4 | Sitemap / cấu trúc URL | Xác định canonical URL pattern, có/không AMP | Eng |
| 5 | Chất lượng metadata | Có `og:title`, `article:published_time`, `og:image` không | Eng |
| 6 | Trung bình số liệu/bài | Đếm số lượng số + đơn vị trên 20 bài mẫu (quyết định giá trị digest) | BA |
| 7 | Rate limit an toàn | Xác định ngưỡng (mặc định ≤ 1 req/s, tôn trọng `Crawl-delay`) | Eng |
| 8 | Liên hệ hợp tác | Đã gửi email xin phép/đề nghị hợp tác nội dung chưa | PM |

**Gate:** nguồn Tier 1/Tier 2 **không được đưa lên production** nếu mục 3 và mục 8 chưa xong. (Xem §8.1 — đây là rủi ro pháp lý số 1.)

---

## 2. Content pipeline end-to-end

### 2.1 Sơ đồ tổng thể

```
[CRON 8h/lần: 06:00 / 14:00 / 22:00 ICT]
        │
        ▼
 ① DISCOVER ──▶ ② FETCH ──▶ ③ EXTRACT ──▶ ④ DEDUPE ──▶ ⑤ TICKER TAGGING
                                                              │
                                                              ▼
 ⑨ PUBLISH ◀── ⑧ SHORT LINK ◀── ⑦ SUMMARIZE ◀── ⑥ CLASSIFY
        │
        ├──▶ [HUMAN REVIEW QUEUE]  (bài low-confidence)
        └──▶ [DIGEST TABLE / Màn hình B]

[FAST-LANE 30 phút/lần: HOSE / HNX / UBCKNN] ──▶ ③ (bỏ qua ①② dạng RSS)
```

Mỗi bước có **trạng thái riêng** trên bản ghi `article` → pipeline là state machine, cho phép retry từng bước độc lập mà không chạy lại toàn bộ.

`discovered → fetched → extracted → deduped → tagged → classified → summarized → linked → published`
(+ các trạng thái lỗi: `fetch_failed`, `extract_failed`, `rejected_duplicate`, `needs_review`, `rejected_quality`)

---

### 2.2 Bước ① DISCOVER — phát hiện bài mới

| Mục | Nội dung |
|---|---|
| **Input** | `source_registry` (danh sách nguồn active + feed URL + sitemap URL + last_seen_cursor) |
| **Output** | Danh sách `candidate_url[]` kèm `{source_id, url, title_hint, published_hint, discovered_at}` |
| **Rule nghiệp vụ** | **R1.1** Ưu tiên RSS/Atom. Chỉ crawl sitemap/HTML listing khi nguồn không có feed.<br>**R1.2** Cửa sổ thời gian: chỉ nhận bài có `published_hint` trong **48 giờ** gần nhất. Bài cũ hơn → bỏ (trừ khi là CBTT chưa từng thấy).<br>**R1.3** Conditional request: dùng `If-Modified-Since` / `ETag` để giảm tải cho nguồn.<br>**R1.4** URL đã có trong `article` (theo `url_hash` canonical, §6.1) → bỏ ngay, không fetch.<br>**R1.5** Blacklist path: bỏ các URL chứa `/video/`, `/podcast/`, `/infographic/`, `/emagazine/`, `/rao-vat/`, `/quang-cao/`, `/tag/`, `/chuyen-muc/` — không phải bài text tóm tắt được.<br>**R1.6** Mỗi run có **budget**: tối đa N bài/nguồn (đề xuất N=60) để tránh bùng nổ khi feed lỗi. |
| **Failure mode** | **F1.a** Feed 404/timeout → mark `source.health=degraded`, retry ở run sau, alert nếu fail 3 run liên tiếp.<br>**F1.b** Feed trả về toàn bài cũ (CMS lỗi ngày) → phát hiện bằng check "0 bài mới trong 3 run liên tiếp" → alert.<br>**F1.c** Feed đổi URL âm thầm → redirect 301 phải được follow và **cập nhật lại** `source_registry`, không hard-code.<br>**F1.d** Nguồn chặn IP/rate limit (429) → backoff mũ, giảm concurrency, KHÔNG đổi IP để né (vi phạm ToS, §8). |

---

### 2.3 Bước ② FETCH — tải nội dung gốc

| Mục | Nội dung |
|---|---|
| **Input** | `candidate_url[]` |
| **Output** | `raw_html` (lưu tạm, có TTL), `http_status`, `final_url` (sau redirect), `fetched_at`, `content_hash` |
| **Rule nghiệp vụ** | **R2.1** User-Agent **định danh rõ ràng**: `TenPointBot/1.0 (+https://tenpoint.vn/bot)` — trang `/bot` nêu mục đích + email liên hệ + cách opt-out. Tuyệt đối không giả mạo UA trình duyệt.<br>**R2.2** Rate limit theo **domain**: mặc định 1 req/s, tôn trọng `Crawl-delay` trong robots.txt nếu lớn hơn.<br>**R2.3** Timeout 15s, retry tối đa 2 lần với backoff 2s/8s.<br>**R2.4** Follow redirect tối đa 3 hop; ghi lại `final_url` làm canonical candidate.<br>**R2.5** Nếu bài có paywall (phát hiện qua marker DOM hoặc độ dài body < 500 ký tự trong khi title có) → **dừng, không cố vượt**. Mark `paywalled=true`, loại khỏi digest.<br>**R2.6** `raw_html` **TTL 30 ngày** rồi xoá; chỉ giữ text đã extract + metadata (§8.4).<br>**R2.7** Không render JS trừ khi bắt buộc (Tier 0). Headless browser chỉ bật cho whitelist domain — vì tốn tài nguyên và dễ bị coi là hành vi né chặn. |
| **Failure mode** | **F2.a** 403/429 → backoff, alert nếu lặp lại; cân nhắc liên hệ nguồn xin whitelist thay vì né.<br>**F2.b** 200 nhưng trả về trang captcha/anti-bot → phát hiện qua heuristic (body chứa "cloudflare"/"captcha", độ dài bất thường) → mark `blocked`, KHÔNG giải captcha.<br>**F2.c** Nội dung bị cắt (soft paywall) → xem R2.5.<br>**F2.d** Encoding sai (UTF-8 vs Windows-1258) → chuẩn hoá về UTF-8 NFC; bài có tỷ lệ ký tự lỗi > 5% → `extract_failed`.<br>**F2.e** Link chết sau khi publish (nguồn gỡ bài) → job kiểm link hằng ngày, mark `source_dead=true`, hiển thị badge "bài gốc không còn truy cập được" thay vì link vỡ (§9, KPI broken-link). |

---

### 2.4 Bước ③ EXTRACT — bóc tách nội dung & metadata

| Mục | Nội dung |
|---|---|
| **Input** | `raw_html`, `final_url`, `source_id` |
| **Output** | `article{ title, sapo, body_text, author, published_at, updated_at, section, images[], numbers_index[] }` |
| **Rule nghiệp vụ** | **R3.1 Thứ tự ưu tiên metadata:** JSON-LD (`NewsArticle`) → OpenGraph/`article:published_time` → meta chuẩn của CMS → selector riêng của nguồn → heuristic chung. Selector riêng lưu trong `source_registry`, **không hard-code trong code**.<br>**R3.2 Chuẩn hoá thời gian:** mọi `published_at` quy về **UTC lưu trữ**, hiển thị theo **ICT (UTC+7)**. Nguồn VN mặc định ICT nếu không ghi offset.<br>**R3.3** Nếu `updated_at > published_at + 6h` → coi là **bài cập nhật**, cần re-summarize (xem §6.6).<br>**R3.4 Làm sạch body:** loại bỏ box "Tin liên quan", "Đọc thêm", quảng cáo, caption ảnh, footer bản quyền, khối bình luận, disclaimer của nguồn.<br>**R3.5 `numbers_index` — QUAN TRỌNG NHẤT:** trích xuất **toàn bộ** biểu thức số trong bài thành danh sách có cấu trúc `{raw_text, value_normalized, unit, context_snippet, char_offset}`. Ví dụ: `"408.000 tỷ đồng"` → `{value: 408000, unit: "tỷ VND", raw: "408.000 tỷ đồng"}`. Danh sách này là **cơ sở duy nhất** để kiểm tra LLM có bịa số hay không ở §5.4.<br>**R3.6 Chuẩn hoá định dạng số VN:** dấu `.` = phân cách nghìn, dấu `,` = thập phân. `1.821` = một nghìn tám trăm hai mốt, **không phải** 1.821. Đây là bẫy lớn nhất khi dùng thư viện parse số kiểu Anh-Mỹ.<br>**R3.7** Bài quá ngắn (`body_text` < 400 ký tự) hoặc không có số nào → `rejected_quality`, không đưa vào digest (digest của TenPoint có giá trị nhờ số liệu). |
| **Failure mode** | **F3.a** Nguồn đổi layout → selector vỡ → tỷ lệ `extract_failed` tăng đột biến. **Cần alert theo tỷ lệ, không theo từng bài.** Ngưỡng: >20% bài của 1 nguồn fail trong 1 run → alert.<br>**F3.b** `published_at` sai (CMS trả ngày crawl thay vì ngày đăng) → sanity check: `published_at` không được ở tương lai, không được cũ hơn 365 ngày. Vi phạm → fallback sang ngày trong URL slug hoặc `discovered_at`, mark `date_confidence=low`.<br>**F3.c** Bóc nhầm "Tin liên quan" vào body → tóm tắt lẫn nội dung bài khác → **tóm tắt sai lệch nghiêm trọng**. Guard: đoạn cuối body chứa >3 tiêu đề dạng câu ngắn không có động từ → cắt.<br>**F3.d** Bảng số liệu dạng ảnh (infographic) → số không vào được `numbers_index` → tóm tắt nghèo nàn. Chấp nhận: mark `has_image_data=true`, hạ độ ưu tiên. OCR là Phase 2+. |

---

### 2.5 Bước ④ DEDUPE — khử trùng lặp & gom cụm sự kiện

Chi tiết đầy đủ ở **§6**. Tóm tắt vị trí trong pipeline:

| Mục | Nội dung |
|---|---|
| **Input** | `article` đã extract + index của các bài trong **cửa sổ 72 giờ** |
| **Output** | `cluster_id` cho mỗi bài + cờ `is_representative` cho đúng 1 bài/cluster |
| **Rule nghiệp vụ** | 4 tầng L0→L3 (§6.2). Đặt **trước** summarize để không tốn chi phí LLM cho bài trùng. |
| **Failure mode** | Over-merge (gộp nhầm 2 sự kiện khác nhau → mất tin) hoặc under-merge (digest lặp → giảm chất lượng UI). Xem §6.7. |

---

### 2.6 Bước ⑤ TICKER TAGGING — gắn mã CK

Chi tiết đầy đủ ở **§4**.

| Mục | Nội dung |
|---|---|
| **Input** | `article.title`, `article.sapo`, `article.body_text`, `ticker_master` (bảng seed §4.6) |
| **Output** | `tickers[] = [{code, relevance: primary|secondary|mention, confidence, evidence_offsets[]}]` |
| **Rule nghiệp vụ** | 3 pha: (A) nhận diện ứng viên, (B) khử nhập nhằng, (C) chấm mức độ liên quan. Xem §4.2–§4.5. |
| **Failure mode** | False positive (gắn mã không liên quan → người dùng mất niềm tin, đây là lỗi **nghiêm trọng nhất** của sản phẩm) và false negative (thiếu mã → giảm độ phủ). Ưu tiên **precision > recall** (§9). |

---

### 2.7 Bước ⑥ CLASSIFY — gán `Loại tin`

Chi tiết taxonomy ở **§3**.

| Mục | Nội dung |
|---|---|
| **Input** | `article` + `tickers[]` |
| **Output** | `news_type` (1 giá trị chính, bắt buộc) + `news_type_secondary[]` (0–2 giá trị, tuỳ chọn) + `confidence` |
| **Rule nghiệp vụ** | **R6.1** Rule-based trước (§3.3 — cây quyết định), LLM chỉ xử lý case rule không quyết được.<br>**R6.2** Cột `Loại tin` trên UI chỉ hiện **1 nhãn chính** (khớp ảnh tham chiếu). Nhãn phụ dùng cho filter/search.<br>**R6.3** `confidence < 0.7` → vào `needs_review`, không tự publish. |
| **Failure mode** | **F6.a** Nhãn drift theo thời gian (LLM gán không nhất quán giữa các run) → khoá bằng bộ few-shot cố định + eval set 200 bài đã gán tay.<br>**F6.b** Bài đa chủ đề (vừa vĩ mô vừa ngành) → rule ưu tiên §3.4 quyết định, không để LLM tự do chọn. |

---

### 2.8 Bước ⑦ SUMMARIZE — tóm tắt

Chi tiết spec + prompt ở **§5**.

| Mục | Nội dung |
|---|---|
| **Input** | `article.title + sapo + body_text` (bài đại diện cluster), `numbers_index`, `tickers[]`, `news_type` |
| **Output** | `summary_md` (markdown, có `**bold**` cho số chủ chốt), `numbers_used[]`, `quality_flags[]` |
| **Rule nghiệp vụ** | Rubric §5.2 + kiểm chứng số §5.4. **Không LLM nào được publish thẳng khi chưa qua numeric grounding check.** |
| **Failure mode** | **F7.a** Bịa số / sai đơn vị (nghiêm trọng nhất) → chặn bằng §5.4.<br>**F7.b** Tóm tắt mang tính khuyến nghị đầu tư ("nên mua", "tiềm năng tăng giá") → chặn bằng banned-phrase list §5.3, **rủi ro pháp lý** (§8.5).<br>**F7.c** Vượt độ dài → truncate không được phép (mất số liệu); phải re-generate với hint ngắn hơn.<br>**F7.d** LLM API down/quota → hàng đợi retry, digest có thể chạy trễ; **không fallback sang tóm tắt bằng cách cắt 3 câu đầu** (chất lượng không đạt chuẩn TenPoint). |

---

### 2.9 Bước ⑧ SHORT LINK — sinh liên kết rút gọn

Chi tiết ở **§7**.

| Mục | Nội dung |
|---|---|
| **Input** | `article.canonical_url`, `source_id` |
| **Output** | `short_code`, `short_url`, `display_domain` (eTLD+1, ví dụ `vnexpress.net`) |
| **Rule nghiệp vụ** | R7.1–R7.9 ở §7.2. Điểm cốt lõi: **hiển thị domain thật, redirect qua nội bộ, allowlist đích đến.** |
| **Failure mode** | Open redirect (lỗ hổng bảo mật), short code va chạm, click count bị bot làm nhiễu. Xem §7.4. |

---

### 2.10 Bước ⑨ PUBLISH — xuất bản digest

| Mục | Nội dung |
|---|---|
| **Input** | Tập bài đã `summarized` + `linked` trong cửa sổ của run hiện tại |
| **Output** | `digest` (1 bản/run) gồm `digest_items[]` đã sắp xếp + tiêu đề dạng `"N tin mới đáng chú ý — PVS · FPT · VIC"` |
| **Rule nghiệp vụ** | **R9.1 Gate chất lượng:** chỉ item có `summary_ok = true` AND `ticker_confidence ≥ 0.8` AND `news_type_confidence ≥ 0.7` được auto-publish. Còn lại → `needs_review`.<br>**R9.2 Sắp xếp:** (1) độ quan trọng (`importance_score`), (2) thời gian mới nhất. `importance_score` = f(tier nguồn, số mã ảnh hưởng, vốn hoá mã liên quan, độ lớn con số, có/không thuộc VN30).<br>**R9.3 Giới hạn:** tối đa 12 item/digest — bảng dài hơn làm mất giá trị "cô đọng". Bài không lọt top → vào kho tin, truy cập qua filter.<br>**R9.4 Tiêu đề digest** lấy 3 mã xuất hiện nổi bật nhất trong các item (theo `importance_score`), format `PVS · FPT · VIC` khớp ảnh tham chiếu.<br>**R9.5 Idempotent:** chạy lại cùng run không được tạo digest trùng; dùng `digest_key = (date, slot)`.<br>**R9.6 Disclaimer bắt buộc** hiển thị ở chân mỗi digest và mỗi research note (§8.5).<br>**R9.7 Attribution bắt buộc:** mỗi item phải có tên/domain nguồn + link hoạt động. Item không có link hợp lệ → **không được publish** (đây là ranh giới pháp lý, không phải tuỳ chọn UX). |
| **Failure mode** | **F9.a** Digest rỗng (nguồn chết hàng loạt) → publish digest với thông báo "không có tin đáng chú ý" thay vì trang trắng, + alert P1.<br>**F9.b** Digest toàn tin cùng 1 nguồn → rule đa dạng: tối đa 50% item/digest từ cùng 1 domain.<br>**F9.c** Publish tin sai (số liệu sai, gắn nhầm mã) → cần **quy trình đính chính**: unpublish trong ≤15 phút, ghi log `correction`, hiển thị dấu "đã sửa" nếu tin đã có lượt xem. |

---

### 2.11 Ma trận tóm tắt pipeline

| Bước | Input chính | Output chính | Rule then chốt nhất | Failure mode nguy hiểm nhất | Mức độ |
|---|---|---|---|---|---|
| ① Discover | source_registry | candidate_url[] | R1.2 cửa sổ 48h | Feed chết âm thầm | P2 |
| ② Fetch | candidate_url | raw_html | R2.1 UA định danh | Bị chặn IP / vượt paywall | P1 |
| ③ Extract | raw_html | body + **numbers_index** | R3.5 numbers_index | Lẫn "Tin liên quan" vào body | P1 |
| ④ Dedupe | article + index 72h | cluster_id | L0 canonical URL | Over-merge → mất tin | P1 |
| ⑤ Tag ticker | body + ticker_master | tickers[] | R4.B khử nhập nhằng | False positive mã | **P0** |
| ⑥ Classify | article + tickers | news_type | R6.1 rule trước LLM | Nhãn drift | P3 |
| ⑦ Summarize | body + numbers_index | summary_md | §5.4 numeric grounding | **Bịa số / khuyến nghị đầu tư** | **P0** |
| ⑧ Short link | canonical_url | short_url | R7.5 allowlist đích | Open redirect | **P0** (bảo mật) |
| ⑨ Publish | items đã duyệt | digest | R9.7 attribution bắt buộc | Publish thiếu nguồn | **P0** (pháp lý) |

---

## 3. Taxonomy `Loại tin`

### 3.1 Nguyên tắc thiết kế

1. **Một bài = một nhãn chính.** UI (ảnh tham chiếu) chỉ có chỗ cho 1 nhãn.
2. **Nhãn trả lời câu hỏi: "tin này tác động ở tầng nào?"** — không phải "tin này nói về chủ đề gì".
3. **Tập đóng (closed set).** Không cho LLM tự sinh nhãn mới. Thêm nhãn mới phải qua quyết định sản phẩm.
4. **Tiếng Việt, ngắn, hiển thị vừa 1 dòng** trong cột hẹp.
5. Ảnh tham chiếu chỉ hiện `Ngành` — 5/5 dòng mẫu đều là `Ngành`. Điều này cho thấy nhãn `Ngành` đang bị **over-used**; taxonomy dưới đây tách bạch để tránh điều đó (xem §3.5).

### 3.2 Bảng taxonomy chốt (9 nhãn chính)

| # | Nhãn | Mã enum | Định nghĩa | Ví dụ (từ ảnh tham chiếu & mở rộng) | Quy tắc gán |
|---|---|---|---|---|---|
| 1 | **Vĩ mô** | `MACRO` | Tin về nền kinh tế tổng thể hoặc chính sách điều hành cấp quốc gia/quốc tế. Tác động lên **toàn thị trường**, không giới hạn ngành. | CPI tháng 8 tăng 3,2%; NHNN hạ lãi suất điều hành 0,5 điểm %; Fed giữ nguyên lãi suất; GDP quý III +7,1%; tỷ giá USD/VND vượt 26.000. | Chủ thể là **cơ quan quản lý vĩ mô** (Chính phủ, NHNN, Bộ Tài chính, Fed, ECB) HOẶC chỉ số kinh tế tổng hợp (CPI/GDP/PMI/IIP/XNK/FDI). Không có mã CK nào là chủ thể chính. |
| 2 | **Ngành** | `SECTOR` | Tin tác động lên **một nhóm doanh nghiệp cùng ngành/chuỗi giá trị**, không phải một doanh nghiệp cụ thể. | *(ảnh)* CBAM của EU áp thuế carbon → HPG, HSG, NKG, TIS. *(ảnh)* 12 ngân hàng cam kết 408.000 tỷ tín dụng DNNVV → nhóm ngân hàng. *(ảnh)* Nhập khẩu 494.000 tấn thịt → DBC, BAF, MML, HAG. | Bài gắn **≥3 mã cùng ngành** VÀ không có mã nào là chủ thể duy nhất. HOẶC bài nói rõ về giá hàng hoá đầu vào/đầu ra của một ngành (thép, heo hơi, phân bón, điện, xi măng, dầu). |
| 3 | **Doanh nghiệp** | `COMPANY` | Tin về **một doanh nghiệp niêm yết cụ thể**: kết quả kinh doanh, dự án, nhân sự, M&A, hợp đồng. | CMG Q2/2026 doanh thu 2.323 tỷ (+5,1% YoY), LNST CĐ mẹ 75 tỷ (-20,3%); FPT ký hợp đồng 100 triệu USD; Hoà Phát khởi công Dung Quất 3. | Có **đúng 1 mã** ở mức `primary` (§4.5). Chủ thể ngữ pháp của tiêu đề là tên doanh nghiệp/mã đó. |
| 4 | **Thị trường** | `MARKET` | Diễn biến giao dịch, chỉ số, thanh khoản, dòng tiền của thị trường chung. | *(ảnh)* VN-Index đóng cửa 1.821 điểm (+1,67%), lần đầu vượt 1.800, thanh khoản ~20.000 tỷ, 182 mã tăng/125 giảm. | Tiêu đề/nội dung xoay quanh **chỉ số** (VN-Index, VN30, HNX-Index, UPCoM-Index) hoặc **thanh khoản/độ rộng thị trường**. Các mã nêu trong bài chỉ là minh hoạ → `relevance = mention`. |
| 5 | **Khối ngoại** | `FOREIGN` | Giao dịch của nhà đầu tư nước ngoài, quỹ ETF, cơ cấu danh mục, room ngoại, nâng hạng thị trường. | *(ảnh)* ACBS ước 5.588 tỷ (216 triệu USD) chảy vào 117 cổ phiếu VN kỳ cơ cấu FTSE GEIS; VIC hút 80,67 triệu USD. Khối ngoại bán ròng 43 tỷ trên HOSE. FTSE Russell nâng hạng VN. | Chủ thể là **nhà đầu tư/quỹ nước ngoài** hoặc **tổ chức xếp hạng thị trường** (FTSE, MSCI). ⚠️ Nhãn này **ưu tiên cao hơn** `MARKET` khi bài tập trung vào dòng vốn ngoại (xem §3.4). |
| 6 | **Cổ tức / Phát hành** | `CORP_ACTION` | Sự kiện doanh nghiệp làm thay đổi cấu trúc vốn hoặc quyền lợi cổ đông. | Chốt quyền trả cổ tức 15% tiền mặt; phát hành 200 triệu cp riêng lẻ; chia cổ phiếu thưởng 1:1; mua cổ phiếu quỹ; niêm yết bổ sung; chuyển sàn. | Nội dung chứa hành động doanh nghiệp: cổ tức, GDKHQ, phát hành, chào bán, cổ phiếu quỹ, chia tách, ESOP, niêm yết/huỷ niêm yết. **Tách riêng khỏi `COMPANY`** vì đây là loại tin nhà đầu tư lọc riêng. |
| 7 | **Pháp lý** | `REGULATORY` | Quy định, xử phạt, thanh tra, tố tụng liên quan TTCK hoặc doanh nghiệp niêm yết. | UBCKNN phạt CTCP X 1,5 tỷ đồng do CBTT sai hạn; Nghị định mới về giao dịch T+; khởi tố lãnh đạo doanh nghiệp; đình chỉ giao dịch. | Chủ thể là **cơ quan quản lý/tư pháp** (UBCKNN, HOSE/HNX ra quyết định, thanh tra, toà án, cơ quan điều tra) VÀ đối tượng là doanh nghiệp/TTCK. |
| 8 | **Phân tích** | `ANALYSIS` | Báo cáo/khuyến nghị của công ty chứng khoán, tổ chức nghiên cứu, hoặc bài phân tích chuyên sâu. | MBS Research dự phóng thị trường Data Center VN; SSI Research nâng giá mục tiêu; báo cáo triển vọng ngành ngân hàng 2027. | Nguồn thông tin trong bài là **báo cáo của CTCK/tổ chức NC**. ⚠️ **Rủi ro pháp lý cao** — xem §8.5, cần disclaimer mạnh + không được trình bày như quan điểm của TenPoint. |
| 9 | **Trái phiếu / Tín dụng** | `CREDIT` | Thị trường trái phiếu doanh nghiệp, nợ, xếp hạng tín nhiệm, lãi suất huy động ngành ngân hàng. | Doanh nghiệp BĐS phát hành 5.000 tỷ TPDN; chậm trả gốc/lãi trái phiếu; lãi suất huy động kỳ hạn 12 tháng tăng; nợ xấu ngành ngân hàng. | Nội dung chính là **công cụ nợ** (trái phiếu, tín dụng, nợ xấu, xếp hạng tín nhiệm). Tách khỏi `MACRO` vì tác động trực tiếp lên nhóm BĐS/ngân hàng. |

### 3.3 Cây quyết định gán nhãn (rule-based, chạy trước LLM)

```
START
 │
 ├─ Nguồn là UBCKNN/HOSE/HNX ra quyết định xử phạt, đình chỉ, cảnh báo?
 │    └─ CÓ ──▶ REGULATORY                                    [dừng]
 │
 ├─ Nội dung chứa keyword corp-action?
 │   (cổ tức | GDKHQ | ngày đăng ký cuối cùng | chào bán | phát hành riêng lẻ |
 │    cổ phiếu quỹ | cổ phiếu thưởng | ESOP | chia tách | niêm yết bổ sung | chuyển sàn)
 │    └─ CÓ ──▶ CORP_ACTION                                   [dừng]
 │
 ├─ Nội dung chứa keyword khối ngoại?
 │   (khối ngoại | nhà đầu tư nước ngoài | mua ròng | bán ròng | ETF | FTSE | MSCI |
 │    nâng hạng | room ngoại | quỹ ngoại | cơ cấu danh mục)
 │    VÀ keyword đó xuất hiện ở tiêu đề hoặc sapo?
 │    └─ CÓ ──▶ FOREIGN                                       [dừng]
 │
 ├─ Nội dung chứa keyword nợ/trái phiếu?
 │   (trái phiếu doanh nghiệp | TPDN | chậm trả gốc lãi | xếp hạng tín nhiệm |
 │    nợ xấu | lãi suất huy động | room tín dụng)
 │    └─ CÓ ──▶ CREDIT                                        [dừng]
 │
 ├─ Nội dung trích dẫn báo cáo CTCK/tổ chức NC làm nguồn chính?
 │   (… Research | Chứng khoán … dự báo | báo cáo phân tích | giá mục tiêu | khuyến nghị)
 │    └─ CÓ ──▶ ANALYSIS                                      [dừng]
 │
 ├─ Tiêu đề chứa tên chỉ số? (VN-Index | VN30 | HNX-Index | UPCoM | thanh khoản thị trường)
 │    └─ CÓ ──▶ MARKET                                        [dừng]
 │
 ├─ Đếm ticker ở mức `primary` (§4.5):
 │    ├─ == 1  ──▶ COMPANY                                    [dừng]
 │    ├─ >= 3 và cùng 1 ngành ICB level-2 ──▶ SECTOR          [dừng]
 │    └─ == 0  ──▶ chuyển tiếp ↓
 │
 ├─ Chủ thể là cơ quan vĩ mô (Chính phủ/NHNN/Bộ TC/Fed/ECB) HOẶC chỉ số kinh tế?
 │    └─ CÓ ──▶ MACRO                                         [dừng]
 │
 └─ KHÔNG QUYẾT ĐƯỢC ──▶ gọi LLM classifier (tập đóng 9 nhãn, few-shot cố định)
                          └─ confidence < 0.7 ──▶ needs_review (không auto-publish)
```

### 3.4 Thứ tự ưu tiên khi bài đa chủ đề

Khi nhiều nhãn cùng hợp lệ, áp dụng thứ tự (cao → thấp):

```
REGULATORY > CORP_ACTION > FOREIGN > CREDIT > ANALYSIS > COMPANY > SECTOR > MARKET > MACRO
```

**Lý do thứ tự này:** nhãn càng **cụ thể và hành động được** thì càng ưu tiên. Nhà đầu tư lọc `CORP_ACTION` để biết ngày chốt quyền — nếu tin đó bị gán `COMPANY`, họ sẽ bỏ lỡ. Ngược lại, `MACRO` là nhãn rộng nhất nên đứng cuối.

**Ví dụ áp dụng (lấy từ ảnh tham chiếu):**

| Tin | Nhãn ảnh gốc | Nhãn theo taxonomy này | Lý do |
|---|---|---|---|
| 12 ngân hàng cam kết 408.000 tỷ tín dụng DNNVV | `Ngành` | `CREDIT` (phụ: `SECTOR`) | Nội dung chính là tín dụng/lãi suất → CREDIT ưu tiên cao hơn SECTOR. |
| CBAM áp thuế carbon lên thép/nhôm | `Ngành` | `SECTOR` ✅ | Đúng. ≥3 mã cùng ngành thép, tin về chi phí đầu vào ngành. |
| VN-Index đóng cửa 1.821 điểm… khối ngoại mua ròng | `Ngành` | `MARKET` (phụ: `FOREIGN`) | Tiêu đề xoay quanh chỉ số → MARKET. Khối ngoại chỉ là 1 đoạn trong bài, không ở tiêu đề. |
| Nhập khẩu 494.000 tấn thịt | `Ngành` | `SECTOR` ✅ | Đúng. Giá hàng hoá đầu vào ngành chăn nuôi. |
| ACBS ước 5.588 tỷ vào 117 cp kỳ cơ cấu FTSE GEIS | `Ngành` | `FOREIGN` | Chủ thể là quỹ ngoại + FTSE ở tiêu đề → FOREIGN ưu tiên. |

> **Phát hiện quan trọng cho PM:** trong 5 dòng mẫu của ảnh tham chiếu, **chỉ 2/5 thực sự là `Ngành`**. Nghĩa là hệ thống hiện tại (hoặc mock-up) đang gán nhãn quá thô. Taxonomy này sửa đúng vấn đề đó và làm cột `Loại tin` trở nên hữu ích cho filter.

### 3.5 Quy tắc chống over-use nhãn `SECTOR`

Vì `SECTOR` dễ trở thành "thùng rác", thêm 3 guard:

- **G-S1.** `SECTOR` chỉ hợp lệ khi ≥3 mã `primary`/`secondary` **thuộc cùng 1 ngành ICB level-2**. Mã rải rác nhiều ngành → không phải SECTOR.
- **G-S2.** Nếu tỷ lệ `SECTOR` trong 1 digest > 50% → trigger cảnh báo chất lượng phân loại, review thủ công.
- **G-S3.** Monitoring hằng tuần: phân bố nhãn. Nếu 1 nhãn chiếm >40% tổng tin trong tuần → review lại rule.

---

## 4. Ticker tagging rules

> ⚠️ Đây là bước **P0** của pipeline. Gắn sai mã là lỗi người dùng nhìn thấy ngay và mất niềm tin nhanh nhất. Nguyên tắc xuyên suốt: **thà bỏ sót còn hơn gắn sai** (precision > recall).

### 4.1 Thách thức riêng của tiếng Việt

| # | Vấn đề | Ví dụ cụ thể | Hệ quả nếu không xử lý |
|---|---|---|---|
| 1 | Mã CK VN là 3 chữ cái in hoa → **va chạm với vô số từ viết tắt** | `VND` = VNDIRECT **và** đơn vị tiền Việt Nam Đồng | Mọi bài có "tỷ VND" đều bị gắn nhầm mã VND |
| 2 | Tên viết tắt tổ chức trùng mã | `HOSE`, `HNX`, `GDP`, `CPI`, `FDI`, `ETF`, `IPO`, `USD`, `EUR`, `CEO`, `CFO` | Nhiễu nặng |
| 3 | **Tên CTCK trùng/gần mã** và thường xuất hiện ở vị trí "nguồn trích dẫn" | "Theo **MBS** Research…", "**ACBS** ước tính…", "**SSI** Research dự báo…" | Gắn nhầm MBS/ACB/SSI vào mọi bài có trích báo cáo |
| 4 | Doanh nghiệp được gọi bằng **tên thương hiệu**, không phải mã | "Vingroup", "Hoà Phát", "Techcombank", "Thế Giới Di Động", "Vinamilk" | Bỏ sót mã → độ phủ kém |
| 5 | **Một tên → nhiều mã** trong cùng hệ sinh thái | "Vingroup" → VIC; nhưng "Vinhomes" → VHM, "Vincom Retail" → VRE, "VinFast" → không niêm yết HOSE (VFS trên Nasdaq) | Gắn sai công ty con/mẹ |
| 6 | Mã nhắc thoáng qua vs mã là chủ đề chính | Bài về VN-Index liệt kê 20 mã tăng/giảm | Digest ngập mã, cột `Mã CK` vô nghĩa |
| 7 | Dấu tiếng Việt & viết hoa không chuẩn | "hòa phát" / "Hoà Phát" / "HOA PHAT" / "HÒA PHÁT" | Bỏ sót |
| 8 | Mã trùng tên riêng/địa danh | `HAG` (Hoàng Anh Gia Lai) vs "HAGL"; `SAM`, `TIP`, `PAN`, `BOT`, `CAP`, `HOT` là từ tiếng Anh thông dụng | False positive |

### 4.2 Pha A — Nhận diện ứng viên (candidate generation)

Chạy **3 detector song song**, gộp kết quả:

| Detector | Cách hoạt động | Confidence khởi điểm |
|---|---|---|
| **A1 — Exact ticker** | Regex `\b[A-Z]{3}\b` (và `[A-Z]{3}\d?` cho một số mã UPCoM) trên text **giữ nguyên hoa/thường**. Đối chiếu `ticker_master`. | 0.5 (phải qua Pha B) |
| **A2 — Company name** | Tìm kiếm chuỗi theo `alias[]` của từng mã (bao gồm tên đầy đủ, tên thương hiệu, tên không dấu, tên viết tắt phổ biến). Dùng Aho-Corasick / trie để quét 1 lượt. Chuẩn hoá NFC + lowercase + bỏ dấu để tạo key phụ. | 0.85 |
| **A3 — Pattern trong ngoặc** | Regex `Tên công ty \(([A-Z]{3})\)` hoặc `\(mã: ([A-Z]{3})\)` hoặc `\(HOSE: ([A-Z]{3})\)` — pattern báo chí VN rất hay dùng. | 0.98 (gần như chắc chắn) |

**Thứ tự áp dụng:** A3 > A2 > A1. Nếu A3 xác nhận một mã, mọi lần xuất hiện của mã đó trong bài được nâng confidence lên 0.95.

### 4.3 Pha B — Khử nhập nhằng (disambiguation) — 8 rule

| Rule | Nội dung | Ví dụ |
|---|---|---|
| **B1. Stoplist tuyệt đối** | Danh sách từ viết tắt 3 chữ **không bao giờ** được coi là mã, trừ khi có A3 xác nhận: <br>`HSX, HNX, GDP, CPI, PMI, IIP, FDI, ODA, ETF, IPO, M&A, USD, EUR, JPY, CNY, KRW, CEO, CFO, COO, CTO, HĐQT, BKS, ĐHĐCĐ, ĐHCĐ, CBTT, TTCK, UBCK, NHNN, BTC, KCN, KĐT, BĐS, CNTT, CNC, TPHCM, TPDN, TPCP, OMO, NIM, ROE, ROA, EPS, PEG, EBITDA, YOY, QOQ, CAGR, ESG, CBAM, EVFTA, CPTPP, RCEP, ASEAN, OPEC, IMF, WTO, WHO, EVN, TKV, SCIC, VCCI, VPA, VASEP, VSA` | "GDP quý III tăng 7,1%" → **không** gắn mã |
| **B2. `VND` — xử lý đặc biệt** | `VND` chỉ được coi là mã VNDIRECT khi: (a) đứng cạnh từ khoá cổ phiếu (`cổ phiếu VND`, `mã VND`, `VND tăng trần`), HOẶC (b) có A3 xác nhận, HOẶC (c) `VNDIRECT` xuất hiện trong bài. <br>Ngược lại, nếu đứng sau số hoặc cạnh `tỷ/triệu/nghìn/đồng` → là **đơn vị tiền tệ**, bỏ. | "5.000 tỷ VND" → bỏ. "Cổ phiếu VND tăng 3%" → gắn VND. |
| **B3. Chống nhiễu "nguồn trích dẫn"** | Nếu mã xuất hiện **chỉ** trong cụm attribution — `theo X`, `X Research`, `X ước tính`, `X dự báo`, `Chứng khoán X`, `báo cáo của X`, `nhóm phân tích X` — thì **không gắn**. <br>Danh sách CTCK cần guard: `SSI, VND, HCM, VCI, MBS, BSI, SHS, VIX, FTS, CTS, AGR, ORS, TVS, BVS, APS, PSI, DSC, ACBS, VPBankS, TCBS, KBSV, MAS, VDSC, PHS, YSVN, MBKE` | "Theo MBS Research, thị trường DC…" → **không** gắn MBS |
| **B4. Ngữ cảnh doanh nghiệp bắt buộc** | Với mã do A1 phát hiện (không có A2/A3 hỗ trợ), phải có **≥1 từ khoá ngữ cảnh** trong cửa sổ ±80 ký tự: `cổ phiếu, mã, CP, doanh nghiệp, công ty, tập đoàn, ngân hàng, tăng, giảm, trần, sàn, khớp lệnh, thị giá, vốn hoá, lợi nhuận, doanh thu, cổ đông, LNST, CBTT, HOSE, HNX, UPCoM`. | "VIC +4,31%" → có ngữ cảnh % + tăng → gắn. "Kế hoạch VIC-2030 của tỉnh" → không có ngữ cảnh → bỏ. |
| **B5. Ưu tiên tên công ty hơn mã trần trụi** | Nếu bài có tên công ty (A2) thì mã tương ứng được xác nhận chắc chắn; các mã khác chỉ do A1 phát hiện bị hạ confidence 0.2. | Bài về "Hoà Phát" → HPG chắc chắn; "TIS" xuất hiện 1 lần → cần B4 để xác nhận. |
| **B6. Phân giải hệ sinh thái mẹ–con** | Với các tập đoàn có nhiều mã (Vingroup, Masan, FPT, Gelex, Viettel, Sovico, THACO…), map theo **tên pháp nhân cụ thể nhất**, không map theo tên tập đoàn. Nếu bài nói "hệ sinh thái Vingroup" chung chung mà không nêu pháp nhân → chỉ gắn `VIC` ở mức `mention`. | "Vinhomes ghi nhận…" → VHM (không phải VIC). "Vingroup và các công ty thành viên" → VIC mức `mention`. |
| **B7. Chỉ tính mã đang niêm yết** | Đối chiếu `ticker_master.status = active`. Mã đã huỷ niêm yết/đình chỉ → gắn nhưng đánh cờ `delisted`, hiển thị khác màu. | Tránh gắn mã đã rời sàn. |
| **B8. Giới hạn số mã/bài** | Tối đa **8 mã** hiển thị ở cột `Mã CK`. Nếu nhiều hơn, giữ 8 mã có `relevance` + vốn hoá cao nhất, thêm hậu tố `+N`. | Bài tổng kết phiên có 30 mã → hiện 8 + "+22" |

### 4.4 Pha C — Chấm mức độ liên quan (relevance scoring)

Đây là rule giải quyết trực tiếp câu hỏi *"mã nhắc thoáng qua vs mã là chủ đề chính"*.

| Mức | Tiêu chí | Xuất hiện ở cột `Mã CK`? |
|---|---|---|
| **`primary`** | Mã/tên công ty xuất hiện ở **tiêu đề** HOẶC **sapo**, HOẶC xuất hiện ≥3 lần trong body, HOẶC là **chủ ngữ** của câu chứa số liệu chính của bài. | ✅ Luôn hiện, in trước |
| **`secondary`** | Xuất hiện 2 lần trong body, hoặc 1 lần nhưng kèm số liệu riêng của mã đó. | ✅ Hiện nếu tổng số mã ≤ 8 |
| **`mention`** | Xuất hiện đúng 1 lần, trong danh sách liệt kê, không có số liệu riêng. | ❌ Không hiện ở digest. Chỉ lưu để search/filter. |

**Công thức điểm gợi ý** (để tinh chỉnh ngưỡng bằng dữ liệu, không hard-code từ đầu):

```
relevance_score =
    3.0 × (xuất hiện trong title)
  + 2.0 × (xuất hiện trong sapo)
  + 1.0 × log(1 + số lần xuất hiện trong body)
  + 1.5 × (có số liệu gắn trực tiếp với mã, ví dụ "TCB trần lên 33.450 đ")
  + 1.0 × (là chủ ngữ của câu)
  - 2.0 × (chỉ nằm trong danh sách phân cách bằng dấu phẩy ≥5 phần tử)
  - 3.0 × (nằm trong cụm attribution — rule B3)

primary   nếu score ≥ 4.0
secondary nếu 2.0 ≤ score < 4.0
mention   nếu score < 2.0
```

**Áp dụng cho ví dụ ảnh tham chiếu (dòng VN-Index):**
> "TCB trần lên 33.450 đ khớp 39,17 triệu cp, BCM trần 44.450 đ, VIC +4,31%, FPT +2,69%. Khối ngoại… mua FPT 203 tỷ, TCB 136 tỷ, PNJ 100 tỷ; bán CTG 78 tỷ, ACB 78 tỷ, VHM 61 tỷ."

→ `TCB, BCM, VIC, FPT` có số liệu riêng + xuất hiện nhiều → `secondary`/`primary` → **hiện** (khớp đúng ảnh gốc).
→ `PNJ, CTG, ACB, VHM` chỉ xuất hiện 1 lần trong liệt kê khối ngoại → `mention` → **không hiện** (khớp đúng ảnh gốc, ảnh chỉ liệt kê 4 mã TCB/BCM/VIC/FPT).

✅ **Rule này tái tạo chính xác hành vi trong ảnh tham chiếu** — đây là bằng chứng taxonomy relevance đúng hướng.

### 4.5 Ngưỡng & escalation

| Điều kiện | Hành động |
|---|---|
| Tất cả mã có `confidence ≥ 0.85` | Auto-publish |
| Có mã `0.6 ≤ confidence < 0.85` | Gắn nhưng đưa bài vào `needs_review` |
| Có mã `confidence < 0.6` | Bỏ mã đó, ghi log để cải thiện rule |
| Bài **không có mã nào** nhưng `news_type ∈ {MACRO, MARKET}` | Hợp lệ — cột `Mã CK` để trống hoặc ghi `VN-Index` |
| Bài **không có mã nào** và `news_type ∈ {COMPANY, SECTOR}` | ⚠️ Mâu thuẫn logic → `needs_review` bắt buộc |

### 4.6 Bảng seed mapping — 50 mã (VN30 + mã có tính thời sự cao)

> ⚠️ **Rổ VN30 thay đổi định kỳ (review 2 lần/năm).** Bảng dưới là **seed để khởi tạo**, không phải danh sách VN30 chính thức tại thời điểm đọc. Trước khi dùng production, **phải đồng bộ từ HOSE** (§1.3) và set up job cập nhật khi HOSE công bố rổ mới.
>
> Cột `Alias` là **các chuỗi cần đưa vào detector A2**. Cần bổ sung thêm biến thể không dấu và viết thường khi build index.

| # | Mã | Tên pháp nhân | Alias chính (cho A2) | Ngành (ICB L2) | Ghi chú nhập nhằng |
|---|---|---|---|---|---|
| 1 | **VIC** | Tập đoàn Vingroup – CTCP | Vingroup, Tập đoàn Vingroup | Bất động sản / Đa ngành | ⚠️ Không map "Vinhomes"/"Vincom Retail"/"VinFast" vào VIC (rule B6) |
| 2 | **VHM** | CTCP Vinhomes | Vinhomes | Bất động sản | Công ty con của VIC |
| 3 | **VRE** | CTCP Vincom Retail | Vincom Retail, Vincom | Bất động sản bán lẻ | "Vincom" cũng là tên TTTM → cần ngữ cảnh |
| 4 | **VCB** | NH TMCP Ngoại Thương Việt Nam | Vietcombank, NH Ngoại Thương, VCB | Ngân hàng | |
| 5 | **BID** | NH TMCP Đầu tư và Phát triển Việt Nam | BIDV, NH Đầu tư và Phát triển | Ngân hàng | |
| 6 | **CTG** | NH TMCP Công Thương Việt Nam | VietinBank, NH Công Thương | Ngân hàng | ⚠️ "Vietinbank" vs "Vietcombank" dễ nhầm khi fuzzy match → **không dùng fuzzy**, chỉ exact |
| 7 | **TCB** | NH TMCP Kỹ Thương Việt Nam | Techcombank, NH Kỹ Thương | Ngân hàng | ⚠️ "TCBS" (Chứng khoán Kỹ Thương) ≠ TCB → rule B3 |
| 8 | **MBB** | NH TMCP Quân Đội | MB Bank, MBBank, Ngân hàng Quân Đội, MB | Ngân hàng | ⚠️ "MB" quá ngắn → chỉ match "MB Bank"/"MBBank". "MBS" là CTCK con → rule B3 |
| 9 | **VPB** | NH TMCP Việt Nam Thịnh Vượng | VPBank, NH Việt Nam Thịnh Vượng | Ngân hàng | ⚠️ "VPBankS" là CTCK con → rule B3 |
| 10 | **ACB** | NH TMCP Á Châu | ACB, Ngân hàng Á Châu | Ngân hàng | ⚠️ **"ACBS" là CTCK con** — xuất hiện trong ảnh tham chiếu ở vai trò nguồn trích dẫn → rule B3 bắt buộc |
| 11 | **HDB** | NH TMCP Phát triển TP.HCM | HDBank, HD Bank | Ngân hàng | ⚠️ "HD Saison", "HDBS" ≠ HDB |
| 12 | **STB** | NH TMCP Sài Gòn Thương Tín | Sacombank, NH Sài Gòn Thương Tín | Ngân hàng | ⚠️ "SBS"/"SBB" ≠ STB |
| 13 | **SHB** | NH TMCP Sài Gòn – Hà Nội | SHB, NH Sài Gòn Hà Nội | Ngân hàng | ⚠️ "SHS" (Chứng khoán SHS) ≠ SHB |
| 14 | **VIB** | NH TMCP Quốc tế Việt Nam | VIB, NH Quốc tế | Ngân hàng | |
| 15 | **TPB** | NH TMCP Tiên Phong | TPBank, NH Tiên Phong | Ngân hàng | ⚠️ "TPS" (Chứng khoán Tiên Phong) ≠ TPB |
| 16 | **LPB** | NH TMCP Lộc Phát Việt Nam | LPBank, Lộc Phát, LienVietPostBank | Ngân hàng | Đã đổi tên nhiều lần → giữ cả alias cũ |
| 17 | **SSB** | NH TMCP Đông Nam Á | SeABank, NH Đông Nam Á | Ngân hàng | |
| 18 | **MSB** | NH TMCP Hàng Hải Việt Nam | MSB, Maritime Bank, NH Hàng Hải | Ngân hàng | Xuất hiện trong ảnh tham chiếu |
| 19 | **HPG** | CTCP Tập đoàn Hòa Phát | Hòa Phát, Hoà Phát, Tập đoàn Hòa Phát | Tài nguyên cơ bản (Thép) | ⚠️ Cần alias cả "Hòa"(dấu hỏi) và "Hoà"(dấu huyền) — 2 cách gõ khác nhau |
| 20 | **HSG** | CTCP Tập đoàn Hoa Sen | Hoa Sen, Tôn Hoa Sen | Thép | ⚠️ "Hoa Sen" vs "Hòa Phát" khác nhau, **không fuzzy** |
| 21 | **NKG** | CTCP Thép Nam Kim | Nam Kim, Thép Nam Kim | Thép | |
| 22 | **TIS** | CTCP Gang thép Thái Nguyên | Gang thép Thái Nguyên, Tisco | Thép | Xuất hiện trong ảnh tham chiếu (tin CBAM) |
| 23 | **FPT** | CTCP FPT | FPT, Tập đoàn FPT | Công nghệ | ⚠️ "FPT Retail" → FRT, "FPT Telecom" → FOX. Rule B6 |
| 24 | **CMG** | CTCP Tập đoàn Công nghệ CMC | CMC, Tập đoàn CMC, CMC Telecom, CMC Cloud | Công nghệ | Chủ thể của Research Note trong ảnh tham chiếu |
| 25 | **MWG** | CTCP Đầu tư Thế Giới Di Động | Thế Giới Di Động, TGDĐ, Bách Hoá Xanh, Điện Máy Xanh | Bán lẻ | Alias theo chuỗi bán lẻ rất quan trọng |
| 26 | **MSN** | CTCP Tập đoàn Masan | Masan, Tập đoàn Masan, WinMart, WinCommerce | Thực phẩm & đồ uống | ⚠️ "Masan Consumer" → MCH, "Masan MEATLife" → MML. Rule B6 |
| 27 | **MML** | CTCP Masan MEATLife | Masan MEATLife, MEATDeli | Thực phẩm | Xuất hiện trong ảnh tham chiếu |
| 28 | **VNM** | CTCP Sữa Việt Nam | Vinamilk, Sữa Việt Nam | Thực phẩm & đồ uống | |
| 29 | **SAB** | Tổng CTCP Bia – Rượu – NGK Sài Gòn | Sabeco, Bia Sài Gòn, Bia 333 | Đồ uống | |
| 30 | **GAS** | Tổng Công ty Khí Việt Nam – CTCP | PV GAS, Khí Việt Nam, PVGas | Dầu khí | |
| 31 | **PLX** | Tập đoàn Xăng dầu Việt Nam | Petrolimex, Tập đoàn Xăng dầu | Dầu khí | |
| 32 | **PVS** | Tổng CTCP Dịch vụ Kỹ thuật Dầu khí VN | PTSC, Dịch vụ Kỹ thuật Dầu khí | Dầu khí | Xuất hiện ở tiêu đề digest trong ảnh |
| 33 | **PVD** | Tổng CTCP Khoan và Dịch vụ Khoan Dầu khí | PV Drilling, Khoan Dầu khí | Dầu khí | |
| 34 | **POW** | Tổng Công ty Điện lực Dầu khí VN – CTCP | PV Power, Điện lực Dầu khí | Điện | ⚠️ "POW" là từ tiếng Anh → bắt buộc rule B4 |
| 35 | **REE** | CTCP Cơ Điện Lạnh | REE, Cơ Điện Lạnh | Điện / Công nghiệp | |
| 36 | **GEX** | CTCP Tập đoàn GELEX | GELEX, Gelex | Công nghiệp đa ngành | |
| 37 | **GVR** | Tập đoàn Công nghiệp Cao su Việt Nam – CTCP | Cao su Việt Nam, VRG | Cao su / BĐS KCN | |
| 38 | **BCM** | Tổng CT Đầu tư và Phát triển Công nghiệp – CTCP | Becamex, Becamex IDC | BĐS khu công nghiệp | Xuất hiện trong ảnh tham chiếu |
| 39 | **KDH** | CTCP Đầu tư và Kinh doanh Nhà Khang Điền | Khang Điền, Nhà Khang Điền | Bất động sản | |
| 40 | **DXG** | CTCP Tập đoàn Đất Xanh | Đất Xanh, Tập đoàn Đất Xanh | Bất động sản | ⚠️ "Đất Xanh Services" → DXS. Rule B6 |
| 41 | **NLG** | CTCP Đầu tư Nam Long | Nam Long, Đầu tư Nam Long | Bất động sản | |
| 42 | **VJC** | CTCP Hàng không VietJet | VietJet, Vietjet Air | Hàng không | |
| 43 | **ACV** | Tổng CT Cảng hàng không Việt Nam – CTCP | ACV, Cảng hàng không Việt Nam | Hạ tầng hàng không | |
| 44 | **BVH** | Tập đoàn Bảo Việt | Bảo Việt, Tập đoàn Bảo Việt | Bảo hiểm | ⚠️ "Bảo Việt Bank"/"BVBank" → BVB. Rule B6 |
| 45 | **SSI** | CTCP Chứng khoán SSI | SSI, Chứng khoán SSI | Dịch vụ tài chính | ⚠️ **"SSI Research" là nguồn trích dẫn** → rule B3 |
| 46 | **VND** | CTCP Chứng khoán VNDIRECT | VNDIRECT, VnDirect | Dịch vụ tài chính | ⚠️⚠️ **Trùng với đơn vị tiền VND** → rule B2 bắt buộc |
| 47 | **VCI** | CTCP Chứng khoán Vietcap | Vietcap, Bản Việt, VCSC | Dịch vụ tài chính | ⚠️ "Vietcap" vs "Vietcombank" khác nhau |
| 48 | **DGC** | CTCP Tập đoàn Hóa chất Đức Giang | Đức Giang, Hóa chất Đức Giang | Hóa chất | |
| 49 | **DBC** | CTCP Tập đoàn Dabaco Việt Nam | Dabaco, Tập đoàn Dabaco | Chăn nuôi / Thực phẩm | Xuất hiện trong ảnh tham chiếu |
| 50 | **BAF** | CTCP Nông nghiệp BAF Việt Nam | BAF, Nông nghiệp BAF, BAF Meat | Chăn nuôi | ⚠️ "BAF" cũng là từ viết tắt khác → rule B4. Xuất hiện trong ảnh tham chiếu |

**Mã bổ sung cần có ngay (xuất hiện trong ảnh tham chiếu nhưng chưa liệt kê chi tiết):**
`PNJ` (Vàng bạc Đá quý Phú Nhuận – Bán lẻ), `HAG` (Hoàng Anh Gia Lai – Nông nghiệp; alias "HAGL"), `NAB` (Nam A Bank), `NVB` (NCB), `SGB` (Saigonbank), `BVB` (BVBank), `Agribank` (⚠️ **chưa niêm yết** — phải có cờ `not_listed` để không tạo mã ảo).

> **Lưu ý quan trọng:** ảnh tham chiếu có dòng liệt kê "Agribank 70.000 tỷ" — Agribank **không niêm yết**. Bảng `ticker_master` cần có mục `unlisted_entities` để nhận diện và **cố tình không gắn mã**, tránh việc hệ thống "sáng tạo" ra mã không tồn tại.

### 4.7 Cấu trúc `ticker_master` tối thiểu

| Trường | Mô tả |
|---|---|
| `code` | Mã CK (PK) |
| `exchange` | HOSE / HNX / UPCOM |
| `legal_name` | Tên pháp nhân đầy đủ |
| `aliases[]` | Danh sách chuỗi nhận diện (có dấu, không dấu, viết tắt, tên thương hiệu, tên chuỗi bán lẻ) |
| `negative_aliases[]` | Chuỗi **không** được map vào mã này (ví dụ VIC: `["Vinhomes","VinFast","Vincom Retail"]`) |
| `icb_l2`, `icb_l4` | Ngành, dùng cho rule SECTOR (§3.5) |
| `market_cap_bucket` | large/mid/small — dùng cho `importance_score` (R9.2) |
| `in_vn30` | Boolean, đồng bộ theo kỳ review |
| `status` | active / suspended / delisted |
| `ambiguity_level` | none / medium / high — mã `high` (VND, POW, BAF, SAM, TIP…) bắt buộc qua rule B4 |

---

## 5. Chuẩn tóm tắt (Summarization spec)

### 5.1 Đặc điểm rút ra từ ảnh tham chiếu

Phân tích 5 ô `Tóm tắt thông tin` trong ảnh:

| Đặc điểm | Quan sát | Kết luận spec |
|---|---|---|
| Độ dài | ~85–110 từ (≈ 450–650 ký tự), 4–6 câu | Target **80–120 từ**, hard cap 150 |
| Mật độ số | Mỗi tóm tắt chứa **8–18 con số** có đơn vị | Bắt buộc giữ tối đa số liệu — đây là USP |
| Bold | **Chính xác 1 cụm bold/tóm tắt**, luôn là **con số quan trọng nhất**, luôn ở **câu đầu** | Rule: **đúng 1 cụm bold**, đặt ở câu 1 |
| Cấu trúc | Câu 1 = số liệu headline. Câu 2–5 = chi tiết bóc tách, so sánh, bối cảnh | Template 3 tầng (§5.2) |
| Giọng văn | Thuần mô tả, không tính từ cảm thán, không "tuy nhiên có thể thấy…", không khuyến nghị | Trung lập tuyệt đối |
| Đơn vị | `tỷ đồng`, `%`, `điểm %`, `USD/tấn`, `triệu USD`, `tấn CO2/tấn`, `đ/kg`, `triệu cp`, `điểm` | Giữ nguyên đơn vị gốc, **không tự quy đổi** |
| Không có | Không có câu mở đầu thừa, không "theo thông tin từ…", không tên tác giả | Cấm câu dẫn nhập |

### 5.2 Rubric chất lượng (thang chấm, dùng cho cả LLM-judge và review tay)

| # | Tiêu chí | Trọng số | Đạt (2đ) | Một phần (1đ) | Không đạt (0đ) |
|---|---|---|---|---|---|
| **C1** | **Trung thực số liệu** | ×5 | Mọi số trong tóm tắt đều khớp `numbers_index`, đúng giá trị + đơn vị | Có số đúng giá trị nhưng sai/thiếu đơn vị | Có ≥1 số **không tồn tại trong bài gốc** → **REJECT tự động** |
| **C2** | **Mật độ số liệu** | ×4 | ≥6 con số có đơn vị | 3–5 con số | <3 con số |
| **C3** | **Bám đúng ý bài gốc** | ×4 | Không thêm nhận định ngoài bài, không đảo nhân–quả | Có 1 suy diễn nhẹ | Xuyên tạc ý bài gốc |
| **C4** | **Độ dài** | ×2 | 80–120 từ | 60–79 hoặc 121–150 từ | <60 hoặc >150 từ |
| **C5** | **Định dạng bold** | ×2 | Đúng 1 cụm `**...**`, là số quan trọng nhất, ở câu đầu | 1 cụm nhưng không phải số quan trọng nhất | 0 cụm, hoặc ≥2 cụm |
| **C6** | **Giọng văn trung lập** | ×4 | Không từ khuyến nghị, không tính từ cảm thán | Có tính từ cảm thán nhẹ ("ấn tượng", "mạnh mẽ") | Có từ khuyến nghị → **REJECT tự động** (§8.5) |
| **C7** | **Tự đứng vững (standalone)** | ×3 | Đọc riêng vẫn hiểu, không cần mở bài gốc | Thiếu 1 ngữ cảnh nhỏ | Có đại từ/tham chiếu mơ hồ ("công ty này", "như đã nêu") |
| **C8** | **Định dạng số VN** | ×2 | Dấu `.` phân cách nghìn, `,` thập phân, đúng chuẩn VN | Lẫn lộn 1 chỗ | Dùng chuẩn Anh-Mỹ |

**Ngưỡng:** tổng điểm tối đa 52. **Pass ≥ 42** (≈80%). **C1 = 0 hoặc C6 = 0 → reject bất kể tổng điểm.**

### 5.3 Ràng buộc cứng

**Bắt buộc có:**
- Ít nhất 3 con số có đơn vị.
- Đúng 1 cụm `**bold**`.
- Câu đầu tiên chứa con số headline.

**Cấm tuyệt đối (banned patterns — chặn bằng regex sau khi LLM sinh):**

| Nhóm | Chuỗi cấm |
|---|---|
| Khuyến nghị đầu tư (⚠️ pháp lý, §8.5) | `nên mua`, `nên bán`, `khuyến nghị mua`, `khuyến nghị bán`, `khuyến nghị nắm giữ`, `giá mục tiêu` (trừ khi trích dẫn có ghi rõ tổ chức), `tiềm năng tăng giá`, `cơ hội đầu tư`, `đáng để đầu tư`, `nhà đầu tư nên`, `canh mua`, `bắt đáy`, `chốt lời` |
| Dự đoán không có nguồn | `sẽ tăng`, `sẽ giảm`, `dự kiến sẽ bứt phá`, `nhiều khả năng` (khi không gán cho tổ chức cụ thể) |
| Câu dẫn thừa | `Theo thông tin từ`, `Bài viết cho biết`, `Như vậy có thể thấy`, `Tóm lại`, `Nhìn chung` |
| Cảm thán / clickbait | `gây sốc`, `bùng nổ`, `lao dốc không phanh`, `bốc hơi`, `thổi bay`, `ấn tượng`, `đáng kinh ngạc`, emoji |
| Ngôi thứ nhất | `chúng tôi`, `chúng ta`, `tôi` |

### 5.4 Numeric grounding check — cơ chế chống bịa số (BẮT BUỘC)

Đây là guard quan trọng nhất của toàn hệ thống.

```
INPUT : summary_md, numbers_index (từ §2.4 R3.5)
PROCESS:
  1. Trích toàn bộ biểu thức số từ summary_md → summary_numbers[]
     (regex bắt: số + dấu phân cách VN + đơn vị, gồm cả % và điểm %)
  2. Với mỗi n ∈ summary_numbers:
       a. Tìm khớp CHÍNH XÁC (value + unit) trong numbers_index → PASS
       b. Nếu không khớp chính xác, kiểm tra có phải phép quy đổi hợp lệ
          đã được LLM khai báo trong numbers_used[] không
          (ví dụ "216 triệu USD" ↔ "5.588 tỷ đồng" có trong bài gốc cả 2) → PASS
       c. Nếu không → FAIL, ghi rõ số vi phạm
  3. Nếu có bất kỳ FAIL nào:
       → reject summary
       → retry tối đa 2 lần với prompt bổ sung "số X không có trong bài gốc, hãy bỏ"
       → vẫn fail → needs_review, KHÔNG auto-publish
OUTPUT: { grounded: bool, violations: [...] }
```

**Quy tắc phụ:**
- **N1.** LLM **không được tự quy đổi** đơn vị (tỷ ↔ triệu USD, tấn ↔ nghìn tấn). Chỉ dùng cặp giá trị đã có sẵn trong bài gốc.
- **N2.** LLM **không được tự tính** phần trăm thay đổi. Chỉ dùng % đã có trong bài.
- **N3.** LLM **không được làm tròn**. `29,91 điểm` phải giữ nguyên, không thành `~30 điểm`.
- **N4.** Nếu bài gốc ghi khoảng (`500-550 USD/tấn`), giữ nguyên khoảng.

### 5.5 Prompt template cho LLM tóm tắt tiếng Việt

#### System prompt

```
Bạn là biên tập viên tài chính của TenPoint, chuyên tóm tắt tin tức chứng khoán
Việt Nam thành các bản tin cô đọng, giàu số liệu, dành cho nhà đầu tư chuyên nghiệp.

NHIỆM VỤ
Viết một đoạn tóm tắt tiếng Việt từ bài báo được cung cấp.

QUY TẮC BẤT BIẾN — vi phạm bất kỳ quy tắc nào dưới đây là lỗi nghiêm trọng:

1. TRUNG THỰC SỐ LIỆU
   - Chỉ được dùng những con số XUẤT HIỆN NGUYÊN VĂN trong bài gốc.
   - TUYỆT ĐỐI KHÔNG tự tính toán, không tự quy đổi đơn vị, không làm tròn,
     không suy ra tỷ lệ phần trăm mới.
   - Giữ nguyên đơn vị gốc: tỷ đồng, nghìn tỷ đồng, triệu USD, %, điểm %,
     USD/tấn, đồng/kg, triệu cp, điểm, tấn CO2/tấn.
   - Nếu không chắc một con số có trong bài, hãy BỎ con số đó.

2. ĐỊNH DẠNG SỐ KIỂU VIỆT NAM
   - Dấu chấm phân cách hàng nghìn: 408.000 tỷ đồng
   - Dấu phẩy phân cách thập phân: 1,67% ; 29,91 điểm
   - Giữ nguyên cách viết của bài gốc.

3. ĐỘ DÀI
   - 80 đến 120 từ. Tối đa 150 từ. Viết thành 4–6 câu liền mạch, KHÔNG bullet.

4. CẤU TRÚC
   - Câu 1: con số/sự kiện quan trọng nhất, có bối cảnh đủ để hiểu độc lập.
   - Câu 2–5: bóc tách chi tiết — cấu phần, so sánh cùng kỳ, các bên liên quan,
     tác động cụ thể lên doanh nghiệp/ngành.
   - Ưu tiên chi tiết CÓ SỐ hơn chi tiết định tính.

5. IN ĐẬM
   - Dùng ĐÚNG MỘT cụm in đậm bằng markdown **...**
   - Cụm đó phải nằm trong câu đầu tiên và phải là con số quan trọng nhất.
   - Ví dụ: **408.000 tỷ đồng** / **1.821 điểm (+29,91 điểm, +1,67%)**

6. GIỌNG VĂN
   - Trung lập, mô tả, như thông tấn xã. Không tính từ cảm thán.
   - KHÔNG viết câu dẫn nhập ("Theo bài viết...", "Như vậy...", "Nhìn chung...").
   - KHÔNG dùng ngôi thứ nhất. KHÔNG emoji. KHÔNG câu hỏi.

7. CẤM KHUYẾN NGHỊ ĐẦU TƯ — QUY ĐỊNH PHÁP LÝ
   - TUYỆT ĐỐI KHÔNG viết: nên mua, nên bán, nên nắm giữ, khuyến nghị,
     tiềm năng tăng giá, cơ hội đầu tư, canh mua, chốt lời, bắt đáy.
   - Nếu bài gốc có khuyến nghị của một công ty chứng khoán, chỉ được TƯỜNG THUẬT
     kèm tên tổ chức đó (ví dụ: "SSI Research đưa giá mục tiêu 45.000 đồng/cp"),
     không được trình bày như quan điểm của người viết.

8. KHÔNG SUY DIỄN
   - Không giải thích nguyên nhân biến động giá nếu bài gốc không nêu.
   - Không gán nhân–quả không có trong bài.

ĐẦU RA
Trả về DUY NHẤT một JSON hợp lệ, không kèm giải thích:
{
  "summary_md": "<đoạn tóm tắt, có đúng 1 cụm **bold**>",
  "numbers_used": [
    {"raw": "<chuỗi số nguyên văn trong tóm tắt>",
     "source_snippet": "<đoạn ngắn trong bài gốc chứa đúng con số này>"}
  ],
  "headline_number": "<con số đã in đậm>",
  "confidence": <0.0-1.0>,
  "omitted_because_uncertain": ["<số đã bỏ vì không chắc>"]
}
```

#### User prompt (template)

```
=== METADATA ===
Nguồn: {source_name} ({source_domain})
Ngày đăng: {published_at_ict}
Mã CK liên quan (đã xác định): {tickers_primary}
Loại tin (đã phân loại): {news_type}

=== TIÊU ĐỀ ===
{title}

=== SAPO ===
{sapo}

=== NỘI DUNG BÀI GỐC ===
{body_text}

=== DANH SÁCH SỐ LIỆU ĐÃ TRÍCH XUẤT TỪ BÀI GỐC ===
(Bạn CHỈ được dùng các con số trong danh sách này. Không dùng số nào khác.)
{numbers_index_formatted}

=== YÊU CẦU ===
Viết tóm tắt theo đúng quy tắc trong system prompt.
Ưu tiên các số liệu liên quan trực tiếp tới {tickers_primary}.
```

#### Few-shot example (gắn cố định vào prompt để khoá phong cách)

```
VÍ DỤ ĐẦU VÀO (rút gọn):
Tiêu đề: "CBAM: Doanh nghiệp thép, nhôm Việt đối mặt chi phí carbon lớn khi xuất sang EU"
Số liệu trích xuất: [hơn 100 USD/tấn; cả nghìn USD/tấn; 75 USD/tấn CO2;
100.000 tấn; 500-550 USD/tấn; khoảng 300 tỷ đồng; 20%; 2,5 tấn CO2/tấn;
14 tấn; 1,4 tấn; sáu nhóm; đầu 2026]

VÍ DỤ ĐẦU RA:
{
  "summary_md": "CBAM của EU áp dụng từ đầu 2026, doanh nghiệp có nguy cơ bù đắp **hơn 100 USD/tấn thép** xuất sang EU, nhôm tới cả nghìn USD/tấn. Giá chứng chỉ carbon EU trên 75 USD/tấn CO2. Lô 100.000 tấn thép đơn giá 500-550 USD/tấn phải nộp khoảng 300 tỷ đồng, tương đương 20% giá trị đơn hàng. Cường độ phát thải thép Việt 2,5 tấn CO2/tấn, nhôm 14 tấn so chuẩn EU 1,4 tấn. Sáu nhóm chịu CBAM: sắt thép, nhôm, xi măng, phân bón, điện, hydrogen.",
  "headline_number": "hơn 100 USD/tấn thép",
  "confidence": 0.94,
  "omitted_because_uncertain": []
}
```

*(Ví dụ trên lấy trực tiếp từ ảnh tham chiếu — dùng làm "gold standard" cho cả prompt lẫn eval set.)*

### 5.6 Chuẩn riêng cho Research Note (Màn hình A)

Màn hình A là loại nội dung **khác hẳn** digest, cần spec riêng:

| Thuộc tính | Digest (Màn hình B) | Research Note (Màn hình A) |
|---|---|---|
| Độ dài | 80–120 từ | 600–1.200 từ |
| Cấu trúc | 1 đoạn liền | Tiêu đề + "LUẬN ĐIỂM ĐẦU TƯ" + 3–5 luận điểm đánh số, mỗi luận điểm có dòng dẫn **bold** + 1 đoạn body |
| Nguồn | 1 bài báo | **Nhiều bài + BCTC + báo cáo phân tích** |
| Bold | 1 cụm | Dòng dẫn mỗi luận điểm |
| Rủi ro pháp lý | Trung bình | ⚠️ **CAO** — nội dung này rất gần "tư vấn đầu tư" (§8.5) |
| Quy trình | Auto-publish nếu pass gate | **Bắt buộc có người duyệt** trước khi publish. Không auto-publish trong MVP. |

> **Khuyến nghị PM:** Research Note **không nên** auto-publish trong giai đoạn MVP. Rủi ro pháp lý và rủi ro uy tín quá cao so với lợi ích. Đề xuất: LLM sinh draft → analyst duyệt/sửa → publish, có ghi tên người chịu trách nhiệm nội dung.

---

## 6. Dedupe & clustering

### 6.1 Chuẩn hoá URL (canonicalization) — nền tảng của mọi tầng

| Bước | Quy tắc |
|---|---|
| 1 | Ưu tiên `<link rel="canonical">` nếu có và cùng domain |
| 2 | Lowercase scheme + host. Bỏ `www.` |
| 3 | Bỏ query param tracking: `utm_*`, `fbclid`, `gclid`, `zarsrc`, `source`, `ref`, `campaign`, `cmp`, `spm` |
| 4 | Bỏ fragment `#...` |
| 5 | Bỏ hậu tố AMP: `/amp`, `?amp=1`, `.amp` |
| 6 | Bỏ trailing slash |
| 7 | Với CafeF/VCCorp: giữ ID bài trong slug (`-2026082714300123.chn`) làm khoá phụ vì slug có thể đổi |
| 8 | `url_hash = SHA256(canonical_url)` → **unique index trong DB** |

### 6.2 Bốn tầng dedupe

| Tầng | Tên | Kỹ thuật | Cửa sổ | Xử lý khi trùng |
|---|---|---|---|---|
| **L0** | **Trùng URL** | `url_hash` unique | Không giới hạn | Bỏ ngay ở bước ① Discover. Rẻ nhất, bắt được ~40% trùng lặp. |
| **L1** | **Trùng nội dung nguyên văn** | `content_hash = SHA256(normalized_body)`<br>*(normalize: lowercase, bỏ dấu câu, bỏ khoảng trắng thừa)* | 7 ngày | Bài đăng lại y nguyên (rất phổ biến giữa `cafef.vn` ↔ `markettimes.vn`). Merge, giữ bài sớm nhất. |
| **L2** | **Trùng gần (tiêu đề)** | **SimHash 64-bit** trên tiêu đề đã normalize (lowercase, **bỏ dấu tiếng Việt**, bỏ stopword `của, và, trong, cho, với, về, là, các, những, tại, từ`). Trùng khi **Hamming distance ≤ 3**. | **48 giờ** | Tiêu đề xào lại. Merge vào cùng cluster. |
| **L3** | **Cùng sự kiện (clustering)** | **MinHash + LSH** trên shingle 5-gram của body (Jaccard ≥ 0,5)<br>**HOẶC** embedding cosine ≥ 0,86<br>**RÀNG BUỘC BẮT BUỘC:** phải chia sẻ ≥1 ticker `primary`/`secondary` HOẶC cùng `news_type` + chia sẻ ≥3 con số giống nhau | **72 giờ** | Gom vào 1 `cluster`, chọn 1 bài đại diện. |

### 6.3 Vì sao dùng SimHash cho tiêu đề và MinHash cho body

| | SimHash (L2) | MinHash/LSH (L3) |
|---|---|---|
| Phù hợp với | Văn bản **ngắn** (tiêu đề 10–20 từ) | Văn bản **dài** (body 500–2.000 từ) |
| Chi phí | Rất rẻ, 1 số 64-bit/bài | Tốn hơn, cần signature ~128 hash |
| Nhạy với | Thay đổi vài từ | Thay đổi cấu trúc nội dung |
| Vai trò | Lọc thô, nhanh, chạy trên toàn bộ bài | Lọc tinh, chỉ chạy trên nhóm candidate từ LSH bucket |

**Lưu ý tiếng Việt:** bỏ dấu trước khi hash là bắt buộc — cùng một tiêu đề có thể được gõ "Hoà Phát" hoặc "Hòa Phát", hai chuỗi Unicode khác nhau hoàn toàn.

### 6.4 Chọn bài đại diện cluster (representative selection)

Điểm chọn đại diện, xét theo thứ tự ưu tiên:

```
1. Tier nguồn:           Tier 0 (5đ) > Tier 1 (4đ) > Tier 2 (2đ) > Tier 3 (loại)
2. Không paywall:        +3đ nếu truy cập tự do (quan trọng: người dùng phải mở được link)
3. Mật độ số liệu:       + (số lượng entry trong numbers_index) × 0,15, tối đa +3đ
4. Thời gian đăng:       +2đ cho bài SỚM NHẤT trong cluster (ưu tiên nguồn đưa tin đầu)
5. Độ dài body:          +1đ nếu 800–3.000 ký tự (không quá ngắn, không lan man)
6. Có ảnh đại diện:      +0,5đ
7. Trust score nguồn:    + trust_score × 0,5
```

Bài điểm cao nhất → `is_representative = true` → được summarize và lên digest.
Các bài còn lại → lưu làm `also_reported_by[]`, hiển thị dạng "Cũng được đưa bởi: cafef.vn, markettimes.vn" (tăng độ tin cậy cho người đọc, và là cách attribution tốt).

### 6.5 Rule đặc biệt cho nguồn cùng hệ sinh thái

| Cặp nguồn | Lý do | Rule |
|---|---|---|
| `cafef.vn` ↔ `markettimes.vn` | Cùng VCCorp, đăng chéo gần nguyên văn | Hạ ngưỡng L2 xuống Hamming ≤ 5 (dễ merge hơn). Luôn chọn CafeF làm đại diện. |
| `baodautu.vn` ↔ `tinnhanhchungkhoan.vn` | Cùng Báo Đầu tư | Như trên. Chọn theo chuyên môn: tin TTCK → tinnhanhchungkhoan, tin đầu tư/FDI → baodautu. |
| Bất kỳ ↔ `hsx.vn`/`hnx.vn` | Báo đưa lại CBTT | ⚠️ **Luôn chọn Tier 0 làm canonical**, bài báo thành `also_reported_by`. Nhưng vẫn link tới bài báo vì dễ đọc hơn. |

### 6.6 Xử lý bài cập nhật (live/updated article)

| Tình huống | Rule |
|---|---|
| `updated_at > published_at + 6h` và `content_hash` đổi | Coi là **phiên bản mới** của cùng article. Re-run extract → summarize. Giữ nguyên `short_code` (link không đổi). Đánh dấu `revision = n+1`. |
| Bài live-blog cập nhật liên tục | Không đưa vào digest cho tới khi ngừng cập nhật ≥2 giờ, hoặc chỉ lấy ở run cron kế tiếp. |
| Digest cũ đã publish, bài gốc thay đổi số liệu | Nếu thay đổi làm sai lệch tóm tắt → chạy quy trình đính chính (F9.c). |
| Tin sáng và tin "tổng kết cuối ngày" cùng sự kiện | Ràng buộc L3 sẽ gom lại. Bài tổng kết thường mật độ số cao hơn → tự động thắng ở rule 3 của §6.4. |

### 6.7 Failure mode của dedupe

| Mã | Lỗi | Hệ quả | Giảm thiểu |
|---|---|---|---|
| **D1** | **Over-merge** — gộp 2 sự kiện khác nhau | **Mất tin** (nghiêm trọng hơn) | Ràng buộc bắt buộc ở L3 (phải chia sẻ ticker HOẶC số liệu). Ngưỡng cosine cao (0,86). Log mọi merge để audit. |
| **D2** | **Under-merge** — digest có 2 dòng cùng sự kiện | Giảm chất lượng UI, người dùng thấy lặp | Bổ sung tầng kiểm cuối ở bước ⑨ Publish: so sánh `headline_number` giữa các item; trùng số headline + trùng ≥1 ticker → cảnh báo. |
| **D3** | Cluster quá lớn (>15 bài) | Thường là sự kiện lớn thật, hoặc lỗi over-merge | Cảnh báo review tay khi cluster > 15 |
| **D4** | Bài Tier 2 thắng đại diện vì dài hơn | Digest dẫn nguồn kém uy tín | Tier là tiêu chí ưu tiên số 1, không phải độ dài |
| **D5** | Bài đại diện bị paywall | Người dùng click vào không đọc được | Rule "không paywall +3đ" ở §6.4. Nếu **mọi** bài trong cluster đều paywall → vẫn publish nhưng gắn badge "nội dung có thu phí". |

---

## 7. Short link

### 7.1 Yêu cầu nghiệp vụ (từ ảnh tham chiếu)

Ảnh cho thấy cột `Source` hiển thị `vnexpress.net`, `cafef.vn`, `thanhnien.vn` — **gạch chân, dạng link**. Nghĩa là:

> Người dùng **nhìn thấy domain thật của nguồn**, nhưng khi click thì **đi qua redirect nội bộ của TenPoint** để đo lường.

Đây là mô hình đúng: **minh bạch với người đọc + đo được cho sản phẩm + tôn trọng nguồn**.

### 7.2 Rule nghiệp vụ

| Rule | Nội dung |
|---|---|
| **R7.1 — Hiển thị** | Text hiển thị = **eTLD+1** của `canonical_url`, tính theo **Public Suffix List** (không tự parse bằng regex — `.com.vn`, `.gov.vn`, `.org.vn` sẽ sai). Ví dụ: `https://vnexpress.net/kinh-doanh/...` → hiển thị `vnexpress.net`. |
| **R7.2 — Tooltip / a11y** | `title` và `aria-label` của link phải chứa **tiêu đề bài gốc + tên nguồn đầy đủ**, để người dùng và screen reader biết sẽ đi đâu. Bắt buộc cho a11y — link chỉ có domain là thiếu ngữ cảnh. |
| **R7.3 — Đích thực** | `href` trỏ tới `https://tenpoint.vn/r/{code}` (redirect nội bộ). |
| **R7.4 — Sinh mã** | `code` = **base62** (`0-9A-Za-z`), độ dài **7 ký tự** (~3,5 nghìn tỷ khả năng). Sinh bằng **CSPRNG** hoặc **Feistel/hashids trên counter** để mã **không đoán được và không tuần tự** (chống enumeration → lộ toàn bộ nội dung + số liệu traffic). Check collision trước khi ghi. |
| **R7.5 — Allowlist đích đến** ⚠️ | Redirect **chỉ được** trỏ tới domain trong `source_registry` (đã duyệt). Bất kỳ URL nào ngoài allowlist → **từ chối tạo short link**. Đây là biện pháp chống **open redirect** — lỗ hổng bảo mật cho phép kẻ xấu dùng tenpoint.vn để phishing. |
| **R7.6 — Scheme** | Chỉ chấp nhận `https://` và `http://`. Từ chối `javascript:`, `data:`, `file:`, `//evil.com` (protocol-relative). |
| **R7.7 — Loại redirect** | Dùng **HTTP 302 (Found)**, **không dùng 301**. Lý do: 301 bị trình duyệt cache vĩnh viễn → mất tracking từ lần click thứ 2, và không sửa được đích nếu nguồn đổi URL. |
| **R7.8 — Thuộc tính link** | `rel="nofollow noopener noreferrer"` + `target="_blank"`.<br>• `nofollow` — không truyền link equity, thể hiện đây là trích dẫn tự động, tránh bị coi là link scheme.<br>• `noopener` — chặn tab đích truy cập `window.opener` (bảo mật: tabnabbing).<br>• `noreferrer` — chặn rò rỉ referrer. ⚠️ *Cân nhắc:* nếu muốn nguồn tin **thấy được** traffic TenPoint gửi tới (lý do goodwill để xin hợp tác), nên **bỏ `noreferrer`** và thay bằng UTM (R7.9). Cần PM quyết. |
| **R7.9 — UTM goodwill** | Đề xuất gắn `?utm_source=tenpoint&utm_medium=referral` vào URL đích. Lợi ích: nguồn tin thấy TenPoint **đem traffic đến cho họ**, là lập luận đàm phán hợp tác mạnh (§8.1). Nhưng phải đảm bảo không phá vỡ URL gốc. |
| **R7.10 — Bất biến** | Một `article` = một `short_code` **vĩnh viễn**. Không tái sinh mã. Nếu bài gốc đổi URL, cập nhật đích của cùng `short_code` (đây là lý do dùng 302). |
| **R7.11 — Không cloaking** | **Nghiêm cấm** hiển thị domain A nhưng redirect sang domain B khác. Hiển thị phải luôn là domain thật của đích. |

### 7.3 Tracking

| Sự kiện ghi nhận | Trường |
|---|---|
| `link.click` | `short_code`, `article_id`, `source_id`, `digest_id`, `ts`, `user_id` (nếu đăng nhập), `session_id`, `referrer_internal`, `device_type`, `is_bot` |

**Rule chống nhiễu số liệu:**
- **T1.** Lọc bot theo User-Agent (danh sách crawler phổ biến) + heuristic (không có JS, không có cookie, tần suất bất thường).
- **T2.** De-dup click: cùng `session_id` + cùng `short_code` trong **30 giây** → tính 1 click.
- **T3.** Phân biệt `raw_clicks` và `unique_clicks` (theo session/ngày) trong mọi báo cáo.
- **T4.** Redirect phải **nhanh** (< 100ms) — ghi log **bất đồng bộ**, không chặn redirect. Nếu ghi log lỗi → vẫn redirect, không bao giờ để người dùng gặp lỗi vì tracking.

**Quyền riêng tư (§8.6):**
- Không lưu IP đầy đủ; hash hoặc truncate octet cuối.
- Phải có mục trong Chính sách quyền riêng tư nêu rõ TenPoint đo lường click.
- Tuân thủ Nghị định 13/2023/NĐ-CP về bảo vệ dữ liệu cá nhân (cần luật sư xác nhận phạm vi áp dụng).

### 7.4 Chống lạm dụng

| Nguy cơ | Biện pháp |
|---|---|
| **Open redirect** → phishing qua domain tenpoint.vn | R7.5 allowlist (bắt buộc), R7.6 scheme whitelist |
| **Enumeration** (quét `/r/aaaaaaa` → `/r/zzzzzzz`) | Mã ngẫu nhiên 7 ký tự + rate limit theo IP (ví dụ 60 req/phút) + trả 404 cho mã không tồn tại (không phân biệt "không tồn tại" và "bị vô hiệu") |
| **Click fraud** (bơm số liệu) | T1–T3 + phát hiện bất thường (1 IP > 50 click/giờ) |
| **Short link bị dùng ngoài TenPoint** để spam | Log referrer; nếu lưu lượng lớn đến từ domain lạ → cảnh báo |
| **Đích bị hijack** (nguồn mất domain, domain bị mua lại) | Job kiểm tra định kỳ: nếu domain hết hạn/đổi chủ hoặc trả nội dung bất thường → **vô hiệu hoá** mọi short link tới domain đó |
| **Link tới nội dung bị gỡ** | Job kiểm 404 hằng ngày; mark `source_dead`, hiển thị "bài gốc không còn truy cập được" (F2.e) |

---

## 8. Rủi ro pháp lý & đạo đức

> ⚠️ **Đây là phân tích rủi ro nghiệp vụ, KHÔNG phải tư vấn pháp lý.** Mọi nội dung dưới đây phải được luật sư Việt Nam rà soát trước khi TenPoint ra mắt công khai. Một số văn bản có thể đã được sửa đổi/thay thế tại thời điểm đọc.

### 8.1 🔴 RỦI RO #1 — Giấy phép "trang thông tin điện tử tổng hợp" + thỏa thuận với nguồn tin

**Đây là rủi ro lớn nhất, mang tính sống còn với mô hình kinh doanh.**

Pháp luật Việt Nam phân biệt rõ giữa:
- **Báo điện tử** — cơ quan báo chí, tự sản xuất tin, cần giấy phép hoạt động báo chí;
- **Trang thông tin điện tử tổng hợp** — **tổng hợp, trích dẫn lại tin từ các cơ quan báo chí**, cần **giấy phép riêng**;
- **Mạng xã hội** — nội dung do người dùng tạo, giấy phép khác.

**TenPoint rơi vào nhóm thứ hai.** Theo khung pháp lý về quản lý, cung cấp, sử dụng dịch vụ Internet và thông tin trên mạng (Nghị định 72/2013/NĐ-CP và các nghị định sửa đổi, gần nhất là **Nghị định 147/2024/NĐ-CP**), trang thông tin điện tử tổng hợp thông thường phải:

1. **Có giấy phép** do cơ quan quản lý nhà nước về thông tin cấp;
2. **Có thỏa thuận bằng văn bản** với cơ quan báo chí nguồn về việc sử dụng lại tin, bài;
3. **Trích dẫn nguyên văn, chính xác**, ghi rõ nguồn, tác giả, thời gian đăng, và **dẫn đường link tới bài gốc**;
4. **Không được tự ý biên tập, tổng hợp lại** làm sai lệch nội dung.

**⚠️ Xung đột nghiêm trọng với thiết kế sản phẩm:** điểm (3) và (4) yêu cầu **trích dẫn nguyên văn**, trong khi TenPoint **viết lại tóm tắt bằng LLM**. Đây là mâu thuẫn cần được luật sư giải quyết trước, không phải sau.

| Hành động bắt buộc | Chủ trì | Khi nào |
|---|---|---|
| Tham vấn luật sư chuyên ngành báo chí – truyền thông về mô hình "tóm tắt bằng AI" | PM + Legal | **Trước khi viết dòng code đầu tiên** |
| Xin giấy phép trang TTĐT tổng hợp (nếu áp dụng) | Legal | Trước public launch |
| Ký thỏa thuận nội dung với ≥3 nguồn Tier 1 | PM | Trước public launch |
| Nếu chưa có giấy phép/thỏa thuận → giới hạn phạm vi (nội bộ / beta kín / chỉ dùng Tier 0) | PM | Ngay |

**Chiến lược giảm rủi ro:**
- **Ưu tiên Tier 0** (HOSE/HNX/UBCKNN/NHNN/TCTK). Thông tin CBTT và văn bản nhà nước có rủi ro bản quyền thấp nhất. Một MVP chỉ dùng Tier 0 vẫn có giá trị lớn và **gần như sạch về pháp lý**.
- **Đàm phán sớm với 2–3 nguồn Tier 1**, dùng lập luận: TenPoint gửi traffic về cho họ (R7.9 UTM), không cache full-text, luôn dẫn link.
- **Định vị là "công cụ khám phá tin" (discovery)** thay vì "nơi đọc tin" — mọi thiết kế UI phải khuyến khích click sang bài gốc.

### 8.2 🔴 RỦI RO #2 — Bản quyền: ranh giới giữa "sự kiện" và "cách diễn đạt"

Luật Sở hữu trí tuệ Việt Nam có hai nguyên tắc quan trọng cần khai thác đúng:

| Nguyên tắc | Nội dung (cần luật sư xác nhận) | Áp dụng cho TenPoint |
|---|---|---|
| **Tin thời sự thuần tuý đưa tin không được bảo hộ quyền tác giả** | Bản thân **sự kiện và số liệu** không phải đối tượng bảo hộ. Cái được bảo hộ là **cách diễn đạt** của nhà báo. | ✅ TenPoint tóm tắt **số liệu và sự kiện**, viết lại bằng ngôn từ khác → đây là lập luận phòng vệ chính. **Nhưng** phải đảm bảo LLM **không sao chép nguyên văn câu chữ**. |
| **Trích dẫn hợp lý** | Được trích dẫn tác phẩm đã công bố mà không phải xin phép, không phải trả tiền, nếu: không làm sai ý tác giả, **có dẫn nguồn**, và không nhằm mục đích thương mại chủ yếu dựa trên việc trích dẫn. | ⚠️ TenPoint là sản phẩm thương mại → ngoại lệ "trích dẫn hợp lý" **có thể không áp dụng đầy đủ**. Đây là điểm cần luật sư làm rõ nhất. |

**Rule kỹ thuật bắt buộc để hỗ trợ lập luận trên:**

| Rule | Nội dung | Cách kiểm tra |
|---|---|---|
| **L2.1** | Tóm tắt **không được chứa chuỗi ≥ 15 từ liên tiếp** giống nguyên văn bài gốc (trừ tên riêng, tên tổ chức, thuật ngữ, và cụm số liệu). | Chạy longest-common-substring giữa `summary_md` và `body_text`. Vượt ngưỡng → re-generate. |
| **L2.2** | Tỷ lệ trùng lặp n-gram (5-gram) giữa tóm tắt và bài gốc **< 25%**. | Tính tự động ở bước ⑦ |
| **L2.3** | Tóm tắt phải ngắn hơn **15%** độ dài bài gốc. | 80–120 từ vs bài 600–1.500 từ → thường đạt |
| **L2.4** | **Không sao chép tiêu đề nguyên văn** làm tiêu đề của TenPoint. | Kiểm tra bằng L2.1 |
| **L2.5** | Không sao chép ảnh của nguồn. Nếu cần ảnh, dùng ảnh tự tạo/ảnh có giấy phép. | Rule sản phẩm |

### 8.3 🔴 RỦI RO #3 — Không được khai thác kỹ thuật: robots.txt, ToS, rate limit

| Nguyên tắc | Rule cụ thể | Lý do |
|---|---|---|
| **Tôn trọng robots.txt** | Kiểm tra `robots.txt` trước mỗi crawl (cache 24h). Path bị `Disallow` cho `TenPointBot` hoặc `*` → **không crawl**, dù kỹ thuật vẫn làm được. | Vi phạm robots.txt là bằng chứng bất lợi rõ ràng nếu có tranh chấp |
| **Tôn trọng `Crawl-delay`** | Nếu robots.txt khai báo `Crawl-delay`, tuân thủ giá trị đó (kể cả khi lớn hơn mặc định 1s) | Như trên |
| **UA định danh** | `TenPointBot/1.0 (+https://tenpoint.vn/bot)` — trang `/bot` nêu mục đích, email liên hệ, hướng dẫn chặn. **Nghiêm cấm giả mạo UA trình duyệt.** | Minh bạch = thiện chí. Giả mạo UA = hành vi né tránh có chủ ý |
| **Rate limit lịch sự** | Mặc định **≤ 1 request/giây/domain**, tối đa 2 kết nối đồng thời. Giảm mạnh khi nhận 429/503. Tránh giờ cao điểm của nguồn nếu có thể. | Không gây tải cho nguồn |
| **Không né chặn** | Nếu bị chặn (403/429/captcha): **dừng lại và liên hệ nguồn**. ❌ Không xoay IP, không dùng proxy pool, không giải captcha, không giả mạo header. | Né chặn có thể bị coi là **truy cập trái phép** (rủi ro hình sự theo Bộ luật Hình sự về tội truy cập bất hợp pháp vào mạng máy tính) |
| **Không vượt paywall** | Không dùng tài khoản chia sẻ, không khai thác lỗ hổng soft-paywall, không dùng cache của bên thứ ba để đọc bài có phí. | Vi phạm hợp đồng + rủi ro hình sự |
| **Đọc & lưu ToS** | Mỗi nguồn: đọc Điều khoản sử dụng, lưu bản snapshot + ngày đọc vào `source_registry`. Review lại 6 tháng/lần. | Bằng chứng thiện chí |
| **Cơ chế opt-out** | Trang `/bot` phải có form để chủ website yêu cầu TenPoint ngừng crawl. Xử lý trong ≤ 7 ngày. | Giảm rủi ro leo thang thành tranh chấp |

### 8.4 Không cache full-text công khai

| Rule | Nội dung |
|---|---|
| **C1** | `raw_html` và `body_text` **chỉ dùng nội bộ** cho pipeline. **Không bao giờ** trả về qua API public, không index vào search public, không hiển thị trên UI. |
| **C2** | `raw_html` TTL **30 ngày** → tự động xoá. `body_text` giữ lâu hơn (phục vụ eval/retrain) nhưng cũng phải có chính sách xoá. |
| **C3** | Public chỉ thấy: `summary_md` (do TenPoint viết) + `title` (hiển thị dạng tooltip) + `source_domain` + `link`. |
| **C4** | Không tạo bản "reader mode" / "đọc trong app" của bài gốc. |
| **C5** | Khi nguồn gỡ bài, TenPoint phải gỡ tóm tắt tương ứng nếu được yêu cầu (quy trình takedown, SLA ≤ 48h). |
| **C6** | Có trang `/takedown` công khai để chủ sở hữu nội dung gửi yêu cầu. |

### 8.5 🔴 Disclaimer & rủi ro "tư vấn đầu tư không phép"

**Đây là rủi ro riêng của nội dung tài chính, và là rủi ro pháp lý lớn thứ ba.**

Theo pháp luật chứng khoán Việt Nam, **tư vấn đầu tư chứng khoán là nghiệp vụ kinh doanh có điều kiện**, chỉ công ty chứng khoán được UBCKNN cấp phép mới được thực hiện. Nếu TenPoint phát hành nội dung mà người dùng hiểu là khuyến nghị mua/bán, có thể bị coi là **hoạt động tư vấn đầu tư không phép**.

| Biện pháp | Chi tiết |
|---|---|
| **D1. Disclaimer bắt buộc** | Hiển thị ở **chân mỗi digest**, **chân mỗi research note**, và trong **footer toàn site**. Không được ẩn trong modal hay trang riêng. |
| **D2. Nội dung disclaimer đề xuất** | *"Nội dung trên TenPoint được tổng hợp và tóm tắt tự động từ các nguồn tin công khai, chỉ nhằm mục đích cung cấp thông tin tham khảo. **Đây không phải là khuyến nghị đầu tư, tư vấn mua bán chứng khoán hay lời mời chào giao dịch dưới bất kỳ hình thức nào.** TenPoint không phải công ty chứng khoán và không cung cấp dịch vụ tư vấn đầu tư. Nội dung có thể chứa sai sót do quá trình tự động hoá; vui lòng đối chiếu với bài gốc và các nguồn chính thức trước khi ra quyết định. Nhà đầu tư tự chịu trách nhiệm với mọi quyết định đầu tư của mình."* |
| **D3. Chặn ở tầng nội dung** | Banned-phrase list §5.3 — LLM không được sinh từ khuyến nghị. Đây là **kiểm soát kỹ thuật**, không chỉ là disclaimer. |
| **D4. Tường thuật, không đồng tình** | Khi bài gốc có khuyến nghị của CTCK, TenPoint chỉ **tường thuật kèm tên tổ chức** ("SSI Research đưa giá mục tiêu X"), không bao giờ trình bày như quan điểm của TenPoint. |
| **D5. Nhãn "nội dung tạo bởi AI"** | Ghi rõ tóm tắt được tạo tự động. Minh bạch giảm rủi ro và tăng niềm tin. |
| **D6. Research Note phải có người chịu trách nhiệm** | Màn hình A gần với phân tích đầu tư nhất → bắt buộc có analyst duyệt + ghi tên (§5.6). |
| **D7. Không cá nhân hoá lời khuyên** | Tuyệt đối không có tính năng kiểu "dựa trên danh mục của bạn, bạn nên…". Đó là ranh giới rõ ràng sang tư vấn đầu tư. |

### 8.6 Dữ liệu cá nhân

| Rủi ro | Biện pháp |
|---|---|
| Tracking click thu thập dữ liệu người dùng | Chính sách quyền riêng tư công khai; hash/truncate IP; không bán dữ liệu |
| Tài khoản người dùng (nếu có) | Tuân thủ Nghị định 13/2023/NĐ-CP về bảo vệ dữ liệu cá nhân — cần đánh giá tác động xử lý DLCN (cần luật sư xác nhận phạm vi) |
| Tên cá nhân trong tin tức (lãnh đạo doanh nghiệp, bị can) | Chỉ giữ thông tin đã công khai trên báo chí chính thống; có quy trình gỡ khi có yêu cầu hợp pháp |

### 8.7 Đạo đức nội dung

| Nguyên tắc | Rule |
|---|---|
| **Không thao túng thị trường** | Không ưu tiên/bỏ qua tin về mã cụ thể vì lợi ích. Nếu TenPoint hoặc nhân sự nắm giữ cổ phiếu → phải có chính sách công bố xung đột lợi ích. |
| **Không giật tít** | Banned clickbait §5.3. Tiêu đề digest phải mô tả đúng nội dung. |
| **Không khuếch đại tin đồn** | Chỉ lấy từ nguồn trong `source_registry`. Không lấy từ diễn đàn, group Zalo/Telegram, room hô hào. |
| **Minh bạch về sai sót** | Có trang "Đính chính" công khai. Mọi lần sửa nội dung đã publish đều được ghi log và hiển thị. |
| **Không ẩn nguồn** | Không bao giờ để tóm tắt xuất hiện mà thiếu link nguồn (R9.7). |
| **Đối xử công bằng với nguồn** | Không thiên vị một nguồn vì lý do thương mại mà không công bố. |

---

## 9. Data Quality KPI

### 9.1 Bảng KPI chính

| # | KPI | Định nghĩa | Cách đo | Mục tiêu MVP | Mục tiêu 6 tháng | Ngưỡng cảnh báo |
|---|---|---|---|---|---|---|
| **K1** | **Source coverage** | % nguồn active trả về ≥1 bài mới trong 24h | `nguồn có bài / tổng nguồn active` | ≥ 90% | ≥ 95% | < 80% → P2 |
| **K2** | **Event coverage** | % sự kiện quan trọng (đối chiếu checklist tay 30 sự kiện/tuần) mà TenPoint có trong digest | Audit tay hằng tuần | ≥ 75% | ≥ 90% | < 60% → P2 |
| **K3** | **Ingestion success rate** | % candidate_url đi được tới trạng thái `extracted` | `extracted / discovered` | ≥ 92% | ≥ 97% | < 85% → P1 |
| **K4** | **Độ trễ P50** | Trung vị (published_at → publish_at trên TenPoint) | Tính trên mọi digest item | ≤ 4 giờ | ≤ 3 giờ | > 6h → P2 |
| **K5** | **Độ trễ P95** | Phân vị 95 của cùng chỉ số | | ≤ 9 giờ | ≤ 8,5 giờ | > 12h → P2 |
| **K6** | **Độ trễ fast-lane (CBTT)** | published_at → publish_at cho tin Tier 0 | | ≤ 45 phút | ≤ 20 phút | > 2h → P1 |
| **K7** | 🔴 **Tỷ lệ tóm tắt sai số liệu** | % tóm tắt đã publish chứa ≥1 số không khớp bài gốc | Audit tay **40 item/tuần** + numeric grounding log | **≤ 1%** | **≤ 0,3%** | **> 2% → P0, dừng auto-publish** |
| **K8** | **Tỷ lệ bịa số (hallucinated)** | % tóm tắt chứa số **không tồn tại** trong bài gốc | Numeric grounding check (§5.4) | **0%** ở output publish | **0%** | **> 0 → P0** |
| **K9** | 🔴 **Ticker tagging precision** | Trong các mã đã gắn, % gắn đúng | Audit tay 100 bài/tuần | **≥ 96%** | **≥ 99%** | **< 93% → P0** |
| **K10** | **Ticker tagging recall** | Trong các mã đáng lẽ phải gắn, % đã gắn | Audit tay cùng bộ | ≥ 85% | ≥ 92% | < 75% → P2 |
| **K11** | **Ticker F1** | Trung bình điều hoà K9, K10 | | ≥ 0,90 | ≥ 0,95 | — |
| **K12** | **Tỷ lệ trùng lặp trong digest** | % digest có ≥2 item cùng sự kiện | Audit tay mỗi digest | ≤ 3% | ≤ 1% | > 8% → P2 |
| **K13** | **Dedupe over-merge rate** | % cluster gộp nhầm 2 sự kiện khác nhau | Audit tay 50 cluster/tuần | ≤ 2% | ≤ 0,5% | > 5% → P1 |
| **K14** | **Classification accuracy** | % `news_type` gán đúng theo taxonomy §3 | Eval set 200 bài gán tay, chạy lại hằng tuần | ≥ 85% | ≥ 92% | < 78% → P2 |
| **K15** | **Summary rubric pass rate** | % tóm tắt đạt ≥42/52 điểm (§5.2) | LLM-judge toàn bộ + audit tay mẫu | ≥ 88% | ≥ 95% | < 80% → P1 |
| **K16** | **Broken link rate** | % short link trả về 4xx/5xx ở bài gốc | Job kiểm hằng ngày | ≤ 2% | ≤ 0,5% | > 5% → P2 |
| **K17** | **Human review queue size** | Số item chờ duyệt cuối ngày | | ≤ 15% tổng item | ≤ 8% | > 30% → P2 |
| **K18** | **Publish SLA** | % digest publish đúng khung giờ (±15 phút) | | ≥ 95% | ≥ 99% | < 90% → P1 |
| **K19** | **Robots/ToS compliance** | Số lần crawl vi phạm robots.txt | Log tự động | **0** | **0** | **> 0 → P0 pháp lý** |
| **K20** | **Attribution completeness** | % item publish có link nguồn hoạt động | | **100%** | **100%** | **< 100% → P0 pháp lý** |

### 9.2 Vì sao precision ưu tiên hơn recall (K9 vs K10)

| | Gắn sai mã (false positive) | Bỏ sót mã (false negative) |
|---|---|---|
| Người dùng nhìn thấy? | ✅ Có — ngay lập tức | ❌ Không — họ không biết cái họ không thấy |
| Hệ quả | Mất niềm tin vào **toàn bộ** sản phẩm. "Nó gắn nhầm VIC vào tin không liên quan → tôi không tin cái gì khác nữa." | Giảm độ phủ, có thể tìm ở nơi khác |
| Khả năng phục hồi | Khó | Dễ |

→ Vì vậy K9 (precision) đặt ngưỡng **96% cho MVP**, cao hơn hẳn K10 (recall, 85%). Mọi rule ở §4 đều thiết kế theo hướng **bỏ sót còn hơn gắn sai**.

### 9.3 Quy trình audit

| Tần suất | Việc | Người làm | Sản phẩm |
|---|---|---|---|
| **Mỗi run (8h)** | Auto-check: numeric grounding, banned phrase, độ dài, bold count, allowlist link | Hệ thống | Gate pass/fail |
| **Hằng ngày** | Kiểm broken link, health nguồn, kích thước review queue | Hệ thống + on-call | Dashboard |
| **Hằng tuần** | Audit tay: 100 bài cho ticker (K9/K10), 40 tóm tắt (K7), 50 cluster (K13) | BA/Analyst | Báo cáo chất lượng tuần |
| **Hằng tuần** | Chạy eval set 200 bài cho classification (K14) | Hệ thống | Điểm accuracy |
| **Hằng tháng** | Đối chiếu event coverage (K2) với danh sách sự kiện quan trọng do analyst lập | BA | Báo cáo độ phủ |
| **6 tháng/lần** | Review ToS/robots.txt toàn bộ nguồn; đồng bộ rổ VN30; review taxonomy | BA + Legal | Cập nhật `source_registry`, `ticker_master` |

### 9.4 Bộ dữ liệu vàng (golden set) cần xây trước khi code

| Bộ | Kích thước | Dùng cho | Ai gán |
|---|---|---|---|
| **GS-Ticker** | 300 bài, gán tay toàn bộ mã + mức relevance | K9, K10, K11 | 2 analyst gán độc lập, giải quyết bất đồng |
| **GS-Type** | 200 bài, gán tay `news_type` | K14 | Như trên |
| **GS-Summary** | 100 bài + tóm tắt chuẩn viết tay (dùng 5 dòng trong ảnh tham chiếu làm hạt nhân) | K15, tinh chỉnh prompt | Analyst + BA |
| **GS-Dedupe** | 50 cluster đã gán tay | K13 | Analyst |

> **Khuyến nghị PM:** GS-Summary là tài sản quan trọng nhất. 5 dòng tóm tắt trong ảnh tham chiếu đã là "gold standard" chất lượng rất cao — hãy dùng chúng làm chuẩn và mở rộng lên 100 mẫu **trước khi** bắt đầu tinh chỉnh prompt.

---

## Phụ lục A — Các thực thể dữ liệu tối thiểu

| Thực thể | Trường then chốt |
|---|---|
| `source` | `id, name, domain, tier, trust_score, feed_urls[], robots_checked_at, tos_snapshot, tos_reviewed_at, agreement_status, selectors, rate_limit, health, active` |
| `article` | `id, source_id, canonical_url, url_hash, content_hash, title, sapo, body_text, published_at, updated_at, fetched_at, paywalled, numbers_index[], state, revision` |
| `cluster` | `id, representative_article_id, member_article_ids[], event_window_start, event_window_end, merge_method, merge_score` |
| `ticker_master` | xem §4.7 |
| `article_ticker` | `article_id, ticker_code, relevance, confidence, evidence_offsets[], detector` |
| `summary` | `id, article_id, summary_md, headline_number, numbers_used[], rubric_score, grounded, model, prompt_version, generated_at, reviewed_by` |
| `short_link` | `short_code, article_id, target_url, display_domain, created_at, disabled, disabled_reason` |
| `click_event` | `short_code, ts, session_id, user_id, digest_id, device_type, is_bot, ip_hash` |
| `digest` | `id, date, slot(06/14/22), title, item_ids[], published_at, status` |
| `correction` | `id, entity_type, entity_id, reason, before, after, corrected_by, corrected_at, public_note` |

---

## Phụ lục B — Câu hỏi cần PM quyết định

| # | Câu hỏi | Vì sao quan trọng | Đề xuất của BA2 |
|---|---|---|---|
| Q1 | **Có xin giấy phép trang TTĐT tổng hợp + thỏa thuận nguồn không?** | Quyết định toàn bộ mô hình có hợp pháp hay không (§8.1) | **Bắt buộc**. Tham vấn luật sư ngay tuần này. |
| Q2 | MVP có nên **chỉ dùng Tier 0** (HOSE/HNX/UBCKNN/NHNN) không? | Rủi ro pháp lý gần như bằng 0, vẫn có giá trị sản phẩm | **Nên** — làm MVP sạch pháp lý trước, mở Tier 1 khi đã có thoả thuận |
| Q3 | Fast-lane CBTT 30 phút/lần có trong MVP không? | Cron 8h là quá chậm cho tin vật chất | Nên có — đây là điểm khác biệt lớn nhất so với đọc báo thường |
| Q4 | Research Note (Màn hình A) có auto-publish không? | Rủi ro pháp lý cao nhất (§5.6, §8.5) | **Không** — bắt buộc analyst duyệt trong MVP |
| Q5 | Giữ `noreferrer` hay bỏ để nguồn thấy traffic từ TenPoint? | Ảnh hưởng đàm phán hợp tác (R7.8, R7.9) | Bỏ `noreferrer`, thêm UTM — traffic là lá bài đàm phán |
| Q6 | Ngân sách cho human review queue? | K17 dự kiến 8–15% item cần người duyệt | Cần ít nhất 1 analyst part-time từ ngày đầu |
| Q7 | Fireant/Simplize: đối tác hay đối thủ? | Ảnh hưởng chiến lược dữ liệu (§1.5) | Không scrape. Tự build ticker_master từ HOSE/HNX. |
| Q8 | Có build golden set trước khi code không? | Không có golden set thì không đo được K7–K15 | **Có** — 2 tuần công sức analyst, tiết kiệm nhiều tháng sửa sai |

---

*Hết BA2. Tài liệu cần đọc cùng: BA1 (sản phẩm/UX), BA3 (nếu có).*
