import { Plan, Curriculum } from "@/types";

/**
 * Export/Import data structure for study plans
 */
export interface StudyPlanExportData {
  version: string; // Format version for backward compatibility
  exportDate: string; // ISO timestamp
  studyPlan: Plan[];
  exemptions: string[];
  selectedCurriculum: Curriculum | null;
  yearMapping?: { [key: number]: number } | null; // Maps yearOfStudy to academicYear
  metadata?: {
    totalCredits?: number;
    coursesCount?: number;
  };
}

/**
 * Validation result for import
 */
export interface ValidationResult {
  isValid: boolean;
  errors: string[];
  warnings: string[];
}

/**
 * Export study plan data to JSON format
 */
export function exportStudyPlan(
  studyPlan: Plan[],
  exemptions: string[],
  selectedCurriculum: Curriculum | null,
  yearMapping?: { [key: number]: number } | null
): StudyPlanExportData {
  const totalCredits = studyPlan.reduce((sum, plan) => sum + plan.credits, 0);

  return {
    version: "1.0.0",
    exportDate: new Date().toISOString(),
    studyPlan,
    exemptions,
    selectedCurriculum,
    yearMapping,
    metadata: {
      totalCredits,
      coursesCount: studyPlan.length,
    },
  };
}

/**
 * Download study plan as JSON file
 */
export function downloadStudyPlanJSON(data: StudyPlanExportData): void {
  const jsonString = JSON.stringify(data, null, 2);
  const blob = new Blob([jsonString], { type: "application/json" });
  const url = URL.createObjectURL(blob);

  const link = document.createElement("a");
  link.href = url;

  // Generate filename with curriculum and date
  const curriculumName = data.selectedCurriculum?.nameEN || "study-plan";
  const dateStr = new Date().toISOString().split("T")[0];
  link.download = `${curriculumName.replace(/\s+/g, "-")}-${dateStr}.json`;

  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}

/**
 * Validate imported study plan data
 */
export function validateImportData(data: any): ValidationResult {
  const errors: string[] = [];
  const warnings: string[] = [];

  // Check if data is an object
  if (!data || typeof data !== "object") {
    errors.push("Invalid data format: must be a JSON object");
    return { isValid: false, errors, warnings };
  }

  // Check required fields
  if (!data.version) {
    errors.push("Missing version field");
  }

  if (!Array.isArray(data.studyPlan)) {
    errors.push("studyPlan must be an array");
  } else {
    // Validate each plan entry
    data.studyPlan.forEach((plan: any, index: number) => {
      if (!plan.course_code) {
        errors.push(`Plan at index ${index}: missing course_code`);
      }
      if (typeof plan.academicYear !== "number") {
        errors.push(`Plan at index ${index}: academicYear must be a number`);
      }
      if (typeof plan.yearOfStudy !== "number") {
        errors.push(`Plan at index ${index}: yearOfStudy must be a number`);
      }
      if (typeof plan.semester !== "number") {
        errors.push(`Plan at index ${index}: semester must be a number`);
      }
      if (typeof plan.credits !== "number") {
        errors.push(`Plan at index ${index}: credits must be a number`);
      }
    });
  }

  if (!Array.isArray(data.exemptions)) {
    errors.push("exemptions must be an array");
  }

  // Warnings for optional fields
  if (data.selectedCurriculum === null) {
    warnings.push("No curriculum selected in imported data");
  }

  if (!data.exportDate) {
    warnings.push("Export date not found");
  }

  // Version compatibility check
  if (data.version && data.version !== "1.0.0") {
    warnings.push(`Import data version (${data.version}) may not be fully compatible`);
  }

  return {
    isValid: errors.length === 0,
    errors,
    warnings,
  };
}

/**
 * Parse imported JSON file
 */
export function parseImportFile(fileContent: string): {
  data: StudyPlanExportData | null;
  validation: ValidationResult;
} {
  try {
    const data = JSON.parse(fileContent);
    const validation = validateImportData(data);

    if (validation.isValid) {
      return { data: data as StudyPlanExportData, validation };
    }

    return { data: null, validation };
  } catch (error) {
    return {
      data: null,
      validation: {
        isValid: false,
        errors: ["Failed to parse JSON: " + (error as Error).message],
        warnings: [],
      },
    };
  }
}

/**
 * Read file as text
 */
export function readFileAsText(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target?.result as string;
      resolve(content);
    };
    reader.onerror = () => reject(new Error("Failed to read file"));
    reader.readAsText(file);
  });
}

/**
 * ---------- Cookie Management ----------
 */

/**
 * Cookie management for study plan persistence
 */
const COOKIE_NAME = "checkjop_study_plan";
const COOKIE_MAX_AGE = 60 * 60 * 24 * 7; // 30 days

/**
 * Save study plan to cookie
 */
export function saveStudyPlanToCookie(
  studyPlan: Plan[],
  exemptions: string[],
  selectedCurriculum: Curriculum | null,
  yearMapping?: { [key: number]: number } | null
): boolean {
  try {
    const data = exportStudyPlan(studyPlan, exemptions, selectedCurriculum, yearMapping);
    const jsonString = JSON.stringify(data);

    // Check if data is too large (cookies have ~4KB limit per cookie)
    if (jsonString.length > 4000) {
      console.warn("Study plan data too large for cookie storage");
      return false;
    }

    document.cookie = `${COOKIE_NAME}=${encodeURIComponent(jsonString)}; max-age=${COOKIE_MAX_AGE}; path=/; SameSite=Lax`;
    return true;
  } catch (error) {
    console.error("Failed to save study plan to cookie:", error);
    return false;
  }
}

/**
 * Load study plan from cookie
 */
export function loadStudyPlanFromCookie(): {
  data: StudyPlanExportData | null;
  validation: ValidationResult;
} {
  try {
    const cookies = document.cookie.split(";");
    const studyPlanCookie = cookies.find((cookie) =>
      cookie.trim().startsWith(`${COOKIE_NAME}=`)
    );

    if (!studyPlanCookie) {
      return {
        data: null,
        validation: {
          isValid: false,
          errors: ["No saved study plan found"],
          warnings: [],
        },
      };
    }

    const jsonString = decodeURIComponent(
      studyPlanCookie.split("=")[1]
    );

    return parseImportFile(jsonString);
  } catch (error) {
    console.error("Failed to load study plan from cookie:", error);
    return {
      data: null,
      validation: {
        isValid: false,
        errors: ["Failed to parse cookie data"],
        warnings: [],
      },
    };
  }
}

/**
 * Clear study plan cookie
 */
export function clearStudyPlanCookie(): void {
  document.cookie = `${COOKIE_NAME}=; max-age=0; path=/`;
}
