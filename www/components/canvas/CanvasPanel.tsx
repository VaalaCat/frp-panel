'use client'

import React, { useCallback, useState } from 'react'
import { Button } from '@/components/ui/button'
import Canvas from './Canvas'
import type { CanvasNode, ClientNode, ServerNode } from './types'
import { useTranslation } from 'react-i18next'
import { Maximize2, Minimize2, Plus, Trash2 } from 'lucide-react'
import { ReactFlowProvider } from '@xyflow/react'
import { AgentSelector, type Agent } from '@/components/base/agent-selector'
import { Client, Server } from '@/lib/pb/common'

export default function CanvasPanel() {
  const { t } = useTranslation()
  const [fullscreen, setFullscreen] = useState(false)
  const [selectedAgent, setSelectedAgent] = useState<Agent | undefined>()
  const [focusNodeId, setFocusNodeId] = useState<string | undefined>()

  const [nodes, setNodes] = useState<CanvasNode[]>([])
  const [agentNodeCount, setAgentNodeCount] = useState(0)

  const handleAddAgent = useCallback(() => {
    if (!selectedAgent) return

    const nodeId =
      selectedAgent.type === 'client'
        ? `client-${(selectedAgent.original as Client).id}`
        : `server-${(selectedAgent.original as Server).id}`

    // 检查是否已存在
    if (nodes.some((n) => n.id === nodeId)) {
      // 如果已存在，聚焦到该节点
      setFocusNodeId(nodeId)
      setSelectedAgent(undefined)
      return
    }

    // 根据类型创建节点
    if (selectedAgent.type === 'client') {
      const client = selectedAgent.original as Client
      const newNode: ClientNode = {
        id: nodeId,
        type: 'client',
        dragHandle: '.drag-handle',
        data: {
          label: client.id || 'Client',
          original: client,
        },
        position: {
          x: 0,
          y: agentNodeCount * 100 + 100,
        },
      }
      setNodes((prev) => [...prev, newNode])
    } else {
      const server = selectedAgent.original as Server
      const newNode: ServerNode = {
        id: nodeId,
        type: 'server',
        dragHandle: '.drag-handle',
        data: {
          label: server.id || 'Server',
          original: server,
          status: 'offline',
        },
        position: {
          x: 0,
          y: agentNodeCount * 100 + 100,
        },
      }
      setNodes((prev) => [...prev, newNode])
    }

    setAgentNodeCount(agentNodeCount + 1)

    // 聚焦到新节点
    setFocusNodeId(nodeId)
    // 清空选择
    setSelectedAgent(undefined)
  }, [selectedAgent, nodes, setNodes, setFocusNodeId, setSelectedAgent, agentNodeCount, setAgentNodeCount])

  const handleClearCanvas = () => {
    if (confirm(t('canvas.panel.confirmClear'))) {
      // 只保留 client 和 server 节点
      setNodes(nodes.filter((n) => n.type === 'client' || n.type === 'server'))
      setAgentNodeCount(nodes.filter((n) => n.type === 'client' || n.type === 'server').length)
    }
  }

  return (
    <ReactFlowProvider>
      <div className="flex flex-col gap-3">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div className="flex flex-wrap items-center gap-2">
            <h2 className="text-lg font-semibold">{t('canvas.panel.title')}</h2>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" onClick={handleClearCanvas} title={t('canvas.panel.clear')}>
              <Trash2 className="h-4 w-4 mr-2" />
              {t('canvas.panel.clear')}
            </Button>
            <Button
              variant="outline"
              size="icon"
              onClick={() => setFullscreen(!fullscreen)}
              className="hidden md:flex"
              title={fullscreen ? t('canvas.panel.exitFullscreen') : t('canvas.panel.fullscreen')}
            >
              {fullscreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
            </Button>
          </div>
        </div>

        <div className="flex flex-wrap gap-4 p-3 bg-muted/50 rounded-lg items-center">
          <div className="flex items-center gap-2 flex-1 min-w-[300px] max-w-[500px]">
            <AgentSelector
              value={selectedAgent}
              onChange={setSelectedAgent}
              placeholder={t('canvas.panel.selectAgent')}
              className="flex-1"
            />
            <Button variant="default" size="sm" onClick={handleAddAgent} disabled={!selectedAgent}>
              <Plus className="h-4 w-4 mr-2" />
              {t('canvas.panel.addNode')}
            </Button>
          </div>
          <div className="flex-1"></div>

          <div className="text-sm text-muted-foreground">
            {t('canvas.panel.nodeCount')}: {agentNodeCount}
          </div>
        </div>

        <div className={fullscreen ? 'fixed inset-0 z-50 bg-background flex flex-col' : ''}>
          {fullscreen && (
            <div className="flex items-center justify-end p-3 border-b">
              <Button
                variant="outline"
                size="icon"
                onClick={() => setFullscreen(false)}
                title={t('canvas.panel.exitFullscreen')}
              >
                <Minimize2 className="h-4 w-4" />
              </Button>
            </div>
          )}
          <div className={fullscreen ? 'flex-1 min-h-0' : ''}>
            <Canvas
              data={{ nodes }}
              setNodes={setNodes}
              focusNodeId={focusNodeId}
              onFocusComplete={() => setFocusNodeId(undefined)}
              fullscreen={fullscreen}
            />
          </div>
        </div>
      </div>
    </ReactFlowProvider>
  )
}
