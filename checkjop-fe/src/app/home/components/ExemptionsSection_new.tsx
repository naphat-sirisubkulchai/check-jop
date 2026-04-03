"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { X, Plus } from "lucide-react";
import { useAppStore } from "@/store/appStore";

interface ExemptionsSectionProps {
  className?: string;
}

/**
 * Compact component for managing course exemptions
 */
export default function ExemptionsSection({
  className = "",
}: ExemptionsSectionProps) {
  const { exemptions, addExemption, removeExemption } = useAppStore();

  const [exemptionInput, setExemptionInput] = useState("");
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
    <div className={`space-y-3 ${className}`}>
      {/* Header & Input Row */}
      <div className="flex items-center justify-between gap-4">
        <h3 className="text-base font-semibold text-gray-900 whitespace-nowrap">
          Course Exemptions
        </h3>

        {/* Input & Add Button */}
        <div className="flex items-start gap-2 flex-1 max-w-sm">
          <div className="flex-1">
            <Input
              value={exemptionInput}
              onChange={handleInputChange}
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  handleAddExemption();
                }
              }}
              placeholder="Enter course code (e.g., 2301101)"
              className={`text-sm ${
                error ? "border-red-300 focus-visible:ring-red-500" : ""
              }`}
              aria-invalid={!!error}
              aria-describedby={error ? "exemption-error" : undefined}
            />
            {error && (
              <p id="exemption-error" className="text-xs text-red-600 mt-1">
                {error}
              </p>
            )}
          </div>
          <Button
            onClick={handleAddExemption}
            size="sm"
            className="bg-chula-active hover:bg-chula-active/90"
            disabled={!exemptionInput.trim()}
          >
            <Plus className="h-4 w-4 mr-1" />
            Add
          </Button>
        </div>
      </div>

      {/* Exemptions List - Separate Row */}
      {exemptions.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {exemptions.map((code) => (
            <Badge
              key={code}
              variant="secondary"
              className="flex items-center gap-1.5 bg-white border border-chula-active/20 hover:border-chula-active/40 px-2.5 py-1.5 text-xs font-medium text-gray-800 transition-all"
            >
              <span className="font-mono">{code}</span>
              <button
                onClick={() => removeExemption(code)}
                className="text-gray-400 hover:text-red-600 transition-colors"
                aria-label={`Remove ${code} from exemptions`}
              >
                <X className="h-3 w-3" />
              </button>
            </Badge>
          ))}
        </div>
      )}
    </div>
  );
}
