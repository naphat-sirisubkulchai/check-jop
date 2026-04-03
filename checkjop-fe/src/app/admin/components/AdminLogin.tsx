"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import { Lock, User, AlertCircle, Shield } from "lucide-react";

interface AdminLoginProps {
  onLogin: (username: string, password: string) => boolean;
}

export default function AdminLogin({ onLogin }: AdminLoginProps) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setIsLoading(true);

    // Simulate loading delay for better UX
    await new Promise((resolve) => setTimeout(resolve, 500));

    const success = onLogin(username, password);

    if (!success) {
      setError("Invalid username or password");
      setPassword("");
    }

    setIsLoading(false);
  };

  return (
    <div className="flex items-center justify-center h-full">
      <Card className="w-full max-w-md p-8 shadow-2xl border-t-4 border-t-chula-active">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-chula-active to-pink-500 rounded-2xl mb-4 shadow-lg">
            <Shield className="h-8 w-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold bg-gradient-to-r from-chula-active to-pink-500 bg-clip-text text-transparent mb-2">
            Admin Login
          </h1>
          <p className="text-gray-600 text-sm">
            Enter your credentials to access admin panel
          </p>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border-l-4 border-red-500 rounded-lg flex items-start gap-3 animate-in slide-in-from-top-2">
            <AlertCircle className="h-5 w-5 text-red-600 flex-shrink-0 mt-0.5" />
            <p className="text-sm text-red-800 font-medium">{error}</p>
          </div>
        )}

        {/* Login Form */}
        <form onSubmit={handleSubmit} className="space-y-5">
          {/* Username */}
          <div>
            <label htmlFor="username" className="block text-sm font-semibold text-gray-700 mb-2">
              Username
            </label>
            <div className="relative">
              <User className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
              <Input
                id="username"
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="Enter username"
                className="pl-10 h-12"
                required
                autoFocus
                disabled={isLoading}
              />
            </div>
          </div>

          {/* Password */}
          <div>
            <label htmlFor="password" className="block text-sm font-semibold text-gray-700 mb-2">
              Password
            </label>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter password"
                className="pl-10 h-12"
                required
                disabled={isLoading}
              />
            </div>
          </div>

          {/* Submit Button */}
          <Button
            type="submit"
            disabled={isLoading || !username || !password}
            className="w-full h-12 bg-gradient-to-r from-chula-active to-pink-500 hover:from-chula-active/90 hover:to-pink-600 text-white font-semibold shadow-lg hover:shadow-xl transition-all"
          >
            {isLoading ? (
              <div className="flex items-center gap-2">
                <div className="h-5 w-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                Signing in...
              </div>
            ) : (
              "Sign In"
            )}
          </Button>
        </form>

        {/* Footer */}
        <div className="mt-6 text-center">
          <p className="text-xs text-gray-500">
            Protected admin area. Unauthorized access is prohibited.
          </p>
        </div>
      </Card>
    </div>
  );
}
