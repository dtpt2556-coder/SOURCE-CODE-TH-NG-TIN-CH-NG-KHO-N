# TenPoint — PRD v1.0 (PM)

> Tài liệu quyết định của PM, tổng hợp từ 3 báo cáo BA. Đây là **nguồn chốt** khi có tranh chấp về phạm vi.
> Đầu vào: `01-ba/BA1-users-and-usecases.md`, `01-ba/BA2-data-sources-pipeline.md`, `01-ba/BA3-market-differentiation-business.md`, `00-inputs/reference-images.md`.

---

## 1. Sản phẩm là gì

**TenPoint biến hàng trăm bài báo chứng khoán mỗi ngày thành ~10 tin đã chưng cất — giàu số liệu, gắn sẵn mã cổ phiếu, luôn dẫn nguồn gốc bằng một cú chạm — để nhà đầu tư nắm hết thông tin quan trọng trong 3 phút buổi sáng thay vì 45 phút lướt tin.**

Hai lớp nội dung, đúng 2 ảnh tham chiếu:

| Lớp | Màn hình | Vai trò |
|---|---|---|
| **Digest** | Ảnh B — bảng 5 cột | Nhịp hằng ngày. Lý do người dùng quay lại mỗi sáng |
| **Research Note** | Ảnh A — luận điểm đầu tư | Chiều sâu. Lý do người dùng tin tưởng sản phẩm |

---

## 2. Vì sao sản phẩm này có chỗ đứng

BA3 xác định 3 khoảng trống, tôi chốt cả 3 là nền tảng định vị:

1. **Không ai làm lớp chưng cất có kỷ luật số liệu.** Thị trường chỉ có 2 thái cực: bài báo 1.000 từ giật tít, hoặc PDF research 15 trang. CafeF *không thể* làm vì sống bằng pageview — đây là innovator's dilemma và là lợi thế bền của TenPoint.
2. **Không ai coi tin tức là dữ liệu có khoá `ticker`.** Trong ảnh, cột Mã CK đứng **thứ hai, trước cả nội dung**. Đó là tuyên ngôn sản phẩm. Hệ quả kỹ thuật: timeline theo mã, alert watchlist, và hệ thống internal link tự sinh cho SEO.
3. **Không ai vừa nhanh vừa truy vết được nguồn.** Room Zalo/Telegram nhanh nhưng xung đột lợi ích; CTCK uy tín nhưng là PDF không tra cứu được. Ô "nhanh + cô đọng + kiểm chứng được" đang trống.

**North Star Metric (chốt theo BA3):** **WQMR** — số người có ≥3 "buổi sáng chất lượng"/tuần (mở digest 6h–9h, xem ≥3 tin, ≥1 tin trúng watchlist).
Không chọn pageview (chính là thứ ép CafeF giật tít). Không chọn time-on-site (với TenPoint, **đọc xong nhanh là thành công**).

---

## 3. Quyết định của PM — 8 câu hỏi mở của BA1

| # | Vấn đề | **Quyết định** | Lý do |
|---|---|---|---|
| Q1 | Tin trùng lặp nhiều báo | **Gộp cụm, hiển thị 1 nguồn tier cao nhất làm đại diện**; lưu `cluster_id` cho toàn cụm | Giữ bảng sạch đúng tinh thần ảnh. Cụm đã lưu → v1.1 nâng thành "N nguồn cùng đưa tin" (W3 của BA3) mà không cần migrate |
| Q2 | Số liệu mâu thuẫn giữa các báo | **Tin nguồn tier cao hơn, ghi log mâu thuẫn** | Cần nhất quán; log để QA rà |
| Q3 | Mã nhắc thoáng qua | **Lưu 2 mức `primary` / `mentioned`; mặc định lọc `primary`** | Tái tạo đúng hành vi trong ảnh: hiện TCB/BCM/VIC/FPT, ẩn PNJ/CTG/ACB/VHM |
| Q4 | Có cần đăng nhập? | **KHÔNG ở v1.** Watchlist lưu `localStorage` | Bỏ toàn bộ chi phí auth + bảo vệ dữ liệu người dùng khỏi MVP |
| Q5 | LLM bịa số liệu | **Validator bắt buộc: mọi số trong tóm tắt phải truy vết được về bài gốc. Fail → `pending_review`, KHÔNG publish** | BA3 xếp đây vào rủi ro **mức tồn vong**. Nội dung tài chính sai số gây thiệt hại thật |
| Q6 | Short link chết | **Health-check hàng tuần, gắn cờ `target_alive=false`** → đưa vào v1.1 | Không chặn MVP nhưng ảnh hưởng niềm tin |
| Q7 | Nguồn đổi HTML → parser gãy | **Cảnh báo khi 1 nguồn trả 0 bài trong 2 run liên tiếp** | Đã có sẵn `crawl_run_sources` để phát hiện |
| Q8 | Bản quyền: lưu bao nhiêu? | **Chỉ hiển thị công khai tóm tắt tự sinh + link nguồn. KHÔNG public full-text.** `excerpt` chỉ dùng nội bộ cho validator | Xem mục 4 |

