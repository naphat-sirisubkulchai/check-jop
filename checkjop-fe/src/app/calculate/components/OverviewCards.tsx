import { Card } from "@/components/ui/card";
import { BookOpen, CheckCircle, XCircle, Calculator } from "lucide-react";
import { GraduationResult } from "@/types";

interface OverviewCardsProps {
  result: GraduationResult;
  satisfiedCategories: number;
  totalCategories: number;
}

export function OverviewCards({
  result,
  satisfiedCategories,
  totalCategories,
}: OverviewCardsProps) {
  const totalViolations =
    (result.prerequisite_violations?.length || 0) +
    (result.credit_limit_violations?.length || 0);

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {/* Total Credits */}
      <Card
        className="h-min gap-4 border border-gray-100 bg-white p-5 shadow-sm hover:shadow-md transition-shadow"
        role="region"
        aria-label={`Total credits: ${result.total_credits} out of ${result.required_credits} required credits`}
      >
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <h3 className="text-sm font-medium text-gray-500">Total Credits</h3>
            <p className="text-3xl font-bold text-gray-900" aria-live="polite">
              {result.total_credits} <span className="text-xl font-normal text-gray-400">/ {result.required_credits}</span>
            </p>
          </div>
          <BookOpen className="h-6 w-6 text-gray-300" aria-hidden="true" />
        </div>
      </Card>

      {/* GPAX */}
      <Card
        className="h-min gap-4 border border-gray-100 bg-white p-5 shadow-sm hover:shadow-md transition-shadow"
        role="region"
        aria-label={`Current GPAX: ${result.gpax.toFixed(2)}`}
      >
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <h3 className="text-sm font-medium text-gray-500">GPAX</h3>
            <p className="text-3xl font-bold text-gray-900" aria-live="polite">
              {result.gpax.toFixed(2)}
            </p>
          </div>
          <Calculator className="h-6 w-6 text-gray-300" aria-hidden="true" />
        </div>
      </Card>

      {/* Categories Progress */}
      <Card
        className="h-min gap-4 border border-gray-100 bg-white p-5 shadow-sm hover:shadow-md transition-shadow"
        role="region"
        aria-label={`Categories completed: ${satisfiedCategories} out of ${totalCategories} total categories`}
      >
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <h3 className="text-sm font-medium text-gray-500">Categories</h3>
            <p className="text-3xl font-bold text-gray-900" aria-live="polite">
              {satisfiedCategories} <span className="text-xl font-normal text-gray-400">/ {totalCategories}</span>
            </p>
          </div>
          <CheckCircle className="h-6 w-6 text-gray-300" aria-hidden="true" />
        </div>
      </Card>

      {/* Issues */}
      <Card
        className="h-min gap-4 border border-gray-100 bg-white p-5 shadow-sm hover:shadow-md transition-shadow"
        role="region"
        aria-label={`Total issues found: ${totalViolations} ${totalViolations === 1 ? 'violation' : 'violations'}`}
      >
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <h3 className="text-sm font-medium text-gray-500">Issues</h3>
            <p className="text-3xl font-bold text-gray-900" aria-live="polite">
              {totalViolations}
            </p>
          </div>
          <XCircle className="h-6 w-6 text-gray-300" aria-hidden="true" />
        </div>
      </Card>
    </div>
  );
}
