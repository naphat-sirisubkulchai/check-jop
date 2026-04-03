"use client";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import ManualCourseForm from "@/components/ManualCourseForm";

interface ManualCourseDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title?: string;
  description?: string;
  // Pre-filled values from semester card
  semester: number;
  yearOfStudy: number;
  academicYear: number;
}

/**
 * Reusable dialog component for manually adding courses
 *
 * @example
 * ```tsx
 * const [isOpen, setIsOpen] = useState(false);
 *
 * <ManualCourseDialog
 *   open={isOpen}
 *   onOpenChange={setIsOpen}
 *   semester={1}
 *   yearOfStudy={1}
 *   academicYear={2567}
 * />
 * ```
 */
export default function ManualCourseDialog({
  open,
  onOpenChange,
  title = "Add Course Manually",
  description = "Enter course details to add it to your study plan",
  semester,
  yearOfStudy,
  academicYear,
}: ManualCourseDialogProps) {
  const handleClose = () => {
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>
        <ManualCourseForm
          onClose={handleClose}
          semester={semester}
          yearOfStudy={yearOfStudy}
          academicYear={academicYear}
        />
      </DialogContent>
    </Dialog>
  );
}
