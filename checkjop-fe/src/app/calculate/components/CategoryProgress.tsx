"use client";

import { CategoryResult } from "@/types";
import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";

interface CategoryProgressProps {
  category: CategoryResult;
}

export function CategoryProgress({ category }: CategoryProgressProps) {
  const progressPercentage = (category.earned_credits / category.required_credits) * 100;
  const missingCredits = category.required_credits - category.earned_credits;

  // Determine progress bar color based on completion status
  const getProgressColor = () => {
    if (category.is_satisfied) {
      return "bg-gradient-to-r from-green-500 to-emerald-600"; // Green for completed
    } else if (category.earned_credits === 0) {
      return "bg-gradient-to-r from-red-500 to-red-600"; // Red for not started
    } else {
      return "bg-gradient-to-r from-yellow-500 to-orange-500"; // Orange for in progress
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
    <div
      className="flex items-center justify-between p-5 hover:bg-gray-50 transition-colors"
      role="region"
      aria-label={`${category.category_name}: ${category.earned_credits} of ${category.required_credits} credits completed`}
    >
      <div className="flex-1">
        <h4 className="text-base font-bold text-gray-900 truncate mb-2">
          {category.category_name}
        </h4>
        <div className="flex items-center gap-4">
          <span className="w-28 flex-shrink-0 text-sm font-semibold text-gray-700">
            {category.earned_credits} / {category.required_credits} credits
          </span>
          <div className="relative h-3 w-full max-w-xl flex-1 overflow-hidden rounded-full bg-gray-200 shadow-inner">
            <div
              className={cn("h-full transition-all duration-500 shadow-sm", getProgressColor())}
              style={{ width: `${progressPercentage}%` }}
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
  );
}
