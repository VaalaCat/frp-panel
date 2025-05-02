'use client'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { HTTP2HTTPSPluginOptions } from '@/types/plugin'
import { useTranslation } from 'react-i18next'

interface Props {
  config: HTTP2HTTPSPluginOptions
  setConfig: (c: HTTP2HTTPSPluginOptions) => void
}

export function HTTP2HTTPSPluginForm({ config, setConfig }: Props) {
  const { t } = useTranslation()

  return (
    <div className="space-y-4">
      <div>
        <Label htmlFor="localAddr">{t('frpc.client_plugins.http_local_addr')}</Label>
        <Input
          id="localAddr"
          value={config.localAddr ?? ''}
          onChange={(e) => setConfig({ ...config, localAddr: e.target.value })}
          placeholder="127.0.0.1:8080"
        />
      </div>
      <div>
        <Label htmlFor="hostHeaderRewrite">{t('frpc.client_plugins.http_host_header_rewrite')}</Label>
        <Input
          id="hostHeaderRewrite"
          value={config.hostHeaderRewrite ?? ''}
          onChange={(e) => setConfig({ ...config, hostHeaderRewrite: e.target.value })}
          placeholder="example.com"
        />
      </div>
      {/* You could add a custom HeaderOperations component here */}
    </div>
  )
}
