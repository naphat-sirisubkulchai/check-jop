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
        className="h-min gap-4 border-0 border-l-4 border-l-chula-active bg-gradient-to-br from-white to-chula-soft p-5 shadow-md hover:shadow-lg transition-shadow"
        role="region"
        aria-label={`Total credits: ${result.total_credits} out of ${result.required_credits} required credits`}
      >
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <h3 className="text-sm font-semibold text-gray-600">
              Total Credits
            </h3>
            <p className="text-3xl font-bold text-gray-900" aria-live="polite">
              {result.total_credits} <span className="text-xl text-gray-400">/ {result.required_credits}</span>
            </p>
          </div>
          <div className="p-3 bg-gradient-to-br from-chula-active to-pink-500 rounded-xl shadow-sm">
            <BookOpen className="h-6 w-6 text-white" aria-hidden="true" />
          </div>
        </div>
      </Card>

      {/* GPAX */}
      <Card
        className="h-min gap-4 border-0 border-l-4 border-l-purple-500 bg-gradient-to-br from-white to-purple-100 p-5 shadow-md hover:shadow-lg transition-shadow"
        role="region"
        aria-label={`Current GPAX: ${result.gpax.toFixed(2)}`}
      >
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <h3 className="text-sm font-semibold text-gray-600">GPAX</h3>
            <p className="text-3xl font-bold text-gray-900" aria-live="polite">
              {result.gpax.toFixed(2)}
            </p>
          </div>
          <div className="p-3 bg-gradient-to-br from-purple-500 to-purple-600 rounded-xl shadow-sm">
            <Calculator className="h-6 w-6 text-white" aria-hidden="true" />
          </div>
        </div>
      </Card>

      {/* Categories Progress */}
      <Card
        className="h-min gap-4 border-0 border-l-4 border-l-green-500 bg-gradient-to-br from-white to-green-50 p-5 shadow-md hover:shadow-lg transition-shadow"
        role="region"
        aria-label={`Categories completed: ${satisfiedCategories} out of ${totalCategories} total categories`}
      >
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <h3 className="text-sm font-semibold text-gray-600">
              Categories
            </h3>
            <p className="text-3xl font-bold text-gray-900" aria-live="polite">
              {satisfiedCategories} <span className="text-xl text-gray-400">/ {totalCategories}</span>
            </p>
          </div>
          <div className="p-3 bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl shadow-sm">
            <CheckCircle className="h-6 w-6 text-white" aria-hidden="true" />
          </div>
        </div>
      </Card>

      {/* Issues */}
      <Card
        className="h-min gap-4 border-0 border-l-4 border-l-orange-500 bg-gradient-to-br from-white to-orange-50 p-5 shadow-md hover:shadow-lg transition-shadow"
        role="region"
        aria-label={`Total issues found: ${totalViolations} ${totalViolations === 1 ? 'violation' : 'violations'}`}
      >
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <h3 className="text-sm font-semibold text-gray-600">Issues</h3>
            <p className="text-3xl font-bold text-gray-900" aria-live="polite">
              {totalViolations}
            </p>
          </div>
          <div className="p-3 bg-gradient-to-br from-orange-500 to-red-500 rounded-xl shadow-sm">
            <XCircle className="h-6 w-6 text-white" aria-hidden="true" />
          </div>
        </div>
      </Card>
    </div>
  );
}
