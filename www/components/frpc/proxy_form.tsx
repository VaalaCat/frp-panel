import { HTTPProxyConfig, TCPProxyConfig, TypedProxyConfig, UDPProxyConfig, STCPProxyConfig } from '@/types/proxy'
import * as z from 'zod'
import React from 'react'
import { ZodPortSchema, ZodStringSchema } from '@/lib/consts'
import { useEffect, useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { Control, FieldValues, useForm } from 'react-hook-form'
import { Button } from '@/components/ui/button'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { YesIcon } from '@/components/ui/icon'
import { useTranslation } from 'react-i18next'
import { useQuery } from '@tanstack/react-query'
import { getServer } from '@/api/server'
import { ServerConfig } from '@/types/server'
import { ArrowRightIcon } from 'lucide-react'

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
  subDomain: ZodStringSchema,
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
            <Input placeholder={placeholder || '127.0.0.1'} {...field} />
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
            <Input placeholder={placeholder || '1234'} {...field} />
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
            <Input placeholder={placeholder || "secret"} {...field} />
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
            <Input placeholder={placeholder || '127.0.0.1'} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}    
      defaultValue={defaultValue}
    />
  )
}

export const TCPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID, defaultProxyConfig, proxyName, clientProxyConfigs, setClientProxyConfigs }) => {
  const defaultConfig = defaultProxyConfig as TCPProxyConfig
  const [_, setTCPConfig] = useState<TCPProxyConfig | undefined>()
  const form = useForm<z.infer<typeof TCPConfigSchema>>({
    resolver: zodResolver(TCPConfigSchema),
    defaultValues: {
      remotePort: defaultConfig?.remotePort,
      localIP: defaultConfig?.localIP,
      localPort: defaultConfig?.localPort,
    }
  })

  useEffect(() => {
    setTCPConfig(undefined)
    form.reset({})
  }, [])

  const onSubmit = async (values: z.infer<typeof TCPConfigSchema>) => {
    handleSave()
    setTCPConfig({ type: 'tcp', ...values, name: proxyName })
    const newProxiyConfigs = clientProxyConfigs.map((proxyCfg) => {
      if (proxyCfg.name === proxyName) {
        return { ...values, type: 'tcp', name: proxyName } as TCPProxyConfig
      }
      return proxyCfg
    })
    setClientProxyConfigs(newProxiyConfigs)
  }

  const [isSaveDisabled, setSaveDisabled] = useState(false)

  const handleSave = () => {
    setSaveDisabled(true)
    setTimeout(() => {
      setSaveDisabled(false)
    }, 3000)
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
        <Label className="text-sm font-medium">{t('proxy.form.access_method')}</Label>
        {server?.server?.ip && defaultConfig.remotePort && defaultConfig.localIP && defaultConfig.localPort && (
          <div className="flex items-center space-x-2">
            <Input
              value={`${server?.server?.ip}:${defaultConfig?.remotePort}`}
              className="text-sm font-mono"
              disabled
            />{' '}
            <ArrowRightIcon className="h-4 w-4" />{' '}
            <Input
              value={`${defaultConfig?.localIP}:${defaultConfig?.localPort}`}
              className="text-sm font-mono"
              disabled
            />
          </div>
        )}
        <PortField name="localPort" control={form.control} label={t('proxy.form.local_port')} />
        <HostField name="localIP" control={form.control} label={t('proxy.form.local_ip')} />
        <PortField name="remotePort" control={form.control} label={t('proxy.form.remote_port')} placeholder='4321'/>
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'}>
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}

export const STCPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID, defaultProxyConfig, proxyName, clientProxyConfigs, setClientProxyConfigs }) => {
  const defaultConfig = defaultProxyConfig as STCPProxyConfig
  const [_, setSTCPConfig] = useState<STCPProxyConfig | undefined>()
  const form = useForm<z.infer<typeof STCPConfigSchema>>({
    resolver: zodResolver(STCPConfigSchema),
    defaultValues: {
      localPort: defaultConfig?.localPort,
      localIP: defaultConfig?.localIP,
      secretKey: defaultConfig?.secretKey,
    }
  })

  useEffect(() => {
    setSTCPConfig(undefined)
    form.reset({})
  }, [])

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
    setTimeout(() => {
      setSaveDisabled(false)
    }, 3000)
  }

  const { t } = useTranslation()

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 px-0.5">
        <PortField name="localPort" control={form.control} label={t('proxy.form.local_port')} />
        <HostField name="localIP" control={form.control} label={t('proxy.form.local_ip')} />
        <SecretStringField name="secretKey" control={form.control} label={t('proxy.form.secret_key')} />
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'}>
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}

