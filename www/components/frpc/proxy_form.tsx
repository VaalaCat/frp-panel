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
import { useStore } from '@nanostores/react'
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

const IPField = ({
  control,
  name,
  label,
  defaultValue,
}: {
  control: Control<any>
  name: string
  label: string
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
            <Input placeholder="127.0.0.1" {...field} />
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
  defaultValue,
}: {
  control: Control<any>
  name: string
  label: string
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
            <Input placeholder="8080" {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
      defaultValue={defaultValue}
    />
  )
}

const SecretKeyField = ({
  control,
  name,
  label,
  defaultValue,
}: {
  control: Control<any>
  name: string
  label: string
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
            <Input placeholder="secret key" {...field} />
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
  })

  useEffect(() => {
    setTCPConfig(undefined)
    form.reset({})
  }, [form])

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
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
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
        <PortField
          control={form.control}
          name="localPort"
          label="proxy.form.local_port"
          defaultValue={defaultConfig?.localPort || 1234}
        />
        <IPField
          control={form.control}
          name="localIP"
          label="proxy.form.local_ip"
          defaultValue={defaultConfig?.localIP || '127.0.0.1'}
        />
        <PortField
          control={form.control}
          name="remotePort"
          label="proxy.form.remote_port"
          defaultValue={defaultConfig?.remotePort || 4321}
        />
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
  })

  useEffect(() => {
    setSTCPConfig(undefined)
    form.reset({})
  }, [form])

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

  const { data: server } = useQuery({
    queryKey: ['getServer', serverID],
    queryFn: () => {
      return getServer({ serverId: serverID })
    },
  })

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <PortField
          control={form.control}
          name="localPort"
          label="proxy.form.local_port"
          defaultValue={defaultConfig?.localPort || 1234}
        />
        <IPField
          control={form.control}
          name="localIP"
          label="proxy.form.local_ip"
          defaultValue={defaultConfig?.localIP || '127.0.0.1'}
        />
        <SecretKeyField control={form.control} name="secretKey" label="proxy.form.secret_key" defaultValue={defaultConfig?.secretKey} />
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'}>
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}

export const UDPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID, defaultProxyConfig, proxyName, clientProxyConfigs, setClientProxyConfigs }) => {
  const [_, setUDPConfig] = useState<UDPProxyConfig | undefined>()
  const form = useForm<z.infer<typeof UDPConfigSchema>>({
    resolver: zodResolver(UDPConfigSchema),
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
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <Label className="text-sm font-medium">{t('proxy.form.access_method')}</Label>
        <p className="text-sm border rounded p-2 my-2 font-mono overflow-auto">
          {`${server?.server?.ip}:${(defaultProxyConfig as UDPProxyConfig).remotePort} -> ${
            defaultProxyConfig?.localIP
          }:${defaultProxyConfig?.localPort}`}
        </p>
        <FormField
          control={form.control}
          name="localPort"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t('proxy.form.local_port')}</FormLabel>
              <FormControl>
                <Input type="number" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
          defaultValue={defaultProxyConfig === undefined ? 1234 : defaultProxyConfig.localPort}
        />
        <FormField
          control={form.control}
          name="localIP"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t('proxy.form.local_ip')}</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
          defaultValue={defaultProxyConfig === undefined ? '127.0.0.1' : defaultProxyConfig.localIP}
        />
        <FormField
          control={form.control}
          name="remotePort"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t('proxy.form.remote_port')}</FormLabel>
              <FormControl>
                <Input type="number" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
          defaultValue={defaultProxyConfig === undefined ? 4321 : (defaultProxyConfig as UDPProxyConfig).remotePort}
        />
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
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <Label className="text-sm font-medium">{t('proxy.form.access_method')}</Label>
        <p className="text-sm border rounded p-2 my-2 font-mono overflow-auto">
          {`http://${(defaultProxyConfig as HTTPProxyConfig).subdomain}.${serverConfig?.subDomainHost}:${
            serverConfig?.vhostHTTPPort
          } -> ${defaultProxyConfig?.localIP}:${defaultProxyConfig?.localPort}`}
        </p>
        <FormField
          control={form.control}
          name="localPort"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t('proxy.form.local_port')}</FormLabel>
              <FormControl>
                <Input type="number" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
          defaultValue={defaultProxyConfig === undefined ? 1234 : defaultProxyConfig.localPort}
        />
        <FormField
          control={form.control}
          name="localIP"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t('proxy.form.local_ip')}</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
          defaultValue={defaultProxyConfig === undefined ? '127.0.0.1' : defaultProxyConfig.localIP}
        />
        <FormField
          control={form.control}
          name="subDomain"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t('proxy.form.subdomain')}</FormLabel>
              <FormControl>
                <Input type="text" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
          defaultValue={defaultProxyConfig === undefined ? '' : (defaultProxyConfig as HTTPProxyConfig).subdomain}
        />
        <Button type="submit" disabled={isSaveDisabled} variant={'outline'}>
          <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? '' : 'hidden'}`}></YesIcon>
          {t('proxy.form.save_changes')}
        </Button>
      </form>
    </Form>
  )
}
