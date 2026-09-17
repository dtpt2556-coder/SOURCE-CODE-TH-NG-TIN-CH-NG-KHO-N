# BA1 — Người dùng & Use case

> Phân tích nghiệp vụ từ góc nhìn người dùng cho sản phẩm **TenPoint** — nền tảng tổng hợp & tóm tắt tin tức chứng khoán Việt Nam.
> Nguồn đầu vào: `docs/00-inputs/reference-images.md` (2 màn hình tham chiếu) + ý tưởng gốc của khách hàng.

---

## 1. Personas

### P1 — Minh, 34 tuổi, "Nhà đầu tư cá nhân bận rộn" (persona chính, 70% traffic mục tiêu)

| Thuộc tính | Nội dung |
|---|---|
| Bối cảnh | Nhân viên văn phòng tại TP.HCM, danh mục 300–800 triệu, nắm 5–12 mã. Giao dịch qua app CTCK (VPS/SSI/DNSE). Không có thời gian trong giờ hành chính. |
| Thói quen | Đọc tin 6:30–8:00 sáng trên điện thoại; liếc nhanh 11:35 giờ nghỉ trưa; 15:30 sau phiên ATC. |
| Pain point | (1) Phải mở 4–5 tab: CafeF, Vietstock, VnExpress, nhóm Zalo. (2) Tin dài 1.200 chữ nhưng ý chính chỉ 3 dòng. (3) Không biết tin này **liên quan mã nào mình đang cầm**. (4) Tin bị lặp — 5 báo viết cùng 1 sự kiện. (5) Số liệu nằm rải rác trong bài, phải tự nhặt. |
| Jobs-to-be-done | *"Khi tôi có 5 phút buổi sáng, tôi muốn biết chính xác có gì mới ảnh hưởng tới các mã tôi đang nắm, kèm con số cụ thể, để tôi quyết định có cần hành động trong phiên hôm nay không."* |
| Thành công = | Đọc xong trong < 4 phút, biết 100% tin quan trọng của danh mục, không phải mở app khác. |

### P2 — Hương, 28 tuổi, "Môi giới / Chuyên viên tư vấn đầu tư"

| Thuộc tính | Nội dung |
|---|---|
| Bối cảnh | Broker tại CTCK, quản 60–150 tài khoản khách. Mỗi sáng phải gửi bản tin cho nhóm khách hàng. |
| Pain point | (1) Mất 60–90 phút/ngày tự tổng hợp bản tin sáng. (2) Phải tự viết lại tóm tắt cho ngắn gọn. (3) Khách hỏi "tin này ở đâu?" → cần link nguồn ngay để giữ uy tín. (4) Cần lọc tin theo nhóm ngành để tư vấn từng khách. |
| Jobs-to-be-done | *"Khi chuẩn bị bản tin sáng, tôi muốn copy được một bảng tin đã tóm tắt sẵn kèm link nguồn, để gửi khách trong 5 phút thay vì 1 tiếng."* |
| Thành công = | Thời gian soạn bản tin giảm từ 60' → 10'; có nút copy/export. |

### P3 — Tuấn, 23 tuổi, "Nhà đầu tư mới (F0)"

| Thuộc tính | Nội dung |
|---|---|
| Bối cảnh | Mới mở tài khoản 3–6 tháng, vốn 20–80 triệu, học qua YouTube/TikTok. |
| Pain point | (1) Đọc tin tài chính không hiểu thuật ngữ (LNST CĐ mẹ, biên gộp, P/E TTM, khối ngoại). (2) Không phân biệt được tin quan trọng vs tin nhiễu. (3) Dễ bị dẫn dắt bởi tin đồn nhóm chat. |
| Jobs-to-be-done | *"Khi thấy một mã tăng mạnh, tôi muốn hiểu vì sao bằng ngôn ngữ dễ hiểu và có nguồn đáng tin, để tránh mua theo tin đồn."* |
| Thành công = | Hiểu được lý do biến động; tin cậy vì thấy nguồn báo chính thống. |

**Persona phản diện (anti-persona):** trader phái sinh scalping cần dữ liệu realtime từng giây — TenPoint **không** phục vụ nhóm này (cron 8h/lần, không phải realtime feed).

---

## 2. User Journey — Luồng chính "Bản tin buổi sáng"