---

## 4. Quyết định về nguồn dữ liệu — điểm rủi ro lớn nhất của dự án

BA2 nêu một vấn đề nghiêm trọng mà ý tưởng gốc chưa tính đến:

> Tại Việt Nam, tổng hợp tin từ báo chí có thể cần **giấy phép "trang thông tin điện tử tổng hợp"** (NĐ 72/2013, NĐ 147/2024) và **thoả thuận bằng văn bản với từng nguồn**. Quy định yêu cầu **trích dẫn nguyên văn** — điều này mâu thuẫn trực tiếp với việc tóm tắt bằng LLM.

Đây là rủi ro pháp lý thật, không phải rủi ro kỹ thuật, và **không thể giải quyết bằng code**.

### Quyết định của PM

Xây hệ thống **trung lập với nguồn**, điều khiển bằng cấu hình chứ không phải bằng code:

| Môi trường | `sources.enabled` | Mục đích |
|---|---|---|
| **Dev / Demo / QA** | Bật toàn bộ Tier 0 + Tier 1 + Tier 2 | Chạy thử đầy đủ, kiểm chứng pipeline |
| **Production trước khi có thoả thuận** | **Chỉ bật Tier 0** (HOSE, HNX, UBCKNN, NHNN, TCTK) | Nguồn công bố chính thức — sạch pháp lý nhất, và vẫn có giá trị thật với người dùng |
| **Production sau khi ký thoả thuận** | Bật thêm từng nguồn Tier 1 đã ký | Mở dần theo tiến độ pháp lý |

**Cổng chặn phát hành (release gate) — bắt buộc trước khi mở công khai:**
1. Tham vấn luật sư về giấy phép trang TTĐT tổng hợp.
2. Rule kỹ thuật chống vi phạm bản quyền đã bật: không copy ≥15 từ liên tiếp từ bài gốc, trùng 5-gram < 25%, không sao chép tiêu đề gốc nguyên văn, không sao chép ảnh, không cache full-text public.
3. Banned-phrase list chặn ngôn ngữ tư vấn đầu tư **ở tầng sinh nội dung**, không chỉ dựa vào disclaimer — vì tư vấn đầu tư là nghiệp vụ có điều kiện, chỉ CTCK được cấp phép.

> **Ghi chú cho chủ sản phẩm:** phần code của MVP không bị chặn bởi việc này — hệ thống chạy đầy đủ với mọi nguồn trong môi trường dev. Nhưng **trước khi mở cho người dùng thật tại Việt Nam**, ba cổng chặn trên phải được thông qua. Đây là việc của pháp chế, không phải của kỹ thuật.

---

## 5. Taxonomy `Loại tin` — và một phát hiện về ảnh tham chiếu

**Chốt 9 nhãn (theo BA2):**
`Vĩ mô` · `Ngành` · `Doanh nghiệp` · `Thị trường` · `Khối ngoại` · `Cổ tức/Phát hành` · `Pháp lý` · `Phân tích` · `Trái phiếu/Tín dụng`

Thứ tự ưu tiên khi 1 tin thuộc nhiều nhóm:
`Pháp lý > Cổ tức/PH > Khối ngoại > Tín dụng > Phân tích > Doanh nghiệp > Ngành > Thị trường > Vĩ mô`

### Phát hiện của BA2 cần ghi lại

Trong ảnh tham chiếu, **cả 5/5 dòng đều được gán `Ngành`, nhưng chỉ 2/5 thực sự đúng**:

| Tin | Nhãn trong ảnh | Nhãn đúng |
|---|---|---|
| 12 ngân hàng cam kết 408.000 tỷ tín dụng DNNVV | Ngành | **Trái phiếu/Tín dụng** |
| CBAM thép/nhôm xuất EU | Ngành | Ngành ✅ |
| VN-Index đóng cửa 1.821 điểm | Ngành | **Thị trường** |
| Nhập khẩu thịt 6 tháng | Ngành | Ngành ✅ |
| FTSE GEIS hút 5.588 tỷ | Ngành | **Khối ngoại** |

Nghĩa là ở bản mẫu, **cột `Loại tin` đang vô dụng cho việc lọc**. Đây chính là phần TenPoint phải làm tốt hơn bản mẫu.

**Quyết định:** dữ liệu seed demo giữ nguyên nhãn như ảnh để đối chiếu thị giác 1:1; nhưng **classifier của hệ thống phải gán đúng theo taxonomy 9 nhãn**. QA kiểm hai việc này tách bạch.

