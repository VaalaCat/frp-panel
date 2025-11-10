'use client'

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
import { MoreHorizontal } from 'lucide-react'
import { deleteEndpoint } from '@/api/wg'
import { DeleteEndpointRequest } from '@/lib/pb/api_wg'
import EndpointEditDialog from './endpoint-edit-dialog'
import { ColumnDef, Row } from '@tanstack/react-table'
import { EndpointTableSchema } from './endpoint-list'
import { Endpoint } from '@/lib/pb/types_wg'
import { useRouter } from 'next/router'
import { ArrowUpRight } from 'lucide-react'

export function createEndpointColumns({ onChanged }: { onChanged?: () => void }): ColumnDef<EndpointTableSchema>[] {
  return [
    {
      accessorKey: 'type',
      header: function Header() {
        const { t } = useTranslation()
        return t('wg.endpointTable.type')
      },
      cell: ({ row }: { row: Row<EndpointTableSchema> }) => {
        return <span className="text-sm font-mono nowrap">{row.original.type || 'udp'}</span>
      },
    },
    {
      accessorKey: 'uri',
      header: function Header() {
        const { t } = useTranslation()
        return t('wg.endpointTable.uri')
      },
      cell: ({ row }: { row: Row<EndpointTableSchema> }) => {
        return <span className="text-sm font-mono nowrap">{row.original.uri || ''}</span>
      },
    },
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
      cell: ({ row }) => (
        <EndpointActions clientId={row.original.clientId} endpoint={row.original.origin} onChanged={onChanged} />
      ),
    },
  ]
}

export default function EndpointActions({
  clientId,
  endpoint,
  onChanged,
}: {
  clientId: string
  endpoint: Endpoint
  onChanged?: () => void
}) {
  const { t } = useTranslation()
  const [openEdit, setOpenEdit] = useState(false)

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
      <Dialog>
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
                setOpenEdit(true)
              }}
            >
              {t('wg.endpointActions.edit')}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DialogTrigger asChild>
              <DropdownMenuItem className="text-destructive">{t('wg.endpoint.delete')}</DropdownMenuItem>
            </DialogTrigger>
          </DropdownMenuContent>
        </DropdownMenu>

        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t('wg.endpoint.delete')}</DialogTitle>
            <DialogDescription>{t('wg.endpointDelete.confirm')}</DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline">{t('common.cancel')}</Button>
            </DialogClose>
            <DialogClose asChild>
              <Button variant="destructive" onClick={onDelete}>
                {t('wg.endpoint.delete')}
              </Button>
            </DialogClose>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      <EndpointEditDialog
        clientId={clientId}
        endpoint={endpoint}
        onSaved={onChanged}
        open={openEdit}
        onOpenChange={setOpenEdit}
      />
    </>
  )
}
