import { ProxyType, TypedProxyConfig } from '@/types/proxy'
import React, { useEffect } from 'react'
import { useState } from 'react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Label } from '@radix-ui/react-label'
import { HTTPProxyForm, STCPProxyForm, TCPProxyForm, TypedProxyForm, UDPProxyForm } from './proxy_form'
import { Button } from '@/components/ui/button'
import { Client, RespCode } from '@/lib/pb/common'
import { ClientConfig } from '@/types/client'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Input } from '@/components/ui/input'
import { AccordionHeader } from '@radix-ui/react-accordion'
import { QueryObserverResult, RefetchOptions, useMutation } from '@tanstack/react-query'
import { updateFRPC } from '@/api/frp'
import { Card, CardContent } from '@/components/ui/card'
import { GetClientResponse } from '@/lib/pb/api_client'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

export interface FRPCFormProps {
  clientID: string
  serverID: string
  client?: Client
  clientConfig: ClientConfig
  refetchClient: (options?: RefetchOptions) => Promise<QueryObserverResult<GetClientResponse, Error>>
  clientProxyConfigs: TypedProxyConfig[]
  setClientProxyConfigs: React.Dispatch<React.SetStateAction<TypedProxyConfig[]>>
}

export const FRPCForm: React.FC<FRPCFormProps> = ({ clientID, serverID, client, refetchClient, clientProxyConfigs, setClientProxyConfigs }) => {
  const { t } = useTranslation()
  const [proxyType, setProxyType] = useState<ProxyType>('http')
  const [proxyName, setProxyName] = useState<string | undefined>()

  const handleTypeChange = (value: string) => {
    setProxyType(value as ProxyType)
  }

  const handleAddProxy = () => {
    console.log('add proxy', proxyName, proxyType)
    if (!proxyName) return
    if (!proxyType) return
    if (clientProxyConfigs.findIndex((proxy) => proxy.name === proxyName) !== -1) {
      toast(t('proxy.status.create'), {
        description: t('proxy.status.name_exists')
      })
      return
    }
    const newProxy = {
      name: proxyName,
      type: proxyType,
    } as TypedProxyConfig
    setClientProxyConfigs([...clientProxyConfigs, newProxy])
  }

  const handleDeleteProxy = (proxyName: string) => {
    const newProxies = clientProxyConfigs.filter((proxy) => proxy.name !== proxyName)
    setClientProxyConfigs(newProxies)
  }

  const updateFrpc = useMutation({ mutationFn: updateFRPC })

  const handleUpdate = async () => {
    try {
      const res = await updateFrpc.mutateAsync({
        //@ts-ignore
        config: Buffer.from(
          JSON.stringify({
            proxies: clientProxyConfigs,
          } as ClientConfig),
        ),
        serverId: serverID,
        clientId: clientID,
      })
      await refetchClient()
      toast(t('proxy.status.update'), {
        description: res.status?.code === RespCode.SUCCESS ? t('proxy.status.success') : t('proxy.status.error')
      })
    } catch (error) {
      console.error(error)
      toast(t('proxy.status.update'), {
        description: t('proxy.status.error') + JSON.stringify(error)
      })
    }
  }

  return (
    <>
      <Popover>
        <PopoverTrigger asChild>
          <Button className="my-2">{t('proxy.form.add')}</Button>
        </PopoverTrigger>
        <PopoverContent>
          <Label className="text-sm font-medium">{t('proxy.form.name')}</Label>
          <Input
            onChange={(e) => {
              setProxyName(e.target.value)
            }}
          />
          <Select onValueChange={handleTypeChange} defaultValue={proxyType}>
            <Label className="text-sm font-medium">{t('proxy.form.protocol')}</Label>
            <SelectTrigger className="my-2">
              <SelectValue placeholder={t('proxy.form.type')} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="http">{t('proxy.type.http')}</SelectItem>
              <SelectItem value="tcp">{t('proxy.type.tcp')}</SelectItem>
              <SelectItem value="udp">{t('proxy.type.udp')}</SelectItem>
              <SelectItem value="stcp">{t('proxy.type.stcp')}</SelectItem>
            </SelectContent>
          </Select>
          <Button variant={'outline'} onClick={handleAddProxy}>
            {t('proxy.form.confirm')}
          </Button>
        </PopoverContent>
      </Popover>
      <Accordion type="single" defaultValue="proxies" collapsible key={clientID + serverID + client}>
        <AccordionItem value="proxies">
          <AccordionTrigger>
            <AccordionHeader className="flex flex-row justify-between w-full">
              <p>{t('proxy.form.config')}</p>
              <p>{t('proxy.form.expand', { count: clientProxyConfigs.length })}</p>
            </AccordionHeader>
          </AccordionTrigger>
          <AccordionContent className="grid gap-2 grid-cols-1">
            {clientProxyConfigs.map((item, index) => {
              return (
                <Accordion type="single" collapsible key={index}>
                  <AccordionItem value={item.name}>
                    <AccordionTrigger>
                      <div className='flex flex-row justify-start items-center w-full gap-4'>
                        <Button variant={'outline'} onClick={() => { handleDeleteProxy(item.name) }}>
                          {t('proxy.form.delete')}
                        </Button>
                        <div>{t('proxy.form.tunnel_name')}: {item.name}</div>
                        <div>{t('proxy.form.type_label', { type: item.type })}</div>
                      </div>
                    </AccordionTrigger>
                    <AccordionContent className='border rounded-xl p-4'>
                      {serverID && clientID && (
                        <TypedProxyForm
                          enablePreview
                          defaultProxyConfig={item}
                          proxyName={item.name}
                          serverID={serverID}
                          clientID={clientID}
                          clientProxyConfigs={clientProxyConfigs}
                          setClientProxyConfigs={setClientProxyConfigs}
                        />
                      )}
                    </AccordionContent>
                  </AccordionItem>
                </Accordion>
              )
            })}
          </AccordionContent>
        </AccordionItem>
      </Accordion>
      <Button
        className="mt-2"
        onClick={() => {
          handleUpdate()
        }}
      >
        {t('proxy.form.submit')}
      </Button>
    </>
  )
}
