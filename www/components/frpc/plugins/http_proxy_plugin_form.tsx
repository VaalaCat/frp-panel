'use client'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { HTTPProxyPluginOptions } from '@/types/plugin'
import { useTranslation } from 'react-i18next'

interface Props {
  config: HTTPProxyPluginOptions
  setConfig: (c: HTTPProxyPluginOptions) => void
}

export function HTTPProxyPluginForm({ config, setConfig }: Props) {
  const { t } = useTranslation()

  return (
    <div className="space-y-4">
      <div>
        <Label htmlFor="httpUser">{t('frpc.client_plugins.http_user')}</Label>
        <Input
          id="httpUser"
          value={config.httpUser ?? ''}
          onChange={(e) => setConfig({ ...config, httpUser: e.target.value })}
          placeholder="username"
        />
      </div>
      <div>
        <Label htmlFor="httpPassword">{t('frpc.client_plugins.http_password')}</Label>
        <Input
          id="httpPassword"
          type="password"
          value={config.httpPassword ?? ''}
          onChange={(e) => setConfig({ ...config, httpPassword: e.target.value })}
          placeholder="password"
        />
      </div>
    </div>
  )
}
