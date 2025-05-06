'use client'

import React, { useLayoutEffect } from 'react'
import Editor, { loader } from '@monaco-editor/react'

loader.config({
  paths: {
    vs: 'https://fastly.jsdelivr.net/npm/monaco-editor@0.36.1/min/vs',
  },
})

export interface TemplateEditorProps {
  content: string
  onChange: (value: string) => void
}

export function TemplateEditor({ content, onChange }: TemplateEditorProps) {
  useLayoutEffect(() => {}, [])

  return (
    <div className="h-full">
      <Editor height="100%" defaultLanguage="capnp" value={content} onChange={(v) => onChange(v ?? '')} />
    </div>
  )
}
