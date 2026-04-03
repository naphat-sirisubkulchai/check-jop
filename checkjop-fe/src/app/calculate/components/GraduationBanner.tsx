"use client";

import { Check, BookOpen, XCircle, AlertCircle } from "lucide-react";
import { GraduationResult } from "@/types";

interface GraduationBannerProps {
  result: GraduationResult;
}

export function GraduationBanner({ result }: GraduationBannerProps) {
  const { can_graduate, total_credits, required_credits } = result;

  // Determine banner variant based on graduation status
  const getBannerConfig = () => {
    if (can_graduate) {
      return {
        bgColor: "bg-gradient-to-r from-green-50 to-emerald-50",
        borderColor: "border-green-200",
        iconBg: "bg-gradient-to-br from-green-500 to-emerald-600",
        titleColor: "text-green-800",
        textColor: "text-green-700",
        badgeBg: "bg-gradient-to-r from-green-500 to-emerald-600",
        badgeText: "text-white",
        title: "All requirements completed!",
        subtitle: "Congratulations! You are eligible to",
        highlightText: "graduate",
        badgeLabel: "Eligible",
        additionalText: "You can apply for graduation.",
        showIllustration: true,
      };
    } else {
      const totalViolations =
        result.prerequisite_violations.length +
        result.credit_limit_violations.length;
      const unsatisfiedCategories = result.category_results.filter(
        (cat) => !cat.is_satisfied
      ).length;

      return {
        bgColor: totalViolations > 0
          ? "bg-gradient-to-r from-yellow-50 to-orange-50"
          : "bg-gradient-to-r from-chula-soft/50 to-pink-100",
        borderColor: totalViolations > 0
          ? "border-yellow-200"
          : "border-chula-active/30",
        iconBg: totalViolations > 0
          ? "bg-gradient-to-br from-yellow-500 to-orange-600"
          : "bg-gradient-to-br from-chula-active to-pink-500",
        titleColor: totalViolations > 0
          ? "text-orange-800"
          : "text-chula-active",
        textColor: totalViolations > 0
          ? "text-orange-700"
          : "text-pink-700",
        badgeBg: totalViolations > 0
          ? "bg-gradient-to-r from-yellow-500 to-orange-600"
          : "bg-gradient-to-r from-chula-active to-pink-500",
        badgeText: "text-white",
        title: "Requirements incomplete",
        subtitle: `You have ${unsatisfiedCategories} unsatisfied ${unsatisfiedCategories === 1 ? "category" : "categories"} and ${totalViolations}`,
        highlightText: totalViolations === 1 ? "violation" : "violations",
        badgeLabel: "Not Eligible",
        additionalText: "Review your plan and make adjustments.",
        showIllustration: false,
      };
    }
  };

  const config = getBannerConfig();

  return (
    <div
      className={`w-full ${config.bgColor} border-2 ${config.borderColor} rounded-xl overflow-hidden shadow-lg`}
      role="status"
      aria-live="polite"
    >
      <div className="flex items-start p-6 gap-4">
        {/* Icon circle */}
        <div className={`flex-shrink-0 w-12 h-12 ${config.iconBg} rounded-full flex items-center justify-center shadow-md`}>
          {can_graduate ? (
            <Check className="w-7 h-7 text-white" strokeWidth={3} />
          ) : (
            config.badgeLabel === "Not Eligible" && result.prerequisite_violations.length + result.credit_limit_violations.length > 0 ? (
              <AlertCircle className="w-7 h-7 text-white" strokeWidth={3} />
            ) : (
              <XCircle className="w-7 h-7 text-white" strokeWidth={3} />
            )
          )}
        </div>

        {/* Main content */}
        <div className="flex-1 min-w-0">
          <h1 className={`${config.titleColor} font-bold text-xl leading-tight`}>
            {config.title}
          </h1>
          <p className={`${config.textColor} text-base mt-2 font-medium`}>
            {config.subtitle} <span className="font-bold">{config.highlightText}</span>.
          </p>
          <p className={`${config.textColor} text-sm flex items-center gap-2 mt-2 flex-wrap`}>
            Your academic plan fulfills {can_graduate ? "all" : "some"} requirements with
            <span className="inline-flex items-center gap-1.5 font-semibold">
              <BookOpen className={`w-4 h-4 ${config.textColor}`} />
              {total_credits} / {required_credits} credits
            </span>
          </p>

          {/* Eligible badge and text */}
          <div className="flex items-center gap-3 mt-4 flex-wrap">
            <div className={`inline-flex items-center gap-2 ${config.badgeBg} ${config.badgeText} text-sm font-bold px-4 py-2 rounded-lg shadow-md`}>
              <Check className="w-4 h-4" strokeWidth={3} />
              {config.badgeLabel}
            </div>
            <span className="text-gray-600 text-sm font-medium">{config.additionalText}</span>
          </div>
        </div>

        {/* Graduation illustration - only show if eligible */}
        {config.showIllustration && (
          <div className="flex-shrink-0 hidden sm:flex items-end gap-1 mr-2">
            {/* Graduation cap */}
            <svg width="48" height="40" viewBox="0 0 48 40" fill="none" className="-mb-1">
              {/* Cap base */}
              <polygon points="24,8 4,18 24,28 44,18" fill="#1b5e20" />
              {/* Cap top */}
              <polygon points="24,4 4,14 24,24 44,14" fill="#2e7d32" />
              {/* Tassel string */}
              <line x1="44" y1="18" x2="44" y2="28" stroke="#ffc107" strokeWidth="2" />
              {/* Tassel end */}
              <rect x="41" y="28" width="6" height="8" rx="1" fill="#ffc107" />
              {/* Button on top */}
              <circle cx="24" cy="6" r="2" fill="#ffc107" />
            </svg>

            {/* Diploma */}
            <svg width="32" height="36" viewBox="0 0 32 36" fill="none" className="ml-1">
              {/* Diploma roll */}
              <rect x="4" y="4" width="20" height="28" rx="2" fill="#f5e6d3" stroke="#c9a66b" strokeWidth="1" />
              {/* Roll shadow */}
              <ellipse cx="14" cy="4" rx="10" ry="3" fill="#e8d4b8" />
              <ellipse cx="14" cy="32" rx="10" ry="3" fill="#d4b896" />
              {/* Ribbon */}
              <path d="M14 18 L10 28 L14 24 L18 28 L14 18" fill="#c62828" />
              {/* Seal */}
              <circle cx="14" cy="18" r="4" fill="#c62828" />
              <circle cx="14" cy="18" r="2" fill="#ef5350" />
            </svg>
          </div>
        )}
      </div>
    </div>
  );
}
