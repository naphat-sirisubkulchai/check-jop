import apiClient from "./apiClient";
import { Course, Curriculum, Category } from "@/types";

/* eslint-disable @typescript-eslint/no-explicit-any */

// Transform snake_case API response to camelCase
function transformCurriculum(data: any): Curriculum {
  return {
    id: data.id,
    nameTH: data.name_th,
    nameEN: data.name_en,
    year: data.year,
    minTotalCredits: data.min_total_credits,
    isActive: data.is_active,
    categories: data.categories ? data.categories.map((cat: any): Category => ({
      id: cat.id,
      nameTH: cat.name_th,
      nameEN: cat.name_en,
      minCredits: cat.min_credits,
    })) : [],
    courses: data.courses ? data.courses.map((course: any): Course => ({
      code: course.code,
      nameTH: course.name_th,
      nameEN: course.name_en,
      credits: course.credits,
      categoryId: course.category_id,
      prerequisites: course.prerequisites,
      corequisites: course.corequisites,
      categories: course.categories,
      curriculum: course.curriculum,
      year: course.year,
    })) : [],
  };
}

// Convert prerequisite_groups / corequisite_groups to a human-readable string
// Logic: multiple groups = AND between groups; is_or_group=true = OR within group
function groupsToString(groups: any[]): string {
  if (!groups || groups.length === 0) return "";
  const parts = groups.map((g: any) => {
    const codes: string[] = (g.prerequisite_courses ?? [])
      .map((pc: any) => pc.prerequisite_course?.code)
      .filter(Boolean);
    if (codes.length === 0) return null;
    if (codes.length === 1) return codes[0];
    return g.is_or_group ? `(${codes.join(" OR ")})` : `(${codes.join(" AND ")})`;
  }).filter(Boolean);
  return parts.join(" AND ");
}

