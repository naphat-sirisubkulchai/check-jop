"use client";

import { AlertTriangle } from "lucide-react";

interface MissingCatalogYearsCardProps {
  missingCatalogYears: number[];
  admissionYear?: number;
  catalogYearFallbacks?: Record<number, number>;
}

export function MissingCatalogYearsCard({
  missingCatalogYears,
  catalogYearFallbacks,
}: MissingCatalogYearsCardProps) {
  if (!missingCatalogYears || missingCatalogYears.length === 0) return null;

  const sorted = [...missingCatalogYears].sort((a, b) => a - b);

  return (
    <div className="w-full bg-amber-50 border-2 border-amber-200 rounded-xl p-5 shadow-sm">
      <div className="flex items-start gap-3">
        <div className="flex-shrink-0 w-9 h-9 bg-amber-400 rounded-full flex items-center justify-center">
          <AlertTriangle className="w-5 h-5 text-white" strokeWidth={2.5} />
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="text-amber-800 font-bold text-base">
            ข้อมูลรายวิชาไม่ครบถ้วน
          </h3>
          <p className="text-amber-700 text-sm mt-1">
            คุณมีรายวิชาที่เรียนในปีการศึกษา{" "}
            <span className="font-semibold">{sorted.join(", ")}</span>{" "}
            ซึ่งยังไม่มีข้อมูลหลักสูตรในระบบ ผลการตรวจสอบ pre/co req
            ของแต่ละปีจะอิงจากหลักสูตรปีก่อนหน้าที่ใกล้ที่สุดที่มีข้อมูลแทน
          </p>
          <div className="flex flex-wrap gap-2 mt-3">
            {sorted.map((year) => {
              const fallback = catalogYearFallbacks?.[year];
              const usesFallback = fallback !== undefined && fallback !== year;
              return (
                <span
                  key={year}
                  className="inline-flex items-center gap-1 bg-amber-100 border border-amber-300 text-amber-800 text-xs font-semibold px-3 py-1 rounded-full"
                >
                  <AlertTriangle className="w-3 h-3" />
                  ปี {year} ยังไม่มีข้อมูล
                  {usesFallback && (
                    <span className="font-normal text-amber-700">
                      {" "}→ ใช้ปี {fallback} แทน
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
}
