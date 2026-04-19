"use client";
import { useState } from "react";
import CourseListContainer from "./components/CourseListContainer";
import StudyPlanContainer from "./components/StudyPlanContainer";
import { useAppStore } from "@/store/appStore";
import SettingsDialog from "@/components/SettingsDialog";

export default function Page() {
  const { selectedCurriculum } = useAppStore();
  const [settingsOpen, setSettingsOpen] = useState(!selectedCurriculum);

  return (
    <div className="flex flex-col h-full overflow-hidden">
      <main className="flex flex-col xl:flex-row flex-1 xl:overflow-hidden overflow-y-auto overflow-x-hidden">
        <CourseListContainer />
        <StudyPlanContainer />
      </main>

      <SettingsDialog
        open={settingsOpen}
        onOpenChange={setSettingsOpen}
      />
    </div>
  );
}
