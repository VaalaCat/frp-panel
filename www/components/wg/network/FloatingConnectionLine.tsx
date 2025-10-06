'use client'

import React from 'react'
import { BaseEdge, getStraightPath, useNodes } from '@xyflow/react'
import type { ConnectionLineComponentProps } from '@xyflow/react'
import type { WGNode } from './types'
import { NODE_WIDTH, NODE_HEIGHT } from './layout'

function getNodeRect(n?: WGNode) {
  const width = (n?.width ?? (n as any)?.style?.width ?? NODE_WIDTH) as number
  const height = (n?.height ?? (n as any)?.style?.height ?? NODE_HEIGHT) as number
  const x = n?.position?.x ?? 0
  const y = n?.position?.y ?? 0
  return { x, y, width, height, cx: x + width / 2, cy: y + height / 2 }
}

function contains(rect: { x: number; y: number; width: number; height: number }, px: number, py: number) {
  return px >= rect.x && px <= rect.x + rect.width && py >= rect.y && py <= rect.y + rect.height
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

const FloatingConnectionLine: React.FC<ConnectionLineComponentProps> = (props) => {
  const { fromX, fromY, toX, toY } = props
  const nodes = useNodes<WGNode>()

  const srcNode = nodes.find((n) => {
    const r = getNodeRect(n)
    return contains(r, fromX, fromY)
  })
  const tgtNode = nodes.find((n) => {
    const r = getNodeRect(n)
    return contains(r, toX, toY)
  })

  let sx = fromX
  let sy = fromY
  let tx = toX
  let ty = toY

  if (srcNode) {
    const sr = getNodeRect(srcNode)
    const vx = toX - sr.cx
    const vy = toY - sr.cy
    const p = intersectRect(sr, vx, vy)
    sx = p.x
    sy = p.y
  }
  if (tgtNode) {
    const tr = getNodeRect(tgtNode)
    const vx = fromX - tr.cx
    const vy = fromY - tr.cy
    const p = intersectRect(tr, vx, vy)
    tx = p.x
    ty = p.y
  }

  const [path] = getStraightPath({ sourceX: sx, sourceY: sy, targetX: tx, targetY: ty })
  return <BaseEdge path={path} style={{ stroke: '#94a3b8', strokeDasharray: '4 4' }} />
}

export default FloatingConnectionLine