---

## 6. Phạm vi MVP v1 — chốt

**MUST (P0) — không có thì không phải sản phẩm**

| # | Hạng mục | Nguồn yêu cầu |
|---|---|---|
| 1 | Bảng digest 5 cột đúng ảnh B | Ảnh tham chiếu |
| 2 | Tóm tắt 60–150 từ, giàu số liệu, bold 1–3 cụm số | BA1 NUM-1..NUM-8 |
| 3 | **Validator số liệu chặn publish khi bịa số** | BA1 Q5 · BA3 rủi ro tồn vong |
| 4 | Gắn mã CK `primary`/`mentioned`, có stoplist + rule VND + chống nhiễu nguồn trích dẫn | BA2 |
| 5 | Phân loại 9 nhãn Loại tin | BA2 |
| 6 | Short link `/r/{code}` + đếm click + **allowlist chống open redirect** | Yêu cầu gốc · BA1 US-3.3 |
| 7 | Cron 8h/lần tại **06:00 / 14:00 / 22:00 giờ VN** | Yêu cầu gốc · BA1 mục 8 |
| 8 | Lọc theo mã CK và Loại tin, state trong query string | BA1 US-2.1/2.2 |
| 9 | Trang `/ma/{mã}` — timeline theo mã | BA1 US-4.1 · BA3 W1 (nền móng SEO) |
| 10 | Trang Research Note đúng ảnh A + disclaimer | Ảnh tham chiếu · BA1 US-4.2 |
| 11 | Dedupe cụm tin trùng | BA1 Q1 |
| 12 | Responsive mobile: bảng → card, không scroll ngang | BA1 NFR |
| 13 | `Cập nhật lần cuối` + cảnh báo dữ liệu cũ | BA1 US-1.3 |
| 14 | `/healthz`, `/readyz`, `crawl_runs` | BA1 US-6.1/6.2 |
| 15 | Ingest thủ công có token | BA1 US-6.3 |

**SHOULD (P1)** — Tìm kiếm không dấu · Lọc khoảng ngày · Watchlist `localStorage` · Copy bản tin cho môi giới · Phân trang "Xem thêm" · SEO metadata + sitemap + JSON-LD

**COULD (P2)** — Dark mode · RSS output · Export CSV · Sentiment · Lọc theo nguồn báo

**WON'T (v1)** — Đăng nhập · Realtime · Dữ liệu giá & biểu đồ · App native · Bình luận · Thanh toán · Zalo/Telegram bot · API B2B

---

## 7. Điều chỉnh so với đề xuất của BA3

BA3 đề nghị đưa 5 tính năng "wow" vào MVP. Tôi chốt lại:

| BA3 đề xuất | Quyết định PM | Lý do |
|---|---|---|
| W1 Ticker Timeline | ✅ **Vào MVP** | Chi phí gần bằng 0 vì dữ liệu đã có sẵn; đồng thời là toàn bộ nền móng SEO |
| W8 Bản tin 7h00 + 15h30 | ⏸ **Lùi v1.1** | Cần hạ tầng email/list + xử lý opt-out. Không chặn giá trị cốt lõi |
| W2 Alert watchlist qua Zalo/Telegram | ⏸ **Lùi v2** | Cần tài khoản OA, kiểm duyệt, chi phí vận hành |
| W7 Diff Digest ("từ lần đọc gần nhất") | ⏸ **Lùi v1.1** | Rẻ và chống churn tốt, nhưng cần `localStorage` timestamp — làm sau khi có người dùng thật |
| W3 "N nguồn cùng đưa tin" | ⏸ **Lùi v1.1** | `cluster_id` đã lưu sẵn ở v1 → v1.1 chỉ là việc hiển thị |

**Nguyên tắc:** MVP chỉ giữ thứ chứng minh được giả thuyết cốt lõi — *người dùng có quay lại mỗi sáng vì bảng digest chưng cất hay không*. Mọi kênh phân phối (email, Zalo) là khuếch đại, chỉ đáng làm khi lõi đã được chứng minh.

---

## 8. Mô hình kinh doanh — chốt cho giai đoạn đầu

Theo BA3, ưu tiên 2 hướng, **không** làm freemium ở v1:

1. **White-label cho môi giới** ⭐ (500k–1,5tr/môi giới/tháng; 10–30tr/phòng giao dịch) — dòng tiền nhanh nhất, biến đối thủ phân phối nguy hiểm nhất (broker tự soạn bản tin) thành khách hàng. Mỗi môi giới mang theo 200–2.000 người đọc cuối → kênh B2B2C. **Pilot 3–5 môi giới từ tuần 3.**
2. **Affiliate mở tài khoản CK** (~100–400k/tài khoản active) — không tạo ma sát cho người đọc.

