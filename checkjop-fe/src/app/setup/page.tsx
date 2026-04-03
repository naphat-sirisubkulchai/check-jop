"use client";

import { useState, useEffect, useCallback, useMemo } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useAppStore } from "@/store/appStore";
import { courseApi } from "@/api/courseApi";
import { Curriculum } from "@/types";
import { ChevronRight, Loader2, AlertCircle, BookOpen, TriangleAlert } from "lucide-react";

// ==================== TYPES ====================
type YearMapping = {
  [key: number]: string;
};

type SetupErrors = {
  curriculum?: string;
  yearMapping?: string;
  general?: string;
};

// ==================== CONSTANTS ====================
const YEAR_CONFIG = [
  { id: 1, label: "Year 1" },
  { id: 2, label: "Year 2" },
  { id: 3, label: "Year 3" },
  { id: 4, label: "Year 4" },
] as const;

const INITIAL_YEAR_MAPPING: YearMapping = {
  1: "",
  2: "",
  3: "",
  4: "",
};

const START_YEAR = 2565;
const CURRENT_THAI_YEAR = new Date().getFullYear() + 543;
const END_YEAR = CURRENT_THAI_YEAR + 4;

// ==================== UTILS ====================
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

// ==================== COMPONENT ====================
export default function SetupPage() {
  const router = useRouter();

  // Store
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

  // State
  const [curriculums, setCurriculums] = useState<Curriculum[]>([]);
  const [selectedCurriculumId, setSelectedCurriculumId] = useState("");
  const [yearMapping, setLocalYearMapping] = useState<YearMapping>(INITIAL_YEAR_MAPPING);
  const [isLoading, setIsLoading] = useState(false);
  const [isFetchingCurriculums, setIsFetchingCurriculums] = useState(true);
  const [errors, setErrors] = useState<SetupErrors>({});
  const [showResetDialog, setShowResetDialog] = useState(false);

  // Computed values
  const academicYears = useMemo(() => generateAcademicYears(), []);

  // ==================== EFFECTS ====================
  // Fetch curriculums
  useEffect(() => {
    let isMounted = true;

    const fetchCurriculums = async () => {
      try {
        const data = await courseApi.getAllCurriculaWithout();
        if (isMounted) setCurriculums(data);
      } catch (error) {
        console.error("Error fetching curriculums:", error);
        if (isMounted) {
          setErrors({ general: "Failed to load curriculums. Please refresh the page." });
        }
      } finally {
        if (isMounted) setIsFetchingCurriculums(false);
      }
    };

    fetchCurriculums();
    return () => {
      isMounted = false;
    };
  }, []);

  // Load existing data from store
  useEffect(() => {
    if (selectedCurriculum?.id) {
      setSelectedCurriculumId(selectedCurriculum.id);
    }
    if (storedYearMapping) {
      setLocalYearMapping(convertYearMapping(storedYearMapping));
    }
  }, [selectedCurriculum, storedYearMapping]);

  // ==================== HANDLERS ====================
  const handleCurriculumChange = useCallback((value: string) => {
    setSelectedCurriculumId(value);
    setErrors({});
  }, []);

  const handleYearChange = useCallback((yearId: number, value: string) => {
    setLocalYearMapping((prev) => {
      const newMapping = { ...prev, [yearId]: value };

      // Auto-populate subsequent years if Year 1 is selected
      if (yearId === 1 && value) {
        const baseYear = parseInt(value);
        newMapping[2] = (baseYear + 1).toString();
        newMapping[3] = (baseYear + 2).toString();
        newMapping[4] = (baseYear + 3).toString();
      }

      return newMapping;
    });
    setErrors({});
  }, []);

  const validateForm = useCallback((): boolean => {
    const validationErrors: SetupErrors = {};

    if (!selectedCurriculumId) {
      validationErrors.curriculum = "Please select a curriculum to continue";
    }

    const hasAllYears = Object.values(yearMapping).every((year) => year !== "");
    if (!hasAllYears) {
      validationErrors.yearMapping = "Please select all academic years (1-4)";
    }

    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      return false;
    }

    return true;
  }, [selectedCurriculumId, yearMapping]);

  const hasExistingPlan = studyPlan.length > 0;

  const hasSettingsChanged = useCallback((): boolean => {
    if (!selectedCurriculum || !storedYearMapping) return false;
    const curriculumChanged = selectedCurriculumId !== selectedCurriculum.id;
    const yearChanged =
      parseInt(yearMapping[1]) !== storedYearMapping[1] ||
      parseInt(yearMapping[2]) !== storedYearMapping[2] ||
      parseInt(yearMapping[3]) !== storedYearMapping[3] ||
      parseInt(yearMapping[4]) !== storedYearMapping[4];
    return curriculumChanged || yearChanged;
  }, [selectedCurriculumId, yearMapping, selectedCurriculum, storedYearMapping]);

  const proceedWithPlanning = useCallback(async (shouldReset: boolean) => {
    setIsLoading(true);
    setShowResetDialog(false);

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

        if (shouldReset) {
          setStudyPlan([]);
          setExemptions([]);
        }

        router.push("/home");
      }
    } catch (error) {
      console.error("Error starting planning:", error);
      setErrors({
        general: "Failed to load curriculum data. Please check your connection and try again.",
      });
    } finally {
      setIsLoading(false);
    }
  }, [
    curriculums,
    selectedCurriculumId,
    yearMapping,
    setSelectedCurriculum,
    setCourses,
    setCategories,
    setYearMapping,
    setStudyPlan,
    setExemptions,
    router,
  ]);

  const handleStartPlanning = useCallback(async () => {
    if (!validateForm()) return;

    // Has existing plan and settings changed → show confirmation
    if (hasExistingPlan && hasSettingsChanged()) {
      setShowResetDialog(true);
      return;
    }

    // Has existing plan but nothing changed → go back without reset
    if (hasExistingPlan && !hasSettingsChanged()) {
      router.push("/home");
      return;
    }

    // No existing plan → proceed normally (reset)
    await proceedWithPlanning(true);
  }, [validateForm, hasExistingPlan, hasSettingsChanged, proceedWithPlanning, router]);

  // ==================== RENDER ====================
  return (
    <div className="flex items-center justify-center p-8">
      <div className="bg-white rounded-2xl shadow-2xl p-8 space-y-6 sm:p-10 w-full max-w-2xl border border-gray-100 transition-all duration-300 hover:shadow-3xl">
        {/* Header */}
        <Header />

        {/* General Error */}
        {errors.general && <ErrorAlert message={errors.general} />}

        {/* Curriculum Selection */}
        <CurriculumSection
          curriculums={curriculums}
          selectedId={selectedCurriculumId}
          isLoading={isFetchingCurriculums}
          error={errors.curriculum}
          onChange={handleCurriculumChange}
        />

        {/* Year Mapping */}
        <YearMappingSection
          yearMapping={yearMapping}
          academicYears={academicYears}
          error={errors.yearMapping}
          isLoading={isLoading}
          onChange={handleYearChange}
        />

        {/* Submit Button */}
        <SubmitButton
          isLoading={isLoading}
          isFetchingCurriculums={isFetchingCurriculums}
          onClick={handleStartPlanning}
        />

        {/* Help Text */}
        <HelpText />
      </div>

      {/* Reset Confirmation Dialog */}
      <Dialog open={showResetDialog} onOpenChange={setShowResetDialog}>
        <DialogContent>
          <DialogHeader>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-amber-100 flex items-center justify-center flex-shrink-0">
                <TriangleAlert className="h-5 w-5 text-amber-600" />
              </div>
              <DialogTitle>Reset Study Plan?</DialogTitle>
            </div>
            <DialogDescription className="pt-2">
              You have an existing study plan with <strong>{studyPlan.length} courses</strong>.
              Changing the curriculum or academic years will reset your entire plan. This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="gap-2">
            <Button
              variant="outline"
              onClick={() => setShowResetDialog(false)}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={() => proceedWithPlanning(true)}
            >
              Reset & Continue
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}

