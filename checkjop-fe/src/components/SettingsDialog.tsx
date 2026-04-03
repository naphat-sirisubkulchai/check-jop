"use client";

import { useState, useEffect, useCallback, useMemo } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAppStore } from "@/store/appStore";
import { courseApi } from "@/api/courseApi";
import { Curriculum } from "@/types";
import { AlertCircle, BookOpen, Loader2, TriangleAlert } from "lucide-react";

type YearMapping = { [key: number]: string };

const YEAR_CONFIG = [
  { id: 1, label: "Year 1" },
  { id: 2, label: "Year 2" },
  { id: 3, label: "Year 3" },
  { id: 4, label: "Year 4" },
] as const;

const INITIAL_YEAR_MAPPING: YearMapping = { 1: "", 2: "", 3: "", 4: "" };

const START_YEAR = 2565;
const CURRENT_THAI_YEAR = new Date().getFullYear() + 543;
const END_YEAR = CURRENT_THAI_YEAR + 4;

const generateAcademicYears = () =>
  Array.from({ length: END_YEAR - START_YEAR + 1 }, (_, i) => START_YEAR + i);

const convertYearMapping = (mapping: { [key: number]: number } | null): YearMapping => {
  if (!mapping) return INITIAL_YEAR_MAPPING;
  return {
    1: mapping[1]?.toString() || "",
    2: mapping[2]?.toString() || "",
    3: mapping[3]?.toString() || "",
    4: mapping[4]?.toString() || "",
  };
};

interface SettingsDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export default function SettingsDialog({ open, onOpenChange }: SettingsDialogProps) {
  const {
    setSelectedCurriculum,
    setCourses,
    setCategories,
    setYearMapping,
    setStudyPlan,
    setExemptions,
    selectedCurriculum,
    yearMapping: storedYearMapping,
    studyPlan,
  } = useAppStore();

  const [curriculums, setCurriculums] = useState<Curriculum[]>([]);
  const [selectedCurriculumId, setSelectedCurriculumId] = useState("");
  const [yearMapping, setLocalYearMapping] = useState<YearMapping>(INITIAL_YEAR_MAPPING);
  const [isLoading, setIsLoading] = useState(false);
  const [isFetchingCurriculums, setIsFetchingCurriculums] = useState(true);
  const [errors, setErrors] = useState<{ curriculum?: string; yearMapping?: string; general?: string }>({});
  const [showResetConfirm, setShowResetConfirm] = useState(false);

  const academicYears = useMemo(() => generateAcademicYears(), []);

  useEffect(() => {
    if (!open) return;
    let isMounted = true;
    setIsFetchingCurriculums(true);
    courseApi.getAllCurriculaWithout().then((data) => {
      if (isMounted) setCurriculums(data);
    }).catch(() => {
      if (isMounted) setErrors({ general: "Failed to load curriculums." });
    }).finally(() => {
      if (isMounted) setIsFetchingCurriculums(false);
    });
    return () => { isMounted = false; };
  }, [open]);

  useEffect(() => {
    if (selectedCurriculum?.id) setSelectedCurriculumId(selectedCurriculum.id);
    if (storedYearMapping) setLocalYearMapping(convertYearMapping(storedYearMapping));
  }, [selectedCurriculum, storedYearMapping]);

  const handleYearChange = useCallback((yearId: number, value: string) => {
    setLocalYearMapping((prev) => {
      const next = { ...prev, [yearId]: value };
      if (yearId === 1 && value) {
        const base = parseInt(value);
        next[2] = (base + 1).toString();
        next[3] = (base + 2).toString();
        next[4] = (base + 3).toString();
      }
      return next;
    });
    setErrors({});
  }, []);

  const validate = useCallback(() => {
    const errs: typeof errors = {};
    if (!selectedCurriculumId) errs.curriculum = "Please select a curriculum";
    if (!Object.values(yearMapping).every((y) => y !== "")) errs.yearMapping = "Please select all academic years";
    if (Object.keys(errs).length > 0) { setErrors(errs); return false; }
    return true;
  }, [selectedCurriculumId, yearMapping]);

  const hasSettingsChanged = useCallback(() => {
    if (!selectedCurriculum || !storedYearMapping) return false;
    return (
      selectedCurriculumId !== selectedCurriculum.id ||
      parseInt(yearMapping[1]) !== storedYearMapping[1] ||
      parseInt(yearMapping[2]) !== storedYearMapping[2] ||
      parseInt(yearMapping[3]) !== storedYearMapping[3] ||
      parseInt(yearMapping[4]) !== storedYearMapping[4]
    );
  }, [selectedCurriculumId, yearMapping, selectedCurriculum, storedYearMapping]);

