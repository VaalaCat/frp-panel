'use client'

import React, { useLayoutEffect } from 'react'
import Editor, { loader } from '@monaco-editor/react'

loader.config({
  paths: {
    vs: 'https://fastly.jsdelivr.net/npm/monaco-editor@0.36.1/min/vs',
  },
})

export interface WorkerEditorProps {
  code: string
  onChange: (value: string) => void
}

export function WorkerEditor({ code, onChange }: WorkerEditorProps) {
  // 强制客户端渲染
  useLayoutEffect(() => {}, [])

  return (
    <div className="h-full">
      <Editor height="100%" defaultLanguage="javascript" value={code} onChange={(v) => onChange(v ?? '')} />
    </div>
  )
}
