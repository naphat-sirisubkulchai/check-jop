"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import HomeActions from "./header/HomeActions";
import CalculateActions from "./header/CalculateActions";

export default function Header() {
  const pathname = usePathname();

  return (
    <header className="bg-chula-active px-6 py-2.5 flex justify-between items-center flex-shrink-0 shadow-md z-20">
      {/* Logo & Title */}
      <div className="flex items-center gap-1 flex-shrink-0 px-2">
        <Link href="/" className="flex-shrink-0">
          <img
            src="/CheckJop_logo4.png"
            alt="Check Jop Logo"
            className="h-12 cursor-pointer hover:opacity-90 transition-opacity"
          />
        </Link>
      </div>

      {/* Right Section: Action Buttons */}
      <div className="flex items-center gap-3">
        {pathname === "/home" && <HomeActions />}
        {pathname === "/calculate" && <CalculateActions />}
      </div>
    </header>
  );
}
