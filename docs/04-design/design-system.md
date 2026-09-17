# TenPoint — Design System v2

> **v2 thay toàn bộ bảng màu của v1.** Lý do ở mục 0.
> Nguồn: 2 ảnh tham chiếu (`docs/00-inputs/reference-images.md`) + skill anti-slop `design-taste-frontend`.

---

## 0. Vì sao phải làm lại bảng màu (đọc trước)

**Design Read:** *editorial data-product (bảng digest + research note) cho nhà đầu tư cá nhân và môi giới Việt Nam, ngôn ngữ báo chí tài chính, dựng bằng CSS tokens + Tailwind, motion tiết chế.*

**Dials:** `DESIGN_VARIANCE 4` · `MOTION_INTENSITY 3` · `VISUAL_DENSITY 7`
Đây là sản phẩm **để đọc nhanh lúc 6:45 sáng**, không phải landing page. Bố cục đối xứng phục vụ việc quét mắt; motion cao sẽ làm chậm người dùng; mật độ dữ liệu cao là đúng bản chất.

### Bảng màu v1 đã sai ở đâu

v1 dùng `--paper: #F7F5F0` (kem ấm) + mực đen `#111111` + accent xanh rêu `#1F4D3D`.

Kiểm tra lại với danh sách chống-AI-slop: nhóm màu **"warm paper / cream / bone"** (`#f5f1ea`, `#f7f5f1`, `#fbf8f1`, `#efeae0`...) cùng accent **brass/clay/ochre** và text **espresso near-black** là **bảng màu mặc định mà mô hình AI nào cũng sinh ra** cho brief kiểu "editorial / premium / artisan". `#F7F5F0` nằm chính giữa nhóm đó.

Hệ quả: trang trông "có gu" nhưng **không có bản sắc** — giống hàng nghìn trang AI khác. Đúng như phản hồi của khách hàng.

### Hướng thay thế

Chọn hệ **"Bản Tin"** — cool bone + navy editorial + dark mode đầy đủ. Ba lý do:

1. **Cool bone thay cream.** Sắc lạnh nhẹ gợi giấy báo in thật, thoát hẳn nhóm kem-nâu bị bão hoà.
2. **Navy làm accent, không phải xanh lá hay đỏ.** Bắt buộc: trong sản phẩm chứng khoán Việt Nam, **xanh lá = tăng, đỏ = giảm, tím = trần, xanh lơ = sàn, vàng = tham chiếu**. Năm màu đó đã có nghĩa cố định, brand **không được** lấn vào. Navy là màu duy nhất còn trống mà vẫn có chiều sâu — và cũng là màu của báo chí tài chính lâu đời.
3. **Dark mode là bắt buộc**, không phải "để v2". Nhà đầu tư đọc lúc 5:30 sáng và 22:00 tối.

---

## 1. Ba nguyên tắc thị giác (giữ từ v1)

1. **Chữ là giao diện.** Không card, không shadow, không bo góc. Phân cấp bằng size / weight / case / spacing.
2. **Số liệu là điểm nhấn duy nhất được phép.** Nhấn bằng **font-weight**, không bằng màu.
3. **Đường kẻ mảnh, không khối màu.** Hairline ngăn hàng. Không zebra, không border dọc, không header có nền.

---

## 2. Color tokens — hệ "Bản Tin"

### 2.1 Light mode

```css
:root {
  /* Surface — cool bone, KHÔNG phải cream */
  --surface:          #F4F5F2;
  --surface-raised:   #FCFCFB;
  --surface-sunken:   #E7E9E4;

  /* Ink — cool near-black, không dùng #000 */
  --ink:              #14181B;
  --ink-2:            #39424A;
  --ink-muted:        #666F79;   /* sửa 16/09: xem mục 2.4 */
  --ink-faint:        #9BA3AD;

  /* Rules */
  --rule:             #D5D8D2;
  --rule-strong:      #14181B;

  /* Accent — MỘT màu duy nhất, navy editorial */
  --accent:           #1E4B8F;
  --accent-hover:     #163A70;
  --accent-soft:      #E3EAF4;

  /* Focus */
  --focus:            #1E4B8F;
}
```

