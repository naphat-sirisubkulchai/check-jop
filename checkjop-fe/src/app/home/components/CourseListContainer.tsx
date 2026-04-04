"use client";
import { useEffect, useMemo, useState } from "react";
import { useAppStore } from "@/store/appStore";
import { courseApi } from "@/api/courseApi";
import { toast } from "sonner";
import { uniqueCoursesByCode } from "@/utils";
import { ChevronsLeft, ChevronsRight, Search } from "lucide-react";
import { Input } from "@/components/ui/input";
import CourseLibraryCard from "./CourseLibraryCard";
import { Button } from "@/components/ui/button";

export default function CourseListContainer() {
  const { courses, studyPlan, categories } = useAppStore();

  const electiveCodes = useMemo(() => {
    const s = new Set<string>();
    for (const cat of categories) {
      const name = (cat as any).nameTH ?? (cat as any).name_th ?? "";
      if (name.includes("เลือก") && !name.includes("เสรี")) {
        for (const c of (cat as any).courses ?? []) {
          if (c.code) s.add(c.code);
        }
      }
    }
    return s;
  }, [categories]);

  const [searchQuery, setSearchQuery] = useState("");
  const [isCollapsed, setIsCollapsed] = useState(false);
  useEffect(() => {
    loadCurricula();
  }, []);

  // Filter courses based on search query
  const filteredCourses = uniqueCoursesByCode(courses).filter((course: any) => {
    const query = searchQuery.toLowerCase();
    return (
      course.code?.toLowerCase().includes(query) ||
      course.name_en?.toLowerCase().includes(query) ||
      course.name_th?.toLowerCase().includes(query)
    );
  });

  // Render course list or empty state
  const renderCourseList = useMemo(() => {
    if (filteredCourses.length === 0) {
      return (
        <div className="flex flex-col items-center justify-center px-4 py-16">
          <div className="mb-4 rounded-full bg-gray-100 p-4">
            <Search className="h-8 w-8 text-gray-400" />
          </div>
          <p className="mb-1 text-sm font-semibold text-gray-700">
            {searchQuery ? "No courses found" : "No courses available"}
          </p>
          {searchQuery ? (
            <p className="text-center text-xs text-gray-500">
              Try searching with a different course code or name
            </p>
          ) : (
            <p className="text-center text-xs text-gray-500">
              Select a curriculum to view available courses
            </p>
          )}
        </div>
      );
    }
    const sorted = filteredCourses
      .map((course: any) => ({
        course,
        isInPlan: studyPlan.some((p) => p.course_code === course.code && p.grade !== "F"),
      }))
      .sort((a, b) => {
        if (a.isInPlan !== b.isInPlan) return a.isInPlan ? 1 : -1;
        return a.course.code.localeCompare(b.course.code);
      });

    return (
      <div className="flex flex-col space-y-2">
        {sorted.map(({ course, isInPlan }) => (
          <CourseLibraryCard
            key={course.code}
            course={course}
            isInPlan={isInPlan}
            isElective={electiveCodes.has(course.code)}
          />
        ))}
      </div>
    );
  }, [filteredCourses, studyPlan]);

  return (
    <div
      className={`flex flex-col p-4 shadow-sm z-10 transition-all duration-300 ${isCollapsed ? "w-16" : "w-90"}`}
    >
      {/* Search and Course List */}
      <div className="flex flex-col flex-1 overflow-hidden">
        {/* โชว์ side bar */}
        {!isCollapsed && (
          <div className="flex flex-col flex-1 space-y-4 overflow-hidden">
            {/* Header Section with Search */}
            <div className="flex items-center justify-between">
              <h3 className="text-xl font-bold">Course Library</h3>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setIsCollapsed(!isCollapsed)}
                className="h-8 w-8 shrink-0 hover:bg-background"
              >
                <ChevronsLeft className="h-5 w-5 text-gray-400 hover:bg-background" />
              </Button>
            </div>

            {/* Search Bar */}
            <div className="relative">
              <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
              <Input
                placeholder="Search courses..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="h-10 border-gray-200 bg-surface-soft pl-9 pr-4 focus:bg-white"
              />
            </div>

            {/* Course List */}
            <div className="flex-1 overflow-x-hidden overflow-y-auto">
              {renderCourseList}
            </div>
          </div>
        )}

        {/* ซ่อน side bar */}
        {isCollapsed && (
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setIsCollapsed(!isCollapsed)}
            className="hover:bg-background"
          >
            <ChevronsRight className="h-5 w-5 text-gray-400" />
          </Button>
        )}
      </div>
    </div>
  );
}

// โหลดหลักสูตรทั้งหมด
async function loadCurricula() {
  try {
    const data = await courseApi.getAllCurriculaWithout();
    useAppStore.getState().setCurriculums(data);
  } catch (error) {
    const message =
      error instanceof Error ? error.message : "Unknown error occurred";
    toast.error("Failed to load curricula", { description: message });
    console.error("Failed to load curricula:", error);
  }
}
