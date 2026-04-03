/**
 * CourseGraph Component
 * 
 * Workflow:
 * 1. `buildCourseGraph`:
 *    - สร้างกราฟเบื้องต้น (pure graph) โดยสร้าง nodes และ edges จากข้อมูล courses
 *    - ตำแหน่งเริ่มต้นของ nodes (`x`, `y`) จะถูกตั้งค่าเป็นค่าเริ่มต้น เช่น `{ x: 0, y: 0 }`
 * 
 * 2. `graphData`:
 *    - เก็บข้อมูล nodes และ edges ที่สร้างจาก `buildCourseGraph` ไว้ในตัวแปร memoized (`useMemo`)
 *    - ป้องกันการคำนวณซ้ำเมื่อ `courses` ไม่เปลี่ยนแปลง
 * 
 * 3. `generateGraphLayout`:
 *    - ใช้ข้อมูลจาก `graphData` เพื่อจัดตำแหน่ง (`x`, `y`) ของ nodes และ edges ให้สวยงาม
 *    - ใช้ Dagre layout algorithm ในการคำนวณ
 * 
 * 4. State Management:
 *    - เก็บ nodes และ edges ที่ได้จาก `generateGraphLayout` ไว้ใน `nodesState` และ `edgesState` ผ่าน `useState`
 *    - ใช้ state เหล่านี้ในการ render กราฟ
 * 
 * 5. Hover Effects:
 *    - ใช้ `hoveredNode` เพื่อจัดการ hover state ของ nodes
 *    - `styledNodes` และ `styledEdges` จะถูกสร้างใหม่เมื่อ `nodesState`, `hoveredNode`, หรือ `getNodeStyle` เปลี่ยนแปลง
 *    - ปรับแต่งสไตล์ (เช่น สี, border, opacity) ของ nodes และ edges ตามสถานะของ `hoveredNode`
 * 
 * Summary:
 * - กราฟถูกสร้างและจัด layout ใน 2 ขั้นตอน: `buildCourseGraph` และ `generateGraphLayout`
 * - ใช้ state และ memoization เพื่อเพิ่มประสิทธิภาพ
 * - รองรับ interactive hover effects เพื่อแสดงความสัมพันธ์ระหว่าง
**/

"use client";