### 2.2 Dark mode

```css
@media (prefers-color-scheme: dark) {
  :root {
    --surface:        #101317;
    --surface-raised: #171B21;
    --surface-sunken: #0B0E11;

    --ink:            #E9ECEF;
    --ink-2:          #B4BCC6;
    --ink-muted:      #7E8894;
    --ink-faint:      #5A636E;

    --rule:           #262C34;
    --rule-strong:    #E9ECEF;

    --accent:         #6FA0F0;
    --accent-hover:   #96BAF6;
    --accent-soft:    #16233A;

    --focus:          #6FA0F0;
  }
}
```

`:root[data-theme="dark"]` và `:root[data-theme="light"]` phải ghi đè được cả hai chiều (cho nút chuyển chế độ thủ công).

### 2.3 Màu giá — chuẩn thị trường chứng khoán Việt Nam

**Đây là quy ước cứng của thị trường VN. Không được đổi vì lý do thẩm mỹ.**

| Ý nghĩa | Light | Dark | Dùng ở đâu |
|---|---|---|---|
| **Tăng** | `#0F7A43` | `#42C47E` | Chỉ ở ô biến động giá |
| **Giảm** | `#C02A1D` | `#FF7A6B` | Chỉ ở ô biến động giá |
| **Trần** | `#8A2BB0` | `#D08BEF` | Chỉ ở ô biến động giá |
| **Sàn** | `#0E7490` | `#5BC8E8` | Chỉ ở ô biến động giá |
| **Tham chiếu** | `#A07800` | `#E0B83C` | Chỉ ở ô biến động giá |

> **Quy tắc cứng:** năm màu này **chỉ** dùng cho dữ liệu giá. Không dùng cho nhãn, nút, badge, trang trí. Nếu đưa xanh lá vào một nút bất kỳ, người dùng Việt sẽ đọc nhầm là "tăng".
>
> **Không truyền đạt chỉ bằng màu** — luôn kèm dấu `+` / `−` (WCAG 1.4.1).

### 2.4 Kiểm chứng contrast (WCAG 2.2 AA)

| Cặp | Tỉ lệ | Kết luận |
|---|---|---|
| `--ink` `#14181B` / `--surface` `#F4F5F2` | ~15.9:1 | ✅ AAA |
| `--ink-2` `#39424A` / `--surface` | ~9.2:1 | ✅ AAA |
| `--ink-muted` `#666F79` / `--surface` | **4.66:1** | ✅ AA |
| `--ink-faint` `#9BA3AD` / `--surface` | ~2.5:1 | ⚠️ **chỉ non-text** |
| `--accent` `#1E4B8F` / `--surface` | ~8.1:1 | ✅ AAA |
| Tăng `#0F7A43` / `--surface` | ~4.9:1 | ✅ AA |
| Giảm `#C02A1D` / `--surface` | ~5.6:1 | ✅ AA |
| **Dark:** `--ink` `#E9ECEF` / `--surface` `#101317` | ~14.8:1 | ✅ AAA |
| **Dark:** `--ink-muted` `#7E8894` / `--surface` | ~4.7:1 | ✅ AA |
| **Dark:** `--accent` `#6FA0F0` / `--surface` | ~7.4:1 | ✅ AAA |
| **Dark:** Tăng `#42C47E` / `--surface` | ~8.6:1 | ✅ AAA |
| **Dark:** Giảm `#FF7A6B` / `--surface` | ~7.2:1 | ✅ AAA |

**Theme Lock:** toàn trang một chế độ. Không có section nào đảo màu giữa chừng.
**Color Consistency Lock:** đúng **một** accent (`--accent`) cho toàn bộ trang. Không có nút xanh lơ ở footer, không có link tím ở trang mã.

---

## 3. Typography

Giữ nguyên nguyên tắc v1: **hai họ chữ cho hai màn hình** — đây là chủ ý, có trong ảnh gốc.

```css
--font-sans:  "Be Vietnam Pro", ui-sans-serif, -apple-system, "Segoe UI", Roboto, sans-serif;
--font-serif: "Source Serif 4", "Noto Serif", Charter, Georgia, serif;
--font-mono:  "JetBrains Mono", ui-monospace, "SF Mono", Menlo, monospace;
```