| # | Giai đoạn | Hành động của Minh | Điều hệ thống phải làm | Cảm xúc / Rủi ro rớt |
|---|---|---|---|---|
| 1 | **Mở web 6:45 sáng** (mobile) | Gõ domain / mở bookmark / bấm link từ Zalo | Trang digest render < 1.5s, hiển thị ngay tin mới nhất. Header ghi rõ "Cập nhật lúc 06:00 hôm nay" | ⚠️ Rớt nếu load chậm hoặc không rõ dữ liệu có mới không |
| 2 | **Quét bảng digest** | Đọc lướt cột `Tóm tắt thông tin` | Số liệu then chốt **in đậm** để mắt bắt được trong 1 giây; mỗi dòng ≤ 6 dòng text trên mobile | ⚠️ Rớt nếu tóm tắt dài dòng như bài báo gốc |
| 3 | **Lọc theo danh mục** | Bấm chip mã `FPT`, `HPG` (hoặc dùng watchlist đã lưu) | Filter theo mã, phản hồi tức thì, URL thay đổi (`/?ma=FPT,HPG`) để chia sẻ được | ⚠️ Rớt nếu filter reload cả trang |
| 4 | **Đọc tóm tắt** | Đọc 3–5 tin liên quan | Tóm tắt tự đứng độc lập — không cần đọc bài gốc vẫn nắm đủ số liệu | — |
| 5 | **Kiểm chứng nguồn** | Bấm `cafef.vn` ở cột Source | Short link `/r/{code}` → 302 về bài gốc, mở tab mới, ghi nhận click | ⚠️ Rớt nếu link chết → mất niềm tin ngay lập tức |
| 6 | **Đào sâu 1 mã** | Bấm mã `CMG` → trang `/ma/CMG` | Dòng thời gian tin theo mã + Research Note nếu có | — |
| 7 | **Đọc Research Note** | Đọc luận điểm đầu tư 1→4 | Layout editorial 1 cột, tiêu đề in đậm, luận điểm đánh số, có disclaimer cuối bài | — |
| 8 | **Quay lại ngày mai** | — | Email/Zalo 7:00 sáng (V2) hoặc thói quen bookmark | Mục tiêu: D7 retention ≥ 35% |

**Journey phụ P2 (Broker):** bước 3 → chọn nhiều mã → bước 5' **"Copy bản tin"** → dán vào Zalo/Email đã kèm sẵn link nguồn.

---

## 3. Epic & User Stories

### EPIC 1 — Bảng tin tổng hợp (News Digest)

**US-1.1** — Là *nhà đầu tư cá nhân*, tôi muốn xem bảng tin mới nhất ngay khi mở trang chủ, để nắm bắt thị trường mà không cần thao tác gì thêm.
- **AC1:** Given tôi truy cập `/`, When trang render xong, Then tôi thấy bảng tin sắp xếp giảm dần theo ngày đăng, mặc định 7 ngày gần nhất.
- **AC2:** Given có ≥ 1 tin trong CSDL, When trang load, Then mỗi dòng hiển thị đủ 5 cột: Ngày, Mã CK, Tóm tắt thông tin, Source, Loại tin.
- **AC3:** Given không có tin nào, When trang load, Then hiện empty state "Chưa có tin trong khoảng thời gian này" — không hiện bảng rỗng.

**US-1.2** — Là *nhà đầu tư*, tôi muốn thấy các con số then chốt được in đậm trong tóm tắt, để nhận ra thông tin quan trọng trong 1 giây.
- **AC1:** Given tin có số liệu định lượng, When render tóm tắt, Then ít nhất 1 và tối đa 3 cụm số liệu được bọc `<strong>`.
- **AC2:** Given tóm tắt chứa markdown `**...**`, When render, Then chỉ `<strong>` và `<em>` được phép — mọi HTML khác bị sanitize.

**US-1.3** — Là *nhà đầu tư*, tôi muốn biết dữ liệu được cập nhật lần cuối lúc nào, để tin rằng mình đang đọc tin mới.
- **AC1:** Given trang digest, When render, Then header hiển thị "Cập nhật lần cuối: HH:mm DD/MM/YYYY" theo giờ VN (UTC+7).
- **AC2:** Given lần crawl gần nhất > 12 giờ trước, When render, Then hiện cảnh báo nhẹ "Dữ liệu có thể chưa mới nhất".

