"use client"

import React, { useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import {
  getCoreRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  getFilteredRowModel,
  useReactTable,
  SortingState,
  PaginationState,
  ColumnFiltersState,
} from '@tanstack/react-table'
import { DataTable } from '@/components/base/data_table'
import { listWireGuards } from '@/api/wg'
import { ListWireGuardsRequest } from '@/lib/pb/api_wg'
import { WireGuardConfig } from '@/lib/pb/types_wg'
import { WireGuardColumns } from './wireguard-row'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { DropdownMenu, DropdownMenuCheckboxItem, DropdownMenuContent, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { Columns3, RotateCcw } from 'lucide-react'
import JoinNetworkDialog from './join-network-dialog'
import { useRouter } from 'next/router'

export type WireGuardTableSchema = {
  id: number
  interfaceName: string
  networkId?: number
  clientId: string
  localAddress: string
  listenPort?: number
  tags?: string[]
  origin: WireGuardConfig
}

export function WireGuardList({ clientId, networkId, keyword, onChanged }: { clientId?: string; networkId?: number; keyword?: string; onChanged?: (wireguards: WireGuardConfig[]) => void }) {
  const { t } = useTranslation()
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([])
  const [{ pageIndex, pageSize }, setPagination] = React.useState<PaginationState>({ pageIndex: 0, pageSize: 10 })
  const [columnVisibility, setColumnVisibility] = React.useState<Record<string, boolean>>({})
  const [refreshKey, setRefreshKey] = React.useState<number>(0)
  const router = useRouter()
  const { data, isFetching, refetch } = useQuery({
    queryKey: ['listWireGuards', clientId, networkId, keyword, pageIndex, pageSize, refreshKey],
    queryFn: () =>
      listWireGuards(
        ListWireGuardsRequest.create({
          page: pageIndex + 1,
          pageSize,
          clientId: clientId || undefined,
          networkId: networkId || undefined,
          keyword: keyword || undefined,
        }),
      ),
    placeholderData: keepPreviousData,
  })


  const rows: WireGuardTableSchema[] = (data?.wireguardConfigs ?? []).map((w) => ({
    id: w.id!,
    interfaceName: w.interfaceName ?? '',
    networkId: w.networkId,
    clientId: w.clientId ?? '',
    localAddress: w.localAddress ?? '',
    listenPort: w.listenPort,
    tags: w.tags,
    origin: w,
  }))

  const total = data?.total ?? 0

  const handleMutated = React.useCallback(() => {
    setRefreshKey((x) => x + 1)
  }, [])

  const table = useReactTable({
    data: rows,
    columns: WireGuardColumns,
    state: { sorting, pagination: { pageIndex, pageSize }, columnFilters, columnVisibility },
    manualPagination: true,
    pageCount: Math.ceil((data?.total ?? 0) / pageSize),
    onSortingChange: setSorting,
    onPaginationChange: setPagination,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    meta: {
      onChanged: handleMutated,
    },
  })

  const toolbar = (
    <div className="flex items-center gap-2">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline" size="sm" className="gap-2">
            <Columns3 className="h-4 w-4" />
            {t('table.columns')}
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="min-w-[200px]">
          {table.getAllLeafColumns().map((column) => {
            if (!column.getCanHide()) return null
            // const labelKey = column.columnDef.meta?.label as string | undefined
            return (
              <DropdownMenuCheckboxItem
                key={column.id}
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
      <Button variant="outline" size="sm" onClick={() => refetch()} disabled={isFetching}>
        <RotateCcw className="h-4 w-4" />
      </Button>
      <JoinNetworkDialog networkId={networkId} clientId={clientId} onJoined={() => setRefreshKey((x) => x + 1)}>
        <Button size="sm">{t('wg.joinNetwork.label')}</Button>
      </JoinNetworkDialog>
    </div>
  )

  return (
    <div className="grid gap-4 grid-cols-1">
      <Card>
        <CardHeader className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div className="space-y-1">
            <CardTitle>{t('wg.wireguardList.headerTitle')}</CardTitle>
            <CardDescription>{t('wg.wireguardList.headerDesc')}</CardDescription>
          </div>
          <div className="flex items-center gap-3">
            <Badge variant="secondary">{t('wg.wireguardList.headerTotal', { count: total })}</Badge>
            {toolbar}
          </div>
        </CardHeader>
      </Card>
      <DataTable table={table} columns={WireGuardColumns} />
    </div>
  )
}