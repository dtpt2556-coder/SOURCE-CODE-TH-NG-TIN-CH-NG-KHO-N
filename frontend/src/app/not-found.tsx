import type { Metadata } from "next";
import Link from "next/link";

export const metadata: Metadata = {
  title: "Không tìm thấy trang",
  robots: { index: false, follow: false },
};

export default function NotFound() {
  return (
    <div className="mx-auto max-w-content px-4 py-20 md:px-6">
      <h1 className="font-sans text-[30px] font-bold leading-[1.2] text-ink">
        Không tìm thấy trang này.
      </h1>
      <p className="mt-3 max-w-prose font-serif text-[17px] leading-[1.62] text-ink-2">
        Đường dẫn có thể đã thay đổi hoặc nội dung không còn tồn tại. Bạn có thể quay về
        bảng tin để xem các tin mới nhất.
      </p>
      <p className="mt-6">
        <Link href="/" className="btn btn-strong inline-block no-underline">
          Về bảng tin
        </Link>
      </p>
    </div>
  );
}
