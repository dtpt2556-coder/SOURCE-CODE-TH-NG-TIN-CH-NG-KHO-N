import type { Metadata } from "next";
import { Be_Vietnam_Pro, JetBrains_Mono, Source_Serif_4 } from "next/font/google";

import { SiteFooter, SiteHeader } from "@/components/SiteChrome";
import { SITE_NAME, SITE_URL } from "@/lib/site";
import { THEME_INIT_SCRIPT } from "@/lib/theme";

import "./globals.css";

/**
 * Bat buoc `subsets: ["latin", "vietnamese"]` - design-system v2 muc 3.
 * Chuoi kiem tra dau: Ừ Ữ Ỡ Ợ Ặ Ẫ Ỹ ọ ự ẳ
 */
const beVietnamPro = Be_Vietnam_Pro({
  variable: "--font-be-vietnam-pro",
  subsets: ["latin", "vietnamese"],
  weight: ["400", "500", "600", "700", "800"],
  display: "swap",
});

const sourceSerif = Source_Serif_4({
  variable: "--font-source-serif-4",
  subsets: ["latin", "vietnamese"],
  display: "swap",
});

const jetBrainsMono = JetBrains_Mono({
  variable: "--font-jetbrains-mono",
  subsets: ["latin", "vietnamese"],
  display: "swap",
});

const SITE_TAGLINE = "Tổng hợp & tóm tắt tin chứng khoán Việt Nam";
const SITE_DESCRIPTION =
  "Bảng tin chứng khoán Việt Nam được tóm tắt cô đọng, giàu số liệu, gắn mã CK và luôn kèm liên kết về bài gốc. Cập nhật 3 lần mỗi ngày.";

export const metadata: Metadata = {
  metadataBase: new URL(SITE_URL),
  title: {
    default: `${SITE_NAME}: ${SITE_TAGLINE}`,
    template: `%s · ${SITE_NAME}`,
  },
  description: SITE_DESCRIPTION,
  applicationName: SITE_NAME,
  alternates: { canonical: "/" },
  openGraph: {
    type: "website",
    locale: "vi_VN",
    siteName: SITE_NAME,
    title: `${SITE_NAME}: ${SITE_TAGLINE}`,
    description: SITE_DESCRIPTION,
    url: "/",
  },
  twitter: {
    card: "summary",
    title: `${SITE_NAME}: ${SITE_TAGLINE}`,
    description: SITE_DESCRIPTION,
  },
  robots: { index: true, follow: true },
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html
      lang="vi"
      className={`${beVietnamPro.variable} ${sourceSerif.variable} ${jetBrainsMono.variable} h-full`}
      suppressHydrationWarning
    >
      <head>
        {/*
          Chan FOUC: chay dong bo khi trinh duyet parse HTML, truoc paint dau tien,
          nen khong bao gio thay chop nen sang roi moi doi sang toi.
          Noi dung la hang so do FE viet, khong nhan du lieu tu ben ngoai.
        */}
        <script dangerouslySetInnerHTML={{ __html: THEME_INIT_SCRIPT }} />
      </head>
      <body className="flex min-h-full flex-col bg-surface text-ink">
        <a className="skip-link" href="#noi-dung">
          Bỏ qua tới nội dung chính
        </a>
        <SiteHeader />
        <main id="noi-dung" className="flex-1">
          {children}
        </main>
        <SiteFooter />
      </body>
    </html>
  );
}
