import { ProxyType, TypedProxyConfig } from '@/types/proxy'
import React, { useEffect } from 'react'
import { useState } from 'react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Label } from '@radix-ui/react-label'
import { HTTPProxyForm, STCPProxyForm, TCPProxyForm, UDPProxyForm } from './proxy_form'
import { Button } from '@/components/ui/button'
import { Client, RespCode } from '@/lib/pb/common'
import { ClientConfig } from '@/types/client'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Input } from '@/components/ui/input'
import { AccordionHeader } from '@radix-ui/react-accordion'
import { useToast } from '@/components/ui/use-toast'
import { QueryObserverResult, RefetchOptions, useMutation } from '@tanstack/react-query'
import { updateFRPC } from '@/api/frp'
import { Card, CardContent } from '@/components/ui/card'
import { GetClientResponse } from '@/lib/pb/api_client'
import { useTranslation } from 'react-i18next'

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
  const { toast } = useToast()
  
  const handleTypeChange = (value: string) => {
    setProxyType(value as ProxyType)
  }

  const handleAddProxy = () => {
    console.log('add proxy', proxyName, proxyType)
    if (!proxyName) return
    if (!proxyType) return
    if (clientProxyConfigs.findIndex((proxy) => proxy.name === proxyName) !== -1) {
      toast({ 
        title: t('proxy.status.create'), 
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
      toast({ 
        title: t('proxy.status.update'), 
        description: res.status?.code === RespCode.SUCCESS ? t('proxy.status.success') : t('proxy.status.error') 
      })
    } catch (error) {
      console.error(error)
      toast({ 
        title: t('proxy.status.update'), 
        description: t('proxy.status.error') 
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
          <AccordionContent className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
            {clientProxyConfigs.map((item) => {
              return (
                <Card key={item.name}>
                  <CardContent>
                    <div className="flex flex-col w-full pt-2">
                      <Accordion type="single" collapsible>
                        <AccordionItem value={item.name}>
                          <AccordionHeader className="flex flex-row justify-between">
                            <div>{t('proxy.form.tunnel_name')}: {item.name}</div>
                            <Button
                              variant={'outline'}
                              onClick={() => {
                                handleDeleteProxy(item.name)
                              }}
                            >
                              {t('proxy.form.delete')}
                            </Button>
                          </AccordionHeader>
                          <AccordionTrigger>{t('proxy.form.type_label', { type: item.type })}</AccordionTrigger>
                          <AccordionContent>
                            {item.type === 'tcp' && serverID && clientID && (
                              <TCPProxyForm
                                defaultProxyConfig={item}
                                proxyName={item.name}
                                serverID={serverID}
                                clientID={clientID}
                                clientProxyConfigs={clientProxyConfigs}
                                setClientProxyConfigs={setClientProxyConfigs}
                              />
                            )}
                            {item.type === 'udp' && serverID && clientID && (
                              <UDPProxyForm
                                defaultProxyConfig={item}
                                proxyName={item.name}
                                serverID={serverID}
                                clientID={clientID}
                                clientProxyConfigs={clientProxyConfigs}
                                setClientProxyConfigs={setClientProxyConfigs}
                              />
                            )}
                            {item.type === 'http' && serverID && clientID && (
                              <HTTPProxyForm
                                defaultProxyConfig={item}
                                proxyName={item.name}
                                serverID={serverID}
                                clientID={clientID}
                                clientProxyConfigs={clientProxyConfigs}
                                setClientProxyConfigs={setClientProxyConfigs}
                              />
                            )}
                            {item.type === 'stcp' && serverID && clientID && (
                              <STCPProxyForm
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
                    </div>
                  </CardContent>
                </Card>
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
