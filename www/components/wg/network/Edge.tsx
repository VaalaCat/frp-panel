'use client'

import React from 'react'
import type { EdgeProps } from '@xyflow/react'
import { BaseEdge, getStraightPath, useNodes } from '@xyflow/react'
import type { WGEdge, WGNode } from './types'
import { NODE_WIDTH, NODE_HEIGHT } from './layout'

/**
 * 计算矩形和点的交点
 */
function getRectIntersection(
  rectX: number,
  rectY: number,
  rectWidth: number,
  rectHeight: number,
  pointX: number,
  pointY: number,
  dirX: number,
  dirY: number
): { x: number; y: number } {
  const centerX = rectX + rectWidth / 2
  const centerY = rectY + rectHeight / 2

  // 计算从矩形中心到目标点的方向
  const dx = dirX
  const dy = dirY
  const length = Math.sqrt(dx * dx + dy * dy)
  
  if (length === 0) return { x: centerX, y: centerY }

  const normalizedDx = dx / length
  const normalizedDy = dy / length

  // 计算与矩形边界的交点
  const halfWidth = rectWidth / 2
  const halfHeight = rectHeight / 2

  let t = Infinity
  let intersection = { x: centerX, y: centerY }

  // 检查与左右边的交点
  if (normalizedDx !== 0) {
    const tLeft = -halfWidth / normalizedDx
    const tRight = halfWidth / normalizedDx
    for (const tSide of [tLeft, tRight]) {
      if (tSide > 0 && tSide < t) {
        const y = normalizedDy * tSide
        if (Math.abs(y) <= halfHeight) {
          t = tSide
          intersection = { x: centerX + normalizedDx * tSide, y: centerY + normalizedDy * tSide }
        }
      }
    }
  }

  // 检查与上下边的交点
  if (normalizedDy !== 0) {
    const tTop = -halfHeight / normalizedDy
    const tBottom = halfHeight / normalizedDy
    for (const tSide of [tTop, tBottom]) {
      if (tSide > 0 && tSide < t) {
        const x = normalizedDx * tSide
        if (Math.abs(x) <= halfWidth) {
          t = tSide
          intersection = { x: centerX + normalizedDx * tSide, y: centerY + normalizedDy * tSide }
        }
      }
    }
  }

  return intersection
}

