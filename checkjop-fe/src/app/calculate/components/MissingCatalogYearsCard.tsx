"use client";

import { AlertTriangle, X } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";

interface MissingCatalogYearsCardProps {
  missingCatalogYears: number[];
  admissionYear?: number;
  catalogYearFallbacks?: Record<number, number>;
}

export function MissingCatalogYearsCard({
  missingCatalogYears,
  catalogYearFallbacks,
}: MissingCatalogYearsCardProps) {
  const [dismissed, setDismissed] = useState(false);

  if (!missingCatalogYears || missingCatalogYears.length === 0) return null;

  const sorted = [...missingCatalogYears].sort((a, b) => a - b);

  const banner = (
    <div className="w-full bg-red-50 border-2 border-red-300 rounded-xl p-5 shadow-sm">
      <div className="flex items-start gap-3">
        <div className="flex-shrink-0 w-9 h-9 bg-red-500 rounded-full flex items-center justify-center">
          <AlertTriangle className="w-5 h-5 text-white" strokeWidth={2.5} />
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="text-red-800 font-bold text-base">
            Incomplete Course Data
          </h3>
          <p className="text-red-700 text-sm mt-1">
            You have courses taken in academic year{" "}
            <span className="font-semibold">{sorted.join(", ")}</span>{" "}
            which has no curriculum data in the system. Pre/co req checks for
            those years will use the nearest earlier available catalog year instead.
          </p>
          <div className="flex flex-wrap gap-2 mt-3">
            {sorted.map((year) => {
              const fallback = catalogYearFallbacks?.[year];
              const usesFallback = fallback !== undefined && fallback !== year;
              return (
                <span
                  key={year}
                  className="inline-flex items-center gap-1 bg-red-100 border border-red-300 text-red-800 text-xs font-semibold px-3 py-1 rounded-full"
                >
                  <AlertTriangle className="w-3 h-3" />
                  Year {year} not available
                  {usesFallback && (
                    <span className="font-normal text-red-700">
                      {" "}→ using {fallback} instead
                    </span>
                  )}
                </span>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );

  return (
    <>
      {/* Inline banner */}
      {banner}

      {/* Modal overlay */}
      {!dismissed && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 print:hidden">
          <div className="relative w-full max-w-lg mx-4 bg-white rounded-2xl shadow-2xl border-2 border-red-300 p-6">
            <button
              onClick={() => setDismissed(true)}
              className="absolute top-4 right-4 text-gray-400 hover:text-gray-600"
            >
              <X className="w-5 h-5" />
            </button>
            <div className="flex items-center gap-3 mb-4">
              <div className="w-11 h-11 bg-red-500 rounded-full flex items-center justify-center flex-shrink-0">
                <AlertTriangle className="w-6 h-6 text-white" strokeWidth={2.5} />
              </div>
              <h2 className="text-red-800 font-bold text-lg">Incomplete Course Data</h2>
            </div>
            <p className="text-red-700 text-sm mb-4">
              You have courses taken in academic year{" "}
              <span className="font-semibold">{sorted.join(", ")}</span>{" "}
              which has no curriculum data in the system. Pre/co req checks for
              those years will use the nearest earlier available catalog year instead.
            </p>
            <div className="flex flex-wrap gap-2 mb-5">
              {sorted.map((year) => {
                const fallback = catalogYearFallbacks?.[year];
                const usesFallback = fallback !== undefined && fallback !== year;
                return (
                  <span
                    key={year}
                    className="inline-flex items-center gap-1 bg-red-100 border border-red-300 text-red-800 text-xs font-semibold px-3 py-1 rounded-full"
                  >
                    <AlertTriangle className="w-3 h-3" />
                    Year {year} not available
                    {usesFallback && (
                      <span className="font-normal text-red-700">
                        {" "}→ using {fallback} instead
                      </span>
                    )}
                  </span>
                );
              })}
            </div>
            <Button
              onClick={() => setDismissed(true)}
              className="w-full bg-red-500 hover:bg-red-600 text-white"
            >
              Acknowledge
            </Button>
          </div>
        </div>
      )}
    </>
  );
}
