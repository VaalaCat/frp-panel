import { Server } from '@/lib/pb/common'
import { ServerTableSchema, columns as serverColumnsDef } from './server_item'
import { DataTable } from './data_table'

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
import { useQuery } from '@tanstack/react-query'
import { listServer } from '@/api/server'

export interface ServerListProps {
  Servers: Server[]
  Keyword?: string
  TriggerRefetch?: string
}

export const ServerList: React.FC<ServerListProps> = ({ Servers, Keyword, TriggerRefetch }) => {
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([])
  const data = Servers.map(
    (server) =>
      ({
        id: server.id == undefined ? '' : server.id,
        status: server.config == undefined || server.config == '' ? 'invalid' : 'valid',
        secret: server.secret == undefined ? '' : server.secret,
        config: server.config,
      }) as ServerTableSchema,
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
  }
  const pagination = React.useMemo(
    () => ({
      pageIndex,
      pageSize,
    }),
    [pageIndex, pageSize],
  )

  const dataQuery = useQuery({
    queryKey: ['listServer', fetchDataOptions],
    queryFn: async () => {
      return await listServer({ page: fetchDataOptions.pageIndex + 1, pageSize: fetchDataOptions.pageSize, keyword: fetchDataOptions.Keyword })
    },
  })

  const table = useReactTable({
    data:
      dataQuery.data?.servers.map((server) => {
        return {
          id: server.id == undefined ? '' : server.id,
          status: server.config == undefined || server.config == '' ? 'invalid' : 'valid',
          secret: server.secret == undefined ? '' : server.secret,
          ip: server.ip,
          config: server.config,
        } as ServerTableSchema
      }) ?? data,
    pageCount: Math.ceil(
      //@ts-ignore
      (dataQuery.data?.total == undefined ? 0 : dataQuery.data?.total) / fetchDataOptions.pageSize ?? 0,
    ),
    columns: serverColumnsDef,
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
  return <DataTable table={table} columns={serverColumnsDef} />
}
