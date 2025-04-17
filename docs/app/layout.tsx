import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "UQL",
  description: "UQL",
};

type Props = Readonly<{
  children: React.ReactNode;
}>;

export default function RootLayout({ children }: Props) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`antialiased`}>{children}</body>
    </html>
  );
}
