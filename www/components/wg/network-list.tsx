"use client"

import React from 'react'
import { useTranslation } from 'react-i18next'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import { getCoreRowModel, getPaginationRowModel, getSortedRowModel, getFilteredRowModel, useReactTable, SortingState, PaginationState, ColumnFiltersState } from '@tanstack/react-table'
import { DataTable } from '@/components/base/data_table'
import { listNetworks } from '@/api/wg'
import { ListNetworksRequest } from '@/lib/pb/api_wg'
import { AclConfig, Network } from '@/lib/pb/types_wg'
import { createNetworkColumns } from './network-row'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { RotateCcw, Columns3 } from 'lucide-react'
import { DropdownMenu, DropdownMenuCheckboxItem, DropdownMenuContent, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'

export type NetworkRow = {
  id: number;
  name: string;
  cidr: string,
  acl?: AclConfig
  origin: Network
}

export function NetworkList({ keyword, refreshToken, onChanged, onSummary }: { keyword?: string; refreshToken?: string; onChanged?: () => void; onSummary?: (info: { total: number }) => void }) {
  const { t } = useTranslation()
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([])
  const [columnVisibility, setColumnVisibility] = React.useState<Record<string, boolean>>({})
  const [{ pageIndex, pageSize }, setPagination] = React.useState<PaginationState>({ pageIndex: 0, pageSize: 10 })

  const { data, isFetching, refetch } = useQuery({
    queryKey: ['listNetworks', keyword, pageIndex, pageSize, refreshToken],
    queryFn: () => listNetworks(ListNetworksRequest.create({ page: pageIndex + 1, pageSize, keyword })),
    placeholderData: keepPreviousData,
  })

  const handleMutated = React.useCallback(() => {
    refetch()
    onChanged?.()
  }, [refetch, onChanged])

  const columns = React.useMemo(() => createNetworkColumns({ onChanged: handleMutated, t }), [handleMutated, t])

  const rows: NetworkRow[] = (data?.networks ?? []).map((n) => ({
    id: n.id!,
    name: n.name!,
    cidr: n.cidr!,
    acl: n.acl,
    origin: n
  }))

  const total = data?.total ?? 0

  React.useEffect(() => {
    onSummary?.({ total })
  }, [onSummary, total])

  const table = useReactTable({
    data: rows,
    columns,
    state: { sorting, pagination: { pageIndex, pageSize }, columnFilters, columnVisibility },
    manualPagination: true,
    pageCount: Math.ceil((data?.total ?? 0) / pageSize),
    onSortingChange: setSorting,
    onPaginationChange: setPagination,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onColumnVisibilityChange: setColumnVisibility,
  })

  const toolbar = (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="gap-2">
          <Columns3 className="h-4 w-4" />
          {t('table.columns')}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="min-w-[200px]">
        {table.getAllLeafColumns().map((column) => {
          if (column.getCanHide() === false) return null
          // const labelKey = column.columnDef.meta?.label as string | undefined
          return (
            <DropdownMenuCheckboxItem
              key={column.id}
              className="capitalize"
              checked={column.getIsVisible()}
              onCheckedChange={(value) => column.toggleVisibility(!!value)}
            >
              {/* {labelKey ? t(labelKey) : column.id} */}
              {column.id}
            </DropdownMenuCheckboxItem>
          )
        })}
      </DropdownMenuContent>
    </DropdownMenu>
  )

  return (
    <Card>
      <CardHeader className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="space-y-1">
          <CardTitle>{t('wg.networkList.headerTitle')}</CardTitle>
          <CardDescription>{t('wg.networkList.headerDesc')}</CardDescription>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="secondary" className="px-3 py-1">
            {t('wg.networkList.headerTotal', { count: total })}
          </Badge>
          <Button variant="outline" size="sm" onClick={() => refetch()} disabled={isFetching}>
            <RotateCcw className="h-4 w-4" />
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <DataTable table={table} columns={columns} toolbar={toolbar} />
      </CardContent>
    </Card>
  )
}
