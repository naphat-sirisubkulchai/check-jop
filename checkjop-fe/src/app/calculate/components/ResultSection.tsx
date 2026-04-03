"use client";

import { GraduationResult } from "@/types";
import { useMemo, useEffect } from "react";
import { OverviewCards } from "./OverviewCards";
import { CategoryResultsCard } from "./CategoryResultsCard";
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

export function ResultSection({
  result,
}: ResultSectionProps) {
  const { trackEvent } = useAnalytics();
  const { studyPlan } = useAppStore();

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

  // Track when results are viewed
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
        {/* Graduation Eligibility Banner */}
        <GraduationBanner result={result} />

        {/* Missing Catalog Years Warning */}
        <MissingCatalogYearsCard
          missingCatalogYears={result.missing_catalog_years}
          admissionYear={admissionYear}
          catalogYearFallbacks={result.catalog_year_fallbacks}
        />

        {/* Progress Overview Cards */}
        <OverviewCards
          result={result}
          satisfiedCategories={satisfiedCategories}
          totalCategories={totalCategories}
        />

        {/* Two Column Layout: Study Plan Record & Analysis Results */}
        <div className="result-two-col grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Left Column - Study Plan Record */}
          <div className="lg:sticky lg:top-6 lg:self-start print:static">
            <StudyPlanRecordCard studyPlan={studyPlan} />
          </div>

          {/* Right Column - Analysis Results */}
          <div className="space-y-6">
            {/* Category Details */}
            <CategoryResultsCard categoryResults={result.category_results} />

            {/* Pre and Co requisite Violations */}
            <PrerequisiteViolationsCard
              violations={result.prerequisite_violations}
            />

            {/* Credit Limit Violations */}
            <CreditLimitViolationsCard
              violations={result.credit_limit_violations}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
