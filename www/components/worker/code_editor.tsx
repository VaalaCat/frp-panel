'use client'

import React from 'react'
import dynamic from 'next/dynamic'

// 动态加载编辑器组件
const MonacoEditor = dynamic(() => import('@/components/base/monaco').then((m) => m.WorkerEditor), {
  ssr: false,
})

interface WorkerCodeEditorProps {
  code: string
  onChange: (code: string) => void
}

export function WorkerCodeEditor({ code, onChange }: WorkerCodeEditorProps) {
  return (
    <div className="h-full w-full">
      <MonacoEditor code={code} onChange={onChange} />
    </div>
  )
}
