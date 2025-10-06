'use client'

import React, { useCallback, useRef } from 'react'
import {
  applyNodeChanges,
  applyEdgeChanges,
  Background,
  Controls,
  MiniMap,
  ReactFlow,
  NodeChange,
  EdgeChange,
  ConnectionLineType,
  addEdge,
  ReactFlowInstance,
} from '@xyflow/react'
import type { TopologyData, TopologyNode, WGEdge, OnEdgeClick } from './types'
import CustomNode from './Node'
import CustomEdge from './Edge'
import FloatingConnectionLine from './FloatingConnectionLine'
import { nanoid } from 'nanoid'
import { TerminalNode } from '@/components/canvas'

export interface TopologyCanvasProps {
  data: TopologyData
  onEdgeClick?: OnEdgeClick
  setNodes: React.Dispatch<React.SetStateAction<TopologyNode[]>>
  setEdges: React.Dispatch<React.SetStateAction<WGEdge[]>>
  onNodeConnect?: (fromId: string, toId: string) => void
  onDeleteNode?: (nodeId: string) => void
}

export default function TopologyCanvas({
  data,
  onEdgeClick,
  setNodes,
  setEdges,
  onNodeConnect,
  onDeleteNode,
}: TopologyCanvasProps) {
  const reactFlowInstance = useRef<ReactFlowInstance | null>(null)

  const onNodesChange = useCallback(
    (changes: NodeChange[]) => setNodes((nodesSnapshot) => applyNodeChanges(changes, nodesSnapshot as any) as any),
    [setNodes],
  )

  const onEdgesChange = useCallback(
    (changes: EdgeChange[]) => setEdges((edgesSnapshot) => applyEdgeChanges(changes, edgesSnapshot as any) as any),
    [setEdges],
  )

  const onConnect = useCallback(
    (params: any) => {
      setEdges((eds) => addEdge({ ...params, type: 'wgEdge' } as any, eds as any) as any)
      onNodeConnect?.(params.source, params.target)
    },
    [setEdges, onNodeConnect],
  )

  const onInit = useCallback((instance: any) => {
    reactFlowInstance.current = instance
  }, [])

  const handleOpenTerminal = useCallback(
    (clientId: string, clientType: number, sourceNodeId?: string) => {
      if (!reactFlowInstance.current) return

      const instance = reactFlowInstance.current

      // 找到源节点的位置
      let position = { x: 0, y: 0 }

      if (sourceNodeId) {
        const sourceNode = instance.getNodes().find((n: any) => n.id === sourceNodeId)
        if (sourceNode) {
          // 在源节点右侧创建终端节点
          position = {
            x: sourceNode.position.x + (sourceNode.width || 200) + 50,
            y: sourceNode.position.y,
          }
        }
      }

      // 如果找不到源节点，使用画布中心
      if (position.x === 0 && position.y === 0) {
        const bounds = document.querySelector('.react-flow__viewport')?.getBoundingClientRect()

        const centerX = bounds ? bounds.width / 2 : 400
        const centerY = bounds ? bounds.height / 2 : 300

        position = instance.screenToFlowPosition({
          x: centerX,
          y: centerY,
        })
      }

      const newNode = {
        id: `terminal-${nanoid()}`,
        type: 'terminal',
        position,
        dragHandle: '.drag-handle', // 只允许通过 drag-handle 拖拽
        style: { width: 650, height: 500 },
        data: {
          label: `${clientId}`,
          clientId,
          clientType,
        },
      }

      // 添加新的终端节点
      setNodes((nds) => [...nds, newNode as any])
    },
    [setNodes],
  )

  const handleDeleteNode = useCallback(
    (nodeId: string) => {
      setNodes((nds) => nds.filter((n) => n.id !== nodeId))
      onDeleteNode?.(nodeId)
    },
    [setNodes, onDeleteNode],
  )

  // 创建自定义节点类型的工厂函数，注入回调
  const nodeTypes = React.useMemo(
    () => ({
      wg: (props: any) => (
        <CustomNode
          {...props}
          onOpenTerminal={(clientId: string, clientType: number) => handleOpenTerminal(clientId, clientType, props.id)}
        />
      ),
      terminal: (props: any) => <TerminalNode {...props} onDelete={handleDeleteNode} />,
    }),
    [handleOpenTerminal, handleDeleteNode],
  )

  return (
    <div className="h-[600px] md:h-[700px] lg:h-[800px] w-full border rounded-lg overflow-hidden shadow-sm">
      <ReactFlow
        nodes={data.nodes}
        edges={data.edges.map((e) => ({ ...e, sourceHandle: undefined, targetHandle: undefined }))}
        nodeTypes={nodeTypes}
        edgeTypes={{ wgEdge: CustomEdge as any }}
        defaultEdgeOptions={{ type: 'wgEdge' as any }}
        connectionLineType={ConnectionLineType.Straight}
        connectionLineComponent={FloatingConnectionLine as any}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onInit={onInit}
        fitView
        onEdgeClick={(_, edge) => onEdgeClick?.(edge.id)}
      >
        <MiniMap nodeStrokeWidth={3} zoomable pannable />
        <Controls />
        <Background gap={12} size={1} />
      </ReactFlow>
    </div>
  )
}