**US-1.4** — Là *nhà đầu tư*, tôi muốn phân trang / tải thêm tin cũ, để xem lại diễn biến những ngày trước.
- **AC1:** Given > 20 tin, When tôi cuộn xuống cuối, Then có nút "Xem thêm" tải tiếp 20 tin, không reload trang.

**US-1.5** — Là *nhà đầu tư*, tôi muốn thấy tiêu đề động kiểu "8 tin mới đáng chú ý — PVS · FPT · VIC", để biết ngay hôm nay có gì nổi bật.
- **AC1:** Given bộ lọc hiện tại trả về N tin, When render, Then tiêu đề ghi đúng N và liệt kê tối đa 3 mã xuất hiện nhiều nhất, nối bằng dấu `·`.

### EPIC 2 — Lọc & tìm kiếm theo mã / loại tin

**US-2.1** — Là *nhà đầu tư*, tôi muốn lọc tin theo mã chứng khoán tôi đang nắm giữ, để bỏ qua tin không liên quan.
- **AC1:** Given tôi nhập/chọn `FPT`, When áp dụng filter, Then chỉ hiện tin có `FPT` trong danh sách mã.
- **AC2:** Given tôi chọn nhiều mã, When áp dụng, Then dùng logic OR (tin khớp bất kỳ mã nào).
- **AC3:** Given filter đang áp dụng, When tôi copy URL và mở ở máy khác, Then filter được giữ nguyên (state nằm trong query string).

**US-2.2** — Là *nhà đầu tư*, tôi muốn lọc theo Loại tin (Vĩ mô / Ngành / Doanh nghiệp...), để tập trung vào nhóm thông tin mình quan tâm.
- **AC1:** Given tôi chọn `Doanh nghiệp`, Then chỉ hiện tin thuộc loại đó; số lượng tin hiển thị cạnh mỗi nhãn.

**US-2.3** — Là *nhà đầu tư*, tôi muốn lọc theo khoảng ngày, để xem lại tin của một tuần cụ thể.
- **AC1:** Given tôi chọn `từ 20/08` `đến 27/08`, Then chỉ hiện tin có ngày đăng trong khoảng (bao gồm 2 đầu mút, theo giờ VN).

**US-2.4** — Là *nhà đầu tư*, tôi muốn tìm kiếm toàn văn theo từ khoá (VD "data center", "khối ngoại"), để tra cứu chủ đề.
- **AC1:** Given tôi gõ ≥ 2 ký tự, When submit, Then kết quả khớp tiêu đề hoặc tóm tắt, **có dấu và không dấu đều tìm được** (unaccent).

**US-2.5** — Là *nhà đầu tư*, tôi muốn xoá nhanh toàn bộ bộ lọc, để quay về trạng thái ban đầu.

### EPIC 3 — Trích dẫn nguồn & Short link

**US-3.1** — Là *nhà đầu tư*, tôi muốn thấy tên nguồn báo rút gọn (`cafef.vn`), để đánh giá độ tin cậy trước khi bấm.
- **AC1:** Given tin có `source_url`, Then cột Source hiển thị domain đã bỏ `www.`, có gạch chân.

**US-3.2** — Là *nhà đầu tư*, tôi muốn bấm vào nguồn để đọc bài gốc, để kiểm chứng thông tin.
- **AC1:** Given tôi bấm Source, Then mở tab mới (`target="_blank" rel="nofollow noopener noreferrer"`).
- **AC2:** Given short link `/r/{code}` tồn tại, When gọi, Then trả `302` về đúng `source_url` gốc.
- **AC3:** Given `{code}` không tồn tại, When gọi, Then trả `404` kèm trang lỗi thân thiện tiếng Việt — **không** redirect về trang chủ (tránh mở redirect).
- **AC4:** Given người dùng bấm short link, Then hệ thống tăng `click_count` bất đồng bộ, **không** làm chậm redirect (< 50ms).

**US-3.3** — Là *chủ sản phẩm*, tôi muốn short link chỉ trỏ tới domain trong danh sách nguồn đã duyệt, để tránh bị lợi dụng làm open redirect.
- **AC1:** Given URL đích không thuộc allowlist domain, When tạo short link, Then từ chối tạo và ghi log cảnh báo.

