# TenPoint — Báo cáo nghiệm thu QA

> Thực hiện: QA Engineer · Ngày: 16/09/2026
> Đối chiếu: `docs/05-qa/test-plan.md`, `docs/02-pm/PRD.md` mục 10, `docs/03-sa/ingest-edge-cases.md`, `docs/04-design/design-system.md`, `docs/00-inputs/reference-images.md`
> Môi trường: macOS darwin 25.4.0 · Go 1.26.2 · Node v25.9.0 · Docker 29.4.0 / Compose v5.1.2
> Toàn bộ kết quả dưới đây là output thật từ lệnh đã chạy. Case không chạy được ghi `BLOCKED`, không suy đoán.

**Tài liệu này có 3 vòng.** [Vòng 3](#vong-3) ở đầu là kết quả nghiệm thu hiện hành. [Vòng 2](#vong-2) và [Vòng 1](#vong-1) giữ nguyên bên dưới để đối chiếu.

---

<a id="vong-3"></a>

# ══════ VÒNG 3 — NGHIỆM THU CUỐI ══════

> Ngày 16/09/2026 · Môi trường: **compose PRODUCTION** (`deploy/docker-compose.yml`, không dev override)

## V3.0 ✅ KẾT LUẬN: ĐẠT — ĐỀ NGHỊ PHÁT HÀNH MVP v1

Toàn bộ điều kiện chặn phát hành đã được gỡ. Tôi tự đo lại độc lập, không dùng số BE báo.

| Hạng mục | Ngưỡng | Vòng 1 | Vòng 2 | **Vòng 3 (tôi tự đo)** | |
|---|---|---|---|---|---|
| Tin publish có số bịa | **0** | 0 | 0 | **0** (76 số/6 bài ngẫu nhiên) | ✅ |
| Tóm tắt đúng bài | 100% | 81% | 100% | **100%** | ✅ |
| **Precision tag mã** | **≥96%** | ~52% | 80,5–83,2% | **100,00%** (96/96, 0 FP) | ✅ |
| Ngôn ngữ khuyến nghị đầu tư | 0 | 1/20 | 0 | **0** bài published | ✅ |
| Cụm bold đúng 1–3 | | 1 tin có 0 | 3 tin có 0 | **0 vi phạm** (`{1: 100}`) | ✅ |
| Độ dài tóm tắt 60–150 | | đạt | đạt | **đạt** | ✅ |
| 4 service healthy (production) | 4/4 | **3/4** | 4/4 | **4/4** | ✅ |
| 3/3 trang chạy | | **1/3** | 3/3 | **3/3** | ✅ |
| WCAG 2.2 AA | đạt | 8/10 | 8/10 | **10/10** | ✅ |
| Khớp 2 ảnh tham chiếu | 1:1 | Ảnh A hỏng | đạt | **đạt, không hồi quy** | ✅ |
| Open redirect | 0 | 0 | 0 | **0** | ✅ |
| Pipeline 1 run | <10ph | 2m53s | 2m55s | **2m54s** | ✅ |
| Bảng mobile scroll ngang | không | không | không | **không** | ✅ |
| Console error | 0 | 0 | 0 | **0** | ✅ |

**Không còn defect P0. Không còn defect P1 chặn phát hành.**

---

## V3.1 Phân xử `#108 MWG` → **`expected` / `mentioned`**

BE làm đúng khi không tự quyết mà chuyển cho QA. Phán quyết:

**Thân bài (trích nguyên văn):**
> *"Các khoản đầu tư trong lĩnh vực bán lẻ cũng diễn biến tích cực, dẫn đầu là FPT Retail (FRT) tăng 18,3% và Digiworld (DGW) tăng 16,3%, trong khi **Mobile World Group (MWG)** cũng tăng giá."*

**Lý do xếp `expected`:**
1. Nêu **đích danh tên công ty kèm mã trong ngoặc** — `Mobile World Group (MWG)`.
2. **Giống hệt dạng của TCB/MBB/VPB trong CHÍNH câu liền trước** cùng bài: *"Techcombank (TCB) tăng 16,4%, MB Bank (MBB) tăng 13,9% và VPBank (VPB) tăng 13,1%"* — ba mã này tôi đã xếp `expected` từ vòng 2. Xếp MWG là FP sẽ **mâu thuẫn nội bộ ngay trong một đoạn văn**.
3. Cùng dạng với BID/MBB/TCB/VCB ở bài 191 mà tôi đã xếp `expected`.
4. `MWG` **có trong từ điển seed** (`CTCP Đầu tư Thế Giới Di Động`).

**Mức `mentioned` chứ không `primary`:** MWG chỉ là một khoản nắm giữ trong danh sách và là mã duy nhất **không có số % riêng** (*"cũng tăng giá"*). Bài này đã có tiền lệ chấm hỗn hợp (ACB `mentioned`).

**Vì sao vòng 2 thiếu nhãn này:** tôi chỉ adjudicate các tag mà tagger sinh ra; lúc đó tagger chưa sinh MWG. Đúng như `recall_caveat` tôi đã ghi sẵn trong `_meta`. Đây là bổ sung recall, không phải sửa sai.

**Ghi chú recall (đã ghi vào `note` của mẫu):** cùng câu còn nêu `FRT` và `DGW`, nhưng **cả hai KHÔNG có trong từ điển seed** (0 lần trong `0003_seed_tickers.sql`) nên tagger không thể gắn. Đó là thiếu độ phủ từ điển, không phải lỗi tagger → **không** đưa vào `expected` để không làm sai mẫu số của cổng.

### Corpus đã cập nhật (tôi là người duy nhất sửa)

```
md5 TRƯỚC khi tôi sửa : 73b202fbee958c40cbe3d0dfbb41ff0e   ← khớp md5 BE báo (73b202fb…)
md5 SAU khi tôi sửa   : a7b2550f3a00705c46b5e94967f70ca2
```
Việc md5 trước khi sửa khớp đúng con BE báo **chứng minh BE không hề đụng vào file**.

| | trước | sau |
|---|---|---|
| samples | 27 | **27** |
| expected | 95 | **96** |
| must_not_tag | 42 | **42** |
| uncertain | 3 | **3** |

> ⚠️ **BE phải cập nhật hằng số chốt cứng** trong `TestTaggerPrecisionOnQACorpus`: `ExpectedTags` **95 → 96**. Hiện 4 test fail đúng như thiết kế — xem V3.3. Đã ghi `revision` + `changelog` vào `_meta`.

## V3.2 Precision — TÔI TỰ ĐO, KHÔNG DÙNG SỐ BE BÁO

Tôi **không** chạy test của BE. Tôi viết scorer riêng theo đúng `scoring_rule` trong `_meta` của chính mình, và **nạp từ điển mã xuất trực tiếp từ DB đang chạy** (57 mã, 11 mã có `negative_aliases`) thay vì fixture của BE — vì vòng 3 BE tự thừa nhận cổng cũ chạy trên từ điển viết tay thiếu `negative_aliases`.

```
=== QA ĐỘC LẬP — precision trên corpus QA (từ điển lấy từ DB đang chạy, 57 mã) ===
  TP=96  FP=0  uncertain-skipped=3
  PRECISION = 96/96 = 100.00%   (ngưỡng PRD >= 96%)
  primary  : TP=65 FP=0 -> 100.00%
  mentioned: TP=31 FP=0 -> 100.00%
  relevance grade: khớp=92 lệch=4
  FALSE POSITIVES: KHÔNG CÓ
  vi phạm trên bài nhãn RỖNG: 0 []
  recall (THAM KHẢO, corpus không vét cạn) = 96/96 = 100.00%
```

**100,00% — 0 false positive.** Cao hơn con số 98,96% BE báo, vì BE tính `#108 MWG` là FP theo đúng luật của tôi, còn tôi đã phân xử nó là `expected`. **Cả hai cách tính đều vượt ngưỡng 96%.**

*Ghi nhận trung thực:* `relevance grade: lệch=4` — có 4 tag đúng mã nhưng mức `primary`/`mentioned` khác nhãn của tôi. Theo `scoring_rule` tôi tự viết, precision chỉ tính **mã**, không tính mức, nên không ảnh hưởng kết quả. Ghi lại để minh bạch.

### Tầng T3 — tự đếm trên DB, không tin báo cáo

```
score|relevance|count
0.80 |mentioned| 23      0.90|mentioned| 82      0.95|primary| 19
0.80 |primary  | 26      0.90|primary  | 98
```
**Tầng `score=0.40` biến mất hoàn toàn: 0 bản ghi** (vòng 2: 38 tag / 12 bài). Trần cứng 86,7% của vòng 2 không còn.

### 4 kiểm chứng riêng — tôi tự chạy lại, 4/4 PASS

| Bài | Kỳ vọng | Kết quả |
|---|---|---|
| `#108` | hết `FPT` và `GAS` | ✅ (nằm trong 0 FP) |
| `#129` | hết `GAS` | ✅ |
| `#168` | hết `VND` | ✅ |
| `#191` | **giữ đủ** `BID, MBB, TCB, VCB`, chỉ bỏ 4 mã dầu khí T3 | ✅ |
| 10 bài `verified_empty` | 0 tag | ✅ `vi phạm trên bài nhãn RỖNG: 0` |

## V3.3 Cổng tự động — sạch, trừ 4 test fail CÓ CHỦ Ý

```
go build ./...  → sạch      go vet ./...  → sạch      gofmt -l .  → sạch
npx tsc --noEmit → sạch     npm run lint  → sạch      npm run build → ✓ 7/7 pages
```

`go test ./...` → **4 test FAIL trong package `tagger`**, tất cả từ **cùng một dòng**:

```
corpus_qa_test.go:66   require.Equal(t, 95, c.Meta.Counts.ExpectedTags)
                       expected: 95   actual: 96
```

Đây là **cơ chế chống sửa trộm hoạt động đúng như thiết kế** — tôi thêm 1 nhãn thì hằng số vỡ. 4 test (`TestTaggerPrecisionOnQACorpus`, `TestQACitedFalsePositivesAreFixed`, `TestArticle191KeepsCorrectHighScoreTags`, `TestVerifiedEmptySamplesGetNoTags`) đều gọi chung helper `loadQACorpus`. **Không phải fail về precision.** Sửa 1 số `95 → 96` là xanh.

## V3.4 Compose PRODUCTION — 4/4 healthy

```
$ docker compose -f deploy/docker-compose.yml up -d --build
 Container tenpoint-api Healthy        ← vòng 1 không bao giờ tới được
 Container tenpoint-web Started
UP_EXIT=0

NAME             STATUS
tenpoint-api     Up (healthy)      tenpoint-db     Up (healthy)
tenpoint-caddy   Up (healthy)      tenpoint-web    Up (healthy)
```

Cả 3 trang `200` với backend thật: `/`, `/nhan-dinh`, `/ma/FPT`, `/ma/CMG`, `/ma/PNJ`, `/nhan-dinh/cmg-quy-2-2026`.

> Ghi chú không phải lỗi: qua Caddy, `/healthz` trả `404` vì Caddy **không expose endpoint health ra ngoài** (chỉ route `/api/v1/*`). Kiểm trong container: `HEAD /healthz → 200`, `HEAD /readyz → 200`, `wget_exit=0`. Đây là lựa chọn đúng về bảo mật.

## V3.5 Hai lỗi a11y — cả hai PASS

### TC-A11Y-03 Contrast ✅ PASS (đo computed style thật)

| | Light | Dark |
|---|---|---|
| `--ink-muted` | `#666f79` → **4,66:1** ✅ (vòng 2: `#6B7580` = 4,28 ❌) | `#7e8894` → **5,18:1** ✅ |
| `.digest-date` / `.digest-type` 13.5px | 4,66 ✅ | 5,18 ✅ |
| `.digest-summary` | 9,35 ✅ | 9,71 ✅ |
| `.source-link` | 9,35 ✅ | 9,71 ✅ |
| `.ticker-cell` | 16,31 ✅ | 15,71 ✅ |

**0/5 mẫu fail AA ở cả hai chế độ.** Cảm ơn anh đã tính lại và đính chính tài liệu — con số 4,28 của tôi được xác nhận.

### TC-A11Y-08 Target size ✅ PASS — kiểm bằng hit-test thật, không chỉ đọc CSS

Tôi không chỉ đọc CSS mà dùng `document.elementFromPoint()` bắn vào 4 cạnh của vùng mở rộng:

| | Hộp thị giác | `::after` inset | **Vùng chạm thật** | ≥24px | 4 cạnh có nhận click |
|---|---|---|---|---|---|
| `.source-link` | 54,3 × 17 | `-4px -3px` | **60,3 × 25** | ✅ | ✅ 4/4 |
| `.ticker-link` | 25,9 × 18 | `-4px -3px` | **31,9 × 26** | ✅ | ✅ 4/4 |

`position: relative` trên phần tử gốc → pseudo neo đúng. Bắn ra **ngoài** vùng mở rộng thì rơi xuống `TD` → vùng chạm có biên thật, không nuốt cả ô. **55 link kiểm tra, 0 link dưới 24×24.**

> *Tự đính chính:* lần đo đầu tôi kiểm `::before` (CSS thật dùng `::after`) nên báo "không có pseudo"; và 4 probe của `.ticker-link` trả `false` chỉ vì phần tử **nằm ngoài viewport** — `elementFromPoint` trả null cho toạ độ ngoài màn hình. Sau khi `scrollIntoView`, cả 5 probe đều trúng `A.ticker-link`. Cả hai đều là lỗi của tôi, không phải của bản vá.

### Không hồi quy thị giác

`padding` ô `20px / 24px 0` (đúng `--space-5`/`--space-6`), bảng rộng `1132px` trong container `1180px`, chiều cao hàng ~214px. Kiểm cấu trúc: 5 cột đúng thứ tự · `caption.sr-only` · hairline `#D5D8D2` · kẻ header `#14181B` · **không** border dọc · **không** zebra · `radii` chỉ `2px ×14` (chip) · **0 shadow** · **0 gradient** · `strong` chỉ `weight:700`, không nền, không gạch chân. **0 sai lệch token ở CẢ hai chế độ.**

## V3.6 Cổng `off_topic` — hiệu quả, nhưng loại oan một ít (P1 KHÔNG chặn)

**Kết quả tốt:** trang chủ nay **20/20 là tin tài chính**. Ba tin "kênh đào Trung Quốc / lên men cà phê / tuỳ bút Nhà xuất bản Trẻ" của vòng 2 **đã hết**. Ba tin đầu hiện tại: *thị giá DCL tăng trần lên 41.650 đồng* · *Công ty Quản lý Quỹ Nhân Việt* · *lợi suất trái phiếu chính phủ Hy Lạp*.

**Kiểm hồi quy theo yêu cầu — soi 35 bài Tuổi Trẻ bị loại.** Tôi tự tải lại RSS `tuoitre.vn/rss/kinh-doanh.rss`, đối chiếu `pubDate` với thời điểm crawl để chỉ xét bài **thực sự được đưa qua cổng**, rồi tự phân loại:

| Nhóm | Số bài | |
|---|---|---|
| **A. Loại SAI rõ ràng** — tin chứng khoán đích thực | **2** (6%) | ❌ |
| B. Ranh giới — tin doanh nghiệp/hàng hoá, tranh luận được | 5 (14%) | ⚠️ |
| C. Loại ĐÚNG | 28 (80%) | ✅ |

Hai bài loại sai (đã xác minh `0 rows` trong DB, và `pubDate` trước giờ crawl nên chắc chắn đã qua cổng):

```
✗ "Chân dung những quỹ ngoại trăm tỉ USD bắt đầu tính đến cổ phiếu Việt Nam"   16/09 06:00
✗ "Khối ngoại bỏ nghìn tỉ 'gom hàng' trước ngày chứng khoán nâng hạng"          15/09 19:52
```

Cả hai đúng chủ đề `Khối ngoại` — **một trong 9 nhãn taxonomy**.

**Vì sao KHÔNG chặn phát hành:** đây là lỗi **recall**, không phải lỗi đúng/sai. Không có thông tin sai nào được hiển thị. 124 bài tài chính vẫn publish, cả 9 nhãn đều có dữ liệu (`khoi_ngoai` = 15 bài). PRD mục 10 nói rõ **ưu tiên precision hơn recall**. Đây là ngưỡng tinh chỉnh được. → **P1 không chặn, xử lý ở v1.1.**

## V3.7 Kiểm chứng lớp lỗi bỏ dấu — sửa tận gốc, có bằng chứng dữ liệu thật

| Kiểm | Kết quả |
|---|---|
| Bài published chứa ngôn ngữ tư vấn | **0** |
| Bài bị cách ly | **3** (vòng 2: 10) |
| Lý do cách ly `nen tang` / `chung ta` / `can ban` | **0** (vòng 2: 8/10 là oan) |
| Bài chứa `nền tảng` / `căn bản` | **published**, không bị cách ly |

3 bài còn bị cách ly đều **đúng là tin tư vấn**:
```
#8  Nhận định thị trường phiên giao dịch ngày 17/9: Thận trọng ngắn hạn
#37 Dòng tiền ngoại mua vào nâng hạng sẽ "đỡ" cho thông tin…
#74 "VnEconomy giới thiệu nhận định và khuyến nghị đầu tư của một số công ty chứng khoán…"
```

Bằng chứng trực quan nhất cho bản vá: bài **"CEO Kafi: 'Sau giai đoạn xây nền tảng, Kafi bước vào bài toán tăng trưởng'"** nay **published** — đúng cụm `nền tảng` từng làm vòng 2 cách ly oan 5 bài.

Migration `0007` **không thả nhầm**: 0 bài published chứa ngôn ngữ tư vấn.

## V3.8 Các case FAIL/BLOCKED còn lại

| ID | Vòng 2 | **Vòng 3** | Bằng chứng |
|---|---|---|---|
| TC-VAL-09 | FAIL | ✅ **PASS** | Phân bố cụm bold trên 100 tin: **`{1: 100}`**, 0 vi phạm (vòng 2: `{0:10, 1:90}`) |
| TC-A11Y-03 | FAIL | ✅ **PASS** | 4,66:1 light · 5,18:1 dark |
| TC-A11Y-08 | FAIL | ✅ **PASS** | 55 link, 0 dưới 24px; hit-test 4/4 cạnh |
| TC-PIPE-03/04 | PASS | ✅ **PASS** | 3 cụm đa bài, mỗi cụm đúng 1 canonical |
| **P2-8 (vòng 1)** | mở | ✅ **ĐÃ SỬA** | Cột `truncated` nay ghi `t` cho 3 nguồn chạm trần 40 (vòng 1&2 đều `f`) |
| TC-PIPE-01 | BLOCKED | **BLOCKED** | Không nguồn nào trả 5xx; `error_message` rỗng cả 7 nguồn |
| TC-PIPE-08 | BLOCKED | **BLOCKED** | Không nguồn nào trả 0 bài |
| TC-PIPE-09 | BLOCKED | **BLOCKED** | `skipped_by_robots = 0` toàn bộ |
| PRD LCP 4G | BLOCKED | **BLOCKED** | Chưa chạy Lighthouse throttle 4G |

5 case BLOCKED đều **không phải lỗi sản phẩm** — là tình huống không xảy ra tự nhiên trong run thật. Không chặn phát hành.

## V3.9 Regression sweep — không có hồi quy

| Kiểm | Kết quả |
|---|---|
| Short link `/r/{code}` → 302 → đích | **3/3 dest HTTP 200** |
| Tin demo | **410** · Mã không tồn tại: **404** |
| Số liệu bịa (6 bài ngẫu nhiên, 76 số) | **0** |
| Console error/warning | **0** |
| Ảnh A (Research Note) | LUẬN ĐIỂM ĐẦU TƯ ✓ · đủ 4 luận điểm ✓ · disclaimer ✓ · JSON-LD `NewsArticle` ✓ |
| Ảnh B (Digest) | 5 cột · hairline · không zebra/border dọc/shadow/gradient · bold-by-weight ✓ |
| Dark + Light | **0 sai lệch token** cả hai chế độ, cả hai trang |
| `loai=abc` → 400 · `limit=999` → 100 · `/api/v1/meta` → 9 nhãn | ✅ |

## V3.10 Tồn đọng cho v1.1 (KHÔNG chặn phát hành)

Theo chỉ đạo, liệt kê để đưa vào v1.1 chứ không dùng để chặn:

| ID | Mức | Mô tả |
|---|---|---|
| **P1-R3-1** | P1 | Cổng `off_topic` loại oan ~6% (2/35 mẫu Tuổi Trẻ) — mất tin `Khối ngoại` hợp lệ. Nới ngưỡng hoặc whitelist từ khoá taxonomy |
| P2-R3-2 | P2 | Độ phủ từ điển mã: 57/~1.600 mã. 11/20 hàng trang chủ không có mã; `FRT`, `DGW`, `DCL`, `BSR` xuất hiện trong bài nhưng không có trong từ điển |
| P2-R3-3 | P2 | `relevance grade` lệch 4/96 (đúng mã, sai mức `primary`/`mentioned`) |
| P2-R3-4 | P2 | Chất lượng cụm bold: vẫn còn bold năm trần (`3/2009`, `3/2012`) và cắt đơn vị (`18,7 triệu` bỏ `đồng`). NUM-3 đã đạt nhưng B-3 chưa hoàn hảo |
| P2 (vòng 2) | P2 | `duplicate_summary` chỉ so tuyệt đối + chỉ trong cùng nguồn; chưa có unit test |
| P2 (vòng 2) | P2 | `is_stale` API 26h vs cảnh báo FE 12h — hai ngưỡng khác nhau, dễ gây hiểu nhầm |
| P2 (vòng 2) | P2 | Timeout lẻ tẻ khi fetch `tinnhanhchungkhoan.vn` |
| P1 (vòng 1) | P1 | P1-1 fixture render link chết · P1-2 host `robots.txt` build-time · P1-5 bypass rate limit bằng XFF |
| DOC | P2 | Corpus cũ `tagger_corpus.json` có trường `id` không ánh xạ đúng `articles.id` |

## V3.11 Bằng chứng & dọn dẹp

Ảnh trong `docs/05-qa/evidence/`: `qa-r3-home-light.png`, `qa-r3-home-dark.png` *(vòng 2: `qa-r2-*`; vòng 1: `qa-reference-compare.png`…)*

**Không sửa một dòng code sản phẩm nào.** Chỉ: tạo `deploy/.env` từ `.env.example`; sửa **1 nhãn** trong corpus QA (quyền của QA); 1 file probe test tạm trong package `tagger` **đã xoá ngay sau khi chạy**. Docker đã `down -v`.

---

<a id="vong-2"></a>

# ══════ VÒNG 2 — Nghiệm thu lại sau khi team sửa 5 P0 ══════

> Ngày: 16/09/2026 (buổi tối) · Môi trường: **compose PRODUCTION** (`deploy/docker-compose.yml`, KHÔNG dùng dev override)

## V2.0 Kết luận vòng 2

### ❌ VẪN CHƯA ĐẠT — nhưng tiến bộ rất lớn

**Cả 5 lỗi P0 của vòng 1 đều đã được sửa và tôi đã kiểm chứng độc lập từng lỗi.** Đây là kết quả thật, không phải nhận xét xã giao.

| P0 vòng 1 | Trạng thái vòng 2 | Bằng chứng tôi tự đo |
|---|---|---|
| **P0-1** 2/3 trang HTTP 500 | ✅ **ĐÃ SỬA** | 7/7 URL trả `200` với backend thật; API tắt thì rơi fixture êm, không 500 |
| **P0-2** 36 tin VnEconomy chung 1 tóm tắt | ✅ **ĐÃ SỬA** | `vneconomy.vn: 36 bài / 36 tóm tắt phân biệt` (vòng 1: 36/1); text hội thảo ô-dôn còn **0** bản ghi |
| **P0-3** Precision tag ≈52% | ⚠️ **CẢI THIỆN NHIỀU nhưng CHƯA ĐẠT** | Đo độc lập: **80,5–83,2%** (vòng 1: ~52%; ngưỡng PRD 96%). `primary` **96,8–98,4% ✅ đạt**; `mentioned` 60–64% ❌ |
| **P0-4** Lọt ngôn ngữ khuyến nghị | ✅ **ĐÃ SỬA** | 15/15 payload vòng 1 nay bị chặn; bài "Nhận định 17/9" nay `pending_review`; **0** bài published còn ngôn ngữ tư vấn |
| **P0-5** Healthcheck không bao giờ pass | ✅ **ĐÃ SỬA** | Compose **production**: `4/4 healthy`, `web` khởi động được |

**Chưa thể nghiệm thu, vì 3 lý do:**

1. **P0-3 chưa đóng được.** Precision đo độc lập **80,5–83,2%**, ngưỡng PRD **≥96%**. Quan trọng hơn con số mẫu: tầng T3 (`score=0.40`) có **38 tag trên 12 bài, và cả 38 đều sai — đây là phép ĐẾM TOÀN BỘ, không phải lấy mẫu**. Riêng nó đã tạo **trần cứng 86,7%** cho toàn quần thể, thấp hơn ngưỡng bất kể sai số thống kê.
2. **P1-NEW-1 (mới, do chính bản vá P0-4 gây ra): bộ lọc tư vấn bắt nhầm vì va chạm khi bỏ dấu.** `nền tảng` (platform) → `nen tang` trùng `nên tăng`; `căn bản` (fundamental) → `can ban` trùng `cần bán`. **6/10 bài bị cách ly là oan.**
3. **P1-NEW-2 (mới, lộ ra khi dữ liệu sạch): 3 tin đầu trang chủ không phải tin chứng khoán** — kênh đào Trung Quốc, lên men cà phê, tuỳ bút về Nhà xuất bản Trẻ. Tất cả gắn nhãn `Thị trường`, không có mã nào.

Cộng thêm các P1 vòng 1 chưa sửa (ngoài phạm vi vòng này nhưng vẫn mở): P1-1 fixture link chết, P1-2 host robots.txt, P1-3 contrast 4.28:1, P1-5 bypass rate limit bằng XFF.

> **Một lớp lỗi duy nhất gây ra 2 defect ở 2 nơi khác nhau — đáng chú ý cho team:** cả bộ lọc tư vấn (`nền tảng`→`nên tăng`) lẫn tagger (alias GAS `Khí Việt Nam`→`Khi Việt Nam`) đều **bỏ dấu tiếng Việt trước khi so khớp**. Dấu trong tiếng Việt là phân biệt nghĩa. Nên rà soát **mọi** chỗ dùng `Normalize`/`NormalizeWords` để so khớp cụm từ, không chỉ 2 chỗ đã tìm thấy.

### Bảng ngưỡng PRD mục 10 — trạng thái vòng 2

| Chỉ tiêu | Ngưỡng | Vòng 1 | **Vòng 2** | |
|---|---|---|---|---|
| Tin publish có số bịa | **0** | 0 | **0** (75 số trên 6 bài vneconomy đã sửa, 0 sai) | ✅ |
| Tóm tắt đúng bài | 100% | 81% | **100%** (7/7 nguồn 1:1) | ✅ |
| Precision tag mã | ≥96% | ~52% | **80,5–83,2%** | ❌ |
| — riêng `primary` | — | — | **96,8–98,4%** | ✅ |
| Recall tag mã | ≥85% | không đo được | **không đo được** (vẫn chưa có tập gán nhãn) | ⚠️ |
| Tóm tắt 60–150 từ | | 20/20 | **20/20** (80–113 từ) | ✅ |
| Cụm bold đúng 1–3 | | 1 tin có 0 | **3 tin có 0** | ❌ |
| 0 cụm khuyến nghị đầu tư | 0 | 1/20 vi phạm | **0 tin published vi phạm** | ✅ |
| Pipeline 1 run | <10 phút | 2m53s | **2m55s** | ✅ |
| p95 `/api/v1/news` | <200ms | 2–9ms | **2–9ms** | ✅ |
| Short link redirect | <50ms | 0–2ms | **0–2ms** | ✅ |
| Open redirect | 0 | 0 | **0** | ✅ |
| Bảng mobile scroll ngang | không | không | **không** | ✅ |
| 4 service healthy | 4/4 | 3/4 | **4/4 (production compose)** | ✅ |
| Khớp ảnh tham chiếu | 1:1 | Ảnh A không render được | **Ảnh A ✅ + Ảnh B ✅** | ✅ |
| WCAG 2.2 AA | đạt | 8/10 | **8/10** (contrast + target size vẫn FAIL) | ❌ |
| LCP mobile 4G | <2.0s | không đo | **không đo** | ⚠️ |

---

## V2.1 Cổng tự động — sạch toàn bộ

```
$ cd backend
$ go build ./...          → (sạch)
$ go vet ./...            → (sạch)
$ gofmt -l .              → (không file nào)
$ go test ./... -count=1  → 10/10 package ok
$ go test ./... -race -count=1 → không data race

$ cd frontend
$ npx tsc --noEmit  → (sạch)
$ npm run lint      → (sạch)
$ npm run build     → ✓ Compiled successfully; 7/7 static pages
```

> Ghi nhận tích cực về quy trình: team báo rằng ngay khi bỏ ép kiểu `as T`, **`tsc` tự bắt được một chỗ lệch thứ hai mà vòng 1 tôi không phát hiện** (`ResearchNoteSummary.id` bắt buộc trong khi BE định danh bằng `slug`). Đây đúng là điều khuyến nghị "cổng tự động không bắt được P0 nào" ở vòng 1 nhắm tới — cổng đã bắt đầu làm được việc của nó.

## V2.2 P0-5 — compose PRODUCTION, 4/4 healthy

Chạy đúng như yêu cầu, **không dùng dev override**:

```
$ docker compose -f deploy/docker-compose.yml up -d --build
 Container tenpoint-db Healthy
 Container tenpoint-api Starting
 Container tenpoint-api Waiting
 Container tenpoint-api Healthy      ← vòng 1 KHÔNG BAO GIỜ tới được dòng này
 Container tenpoint-web Starting
 Container tenpoint-web Started      ← vòng 1 web không bao giờ khởi động
 Container tenpoint-caddy Started
UP_EXIT=0
```

```
NAME             STATUS                        PORTS
tenpoint-api     Up About a minute (healthy)   8080/tcp
tenpoint-caddy   Up 55 seconds (healthy)       0.0.0.0:80->80/tcp, 443->443
tenpoint-db      Up About a minute (healthy)   5432/tcp
tenpoint-web     Up 56 seconds (healthy)       3000/tcp
```

Probe chính xác chỗ vòng 1 hỏng:

```
$ docker exec tenpoint-api sh -c 'wget -q --spider http://127.0.0.1:8080/healthz; echo $?'
wget_exit=0                       (vòng 1: wget_exit=8)
HEAD /healthz  →  HTTP/1.1 200 OK (vòng 1: 405 Method Not Allowed)
HEAD /readyz   →  HTTP/1.1 200 OK
health=healthy failing_streak=0   (vòng 1: unhealthy, streak=6)
```

✅ **`make up` trên máy sạch nay chạy được.** TC-DEPLOY-02 chuyển FAIL → **PASS**.

## V2.3 P0-1 — 3/3 trang hoạt động với backend thật

Tất cả qua Caddy `http://127.0.0.1` (production compose không expose api/web trực tiếp).

| URL | HTTP | Nguồn dữ liệu |
|---|---|---|
| `/` | **200** | backend thật (không có banner fixture) |
| `/nhan-dinh` | **200** | backend thật |
| `/ma/FPT` | **200** | backend thật |
| `/ma/CMG` | **200** | backend thật |
| `/ma/HPG` | **200** | backend thật |
| `/ma/VIC` | **200** | backend thật |
| `/nhan-dinh/cmg-quy-2-2026` | **200** | backend thật |

Kiểm nội dung thật (không phải fixture):

```
/ma/FPT      title: "FPT · CTCP FPT · TenPoint"   company rendered: True
/nhan-dinh/{slug}:
  LUẬN ĐIỂM ĐẦU TƯ present: True
  đủ 4 luận điểm (2.323 tỷ đồng / 24.150 / 250 triệu USD): True
  disclaimer present: True
  JSON-LD NewsArticle: True
```

**Suy giảm êm khi tắt `api`** (yêu cầu riêng của PM):

```
$ docker stop tenpoint-api
/                          200  (fixture banner shown)
/nhan-dinh                 200  (fixture banner shown)
/ma/FPT                    200  (fixture banner shown)
/ma/CMG                    200  (fixture banner shown)
/nhan-dinh/cmg-quy-2-2026  200  (fixture banner shown)
/ma/KHONGTONTAI            404  (đúng — không bịa trang)
/nhan-dinh/khong-ton-tai   404  (đúng)
```

✅ Không còn `500` ở bất kỳ đường nào. Adapter có kiểm chứng runtime hoạt động đúng như team mô tả.

## V2.4 P0-2 — VnEconomy đã sửa thật, không chỉ "khác nhau"

Tôi kiểm 2 tầng, vì "36 tóm tắt phân biệt" chưa chứng minh là "36 tóm tắt ĐÚNG".

**Tầng 1 — tính phân biệt:**

```sql
select s.domain, count(*) articles, count(distinct a.summary_md) distinct_summaries ...
```
```
tuoitre.vn            | 40 | 40
cafef.vn              | 40 | 40
tinnhanhchungkhoan.vn | 39 | 39
vneconomy.vn          | 36 | 36   ← vòng 1: 36 | 1
vnexpress.net         | 28 | 28   ← regression check: KHÔNG bị selector xoá sạch
vietstock.vn          |  4 |  4
thanhnien.vn          |  2 |  2
```

**Tầng 2 — tóm tắt có khớp đúng bài không:**

```
id 10 | Từng vỡ nợ, Hy Lạp bây giờ đi vay…  → "Đây là một sự thay đổi đáng kể… khi Hy Lạp…"     ✓
id 14 | Khối ngoại mua ròng trước thềm Fed  → "Khối ngoại có phiên mua ròng thứ 3 liên tiếp…"   ✓
id 17 | Vận hành thị trường tài sản mã hóa  → "Việt Nam mong muốn tham khảo kinh nghiệm FMA…"   ✓
id 21 | Blog chứng khoán: Dòng tiền nín thở → "…thanh khoản khớp lệnh hai sàn giảm còn 12,6k tỷ" ✓
id 27 | Cổ phiếu vừa và nhỏ điều chỉnh mạnh → "…VIX giảm 1,94% khớp 317,2 t…"                    ✓
id 31 | Quỹ ngoại: Nâng hạng là chất xúc tác→ "…theo LVF, thị trường cổ phiếu Việt Nam…"         ✓
```

Text hội thảo tầng ô-dôn: `select count(*) ... like '%tầng ô-dôn%' or '%Kigali%'` → **0**.

**Truy vết số liệu trên chính nguồn từng hỏng** (6 bài vneconomy, fetch bài gốc, đối chiếu từng số):

```
  id 10   nums_in_summary=8    UNVERIFIED=none
  id 14   nums_in_summary=10   UNVERIFIED=none
  id 17   nums_in_summary=3    UNVERIFIED=none
  id 21   nums_in_summary=21   UNVERIFIED=none
  id 27   nums_in_summary=24   UNVERIFIED=none
  id 31   nums_in_summary=9    UNVERIFIED=none
TOTAL unverified numbers across sample: 0
```

✅ **75 con số, 0 số bịa.** Chỉ tiêu tồn vong của PRD tiếp tục đạt sau khi sửa.

**Hàng rào mới (safety net)** — `pipeline.go:482-487` gọi `SummaryAlreadyUsed` → `pending_review` reason `duplicate_summary`. Ghi nhận trung thực: **hàng rào này KHÔNG kích hoạt lần nào trong run vòng 2**, vì nguyên nhân gốc đã được sửa nên không còn trùng để bắt. Tức là nó chưa được kiểm chứng bằng dữ liệu thật. Xem P2-NEW-1.

## V2.4b P0-3 — Precision tag mã: ĐO ĐỘC LẬP, **không tái lập được con số 100%**

> PM yêu cầu đúng trọng tâm: corpus `backend/testdata/tagger_corpus.json` do chính BE gán nhãn, tức là **người viết code tự chấm bài mình**. Dưới đây là phép đo **độc lập**, mẫu rút ngẫu nhiên từ **DB thật**, không dùng corpus làm chuẩn.

### Phương pháp
- Quần thể: 179 bài thật đã publish; **67 bài có ≥1 mã** (tổng **286 tag**).
- **Loại trừ trùng corpus:** khớp theo tiêu đề chuẩn hoá (phát hiện phụ: trường `id` trong corpus **không** ánh xạ đúng sang `id` của DB — 39/50 trỏ nhầm bài). 47/50 mục corpus khớp 48 bài DB → **loại 37/67 bài có mã**, còn **30 bài đủ điều kiện**.
- Rút ngẫu nhiên **20 bài** (seed cố định) trong 30 bài đó → **113 tag**.
- Over-sample vneconomy theo yêu cầu: **13/20 bài** là vneconomy (tự nhiên đạt, vì corpus đã loại hết vneconomy nên nguồn này chiếm 17/30 bài đủ điều kiện).
- Mỗi bài: resolve `/r/{code}` → fetch bài gốc (`--compressed`) → **bóc riêng thân bài** theo selector từng nguồn, **loại bỏ sidebar/related/tag/share** → đối chiếu từng mã và alias.
- **20/20 trang fetch được HTTP 200. 0 tag BLOCKED**, không tag nào bị loại khỏi mẫu số.

### Kết quả

| Chỉ số | Giá trị |
|---|---|
| Tag đã chấm | **113** |
| Đúng | **91** · Sai: **19** · Tranh cãi: **3** |
| **Precision (tính 3 tranh cãi là sai)** | **80,5%** |
| **Precision (tính 3 tranh cãi là đúng)** | **83,2%** |
| KTC 95% bootstrap theo cụm bài | **64,0% – 94,2%** |
| Ngưỡng PRD | **≥ 96%** → ❌ **KHÔNG ĐẠT** |
| BE công bố (corpus của chính BE) | 100% → **không tái lập được** |

### Tách theo mức relevance — đây là chỗ mấu chốt

| relevance | n | đúng | sai | precision | |
|---|---|---|---|---|---|
| **`primary`** | 63 | 61 | 1 | **96,8 – 98,4%** | ✅ **ĐẠT ngưỡng PRD** |
| **`mentioned`** | 50 | 30 | 18 | **60,0 – 64,0%** | ❌ hỏng nặng |

**Tin tốt thật sự: mức `primary` — thứ hiển thị mặc định trên bảng digest — đã đạt chuẩn.** Toàn bộ sai sót dồn vào `mentioned`.

### Tách theo điểm tin cậy của tagger — nguyên nhân gốc lộ rõ

| score | n | đúng | sai | precision |
|---|---|---|---|---|
| 0.95 | 5 | 5 | 0 | 100% |
| 0.90 | 74 | 73 | 1 | 98,6% |
| 0.80 | 18 | 13 | 2 | 72–89% |
| **0.40** | **16** | **0** | **16** | **0,0%** |

**Mọi tag ở tầng 0.40 trong mẫu đều sai. Không có ngoại lệ.**

### Tầng T3 chưa được sửa — chỉ bị thu hẹp. Đây là phép ĐẾM TOÀN BỘ, không phải lấy mẫu

Tầng T3 nhận diện được chính xác trong DB: `score = 0.400`, luôn `mentioned`, luôn là một nhóm ngành cố định. Tôi tự xác minh lại bằng query:

```
score|relevance|count          articles carrying 0.40 tags
0.40 |mentioned| 38     ←      32|CMG,FPT|Tập đoàn Trung Quốc muốn đầu tư năng lượng…
0.80 |mentioned| 22            62|HPG,HSG,NKG,TIS|Canada và thế khó của 'nền kinh tế hướng Mỹ'
0.80 |primary  | 25            70|HPG,HSG,NKG,TIS|Giá vật liệu xi măng và thép tăng…
0.90 |mentioned| 83            82|POW,REE|EVN hết lỗ lũy kế, lãi hơn 12.200 tỷ…
0.90 |primary  | 98           105|CMG,FPT|Bắc Ninh hỗ trợ 100% tiền thuê đất…
0.95 |primary  | 20           118|POW,REE|Xây dựng cơ chế cảnh báo sớm cho ngành xuất khẩu
                              121|BAF,DBC,HAG,MML|Việt Nam xây dựng thương hiệu bò H'Mông…
                              130|GAS,PLX,PVD,PVS|Lãnh đạo ngành dầu khí Mỹ cảnh báo…
                              159|CMG,FPT|Trình diễn công nghệ in kỹ thuật số…
                              162|ACB,BID,CTG,HDB|Tỷ giá nửa cuối năm…
                              168|ACB,BID,CTG,HDB|'Hồi hộp' Fed trước giờ quyết định lãi suất
                              191|GAS,PLX,PVD,PVS|Chứng khoán 16-9: Tiền đổ vào dầu khí…
```

**38 tag / 12 bài / 13,3% tổng số tag.** Cả 12 bài đã được fetch và kiểm — **38/38 tag đều không xuất hiện trong thân bài**. Ví dụ hiển nhiên: *"Việt Nam xây dựng thương hiệu **bò H'Mông** để cạnh tranh với Wagyu"* gắn `BAF, DBC, HAG, MML`; *"**Canada** và thế khó của nền kinh tế hướng Mỹ"* gắn `HPG, HSG, NKG, TIS`; *"**Trình diễn công nghệ in** kỹ thuật số"* gắn `CMG, FPT`.

> **Tôi tự kiểm chứng lại bài 191** (*"Chứng khoán 16-9: Tiền đổ vào dầu khí, ngân hàng"* — trường hợp trông hợp lệ nhất vì tiêu đề có "dầu khí"): grep **phân biệt hoa thường, có biên từ** trên toàn bộ HTML → `GAS: 0, PLX: 0, PVD: 0, PVS: 0`. Lần grep đầu của tôi báo 7 lần là **sai của chính tôi** — nó khớp `Megastory` và chuỗi base64 do dùng `-i` và không có biên từ. Kết luận của phép đo độc lập đứng vững.

**Hệ quả:** vì 38 tag này đã được xác minh là sai bằng phép đếm toàn bộ (không phải suy ra từ mẫu), tồn tại **trần cứng cho toàn quần thể: precision ≤ 248/286 = 86,7%** — trước khi tính các lỗi khác ở tầng 0.80–0.90. **Trần này thấp hơn ngưỡng 96% bất kể sai số lấy mẫu.**

### Lỗi mới phát hiện ở tầng điểm CAO (không phải T3)

1. **Va chạm alias khi bỏ dấu** — cùng một lớp lỗi với P1-NEW-1 ở bộ lọc tư vấn: alias `"Khí Việt Nam"` của GAS khớp cụm thường `"Khi Việt Nam"` → bài 108 gắn **GAS `primary` 0.80**, bài 129 gắn GAS `mentioned` 0.80. Câu thật: *"**Khi Việt Nam** chuẩn bị chính thức gia nhập câu lạc bộ các thị trường mới nổi FTSE…"*.
2. **`negative_aliases` không được áp ở tầng 0.90** — bài 108 gắn **FPT 0.90** trong khi thân bài chỉ có *"**FPT Retail (FRT)** tăng 18,3%"*. `FPT Retail` **đã nằm trong `negative_aliases` của FPT trong chính DB**, và bài còn ghi rõ mã đúng là FRT. Đúng lớp lỗi `Điện Máy Xanh→MWG` của vòng 1, sống sót ở tầng điểm cao.
3. **Rò rỉ sidebar vẫn còn** — bài 162 gắn BID và HDB; hai mã này chỉ xuất hiện trong `div.box-sidebar-*` bên trong `div.sub-column` (rail bài liên quan), không có trong thân bài.

### Vì sao BE thấy 100%
Không phải do gian lận — do **corpus chính là tập tuning**: 37 trong 67 bài có mã của DB nằm trong corpus, và corpus đã loại toàn bộ vneconomy. Mẫu của tôi cố ý loại đúng 37 bài đó. Đo trên tập đã dùng để chỉnh thì 100% là kết quả dự kiến, không phải bằng chứng về chất lượng.

### Giới hạn của phép đo này (nêu rõ, không làm tròn có lợi cho sản phẩm)
- n = 20 bài / 113 tag. **Không đủ để phân biệt 94% với 97%.** Nhưng **đủ để phân biệt 83% với 96%** — khoảng cách lớn hơn nhiều sai số lấy mẫu.
- Tag bị gom cụm theo bài (1 bài xấu đóng góp 4–8 tag tương quan) nên tôi báo KTC bootstrap theo cụm **64–94%** (rộng hơn), không dùng Wilson naive 72–89%.
- Ước lượng này áp cho **quần thể 30 bài không thuộc corpus**, không phải toàn bộ 179 bài.
- Số theo nguồn ngoài vneconomy (n=16, 11, 2) **quá nhỏ để xếp hạng nguồn** — đừng đọc "tuoitre 0%" như một sự thật về nguồn.
- 3 tag tranh cãi (mã chỉ xuất hiện trong cụm trích dẫn tổ chức phân tích / công thức lãi suất tham chiếu) được tách riêng và báo dạng khoảng. Theo **chính rule B3 trong corpus của BE** ("KHÔNG gắn khi mã chỉ nằm trong cụm trích dẫn nguồn phân tích") thì ít nhất 1 trong 3 là sai — tức điểm ước lượng sẽ **giảm**, không tăng.

### Chỉ số quần thể liên quan
- Bài **không có mã nào**: **112/179 (62,6%)**
- Bài có ≥4 mã toàn `mentioned` (vân tay T3): 12/179 — 6 là nhóm ngành 0.40, 6 là bản tin thị trường đa mã hợp lệ
- Bài có mã nhưng **không có `primary` nào**: 23/179

## V2.5 P0-4 — bộ lọc tư vấn đã bịt, nhưng sinh lỗi mới

**Đã bịt hoàn toàn 15/15 payload vòng 1** (probe tạm, đã xoá sau khi chạy):

```
BLOCKED  phrase="khuyen nghi"        | Chuyên gia khuyến nghị nắm giữ.
BLOCKED  phrase="khuyen nghi"        | SSI khuyến nghị mua.
BLOCKED  phrase="khuyen nghi"        | VCSC khuyến nghị bán.
BLOCKED  phrase="co hoi dau tu"      | Đây là cơ hội đầu tư.
BLOCKED  phrase="chot loi"           | Nhà đầu tư có thể chốt lời.
BLOCKED  phrase="bat day"            | Cổ phiếu đang bắt đáy.
BLOCKED  phrase="tiem nang tang gia" | Cổ phiếu có tiềm năng tăng giá.
BLOCKED  phrase="canh mua"           | Nhà đầu tư có thể canh mua.
BLOCKED  phrase="khuyen nghi"        | Quan điểm: khuyến nghị nắm giữ.
BLOCKED  phrase="khuyen nghi"        | **Khuyến nghị nắm giữ.**
BLOCKED  phrase="khuyen nghi"        | …nắm giữ;  /  …nắm giữ!  /  …nắm giữ?  /  …nắm giữ)
BLOCKED  phrase="nen duy tri"        | [chính câu đã lọt ở vòng 1] "…nhà đầu tư ngắn hạn nên duy trì
                                       tỷ trọng cổ phiếu ở mức trung bình, tránh mua đuổi…"
ADVICE CASES 1-15: blocked=15 leaked=0     (vòng 1: blocked=1 leaked=14)
```

**Bài từng lọt ở vòng 1 nay đã bị cách ly đúng:**

```
id 7 | pending_review | Nhận định thị trường phiên giao dịch ngày 17/9: Thận trọng t…
```

**Quét toàn bộ tin đang publish:** `published articles still containing advice language: 0` → migration `0006` đã dọn sạch.

**Regression check của PM (#6b) — ✅ PASS:** tin vi phạm nằm ở `pending_review`, **không bị vứt**:

```
status         | reject_reason     | count
pending_review | summary_rejected  | 10
```
(0 bản ghi ở trạng thái `rejected` vì lý do tư vấn.)

### ⚠️ P1-NEW-1 — lỗi MỚI do chính bản vá gây ra

Thống kê cụm nào kích hoạt 10 lần cách ly đó:

```
5  banned_phrase (nen tang)     ← "nền tảng" = platform/nền móng
2  banned_phrase (chung ta)
1  banned_phrase (nen duy tri)  ← đúng (tư vấn thật)
1  banned_phrase (khuyen nghi)  ← đúng (tư vấn thật)
1  banned_phrase (can ban)      ← "căn bản" = fundamental
```

Nguyên nhân: bộ lọc **bỏ dấu trước khi so khớp**, mà dấu trong tiếng Việt là phân biệt nghĩa:

| Từ thường gặp | Bỏ dấu | Va chạm với cụm cấm | Nghĩa thật |
|---|---|---|---|
| **nền tảng** | `nen tang` | `nên tăng` (should increase) | platform / nền móng |
| **căn bản** | `can ban` | `cần bán` (should sell) | fundamental / cơ bản |

Probe xác nhận (đã xoá sau khi chạy):

```
Công ty xây dựng nền tảng số cho doanh nghiệp.      | *** FALSE POSITIVE: nen tang ***
Đây là nền tảng công nghệ lõi của tập đoàn.         | *** FALSE POSITIVE: nen tang ***
Nền tảng thương mại điện tử tăng trưởng 20%.        | *** FALSE POSITIVE: nen tang ***
Về căn bản, kết quả quý II đạt kế hoạch.            | *** FALSE POSITIVE: can ban ***
Phân tích căn bản cho thấy biên lợi nhuận ổn định.  | *** FALSE POSITIVE: can ban ***
TRUE-POSITIVE CHECK "Nhà đầu tư nên tăng tỷ trọng"  | found=true  (vẫn bắt đúng)
TRUE-POSITIVE CHECK "Nhà đầu tư cần bán ra ngay."   | found=true  (vẫn bắt đúng)
```

**Mức độ: P1.** `nền tảng` là từ cực kỳ phổ biến trong tin kinh tế/công nghệ Việt Nam. **6/10 bài bị cách ly ở run này là oan (60%).** Không gây sai số liệu, nhưng làm mất tin thật một cách âm thầm và sẽ tăng theo thời gian.

**Đề xuất:** với các cụm mà việc bỏ dấu tạo va chạm, hãy so khớp trên text **còn dấu**; hoặc thêm danh sách miễn trừ (`nền tảng`, `căn bản`, `nền kinh tế`…). Lưu ý cụm `chung ta`/`tom lai`/`nhin chung` (P1-6 vòng 1) vẫn còn trong danh sách và vẫn chặn văn xuôi bình thường — chưa sửa.

## V2.5b Rubric NUM trên dữ liệu vòng 2 (20 tin)

| Rubric | Ngưỡng | Vòng 1 | Vòng 2 | |
|---|---|---|---|---|
| **NUM-1 Độ dài** | 60–150 từ, reject >180 | 20/20 đạt | **20/20 đạt** (80–113 từ) | ✅ |
| **NUM-3 Số cụm bold** | đúng 1–3 | 1 tin có 0 cụm | **3 tin có 0 cụm** (id 141, 146, 32) | ❌ **xấu đi** |
| **B-3 Chất lượng cụm** | số + đơn vị trọn vẹn | 73,7% vi phạm | **47% vi phạm** (8/17) | ⚠️ khá hơn, chưa đạt |

Các cụm còn vi phạm B-3: `2009`, `2012`, `2026` ×2 (năm trần) · `16` (cắt đôi ngày) · `55.100`, `16.000`, `27.200` (số trần không đơn vị).

> **Ghi chú phương pháp — quan trọng để không so sai:** con số 73,7% ở vòng 1 là tôi chấm **thủ công, nghiêm hơn** (tính cả `52 triệu` là cắt cụm từ `52 triệu đồng`). Con số 47% ở vòng 2 là **script tự động, lỏng hơn** (chấp nhận mọi cụm có chữ số kèm một đơn vị bất kỳ). **Hai con số không so trực tiếp được.** Điều chắc chắn: cả hai vòng đều KHÔNG đạt, và lỗi "bold một năm trần" vẫn còn nguyên.

## V2.6 Kiểm tra hồi quy 2 lỗi BE tự báo — ✅ cả hai PASS

| # | Rủi ro | Kết quả |
|---|---|---|
| 6a | VnExpress bị selector `sidebar` xoá sạch bài | ✅ **28 bài** (vòng 1: 27). Không mất nguồn |
| 6b | Tin có ngôn ngữ tư vấn bị vứt thay vì giữ lại | ✅ 10 bài ở `pending_review`, 0 bài bị xoá |

## V2.7 Chạy lại 8 FAIL + 8 BLOCKED của vòng 1

| ID | Vòng 1 | Vòng 2 | Bằng chứng |
|---|---|---|---|
| TC-VAL-09 | FAIL | **FAIL** (vẫn) | Phân bố cụm bold trên 100 tin: `{0: 10, 1: 90}` → 10 tin có **0 cụm**, vi phạm "đúng 1-3". Vẫn không enforce ở backend |
| TC-VAL-10 | FAIL | ✅ **PASS** | 15/15 payload bị chặn; 0 tin published còn ngôn ngữ tư vấn |
| TC-TAG-17 | FAIL | ✅ **PASS** | `"Nguồn: https://cafef.vn/tin-tuc/FPT-abc-123.html"` → `(none)` |
| TC-UI-10 | FAIL | ✅ **PASS** | Trang note render đủ; xem V2.8 |
| TC-UI-11 | FAIL | ✅ **PASS** | Disclaimer hiện trên trang note: *"không phải khuyến nghị mua/bán chứng khoán…"* |
| TC-A11Y-03 | FAIL | **FAIL** (vẫn) | `--ink-muted #6B7580` trên `#F4F5F2` = **4.28:1** < 4.5. Dùng cho `.digest-date`, `.digest-type` 13.5px |
| TC-A11Y-08 | FAIL | **FAIL** (vẫn) | **57/81** target < 24×24px. *(Cải thiện một phần: chip filter nay `44×30` ✓; link nguồn + chip mã trong bảng vẫn ~17-18px cao)* |
| TC-DEPLOY-02 | FAIL | ✅ **PASS** | 4/4 healthy trên compose production |
| TC-PIPE-01 | BLOCKED | **BLOCKED** (vẫn) | Không nguồn nào trả 5xx; `error_message` rỗng cả 7 nguồn |
| TC-PIPE-03 | BLOCKED | ✅ **PASS** | 3 cụm đa bài hình thành, mỗi cụm đúng 1 canonical: `cluster 47→3 bài/1 canon`, `59→2/1`, `6→2/1` |
| TC-PIPE-04 | BLOCKED | ✅ **PASS** | 185 cluster phân biệt / 184 bài → tiêu đề khác thì khác cụm |
| TC-PIPE-08 | BLOCKED | **BLOCKED** (vẫn) | Cả 7 nguồn `consecutive_empty_runs=0`, không nguồn nào trả 0 bài |
| TC-PIPE-09 | BLOCKED | **BLOCKED** (vẫn) | `skipped_by_robots=0` toàn bộ; không đường dẫn nào bị robots chặn |
| TC-UI-09 | BLOCKED | ✅ **PASS** | Lùi `finished_at` 14h (chỉ sửa dữ liệu test, đã revert) → trang chủ hiện `⚠ Cảnh báo: dữ liệu có thể chưa mới nhất.` |
| TC-SEO-05 | BLOCKED | ✅ **PASS** | JSON-LD `@type: NewsArticle`, đủ `headline/author/publisher/datePublished/dateModified/mainEntityOfPage/about/inLanguage` |
| PRD LCP 4G | BLOCKED | **BLOCKED** (vẫn) | Chưa chạy Lighthouse throttling 4G |

**Kết quả: 8 case chuyển sang PASS, 3 FAIL còn lại, 5 BLOCKED còn lại.**

### Phát hiện phụ ở TC-UI-09
Khi dữ liệu cũ 14h: FE hiện cảnh báo đúng (ngưỡng 12h theo design-system §5.7), **nhưng API vẫn trả `is_stale: false`** (dùng `READY_MAX_AGE_HOURS=26`). Hai ngưỡng khác nhau cho 2 mục đích khác nhau là hợp lý, nhưng một trường tên `is_stale` báo `false` trong khi giao diện đang cảnh báo dữ liệu cũ thì dễ gây hiểu nhầm cho người vận hành. → **P2-NEW-2**.

## V2.8 Đối chiếu giao diện — lần đầu so được Ảnh A

### Ảnh A (Research Note) — ✅ ĐẠT

Vòng 1 không so được vì trang 500. Nay đo computed style thật:

| Yếu tố | Spec design-system §5.5 | Đo được | |
|---|---|---|---|
| Tiêu đề | display, sans 800, 52px desktop, lh 1.12, UPPERCASE | `Be Vietnam Pro` / `800` / `52px` / `58.24px` (=1.12) / `uppercase` | ✅ |
| `LUẬN ĐIỂM ĐẦU TƯ` | h2, uppercase, tracking .06em, 21px, 700 | `21px` / `700` / `uppercase` / `1.26px` (=.06em) | ✅ |
| Body luận điểm | serif 17px / 1.62 | `Source Serif 4` / `17px` / `27.54px` (=1.62) / `400` | ✅ |
| `<ol>` list-style | none, số inline trong dòng dẫn | `none`, 4 mục, dòng dẫn `weight 700` ×4 | ✅ |
| Không card/shadow/bo góc | — | `radii={}`, `shadows={}`, `grads=[]` | ✅ |
| Container | 1180px | `1180` | ✅ |
| Disclaimer cuối bài | bắt buộc | có | ✅ |

Đối chiếu với ảnh gốc: tiêu đề UPPERCASE đậm trải rộng, section heading, 4 luận điểm đánh số với dòng dẫn bold rồi đoạn body, 1 cột, không viền — **khớp**. Nội dung số liệu khớp nguyên văn ảnh (2.323 tỷ đồng, 17,7%→17,9%, 101 tỷ, 75 tỷ, 48%, 12%, 25%, 250 triệu USD, 30 MW, 100 MW, 24.150 đồng/cp, P/E 14, ROE 12%).

*Khác biệt có chủ đích, đúng spec:* body dùng serif (design-system §3 chọn 2 họ chữ cho 2 màn hình) trong khi ảnh gốc là sans; nền `#F4F5F2` thay vì `#F7F5F0` của ảnh (design-system v2 §0 cố ý đổi).

### Ảnh B (Digest) — ✅ ĐẠT (không hồi quy)

5 cột đúng thứ tự · hairline `#D5D8D2` · kẻ đậm header `#14181B` · không zebra · không border dọc · `radii` chỉ `2px` (chip) · `shadows={}` · `grads=[]` · số đậm chỉ bằng weight · Source = domain gạch chân · `Loại tin` text thuần.

### Dark mode — ✅ ĐẠT cả 2 trang

0 sai lệch token trên cả trang chủ lẫn trang note. Contrast trong dark: tiêu đề `15.71:1`, body `9.71:1` (AAA).

### Console — ✅ 0 lỗi, 0 cảnh báo

### ⚠️ P1-NEW-2 — 3 tin đầu trang chủ KHÔNG phải tin chứng khoán

Đây là thứ đập vào mắt khách hàng đầu tiên. 12 tiêu đề đầu của feed mặc định:

```
[Thị trường        ] tickers=-   Trung Quốc mở kênh đào dài 134km ở Quảng Tây, nối sông với Biển Đông
[Thị trường        ] tickers=-   Lên men cà phê: Ở đâu, khi nào và vì sao?
[Thị trường        ] tickers=-   Thương hiệu tôi yêu: Những ấn phẩm của Nhà xuất bản Trẻ
[Pháp lý           ] tickers=-   Một công ty quản lý quỹ bị xử phạt do vi phạm công bố thông tin
[Thị trường        ] tickers=-   Ông chủ Công ty Triệu Nụ Cười thu gần 34 tỉ từ chiêu quảng bá…
[Vĩ mô             ] tickers=-   Mới có hai địa phương hỗ trợ người dân lắp điện mặt trời mái nhà
[Trái phiếu/Tín dụng] tickers=-  Từng vỡ nợ, Hy Lạp bây giờ đi vay với lãi suất còn thấp hơn cả Pháp
[Cổ tức/Phát hành  ] tickers=PNJ Hai quỹ Dragon Capital mua lại gần 750.000 cổ phiếu PNJ…
```

Ba tin đầu — **kênh đào Trung Quốc, lên men cà phê, tuỳ bút về Nhà xuất bản Trẻ** — đều gắn nhãn `Thị trường` và không có mã nào.

Nguyên nhân (suy ra từ dữ liệu, cần BE xác nhận):
1. Pipeline nạp nguyên feed RSS tổng hợp của báo phổ thông (tuoitre.vn đóng góp 40 bài) mà **không có tầng lọc chủ đề tài chính**.
2. Classifier dùng `Thị trường` như nhãn mặc định khi không khớp nhóm nào.
3. Bài **không có mã CK nào vẫn được publish và đẩy lên đầu feed** vì sắp xếp thuần theo `published_at`.

Mức độ: **P1**, nhưng xét theo giá trị sản phẩm thì đây là vấn đề nghiêm trọng nhất còn lại. PRD mục 1 hứa *"biến hàng trăm bài báo chứng khoán mỗi ngày thành ~10 tin đã chưng cất"*; hiện trang chủ ghi "179 tin mới đáng chú ý" và mở đầu bằng bài về cà phê.

*Ghi nhận trung thực:* vấn đề này đã tồn tại ở vòng 1 nhưng bị che bởi 2 lỗi khác (36 tóm tắt hỏng của VnEconomy và tình trạng gắn mã tràn lan). Vòng 1 tôi chỉ ghi nhận "40% hàng không có mã" mà chưa nhận ra đây là vấn đề **độ liên quan chủ đề**. Dữ liệu sạch làm nó lộ rõ.

*Chỉ số liên quan:* tỉ lệ hàng không có mã CK **tăng** từ 8/20 (40%) lên **11/20 (55%)** sau khi siết T3 — đúng như đánh đổi precision/recall dự kiến.

## V2.9 Các P1 vòng 1 chưa sửa (ngoài phạm vi vòng này, vẫn mở)

| ID | Trạng thái | Bằng chứng vòng 2 |
|---|---|---|
| **P1-1** fixture render link chết | **VẪN MỞ** | Tắt api → 5 tin mẫu vẫn là `<a href="/r/…">`, `Dữ liệu mẫu` xuất hiện **0** lần. Phân giải: `a7Kx2p→410`, `b3Qm9d→404`, `c8Tz1k→404`, `d2Wr7v→404`, `e6Yn4s→404` (4/5 chết) |
| **P1-2** host build-time vs runtime | **VẪN MỞ** | `robots.txt: Host: http://localhost:3000` nhưng `sitemap <loc>http://localhost/</loc>` và `canonical http://localhost` |
| **P1-3** contrast 4.28:1 | **VẪN MỞ** | đo lại: `date 4.28`, `type 4.28` |
| **P1-5** bypass rate limit bằng XFF | **VẪN MỞ** | Qua Caddy (đường production): không XFF `16/70` bị 429; XFF xoay vòng **`0/70`** |
| **P1-6** chặn nhầm văn xuôi | **VẪN MỞ**, nay rõ hơn | `chung ta` ×2 lần cách ly ở run này; hợp nhất với P1-NEW-1 |
| P1-4, P1-7 | Chưa kiểm lại vòng này | — |

### Ghi nhận tiến bộ ở P1-7 (negative_aliases)
Probe cho thấy đã sửa được phần lớn:

```
"CTCP Đầu tư Điện Máy Xanh (DMX - HOSE)…"  → DMX/primary          ✅ (vòng 1: MWG sai)
"Quỹ VanEck Vectors Vietnam ETF (VNM ETF)…" → (none)              ✅ (vòng 1: VNM=Vinamilk sai)
"CMG trúng thầu…"                          → CMG/primary          ✅ (vòng 1: CMG + FPT thừa)
"Ngành thép… CBAM"                         → HPG,HSG,NKG,TIS mentioned ✅ (G-08 vẫn đúng)
"FPT Telecom báo lãi quý II…"              → FPT/primary FOX/primary  ⚠️ còn thừa FPT
"Masan Consumer chia cổ tức…"              → MCH/primary MSN/primary  ⚠️ còn thừa MSN
```
2 trường hợp "công ty con có tên chứa tên công ty mẹ" vẫn gắn cả mã mẹ. → hạ xuống **P2-NEW-3**.
*(Lưu ý phương pháp: probe này tôi tự nạp thêm FOX/MCH/DMX vào rổ mã; rổ thật trong DB có 56 mã, nếu thiếu FOX/MCH thì kết quả thực tế sẽ khác.)*

## V2.10 Defect mới của vòng 2

| ID | Mức | Mô tả | Vị trí |
|---|---|---|---|
| **P1-NEW-1** | **P1** | Bộ lọc tư vấn bắt nhầm do va chạm khi bỏ dấu: `nền tảng`→`nen tang` trùng `nên tăng`; `căn bản`→`can ban` trùng `cần bán`. 6/10 bài bị cách ly là oan | `summarizer/advice.go` + biên từ trong `validator.go` |
| **P1-NEW-2** | **P1** | Feed không lọc theo chủ đề tài chính: 3 tin đầu trang chủ là kênh đào TQ / lên men cà phê / tuỳ bút sách, đều nhãn `Thị trường`, không mã | tầng discovery + classifier |
| **P2-NEW-1** | P2 | Hàng rào `duplicate_summary` **không có unit test** và chưa từng kích hoạt bằng dữ liệu thật; chỉ so khớp **tuyệt đối** và chỉ **trong cùng một nguồn** (`source_id = $1`) → gần-trùng hoặc trùng chéo nguồn vẫn lọt | `store/revisions.go:246-257` |
| **P2-NEW-2** | P2 | `is_stale` của API (ngưỡng 26h) mâu thuẫn với cảnh báo hiển thị của FE (ngưỡng 12h): dữ liệu 14h → FE cảnh báo nhưng API báo `is_stale:false` | `handlers_meta.go` / `READY_MAX_AGE_HOURS` |
| **P2-NEW-3** | P2 | Công ty con có tên chứa tên công ty mẹ vẫn gắn kèm mã mẹ: `FPT Telecom`→FPT+FOX, `Masan Consumer`→MSN+MCH | `tagger.go` negative_aliases |
| **P2-8 (vòng 1)** | P2 | **VẪN MỞ**: 3 nguồn đạt trần 40 bài và log 3 cảnh báo `feed truncated`, nhưng cột `truncated` vẫn `f` toàn bộ | `store/runs.go` |
| **P1-NEW-3** | **P1** | Va chạm alias khi bỏ dấu trong tagger: alias GAS `"Khí Việt Nam"` khớp cụm thường `"Khi Việt Nam"` → gắn GAS `primary` 0.80 cho bài không liên quan (bài 108, 129). Cùng lớp lỗi với P1-NEW-1 | `tagger.go` scanAliases |
| **P1-NEW-4** | **P1** | `negative_aliases` không được áp ở tầng điểm 0.90: bài 108 gắn `FPT 0.90` dù thân bài chỉ có `"FPT Retail (FRT)"`, mà `FPT Retail` đã có trong `negative_aliases` của FPT trong DB | `tagger.go` |
| **P1-NEW-5** | **P1** | Rò rỉ sidebar vào tagger: bài 162 gắn BID, HDB — hai mã chỉ nằm trong `div.box-sidebar-*` (rail bài liên quan), không có trong thân bài | extractor / tagger input |
| **P2-NEW-4** | P2 | Trường `id` trong `backend/testdata/tagger_corpus.json` **không ánh xạ đúng** sang `articles.id` của DB (39/50 trỏ nhầm bài). Corpus phải khớp bằng tiêu đề mới dùng được → dễ gây hiểu nhầm khi đối chiếu | `backend/testdata/tagger_corpus.json` |

---

## V2.11 Việc cần làm để nghiệm thu vòng 3

**Bắt buộc (chặn phát hành):**

1. **Đóng P0-3 — tầng T3.** Cách rẻ nhất và hiệu quả nhất: **tắt hẳn tầng `score=0.40`**. Nó đóng góp 38 tag, **100% sai**, và xoá nó nâng trần precision từ 86,7% lên mức do các tầng 0.80–0.90 quyết định (hiện ~98% ở 0.90). Nếu muốn giữ tính năng "gợi ý theo ngành" thì tách ra khỏi trường `tickers` của API, đừng hiển thị như mã được gắn.
2. **Sửa lớp lỗi bỏ dấu (P1-NEW-1 + P1-NEW-3) cùng lúc.** Rà mọi chỗ so khớp cụm từ sau `Normalize`. Với cụm mà bỏ dấu gây va chạm thì so trên text **còn dấu**.
3. **P1-NEW-2 — lọc chủ đề tài chính.** Tin không có mã `primary` và không thuộc nhóm tài chính thì không nên lên đầu digest. Đây là thứ quyết định sản phẩm có đúng lời hứa PRD hay không.

**Nên làm:** P1-NEW-4, P1-NEW-5, và các P1 vòng 1 còn mở (P1-1, P1-2, P1-3, P1-5, P1-6).

**Để đo được chỉ tiêu PRD ở vòng 3:**
- **Corpus gán nhãn phải do người KHÁC viết tagger tạo ra**, hoặc tối thiểu tách tập tuning và tập đánh giá. Con số 100% hiện nay là đo trên tập tuning nên không có giá trị nghiệm thu. Test-plan mục 4.2 yêu cầu 50 bài gán nhãn thủ công — vẫn chưa có tập độc lập.
- Sửa ánh xạ `id` trong corpus (P2-NEW-4).
- Chạy Lighthouse throttle 4G cho LCP.
- Dựng 1 nguồn hỏng cố ý để mở khoá TC-PIPE-01/08/09 (3 case còn BLOCKED).

## V2.12 Bàn giao: corpus gán nhãn độc lập của QA

Xuất bộ mẫu QA đã gán nhãn tay ở vòng 2 thành cổng nghiệm thu:

**`backend/testdata/tagger_corpus_qa.json`** (195 KB, 27 mẫu)

| | |
|---|---|
| Mẫu | **27 bài** (nhóm A: 20 bài rút ngẫu nhiên · nhóm B: 7 bài tầng T3) |
| Tag `expected` (nhãn đúng) | **95** (61 `primary` · 34 `mentioned`) |
| Tag `must_not_tag` (tagger đang gắn, QA xác định SAI) | **42** |
| Tag `uncertain` (loại khỏi mọi phép tính) | **3** |
| Bài có **nhãn rỗng** | **12** — trong đó **10 bài rỗng-đã-xác-minh** |
| Phân bố nguồn | vneconomy 14 · vnexpress 6 · tinnhanhchungkhoan 4 · tuoitre 3 |
| Thân bài | 1.736–8.720 ký tự, **đã loại sidebar**, chạy tagger offline được, không cần mạng/DB |

Mỗi mẫu có `id`, `title`, `body`, `source_domain`, `url`, `expected[]`, `must_not_tag[]`, `uncertain[]`, `group`, `note`.

**Baseline của tagger hiện tại trên corpus này: 95/137 = 69,3%.**

> ⚠️ **Con số 69,3% KHÔNG phải precision của sản phẩm.** Nó thấp hơn ước lượng quần thể (80,5–83,2%) **một cách có chủ ý**: nhóm B gồm **toàn bộ** 7 bài tầng T3 nên corpus được làm giàu false positive để cổng chặn nhạy hơn. Dùng làm ngưỡng gate, **đừng trích như chỉ số báo cáo**. Cảnh báo này đã ghi trong `_meta`.

### Phát hiện hành động được — toàn bộ khoảng cách tới ngưỡng nằm ở 2 việc

42 tag sai phân bố theo score: **`0.40`: 38** · `0.80`: 2 · `0.90`: 2

| Việc | Precision trên corpus |
|---|---|
| Hiện tại | 95/137 = **69,3%** |
| **Tắt hẳn tầng T3 (`score=0.40`)** | 95/99 = **95,96%** — chạm ngưỡng, chưa vượt |
| **+ sửa nốt 4 tag sai còn lại** | 95/95 = **100%** |

4 tag sai ngoài T3, khớp đúng 3 nguyên nhân gốc đã ghi ở V2.4b:

```
art 108  FPT  score 0.90  vneconomy   ← negative_aliases không áp ("FPT Retail (FRT)")
art 108  GAS  score 0.80  vneconomy   ← va chạm alias bỏ dấu "Khí Việt Nam" ↔ "Khi Việt Nam"
art 129  GAS  score 0.80  tnck        ← cùng lỗi trên
art 168  VND  score 0.90  tuoitre     ← VND là ĐƠN VỊ TIỀN TỆ, vi phạm rule G-02
```

### Phát hiện bổ sung khi dựng corpus

Khi soát lại để gán nhãn, tôi tìm thêm **2 lỗi mà phép đo vòng 2 chưa bắt**, vì báo cáo trước chỉ đếm 22 tag T3 của 7 bài mà bỏ qua các tag điểm cao trên cùng bài đó:

1. **`VND` gắn `score 0.900` cho bài 168 là false positive** — cả 4 lần xuất hiện đều là tiền tệ (`cặp USD/VND`, `chênh lệch lãi suất VND - USD`, `tỉ giá USD/VND`), và **0 từ khoá ngữ cảnh cổ phiếu** trong toàn bài. Đây là **rule G-02 / TC-TAG-07 hỏng ở tầng điểm cao nhất**, không phải lỗi T3.
2. **Bài 191 KHÔNG phải nhãn rỗng** — 4 mã ngân hàng `BID, MBB, TCB, VCB` ở score 0.900 là **ĐÚNG** (bài nêu đích danh giá trị mua ròng: *"VCB 229 tỉ đồng, BSR 158 tỉ đồng, BID 125 tỉ đồng, TCB 105 tỉ đồng và MBB 98 tỉ đồng"*). Chỉ 4 mã dầu khí ở T3 là sai. Nếu tôi gán nhãn rỗng cho cả bài như giả định ban đầu thì corpus đã sai.

> Ghi nhận phương pháp: đây là lý do phải **tự kiểm lại từng bài khi gán nhãn** thay vì tin số tổng hợp. Con số "22 tag, tất cả sai" ở V2.4b vẫn đúng **cho riêng tầng T3**, nhưng 7 bài đó mang tổng cộng 27 tag chứ không phải 22.

### Giới hạn đã ghi rõ trong `_meta`
- `expected` **đầy đủ cho precision** (mọi tag tagger sinh ra đều đã chấm) nhưng **chưa quét vét cạn** mọi mã có thể có → **recall chỉ dùng tham khảo**, sẽ hơi lạc quan. Ví dụ đã biết: bài 191 có nêu `BSR` nhưng tagger không gắn và tôi chưa đưa vào `expected`.
- Bài 78 có `expected` rỗng **vì tag duy nhất bị xếp `uncertain`**, không phải vì đã xác minh là không có mã → đã ghi `note` riêng, **không dùng làm mẫu bắt false positive**. 10 bài rỗng-đã-xác-minh liệt kê ở `_meta.counts.verified_empty_expected_ids`.
- 3 tag `uncertain` (mã chỉ nằm trong cụm trích dẫn tổ chức phân tích hoặc công thức lãi suất tham chiếu) **phải loại khỏi cả tử số lẫn mẫu số**.

## V2.13 Bằng chứng & dọn dẹp vòng 2

Ảnh chụp trong `docs/05-qa/evidence/`:

| File | Nội dung |
|---|---|
| `qa-r2-note-light.png` | **Trang Research Note (Ảnh A) — lần đầu render được**, full page, nền sáng |
| `qa-r2-home-light.png` | Trang chủ 1440px nền sáng, dữ liệu thật |
| `qa-r2-home-dark.png` | Trang chủ 1440px nền tối |
| *(vòng 1)* `qa-reference-compare.png`, `qa-home-1440-light.png`, `qa-home-390-light.png`, `qa-mobile-card.png` | giữ nguyên để đối chiếu |

**Không sửa một dòng code sản phẩm nào.** Chỉ tạo `deploy/.env` từ `.env.example` (đặt `ADMIN_TOKEN` để test) và 2 file probe test tạm trong package `summarizer`/`tagger`, **đã xoá ngay sau khi chạy** (đã xác minh `zz files remaining: 0`). Dữ liệu test duy nhất bị thay đổi là `crawl_runs.finished_at` lùi 14h để mở khoá TC-UI-09, **đã revert**.

---

<a id="vong-1"></a>

# ══════ VÒNG 1 — Nghiệm thu lần đầu (giữ để đối chiếu) ══════

## 0. Kết luận nghiệm thu (vòng 1)

### ❌ KHÔNG ĐẠT — chặn phát hành

Sản phẩm có nền tảng kỹ thuật tốt (96/112 case PASS, toàn bộ cổng tự động sạch, bảo mật short link vững, giao diện bám sát design system gần như tuyệt đối, **0 số liệu bịa**). Nhưng có **5 lỗi P0** khiến không thể phát hành:

| # | Lỗi P0 | Hệ quả |
|---|---|---|
| **P0-1** | `/ma/{mã}` và `/nhan-dinh/{slug}` trả **HTTP 500** | **2 trong 3 trang của sản phẩm không chạy.** Mất MVP MUST #9 và #10 |
| **P0-2** | **36/36 tin VnEconomy dùng chung MỘT tóm tắt lạc đề** | 19% kho tin hiển thị nội dung hoàn toàn sai bài |
| **P0-3** | **Precision tag mã ≈ 52%** (ngưỡng PRD ≥ 96%) | Gắn sai mã cho tin — đúng thứ PRD nói "gắn sai mã tệ hơn bỏ sót" |
| **P0-4** | Bộ lọc ngôn ngữ khuyến nghị đầu tư **bị vượt qua** — và đã lọt tin thật | Tin id 8 publish câu *"nhà đầu tư ngắn hạn nên duy trì tỷ trọng cổ phiếu…"*. Rủi ro pháp lý (PRD mục 11) |
| **P0-5** | Healthcheck `api` **không bao giờ pass** → compose production không khởi động nổi `web` | `make up` trên máy sạch cho ra site chết hoàn toàn |

P0-1, P0-4, P0-5 sửa nhỏ (1 dòng type, 1 dòng biên từ, 1 dòng healthcheck). **P0-2 và P0-3 là lỗi chất lượng nội dung, cần làm lại logic bóc tách và gắn mã — đây là phần nặng nhất.**

> **Điểm sáng quan trọng nhất:** kiểm 20/20 bài, **không có một số liệu bịa nào** (0 fabrication), 58/60 câu khớp nguyên văn bài gốc. Chỉ tiêu "tồn vong" của PRD (Q5) đạt.

### Đối chiếu tiêu chí phát hành (test-plan mục 6)

| Tiêu chí | Ngưỡng | Thực tế | Kết luận |
|---|---|---|---|
| Case P0 (TC-VAL, TC-SEC) PASS | 100% | TC-SEC 10/10 ✅ · TC-VAL 8/10 ❌ | ❌ |
| Case P1 PASS | ≥ 95% | 88/96 = 91,7% | ❌ |
| **Tin publish có số bịa** | **0** | **0** — kiểm 20/20 bài | ✅ |
| **Tóm tắt đúng bài** | 100% | **81%** — 36/190 tin VnEconomy sai hoàn toàn | ❌ |
| `go build` / `vet` / `gofmt` / `test` | sạch | sạch, kể cả `-race` | ✅ |
| `tsc --noEmit` / `lint` / `build` | sạch | sạch | ✅ |
| `docker compose up` → 4 service healthy | 4/4 | **3/4** (`api` unhealthy) | ❌ |
| Đối chiếu thị giác 2 ảnh | đạt | Ảnh B đạt ✅ · Ảnh A **không render được** ❌ | ❌ |
| Precision tag mã ≥ 96% | ≥96% | **Không đủ dữ liệu để kết luận** | ⚠️ |
| Lỗ hổng bảo mật P0 | 0 | 0 | ✅ |

**Điều kiện chặn phát hành ngay (test-plan mục 6) — kiểm riêng:**

| Điều kiện chặn | Trạng thái |
|---|---|
| Tin publish với số liệu bịa | ✅ Không phát hiện (xem mục 5) |
| Short link redirect ra ngoài allowlist | ✅ Không — 12/12 input tấn công bị chặn |
| Bảng scroll ngang trên mobile | ✅ Không — `scrollWidth <= innerWidth` ở mọi viewport |
| Thiếu disclaimer trên research note | ⚠️ Có trong API + footer, nhưng **trang không render được** |

---

## 1. Tổng hợp kết quả

**112 case · 96 PASS · 8 FAIL · 8 BLOCKED**

| Nhóm | Tổng | PASS | FAIL | BLOCKED |
|---|---|---|---|---|
| TC-VAL — Kỷ luật số liệu (P0) | 10 | 8 | **2** | 0 |
| TC-TAG — Gắn mã CK | 18 | 17 | 1 | 0 |
| TC-API — Hợp đồng API | 20 | 20 | 0 | 0 |
| TC-SEC — Bảo mật (P0) | 10 | **10** | 0 | 0 |
| TC-PIPE — Pipeline & Cron | 10 | 5 | 0 | 5 |
| TC-UI — Giao diện | 14 | 11 | **2** | 1 |
| TC-RWD — Responsive | 6 | 6 | 0 | 0 |
| TC-A11Y — Tiếp cận | 10 | 8 | 2 | 0 |
| TC-SEO | 6 | 5 | 0 | 1 |
| TC-DEPLOY | 8 | 7 | **1** | 0 |

---

## 2. Việc 1 — Cổng tự động (output thật)

### Backend — tất cả sạch

```
$ go version
go version go1.26.2 darwin/arm64

$ go build ./...
(không output — sạch)

$ go vet ./...
(không output — sạch)

$ gofmt -l .
(không output — không file nào cần format)

$ go test ./... -count=1
ok  	github.com/tenpoint/tenpoint-api/internal/config	0.550s
ok  	github.com/tenpoint/tenpoint-api/internal/http	0.994s
ok  	github.com/tenpoint/tenpoint-api/internal/ingest/classifier	3.388s
ok  	github.com/tenpoint/tenpoint-api/internal/ingest/dedupe	2.514s
ok  	github.com/tenpoint/tenpoint-api/internal/ingest/extractor	1.476s
ok  	github.com/tenpoint/tenpoint-api/internal/ingest/summarizer	2.001s
ok  	github.com/tenpoint/tenpoint-api/internal/ingest/tagger	3.851s
ok  	github.com/tenpoint/tenpoint-api/internal/scheduler	5.190s
ok  	github.com/tenpoint/tenpoint-api/internal/shortlink	2.978s
ok  	github.com/tenpoint/tenpoint-api/internal/store	4.364s
?   	.../cmd/api, .../cmd/ingest, .../internal/domain, .../internal/ingest,
    .../internal/ingest/fetcher, .../internal/textutil, .../migrations	[no test files]

$ go test ./... -race -count=1
(toàn bộ ok — không có data race)
RACE_EXIT=0
```

> ⚠️ Ghi nhận: `internal/ingest` (pipeline) và `internal/textutil` **không có test file nào**. `textutil` chứa `NormalizeVNNumber` — lõi của toàn bộ an toàn số liệu.

### Frontend — tất cả sạch

```
$ npx tsc --noEmit
(không output — sạch)

$ npm run lint
> eslint
(không output — sạch)

$ npm run build
▲ Next.js 16.3.5 (Turbopack)
✓ Compiled successfully in 383ms
✓ Generating static pages using 10 workers (7/7) in 271ms
Route (app)
┌ ƒ /            ├ ○ /_not-found   ├ ƒ /ma/[symbol]
├ ƒ /nhan-dinh   ├ ƒ /nhan-dinh/[slug]   ├ ○ /robots.txt   └ ƒ /sitemap.xml
```

> ⚠️ **Quan trọng:** cổng tự động sạch **nhưng không bắt được P0-1**. Lý do: `/ma/[symbol]` và `/nhan-dinh/[slug]` là route động (`ƒ`), không prerender lúc build, nên lỗi runtime không lộ ra. `tsc` cũng không bắt vì `lib/api.ts:50` ép kiểu không kiểm chứng (`as T`). Đây là bài học về cổng chất lượng, không chỉ là một bug.

### Hạ tầng

```
$ docker compose -f deploy/docker-compose.yml config
CONFIG_EXIT=0   (parse sạch, không warning)
```

---

## 3. Việc 2 — Dựng hệ thống & test end-to-end

### Sửa cấu hình để chạy được (ghi rõ theo yêu cầu)

Chỉ sửa `deploy/.env` (file môi trường cục bộ, đã nằm trong `.gitignore`). **Không sửa code sản phẩm.**

| Sửa | Lý do |
|---|---|
| `ADMIN_TOKEN=` → `ADMIN_TOKEN=qa-test-token-9f3a2b` | `.env.example` để trống; cần token để test TC-API-15/16/17 |
| `DEV_API_PORT=8080` → `8090` | Cổng 8080 trên máy test đã bị container `job-bot` không liên quan chiếm. Lỗi môi trường, không phải lỗi sản phẩm |

Ngoài ra `docker compose` được gọi với `env -u ANTHROPIC_MODEL ...` vì shell của máy test có sẵn biến `ANTHROPIC_*` sẽ ghi đè `.env`.

### Trạng thái service

```
NAME             STATUS
tenpoint-api     Up 3 minutes (unhealthy)     ← ❌ P0-5
tenpoint-caddy   Up 16 minutes (healthy)
tenpoint-db      Up 17 minutes (healthy)
tenpoint-web     Up 16 minutes (healthy)
```

Bootstrap T1 hoạt động đúng: khởi động xong tự chạy ingest thật ngay.

```
{"msg":"bootstrap: no real articles yet, running first ingest now"}
{"msg":"crawl run started","run_id":1,"trigger":"manual"}
{"msg":"crawl run finished","run_id":1,"found":190,"new":187,"updated":0,
 "unchanged":0,"rejected":0,"errors":0,"duration":"2m53.583227322s"}
```

### TC-API — 20/20 PASS

| ID | Kết quả | Bằng chứng |
|---|---|---|
| TC-API-01 | ✅ PASS | `meta:{total:93,limit:20,offset:0,has_more:true}`; `sorted_desc=True`; cửa sổ mặc định chứa 14–16/09 (trong 7 ngày) |
| TC-API-02 | ✅ PASS | `?ma=FPT` → `total=2`, `violations=[]` |
| TC-API-03 | ✅ PASS | `?ma=FPT,HPG` → `total=4`, `violations=[]` (logic OR đúng) |
| TC-API-04 | ✅ PASS | `relevance=primary`→2 tin, relevance thấy `['primary']`; `relevance=all`→18 tin, thấy `['mentioned','primary']` |
| TC-API-05 | ✅ PASS | **Cả 9 giá trị** trả 200 và `pure=True` (không lẫn loại khác) |
| TC-API-06 | ✅ PASS | `q=chung khoan` (không dấu) → `total=29`, bằng đúng `q=chứng khoán` → `total=29` |
| TC-API-07 | ✅ PASS | `?limit=999` → `limit=100, n=100` (bị chặn) |
| TC-API-08 | ✅ PASS | `tu=den=15/09`→43 tin, tất cả `dates_vn=['15/09/2026']`; `16/09`→130 tin, tất cả `['16/09/2026']` |
| TC-API-09 | ✅ PASS | `404` + `{"error":{"code":"not_found","message":"Không tìm thấy tin này."}}` |
| TC-API-10 | ✅ PASS | `/tickers/FPT` → 200, `recent_news n=20` |
| TC-API-11 | ✅ PASS | `/notes/cmg-quy-2-2026` → 200, `n points=4 ordinals=[1,2,3,4]`, có `disclaimer` |
| TC-API-12 | ✅ PASS | `/meta` trả **đủ 9 nhãn kèm count**; `last_crawl_at:"2026-09-16T14:22:36Z"`, `is_stale:false` |
| TC-API-13 | ✅ PASS | `/healthz` → `{"status":"ok"}` 200 |
| TC-API-14 | ✅ PASS | Dừng `tenpoint-db` → `/readyz` `503` `{"db":"down","reason":"Không kết nối được cơ sở dữ liệu."}`; `/healthz` vẫn 200 (tách liveness/readiness đúng); phục hồi sau khi db lên |
| TC-API-15 | ✅ PASS | Không token → `401` |
| TC-API-16 | ✅ PASS | Token sai → `401` `{"code":"unauthorized","message":"Token quản trị không hợp lệ."}` |
| TC-API-17 | ✅ PASS | Token đúng → `202` `{"status":"accepted","trigger":"manual"}` |
| TC-API-18 | ✅ PASS | Gọi khi run đang chạy → **`409`** ×2 lần liên tiếp (advisory lock đúng) |
| TC-API-19 | ✅ PASS | `2026-09-16T13:35:00Z` → `16/09/2026`, đúng `DD/MM/YYYY` |
| TC-API-20 | ✅ PASS | `summary_md` có `**...**`, `has <tag>: False` — markdown thô, không phải HTML |

**Tham số `loai` không hợp lệ** — đúng yêu cầu, không im lặng trả rỗng:

```
$ curl "…/api/v1/news?loai=abc"     → status=400
{"error":{"code":"invalid_request","message":"Giá trị loai không hợp lệ: abc.
 Các giá trị hợp lệ: vi_mo, nganh, doanh_nghiep, thi_truong, khoi_ngoai,
 co_tuc_phat_hanh, phap_ly, phan_tich, trai_phieu_tin_dung."}}
```

### Short link — điểm khách hàng phàn nàn: ĐÃ SỬA ĐƯỢC

Đây là phần khách hàng phản ánh *"phần trích dẫn link chưa đưa đến thông tin trang cụ thể"*. **Chứng minh bằng 2 bước: redirect, rồi curl chính URL đích.**

```
/r/SCsju3 → 302 → https://cafef.vn/mot-cong-ty-quan-ly-quy-bi-xu-phat-...-188260916203508002.chn
/r/ur6sSx → 302 → https://vnexpress.net/moi-co-hai-dia-phuong-ho-tro-...-5121075.html
/r/wzUJLS → 302 → https://vneconomy.vn/tung-vo-no-hy-lap-bay-gio-di-vay-...htm
/r/fnuBA9 → 302 → https://tuoitre.vn/hai-quy-dragon-capital-mua-lai-...-100260916180214071.htm
/r/IHvecw → 302 → https://tuoitre.vn/mac-can-voi-khoan-vay-mua-nha-mua-xe-100260916182312457.htm
```

Curl trực tiếp từng URL đích — **5/5 trả 200 với đúng bài**:

```
final_status=200 bytes=104272  <title>: Một công ty quản lý quỹ bị xử phạt do vi phạm công bố thông tin
final_status=200 bytes=269711  <title>: Mới có hai địa phương hỗ trợ người dân lắp điện mặt trời mái nhà - VnExpress
final_status=200 bytes=40437   <title>: (trang không có thẻ title, HTTP 200, 40KB nội dung)
final_status=200 bytes=472809  <title>: Hai quỹ Dragon Capital mua lại gần 750.000 cổ phiếu PNJ… - Tuổi Trẻ
final_status=200 bytes=479338  <title>: 'Mắc cạn' với khoản vay mua nhà, mua xe - Tuổi Trẻ Online
```

`<title>` của trang đích **khớp** với `source.article_title` lưu trong DB → link dẫn đúng bài cụ thể, không phải trang chủ.

**Tin demo → `410`, không redirect:**

```
$ curl -D - http://…/r/e5Nh1s
HTTP/1.1 410 Gone
Content-Type: text/html; charset=utf-8
<html lang="vi"> … <title>Bài gốc không còn truy cập được — TenPoint</title>
<h1>Bài gốc không còn truy cập được</h1>
```

**Mã không tồn tại → `404`, không redirect:**

```
$ curl -D - http://…/r/khongtontai
HTTP/1.1 404 Not Found
{"error":{"code":"not_found","message":"Liên kết không tồn tại."}}
```

---

## 4. Việc 3 — Giao diện

### 4.1 Đối chiếu ảnh tham chiếu

**Ảnh B (News Digest) — ĐẠT.** Chụp ở chế độ dữ liệu mẫu (đúng 5 tin từ ảnh gốc) để so 1:1.

| Điểm đối chiếu | Ảnh gốc | Thực tế | |
|---|---|---|---|
| 5 cột đúng thứ tự | Ngày·Mã CK·Tóm tắt·Source·Loại tin | y hệt | ✅ |
| Hairline ngăn hàng | 1px xám nhạt | `1px solid rgb(213,216,210)` = `--rule` | ✅ |
| Kẻ đậm dưới header | có | `1px solid rgb(20,24,27)` = `--rule-strong` | ✅ |
| Zebra striping | không | 6 hàng đầu đều `rgba(0,0,0,0)` | ✅ |
| Border dọc | không | `borderLeft/Right = 0px/0px` cả 5 ô | ✅ |
| Bo góc | không | toàn trang chỉ `2px` ×14 (chip mã CK) | ✅ |
| Shadow | không | `shadowsUsed = {}` — **không có shadow nào** | ✅ |
| Gradient | không | `gradients = []` | ✅ |
| Số in đậm chỉ bằng weight | có | `weight:700`, `bg:transparent`, `deco:none` | ✅ |
| Cột Source = domain gạch chân | có | `cafef.vn`, `vnexpress.net`, `thanhnien.vn` | ✅ |
| Loại tin = text thuần | có | không nền, không badge màu | ✅ |
| Mã CK wrap, >8 mã → +N | có | hàng 11 mã hiện 8 + `+3` | ✅ |
| Tiêu đề động | `8 tin mới đáng chú ý — PVS · FPT · VIC.` | `5 tin mới đáng chú ý — VIC · HPG · VHM.` | ✅ đúng format |

Các cụm bold tái tạo **chính xác** ảnh gốc: `408.000 tỷ đồng` · `hơn 100 USD/tấn thép` · `1.821 điểm (+29,91 điểm, +1,67%)` · `494.000 tấn thịt, ~1,5 tỷ USD (+10% lượng, +36,4% giá trị)` · `5.588 tỷ đồng (216 triệu USD)`.

> Khác biệt duy nhất: thứ tự hàng sắp theo `published_at` giảm dần (đúng spec), ảnh gốc không sắp theo ngày. Không phải lỗi.

**Ảnh A (Research Note) — KHÔNG ĐỐI CHIẾU ĐƯỢC.** Trang `/nhan-dinh/cmg-quy-2-2026` trả **HTTP 500** (P0-1). Đây là research note duy nhất trong sản phẩm.

### 4.2 Dark mode — token khớp tuyệt đối 26/26

Đọc `getComputedStyle` thật, không đọc file CSS.

**Light** (`data-theme="light"`, ép thủ công khi OS đang dark → chứng minh nút chuyển ghi đè được `prefers-color-scheme`):

```
--surface #f4f5f2 ✓   --surface-raised #fcfcfb ✓   --surface-sunken #e7e9e4 ✓
--ink #14181b ✓       --ink-2 #39424a ✓            --ink-muted #6b7580 ✓
--ink-faint #9ba3ad ✓ --rule #d5d8d2 ✓             --rule-strong #14181b ✓
--accent #1e4b8f ✓    --accent-hover #163a70 ✓     --accent-soft #e3eaf4 ✓
--focus #1e4b8f ✓                                   mismatches: []
```

**Dark** (theo `prefers-color-scheme: dark`, không cần thao tác):

```
--surface #101317 ✓   --surface-raised #171b21 ✓   --surface-sunken #0b0e11 ✓
--ink #e9ecef ✓       --ink-2 #b4bcc6 ✓            --ink-muted #7e8894 ✓
--ink-faint #5a636e ✓ --rule #262c34 ✓             --rule-strong #e9ecef ✓
--accent #6fa0f0 ✓    --accent-hover #96baf6 ✓     --accent-soft #16233a ✓
--focus #6fa0f0 ✓
```

Cả 2 cơ chế đều hoạt động: `prefers-color-scheme` (mặc định) và nút chuyển thủ công (ghi đè được cả 2 chiều).

Màu kem bị cấm `#F7F5F0` / `#f5f1ea` / `#fbf8f1` / `#efeae0`: **0 lần xuất hiện** trong toàn bộ `frontend/`.

### 4.3 Cột Source

```json
{
  "isAnchor": true,
  "href": "/r/SCsju3",
  "target": "_blank",
  "rel": "nofollow noopener noreferrer",
  "title": "Một công ty quản lý quỹ bị xử phạt do vi phạm công bố thông tin · cafef.vn · 16/09/2026",
  "ariaLabel": "Đọc bài gốc: Một công ty quản lý quỹ bị xử phạt do vi phạm công bố thông tin. Nguồn CafeF, đăng 16/09/2026 (mở tab mới)"
}
```

✅ `title` và `aria-label` đều chứa **tiêu đề bài gốc** + tên báo + ngày + "(mở tab mới)". Đúng design-system 5.3 — trả lời trực tiếp phản hồi của khách hàng.

**Tin `is_demo`:** nhánh code render `<span>` (không phải `<a>`) kèm chữ **"Dữ liệu mẫu, chưa có bài gốc"** — `SourceLink.tsx:29-41`. Với dữ liệu thật, demo bị ẩn hoàn toàn nên không quan sát được trên trang. **Nhưng ở chế độ fixture (API chết) thì nhánh này KHÔNG chạy** → xem defect P1-1.

### 4.4 Nút "Copy bản tin" — PASS cả 2 tiêu chí

```
blockCount: 20
linesPerBlock: [3]          ← đúng 3 dòng/tin, không có ngoại lệ
containsEmDash: false       ← KHÔNG có ký tự "—"
emDashCount: 0
```

Mẫu khối:

```
[16/09/2026] [SSI, VCI, VND]
Mỗi QFS dành cho khách hàng cá nhân có giá quy đổi 3,5 - 6 triệu đồng…
Nguồn: http://127.0.0.1:3000/r/t4BFzT
```

⚠️ Khiếm khuyết: **8/20 khối có ngoặc mã rỗng `[]`** (tin không gắn được mã nào). Dán vào Zalo trông lỗi. → defect P2-4.

### 4.5 Responsive — 6/6 PASS, không scroll ngang ở mọi viewport

| Viewport | `body.scrollWidth` | `innerWidth` | Kết quả | Bố cục |
|---|---|---|---|---|
| 375px | 375 | 375 | ✅ | card dọc (`display:grid`) |
| 390px | 390 | 390 | ✅ | card dọc |
| 900px | 900 | 900 | ✅ | 4 cột, `Loại tin` gộp vào dòng meta |
| 1440px | 1440 | 1440 | ✅ | 5 cột |
| 1920px | 1920 | 1920 | ✅ | 5 cột, container khoá **1180px** |

Không tìm thấy phần tử nào tràn (`overflowCulprit: null`), không có wrapper `overflow-x: auto` (`xScrollers: []`).

Thứ tự trong card ở 390px: `16/09/2026 · Pháp lý` → `Mã CK` → `Tóm tắt` → `cafef.vn` — đúng spec mục 6.

> Đã kiểm riêng nghi vấn "Loại tin hiện 2 lần ở mobile": **không xảy ra**. `.digest-type-inline` có `display:none` ở <768px, chỉ ô `td[data-col=type]` hiển thị. Ảnh chụp card xác nhận chỉ có một "Pháp lý".

### 4.6 Console

```
Trang chủ (/):        Total messages: 0 (Errors: 0, Warnings: 0)   ✅
/nhan-dinh:           0 lỗi                                          ✅
/nhan-dinh/{slug}:    3 errors  ← do P0-1
```

### 4.7 A11y — checklist design-system mục 8

| Mục | Kết quả | Bằng chứng |
|---|---|---|
| `<html lang="vi">` | ✅ | `htmlLang: "vi"` |
| `<caption>` sr-only + `<th scope="col">` | ✅ | caption `"Bảng tin chứng khoán tổng hợp, cập nhật 21:19 16/09/2026"` class `sr-only`; `thScopes: ["col"×5]` |
| Focus ring, không bao giờ `outline:none` | ✅ | `:focus-visible { outline: 2px solid var(--focus); outline-offset: 2px; }`; `outlineNoneCount: 0` |
| Link tab mới có aria-label + tiêu đề bài | ✅ | mục 4.3 |
| Filter chips là `<button aria-pressed>` | ✅ | 14 phần tử `aria-pressed`, đều `tag:"BUTTON"`; `div[onclick]: 0` |
| Skip link | ✅ | `"Bỏ qua tới nội dung chính"` → `#noi-dung` |
| **Contrast ≥ 4.5:1** | ❌ **FAIL** | `--ink-muted #6B7580` trên `#F4F5F2` = **4.28:1** (< 4.5). Dùng cho `.digest-date` và `.digest-type` cỡ 13.5px |
| **Target ≥ 24×24px** | ❌ **FAIL** | chip mã CK trong bảng `26×18`; link Source `~17px` cao; nav `18px`; "Lưu danh sách theo dõi" `141×20` |
| Font đủ dấu tiếng Việt | ✅ | `Ừ Ữ Ỡ Ợ Ặ Ẫ Ỹ ọ ự ẳ` — `missing: []` ở cả 3 font (Source Serif 4, Be Vietnam Pro, JetBrains Mono) |
| `prefers-reduced-motion` | ✅ | block `@media` có, khớp spec mục 7 |

> Lưu ý về contrast: design-system mục 2.4 **tự ghi sai** — bảng ghi `--ink-muted` là "~4.6:1 ✅ AA", đo thật là 4.28:1. Cần sửa cả token lẫn tài liệu.

---

## 5. Việc 4 — Kiểm chứng nội dung

Mẫu: 20 tin thật (`is_demo=false`) lấy từ `/api/v1/news?limit=20`.

### 5.1 Rubric NUM-1..NUM-8

| ID | Số từ | Cụm bold | Nội dung cụm bold | B-3 |
|---|---|---|---|---|
| 7 | 101 | 1 | `2009` | ❌ năm trần, cắt đôi "tháng 3/2009" |
| 195 | 107 | 1 | `52 triệu` | ❌ gốc là "39 - 52 triệu **đồng**" |
| 25 | 109 | 1 | `16` | ❌ cắt đôi ngày "Ngày 16/9" |
| 10 | 85 | 1 | `2026` | ❌ năm trần |
| 128 | 97 | 1 | `55.100` | ❌ gốc "55.100 **cổ phiếu** PNJ" |
| 132 | 92 | 1 | `18,7 triệu` | ❌ gốc "gần 18,7 triệu **đồng**" |
| 12 | 98 | 1 | `2.848 tỷ đồng` | ✅ |
| 9 | 97 | 1 | `2026` | ❌ năm trần |
| 14 | 85 | 1 | `2026` | ❌ năm trần |
| 29 | 107 | **0** | — | ❌ **NUM-3: 0 cụm bold** |
| 134 | 103 | 1 | `16.000` | ❌ gốc "16.000**m²**" |
| 33 | 113 | 1 | `59,3 tỷ USD` | ✅ |
| 15 | 93 | 1 | `2,26%` | ✅ |
| 6 | 81 | 1 | `2026` | ❌ năm trần |
| 8 | 117 | 1 | `17` | ❌ cắt đôi "ngày 17/9" |
| 11 | 81 | 1 | `1.800 điểm` | ⚠️ cắt đôi khoảng "1.792-1.800 điểm" |
| 16 | 97 | 1 | `600.000 đồng` | ⚠️ rơi mất "/lượng" |
| 138 | 97 | 1 | `27.200` | ❌ gốc "27.200 **nhân viên**" |
| 19 | 90 | 1 | `2026` | ❌ năm trần |
| 139 | 131 | 1 | `16` | ❌ cắt đôi "Ngày 16-9" |

**Tổng hợp:**

| Rubric | Ngưỡng | Kết quả | |
|---|---|---|---|
| **NUM-1 Độ dài** | 60–150 từ, reject >180 | **20/20 đạt** (81–131 từ, median ~97). 0 bài vượt 180, 0 bài dưới 60 | ✅ |
| **NUM-3 Số cụm bold** | đúng 1–3 | 19/20 có đúng 1 cụm; **id 29 có 0 cụm** | ❌ 1 vi phạm |
| **B-3 Chất lượng cụm bold** | số + đơn vị trọn vẹn | **14/19 vi phạm = 73,7%** (+2 biên → 16/19). Chỉ 3 cụm sạch | ❌ |
| **NUM-7 Không khuyến nghị đầu tư** | 0 cụm | **1/20 vi phạm** (id 8) | ❌ P0-4 |

Phân bố cụm bold: `0 cụm ×1, 1 cụm ×19, 2 cụm ×0, 3 cụm ×0` — bộ sinh bị khoá cứng ở đúng 1 cụm, không bao giờ dùng 2–3.

Phân loại 14 vi phạm B-3: năm trần ×6 (`2009`, `2026`×5) · số trần cắt đôi ngày ×3 (`16`,`17`,`16`) · số trần mất đơn vị ×3 (`55.100`,`16.000`,`27.200`) · cụm ghép bị cắt ×2 (`52 triệu`, `18,7 triệu`). Heuristic đang chọn **dãy chữ số đầu tiên** thay vì chọn cụm số + đơn vị.

### 5.2 Truy vết số liệu — ✅ 0 SỐ BỊA (chỉ tiêu tồn vong ĐẠT)

8 bài kiểm sâu (id 7, 9, 12, 16, 33, 128, 132, 195): resolve `/r/{code}` → curl URL đích → strip HTML → đối chiếu **từng con số** trong `summary_md` với text bài gốc. Cả 8 trang fetch được (200, không paywall).

| id | Số lượng con số | Kết quả |
|---|---|---|
| 7 | 11 | ✅ khớp hết, 3/3 câu nguyên văn |
| 9 | 12 | ✅ khớp hết |
| 12 | 14 | ✅ khớp hết |
| 16 | 17 | ✅ khớp hết, 5/5 câu nguyên văn |
| 33 | 18 | ✅ khớp hết |
| 128 | 15 | ✅ khớp hết, 4/4 câu nguyên văn |
| 132 | 10 | ✅ khớp hết |
| 195 | 11 | ✅ khớp hết |

Mở rộng ra **cả 20 bài**: **mọi con số trong mọi tóm tắt đều xuất hiện nguyên văn trong bài gốc.**

> **SỐ BỊA ĐƯỢC XÁC NHẬN: 0. BLOCKED (không fetch được): 0.**
> **58/60 câu tóm tắt khớp byte-verbatim** với bài gốc. 2 câu lệch đều lành tính (1 thừa khoảng trắng trước dấu chấm; 1 ghép 2 gạch đầu dòng, cả hai đều nguyên văn).

Nguyên nhân kết quả tốt: summarizer mặc định là **extractive** — copy nguyên câu chứ không viết lại, nên số không thể trôi.

> ⚠️ **Nhưng chính điều đó che giấu P0-2:** vì chỉ copy, cổng kiểm số *luôn* pass kể cả khi copy **nhầm khối text**. Một validator chỉ đối chiếu số sẽ **không bao giờ** bắt được lỗi VnEconomy. Đây là khoảng trống thiết kế của validator, không chỉ là một bug.

*Lưu ý kỹ thuật cho lần chạy sau:* `tinnhanhchungkhoan.vn` trả gzip; thiếu `curl --compressed` sẽ thấy trang rỗng và báo động giả 17 số bịa ở id 16.

### 5.3 Ẩn dữ liệu demo khi có tin thật (edge-case mục 0) — ✅ PASS

| Tầng | Kết quả | Bằng chứng |
|---|---|---|
| **T1** Bootstrap ingest thật khi khởi động | ✅ | `"bootstrap: no real articles yet, running first ingest now"` → 187 tin thật sau 2m53s |
| **T2** Đánh dấu dữ liệu demo | ✅ | 5 bài seed `is_demo=t`, `short_links.target_alive=f`; API trả `is_demo:true, short_link_alive:false` |
| **T3** Ẩn demo khi có tin thật | ✅ | `demo items in default feed = []`. DB có 110 tin thật + 5 demo; feed mặc định **0 demo** |

Không có tham số nào (`is_demo=true`, `demo=true`, `include_demo=true`) làm lộ demo ra feed.

---

## 6. Việc 5 — Bảo mật: TC-SEC 10/10 PASS

| ID | Kết quả | Bằng chứng |
|---|---|---|
| TC-SEC-01 | ✅ | `https://evil.example.com` → `REJECTED \| host "evil.example.com": target domain is not in the allowlist`; có `log.Warn` tại `pipeline.go:330` |
| TC-SEC-02 | ✅ | `https://sub.cafef.vn/x` → `ACCEPTED` |
| TC-SEC-03 | ✅ | `https://cafef.vn.evil.com/x` → **REJECTED**. Khớp theo nhãn domain (`strings.HasSuffix(host, "."+d)`), không khớp chuỗi con |
| TC-SEC-04 | ✅ | `/r/khongtontai` → `404`, **không redirect** |
| TC-SEC-05 | ✅ | `' OR 1=1--`, `'; DROP TABLE articles;--`, `%' UNION SELECT NULL--` → 200, `total=0`, bảng `articles` còn nguyên 195 dòng |
| TC-SEC-06 | ✅ | `parseInlineMarkdown("<script>alert(1)</script>")` → `[{"type":"text","value":"alert(1)"}]`. Render bằng React node, không `dangerouslySetInnerHTML` |
| TC-SEC-07 | ✅ | `<img src=x onerror=alert(1)>` → `[]` (rỗng); `<svg/onload=…>` → `[]`. Chỉ tồn tại token `strong`/`em`/`text` |
| TC-SEC-08 | ✅ | `rel="nofollow noopener noreferrer"` + `target="_blank"` (đo trên DOM thật) |
| TC-SEC-09 | ✅ | Quét `sk-ant-*`, `sk-*`, `AKIA*`, `ghp_*`, `glpat-`, `xox[baprs]-`, PRIVATE KEY → **0 match**. Secret chỉ đọc từ env |
| TC-SEC-10 | ✅ | Runtime: `api uid=10001(tenpoint)`, `web uid=1001(nextjs)` — **non-root cả hai** |

**Bộ input tấn công đầy đủ đã chạy — 12/12 xử lý đúng:**

```
https://evil.example.com      → REJECTED    https://cafef.vn.evil.com/x  → REJECTED
https://CAFEF.VN.EVIL.COM/x   → REJECTED    https://xcafef.vn/x          → REJECTED
https://cafef.vn@evil.com/x   → REJECTED    https://evil.com#cafef.vn    → REJECTED
javascript:alert(1)           → REJECTED    JaVaScRiPt:alert(1)          → REJECTED
//evil.com/x                  → REJECTED    data:text/html,<script>…     → REJECTED
https:/\/\evil.com/x          → REJECTED    http://127.0.0.1:8080/x      → REJECTED
https://sub.cafef.vn/x        → ACCEPTED    https://deep.sub.cafef.vn/x  → ACCEPTED
https://www.cafef.vn/x        → ACCEPTED
```

Allowlist **được gọi thật** ở `pipeline.go:329` (trước fetch) và `pipeline.go:354-357` (khi nhận `<link rel=canonical>`) — không phải hàm chết.

**Điểm cộng:** `handlers_admin.go:73` dùng `subtle.ConstantTimeCompare` cho token, có guard fail-closed khi token rỗng.

---

## 7. Danh sách defect

### P0 — chặn phát hành

---

#### **P0-1 · `/ma/{mã}` và `/nhan-dinh/{slug}` trả HTTP 500 — 2/3 trang chết**

**Tái hiện:**
```bash
curl -o /dev/null -w "%{http_code}\n" http://127.0.0.1:3000/ma/FPT            # 500
curl -o /dev/null -w "%{http_code}\n" http://127.0.0.1:3000/nhan-dinh/cmg-quy-2-2026   # 500
```
Log container `tenpoint-web`:
```
⨯ TypeError: Cannot read properties of undefined (reading 'slice')
  digest: '2749032436'
```

**Nguyên nhân gốc — lệch hợp đồng API hoàn toàn ở `/api/v1/tickers/{symbol}`:**

Backend trả (`backend/internal/http/dto.go:108-109`) — các trường ticker nằm phẳng ở cấp trên, kèm:
```go
RecentNews    []NewsItem `json:"recent_news"`
ResearchNotes []NoteItem `json:"research_notes"`
```
Xác nhận bằng response thật: `keys = ['article_count','company_name','exchange','in_vn30','recent_news','research_notes','sector','short_name','symbol']`

Frontend mong đợi (`frontend/src/lib/types.ts:141-145`):
```ts
export interface TickerDetail {
  ticker: Ticker;            // ❌ API không có 'ticker' (các trường nằm phẳng)
  news: NewsItem[];          // ❌ API tên là 'recent_news'
  notes: ResearchNoteSummary[];  // ❌ API tên là 'research_notes'
}
```

**Điểm nổ:**
- `frontend/src/app/nhan-dinh/[slug]/page.tsx:76` — `related?.data.news.slice(0,5)` → `news` là `undefined`, `?.` chỉ bảo vệ `related` chứ không bảo vệ `data.news`
- `frontend/src/app/ma/[symbol]/page.tsx:49` — `const { ticker, news, notes } = detail.data;` cả 3 đều `undefined`, nổ ở `:96` `news.length`

**Vì sao cổng tự động không bắt được:** `frontend/src/lib/api.ts:50` ép kiểu không kiểm chứng:
```ts
return (await response.json()) as T;
```
TypeScript tin lời khai; runtime thì không. Và vì `response?.data` vẫn truthy nên fallback fixture không bao giờ kích hoạt.

**Hệ quả:** mất MVP MUST #9 (trang `/ma/{mã}` — nền móng SEO theo BA3 W1) và MUST #10 (Research Note đúng ảnh A). Kéo theo FAIL: TC-UI-10, TC-UI-11, TC-SEO-05, và toàn bộ đối chiếu Ảnh A.

**Đề xuất:** sửa `types.ts` cho khớp DTO backend, và thay `as T` bằng một lớp kiểm chứng (zod hoặc type guard thủ công) ở `api.ts:50` để lỗi hợp đồng lộ ra lúc build/test thay vì lúc người dùng bấm vào.

---

#### **P0-2 · 36/36 tin VnEconomy dùng chung MỘT tóm tắt lạc đề — 19% kho tin hiển thị sai nội dung**

**Tái hiện (query thật trên DB đang chạy):**
```sql
select s.domain, count(*) as articles, count(distinct a.summary_md) as distinct_summaries
from articles a join sources s on s.id=a.source_id
where a.is_demo=false and a.status='published' group by s.domain order by articles desc;
```
```
tuoitre.vn             | 41 | 41      ← 1:1, bình thường
tinnhanhchungkhoan.vn  | 40 | 40      ← bình thường
cafef.vn               | 40 | 40      ← bình thường
vneconomy.vn           | 36 |  1      ← ❌ 36 bài, ĐÚNG 1 tóm tắt
vnexpress.net          | 27 | 27      ← bình thường
vietstock.vn           |  4 |  4
thanhnien.vn           |  2 |  2
```

**Nội dung sai như thế nào** — tiêu đề thật của 6 trong 36 bài:

```
10 || Từng vỡ nợ, Hy Lạp bây giờ đi vay với lãi suất còn thấp hơn cả Pháp
14 || Khối ngoại mua ròng ngay trước thềm Fed nâng lãi suất
18 || Vận hành thị trường tài sản mã hóa trong năm nay
22 || Blog chứng khoán: Dòng tiền "nín thở"
28 || Cổ phiếu vừa và nhỏ điều chỉnh mạnh, khối ngoại duy trì mua ròng
32 || Quỹ ngoại: Nâng hạng là chất xúc tác quan trọng cho thị trường
```

Tóm tắt **duy nhất** được phục vụ cho cả 36 bài trên:

> "Nhân Ngày Quốc tế Bảo vệ tầng ô-dôn năm **2026** và kỷ niệm 10 năm Bản sửa đổi, bổ sung Kigali được thông qua, Cục Biến đổi khí hậu, Bộ Nông nghiệp và Môi trường phối hợp với Hội Khoa học Kinh tế Việt Nam và Tạp chí Kinh tế Việt Nam/VnEconomy tổ chức Chương trình tọa đàm…"

Text này được cào từ **widget quảng bá ở sidebar** của VnEconomy (`[Trực tuyến]: Tọa đàm…`), không phải thân bài. Không liên quan gì tới bất kỳ tiêu đề nào.

**Nguyên nhân:** bóc tách thân bài (edge case E-05 "loại bỏ box Tin liên quan, quảng cáo") **hỏng hoàn toàn** với layout vneconomy.vn — selector bắt trúng khối sidebar thay vì `article body`.

**Vì sao validator không chặn:** các con số trong khối sidebar đó *có thật* trên trang, nên cổng đối chiếu số vẫn PASS. Xem ghi chú ở mục 5.2.

**Hệ quả:** người dùng đọc tóm tắt về hội thảo tầng ô-dôn khi bấm vào tin "Khối ngoại mua ròng trước thềm Fed nâng lãi suất". Phá huỷ trực tiếp thứ PRD nói sản phẩm bán bằng: **độ tin cậy của trích dẫn**.

---

#### **P0-3 · Precision tag mã chứng khoán ≈ 52% — ngưỡng PRD là ≥ 96%**

Mẫu 64 tag trên 20 bài (8 `primary`, 56 `mentioned`). Với mỗi tag: fetch bài gốc thật, grep **raw HTML** tìm ký hiệu mã và alias tên công ty.

| | |
|---|---|
| Tổng tag kiểm | **64** |
| Đúng | **33** |
| False positive | **31** |
| **Precision ước lượng** | **33/64 ≈ 52%** |
| theo mức | `primary` 6/8 = 75% · `mentioned` 27/56 = 48% |
| Ngưỡng PRD | **≥ 96%** — thiếu ~44 điểm |

**Bằng chứng tôi tự kiểm chứng lại (id 8):**
```
tickers: ACB,BID,BVB,CTG,HDB,LPB,MBB,MSB  (tất cả 'mentioned')
summary có nhắc mã nào? []
bài gốc: https://tinnhanhchungkhoan.vn/nhan-dinh-thi-truong-phien-...-post397703.html (88.465 bytes)
  ACB trong raw HTML: 0      BID: 0      BVB: 0      CTG: 0
  HDB: 2 (nằm ở sidebar "bài liên quan")  LPB: 0   MBB: 0   MSB: 0
```
7/8 mã xuất hiện **0 lần** trong bài. Câu duy nhất liên quan ngành ngân hàng là *"Sức nâng đỡ chủ yếu đến từ nhóm ngân hàng, năng lượng và công nghệ"* — không nêu tên công ty nào.

**Nguyên nhân chính — tier T3 (ngữ cảnh ngành) chạy vô điều kiện** (`tagger.go:188-193`). Edge case G-08 quy định T3 chỉ dùng khi bài **không nhắc mã nào**; thực tế nó chạy cho mọi bài. **29/31 false positive đến từ các bài không nêu tên công ty nào.**

Dấu vân tay ở mức toàn kho: cùng một bộ mã lặp lại trên các bài không liên quan —
`{CMG, FPT}` trên **8** bài (tập đoàn TQ đầu tư đường sắt, triển lãm in ấn, hồ sơ SEC của Meey Global, cuộc gặp Quảng Tây…); `{SSI, VCI, VND}` trên **7** bài; `{HPG, HSG, NKG, TIS}` trên **4**; `{POW, REE}` trên **3**.
**39/77 bài có tag (51%) không có mã `primary` nào** — toàn bộ là `mentioned`. **17 bài mang ≥4 mã đều ở mức `mentioned`** — chữ ký của T3.

**2 lỗi ở mức `primary` (ảnh hưởng trực tiếp feed người dùng):**
1. **id 26 — trùng ký hiệu.** Bài về SeABank. `VNM` gắn `primary`, nhưng `VNM` trong bài là *"VanEck Vectors Vietnam ETF (**VNM** ETF)"* — quỹ ETF Mỹ, không phải Vinamilk. "Vinamilk" xuất hiện 0 lần.
2. **id 6 — alias sai.** `MWG` gắn `primary`; MWG xuất hiện 0 lần. Bài nêu *"CTCP Đầu tư Điện Máy Xanh (**DMX** – HOSE)"*. Alias "Điện Máy Xanh" → MWG bắn nhầm, mã thật là DMX.

**Giới hạn của con số này (nêu rõ):** đây là **đánh giá thủ công trên mẫu tiện lợi 64 tag, không phải đo trên tập gán nhãn chuẩn**, không có khoảng tin cậy. Tôi tính CORRECT cho cả những lần nêu tên thoáng qua; nếu siết theo tiêu chí "có đáng đọc với người nắm mã này không" thì precision **xuống dưới 40%**. Lật ngược 4 phán đoán gây tranh cãi nhất (BID/VCB/MWG/CTG) cũng chỉ đưa 52% → ~58%.
**Hướng của kết luận thì không có gì phải nghi ngờ: precision nằm đâu đó trong khoảng 40–60% so với ngưỡng 96%.**

---

#### **P0-4 · Bộ lọc ngôn ngữ khuyến nghị đầu tư bị vượt qua bởi dấu câu — và đã lọt tin thật**

Vi phạm TC-VAL-10 · NUM-7 · PRD mục 11 (rủi ro "bị quy là tư vấn đầu tư không phép" — mức Cao) · edge case S-06.

**File:** `backend/internal/ingest/summarizer/validator.go:85-87`
```go
norm := " " + textutil.Normalize(strings.ReplaceAll(summaryMD, "*", "")) + " "
for _, p := range bannedPhrases {
    if strings.Contains(norm, " "+p+" ") || strings.Contains(norm, " "+p+",") {
```
`textutil.Normalize` (`normalize.go:40`) chỉ bỏ dấu, hạ chữ thường, gom khoảng trắng — **giữ nguyên dấu câu**. Nên biên từ chỉ chấp nhận khoảng trắng hoặc dấu phẩy phía sau. Mọi cụm cấm đứng **cuối câu** đều lọt.

**Tái hiện:** `summarizer.FindBannedPhrase("SSI khuyến nghị mua.")` → `("", false)` (đáng lẽ phải bắt được).

**Đo thật — lọt 14/16:**
```
BLOCKED      "Chuyên gia khuyến nghị nắm giữ cổ phiếu."
BLOCKED      "Chuyên gia khuyến nghị nắm giữ, theo báo cáo."
*** LEAK *** "Chuyên gia khuyến nghị nắm giữ."     *** LEAK *** "SSI khuyến nghị mua."
*** LEAK *** "VCSC khuyến nghị bán."               *** LEAK *** "Đây là cơ hội đầu tư."
*** LEAK *** "Nhà đầu tư có thể chốt lời."         *** LEAK *** "Cổ phiếu đang bắt đáy."
*** LEAK *** "Cổ phiếu có tiềm năng tăng giá."     *** LEAK *** "Nhà đầu tư có thể canh mua."
*** LEAK *** "Quan điểm: khuyến nghị nắm giữ."     *** LEAK *** "**Khuyến nghị nắm giữ.**"
*** LEAK *** "…nắm giữ;" / "…nắm giữ!" / "…nắm giữ?" / "…nắm giữ)"
TOTAL LEAKS: 14 / 16
```
Xác nhận end-to-end qua `SummarizeWithFallback`: tóm tắt kết thúc bằng `"...: khuyến nghị nắm giữ."` được **chấp nhận** với `provider=anthropic, calls=1, err=<nil>` — không retry, không fallback, sẽ được lưu và publish.

**Vì sao test hiện có không bắt:** cả 5 fixture trong `validator_test.go:109-114` và `fallback_test.go:124` đều đặt cụm cấm ở **giữa câu**, có khoảng trắng phía sau.

**Đây không còn là lỗ hổng lý thuyết — đã lọt ra tin thật đang publish.** Tin `id 8` ("Nhận định thị trường phiên giao dịch ngày 17/9"):

> "Trong bối cảnh này, nhà đầu tư ngắn hạn **nên duy trì tỷ trọng cổ phiếu ở mức trung bình, tránh mua đuổi** và tập trung giao dịch theo vùng hỗ trợ 1.810-1.820 điểm và kháng cự 1.830-1.840 điểm."

Đây là lời khuyên phân bổ danh mục có thể hành động được — đúng loại nghiệp vụ mà PRD mục 4 nói **chỉ CTCK được cấp phép** mới được làm. Câu này copy nguyên văn từ báo cáo môi giới, nhưng sản phẩm **tái xuất bản như tóm tắt của chính mình**.

**Lỗ hổng thứ hai, độc lập — thiếu độ phủ:** cụm `"nên duy trì"` và `"tránh mua đuổi"` **hoàn toàn không có** trong `bannedPhrases` (`validator.go:74-81`). Nên kể cả khi sửa lỗi biên từ, câu này vẫn lọt. Danh sách hiện chỉ có 13 cụm tư vấn, không phủ hết dạng "nhà đầu tư nên + <hành động>".

**Đề xuất:** (a) dựng biên từ bằng `textutil.NormalizeWords` (`normalize.go:52`) — đã có sẵn, biến mọi ký tự không phải chữ/số thành khoảng trắng; thử nghiệm bắt được 7/7 mẫu đang lọt. (b) Mở rộng danh sách theo **mẫu** (`nhà đầu tư nên …`, `nên (mua|bán|duy trì|giảm|tăng) tỷ trọng`) thay vì liệt kê cứng từng cụm.

---

#### **P0-5 · Healthcheck `api` không bao giờ pass → compose production không khởi động được `web`**

Vi phạm TC-DEPLOY-02.

**Tái hiện:**
```bash
$ docker exec tenpoint-api sh -c 'wget -q --spider http://127.0.0.1:8080/healthz; echo $?'
wget_exit=8

$ curl -I http://127.0.0.1:8090/healthz
HTTP/1.1 405 Method Not Allowed        ← HEAD không được nhận

$ curl http://127.0.0.1:8090/healthz
{"status":"ok"}                         ← GET thì bình thường

$ docker inspect tenpoint-api --format '{{.State.Health.Status}} {{.State.Health.FailingStreak}}'
unhealthy 6
```

`wget --spider` gửi **HEAD**; router chỉ đăng ký **GET** cho `/healthz` → 405 → `wget` exit 8 → healthcheck fail vĩnh viễn. Log đầy `"method":"HEAD","path":"/healthz","status":405` mỗi 5 giây.

**Hệ quả nghiêm trọng ở production** — compose production (rendered):
```yaml
web:
  depends_on:
    api:
      condition: service_healthy     # ← api không bao giờ healthy
```
so với dev override:
```yaml
web:
  depends_on:
    api:
      condition: service_started     # ← chỉ vì dòng này mà QA test được
```

⇒ Chạy `make up` / `docker compose -f deploy/docker-compose.yml up -d` trên máy sạch: `web` **không bao giờ khởi động**, site chết hoàn toàn. Toàn bộ quá trình test này chỉ chạy được nhờ dev override nới điều kiện.

**Đề xuất:** đổi healthcheck sang `wget -qO- .../healthz` (dùng GET), hoặc đăng ký HEAD cho `/healthz` ở router. Nên làm cả hai.

---

### P1 — sửa trước khi phát hành

**P1-1 · Chế độ fixture render tin mẫu thành link bấm được, 4/5 dẫn tới 404**
Đúng thứ mà edge-case mục 0 đặt ra để loại bỏ: *"Không bao giờ để người dùng bấm vào một URL bịa."*
Tái hiện: `docker stop tenpoint-api` rồi mở trang chủ → 200, hiện 5 tin mẫu **dưới dạng `<a href="/r/...">`**, không có chữ "Dữ liệu mẫu".
```
fixture codes: a7Kx2p b3Qm9d c8Tz1k d2Wr7v e6Yn4s
DB demo codes: a7Kx2p b3Qw9m c8Zt4r d2Lp6v e5Nh1s
/r/a7Kx2p → 410    /r/b3Qm9d → 404   /r/c8Tz1k → 404
/r/d2Wr7v → 404    /r/e6Yn4s → 404
```
Nguyên nhân: `frontend/src/lib/fixtures.ts:17-114` — cả 5 item **thiếu `is_demo: true`**, nên guard `SourceLink.tsx:29-41` không kích hoạt. Mã short link trong fixture cũng lệch với DB.
Sửa: thêm `is_demo: true` (và đồng bộ code) vào 5 fixture.

**P1-2 · `NEXT_PUBLIC_SITE_URL` không phải build arg → mọi URL tuyệt đối sai host ở production**
`frontend/src/lib/site.ts:8-10` đọc `NEXT_PUBLIC_*`, giá trị bị "nướng" vào bundle lúc **build**. `deploy/Dockerfile.web:38` chỉ khai báo `ARG API_BASE_URL`; compose cấp `NEXT_PUBLIC_SITE_URL` lúc **runtime** — quá muộn cho route tĩnh.
Bằng chứng runtime, hai nguồn lệch nhau:
```
robots.txt (static, build-time):  Host: http://localhost:3000
                                  Sitemap: http://localhost:3000/sitemap.xml
sitemap.xml (dynamic, runtime):   <loc>http://localhost/</loc>
canonical:                        <link rel="canonical" href="http://localhost"
```
Sửa: thêm `ARG`/`ENV NEXT_PUBLIC_SITE_URL` trước `RUN npm run build` và truyền qua `build.args`.

**P1-3 · Contrast `--ink-muted` 4.28:1 < 4.5:1 (WCAG AA)** — TC-A11Y-03 FAIL
Đo thật: `#6B7580` trên `#F4F5F2` = **4.28:1**. Dùng cho `.digest-date` và `.digest-type` (13.5px, là text thường nên ngưỡng là 4.5:1).
Tài liệu `design-system.md` mục 2.4 ghi "~4.6:1 ✅ AA" — **sai**, cần sửa cùng lúc.
Ngoài ra `--ink-faint` (~2.5:1, mục 2.4 đã ghi "chỉ non-text") đang được dùng làm text ở `globals.css:494-501` (`.source-dead`), `:223-225` (`::placeholder`), `TickerList.tsx:31`.

**P1-4 · Lỗi cấp item bị nuốt — không phân biệt được run sạch với run hỏng** (US-6.1 AC2)
Run 1 ghi `errors=0, articles_rejected=0` nhưng mất 3/190 bài; chỉ hiện ở log WARN.
`backend/internal/ingest/pipeline.go:277-291`: khi `ingestItem` trả lỗi thì `continue` ở `:280`, thoát **trước** `switch` ở `:289`, nên không bao giờ chạm `stat.rejected++`. `run.Errors++` (`:171`) chỉ tăng khi lỗi **cấp nguồn**.
Thêm nữa, `articles_rejected` là bộ đếm chết: cả 3 đường sinh `outcomeRejected` (`:331` allowlist, `:349` extract, `:448` cửa sổ thời gian) đều trả lỗi non-nil nên đều bị `continue` chặn trước. Cột này chỉ đếm được held/gone — gây hiểu nhầm cho người vận hành.

**P1-5 · Rate limiter `/r/{code}` bị vô hiệu bằng một header**
`backend/internal/http/handlers_redirect.go:138-144` — `clientIP()` lấy phần tử **trái nhất** của `X-Forwarded-For`, không kiểm trusted proxy. Caddy *append* vào XFF do client gửi, nên giá trị trái nhất là do kẻ tấn công kiểm soát.
Đo thật:
```
70 request cùng IP, không XFF     → 429 count = 10   (limit 60/phút hoạt động)
70 request, XFF xoay vòng          → 429 count = 0    ← bypass hoàn toàn
```
`Caddyfile:93` đã set `X-Real-IP {remote_host}` nhưng code Go không đọc. Đây là phòng tuyến duy nhất chống dò mã short link (BA2 §7.4).

**P1-6 · Banned-phrase chặn nhầm tiếng Việt thông thường, và là chặn vĩnh viễn**
`validator.go:74-81` trộn 2 nhóm: 13 cụm đúng là ngôn ngữ tư vấn, và 11 cụm chỉ là **văn phong** — `tom lai`, `nhin chung`, `boc hoi`, `chung toi`, `chung ta`.
Quan sát thật: 1 bài vnexpress bị loại vì `banned_phrase (chung ta)` — "chúng ta" nghĩa là "we/us", không phải khuyến nghị đầu tư.
Nghiêm trọng hơn vì provider mặc định là `extractive` (lấy nguyên câu từ bài gốc), nên bộ lọc đang áp lên văn của nhà báo; và khi primary = backup, `fallback.go` trả lỗi luôn — **không retry, không fallback**. Bài bị loại vĩnh viễn, lặp lại y hệt ở cả 2 run.

**P1-7 · `negative_aliases` được seed nhưng không bao giờ được đọc** (edge case G-04)
Cột có ở `migrations/0005_edge_cases.sql:53-57` nhưng **0 tham chiếu trong Go**; `store/tickers.go:56-63` không SELECT, `tagger.Ticker` không có field.
False positive đã chứng minh: `FPT Telecom báo lãi` → `FPT/primary` (đúng phải là FOX) · `Masan Consumer chia cổ tức` → `MSN/primary` (đúng phải là MCH) · `Masan MEATLife` → `MML` **+ `MSN`**.
(TC-TAG-12 `Vinhomes → VHM` vẫn PASS, nhưng nhờ khớp biên từ chứ không nhờ cơ chế này.)

### P2 — ghi nhận cho v1.1

| ID | Mô tả | Vị trí |
|---|---|---|
| P2-1 | **TC-VAL-09 FAIL** — số cụm bold 1–3 không được enforce ở backend. `BoldSpans` có tính nhưng không ai kiểm. Chỉ có `console.warn` dev-only ở FE | `validator.go:143-154`, `SummaryText.tsx:17-25` |
| P2-2 | **TC-TAG-17 FAIL** — mã trong URL hiện dạng text vẫn bị gắn. `Nguồn: https://cafef.vn/tin-tuc/FPT-abc.html` → `FPT/mentioned`. Phòng vệ hiện tại chỉ là tình cờ ở `extractor.go:224`, bị bypass bởi fallback `:241` | `tagger.go:261-300` |
| P2-3 | **Chất lượng cụm bold kém — 73,7% vi phạm B-3** (14/19 cụm, đo trên 20 bài ở mục 5.1). Heuristic chọn dãy chữ số đầu tiên thay vì cụm số+đơn vị: năm trần ×6, cắt đôi ngày ×3, mất đơn vị ×3, cắt cụm ghép ×2. Một tin (id 29) có **0 cụm bold**, vi phạm NUM-3. *Nâng lên P1 nếu tính cả việc `**2026**` in đậm trên sản phẩm tài chính đọc như một chỉ số* | `extractive.go:180-234` |
| P2-4 | Copy bản tin in ngoặc rỗng `[]` khi tin không có mã — 8/20 khối | `DigestView.tsx:98-107` |
| P2-5 | `/tickers/FPT` trả `article_count: 0` trong khi `recent_news` có 20 mục | `handlers_tickers.go` |
| P2-6 | Tier T3 (ngữ cảnh ngành) chạy **vô điều kiện**, đáng lẽ chỉ khi bài không nhắc mã nào (G-08). `CMG trúng thầu…` → `CMG/primary` + `FPT/mentioned` thừa | `tagger.go:188-193` |
| P2-7 | Cắt còn 8 mã theo **thứ tự alphabet**, im lặng. 20 mã → giữ `ACB BID CTG HDB LPB MBB SHB STB`, rơi FPT/HPG/VCB. Cảnh báo G-11 (>15 mã) **không bao giờ chạy được** vì đã cắt còn 8 trước đó | `tagger.go:204-215` |
| P2-8 | `crawl_run_sources.truncated` và `skipped_by_robots` **không bao giờ được ghi** (INSERT thiếu cột). Log có 6 lần "feed truncated" nhưng DB vẫn `truncated=f` | `store/runs.go:42-45` |
| P2-9 | Stoplist thiếu UPC, VAT, EVN, CSR so với chính G-01. Rủi ro tiềm ẩn: `HCM`, `BTC`, `SSC`, `BOT`, `PPP` là mã thật, sẽ vĩnh viễn không gắn được khi mở rộng từ điển | `stoplist.go:10-38` |
| P2-10 | Không có `unknown_ticker_candidate` (G-09); `delisted_at` có cột nhưng `store/tickers.go:61` lọc `status='active'` nên mã huỷ niêm yết biến mất khỏi từ điển (G-10) | |
| P2-11 | Không re-validate `target_url` lúc redirect — allowlist chỉ chạy lúc ingest. Một lần sửa DB thủ công là thành open redirect | `handlers_redirect.go:63` |
| P2-12 | Caddy ghi đè `Referrer-Policy: no-referrer` của handler redirect thành `strict-origin-when-cross-origin`. Đo thật: direct api = `no-referrer`, qua Caddy = `strict-origin-when-cross-origin` | `Caddyfile:47` |
| P2-13 | CORS `*` áp cả route admin (`AllowedMethods` có POST, `AllowedHeaders` có Authorization) | `router.go:68-74`, `:96` |
| P2-14 | CSP `script-src 'unsafe-inline'` — CSP không còn là lớp phòng vệ XSS. Đã được ghi chú và có hướng sửa sẵn trong Caddyfile | `Caddyfile:72` |
| P2-15 | Mật khẩu DB mặc định `tenpoint_dev_only_change_me` nằm trong compose production dạng `:-default`. Thiếu `.env` là boot im lặng bằng mật khẩu công khai | `docker-compose.yml:39,100` |
| P2-16 | I-04 chưa hiện thực: không có job dọn `crawl_runs` treo, không có cột `status`. Ghi chú: ERD `architecture.md` mục 4 **cũng không** định nghĩa `status` → mâu thuẫn giữa 2 tài liệu, không phải lệch code-vs-ERD | |
| P2-17 | `FinishRun` gọi thẳng, không trong `defer` (trái I-02). 3 early return ở `pipeline.go:135,139,143` thoát sau `StartRun` mà không đóng run | `pipeline.go:186` |
| P2-18 | `last_crawl_at` dùng `MAX(finished_at) WHERE finished_at IS NOT NULL` → không phân biệt "chưa crawl lần nào" với "đã chạy N lần, chưa lần nào xong". Hệ thống hỏng kinh niên vẫn báo `is_stale=false`, HTTP 200 | `store/runs.go:53-61` |
| P2-19 | Thứ tự DOM ≠ thứ tự thị giác ở mobile: `Loại tin` thứ 5 trong DOM nhưng thứ 2 khi nhìn (WCAG 1.3.2). Thị giác vẫn đúng TC-RWD-02 | `DigestTable.tsx:51-77` vs `globals.css:592-596` |
| P2-20 | Toast và vùng `aria-live` được mount có điều kiện → screen reader có thể không đọc | `DigestView.tsx:178-182`, `FilterBar.tsx:256` |
| P2-21 | Copy bản tin vẫn in `Nguồn: {origin}{short_link}` cho cả tin demo/chết, mâu thuẫn với chính `SourceLink` | `DigestView.tsx:104` |

### Defect tài liệu

**DOC-1 (P1 tài liệu) · `test-plan.md` mục 4.6 đã lỗi thời so với design-system v2.**
TC-UI-03 yêu cầu nền `#F7F5F0` và TC-UI-02 yêu cầu hairline `#DDD8CE` — đều là giá trị **v1**. Design-system v2 mục 0 và mục 10 **cấm đích danh** `#F7F5F0` (là bảng màu AI mặc định). Chạy đúng chữ trong test-plan sẽ cho FAIL giả trên code đúng.
QA đã nghiệm thu theo **design-system v2** (`#F4F5F2` / `#D5D8D2`) và ghi TC-UI-02/03 là PASS. Cần sửa test-plan.

**DOC-2 (P1 tài liệu) · `design-system.md` mục 2.4 ghi sai tỉ lệ contrast.** `--ink-muted` ghi "~4.6:1 ✅ AA", đo thật 4.28:1 (FAIL AA).

**DOC-3 (P2 tài liệu) · Trang lỗi 410 dùng màu bị cấm.** `handlers_redirect.go` render trang 410 với `background:#F7F5F0` — đúng màu kem mà design-system mục 10 cấm.

---

## 8. Đối chiếu bảng ngưỡng PRD mục 10

| Nhóm | Chỉ tiêu | Ngưỡng | Đo được | Kết luận |
|---|---|---|---|---|
| **Độ chính xác số liệu** | Tin publish có số bịa | **0** | **0** — kiểm 20/20 bài, mọi con số truy vết được về bài gốc | ✅ |
| Tag mã CK | Precision | ≥ 96% | **≈ 52%** (33/64 tag, mẫu thủ công) | ❌ P0-3 |
| Tag mã CK | Recall | ≥ 85% | **Không đo được** (không có tập gán nhãn) | ⚠️ BLOCKED |
| Tóm tắt | Độ dài 60–150, reject >180 | | **20/20 đạt** (81–131 từ) | ✅ |
| Tóm tắt | Cụm bold đúng 1–3 | | 19/20 có 1 cụm; id 29 có **0** cụm. Không enforce ở backend (P2-1). Chất lượng cụm vi phạm B-3 73,7% | ❌ |
| Tóm tắt | 0 cụm khuyến nghị đầu tư | 0 | **1/20 vi phạm** (id 8) + bộ lọc lọt 14/16 mẫu | ❌ P0-4 |
| **Nội dung tóm tắt** | Đúng bài | (ngầm định) | **36/190 = 19% sai hoàn toàn** | ❌ P0-2 |
| Pipeline | 1 run | < 10 phút | **2m53s** | ✅ |
| Pipeline | 1 nguồn lỗi không làm hỏng run | | Run hoàn tất dù có item lỗi; nhưng lỗi bị nuốt (P1-4). Chưa thử nguồn trả 500 | ⚠️ BLOCKED |
| API | p95 `/api/v1/news` | < 200ms | **2–9ms** (log: `duration_ms 2,4,5,9`) | ✅ |
| Short link | Thời gian redirect | < 50ms | **0–2ms** (log `duration_ms:0`) | ✅ |
| Short link | Open redirect | **0** | **0** — 12/12 input tấn công bị chặn | ✅ |
| FE | LCP mobile 4G | < 2.0s | **Không đo** (chưa chạy Lighthouse throttle) | ⚠️ BLOCKED |
| FE | Bảng mobile không scroll ngang | không | `scrollWidth<=innerWidth` ở 375/390/900/1440/1920 | ✅ |
| A11y | WCAG 2.2 AA theo mục 8 | đạt | 8/10 — FAIL contrast + target size | ❌ |
| Giao diện | Khớp ảnh tham chiếu | 1:1 | Ảnh B ✅ · Ảnh A ❌ (500) | ❌ |

### Ghi chú về cách đọc con số precision 52%

Con số này là **đánh giá thủ công 64 tag trên mẫu tiện lợi**, không phải đo trên tập gán nhãn chuẩn — chi tiết giới hạn ghi ở P0-3. Nó **không** đủ tư cách làm số nghiệm thu chính thức, nhưng **đủ để kết luận là KHÔNG ĐẠT**: khoảng cách tới ngưỡng 96% là ~44 điểm, lớn hơn nhiều so với mọi sai số của phương pháp lấy mẫu, và nguyên nhân gốc (T3 chạy vô điều kiện) đã được xác định cụ thể ở mức code.

### Vì sao KHÔNG kết luận được recall

- Từ điển chỉ có **56 mã / 164 alias** (`migrations/0003_seed_tickers.sql`), trên ~1.600 mã niêm yết thực tế (**~3,5% độ phủ**).
- Mọi mã ngoài từ điển bị bỏ im lặng (P2-10), nên recall bị chặn cứng bởi độ phủ từ điển chứ không phản ánh chất lượng thuật toán.
- Test-plan mục 4.2 yêu cầu "50 bài mẫu đã gán nhãn thủ công" — **tập này không tồn tại trong repo**, không có harness đánh giá offline.

Mọi con số recall tính hôm nay sẽ mô tả **độ phủ từ điển**, không mô tả hệ thống. Không đưa ra con số.

**Tín hiệu định lượng liên quan:** trên 20 hàng đầu trang chủ với dữ liệu thật, **8/20 (40%) không gắn được mã nào** (cột Mã CK hiện `-`), và **39/77 bài có tag (51%) không có mã `primary` nào**. Với sản phẩm mà PRD mục 2 tuyên bố *"cột Mã CK đứng thứ hai, trước cả nội dung — đó là tuyên ngôn sản phẩm"*, cần corpus gán nhãn để lượng hoá.

---

## 9. Case BLOCKED — nêu rõ lý do

| ID | Lý do không chạy được |
|---|---|
| TC-PIPE-01 | Không nguồn nào trả HTTP 500 trong 2 run quan sát được. Nửa "các nguồn khác vẫn chạy xong" đã PASS; nửa "ghi vào `crawl_run_sources.error_message`" chưa được kích hoạt thật. Cần cố ý cấu hình một nguồn hỏng |
| TC-PIPE-03 | Không quan sát được 2 bài khác nguồn cùng sự kiện trong dữ liệu thu được |
| TC-PIPE-04 | Như trên — không dựng được cặp đối chứng |
| TC-PIPE-08 | Cơ chế có thật (`store/revisions.go:242-253` tăng/reset `consecutive_empty_runs`, `pipeline.go:265-270` log ERROR khi streak≥2) nhưng không nguồn nào trả 0 bài. Live: cả 7 nguồn `consecutive_empty_runs=0` |
| TC-PIPE-09 | Không nguồn nào bị `robots.txt` chặn đường dẫn trong run thật (`grep -c robots` trong log = 0) |
| TC-UI-09 | Không ép được `last_crawl_at` cũ hơn 12h mà không sửa DB |
| TC-SEO-05 | JSON-LD nằm trên trang research note — trang này 500 (P0-1) |
| PRD: LCP mobile 4G | Chưa chạy Lighthouse với throttling 4G |

---

## 10. Điểm làm tốt (ghi nhận)

Để bức tranh cân bằng — đây là những phần làm chắc tay:

- **Bảo mật short link xuất sắc.** 12/12 input tấn công bị chặn, khớp theo nhãn domain chứ không theo chuỗi, chặn cả `javascript:`, `data:`, protocol-relative, userinfo-confusion (`cafef.vn@evil.com`). Allowlist được gọi ở cả 2 đường ghi. Token admin dùng `subtle.ConstantTimeCompare`.
- **Sạch tuyệt đối 26/26 color token** so với design-system v2, cả light lẫn dark, đo bằng computed style. Không có màu kem AI bị cấm, không gradient, **không một shadow nào**, không zebra, không border dọc, một hệ bo góc duy nhất.
- **Validator số liệu bắt đúng 2 case tồn vong.** Bịa số `2.323→2.523` FAIL; làm tròn `29,91→30` FAIL. Đây là phần quan trọng nhất của sản phẩm tài chính và nó hoạt động.
- **Link nguồn đã sửa đúng phản hồi khách hàng.** 5/5 URL đích trả 200, `<title>` khớp `article_title`; `title`/`aria-label` ghi rõ tiêu đề bài gốc.
- **Không scroll ngang ở cả 5 viewport**, container khoá đúng 1180px, bảng chuyển card đúng thứ tự.
- **0 console error/warning** trên trang chủ.
- **Tách liveness/readiness đúng chuẩn**: DB chết → `/readyz` 503 tiếng Việt, `/healthz` vẫn 200, tự phục hồi.
- **API nhanh**: p95 quan sát 2–9ms, redirect 0–2ms — vượt xa ngưỡng PRD.
- **Bootstrap T1/T2/T3 hoạt động đúng thiết kế**: tự ingest thật lúc khởi động, đánh dấu demo, ẩn demo khi có tin thật.

---

## 11. Việc cần làm để nghiệm thu lại

**Bắt buộc (P0) — theo thứ tự công sức tăng dần:**

1. **P0-5 (nhỏ nhất, làm ngay).** Đổi healthcheck `api` sang GET (và/hoặc đăng ký HEAD ở router). Nghiệm thu: `docker compose -f deploy/docker-compose.yml up -d` trên máy sạch cho **4/4 healthy**.
2. **P0-1.** Sửa `TickerDetail` (`frontend/src/lib/types.ts:141-145`) khớp DTO backend; thay `as T` (`api.ts:50`) bằng kiểm chứng runtime. Thêm smoke test HTTP cho **cả 3 trang** vào cổng tự động — đây là lỗ hổng quy trình, không chỉ là 1 bug.
3. **P0-4.** Sửa biên từ (`validator.go:85-87` → `NormalizeWords`) **và** mở rộng danh sách cấm theo mẫu. Thêm test cụm cấm ở cuối câu. Gỡ tin id 8 khỏi publish.
4. **P0-2.** Sửa bóc tách thân bài cho `vneconomy.vn` (đang bắt trúng widget sidebar). **Đồng thời bổ sung một cổng kiểm mới**: phát hiện tóm tắt trùng lặp giữa các bài khác nhau — hiện tại 36 bài dùng chung 1 tóm tắt mà không hệ thống nào kêu. Đây là cổng rẻ và bắt được cả lớp lỗi này về sau.
5. **P0-3 (nặng nhất).** Giới hạn tier T3 đúng theo G-08 (chỉ chạy khi bài không nhắc mã nào) — riêng việc này xử lý 29/31 false positive đã tìm thấy. Sau đó xử lý trùng ký hiệu (VNM ETF) và alias sai (Điện Máy Xanh → DMX). **Phải đo lại trên tập gán nhãn trước khi nghiệm thu.**

**Nên làm trước phát hành (P1):** P1-1 → P1-7 ở mục 7.

**Cần có để nghiệm thu chỉ tiêu PRD:**
- **Corpus 50 bài gán nhãn thủ công** (test-plan mục 4.2 đã yêu cầu, chưa ai làm) — không có tập này thì precision/recall **không thể nghiệm thu** một cách có căn cứ, lần này hay lần sau.
- Chạy Lighthouse throttle 4G cho LCP.
- Dựng 1 nguồn hỏng cố ý để mở khoá TC-PIPE-01/08/09.

**Sửa tài liệu:** DOC-1 (test-plan còn màu v1 đã bị cấm), DOC-2 (contrast ghi sai), DOC-3 (trang 410 dùng màu cấm).

---

## 12. Ghi chú về môi trường test

- Hai lỗi P0 nội dung (P0-2, P0-3) **chỉ lộ ra khi chạy với dữ liệu thật**. Với 5 tin seed demo, mọi thứ trông hoàn hảo. Nghiệm thu lần sau bắt buộc phải chạy ingest thật rồi mới kiểm.
- `fetch` nguồn `tinnhanhchungkhoan.vn` cần `curl --compressed`, nếu không sẽ báo động giả số bịa.
- Cổng tự động hiện tại (`build`/`vet`/`test`/`tsc`/`lint`/`next build`) **sạch 100% nhưng không bắt được 1 lỗi P0 nào** trong số 5 lỗi. Khuyến nghị bổ sung vào CI: smoke test HTTP 3 trang · kiểm trùng `summary_md` · kiểm `4/4 healthy` sau `compose up`.

---

## 13. Bằng chứng kèm theo

Ảnh chụp màn hình trong `docs/05-qa/evidence/`:

| File | Nội dung |
|---|---|
| `qa-reference-compare.png` | Trang chủ full-page ở chế độ 5 tin mẫu — dùng để đối chiếu 1:1 với **Ảnh B** |
| `qa-home-1440-light.png` | Trang chủ 1440px, dữ liệu thật |
| `qa-home-390-light.png` | Trang chủ 390px — bảng đã thành card, không scroll ngang |
| `qa-mobile-card.png` | Một card ở 390px — xác nhận thứ tự trường và không lặp `Loại tin` |

Không có ảnh của **Ảnh A** (Research Note) vì trang trả HTTP 500 — xem P0-1.

### Dọn dẹp

```
$ docker compose -f deploy/docker-compose.yml down -v
 Container tenpoint-api/db/web/caddy  Removed
 Volume tenpoint_pgdata / caddy_data / caddy_config  Removed
 Network tenpoint  Removed
```

Đã xoá `deploy/.env` (chứa token test entropy thấp — tạo lại từ `deploy/.env.example`). Đã kiểm tra **không còn file probe/test tạm nào** trong `backend/` hay `frontend/src/`; chạy lại `go build` + `go vet` + `gofmt -l` + `go test ./...` sau khi dọn: **vẫn sạch, 10/10 package ok**. Không sửa một dòng code sản phẩm nào.

---

*Báo cáo lập bởi QA. Mọi kết quả PASS đều kèm bằng chứng thật; case không chạy được ghi BLOCKED chứ không suy đoán. Không có con số precision/recall nào được bịa — chỗ không đo được đã ghi rõ lý do.*
