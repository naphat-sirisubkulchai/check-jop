import apiClient from "./apiClient";

/* eslint-disable @typescript-eslint/no-explicit-any */

export const gradeService = {
  async checkGraduation(payload: any): Promise<any> {
    console.log("payload",payload);

    const res = await apiClient.post("/graduation/check/name", payload);
    return res.data;
  },

  // Check if a course allows C.F. exemption
  async checkCFOption(
    courseCode: string,
    curriculumId: string,
    year: number
  ): Promise<{ has_cf_option: boolean; message: string; course_code: string }> {
    const res = await apiClient.get(
      `/courses/code/${courseCode}/cf-option`,
      {
        params: {
          curriculum_id: curriculumId,
          year: year,
        },
      }
    );
    if (res.status !== 200) {
      throw new Error(res.statusText);
    }
    return {
      has_cf_option: res.data.has_cf_option,
      message: res.data.message,
      course_code: res.data.course_code,
    };
  },
};

