"use client";

import { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { FileText, TableProperties } from "lucide-react";
export type PrintFormat = "transcript" | "summary";

interface PrintDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: (format: PrintFormat) => void;
}

export function PrintDialog({ open, onClose, onConfirm }: PrintDialogProps) {
  const [selected, setSelected] = useState<PrintFormat>("transcript");

  const options: { value: PrintFormat; label: string; desc: string; icon: React.ReactNode }[] = [
    {
      value: "transcript",
      label: "ตารางบันทึกผลการศึกษา",
      desc: "3 คอลัมน์แบบฟอร์มจริง พร้อมรายวิชาทุกหมวด ปีการศึกษา และเกรด",
      icon: <TableProperties className="w-6 h-6" />,
    },
    {
      value: "summary",
      label: "สรุปผลการตรวจสอบ",
      desc: "แสดงผลการตรวจสอบ ความคืบหน้าแต่ละหมวด และรายการวิชาที่เรียน",
      icon: <FileText className="w-6 h-6" />,
    },
  ];

  return (
    <Dialog open={open} onOpenChange={(o) => { if (!o) onClose(); }}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>เลือกรูปแบบการพิมพ์</DialogTitle>
        </DialogHeader>
        <div className="space-y-3 py-2">
          {options.map((opt) => (
            <button
              key={opt.value}
              onClick={() => setSelected(opt.value)}
              className={`w-full text-left flex items-start gap-4 p-4 rounded-xl border-2 transition-all ${
                selected === opt.value
                  ? "border-chula-active bg-chula-soft/30"
                  : "border-gray-200 hover:border-gray-300 bg-white"
              }`}
            >
              <div className={`mt-0.5 flex-shrink-0 ${selected === opt.value ? "text-chula-active" : "text-gray-400"}`}>
                {opt.icon}
              </div>
              <div>
                <p className={`font-semibold text-sm ${selected === opt.value ? "text-chula-active" : "text-gray-700"}`}>
                  {opt.label}
                </p>
                <p className="text-xs text-gray-500 mt-0.5">{opt.desc}</p>
              </div>
            </button>
          ))}
        </div>
        <div className="flex gap-2 pt-2">
          <Button variant="outline" className="flex-1" onClick={onClose}>
            ยกเลิก
          </Button>
          <Button
            className="flex-1 bg-chula-active hover:bg-chula text-white"
            onClick={() => onConfirm(selected)}
          >
            พิมพ์
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
