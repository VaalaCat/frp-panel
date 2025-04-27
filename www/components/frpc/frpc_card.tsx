"use client"

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
import { useTranslation } from 'react-i18next'
import { Input } from '../ui/input'
import { Server } from '@/lib/pb/common'
import { SuggestiveInput } from '../base/suggestive-input'

export interface FRPCFormCardProps {
  clientID?: string
  serverID?: string
}

export const FRPCFormCard: React.FC<FRPCFormCardProps> = ({
  clientID: defaultClientID,
  serverID: defaultServerID,
}: FRPCFormCardProps) => {
  const { t } = useTranslation()
  const [advanceMode, setAdvanceMode] = useState<boolean>(false)
  const [clientID, setClientID] = useState<string | undefined>()
  const [serverID, setServerID] = useState<string | undefined>()
  const searchParams = useSearchParams()
  const paramClientID = searchParams.get('clientID')
  const [clientProxyConfigs, setClientProxyConfigs] = useState<TypedProxyConfig[]>([])
  const [frpsUrl, setFrpsUrl] = useState<string | undefined>()
  const [selectedServer, setSelectedServer] = useState<Server | undefined>(undefined)

  useEffect(() => {
    if (defaultClientID) {
      setClientID(defaultClientID)
    }
    if (defaultServerID) {
      setServerID(defaultServerID)
    }
  }, [defaultClientID, defaultServerID])

  const { data: client, refetch: refetchClient, error } = useQuery({
    queryKey: ['getClient', clientID, serverID],
    queryFn: () => {
      return getClient({ clientId: clientID, serverId: serverID })
    },
    retry: false,
  })

  useEffect(() => {
    if (error) {
      setClientProxyConfigs([])
    }
  }, [error])

  useEffect(() => {
    if (!client || !client?.client) return
    if (client?.client?.config == undefined) return

    const clientConf = JSON.parse(client?.client?.config || '{}') as ClientConfig

    const proxyConfs = clientConf.proxies
    if (proxyConfs) {
      setClientProxyConfigs(proxyConfs)
    }
    if (clientConf != undefined && clientConf.proxies == undefined) {
      setClientProxyConfigs([])
    }

    if (client?.client?.frpsUrl) {
      setFrpsUrl(client?.client?.frpsUrl)
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
        <CardTitle>{t('frpc.form.title')}</CardTitle>
        <CardDescription>
          <div>{t('frpc.form.description.warning')}</div>
          <div>{t('frpc.form.description.instruction')}</div>
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex items-center space-x-4 rounded-md border p-4">
          <div className="flex-1 space-y-1">
            <p className="text-sm font-medium leading-none">{t('frpc.form.advanced.title')}</p>
            <p className="text-sm text-muted-foreground">{t('frpc.form.advanced.description')}</p>
          </div>
          <Switch onCheckedChange={setAdvanceMode} />
        </div>
        <div className="flex flex-col w-full pt-2 space-y-2">
          <Label className="text-sm font-medium">{t('frpc.form.server')}</Label>
          <ServerSelector serverID={serverID} setServerID={setServerID} setServer={setSelectedServer} />
          <Label className="text-sm font-medium">{t('frpc.form.client')}</Label>
          <ClientSelector clientID={clientID} setClientID={setClientID} />
          <Label className="text-sm font-medium">{t('frpc.form.frps_url.title')}</Label>
          <p className="text-sm text-muted-foreground">{t('frpc.form.frps_url.hint')}</p>
          <SuggestiveInput value={frpsUrl || ''} onChange={setFrpsUrl} suggestions={selectedServer?.frpsUrls || []} />
        </div>
        {clientID && !advanceMode && <div className='flex flex-col w-full pt-2 space-y-2'>
          <Label className="text-sm font-medium">{t('frpc.form.comment.title', { id: clientID })}</Label>
          <p className="text-sm text-muted-foreground">{t('frpc.form.comment.hint')}</p>
          <p className="text-sm border rounded p-2 my-2">
            {client?.client?.comment == undefined || client?.client?.comment === '' ? t('frpc.form.comment.empty') : client?.client?.comment}
          </p></div>}
        {clientID && serverID && !advanceMode && <FRPCForm
          client={client?.client}
          clientConfig={JSON.parse(client?.client?.config || '{}') as ClientConfig} refetchClient={refetchClient}
          clientID={clientID} serverID={serverID}
          clientProxyConfigs={clientProxyConfigs}
          setClientProxyConfigs={setClientProxyConfigs}
          frpsUrl={frpsUrl}
        />
        }
        {clientID && serverID && advanceMode && <FRPCEditor
          client={client?.client}
          clientConfig={JSON.parse(client?.client?.config || '{}') as ClientConfig} refetchClient={refetchClient}
          clientID={clientID} serverID={serverID}
          clientProxyConfigs={clientProxyConfigs}
          setClientProxyConfigs={setClientProxyConfigs}
          frpsUrl={frpsUrl}
        />
        }
      </CardContent>
    </Card>
  )
}
