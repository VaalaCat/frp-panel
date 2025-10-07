'use client'

import React from 'react'
import type { EdgeProps } from '@xyflow/react'
import { BaseEdge, getStraightPath, useNodes, useEdges } from '@xyflow/react'
import type { WGEdge, WGNode } from './types'
import { NODE_WIDTH, NODE_HEIGHT } from './layout'

interface NodeRect {
  x: number
  y: number
  width: number
  height: number
  cx: number
  cy: number
}

interface Point {
  x: number
  y: number
}

interface Vector {
  x: number
  y: number
}

function getNodeRect(n?: WGNode): NodeRect {
  const width = (n?.width ?? (n as any)?.style?.width ?? NODE_WIDTH) as number
  const height = (n?.height ?? (n as any)?.style?.height ?? NODE_HEIGHT) as number
  const x = n?.position?.x ?? 0
  const y = n?.position?.y ?? 0
  return { x, y, width, height, cx: x + width / 2, cy: y + height / 2 }
}

function intersectRectFromPoint(rect: NodeRect, px: number, py: number, vx: number, vy: number): Point {
  const { left, right, top, bottom } = {
    left: rect.x,
    right: rect.x + rect.width,
    top: rect.y,
    bottom: rect.y + rect.height,
  }

  const points: Array<{ t: number; x: number; y: number }> = []

  if (vx !== 0) {
    let t = (left - px) / vx
    let y = py + t * vy
    if (t >= 0 && y >= top && y <= bottom) points.push({ t, x: left, y })

    t = (right - px) / vx
    y = py + t * vy
    if (t >= 0 && y >= top && y <= bottom) points.push({ t, x: right, y })
  }

  if (vy !== 0) {
    let t = (top - py) / vy
    let x = px + t * vx
    if (t >= 0 && x >= left && x <= right) points.push({ t, x, y: top })

    t = (bottom - py) / vy
    x = px + t * vx
    if (t >= 0 && x >= left && x <= right) points.push({ t, x, y: bottom })
  }

  if (points.length === 0) return { x: px, y: py }
  points.sort((a, b) => a.t - b.t)
  return { x: points[0].x, y: points[0].y }
}

function intersectRect(rect: NodeRect, vx: number, vy: number): Point {
  return intersectRectFromPoint(rect, rect.cx, rect.cy, vx, vy)
}

function getPerpendicularUnit(v: Vector): Vector {
  const len = Math.sqrt(v.x * v.x + v.y * v.y)
  if (len === 0) return { x: 0, y: 0 }
  return { x: -v.y / len, y: v.x / len }
}

function getBidirectionalOffset(
  edges: WGEdge[],
  currentId: string,
  source: string,
  target: string
): number {
  const reverseEdge = edges.find(
    (e) => e.source === target && e.target === source && e.id !== currentId
  )
  if (!reverseEdge) return 0
  return currentId.localeCompare(reverseEdge.id) > 0 ? -1 : 1
}