### EPIC 4 — Trang theo mã & Research Note

**US-4.1** — Là *nhà đầu tư*, tôi muốn bấm vào mã CK để xem tất cả tin của mã đó theo dòng thời gian, để hiểu câu chuyện của doanh nghiệp.
- **AC1:** Given tôi bấm `CMG`, Then điều hướng tới `/ma/CMG` hiển thị tên công ty, ngành, và danh sách tin giảm dần theo thời gian.

**US-4.2** — Là *nhà đầu tư*, tôi muốn đọc Research Note dạng luận điểm đầu tư có cấu trúc, để hiểu bức tranh dài hạn chứ không chỉ tin rời rạc.
- **AC1:** Given research note tồn tại, When mở `/nhan-dinh/{slug}`, Then hiển thị tiêu đề uppercase in đậm, section "LUẬN ĐIỂM ĐẦU TƯ", các luận điểm đánh số với dòng dẫn bold + đoạn phân tích.
- **AC2:** Given mọi research note, Then cuối bài **bắt buộc** có disclaimer "Nội dung chỉ mang tính thông tin, không phải khuyến nghị đầu tư."
- **AC3:** Given research note có mã liên quan, Then hiển thị link chéo về `/ma/{mã}`.

**US-4.3** — Là *nhà đầu tư*, tôi muốn từ Research Note xem được các tin gần đây của cùng mã, để đối chiếu luận điểm với diễn biến thực tế.

### EPIC 5 — Dành cho môi giới (P2)

**US-5.1** — Là *môi giới*, tôi muốn copy toàn bộ bản tin đang lọc dưới dạng text kèm link nguồn, để gửi khách qua Zalo.
- **AC1:** Given tôi bấm "Copy bản tin", Then clipboard chứa text thuần, mỗi tin 1 khối 3 dòng:
  ```
  [27/08/2026] [VIC, VHM, HPG]
  ACBS ước 5.588 tỷ đồng chảy vào 117 cổ phiếu Việt Nam...
  Nguồn: https://tenpoint.vn/r/b3Qm9d
  ```
  *(Sửa 16/09: bản đầu dùng `Tóm tắt — Nguồn: <link>` trên một dòng. Đổi vì hai lý do: gạch ngang dài `—` nằm trong danh sách cấm chống AI-slop, và khối 3 dòng đọc dễ hơn hẳn khi dán vào Zalo.)*
- **AC2:** Given copy thành công, Then hiện toast xác nhận trong 2 giây.

**US-5.2** — Là *môi giới*, tôi muốn lưu danh sách mã theo dõi ngay trên trình duyệt (không cần đăng nhập), để mỗi sáng mở lên là có sẵn bộ lọc.
- **AC1:** Given tôi lưu watchlist, When quay lại sau, Then bộ lọc tự áp dụng từ `localStorage`.

### EPIC 6 — Vận hành & Tin cậy (nội bộ / admin)

**US-6.1** — Là *vận hành viên*, tôi muốn cronjob tự chạy 8 giờ/lần và ghi log kết quả, để biết pipeline có khoẻ không.
- **AC1:** Given tới giờ chạy, Then hệ thống crawl toàn bộ nguồn đang bật và ghi 1 bản ghi `crawl_run` (thời điểm, số bài mới, số lỗi, thời lượng).
- **AC2:** Given 1 nguồn lỗi, Then các nguồn còn lại **vẫn chạy tiếp** — không fail toàn bộ job.

**US-6.2** — Là *vận hành viên*, tôi muốn có endpoint kiểm tra sức khoẻ hệ thống, để tích hợp monitoring.
- **AC1:** `GET /healthz` trả `200` kèm trạng thái DB và thời điểm crawl gần nhất.

**US-6.3** — Là *vận hành viên*, tôi muốn kích hoạt crawl thủ công khi có tin nóng, để không phải chờ đủ 8 tiếng.
- **AC1:** Given tôi gọi `POST /api/v1/admin/ingest` kèm token hợp lệ, Then job chạy ngay; token sai → `401`.

---

## 4. Ưu tiên MoSCoW cho MVP v1

