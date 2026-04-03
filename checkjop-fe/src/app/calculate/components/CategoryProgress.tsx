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

  const getProgressColor = () => {
    if (category.is_satisfied) {
      return "bg-gradient-to-r from-green-500 to-emerald-600";
    } else if (category.earned_credits === 0) {
      return "bg-gradient-to-r from-red-500 to-red-600";
    } else {
      return "bg-gradient-to-r from-yellow-500 to-orange-500";
    }
  };

  const getBadgeStyle = () => {
    if (category.is_satisfied) {
      return "bg-gradient-to-r from-green-500 to-emerald-600 text-white";
    } else if (category.earned_credits === 0) {
      return "bg-gradient-to-r from-red-500 to-red-600 text-white";
    } else {
      return "bg-gradient-to-r from-yellow-500 to-orange-400 text-white";
    }
  };

  return (
    <div className="hover:bg-gray-50 transition-colors">
      <div
        className={cn("flex items-center justify-between p-5", hasMissingCourses && "cursor-pointer")}
        onClick={() => hasMissingCourses && setExpanded(!expanded)}
        role="region"
        aria-label={`${category.category_name}: ${category.earned_credits} of ${category.required_credits} credits completed`}
      >
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <h4 className="text-base font-bold text-gray-900 truncate">
              {category.category_name}
            </h4>
            {hasMissingCourses && (
              expanded ? <ChevronUp className="h-4 w-4 text-gray-400" /> : <ChevronDown className="h-4 w-4 text-gray-400" />
            )}
          </div>
          <div className="flex items-center gap-4">
            <span className="w-28 flex-shrink-0 text-sm font-semibold text-gray-700">
              {category.earned_credits} / {category.required_credits} credits
            </span>
            <div className="relative h-3 w-full max-w-xl flex-1 overflow-hidden rounded-full bg-gray-200 shadow-inner">
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
            <span className="text-xs font-semibold text-gray-500 w-12 text-right">
              {Math.round(progressPercentage)}%
            </span>
          </div>
        </div>
        <div className="ml-6 flex-shrink-0">
          {category.is_satisfied ? (
            <Badge className={cn("font-bold text-xs min-w-24 py-2 rounded-lg shadow-sm", getBadgeStyle())} role="status">
              ✓ Completed
            </Badge>
          ) : category.earned_credits === 0 ? (
            <Badge className={cn("font-bold text-xs min-w-24 py-2 rounded-lg shadow-sm", getBadgeStyle())} role="status">
              {category.required_credits} Missing
            </Badge>
          ) : (
            <Badge className={cn("font-bold text-xs min-w-24 py-2 rounded-lg shadow-sm", getBadgeStyle())} role="status">
              {missingCredits} Missing
            </Badge>
          )}
        </div>
      </div>
      {hasMissingCourses && expanded && (
        <div className="px-5 pb-4">
          <p className="text-xs font-semibold text-gray-500 mb-2">วิชาที่ยังไม่ได้ลง:</p>
          <div className="flex flex-wrap gap-2">
            {category.missing_courses.map((code) => (
              <span key={code} className="text-xs bg-red-50 text-red-700 border border-red-200 rounded px-2 py-1 font-mono">
                {code}
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
