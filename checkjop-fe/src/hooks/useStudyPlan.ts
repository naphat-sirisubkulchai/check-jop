import { useCallback } from "react";
import { useAppStore } from "@/store/appStore";
import { gradeService } from "@/api/gradApi";

/**
 * Hook for managing study plan operations
 * Handles export/import, persistence, and graduation calculations
 */
export function useStudyPlan() {
  const {
    studyPlan,
    exemptions,
    selectedCurriculum,
    setResult,
    setIsLoading,
    setError,
    getCourseByCode,
    saveStudyPlanToCookie,
  } = useAppStore();

  const handleCalculate = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    setResult(null);

    try {
      // Validation
      if (!selectedCurriculum) {
        throw new Error("Please select a curriculum before calculating.");
      }

      if (!studyPlan.length) {
        throw new Error("Please add courses to your study plan before calculating.");
      }

      // Save to cookie before calculating
      saveStudyPlanToCookie();
      console.log("Study plan saved to cookie");

      // Prepare API payload
      const transformedStudyPlan = studyPlan.map(({academicYear, yearOfStudy, ...rest})=> ({
        ...rest,
        year: academicYear,
      }));

      // Calculate admission year as the earliest academic year in the study plan
      const admissionYear = Math.min(...studyPlan.map(course => course.academicYear));

      
      const apiPayload = {
        name_th: selectedCurriculum.nameTH,
        admission_year: admissionYear,
        courses: transformedStudyPlan,
        manual_credits: {},
        exemptions: exemptions,
      };

      console.log("API Payload:", JSON.stringify(apiPayload, null, 2));

      const response = await gradeService.checkGraduation(apiPayload);
      console.log(response);
      
      setResult(response);

    } catch (error) {
      const message = error instanceof Error
        ? error.message
        : "Failed to check graduation eligibility. Please try again.";
      
      console.error("Graduation check error:", error);
      setError(message);
    } finally {
      setIsLoading(false);
    }
  }, [selectedCurriculum, studyPlan, exemptions, getCourseByCode, setResult, setIsLoading, setError, saveStudyPlanToCookie]);

  return {
    handleCalculate,
  };
}
