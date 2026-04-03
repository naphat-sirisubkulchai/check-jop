"use client";

import { useRef } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Download, Upload, Settings } from "lucide-react";
import { useAppStore } from "@/store/appStore";
import { toast } from "sonner";
import { readFileAsText, parseImportFile } from "@/utils/exportImport";
import { courseApi } from "@/api/courseApi";

export default function HomeActions() {
  const router = useRouter();
  const { studyPlan, exportStudyPlanToJSON, importStudyPlanFromData, setCourses } = useAppStore();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleExport = () => {
    if (studyPlan.length === 0) {
      toast.error("Cannot export", {
        description: "No study plan to export. Please add some courses first.",
      });
      return;
    }
    exportStudyPlanToJSON();
  };

  const handleImport = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    try {
      const content = await readFileAsText(file);
      const { data, validation } = parseImportFile(content);

      if (!validation.isValid) {
        toast.error("Import failed", {
          description: validation.errors.join(", "),
        });
        return;
      }

      if (data) {
        if (data.selectedCurriculum?.nameTH) {
          const curriculumData = await courseApi.getCurriculumByName(
            data.selectedCurriculum.nameTH
          );
          if (curriculumData?.courses) {
            setCourses(curriculumData.courses);
          }
        }

        importStudyPlanFromData(data);
        toast.success("Study plan imported successfully", {
          description: `Imported ${data.studyPlan.length} courses`,
        });
      }
    } catch (error) {
      toast.error("Import failed", {
        description: error instanceof Error ? error.message : "Unknown error occurred",
      });
    }

    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  return (
    <>
      <Button
        onClick={() => router.push("/setup")}
        variant="outline"
        className="bg-white shadow-sm"
      >
        <Settings className="h-4 w-4 mr-2" />
        Settings
      </Button>
      <Button
        onClick={() => fileInputRef.current?.click()}
        variant="outline"
        className="bg-white shadow-sm"
      >
        <Upload className="h-4 w-4 mr-2" />
        Import
      </Button>
      <Button
        onClick={handleExport}
        variant="outline"
        className="bg-white shadow-sm"
      >
        <Download className="h-4 w-4 mr-2" />
        Export
      </Button>
      <input
        ref={fileInputRef}
        type="file"
        accept=".json"
        onChange={handleImport}
        className="hidden"
      />
    </>
  );
}