**Hoãn:** freemium (v2, cần ≥10k WAU) · API B2B (v3 — ARPU cao nhưng rủi ro bản quyền nhảy vọt, phải qua cổng pháp lý mục 4 trước) · **IR sponsored content: không làm** — phá vỡ định vị trung lập, là tài sản duy nhất của sản phẩm.

---

## 9. Roadmap

| Giai đoạn | Thời lượng | Mục tiêu | Tiêu chí thoát |
|---|---|---|---|
| **MVP v1** | 4–6 tuần | Chứng minh người dùng quay lại mỗi sáng | Digest chạy ổn định 3 run/ngày · tỷ lệ tóm tắt sai số = 0 · precision tag mã ≥ 96% · WQMR có tín hiệu dương |
| **v1.1** | +3 tuần | Giữ chân & phân phối | Bản tin email 7h · Diff Digest · "N nguồn cùng đưa tin" · health-check short link |
| **v2** | +6 tuần | Mở rộng người dùng | Alert Zalo/Telegram · watchlist đồng bộ (cần login) · freemium · dark mode |
| **v3** | — | Doanh thu B2B | API cho CTCK/quỹ — **chỉ sau khi qua cổng pháp lý mục 4** |

---

## 10. Chỉ tiêu chất lượng — QA nghiệm thu theo bảng này

| Nhóm | Chỉ tiêu | Ngưỡng |
|---|---|---|
| **Độ chính xác số liệu** | Tin publish có số bịa | **0** — tuyệt đối |
| **Tag mã CK** | Precision | **≥ 96%** (ưu tiên precision hơn recall — gắn sai mã tệ hơn bỏ sót) |
| **Tag mã CK** | Recall | ≥ 85% |
| **Tóm tắt** | Độ dài | 60–150 từ; reject nếu > 180 |
| **Tóm tắt** | Cụm bold | đúng 1–3 |
| **Tóm tắt** | Giọng văn | 0 cụm từ khuyến nghị đầu tư |
| **Pipeline** | 1 run | < 10 phút |
| **Pipeline** | 1 nguồn lỗi | không làm hỏng run |
| **API** | p95 `/api/v1/news` | < 200ms |
| **Short link** | Thời gian redirect | < 50ms |
| **Short link** | Open redirect | **0** — allowlist chặn 100% domain lạ |
| **FE** | LCP mobile 4G | < 2.0s |
| **FE** | Bảng trên mobile | **không** scroll ngang |
| **A11y** | WCAG 2.2 AA | đạt checklist `04-design/design-system.md` mục 8 |
| **Giao diện** | Khớp ảnh tham chiếu | Đối chiếu 1:1 bằng dữ liệu seed demo |

---

## 11. Rủi ro dự án

| Rủi ro | Mức | Giảm thiểu |
|---|---|---|
| **Pháp lý về bản quyền / giấy phép** | 🔴 Tồn vong | Cổng chặn mục 4. Mặc định production chỉ Tier 0 |
| **Ảo giác số liệu của LLM** | 🔴 Tồn vong | Validator bắt buộc (Q5). Provider mặc định là `extractive` không dùng LLM |
| **Bị quy là tư vấn đầu tư không phép** | 🔴 Cao | Banned-phrase ở tầng sinh nội dung + disclaimer + Research Note phải có người duyệt |
| Nguồn đổi HTML → parser gãy | 🟠 TB | Cảnh báo 0 bài trong 2 run (Q7) |
| Chi phí LLM tăng theo lượng tin | 🟠 TB | Fallback `extractive` mặc định; chỉ bật LLM cho tin có mã `primary` |
| Nguồn chặn IP do crawl | 🟠 TB | robots.txt, rate limit theo nguồn, User-Agent định danh |
| Chính các báo tự làm tính năng này | 🟡 Thấp | Innovator's dilemma — mô hình pageview của họ xung đột với việc chưng cất |

---

## 12. Phân công & trạng thái

| Vai trò | Sản phẩm bàn giao | Trạng thái |
|---|---|---|
| BA1 | `01-ba/BA1-users-and-usecases.md` | ✅ |
| BA2 | `01-ba/BA2-data-sources-pipeline.md` | ✅ |
| BA3 | `01-ba/BA3-market-differentiation-business.md` | ✅ |
| PM | `02-pm/PRD.md` (tài liệu này) | ✅ |
| SA | `03-sa/architecture.md` | ✅ |
| UI/UX | `04-design/design-system.md` | ✅ |
| BE Golang | `backend/` | 🔄 đang triển khai |
| FE | `frontend/` | 🔄 đang triển khai |
| DevOps | `deploy/` + Dockerfile | ⏳ chờ BE/FE |
| QA | `05-qa/` | ⏳ chờ BE/FE |
