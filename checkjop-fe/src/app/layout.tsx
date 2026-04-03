import type { Metadata } from "next";
import { Bai_Jamjuree } from "next/font/google";
import "./globals.css";
import { CustomToaster } from "@/components/custom/CustomToaster";
import StudyPlanInitializer from "@/components/StudyPlanInitializer";
import Header from "@/components/Header";

const baiJamjuree = Bai_Jamjuree({
  weight: ["200", "300", "400", "500", "600", "700"],
  variable: "--font-baijamjuree",
  subsets: ["thai", "latin"],
  fallback: ["sans-serif"],
});

export const metadata: Metadata = {
  title: "CheckJop",
  description: "ระบบตรวจสอบการสำเร็จการศึกษา",
  icons: {
    icon: "/icon.png",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="h-full">
      <body
        className={`${baiJamjuree.variable} antialiased h-full flex flex-col bg-gradient-to-br from-chula-soft via-pink-50 to-white`}
      >
        <StudyPlanInitializer />
        <Header />
        <div className="flex-1 overflow-hidden">
          {children}
        </div>
        <CustomToaster />
      </body>
    </html>
  );
}
