import { ColumnDef, Row } from '@tanstack/react-table'
import { WireGuardTableSchema } from './wireguard-list'
import { useTranslation } from 'react-i18next'
import { useState } from 'react'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
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
import { Button } from '@/components/ui/button'
import { deleteWireGuard } from '@/api/wg'
import { DeleteWireGuardRequest } from '@/lib/pb/api_wg'
import { toast } from 'sonner'
import WireGuardEditDialog from './wireguard-edit-dialog'
import { MoreHorizontal, ArrowUpRight } from 'lucide-react'
import { WireGuardConfig } from '@/lib/pb/types_wg'
import { useRouter } from 'next/router'

export const WireGuardColumns: ColumnDef<WireGuardTableSchema>[] = [
	{
		accessorKey: 'id',
		meta: { label: 'wg.wireguard.list.columns.id' },
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.interface.id')
		},
		cell: ({ row }) => <span className="text-sm">{row.original.id}</span>,
	},
	{
		accessorKey: 'interfaceName',
		meta: { label: 'wg.wireguard.list.columns.interface' },
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.interface.name')
		},
		cell: ({ row }: { row: Row<WireGuardTableSchema> }) => {
			// eslint-disable-next-line react-hooks/rules-of-hooks
			const router = useRouter()
			// eslint-disable-next-line react-hooks/rules-of-hooks
			const { t } = useTranslation()
			return (
				<div className="flex items-center gap-2">
					<div className="flex flex-col">
						<span className="font-medium" title={row.original.interfaceName}>
							{row.original.interfaceName}
						</span>
						<span className="text-muted-foreground text-xs">{row.original.localAddress}</span>
					</div>
					<Button
						variant="ghost"
						size="icon"
						className="h-6 w-6 rounded-full"
						onClick={(e) => {
							e.preventDefault()
							e.stopPropagation()
							router.push({ pathname: '/wg/wireguard-detail', query: { id: row.original.id } })
						}}
						aria-label={t('wg.wireguardActions.view')}
					>
						<ArrowUpRight className="h-3.5 w-3.5" />
					</Button>
				</div>
			)
		},
	},
	{
		accessorKey: 'listenPort',
		meta: { label: 'wg.wireguard.list.columns.port' },
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.interface.port')
		},
		cell: ({ row }) => <span className="text-sm">{row.original.listenPort}</span>,
	},
	{
		accessorKey: 'clientId',
		meta: { label: 'wg.wireguard.list.columns.client' },
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.interface.client')
		},
		cell: ({ row }) => <span className="text-sm font-mono">{row.original.clientId}</span>,
	},
	{
		accessorKey: 'tags',
		meta: { label: 'wg.wireguard.list.columns.tags' },
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.interface.tags')
		},
		cell: ({ row }) => <span className="text-sm">{row.original.tags?.join(',')}</span>,
	},
	{
		accessorKey: 'networkId',
		meta: { label: 'wg.wireguard.list.columns.network' },
		header: function Header() {
			const { t } = useTranslation()
			return t('wg.interface.network')
		},
		cell: ({ row }) => {
			// eslint-disable-next-line react-hooks/rules-of-hooks
			const router = useRouter()
			// eslint-disable-next-line react-hooks/rules-of-hooks
			const { t } = useTranslation()
			return row.original.networkId ? (
				<Button
					variant="link"
					size="sm"
					className="px-0 text-sm"
					onClick={(e) => {
						e.preventDefault()
						e.stopPropagation()
						router.push({ pathname: '/wg/network-detail', query: { networkId: row.original.networkId } })
					}}
				>
					#{row.original.networkId}
				</Button>
			) : (
				<span className="text-sm text-muted-foreground">{t('wg.interface.network_unassigned')}</span>
			)
		},
	},
	{
		id: 'actions',
		cell: ({ row, table }) => <WireGuardActions clientId={row.original.clientId} wg={row.original.origin} onChanged={(table.options.meta as any)?.onChanged} />,
	},
]


function WireGuardActions({ clientId, wg, onChanged }: { clientId: string; wg: WireGuardConfig; onChanged?: () => void }) {
	const { t } = useTranslation()
	const [openEdit, setOpenEdit] = useState(false)
	const onDelete = async () => {
		if (!wg.id) return
		try {
			await deleteWireGuard(DeleteWireGuardRequest.create({ id: wg.id }))
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
						<DropdownMenuLabel>{t('wg.wireguardActions.title')}</DropdownMenuLabel>
						<DropdownMenuItem
							onClick={(e) => {
								e.preventDefault()
								e.stopPropagation()
								setOpenEdit(true)
							}}
						>
							{t('wg.wireguardActions.edit')}
						</DropdownMenuItem>
						<DropdownMenuSeparator />
						<DialogTrigger asChild>
							<DropdownMenuItem className="text-destructive">
								{t('wg.wireguardActions.delete')}
							</DropdownMenuItem>
						</DialogTrigger>
					</DropdownMenuContent>
				</DropdownMenu>

				<DialogContent>
					<DialogHeader>
						<DialogTitle>{t('wg.wireguardActions.delete')}</DialogTitle>
						<DialogDescription>{t('wg.wireguardDelete.confirm')}</DialogDescription>
					</DialogHeader>
					<DialogFooter>
						<DialogClose asChild>
							<Button variant="outline">
								{t('common.cancel')}
							</Button>
						</DialogClose>
						<DialogClose asChild>
							<Button variant="destructive" onClick={onDelete}>
								{t('wg.wireguardActions.delete')}
							</Button>
						</DialogClose>
					</DialogFooter>
				</DialogContent>
			</Dialog>
			<WireGuardEditDialog clientId={clientId} wg={wg} onUpdated={onChanged} open={openEdit} onOpenChange={setOpenEdit} />
		</>
	)
}

