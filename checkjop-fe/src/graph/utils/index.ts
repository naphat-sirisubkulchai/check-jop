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

  const relString = rel.toString().trim();
  if (!relString) return [];

  // Strip all parentheses, then split by OR and AND and comma
  const cleaned = relString.replace(/[()]/g, "");
  const codes = cleaned
    .split(/\s+(?:OR|AND)\s+|,/i)
    .map((s) => s.trim())
    .filter(Boolean);

  return [...new Set(codes)];
}

export { generateGraphLayout, parseRelations };