  const proceed = useCallback(async (shouldReset: boolean) => {
    setIsLoading(true);
    setShowResetConfirm(false);
    try {
      const curriculum = curriculums.find((c) => c.id === selectedCurriculumId);
      if (!curriculum) throw new Error("Curriculum not found");
      const curriculumData = await courseApi.getCurriculumByName(curriculum.nameTH);
      if (curriculumData) {
        setSelectedCurriculum(curriculum);
        setCourses(curriculumData.courses || []);
        setCategories(curriculumData.categories || []);
        setYearMapping({
          1: parseInt(yearMapping[1]),
          2: parseInt(yearMapping[2]),
          3: parseInt(yearMapping[3]),
          4: parseInt(yearMapping[4]),
        });
        if (shouldReset) { setStudyPlan([]); setExemptions([]); }
        onOpenChange(false);
      }
    } catch {
      setErrors({ general: "Failed to load curriculum data." });
    } finally {
      setIsLoading(false);
    }
  }, [curriculums, selectedCurriculumId, yearMapping, setSelectedCurriculum, setCourses, setCategories, setYearMapping, setStudyPlan, setExemptions, onOpenChange]);

  const handleSave = useCallback(() => {
    if (!validate()) return;
    if (studyPlan.length > 0 && hasSettingsChanged()) {
      setShowResetConfirm(true);
      return;
    }
    proceed(false);
  }, [validate, studyPlan, hasSettingsChanged, proceed]);

  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Settings</DialogTitle>
          </DialogHeader>

          <div className="space-y-4">
            {errors.general && (
              <div className="p-3 bg-red-50 border-l-4 border-red-500 rounded flex items-center gap-2">
                <AlertCircle className="h-4 w-4 text-red-600 shrink-0" />
                <p className="text-sm text-red-800">{errors.general}</p>
              </div>
            )}

            {/* Curriculum */}
            <div className="space-y-2">
              <label className="text-sm font-semibold text-gray-700">Curriculum</label>
              <Select value={selectedCurriculumId} onValueChange={(v) => { setSelectedCurriculumId(v); setErrors({}); }} disabled={isFetchingCurriculums}>
                <SelectTrigger className={`w-full h-10 ${errors.curriculum ? "border-red-500" : ""}`}>
                  <SelectValue placeholder={isFetchingCurriculums ? "Loading..." : "Choose curriculum..."} />
                </SelectTrigger>
                <SelectContent>
                  {curriculums.length === 0 && !isFetchingCurriculums ? (
                    <div className="px-4 py-6 text-center text-sm text-gray-500">
                      <BookOpen className="h-8 w-8 mx-auto mb-2 text-gray-400" />
                      <p>No curriculums available</p>
                    </div>
                  ) : (
                    curriculums.map((c) => (
                      <SelectItem key={c.id} value={c.id}>{c.nameTH}</SelectItem>
                    ))
                  )}
                </SelectContent>
              </Select>
              {errors.curriculum && <p className="text-xs text-red-600 flex items-center gap-1"><AlertCircle className="h-3 w-3" />{errors.curriculum}</p>}
            </div>

            {/* Year Mapping */}
            <div className="space-y-2">
              <label className="text-sm font-semibold text-gray-700">Academic Year Mapping (พ.ศ.)</label>
              <p className="text-xs text-gray-500">Select Year 1 and other years will auto-populate</p>
              <div className="grid grid-cols-2 gap-3">
                {YEAR_CONFIG.map((year) => (
                  <div key={year.id} className="flex items-center gap-2 bg-gray-50 rounded-lg p-3 border border-gray-200">
                    <span className="text-sm font-medium text-gray-700 w-14 shrink-0">{year.label}</span>
                    <Select value={yearMapping[year.id]} onValueChange={(v) => handleYearChange(year.id, v)} disabled={isLoading}>
                      <SelectTrigger className={`flex-1 h-8 text-sm ${errors.yearMapping && !yearMapping[year.id] ? "border-red-300" : ""}`}>
                        <SelectValue placeholder="Year" />
                      </SelectTrigger>
                      <SelectContent>
                        {academicYears.map((y) => (
                          <SelectItem key={y} value={y.toString()}>{y}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                ))}
              </div>
              {errors.yearMapping && <p className="text-xs text-red-600 flex items-center gap-1"><AlertCircle className="h-3 w-3" />{errors.yearMapping}</p>}
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => onOpenChange(false)}>Cancel</Button>
            <Button
              onClick={handleSave}
              disabled={isLoading || isFetchingCurriculums}
              className="bg-gradient-to-r from-chula-active to-pink-500 text-white"
            >
              {isLoading ? <><Loader2 className="h-4 w-4 animate-spin mr-2" />Saving...</> : "Save"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Reset Confirmation */}
      <Dialog open={showResetConfirm} onOpenChange={setShowResetConfirm}>
        <DialogContent>
          <DialogHeader>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-amber-100 flex items-center justify-center shrink-0">
                <TriangleAlert className="h-5 w-5 text-amber-600" />
              </div>
              <DialogTitle>Reset Study Plan?</DialogTitle>
            </div>
          </DialogHeader>
          <p className="text-sm text-gray-600">
            You have <strong>{studyPlan.length} courses</strong> in your plan. Changing curriculum or years will reset everything.
          </p>
          <DialogFooter className="gap-2">
            <Button variant="outline" onClick={() => setShowResetConfirm(false)}>Cancel</Button>
            <Button variant="destructive" onClick={() => proceed(true)}>Reset & Save</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
