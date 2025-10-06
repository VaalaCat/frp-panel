import type { Edge as XEdge, Node as XNode } from '@xyflow/react'
import type { WireGuardLink, WireGuardConfig } from '@/lib/pb/types_wg'
import { ClientStatus } from '@/lib/pb/api_master'
import { TerminalNode } from '@/components/canvas/types'

export type WGNodeData = {
  label: string
  original?: WireGuardConfig
  clientStatus?: ClientStatus
  runtimeInfo?: any
}

export type WGEdgeData = {
  original: WireGuardLink
}

export type WGNode = XNode<WGNodeData, 'wg'>
export type WGEdge = XEdge<WGEdgeData, 'wgEdge'>

export type TopologyNode = WGNode | TerminalNode

export type TopologyData = {
  nodes: TopologyNode[]
  edges: WGEdge[]
}

export type OnEdgeClick = (edgeId: string) => void
export type OnNodeClick = (nodeId: string) => void
