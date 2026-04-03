"use client";

import { Button } from "@/components/ui/button";
import { Download, Printer } from "lucide-react";
import { useAppStore } from "@/store/appStore";
import { useAnalytics } from "@/hooks/useAnalytics";

export default function CalculateActions() {
  const { result, exportStudyPlanToJSON } = useAppStore();
  const { trackEvent } = useAnalytics();

  if (!result) return null;

  const handlePrint = () => {
    trackEvent("result_printed", {
      canGraduate: result.can_graduate,
    });
    window.print();
  };

  const handleExport = () => {
    exportStudyPlanToJSON();
  };

  return (
    <>
      <Button
        onClick={handlePrint}
        variant="outline"
        className="bg-white shadow-sm"
      >
        <Printer className="h-4 w-4 mr-2" />
        Print
      </Button>
      <Button
        onClick={handleExport}
        variant="outline"
        className="bg-white shadow-sm"
      >
        <Download className="h-4 w-4 mr-2" />
        Export
      </Button>
    </>
  );
}
