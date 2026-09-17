import Link from "next/link";

/** Masthead kiểu ấn phẩm nghiên cứu - không sidebar, không icon nav. */
export function SiteHeader() {
  return (
    <header className="border-b border-rule">
      <div className="mx-auto flex max-w-content flex-wrap items-baseline justify-between gap-x-6 gap-y-2 px-4 py-4 md:px-6">
        <Link
          href="/"
          className="font-sans text-[20px] font-extrabold uppercase tracking-[-0.01em] text-ink no-underline"
        >
          TenPoint
        </Link>
        <p className="order-3 basis-full text-[13px] text-ink-muted md:order-2 md:basis-auto">
          Tổng hợp &amp; tóm tắt tin chứng khoán Việt Nam
        </p>
        <nav aria-label="Điều hướng chính" className="order-2 md:order-3">
          <ul className="flex gap-5 text-[14px]">
            <li>
              <Link href="/" className="text-ink no-underline hover:underline">
                Bản tin
              </Link>
            </li>
            <li>
              <Link href="/nhan-dinh" className="text-ink no-underline hover:underline">
                Nhận định
              </Link>
            </li>
          </ul>
        </nav>
      </div>
    </header>
  );
}

export function SiteFooter() {
  return (
    <footer className="mt-20 border-t border-rule">
      <div className="mx-auto max-w-content px-4 py-8 text-[13px] leading-[1.6] text-ink-muted md:px-6">
        <p className="max-w-prose">
          Bản tóm tắt được hệ thống tạo tự động từ các nguồn báo chí chính thống và luôn
          kèm liên kết về bài gốc. Nội dung chỉ mang tính thông tin, không phải khuyến
          nghị đầu tư.
        </p>
        <p className="mt-3">© {new Date().getFullYear()} TenPoint</p>
      </div>
    </footer>
  );
}
