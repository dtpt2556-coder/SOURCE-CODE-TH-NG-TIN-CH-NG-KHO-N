"use client";

import { useEffect } from "react";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="mx-auto max-w-content px-4 py-20 md:px-6">
      <h1 className="font-sans text-[30px] font-bold leading-[1.2] text-ink">
        Đã có lỗi khi tải nội dung.
      </h1>
      <p className="mt-3 max-w-prose font-serif text-[17px] leading-[1.62] text-ink-2">
        Hệ thống không lấy được dữ liệu lúc này. Vui lòng thử lại sau ít phút.
      </p>
      <p className="mt-6">
        <button className="btn btn-strong" type="button" onClick={reset}>
          Thử lại
        </button>
      </p>
    </div>
  );
}
