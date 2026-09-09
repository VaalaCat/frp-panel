import {
  forceSimulation,
  forceLink,
  forceManyBody,
  forceX,
  forceY,
  forceCollide,
  SimulationNodeDatum,
  SimulationLinkDatum,
  Simulation
} from 'd3-force';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import {
  useReactFlow,
  useNodesInitialized,
  Node,
  Edge,
  ReactFlowState,
  useStore
} from '@xyflow/react';
import { NODE_WIDTH, NODE_HEIGHT } from './layout';

// D3 Node extension
interface D3Node extends SimulationNodeDatum {
  id: string;
  width?: number;
  height?: number;
  type?: string;
  fx?: number | null;
  fy?: number | null;
  x?: number;
  y?: number;
}

interface D3Link extends SimulationLinkDatum<D3Node> {
  source: string | D3Node;
  target: string | D3Node;
  id: string;
}

const simulation = forceSimulation<D3Node, D3Link>()
  .force('charge', forceManyBody().strength(-3000))
  .force('x', forceX().x(0).strength(0.02))
  .force('y', forceY().y(0).strength(0.02))
  .force('collide', forceCollide().radius((d: any) => Math.max(d.width || NODE_WIDTH, d.height || NODE_HEIGHT) * 0.8))
  .alphaTarget(0.05)
  .stop();

export const useForceLayout = () => {
  const { getNodes, setNodes, getEdges, fitView } = useReactFlow();
  const initialized = useNodesInitialized();

  // Ref to track dragging node
  const draggingNodeRef = useRef<Node | null>(null);

  // Drag events handlers
  const dragEvents = useMemo(
    () => ({
      start: (_event: React.MouseEvent, node: Node) => {
        draggingNodeRef.current = node;
        simulation.alphaTarget(0.3).restart();
      },
      drag: (_event: React.MouseEvent, node: Node) => {
        draggingNodeRef.current = node;
      },
      stop: () => {
        draggingNodeRef.current = null;
        simulation.alphaTarget(0);
      },
    }),
    []
  );

  useEffect(() => {
    let nodes = getNodes().map((node) => ({
      ...node,
      x: node.position.x,
      y: node.position.y,
      width: node.width || NODE_WIDTH,
      height: node.height || NODE_HEIGHT,
    })) as D3Node[];

    let edges = getEdges().map((edge) => ({ ...edge })) as unknown as D3Link[];
    let running = false;

    // Only run if initialized and we have nodes
    if (!initialized || nodes.length === 0) return;

    // Filter only WG nodes for simulation
    // TerminalNode and other custom nodes should NOT participate in force layout
    const simulationNodes = nodes.filter(n => n.type === 'wg');
    const simulationEdges = edges.filter(edge => {
      const sourceId = typeof edge.source === 'string' ? edge.source : (edge.source as any).id;
      const targetId = typeof edge.target === 'string' ? edge.target : (edge.target as any).id;
      const sourceExists = simulationNodes.some(n => n.id === sourceId);
      const targetExists = simulationNodes.some(n => n.id === targetId);
      return sourceExists && targetExists;
    });

    simulation.nodes(simulationNodes).force(
      'link',
      forceLink<D3Node, D3Link>(simulationEdges)
        .id((d) => d.id)
        .strength(0.5)
        .distance(500)
    );

    // Pre-tick the simulation to get a better initial layout
    // This prevents the "explosion from center" effect and speeds up initial stabilization
    for (let i = 0; i < 2000; i++) {
      simulation.tick();
    }

    // Sync initial positions after pre-tick
    setNodes((prevNodes) =>
      prevNodes.map((node) => {
        if (node.type !== 'wg') return node;
        const d3Node = simulationNodes.find((n) => n.id === node.id);
        if (!d3Node) return node;

        return {
          ...node,
          position: {
            x: d3Node.x ?? node.position.x,
            y: d3Node.y ?? node.position.y
          },
        };
      })
    );

    // Fit view after initial layout
    // Use setTimeout to ensure state is updated
    setTimeout(() => {
      fitView({ duration: 0, padding: 0.2 });
    }, 0);

    const tick = () => {
      // Update fx/fy for the dragging node
      nodes.forEach((node, i) => {
        if (node.type !== 'wg') return;

        const isDragging = draggingNodeRef.current?.id === node.id;
        if (isDragging && draggingNodeRef.current) {
          node.fx = draggingNodeRef.current.position.x;
          node.fy = draggingNodeRef.current.position.y;
        } else {
          node.fx = null;
          node.fy = null;
        }
      });

      simulation.tick();

      // Update React Flow nodes with new positions from simulation
      setNodes((prevNodes) =>
        prevNodes.map((node) => {
          if (node.type !== 'wg') return node;

          const d3Node = nodes.find((n) => n.id === node.id);
          if (!d3Node) return node;

          if (draggingNodeRef.current?.id === node.id) return node;

          const newX = d3Node.x ?? node.position.x;
          const newY = d3Node.y ?? node.position.y;

          return {
            ...node,
            position: { x: newX, y: newY },
          };
        })
      );

      window.requestAnimationFrame(() => {
        if (running) tick();
      });
    };

    // Start the simulation
    running = true;
    simulation.alpha(1).restart();
    window.requestAnimationFrame(tick);

    return () => {
      running = false;
      simulation.stop();
    };
  }, [initialized, getNodes, getEdges, setNodes, fitView]);

  return { dragEvents };
};
