'use client'

import { useState, useEffect, useMemo } from 'react'

import {
  TypedClientPluginOptions,
  ClientPluginType,
  HTTPProxyPluginOptions,
  HTTP2HTTPSPluginOptions,
  HTTPS2HTTPPluginOptions,
  HTTPS2HTTPSPluginOptions,
  Socks5PluginOptions,
  StaticFilePluginOptions,
  UnixDomainSocketPluginOptions,
} from '@/types/plugin'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { HTTPProxyPluginForm } from './plugins/http_proxy_plugin_form'
import { HTTP2HTTPSPluginForm } from './plugins/http_2_https_plugin_form'
import { HTTPS2HTTPPluginForm } from './plugins/https_2_http_plugin_form'
import { HTTPS2HTTPSPluginForm } from './plugins/https_2_https_plugin_form'
import { Socks5PluginForm } from './plugins/socks5_plugin_form'
import { StaticFilePluginForm } from './plugins/static_file_plugin_form'
import { UnixDomainSocketPluginForm } from './plugins/unix_domain_socket_plugin_form'
import { useTranslation } from 'react-i18next'

const pluginTypeMap: Record<string, string> = {
  http_proxy: 'HTTP Proxy',
  http2https: 'HTTP→HTTPS',
  https2http: 'HTTPS→HTTP',
  https2https: 'HTTPS→HTTPS',
  socks5: 'SOCKS5',
  static_file: 'Static File',
  unix_domain_socket: 'Unix Domain Socket',
}

interface Props {
  defaultPluginConfig?: TypedClientPluginOptions
  setPluginConfig: (c: TypedClientPluginOptions) => void
  supportedPlugins?: ClientPluginType[]
}

export function PluginConfigForm({ defaultPluginConfig, setPluginConfig, supportedPlugins }: Props) {
  const { t } = useTranslation()
  const [config, updateConfig] = useState(defaultPluginConfig)

  // sync out
  useEffect(() => {
    config && setPluginConfig(config)
  }, [config, setPluginConfig])

  const handleTypeChange = (type: ClientPluginType) => {
    updateConfig({ ...config, type })
  }

  const pluginsToSelect = useMemo(() => {
    const plugins = supportedPlugins || Object.keys(pluginTypeMap)
    return (
      plugins &&
      plugins.map((plugin, i) => (
        <SelectItem key={`${i}`} value={plugin}>
          {pluginTypeMap[plugin]}
        </SelectItem>
      ))
    )
  }, [supportedPlugins])

  return (
    <div className="space-y-6 p-6 bg-white rounded-lg shadow">
      <div>
        <Label htmlFor="plugin-type">{t('frpc.client_plugins.plugin_type')}</Label>
        <Select
          defaultValue={defaultPluginConfig?.type}
          onValueChange={(value) => handleTypeChange(value as ClientPluginType)}
        >
          <SelectTrigger id="plugin-type">
            <SelectValue placeholder={t('frpc.client_plugins.select_plugin_type')} />
          </SelectTrigger>
          <SelectContent>{pluginsToSelect}</SelectContent>
        </Select>
      </div>

      {config && config.type && config.type.length > 0 ? (
        <>
          {config.type === 'http_proxy' && (
            <HTTPProxyPluginForm config={config as HTTPProxyPluginOptions} setConfig={(c) => updateConfig(c)} />
          )}
          {config.type === 'http2https' && (
            <HTTP2HTTPSPluginForm config={config as HTTP2HTTPSPluginOptions} setConfig={(c) => updateConfig(c)} />
          )}
          {config.type === 'https2http' && (
            <HTTPS2HTTPPluginForm config={config as HTTPS2HTTPPluginOptions} setConfig={(c) => updateConfig(c)} />
          )}
          {config.type === 'https2https' && (
            <HTTPS2HTTPSPluginForm config={config as HTTPS2HTTPSPluginOptions} setConfig={(c) => updateConfig(c)} />
          )}
          {config.type === 'socks5' && (
            <Socks5PluginForm config={config as Socks5PluginOptions} setConfig={(c) => updateConfig(c)} />
          )}
          {config.type === 'static_file' && (
            <StaticFilePluginForm config={config as StaticFilePluginOptions} setConfig={(c) => updateConfig(c)} />
          )}
          {config.type === 'unix_domain_socket' && (
            <UnixDomainSocketPluginForm
              config={config as UnixDomainSocketPluginOptions}
              setConfig={(c) => updateConfig(c)}
            />
          )}
        </>
      ) : null}
    </div>
  )
}

export default PluginConfigForm
