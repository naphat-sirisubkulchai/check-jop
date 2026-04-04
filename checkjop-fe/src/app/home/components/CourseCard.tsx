import { Button } from "@/components/ui/button";
import { useAppStore } from "@/store/appStore";
import {
  Trash2,
  GraduationCap,
  PenLine,
  FolderTree,
  GripVertical,
} from "lucide-react";
import { MouseEvent, useState } from "react";
import GradeDialog from "@/components/GradeDialog";
import { Badge } from "@/components/ui/badge";
import { Toggle } from "@/components/ui/toggle";
import { gradeService } from "@/api/gradApi";
import { toast } from "sonner";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import ManualCourseForm from "@/components/ManualCourseForm";

interface CourseCardProps {
  courseCode: string;
  yearOfStudy: number;
  semester: number;
  academicYear: number;
}

export default function CourseCard({
  courseCode,
  yearOfStudy,
  semester,
  academicYear,
}: CourseCardProps) {
  const {
    getCourseByCode,
    removeCoursePlan,
    editCoursePlan,
    studyPlan,
    exemptions,
    addExemption,
    removeExemption,
    selectedCurriculum,
    yearMapping,
    categories,
  } = useAppStore();

  // Build elective code set
  const electiveCodes = new Set<string>();
  for (const cat of categories) {
    const name = (cat as any).nameTH ?? (cat as any).name_th ?? "";
    if (name.includes("เลือก") && !name.includes("เสรี")) {
      for (const c of (cat as any).courses ?? []) {
        if (c.code) electiveCodes.add(c.code);
      }
    }
  }
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isCFLoading, setIsCFLoading] = useState(false);

  // Find the course plan from study plan based on all identifiers
  const coursePlan = studyPlan.find(
    (p) =>
      p.course_code === courseCode &&
      p.yearOfStudy === yearOfStudy &&
      p.semester === semester &&
      p.academicYear === academicYear,
  );

  // Try to get course details from course list (for additional info like name_en)
  const courseFromList = getCourseByCode(courseCode);

  if (!coursePlan) {
    console.log(
      `Course plan not found: ${courseCode} (Year ${yearOfStudy}, Semester ${semester})`,
    );
    return null;
  }

  // Use course name from plan (for manual courses) or from course list
  const courseName =
    coursePlan.course_name ||
    courseFromList?.name_en ||
    courseFromList?.course_name ||
    "Unknown Course";
  const courseCredits = coursePlan.credits || courseFromList?.credits || 0;
  const existingGrade = coursePlan.grade;
  const categoryName = coursePlan.category_name?.trim() || undefined;

  // Check if this is a manually added course
  const isManualCourse = !courseFromList || !!coursePlan.course_name;

  // Check if this course is in exemptions
  const isExempted = exemptions.includes(courseCode);

  function handleOpenGradeDialog(event: MouseEvent<HTMLButtonElement>): void {
    event.stopPropagation();
    setIsDialogOpen(true);
  }

  async function handleToggleExemption(pressed: boolean): Promise<void> {
    if (pressed) {
      // Check if we have required data
      if (!selectedCurriculum?.id || !yearMapping) {
        toast.error("Cannot add CF", {
          description: "Curriculum or year mapping is not available",
        });
        return;
      }

      // Get academic year from yearMapping
      const academicYear = yearMapping[yearOfStudy];
      if (!academicYear) {
        toast.error("Cannot add CF", {
          description: `No academic year found for Year ${yearOfStudy}`,
        });
        return;
      }

      // Check CF option via API
      setIsCFLoading(true);
      try {
        const result = await gradeService.checkCFOption(
          courseCode,
          selectedCurriculum.id,
          academicYear
        );

        if (result.has_cf_option) {
          // Course allows CF - add to exemptions
          addExemption(courseCode);
          // toast.success("CF Added", {
          //   description: result.message,
          // });
        } else {
          // Course does NOT allow CF
          toast.error("CF Not Allowed", {
            description: result.message,
          });
        }
      } catch (error) {
        console.error("Error checking CF option:", error);
        toast.error("Failed to check CF option", {
          description: error instanceof Error ? error.message : "Unknown error occurred",
        });
      } finally {
        setIsCFLoading(false);
      }
    } else {
      // Remove CF
      removeExemption(courseCode);
      // toast.info("CF Removed", {
      //   description: `Removed CF for ${courseCode}`,
      // });
    }
  }

  function handleSaveGrade(grade: string): void {
    editCoursePlan(courseCode, { grade: grade || undefined }, yearOfStudy, semester);
  }

  function handleClearGrade(): void {
    editCoursePlan(courseCode, { grade: undefined }, yearOfStudy, semester);
  }

  return (
    <>
      <div
        className={`group relative flex gap-2 rounded-lg border bg-white p-2.5 shadow-sm transition-all hover:shadow-md cursor-grab active:cursor-grabbing overflow-hidden ${
          isManualCourse
            ? "border-l-4 border-l-chula-active bg-gradient-to-r from-chula-soft/30 to-white"
            : electiveCodes.has(courseCode)
              ? "border-l-4 border-l-sci-normal bg-gradient-to-r from-sci-soft/30 to-white hover:border-sci-hover/30"
              : "border-l-4 border-l-chula-soft bg-gradient-to-r from-chula-soft/30 to-white hover:border-chula-active/30"
        }`}
        draggable={true}
        onDragStart={(e) => {
          e.dataTransfer.setData("text/plain", `PLAN:${courseCode}`);
          e.currentTarget.classList.add("opacity-50");
        }}
        onDragEnd={(e) => {
          e.currentTarget.classList.remove("opacity-50");
        }}
      >
        {/* Manual Ribbon */}
        {isManualCourse && (
          <div className="absolute top-2 right-[-18px] rotate-45 bg-chula-active text-white text-[9px] font-bold px-6 py-0.5 shadow-sm pointer-events-none">
            Manual
          </div>
        )}

        {/* Drag Handle */}
        <div className="flex shrink-0 items-center">
          <GripVertical className="h-4 w-4 text-gray-300 transition-colors group-hover:text-chula-active" />
        </div>

        {/* Content */}
        <div className="min-w-0 flex-1">
          {/* Course Code & Name */}
          <div className="flex items-baseline gap-1.5">
            <span className="text-sm font-bold text-gray-900 shrink-0">{courseCode}</span>
            {courseName !== courseCode && (
              <span className="min-w-0 truncate text-sm text-gray-700">{courseName}</span>
            )}
          </div>

          {/* Metadata Row */}
          <div className="flex items-center gap-1.5 mt-0.5 flex-wrap">
            <span className="text-xs text-gray-500 font-medium">{courseCredits} cr</span>

            {existingGrade && (
              <Badge className="min-w-8 rounded-full bg-sci-soft text-xs font-bold text-sci-active">
                {existingGrade}
              </Badge>
            )}

            {categoryName && (
              <span className="inline-flex items-center gap-1 text-xs text-gray-500 font-medium bg-gray-100 rounded px-1 max-w-[140px]">
                <FolderTree className="h-3 w-3 shrink-0" />
                <span className="truncate">{categoryName}</span>
              </span>
            )}

            {isExempted && (
              <Badge className="rounded-full bg-green-100 text-xs font-semibold text-green-700 border-green-200">
                CF
              </Badge>
            )}
          </div>
        </div>

        {/* Action Buttons - Visible on Hover */}
        <div className="flex shrink-0 items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100">
          {/* Edit Manual Course */}
          {isManualCourse && (
            <Button
              variant="ghost"
              size="icon"
              onClick={(e) => { e.stopPropagation(); setIsEditDialogOpen(true); }}
              className="h-7 w-7 text-gray-500 hover:bg-gray-100"
              title="Edit course"
            >
              <PenLine className="h-4 w-4" />
            </Button>
          )}
          {/* CF Toggle */}
          <Toggle
            pressed={isExempted}
            onPressedChange={handleToggleExemption}
            size="sm"
            variant="outline"
            title="Credit for exemption (CF)"
            onClick={(e) => e.stopPropagation()}
            disabled={isCFLoading}
            className={`gap-1.5 rounded-full border-0 shadow-none text-green-700 hover:bg-green-50 data-[state=on]:bg-green-100 data-[state=on]:text-green-700 ${isCFLoading ? "opacity-50 cursor-not-allowed" : ""}`}
          >
            {isCFLoading ? (
              <div className="h-3 w-3 border-2 border-green-700 border-t-transparent rounded-full animate-spin"></div>
            ) : (
              <span className="text-xs font-medium">CF</span>
            )}
          </Toggle>
          <Button
            variant="ghost"
            size="icon"
            onClick={handleOpenGradeDialog}
            className="h-7 w-7 text-chula-active hover:bg-chula-soft hover:text-chula-active"
            title={existingGrade ? "Change grade" : "Add grade"}
          >
            <GraduationCap className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={(e) => {
              e.stopPropagation();
              removeCoursePlan(courseCode);
            }}
            className="h-7 w-7 text-red-600 hover:bg-red-50 hover:text-red-700"
            title="Remove course"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Edit Manual Course Dialog */}
      {isManualCourse && (
        <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
          <DialogContent className="max-w-lg">
            <DialogHeader>
              <DialogTitle>Edit Course</DialogTitle>
            </DialogHeader>
            <ManualCourseForm
              onClose={() => setIsEditDialogOpen(false)}
              semester={semester}
              yearOfStudy={yearOfStudy}
              academicYear={academicYear}
              editPlan={{
                course_code: courseCode,
                course_name: coursePlan.course_name,
                credits: courseCredits,
                category_name: categoryName,
                grade: existingGrade,
              }}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Grade Dialog */}
      <GradeDialog
        open={isDialogOpen}
        onOpenChange={setIsDialogOpen}
        courseCode={courseCode}
        courseName={courseName}
        currentGrade={existingGrade}
        onSave={handleSaveGrade}
        onClear={handleClearGrade}
      />
    </>
  );
}
