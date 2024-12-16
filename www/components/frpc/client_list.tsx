import { Client } from '@/lib/pb/common'
import { ClientTableSchema, columns as clientColumnsDef } from './client_item'
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
import { listClient } from '@/api/client'
import { ClientConfigured } from '@/lib/consts'
import { useStore } from '@nanostores/react'
import { $clientTableRefetchTrigger } from '@/store/refetch-trigger'

export interface ClientListProps {
  Clients: Client[]
  Keyword?: string
  TriggerRefetch?: string
}

export const ClientList: React.FC<ClientListProps> = ({ Clients, Keyword, TriggerRefetch }) => {
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([])
  const globalRefetchTrigger = useStore($clientTableRefetchTrigger)

  const data = Clients.map(
    (client) =>
      ({
        id: client.id == undefined ? '' : client.id,
        status: ClientConfigured(client) ? 'valid' : 'invalid',
        secret: client.secret == undefined ? '' : client.secret,
        config: client.config,
        originClient: client,
      }) as ClientTableSchema,
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
    queryKey: ['listClient', fetchDataOptions],
    queryFn: async () => {
      return await listClient({ page: fetchDataOptions.pageIndex + 1, pageSize: fetchDataOptions.pageSize, keyword: fetchDataOptions.Keyword })
    },
    placeholderData: keepPreviousData,
  })

  const table = useReactTable({
    data:
      dataQuery.data?.clients.map((client) => {
        return {
          id: client.id == undefined ? '' : client.id,
          status: ClientConfigured(client) ? 'valid' : 'invalid',
          secret: client.secret == undefined ? '' : client.secret,
          config: client.config,
          stopped: client.stopped,
          originClient: client,
        } as ClientTableSchema
      }) ?? data,
    pageCount: Math.ceil(
      // @ts-ignore
      (dataQuery.data?.total == undefined ? 0 : dataQuery.data?.total) / fetchDataOptions.pageSize ?? 0,
    ),
    columns: clientColumnsDef,
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
  return <DataTable table={table} columns={clientColumnsDef} />
}
