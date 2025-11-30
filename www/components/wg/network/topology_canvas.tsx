'use client'

import React, { useCallback, useRef, useMemo, useState } from 'react'
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
  ReactFlowInstance,
  Panel,
  Connection,
  ReactFlowProvider,
} from '@xyflow/react'
import type { TopologyData, TopologyNode, WGEdge } from './types'
import WGNodeComponent from './Node'
import WGEdgeComponent from './Edge'
import FloatingConnectionLine from './FloatingConnectionLine'
import { TerminalNode } from '@/components/canvas'
import { nanoid } from 'nanoid'
import { ClientType } from '@/lib/pb/common'
import { Maximize2, Minimize2, ZoomIn, ZoomOut, Locate } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useForceLayout } from './useForceLayout'
import { useTranslation } from 'react-i18next'
import WireGuardLinkForm from '../wireguard-link-form'

export interface TopologyCanvasProps {
  data: TopologyData
  onEdgeClick?: (edgeId: string) => void
  onPaneClick?: () => void
  setNodes: React.Dispatch<React.SetStateAction<TopologyNode[]>>
  setEdges: React.Dispatch<React.SetStateAction<WGEdge[]>>
  fullscreen?: boolean
  onFullscreenToggle?: () => void
  onLinkCreated?: () => void
  onNodeDragStart?: (event: React.MouseEvent, node: any) => void
  onNodeDrag?: (event: React.MouseEvent, node: any) => void
  onNodeDragStop?: (event: React.MouseEvent, node: any) => void
}

