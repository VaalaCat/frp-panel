"use client"

import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
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
import { toast } from 'sonner'
import { MoreHorizontal } from 'lucide-react'
import { deleteEndpoint } from '@/api/wg'
import { DeleteEndpointRequest } from '@/lib/pb/api_wg'
import EndpointEditDialog from './endpoint-edit-dialog'
import { ColumnDef, Row } from '@tanstack/react-table'
import { EndpointTableSchema } from './endpoint-list'
import { Endpoint } from '@/lib/pb/types_wg'
import { useRouter } from 'next/router'
import { ArrowUpRight } from 'lucide-react'

export const EndpointColumns: ColumnDef<EndpointTableSchema>[] = [
	{
		accessorKey: 'host',
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.endpointTable.host')
		},
		cell: ({ row }: { row: Row<EndpointTableSchema> }) => {
			// eslint-disable-next-line react-hooks/rules-of-hooks
			const router = useRouter()
			// eslint-disable-next-line react-hooks/rules-of-hooks
			const { t } = useTranslation()
			return (
				<div className="flex items-center gap-2">
					<span className="font-medium">
						{row.original.host}:{row.original.port}
					</span>
					<Button
						variant="ghost"
						size="icon"
						className="h-6 w-6 rounded-full"
						onClick={(e) => {
							e.preventDefault()
							e.stopPropagation()
							router.push({ pathname: '/wg/endpoint-detail', query: { id: row.original.id } })
						}}
						aria-label={t('wg.endpointActions.view')}
					>
						<ArrowUpRight className="h-3.5 w-3.5" />
					</Button>
				</div>
			)
		},
	},
	{
		accessorKey: 'wireguardId',
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.endpointTable.wireguardId')
		},
		cell: ({ row }: { row: Row<EndpointTableSchema> }) => (
			<div className="flex items-center gap-1">
				<span className="font-medium">{row.original.origin.wireguardId}</span>
			</div>
		),
	},
	{
		accessorKey: 'clientId',
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.endpointTable.clientId')
		},
	},
	{
		id: 'actions',
		cell: ({ row }) => <EndpointActions clientId={row.original.clientId} endpoint={row.original.origin} />,
	},
]


export default function EndpointActions({ clientId, endpoint, onChanged }: { clientId: string; endpoint: Endpoint; onChanged?: () => void }) {
	const { t } = useTranslation()
	const [openEdit, setOpenEdit] = useState(false)
	const [openDelete, setOpenDelete] = useState(false)

	const onDelete = async () => {
		if (!endpoint.id) return
		try {
			await deleteEndpoint(DeleteEndpointRequest.create({ id: endpoint.id }))
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
					<DropdownMenuLabel>{t('wg.endpointActions.title')}</DropdownMenuLabel>
					<DropdownMenuItem
						onClick={(e) => {
							e.preventDefault()
							e.stopPropagation()
							setOpenDelete(false)
							setOpenEdit(true)
						}}
					>
						{t('wg.endpointActions.edit')}
					</DropdownMenuItem>
					<DropdownMenuSeparator />
					<DropdownMenuItem
						className="text-destructive"
						onClick={(e) => {
							e.preventDefault()
							e.stopPropagation()
							setOpenEdit(false)
							setOpenDelete(true)
						}}
					>
						{t('wg.endpoint.delete')}
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
			<EndpointEditDialog clientId={clientId} endpoint={endpoint} onSaved={onChanged} open={openEdit} onOpenChange={setOpenEdit} />
			<AlertDialog open={openDelete} onOpenChange={setOpenDelete}>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle>{t('wg.endpoint.delete')}</AlertDialogTitle>
						<AlertDialogDescription>{t('wg.endpointDelete.confirm')}</AlertDialogDescription>
					</AlertDialogHeader>
					<AlertDialogFooter>
						<AlertDialogCancel>{t('common.cancel')}</AlertDialogCancel>
						<AlertDialogAction
							className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
							onClick={onDelete}
						>
							{t('wg.endpoint.delete')}
						</AlertDialogAction>
					</AlertDialogFooter>
				</AlertDialogContent>
			</AlertDialog>
		</>
	)
}


