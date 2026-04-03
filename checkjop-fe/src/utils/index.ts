// filter วิชาที่ซ้ำกันโดยใช้รหัสวิชา
export function uniqueCoursesByCode(courses: any[]) {
  return Array.from(new Map(courses.map((c) => [c.code, c])).values());
}