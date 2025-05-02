'use client'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { UnixDomainSocketPluginOptions } from '@/types/plugin'
import { useTranslation } from 'react-i18next'

interface Props {
  config: UnixDomainSocketPluginOptions
  setConfig: (c: UnixDomainSocketPluginOptions) => void
}

export function UnixDomainSocketPluginForm({ config, setConfig }: Props) {
  const { t } = useTranslation()

  return (
    <div>
      <Label htmlFor="unixPath">{t('frpc.plugins.unix_domain_socket_path')}</Label>
      <Input
        id="unixPath"
        value={config.unixPath ?? ''}
        onChange={(e) => setConfig({ ...config, unixPath: e.target.value })}
        placeholder="/tmp/frp.sock"
      />
    </div>
  )
}
