import type { Node as XNode } from '@xyflow/react'
import type { Client, Server } from '@/lib/pb/common'
import { ClientStatus } from '@/lib/pb/api_master'

// ==================== 节点数据类型 ====================

// Client 节点数据
export type ClientNodeData = {
  label: string
  original?: Client
  clientStatus?: ClientStatus
}

// Server 节点数据
export type ServerNodeData = {
  label: string
  original?: Server
  status?: 'online' | 'offline'
}

// 终端节点数据
export type TerminalNodeData = {
  label: string
  clientId: string
  clientType: number
  minimized?: boolean
}

// 日志节点数据
export type LogNodeData = {
  label: string
  clientId: string
  clientType: number
  minimized?: boolean
  pkgs?: string[]
}

// ==================== 节点类型 ====================

export type ClientNode = XNode<ClientNodeData, 'client'>
export type ServerNode = XNode<ServerNodeData, 'server'>
export type TerminalNode = XNode<TerminalNodeData, 'terminal'>
export type LogNode = XNode<LogNodeData, 'log'>

export type CanvasNode = ClientNode | ServerNode | TerminalNode | LogNode

// ==================== 画布数据类型 ====================

export type CanvasData = {
  nodes: CanvasNode[]
}

// ==================== 回调函数类型 ====================

export type OnNodeClick = (nodeId: string) => void
export type OnNodeDelete = (nodeId: string) => void

// ==================== 节点操作接口 ====================

export interface NodeOperations {
  onOpenTerminal?: (clientId: string, clientType: number, sourceNodeId?: string) => void
  onOpenLog?: (clientId: string, clientType: number, sourceNodeId?: string) => void
  onDelete?: (nodeId: string) => void
}
