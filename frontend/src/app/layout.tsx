import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "PageMail - 网页抓取工具",
  description: "抓取网页内容并发送到您的邮箱。支持HTML、PDF、截图格式。",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="font-sans antialiased">
        {children}
      </body>
    </html>
  );
}
