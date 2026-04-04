"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Download, Printer } from "lucide-react";
import { useAppStore } from "@/store/appStore";
import { useAnalytics } from "@/hooks/useAnalytics";
import { PrintDialog, PrintFormat } from "./PrintDialog";

export default function CalculateActions() {
  const { result, exportStudyPlanToJSON, setPrintFormat } = useAppStore();
  const { trackEvent } = useAnalytics();
  const [dialogOpen, setDialogOpen] = useState(false);

  if (!result) return null;

  const handlePrintConfirm = (format: PrintFormat) => {
    setDialogOpen(false);
    setPrintFormat(format);
    trackEvent("result_printed", { canGraduate: result.can_graduate, format });
    // Allow React to re-render with new format before printing
    setTimeout(() => window.print(), 100);
  };

  return (
    <>
      <Button
        onClick={() => setDialogOpen(true)}
        variant="outline"
        className="bg-white shadow-sm"
      >
        <Printer className="h-4 w-4 mr-2" />
        Print
      </Button>
      <Button
        onClick={() => exportStudyPlanToJSON()}
        variant="outline"
        className="bg-white shadow-sm"
      >
        <Download className="h-4 w-4 mr-2" />
        Export
      </Button>

      <PrintDialog
        open={dialogOpen}
        onClose={() => setDialogOpen(false)}
        onConfirm={handlePrintConfirm}
      />
    </>
  );
}
