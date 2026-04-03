"use client";

import { useEffect } from "react";
import { useAppStore } from "@/store/appStore";
import { courseApi } from "@/api/courseApi";

/**
 * Component that auto-loads study plan from cookie on app initialization
 * Handles the flow:
 * 1. Load study plan from cookie
 * 2. Fetch curriculum courses from API if curriculum is found
 * 3. Set courses in store
 */
export default function StudyPlanInitializer() {
  const { loadStudyPlanFromCookie, setCourses, setCategories, setSelectedCurriculum } = useAppStore();

  useEffect(() => {
    const initializeStudyPlan = async () => {
      // Step 1: Load study plan from cookie
      const loaded = loadStudyPlanFromCookie();

      if (!loaded) {
        console.log("No saved study plan found in cookie");
        return;
      }

      console.log("Study plan loaded from cookie");

      // Step 2: Get curriculum name from the loaded data
      const curriculumName = useAppStore.getState().selectedCurriculum?.nameTH;

      if (!curriculumName) {
        console.log("No curriculum found in saved plan");
        return;
      }

      // Step 3: Fetch curriculum courses from API
      try {
        console.log("Fetching curriculum courses for:", curriculumName);
        const curriculumData = await courseApi.getCurriculumByName(curriculumName);

        if (!curriculumData) {
          console.log("No curriculum data returned from API for:", curriculumName);
          return;
        }

        if (curriculumData.courses && curriculumData.categories) {
          setCourses(curriculumData.courses);
          setCategories(curriculumData.categories);
          // Update selectedCurriculum with fresh ID from API (cookie may have stale ID)
          const stale = useAppStore.getState().selectedCurriculum;
          if (stale && curriculumData.id && stale.id !== curriculumData.id) {
            setSelectedCurriculum({ ...stale, id: curriculumData.id });
          }
          console.log("Courses loaded from API:", curriculumData.courses.length);
        }
      } catch (error) {
        console.error("Error loading curriculum courses:", error);
      }
    };

    initializeStudyPlan();
  }, [loadStudyPlanFromCookie, setCourses, setCategories]);

  // This component doesn't render anything
  return null;
}