export const UDPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID, defaultProxyConfig, proxyName, clientProxyConfigs, setClientProxyConfigs }) => {
  const defaultConfig = defaultProxyConfig as UDPProxyConfig
  const [_, setUDPConfig] = useState<UDPProxyConfig | undefined>()
  const form = useForm<z.infer<typeof UDPConfigSchema>>({
    resolver: zodResolver(UDPConfigSchema),
    defaultValues: {
      localPort: defaultConfig?.localPort,
      localIP: defaultConfig?.localIP,
      remotePort: defaultConfig?.remotePort,
    }
  })

  useEffect(() => {
    setUDPConfig(undefined)
    form.reset({})
  }, [])

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
    setTimeout(() => {
      setSaveDisabled(false)
    }, 3000)
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
        {server?.server?.ip && defaultConfig.remotePort && defaultConfig.localIP && defaultConfig.localPort && (
          <div className="flex items-center space-x-2">
            <Input
              value={`${server?.server?.ip}:${defaultConfig?.remotePort}`}
              className="text-sm font-mono"
              disabled
            />{' '}
            <ArrowRightIcon className="h-4 w-4" />{' '}
            <Input
              value={`${defaultConfig?.localIP}:${defaultConfig?.localPort}`}
              className="text-sm font-mono"
              disabled
            />
          </div>
        )}
        <PortField name="localPort" control={form.control} label={t('proxy.form.local_port')} />
        <HostField name="localIP" control={form.control} label={t('proxy.form.local_ip')} />
        <PortField name="remotePort" control={form.control} label={t('proxy.form.remote_port')} />
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'}>
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}

export const HTTPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID, defaultProxyConfig, proxyName, clientProxyConfigs, setClientProxyConfigs }) => {
  const [_, setHTTPConfig] = useState<HTTPProxyConfig | undefined>()
  const [serverConfig, setServerConfig] = useState<ServerConfig | undefined>()
  const form = useForm<z.infer<typeof HTTPConfigSchema>>({
    resolver: zodResolver(HTTPConfigSchema),
  })

  useEffect(() => {
    setHTTPConfig(undefined)
    form.reset({})
  }, [])

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

  const handleSave = () => {
    setSaveDisabled(true)
    setTimeout(() => {
      setSaveDisabled(false)
    }, 3000)
  }

  const { t } = useTranslation()

  const { data: server } = useQuery({
    queryKey: ['getServer', serverID],
    queryFn: () => {
      return getServer({ serverId: serverID })
    },
  })

  useEffect(() => {
    if (server && server.server?.config) {
      setServerConfig(JSON.parse(server.server?.config) as ServerConfig)
    }
  }, [server])

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 px-0.5">
        <Label className="text-sm font-medium">{t('proxy.form.access_method')}</Label>
        <p className="text-sm border rounded p-2 my-2 font-mono overflow-auto">
          {`http://${(defaultProxyConfig as HTTPProxyConfig).subdomain}.${serverConfig?.subDomainHost}:${serverConfig?.vhostHTTPPort
            } -> ${defaultProxyConfig?.localIP}:${defaultProxyConfig?.localPort}`}
        </p>
        <PortField name="localPort" control={form.control} label={t('proxy.form.local_port')} />
        <HostField name="localIP" control={form.control} label={t('proxy.form.local_ip')} />
        <StringField name="subDomain" control={form.control} label={t('proxy.form.subdomain')} placeholder={"your_sub_domain"} />
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'}>
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}
