import ELK from 'elkjs'
import type { TopologyNode, WGEdge } from './types'

export const NODE_WIDTH = 200
export const NODE_HEIGHT = 130

type ElkNode = {
  id: string
  width: number
  height: number
  ports: { id: string; properties: Record<string, string> }[]
  properties: Record<string, string>
}

type ElkEdge = {
  id: string
  sources: string[]
  targets: string[]
}

export async function layoutNetwork(nodes: TopologyNode[], edges: WGEdge[]): Promise<{ nodes: TopologyNode[]; edges: WGEdge[] }> {
  const elk = new ELK()

  const elkNodes: ElkNode[] = nodes.map((n) => ({
    id: n.id,
    width: NODE_WIDTH,
    height: NODE_HEIGHT,
    ports: [{ id: `${n.id}:center`, properties: { side: 'SOUTH' } }],
    properties: {
      'org.eclipse.elk.portConstraints': 'FIXED_ORDER',
      'org.eclipse.elk.algorithm': 'org.eclipse.elk.layered',
    },
  }))

  // const idToNode = new Map(nodes.map((n) => [n.id, n]))

  const elkEdges: ElkEdge[] = edges.map((e) => ({
    id: e.id,
    sources: [`${e.source}:center`],
    targets: [`${e.target}:center`],
  }))

  const elkGraph = {
    id: 'root',
    layoutOptions: {
      'elk.direction': 'RIGHT',
      'elk.layered.spacing.nodeNodeBetweenLayers': '150',
      'elk.spacing.nodeNode': '100',
      'elk.spacing.edgeNode': '80',
      'elk.spacing.edgeEdge': '30',
      'elk.edgeRouting': 'ORTHOGONAL',
      'elk.layered.nodePlacement.strategy': 'NETWORK_SIMPLEX',
    },
    children: elkNodes,
    edges: elkEdges,
  }

  const res = await elk.layout(elkGraph as any)

  const nextNodes: TopologyNode[] = nodes.map((n) => {
    const laid = res.children?.find((c: any) => c.id === n.id)
    const data: TopologyNode = {
      ...n,
      position: {
        x: (laid?.x ?? 0) + (n.position?.x ? 0 : 0),
        y: (laid?.y ?? 0) + (n.position?.y ? 0 : 0),
      },
    }

    if (n.type == 'wg') {
      data.style = { width: NODE_WIDTH, height: NODE_HEIGHT, ...(n as any).style }
    } else {
      // data.style = { ...(n as any).style }
    }

    return data
  })

  const nextEdges: WGEdge[] = edges.map((e) => ({ ...e, sourceHandle: undefined, targetHandle: undefined }))

  return { nodes: nextNodes, edges: nextEdges }
}
