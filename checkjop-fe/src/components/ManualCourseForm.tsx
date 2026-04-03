"use client";

import { useState, useMemo, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Plus, X } from "lucide-react";
import { Plan } from "@/types";
import { useAppStore } from "@/store/appStore";
import { toast } from "sonner";
import { GRADE_OPTIONS } from "@/lib/constants";
interface ManualCourseFormProps {
  onClose?: () => void;
  semester: number;
  yearOfStudy: number;
  academicYear: number;
  editPlan?: {
    course_code: string;
    course_name?: string;
    credits: number;
    category_name?: string;
    grade?: string;
  };
}

export default function ManualCourseForm({
  onClose,
  semester,
  yearOfStudy,
  academicYear,
  editPlan,
}: ManualCourseFormProps) {
  const addCoursePlan = useAppStore((state) => state.addCoursePlan);
  const editCoursePlan = useAppStore((state) => state.editCoursePlan);
  const categories = useAppStore((state) => state.categories);
  const isEditMode = !!editPlan;

  const [formData, setFormData] = useState({
    course_code: editPlan?.course_code || "",
    course_name: editPlan?.course_name || "",
    credits: editPlan?.credits ?? 3,
    category_name: editPlan?.category_name || "",
    grade: editPlan?.grade || "",
    isCF: false,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  const validateForm = useCallback((): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.course_code.trim()) {
      newErrors.course_code = "Course code is required";
    }

    if (formData.credits <= 0 || formData.credits > 20) {
      newErrors.credits = "Credits must be between 1-20";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  }, [formData]);

  const handleSubmit = useCallback(
    (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault();

      if (!validateForm()) {
        return;
      }

      setIsSubmitting(true);

      try {
        const courseCode = formData.course_code.trim().toUpperCase();

        if (isEditMode && editPlan) {
          editCoursePlan(editPlan.course_code, {
            course_code: courseCode,
            course_name: formData.course_name.trim() || undefined,
            credits: formData.credits,
            category_name: formData.category_name || undefined,
            grade: formData.grade || undefined,
          }, yearOfStudy, semester);
        } else {
          const newPlan: Plan = {
            course_code: courseCode,
            course_name: formData.course_name.trim() || undefined,
            academicYear: academicYear,
            yearOfStudy: yearOfStudy,
            semester: semester,
            credits: formData.credits,
            category_name: formData.category_name || undefined,
            grade: formData.grade || undefined,
          };
          addCoursePlan(newPlan);

          setFormData({ course_code: "", course_name: "", credits: 3, category_name: "", grade: "", isCF: false });
          setErrors({});
        }

        if (onClose) onClose();
      } catch (error) {
        console.error("Error adding course:", error);
        toast.error("Failed to add course", {
          description:
            error instanceof Error ? error.message : "Please try again",
        });
        setErrors({ submit: "Failed to add course. Please try again." });
      } finally {
        setIsSubmitting(false);
      }
    },
    [
      formData,
      validateForm,
      addCoursePlan,
      academicYear,
      yearOfStudy,
      semester,
      onClose,
    ],
  );

  const handleInputChange = useCallback(
    (field: string, value: string | number | boolean) => {
      setFormData((prev) => ({ ...prev, [field]: value }));
      // Clear error for this field when user starts typing
      setErrors((prev) => {
        if (prev[field]) {
          const newErrors = { ...prev };
          delete newErrors[field];
          return newErrors;
        }
        return prev;
      });
    },
    [],
  );

  // Memoize category options to prevent re-rendering
  const categoryOptions = useMemo(
    () =>
      categories.map((category) => (
        <SelectItem key={category.id} value={category.name_th}>
          {category.name_th}
        </SelectItem>
      )),
    [categories],
  );

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      {/* Context Info Banner */}
      <div className="p-4 bg-chula-soft border-2 border-chula/20 rounded-xl">
        <p className="text-sm text-gray-900 font-medium">
          <span className="text-chula font-semibold">Adding to:</span> Year {yearOfStudy}, Semester {semester}
          {semester === 3 ? " (Summer)" : ""} <span className="text-gray-500">(พ.ศ. {academicYear})</span>
        </p>
      </div>

      {/* Course Code */}
      <div className="space-y-2">
        <Label htmlFor="course_code" className="text-base font-medium">
          Course Code <span className="text-danger">*</span>
        </Label>
        <Input
          id="course_code"
          placeholder="e.g., 2301102"
          value={formData.course_code}
          onChange={(e) => handleInputChange("course_code", e.target.value)}
          className={`h-11 text-base border-2 bg-surface-soft ${errors.course_code ? "border-danger focus:border-danger" : "focus:border-chula"}`}
          autoFocus
        />
        {errors.course_code && (
          <p className="text-sm text-danger flex items-center gap-1">
            {errors.course_code}
          </p>
        )}
      </div>

      {/* Course Name */}
      <div className="space-y-2">
        <Label htmlFor="course_name" className="text-base font-medium">
          Course Name <span className="text-gray-400 text-sm font-normal">(Optional)</span>
        </Label>
        <Input
          id="course_name"
          placeholder="e.g., Computer Programming"
          value={formData.course_name}
          onChange={(e) => handleInputChange("course_name", e.target.value)}
          className="h-11 text-base border-2 focus:border-chula bg-surface-soft"
        />
      </div>

      <div className="grid grid-cols-2 gap-4">
        {/* Credits */}
        <div className="space-y-2">
          <Label htmlFor="credits" className="text-base font-medium">
            Credits <span className="text-danger">*</span>
          </Label>
          <Input
            id="credits"
            type="number"
            step="0.5"
            placeholder="3"
            value={formData.credits}
            onChange={(e) =>
              handleInputChange("credits", parseFloat(e.target.value) || 0)
            }
            className={`h-11 text-base border-2 bg-surface-soft ${errors.credits ? "border-danger focus:border-danger" : "focus:border-chula"}`}
          />
          {errors.credits && (
            <p className="text-sm text-danger">{errors.credits}</p>
          )}
        </div>

        {/* Category Name */}
        <div className="space-y-2">
          <Label htmlFor="category_name" className="text-base font-medium">
            Category <span className="text-gray-400 text-sm font-normal">(Optional)</span>
          </Label>
          <Select
            value={formData.category_name || "none"}
            onValueChange={(value) =>
              handleInputChange("category_name", value === "none" ? "" : value)
            }
          >
            <SelectTrigger className="w-full min-h-11 text-base border-2 focus:border-chula mb-0 bg-surface-soft">
              <SelectValue placeholder="Select category" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="none" className="text-gray-500">None</SelectItem>
              {categoryOptions}
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Grade */}
      <div className="space-y-2">
        <Label htmlFor="grade" className="text-base font-medium">
          Grade <span className="text-gray-400 text-sm font-normal">(Optional)</span>
        </Label>
        <Select
          value={formData.grade || "none"}
          onValueChange={(value) =>
            handleInputChange("grade", value === "none" ? "" : value)
          }
        >
          <SelectTrigger className="w-full h-11 text-base border-2 focus:border-chula bg-surface-soft">
            <SelectValue placeholder="Select grade" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="none" className="text-gray-500">None</SelectItem>
            {GRADE_OPTIONS.map((option) => (
              <SelectItem key={option.value} value={option.value} className="text-base">
                <span>{option.label}</span>
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {/* Is CF */}
      {/* <div className="flex items-center justify-between p-4 bg-surface-soft rounded-lg border-2 border-gray-200">
        <div className="flex-1">
          <Label htmlFor="isCF" className="text-base font-medium cursor-pointer">
            TODO: CF ย่อจากไรนะ
          </Label>
          <p className="text-sm text-gray-500 mt-1">description</p>
        </div>
        <Switch
          id="isCF"
          checked={formData.isCF}
          onCheckedChange={(checked) => handleInputChange("isCF", checked)}
          className="data-[state=checked]:bg-chula"
        />
      </div> */}

      {/* Submit Error */}
      {errors.submit && (
        <div className="p-4 rounded-lg bg-danger-soft border-2 border-danger text-danger text-sm font-medium">
          {errors.submit}
        </div>
      )}

      {/* Action Buttons */}
      <div className="flex gap-3 pt-4 border-t-2">
        <Button
          type="submit"
          disabled={isSubmitting}
          className="flex-1 h-12 bg-chula hover:bg-chula-hover text-white font-medium shadow-sm text-base disabled:opacity-50"
        >
          {!isEditMode && <Plus className="h-5 w-5 mr-2" />}
          {isSubmitting ? (isEditMode ? "Saving..." : "Adding...") : (isEditMode ? "Save Changes" : "Add Course")}
        </Button>
        {onClose && (
          <Button
            type="button"
            variant="outline"
            onClick={onClose}
            className="h-12 border-2 px-6 text-base font-medium"
          >
            <X className="h-5 w-5 mr-2" />
            Cancel
          </Button>
        )}
      </div>
    </form>
  );
}
