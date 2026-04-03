"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import AdminLogin from "@/app/admin/components/AdminLogin";

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isChecking, setIsChecking] = useState(true);
  const router = useRouter();

  useEffect(() => {
    // Check if user is authenticated
    const authStatus = sessionStorage.getItem("admin_authenticated");
    setIsAuthenticated(authStatus === "true");
    setIsChecking(false);
  }, []);

  const handleLogin = (username: string, password: string): boolean => {
    // Get credentials from environment variables
    const validUsername = process.env.NEXT_PUBLIC_ADMIN_USERNAME;
    const validPassword = process.env.NEXT_PUBLIC_ADMIN_PASSWORD;

    if (username === validUsername && password === validPassword) {
      sessionStorage.setItem("admin_authenticated", "true");
      setIsAuthenticated(true);
      return true;
    }
    return false;
  };

  const handleLogout = () => {
    sessionStorage.removeItem("admin_authenticated");
    setIsAuthenticated(false);
    router.push("/admin");
  };

  // Show loading state while checking authentication
  if (isChecking) {
    return (
      <div className="flex items-center justify-center h-full bg-gray-50">
        <div className="text-center">
          <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-chula-active border-r-transparent"></div>
          <p className="mt-4 text-gray-600">Checking authentication...</p>
        </div>
      </div>
    );
  }

  // Show login page if not authenticated
  if (!isAuthenticated) {
    return <AdminLogin onLogin={handleLogin} />;
  }

  // Show admin content if authenticated
  return (
    <div className="h-full">
      {/* Logout button - could be integrated into header */}
      <div className="fixed top-4 right-4 z-50">
        <button
          onClick={handleLogout}
          className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors shadow-lg"
        >
          Logout
        </button>
      </div>
      {children}
    </div>
  );
}
