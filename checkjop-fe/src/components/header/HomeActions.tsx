"use client";

import { useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Download, Upload, Settings } from "lucide-react";
import { useAppStore } from "@/store/appStore";
import { toast } from "sonner";
import { parseImportFile, readFileAsText } from "@/utils/exportImport";
import { courseApi } from "@/api/courseApi";
import SettingsDialog from "@/components/SettingsDialog";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

const SAMPLES = [
  { label: "Single Major – Intern", file: "/samples/course-1-intern.json" },
  { label: "Single Major – CoopEd", file: "/samples/course-2-cooped.json" },
  { label: "Major/Minor – Intern", file: "/samples/course-3-minor-intern.json" },
  { label: "Major/Minor – Coop", file: "/samples/course-4-minor-coop.json" },
];

export default function HomeActions() {
  const { studyPlan, exportStudyPlanToJSON, importStudyPlanFromData, setCourses } = useAppStore();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [settingsOpen, setSettingsOpen] = useState(false);

  const handleExport = () => {
    if (studyPlan.length === 0) {
      toast.error("Cannot export", { description: "No study plan to export." });
      return;
    }
    exportStudyPlanToJSON();
  };

  const handleImport = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;
    try {
      const content = await readFileAsText(file);
      await importFromContent(content);
    } catch (error) {
      toast.error("Import failed", {
        description: error instanceof Error ? error.message : "Unknown error",
      });
    }
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const importFromContent = async (content: string) => {
    const { data, validation } = parseImportFile(content);
    if (!validation.isValid) {
      toast.error("Import failed", { description: validation.errors.join(", ") });
      return;
    }
    if (data) {
      if (data.selectedCurriculum?.nameTH) {
        const curriculumData = await courseApi.getCurriculumByName(data.selectedCurriculum.nameTH);
        if (curriculumData?.courses) setCourses(curriculumData.courses);
      }
      importStudyPlanFromData(data);
      toast.success("Imported successfully", {
        description: `Loaded ${data.studyPlan.length} courses`,
      });
    }
  };

  const handleLoadSample = async (file: string) => {
    try {
      const res = await fetch(file);
      const content = await res.text();
      await importFromContent(content);
    } catch {
      toast.error("Failed to load sample");
    }
  };

  return (
    <>
      <TooltipProvider delayDuration={300}>
        <Tooltip>
          <TooltipTrigger asChild>
            <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-gray-200 bg-white shadow-sm text-xs text-gray-600 cursor-default select-none">
              <span className="inline-block w-2.5 h-2.5 rounded-sm border-l-2 border-sci-normal bg-sci-soft/40" />
              <span>วิชาเลือก</span>
              <span className="mx-1 text-gray-300">|</span>
              <span className="inline-block w-2.5 h-2.5 rounded-sm border-l-2 border-chula-active bg-chula-soft/40" />
              <span>วิชาอื่น</span>
            </div>
          </TooltipTrigger>
          <TooltipContent side="left" className="max-w-52 text-xs">
            <p className="font-semibold mb-1">ความหมายของสี</p>
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <span className="inline-block w-3 h-3 rounded-sm border-l-2 border-sci-normal bg-sci-soft/40 flex-shrink-0" />
                <span>วิชาเลือก / วิชาบังคับเลือก</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="inline-block w-3 h-3 rounded-sm border-l-2 border-chula-active bg-chula-soft/40 flex-shrink-0" />
                <span>วิชาทั่วไป / วิชาแกน / วิชาอื่นๆ</span>
              </div>
            </div>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>

      <Button onClick={() => setSettingsOpen(true)} variant="outline" className="bg-white shadow-sm">
        <Settings className="h-4 w-4 mr-2" />
        Settings
      </Button>

      <Select
        value=""
        onValueChange={(file) => { handleLoadSample(file); }}
      >
        <SelectTrigger className="bg-white shadow-sm border-input w-36">
          <SelectValue placeholder="Load Sample" />
        </SelectTrigger>
        <SelectContent>
          {SAMPLES.map((s) => (
            <SelectItem key={s.file} value={s.file}>
              {s.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Button onClick={() => fileInputRef.current?.click()} variant="outline" className="bg-white shadow-sm">
        <Upload className="h-4 w-4 mr-2" />
        Import
      </Button>
      <Button onClick={handleExport} variant="outline" className="bg-white shadow-sm">
        <Download className="h-4 w-4 mr-2" />
        Export
      </Button>
      <input ref={fileInputRef} type="file" accept=".json" onChange={handleImport} className="hidden" />

      <SettingsDialog open={settingsOpen} onOpenChange={setSettingsOpen} />
    </>
  );
}
