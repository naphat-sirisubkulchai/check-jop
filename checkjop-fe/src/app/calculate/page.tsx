"use client";

import { useAppStore } from "@/store/appStore";
import { ResultSection } from "./components/ResultSection";
import { Loader2, AlertCircle, Calculator } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";

export default function CalculatePage() {
  const { result, isLoading } = useAppStore();
  const router = useRouter();

  if (isLoading) {
    return (
      <div className="h-full flex items-center justify-center bg-gradient-to-br from-chula-soft/30 via-pink-50 to-white">
        <div className="text-center space-y-4">
          <div className="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-br from-chula-active to-pink-500 rounded-2xl shadow-lg animate-pulse">
            <Loader2 className="h-10 w-10 text-white animate-spin" />
          </div>
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Calculating Your Progress</h2>
            <p className="text-gray-600 mt-2">Analyzing your graduation eligibility...</p>
          </div>
        </div>
      </div>
    );
  }

  if (!result) {
    return (
      <div className="h-full flex items-center justify-center bg-gradient-to-br from-chula-soft/30 via-pink-50 to-white px-4">
        <div className="max-w-md w-full bg-white rounded-2xl shadow-xl p-8 border border-gray-200">
          <div className="text-center space-y-4">
            <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-gray-100 to-gray-200 rounded-2xl">
              <AlertCircle className="h-8 w-8 text-gray-400" />
            </div>
            <div>
              <h2 className="text-2xl font-bold text-gray-900 mb-2">No Results Available</h2>
              <p className="text-gray-600 text-sm">
                Add courses to your study plan and calculate your graduation eligibility to see results here.
              </p>
            </div>
            <Button
              onClick={() => router.push("/home")}
              className="w-full bg-gradient-to-r from-chula-active to-pink-500 hover:from-chula-active/90 hover:to-pink-600 text-white font-semibold shadow-lg hover:shadow-xl"
            >
              <Calculator className="h-4 w-4 mr-2" />
              Go to Study Plan
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full overflow-hidden">
      {/* Main Content */}
      <main className="flex-1 overflow-y-auto">
        <div className="mx-8 my-6">
          <ResultSection result={result} />
        </div>
      </main>
    </div>
  );
}
