import { useState, useCallback, useEffect } from 'react';
import { Node, Edge, Position } from '@xyflow/react';
import { parseRelations } from '../utils/index';
import { Course } from '@/types/index';

/* eslint-disable @typescript-eslint/no-explicit-any */

interface CourseNode extends Node {
  data: {
    label: string;
    code: string;
    credits: number;
    nameEN: string;
    nameTH: string;
  };
}

interface CourseGraphData {
  nodes: CourseNode[];
  edges: Edge[];
  courses: Course[];
  loading: boolean;
  error: string | null;
}

export const useCourseGraph = (courses: any[] | null) => {
  const [graphData, setGraphData] = useState<CourseGraphData>({
    nodes: [],
    edges: [],
    courses: [],
    loading: false,
    error: null,
  });

  const createNode = useCallback((course: any, index: number): CourseNode => {
    // Calculate position in a grid layout with better spacing
    const nodesPerRow = 5;
    const nodeWidth = 220;
    const nodeHeight = 80;
    const horizontalSpacing = 280;
    const verticalSpacing = 140;
    
    const x = (index % nodesPerRow) * horizontalSpacing;
    const y = Math.floor(index / nodesPerRow) * verticalSpacing;

    const courseCode = course.code.toString();
    const courseName = course.courseNameEN || course.nameEN || course.courseName || '';
    const courseNameTH = course.courseNameTH || course.nameTH || '';
    
    return {
      id: courseCode,
      type: 'default',
      position: { x, y },
      data: {
        label: `${courseCode}\n${courseName}\n${courseNameTH}`,
        code: courseCode,
        credits: course.credit || course.credits || 0,
        nameEN: courseName,
        nameTH: courseNameTH,
      },
      sourcePosition: Position.Right,
      targetPosition: Position.Left,
      style: {
        width: nodeWidth,
        height: nodeHeight,
        fontSize: '11px',
        padding: '8px',
      },
    };
  }, []);

  const createEdges = useCallback((courses: any[]): Edge[] => {
    const edges: Edge[] = [];
    const courseCodeSet = new Set(courses.map(c => c.code.toString()));

    courses.forEach((course) => {
      const sourceId = course.code.toString();

      // Parse prerequisites (red edges)
      if (course.prerequisites) {
        const prerequisites = parseRelations(course.prerequisites.toString());
        prerequisites.forEach((prereq, index) => {
          // Only create edge if the prerequisite course exists in our dataset
          if (courseCodeSet.has(prereq)) {
            edges.push({
              id: `prereq-${prereq}-${sourceId}-${index}`,
              source: prereq,
              target: sourceId,
              type: 'smoothstep',
              style: { 
                stroke: '#ef4444', 
                strokeWidth: 2,
                strokeDasharray: '0',
              },
              markerEnd: {
                type: 'arrowclosed',
                color: '#ef4444',
              },
              label: 'prerequisite',
              labelStyle: { 
                fontSize: 10, 
                fill: '#ef4444',
                fontWeight: 'bold',
              },
              labelBgStyle: { fill: '#ffffff', fillOpacity: 0.8 },
            });
          }
        });
      }

      // Parse corequisites (green animated edges)
      if (course.corequisites) {
        const corequisites = parseRelations(course.corequisites.toString());
        corequisites.forEach((coreq, index) => {
          // Only create edge if the corequisite course exists in our dataset
          if (courseCodeSet.has(coreq)) {
            edges.push({
              id: `coreq-${coreq}-${sourceId}-${index}`,
              source: coreq,
              target: sourceId,
              type: 'straight',
              style: { 
                stroke: '#22c55e', 
                strokeWidth: 2,
                strokeDasharray: '5,5',
              },
              animated: true,
              markerEnd: {
                type: 'arrowclosed',
                color: '#22c55e',
              },
              label: 'corequisite',
              labelStyle: { 
                fontSize: 10, 
                fill: '#22c55e',
                fontWeight: 'bold',
              },
              labelBgStyle: { fill: '#ffffff', fillOpacity: 0.8 },
            });
          }
        });
      }
    });

    return edges;
  }, []);

  const loadCourseGraph = useCallback(async () => {
    if (!courses) {
      setGraphData(prev => ({ ...prev, nodes: [], edges: [], courses: [] }));
      return;
    }
    setGraphData(prev => ({ ...prev, loading: true, error: null }));

    try {
      const nodes = courses.map((course, index) => createNode(course, index));
      const edges = createEdges(courses);

      setGraphData({
        nodes,
        edges,
        courses: courses,
        loading: false,
        error: null,
      });
    } catch (error) {
      console.error('Error loading course graph:', error);
    }
  }, [courses, createNode, createEdges]);

  useEffect(() => {
    loadCourseGraph();
  }, [loadCourseGraph]);

  return {
    ...graphData,
    refresh: loadCourseGraph,
  };
};