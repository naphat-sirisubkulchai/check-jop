import { Plan } from "@/types";

export interface TestCase {
  id: string;
  name: string;
  description: string;
  curriculumName: string;
  admissionYear: number;
  studyPlan: Plan[];
  exemptions: string[];
  manualCredits?: Record<string, number>;
  expectedResult?: {
    canGraduate: boolean;
    hasViolations: boolean;
    description: string;
  };
  implemented: boolean;
}

export const testCases: TestCase[] = [
  {
    id: "case-1",
    name: "กรณีพื้นฐาน - วิชาไม่มี Prerequisites",
    description: "วิชา 2301170 (คอมพิวเตอร์และการโปรแกรม) ไม่มี prerequisite สามารถลงได้โดยไม่มีการละเมิด",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301173",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 4,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: false,
      description: `{
  "can_graduate": false,
  "gpax": 0,
  "total_credits": 4,
  "required_credits": 136,
  "prerequisite_violations": [],
  "credit_limit_violations": []
}`,
    },
    implemented: true,
  },
  {
    id: "case-2a",
    name: "Complex Prerequisites - Violation (missing all)",
    description: "2301260 requires (2301170 AND 2301172) OR 2301173 — none provided",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301260",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 4,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: true,
      description: `{
  "can_graduate": false,
  "prerequisite_violations": [
    {
      "course_code": "2301260",
      "missing_prereqs": ["2301170", "2301173", "2301172"],
      "prereqs_taken_in_wrong_term": [],
      "taken_in_wrong_term": false,
      "missing_coreqs": [],
      "coreqs_taken_in_wrong_term": []
    }
  ]
}`,
    },
    implemented: true,
  },
  {
    id: "case-2b",
    name: "Complex Prerequisites - Pass (OR path)",
    description: "2301260 requires (2301170 AND 2301172) OR 2301173 — satisfied via 2301173",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301173",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 4,
      },
      {
        course_code: "2301260",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 2,
        credits: 4,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: false,
      description: `"prerequisite_violations": []`,
    },
    implemented: true,
  },
  {
    id: "case-2c",
    name: "Complex Prerequisites - Pass (AND path)",
    description: "2301260 requires (2301170 AND 2301172) OR 2301173 — satisfied via 2301170 AND 2301172",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301170",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
      },
      {
        course_code: "2301172",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 1,
      },
      {
        course_code: "2301260",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 2,
        credits: 4,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: false,
      description: `"prerequisite_violations": []`,
    },
    implemented: true,
  },
  {
    id: "case-3",
    name: "Strict Term Requirement - Violation",
    description: "2301173 and 2301260 in same term — prerequisite must be taken before",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301173",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 4,
      },
      {
        course_code: "2301260",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 4,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: true,
      description: `{
  "prerequisite_violations": [
    {
      "course_code": "2301260",
      "missing_prereqs": ["2301170", "2301172"],
      "prereqs_taken_in_wrong_term": ["2301173"],
      "taken_in_wrong_term": false,
      "missing_coreqs": [],
      "coreqs_taken_in_wrong_term": []
    }
  ]
}`,
    },
    implemented: true,
  },
  {
    id: "case-4a",
    name: "Transitive Prerequisites - All Missing",
    description: "2301263 requires 2301260, which requires (2301170 AND 2301172) OR 2301173 — all missing",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301263",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 4,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: true,
      description: `{
  "prerequisite_violations": [
    {
      "course_code": "2301263",
      "missing_prereqs": ["2301260", "2301170", "2301173", "2301172"],
      "prereqs_taken_in_wrong_term": [],
      "taken_in_wrong_term": false,
      "missing_coreqs": [],
      "coreqs_taken_in_wrong_term": []
    }
  ]
}`,
    },
    implemented: true,
  },
  {
    id: "case-4b",
    name: "Transitive Prerequisites - Partial",
    description: "Has 2301173 but missing 2301260 (direct prerequisite of 2301263)",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301173",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 4,
      },
      {
        course_code: "2301263",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 2,
        credits: 4,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: true,
      description: `{
  "prerequisite_violations": [
    {
      "course_code": "2301263",
      "missing_prereqs": ["2301260"],
      "prereqs_taken_in_wrong_term": [],
      "taken_in_wrong_term": false,
      "missing_coreqs": [],
      "coreqs_taken_in_wrong_term": []
    }
  ]
}`,
    },
    implemented: true,
  },
  {
    id: "case-5a",
    name: "Corequisites - Missing",
    description: "2301362 requires corequisite 2301279 OR 2301369 — none provided",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301362",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: true,
      description: `{
  "prerequisite_violations": [
    {
      "course_code": "2301362",
      "missing_prereqs": [],
      "prereqs_taken_in_wrong_term": [],
      "taken_in_wrong_term": false,
      "missing_coreqs": ["2301279", "2301369"],
      "coreqs_taken_in_wrong_term": []
    }
  ]
}`,
    },
    implemented: true,
  },
  {
    id: "case-5b",
    name: "Corequisites - Wrong Term",
    description: "2301279 taken in different term from 2301362 — corequisite not satisfied",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301279",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
      },
      {
        course_code: "2301362",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 2,
        credits: 3,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: true,
      description: `{
  "prerequisite_violations": [
    {
      "course_code": "2301362",
      "missing_prereqs": [],
      "prereqs_taken_in_wrong_term": [],
      "taken_in_wrong_term": true,
      "missing_coreqs": ["2301369"],
      "coreqs_taken_in_wrong_term": ["2301279"]
    }
  ]
}`,
    },
    implemented: true,
  },
  {
    id: "case-5c",
    name: "Corequisites - Pass",
    description: "2301279 and 2301362 taken in same term — corequisite satisfied",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301279",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
      },
      {
        course_code: "2301362",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: false,
      description: `{
  "prerequisite_violations": []
}`,
    },
    implemented: true,
  },
  {
    id: "case-5d",
    name: "Corequisites - Pass with Grade F",
    description: "2301279 (grade F) and 2301362 in same term — corequisite still satisfied",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301279",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
        grade: "F",
      },
      {
        course_code: "2301362",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: false,
      description: `{
  "prerequisite_violations": []
}`,
    },
    implemented: true,
  },
  {
    id: "case-6a",
    name: "Credit Limit - Normal Semester",
    description: "Exceeding 22 credits in normal semester (23 credits in sem 1)",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301173",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 15,
      },
      {
        course_code: "2301260",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 8,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: true,
      description: `{
  "credit_limit_violations": [
    {
      "year": 2566,
      "semester": 1,
      "credits": 23,
      "max_credits": 22
    }
  ]
}`,
    },
    implemented: true,
  },
  {
    id: "case-6b",
    name: "Credit Limit - Summer Semester",
    description: "Exceeding 10 credits in summer semester (12 credits in sem 3)",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301173",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 3,
        credits: 8,
      },
      {
        course_code: "2301260",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 3,
        credits: 4,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: true,
      description: `{
  "credit_limit_violations": [
    {
      "year": 2566,
      "semester": 3,
      "credits": 12,
      "max_credits": 10
    }
  ]
}`,
    },
    implemented: true,
  },
  {
    id: "case-9",
    name: "Cross Curriculum Prerequisites",
    description: "2301172 is from เอกเดี่ยว-2561, not in เอกเดี่ยว-ฝึกงาน-2566 — no corequisite check",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301170",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
      },
      {
        course_code: "2301172",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 2,
        credits: 1,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: false,
      description: `{
  "prerequisite_violations": []
}`,
    },
    implemented: true,
  },
  {
    id: "case-10",
    name: "Free Electives",
    description: "Unknown course with category_name 'วิชาเสรี' — skips prerequisite check, counts credits",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "9999999",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
        category_name: "วิชาเสรี",
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: false,
      description: `{
  "prerequisite_violations": [],
  "category_results": [
    {
      "category_name": "วิชาเสรี",
      "earned_credits": 3,
      "required_credits": 6,
      "is_satisfied": false
    }
  ]
}`,
    },
    implemented: true,
  },
  {
    id: "case-11a",
    name: "C.F. Permission - No Violation (C.F. only)",
    description: "2301290 has prerequisite C.F. only — no violation regardless of exemptions",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301290",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 1,
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: false,
      description: `{
  "prerequisite_violations": []
}`,
    },
    implemented: true,
  },
  {
    id: "case-11b",
    name: "C.F. Permission - With Exemption",
    description: "2301290 with exemption — explicitly granted C.F. permission",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301290",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 1,
      },
    ],
    exemptions: ["2301290"],
    expectedResult: {
      canGraduate: false,
      hasViolations: false,
      description: `{
  "prerequisite_violations": []
}`,
    },
    implemented: true,
  },
  {
    id: "case-12a",
    name: "GPAX - Normal Grades",
    description: "A(3cr) + B(3cr) + C+(2cr) = GPAX 3.25",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "A",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
        grade: "A",
      },
      {
        course_code: "B",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
        grade: "B",
      },
      {
        course_code: "C",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 2,
        grade: "C+",
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: false,
      description: `{
  "gpax": 3.25,
  "total_credits": 8
}`,
    },
    implemented: true,
  },
  {
    id: "case-12b",
    name: "GPAX - W and S Grades",
    description: "A(3cr) + W(3cr) + S(3cr) — W and S not counted in GPAX, GPAX = 4.0",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "A",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
        grade: "A",
      },
      {
        course_code: "B",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
        grade: "W",
      },
      {
        course_code: "C",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
        grade: "S",
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: false,
      description: `{
  "gpax": 4,
  "total_credits": 9
}`,
    },
    implemented: true,
  },
  {
    id: "case-12c",
    name: "GPAX - Grade F",
    description: "A(3cr) + F(3cr) — F counts as 0.0 in GPAX, GPAX = 2.0",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "A",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
        grade: "A",
      },
      {
        course_code: "B",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 3,
        grade: "F",
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: false,
      description: `{
  "gpax": 2,
  "total_credits": 6
}`,
    },
    implemented: true,
  },
  {
    id: "case-12d",
    name: "GPAX - Grade F Not Counting as Prerequisite",
    description: "2301173(F) then 2301260 — F grade means prerequisite not satisfied",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301173",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 4,
        grade: "F",
      },
      {
        course_code: "2301260",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 2,
        credits: 4,
        grade: "A",
      },
    ],
    exemptions: [],
    expectedResult: {
      canGraduate: false,
      hasViolations: true,
      description: `{
  "prerequisite_violations": [
    {
      "course_code": "2301260",
      "missing_prereqs": ["2301170", "2301173", "2301172"],
      "prereqs_taken_in_wrong_term": [],
      "taken_in_wrong_term": false,
      "missing_coreqs": [],
      "coreqs_taken_in_wrong_term": []
    }
  ],
  "gpax": 2,
  "total_credits": 8
}`,
    },
    implemented: true,
  },
  {
    id: "case-13",
    name: "Manual Credits",
    description: "Adding manual credits for transfer credits to วิชาพื้นฐานวิทยาศาสตร์",
    curriculumName: "เอกเดี่ยว-ฝึกงาน-2566",
    admissionYear: 2566,
    studyPlan: [
      {
        course_code: "2301173",
        yearOfStudy: 1,
        academicYear: 2566,
        semester: 1,
        credits: 12,
      },
    ],
    exemptions: [],
    manualCredits: {
      "วิชาพื้นฐานวิทยาศาสตร์": 12,
    },
    expectedResult: {
      canGraduate: false,
      hasViolations: false,
      description: `{
  "category_results": [
    {
      "category_name": "วิชาพื้นฐานวิทยาศาสตร์",
      "earned_credits": 12,
      "required_credits": 12,
      "is_satisfied": true
    }
  ]
}`,
    },
    implemented: true,
  },
];