| Mức | Hạng mục |
|---|---|
| **MUST** | Bảng digest 5 cột đúng ảnh tham chiếu · Tóm tắt giàu số liệu, bold số then chốt · Short link `/r/{code}` + tracking click · Cron 8h/lần + log run · Lọc theo mã CK · Lọc theo Loại tin · Trang `/ma/{mã}` · Trang Research Note · Hiển thị "Cập nhật lần cuối" · Disclaimer đầu tư · Responsive mobile · `/healthz` |
| **SHOULD** | Tìm kiếm toàn văn (có/không dấu) · Lọc khoảng ngày · Watchlist lưu localStorage · Copy bản tin cho môi giới · Phân trang "Xem thêm" · Dedupe tin trùng · SEO metadata + sitemap |
| **COULD** | Dark mode · RSS feed đầu ra · Export CSV · Sentiment tin · Trang chủ đề/ngành · Bộ lọc nguồn báo |
| **WON'T (v1)** | Đăng nhập/tài khoản người dùng · Realtime websocket · Dữ liệu giá & biểu đồ · App mobile native · Bình luận · Thanh toán/subscription · Telegram/Zalo bot |

---

## 5. Phạm vi MVP v1

| # | Tính năng | Lý do (gắn với persona & bằng chứng từ ảnh) | Ưu tiên |
|---|---|---|---|
| 1 | Digest table 5 cột | Chính là Màn hình B trong ảnh tham chiếu — trái tim sản phẩm | P0 |
| 2 | Tóm tắt cô đọng, bold số liệu | Đặc trưng khác biệt rõ nhất trong ảnh; JTBD của P1 "biết con số trong 5 phút" | P0 |
| 3 | Gắn mã CK cho từng tin | Biến tin tức → tín hiệu theo danh mục; P1 pain #3 | P0 |
| 4 | Phân loại Loại tin | Cột có sẵn trong ảnh; phục vụ lọc của P2 | P0 |
| 5 | Short link + đếm click | Yêu cầu gốc của khách hàng; xây niềm tin (P3) & uy tín môi giới (P2) | P0 |
| 6 | Cronjob 8h/lần | Yêu cầu gốc; đảm bảo nội dung luôn mới | P0 |
| 7 | Lọc theo mã & loại tin | JTBD chính của P1 | P0 |
| 8 | Trang `/ma/{mã}` | Dòng thời gian theo mã; cũng là cửa ngõ SEO | P0 |
| 9 | Research Note | Chính là Màn hình A trong ảnh; chiều sâu giữ chân người dùng | P0 |
| 10 | Responsive mobile-first | 6:45 sáng người dùng đọc trên điện thoại | P0 |
| 11 | Dedupe tin trùng | P1 pain #4 — 5 báo cùng 1 sự kiện | P1 |
| 12 | Tìm kiếm không dấu | Thói quen gõ tiếng Việt không dấu rất phổ biến | P1 |
| 13 | Copy bản tin | Mở khoá persona P2, kênh lan truyền tự nhiên | P1 |
| 14 | Watchlist localStorage | Giảm ma sát mỗi sáng, không cần login | P1 |

---

## 6. Tình huống biên & Câu hỏi mở cần PM quyết

| # | Vấn đề | Phương án | **Khuyến nghị của BA** |
|---|---|---|---|
| Q1 | **Tin trùng lặp** giữa nhiều báo về cùng sự kiện | (a) hiện tất cả (b) gộp thành 1 dòng, cột Source liệt kê nhiều nguồn (c) chọn 1 nguồn Tier-1 đại diện, ẩn phần còn lại | **(c) cho v1** — giữ bảng sạch đúng tinh thần ảnh; lưu cluster để v2 nâng cấp thành (b) |
| Q2 | **Số liệu mâu thuẫn** giữa các báo | (a) tin nguồn Tier-1 (b) gắn cờ "cần kiểm chứng" (c) hiện cả hai | **(a) + lưu log mâu thuẫn** để QA rà soát |
| Q3 | **Mã CK chỉ được nhắc thoáng qua** (VD bài về VN-Index liệt kê 20 mã) | (a) tag hết (b) chỉ tag mã là chủ thể chính (c) tag có trọng số primary/mentioned | **(c)** — lưu cả hai mức, mặc định lọc theo `primary`, cho phép bật "gồm cả mã được nhắc tới" |
| Q4 | **Có cần đăng nhập không?** | (a) không login ở v1 (b) login để lưu watchlist | **(a)** — watchlist qua `localStorage`; login chỉ khi cần đồng bộ đa thiết bị (V2) |
| Q5 | **Tóm tắt do LLM sinh có thể sai số liệu** | (a) publish thẳng (b) kiểm tra tự động: mọi số trong tóm tắt phải xuất hiện trong bài gốc (c) người duyệt thủ công | **(b) bắt buộc cho v1** — validator đối chiếu số; tin không đạt → xuống hàng chờ duyệt |
| Q6 | Bài gốc bị xoá / đổi URL → short link chết | (a) kệ (b) health-check định kỳ, gắn nhãn "nguồn không còn khả dụng" | **(b)** — kiểm tra hàng tuần |
| Q7 | Nguồn đổi cấu trúc HTML → parser gãy | (a) phát hiện khi có người báo (b) cảnh báo khi 1 nguồn trả 0 bài trong 2 lần chạy liên tiếp | **(b)** — đưa vào `crawl_run` metrics |
| Q8 | Bản quyền: lưu và hiển thị bao nhiêu nội dung gốc? | — | **Chỉ hiển thị tóm tắt do hệ thống sinh + link nguồn. Không hiển thị công khai full-text.** Xem BA2. |

