import { HTTPProxyConfig, TCPProxyConfig, TypedProxyConfig, UDPProxyConfig, STCPProxyConfig } from '@/types/proxy'
import * as z from 'zod'
import React from 'react'
import { TypedProxyConfigValid, ZodPortSchema, ZodStringOptionalSchema, ZodStringSchema } from '@/lib/consts'
import { useEffect, useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { Button } from '@/components/ui/button'
import { Form, FormDescription } from '@/components/ui/form'
import { Label } from '@/components/ui/label'
import { YesIcon } from '@/components/ui/icon'
import { useTranslation } from 'react-i18next'
import { useQuery } from '@tanstack/react-query'
import { getServer } from '@/api/server'
import { Switch } from '@/components/ui/switch'
import { VisitPreview } from '../base/visit-preview'
import { HostField, PortField, SecretStringField, StringArrayField, StringField } from '../base/form-field'
import PluginConfigForm from './client_plugins'
import { TypedClientPluginOptions } from '@/types/plugin'
import { toast } from 'sonner'

export const TCPConfigSchema = z.object({
  remotePort: ZodPortSchema.optional(),
  localIP: ZodStringSchema.default('127.0.0.1').optional(),
  localPort: ZodPortSchema.optional(),
})

export const UDPConfigSchema = z.object({
  remotePort: ZodPortSchema.optional(),
  localIP: ZodStringSchema.default('127.0.0.1').optional(),
  localPort: ZodPortSchema.optional(),
})

export const HTTPConfigSchema = z.object({
  localPort: ZodPortSchema.optional(),
  localIP: ZodStringSchema.default('127.0.0.1').optional(),
  subdomain: ZodStringOptionalSchema,
  locations: z.array(ZodStringSchema).optional(),
  customDomains: z.array(ZodStringSchema).optional(),
  httpUser: ZodStringOptionalSchema,
  httpPassword: ZodStringOptionalSchema,
})

export const STCPConfigSchema = z.object({
  localIP: ZodStringSchema.default('127.0.0.1').optional(),
  localPort: ZodPortSchema.optional(),
  secretKey: ZodStringSchema.optional(),
})

export interface ProxyFormProps {
  clientID: string
  serverID: string
  proxyName: string
  enablePreview?: boolean
  defaultProxyConfig?: TypedProxyConfig
  clientProxyConfigs: TypedProxyConfig[]
  setClientProxyConfigs: React.Dispatch<React.SetStateAction<TypedProxyConfig[]>>
}

export const TypedProxyForm: React.FC<ProxyFormProps> = ({
  serverID,
  clientID,
  defaultProxyConfig,
  proxyName,
  clientProxyConfigs,
  setClientProxyConfigs,
  enablePreview,
}) => {
  if (!defaultProxyConfig) {
    return <></>
  }

  return (
    <>
      {defaultProxyConfig.type === 'tcp' && serverID && clientID && (
        <TCPProxyForm
          defaultProxyConfig={defaultProxyConfig}
          proxyName={proxyName}
          serverID={serverID}
          clientID={clientID}
          clientProxyConfigs={clientProxyConfigs}
          setClientProxyConfigs={setClientProxyConfigs}
          enablePreview={enablePreview}
        />
      )}
      {defaultProxyConfig.type === 'udp' && serverID && clientID && (
        <UDPProxyForm
          defaultProxyConfig={defaultProxyConfig}
          proxyName={proxyName}
          serverID={serverID}
          clientID={clientID}
          clientProxyConfigs={clientProxyConfigs}
          setClientProxyConfigs={setClientProxyConfigs}
          enablePreview={enablePreview}
        />
      )}
      {defaultProxyConfig.type === 'http' && serverID && clientID && (
        <HTTPProxyForm
          defaultProxyConfig={defaultProxyConfig}
          proxyName={proxyName}
          serverID={serverID}
          clientID={clientID}
          clientProxyConfigs={clientProxyConfigs}
          setClientProxyConfigs={setClientProxyConfigs}
          enablePreview={enablePreview}
        />
      )}
      {defaultProxyConfig.type === 'stcp' && serverID && clientID && (
        <STCPProxyForm
          defaultProxyConfig={defaultProxyConfig}
          proxyName={proxyName}
          serverID={serverID}
          clientID={clientID}
          clientProxyConfigs={clientProxyConfigs}
          setClientProxyConfigs={setClientProxyConfigs}
          enablePreview={enablePreview}
        />
      )}
    </>
  )
}

export const TCPProxyForm: React.FC<ProxyFormProps> = ({
  serverID,
  clientID,
  defaultProxyConfig,
  proxyName,
  clientProxyConfigs,
  setClientProxyConfigs,
  enablePreview,
}) => {
  const defaultConfig = defaultProxyConfig as TCPProxyConfig
  const [timeoutID, setTimeoutID] = useState<NodeJS.Timeout | undefined>()
  const form = useForm<z.infer<typeof TCPConfigSchema>>({
    resolver: zodResolver(TCPConfigSchema),
    defaultValues: {
      remotePort: defaultConfig?.remotePort,
      localIP: defaultConfig?.localIP,
      localPort: defaultConfig?.localPort,
    },
  })

  const [usePlugin, setUsePlugin] = useState<boolean>(
    (defaultConfig.plugin && defaultConfig.plugin.type.length > 0) || false,
  )
  const [pluginConfig, setPluginConfig] = useState<TypedClientPluginOptions | undefined>(defaultConfig.plugin)

  const onSubmit = async (values: z.infer<typeof TCPConfigSchema>) => {
    const cfgToSubmit = { ...values, plugin: pluginConfig, type: 'tcp', name: proxyName } as TCPProxyConfig
    if (!TypedProxyConfigValid(cfgToSubmit)) {
      toast.error('Invalid configuration')
      return
    }
    handleSave()
    const newProxiyConfigs = clientProxyConfigs.map((proxyCfg) => {
      if (proxyCfg.name === proxyName) {
        return cfgToSubmit
      }
      return proxyCfg
    })
    setClientProxyConfigs(newProxiyConfigs)
  }

  const [isSaveDisabled, setSaveDisabled] = useState(false)

  const handleSave = () => {
    setSaveDisabled(true)
    if (timeoutID) {
      clearTimeout(timeoutID)
    }
    setTimeoutID(
      setTimeout(() => {
        setSaveDisabled(false)
      }, 3000),
    )
  }

  const { t } = useTranslation()

  const { data: server } = useQuery({
    queryKey: ['getServer', serverID],
    queryFn: () => {
      return getServer({ serverId: serverID })
    },
  })

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 px-0.5">
        {server?.server?.ip &&
          defaultConfig.remotePort &&
          defaultConfig.localIP &&
          defaultConfig.localPort &&
          enablePreview && (
            <div className="flex items-center space-x-2 flex-col justify-start w-full">
              <Label className="text-sm font-medium text-start w-full">{t('proxy.form.access_method')}</Label>
              <div className="w-full justify-start overflow-x-scroll">
                <VisitPreview server={server?.server} typedProxyConfig={defaultConfig} />
              </div>
            </div>
          )}
        <PortField name="localPort" control={form.control} label={t('proxy.form.local_port') + '*'} />
        <HostField name="localIP" control={form.control} label={t('proxy.form.local_ip') + '*'} />
        <PortField
          name="remotePort"
          control={form.control}
          label={t('proxy.form.remote_port') + '*'}
          placeholder="4321"
        />
        <SwitchWithLabel
          name="usePlugin"
          label={t('proxy.form.use_plugin')}
          defaultValue={usePlugin}
          setValue={setUsePlugin}
        />
        {usePlugin ? (
          <PluginConfigForm defaultPluginConfig={defaultConfig.plugin} setPluginConfig={setPluginConfig} />
        ) : null}
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'} className="w-full">
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}

export const STCPProxyForm: React.FC<ProxyFormProps> = ({
  serverID,
  clientID,
  defaultProxyConfig,
  proxyName,
  clientProxyConfigs,
  setClientProxyConfigs,
  enablePreview,
}) => {
  const defaultConfig = defaultProxyConfig as STCPProxyConfig
  const [timeoutID, setTimeoutID] = useState<NodeJS.Timeout | undefined>()
  const form = useForm<z.infer<typeof STCPConfigSchema>>({
    resolver: zodResolver(STCPConfigSchema),
    defaultValues: {
      localPort: defaultConfig?.localPort,
      localIP: defaultConfig?.localIP,
      secretKey: defaultConfig?.secretKey,
    },
  })

  const [usePlugin, setUsePlugin] = useState<boolean>(
    (defaultConfig.plugin && defaultConfig.plugin.type.length > 0) || false,
  )

  const [pluginConfig, setPluginConfig] = useState<TypedClientPluginOptions | undefined>(defaultConfig.plugin)

  const onSubmit = async (values: z.infer<typeof STCPConfigSchema>) => {
    const cfgToSubmit = { ...values, plugin: pluginConfig, type: 'stcp', name: proxyName } as STCPProxyConfig
    if (!TypedProxyConfigValid(cfgToSubmit)) {
      toast.error('Invalid configuration')
      return
    }
    handleSave()
    const newProxiyConfigs = clientProxyConfigs.map((proxyCfg) => {
      if (proxyCfg.name === proxyName) {
        return cfgToSubmit
      }
      return proxyCfg
    })
    setClientProxyConfigs(newProxiyConfigs)
  }

  const [isSaveDisabled, setSaveDisabled] = useState(false)

  const handleSave = () => {
    setSaveDisabled(true)
    if (timeoutID) {
      clearTimeout(timeoutID)
    }
    setTimeoutID(
      setTimeout(() => {
        setSaveDisabled(false)
      }, 3000),
    )
  }

  const { t } = useTranslation()

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 px-0.5">
        <PortField name="localPort" control={form.control} label={t('proxy.form.local_port') + '*'} />
        <HostField name="localIP" control={form.control} label={t('proxy.form.local_ip') + '*'} />
        <SwitchWithLabel
          name="usePlugin"
          defaultValue={usePlugin}
          setValue={setUsePlugin}
          label={t('proxy.form.use_plugin')}
        />
        {usePlugin ? (
          <PluginConfigForm defaultPluginConfig={defaultConfig.plugin} setPluginConfig={setPluginConfig} />
        ) : null}
        <SecretStringField name="secretKey" control={form.control} label={t('proxy.form.secret_key') + '*'} />
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'} className="w-full">
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}

export const UDPProxyForm: React.FC<ProxyFormProps> = ({
  serverID,
  clientID,
  defaultProxyConfig,
  proxyName,
  clientProxyConfigs,
  setClientProxyConfigs,
  enablePreview,
}) => {
  const defaultConfig = defaultProxyConfig as UDPProxyConfig
  const [timeoutID, setTimeoutID] = useState<NodeJS.Timeout | undefined>()
  const form = useForm<z.infer<typeof UDPConfigSchema>>({
    resolver: zodResolver(UDPConfigSchema),
    defaultValues: {
      localPort: defaultConfig?.localPort,
      localIP: defaultConfig?.localIP,
      remotePort: defaultConfig?.remotePort,
    },
  })

  const [usePlugin, setUsePlugin] = useState<boolean>(
    (defaultConfig.plugin && defaultConfig.plugin.type.length > 0) || false,
  )

  const [pluginConfig, setPluginConfig] = useState<TypedClientPluginOptions | undefined>(defaultConfig.plugin)

  const onSubmit = async (values: z.infer<typeof UDPConfigSchema>) => {
    const cfgToSubmit = { ...values, plugin: pluginConfig, type: 'udp', name: proxyName } as UDPProxyConfig
    if (!TypedProxyConfigValid(cfgToSubmit)) {
      toast.error('Invalid configuration')
      return
    }
    handleSave()
    const newProxiyConfigs = clientProxyConfigs.map((proxyCfg) => {
      if (proxyCfg.name === proxyName) {
        return cfgToSubmit
      }
      return proxyCfg
    })
    setClientProxyConfigs(newProxiyConfigs)
  }

  const [isSaveDisabled, setSaveDisabled] = useState(false)

  const handleSave = () => {
    setSaveDisabled(true)
    if (timeoutID) {
      clearTimeout(timeoutID)
    }
    setTimeoutID(
      setTimeout(() => {
        setSaveDisabled(false)
      }, 3000),
    )
  }

  const { t } = useTranslation()

  const { data: server } = useQuery({
    queryKey: ['getServer', serverID],
    queryFn: () => {
      return getServer({ serverId: serverID })
    },
  })

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 px-0.5">
        {server?.server?.ip &&
          defaultConfig.remotePort &&
          defaultConfig.localIP &&
          defaultConfig.localPort &&
          enablePreview && (
            <div className="flex items-center space-x-2 flex-col justify-start w-full">
              <Label className="text-sm font-medium text-start w-full">{t('proxy.form.access_method')}</Label>
              <div className="w-full justify-start overflow-x-scroll">
                <VisitPreview server={server?.server} typedProxyConfig={defaultConfig} />
              </div>
            </div>
          )}
        <PortField name="localPort" control={form.control} label={t('proxy.form.local_port') + '*'} />
        <HostField name="localIP" control={form.control} label={t('proxy.form.local_ip') + '*'} />
        <PortField name="remotePort" control={form.control} label={t('proxy.form.remote_port') + '*'} />
        <SwitchWithLabel
          name="usePlugin"
          defaultValue={usePlugin}
          setValue={setUsePlugin}
          label={t('proxy.form.use_plugin')}
        />
        {usePlugin ? (
          <PluginConfigForm defaultPluginConfig={defaultConfig.plugin} setPluginConfig={setPluginConfig} />
        ) : null}
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'} className="w-full">
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}

export const HTTPProxyForm: React.FC<ProxyFormProps> = ({
  serverID,
  clientID,
  defaultProxyConfig,
  proxyName,
  clientProxyConfigs,
  setClientProxyConfigs,
  enablePreview,
}) => {
  const defaultConfig = defaultProxyConfig as HTTPProxyConfig
  const [timeoutID, setTimeoutID] = useState<NodeJS.Timeout | undefined>()
  const [moreSettings, setMoreSettings] = useState(false)
  const [useAuth, setUseAuth] = useState(false)

  const form = useForm<z.infer<typeof HTTPConfigSchema>>({
    resolver: zodResolver(HTTPConfigSchema),
    defaultValues: {
      localIP: defaultConfig?.localIP,
      localPort: defaultConfig?.localPort,
      subdomain: defaultConfig?.subdomain,
      locations: defaultConfig?.locations,
      customDomains: defaultConfig?.customDomains,
      httpPassword: defaultConfig?.httpPassword,
      httpUser: defaultConfig?.httpUser,
    },
  })

  const [usePlugin, setUsePlugin] = useState<boolean>(
    (defaultConfig.plugin && defaultConfig.plugin.type.length > 0) || false,
  )

  const [pluginConfig, setPluginConfig] = useState<TypedClientPluginOptions | undefined>(defaultConfig.plugin)

  const onSubmit = async (values: z.infer<typeof HTTPConfigSchema>) => {
    const cfgToSubmit = { ...values, plugin: pluginConfig, type: 'http', name: proxyName } as HTTPProxyConfig
    if (!TypedProxyConfigValid(cfgToSubmit)) {
      toast.error('Invalid configuration')
      return
    }
    if (!values.customDomains && !values.subdomain) {
      toast.error('Please provide a subdomain or custom domains')
      return
    }
    handleSave()
    const newProxiyConfigs = clientProxyConfigs.map((proxyCfg) => {
      if (proxyCfg.name === proxyName) {
        return cfgToSubmit
      }
      return proxyCfg
    })
    setClientProxyConfigs(newProxiyConfigs)
  }

  const [isSaveDisabled, setSaveDisabled] = useState(false)

  useEffect(() => {
    if (defaultConfig?.httpPassword || defaultConfig?.httpUser) {
      setUseAuth(true)
    }
  }, [defaultConfig?.httpPassword, defaultConfig?.httpUser])

  const handleSave = () => {
    setSaveDisabled(true)
    if (timeoutID) {
      clearTimeout(timeoutID)
    }
    setTimeoutID(
      setTimeout(() => {
        setSaveDisabled(false)
      }, 3000),
    )
  }

  const { t } = useTranslation()

  const { data: server } = useQuery({
    queryKey: ['getServer', serverID],
    queryFn: () => {
      return getServer({ serverId: serverID })
    },
  })

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 px-0.5">
        {server &&
          server.server &&
          server.server.ip &&
          defaultConfig &&
          defaultConfig.localIP &&
          defaultConfig.localPort &&
          defaultConfig.subdomain &&
          enablePreview && (
            <div className="flex items-center space-x-2 flex-col justify-start w-full">
              <Label className="text-sm font-medium text-start w-full">{t('proxy.form.access_method')}</Label>
              <div className="w-full justify-start overflow-x-scroll">
                <VisitPreview server={server?.server} typedProxyConfig={defaultConfig} />
              </div>
            </div>
          )}
        <PortField name="localPort" control={form.control} label={t('proxy.form.local_port') + '*'} />
        <HostField name="localIP" control={form.control} label={t('proxy.form.local_ip') + '*'} />
        <StringField
          name="subdomain"
          control={form.control}
          label={t('proxy.form.subdomain')}
          placeholder={'your_sub_domain'}
        />
        <StringArrayField
          name="customDomains"
          control={form.control}
          label={t('proxy.form.custom_domains')}
          placeholder={'your.example.com'}
        />
        <FormDescription>{t('proxy.form.domain_description')}</FormDescription>
        <SwitchWithLabel
          name="usePlugin"
          defaultValue={usePlugin}
          setValue={setUsePlugin}
          label={t('proxy.form.use_plugin')}
        />
        {usePlugin ? (
          <PluginConfigForm defaultPluginConfig={defaultConfig.plugin} setPluginConfig={setPluginConfig} />
        ) : null}
        <SwitchWithLabel
          name="moreSettings"
          label={t('proxy.form.more_settings')}
          defaultValue={moreSettings}
          setValue={setMoreSettings}
        />
        {moreSettings && (
          <div className="p-4 space-y-4 border rounded-md">
            <StringArrayField
              name="locations"
              control={form.control}
              label={t('proxy.form.route')}
              placeholder={'/path'}
            />
            <SwitchWithLabel
              name="enableHttpAuth"
              label={t('proxy.form.enable_http_auth')}
              defaultValue={useAuth}
              setValue={setUseAuth}
            />
            {useAuth && (
              <div className="p-4 space-y-4 border rounded-md">
                <StringField
                  name="httpUser"
                  control={form.control}
                  label={t('proxy.form.username')}
                  placeholder={'username'}
                />
                <StringField
                  name="httpPassword"
                  control={form.control}
                  label={t('proxy.form.password')}
                  placeholder={'password'}
                />
              </div>
            )}
          </div>
        )}
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'} className="w-full">
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}

const SwitchWithLabel = ({
  name,
  label,
  defaultValue,
  setValue,
}: {
  name: string
  label: string
  defaultValue?: boolean
  setValue: (value: boolean) => void
}) => {
  const { t } = useTranslation()
  return (
    <div className="flex items-center space-x-2 justify-between">
      <Label htmlFor={name}>{t(label)}</Label>
      <Switch id={`switch-with-label-${name}-switch`} checked={defaultValue} onCheckedChange={setValue} />
    </div>
  )
}
