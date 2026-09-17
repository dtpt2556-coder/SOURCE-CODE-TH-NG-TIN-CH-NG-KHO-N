# BA3 — Cạnh tranh, Khác biệt hoá & Mô hình kinh doanh

**Sản phẩm:** TenPoint — nền tảng tổng hợp & tóm tắt tin tức chứng khoán Việt Nam
**Góc nhìn:** BA senior #3 (Competitive / Differentiation / Business Model)
**Ngày:** 16/09/2026
**Input bắt buộc:** `/Users/vinh/Desktop/Tho/tenpoint/docs/00-inputs/reference-images.md` + 2 ảnh gốc tại root workspace

**Ý tưởng gốc (verbatim):**
> "cập nhật tin tức chứng khoán, tổng hợp nội dung và trích dẫn short link, tin tức có cronjob cập nhật 8h 1 lần, nguồn data là những trang uy tín"

**Quy ước tin cậy trong tài liệu này:**
- `[ẢNH]` = suy luận trực tiếp từ 2 reference image → độ tin cậy cao, dùng làm nền tảng thiết kế.
- `[GIẢ ĐỊNH]` = phán đoán thị trường của BA, **cần validate** bằng desk research / phỏng vấn trước khi đưa vào quyết định đầu tư.
- Mọi con số giá tham chiếu VN đều ở mức `[GIẢ ĐỊNH]` và phải được verify lại bằng bảng giá công khai của đối thủ tại thời điểm ra quyết định.

---

## 0. Tóm tắt điều hành (đọc 60 giây)

TenPoint không cạnh tranh ở **lượng tin** (CafeF/Vietstock thắng tuyệt đối) và cũng không cạnh tranh ở **độ sâu dữ liệu** (Vietstock Finance/Wichart thắng tuyệt đối). TenPoint cạnh tranh ở **lớp chưng cất (distillation layer)**: biến 200–400 tin/ngày thành ~10 tin có số liệu, gắn mã CK, trích nguồn minh bạch, đọc hết trong 3 phút.

Ba khoảng trống lớn nhất của thị trường VN hiện nay:
1. **Không ai làm lớp distill có kỷ luật số liệu.** Tin tức VN hiện tồn tại ở 2 thái cực: bài báo dài + giật tít (CafeF/VnExpress/Thanh Niên) hoặc PDF research 15 trang (SSI/HSC/VND/MBS). Không có "3 phút đọc xong, biết đủ số".
2. **Không ai coi tin tức là dữ liệu có cấu trúc theo mã.** Tin ở các portal là bài viết; ở TenPoint (theo ảnh) tin là **record** có `ngày | [mã CK] | tóm tắt số liệu | nguồn | loại tin` → có thể truy vấn, cảnh báo, dựng timeline, ghép watchlist.
3. **Không ai vừa nhanh vừa truy vết được nguồn.** Room Zalo/Telegram nhanh nhưng không kiểm chứng và có xung đột lợi ích; báo chí kiểm chứng được nhưng chậm/nhiễu. TenPoint đứng giữa: tóm tắt nhanh + **short link về bài gốc trên mọi dòng** `[ẢNH]`.

Khuyến nghị: **MVP miễn phí hoàn toàn, tiền đến từ affiliate mở tài khoản CK + hợp đồng white-label với môi giới** (2 dòng tiền sớm nhất, không đặt paywall lên sản phẩm chưa có thói quen). Freemium và API B2B để dành V2/V3.

---

## 1. Bản đồ cạnh tranh Việt Nam

### 1.0 Khung phân nhóm

Trước khi đi từng đối thủ, cần thấy rõ thị trường chia làm **5 nhóm theo "công việc" người dùng thuê sản phẩm làm (JTBD)**:

| Nhóm | Đối thủ tiêu biểu | JTBD họ phục vụ | TenPoint có đụng không? |
|---|---|---|---|
| A. Portal tin tức tài chính | CafeF, Vietstock (mảng tin), VnExpress Kinh doanh, Thanh Niên | "Cho tôi biết hôm nay thị trường có gì" | **Đụng trực tiếp** (nhưng ở lớp trên, không tranh nguồn) |
| B. Nền tảng dữ liệu & phân tích | Vietstock Finance, Wichart, Fireant, Simplize | "Cho tôi số liệu để tự phân tích" | Bổ trợ — có thể là đối tác/khách hàng API |
| C. Kênh của CTCK | DNSE Senses, TCInvest, VPS bản tin | "Cho tôi lý do để đặt lệnh (và mở TK ở đây)" | **Đối thủ về sự chú ý, đối tác về doanh thu** |
| D. Research house | Bản tin sáng SSI/HSC/VND/MBS | "Cho tôi quan điểm chuyên gia có mô hình" | Không đụng trực tiếp — TenPoint là lớp tổng hợp phía trên |
| E. Kênh xã hội | Nhóm Zalo/Telegram bản tin sáng | "Cho tôi cảm giác được chăm sóc + tin sớm" | **Đụng trực tiếp về thói quen buổi sáng** |

Insight quan trọng: **đối thủ nguy hiểm nhất của TenPoint không phải CafeF, mà là nhóm Zalo/Telegram bản tin sáng** — vì họ đã chiếm được đúng *thời điểm* và *thiết bị* mà TenPoint muốn chiếm (7h–8h30 sáng, trên điện thoại).

---

### 1.1 CafeF (VCCorp)

| | |
|---|---|
| **Mô hình** | Portal tin tức tài chính, doanh thu quảng cáo display + native/booking bài. Traffic-driven. |
| **Điểm mạnh** | (1) Quy mô traffic lớn nhất mảng tài chính VN `[GIẢ ĐỊNH – cần verify SimilarWeb]`; (2) SEO cực mạnh, gần như luôn top 1–3 cho query "tin tức + tên mã"; (3) Tốc độ đăng tin rất nhanh, có phóng viên thường trú tại sàn; (4) Thương hiệu mặc định trong đầu NĐT F0 — "có gì thì lên CafeF xem"; (5) Có sẵn bảng giá, dữ liệu cơ bản theo mã. |
| **Điểm yếu** | (1) **Volume là bug, không phải feature** — 200+ bài/ngày, người đọc không có cách nào biết đâu là tin quan trọng; (2) Giật tít để lấy click → nhiễu tín hiệu, tiêu đề thường không chứa con số; (3) Trải nghiệm đọc nặng quảng cáo, đặc biệt trên mobile; (4) Gắn mã CK chỉ ở mức tag thủ công, không đầy đủ, không phải mapping có kỷ luật (một tin ngành ngân hàng hiếm khi được tag đủ 11 mã như trong `[ẢNH]` của TenPoint); (5) Không cá nhân hoá theo watchlist; (6) Không có bản tóm tắt — người dùng phải đọc hết 800–1.500 từ để lấy 5 con số. |
| **Khoảng trống bỏ ngỏ** | **Lớp lọc + chưng cất.** CafeF tối ưu cho pageview nên về mặt cấu trúc kinh doanh họ *không thể* làm sản phẩm "chỉ 10 tin/ngày" — điều đó giết chính mô hình doanh thu của họ. Đây là classic innovator's dilemma và là hào phòng thủ (moat) hành vi tốt nhất của TenPoint. |

---

### 1.2 Vietstock (vietstock.vn — mảng tin & cộng đồng)

| | |
|---|---|
| **Mô hình** | Portal tin + dữ liệu + diễn đàn; doanh thu quảng cáo, gói dữ liệu trả phí, đào tạo/sự kiện. |
| **Điểm mạnh** | (1) Độ phủ **công bố thông tin (CBTT)** và tin doanh nghiệp tốt nhất — nghị quyết HĐQT, giao dịch nội bộ, chốt quyền; (2) Lịch sự kiện doanh nghiệp đầy đủ; (3) Uy tín lâu năm với NĐT có kinh nghiệm; (4) Kho dữ liệu lịch sử sâu. |
| **Điểm yếu** | (1) UI/UX cũ, mật độ thông tin dày, mobile experience kém; (2) Tin vẫn ở dạng nguyên bản hoặc copy nguyên văn CBTT — không có lớp diễn giải "điều này nghĩa là gì"; (3) Cấu trúc paywall phân mảnh, người dùng khó biết mình đang mua gì; (4) Không có tóm tắt AI, không có cá nhân hoá. |
| **Khoảng trống bỏ ngỏ** | **Diễn giải CBTT.** Một nghị quyết HĐQT phát hành riêng lẻ 200 triệu cp là dữ kiện; "pha loãng ~15%, giá phát hành thấp hơn thị giá 22%, tiền dùng trả nợ chứ không mở rộng SXKD" là thông tin. Vietstock dừng ở dữ kiện. TenPoint có thể chiếm lớp thứ hai — và đây là loại tin **có tác động giá cao nhất**. |

---

### 1.3 Fireant

| | |
|---|---|
| **Mô hình** | Nền tảng dữ liệu + mạng xã hội đầu tư; freemium (gói Pro theo tháng/năm) + API `[GIẢ ĐỊNH về cơ cấu giá]`. |
| **Điểm mạnh** | (1) **Trang theo mã (`/symbol/FPT`) làm tốt nhất thị trường** — gom giá, dữ liệu, tin, thảo luận về một chỗ; (2) Đã có khái niệm "dòng tin theo mã" — gần TenPoint nhất về mặt cấu trúc; (3) Screener và bộ lọc mạnh; (4) Cộng đồng tạo nội dung liên tục → chi phí nội dung thấp. |
| **Điểm yếu** | (1) Dòng tin theo mã chủ yếu là **headline aggregation** — link + tiêu đề, người dùng vẫn phải click ra ngoài đọc; (2) Nội dung cộng đồng nhiễu và không kiểm chứng, lẫn nhiều "phím hàng"; (3) Không có lớp tóm tắt chuẩn hoá, càng không có research note dạng luận điểm; (4) Sản phẩm thiên về trader kỹ thuật hơn nhà đầu tư đọc tin. |
| **Khoảng trống bỏ ngỏ** | **Chất lượng của dòng tin theo mã.** Fireant đã chứng minh *nhu cầu* tồn tại (người dùng vào trang mã để xem tin), nhưng chưa giải quyết *chất lượng*. TenPoint đi đúng chỗ đó với tóm tắt giàu số liệu — nghĩa là thị trường đã được Fireant "làm nóng" hộ, rủi ro giáo dục thị trường thấp. |

---

### 1.4 Simplize

| | |
|---|---|
| **Mô hình** | SaaS phân tích cổ phiếu cho NĐT cá nhân; freemium với gói trả phí theo tháng/năm. |
| **Điểm mạnh** | (1) **UX/UI hiện đại nhất nhóm** — đặt chuẩn mới cho fintech VN, người dùng trẻ thích; (2) Chấm điểm doanh nghiệp, định giá, "báo cáo phân tích tự động" dễ hiểu cho người không chuyên; (3) Tổng hợp khuyến nghị & giá mục tiêu từ nhiều CTCK — bước đầu của tư duy "tổng hợp đa nguồn"; (4) Đã chứng minh **NĐT cá nhân VN sẵn sàng trả tiền** cho công cụ phân tích — validation quan trọng cho mô hình freemium. |
| **Điểm yếu** | (1) Lõi là **dữ liệu & định giá**, không phải **nhịp tin**; nội dung cập nhật theo quý (mùa BCTC) chứ không theo ngày; (2) Không phải sản phẩm mở-mỗi-sáng → tần suất sử dụng thấp hơn, khó tạo thói quen; (3) Coverage phân tích sâu tập trung nhóm large/mid-cap; (4) Không có kênh push (Zalo/Telegram/email sáng) mạnh. |
| **Khoảng trống bỏ ngỏ** | **Tần suất.** Simplize là sản phẩm "mở 1 lần/tuần khi cân nhắc mua"; TenPoint là sản phẩm "mở mỗi sáng". Hai tần suất khác nhau → **có thể cùng tồn tại, thậm chí Simplize là đối tác phân phối tiềm năng** hơn là đối thủ trực diện. |

---

### 1.5 Wichart