---

## 7. Yêu cầu phi chức năng (từ góc nhìn người dùng)

| Nhóm | Yêu cầu | Chỉ tiêu đo được |
|---|---|---|
| **Hiệu năng** | Trang digest tải nhanh trên 4G | LCP < 2.0s (mobile 4G), TTFB < 400ms, JS bundle trang digest < 150KB gzip |
| **Thiết bị** | **Mobile-first** (6:45 sáng đọc trên điện thoại), nhưng bảng 5 cột cần desktop để đọc thoải mái | ≥ 375px: bảng chuyển sang dạng **card xếp dọc**, giữ nguyên thứ tự thông tin. ≥ 1024px: bảng 5 cột đúng ảnh |
| **Tính mới của dữ liệu** | Người dùng luôn biết độ tươi của dữ liệu | Luôn hiện timestamp "Cập nhật lần cuối"; cảnh báo khi > 12h |
| **Offline** | Không yêu cầu offline đầy đủ ở v1 | Chỉ cần cache trang đã xem qua HTTP cache; PWA để V2 |
| **Accessibility (tiếng Việt)** | Đọc được, tương phản tốt, hỗ trợ screen reader | WCAG 2.2 AA: contrast ≥ 4.5:1 (chữ đen `#111` trên nền giấy `#F7F5F0` đạt ~16:1) · `lang="vi"` · bảng có `<th scope="col">` + `<caption>` · focus ring rõ · font hiển thị đầy đủ dấu tiếng Việt (ư, ơ, ừ, ỗ) — **bắt buộc kiểm tra font có bộ ký tự Việt trước khi chốt** |
| **SEO** | Tin tức tài chính có lưu lượng tìm kiếm cao theo mã | SSR toàn bộ trang công khai · thẻ `title`/`description` động theo mã và ngày · Open Graph · `sitemap.xml` (gồm `/ma/{mã}` cho toàn bộ mã) · JSON-LD `NewsArticle` cho research note · canonical URL |
| **Tin cậy** | Nội dung tài chính cần minh bạch | 100% tin có link nguồn hoạt động · disclaimer trên mọi trang có nội dung phân tích · ghi rõ tóm tắt được tạo tự động |
| **Ngôn ngữ** | Toàn bộ giao diện tiếng Việt | Định dạng số kiểu VN: `2.323 tỷ đồng`, `17,9%` (dấu `.` phân nhóm nghìn, dấu `,` thập phân) · ngày `DD/MM/YYYY` · múi giờ `Asia/Ho_Chi_Minh` |

---

## 8. Khung giờ Cronjob (8 giờ/lần) — đề xuất & lập luận

Giờ giao dịch VN: **9:00–11:30** (phiên sáng) và **13:00–14:45** (phiên chiều, ATC kết thúc 14:45).