// ==================== SUB-COMPONENTS ====================
function Header() {
  return (
    <div className="text-center">
      {/* <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-chula-active to-pink-500 rounded-2xl mb-4 shadow-lg">
        <BookOpen className="h-8 w-8 text-white" />
      </div> */}
      <h1 className="text-3xl font-bold text-gray-900 mb-3">Setup Your Study Plan</h1>
      <p className="text-sm text-gray-600 max-w-xl mx-auto">
        Select your curriculum and map your academic years to begin planning your journey
      </p>
    </div>
  );
}

function ErrorAlert({ message }: { message: string }) {
  return (
    <div className="mb-6 p-4 bg-red-50 border-l-4 border-red-500 rounded-lg flex items-start gap-3 animate-in slide-in-from-top-2">
      <AlertCircle className="h-5 w-5 text-red-600 flex-shrink-0 mt-0.5" />
      <p className="text-sm text-red-800 font-medium">{message}</p>
    </div>
  );
}

function CurriculumSection({
  curriculums,
  selectedId,
  isLoading,
  error,
  onChange,
}: {
  curriculums: Curriculum[];
  selectedId: string;
  isLoading: boolean;
  error?: string;
  onChange: (value: string) => void;
}) {
  return (
    <div className="bg-gray-50 rounded-xl px-6 py-4 border border-gray-200">
      <div className="flex items-center gap-2 mb-4">
        <div className="w-8 h-8 rounded-lg bg-chula-active/10 flex items-center justify-center">
          <span className="text-chula-active font-bold">1</span>
        </div>
        <label htmlFor="curriculum-select" className="text-lg font-semibold text-gray-900">
          Select Your Curriculum
        </label>
      </div>

      <Select value={selectedId} onValueChange={onChange} disabled={isLoading}>
        <SelectTrigger
          id="curriculum-select"
          className={`w-full h-12 bg-white transition-all ${
            error
              ? "border-red-500 focus:ring-red-500"
              : "border-gray-300 hover:border-chula-active focus:border-chula-active focus:ring-chula-active"
          }`}
          aria-invalid={!!error}
          aria-describedby={error ? "curriculum-error" : undefined}
        >
          <SelectValue
            placeholder={isLoading ? "Loading curriculums..." : "Choose your curriculum..."}
          />
        </SelectTrigger>
        <SelectContent>
          {curriculums.length === 0 && !isLoading ? (
            <div className="px-4 py-8 text-center text-sm text-gray-500">
              <BookOpen className="h-10 w-10 mx-auto mb-3 text-gray-400" />
              <p className="font-medium">No curriculums available</p>
              <p className="text-xs mt-1">Please contact support</p>
            </div>
          ) : (
            curriculums.map((curriculum) => (
              <SelectItem key={curriculum.id} value={curriculum.id} className="py-3">
                <span className="font-medium">{curriculum.nameTH}</span>
              </SelectItem>
            ))
          )}
        </SelectContent>
      </Select>

      {error && (
        <p
          id="curriculum-error"
          className="mt-2 text-sm text-red-600 flex items-center gap-1.5 animate-in slide-in-from-top-1"
        >
          <AlertCircle className="h-4 w-4" />
          {error}
        </p>
      )}
    </div>
  );
}

