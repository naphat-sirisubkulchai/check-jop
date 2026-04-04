"use client";

import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Plan } from "@/types";
import { BookOpen, GraduationCap } from "lucide-react";
import { useAppStore } from "@/store/appStore";

interface StudyPlanRecordCardProps {
  studyPlan: Plan[];
}

export function StudyPlanRecordCard({ studyPlan }: StudyPlanRecordCardProps) {
  const { getCourseByCode, categories } = useAppStore();

  // Build set of course codes in elective categories
  const electiveCodes = new Set<string>();
  for (const cat of categories) {
    const name = (cat as any).nameTH ?? (cat as any).name_th ?? "";
    if (name.includes("เลือก") && !name.includes("เสรี")) {
      for (const c of (cat as any).courses ?? []) {
        if (c.code) electiveCodes.add(c.code);
      }
    }
  }

  // Group courses by year and semester
  const groupedPlan = studyPlan.reduce((acc, plan) => {
    const key = `${plan.yearOfStudy}`;
    if (!acc[key]) {
      acc[key] = { 1: [], 2: [], 3: [] };
    }
    acc[key][plan.semester].push(plan);
    console.log(acc);
    
    return acc;
  }, {} as Record<string, Record<number, Plan[]>>);

  const years = Object.keys(groupedPlan)
    .map(Number)
    .sort((a, b) => a - b);

  // Calculate total credits per semester
  const getSemesterCredits = (year: number, semester: number) => {
    const courses = groupedPlan[year]?.[semester] || [];
    return courses.reduce((sum, plan) => sum + plan.credits, 0);
  };

  if (studyPlan.length === 0) {
    return (
      <Card className="p-8 text-center bg-gray-50 border-2 border-dashed">
        <BookOpen className="h-12 w-12 text-gray-300 mx-auto mb-3" />
        <p className="text-gray-500 font-medium">No study plan records</p>
      </Card>
    );
  }

  return (
    <Card className="p-6 shadow-md border-gray-200">
      <div className="flex items-center gap-3 mb-2">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-gradient-to-br from-chula-active to-pink-500 shadow-sm">
            <GraduationCap className="h-5 w-5 text-white" />
          </div>
          <div>
            <h3 className="text-2xl font-bold bg-gradient-to-r from-chula-active to-pink-500 bg-clip-text text-transparent">
              Study Plan Record
            </h3>
            <p className="text-sm text-gray-600 mt-0.5">
              Review the courses used in this calculation
            </p>
          </div>
        </div>

      <div className="space-y-6">
        {years.map((year) => {
          const academicYear = groupedPlan[year][1]?.[0]?.academicYear || 2566 + year - 1;

          return (
            <div key={year} className="space-y-3">
              {/* Year Header */}
              <div className="flex items-center gap-2 pb-2 border-b border-gray-200">
                <h4 className="text-lg font-bold text-gray-900">
                  Year {year}
                </h4>
                <span className="text-sm font-medium text-gray-500">
                  ({academicYear})
                </span>
              </div>

              {/* Semesters Grid */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {[1, 2, 3].map((semester) => {
                  const courses = groupedPlan[year][semester] || [];
                  const totalCredits = getSemesterCredits(year, semester);

                  if (courses.length === 0) return null;

                  return (
                    <div
                      key={semester}
                      className="bg-gray-50 rounded-lg p-4 border border-gray-200"
                    >
                      {/* Semester Header */}
                      <div className="flex items-center justify-between mb-3">
                        <h5 className="font-semibold text-sm text-gray-900">
                          Semester {semester}
                          {semester === 3 && (
                            <span className="text-chula-active ml-1">(Summer)</span>
                          )}
                        </h5>
                        <Badge className="bg-chula-soft text-chula-active text-xs">
                          {totalCredits} cr.
                        </Badge>
                      </div>

                      {/* Course List */}
                      <div className="space-y-2">
                        {courses.map((plan, idx) => {
                          const course = getCourseByCode(plan.course_code);
                          const courseName = course?.name_en || plan.course_name || "Manual Course";
                          // const courseCategory = course?.category_name || null;

                          const isElective = electiveCodes.has(plan.course_code);
                          return (
                            <div
                              key={`${plan.course_code}-${idx}`}
                              className={`bg-white rounded-md p-2.5 border shadow-sm ${isElective ? "border-l-4 border-l-sci-normal border-gray-200" : "border border-gray-200"}`}
                            >
                              <div className="flex items-start justify-between gap-2">
                                <div className="flex-1 min-w-0">
                                  <div className="font-mono text-xs font-semibold text-chula-active mb-0.5">
                                    {plan.course_code}
                                  </div>
                                  <div className="text-xs text-gray-700 line-clamp-2">
                                    {courseName}
                                  </div>
                                </div>
                                <div className="flex flex-row-reverse items-end gap-1 flex-shrink-0">
                                  <Badge
                                    variant="secondary"
                                    className="bg-gray-100 text-gray-700 text-xs px-1.5 py-0.5"
                                  >
                                    {plan.credits} cr.
                                  </Badge>
                                  {plan.grade && (
                                    <Badge
                                      variant="secondary"
                                      className="bg-green-100 text-green-700 text-xs px-1.5 py-0.5"
                                    >
                                      {plan.grade}
                                    </Badge>
                                  )}
                                </div>
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          );
        })}
      </div>
    </Card>
  );
}
