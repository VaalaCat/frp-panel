"use client"

import React from 'react'
import { useTranslation } from 'react-i18next'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import { getCoreRowModel, getPaginationRowModel, getSortedRowModel, getFilteredRowModel, useReactTable, SortingState, PaginationState, ColumnFiltersState, ColumnDef, Row } from '@tanstack/react-table'
import { DataTable } from '@/components/base/data_table'
import { Button } from '@/components/ui/button'
import { listEndpoints } from '@/api/wg'
import { ListEndpointsRequest } from '@/lib/pb/api_wg'
import EndpointEditDialog from './endpoint-edit-dialog'
import { createEndpointColumns } from './endpoint-row'
import { Endpoint } from '@/lib/pb/types_wg'

export type EndpointTableSchema = {
  id: number
  host: string
  port: number
  origin: Endpoint
  clientId: string
}

export function EndpointList({ clientId, wireguardId, keyword }: { clientId?: string; wireguardId?: number; keyword?: string }) {
  const { t } = useTranslation()
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([])
  const [{ pageIndex, pageSize }, setPagination] = React.useState<PaginationState>({ pageIndex: 0, pageSize: 10 })
  const [refreshKey, setRefreshKey] = React.useState(0)
  const [openAdd, setOpenAdd] = React.useState(false)

  const { data } = useQuery({
    queryKey: ['listEndpoints', clientId, wireguardId, keyword, pageIndex, pageSize, refreshKey],
    queryFn: () =>
      listEndpoints(
        ListEndpointsRequest.create({
          page: pageIndex + 1,
          pageSize,
          clientId: clientId || undefined,
          wireguardId: wireguardId || undefined,
          keyword: keyword || undefined,
        }),
      ),
    placeholderData: keepPreviousData,
  })

  const rows: EndpointTableSchema[] = (data?.endpoints ?? []).map((e) => ({
    id: e.id!,
    host: e.host!,
    port: e.port!,
    clientId: e.clientId!,
    origin: e,
  }))

  const handleMutated = React.useCallback(() => {
    setRefreshKey((x) => x + 1)
  }, [])

  const columns = React.useMemo(() => createEndpointColumns({ onChanged: handleMutated }), [handleMutated])

  const table = useReactTable({
    data: rows,
    columns,
    state: { sorting, pagination: { pageIndex, pageSize }, columnFilters },
    manualPagination: true,
    pageCount: Math.ceil((data?.total ?? 0) / pageSize),
    onSortingChange: setSorting,
    onPaginationChange: setPagination,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
  })

  return (
    <div className="space-y-3">
      <EndpointEditDialog clientId={clientId || ""} onSaved={() => setRefreshKey((x) => x + 1)} open={openAdd} onOpenChange={setOpenAdd}>
        <Button size="sm">{t('wg.endpointCreate.button')}</Button>
      </EndpointEditDialog>
      <DataTable table={table} columns={columns} />
    </div>
  )
}