| | |
|---|---|
| **Mô hình** | Công cụ dữ liệu & biểu đồ (vĩ mô, ngành, BCTC); freemium/subscription hướng người dùng chuyên. |
| **Điểm mạnh** | (1) Kho dữ liệu vĩ mô và dữ liệu ngành sâu, chuẩn hoá tốt; (2) Biểu đồ tuỳ biến mạnh, được analyst và người viết báo cáo dùng nhiều; (3) Uy tín trong giới chuyên nghiệp — "nguồn biểu đồ" mặc định trên nhiều bài phân tích. |
| **Điểm yếu** | (1) **Gần như không có lớp tin tức** — thuần dữ liệu; (2) Learning curve cao, không dành cho NĐT phổ thông; (3) Không cá nhân hoá theo danh mục; (4) Không có kênh push hằng ngày. |
| **Khoảng trống bỏ ngỏ** | **Nối dữ liệu với câu chuyện.** Wichart cho biết "giá thép HRC giảm 12% trong 3 tháng"; không cho biết "vì sao" và "tin nào liên quan HPG/HSG/NKG". TenPoint `[ẢNH – dòng CBAM gắn HPG, HSG, NKG, TIS]` đã làm đúng nối đó. **Cơ hội đối tác:** nhúng biểu đồ Wichart vào Research Note, hoặc ngược lại bán feed tin của TenPoint cho Wichart. |

---

### 1.6 DNSE Senses

| | |
|---|---|
| **Mô hình** | Sản phẩm tin tức/insight thuộc hệ sinh thái CTCK DNSE (app Entrade X). **Doanh thu thực = phí giao dịch + margin**, nội dung là chi phí acquisition. |
| **Điểm mạnh** | (1) Đã ứng dụng AI cho tóm tắt/sentiment — đối thủ gần TenPoint nhất về *công nghệ*; (2) Miễn phí (được trợ giá bởi mảng môi giới) → TenPoint **không thể thắng bằng giá**; (3) Gắn thẳng vào luồng đặt lệnh: đọc tin → bấm mua, conversion cực ngắn; (4) Có dữ liệu hành vi giao dịch thật của người dùng — thứ TenPoint không có. |
| **Điểm yếu** | (1) **Bị khoá trong hệ sinh thái DNSE** — người dùng VPS/SSI/TCBS/VNDirect không dùng được, mà đó là phần lớn thị trường; (2) Không thể trung lập: mọi nội dung đều phục vụ mục tiêu tăng giao dịch (nghiêng về tin kích thích giao dịch, hiếm khi khuyên "ngồi im"); (3) Web/SEO yếu vì sản phẩm sống trong app → gần như không có kênh organic; (4) Ưu tiên phát triển sản phẩm luôn nằm sau lõi giao dịch. |
| **Khoảng trống bỏ ngỏ** | **Tính trung lập và tính đa-CTCK.** Chính xác vì Senses thuộc DNSE mà nó không bao giờ trở thành "nguồn tin chung của thị trường". TenPoint có thể định vị là **Switzerland của tin tức chứng khoán VN** — trung lập, đa nguồn, ai dùng CTCK nào cũng được. |

---

### 1.7 TCInvest (TCBS) / VPS bản tin

| | |
|---|---|
| **Mô hình** | Bản tin & khuyến nghị đính kèm tài khoản giao dịch; miễn phí, thu từ phí GD/margin/trái phiếu/chứng quyền. |
| **Điểm mạnh** | (1) **Kênh phân phối khổng lồ** — trực tiếp tới hàng triệu tài khoản đang mở `[GIẢ ĐỊNH về quy mô]`; (2) Dữ liệu thị trường realtime và chính xác (họ là nguồn gốc); (3) Có đội analyst in-house và tiếp cận doanh nghiệp; (4) TCBS có lượng người dùng tự phục vụ rất lớn với chi phí biên gần 0. |
| **Điểm yếu** | (1) Định dạng PDF/bài dài trong app — không tìm kiếm được, không lưu trữ tra cứu được, không chia sẻ được ra ngoài; (2) **Xung đột lợi ích cấu trúc** — khuyến nghị "bán"/"đứng ngoài" gần như không tồn tại; (3) Chỉ nói về thị trường và nhóm mã họ coverage, không tổng hợp tin từ báo chí ngoài; (4) Không SEO, không truy cập được nếu không có tài khoản → nội dung "chết" sau 1 ngày; (5) Không có khái niệm short link về nguồn gốc — người đọc phải tin, không kiểm chứng được. |
| **Khoảng trống bỏ ngỏ** | **Tính kiểm chứng và tính lưu trữ.** Bản tin CTCK là nội dung dùng-một-lần. TenPoint xây **kho lưu trữ có cấu trúc, tra cứu được, có nguồn gốc** — tài sản tích luỹ theo thời gian và là nền tảng cho toàn bộ chiến lược SEO (mục 5). |

---

### 1.8 Vietstock Finance (gói dữ liệu trả phí)

| | |
|---|---|
| **Mô hình** | Subscription dữ liệu tài chính chuyên sâu (BCTC chuẩn hoá, dữ liệu lịch sử, công cụ so sánh, export). Khách hàng: analyst, sinh viên tài chính, NĐT chuyên nghiệp. |
| **Điểm mạnh** | (1) **BCTC chuẩn hoá sâu và dài nhất thị trường** — tài sản khó sao chép; (2) Export Excel phục vụ dựng mô hình; (3) Đã chứng minh mô hình subscription dữ liệu tài chính VN có người trả tiền; (4) Là nguồn dữ liệu nền cho nhiều sản phẩm khác. |
| **Điểm yếu** | (1) Sản phẩm cho **người dựng mô hình**, không phải người đọc tin — sai persona so với TenPoint; (2) UI cũ, đường cong học tập dốc; (3) Không phải sản phẩm dùng hằng ngày; (4) Tách biệt khỏi dòng tin — dữ liệu và tin không nói chuyện với nhau. |
| **Khoảng trống bỏ ngỏ** | **Không phải đối thủ — là nguồn dữ liệu bổ trợ tiềm năng.** Research Note `[ẢNH – CMG: P/E TTM ~14 lần, ROE ~12%, giá 24.150đ/cp]` cần đúng loại dữ liệu này để neo luận điểm định giá. Cân nhắc mua gói dữ liệu thay vì tự crawl BCTC ở giai đoạn V2. |

---

### 1.9 Nhóm Zalo / Telegram "bản tin sáng"

