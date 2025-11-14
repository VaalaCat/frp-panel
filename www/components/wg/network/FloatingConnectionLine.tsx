'use client'

import React from 'react'
import { getStraightPath } from '@xyflow/react'
import type { ConnectionLineComponentProps } from '@xyflow/react'

/**
 * 浮动连接线组件
 * 在用户拖拽创建新连接时显示
 */
const FloatingConnectionLine: React.FC<ConnectionLineComponentProps> = ({
  fromX,
  fromY,
  toX,
  toY,
  fromNode,
}) => {
  const [edgePath] = getStraightPath({
    sourceX: fromX,
    sourceY: fromY,
    targetX: toX,
    targetY: toY,
  })

  return (
    <g>
      <path
        fill="none"
        stroke="#6366f1"
        strokeWidth={2.5}
        className="animated"
        d={edgePath}
        strokeDasharray="8 4"
        strokeLinecap="round"
        style={{
          animation: 'dashdraw 0.5s linear infinite',
        }}
      />
      <circle
        cx={toX}
        cy={toY}
        fill="#6366f1"
        r={4}
        stroke="#fff"
        strokeWidth={2}
      />
      <style>{`
        @keyframes dashdraw {
          to {
            stroke-dashoffset: -12;
          }
        }
      `}</style>
    </g>
  )
}

export default FloatingConnectionLine