function YearMappingSection({
  yearMapping,
  academicYears,
  error,
  isLoading,
  onChange,
}: {
  yearMapping: YearMapping;
  academicYears: number[];
  error?: string;
  isLoading: boolean;
  onChange: (yearId: number, value: string) => void;
}) {
  return (
    <div className="bg-gray-50 rounded-xl px-6 py-4 border border-gray-200">
      <div className="flex items-center gap-2 mb-2">
        <div className="w-8 h-8 rounded-lg bg-chula-active/10 flex items-center justify-center">
          <span className="text-chula-active font-bold">2</span>
        </div>
        <h2 className="text-lg font-semibold text-gray-900">Academic Year Mapping (พ.ศ.)</h2>
      </div>
      <p className="text-sm text-gray-600 mb-5 ml-10">
        Select Year 1 and other years will auto-populate
      </p>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {YEAR_CONFIG.map((year) => (
          <div
            key={year.id}
            className="flex items-baseline gap-2 bg-white rounded-lg p-4 border border-gray-200 transition-all hover:border-chula-active hover:shadow-sm"
          >
            <label
              htmlFor={`year-${year.id}-select`}
              className="block text-sm font-semibold text-gray-700 mb-2 shrink-0"
            >
              {year.label}
            </label>
            <Select
              value={yearMapping[year.id]}
              onValueChange={(value) => onChange(year.id, value)}
              disabled={isLoading}
            >
              <SelectTrigger
                id={`year-${year.id}-select`}
                className={`w-full transition-all ${
                  error && !yearMapping[year.id]
                    ? "border-red-300 focus:ring-red-500"
                    : "hover:border-chula-active focus:border-chula-active focus:ring-chula-active"
                } ${yearMapping[year.id] ? "bg-chula-soft/20" : "bg-white"}`}
                aria-label={`Select academic year for ${year.label}`}
              >
                <SelectValue placeholder="Select year" />
              </SelectTrigger>
              <SelectContent>
                {academicYears.map((academicYear) => (
                  <SelectItem key={academicYear} value={academicYear.toString()}>
                    {academicYear}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        ))}
      </div>

      {error && (
        <p className="mt-4 text-sm text-red-600 flex items-center gap-1.5 animate-in slide-in-from-top-1">
          <AlertCircle className="h-4 w-4" />
          {error}
        </p>
      )}
    </div>
  );
}

function SubmitButton({
  isLoading,
  isFetchingCurriculums,
  onClick,
}: {
  isLoading: boolean;
  isFetchingCurriculums: boolean;
  onClick: () => void;
}) {
  return (
    <Button
      onClick={onClick}
      disabled={isLoading || isFetchingCurriculums}
      className="w-full bg-gradient-to-r from-chula-active to-pink-500 hover:from-chula-active/90 hover:to-pink-600 disabled:from-gray-300 disabled:to-gray-300 disabled:cursor-not-allowed text-white py-6 rounded-xl text-lg font-semibold transition-all duration-200 flex items-center justify-center gap-2 shadow-lg hover:shadow-xl hover:scale-[1.02] active:scale-[0.98]"
      aria-label="Start planning your study schedule"
    >
      {isLoading ? (
        <>
          <Loader2 className="h-5 w-5 animate-spin" />
          Loading Curriculum...
        </>
      ) : (
        <>
          Start Planning
          <ChevronRight className="h-5 w-5" />
        </>
      )}
    </Button>
  );
}

function HelpText() {
  return (
    <p className="text-center text-xs text-gray-500">
      Need help? Check our{" "}
      <button className="text-chula-active hover:underline font-medium">documentation</button>
    </p>
  );
}
