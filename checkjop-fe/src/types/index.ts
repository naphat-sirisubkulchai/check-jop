/* eslint-disable @typescript-eslint/no-explicit-any */
type Course = {
  code: string;
  nameEN: string;
  nameTH: string;
  credits: number;
  categoryId?: string; // Assuming category_id is part of the course data
  // Additional fields for course dependency graph
  prerequisites?: string;
  corequisites?: string;
  categories?: string;
  curriculum?: string;
  year?: number;
};

type Plan = {
  course_code: string;
  course_name?: string; // Optional field for manually added courses without a code match
  academicYear: number; // 2025,2026
  yearOfStudy: number; //1,2,3,4
  semester: number; // 1,2,3
  grade?: string; // "A", "B+", etc. (optional)
  credits: number;
  category_name?: string; // (optional)
};

type Category = {
  id: string;
  nameTH: string;
  nameEN: string;
  minCredits: number; // จำนวนหน่วยกิตขั้นต่ำในหมวดนี้
  courses?: Course[]; // รายวิชาที่อยู่ในหมวดนี้
};

type Curriculum = {
  id: string;
  nameTH: string;
  nameEN: string;
  year: number;
  minTotalCredits: number; // จำนวนหน่วยกิตขั้นต่ำทั้งหมด
  isActive: boolean;
  categories?: any[];
  courses?: Course[];
};

// Graduation calculation result types
type CategoryResult = {
  category_name: string;
  earned_credits: number;
  required_credits: number;
  is_satisfied: boolean;
};

type PrerequisiteViolation = {
  course_code: string;
  missing_prereqs: string[];
  prereqs_taken_in_wrong_term: string[];
  taken_in_wrong_term: boolean;
  missing_coreqs: string[];
  coreqs_taken_in_wrong_term: string[];
};

type CreditLimitViolation = {
  year: number;
  semester: number;
  credits: number;
  max_credits: number;
};

type GraduationResult = {
  can_graduate: boolean;
  gpax: number;
  total_credits: number;
  required_credits: number;
  category_results: CategoryResult[];
  missing_courses: string[];
  unrecognized_courses: string[];
  missing_catalog_years: number[];
  catalog_year_fallbacks: Record<number, number>; // missing year → catalog year actually used
  prerequisite_violations: PrerequisiteViolation[];
  credit_limit_violations: CreditLimitViolation[];
};

type CreateCurriculumForm = {
  curriculumFile: File | null;
  categoryFile: File | null;
  courseFiles: Array<{ file: File; year: number }>;
  previewData: null | {
    curriculums: any[];
    categories: any[];
    courses: any[];
    totalCurriculums: number;
    totalCategories: number;
    totalCourses: number;
  };
};

export type { Course, Plan, Category, Curriculum, GraduationResult, CategoryResult, PrerequisiteViolation, CreditLimitViolation, CreateCurriculumForm };
