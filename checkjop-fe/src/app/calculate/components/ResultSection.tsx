"use client";

import { GraduationResult } from "@/types";
import { useMemo, useEffect } from "react";
import { OverviewCards } from "./OverviewCards";
import { CategoryResultsCard } from "./CategoryResultsCard";
import { PrintTranscriptView } from "./PrintTranscriptView";
import { PrerequisiteViolationsCard } from "./PrerequisiteViolationsCard";
import { CreditLimitViolationsCard } from "./CreditLimitViolationsCard";
import { StudyPlanRecordCard } from "./StudyPlanRecordCard";
import { GraduationBanner } from "./GraduationBanner";
import { MissingCatalogYearsCard } from "./MissingCatalogYearsCard";
import { useAnalytics } from "@/hooks/useAnalytics";
import { useAppStore } from "@/store/appStore";

interface ResultSectionProps {
  result: GraduationResult;
}

export function ResultSection({ result }: ResultSectionProps) {
  const { trackEvent } = useAnalytics();
  const { studyPlan, printFormat } = useAppStore();

  const admissionYear = useMemo(
    () => studyPlan.length > 0 ? Math.min(...studyPlan.map((p) => p.academicYear)) : 0,
    [studyPlan]
  );

  const satisfiedCategories = useMemo(
    () => result.category_results.filter((cat) => cat.is_satisfied).length || 0,
    [result]
  );

  const totalCategories = useMemo(
    () => result.category_results.length || 0,
    [result]
  );

  useEffect(() => {
    if (result) {
      trackEvent("result_viewed", {
        canGraduate: result.can_graduate,
        totalCredits: result.total_credits,
        gpax: result.gpax,
        violationCount:
          result.prerequisite_violations.length +
          result.credit_limit_violations.length,
      });
    }
  }, [result, trackEvent]);

  return (
    <div className="space-y-6">
      <div aria-label="Graduation analysis results" className="space-y-6">

        {/* ── Print-only content ── */}
        <div className="hidden print:block">
          {printFormat === "transcript" ? (
            <PrintTranscriptView result={result} />
          ) : (
            /* summary: banner + overview + full two-col layout rendered inline */
            <div className="space-y-4">
              <GraduationBanner result={result} />
              <OverviewCards
                result={result}
                satisfiedCategories={satisfiedCategories}
                totalCategories={totalCategories}
              />
              <CategoryResultsCard categoryResults={result.category_results} />
              <div className="grid grid-cols-2 gap-4">
                <PrerequisiteViolationsCard violations={result.prerequisite_violations} />
                <CreditLimitViolationsCard violations={result.credit_limit_violations} />
              </div>
              <StudyPlanRecordCard studyPlan={studyPlan} />
            </div>
          )}
        </div>

        {/* ── Screen-only content ── */}
        <div className="print:hidden space-y-6">
          <GraduationBanner result={result} />

          <MissingCatalogYearsCard
            missingCatalogYears={result.missing_catalog_years}
            admissionYear={admissionYear}
            catalogYearFallbacks={result.catalog_year_fallbacks}
          />

          <OverviewCards
            result={result}
            satisfiedCategories={satisfiedCategories}
            totalCategories={totalCategories}
          />

          <div className="result-two-col grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="lg:sticky lg:top-6 lg:self-start">
              <StudyPlanRecordCard studyPlan={studyPlan} />
            </div>
            <div className="space-y-6">
              <CategoryResultsCard categoryResults={result.category_results} />
              <PrerequisiteViolationsCard violations={result.prerequisite_violations} />
              <CreditLimitViolationsCard violations={result.credit_limit_violations} />
            </div>
          </div>
        </div>

      </div>
    </div>
  );
}
