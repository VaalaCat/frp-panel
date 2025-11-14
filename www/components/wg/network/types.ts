import type { Edge as XEdge, Node as XNode } from '@xyflow/react'
import type { WireGuardLink, WireGuardConfig, WGDeviceRuntimeInfo, WGPeerRuntimeInfo } from '@/lib/pb/types_wg'
import type { ClientStatus } from '@/lib/pb/api_master'
import { TerminalNode } from '@/components/canvas/types'

/**
 * WireGuard 节点数据
 */
export type WGNodeData = {
  label: string
  config?: WireGuardConfig
  runtime?: WGDeviceRuntimeInfo
  clientStatus?: ClientStatus
}

/**
 * WireGuard 边数据
 */
export type WGEdgeData = {
  link: WireGuardLink
  // 计算出的性能指标
  quality?: 'excellent' | 'good' | 'fair' | 'poor'
  // 是否为双向连接
  isBidirectional?: boolean
}

/**
 * WireGuard 节点类型
 */
export type WGNode = XNode<WGNodeData, 'wg'>

/**
 * WireGuard 边类型
 */
export type WGEdge = XEdge<WGEdgeData, 'wgEdge'>

/**
 * 拓扑节点（包括WG节点和终端节点）
 */
export type TopologyNode = WGNode | TerminalNode

/**
 * 拓扑数据
 */
export type TopologyData = {
  nodes: TopologyNode[]
  edges: WGEdge[]
}

/**
 * 节点统计信息
 */
export type NodeStats = {
  peerCount: number
  totalTx: number
  totalRx: number
  avgLatency?: number
  lastHandshake?: number
  isOnline: boolean
}

/**
 * 边统计信息
 */
export type EdgeStats = {
  latency: number
  upBandwidth: number
  downBandwidth: number
  isActive: boolean
  packetLoss?: number
}

/**
 * 回调函数类型
 */
export type OnEdgeClick = (edgeId: string) => void
export type OnNodeClick = (nodeId: string) => void
export type OnNodeDoubleClick = (nodeId: string) => void