import React, { useState, useEffect, useMemo, useCallback } from "react";
import {
  ReactFlow,
  Background,
  Controls,
  BackgroundVariant,
  Node,
  Edge,
  MarkerType,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import { Course } from "@/types/index";
import CustomNode from "@/graph/components/CustomNode";
import { parseRelations, generateGraphLayout } from "../utils/index";

// Constants
// const _CATEGORY_COLORS = {
//   "วิชาแกน": "#3B82F6",
//   "วิชาเฉพาะด้าน": "#10B981",
//   "วิชาพื้นฐานวิทยาศาสตร์": "#F59E0B",
//   "กลุ่มวิชาภาษา": "#8B5CF6",
//   "วิชาศึกษาทั่วไป": "#EF4444",
//   "วิชาเลือกเสรี": "#6B7280",
//   default: "#64748B",
// } as const;

const EDGE_COLORS = {
  prereq: ['#64748B', '#374151', '#111827', '#6B7280', '#4B5563'],
  coreq: ['#F59E0B', '#D97706', '#92400E'],
} as const;

const HOVER_STYLES = {
  active: { border: "3px solid #ff5722", bg: "#ffe0b2", z: 100 },
  coreq: { border: "3px solid #39883dff", bg: "#dbf3dcff", z: 50 },
  connected: { border: "3px solid #2196f3", bg: "#e3f2fd", z: 50 },
  default: { border: "2px solid transparent", bg: "#fff", z: 0 },
} as const;

const nodeTypes = { custom: CustomNode };

// Types
type NodeStyle = 'active' | 'coreq' | 'connected' | 'default';

//-----สำหรับสร้างกราฟ-----
const topologicalSort = (courses: Course[], courseMap: Map<string, Course>): Course[] => {
  const inDegree = new Map<string, number>();
  const adjList = new Map<string, string[]>();
  
  // Initialize
  courses.forEach(course => {
    inDegree.set(course.code, 0);
    adjList.set(course.code, []);
  });

  // Build graph
  courses.forEach(course => {
    [course.prerequisites, course.corequisites].forEach((relations, isCoreq) => {
      if (!relations) return;
      parseRelations(relations).forEach(relCode => {
        if (!courseMap.has(relCode)) return;
        adjList.get(relCode)!.push(course.code);
        if (!isCoreq) inDegree.set(course.code, (inDegree.get(course.code) || 0) + 1);
      });
    });
  });

  // Topological sort
  const queue = courses.filter(c => (inDegree.get(c.code) || 0) === 0).map(c => c.code);
  const sorted: Course[] = [];
  const visited = new Set<string>();

  while (queue.length > 0) {
    const currentCode = queue.shift()!;
    if (visited.has(currentCode)) continue;
    
    visited.add(currentCode);
    const course = courseMap.get(currentCode);
    if (course) sorted.push(course);

    adjList.get(currentCode)?.forEach(depCode => {
      if (visited.has(depCode)) return;
      const newInDegree = (inDegree.get(depCode) || 0) - 1;
      inDegree.set(depCode, newInDegree);
      if (newInDegree === 0) queue.push(depCode);
    });
  }

  // Add remaining courses
  courses.forEach(course => {
    if (!visited.has(course.code)) sorted.push(course);
  });

  return sorted;
};

// Create nodes with ข้อมูลที่จำเป็น
const createNodes = (courses: Course[]): Node[] =>
  courses.map(course => ({
    id: course.code,
    type: 'custom',
    position: { x: 0, y: 0 },
    data: { ...course }
  }));

const createEdges = (courses: Course[], courseMap: Map<string, Course>): Edge[] => {
  const edges: Edge[] = [];
  
  courses.forEach(course => {
    // Prerequisites
    if (course.prerequisites) {
      parseRelations(course.prerequisites).forEach((prereqCode, index) => {
        if (!courseMap.has(prereqCode)) return;
        const color = EDGE_COLORS.prereq[index % EDGE_COLORS.prereq.length];
        edges.push(createEdge(`prereq-${prereqCode}-${course.code}`, prereqCode, course.code, color));
      });
    }
    
    // Corequisites
    if (course.corequisites) {
      parseRelations(course.corequisites).forEach((coreqCode, index) => {
        if (!courseMap.has(coreqCode)) return;
        const color = EDGE_COLORS.coreq[index % EDGE_COLORS.coreq.length];
        edges.push(createEdge(`coreq-${coreqCode}-${course.code}`, coreqCode, course.code, color, true));
      });
    }
  });
  
  return edges;
};

const createEdge = (id: string, source: string, target: string, color: string, isDashed = false): Edge => ({
  id,
  source,
  target,
  type: 'smoothstep',
  animated: false,
  style: {
    stroke: color,
    strokeWidth: 2.5,
    strokeOpacity: isDashed ? 0.9 : 0.8,
    ...(isDashed && { strokeDasharray: '8,4' })
  },
  markerEnd: {
    type: MarkerType.ArrowClosed,
    width: 18,
    height: 18,
    color
  }
});

const CourseGraph = ({ courses }: { courses: Course[] }) => {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);
  const [isLayouting, setIsLayouting] = useState(false);
  const [hoveredNode, setHoveredNode] = useState<string | null>(null);

  const buildCourseGraph = useCallback((courses: Course[]) => {
    const courseMap = new Map(courses.map(c => [c.code, c]));
    const sortedCourses = topologicalSort(courses, courseMap);
    const nodes = createNodes(sortedCourses);
    const edges = createEdges(sortedCourses, courseMap);
    return { nodes, edges };
  }, []);

  // Graph data with optimized dependency sorting
  const graphData = useMemo(() => buildCourseGraph(courses), [courses, buildCourseGraph]);

  // Layout effect
  useEffect(() => {
    if (graphData.nodes.length === 0) return;
    setIsLayouting(true);
    const layouted = generateGraphLayout(graphData.nodes, graphData.edges);
    setNodes(layouted.nodes);
    setEdges(layouted.edges);
    setIsLayouting(false);
  }, [graphData]);
  
  // ----- For hover effects -----
  const handleNodeHover = useCallback((_event: React.MouseEvent, node: Node) => setHoveredNode(node.id), []);
  const handleNodeLeave = useCallback(() => setHoveredNode(null), []);
  
  // Node connection utilities
  const getNodeStyle = useCallback((node: Node): NodeStyle => {
    if (!hoveredNode) return 'default';
    if (node.id === hoveredNode) return 'active';
    
    const isCoreqConnected = edges.some(e => 
      ((e.source === hoveredNode && e.target === node.id) || (e.target === hoveredNode && e.source === node.id)) 
      && e.style?.strokeDasharray
    );
    if (isCoreqConnected) return 'coreq';
    
    const isConnected = edges.some(e => 
      (e.source === hoveredNode && e.target === node.id) || (e.target === hoveredNode && e.source === node.id)
    );
    return isConnected ? 'connected' : 'default';
  }, [hoveredNode, edges]);

  // Styled nodes and edges
  const styledNodes = useMemo(() => {
    if (!hoveredNode) return nodes;
    return nodes.map(node => {
      const styleType = getNodeStyle(node);
      const styles = HOVER_STYLES[styleType];
      
      return {
        ...node,
        style: {
          ...node.style,
          border: styles.border,
          backgroundColor: styles.bg,
          zIndex: styles.z,
          opacity: hoveredNode && styleType === 'default' ? 0.3 : 1,
          transition: "all 0.2s ease-in-out",
        },
      };
    });
  }, [nodes, hoveredNode, getNodeStyle]);

  const styledEdges = useMemo(() => {
    if (!hoveredNode) return edges;
    return edges.map(edge => {
      const isConnected = edge.source === hoveredNode || edge.target === hoveredNode;
      const currentWidth = typeof edge.style?.strokeWidth === 'number' ? edge.style.strokeWidth : 2.5;
      const currentMarker = edge.markerEnd as any;
      
      return {
        ...edge,
        style: {
          ...edge.style,
          strokeOpacity: isConnected ? 1 : 0.1,
          strokeWidth: isConnected ? currentWidth + 0.5 : currentWidth,
          zIndex: isConnected ? 10 : 0,
        },
        markerEnd: currentMarker ? {
          ...currentMarker,
          color: isConnected ? currentMarker.color : '#D1D5DB',
        } : undefined
      };
    });
  }, [edges, hoveredNode]);

  return (
    <div className="w-full h-full relative">
      {isLayouting && (
        <div className="absolute top-4 left-4 z-10 bg-white px-3 py-2 rounded shadow-lg">
          <div className="flex items-center space-x-2">
            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
            <span className="text-sm text-gray-600">Calculating layout...</span>
          </div>
        </div>
      )}
      
      {/* Legend */}
      {/* <div className="absolute top-4 right-4 z-10 bg-white p-4 rounded-lg shadow-lg max-w-xs">
        <h3 className="font-semibold text-sm mb-3 text-gray-800">Course Categories</h3>
        <div className="space-y-2 text-xs">
          {Object.entries(CATEGORY_COLORS).map(([category, color]) => {
            if (category === 'default') return null;
            return (
              <div key={category} className="flex items-center space-x-2">
                <div 
                  className="w-3 h-3 rounded-full"
                  style={{ backgroundColor: color }}
                ></div>
                <span className="text-gray-700">{category}</span>
              </div>
            );
          })}
        </div>
        
        <div className="mt-4 pt-3 border-t border-gray-200">
          <h4 className="font-semibold text-xs mb-2 text-gray-800">การเชื่อมต่อ</h4>
          <div className="space-y-1 text-xs text-gray-600">
            <div className="flex items-center space-x-2">
              <div className="w-4 h-0.5 bg-gray-600"></div>
              <span>วิชาที่ต้องเรียนก่อน</span>
            </div>
            <div className="flex items-center space-x-2">
              <div className="w-4 h-0.5 bg-yellow-500" style={{ backgroundImage: 'repeating-linear-gradient(90deg, transparent, transparent 2px, #F59E0B 2px, #F59E0B 4px)' }}></div>
              <span>วิชาที่ต้องเรียนพร้อมกัน</span>
            </div>
          </div>
        </div>

        <div className="mt-3 pt-3 border-t border-gray-200">
          <h4 className="font-semibold text-xs mb-2 text-gray-800">คำแนะนำ</h4>
          <div className="text-xs text-gray-600 space-y-1">
            <p>• วางเมาส์บน node เพื่อดูการเชื่อมต่อ</p>
            <p>• Node จัดกลุ่มตามความสัมพันธ์</p>
            <p>• ระดับบนคือวิชาพื้นฐาน</p>
          </div>
        </div>
      </div> */}

      <ReactFlow
        nodes={styledNodes}
        edges={styledEdges}
        nodeTypes={nodeTypes}
        onNodeMouseEnter={handleNodeHover}
        onNodeMouseLeave={handleNodeLeave}
        nodesDraggable={true}
        fitView
        defaultViewport={{ x: 0, y: 0, zoom: 0.7 }}
        minZoom={0.1}
        maxZoom={1.5}
        attributionPosition="bottom-left"
        proOptions={{ hideAttribution: true }}
      >
        <Background
          variant={BackgroundVariant.Dots}
          gap={30}
          size={1}
          color="#e5e7eb"
        />
        <Controls
          position="bottom-right"
          showInteractive={false}
          showZoom={true}
          showFitView={true}
          className="shadow-lg"
        />
      </ReactFlow>
    </div>
  );
};

export default CourseGraph;
