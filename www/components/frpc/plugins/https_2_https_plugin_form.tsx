'use client'

import { HTTPS2HTTPPluginOptions, HTTPS2HTTPSPluginOptions } from '@/types/plugin'
import { HTTPS2HTTPPluginForm } from './https_2_http_plugin_form'

interface Props {
  config: HTTPS2HTTPSPluginOptions
  setConfig: (c: HTTPS2HTTPSPluginOptions) => void
}

export function HTTPS2HTTPSPluginForm({ config, setConfig }: Props) {
  const { type, ...rest } = config
  const setTyped = (c: HTTPS2HTTPPluginOptions) => {
    setConfig({
      ...rest,
      ...c,
      type: 'https2https',
    })
  }
  return (
    <div className="space-y-4">
      {/* same as HTTPS2HTTP + additional fields */}
      <HTTPS2HTTPPluginForm
        config={{
          type: 'https2http',
          ...rest,
        }}
        setConfig={setTyped}
      />
    </div>
  )
}
