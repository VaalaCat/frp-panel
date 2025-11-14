"use client"

import React from 'react'
import { useTranslation } from 'react-i18next'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import { getCoreRowModel, getPaginationRowModel, getSortedRowModel, getFilteredRowModel, useReactTable, SortingState, PaginationState, ColumnFiltersState, ColumnDef, Row } from '@tanstack/react-table'
import { useRouter } from 'next/router'
import { DataTable } from '@/components/base/data_table'
import { Button } from '@/components/ui/button'
import { deleteWireGuardLink, listWireGuardLinks } from '@/api/wg'
import { DeleteWireGuardLinkRequest, ListWireGuardLinksRequest } from '@/lib/pb/api_wg'
import WireGuardLinkEditDialog from './wireguard-link-edit-dialog'
import { WireGuardLink } from '@/lib/pb/types_wg'
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuItem } from '@/components/ui/dropdown-menu'
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'
import { MoreHorizontal } from 'lucide-react'
import { toast } from 'sonner'

export type WireGuardLinkRow = {
	id: number
	fromWireguardId: number
	toWireguardId: number
	toEndpoint?: string
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
		cell: function Cell({ row }: { row: Row<WireGuardLinkRow> }) {
			// eslint-disable-next-line react-hooks/rules-of-hooks
			const router = useRouter()
			return (
				<Button
					variant="link"
					size="sm"
					className="px-0 text-sm"
					onClick={(e) => {
						e.preventDefault()
						e.stopPropagation()
						router.push(`/wg/wireguard-detail?id=${row.original.fromWireguardId}`)
					}}
				>
					#{row.original.fromWireguardId}
				</Button>
			)
		},
	},
	{
		accessorKey: 'toWireguardId',
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.link.to')
		},
		cell: function Cell({ row }: { row: Row<WireGuardLinkRow> }) {
			// eslint-disable-next-line react-hooks/rules-of-hooks
			const router = useRouter()
			return (
				<Button
					variant="link"
					size="sm"
					className="px-0 text-sm"
					onClick={(e) => {
						e.preventDefault()
						e.stopPropagation()
						router.push(`/wg/wireguard-detail?id=${row.original.toWireguardId}`)
					}}
				>
					#{row.original.toWireguardId}
				</Button>
			)
		},
	},
	{
		accessorKey: 'toEndpoint',
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.link.toEndpoint')
		},
		cell: ({ row }) => {
			const endpoint = row.original.toEndpoint
			return endpoint ? (
				<span className="text-xs font-mono text-muted-foreground">{endpoint}</span>
			) : (
				<span className="text-xs text-muted-foreground italic">-</span>
			)
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
		cell: ({ row, table }) => <WireGuardLinkActions link={row.original.origin} onChanged={(table.options.meta as any)?.onChanged} />,
	},
]

function WireGuardLinkActions({ link, onChanged }: { link: WireGuardLink; onChanged?: () => void }) {
	const { t } = useTranslation()
	const [openEdit, setOpenEdit] = React.useState(false)
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
			<Dialog>
				<DropdownMenu>
					<DropdownMenuTrigger asChild>
						<Button variant="ghost" size="icon">
							<MoreHorizontal className="h-4 w-4" />
						</Button>
					</DropdownMenuTrigger>
					<DropdownMenuContent align="end">
						<DropdownMenuItem onClick={() => setOpenEdit(true)}>{t('wg.linkActions.edit')}</DropdownMenuItem>
						<DialogTrigger asChild>
							<DropdownMenuItem className="text-destructive">
								{t('wg.linkActions.delete')}
							</DropdownMenuItem>
						</DialogTrigger>
					</DropdownMenuContent>
				</DropdownMenu>

				<DialogContent>
					<DialogHeader>
						<DialogTitle>{t('wg.linkActions.delete')}</DialogTitle>
						<DialogDescription>{t('wg.linkDelete.confirm')}</DialogDescription>
					</DialogHeader>
					<DialogFooter>
						<DialogClose asChild>
							<Button variant="outline">
								{t('common.cancel')}
							</Button>
						</DialogClose>
						<DialogClose asChild>
							<Button variant="destructive" onClick={onDelete}>
								{t('wg.linkActions.delete')}
							</Button>
						</DialogClose>
					</DialogFooter>
				</DialogContent>
			</Dialog>
			<WireGuardLinkEditDialog link={link} onSaved={onChanged} open={openEdit} onOpenChange={setOpenEdit} />
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
		toEndpoint: l.toEndpoint ? `${l.toEndpoint.host}:${l.toEndpoint.port}` : undefined,
		upBandwidthMbps: l.upBandwidthMbps,
		downBandwidthMbps: l.downBandwidthMbps,
		latencyMs: l.latencyMs,
		active: !!l.active,
		origin: l,
	}))

	const handleMutated = React.useCallback(() => {
		setRefreshKey((x) => x + 1)
	}, [])

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
		meta: {
			onChanged: handleMutated,
		},
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


