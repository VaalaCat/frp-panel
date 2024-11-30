import { ProxyType, TypedProxyConfig } from '@/types/proxy'
import React, { useEffect } from 'react'
import { useState } from 'react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Label } from '@radix-ui/react-label'
import { HTTPProxyForm, STCPProxyForm, TCPProxyForm, UDPProxyForm } from './proxy_form'
import { Button } from './ui/button'
import { Client, RespCode } from '@/lib/pb/common'
import { ClientConfig } from '@/types/client'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Input } from './ui/input'
import { AccordionHeader } from '@radix-ui/react-accordion'
import { useToast } from './ui/use-toast'
import { QueryObserverResult, RefetchOptions, useMutation } from '@tanstack/react-query'
import { updateFRPC } from '@/api/frp'
import { Card, CardContent } from './ui/card'
import { GetClientResponse } from '@/lib/pb/api_client'

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
      toast({ title: '创建隧道状态', description: '名称重复' })
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
      toast({ title: '更新隧道状态', description: res.status?.code === RespCode.SUCCESS ? '更新成功' : '更新失败' })
    } catch (error) {
      console.error(error)
      toast({ title: '更新隧道状态', description: '更新失败' })
    }
  }

  return (
    <>
      <Popover>
        <PopoverTrigger asChild>
          <Button className="my-2">新增隧道</Button>
        </PopoverTrigger>
        <PopoverContent>
          <Label className="text-sm font-medium">名称</Label>
          <Input
            onChange={(e) => {
              setProxyName(e.target.value)
            }}
          />
          <Select onValueChange={handleTypeChange} defaultValue={proxyType}>
            <Label className="text-sm font-medium">协议</Label>
            <SelectTrigger className="my-2">
              <SelectValue placeholder="类型" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="http">http</SelectItem>
              <SelectItem value="tcp">tcp</SelectItem>
              <SelectItem value="udp">udp</SelectItem>
              <SelectItem value="stcp">stcp</SelectItem>
            </SelectContent>
          </Select>
          <Button variant={'outline'} onClick={handleAddProxy}>
            确定
          </Button>
        </PopoverContent>
      </Popover>
      <Accordion type="single" defaultValue="proxies" collapsible key={clientID + serverID + client}>
        <AccordionItem value="proxies">
          <AccordionTrigger>
            <AccordionHeader className="flex flex-row justify-between w-full"><p>隧道配置</p> <p>点击展开{`${clientProxyConfigs.length}条隧道`}</p></AccordionHeader>
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
                            <div>隧道名称：{item.name}</div>
                            <Button
                              variant={'outline'}
                              onClick={() => {
                                handleDeleteProxy(item.name)
                              }}
                            >
                              删除
                            </Button>
                          </AccordionHeader>
                          <AccordionTrigger>类型:「{item.type}」</AccordionTrigger>
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
        {/* <AccordionItem value="visitors">
          <AccordionTrigger>
            <AccordionHeader className="flex flex-row justify-between">Visitor 配置</AccordionHeader>
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
                            <div>隧道名称：{item.name}</div>
                            <Button
                              variant={'outline'}
                              onClick={() => {
                                handleDeleteProxy(item.name)
                              }}
                            >
                              删除
                            </Button>
                          </AccordionHeader>
                          <AccordionTrigger>类型:「{item.type}」</AccordionTrigger>
                          <AccordionContent>
                            {item.type === 'tcp' && serverID && clientID && (
                              <TCPProxyForm
                                defaultProxyConfig={item}
                                proxyName={item.name}
                                serverID={serverID}
                                clientID={clientID}
                              />
                            )}
                            {item.type === 'udp' && serverID && clientID && (
                              <UDPProxyForm
                                defaultProxyConfig={item}
                                proxyName={item.name}
                                serverID={serverID}
                                clientID={clientID}
                              />
                            )}
                            {item.type === 'http' && serverID && clientID && (
                              <HTTPProxyForm
                                defaultProxyConfig={item}
                                proxyName={item.name}
                                serverID={serverID}
                                clientID={clientID}
                              />
                            )}
                            {item.type === 'stcp' && serverID && clientID && (
                              <STCPProxyForm
                                defaultProxyConfig={item}
                                proxyName={item.name}
                                serverID={serverID}
                                clientID={clientID}
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
        </AccordionItem> */}
      </Accordion>
      <Button
        className="mt-2"
        onClick={() => {
          handleUpdate()
        }}
      >
        提交变更
      </Button>
    </>
  )
}
