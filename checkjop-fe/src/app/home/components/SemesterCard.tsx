"use client";
import { Plan } from "@/types";
import CourseCard from "./CourseCard";
import { QuickAdd } from "./QuickAdd";
import { useAppStore } from "@/store/appStore";
import ManualCourseDialog from "@/components/ManualCourseDialog";
import { useMemo, useState } from "react";
import { Button } from "@/components/ui/button";
import { PenLine, BookOpen } from "lucide-react";
import { toast } from "sonner";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
interface SemesterCardProps {
  sem: number;
  yearOfStudy: number;
  academicYear: number;
  className?: string; // Optional className prop for custom styling
}

export default function SemesterCard({
  sem,
  yearOfStudy,
  academicYear,
  className,
}: SemesterCardProps) {
  const { studyPlan, addCoursePlan, getCourseByCode, editCoursePlan } =
    useAppStore();
  const [isManualFormOpen, setIsManualFormOpen] = useState(false);

  const onFocus = () => {
    const el = document.getElementById(`Y:${yearOfStudy}S:${sem}`);
    el?.classList.remove("border-gray-200", "bg-surface-soft");
    el?.classList.add(
      "border-chula",
      "bg-chula-soft",
      "shadow-md",
      "scale-[1.01]",
    );
  };

  const onBlur = () => {
    const el = document.getElementById(`Y:${yearOfStudy}S:${sem}`);
    el?.classList.remove(
      "border-chula",
      "bg-chula-soft",
      "shadow-md",
      "scale-[1.01]",
    );
    el?.classList.add("border-gray-200", "bg-surface-soft");
  };

  // Handle adding course from library
  const handleAddFromLibrary = (courseCode: string) => {
    // Check if course already exists in study plan (allow re-adding if previous attempt was F)
    const isDuplicate = studyPlan.find((p) => p.course_code === courseCode);

    if (isDuplicate && isDuplicate.grade !== "F") {
      toast.error("Course already added", {
        description: `${courseCode} is already in Year ${isDuplicate.yearOfStudy}, Semester ${isDuplicate.semester}`,
      });
      return;
    }

    const course = getCourseByCode(courseCode);
    addCoursePlan({
      course_code: courseCode,
      yearOfStudy: yearOfStudy,
      semester: sem,
      academicYear: academicYear,
      credits: course?.credits || 0,
    } as Plan);
  };

  // Handle moving course from another semester
  const handleMoveCourse = (courseCode: string) => {
    // Find the course in the study plan
    const existingCourse = studyPlan.find((p) => p.course_code === courseCode);
    if (!existingCourse) return;

    // Check if trying to drop in the same semester
    if (
      existingCourse.yearOfStudy === yearOfStudy &&
      existingCourse.semester === sem
    ) {
      toast.info("Same semester", {
        description: `${courseCode} is already in this semester`,
      });
      return;
    }

    // Update the course's semester and year
    editCoursePlan(courseCode, {
      yearOfStudy: yearOfStudy,
      semester: sem,
      academicYear: academicYear,
    });

    // toast.success("Course moved", {
    //   description: `${courseCode} moved to Year ${yearOfStudy}, Semester ${sem}`,
    // });
  };

  // Handle drop event
  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    onBlur();

    const payload = e.dataTransfer?.getData("text/plain");
    if (!payload) return;

    // Handle drag from course library
    if (payload.startsWith("LIB:")) {
      const courseCode = payload.replace("LIB:", "");
      handleAddFromLibrary(courseCode);
      return;
    }

    // Handle drag from another semester (move course)
    if (payload.startsWith("PLAN:")) {
      const courseCode = payload.replace("PLAN:", "");
      handleMoveCourse(courseCode);
      return;
    }
  };

  // Calculate total credits for the semester
  function getTotalCreditsForSemester(sem: number, year: number) {
    let total = 0;
    studyPlan
      .filter((item) => item.semester === sem && item.yearOfStudy === year)
      .forEach((item) => {
        const course = getCourseByCode(item.course_code);
        if (course) {
          total += course.credits;
        }
      });
    return total;
  }

  if (!studyPlan) console.log("don't have any course plan");

  const totalCredits = getTotalCreditsForSemester(sem, yearOfStudy);
  const coursesInSemester = useMemo(
    () =>
      studyPlan.filter(
        (item) => item.semester === sem && item.yearOfStudy === yearOfStudy
      ),
    [studyPlan, sem, yearOfStudy]
  );

  return (
    <div
      key={`Y:${yearOfStudy}S:${sem}`}
      id={`Y:${yearOfStudy}S:${sem}`}
      className={`semester-box rounded-xl p-4 bg-surface-soft shadow-sm hover:shadow-md transition-all duration-200 flex flex-col ${className}`}
      onDragOver={(e) => {
        e.preventDefault();
        onFocus();
      }}
      onDragLeave={onBlur}
      onDrop={handleDrop}
    >
      {/* Header */}
      <div className="flex justify-between items-center mb-3">
        <h3 className="font-bold text-lg text-gray-900">
          Semester {sem}{" "}
          {sem === 3 ? <span className="text-chula-active">(Summer)</span> : ""}
        </h3>
        <Badge className="rounded-full bg-chula-soft text-chula-active">
          {totalCredits} cr.
        </Badge>
      </div>

      {/* Courses List */}
      <div className="semester-courses space-y-2 overflow-y-auto flex-col flex flex-1 min-h-[120px]">
        {coursesInSemester.length > 0 ? (
          coursesInSemester
            .filter((item) => item.course_code) // Ensure item.course exists
            .map((item) => (
              <CourseCard
                key={`${item.course_code}-${yearOfStudy}-${sem}-${academicYear}`}
                courseCode={item.course_code}
                yearOfStudy={yearOfStudy}
                semester={sem}
                academicYear={academicYear}
              />
            ))
        ) : (
          <div className="flex justify-center items-center flex-1 py-4">
            <div className="empty-state text-center space-y-2 px-4 py-4 rounded-lg border border-dashed border-gray-200 bg-gray-50/50">
              <div className="inline-flex items-center justify-center w-8 h-8 rounded-full bg-gray-100 mb-1">
                <BookOpen className="h-4 w-4 text-gray-400" />
              </div>
              <div>
                <div className="text-gray-600 text-sm font-medium mb-0.5">
                  No courses yet
                </div>
                <p className="text-gray-400 text-xs">
                  Drag courses or use Quick Add
                </p>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Separator */}
      <Separator className="my-4" />

      {/* Quick Add at bottom */}
      <div className="flex gap-2 border-gray-200 mt-auto">
        <div className="flex-1">
          <QuickAdd
            sem={sem}
            yearOfStudy={yearOfStudy}
            academicYear={academicYear}
          />
        </div>
        <Button
          onClick={() => setIsManualFormOpen(true)}
          title="Add course manually"
          className="px-4 bg-chula-soft hover:bg-chula-soft/80 rounded-md shadow-sm p-3 text-xs text-chula-active font-medium border border-chula-active/30"
        >
          <PenLine className="w-4 h-4 mr-2" />
          Manual
        </Button>
      </div>

      {/* Manual Course Form Dialog */}
      <ManualCourseDialog
        open={isManualFormOpen}
        onOpenChange={setIsManualFormOpen}
        semester={sem}
        yearOfStudy={yearOfStudy}
        academicYear={academicYear}
      />
    </div>
  );
}
