"use client";

import { toggleTheme } from "@/lib/theme";

/**
 * Nut chuyen sang/toi - muc 5.6.
 *
 * Nhan hien thi do CSS quyet dinh (`.theme-label-to-dark` / `.theme-label-to-light`)
 * chu khong do React state, nen HTML server va client giong het nhau: khong lech
 * hydration, khong nhap nhay, va nhan van dung ca khi JS chua chay.
 * Chi mot nhan duoc `display: inline` nen screen reader doc dung mot nhan.
 */
export function ThemeToggle() {
  return (
    <button
      type="button"
      className="btn"
      onClick={() => {
        toggleTheme();
      }}
    >
      <span className="theme-label-to-dark">Chuyển sang nền tối</span>
      <span className="theme-label-to-light">Chuyển sang nền sáng</span>
    </button>
  );
}