const Edge: React.FC<EdgeProps<WGEdge>> = (props) => {
  const { id, source, target, selected, markerEnd, data, label } = props
  const nodes = useNodes<WGNode>()
  const edges = useEdges<WGEdge>()

  const sourceRect = getNodeRect(nodes.find((n) => n.id === source))
  const targetRect = getNodeRect(nodes.find((n) => n.id === target))

  const baseVector: Vector = {
    x: targetRect.cx - sourceRect.cx,
    y: targetRect.cy - sourceRect.cy,
  }

  const offsetDirection = getBidirectionalOffset(edges, id, source, target)
  const hasBidirectional = offsetDirection !== 0

  const LINE_OFFSET = 12
  let edgeStart: Point
  let edgeEnd: Point
  let perpUnit: Vector | null = null

  if (hasBidirectional) {
    // 使用统一的基准向量确保双向边偏移方向相反
    const [nodeA, nodeB] = source < target ? [sourceRect, targetRect] : [targetRect, sourceRect]
    const canonical: Vector = { x: nodeB.cx - nodeA.cx, y: nodeB.cy - nodeA.cy }
    perpUnit = getPerpendicularUnit(canonical)

    const offset = LINE_OFFSET * offsetDirection
    const virtualSrc = {
      cx: sourceRect.cx + perpUnit.x * offset,
      cy: sourceRect.cy + perpUnit.y * offset,
    }
    const virtualTgt = {
      cx: targetRect.cx + perpUnit.x * offset,
      cy: targetRect.cy + perpUnit.y * offset,
    }

    const offsetVec: Vector = { x: virtualTgt.cx - virtualSrc.cx, y: virtualTgt.cy - virtualSrc.cy }
    edgeStart = intersectRectFromPoint(sourceRect, virtualSrc.cx, virtualSrc.cy, offsetVec.x, offsetVec.y)
    edgeEnd = intersectRectFromPoint(targetRect, virtualTgt.cx, virtualTgt.cy, -offsetVec.x, -offsetVec.y)
  } else {
    edgeStart = intersectRect(sourceRect, baseVector.x, baseVector.y)
    edgeEnd = intersectRect(targetRect, -baseVector.x, -baseVector.y)
  }

  const [edgePath] = getStraightPath({
    sourceX: edgeStart.x,
    sourceY: edgeStart.y,
    targetX: edgeEnd.x,
    targetY: edgeEnd.y,
  })

  const isActive = data?.original?.active ?? false
  const strokeColor = selected ? '#6366f1' : isActive ? '#10b981' : '#94a3b8'
  const strokeWidth = selected ? 2.5 : isActive ? 2 : 1.5

  const toEndpoint = data?.original?.toEndpoint
  const edgeLabel = toEndpoint
    ? `${toEndpoint.host}:${toEndpoint.port}${label ? ` (${label})` : ''}`
    : label

  // 计算标签位置和角度
  const labelCenter: Point = {
    x: (edgeStart.x + edgeEnd.x) / 2,
    y: (edgeStart.y + edgeEnd.y) / 2,
  }

  const edgeVec: Vector = { x: edgeEnd.x - edgeStart.x, y: edgeEnd.y - edgeStart.y }
  const edgeLen = Math.sqrt(edgeVec.x * edgeVec.x + edgeVec.y * edgeVec.y)

  let labelAngle = Math.atan2(edgeVec.y, edgeVec.x) * (180 / Math.PI)
  let flipText = false
  if (labelAngle > 90 || labelAngle < -90) {
    labelAngle += 180
    flipText = true
  }

  let labelPos = labelCenter
  if (edgeLen > 0) {
    const labelPerp = getPerpendicularUnit(edgeVec)
    // 对于双向边，使用与线条相同的偏移方向
    const labelOffsetDir = hasBidirectional ? offsetDirection : (flipText ? -1 : 1)
    labelPos = {
      x: labelCenter.x + labelPerp.x * 12 * labelOffsetDir,
      y: labelCenter.y + labelPerp.y * 12 * labelOffsetDir,
    }
  }

  return (
    <>
      <BaseEdge
        id={id}
        path={edgePath}
        markerEnd={markerEnd}
        style={{
          stroke: strokeColor,
          strokeWidth,
          strokeDasharray: isActive ? undefined : '5 5',
        }}
      />
      {edgeLabel && (
        <g transform={`translate(${labelPos.x}, ${labelPos.y}) rotate(${labelAngle})`}>
          <text
            className="text-[10px] fill-foreground pointer-events-none select-none"
            textAnchor="middle"
            dominantBaseline="middle"
          >
            <tspan
              style={{
                fill: isActive ? '#10b981' : '#94a3b8',
                fontWeight: selected ? 600 : 400,
              }}
            >
              {edgeLabel}
            </tspan>
          </text>
        </g>
      )}
    </>
  )
}

export default Edge
