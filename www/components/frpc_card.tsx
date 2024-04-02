import React, { useEffect } from 'react'
import { useState } from 'react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Label } from '@radix-ui/react-label'
import { useQuery } from '@tanstack/react-query'
import { listServer } from '@/api/server'
import { getClient, listClient } from '@/api/client'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Switch } from './ui/switch'
import { FRPCEditor } from './frpc_editor'
import { FRPCForm } from './frpc_form'
import { useSearchParams } from 'next/navigation'
import { $clientProxyConfigs } from '@/store/proxy'

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
  const handleServerChange = (value: string) => {
    setServerID(value)
  }

  const handleClientChange = (value: string) => {
    setClientID(value)
  }

  useEffect(() => {
    if (defaultClientID) {
      setClientID(defaultClientID)
    }
    if (defaultServerID) {
      setServerID(defaultServerID)
    }
  }, [defaultClientID, defaultServerID])

  const { data: serverList, refetch: refetchServers } = useQuery({
    queryKey: ['listServer'],
    queryFn: () => {
      return listServer({ page: 1, pageSize: 500 })
    },
  })

  const { data: clientList, refetch: refetchClients } = useQuery({
    queryKey: ['listClient'],
    queryFn: () => {
      return listClient({ page: 1, pageSize: 500 })
    },
  })

  const { data: client, refetch: refetchClient } = useQuery({
    queryKey: ['getClient', clientID],
    queryFn: () => {
      return getClient({ clientId: clientID })
    },
  })

  useEffect(() => {
    if (paramClientID) {
      setClientID(paramClientID)
      setServerID(clientList?.clients?.find((client) => client.id == paramClientID)?.serverId)
      refetchClient()
    }
    $clientProxyConfigs.set([])
  }, [paramClientID, clientList])

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
        <div className="flex flex-col w-full pt-2">
          <Label className="text-sm font-medium">服务端</Label>
          <Select
            onValueChange={handleServerChange}
            value={serverID}
            onOpenChange={() => {
              refetchServers()
              refetchClients()
            }}
          >
            <SelectTrigger className="my-2">
              <SelectValue placeholder="节点名称" />
            </SelectTrigger>
            <SelectContent>
              {serverList?.servers.map(
                (server) =>
                  server.id && (
                    <SelectItem key={server.id} value={server.id}>
                      {server.id}
                    </SelectItem>
                  ),
              )}
            </SelectContent>
          </Select>
          <Label className="text-sm font-medium">客户端</Label>
          <Select
            onValueChange={handleClientChange}
            value={clientID}
            onOpenChange={() => {
              refetchServers()
              refetchClients()
            }}
          >
            <SelectTrigger className="my-2">
              <SelectValue placeholder="节点名称" />
            </SelectTrigger>
            <SelectContent>
              {clientList?.clients.map(
                (client) =>
                  client.id && (
                    <SelectItem key={client.id} value={client.id}>
                      {client.id}
                    </SelectItem>
                  ),
              )}
            </SelectContent>
          </Select>
        </div>
        {clientID && !advanceMode && <>
        <Label className="text-sm font-medium">节点 {clientID} 的备注</Label>
        <p className="text-sm text-muted-foreground">可以到高级模式修改备注哦！</p>
        <p className="text-sm border rounded p-2 my-2">
          {client?.client?.comment == undefined || client?.client?.comment === '' ? '空空如也' : client?.client?.comment}
        </p></>}
        {clientID && serverID && !advanceMode && <FRPCForm clientID={clientID} serverID={serverID} />}
        {clientID && serverID && advanceMode && <FRPCEditor clientID={clientID} serverID={serverID} />}
      </CardContent>
    </Card>
  )
}
