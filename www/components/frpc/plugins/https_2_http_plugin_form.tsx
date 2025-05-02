'use client'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { HTTPS2HTTPPluginOptions } from '@/types/plugin'
import { useTranslation } from 'react-i18next'

interface Props {
  config: HTTPS2HTTPPluginOptions
  setConfig: (c: HTTPS2HTTPPluginOptions) => void
}

export function HTTPS2HTTPPluginForm({ config, setConfig }: Props) {
  const { t } = useTranslation()

  return (
    <div className="space-y-4">
      <div>
        <Label htmlFor="localAddr">{t('frpc.plugins.local_addr')}</Label>
        <Input
          id="localAddr"
          value={config.localAddr ?? ''}
          onChange={(e) => setConfig({ ...config, localAddr: e.target.value })}
          placeholder="127.0.0.1:8080"
        />
      </div>
      <div>
        <Label htmlFor="hostHeaderRewrite">{t('frpc.plugins.host_header_rewrite')}</Label>
        <Input
          id="hostHeaderRewrite"
          value={config.hostHeaderRewrite ?? ''}
          onChange={(e) => setConfig({ ...config, hostHeaderRewrite: e.target.value })}
          placeholder="example.com"
        />
      </div>
      <div>
        <Label htmlFor="crtPath">{t('frpc.plugins.crt_path')}</Label>
        <Input
          id="crtPath"
          value={config.crtPath ?? ''}
          onChange={(e) => setConfig({ ...config, crtPath: e.target.value })}
          placeholder="/path/to/cert.pem"
        />
      </div>
      <div>
        <Label htmlFor="keyPath">{t('frpc.plugins.key_path')}</Label>
        <Input
          id="keyPath"
          value={config.keyPath ?? ''}
          onChange={(e) => setConfig({ ...config, keyPath: e.target.value })}
          placeholder="/path/to/key.pem"
        />
      </div>
    </div>
  )
}
