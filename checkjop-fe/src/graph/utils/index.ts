import dagre from "dagre";
import { Node, Edge } from "@xyflow/react";

function generateGraphLayout(
  nodes: Node[],
  edges: Edge[],
  options: { direction?: "TB" | "BT" | "LR" | "RL" } = {}
): { nodes: Node[]; edges: Edge[] } {
  const nodeWidth = 172;
  const nodeHeight = 36;
  const { direction = "TB" } = options;

  const dagreGraph = new dagre.graphlib.Graph();
  dagreGraph.setDefaultEdgeLabel(() => ({}));

  // Improved graph settings for better layout
  dagreGraph.setGraph({
    rankdir: direction,
    nodesep: 50, // Horizontal spacing between nodes
    ranksep: 120, // Vertical spacing between ranks
    marginx: 20,
    marginy: 20,
    acyclicer: "greedy", // Handle cycles better
    ranker: "tight-tree", // Better ranking algorithm
  });

  // set nodes
  nodes.forEach((node) => {
    dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
  });

  // set edges
  edges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });

  dagre.layout(dagreGraph);

  const layoutedNodes = nodes.map((node) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    return {
      ...node,
      position: {
        x: nodeWithPosition.x - nodeWidth / 2,
        y: nodeWithPosition.y - nodeHeight / 2,
      },
      // important: prevent React Flow from re-calculating position
      draggable: true,
    };
  });

  return { nodes: layoutedNodes, edges };
}

function parseRelations(rel: string): string[] {
  if (!rel) return [];

  // Convert to string if it's a number
  const relString = rel.toString().trim();
  if (!relString) return [];

  // Remove outer parentheses if they wrap the entire string
  const cleaned = relString.replace(/^\(|\)$/g, "");

  // Split by OR first, then by comma within each OR group
  const orGroups = cleaned.split(/\s+OR\s+/i);
  const result: string[] = [];

  orGroups.forEach((group) => {
    // Remove any remaining parentheses and split by comma
    const courses = group
      .replace(/[()]/g, "")
      .split(",")
      .map((s) => s.trim())
      .filter(Boolean);

    result.push(...courses);
  });

  return [...new Set(result)]; // Remove duplicates
}

export { generateGraphLayout, parseRelations };
