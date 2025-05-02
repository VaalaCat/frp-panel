'use client'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Socks5PluginOptions } from '@/types/plugin'
import { useTranslation } from 'react-i18next'

interface Props {
  config: Socks5PluginOptions
  setConfig: (c: Socks5PluginOptions) => void
}

export function Socks5PluginForm({ config, setConfig }: Props) {
  const { t } = useTranslation()
  return (
    <div className="space-y-4">
      <div>
        <Label htmlFor="username">{t('frpc.plugins.username')}</Label>
        <Input
          id="username"
          value={config.username ?? ''}
          onChange={(e) => setConfig({ ...config, username: e.target.value })}
          placeholder="username"
        />
      </div>
      <div>
        <Label htmlFor="password">{t('frpc.plugins.password')}</Label>
        <Input
          id="password"
          type="password"
          value={config.password ?? ''}
          onChange={(e) => setConfig({ ...config, password: e.target.value })}
          placeholder="password"
        />
      </div>
    </div>
  )
}
