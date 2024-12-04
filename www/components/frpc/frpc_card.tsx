import React, { useEffect } from 'react'
import { useState } from 'react'
import { Label } from '@radix-ui/react-label'
import { useQuery } from '@tanstack/react-query'
import { getClient } from '@/api/client'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'
import { FRPCEditor } from './frpc_editor'
import { FRPCForm } from './frpc_form'
import { useSearchParams } from 'next/navigation'
import { ClientConfig } from '@/types/client'
import { TypedProxyConfig } from '@/types/proxy'
import { ClientSelector } from '../base/client-selector'
import { ServerSelector } from '../base/server-selector'

export interface FRPCFormCardProps {
  clientID?: string
  serverID?: string
}
export const FRPCFormCard: React.FC<FRPCFormCardProps> = ({
  clientID: defaultClientID,
  serverID: defaultServerID,
}: FRPCFormCardProps) => {
  const [advanceMode, setAdvanceMode] = useState<boolean>(false)
  const [clientID, setClientID] = useState<string | undefined>()
  const [serverID, setServerID] = useState<string | undefined>()
  const searchParams = useSearchParams()
  const paramClientID = searchParams.get('clientID')
  const [clientProxyConfigs, setClientProxyConfigs] = useState<TypedProxyConfig[]>([])

  useEffect(() => {
    if (defaultClientID) {
      setClientID(defaultClientID)
    }
    if (defaultServerID) {
      setServerID(defaultServerID)
    }
  }, [defaultClientID, defaultServerID])

  const { data: client, refetch: refetchClient } = useQuery({
    queryKey: ['getClient', clientID],
    queryFn: () => {
      return getClient({ clientId: clientID })
    },
  })

  useEffect(() => {
    if (!client || !client?.client) return
    if (client?.client?.config == undefined) return

    const clientConf = JSON.parse(client?.client?.config || '{}') as ClientConfig

    const proxyConfs = clientConf.proxies
    console.log('proxyConfs', proxyConfs)
    if (proxyConfs) {
      setClientProxyConfigs(proxyConfs)
    }
    if (clientConf != undefined && clientConf.proxies == undefined) {
      setClientProxyConfigs([])
    }
  }, [client, refetchClient, setClientProxyConfigs])

  useEffect(() => {
    if (paramClientID) {
      setClientID(paramClientID)
      if (client?.client?.serverId) {
        setServerID(client?.client?.serverId)
      }
    }
  }, [paramClientID])

  useEffect(() => {
    if (clientID && client?.client?.serverId) {
      setServerID(client?.client?.serverId)
    }
  }, [clientID, paramClientID, client])

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>编辑隧道</CardTitle>
        <CardDescription>
          <div>注意⚠️：选择的「服务端」必须提前配置！</div>
          <div>选择客户端和服务端以编辑隧道</div>
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className=" flex items-center space-x-4 rounded-md border p-4">
          <div className="flex-1 space-y-1">
            <p className="text-sm font-medium leading-none">高级模式</p>
            <p className="text-sm text-muted-foreground">编辑客户端原始配置文件</p>
          </div>
          <Switch onCheckedChange={setAdvanceMode} />
        </div>
        <div className="flex flex-col w-full pt-2 space-y-2">
          <Label className="text-sm font-medium">服务端</Label>
          <ServerSelector serverID={serverID} setServerID={setServerID} />
          <Label className="text-sm font-medium">客户端</Label>
          <ClientSelector clientID={clientID} setClientID={setClientID} />
        </div>
        {clientID && !advanceMode && <div className='flex flex-col w-full pt-2 space-y-2'>
          <Label className="text-sm font-medium">节点 {clientID} 的备注</Label>
          <p className="text-sm text-muted-foreground">可以到高级模式修改备注哦！</p>
          <p className="text-sm border rounded p-2 my-2">
            {client?.client?.comment == undefined || client?.client?.comment === '' ? '空空如也' : client?.client?.comment}
          </p></div>}
        {clientID && serverID && !advanceMode && <FRPCForm
          client={client?.client}
          clientConfig={JSON.parse(client?.client?.config || '{}') as ClientConfig} refetchClient={refetchClient}
          clientID={clientID} serverID={serverID}
          clientProxyConfigs={clientProxyConfigs}
          setClientProxyConfigs={setClientProxyConfigs} />
        }
        {clientID && serverID && advanceMode && <FRPCEditor
          client={client?.client}
          clientConfig={JSON.parse(client?.client?.config || '{}') as ClientConfig} refetchClient={refetchClient}
          clientID={clientID} serverID={serverID}
          clientProxyConfigs={clientProxyConfigs}
          setClientProxyConfigs={setClientProxyConfigs} />
        }
      </CardContent>
    </Card>
  )
}
