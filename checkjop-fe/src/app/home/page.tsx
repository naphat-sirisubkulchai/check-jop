"use client";
import CourseListContainer from "./components/CourseListContainer";
import StudyPlanContainer from "./components/StudyPlanContainer";
import { useAppStore } from "@/store/appStore";
import { useEffect } from "react";
import { useRouter } from "next/navigation";

export default function Page() {
  const router = useRouter();
  const { selectedCurriculum } = useAppStore();

  // Redirect to setup page if no curriculum is selected
  useEffect(() => {
    if (!selectedCurriculum) {
      router.push("/setup");
    }
  }, [selectedCurriculum, router]);

  // Show nothing while redirecting
  if (!selectedCurriculum) {
    return null;
  }

  return (
    <div className="flex flex-col h-full overflow-hidden">
      {/* Main Content */}
      <main className="flex flex-1 overflow-hidden">
        {/* Left Panel - Course List */}
        <CourseListContainer />
        {/* Right Panel - Semester Grid */}
        <StudyPlanContainer />
      </main>
    </div>
  );
}