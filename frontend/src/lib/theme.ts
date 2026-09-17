/**
 * Che do sang/toi - design-system v2 muc 2.2.
 *
 * Mac dinh theo `prefers-color-scheme`. Nguoi dung chon thu cong thi ghi
 * `data-theme` len <html> va luu localStorage; CSS da viet sao cho
 * `:root[data-theme="..."]` thang duoc media query o ca hai chieu.
 *
 * Theme Lock: chi co MOT thuoc tinh `data-theme` tren <html>, khong component
 * nao duoc tu dao mau cuc bo.
 */

export type Theme = "light" | "dark";

export const THEME_STORAGE_KEY = "tenpoint.theme";

/**
 * Script chan FOUC: chay dong bo trong <head>, truoc khi trinh duyet paint.
 * Chi gan `data-theme` khi nguoi dung da chon thu cong, de mac dinh van la
 * `prefers-color-scheme`.
 */
export const THEME_INIT_SCRIPT = `(function(){try{var t=localStorage.getItem(${JSON.stringify(
  THEME_STORAGE_KEY,
)});if(t==="dark"||t==="light"){document.documentElement.setAttribute("data-theme",t)}}catch(e){}})();`;

export function readSystemTheme(): Theme {
  return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
}

export function readCurrentTheme(): Theme {
  const explicit = document.documentElement.getAttribute("data-theme");
  if (explicit === "dark" || explicit === "light") return explicit;
  return readSystemTheme();
}

/** Dao che do hien tai, ghi DOM + localStorage, tra ve che do moi. */
export function toggleTheme(): Theme {
  const next: Theme = readCurrentTheme() === "dark" ? "light" : "dark";
  document.documentElement.setAttribute("data-theme", next);
  try {
    window.localStorage.setItem(THEME_STORAGE_KEY, next);
  } catch {
    // Trinh duyet chan storage (che do rieng tu): van doi duoc trong phien nay.
  }
  return next;
}
