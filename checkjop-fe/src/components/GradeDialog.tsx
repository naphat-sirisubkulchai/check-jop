"use client";

import { useState, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { GraduationCap } from "lucide-react";
import { GRADE_OPTIONS } from "@/lib/constants";

interface GradeDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  courseCode: string;
  courseName: string;
  currentGrade?: string;
  onSave: (grade: string) => void;
  onClear?: () => void;
}

export default function GradeDialog({
  open,
  onOpenChange,
  courseCode,
  courseName,
  currentGrade,
  onSave,
  onClear,
}: GradeDialogProps) {
  const [selectedGrade, setSelectedGrade] = useState<string>("");

  // Update selected grade when dialog opens or current grade changes
  useEffect(() => {
    if (open) {
      setSelectedGrade(currentGrade || "");
    }
  }, [open, currentGrade]);

  const handleSave = () => {
    if (selectedGrade) {
      onSave(selectedGrade);
      onOpenChange(false);
    }
  };

  const handleClear = () => {
    setSelectedGrade("");
    if (onClear) {
      onClear();
    }
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader className="space-y-3">
          <DialogTitle className="flex items-center gap-2 text-xl">
            <div className="p-2 bg-chula-soft rounded-lg">
              <GraduationCap className="h-5 w-5 text-chula" />
            </div>
            Set Course Grade
          </DialogTitle>
          <DialogDescription className="text-base">
            <span className="font-semibold text-gray-900">{courseCode}</span>
            <span className="text-gray-500 mx-2">•</span>
            <span className="text-gray-600">{courseName}</span>
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6 py-6">
          <div className="space-y-3">
            <Label htmlFor="grade" className="text-base font-medium">
              Select Letter Grade
            </Label>
            <Select value={selectedGrade} onValueChange={setSelectedGrade}>
              <SelectTrigger
                id="grade"
                className="w-full h-12 text-base border-2 focus:border-chula focus:ring-chula"
              >
                <SelectValue placeholder="Choose a grade..." />
              </SelectTrigger>
              <SelectContent>
                {GRADE_OPTIONS.map((option) => (
                  <SelectItem
                    key={option.value}
                    value={option.value}
                    className="text-base"
                  >
                    <span className="font-semibold">{option.label}</span>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {selectedGrade && (
              <p className="text-sm text-gray-600">
                Current selection: <span className="font-semibold text-chula">{selectedGrade}</span>
              </p>
            )}
          </div>
        </div>

        <div className="flex gap-3 pt-4 border-t">
          {currentGrade ? (
            <>
              <Button
                type="button"
                variant="outline"
                onClick={handleClear}
                className="flex-1 h-11 border-2 border-danger text-danger hover:bg-danger-soft hover:border-danger"
              >
                Clear Grade
              </Button>
              <Button
                type="button"
                onClick={handleSave}
                disabled={!selectedGrade}
                className="flex-1 h-11 bg-chula hover:bg-chula-hover text-white font-medium shadow-sm disabled:opacity-50"
              >
                Save Grade
              </Button>
            </>
          ) : (
            <>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                className="flex-1 h-11 border-2"
              >
                Cancel
              </Button>
              <Button
                type="button"
                onClick={handleSave}
                disabled={!selectedGrade}
                className="flex-1 h-11 bg-chula hover:bg-chula-hover text-white font-medium shadow-sm disabled:opacity-50"
              >
                Save Grade
              </Button>
            </>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
