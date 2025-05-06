'use client'

import React from 'react'
import dynamic from 'next/dynamic'

// 动态加载模板编辑器组件
const TemplateEditorComponent = dynamic(
  () => import('@/components/worker/template_edit').then((m) => m.TemplateEditor),
  { ssr: false },
)

interface WorkerTemplateEditorProps {
  content: string
  onChange: (content: string) => void
}

export function WorkerTemplateEditor({ content, onChange }: WorkerTemplateEditorProps) {
  return (
    <div className="h-full w-full">
      <TemplateEditorComponent content={content} onChange={onChange} />
    </div>
  )
}
