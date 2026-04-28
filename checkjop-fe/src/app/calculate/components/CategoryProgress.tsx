"use client";

import { useState } from "react";
import { CategoryResult } from "@/types";
import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import { ChevronDown, ChevronUp } from "lucide-react";

interface CategoryProgressProps {
  category: CategoryResult;
}

export function CategoryProgress({ category }: CategoryProgressProps) {
  const [expanded, setExpanded] = useState(false);
  const progressPercentage = (category.earned_credits / category.required_credits) * 100;
  const missingCredits = category.required_credits - category.earned_credits;
  const hasMissingCourses = !category.is_satisfied && category.missing_courses?.length > 0;
  const hasViolatedCourses = category.violated_courses?.length > 0;
  const completedWithWarning = category.is_satisfied && hasViolatedCourses;

  const getProgressColor = () => {
    if (completedWithWarning) return "bg-gradient-to-r from-yellow-400 to-orange-400";
    if (category.is_satisfied) return "bg-gradient-to-r from-green-500 to-emerald-600";
    if (category.earned_credits === 0) return "bg-gradient-to-r from-red-500 to-red-600";
    return "bg-gradient-to-r from-yellow-500 to-orange-500";
  };

  const getBadgeStyle = () => {
    if (completedWithWarning) return "bg-gradient-to-r from-yellow-400 to-orange-400 text-white";
    if (category.is_satisfied) return "bg-gradient-to-r from-green-500 to-emerald-600 text-white";
    if (category.earned_credits === 0) return "bg-gradient-to-r from-red-500 to-red-600 text-white";
    return "bg-gradient-to-r from-yellow-500 to-orange-400 text-white";
  };

  return (
    <div className="hover:bg-gray-50 transition-colors">
      <div
        className={cn("flex items-center justify-between p-4", (hasMissingCourses || hasViolatedCourses) && "cursor-pointer")}
        onClick={() => (hasMissingCourses || hasViolatedCourses) && setExpanded(!expanded)}
        role="region"
        aria-label={`${category.category_name}: ${category.earned_credits} of ${category.required_credits} credits completed`}
      >
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <h4 className="text-sm font-bold text-gray-900 truncate">
              {category.category_name}
            </h4>
            {(hasMissingCourses || hasViolatedCourses) && (
              expanded ? <ChevronUp className="h-4 w-4 text-gray-400 shrink-0" /> : <ChevronDown className="h-4 w-4 text-gray-400 shrink-0" />
            )}
          </div>
          <div className="flex items-center gap-2">
            <span className="text-xs font-semibold text-gray-700 shrink-0">
              {category.earned_credits} / {category.required_credits} credits
            </span>
            <div className="relative h-2 w-full flex-1 overflow-hidden rounded-full bg-gray-200 shadow-inner hidden sm:block">
              <div
                className={cn("h-full transition-all duration-500 shadow-sm", getProgressColor())}
                style={{ width: `${Math.min(progressPercentage, 100)}%` }}
                role="progressbar"
                aria-valuenow={category.earned_credits}
                aria-valuemin={0}
                aria-valuemax={category.required_credits}
                aria-label={`Progress: ${Math.round(progressPercentage)}%`}
              />
            </div>
            <span className="text-xs font-semibold text-gray-500 shrink-0">
              {Math.round(progressPercentage)}%
            </span>
          </div>
        </div>
        <div className="ml-3 flex-shrink-0">
          {completedWithWarning ? (
            <Badge className={cn("font-bold text-xs px-2 py-1 rounded-lg shadow-sm whitespace-nowrap", getBadgeStyle())} role="status">
              ⚠ Completed
            </Badge>
          ) : category.is_satisfied ? (
            <Badge className={cn("font-bold text-xs px-2 py-1 rounded-lg shadow-sm whitespace-nowrap", getBadgeStyle())} role="status">
              ✓ Completed
            </Badge>
          ) : category.earned_credits === 0 ? (
            <Badge className={cn("font-bold text-xs px-2 py-1 rounded-lg shadow-sm whitespace-nowrap", getBadgeStyle())} role="status">
              {category.required_credits} Missing
            </Badge>
          ) : (
            <Badge className={cn("font-bold text-xs px-2 py-1 rounded-lg shadow-sm whitespace-nowrap", getBadgeStyle())} role="status">
              {missingCredits} Missing
            </Badge>
          )}
        </div>
      </div>
      {(hasMissingCourses || hasViolatedCourses) && expanded && (
        <div className="px-5 pb-4 space-y-3">
          {hasViolatedCourses && (
            <div>
              <p className="text-xs font-semibold text-orange-600 mb-2">Condition violations (prereq/coreq):</p>
              <div className="flex flex-wrap gap-2">
                {category.violated_courses.map((code) => (
                  <span key={code} className="text-xs bg-orange-50 text-orange-700 border border-orange-300 rounded px-2 py-1 font-mono">
                    {code}
                  </span>
                ))}
              </div>
            </div>
          )}
          {hasMissingCourses && (
            <div>
              <p className="text-xs font-semibold text-gray-500 mb-2">Missing courses:</p>
              <div className="flex flex-wrap gap-2">
                {category.missing_courses
                  .filter((code) => !category.violated_courses?.includes(code))
                  .map((code) => (
                    <span key={code} className="text-xs bg-red-50 text-red-700 border border-red-200 rounded px-2 py-1 font-mono">
                      {code}
                    </span>
                  ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
