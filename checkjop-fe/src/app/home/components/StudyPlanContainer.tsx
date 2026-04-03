"use client";
import SemesterCard from "./SemesterCard";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Trash2, ChevronDown, ChevronRight } from "lucide-react";
import { useStudyPlan } from "@/hooks/useStudyPlan";
import { useRouter } from "next/navigation";
import { useAppStore } from "@/store/appStore";
import { useState } from "react";

export default function StudyPlanContainer() {
  const router = useRouter();
  const { handleCalculate } = useStudyPlan();
  const {
    selectedCurriculum,
    clearStudyPlan,
    clearExemptions,
    yearMapping,
    totalCredits,
  } = useAppStore();

  const [collapsedYears, setCollapsedYears] = useState<Record<number, boolean>>({});
  const [showClearDialog, setShowClearDialog] = useState(false);

  const toggleYearCollapse = (year: number) => {
    setCollapsedYears((prev) => ({
      ...prev,
      [year]: !prev[year],
    }));
  };

  function handleClearStudyPlan(): void {
    setShowClearDialog(true);
  }

  function confirmClear(): void {
    clearStudyPlan();
    clearExemptions();
    setShowClearDialog(false);
  }

  // Get years array from yearMapping
  const years = yearMapping
    ? Object.keys(yearMapping).map(Number).sort()
    : [1, 2, 3, 4];

  function handleSubmit(): void {
    try {
      handleCalculate();
      router.push("/calculate");
    } catch (error) {
      alert(error instanceof Error ? error.message : "An unknown error occurred.");
    }
  }

  return (
    <div className="flex flex-col flex-1 h-full px-4 pt-4">
      {/* Top section */}
      <div className="flex justify-between items-center mb-4">
        {/* Academic Plan & Record Header */}
        <div className="flex items-baseline">
          <h3 className="text-xl font-bold">Academic Plan & Record</h3>
          <p className="ml-4 text-gray-600">{selectedCurriculum?.nameTH}</p>
        </div>

        {/* Summary Course Credits */}
        <div className="ml-4">
          <p className="text-lg text-gray-600 truncate">
            Total : <span className="text-chula-active font-semibold text-2xl">{totalCredits}</span> /{" "}
            {selectedCurriculum?.minTotalCredits} credits
          </p>
        </div>
      </div>

      {/* Exemptions Section */}
      {/* <ExemptionsSection /> */}

      {/* Plan area - Scrollable all years view */}
      <section className="flex-1 overflow-y-auto space-y-6 overflow-x-hidden">
        {/* Years Section */}
        {years.map((year) => {
          const academicYear = yearMapping?.[year] || 2566 + year - 1;

          return (
            <div key={year} className="space-y-4">
              {/* Year Header */}
              <button
                className="w-full flex items-center justify-between cursor-pointer -mx-2 px-2 py-2 transition-colors hover:bg-chula-soft"
                onClick={() => toggleYearCollapse(year)}
                aria-expanded={!collapsedYears[year]}
                aria-controls={`year-${year}-content`}
              >
                <div className="flex items-center gap-2">
                  <h2 className="text-xl font-bold text-gray-900">
                    Year {year}{" "}
                    <span className="text-lg font-light text-gray-600">({academicYear})</span>
                  </h2>
                  {collapsedYears[year] ? (
                    <ChevronRight className="size-4 text-gray-400" />
                  ) : (
                    <ChevronDown className="size-4 text-gray-400" />
                  )}
                </div>
              </button>

              {/* Collapsible Content */}
              {!collapsedYears[year] && (
                <div id={`year-${year}-content`}>
                  {/* Semester Grid - 2 or 3 columns based on summer term */}
                  <div
                    className={`grid grid-cols-1 gap-4 md:grid-cols-3`}
                  >
                    {/* Semester 1 */}
                    <div className="flex flex-col">
                      <SemesterCard
                        sem={1}
                        yearOfStudy={year}
                        academicYear={academicYear}
                        className="flex-1 min-h-[300px]"
                      />
                    </div>

                    {/* Semester 2 */}
                    <div className="flex flex-col">
                      <SemesterCard
                        sem={2}
                        yearOfStudy={year}
                        academicYear={academicYear}
                        className="flex-1 min-h-[300px]"
                      />
                    </div>

                    {/* Semester 3 (Summer) */}
                    <div className="flex flex-col">
                      <SemesterCard
                        sem={3}
                        yearOfStudy={year}
                        academicYear={academicYear}
                        className="flex-1 min-h-[300px]"
                      />
                    </div>
                  </div>
                </div>
              )}
            </div>
          );
        })}
      </section>

      {/* Action Buttons */}
      <div className="p-4 flex justify-end gap-3">
        <div className="flex gap-2">
          <Button onClick={handleClearStudyPlan} variant={"outline"} size="default" className="bg-white shadow-sm px-8 py-5">
            <Trash2 className="h-5 w-5 mr-2" />
            Clear Table
          </Button>
        </div>
        <Button
          disabled={!selectedCurriculum}
          onClick={handleSubmit}
          className="bg-gradient-to-r from-chula-active to-pink-500 hover:from-chula-active/90 hover:to-pink-600 text-white px-8 py-5 shadow-md hover:shadow-lg font-semibold text-base"
          size="default"
        >
          Calculate Eligibility
        </Button>
      </div>

      <Dialog open={showClearDialog} onOpenChange={setShowClearDialog}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-red-100 flex items-center justify-center shrink-0">
                <Trash2 className="h-5 w-5 text-red-600" />
              </div>
              <DialogTitle>Clear Study Plan?</DialogTitle>
            </div>
          </DialogHeader>
          <p className="text-sm text-gray-600">
            วิชาทั้งหมดและ exemptions จะถูกลบออก ไม่สามารถย้อนกลับได้
          </p>
          <DialogFooter className="gap-2">
            <Button variant="outline" onClick={() => setShowClearDialog(false)}>ยกเลิก</Button>
            <Button variant="destructive" onClick={confirmClear}>ลบทั้งหมด</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
