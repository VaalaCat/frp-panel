import React from 'react'
import { useRouter } from 'next/router'
import { useQuery, useMutation, keepPreviousData } from '@tanstack/react-query'
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
import { useStore } from '@nanostores/react'

import { listWorkers } from '@/api/worker'
import { Worker as PbWorker } from '@/lib/pb/common'
import { DataTable } from '../base/data_table'
import { WorkerTableSchema, columns as workerColumnsDef } from './worker_item'
import { $workerTableRefetchTrigger } from '@/store/refetch-trigger'

export interface WorkerListProps {
  initialWorkers: PbWorker[]
  initialTotal: number
  triggerRefetch?: string
  keyword?: string
}

export const WorkerList: React.FC<WorkerListProps> = ({ initialWorkers, initialTotal, triggerRefetch, keyword }) => {
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([])
  const [{ pageIndex, pageSize }, setPagination] = React.useState<PaginationState>({
    pageIndex: 0,
    pageSize: 10,
  })
  const globalTrigger = useStore($workerTableRefetchTrigger)

  const fetchOptions = { pageIndex, pageSize, triggerRefetch, globalTrigger, keyword }

  const { data, isFetching } = useQuery({
    queryKey: ['listWorkers', fetchOptions],
    queryFn: () =>
      listWorkers({
        page: pageIndex + 1,
        pageSize,
        keyword,
      }),
    placeholderData: keepPreviousData,
  })

  const dataRows: WorkerTableSchema[] =
    data?.workers.map((w) => ({
      workerId: w.workerId ?? '',
      name: w.name ?? '',
      userId: w.userId ?? 0,
      tenantId: w.tenantId ?? 0,
      socketAddress: w.socket?.address ?? '',
      origin: w,
    })) ?? []

  const table = useReactTable({
    data: dataRows,
    columns: workerColumnsDef,
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

  return <DataTable table={table} columns={workerColumnsDef} />
}
