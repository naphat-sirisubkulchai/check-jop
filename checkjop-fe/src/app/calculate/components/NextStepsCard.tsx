"use client";

import { Card } from "@/components/ui/card";
import { Lightbulb, CheckCircle2, AlertCircle } from "lucide-react";
import { GraduationResult } from "@/types";

interface NextStepsCardProps {
  result: GraduationResult;
}

export function NextStepsCard({ result }: NextStepsCardProps) {
  const hasUnsatisfiedCategories = result.category_results.some(
    (cat) => !cat.is_satisfied
  );
  const hasPrereqViolations = result.prerequisite_violations.length > 0;
  const hasCreditLimitViolations = result.credit_limit_violations.length > 0;

  return (
    <Card className={`border-2 ${result.can_graduate ? 'border-green-300 bg-gradient-to-br from-green-50 to-emerald-50' : 'border-chula-active/30 bg-gradient-to-br from-chula-soft/50 to-pink-50'} p-6 shadow-lg`}>
      <div className="flex items-start space-x-4">
        <div className="flex-shrink-0">
          <div className={`flex h-14 w-14 items-center justify-center rounded-xl shadow-md ${result.can_graduate ? 'bg-gradient-to-br from-green-500 to-emerald-600' : 'bg-gradient-to-br from-chula-active to-pink-500'}`}>
            {result.can_graduate ? (
              <CheckCircle2 className="h-7 w-7 text-white" aria-hidden="true" />
            ) : (
              <Lightbulb className="h-7 w-7 text-white" aria-hidden="true" />
            )}
          </div>
        </div>
        <div className="flex-1">
          <h3 className={`mb-3 text-2xl font-bold ${result.can_graduate ? 'bg-gradient-to-r from-green-700 to-emerald-700' : 'bg-gradient-to-r from-chula-active to-pink-500'} bg-clip-text text-transparent`}>
            {result.can_graduate ? "🎓 Ready to Graduate!" : "📋 Next Steps"}
          </h3>
          <div className="space-y-3">
            {result.can_graduate ? (
              <div>
                <p className="mb-3 text-base font-medium text-gray-800">
                  Congratulations! You&apos;ve met all graduation requirements.
                  Here&apos;s what to do next:
                </p>
                <ul className="space-y-2">
                  <li className="flex items-start gap-2 text-sm text-gray-700">
                    <span className="mt-0.5 h-1.5 w-1.5 rounded-full bg-green-600 flex-shrink-0"></span>
                    <span>Apply for graduation through the registrar&apos;s office</span>
                  </li>
                  <li className="flex items-start gap-2 text-sm text-gray-700">
                    <span className="mt-0.5 h-1.5 w-1.5 rounded-full bg-green-600 flex-shrink-0"></span>
                    <span>Double-check all course grades are finalized</span>
                  </li>
                  <li className="flex items-start gap-2 text-sm text-gray-700">
                    <span className="mt-0.5 h-1.5 w-1.5 rounded-full bg-green-600 flex-shrink-0"></span>
                    <span>Plan your graduation ceremony attendance</span>
                  </li>
                  <li className="flex items-start gap-2 text-sm text-gray-700">
                    <span className="mt-0.5 h-1.5 w-1.5 rounded-full bg-green-600 flex-shrink-0"></span>
                    <span>Your final GPAX: <span className="font-bold text-green-700">{result.gpax.toFixed(2)}</span></span>
                  </li>
                </ul>
              </div>
            ) : (
              <div>
                <p className="mb-3 text-base font-medium text-gray-800">
                  To complete your graduation requirements:
                </p>
                <ul className="space-y-2">
                  {hasUnsatisfiedCategories && (
                    <li className="flex items-start gap-2 text-sm text-gray-700">
                      <AlertCircle className="h-4 w-4 text-chula-active flex-shrink-0 mt-0.5" />
                      <span>Complete missing category requirements</span>
                    </li>
                  )}
                  {hasPrereqViolations && (
                    <li className="flex items-start gap-2 text-sm text-gray-700">
                      <AlertCircle className="h-4 w-4 text-yellow-600 flex-shrink-0 mt-0.5" />
                      <span>
                        Resolve <span className="font-bold text-yellow-700">{result.prerequisite_violations.length}</span>{" "}
                        prerequisite/corequisite violation(s)
                      </span>
                    </li>
                  )}
                  {hasCreditLimitViolations && (
                    <li className="flex items-start gap-2 text-sm text-gray-700">
                      <AlertCircle className="h-4 w-4 text-red-600 flex-shrink-0 mt-0.5" />
                      <span>
                        Address <span className="font-bold text-red-700">{result.credit_limit_violations.length}</span> credit
                        limit violation(s)
                      </span>
                    </li>
                  )}
                  <li className="flex items-start gap-2 text-sm text-gray-700">
                    <span className="mt-0.5 h-1.5 w-1.5 rounded-full bg-gray-400 flex-shrink-0"></span>
                    <span>Consider meeting with your academic advisor</span>
                  </li>
                  <li className="flex items-start gap-2 text-sm text-gray-700">
                    <span className="mt-0.5 h-1.5 w-1.5 rounded-full bg-gray-400 flex-shrink-0"></span>
                    <span>Update your study plan and recalculate</span>
                  </li>
                </ul>
              </div>
            )}
          </div>
        </div>
      </div>
    </Card>
  );
}
