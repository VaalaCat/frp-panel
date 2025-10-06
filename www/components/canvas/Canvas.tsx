'use client'

import React, { useCallback, useRef } from 'react'
import {
  applyNodeChanges,
  Background,
  Controls,
  MiniMap,
  ReactFlow,
  NodeChange,
  ReactFlowInstance,
} from '@xyflow/react'
import type { CanvasData, CanvasNode } from './types'
import ClientNode from './ClientNode'
import ServerNode from './ServerNode'
import TerminalNode from './TerminalNode'
import LogNode from './LogNode'
import { nanoid } from 'nanoid'

export interface CanvasProps {
  data: CanvasData
  setNodes: React.Dispatch<React.SetStateAction<CanvasNode[]>>
  onDeleteNode?: (nodeId: string) => void
  focusNodeId?: string
  onFocusComplete?: () => void
  fullscreen?: boolean
}

export default function Canvas({
  data,
  setNodes,
  onDeleteNode,
  focusNodeId,
  onFocusComplete,
  fullscreen = false,
}: CanvasProps) {
  const reactFlowInstance = useRef<ReactFlowInstance | null>(null)

  const onNodesChange = useCallback(
    (changes: NodeChange[]) => setNodes((nodesSnapshot) => applyNodeChanges(changes, nodesSnapshot as any) as any),
    [setNodes],
  )

  const onInit = useCallback((instance: any) => {
    reactFlowInstance.current = instance
  }, [])

  // 聚焦到指定节点
  React.useEffect(() => {
    if (focusNodeId && reactFlowInstance.current) {
      const node = reactFlowInstance.current.getNodes().find((n) => n.id === focusNodeId)
      if (node) {
        reactFlowInstance.current.setCenter(
          node.position.x + (node.width || 200) / 2,
          node.position.y + (node.height || 100) / 2,
        )

        // 高亮节点
        reactFlowInstance.current.setNodes((nds: any) =>
          nds.map((n: any) => ({
            ...n,
            selected: n.id === focusNodeId,
          })),
        )
        onFocusComplete?.()
      }
    }
  }, [focusNodeId, onFocusComplete])

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
        dragHandle: '.drag-handle',
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

  const handleOpenLog = useCallback(
    (clientId: string, clientType: number, sourceNodeId?: string) => {
      if (!reactFlowInstance.current) return

      const instance = reactFlowInstance.current

      let position = { x: 0, y: 0 }

      if (sourceNodeId) {
        const sourceNode = instance.getNodes().find((n: any) => n.id === sourceNodeId)
        if (sourceNode) {
          position = {
            x: sourceNode.position.x + (sourceNode.width || 200) + 50,
            y: sourceNode.position.y + 200,
          }
        }
      }

      if (position.x === 0 && position.y === 0) {
        const viewport = instance.getViewport()
        const bounds = document.querySelector('.react-flow__viewport')?.getBoundingClientRect()

        const centerX = bounds ? bounds.width / 2 : 400
        const centerY = bounds ? bounds.height / 2 : 400

        position = instance.screenToFlowPosition({
          x: centerX,
          y: centerY,
        })
      }

      const newNode = {
        id: `log-${nanoid()}`,
        type: 'log',
        position,
        dragHandle: '.drag-handle',
        style: { width: 650, height: 500 },
        data: {
          label: `Log - ${clientId}`,
          clientId,
          clientType,
          minimized: false,
          pkgs: ['all'],
        },
      }

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
      client: (props: any) => (
        <ClientNode
          {...props}
          onOpenTerminal={(clientId: string, clientType: number) => handleOpenTerminal(clientId, clientType, props.id)}
          onOpenLog={(clientId: string, clientType: number) => handleOpenLog(clientId, clientType, props.id)}
        />
      ),
      server: (props: any) => (
        <ServerNode
          {...props}
          onOpenTerminal={(clientId: string, clientType: number) => handleOpenTerminal(clientId, clientType, props.id)}
          onOpenLog={(clientId: string, clientType: number) => handleOpenLog(clientId, clientType, props.id)}
        />
      ),
      terminal: (props: any) => <TerminalNode {...props} onDelete={handleDeleteNode} />,
      log: (props: any) => <LogNode {...props} onDelete={handleDeleteNode} />,
    }),
    [handleOpenTerminal, handleOpenLog, handleDeleteNode],
  )

  const canvasClassName = fullscreen ? 'h-full w-full' : 'w-full border rounded-lg overflow-hidden shadow-sm'

  const canvasStyle = fullscreen ? undefined : { height: '80dvh' }

  return (
    <div className={canvasClassName} style={canvasStyle}>
      <ReactFlow
        nodes={data.nodes}
        edges={[]}
        nodeTypes={nodeTypes}
        onNodesChange={onNodesChange}
        onInit={onInit}
        fitView
        minZoom={0.1}
        maxZoom={2}
        defaultViewport={{ x: 0, y: 0, zoom: 0.8 }}
      >
        <MiniMap nodeStrokeWidth={3} zoomable pannable />
        <Controls />
        <Background gap={12} size={1} />
      </ReactFlow>
    </div>
  )
}
