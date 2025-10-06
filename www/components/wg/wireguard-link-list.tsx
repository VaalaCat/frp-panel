"use client"

import React from 'react'
import { useTranslation } from 'react-i18next'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import { getCoreRowModel, getPaginationRowModel, getSortedRowModel, getFilteredRowModel, useReactTable, SortingState, PaginationState, ColumnFiltersState, ColumnDef, Row } from '@tanstack/react-table'
import { DataTable } from '@/components/base/data_table'
import { Button } from '@/components/ui/button'
import { deleteWireGuardLink, listWireGuardLinks } from '@/api/wg'
import { DeleteWireGuardLinkRequest, ListWireGuardLinksRequest } from '@/lib/pb/api_wg'
import WireGuardLinkEditDialog from './wireguard-link-edit-dialog'
import { WireGuardLink } from '@/lib/pb/types_wg'
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuItem } from '@/components/ui/dropdown-menu'
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { MoreHorizontal } from 'lucide-react'
import { toast } from 'sonner'

export type WireGuardLinkRow = {
	id: number
	fromWireguardId: number
	toWireguardId: number
	upBandwidthMbps?: number
	downBandwidthMbps?: number
	latencyMs?: number
	active: boolean
	origin: WireGuardLink
}

export const WireGuardLinkColumns: ColumnDef<WireGuardLinkRow>[] = [
	{
		accessorKey: 'fromWireguardId',
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.link.from')
		},
	},
	{
		accessorKey: 'toWireguardId',
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.link.to')
		},
	},
	{
		accessorKey: 'upBandwidthMbps',
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.link.up_bw')
		},
	},
	{
		accessorKey: 'downBandwidthMbps',
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.link.down_bw')
		},
	},
	{
		accessorKey: 'latencyMs',
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.link.latency')
		},
	},
	{
		accessorKey: 'active',
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.link.active')
		},
		cell: ({ row }) => {
			// eslint-disable-next-line react-hooks/rules-of-hooks
			const { t } = useTranslation()
			return <span className="text-sm">{row.original.active ? t('wg.linkState.active') : t('wg.linkState.inactive')}</span>
		},
	},
	{
		id: 'actions',
		cell: ({ row }) => <WireGuardLinkActions link={row.original.origin} />,
	},
]

function WireGuardLinkActions({ link, onChanged }: { link: WireGuardLink; onChanged?: () => void }) {
	const { t } = useTranslation()
	const [openEdit, setOpenEdit] = React.useState(false)
	const [openDelete, setOpenDelete] = React.useState(false)
	const onDelete = async () => {
		if (!link.id) return
		try {
			await deleteWireGuardLink(DeleteWireGuardLinkRequest.create({ id: link.id }))
			toast.success(t('common.success'))
			onChanged?.()
		} catch (e: any) {
			toast.error(e.message)
		}
	}

	return (
		<>
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<Button variant="ghost" size="icon">
						<MoreHorizontal className="h-4 w-4" />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end">
					<DropdownMenuItem onClick={() => setOpenEdit(true)}>{t('wg.linkActions.edit')}</DropdownMenuItem>
					<DropdownMenuItem
						onClick={() => setOpenDelete(true)}
						className="text-destructive"
					>
						{t('wg.linkActions.delete')}
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
			<WireGuardLinkEditDialog link={link} onSaved={onChanged} open={openEdit} onOpenChange={setOpenEdit} />
			<AlertDialog open={openDelete} onOpenChange={setOpenDelete}>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle>{t('wg.linkActions.delete')}</AlertDialogTitle>
						<AlertDialogDescription>{t('wg.linkDelete.confirm')}</AlertDialogDescription>
					</AlertDialogHeader>
					<AlertDialogFooter>
						<AlertDialogCancel>{t('common.cancel')}</AlertDialogCancel>
						<AlertDialogAction
							className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
							onClick={onDelete}
						>
							{t('wg.linkActions.delete')}
						</AlertDialogAction>
					</AlertDialogFooter>
				</AlertDialogContent>
			</AlertDialog>
		</>
	)
}

export function WireGuardLinkList({ networkId, keyword }: { networkId?: number; keyword?: string }) {
	const { t } = useTranslation()
	const [sorting, setSorting] = React.useState<SortingState>([])
	const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([])
	const [{ pageIndex, pageSize }, setPagination] = React.useState<PaginationState>({ pageIndex: 0, pageSize: 10 })
	const [refreshKey, setRefreshKey] = React.useState(0)
	const [openAdd, setOpenAdd] = React.useState(false)

	const { data } = useQuery({
		queryKey: ['listWireGuardLinks', networkId, keyword, pageIndex, pageSize, refreshKey],
		queryFn: () =>
			listWireGuardLinks(
				ListWireGuardLinksRequest.create({
					page: pageIndex + 1,
					pageSize,
					networkId,
					keyword,
				}),
			),
		placeholderData: keepPreviousData,
	})

	const rows: WireGuardLinkRow[] = (data?.wireguardLinks ?? []).map((l) => ({
		id: l.id!,
		fromWireguardId: l.fromWireguardId!,
		toWireguardId: l.toWireguardId!,
		upBandwidthMbps: l.upBandwidthMbps,
		downBandwidthMbps: l.downBandwidthMbps,
		latencyMs: l.latencyMs,
		active: !!l.active,
		origin: l,
	}))

	const table = useReactTable({
		data: rows,
		columns: WireGuardLinkColumns,
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
			<WireGuardLinkEditDialog onSaved={() => setRefreshKey((x) => x + 1)} open={openAdd} onOpenChange={setOpenAdd}>
				<Button size="sm">{t('wg.linkCreate.button')}</Button>
			</WireGuardLinkEditDialog>
			<DataTable table={table} columns={WireGuardLinkColumns} />
		</div>
	)
}


