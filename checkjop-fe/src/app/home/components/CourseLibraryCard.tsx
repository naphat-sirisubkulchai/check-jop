import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import { GripVertical } from "lucide-react";

interface CourseLibraryCardProps {
  course: {
    code: string;
    name_en: string;
    name_th?: string;
    credits: number;
  };
  isInPlan: boolean;
}

export default function CourseLibraryCard({ course, isInPlan }: CourseLibraryCardProps) {
  return (
    <Card
      className={`p-2 transition-all border ${
        isInPlan
          ? "bg-gray-50 border-gray-300 opacity-60 cursor-not-allowed"
          : "bg-surface border-gray-200 hover:shadow-md hover:border-chula-active hover:scale-[1.01] cursor-grab active:cursor-grabbing"
      }`}
      draggable={!isInPlan}
      onDragStart={(e) => {
        if (!isInPlan) {
          e.dataTransfer.setData("text/plain", `LIB:${course.code}`);
          e.currentTarget.classList.add("opacity-50");
        } else {
          e.preventDefault();
        }
      }}
      onDragEnd={(e) => {
        if (!isInPlan) {
          e.currentTarget.classList.remove("opacity-50");
        }
      }}
    >
      <div className="flex items-center gap-1 px-1">
        {!isInPlan && (
          <GripVertical className="h-5 w-5 text-gray-300 flex-shrink-0 -ml-1" />
        )}

        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1.5 flex-wrap">
            <span className={`text-sm font-semibold ${
              isInPlan ? "text-gray-400" : "text-gray-900"
            }`}>
              {course.code}
            </span>
          </div>
          <p className={`text-sm line-clamp-1 leading-snug ${
            isInPlan ? "text-gray-500" : "text-gray-900"
          }`}>
            {course.name_en}
          </p>
        </div>

        {/* Credits */}
        <Badge className="rounded-full bg-chula-soft text-chula-active">{course.credits} cr.</Badge>
      </div>
    </Card>
  );
}