| Lần chạy | Giờ VN (UTC+7) | Cron (UTC) | Lý do |
|---|---|---|---|
| **Run A** | **06:00** | `0 23 * * *` | Quét toàn bộ tin đêm qua + tin sáng sớm. Sẵn sàng **trước** khi người dùng P1 mở web lúc 6:45 và **trước** khi broker P2 soạn bản tin sáng. Đây là lần chạy quan trọng nhất. |
| **Run B** | **14:00** | `0 7 * * *` | Bắt tin phát sinh trong phiên sáng và đầu phiên chiều; kịp trước ATC 14:45 để người dùng còn hành động được trong ngày. |
| **Run C** | **22:00** | `0 15 * * *` | Gom tin tổng kết phiên, bản tin chiều tối của các báo, tin quốc tế ảnh hưởng phiên hôm sau. |

**Cấu hình:** `CRON_SPEC="0 23,7,15 * * *"` chạy theo UTC, hoặc đặt `TZ=Asia/Ho_Chi_Minh` và dùng `0 6,14,22 * * *`.

**Lập luận chọn mốc:** ràng buộc 8h/lần chỉ cho 3 lần chạy/ngày, nên phải đặt đúng vào 3 thời điểm người dùng thực sự cần: **trước phiên** (06:00), **giữa phiên/trước ATC** (14:00), **sau phiên** (22:00). Cách đặt "00:00 / 08:00 / 16:00" là sai nhịp — 08:00 là quá muộn cho bản tin sáng của broker, còn 16:00 thì đã lỡ mất cửa sổ hành động trong phiên.

**Yêu cầu bổ sung:** cho phép kích hoạt thủ công (US-6.3) để xử lý tin nóng ngoài lịch — đây là van an toàn cho ràng buộc 8h.

---

## 9. Đặc tả yêu cầu "tóm tắt giàu số liệu" (đo lường được)

Đây là đặc trưng nổi bật nhất trong ảnh tham chiếu. Chuyển thành yêu cầu sản phẩm:

| Mã YC | Yêu cầu | Tiêu chí chấp nhận (đo được) |
|---|---|---|
| **NUM-1** | Độ dài tóm tắt cô đọng | 60–150 từ (ảnh: ~70–110 từ/tin). QA reject nếu > 180 từ |
| **NUM-2** | Mật độ số liệu | ≥ 3 dữ kiện định lượng (số + đơn vị) cho tin có số. Ảnh tham chiếu đạt 8–15 |
| **NUM-3** | Nhấn mạnh số then chốt | Đúng **1–3** cụm được bold (`**...**`). 0 cụm → thiếu điểm nhấn; > 3 → loãng |
| **NUM-4** | Không bịa số | **100%** con số trong tóm tắt phải truy vết được về bài gốc. Validator tự động đối chiếu; sai → chặn publish |
| **NUM-5** | Giữ nguyên đơn vị | Bắt buộc kèm đơn vị đầy đủ: `tỷ đồng`, `%`, `điểm %`, `USD/tấn`, `triệu USD`, `MW`, `tấn` |
| **NUM-6** | Định dạng số kiểu VN | `408.000 tỷ đồng`, `+1,67%`, `24.150 đồng/cp` |
| **NUM-7** | Giọng văn trung lập | Không dùng từ khuyến nghị: "nên mua", "khuyến nghị", "chắc chắn tăng". Danh sách từ cấm được QA kiểm tra tự động |
| **NUM-8** | Tự đứng độc lập | Người đọc hiểu đủ ý mà không cần mở bài gốc (đánh giá thủ công trên mẫu 30 tin/tuần, đạt ≥ 90%) |

---

## 10. Bàn giao cho các vai trò tiếp theo

- **PM:** chốt 8 câu hỏi mở ở mục 6; chốt phạm vi MVP mục 5.
- **SA:** cần thiết kế cho `primary` vs `mentioned` ticker (Q3), cluster dedupe (Q1), validator số liệu (Q5/NUM-4), allowlist domain cho short link (US-3.3).
- **UI/UX:** mục 7 (responsive breakpoint bảng → card), mục 9 (quy tắc bold), accessibility tiếng Việt.
- **FE:** filter state nằm trong query string (US-2.1 AC3); SSR cho SEO; sanitize markdown (US-1.2 AC2).
- **BE:** cron 3 mốc mục 8; `/healthz`; short link 302 < 50ms + đếm click bất đồng bộ; nguồn lỗi không làm hỏng cả job.
- **QA:** mục 9 là bộ tiêu chí test nội dung; mục 6 là bộ test tình huống biên.
