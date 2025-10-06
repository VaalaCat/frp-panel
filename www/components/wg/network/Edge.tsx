'use client'

import React from 'react'
import type { EdgeProps } from '@xyflow/react'
import { BaseEdge, getStraightPath, useNodes } from '@xyflow/react'
import type { WGEdge, WGNode } from './types'
import { NODE_WIDTH, NODE_HEIGHT } from './layout'

function getNodeRect(n?: WGNode) {
  const width = (n?.width ?? (n as any)?.style?.width ?? NODE_WIDTH) as number
  const height = (n?.height ?? (n as any)?.style?.height ?? NODE_HEIGHT) as number
  const x = n?.position?.x ?? 0
  const y = n?.position?.y ?? 0
  return { x, y, width, height, cx: x + width / 2, cy: y + height / 2 }
}

function intersectRect(
  rect: { x: number; y: number; width: number; height: number; cx: number; cy: number },
  vx: number,
  vy: number,
) {
  const left = rect.x
  const right = rect.x + rect.width
  const top = rect.y
  const bottom = rect.y + rect.height
  const cx = rect.cx
  const cy = rect.cy

  const points: Array<{ t: number; x: number; y: number }> = []
  if (vx !== 0) {
    let t = (left - cx) / vx
    let y = cy + t * vy
    if (t >= 0 && y >= top && y <= bottom) points.push({ t, x: left, y })
    t = (right - cx) / vx
    y = cy + t * vy
    if (t >= 0 && y >= top && y <= bottom) points.push({ t, x: right, y })
  }
  if (vy !== 0) {
    let t = (top - cy) / vy
    let x = cx + t * vx
    if (t >= 0 && x >= left && x <= right) points.push({ t, x, y: top })
    t = (bottom - cy) / vy
    x = cx + t * vx
    if (t >= 0 && x >= left && x <= right) points.push({ t, x, y: bottom })
  }
  if (points.length === 0) return { x: cx, y: cy }
  points.sort((a, b) => a.t - b.t)
  return { x: points[0].x, y: points[0].y }
}

const Edge: React.FC<EdgeProps<WGEdge>> = (props) => {
  const { id, source, target, selected, markerEnd, data, label } = props
  const nodes = useNodes<WGNode>()
  const sn = nodes.find((n) => n.id === source)
  const tn = nodes.find((n) => n.id === target)

  const sr = getNodeRect(sn)
  const tr = getNodeRect(tn)

  const vx = tr.cx - sr.cx
  const vy = tr.cy - sr.cy

  const s = intersectRect(sr, vx, vy)
  const t = intersectRect(tr, -vx, -vy)

  const [edgePath, labelX, labelY] = getStraightPath({ sourceX: s.x, sourceY: s.y, targetX: t.x, targetY: t.y })

  const isActive = data?.original?.active ?? false
  const strokeColor = selected ? '#6366f1' : isActive ? '#10b981' : '#94a3b8'
  const strokeWidth = selected ? 2.5 : isActive ? 2 : 1.5

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
      {label && (
        <text
          x={labelX}
          y={labelY}
          className="text-[10px] fill-foreground pointer-events-none select-none"
          textAnchor="middle"
          dy={-5}
        >
          <tspan
            className="bg-background px-1 py-0.5 rounded"
            style={{
              fill: isActive ? '#10b981' : '#94a3b8',
              fontWeight: selected ? 600 : 400,
            }}
          >
            {label}
          </tspan>
        </text>
      )}
    </>
  )
}

export default Edge
