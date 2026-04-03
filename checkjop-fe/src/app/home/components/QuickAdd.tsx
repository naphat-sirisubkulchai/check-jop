"use client";

import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Plan } from "@/types";
import { useState, useRef, useMemo, useEffect, useCallback } from "react";
import { useAppStore } from "@/store/appStore";
import { uniqueCoursesByCode } from "@/utils";
import { Plus } from "lucide-react";

interface QuickAddProps {
  sem: number;
  yearOfStudy: number;
  academicYear: number;
}

export function QuickAdd({ sem, yearOfStudy, academicYear }: QuickAddProps) {
  const { courses, addCoursePlan, studyPlan } = useAppStore();
  const [query, setQuery] = useState("");
  const [isOpen, setIsOpen] = useState(false);
  const [highlightIdx, setHighlightIdx] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Course codes that are planned and NOT failed (F) — these are truly "taken"
  const nonFailedPlanned = useMemo(() => {
    return new Set(
      studyPlan.filter((p) => p.grade !== "F").map((p) => p.course_code)
    );
  }, [studyPlan]);

  // Filter courses
  const filtered = useMemo(() => {
    const q = query.toLowerCase().trim();
    if (!q) return [];
    return uniqueCoursesByCode(courses)
      .filter(
        (c) =>
          !nonFailedPlanned.has(c.code) &&
          (c.code?.toLowerCase().includes(q) ||
            c.name_en?.toLowerCase().includes(q) ||
            c.name_th?.toLowerCase().includes(q))
      )
      .slice(0, 8);
  }, [courses, query, nonFailedPlanned]);

  // Reset highlight when results change
  useEffect(() => {
    setHighlightIdx(0);
  }, [filtered.length]);

  // Close dropdown when clicking outside
  useEffect(() => {    
    const handleClickOutside = (e: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(e.target as Node) &&
        inputRef.current &&
        !inputRef.current.contains(e.target as Node)
      ) {
        setIsOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleSelect = useCallback(
    (course: any) => {
      addCoursePlan({
        course_code: course.code,
        academicYear: academicYear,
        semester: sem,
        yearOfStudy: yearOfStudy,
        credits: course.credits,
      } as Plan);
      setQuery("");
      setIsOpen(false);
      // Keep focus on input for rapid adding
      inputRef.current?.focus();
    },
    [addCoursePlan, academicYear, sem, yearOfStudy]
  );

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!isOpen || filtered.length === 0) return;

    if (e.key === "ArrowDown") {
      e.preventDefault();
      setHighlightIdx((prev) => Math.min(prev + 1, filtered.length - 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setHighlightIdx((prev) => Math.max(prev - 1, 0));
    } else if (e.key === "Enter") {
      e.preventDefault();
      handleSelect(filtered[highlightIdx]);
    } else if (e.key === "Escape") {
      setIsOpen(false);
    }
  };

  // Scroll highlighted item into view
  useEffect(() => {
    if (!dropdownRef.current) return;
    const item = dropdownRef.current.children[highlightIdx] as HTMLElement;
    if (item) {
      item.scrollIntoView({ block: "nearest" });
    }
  }, [highlightIdx]);

  return (
    <div className="relative">
      <div className="relative">
        <Plus className="absolute left-2.5 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-gray-400" />
        <Input
          ref={inputRef}
          placeholder="Type course code or name..."
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            setIsOpen(e.target.value.trim().length > 0);
          }}
          onFocus={() => {
            if (query.trim().length > 0) setIsOpen(true);
          }}
          onKeyDown={handleKeyDown}
          className="pl-8 text-xs h-9 bg-gray-50 border-gray-200 text-gray-600 placeholder:text-gray-400"
        />
      </div>

      {isOpen && filtered.length > 0 && (
        <div
          ref={dropdownRef}
          className="absolute z-50 top-full mt-1 left-0 right-0 bg-white border border-gray-200 rounded-md shadow-sm max-h-[280px] overflow-y-auto"
        >
          {filtered.map((course, idx) => (
            <button
              key={course.code}
              type="button"
              className={`w-full text-left px-3 py-2 text-xs flex items-center gap-2 transition-colors cursor-pointer border-b border-gray-100 last:border-0 ${
                idx === highlightIdx
                  ? "bg-gray-100"
                  : "hover:bg-gray-50"
              }`}
              onMouseEnter={() => setHighlightIdx(idx)}
              onMouseDown={(e) => {
                e.preventDefault(); // Prevent input blur
                handleSelect(course);
              }}
            >
              <span className="font-mono font-medium shrink-0 text-gray-700">
                {course.code}
              </span>
              <span className="truncate flex-1 text-gray-600">
                {course.name_en || course.name_th}
              </span>
              {course.categoryId && (
                <Badge
                  variant="secondary"
                  className="text-[10px] px-1.5 py-0 shrink-0 bg-gray-100 text-gray-600"
                >
                  {course.categoryId}
                </Badge>
              )}
              <span className="text-[10px] text-gray-500 shrink-0">
                {course.credits} cr
              </span>
            </button>
          ))}
        </div>
      )}

      {isOpen && query.trim().length > 0 && filtered.length === 0 && (
        <div className="absolute z-50 top-full mt-1 left-0 right-0 bg-white border border-gray-200 rounded-md shadow-sm p-3 text-xs text-gray-500 text-center">
          No courses found for &quot;{query}&quot;
        </div>
      )}
    </div>
  );
}