**Ghi chú về serif:** skill anti-slop coi serif là lựa chọn mặc định đáng ngờ. Ở đây serif **có lý do thật**: sản phẩm là bản tin/research note, ảnh tham chiếu của khách hàng dùng serif, và cột tóm tắt là khối text dài cần độ đọc cao. Không dùng `Fraunces` và `Instrument Serif` (hai serif mà mô hình AI hay chọn). `Source Serif 4` có bộ dấu tiếng Việt đầy đủ và self-host được.

**Bắt buộc:** `next/font/google` với `subsets: ["latin", "vietnamese"]`.

### Type scale

| Token | Size / LH | Weight | Dùng ở đâu |
|---|---|---|---|
| `display` | 40px / 1.12 · desktop 52px | 800 | Tiêu đề Research Note (UPPERCASE, `tracking -0.015em`) |
| `h1` | 30px / 1.2 | 700 | Tiêu đề trang digest |
| `h2` | 21px / 1.3 | 700 | `LUẬN ĐIỂM ĐẦU TƯ` (uppercase, `tracking .06em`) |
| `h3` | 18px / 1.4 | 700 | Dòng dẫn mỗi luận điểm |
| `body` | 17px / 1.62 | 400 | Body Research Note (serif) |
| `body-table` | 15.5px / 1.6 | 400 | Cột "Tóm tắt thông tin" (serif) |
| `label` | 13px / 1.4 | 600 | Header bảng |
| `meta` | 13.5px / 1.5 | 400 | Cột Ngày, Source |
| `ticker` | 13.5px / 1.55 | 500 | Mã CK (mono, `tracking .04em`) |
| `num` | — | 500 | **Mọi con số dạng bảng dùng `--font-mono`** (density 7) |

### Nhấn mạnh số liệu

```css
.summary strong {
  font-weight: 700;
  color: var(--ink);
  /* KHÔNG đổi màu. KHÔNG highlight nền. KHÔNG gạch chân. */
}
```

| Rule | Nội dung |
|---|---|
| B-1 | Chỉ `<strong>` và `<em>` được phép trong tóm tắt. HTML khác bị sanitize. |
| B-2 | Đúng **1–3** cụm bold mỗi tóm tắt. |
| B-3 | Cụm bold là **số + đơn vị trọn vẹn**: `408.000 tỷ đồng`, `+1,67%`. Không bold nửa cụm. |
| B-4 | **Không** nhấn số bằng màu. Màu giá chỉ dành cho ô biến động. |

---

## 4. Spacing & Layout

```css
--space-1: 4px;  --space-2: 8px;  --space-3: 12px; --space-4: 16px;
--space-5: 20px; --space-6: 24px; --space-8: 32px; --space-10: 40px;
--space-12: 48px; --space-16: 64px;
```

| Vị trí | Giá trị |
|---|---|
| Padding dọc cell bảng | `--space-5` |
| Padding ngang giữa cột | `--space-6` |
| Giữa các luận điểm | `--space-8` |
| Dòng dẫn → đoạn phân tích | `--space-3` |
| Tiêu đề trang → bảng | `--space-10` |

**Shape Consistency Lock:** `border-radius: 0` ở mọi nơi, **trừ** chip mã CK (`2px`) và focus ring. Một hệ duy nhất, không trộn.
**Shadow:** không dùng. Sticky header chỉ được `box-shadow: 0 1px 0 var(--rule)` — là một đường kẻ, không phải bóng đổ.

---

## 5. Component specs

### 5.1 `DigestTable`

```html
<table>
  <caption class="sr-only">Bảng tin chứng khoán tổng hợp, cập nhật 06:00 16/09/2026</caption>
  <thead>
    <tr>
      <th scope="col">Ngày</th><th scope="col">Mã CK</th>
      <th scope="col">Tóm tắt thông tin</th>
      <th scope="col">Source</th><th scope="col">Loại tin</th>
    </tr>
  </thead>
  <tbody>…</tbody>
</table>
```

