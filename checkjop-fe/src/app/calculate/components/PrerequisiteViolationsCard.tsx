"use client";

import {
  AlertTriangle,
  CheckCircle,
  ChevronDown,
  ChevronUp,
} from "lucide-react";
import { PrerequisiteViolation } from "@/types";
import { ViolationCard, SuccessMessage } from "@/components/ViolationCard";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { useAnalytics } from "@/hooks/useAnalytics";

interface PrerequisiteViolationsCardProps {
  violations: PrerequisiteViolation[];
}

const helpContent = (
  <>
    <p className="font-medium mb-1">How to fix:</p>
    <ul className="list-disc list-inside text-xs space-y-1">
      <li>Ensure prerequisites are taken before the course</li>
      <li>Corequisites must be taken in the same semester</li>
      <li>Adjust your study plan to resolve violations</li>
    </ul>
  </>
);

export function PrerequisiteViolationsCard({
  violations,
}: PrerequisiteViolationsCardProps) {
  const [expandedItems, setExpandedItems] = useState<Set<number>>(new Set());
  const [showAll, setShowAll] = useState(false);
  const { trackEvent } = useAnalytics();

  const toggleItem = (index: number) => {
    const newExpanded = new Set(expandedItems);
    const isExpanding = !newExpanded.has(index);

    if (newExpanded.has(index)) {
      newExpanded.delete(index);
    } else {
      newExpanded.add(index);
    }
    setExpandedItems(newExpanded);

    // Track expansion/collapse
    trackEvent(isExpanding ? "violation_expanded" : "violation_collapsed", {
      type: "prerequisite",
      courseCode: violations[index]?.course_code,
    });
  };

  const toggleAll = () => {
    if (showAll) {
      setExpandedItems(new Set());
    } else {
      setExpandedItems(new Set(violations.map((_, i) => i)));
    }
    setShowAll(!showAll);

    trackEvent(showAll ? "violation_collapsed" : "violation_expanded", {
      type: "prerequisite",
      action: "all",
      count: violations.length,
    });
  };

  if (!violations || violations.length === 0) {
    return (
      <ViolationCard
        title="No Prerequisite or Corequisite Violations"
        description="All course prerequisites and corequisites are properly satisfied."
        icon={CheckCircle}
        variant="success"
      />
    );
  }

  return (
    <ViolationCard
      title="Prerequisite & Corequisite Violations"
      description="Some courses have prerequisite or corequisite requirements that need to be resolved."
      icon={AlertTriangle}
      count={violations.length}
      variant="warning"
      helpContent={helpContent}
    >
      <div className="flex flex-col gap-2 -mt-4">
        {violations.length > 3 && (
          <div className="flex justify-end">
            <Button
              variant="ghost"
              size="sm"
              onClick={toggleAll}
              className="text-sm hover:bg-yellow-100 hover:text-yellow-800 font-medium"
            >
              {showAll ? (
                <>
                  <ChevronUp className="h-4 w-4 mr-1.5" />
                  Collapse All
                </>
              ) : (
                <>
                  <ChevronDown className="h-4 w-4 mr-1.5" />
                  Expand All
                </>
              )}
            </Button>
          </div>
        )}
        <div
          className="space-y-3"
          role="list"
          aria-label="Prerequisite and corequisite violations"
        >
          {violations.map((violation, index) => {
            const isExpanded = expandedItems.has(index);
            const hasDetails =
              violation.missing_prereqs.length > 0 ||
              violation.prereqs_taken_in_wrong_term.length > 0 ||
              violation.missing_coreqs.length > 0 ||
              violation.coreqs_taken_in_wrong_term.length > 0;

            return (
              <div
                key={index}
                className="rounded-lg border-2 border-yellow-300 bg-gradient-to-r from-yellow-50 to-orange-50 p-4 shadow-sm hover:shadow-md transition-shadow"
                role="listitem"
              >
                <button
                  onClick={() => toggleItem(index)}
                  className="w-full flex items-center justify-between text-left"
                  aria-expanded={isExpanded}
                >
                  <div className="font-bold text-lg text-gray-900 truncate flex-1">
                    Course:{" "}
                    <span className="text-orange-700 font-mono">
                      {violation.course_code}
                    </span>
                  </div>
                  {hasDetails && (
                    <div className="ml-3 flex-shrink-0 p-1 hover:bg-yellow-100 rounded-md transition-colors">
                      {isExpanded ? (
                        <ChevronUp className="h-5 w-5 text-yellow-700" />
                      ) : (
                        <ChevronDown className="h-5 w-5 text-yellow-700" />
                      )}
                    </div>
                  )}
                </button>

                {isExpanded && hasDetails && (
                  <div className="space-y-2 mt-3 pl-4 border-l-4 border-yellow-400">
                    {violation.missing_prereqs.length > 0 && (
                      <div className="text-sm text-gray-800">
                        <span className="font-bold text-orange-700">
                          • Missing Prerequisites :
                        </span>{" "}
                        <span className="font-medium block sm:inline pl-4 sm:pl-0">
                          {violation.missing_prereqs.join(", ")}
                        </span>
                      </div>
                    )}

                    {violation.prereqs_taken_in_wrong_term.length > 0 && (
                      <div className="text-sm text-gray-800">
                        <span className="font-bold text-orange-700">
                          • Wrong Term (Prerequisites) :
                        </span>{" "}
                        <span className="font-medium block sm:inline pl-4 sm:pl-0">
                          {violation.prereqs_taken_in_wrong_term.join(", ")}
                        </span>
                      </div>
                    )}

                    {violation.missing_coreqs.length > 0 && (
                      <div className="text-sm text-gray-800">
                        <span className="font-bold text-orange-700">
                          • Missing Corequisites :
                        </span>{" "}
                        <span className="font-medium block sm:inline pl-4 sm:pl-0">
                          {violation.missing_coreqs.join(", ")}
                        </span>
                      </div>
                    )}

                    {violation.coreqs_taken_in_wrong_term.length > 0 && (
                      <div className="text-sm text-gray-800">
                        <span className="font-bold text-orange-700">
                          • Wrong Term (Corequisites) :
                        </span>{" "}
                        <span className="font-medium block sm:inline pl-4 sm:pl-0">
                          {violation.coreqs_taken_in_wrong_term.join(", ")}
                        </span>
                      </div>
                    )}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </ViolationCard>
  );
}
