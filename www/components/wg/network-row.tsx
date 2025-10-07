"use client"

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
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
import { toast } from 'sonner'
import { MoreHorizontal, ArrowUpRight } from 'lucide-react'
import { deleteNetwork } from '@/api/wg'
import { DeleteNetworkRequest } from '@/lib/pb/api_wg'
import NetworkEditDialog from './network-edit-dialog'
import { ColumnDef } from '@tanstack/react-table'
import { NetworkRow } from './network-list'
import { useRouter } from 'next/router'
import { useTranslation } from 'react-i18next'
import { TFunction } from 'i18next'

export function createNetworkColumns({ onChanged, t }: { onChanged?: () => void; t: TFunction }): ColumnDef<NetworkRow>[] {
	return [
		{
			accessorKey: 'name',
			meta: { label: 'wg.networkField.name' },
			header: () => t('wg.networkField.name'),
			cell: ({ row }) => {
				// eslint-disable-next-line react-hooks/rules-of-hooks
				const router = useRouter()
				return (
					<div className="flex items-center gap-2">
						<span className="font-medium truncate max-w-[180px]" title={row.original.name}>
							{row.original.name}
						</span>
						<Button
							variant="ghost"
							size="icon"
							className="h-6 w-6 rounded-full"
							onClick={(e) => {
								e.preventDefault()
								e.stopPropagation()
								router.push({ pathname: '/wg/network-detail', query: { networkId: row.original.id } })
							}}
							aria-label={t('wg.networkActions.view')}
						>
							<ArrowUpRight className="h-3.5 w-3.5" />
						</Button>
					</div>
				)
			}
		},
		{
			accessorKey: 'cidr',
			meta: { label: 'wg.networkField.cidr' },
			header: () => t('wg.networkField.cidr'),
			cell: ({ row }) => row.original.cidr
		},
		{
			id: 'acl',
			meta: { label: 'wg.networkField.acl' },
			header: () => t('wg.networkField.acl'),
			cell: ({ row }) => {
				const aclCount = row.original.acl?.acls?.length ?? 0
				return <span className="text-sm text-muted-foreground">{aclCount}</span>
			},
		},
		{
			id: 'actions',
			cell: ({ row }) => <NetworkActions row={row.original} onChanged={onChanged} />,
		},
	]
}

export function NetworkActions({ row, onChanged }: { row: NetworkRow; onChanged?: () => void }) {
	const { t } = useTranslation()
	const [openEdit, setOpenEdit] = useState(false)

	const onDelete = async () => {
		try {
			await deleteNetwork(DeleteNetworkRequest.create({ id: row.id }))
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
						<DropdownMenuLabel>{t('wg.networkActions.title')}</DropdownMenuLabel>
						<DropdownMenuItem
							onClick={(e) => {
								e.preventDefault()
								e.stopPropagation()
								setOpenEdit(true)
							}}
						>
							{t('wg.networkActions.edit')}
						</DropdownMenuItem>
						<DropdownMenuSeparator />
						<DialogTrigger asChild>
							<DropdownMenuItem className="text-destructive">
								{t('wg.networkActions.delete')}
							</DropdownMenuItem>
						</DialogTrigger>
					</DropdownMenuContent>
				</DropdownMenu>

				<DialogContent>
					<DialogHeader>
						<DialogTitle>{t('wg.networkActions.delete')}</DialogTitle>
						<DialogDescription>{t('wg.networkDelete.confirm')}</DialogDescription>
					</DialogHeader>
					<DialogFooter>
						<DialogClose asChild>
							<Button variant="outline">
								{t('common.cancel')}
							</Button>
						</DialogClose>
						<DialogClose asChild>
							<Button variant="destructive" onClick={onDelete}>
								{t('wg.networkActions.delete')}
							</Button>
						</DialogClose>
					</DialogFooter>
				</DialogContent>
			</Dialog>
			<NetworkEditDialog network={row.origin} onSaved={onChanged} open={openEdit} onOpenChange={setOpenEdit} />
		</>
	)
}