| Cột | Width | Ghi chú |
|---|---|---|
| Ngày | `11%` (min 96px) | mono, `nowrap`, `DD/MM/YYYY` |
| Mã CK | `16%` | wrap, nối bằng `, ` |
| Tóm tắt thông tin | `48%` | cột chủ đạo, `max-width: 62ch` |
| Source | `14%` | domain, gạch chân |
| Loại tin | `11%` | text thuần |

```css
table { width: 100%; border-collapse: collapse; font-family: var(--font-serif); }

thead th {
  font-family: var(--font-sans);
  font-size: 13px; font-weight: 700; color: var(--ink);
  text-align: left; padding: 0 var(--space-6) var(--space-3) 0;
  border-bottom: 1px solid var(--rule-strong);
}

tbody td {
  padding: var(--space-5) var(--space-6) var(--space-5) 0;
  border-bottom: 1px solid var(--rule);
  vertical-align: top;
}

tbody tr:last-child td { border-bottom: none; }
tbody tr:hover { background: var(--surface-sunken); }
```

**Cấm:** border dọc · zebra striping · nền màu cho header · bo góc · shadow · badge nhiều màu.

### 5.2 `TickerChip`

- Thường: mono 13.5px, weight 500, `--ink`, không nền.
- Hover: `--accent` + underline.
- Đang chọn: nền `--accent-soft`, chữ `--accent`, `padding: 2px 6px`, `radius: 2px`.
- Luôn là `<a href="/ma/FPT">`.
- \> 8 mã: hiện 8 + nút `+N`.

### 5.3 `SourceLink` — **đã sửa theo phản hồi khách hàng**

Vấn đề cũ: link chỉ hiện domain nên người dùng không biết nó dẫn tới bài nào, và nhiều link không dẫn tới trang cụ thể.

```html
<a href="/r/a7Kx2p"
   target="_blank" rel="nofollow noopener noreferrer"
   title="VN-Index đóng cửa 1.821 điểm, lần đầu vượt mốc 1.800 — cafef.vn, 26/08/2026"
   aria-label="Đọc bài gốc: VN-Index đóng cửa 1.821 điểm — trên CafeF, đăng 26/08/2026 (mở tab mới)">
  cafef.vn
  <span class="sr-only">, bài: VN-Index đóng cửa 1.821 điểm</span>
</a>
```

| Yêu cầu | Chi tiết |
|---|---|
| **Deep link** | `href` trỏ `/r/{code}` → `302` tới **URL bài cụ thể** (`canonical_url`), không bao giờ về trang chủ của báo |
| **Cho biết dẫn đi đâu** | `title` + `aria-label` chứa **tiêu đề bài gốc + tên báo + ngày đăng** |
| Trạng thái chết | `target_alive=false` → hiện `--ink-faint`, `cursor: not-allowed`, chú thích "Bài gốc không còn khả dụng", **không** render thành link |
| Thị giác | `underline`, `text-underline-offset: 2px`, `thickness: 1px`, màu `--ink-2`; hover `--accent` |
| Bắt buộc | `rel="nofollow noopener noreferrer"` |

### 5.4 `NewsTypeLabel`

Text thuần, **không phải badge có nền**, không chấm màu:
`font-family: var(--font-sans); font-size: 13.5px; color: var(--ink-muted);`

`Vĩ mô` · `Ngành` · `Doanh nghiệp` · `Thị trường` · `Khối ngoại` · `Cổ tức/Phát hành` · `Pháp lý` · `Phân tích` · `Trái phiếu/Tín dụng`

### 5.5 `ResearchNote`

```
container max-width 1180px, padding-top 48px

  CMG – QUÝ 2/2026: DOANH THU TIẾP TỤC TĂNG NHƯNG LỢI NHUẬN     display / sans 800
  CĐ MẸ GIẢM – CHU KỲ ĐẦU TƯ DATA CENTER CHƯA TẠO RA LỢI        UPPERCASE, max-width 92ch
  NHUẬN TƯƠNG XỨNG                                              line-height 1.12
  ── 40px ──
  LUẬN ĐIỂM ĐẦU TƯ                                              h2 / uppercase / tracking .06em
  ── 24px ──
  1. Hoạt động kinh doanh cốt lõi chưa xấu đi, nhưng Q2/2026…   h3 / bold / serif
  ── 12px ──
  Q2/2026, tương ứng quý đầu tiên của niên độ tài chính…        body / serif / max-width 68ch
  ── 32px ──
  2. …
  ── 48px ──  hairline
  Nội dung chỉ mang tính thông tin, không phải khuyến nghị       meta / --ink-muted
  đầu tư. Nguồn tham khảo: …
```

