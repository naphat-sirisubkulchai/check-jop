"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { ChevronDown, ChevronUp, X, CheckCircle2 } from "lucide-react";
import { useAppStore } from "@/store/appStore";

interface ExemptionsSectionProps {
  defaultOpen?: boolean;
  className?: string;
}

/**
 * Reusable component for managing course exemptions
 *
 * @example
 * ```tsx
 * <ExemptionsSection defaultOpen={true} />
 * ```
 */
export default function ExemptionsSection({
  defaultOpen = false,
  className = "",
}: ExemptionsSectionProps) {
  const { exemptions, addExemption, removeExemption, clearExemptions } =
    useAppStore();

  const [exemptionInput, setExemptionInput] = useState("");
  const [isOpen, setIsOpen] = useState(defaultOpen);
  const [error, setError] = useState("");

  const handleAddExemption = () => {
    const trimmedInput = exemptionInput.trim().toUpperCase();

    if (!trimmedInput) {
      setError("Please enter a course code");
      return;
    }

    if (exemptions.includes(trimmedInput)) {
      setError("This course is already in your exemptions");
      return;
    }

    addExemption(trimmedInput);
    setExemptionInput("");
    setError("");
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setExemptionInput(e.target.value);
    if (error) setError("");
  };

  return (
    <div
      className={`mb-4 bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden transition-all duration-200 hover:shadow-md ${className}`}
    >
      {/* Header - Collapsible */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full flex items-center justify-between px-5 py-3 hover:bg-gray-50/80 transition-colors cursor-pointer group"
        aria-expanded={isOpen}
        aria-controls="exemptions-content"
      >
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2">
            <CheckCircle2 className="h-5 w-5 text-chula-active" />
            <h3 className="text-lg font-semibold text-gray-900">
              Course Exemptions
            </h3>
          </div>
          {exemptions.length > 0 && (
            <Badge className="rounded-full bg-chula-soft text-chula-active font-medium px-3 py-1">
              {exemptions.length}
            </Badge>
          )}
        </div>
        <div className="flex items-center gap-2">
          {exemptions.length > 0 && !isOpen && (
            <span className="text-xs text-gray-500 hidden sm:block">
              {exemptions.slice(0, 3).join(", ")}
              {exemptions.length > 3 && "..."}
            </span>
          )}
          {isOpen ? (
            <ChevronUp className="h-5 w-5 text-gray-400 group-hover:text-gray-600 transition-colors" />
          ) : (
            <ChevronDown className="h-5 w-5 text-gray-400 group-hover:text-gray-600 transition-colors" />
          )}
        </div>
      </button>

      {/* Content - Collapsible with Animation */}
      <div
        id="exemptions-content"
        className={`transition-all duration-300 ease-in-out ${
          isOpen ? "max-h-[1000px] opacity-100" : "max-h-0 opacity-0"
        }`}
      >
        <div className="px-5 pb-5 border-t border-gray-100">
          <p className="text-sm text-gray-600 mb-4 mt-4">
            Add courses you&apos;ve already completed or been exempted from. These will be excluded from your graduation requirements.
          </p>

          {/* Input Section */}
          <div className="space-y-3 mb-4">
            <div className="flex flex-col sm:flex-row gap-2">
              <div className="flex-1">
                <Input
                  value={exemptionInput}
                  onChange={handleInputChange}
                  onKeyDown={(e) => {
                    if (e.key === "Enter") {
                      handleAddExemption();
                    }
                  }}
                  placeholder="Enter course code (e.g., CS102, MATH101)"
                  className={`transition-colors ${
                    error ? "border-red-300 focus-visible:ring-red-500" : ""
                  }`}
                  aria-invalid={!!error}
                  aria-describedby={error ? "exemption-error" : undefined}
                />
                {error && (
                  <p id="exemption-error" className="text-xs text-red-600 mt-1.5 ml-1">
                    {error}
                  </p>
                )}
              </div>
              <div className="flex gap-2">
                <Button
                  onClick={handleAddExemption}
                  variant="default"
                  size="sm"
                  className="bg-chula-active hover:bg-chula-active/90 flex-1 sm:flex-none"
                  disabled={!exemptionInput.trim()}
                >
                  Add Course
                </Button>
                {exemptions.length > 0 && (
                  <Button
                    onClick={(e) => {
                      e.stopPropagation();
                      clearExemptions();
                    }}
                    variant="outline"
                    size="sm"
                    className="text-red-600 hover:text-red-700 hover:bg-red-50 border-red-200 flex-1 sm:flex-none"
                  >
                    Clear All
                  </Button>
                )}
              </div>
            </div>
          </div>

          {/* Exemptions List */}
          {exemptions.length > 0 ? (
            <div className="space-y-2">
              <p className="text-xs font-medium text-gray-700 uppercase tracking-wide">
                Exempted Courses ({exemptions.length})
              </p>
              <div className="flex flex-wrap gap-2">
                {exemptions.map((code) => (
                  <Badge
                    key={code}
                    variant="secondary"
                    className="flex items-center gap-2 bg-gradient-to-br from-chula-soft to-pink-50 hover:from-chula-soft/80 hover:to-pink-100 border border-chula-active/20 px-3 py-2 text-sm font-medium text-gray-800 transition-all duration-200 hover:shadow-sm"
                  >
                    <span className="font-mono">{code}</span>
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        removeExemption(code);
                      }}
                      className="text-gray-500 hover:text-red-600 transition-colors rounded-full hover:bg-white/50 p-0.5"
                      aria-label={`Remove ${code} from exemptions`}
                    >
                      <X className="h-3.5 w-3.5" />
                    </button>
                  </Badge>
                ))}
              </div>
            </div>
          ) : (
            <div className="text-center py-6 px-4 bg-gray-50 rounded-lg border-2 border-dashed border-gray-200">
              <CheckCircle2 className="h-10 w-10 text-gray-300 mx-auto mb-2" />
              <p className="text-sm text-gray-500 font-medium">
                No exemptions added yet
              </p>
              <p className="text-xs text-gray-400 mt-1">
                Add course codes above to track completed courses
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
