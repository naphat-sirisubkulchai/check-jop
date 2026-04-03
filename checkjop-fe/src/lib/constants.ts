export const GRADE_OPTIONS = [
  { value: "A", label: "A" },
  { value: "B+", label: "B+" },
  { value: "B", label: "B" },
  { value: "C+", label: "C+" },
  { value: "C", label: "C" },
  { value: "D+", label: "D+" },
  { value: "D", label: "D" },
  { value: "F", label: "F" },
  { value: "S", label: "S (Satisfactory)" },
  { value: "U", label: "U (Unsatisfactory)" },
] as const;

// Credit limit constants
export const CREDIT_LIMITS = {
  REGULAR_SEMESTER: 22,
  SUMMER_SEMESTER: 10,
} as const;

// Semester type identifiers
export const SEMESTER_TYPES = {
  FIRST: 1,
  SECOND: 2,
  SUMMER: 3,
} as const;