`<ol>` với `list-style: none`, số viết inline trong dòng dẫn (đúng ảnh).

### 5.6 `FilterBar`

Sticky top, nền `--surface-raised`, `border-bottom: 1px solid var(--rule)`.
Chứa: ô tìm kiếm · chips mã CK · dropdown Loại tin · khoảng ngày · "Xoá lọc" · "Copy bản tin" · **nút chuyển sáng/tối**.
Control: `--font-sans` 14px, `border: 1px solid var(--rule)`, `radius: 0`, nền `--surface`.

### 5.7 `UpdatedAtBanner`

`Cập nhật lần cuối: 06:00 16/09/2026`, style `meta`.

Quá 12h thì thêm dòng cảnh báo: chữ `--ink` weight 600, tiền tố `⚠ Cảnh báo:`, và một đường kẻ trái `2px solid var(--accent)`.

> **Sửa 16/09.** Bản đầu của mục này bảo dùng màu `--price-down` cho dòng cảnh báo. Sai, và mâu thuẫn với chính quy tắc ở mục 2.3: năm màu giá **chỉ** dùng cho dữ liệu giá. Người dùng Việt nhìn chữ đỏ trong sản phẩm chứng khoán sẽ đọc là "giảm giá", không phải "dữ liệu cũ". Cảnh báo phân biệt bằng **weight + icon + kẻ trái**, không cần màu riêng, và như vậy cũng đạt luôn WCAG 1.4.1 (không truyền đạt chỉ bằng màu).

---

## 6. Responsive

| Breakpoint | Hành vi |
|---|---|
| `< 768px` | **Bảng → card dọc.** Thứ tự: `Ngày · Loại tin` → `Mã CK` → `Tóm tắt` → `Source`. **Không scroll ngang.** |
| `768–1023px` | 4 cột, `Loại tin` gộp vào dòng meta |
| `≥ 1024px` | 5 cột đầy đủ |
| `≥ 1440px` | Container khoá `1180px` |

**Quy tắc cứng:** `document.body.scrollWidth <= window.innerWidth` ở mọi viewport.

---

## 7. Motion (`MOTION_INTENSITY 3`)

Sản phẩm để đọc nhanh. Motion chỉ để **phản hồi thao tác**, không để trang trí.

| Tương tác | Motion | Lý do (bắt buộc có) |
|---|---|---|
| Hover hàng bảng | `background 120ms ease-out` | Phản hồi: cho biết đang ở hàng nào |
| Áp dụng filter | `opacity 0→1, 160ms` | Chuyển trạng thái: nội dung đã đổi |
| Tải thêm tin | fade-in, tối đa 6 dòng đầu | Phân cấp: chỉ ra phần mới |
| Toast copy | trượt 8px + fade, 200ms, ẩn sau 2s | Phản hồi thao tác |

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: .01ms !important;
    transition-duration: .01ms !important;
  }
}
```

**Cấm:** parallax · scroll hijack · marquee · skeleton nhấp nháy toàn trang · animation vòng lặp vô hạn · `window.addEventListener('scroll')`.

---

## 8. Accessibility (WCAG 2.2 AA)

- [ ] `<html lang="vi">`
- [ ] Bảng có `<caption>` sr-only + `<th scope="col">`
- [ ] `outline: 2px solid var(--focus); outline-offset: 2px` — không bao giờ `outline: none`
- [ ] Link tab mới có `aria-label` ghi rõ "(mở tab mới)" **và tiêu đề bài**
- [ ] Filter chips là `<button aria-pressed>`
- [ ] Skip-link "Bỏ qua tới nội dung chính"
- [ ] Target ≥ 24×24px (2.5.8)
- [ ] Tăng/giảm luôn kèm dấu `+`/`−`, không chỉ dùng màu (1.4.1)
- [ ] Font đủ dấu tiếng Việt: `Ừ Ữ Ỡ Ợ Ặ Ẫ Ỹ ọ ự ẳ`
- [ ] `prefers-reduced-motion` được tôn trọng
- [ ] Contrast đã kiểm chứng ở mục 2.4 — **cả hai chế độ**

---

## 9. Ánh xạ token → Tailwind v4

Tailwind v4 dùng `@theme` trong `globals.css`, **không** có `tailwind.config.ts`. Kiểm tra phiên bản thật trước khi viết.

```css
@import "tailwindcss";

