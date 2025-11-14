import ELK from 'elkjs'
import type { TopologyNode, WGEdge } from './types'

// 节点尺寸常量
export const NODE_WIDTH = 240
export const NODE_HEIGHT = 150

// 布局参数
const LAYOUT_CONFIG = {
  // 节点间距
  nodeSpacing: 120,
  // 层间距
  layerSpacing: 200,
  // 边到节点的距离
  edgeNodeSpacing: 60,
  // 边到边的距离
  edgeEdgeSpacing: 20,
}

/**
 * 使用 ELK 算法进行网络布局
 */
export async function layoutNetwork(
  nodes: TopologyNode[],
  edges: WGEdge[]
): Promise<{ nodes: TopologyNode[]; edges: WGEdge[] }> {
  if (nodes.length === 0) {
    return { nodes: [], edges: [] }
  }

  const elk = new ELK()

  // 为非WG节点（如终端节点）保留原始位置
  const wgNodes = nodes.filter((n) => n.type === 'wg')
  const otherNodes = nodes.filter((n) => n.type !== 'wg')

  if (wgNodes.length === 0) {
    return { nodes, edges }
  }

  // 构建 ELK 图结构
  const elkNodes = wgNodes.map((n) => ({
    id: n.id,
    width: NODE_WIDTH,
    height: NODE_HEIGHT,
  }))

  const elkEdges = edges
    .filter((e) => wgNodes.some((n) => n.id === e.source) && wgNodes.some((n) => n.id === e.target))
    .map((e) => ({
      id: e.id,
      sources: [e.source],
      targets: [e.target],
    }))

  // ELK 布局配置
  const elkGraph = {
    id: 'root',
    layoutOptions: {
      'elk.algorithm': 'layered',
      'elk.direction': 'RIGHT',
      'elk.spacing.nodeNode': String(LAYOUT_CONFIG.nodeSpacing),
      'elk.layered.spacing.nodeNodeBetweenLayers': String(LAYOUT_CONFIG.layerSpacing),
      'elk.spacing.edgeNode': String(LAYOUT_CONFIG.edgeNodeSpacing),
      'elk.spacing.edgeEdge': String(LAYOUT_CONFIG.edgeEdgeSpacing),
      'elk.layered.nodePlacement.strategy': 'NETWORK_SIMPLEX',
      'elk.layered.crossingMinimization.strategy': 'LAYER_SWEEP',
      'elk.edgeRouting': 'ORTHOGONAL',
      'elk.layered.unnecessaryBendpoints': 'true',
      'elk.layered.spacing.edgeNodeBetweenLayers': '60',
    },
    children: elkNodes,
    edges: elkEdges,
  }

  try {
    const layoutResult = await elk.layout(elkGraph as any)

    // 应用布局结果
    const layoutedWgNodes: TopologyNode[] = wgNodes.map((node) => {
      const elkNode = layoutResult.children?.find((c: any) => c.id === node.id)
      if (!elkNode) return node

      return {
        ...node,
        position: {
          x: elkNode.x ?? node.position.x,
          y: elkNode.y ?? node.position.y,
        },
        width: NODE_WIDTH,
        height: NODE_HEIGHT,
        style: {
          width: NODE_WIDTH,
          height: NODE_HEIGHT,
        },
      } as TopologyNode
    })

    // 合并布局后的节点
    const allNodes = [...layoutedWgNodes, ...otherNodes]

    return { nodes: allNodes, edges }
  } catch (error) {
    console.error('Layout error:', error)
    // 如果布局失败，返回原始数据但设置固定尺寸
    return {
      nodes: nodes.map((n) => ({
        ...n,
        width: n.type === 'wg' ? NODE_WIDTH : n.width,
        height: n.type === 'wg' ? NODE_HEIGHT : n.height,
        style: {
          ...((n as any).style || {}),
          width: n.type === 'wg' ? NODE_WIDTH : ((n as any).style?.width || n.width),
          height: n.type === 'wg' ? NODE_HEIGHT : ((n as any).style?.height || n.height),
        },
      })),
      edges,
    }
  }
}

/**
 * 计算连接质量等级
 * 主要基于延迟判断，带宽作为参考
 */
export function calculateConnectionQuality(latency: number, bandwidth: number): 'excellent' | 'good' | 'fair' | 'poor' {
  // 优秀：延迟 0-50ms
  if (latency <= 50) return 'excellent'
  // 一般：延迟 51-200ms
  if (latency <= 200) return 'fair'
  // 过高：延迟 200ms 以上
  return 'poor'
}
