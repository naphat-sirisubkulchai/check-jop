"use client";

import React, { useState, useEffect, useMemo } from "react";
import "@xyflow/react/dist/style.css";
import { useCourseGraph } from "@/graph/hooks/useCourseGraph";
import CourseGraph from "@/graph/components/CourseGraph";
import { Course } from "@/types";

const GraphManager = ({
  courses,
  curriculums,
}: {
  courses: Course[];
  curriculums: any[];
}) => {
  const [selectedCurriculum, setSelectedCurriculum] = useState<string | null>(
    null
  );
  const [filteredCourses, setFilteredCourses] = useState<Course[]>(courses);

  useEffect(() => {
    if (selectedCurriculum) {
      const filteredCourses = courses.filter((course) => {
        if (!course.curriculum) {
          return false;
        }
        const splitCurriculum = course.curriculum
          .split(",")
          .map((c: string) => c.trim());
        return splitCurriculum.includes(selectedCurriculum);
      });
      setFilteredCourses(filteredCourses);
    } else {
      setFilteredCourses(courses);
    }
  }, [selectedCurriculum, courses]);

  const { loading, error } = useCourseGraph(courses);

  const handleCurriculumChange = (
    event: React.ChangeEvent<HTMLSelectElement>
  ) => {
    setSelectedCurriculum(event.target.value || null);
  };

  const renderContent = useMemo(() => {
    if (loading) {
      return (
        <div className="absolute inset-0 flex items-center justify-center bg-white bg-opacity-90 z-10">
          <div className="flex flex-col items-center space-y-4">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <div className="text-lg text-gray-600">Loading course graph...</div>
          </div>
        </div>
      );
    }

    if (!selectedCurriculum) {
      return (
        <div className="absolute inset-0 flex items-center justify-center">
          <div className="text-center">
            <div className="text-gray-400 text-6xl mb-4">🎯</div>
            <div className="text-xl text-gray-500 mb-2">
              Select a curriculum
            </div>
            <div className="text-gray-400">
              Choose a curriculum from the dropdown to view the course
              dependency graph
            </div>
          </div>
        </div>
      );
    }

    if (filteredCourses.length === 0) {
      return (
        <div className="absolute inset-0 flex items-center justify-center">
          <div className="text-center">
            <div className="text-gray-400 text-6xl mb-4">📚</div>
            <div className="text-xl text-gray-500 mb-2">No courses found</div>
            <div className="text-gray-400">
              This curriculum doesn't have any courses with dependency data
            </div>
          </div>
        </div>
      );
    }

    if (error) {
      return (
        <div className="absolute top-4 left-4 right-4 z-20 p-4 bg-red-50 border border-red-200 rounded-lg shadow-sm">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <svg
                className="h-5 w-5 text-red-400"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                  clipRule="evenodd"
                />
              </svg>
            </div>
            <div className="ml-3">
              <div className="text-sm text-red-800">
                <strong>Error:</strong> {error}
              </div>
            </div>
          </div>
        </div>
      );
    }

    return <CourseGraph courses={filteredCourses} />;
  }, [loading, error, filteredCourses]);

  return (
    <div className="w-full h-screen flex flex-col bg-gray-50">
      {/* Header with curriculum filter */}
      <div className="p-6 bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto">
          <div className="flex items-center justify-between mb-4">
            <h1 className="text-3xl font-bold text-gray-900">
              Course Dependency Graph
            </h1>
            <div className="flex items-center space-x-4">
              <label
                htmlFor="curriculum-select"
                className="text-sm font-medium text-gray-700"
              >
                Select Curriculum:
              </label>
              <select
                id="curriculum-select"
                value={selectedCurriculum || ""}
                onChange={handleCurriculumChange}
                className="px-4 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white shadow-sm min-w-[300px]"
              >
                <option value="">Choose a curriculum</option>
                {curriculums.map((curriculum) => (
                  <option key={curriculum.nameTH} value={curriculum.nameTH}>
                    {curriculum.nameTH} ({curriculum.year})
                  </option>
                ))}
              </select>
            </div>
          </div>

          {/* Legend and Stats */}
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-8">
              <div className="flex items-center space-x-2">
                <div className="w-6 h-0.5 bg-red-500"></div>
                <span className="text-sm text-gray-600">Prerequisites</span>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-6 h-0.5 bg-green-500 border-dashed border-2 border-green-500"></div>
                <span className="text-sm text-gray-600">Corequisites</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Graph container */}
      <div className="flex-1 relative">{renderContent}</div>
    </div>
  );
};

export default GraphManager;