@theme {
  --color-surface:        var(--surface);
  --color-surface-raised: var(--surface-raised);
  --color-surface-sunken: var(--surface-sunken);
  --color-ink:            var(--ink);
  --color-ink-2:          var(--ink-2);
  --color-ink-muted:      var(--ink-muted);
  --color-ink-faint:      var(--ink-faint);
  --color-rule:           var(--rule);
  --color-rule-strong:    var(--rule-strong);
  --color-accent:         var(--accent);
  --color-accent-soft:    var(--accent-soft);
  --color-up:             var(--price-up);
  --color-down:           var(--price-down);

  --font-sans:  var(--font-sans);
  --font-serif: var(--font-serif);
  --font-mono:  var(--font-mono);

  --radius-none: 0;
  --radius-chip: 2px;
}
```

---

## 10. Danh sách chống AI-slop — FE phải tick từng dòng trước khi báo cáo xong

**Màu & hình khối**
- [ ] **Không** dùng nhóm màu kem/beige `#F7F5F0`, `#f5f1ea`, `#fbf8f1` — đây là bảng màu AI mặc định
- [ ] **Không** dùng brass / clay / ochre / oxblood làm accent
- [ ] **Không** tím/indigo mặc định của Tailwind, không AI-purple glow
- [ ] **Không** gradient (nền, chữ, nút)
- [ ] Đúng **một** accent cho toàn trang
- [ ] Đúng **một** hệ bo góc
- [ ] **Không** `#000000` và `#ffffff` thuần
- [ ] Không card bo góc + shadow
- [ ] Không zebra striping
- [ ] Không badge pill nhiều màu cho Loại tin
- [ ] **Không chấm màu trang trí** trước nav/nhãn/hàng

**Chữ & nội dung**
- [ ] **Không dùng dấu gạch ngang dài `—` ở bất kỳ chuỗi nào do mình viết.** Dùng dấu `-` thường. *(Ngoại lệ duy nhất: tiêu đề digest tái tạo đúng ảnh gốc của khách hàng — `8 tin mới đáng chú ý — PVS · FPT · VIC.` — vì đó là nội dung của khách, không phải trang trí do ta thêm)*
- [ ] Không eyebrow `00 / INDEX`, `001 · Capabilities`
- [ ] Không "Scroll", "↓ scroll", "Scroll to explore"
- [ ] Không strip địa danh/giờ/thời tiết
- [ ] Không footer version `v1.4.2`, `Build 0048`
- [ ] Không số liệu bịa cho đẹp — mọi số phải từ dữ liệu thật
- [ ] Không tên giả kiểu "Acme", "Nexus", "John Doe"
- [ ] Không động từ rỗng: "Elevate", "Seamless", "Unleash", "Next-Gen"

**Cấu trúc**
- [ ] Không hero có ảnh nền — sản phẩm này không có hero
- [ ] Không sidebar icon nav — chỉ có 3 trang
- [ ] Không fake screenshot dựng bằng `<div>`
- [ ] Không SVG icon tự vẽ tay — dùng thư viện icon
- [ ] Không skeleton nhấp nháy toàn trang

**Đặc thù sản phẩm**
- [ ] Nhấn số bằng **font-weight**, không bằng màu
- [ ] Năm màu giá **chỉ** dùng cho dữ liệu giá
- [ ] Bảng **không** scroll ngang trên mobile
- [ ] Mọi link nguồn dẫn tới **bài cụ thể**, có `title` ghi rõ tiêu đề bài
- [ ] Dark mode hoạt động đầy đủ, đã xem thật ở cả hai chế độ
