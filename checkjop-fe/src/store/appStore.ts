import { create } from "zustand";
import { Curriculum, Plan } from "@/types";
import {
  exportStudyPlan,
  downloadStudyPlanJSON,
  StudyPlanExportData,
  saveStudyPlanToCookie,
  loadStudyPlanFromCookie,
  clearStudyPlanCookie,
} from "@/utils/exportImport";
import { toast } from "sonner";

/**
 * Unified app store for managing:
 * - Course data (categories, courses)
 * - Study plans and curriculum selection
 * - Graduation calculation results
 */
export const useAppStore = create(
  (set: any, get: any) => ({
      // Initial State
      courses: [] as any[],
      categories: [] as any[],
      curriculums: [] as Curriculum[],
      studyPlan: [] as Plan[],
      exemptions: [] as string[],
      selectedCurriculum: null as Curriculum | null,
      yearMapping: null as { [key: number]: number } | null, // Maps yearOfStudy to academicYear
      result: null as any,
      isLoading: false as boolean,
      error: null as string | null,
      totalCredits: 0 as number,
      printFormat: "transcript" as "transcript" | "summary",

      // Course Data Actions
      setCourses: (courses: any[]) => set({ courses }),
      setCategories: (categories: any[]) => set({ categories }),
      setCurriculums: (curriculums: Curriculum[]) => set({ curriculums }),

      //---------- Study Plan Actions ----------

      // get Study Plan
      getStudyPlan: () => get().studyPlan,
      // add a Course to Study Plan
      addCoursePlan: (plan: Plan) => {
        const existing = get().studyPlan.find(
          (p: any) => p.course_code === plan.course_code
        );
        if (existing && existing.grade !== "F") {
          console.warn("Course already in study plan for this semester");
          toast.error("Course already added", {
            description: `${plan.course_code} is already in Year ${existing.yearOfStudy}, Semester ${existing.semester}`,
          });
          return get().studyPlan;
        }
        set((state: any) => ({ studyPlan: [...state.studyPlan, plan] }));
        get().calculateTotalCredits();
        return get().studyPlan;
      },
      // remove a Course from Study Plan
      removeCoursePlan: (code: string) => {
        set((state: any) => ({
          studyPlan: state.studyPlan.filter((p: any) => p.course_code !== code),
          exemptions: state.exemptions.filter((e: string) => e !== code),
        }));
        get().calculateTotalCredits();
        return get().studyPlan;
      },
      //edit course in Study Plan
      editCoursePlan: (code: string, newPlan: Partial<Plan>, yearOfStudy?: number, semester?: number) => {
        set((state: any) => ({
          studyPlan: state.studyPlan.map((p: any) =>
            p.course_code === code &&
            (yearOfStudy === undefined || p.yearOfStudy === yearOfStudy) &&
            (semester === undefined || p.semester === semester)
              ? {...p, ...newPlan} : p,
          ),
        }));
        get().calculateTotalCredits();
        return get().studyPlan;
      },
      // set Study Plan
      setStudyPlan: (studyPlan: Plan[]) => {
        set({ studyPlan });
        get().calculateTotalCredits();
        return get().studyPlan;
      },
      // clear Study Plan
      clearStudyPlan: () => {
        set({ studyPlan: [], exemptions: [] });
        get().calculateTotalCredits();
        return get().studyPlan;
      },

      //---------- Exemptions Actions ----------

      // get Exemptions
      getExemptions: () => get().exemptions,
      // set Exemptions
      setExemptions: (exemptions: string[]) => {
        set({ exemptions });
      },
      // add an Exemption
      addExemption: (code: string) => {
        set((state: any) => ({
          exemptions: Array.from(new Set([...state.exemptions, code])),
        }));
      },
      // remove an Exemption
      removeExemption: (code: string) => {
        set((state: any) => ({
          exemptions: state.exemptions.filter((c: any) => c !== code),
        }));
      },
      // clear Exemptions
      clearExemptions: () => {
        set({ exemptions: [] });
      },


      //---------- Curriculum Actions ----------
      setSelectedCurriculum: (curriculum: Curriculum | null) =>
        set({ selectedCurriculum: curriculum }),

      //---------- Year Mapping Actions ----------
      setYearMapping: (yearMapping: { [key: number]: number } | null) =>
        set({ yearMapping }),


      //---------- Additional Actions ----------
      // Result Actions
      setResult: (result: any) => set({ result }),
      setIsLoading: (isLoading: boolean) => set({ isLoading }),
      setError: (error: string | null) => set({ error }),
      setPrintFormat: (printFormat: "transcript" | "summary") => set({ printFormat }),

      // Selectors
      getCourseByCode: (code: string) => {
        const { courses } = get();

          const course = courses.find((c: any) => c.code === code);
          if (course) return course;

          const courseInPlan = get().studyPlan.find((p: any) => p.course_code === code);
          if (courseInPlan) return courseInPlan;

        console.log(`Course with code ${code} not found`);
        return undefined;
      },

      calculateTotalCredits: () => {
        const total = get().studyPlan.reduce((sum: number, course: any) => sum + course.credits, 0);
        set({ totalCredits: total });
        return total;
      },

      //---------- Export/Import Actions ----------

      // Export study plan to JSON and download
      exportStudyPlanToJSON: () => {
        const { studyPlan, exemptions, selectedCurriculum, yearMapping } = get();
        const exportData = exportStudyPlan(
          studyPlan,
          exemptions,
          selectedCurriculum,
          yearMapping
        );
        downloadStudyPlanJSON(exportData);
      },

      // Import study plan from parsed data
      importStudyPlanFromData: (data: StudyPlanExportData) => {
        const totalCredits = data.studyPlan.reduce((sum: number, course: any) => sum + course.credits, 0);
        set({
          studyPlan: data.studyPlan,
          exemptions: data.exemptions,
          selectedCurriculum: data.selectedCurriculum,
          yearMapping: data.yearMapping || null,
          totalCredits,
        });
      },

      //---------- Cookie Actions ----------

      // Save study plan to cookie
      saveStudyPlanToCookie: () => {
        const { studyPlan, exemptions, selectedCurriculum, yearMapping } = get();
        return saveStudyPlanToCookie(studyPlan, exemptions, selectedCurriculum, yearMapping);
      },

      // Load study plan from cookie
      loadStudyPlanFromCookie: () => {
        const { data, validation } = loadStudyPlanFromCookie();

        if (validation.isValid && data) {
          const totalCredits = data.studyPlan.reduce((sum: number, course: any) => sum + course.credits, 0);
          set({
            studyPlan: data.studyPlan,
            exemptions: data.exemptions,
            selectedCurriculum: data.selectedCurriculum,
            yearMapping: data.yearMapping || null,
            totalCredits,
          });
          console.log("Study plan loaded from cookie");
          return true;
        }

        console.log("No valid study plan found in cookie");
        return false;
      },

      // Clear study plan cookie
      clearStudyPlanCookie: () => {
        clearStudyPlanCookie();
      },
  }),
);

// Export type for external use
export type AppStore = ReturnType<typeof useAppStore.getState>;