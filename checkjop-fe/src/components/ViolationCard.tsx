"use client";

import { Card } from "@/components/ui/card";
import { LucideIcon } from "lucide-react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { HelpCircle } from "lucide-react";
import { ReactNode, useMemo } from "react";

interface ViolationCardProps {
  title: string;
  description?: string;
  icon: LucideIcon;
  count?: number;
  variant: "success" | "warning" | "danger";
  helpContent?: ReactNode;
  children?: ReactNode;
}

const variantStyles = {
  success: {
    border: "border-green-200",
    bg: "bg-gradient-to-br from-green-50 to-emerald-50",
    iconBg: "bg-gradient-to-br from-green-500 to-emerald-600",
    text: "text-green-700",
    badgeBg: "bg-green-100",
    badgeText: "text-green-800",
  },
  warning: {
    border: "border-yellow-200",
    bg: "bg-gradient-to-br from-yellow-50 to-orange-50",
    iconBg: "bg-gradient-to-br from-yellow-500 to-orange-500",
    text: "text-orange-700",
    badgeBg: "bg-yellow-100",
    badgeText: "text-yellow-800",
  },
  danger: {
    border: "border-red-200",
    bg: "bg-gradient-to-br from-red-50 to-pink-50",
    iconBg: "bg-gradient-to-br from-red-500 to-red-600",
    text: "text-red-700",
    badgeBg: "bg-red-100",
    badgeText: "text-red-800",
  },
};

export function ViolationCard({
  title,
  description,
  icon: Icon,
  count,
  variant,
  helpContent,
  children,
}: ViolationCardProps) {
  const styles = variantStyles[variant];

  const renderTitle = useMemo(() => {
    if (variant === "success") {
      return (
        <SuccessMessage
          icon={Icon}
          title={title}
          description={description}
        />
      );
    } else {
      return (
        <div className="flex items-start gap-4">
          <div className={`flex h-12 w-12 items-center justify-center rounded-xl ${styles.iconBg} shadow-md flex-shrink-0`}>
            <Icon className="h-6 w-6 text-white" aria-hidden="true" />
          </div>
          <div className="flex-1">
            <h3 className={`flex items-center gap-2 text-xl font-bold ${styles.text} mb-1`}>
              {title}
              {count !== undefined && (
                <span className={`${styles.badgeBg} ${styles.badgeText} text-base font-bold px-2.5 py-0.5 rounded-lg`}>
                  {count}
                </span>
              )}
            </h3>
            <p className={`text-sm ${styles.text} font-medium`}>
              {description}
            </p>
          </div>
          {helpContent && (
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <button
                    className="text-gray-400 hover:text-gray-600 transition-colors flex-shrink-0"
                    aria-label={`Help information about ${title.toLowerCase()}`}
                  >
                    <HelpCircle className="h-5 w-5" />
                  </button>
                </TooltipTrigger>
                <TooltipContent side="right" className="max-w-xs">
                  {helpContent}
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          )}
        </div>
      )
    }
  }, [variant, title, description, count, styles, helpContent]);

  return (
    <Card className={`${styles.bg} border-2 ${styles.border} p-6 gap-6 shadow-md`}>
      {renderTitle}
      {children}
    </Card>
  );
}

interface SuccessMessageProps {
  icon: LucideIcon;
  title: string;
  description?: string;
}

export function SuccessMessage({ icon: Icon, title, description }: SuccessMessageProps) {
  return (
    <div className="flex items-start gap-4" role="status" aria-live="polite">
      <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-green-500 to-emerald-600 shadow-md flex-shrink-0">
        <Icon className="h-6 w-6 text-white" aria-hidden="true" />
      </div>
      <div className="flex-1">
        <h3 className="text-xl font-bold text-green-800 mb-1">
          {title}
        </h3>
        <p className="text-sm text-green-700 font-medium">
          {description}
        </p>
      </div>
    </div>
  );
}
