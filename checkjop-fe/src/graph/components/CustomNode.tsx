import { Handle, Position } from "@xyflow/react";
import { memo } from "react";

const CustomNode = ({ data }: any) => {
  const { code, nameEN, color } = data;

  return (
    <div 
      className="px-4 py-3 shadow-lg rounded-lg bg-white border-2 hover:shadow-xl transition-all duration-200 w-40"
      style={{ 
        borderColor: color,
        borderWidth: 2,
      }}
      title={`${code} - ${nameEN}`}
      
    >
      {/* Course Code - Most important for identification */}
      <div className="text-sm font-bold text-gray-900 mb-1">
        {code}
      </div>
      
      {/* Course Name (English) */}
      <div className="text-xs text-gray-700 font-medium leading-tight mb-1 break-words truncate">
        {nameEN}
      </div>

      {/* Connection handles */}
      <Handle
        type="target"
        position={Position.Top}
      />
      <Handle
        type="source"
        position={Position.Bottom}
      />
      
    </div>
  );
};

export default memo(CustomNode);