const WGEdgeComponent: React.FC<EdgeProps<WGEdge>> = (props) => {
  const { id, source, target, selected, data, label } = props
  const nodes = useNodes<WGNode>()

  const sourceNode = nodes.find((n) => n.id === source)
  const targetNode = nodes.find((n) => n.id === target)

  if (!sourceNode || !targetNode) {
    return null
  }

  // 获取节点尺寸
  const sourceWidth = sourceNode.width ?? NODE_WIDTH
  const sourceHeight = sourceNode.height ?? NODE_HEIGHT
  const targetWidth = targetNode.width ?? NODE_WIDTH
  const targetHeight = targetNode.height ?? NODE_HEIGHT

  // 计算节点中心点
  const sourceCenterX = sourceNode.position.x + sourceWidth / 2
  const sourceCenterY = sourceNode.position.y + sourceHeight / 2
  const targetCenterX = targetNode.position.x + targetWidth / 2
  const targetCenterY = targetNode.position.y + targetHeight / 2

  // 计算方向向量
  let dirX = targetCenterX - sourceCenterX
  let dirY = targetCenterY - sourceCenterY
  
  // 检测双向边：按字典序排列source和target，确定偏移方向
  // 这样同一对节点间的两条边会使用一致的偏移方向
  const hasBidirectional = source !== target
  const isSecondaryEdge = hasBidirectional && id.localeCompare(`${target}-${source}`) > 0
  
  // 如果是双向边，给边添加一个小的垂直偏移
  const OFFSET = 12 // 偏移像素，增加到12以避免重叠
  if (isSecondaryEdge) {
    // 计算垂直于连线的方向
    const length = Math.sqrt(dirX * dirX + dirY * dirY)
    if (length > 0) {
      const perpX = -dirY / length * OFFSET
      const perpY = dirX / length * OFFSET
      
      // 偏移中心点
      const offsetSourceX = sourceCenterX + perpX
      const offsetSourceY = sourceCenterY + perpY
      const offsetTargetX = targetCenterX + perpX
      const offsetTargetY = targetCenterY + perpY
      
      // 重新计算方向
      dirX = offsetTargetX - offsetSourceX
      dirY = offsetTargetY - offsetSourceY
      
      // 使用偏移后的中心点计算交点
      const sourceIntersection = getRectIntersection(
        sourceNode.position.x,
        sourceNode.position.y,
        sourceWidth,
        sourceHeight,
        offsetTargetX,
        offsetTargetY,
        dirX,
        dirY
      )

      const targetIntersection = getRectIntersection(
        targetNode.position.x,
        targetNode.position.y,
        targetWidth,
        targetHeight,
        offsetSourceX,
        offsetSourceY,
        -dirX,
        -dirY
      )
      
      // 再次偏移交点
      const sourceOffsetIntersection = {
        x: sourceIntersection.x + perpX,
        y: sourceIntersection.y + perpY
      }
      const targetOffsetIntersection = {
        x: targetIntersection.x + perpX,
        y: targetIntersection.y + perpY
      }
      
      const [edgePath] = getStraightPath({
        sourceX: sourceOffsetIntersection.x,
        sourceY: sourceOffsetIntersection.y,
        targetX: targetOffsetIntersection.x,
        targetY: targetOffsetIntersection.y,
      })
      
      return renderEdge(edgePath, sourceOffsetIntersection, targetOffsetIntersection, dirX, dirY)
    }
  }

  // 计算边的起点和终点（在矩形边界上）
  const sourceIntersection = getRectIntersection(
    sourceNode.position.x,
    sourceNode.position.y,
    sourceWidth,
    sourceHeight,
    targetCenterX,
    targetCenterY,
    dirX,
    dirY
  )

  const targetIntersection = getRectIntersection(
    targetNode.position.x,
    targetNode.position.y,
    targetWidth,
    targetHeight,
    sourceCenterX,
    sourceCenterY,
    -dirX,
    -dirY
  )

  const [edgePath] = getStraightPath({
    sourceX: sourceIntersection.x,
    sourceY: sourceIntersection.y,
    targetX: targetIntersection.x,
    targetY: targetIntersection.y,
  })
  
  return renderEdge(edgePath, sourceIntersection, targetIntersection, dirX, dirY)
  
  function renderEdge(
    path: string,
    sourcePos: { x: number; y: number },
    targetPos: { x: number; y: number },
    dx: number,
    dy: number
  ) {
    // 边的样式配置
    const isActive = data?.link?.active ?? false
    const quality = data?.quality || 'fair'
    
    const colorMap: Record<string, string> = {
      excellent: '#10b981',
      good: '#22c55e',
      fair: '#94a3b8',
      poor: '#ef4444',
    }

    const strokeColor = selected
      ? '#6366f1'
      : isActive
      ? colorMap[quality]
      : '#cbd5e1'

    const strokeWidth = selected ? 3 : isActive ? 2.5 : 1.5

    // 计算标签位置
    const labelX = (sourcePos.x + targetPos.x) / 2
    const labelY = (sourcePos.y + targetPos.y) / 2

    // 计算标签角度
    const angle = Math.atan2(dy, dx) * (180 / Math.PI)
    const labelAngle = angle > 90 || angle < -90 ? angle + 180 : angle

    // 构建显示标签 - 简化格式
    const displayLabel = String(label || (data?.link 
      ? `${data.link.latencyMs}ms | ${data.link.upBandwidthMbps}↑ ${data.link.downBandwidthMbps}↓`
      : ''))

    return (
      <>
        <BaseEdge
          id={id}
          path={path}
          style={{
            stroke: strokeColor,
            strokeWidth,
            strokeDasharray: isActive ? undefined : '6 4',
            strokeLinecap: 'round',
          }}
        />
        
        {displayLabel && (
          <g transform={`translate(${labelX}, ${labelY})`}>
            <g transform={`rotate(${labelAngle})`}>
              {/* 标签背景 */}
              <rect
                x={-displayLabel.length * 2.8}
                y={-9}
                width={displayLabel.length * 5.6}
                height={18}
                rx={3}
                fill="white"
                fillOpacity={0.95}
                stroke={strokeColor}
                strokeWidth={1}
                className="pointer-events-none"
              />
              {/* 标签文本 */}
              <text
                className="text-[10px] font-medium pointer-events-none select-none"
                textAnchor="middle"
                dominantBaseline="middle"
                fill={strokeColor}
              >
                {displayLabel}
              </text>
            </g>
          </g>
        )}
      </>
    )
  }
}

export default WGEdgeComponent