function TopologyFlow({
  data,
  onEdgeClick,
  onPaneClick,
  setNodes,
  setEdges,
  fullscreen,
  onFullscreenToggle,
  onLinkCreated,
  onNodeDragStart,
  onNodeDrag,
  onNodeDragStop,
}: TopologyCanvasProps) {
  const { t } = useTranslation()
  const reactFlowInstance = useRef<ReactFlowInstance | null>(null)
  const [linkDialogOpen, setLinkDialogOpen] = useState(false)
  const [newLinkConnection, setNewLinkConnection] = useState<{ from: number; to: number } | null>(null)

  // Use force layout hook
  const { dragEvents } = useForceLayout()

  // 节点变化处理
  const onNodesChange = useCallback(
    (changes: NodeChange[]) => {
      setNodes((nodes) => applyNodeChanges(changes, nodes as any) as any)
    },
    [setNodes]
  )

  // 边变化处理
  const onEdgesChange = useCallback(
    (changes: EdgeChange[]) => {
      setEdges((edges) => applyEdgeChanges(changes, edges as any) as any)
    },
    [setEdges]
  )

  // ReactFlow 实例初始化
  const onInit = useCallback((instance: any) => {
    reactFlowInstance.current = instance
  }, [])

  // 处理连接事件 - 拖拽创建边时触发
  const onConnect = useCallback((connection: Connection) => {
    // 只处理WG节点之间的连接
    if (connection.source && connection.target && connection.source !== connection.target) {
      const fromId = parseInt(connection.source)
      const toId = parseInt(connection.target)

      if (!isNaN(fromId) && !isNaN(toId)) {
        setNewLinkConnection({ from: fromId, to: toId })
        setLinkDialogOpen(true)
      }
    }
  }, [])

  // 处理链接创建成功
  const handleLinkCreated = useCallback(() => {
    setLinkDialogOpen(false)
    setNewLinkConnection(null)
    onLinkCreated?.()
  }, [onLinkCreated])

  // 打开终端
  const handleOpenTerminal = useCallback(
    (clientId: string, clientType: number, sourceNodeId?: string) => {
      if (!reactFlowInstance.current) return

      const instance = reactFlowInstance.current
      let position = { x: 100, y: 100 }

      if (sourceNodeId) {
        const sourceNode = instance.getNodes().find((n) => n.id === sourceNodeId)
        if (sourceNode) {
          position = {
            x: sourceNode.position.x + (sourceNode.width || 280) + 50,
            y: sourceNode.position.y,
          }
        }
      } else {
        const center = instance.screenToFlowPosition({
          x: window.innerWidth / 2,
          y: window.innerHeight / 2,
        })
        position = center
      }

      const newNode: TopologyNode = {
        id: `terminal-${nanoid()}`,
        type: 'terminal',
        position,
        dragHandle: '.drag-handle',
        style: { width: 650, height: 500 },
        data: {
          label: clientId,
          clientId,
          clientType,
        },
      } as any

      setNodes((nodes) => [...nodes, newNode])
    },
    [setNodes]
  )

  // 打开日志
  const handleOpenLog = useCallback(
    (clientId: string, clientType: number, sourceNodeId?: string) => {
      if (!reactFlowInstance.current) return

      const instance = reactFlowInstance.current
      let position = { x: 100, y: 100 }

      if (sourceNodeId) {
        const sourceNode = instance.getNodes().find((n) => n.id === sourceNodeId)
        if (sourceNode) {
          position = {
            x: sourceNode.position.x + (sourceNode.width || 280) + 50,
            y: sourceNode.position.y + 200,
          }
        }
      } else {
        const center = instance.screenToFlowPosition({
          x: window.innerWidth / 2,
          y: window.innerHeight / 2,
        })
        position = center
      }

      const newNode: TopologyNode = {
        id: `log-${nanoid()}`,
        type: 'log',
        position,
        dragHandle: '.drag-handle',
        style: { width: 650, height: 500 },
        data: {
          label: clientId,
          clientId,
          clientType,
          pkgs: ['all'],
        },
      } as any

      setNodes((nodes) => [...nodes, newNode])
    },
    [setNodes]
  )

  // 删除节点
  const handleDeleteNode = useCallback(
    (nodeId: string) => {
      setNodes((nodes) => nodes.filter((n) => n.id !== nodeId))
    },
    [setNodes]
  )

  // 节点类型定义
  const nodeTypes = useMemo(
    () => ({
      wg: (props: any) => (
        <WGNodeComponent
          {...props}
          onOpenTerminal={(clientId: string, clientType: number) =>
            handleOpenTerminal(clientId, clientType, props.id)
          }
          onOpenLog={(clientId: string, clientType: number) =>
            handleOpenLog(clientId, clientType, props.id)
          }
        />
      ),
      terminal: (props: any) => <TerminalNode {...props} onDelete={handleDeleteNode} />,
      log: (props: any) => {
        const LogNode = require('@/components/canvas').LogNode
        return <LogNode {...props} onDelete={handleDeleteNode} />
      },
    }),
    [handleOpenTerminal, handleOpenLog, handleDeleteNode]
  )

  // 边类型定义
  const edgeTypes = useMemo(
    () => ({
      wgEdge: WGEdgeComponent,
    }),
    []
  )

  return (
    <div className="relative w-full h-full rounded-xl overflow-hidden border-2 border-border bg-background shadow-lg" style={{ width: '100%', height: '100%' }}>
      <ReactFlow
        nodes={data.nodes}
        edges={data.edges}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onInit={onInit}
        onConnect={onConnect}
        onEdgeClick={(_, edge) => onEdgeClick?.(edge.id)}
        onPaneClick={onPaneClick}
        onNodeDragStart={(e, n) => {
          dragEvents.start(e as any, n)
          onNodeDragStart?.(e, n)
        }}
        onNodeDrag={(e, n) => {
          dragEvents.drag(e as any, n)
          onNodeDrag?.(e, n)
        }}
        onNodeDragStop={(e, n) => {
          dragEvents.stop()
          onNodeDragStop?.(e, n)
        }}
        connectionLineType={ConnectionLineType.Straight}
        connectionLineComponent={FloatingConnectionLine}
        fitView
        fitViewOptions={{ padding: 0.2, maxZoom: 1.2, minZoom: 0.5, duration: 0 }}
        minZoom={0.1}
        maxZoom={2}
        defaultEdgeOptions={{
          type: 'wgEdge',
        }}
        proOptions={{ hideAttribution: true }}
      >
        {/* 背景网格 */}
        <Background
          gap={20}
          size={1}
          color="#e2e8f0"
          className="bg-muted/20"
        />

        {/* 小地图 */}
        <MiniMap
          nodeStrokeWidth={3}
          zoomable
          pannable
          className="!bg-card !border-2 !border-border !rounded-lg !shadow-lg"
          nodeColor={(node) => {
            if (node.type === 'terminal') return '#8b5cf6'
            return '#6366f1'
          }}
        />

        {/* 默认控制按钮 */}
        <Controls
          showZoom
          showFitView
          showInteractive
          className="!bg-card !border-2 !border-border !rounded-lg !shadow-lg"
        />

        {/* 自定义控制面板 */}
        <Panel position="top-right" className="flex gap-2 m-4">
          {onFullscreenToggle && (
            <Button
              variant="secondary"
              size="icon"
              onClick={onFullscreenToggle}
              className="shadow-lg"
              title={fullscreen ? t('wg.topologyActions.exitFullscreen') : t('wg.topologyActions.fullscreen')}
            >
              {fullscreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
            </Button>
          )}
        </Panel>

        {/* 信息面板 */}
        <Panel position="top-left" className="m-4">
          <div className="bg-card/95 backdrop-blur-sm border-2 border-border rounded-lg shadow-lg px-3 py-2">
            <div className="text-sm font-medium">
              {t('wg.topologyInfo.nodeCount', { count: data.nodes.filter(n => n.type === 'wg').length })} · {t('wg.topologyInfo.linkCount', { count: data.edges.length })}
            </div>
          </div>
        </Panel>
      </ReactFlow>

      {/* 新建连接对话框 */}
      <Dialog open={linkDialogOpen} onOpenChange={setLinkDialogOpen}>
        <DialogContent className="sm:max-w-[600px]">
          <DialogHeader>
            <DialogTitle>{t('wg.linkCreate.title')}</DialogTitle>
            <DialogDescription>
              {t('wg.linkCreate.description', { from: newLinkConnection?.from, to: newLinkConnection?.to })}
            </DialogDescription>
          </DialogHeader>
          {newLinkConnection && (
            <WireGuardLinkForm
              link={{
                id: 0,
                fromWireguardId: newLinkConnection.from,
                toWireguardId: newLinkConnection.to,
                upBandwidthMbps: 100,
                downBandwidthMbps: 100,
                latencyMs: 60,
                active: true,
              } as any}
              onSuccess={handleLinkCreated}
              submitText={t('wg.linkCreate.submit')}
            />
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default function TopologyCanvas(props: TopologyCanvasProps) {
  return (
    <ReactFlowProvider>
      <TopologyFlow {...props} />
    </ReactFlowProvider>
  )
}