| | |
|---|---|
| **Mô hình** | Môi giới / KOL / "room" vận hành nhóm miễn phí để gom lead → ăn hoa hồng phí giao dịch & margin; hoặc room VIP thu phí 500k–5 triệu/tháng `[GIẢ ĐỊNH]`. |
| **Điểm mạnh** | (1) **Khớp hoàn hảo với hành vi người Việt** — Zalo là app mặc định, không cần cài thêm gì; (2) Push notification thẳng vào màn hình khoá lúc 7h sáng — không sản phẩm web nào cạnh tranh được về mặt attention; (3) Tương tác 2 chiều, "cảm giác được chăm sóc bởi người thật"; (4) Chi phí công nghệ bằng 0, scale bằng con người; (5) Nội dung ngắn, đúng format mobile — **đây chính là format TenPoint đang nhắm tới** `[ẢNH – bảng digest ngắn]`. |
| **Điểm yếu** | (1) **Không kiểm chứng được** — tin đồn, tin "nội bộ", không có nguồn gốc; (2) Xung đột lợi ích tối đa: người gửi tin kiếm tiền khi bạn giao dịch nhiều; (3) Không lưu trữ, không tìm kiếm — tin trôi mất sau 200 tin nhắn; (4) Không cá nhân hoá — mọi người trong nhóm nhận cùng nội dung bất kể danh mục; (5) Chất lượng dao động hoàn toàn theo người vận hành; (6) Không scale được (1 môi giới chỉ chăm được vài trăm khách). |
| **Khoảng trống bỏ ngỏ** | **Vừa nhanh vừa đáng tin.** Đây là khoảng trống có giá trị nhất. TenPoint giữ nguyên format và kênh (push ngắn gọn buổi sáng) nhưng thay "tin đồn" bằng "tin có số + short link nguồn gốc". **Đồng thời đây là cơ hội white-label lớn nhất**: chính các môi giới này đang phải tự soạn bản tin mỗi sáng — bán cho họ công cụ, đừng đánh nhau với họ (xem Mô hình #5). |

---

### 1.10 Bản tin sáng của CTCK (SSI Research / HSC / VNDirect / MBS)

| | |
|---|---|
| **Mô hình** | Research house — nội dung miễn phí làm thương hiệu & marketing cho mảng môi giới; doanh thu thật từ khách tổ chức, IB, môi giới. |
| **Điểm mạnh** | (1) **Uy tín chuyên môn cao nhất thị trường** — có mô hình định giá, có tiếp cận ban lãnh đạo doanh nghiệp, có site visit; (2) Chất lượng luận điểm đầu tư là chuẩn mực — chính là chuẩn mà Research Note của TenPoint `[ẢNH – màn A]` đang cố mô phỏng; (3) Được truyền thông trích dẫn rộng rãi → hiệu ứng khuếch đại. |
| **Điểm yếu** | (1) **PDF** — không machine-readable, không mobile-friendly, không deep-link tới một luận điểm cụ thể; (2) Phát hành ~7h30–8h00 và thường dài 3–15 trang → sai định dạng cho người đọc lúc uống cà phê; (3) **Coverage hẹp** — mỗi CTCK chỉ theo sát ~50–80 mã `[GIẢ ĐỊNH]`, bỏ trống toàn bộ mid/small-cap nơi NĐT cá nhân hoạt động nhiều nhất; (4) Không cá nhân hoá theo watchlist; (5) Thiên lệch tích cực do quan hệ với doanh nghiệp; (6) Mỗi bản tin chỉ là một góc nhìn — người đọc muốn đối chiếu SSI vs HSC vs MBS phải tự tải 3 file PDF. |
| **Khoảng trống bỏ ngỏ** | (a) **Coverage đuôi dài** — ~1.600 mã niêm yết/ĐKGD, research house phủ chưa tới 10%. TenPoint dùng LLM có thể phủ 100% với chi phí biên gần 0. (b) **Tổng hợp đa góc nhìn** — "SSI, HSC và MBS nói gì về cùng sự kiện này" là sản phẩm chưa ai làm. |

---

### 1.11 Tổng hợp khoảng trống → cơ hội

| # | Khoảng trống | Ai đang bỏ ngỏ | TenPoint khai thác bằng | Độ khó |
|---|---|---|---|---|
| G1 | Chưng cất tin thành tóm tắt giàu số liệu | Toàn bộ nhóm A, D | Pipeline LLM + prompt ép trích số + in đậm số then chốt `[ẢNH]` | Trung bình |
| G2 | Tin tức như dữ liệu có cấu trúc theo mã | Nhóm A, D, E | Bảng `ngày/mã/tóm tắt/nguồn/loại` + mapping n-n `[ẢNH]` | Trung bình |
| G3 | Vừa nhanh vừa truy vết được nguồn | Nhóm E (Zalo/Telegram) | Short link domain hiển thị công khai mọi dòng `[ẢNH]` | Thấp |
| G4 | Cá nhân hoá theo watchlist | Nhóm A, D, E | Tiêu đề động "8 tin mới đáng chú ý — PVS · FPT · VIC" `[ẢNH]` | Thấp |
| G5 | Trung lập, đa-CTCK | Nhóm C | Không sở hữu mảng môi giới; công khai chính sách affiliate | Thấp (chiến lược) |
| G6 | Coverage đuôi dài mid/small-cap | Nhóm D | LLM phủ toàn bộ mã với chi phí biên thấp | Trung bình |
| G7 | Lưu trữ tra cứu được + SEO | Nhóm C, D, E | Trang `/ma/{TICKER}`, `/ngay/{date}`, `/chu-de/{topic}` | Trung bình |
| G8 | Tổng hợp đa góc nhìn về cùng sự kiện | Tất cả | Story clustering + "3 báo nói gì" | Cao |

---

## 2. Định vị khác biệt của TenPoint

### 2.1 Bằng chứng từ 2 ảnh (nền tảng của mọi lập luận bên dưới)

**(a) Tóm tắt cực kỳ cô đọng, giàu số liệu, in đậm số then chốt**

Bằng chứng trực tiếp `[ẢNH – màn B]`: dòng tin về gói tín dụng DNNVV dài ~70 từ nhưng chứa **408.000 tỷ đồng** (in đậm), 220.000 tỷ, Agribank 70.000, BIDV/VCB/CTG mỗi bên 50.000 tỷ, "thấp hơn ít nhất 1 điểm %", 188.000 tỷ, "0,5–2 điểm %", "8,8% DNNVV tiếp cận được vốn so với trên 47% ở doanh nghiệp lớn" — **11 dữ kiện định lượng trong 70 từ**. Bài gốc trên vnexpress.net nhiều khả năng dài 800–1.200 từ.

→ Tỷ lệ nén ước tính **~12–15×**, nhưng **mật độ số liệu tăng** chứ không giảm. Đây không phải "tóm tắt" theo nghĩa thông thường (cắt bớt), mà là **chưng cất** (giữ lại phần có giá trị quyết định, bỏ phần tường thuật). Con số then chốt nhất mỗi dòng được in đậm → mắt bắt được trong 0,5 giây khi lướt.

Đối chiếu: không đối thủ nào trong mục 1 có sản phẩm ở định dạng này. Bản tin CTCK cô đọng nhưng là văn xuôi và dài; headline CafeF ngắn nhưng rỗng số liệu.

**(b) Gắn mã CK cho từng tin → biến tin tức thành *tín hiệu theo mã***

Bằng chứng `[ẢNH – màn B]`: cột `Mã CK` đứng **thứ hai, ngay sau ngày** — trước cả nội dung. Đây là tuyên bố thiết kế: *mã chứng khoán là khoá chính, không phải tiêu đề bài báo*.

- Quan hệ **nhiều-nhiều**: 1 tin ↔ nhiều mã (dòng DNNVV gắn **11 mã**: BID, VCB, CTG, SHB, MSB, STB, BVB, NAB, NVB, SGB, TPB) và 1 mã ↔ nhiều tin (VIC xuất hiện ở 3/5 dòng hiển thị).
- Mapping có suy luận, không phải keyword matching: tin CBAM của EU **không nhắc tên HPG/HSG/NKG/TIS** nhưng vẫn được gắn đúng 4 mã thép chịu tác động. Tương tự, tin nhập khẩu thịt gắn DBC, BAF, MML, HAG. **Đây là năng lực lõi không copy nhanh được** — nó đòi hỏi tri thức ngành, không chỉ NLP.
- Tiêu đề **"8 tin mới đáng chú ý — PVS · FPT · VIC."** chứng minh digest được **cắt theo watchlist người dùng**, không phải feed chung.

→ Hệ quả kiến trúc: một khi tin là record có khoá `ticker`, mọi tính năng cao cấp mở ra miễn phí — timeline theo mã, alert watchlist, sentiment theo mã, đếm mật độ tin theo ngành, event study. **Đối thủ coi tin là bài viết thì không làm được gì trong số đó.**

**(c) Trích dẫn nguồn minh bạch qua short link**

Bằng chứng `[ẢNH – màn B]`: cột `Source` hiển thị **domain rút gọn có gạch chân** (`vnexpress.net`, `cafef.vn`, `thanhnien.vn`) trên **100% số dòng**, không có ngoại lệ.

→ Ba tác dụng, đúng thứ tự tầm quan trọng:
1. **Chống lại nghi ngờ "AI bịa số"** — rào cản niềm tin lớn nhất của sản phẩm tóm tắt bằng LLM trong lĩnh vực tiền bạc. Người dùng kiểm chứng được bất cứ lúc nào bằng 1 cú chạm.
2. **Phân biệt tuyệt đối với room Zalo/Telegram** — nơi tin không có nguồn gốc.
3. **Đo lường được** — CTR short link là proxy trực tiếp cho "tin này đáng quan tâm", dùng để xếp hạng nội dung (xem mục 6).

Lưu ý chiến lược: hiển thị **domain chứ không phải tên bài** là lựa chọn có chủ đích — nó biến nguồn thành **tín hiệu uy tín** (vnexpress.net ≠ một fanpage) thay vì một cái link rỗng.

**(d) Research Note dạng luận điểm đầu tư có cấu trúc**

Bằng chứng `[ẢNH – màn A]`: 4 luận điểm đánh số, mỗi luận điểm = **1 câu khẳng định in đậm + 1 đoạn chứng minh bằng số**. Cấu trúc luận điểm chuẩn sell-side:

1. Kinh doanh cốt lõi (doanh thu 2.323 tỷ, +5,1% YoY; biên gộp 17,7%→17,9%; LNST CĐ mẹ −20,3% còn 75 tỷ; chi phí lãi vay **+48%**)
2. Chất lượng tài sản (thị phần Data Center ~12%, cloud nội địa >25%, khách hàng Samsung SDS/IBM/Honda…)
3. Động lực tăng trưởng có điều kiện (>250 triệu USD, 30 MW → >100 MW)
4. Định giá (24.150 đ/cp, P/E TTM ~14 lần, ROE ~12%)

Điểm quan trọng nhất về **giọng điệu**: tiêu đề nêu thẳng mặt tiêu cực ("LỢI NHUẬN CĐ MẸ GIẢM – CHU KỲ ĐẦU TƯ DATA CENTER CHƯA TẠO RA LỢI NHUẬN TƯƠNG XỨNG"), luận điểm 3 tự phản biện ("chưa thể coi là lợi nhuận đã chắc chắn"), luận điểm 4 kết bằng điều kiện ("Biên an toàn chỉ thực sự xuất hiện **nếu** CMC chứng minh được…").

→ **Đây là tài sản định vị mạnh nhất mà TenPoint có.** Không có khuyến nghị MUA/BÁN, không có giá mục tiêu, không phím hàng. Trong một thị trường mà 100% nội dung miễn phí đều có động cơ khiến bạn đặt lệnh, **sự không-khuyến-nghị chính là điểm khác biệt**. Nó cũng giảm mạnh rủi ro pháp lý (xem mục 8).

---

### 2.2 Value Proposition

> **TenPoint biến hàng trăm bài báo chứng khoán mỗi ngày thành ~10 tin đã chưng cất — giàu số liệu, gắn sẵn mã cổ phiếu bạn nắm giữ, và luôn dẫn nguồn gốc bằng một cú chạm — để bạn nắm hết thông tin quan trọng của thị trường trong 3 phút buổi sáng, thay vì 45 phút lướt tin.**

**Ba gạch đầu dòng chứng minh (mỗi gạch neo vào bằng chứng ảnh):**

- **Chưng cất, không phải rút gọn.** Một dòng digest ~70 từ chứa 11 dữ kiện định lượng và in đậm con số quyết định `[ẢNH – dòng 408.000 tỷ]`. Bạn không đánh đổi độ chính xác để lấy tốc độ — bạn nhận được *nhiều số liệu hơn trên mỗi giây đọc* so với việc đọc bài gốc.

- **Tin tức được sắp theo mã của bạn, không theo lịch xuất bản của toà soạn.** Mỗi tin được gắn đầy đủ các mã chịu tác động — kể cả khi bài gốc không hề nhắc tên mã (tin CBAM của EU → HPG, HSG, NKG, TIS `[ẢNH]`) — và digest được cắt theo watchlist ("8 tin mới đáng chú ý — PVS · FPT · VIC" `[ẢNH]`). Không đối thủ nào ở VN làm mapping suy luận theo tác động ngành ở quy mô này.

- **Mọi khẳng định đều kiểm chứng được, và chúng tôi không bán lệnh cho bạn.** 100% dòng tin có short link về bài gốc trên nguồn uy tín `[ẢNH – cột Source]`; Research Note viết dạng luận điểm có điều kiện, nêu cả mặt xấu, và **không đưa khuyến nghị mua/bán** `[ẢNH – màn A]`. TenPoint không sở hữu công ty chứng khoán nào, nên không kiếm tiền khi bạn giao dịch nhiều hơn.

---

## 3. Positioning statement & Tên gọi / Tagline

### 3.1 Positioning statement (khung Geoffrey Moore)

> **Dành cho** nhà đầu tư chứng khoán cá nhân Việt Nam theo dõi 5–30 mã, **những người** không có 45 phút mỗi sáng để lướt CafeF, Vietstock, 3 nhóm Zalo và 4 bản tin CTCK rồi tự ghép lại xem tin nào liên quan danh mục của mình,
> **TenPoint là** nền tảng chưng cất tin tức chứng khoán **mà** mỗi ngày chỉ đưa khoảng 10 tin đã nén thành số liệu, gắn sẵn mã chịu tác động và dẫn nguồn gốc bằng short link.
> **Khác với** các portal tin tức (CafeF, Vietstock) vốn tối ưu cho pageview nên buộc phải đăng nhiều và giật tít, và khác với các room Zalo/Telegram hay bản tin CTCK vốn nhanh nhưng không truy vết được nguồn và luôn có động cơ khiến bạn đặt lệnh,
> **TenPoint** trung lập với mọi công ty chứng khoán, không đưa khuyến nghị mua/bán, và mọi con số đều kiểm chứng được trong một cú chạm.

### 3.2 Định vị trên bản đồ 2 trục

```
                    ĐỘ TIN CẬY / TRUY VẾT NGUỒN  (cao)
                                 ▲
      SSI/HSC/VND/MBS ●          │          ● TenPoint  ← ô trống
      Vietstock ●                │            (nhanh + đáng tin + cô đọng)
      Vietstock Finance ●        │
                                 │      ● Simplize
      CafeF ●  ───────────────────────────────────────►  TỐC ĐỘ & ĐỘ CÔ ĐỌNG (cao)
                                 │      ● Fireant   ● DNSE Senses
                                 │
                                 │              ● Room Zalo/Telegram
                                 ▼
                             (thấp)
```
Ô phải-trên (nhanh + cô đọng + truy vết được) hiện **trống**. Đó là ô TenPoint phải chiếm và phòng thủ.

### 3.3 Tên gọi & Tagline

**Về cái tên "TenPoint":** đọc được cả 2 nghĩa — *10 tin/ngày* và *10 điểm (luận điểm)*. Trong tiếng Việt còn gợi "mười điểm" = chất lượng tuyệt đối. Giữ nguyên, tên tốt.

| Ngữ cảnh | Tiếng Việt | Tiếng Anh |
|---|---|---|
| **Tagline chính** (đặt cạnh logo) | **Mười tin. Ba phút. Đủ cả phiên.** | **Ten stories. Three minutes. The whole market.** |
| Tagline phụ (landing hero) | Tin tức chứng khoán, chưng cất thành số liệu. | Vietnam stock news, distilled to the numbers. |
| Nhấn khác biệt "theo mã" | Tin tức của thị trường, sắp theo mã của bạn. | The market's news, indexed by your tickers. |
| Nhấn minh bạch | Mỗi con số đều có nguồn. Không phím hàng. | Every number sourced. No stock tips. |
| Nhấn cho Research Note | Luận điểm, không phải khuyến nghị. | Theses, not recommendations. |
| Dòng mở email 7h sáng | *"Sáng nay có 8 tin chạm danh mục của bạn."* | *"8 stories touched your watchlist overnight."* |
| Pitch 1 câu cho nhà đầu tư/đối tác | Bloomberg Brief cho thị trường Việt Nam — nhưng miễn phí, và biết bạn đang nắm mã nào. | A Bloomberg Brief for Vietnam — free, and it knows what you hold. |

**Chống chỉ định đặt tagline:** tránh mọi từ gợi lợi nhuận ("cơ hội", "đón sóng", "khuyến nghị", "cổ phiếu tiềm năng"). Chúng kéo TenPoint về ô room Zalo/Telegram và tạo rủi ro pháp lý (mục 8.6).

---

## 4. Mô hình kinh doanh — 6 phương án

> **Nguyên tắc xuyên suốt:** giai đoạn đầu **không đặt bất kỳ ma sát nào** lên sản phẩm cốt lõi. Thứ TenPoint cần trong 6 tháng đầu là **thói quen buổi sáng**, không phải doanh thu. Mọi đồng doanh thu sớm nên đến từ *bên thứ ba*, không từ người đọc.

### PA1 — Miễn phí + affiliate mở tài khoản chứng khoán

| Hạng mục | Nội dung |
|---|---|
| **Ai trả tiền** | Công ty chứng khoán (DNSE, TCBS, VPS, SSI, VNDirect…) trả cho TenPoint mỗi tài khoản mở mới có phát sinh giao dịch. |
| **Giá tham chiếu VN** `[GIẢ ĐỊNH – phải verify]` | CPA ~100.000–400.000 đ/tài khoản active; hoặc revenue-share 10–30% phí giao dịch trong 6–12 tháng đầu. Với 50.000 MAU và tỷ lệ chuyển đổi 1,5% → ~750 TK/tháng × 250k ≈ **~190 triệu đ/tháng**. |
| **Ưu** | Không ma sát người dùng; triển khai được ngay ở MVP; đòn bẩy tự nhiên (người đọc tin tài chính *chính là* người mở tài khoản); không cần đội sales. |
| **Nhược** | ARPU thấp; phụ thuộc chính sách CTCK có thể thay đổi đơn phương; **rủi ro lớn nhất: xói mòn tính trung lập** — vốn là tài sản định vị số 1 (mục 2.2). |
| **Điều kiện bắt buộc** | (i) Ký **nhiều** CTCK, không độc quyền; (ii) gắn nhãn "Tài trợ" rõ ràng; (iii) **tuyệt đối không** để affiliate ảnh hưởng thứ tự/nội dung tin; (iv) công bố chính sách affiliate trên trang riêng. |
| **Thời điểm** | **MVP → V2.** Dòng tiền đầu tiên. |

### PA2 — Freemium theo watchlist / độ trễ

| Hạng mục | Nội dung |
|---|---|
| **Ai trả tiền** | Nhà đầu tư cá nhân nghiêm túc (danh mục ≥300 triệu, theo dõi >10 mã). |
| **Giá tham chiếu VN** `[GIẢ ĐỊNH – bám khung Simplize/Fireant/Vietstock]` | **Free:** 5 mã watchlist, digest theo lịch cron (8h/lần), Research Note trễ 7 ngày.<br>**Pro 79.000 đ/tháng** (790.000 đ/năm, ~17% chiết khấu): watchlist không giới hạn, alert tức thời, Research Note ngay khi phát hành, timeline đầy đủ theo mã, export.<br>**Pro+ 199.000 đ/tháng:** thêm so sánh đa nguồn, sector radar, truy vấn kho lưu trữ. |
| **Ưu** | Doanh thu định kỳ, biên gộp cao; củng cố tính trung lập (người đọc trả tiền → TenPoint phục vụ người đọc); dễ dự báo. |
| **Nhược** | Tỷ lệ chuyển đổi freemium thực tế thường **1–3%** → cần ≥50.000 MAU mới có ý nghĩa; "độ trễ" là ma sát dễ gây khó chịu; mọi tính năng đặt sau paywall đều làm chậm tăng trưởng organic và SEO. |
| **Thời điểm** | **V2**, chỉ bật khi đã đạt ≥10.000 WAU và retention tuần ≥35%. Bật sớm hơn = giết tăng trưởng. |

### PA3 — API B2B cho CTCK / quỹ / fintech

| Hạng mục | Nội dung |
|---|---|
| **Ai trả tiền** | CTCK (nhúng vào app), công ty quản lý quỹ, ngân hàng (mảng wealth), robo-advisor, doanh nghiệp niêm yết (theo dõi truyền thông). |
| **Sản phẩm** | REST/webhook trả về JSON có cấu trúc đúng schema của `[ẢNH – màn B]`: `{date, tickers[], summary, key_figures[], source_url, source_domain, category, sentiment}`. Đây là **thứ không ai ở VN bán** — bản tin CTCK là PDF, portal chỉ có RSS headline. |
| **Giá tham chiếu VN** `[GIẢ ĐỊNH]` | Starter 20–30 triệu đ/tháng (1 sản phẩm, ≤200 mã); Growth 50–80 triệu đ/tháng (toàn thị trường, webhook realtime); Enterprise 150–300 triệu đ/tháng (white-label + SLA + dữ liệu lịch sử). Hợp đồng năm: 300 triệu – 1,5 tỷ đ. |
| **Ưu** | ARPU cực cao — 5 khách hàng Growth ≈ 4 tỷ đ/năm, bằng ~4.200 thuê bao Pro; doanh thu ổn định theo hợp đồng năm; tận dụng chính hạ tầng đã có, chi phí biên gần 0. |
| **Nhược** | Chu kỳ bán 3–9 tháng; yêu cầu SLA/uptime/pháp chế; **rủi ro bản quyền phóng đại** (bán lại nội dung phái sinh từ báo chí ở quy mô thương mại — xem 8.3); khách hàng CTCK có thể tự xây sau 12 tháng. |
| **Thời điểm** | **V3.** Ngoại lệ: nếu có **1 khách hàng anchor** chịu trả trước, làm pilot ngay từ V2 — nó tài trợ toàn bộ chi phí LLM. |

### PA4 — Bản tin email trả phí (premium newsletter)

| Hạng mục | Nội dung |
|---|---|
| **Ai trả tiền** | Nhà đầu tư cá nhân bận rộn, môi giới, nhân sự IR/tài chính doanh nghiệp. |
| **Sản phẩm** | Free 7h sáng: digest ~10 tin. Paid: thêm bản 15h30 (tổng kết phiên + tin sau giờ), Research Note 2–3 mã/tuần, bản tin chuyên đề ngành hằng tháng. |
| **Giá tham chiếu VN** `[GIẢ ĐỊNH]` | 99.000 đ/tháng hoặc 990.000 đ/năm. Gói doanh nghiệp (≤20 tài khoản): 3–5 triệu đ/tháng. |
| **Ưu** | Chi phí kỹ thuật thấp nhất trong 6 phương án; email là kênh sở hữu (không phụ thuộc thuật toán Google/Zalo); đo lường sạch (open rate, CTR); **danh sách email là tài sản** dùng được cho mọi phương án khác. |
| **Nhược** | Trần doanh thu thấp; phụ thuộc chất lượng biên tập đều đặn; deliverability tiếng Việt vào Gmail/Outlook cần chăm sóc kỹ thuật; khó tạo hiệu ứng lan truyền. |
| **Thời điểm** | Bản **free từ MVP** (đây là công cụ tạo thói quen quan trọng nhất). Bản **paid từ V2**, sau khi open rate ổn định ≥35%. |

### PA5 — White-label cho môi giới / phòng giao dịch ⭐

| Hạng mục | Nội dung |
|---|---|
| **Ai trả tiền** | Môi giới cá nhân và trưởng phòng giao dịch — chính những người đang vận hành nhóm Zalo/Telegram ở mục 1.9. |
| **Bài toán họ đang đau** | Mỗi sáng phải tự đọc báo, tự soạn bản tin, tự gửi cho 200–2.000 khách. Tốn 60–90 phút/ngày, chất lượng không đều, không scale được. **Đây là nỗi đau có sẵn, được định giá sẵn bằng thời gian.** |
| **Sản phẩm** | Bản digest gắn logo + tên + số điện thoại môi giới, tự động đẩy vào Zalo OA/Telegram của họ lúc 7h; kèm dashboard xem khách nào mở, click tin nào (→ **lead scoring** cho môi giới, giá trị cực cao với họ). |
| **Giá tham chiếu VN** `[GIẢ ĐỊNH]` | Môi giới cá nhân: 500.000–1.500.000 đ/tháng. Team/phòng GD (10–30 môi giới): 10–30 triệu đ/tháng. Gói CTCK cấp cho toàn bộ đội môi giới: 100–300 triệu đ/tháng. |
| **Ưu** | **Doanh thu tiền mặt sớm nhất** — bán được ngay khi MVP chạy, không cần đợi quy mô người dùng; biến đối thủ hành vi nguy hiểm nhất (mục 1.9) thành khách hàng; mỗi môi giới mang theo 200–2.000 người dùng cuối → **kênh phân phối B2B2C**; chi phí biên gần 0. |
| **Nhược** | Thương hiệu TenPoint bị che (giảm bằng dòng "Nội dung bởi TenPoint" ở footer); cần hỗ trợ khách hàng; **rủi ro uy tín** nếu môi giới chèn thêm khuyến nghị phím hàng vào bản tin mang nội dung của TenPoint → phải có điều khoản hợp đồng cấm chỉnh sửa nội dung. |
| **Thời điểm** | **Cuối MVP / đầu V2.** Khuyến nghị pilot với 3–5 môi giới ngay tuần 5–6 của MVP. |

### PA6 — Phân phối IR / nội dung tài trợ từ doanh nghiệp niêm yết

| Hạng mục | Nội dung |
|---|---|
| **Ai trả tiền** | Bộ phận IR của doanh nghiệp niêm yết, công ty tư vấn IR. |
| **Giá tham chiếu VN** `[GIẢ ĐỊNH]` | 15–50 triệu đ/bản tin IR được phân phối; 200–500 triệu đ/năm cho gói theo dõi & phân phối thường xuyên. |
| **Ưu** | Biên rất cao, không cần công nghệ mới. |
| **Nhược** | **Rủi ro định vị nghiêm trọng nhất trong 6 phương án.** Nếu người dùng nghi ngờ tin được trả tiền để lên top, toàn bộ value proposition ở mục 2.2 sụp đổ. |
| **Khuyến nghị** | **Không làm ở MVP và V2.** Nếu làm ở V3: đặt ở khu vực tách biệt hoàn toàn, nhãn "IR — doanh nghiệp tự công bố", không bao giờ trộn vào digest chính, không tính vào "10 tin". |

### 4.7 Tổng hợp & khuyến nghị

| PA | Ai trả | ARPU | Time-to-revenue | Rủi ro định vị | Giai đoạn |
|---|---|---|---|---|---|
| PA1 Affiliate | CTCK | Thấp | **Nhanh** | Trung bình | **MVP** ✅ |
| PA2 Freemium | NĐT cá nhân | Trung bình | Chậm | Thấp | V2 |
| PA3 API B2B | CTCK/quỹ | **Rất cao** | Rất chậm | Thấp | V3 |
| PA4 Newsletter | NĐT cá nhân | Thấp–TB | Trung bình | Thấp | Free ở MVP, paid ở V2 |
| PA5 White-label | Môi giới | Cao | **Nhanh nhất** | Trung bình | **Cuối MVP** ✅ |
| PA6 IR | DN niêm yết | Cao | Nhanh | **Rất cao** ⚠️ | V3 hoặc không |

**Khuyến nghị cho MVP: PA1 + PA5.**
Lý do: cả hai đều **không đặt ma sát lên người đọc** (điều kiện sống còn khi chưa có thói quen), cả hai đều có thể bắt đầu thu tiền trong vòng 6–10 tuần, và PA5 còn đóng vai trò kênh phân phối. PA4 chạy song song ở dạng miễn phí vì email là công cụ tạo thói quen, không phải sản phẩm doanh thu ở giai đoạn này.

**Mục tiêu tài chính định hướng 12 tháng** `[GIẢ ĐỊNH – cần mô hình tài chính riêng]`: 40–60% doanh thu từ PA5, 25–35% từ PA1, phần còn lại từ PA2/PA4. PA3 là cú nhảy doanh thu của năm 2.

---

## 5. Chiến lược tăng trưởng & SEO

### 5.1 Vì sao SEO là kênh chiến lược (không chỉ là "làm cho có")

Ba yếu tố cấu trúc khiến tin tức tài chính VN là mảnh đất SEO hiếm có:

1. **Truy vấn có ý định cực rõ và lặp lại vĩnh viễn.** "FPT có tin gì mới", "vì sao HPG giảm", "VIC tin tức hôm nay" — được gõ mỗi ngày, bởi những người có tiền và đang ra quyết định. Giá trị thương mại trên mỗi lượt truy cập cao hơn nhiều so với tin tức thông thường.
2. **Đuôi dài gần như bỏ trống.** ~1.600 mã niêm yết/ĐKGD; research house phủ <10%; báo chí chỉ viết về ~100–200 mã nóng. Với mỗi mã mid/small-cap, một trang tổng hợp đầy đủ và cập nhật có thể lên top mà **gần như không có cạnh tranh**.
3. **Nội dung của TenPoint có cấu trúc sẵn.** Tóm tắt giàu số liệu + mapping mã + phân loại `[ẢNH]` là nguyên liệu lý tưởng cho structured data và cho khối trả lời trích dẫn (featured snippet / AI Overview).

⚠️ **Cảnh báo nghiêm túc:** Google Helpful Content áp dụng rất mạnh với nội dung YMYL (Your Money Your Life) và nội dung sinh bởi AI ở quy mô lớn. Chiến lược bên dưới **chỉ an toàn nếu** mỗi trang có giá trị tổng hợp thật (đa nguồn + phân tích + dữ liệu), dẫn nguồn đầy đủ, và có tín hiệu E-E-A-T (tác giả thật, trang giới thiệu đội ngũ, phương pháp luận công khai). Sinh 1.600 trang rỗng = bị phạt.

### 5.2 Kiến trúc URL & cụm nội dung

| Loại trang | URL | Nội dung | Ưu tiên | Quy mô |
|---|---|---|---|---|
| **Trang theo mã** ⭐ | `/ma/FPT` | Timeline tin theo mã, tóm tắt giàu số, "trong 30 ngày qua có N tin", link Research Note, dữ liệu cơ bản | **P0** | ~1.600 trang (rollout theo đợt: 100 mã có tin → mở rộng) |
| Trang theo ngày | `/ngay/2026-09-16` | Digest hôm đó, có `prev/next` → internal link chain rất mạnh | P0 | ~365/năm |
| Trang chủ đề/ngành | `/chu-de/ngan-hang`, `/chu-de/thep`, `/chu-de/khoi-ngoai` | Tin theo ngành/chủ đề + diễn giải bối cảnh | P1 | 30–60 trang chất lượng cao |
| **Trang sự kiện** ⭐ | `/su-kien/cbam-thep-eu`, `/su-kien/ftse-nang-hang` | Cụm tin về một sự kiện kéo dài, đa nguồn, dòng thời gian — **loại trang không đối thủ nào ở VN có** | P1 | 50–150 trang, giá trị/trang rất cao |
| Research Note | `/research/cmg-q2-2026` | Nội dung `[ẢNH – màn A]`; nội dung sâu nhất, nam châm backlink | P1 | 5–10/tuần ở V2 |
| Trang so sánh | `/so-sanh/HPG-vs-HSG` | Tin & số liệu 2 mã cạnh nhau | P2 | Tổ hợp theo ngành |
| Glossary | `/thuat-ngu/cbam`, `/thuat-ngu/pe-ttm` | Giải thích thuật ngữ xuất hiện trong digest, internal link về trang mã | P2 | 100–200 trang |

**Chiến thuật internal link (quyết định thành bại):** Mỗi dòng digest đã chứa 2–11 mã `[ẢNH]` → mỗi dòng tự động sinh 2–11 internal link tới `/ma/{TICKER}`. Với ~10 tin/ngày × ~4 mã/tin × 365 ngày ≈ **~14.600 internal link/năm được sinh tự động**, phân bổ đúng theo mức độ được nhắc đến của từng mã. Đây là hệ quả miễn phí của quyết định kiến trúc "tin là record có khoá ticker" ở mục 2.1(b) — đối thủ coi tin là bài viết không có được.

### 5.3 Structured data & kỹ thuật

- `NewsArticle` / `Dataset` schema cho mỗi digest item; `FinancialProduct` / `Organization` cho trang mã.
- **`Article.citation` + link `rel="nofollow"` có chủ đích tới bài gốc** — vừa minh bạch, vừa tín hiệu tổng hợp hợp pháp.
- `FAQPage` trên trang mã cho các câu hỏi "FPT có tin gì mới nhất?", "Vì sao FPT tăng hôm nay?".
- News Sitemap riêng, cập nhật theo cron; `lastmod` chính xác.
- SSR/SSG bắt buộc (không CSR) — Core Web Vitals là lợi thế cạnh tranh trực tiếp vì CafeF nặng quảng cáo và chậm trên mobile.
- Cron 8h/lần `[ý tưởng gốc]` → **lệch giờ chạy so với 8h/14h/22h chẵn** (ví dụ 5h45 / 13h45 / 21h45) để digest sáng sẵn sàng **trước** 7h khi email gửi đi.

### 5.4 Kênh phân phối

| Kênh | Vai trò | Cơ chế | Ưu tiên |
|---|---|---|---|
| **Email 7h00 sáng** ⭐ | **Cỗ máy tạo thói quen** | Digest ~10 tin, cá nhân hoá theo watchlist; tiêu đề động `[ẢNH – "8 tin mới đáng chú ý — PVS · FPT · VIC"]` | **P0 — MVP** |
| **Zalo OA** ⭐ | Kênh có tiếp cận cao nhất tại VN | Push 7h + alert watchlist; là kênh gốc của đối thủ 1.9 nên người dùng đã quen | **P0 — MVP** (lưu ý: ZNS có chi phí/tin và quy định nội dung — cần verify) |
| **Telegram bot** | Người dùng chuyên/trader | Bot `/watch FPT`, `/tin FPT`, push tức thời, miễn phí, API dễ | P0 — MVP |
| **Web (SEO)** | Tăng trưởng dài hạn, chi phí biên giảm dần | Mục 5.2 | P0 — MVP (kiến trúc), thu hoạch từ tháng 4–6 |
| **RSS/JSON feed** | Tín hiệu trung lập + nguyên liệu cho tích hợp | Feed công khai theo mã `/ma/FPT/rss` | P1 — V2 |
| Facebook Group/Page | Thu hút F0 | Đăng 1 tin "wow" nhất mỗi ngày dưới dạng ảnh bảng, link về trang ngày | P1 |
| Embed widget | Backlink & phân phối | Widget "Tin mới nhất về {MÃ}" nhúng vào blog/diễn đàn, có backlink | P2 — V2 |

### 5.5 Vòng lặp tăng trưởng

**Vòng chính (content loop):** Cron 8h/lần → thêm tin có cấu trúc → mỗi tin sinh 2–11 internal link tới trang mã → trang mã dày lên → xếp hạng cao hơn cho query "{mã} tin tức" → traffic organic → người dùng lập watchlist → nhận email/Zalo mỗi sáng → quay lại → **dữ liệu click cho biết tin nào quan trọng** → cải thiện xếp hạng nội dung → digest tốt hơn → giữ chân tốt hơn.

**Vòng phụ (broker loop, gắn với PA5):** Môi giới dùng white-label → gửi cho 200–2.000 khách mỗi sáng → khách thấy footer "Nội dung bởi TenPoint" → một phần tự đăng ký trực tiếp → tăng người dùng mà không tốn chi phí marketing.

**Vòng phụ (citation loop):** Research Note chất lượng `[ẢNH – màn A]` → báo chí/diễn đàn trích dẫn → backlink chất lượng cao → tăng authority toàn domain → toàn bộ 1.600 trang mã được hưởng lợi.

---

## 6. North Star Metric + KPI Tree

### 6.1 North Star Metric

> **NSM = Số người dùng có "buổi sáng chất lượng" trong tuần (Weekly Quality Morning Readers — WQMR)**
>
> Định nghĩa một *buổi sáng chất lượng*: người dùng mở digest **trong khoảng 6h–9h**, **xem ≥3 tin**, trong đó **≥1 tin gắn mã thuộc watchlist** của họ.
> NSM = số người dùng có **≥3 buổi sáng chất lượng trong 7 ngày gần nhất**.

**Vì sao chọn chỉ số này:**
- Nó gói trọn cả 4 yếu tố khác biệt trong mục 2: *đọc được nhanh* (≥3 tin trong một lần mở), *liên quan cá nhân* (≥1 tin trúng watchlist), *đúng thời điểm* (6h–9h), *thành thói quen* (≥3 ngày/tuần).
- Nó là **chỉ số dẫn (leading) cho mọi mô hình doanh thu**: người có thói quen buổi sáng mới là người trả tiền Pro (PA2), mới click affiliate (PA1), và mới khiến môi giới muốn mua white-label (PA5).
- **Không thể gian lận bằng clickbait** — khác hẳn pageview, thứ mà CafeF buộc phải tối ưu.

**Vì sao KHÔNG chọn:**
- ❌ *DAU/MAU:* đếm cả người ghé một lần từ Google, không phản ánh giá trị.
- ❌ *Pageview:* đây chính là chỉ số khiến CafeF phải giật tít — tối ưu nó là tự biến thành đối thủ mình đang muốn thay thế.
- ❌ *Thời gian trên trang:* với TenPoint, **đọc xong nhanh là thành công**, không phải thất bại. Chỉ số này ngược dấu với value proposition.

### 6.2 KPI Tree

```
                    NSM: Weekly Quality Morning Readers (WQMR)
                                     │
   ┌────────────┬──────────────┬─────┴──────┬──────────────┬─────────────┐
   │ ACQUISITION│  ACTIVATION  │ ENGAGEMENT │  RETENTION   │ MONETIZATION│
   └────────────┴──────────────┴────────────┴──────────────┴─────────────┘
         │              │              │             │              │
   ┌─────┴─────┐  ┌─────┴─────┐  ┌─────┴─────┐ ┌─────┴─────┐  ┌─────┴─────┐
   │Organic/th │  │% tạo      │  │Số tin đọc │ │Tỷ lệ quay │  │TK mở qua  │
   │Trang/phiên│  │watchlist  │  │/phiên     │ │lại buổi   │  │affiliate  │
   │Đăng ký    │  │Số mã TB   │  │CTR short  │ │sáng D1/D7 │  │/tháng     │
   │email/Zalo │  │/user      │  │link       │ │/D30       │  │Free→Pro % │
   │Vòng môi   │  │% bật email│  │% tin trúng│ │Email open │  │MRR        │
   │giới (PA5) │  │/Zalo push │  │watchlist  │ │rate       │  │Hợp đồng   │
   │Backlink   │  │Time-to-   │  │Lượt xem   │ │Weekly     │  │white-label│
   │           │  │first-value│  │/mã        │ │churn      │  │ARPU       │
   └───────────┘  └───────────┘  └───────────┘ └───────────┘  └───────────┘
                                     │
                          ┌──────────┴──────────┐
                          │  QUALITY (nền móng) │
                          └──────────┬──────────┘
                    ┌────────────────┼────────────────┐
              Độ chính xác      Chất lượng        Vận hành
              ticker tagging    tóm tắt           pipeline
              (precision/       (số liệu/tin,     (độ trễ crawl→
               recall)          tỷ lệ báo sai,    publish, % nguồn
                                % trùng lặp)      lỗi, chi phí LLM/tin)
```

### 6.3 Bảng chỉ số chi tiết & mục tiêu

| Tầng | Chỉ số | Định nghĩa | Mục tiêu MVP | V2 | Vì sao quan trọng |
|---|---|---|---|---|---|
| **NSM** | WQMR | Xem 6.1 | 300 | 3.500 | Chỉ số duy nhất cả team nhìn |
| Acq | Organic sessions/tháng | GA4 | 5.000 | 60.000 | Kênh chi phí biên giảm dần |
| Acq | Đăng ký email/Zalo | Số mới/tuần | 150 | 1.200 | Kênh sở hữu |
| Act | **% user tạo watchlist** | Trong 7 ngày đầu | **≥45%** | ≥60% | Cổng vào mọi giá trị cá nhân hoá |
| Act | **Số mã theo dõi TB/user** | Median | **≥5** | ≥8 | Càng nhiều mã → càng nhiều tin trúng → càng dính |
| Act | Time-to-first-value | Từ vào trang → thấy tin trúng watchlist đầu tiên | <90 giây | <45 giây | Quyết định tỷ lệ bỏ |
| Eng | **Số tin đọc/phiên** | Số dòng mở rộng/click | **≥3,0** | ≥4,5 | Đo "digest có đáng đọc hết không" |
| Eng | **CTR short link** | Click nguồn / lượt hiển thị dòng | **8–15%** | 10–18% | Xem ghi chú 6.4 ⚠️ |
| Eng | % tin trúng watchlist | Số tin gắn mã trong watchlist / tổng tin digest | ≥25% | ≥40% | Đo chất lượng cá nhân hoá |
| Ret | **Tỷ lệ quay lại buổi sáng D7** | % user mở lại trong khung 6h–9h sau 7 ngày | **≥40%** | ≥55% | Chỉ số dẫn mạnh nhất cho NSM |
| Ret | Tỷ lệ quay lại D30 | | ≥20% | ≥35% | Giá trị dài hạn |
| Ret | Email open rate | | ≥35% | ≥45% | Sức khoẻ kênh sở hữu |
| Ret | Weekly churn | | <12% | <7% | |
| Mon | TK mở qua affiliate | /tháng | 20 (pilot) | 500 | PA1 |
| Mon | Hợp đồng white-label | Số môi giới trả phí | 3 (pilot) | 40 | PA5 |
| Mon | MRR | | >0 | 200 triệu đ | |
| **Qual** | **Ticker tagging precision** | % mã gắn đúng (đánh giá thủ công 100 tin/tuần) | **≥92%** | ≥96% | Gắn sai mã = mất niềm tin ngay lập tức |
| Qual | Ticker tagging recall | % mã đáng lẽ phải gắn mà đã gắn | ≥75% | ≥88% | Sót mã = sót giá trị |
| Qual | Mật độ số liệu | Số dữ kiện định lượng TB/tóm tắt | **≥6** | ≥8 | Đặc trưng lõi `[ẢNH: đạt ~11]` |
| Qual | Tỷ lệ tóm tắt sai số liệu | Kiểm thủ công mẫu ngẫu nhiên | **<1%** | <0,3% | Rủi ro tồn vong |
| Qual | Tỷ lệ trùng lặp | % tin trùng nội dung với tin khác trong digest | <5% | <2% | Trùng lặp phá vỡ lời hứa "10 tin" |
| Qual | Độ trễ crawl→publish | p95 | <45 phút | <15 phút | |
| Qual | **Chi phí LLM/tin đã publish** | Tổng chi phí LLM / số tin | **<1.500 đ** | <800 đ | Quyết định đơn vị kinh tế (8.2) |

### 6.4 ⚠️ Ghi chú quan trọng về CTR short link

CTR short link là chỉ số **hai mặt, phải đọc kèm ngữ cảnh**:
- CTR **quá thấp (<3%)** → có thể người dùng không quan tâm, **hoặc** tóm tắt đã đủ tốt đến mức không cần đọc gốc (đây là *thành công*).
- CTR **quá cao (>25%)** → tóm tắt thiếu thông tin, người dùng buộc phải click ra (đây là *thất bại*).

→ **Không đặt CTR làm mục tiêu tối ưu hoá.** Dùng nó như (a) *tín hiệu xếp hạng* mức độ quan tâm theo tin/theo mã, và (b) *chỉ báo cảnh giới* cho chất lượng tóm tắt. Vùng lành mạnh: **8–15%**.

---

## 7. Roadmap 3 giai đoạn

### Giai đoạn 1 — MVP (4–6 tuần)

**Mục tiêu:** Chứng minh **một người thật sẽ mở TenPoint mỗi sáng thay vì lướt CafeF**. Không phải chứng minh doanh thu, không phải chứng minh quy mô.

**Phạm vi tính năng:**

*Pipeline dữ liệu*
- Crawl 6–8 nguồn uy tín: cafef.vn, vietstock.vn, vnexpress.net (Kinh doanh), thanhnien.vn (Tài chính), tuoitre.vn, baodautu.vn, tinnhanhchungkhoan.vn + CBTT HOSE/HNX `[nguồn trong ẢNH: vnexpress.net, cafef.vn, thanhnien.vn → xác nhận hướng chọn nguồn]`
- Cronjob **8h/lần** `[ý tưởng gốc]` — lệch giờ chạy 5h45 / 13h45 / 21h45
- Khử trùng lặp (dedupe) đa nguồn cho cùng một sự kiện
- LLM tóm tắt với prompt **ép trích số liệu** + đánh dấu số then chốt để in đậm `[ẢNH]`
- Gắn mã CK (n-n) + phân loại `Ngành / Doanh nghiệp / Vĩ mô` `[ẢNH]`
- Sinh short link + lưu domain nguồn `[ẢNH]`

*Sản phẩm người dùng*
- Trang digest theo ngày, thiết kế đúng `[ẢNH – màn B]`: nền giấy ấm, serif, hairline rule, số in đậm
- Watchlist (MVP: localStorage, không bắt đăng nhập — giảm ma sát tối đa)
- Tiêu đề động theo watchlist `[ẢNH – "8 tin mới đáng chú ý — PVS · FPT · VIC"]`
- Trang theo mã `/ma/{TICKER}` cho ~100 mã có tin (nền móng SEO)
- Email 7h00 sáng (đăng ký chỉ cần email)
- Telegram bot cơ bản (`/watch`, `/tin`)
- **5 tính năng "wow" ưu tiên**: W1, W2, W3 (bản nhẹ), W7, W8 (xem mục 9)

*Kinh doanh*
- Pilot PA5 white-label với 3–5 môi giới (bắt đầu tiếp cận từ tuần 3)
- Đàm phán PA1 affiliate với 2–3 CTCK (chưa cần go-live)

**Chưa làm ở MVP:** Research Note tự động (viết tay 2–3 bài để test phản ứng), sentiment score, so sánh đa nguồn đầy đủ, Zalo OA (đợi duyệt), app mobile, thanh toán, phủ toàn bộ 1.600 mã.

**Tiêu chí thoát giai đoạn (phải đạt ĐỒNG THỜI):**

| # | Tiêu chí | Ngưỡng |
|---|---|---|
| 1 | NSM — WQMR | **≥300** |
| 2 | Tỷ lệ quay lại buổi sáng D7 | **≥40%** |
| 3 | Ticker tagging precision (kiểm thủ công 100 tin) | **≥92%** |
| 4 | Tỷ lệ tóm tắt sai số liệu | **<1%** |
| 5 | Số tin đọc/phiên | **≥3,0** |
| 6 | % user tạo watchlist trong 7 ngày đầu | **≥45%** |
| 7 | Chi phí LLM/tin đã publish | **<1.500 đ** |
| 8 | Ít nhất **1 môi giới trả tiền thật** cho white-label | ≥1 |
| 9 | Uptime pipeline (không lỡ chu kỳ cron) | ≥98% |

> **Cổng chặn (kill gate):** Nếu tiêu chí #2 (<25%) hoặc #4 (>2%) không đạt, **dừng mở rộng tính năng** và quay lại sửa chất lượng nội dung. Không được đánh đổi hai chỉ số này để chạy theo tiêu chí khác.

---

### Giai đoạn 2 — V2 "Từ bản tin thành nền tảng" (8–12 tuần sau MVP)

**Mục tiêu:** Chuyển từ *sản phẩm đọc* sang *nền tảng theo dõi*, và bật dòng doanh thu định kỳ đầu tiên.

**Phạm vi tính năng:**
- Tài khoản người dùng + watchlist đồng bộ đa thiết bị; import danh mục từ file CTCK
- **Zalo OA chính thức** + alert watchlist tức thời (W2 bản đầy đủ)
- **Research Note** `[ẢNH – màn A]`: bán tự động, 5–10 mã/tuần, cấu trúc 4 luận điểm, có biên tập viên con người review trước khi publish
- Timeline theo mã đầy đủ (W1) + kho lưu trữ tra cứu được
- Sentiment & thang tác động (W4)
- So sánh đa nguồn cho cùng sự kiện (W3 bản đầy đủ) + trang `/su-kien/*`
- Sector radar (W10)
- Mở rộng SEO: toàn bộ ~1.600 trang mã (rollout theo đợt, chỉ publish trang có ≥1 tin thật), trang chủ đề, trang sự kiện, glossary
- Bật **PA2 freemium** (Free 5 mã / Pro 79k) và **PA4 newsletter trả phí**
- Mở rộng PA5 lên 30–50 môi giới; go-live PA1 affiliate
- Bảng điều khiển vận hành nội bộ: theo dõi chi phí LLM, độ chính xác tagging, sức khoẻ từng nguồn

**Tiêu chí thoát giai đoạn:**

| # | Tiêu chí | Ngưỡng |
|---|---|---|
| 1 | NSM — WQMR | **≥3.500** |
| 2 | Organic sessions/tháng | **≥60.000** |
| 3 | Top 10 Google cho ≥200 truy vấn dạng "{mã} tin tức" | ≥200 |
| 4 | Tỷ lệ chuyển đổi Free→Pro | **≥2%** |
| 5 | MRR | **≥200 triệu đ** |
| 6 | Tỷ lệ quay lại buổi sáng D30 | **≥35%** |
| 7 | Ticker tagging precision | **≥96%** |
| 8 | Hợp đồng white-label đang hoạt động | **≥30** |
| 9 | Đã ký ≥1 LOI/hợp đồng pilot API B2B | ≥1 |

---

### Giai đoạn 3 — V3 "Hạ tầng tin tức thị trường" (3–6 tháng sau V2)

**Mục tiêu:** Trở thành **lớp hạ tầng** mà các sản phẩm tài chính khác xây dựng trên đó — nơi biên lợi nhuận thật nằm.

**Phạm vi tính năng:**
- **API B2B (PA3)** với schema chuẩn, webhook realtime, SLA, dữ liệu lịch sử
- **W5 — "Tin này ảnh hưởng gì đến giá?"**: event study tự động, đối chiếu với lịch sử phản ứng giá T+1/T+3/T+5 của các tin cùng loại trên cùng mã
- **W9 — Living Thesis**: Research Note tự gắn cờ khi có tin mới xác nhận hoặc phản bác một luận điểm
- **W11 — Hỏi kho lưu trữ** bằng ngôn ngữ tự nhiên
- Phủ 100% mã niêm yết & ĐKGD, bao gồm UPCoM
- App mobile (iOS/Android) với push notification gốc
- Bản tiếng Anh cho nhà đầu tư nước ngoài (đón dòng vốn nâng hạng thị trường)
- Bảng điều khiển IR cho doanh nghiệp niêm yết (PA6, khu vực tách biệt)
- Nền tảng white-label đầy đủ cho CTCK

**Tiêu chí thoát giai đoạn / thành công V3:**

| # | Tiêu chí | Ngưỡng |
|---|---|---|
| 1 | NSM — WQMR | **≥25.000** |
| 2 | Tổng ARR | **≥25 tỷ đ** |
| 3 | Khách hàng API B2B trả phí | **≥5** |
| 4 | % doanh thu từ B2B (API + white-label) | **≥50%** |
| 5 | Biên gộp | **≥70%** |
| 6 | Được ≥20 bài báo/báo cáo trích dẫn như nguồn | ≥20 |
| 7 | Tỷ lệ quay lại buổi sáng D30 | ≥45% |

---

## 8. Rủi ro kinh doanh

> Thang: **Xác suất** (Thấp/TB/Cao) × **Tác động** (Thấp/TB/Cao/**Tồn vong**)

### 8.1 Phụ thuộc nguồn tin — `Cao × Cao`

**Rủi ro:** Toàn bộ sản phẩm là hàm số của 6–8 nguồn. Một nguồn chặn IP, dựng paywall, thêm Cloudflare bot protection, hoặc đơn giản là sập trong 6 giờ → digest sáng hôm đó rỗng hoặc lệch. Rủi ro tập trung: `[ẢNH]` cho thấy 5 dòng mẫu chỉ đến từ 3 domain (vnexpress, cafef, thanhnien).

**Giảm thiểu:**
- Đa dạng hoá ngay từ MVP: **≥8 nguồn**, không nguồn nào vượt **25%** số tin publish (đặt ngưỡng cảnh báo tự động).
- Ưu tiên nguồn **gốc, không ai chặn được**: CBTT HOSE/HNX/UPCoM, website IR doanh nghiệp, công bố NHNN/GSO/Bộ Tài chính. Đây cũng là nguồn *chất lượng cao nhất và không ai tóm tắt*.
- Ưu tiên RSS chính thức trước, crawl HTML chỉ là phương án dự phòng.
- Chế độ suy giảm mềm: nếu 1 nguồn chết, publish 6 tin thay vì 10 và **nói rõ** — không bao giờ lấp chỗ trống bằng tin kém chất lượng.
- Trung hạn: tiếp cận ký **thoả thuận nội dung chính thức** với 2–3 tờ báo (biến rủi ro thành hào phòng thủ — xem 8.3).

### 8.2 Chi phí LLM — `Cao × TB`

**Rủi ro:** ~300–400 tin thô/ngày × 3 chu kỳ cron. Nếu chạy model mạnh trên toàn bộ tin thô, chi phí có thể vượt xa doanh thu giai đoạn đầu. Chi phí tăng tuyến tính theo số nguồn và theo tần suất — nhưng doanh thu MVP gần bằng 0.

**Giảm thiểu:**
- **Kiến trúc phân tầng (quan trọng nhất):** Tầng 1 — lọc bằng heuristic/embedding rẻ (từ khoá mã, khớp danh mục mã, điểm liên quan) loại bỏ ~70% tin thô **trước khi** gọi LLM. Tầng 2 — model nhỏ/rẻ cho phân loại & gắn mã. Tầng 3 — model mạnh **chỉ** cho ~10–15 tin lọt vào digest và cho Research Note.
- Cache theo hash nội dung — cùng một sự kiện được 5 báo đăng chỉ tóm tắt **một lần**.
- Batch API nếu độ trễ cho phép (cron 8h/lần vốn không yêu cầu realtime → tận dụng được giá batch).
- **Đưa "chi phí LLM/tin publish" vào KPI tree (mục 6.3) và theo dõi hằng ngày**, có cảnh báo khi vượt ngưỡng.
- Mô hình đơn vị kinh tế: nếu chi phí/tin publish < 1.500 đ và 1 user tiêu thụ ~300 tin/tháng nhưng dùng chung tin với tất cả user khác → **chi phí biên/user gần 0**. Đây là điểm mấu chốt: chi phí là *cố định theo số tin*, không theo số user. Càng nhiều user, đơn vị kinh tế càng tốt.

### 8.3 Bản quyền & pháp lý nội dung — `TB × Tồn vong` ⚠️

**Rủi ro:** TenPoint tạo tác phẩm phái sinh từ nội dung báo chí có bản quyền, ở quy mô lớn, cho mục đích thương mại. Việt Nam không có học thuyết "fair use" rộng như Mỹ. Nếu CafeF/VCCorp hoặc một tờ báo lớn khởi kiện hoặc gửi yêu cầu gỡ bỏ, và đặc biệt nếu đang bán API B2B (PA3) — đây là rủi ro **có thể chấm dứt doanh nghiệp**, không chỉ gây thiệt hại.

**Giảm thiểu (nhiều lớp, phải làm từ MVP):**
1. **Không bao giờ copy nguyên văn.** Tóm tắt phải là **dữ kiện đã được tái cấu trúc**, không phải câu văn bị rút gọn. Quan sát quan trọng: `[ẢNH]` cho thấy đầu ra đã đúng hướng này — nó là *bảng dữ kiện định lượng*, không phải đoạn văn của bài gốc. **Dữ kiện (số liệu, sự kiện) không được bảo hộ bản quyền; cách diễn đạt mới được bảo hộ.** Đây là lập luận pháp lý mạnh nhất của TenPoint và phải được bảo vệ bằng thiết kế prompt.
2. **Ghi nguồn + link trên 100% dòng** `[ẢNH – cột Source]` — đã là thiết kế sẵn, hãy giữ tuyệt đối, không có ngoại lệ.
3. Giới hạn cứng độ dài tóm tắt (ví dụ ≤120 từ) và **không bao giờ hiển thị ảnh của bài gốc**.
4. Tôn trọng `robots.txt`; dùng RSS chính thức khi có; đặt User-Agent định danh và có trang liên hệ.
5. **Quy trình gỡ bỏ trong 24h** công khai + đầu mối liên hệ rõ ràng.
6. Tham vấn luật sư sở hữu trí tuệ **trước khi** ra mắt PA3 (API B2B) — đây là thời điểm rủi ro nhảy vọt.
7. **Chuyển hoá thành cơ hội:** đề xuất chia sẻ doanh thu hoặc trả phí cấp phép cho 2–3 tờ báo lớn. Chi phí này vừa là bảo hiểm pháp lý vừa là **rào cản gia nhập cho đối thủ sau** — TenPoint có giấy phép, người sao chép thì không.

### 8.4 Thay đổi HTML nguồn (crawler vỡ) — `Cao × TB`

**Rủi ro:** Báo VN đổi layout không báo trước. Crawler hỏng âm thầm → cron vẫn chạy, digest vẫn publish, nhưng thiếu tin hoặc parse sai nội dung. Nguy hiểm nhất là **hỏng âm thầm**, không phải hỏng ồn ào.

**Giảm thiểu:**
- **Kiểm tra dị thường theo nguồn:** nếu một nguồn trả về số bài lệch >50% so với trung bình 7 ngày → cảnh báo ngay, không đợi ai phát hiện.
- Kiểm tra chất lượng sau parse: độ dài nội dung, có ngày tháng, có ≥1 dữ kiện số → nếu không đạt, đánh dấu "cần xem lại" thay vì publish.
- Parser theo tầng: RSS → JSON-LD/schema.org nhúng trong trang → chọn lọc CSS → **LLM trích xuất làm lớp cuối** (chậm/đắt hơn nhưng bền với thay đổi HTML).
- Bộ test crawler chạy hằng ngày với trang mẫu đã lưu (golden file).
- Nguyên tắc vận hành: **thà publish 6 tin đúng còn hơn 10 tin có 2 tin hỏng.**

### 8.5 Chính các tờ báo tự làm tính năng này — `TB × Cao`

**Rủi ro:** CafeF/VnExpress thêm "Tóm tắt bằng AI" vào mỗi bài (chi phí gần bằng 0 với họ, và họ có nội dung gốc). Hoặc DNSE Senses mở rộng ra ngoài hệ sinh thái. Hoặc một CTCK lớn mua đứt một startup tương tự.

**Giảm thiểu — hào phòng thủ phải xây, theo thứ tự ưu tiên:**
1. **Tổng hợp đa nguồn** — CafeF sẽ không bao giờ tóm tắt bài của VnExpress và dẫn link sang đó. Giá trị "một nơi thấy hết" là thứ không tờ báo nào tái tạo được, vì nó đi ngược lợi ích cạnh tranh của họ.
2. **Innovator's dilemma** — CafeF sống bằng pageview; sản phẩm "10 tin, đọc 3 phút, rồi thôi" trực tiếp phá huỷ doanh thu của họ. Họ *có thể* làm nhưng *không nên* làm, và tổ chức lớn hiếm khi tự phá mô hình doanh thu đang chạy.
3. **Đồ thị watchlist + dữ liệu hành vi** — càng nhiều người dùng đặt watchlist và click, xếp hạng nội dung càng tốt. Đây là tài sản tích luỹ, đối thủ khởi động lại từ 0.
4. **Kho lưu trữ có cấu trúc + authority SEO** — 2 năm tin đã gắn mã là tài sản không thể mua bằng tiền trong ngắn hạn.
5. **Tốc độ.** Đây là lợi thế thực tế nhất của một đội nhỏ. Phải đạt product-market fit **trước khi** các tổ chức lớn kịp chú ý.

### 8.6 Rủi ro pháp lý tư vấn đầu tư — `TB × Cao` ⚠️

**Rủi ro:** Research Note `[ẢNH – màn A]` có thể bị diễn giải là hoạt động **tư vấn đầu tư chứng khoán** — hoạt động có điều kiện, cần giấy phép của UBCKNN theo Luật Chứng khoán. Nếu người dùng thua lỗ và khiếu nại, rủi ro pháp lý và uy tín đều lớn.

**Giảm thiểu:**
- **Tuyệt đối không** đưa khuyến nghị MUA/BÁN/NẮM GIỮ, không đưa giá mục tiêu, không đưa tỷ trọng danh mục. Quan sát quan trọng: `[ẢNH – màn A]` **đã tự nhiên tuân thủ** điều này — nó nêu luận điểm có điều kiện ("Biên an toàn chỉ thực sự xuất hiện nếu…") thay vì kết luận hành động. **Giữ nguyên giọng điệu này như một ràng buộc thiết kế bắt buộc, không phải lựa chọn phong cách.**
- Disclaimer rõ ràng, đặt ở vị trí nhìn thấy được (không giấu ở footer).
- Định vị trong mọi tài liệu marketing là **"thông tin và dữ liệu"**, không phải "tư vấn".
- Cấm hợp đồng: đối tác white-label (PA5) **không được chỉnh sửa hoặc chèn khuyến nghị** vào nội dung TenPoint.
- Tham vấn luật sư chứng khoán trước khi ra mắt Research Note tự động ở V2.

### 8.7 Sai sót do LLM (ảo giác số liệu) — `TB × Tồn vong` ⚠️

**Rủi ro:** Một con số sai trong tóm tắt — ví dụ ghi "lãi 101 tỷ" thay vì "lỗ 101 tỷ", hoặc sai một chữ số ở "408.000 tỷ" — có thể khiến người dùng ra quyết định sai. **Trong lĩnh vực tài chính, niềm tin mất một lần là mất vĩnh viễn**, và một ảnh chụp màn hình lan truyền trên mạng xã hội có thể xoá sổ thương hiệu trong 48 giờ.

**Giảm thiểu (đây là rủi ro cần đầu tư kỹ thuật nhiều nhất):**
- **Kiểm tra căn cứ (grounding check) tự động:** mọi con số trong tóm tắt phải **khớp chuỗi (string-match) với con số xuất hiện trong bài gốc**. Con số nào không khớp → chặn publish dòng đó. Đây là biện pháp hiệu quả nhất và rẻ nhất, **bắt buộc có từ MVP**.
- Kiểm tra dấu/chiều: các từ tăng/giảm, lãi/lỗ, mua ròng/bán ròng phải đối chiếu được với bài gốc.
- Kiểm tra tính hợp lý: cảnh báo khi giá trị vượt ngưỡng vô lý (ví dụ doanh thu > 500.000 tỷ, tăng trưởng > 1.000%).
- Con người review **100% Research Note** ở MVP và V2 (đây là nội dung rủi ro nhất, tần suất thấp nên chi phí review chấp nhận được).
- Nút "Báo lỗi" trên mỗi dòng digest + cam kết sửa trong 2 giờ.
- **Short link là lưới an toàn cuối cùng** `[ẢNH]` — người dùng luôn kiểm chứng được. Đây là lý do nữa để không bao giờ bỏ cột Source.

### 8.8 Bảng tổng hợp rủi ro

| # | Rủi ro | Xác suất | Tác động | Mức ưu tiên | Chủ động từ giai đoạn |
|---|---|---|---|---|---|
| 8.3 | Bản quyền & pháp lý nội dung | TB | **Tồn vong** | **P0** | MVP (thiết kế) + V2 (luật sư) |
| 8.7 | Ảo giác số liệu LLM | TB | **Tồn vong** | **P0** | MVP (grounding check) |
| 8.1 | Phụ thuộc nguồn tin | Cao | Cao | **P0** | MVP |
| 8.2 | Chi phí LLM | Cao | TB | P1 | MVP (kiến trúc phân tầng) |
| 8.4 | Crawler vỡ do đổi HTML | Cao | TB | P1 | MVP (cảnh báo dị thường) |
| 8.5 | Báo/CTCK tự làm | TB | Cao | P1 | V2 (xây hào) |
| 8.6 | Pháp lý tư vấn đầu tư | TB | Cao | P1 | V2 (trước Research Note tự động) |

---

## 9. 10 ý tưởng tính năng "wow"

> Tất cả đều là **hệ quả tự nhiên của quyết định kiến trúc "tin = record có khoá ticker"** `[ẢNH – mục 2.1(b)]`. Khách hàng mô tả một cronjob tóm tắt tin; thứ họ thực sự đang xây là **cơ sở dữ liệu sự kiện thị trường có cấu trúc**. Mười tính năng dưới đây là những thứ chỉ cơ sở dữ liệu đó mới làm được.

---

**W1 — Dòng thời gian tin theo mã (Ticker Timeline)**
Trang `/ma/FPT` hiển thị toàn bộ lịch sử tin của mã dưới dạng dòng thời gian dọc, mỗi mốc là một tóm tắt giàu số liệu, có thể lọc theo loại tin và khoảng thời gian. Người dùng đọc được "câu chuyện của FPT trong 6 tháng qua" trong 2 phút — thứ mà hiện nay phải Google 20 lần mới ghép lại được.
`Effort: S` · `Impact: Cao` · *Gần như miễn phí vì dữ liệu đã có sẵn cấu trúc; đồng thời là tài sản SEO số 1 (mục 5.2).*

**W2 — Cảnh báo watchlist tức thời (Zalo / Telegram)**
Khi một tin mới được gắn mã thuộc watchlist của người dùng, đẩy ngay tóm tắt 1 dòng + short link qua Zalo OA hoặc Telegram. Đây là tính năng biến TenPoint từ "trang web tôi thỉnh thoảng ghé" thành "dịch vụ tôi không thể bỏ".
`Effort: S` · `Impact: Cao` · *Đánh thẳng vào chỗ mạnh nhất của đối thủ room Zalo/Telegram (1.9) nhưng với nội dung có nguồn gốc.*

**W3 — Cùng một sự kiện, ba tờ báo nói gì**
Gom các bài từ nhiều nguồn về cùng một sự kiện thành một cụm, hiển thị "CafeF nhấn mạnh X — VnExpress nhấn mạnh Y — Vietstock bổ sung Z", kèm các con số bị lệch giữa các nguồn. Không tờ báo nào có thể làm tính năng này vì nó đòi hỏi trích dẫn đối thủ trực tiếp của họ.
`Effort: M` · `Impact: Cao` · *Hào phòng thủ cấu trúc mạnh nhất (8.5); bản nhẹ "N nguồn cùng đưa tin này" làm được ngay ở MVP.*

**W4 — Điểm sentiment & thang tác động**
Mỗi tin được chấm 2 chiều: sắc thái (tích cực/trung tính/tiêu cực đối với từng mã được gắn) và mức độ tác động ước tính (cao/trung bình/thấp). Hiển thị dạng chấm màu nhỏ ở cột Mã CK, cho phép lọc "chỉ xem tin tác động cao".
`Effort: M` · `Impact: TB` · *Lưu ý: một tin có thể tích cực với mã này và tiêu cực với mã khác (ví dụ CBAM tiêu cực với HPG nhưng trung tính với TIS) — phải chấm theo cặp (tin, mã), không chấm theo tin.*

**W5 — "Tin này ảnh hưởng gì đến giá?"**
Với mỗi loại tin, đối chiếu kho lưu trữ lịch sử: "Trong 18 lần mã này có tin thuộc nhóm *phát hành riêng lẻ*, giá T+3 giảm trung bình 4,2%, 13/18 lần giảm." Không phải dự báo — là **thống kê lịch sử có ngữ cảnh**, tránh hoàn toàn rủi ro tư vấn đầu tư (8.6).
`Effort: L` · `Impact: Cao` · *Tính năng "wow" nhất và khó sao chép nhất — cần ≥12 tháng dữ liệu tích luỹ, nên là phần thưởng cho việc bắt đầu sớm. V3.*

**W6 — Bảng số liệu tự sinh từ tin**
Trích mọi dữ kiện định lượng trong tóm tắt thành bảng chuẩn hoá `{chỉ tiêu, giá trị, đơn vị, kỳ, so sánh YoY}` và tự động so với số liệu cùng chỉ tiêu ở các tin trước. Ví dụ: chi phí lãi vay CMG +48% được đặt cạnh cùng chỉ tiêu của 3 quý trước.
`Effort: M` · `Impact: TB` · *Tận dụng việc số liệu đã được trích và in đậm sẵn `[ẢNH]`; đồng thời chính là hạ tầng cho grounding check ở 8.7 — làm một lần, dùng hai việc.*

**W7 — Digest chênh lệch ("Từ lần bạn đọc gần nhất")**
Ghi nhớ mốc đọc cuối của mỗi người dùng và chỉ hiển thị tin phát sinh sau đó, kèm dòng đếm "12 tin mới kể từ 7h02 sáng nay, 3 tin chạm danh mục của bạn". Người đi công tác 3 ngày về vẫn bắt kịp trong 4 phút thay vì bỏ cuộc.
`Effort: S` · `Impact: Cao` · *Khớp trực tiếp với nhịp cronjob 8h/lần trong ý tưởng gốc; là tính năng chống churn rẻ nhất và hiệu quả nhất.*

**W8 — Bản tin trước phiên ATO (7h00) & tổng kết phiên (15h30)**
Hai bản tin cố định theo nhịp giao dịch VN: bản 7h00 chỉ gồm tin phát sinh sau giờ đóng cửa hôm trước + lịch sự kiện doanh nghiệp trong ngày; bản 15h30 tổng kết phiên + tin trong giờ. Bám đúng nhịp sinh hoạt của NĐT Việt Nam.
`Effort: S` · `Impact: Cao` · *Chính là tính năng chiếm thói quen buổi sáng — tức là chiếm NSM (mục 6.1). Bắt buộc có ở MVP.*

**W9 — Luận điểm sống (Living Thesis)**
Research Note `[ẢNH – màn A]` không "chết" sau khi publish: mỗi luận điểm được gắn với các mã và chủ đề liên quan, và khi có tin mới, hệ thống tự gắn cờ "Tin ngày 14/10 **củng cố** luận điểm 3" hoặc "**mâu thuẫn** với luận điểm 1". Người dùng thấy luận điểm đầu tư của mình đang đúng hay sai theo thời gian thực.
`Effort: L` · `Impact: Cao` · *Không tổ chức research nào ở VN làm được vì họ không có dòng tin có cấu trúc; V3.*

**W10 — Radar ngành (Sector Heat)**
Đếm mật độ tin và sắc thái trung bình theo ngành trong 7/30 ngày, hiển thị dạng bản đồ nhiệt: "Ngành thép: 23 tin trong 7 ngày (+180% so với trung bình), sắc thái đang xấu đi". Phát hiện sớm ngành nào "đang được thị trường kể chuyện" trước khi giá phản ánh hết.
`Effort: M` · `Impact: TB` · *Rất dễ chia sẻ trên mạng xã hội → kênh tăng trưởng lan truyền tự nhiên; là ứng viên số 1 cho nội dung marketing hằng tuần.*

---

### 9.1 Ma trận ưu tiên & khuyến nghị đưa vào MVP

| | **Impact Cao** | **Impact TB** |
|---|---|---|
| **Effort S** | **W1** Ticker Timeline · **W2** Alert watchlist · **W7** Diff Digest · **W8** ATO/Tổng kết phiên | — |
| **Effort M** | **W3** Ba tờ báo nói gì | W4 Sentiment · W6 Bảng số liệu · W10 Sector Heat |
| **Effort L** | W5 Ảnh hưởng đến giá · W9 Living Thesis | — |

**Top 5 đưa vào MVP (4 tính năng effort S + 1 bản nhẹ của W3):**

| # | Tính năng | Vì sao phải có ngay ở MVP |
|---|---|---|
| 1 | **W8 — Bản tin 7h00 & 15h30** | Không có nhịp cố định thì không có thói quen; không có thói quen thì NSM = 0. Đây là tính năng *định nghĩa* sản phẩm. |
| 2 | **W2 — Cảnh báo watchlist** | Biến sản phẩm từ "trang web" thành "dịch vụ"; là lý do người dùng đăng ký và để lại liên hệ. |
| 3 | **W1 — Ticker Timeline** | Chi phí gần bằng 0 vì dữ liệu đã có cấu trúc; đồng thời là toàn bộ nền móng SEO (mục 5.2). Hai giá trị, một lần làm. |
| 4 | **W7 — Diff Digest** | Tính năng chống churn rẻ nhất; trực tiếp bảo vệ tiêu chí thoát MVP "quay lại D7 ≥40%". |
| 5 | **W3 (bản nhẹ) — "N nguồn cùng đưa tin này"** | Chỉ cần hiển thị số nguồn và danh sách link đã đủ tạo cảm giác "đây là nơi thấy hết"; nâng cấp thành so sánh đầy đủ ở V2. |

**Để lại V2:** W4, W6, W10 — có giá trị nhưng không ảnh hưởng tiêu chí thoát MVP.
**Để lại V3:** W5, W9 — cần dữ liệu lịch sử tích luỹ (≥12 tháng), làm sớm cũng không có nguyên liệu.

---

## 10. Ba câu hỏi cần chốt với khách hàng trước khi build

1. **Persona chính là ai — nhà đầu tư cá nhân hay môi giới?** Quyết định này thay đổi thứ tự ưu tiên mô hình kinh doanh (PA1 vs PA5), thiết kế onboarding, và cả kênh phân phối ưu tiên. Khuyến nghị của BA: **xây cho NĐT cá nhân, bán cho môi giới trước** — vì môi giới trả tiền nhanh hơn nhưng NĐT cá nhân mới tạo được tài sản dài hạn.

2. **"Cronjob 8h/lần" là ràng buộc cứng hay chỉ là mô tả ban đầu?** Tần suất này phù hợp với chi phí và với bản tin sáng, nhưng **mâu thuẫn với W2 (cảnh báo tức thời)** — vốn là tính năng giữ chân mạnh nhất. Đề xuất kiến trúc lai: cron 3 lần/ngày cho digest đầy đủ + kênh realtime nhẹ riêng cho CBTT và tin tác động cao.

3. **Research Note `[ẢNH – màn A]` là sản phẩm chính hay sản phẩm marketing?** Nếu là sản phẩm chính → cần biên tập viên/analyst con người, chi phí vận hành thay đổi hoàn toàn, và rủi ro 8.6 tăng mạnh. Nếu là sản phẩm marketing (5–10 bài/tuần cho các mã nóng, để tạo backlink và uy tín) → giữ ở V2 và chấp nhận quy mô nhỏ. Khuyến nghị của BA: **sản phẩm marketing ở V2, cân nhắc thành sản phẩm chính ở V3 khi đã có W9 (Living Thesis) làm điểm khác biệt thật sự.**

---

## Phụ lục — Danh mục cần verify trước khi ra quyết định

| # | Giả định cần kiểm chứng | Cách verify | Ảnh hưởng tới |
|---|---|---|---|
| 1 | Bảng giá và cơ cấu gói của Simplize, Fireant, Vietstock Finance, Wichart | Truy cập trang giá công khai | PA2 (định giá freemium) |
| 2 | Chính sách hoa hồng affiliate của các CTCK | Liên hệ trực tiếp phòng phát triển đối tác | PA1 (mô hình doanh thu MVP) |
| 3 | Mức phí room VIP và quy mô nhóm Zalo/Telegram | Phỏng vấn 5–10 môi giới | PA5 (định giá white-label) |
| 4 | Số lượng tin/ngày thực tế từ 8 nguồn dự kiến | Chạy crawler thử 7 ngày | 8.2 (mô hình chi phí LLM) |
| 5 | Chi phí LLM thực tế/tin sau khi lọc tầng 1 | Prototype pipeline, đo trên 1.000 tin | 8.2, tiêu chí thoát MVP #7 |
| 6 | Độ khó crawl từng nguồn (robots.txt, Cloudflare, paywall) | Kiểm tra kỹ thuật từng domain | 8.1, 8.4 |
| 7 | Quy định & chi phí Zalo OA/ZNS cho nội dung tài chính | Tài liệu Zalo for Business | 5.4 (kênh phân phối) |
| 8 | Ranh giới pháp lý "tư vấn đầu tư" theo Luật Chứng khoán | Tham vấn luật sư chứng khoán | 8.6, thiết kế Research Note |
| 9 | Số mã niêm yết/ĐKGD thực tế & khối lượng CBTT hằng ngày | HOSE/HNX/VSD | 5.2 (quy mô SEO), 8.2 |
| 10 | Traffic thực tế của CafeF/Vietstock/Fireant/Simplize | SimilarWeb / Ahrefs | Mục 1 (quy mô đối thủ), 5.1 (tiềm năng SEO) |

---

*Tài liệu này dựa trên ý tưởng gốc của khách hàng và bằng chứng từ 2 reference image. Mọi mục gắn `[GIẢ ĐỊNH]` cần được kiểm chứng theo phụ lục trên trước khi dùng làm cơ sở ra quyết định đầu tư.*
