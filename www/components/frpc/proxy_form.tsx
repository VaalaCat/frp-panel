import { HTTPProxyConfig, TCPProxyConfig, TypedProxyConfig, UDPProxyConfig, STCPProxyConfig } from '@/types/proxy'
import * as z from 'zod'
import React from 'react'
import { ZodPortSchema, ZodStringOptionalSchema, ZodStringSchema } from '@/lib/consts'
import { useEffect, useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { Control, useForm } from 'react-hook-form'
import { Button } from '@/components/ui/button'
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { YesIcon } from '@/components/ui/icon'
import { useTranslation } from 'react-i18next'
import { useQuery } from '@tanstack/react-query'
import { getServer } from '@/api/server'
import { Switch } from "@/components/ui/switch"
import { VisitPreview } from '../base/visit-preview'
import StringListInput from '../base/list-input'

export const TCPConfigSchema = z.object({
  remotePort: ZodPortSchema,
  localIP: ZodStringSchema.default('127.0.0.1'),
  localPort: ZodPortSchema,
})

export const UDPConfigSchema = z.object({
  remotePort: ZodPortSchema.optional(),
  localIP: ZodStringSchema.default('127.0.0.1'),
  localPort: ZodPortSchema,
})

export const HTTPConfigSchema = z.object({
  localPort: ZodPortSchema,
  localIP: ZodStringSchema.default('127.0.0.1'),
  subdomain: ZodStringOptionalSchema,
  locations: z.array(ZodStringSchema).optional(),
  customDomains: z.array(ZodStringSchema).optional(),
  httpUser: ZodStringOptionalSchema,
  httpPassword: ZodStringOptionalSchema,
})

export const STCPConfigSchema = z.object({
  localIP: ZodStringSchema.default('127.0.0.1'),
  localPort: ZodPortSchema,
  secretKey: ZodStringSchema,
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

const HostField = ({
  control,
  name,
  label,
  placeholder,
  defaultValue,
}: {
  control: Control<any>
  name: string
  label: string
  placeholder?: string
  defaultValue?: string
}) => {
  const { t } = useTranslation()
  return (
    <FormField
      name={name}
      control={control}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{t(label)}</FormLabel>
          <FormControl>
            <Input className='text-sm' placeholder={placeholder || '127.0.0.1'} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
      defaultValue={defaultValue}
    />
  )
}
const PortField = ({
  control,
  name,
  label,
  placeholder,
  defaultValue,
}: {
  control: Control<any>
  name: string
  label: string
  placeholder?: string
  defaultValue?: number
}) => {
  const { t } = useTranslation()
  return (
    <FormField
      name={name}
      control={control}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{t(label)}</FormLabel>
          <FormControl>
            <Input className='text-sm' placeholder={placeholder || '1234'} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
      defaultValue={defaultValue}
    />
  )
}
const SecretStringField = ({
  control,
  name,
  label,
  placeholder,
  defaultValue,
}: {
  control: Control<any>
  name: string
  label: string
  placeholder?: string
  defaultValue?: string
}) => {
  const { t } = useTranslation()
  return (
    <FormField
      name={name}
      control={control}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{t(label)}</FormLabel>
          <FormControl>
            <Input className='text-sm' placeholder={placeholder || "secret"} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
      defaultValue={defaultValue}
    />
  )
}

const StringField = ({
  control,
  name,
  label,
  placeholder,
  defaultValue,
}: {
  control: Control<any>
  name: string
  label: string
  placeholder?: string
  defaultValue?: string
}) => {
  const { t } = useTranslation()
  return (
    <FormField
      name={name}
      control={control}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{t(label)}</FormLabel>
          <FormControl>
            <Input className='text-sm' placeholder={placeholder || '127.0.0.1'} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
      defaultValue={defaultValue}
    />
  )
}

const StringArrayField = ({
  control,
  name,
  label,
  placeholder,
  defaultValue,
}: {
  control: Control<any>
  name: string
  label: string
  placeholder?: string
  defaultValue?: string[]
}) => {
  const { t } = useTranslation()
  return (
    <FormField
      name={name}
      control={control}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{t(label)}</FormLabel>
          <FormControl>
            <StringListInput placeholder={placeholder || '/path'} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
      defaultValue={defaultValue}
    />
  )
}

export const TypedProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID, defaultProxyConfig, proxyName, clientProxyConfigs, setClientProxyConfigs, enablePreview }) => {
  if (!defaultProxyConfig) {
    return <></>
  }

  return (<> {defaultProxyConfig.type === 'tcp' && serverID && clientID && (
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
    )}</>)
}

export const TCPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID, defaultProxyConfig, proxyName, clientProxyConfigs, setClientProxyConfigs, enablePreview }) => {
  const defaultConfig = defaultProxyConfig as TCPProxyConfig
  const [_, setTCPConfig] = useState<TCPProxyConfig | undefined>()
  const [timeoutID, setTimeoutID] = useState<NodeJS.Timeout | undefined>()
  const form = useForm<z.infer<typeof TCPConfigSchema>>({
    resolver: zodResolver(TCPConfigSchema),
    defaultValues: {
      remotePort: defaultConfig?.remotePort,
      localIP: defaultConfig?.localIP,
      localPort: defaultConfig?.localPort,
    }
  })

  const onSubmit = async (values: z.infer<typeof TCPConfigSchema>) => {
    handleSave()
    setTCPConfig({ type: 'tcp', ...values, name: proxyName })
    const newProxiyConfigs = clientProxyConfigs.map((proxyCfg) => {
      if (proxyCfg.name === proxyName) {
        return { ...values, type: 'tcp', name: proxyName } as TCPProxyConfig
      }
      return proxyCfg
    })
    console.log('newProxiyConfigs', newProxiyConfigs)
    setClientProxyConfigs(newProxiyConfigs)
  }

  const [isSaveDisabled, setSaveDisabled] = useState(false)

  const handleSave = () => {
    setSaveDisabled(true)
    if (timeoutID) {
      clearTimeout(timeoutID)
    }
    setTimeoutID(setTimeout(() => {
      setSaveDisabled(false)
    }, 3000))
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
        {server?.server?.ip && defaultConfig.remotePort && defaultConfig.localIP && defaultConfig.localPort && enablePreview && (
          <div className="flex items-center space-x-2 flex-col justify-start w-full">
            <Label className="text-sm font-medium text-start w-full">{t('proxy.form.access_method')}</Label>
            <div className='w-full justify-start overflow-x-scroll'>
              <VisitPreview server={server?.server} typedProxyConfig={defaultConfig} />
            </div>
          </div>
        )}
        <PortField name="localPort" control={form.control} label={t('proxy.form.local_port') + "*"} />
        <HostField name="localIP" control={form.control} label={t('proxy.form.local_ip') + "*"} />
        <PortField name="remotePort" control={form.control} label={t('proxy.form.remote_port') + "*"} placeholder='4321' />
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'} className='w-full'>
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}

export const STCPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID, defaultProxyConfig, proxyName, clientProxyConfigs, setClientProxyConfigs, enablePreview }) => {
  const defaultConfig = defaultProxyConfig as STCPProxyConfig
  const [_, setSTCPConfig] = useState<STCPProxyConfig | undefined>()
  const [timeoutID, setTimeoutID] = useState<NodeJS.Timeout | undefined>()
  const form = useForm<z.infer<typeof STCPConfigSchema>>({
    resolver: zodResolver(STCPConfigSchema),
    defaultValues: {
      localPort: defaultConfig?.localPort,
      localIP: defaultConfig?.localIP,
      secretKey: defaultConfig?.secretKey,
    }
  })

  const onSubmit = async (values: z.infer<typeof STCPConfigSchema>) => {
    handleSave()
    setSTCPConfig({ type: 'stcp', ...values, name: proxyName })
    const newProxiyConfigs = clientProxyConfigs.map((proxyCfg) => {
      if (proxyCfg.name === proxyName) {
        return { ...values, type: 'stcp', name: proxyName } as STCPProxyConfig
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
    setTimeoutID(setTimeout(() => {
      setSaveDisabled(false)
    }, 3000))
  }

  const { t } = useTranslation()

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 px-0.5">
        <PortField name="localPort" control={form.control} label={t('proxy.form.local_port') + "*"} />
        <HostField name="localIP" control={form.control} label={t('proxy.form.local_ip') + "*"} />
        <SecretStringField name="secretKey" control={form.control} label={t('proxy.form.secret_key') + "*"} />
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'} className='w-full'>
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}

export const UDPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID, defaultProxyConfig, proxyName, clientProxyConfigs, setClientProxyConfigs, enablePreview }) => {
  const defaultConfig = defaultProxyConfig as UDPProxyConfig
  const [_, setUDPConfig] = useState<UDPProxyConfig | undefined>()
  const [timeoutID, setTimeoutID] = useState<NodeJS.Timeout | undefined>()
  const form = useForm<z.infer<typeof UDPConfigSchema>>({
    resolver: zodResolver(UDPConfigSchema),
    defaultValues: {
      localPort: defaultConfig?.localPort,
      localIP: defaultConfig?.localIP,
      remotePort: defaultConfig?.remotePort,
    }
  })

  const onSubmit = async (values: z.infer<typeof UDPConfigSchema>) => {
    handleSave()
    setUDPConfig({ type: 'udp', ...values, name: proxyName })
    const newProxiyConfigs = clientProxyConfigs.map((proxyCfg) => {
      if (proxyCfg.name === proxyName) {
        return { ...values, type: 'udp', name: proxyName } as UDPProxyConfig
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
    setTimeoutID(setTimeout(() => {
      setSaveDisabled(false)
    }, 3000))
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
        {server?.server?.ip && defaultConfig.remotePort && defaultConfig.localIP && defaultConfig.localPort && enablePreview && (
          <div className="flex items-center space-x-2 flex-col justify-start w-full">
            <Label className="text-sm font-medium text-start w-full">{t('proxy.form.access_method')}</Label>
            <div className='w-full justify-start overflow-x-scroll'>
              <VisitPreview server={server?.server} typedProxyConfig={defaultConfig} />
            </div>
          </div>
        )}
        <PortField name="localPort" control={form.control} label={t('proxy.form.local_port') + "*"} />
        <HostField name="localIP" control={form.control} label={t('proxy.form.local_ip') + "*"} />
        <PortField name="remotePort" control={form.control} label={t('proxy.form.remote_port') + "*"} />
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'} className='w-full'>
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}

export const HTTPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID, defaultProxyConfig, proxyName, clientProxyConfigs, setClientProxyConfigs, enablePreview }) => {
  const defaultConfig = defaultProxyConfig as HTTPProxyConfig
  const [_, setHTTPConfig] = useState<HTTPProxyConfig | undefined>()
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
      httpUser: defaultConfig?.httpUser
    }
  })

  const onSubmit = async (values: z.infer<typeof HTTPConfigSchema>) => {
    handleSave()
    setHTTPConfig({ ...values, type: 'http', name: proxyName })
    const newProxiyConfigs = clientProxyConfigs.map((proxyCfg) => {
      if (proxyCfg.name === proxyName) {
        return { ...values, type: 'http', name: proxyName } as HTTPProxyConfig
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
    setTimeoutID(setTimeout(() => {
      setSaveDisabled(false)
    }, 3000))
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
        {server && server.server && server.server.ip && defaultConfig &&
          defaultConfig.localIP && defaultConfig.localPort &&
          defaultConfig.subdomain
          && enablePreview && <div className="flex items-center space-x-2 flex-col justify-start w-full">
            <Label className="text-sm font-medium text-start w-full">{t('proxy.form.access_method')}</Label>
            <div className='w-full justify-start overflow-x-scroll'>
              <VisitPreview server={server?.server} typedProxyConfig={defaultConfig} />
            </div>
          </div>}
        <PortField name="localPort" control={form.control} label={t('proxy.form.local_port') + "*"} />
        <HostField name="localIP" control={form.control} label={t('proxy.form.local_ip') + "*"} />
        <StringField name="subdomain" control={form.control} label={t('proxy.form.subdomain')} placeholder={"your_sub_domain"} />
        <StringArrayField name="customDomains" control={form.control} label={t('proxy.form.custom_domains')} placeholder={"your.example.com"} />
        <FormDescription>
          {t('proxy.form.domain_description')}
        </FormDescription>
        <div className="flex items-center space-x-2 justify-between">
          <Label htmlFor="more-settings">{t('proxy.form.more_settings')}</Label>
          <Switch id="more-settings" checked={moreSettings} onCheckedChange={setMoreSettings} />
        </div>
        {moreSettings && <div className='p-4 space-y-4 border rounded-md'>
          <StringArrayField name="locations" control={form.control} label={t('proxy.form.route')} placeholder={"/path"} />
          <div className="flex items-center space-x-2 justify-between">
            <Label htmlFor="enable-http-auth">{t('proxy.form.enable_http_auth')}</Label>
            <Switch id="enable-http-auth" checked={useAuth} onCheckedChange={setUseAuth} />
          </div>
          {useAuth && <div className='p-4 space-y-4 border rounded-md'>
            <StringField name="httpUser" control={form.control} label={t('proxy.form.username')} placeholder={"username"} />
            <StringField name="httpPassword" control={form.control} label={t('proxy.form.password')} placeholder={"password"} />
          </div>}
        </div>}
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'} className='w-full'>
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}
