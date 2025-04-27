import { ColumnDef, Row, Table } from '@tanstack/react-table'
import React from 'react'
import { useTranslation } from 'react-i18next'
import { VisitPreview } from '../base/visit-preview'
import { useQuery } from '@tanstack/react-query'
import { getServer } from '@/api/server'
import { ProxyType, TypedProxyConfig } from '@/types/proxy'
import { ProxyConfigActions } from './proxy_config_actions'
import { ProxyConfig } from '@/lib/pb/common'
import { getProxyConfig } from '@/api/proxy'
import { Badge } from '../ui/badge'
import { useStore } from '@nanostores/react'
import { $proxyTableRefetchTrigger } from '@/store/refetch-trigger'

export type ProxyConfigTableSchema = {
  serverID: string
  clientID: string
  name: string
  type: ProxyType
  localIP?: string
  localPort?: number
  visitPreview: string
  config?: string
  originalProxyConfig: ProxyConfig
}

export const columns: ColumnDef<ProxyConfigTableSchema>[] = [
  {
    accessorKey: 'clientID',
    header: function Header() {
      const { t } = useTranslation()
      return t('proxy.item.client_id')
    },
    cell: ({ row }) => {
      return <div className='font-mono text-nowrap'>{row.original.originalProxyConfig.originClientId}</div>
    },
  },
  {
    accessorKey: 'name',
    header: function Header() {
      const { t } = useTranslation()
      return t('proxy.item.proxy_name')
    },
    cell: ({ row }) => {
      return <div className='font-mono text-nowrap'>{row.original.name}</div>
    },
  },
  {
    accessorKey: 'type',
    header: function Header() {
      const { t } = useTranslation()
      return t('proxy.item.proxy_type')
    },
    cell: ({ row }) => {
      return <div className='font-mono text-nowrap'>{row.original.type}</div>
    },
  },
  {
    accessorKey: 'serverID',
    header: function Header() {
      const { t } = useTranslation()
      return t('proxy.item.server_id')
    },
    cell: ({ row }) => {
      return <div className='font-mono text-nowrap'>{row.original.serverID}</div>
    },
  },
  {
    id: 'status',
    header: function Header() {
      const { t } = useTranslation()
      return t('proxy.item.status')
    },
    cell: ProxyStatus,
  },
  {
    accessorKey: 'visitPreview',
    header: function Header() {
      const { t } = useTranslation()
      return t('proxy.item.visit_preview')
    },
    cell: ({ row }) => {
      return <VisitPreviewField row={row} />
    },
  },
  {
    id: 'action',
    cell: ({ row }) => {
      return <ProxyConfigActions
        row={row}
        serverID={row.original.serverID}
        clientID={row.original.clientID}
        name={row.original.name}
      />
    },
  }
]

function VisitPreviewField({ row }: { row: Row<ProxyConfigTableSchema> }) {
  const { data: server } = useQuery({
    queryKey: ['getServer', row.original.serverID],
    queryFn: () => {
      return getServer({ serverId: row.original.serverID })
    },
  })

  const typedProxyConfig = JSON.parse(row.original.config || '{}') as TypedProxyConfig

  return <VisitPreview
    server={server?.server || {frpsUrls: []}}
    typedProxyConfig={typedProxyConfig} />
}

function ProxyStatus({ row }: { row: Row<ProxyConfigTableSchema> }) {
  const refetchTrigger = useStore($proxyTableRefetchTrigger)
  const { data } = useQuery({
    queryKey: ['getProxyConfig', row.original.clientID, row.original.serverID, row.original.name, refetchTrigger],
    queryFn: () => {
      return getProxyConfig({
        clientId: row.original.clientID,
        serverId: row.original.serverID,
        name: row.original.name
      })
    },
    refetchInterval: 10000
  })

  function getStatusColor(status: string): string {
    switch (status) {
      case 'new':
        return 'text-blue-500'
      case 'wait start':
        return 'text-yellow-400';
      case 'start error':
        return 'text-red-500';
      case 'running':
        return 'text-green-500';
      case 'check failed':
        return 'text-orange-500';
      case 'error':
        return 'text-red-600';
      default:
        return 'text-gray-500';
    }
  }

  return <div className="flex items-center gap-2 flex-row text-nowrap">
    <Badge variant={"secondary"} className={`p-2 border rounded font-mono w-fit ${getStatusColor(data?.workingStatus?.status || 'unknown')} text-nowrap rounded-full h-6`}>
      {data?.workingStatus?.status || "loading"}
    </Badge>
  </div>
}