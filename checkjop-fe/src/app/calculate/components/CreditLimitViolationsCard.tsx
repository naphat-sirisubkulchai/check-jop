"use client";

import { XCircle, CheckCircle } from "lucide-react";
import { CreditLimitViolation } from "@/types";
import { CREDIT_LIMITS, SEMESTER_TYPES } from "@/lib/constants";
import { ViolationCard } from "@/components/ViolationCard";

interface CreditLimitViolationsCardProps {
  violations: CreditLimitViolation[];
}

const helpContent = (
  <>
    <p className="font-medium mb-1">How to fix:</p>
    <ul className="list-disc list-inside text-xs space-y-1">
      <li>Regular semesters: Max {CREDIT_LIMITS.REGULAR_SEMESTER} credits</li>
      <li>Summer semester: Max {CREDIT_LIMITS.SUMMER_SEMESTER} credits</li>
      <li>Move some courses to other semesters to reduce load</li>
    </ul>
  </>
);

export function CreditLimitViolationsCard({
  violations,
}: CreditLimitViolationsCardProps) {
  if (!violations || violations.length === 0) {
    return (
      <ViolationCard
        title="No Credit Limit Violations"
        description="All semesters are within the allowed credit limits."
        icon={CheckCircle}
        variant="success"
      />
    );
  }

  return (
    <ViolationCard
      title="Credit Limit Violations"
      description="Some semesters exceed the maximum allowed credit limit."
      icon={XCircle}
      count={violations.length}
      variant="danger"
      helpContent={helpContent}
    >
      <div
        className="space-y-3"
        role="list"
        aria-label="Credit limit violations"
      >
        {violations.map((violation, index) => {
          const exceededBy = violation.credits - violation.max_credits;
          const isSummerSemester = violation.semester === SEMESTER_TYPES.SUMMER;

          return (
            <div
              key={index}
              className="rounded-lg border-2 border-red-300 bg-gradient-to-br from-red-50 to-pink-50 p-4 shadow-sm hover:shadow-md transition-shadow"
              role="listitem"
            >
              <div className="flex items-center justify-between">
                <div className="flex-1">
                  <div className="font-bold text-base text-gray-900 mb-1">
                    Year {violation.year}, Semester {violation.semester}
                    <span className="ml-2 text-xs font-medium text-gray-500">
                      ({isSummerSemester ? "Summer" : "Regular"})
                    </span>
                  </div>
                  <div className="text-sm text-gray-700">
                    <span className="font-bold text-red-700">{violation.credits}</span>
                    <span className="text-gray-500"> / {violation.max_credits} credits</span>
                    <span className="ml-2 text-xs text-red-600">
                      (over by <span className="font-bold">{exceededBy}</span>)
                    </span>
                  </div>
                </div>
                <div className="flex-shrink-0 px-3 py-1.5 bg-red-100 rounded-lg">
                  <span className="text-xl font-bold text-red-700">
                    +{exceededBy}
                  </span>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </ViolationCard>
  );
}
