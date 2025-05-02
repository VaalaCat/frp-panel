import { ProxyConfig } from '@/lib/pb/common'
import { ProxyConfigTableSchema, columns as proxyConfigColumnsDef } from './proxy_config_item'
import { DataTable } from '../base/data_table'

import {
  getSortedRowModel,
  getCoreRowModel,
  ColumnFiltersState,
  useReactTable,
  getFilteredRowModel,
  getPaginationRowModel,
  SortingState,
  PaginationState,
} from '@tanstack/react-table'

import React from 'react'
import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { listProxyConfig } from '@/api/proxy'
import { TypedProxyConfig } from '@/types/proxy'
import { $proxyTableRefetchTrigger } from '@/store/refetch-trigger'
import { useStore } from '@nanostores/react'

export interface ProxyConfigListProps {
  ProxyConfigs: ProxyConfig[]
  Keyword?: string
  ClientID?: string
  ServerID?: string
  TriggerRefetch?: string
}

export const ProxyConfigList: React.FC<ProxyConfigListProps> = ({
  ProxyConfigs,
  Keyword,
  TriggerRefetch,
  ClientID,
  ServerID,
}) => {
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([])
  const globalRefetchTrigger = useStore($proxyTableRefetchTrigger)

  const data = ProxyConfigs.map(
    (proxy_config) =>
      ({
        id: proxy_config.id || '',
        clientID: proxy_config.clientId || '',
        serverID: proxy_config.serverId || '',
        name: proxy_config.name || '',
        type: proxy_config.type || '',
        visitPreview: 'for test',
        originalProxyConfig: proxy_config,
        stopped: proxy_config.stopped,
      }) as ProxyConfigTableSchema,
  )

  const [{ pageIndex, pageSize }, setPagination] = React.useState<PaginationState>({
    pageIndex: 0,
    pageSize: 10,
  })

  const fetchDataOptions = {
    pageIndex,
    pageSize,
    Keyword,
    TriggerRefetch,
    ClientID,
    ServerID,
    globalRefetchTrigger,
  }
  const pagination = React.useMemo(
    () => ({
      pageIndex,
      pageSize,
    }),
    [pageIndex, pageSize],
  )

  const dataQuery = useQuery({
    queryKey: ['listProxyConfigs', fetchDataOptions],
    queryFn: async () => {
      return await listProxyConfig({
        page: fetchDataOptions.pageIndex + 1,
        pageSize: fetchDataOptions.pageSize,
        keyword: fetchDataOptions.Keyword,
        clientId: fetchDataOptions.ClientID,
        serverId: fetchDataOptions.ServerID,
      })
    },
    placeholderData: keepPreviousData,
  })

  const table = useReactTable({
    data:
      dataQuery.data?.proxyConfigs.map((proxy_config) => {
        return {
          id: proxy_config.id || '',
          name: proxy_config.name || '',
          clientID: proxy_config.clientId || '',
          serverID: proxy_config.serverId || '',
          type: proxy_config.type || '',
          config: proxy_config.config || '',
          localIP: proxy_config.config && ParseProxyConfig(proxy_config.config).localIP,
          localPort: proxy_config.config && ParseProxyConfig(proxy_config.config).localPort,
          visitPreview: '',
          originalProxyConfig: proxy_config,
          stopped: proxy_config.stopped || false,
        } as ProxyConfigTableSchema
      }) ?? data,
    pageCount: Math.ceil(
      //@ts-ignore
      (dataQuery.data?.total == undefined ? 0 : dataQuery.data?.total) / fetchDataOptions.pageSize ?? 0,
    ),
    columns: proxyConfigColumnsDef,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onSortingChange: setSorting,
    onPaginationChange: setPagination,
    onColumnFiltersChange: setColumnFilters,
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
    manualPagination: true,
    state: {
      sorting,
      pagination,
      columnFilters,
    },
  })
  return <DataTable table={table} columns={proxyConfigColumnsDef} />
}

function ParseProxyConfig(cfg: string): TypedProxyConfig {
  return JSON.parse(cfg)
}