export const courseApi = {
  // -------- Curricula --------
  async getCurriculumByName(curriculumName: string | null): Promise<any> {
    if (!curriculumName) {
      console.warn("No curriculum provided");
      return null;
    }

    try {
      console.log("Fetching curriculum:", curriculumName);
      const res = await apiClient.get(`/curricula/name/${curriculumName}`);
      if (res.status !== 200) {
        throw new Error(`Failed to fetch curriculum: ${res.statusText}`);
      }

      // TODO: ยกเลิก transform ไปก่อน
      // const curriculum = transformCurriculum(res.data);
      const curriculum = res.data;
      if (curriculum) {
        curriculum.categories?.forEach((cat: any) => {
          // Filter courses by category and remove duplicates by course code
          const coursesForCategory = curriculum
            .courses!.filter((course: any) => (course.categoryId ?? course.category_id) === cat.id);

          // Deduplicate by course code using Map
          const uniqueCoursesMap = new Map(
            coursesForCategory.map((course: any) => [course.code, course])
          );

          cat.courses = Array.from(uniqueCoursesMap.values())
            .sort((a: any, b: any) => a.code.localeCompare(b?.code));
        });
      }
      console.log("fetching curriculum: ",curriculum);
      
      return curriculum;
    } catch (error) {
      console.error("Error fetching and transforming curriculum:", error);
      return null;
    }
  },

  async getCurriculumById(id: string): Promise<Curriculum | null> {
    try {
      const [currRes, coursesRes] = await Promise.all([
        apiClient.get(`/curricula/${id}`),
        apiClient.get(`/courses/`),
      ]);
      if (currRes.status !== 200) throw new Error(currRes.statusText);

      const curriculum = transformCurriculum(currRes.data);

      // Filter courses for this curriculum, deduplicate by code (same course exists per year), convert groups → string
      const allCourses: any[] = Array.isArray(coursesRes.data) ? coursesRes.data : [];
      const seenCodes = new Set<string>();
      const curriculumCourses: Course[] = allCourses
        .filter((c: any) => c.curriculum_id === id)
        .filter((c: any) => {
          if (seenCodes.has(c.code)) return false;
          seenCodes.add(c.code);
          return true;
        })
        .map((c: any): Course => ({
          code: c.code,
          nameTH: c.name_th,
          nameEN: c.name_en,
          credits: c.credits,
          categoryId: c.category_id,
          prerequisites: groupsToString(c.prerequisite_groups ?? []),
          corequisites: groupsToString(c.corequisite_groups ?? []),
          curriculum: curriculum.nameTH,
          year: c.year,
        }));

      // Attach courses to categories
      if (curriculum.categories) {
        curriculum.categories = curriculum.categories.map((cat: any) => {
          const courses = curriculumCourses
            .filter((c) => c.categoryId === cat.id)
            .sort((a, b) => a.code.localeCompare(b.code));
          return { ...cat, courses };
        });
      }

      curriculum.courses = curriculumCourses;
      return curriculum;
    } catch {
      return null;
    }
  },

  async getAllCurriculum(): Promise<Curriculum[]> {
    const res = await apiClient.get("/curricula/");
    if (res.status !== 200) throw new Error(res.statusText);
    return Array.isArray(res.data)
      ? res.data.map((c: any) => transformCurriculum(c))
      : [];
  },
  
  async getAllCurriculaWithout(): Promise<Curriculum[]> {
    const res = await apiClient.get("/curricula/allwithout");
    if (res.status !== 200) throw new Error(res.statusText);
    return Array.isArray(res.data)
      ? res.data.map((c: any) => transformCurriculum(c))
      : [];
  },

  async createCurriculum(payload: {
    fileCurriculum: File;
    fileCategory: File;
    courseFiles: Array<{ file: File; year: number }>
  }): Promise<Curriculum> {
    // Helper to upload a file with a given key and endpoint
    const uploadFile = async (file: File, key: string, endpoint: string) => {
      const formData = new FormData();
      formData.append('file', file);
      const res = await apiClient.post(endpoint, formData);
      if (![200, 201].includes(res.status)) {
        throw new Error(`Failed to upload ${key}: ${res.statusText}`);
      }
      return res;
    };

    // Helper to upload course file with year parameter (Version 4)
    const uploadCourseFileWithYear = async (file: File, year: number) => {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('year', year.toString());
      const res = await apiClient.post("/import/course-csv-with-year", formData);
      if (![200, 201].includes(res.status)) {
        throw new Error(`Failed to upload course file for year ${year}: ${res.statusText}`);
      }
      return res;
    };

    // Upload curriculum file
    const resCurriculum = await uploadFile(payload.fileCurriculum, "file", "/import/curriculum-csv");

    // Upload category file
    await uploadFile(payload.fileCategory, "fileCategory", "/import/category-csv");

    // Upload course files for each year
    for (const courseFile of payload.courseFiles) {
      await uploadCourseFileWithYear(courseFile.file, courseFile.year);
    }

    // Return the curriculum from the first response
    return transformCurriculum(resCurriculum.data);
  },

  async deleteCurriculum(id: string): Promise<void> {
    const res = await apiClient.delete(`/curricula/${id}`);
    if (res.status !== 200 && res.status !== 204)
      throw new Error(res.statusText);
  },

  // -------- Courses --------
  async getAllCourse(): Promise<Course[]> {
    const res = await apiClient.get("/courses");
    if (res.status !== 200) throw new Error(res.statusText);
    return Array.isArray(res.data)
      ? res.data.map((r: any) => ({
          code: r.code,
          nameEN: r.name_en ?? r.nameEN ?? r.nameEn,
          nameTH: r.name_th ?? r.nameTH ?? r.nameTh,
          credits: Number(r.credits),
          categoryId: r.category_id ?? r.categoryId ?? undefined,
        }))
      : [];
  },

  async createCourse(payload: Course): Promise<Course> {
    const body = {
      code: payload.code,
      name_en: payload.nameEN,
      name_th: payload.nameTH,
      credits: payload.credits,
      category_id: payload.categoryId,
    };
    const res = await apiClient.post("/courses", body);
    if (res.status !== 201 && res.status !== 200)
      throw new Error(res.statusText);
    const r = res.data;
    return {
      code: r.code,
      nameEN: r.name_en ?? r.nameEN,
      nameTH: r.name_th ?? r.nameTH,
      credits: Number(r.credits),
      categoryId: r.category_id ?? r.categoryId,
    };
  },

  async updateCourse(
    code: string,
    payload: Partial<Omit<Course, "code">>
  ): Promise<Course> {
    const body: any = {};
    if (payload.nameEN !== undefined) body.name_en = payload.nameEN;
    if (payload.nameTH !== undefined) body.name_th = payload.nameTH;
    if (payload.credits !== undefined) body.credits = payload.credits;
    if (payload.categoryId !== undefined) body.category_id = payload.categoryId;
    const res = await apiClient.put(`/courses/${code}`, body);
    if (res.status !== 200) throw new Error(res.statusText);
    const r = res.data;
    return {
      code: r.code,
      nameEN: r.name_en ?? r.nameEN,
      nameTH: r.name_th ?? r.nameTH,
      credits: Number(r.credits),
      categoryId: r.category_id ?? r.categoryId,
    };
  },

  async deleteCourse(code: string): Promise<void> {
    const res = await apiClient.delete(`/courses/${code}`);
    if (res.status !== 200 && res.status !== 204)
      throw new Error(res.statusText);
  },
};
