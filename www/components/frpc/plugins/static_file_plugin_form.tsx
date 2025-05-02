'use client'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { StaticFilePluginOptions } from '@/types/plugin'
import { useTranslation } from 'react-i18next'

interface Props {
  config: StaticFilePluginOptions
  setConfig: (c: StaticFilePluginOptions) => void
}

export function StaticFilePluginForm({ config, setConfig }: Props) {
  const { t } = useTranslation()

  return (
    <div className="space-y-4">
      <div>
        <Label htmlFor="localPath">{t('frpc.plugins.local_path')}</Label>
        <Input
          id="localPath"
          value={config.localPath ?? ''}
          onChange={(e) => setConfig({ ...config, localPath: e.target.value })}
          placeholder="/var/www"
        />
      </div>
      <div>
        <Label htmlFor="stripPrefix">{t('frpc.plugins.strip_prefix')}</Label>
        <Input
          id="stripPrefix"
          value={config.stripPrefix ?? ''}
          onChange={(e) => setConfig({ ...config, stripPrefix: e.target.value })}
          placeholder="/static"
        />
      </div>
      <div>
        <Label htmlFor="httpUser">{t('frpc.plugins.http_user')}</Label>
        <Input
          id="httpUser"
          value={config.httpUser ?? ''}
          onChange={(e) => setConfig({ ...config, httpUser: e.target.value })}
        />
      </div>
      <div>
        <Label htmlFor="httpPassword">{t('frpc.plugins.http_password')}</Label>
        <Input
          id="httpPassword"
          type="password"
          value={config.httpPassword ?? ''}
          onChange={(e) => setConfig({ ...config, httpPassword: e.target.value })}
        />
      </div>
    </div>
  )
}
