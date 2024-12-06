import { ColumnDef, Table, TableMeta } from '@tanstack/react-table'
import { MoreHorizontal } from 'lucide-react'
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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useToast } from '@/components/ui/use-toast'
import React, { useState } from 'react'
import { ClientEnvFile, ExecCommandStr, LinuxInstallCommand, WindowsInstallCommand } from '@/lib/consts'
import { useMutation, useQuery } from '@tanstack/react-query'
import { deleteClient, listClient } from '@/api/client'
import { useRouter } from 'next/router'
import { useStore } from '@nanostores/react'
import { $platformInfo } from '@/store/user'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { getClientsStatus } from '@/api/platform'
import { ClientType } from '@/lib/pb/common'
import { ClientStatus, ClientStatus_Status } from '@/lib/pb/api_master'
import { startFrpc, stopFrpc } from '@/api/frp'
import { Badge } from '../ui/badge'
import { Label } from '@/components/ui/label'
import { ClientDetail } from '../base/client_detail'
import { Input } from '../ui/input'
import { useTranslation } from 'react-i18next'

export type ClientTableSchema = {
  id: string
  status: 'invalid' | 'valid'
  secret: string
  stopped: boolean
  info?: string
  config?: string
}

export interface TableMetaType extends TableMeta<ClientTableSchema> {
  refetch: () => void
}

export const columns: ColumnDef<ClientTableSchema>[] = [
  {
    accessorKey: 'id',
    header: function Header() {
      const { t } = useTranslation()
      return t('client.id')
    },
    cell: ({ row }) => {
      return <ClientID client={row.original} />
    },
  },
  {
    accessorKey: 'status',
    header: function Header() {
      const { t } = useTranslation()
      return t('client.status')
    },
    cell: ({ row }) => {
      function Cell({ client }: { client: ClientTableSchema }) {
        const { t } = useTranslation()
        return (
          <div className={`font-medium ${client.status === 'valid' ? 'text-green-500' : 'text-red-500'} min-w-12`}>
            {client.status === 'valid' ? t('client.status_configured') : t('client.status_unconfigured')}
          </div>
        )
      }
      return <Cell client={row.original} />
    },
  },
  {
    accessorKey: 'info',
    header: function Header() {
      const { t } = useTranslation()
      return t('client.info')
    },
    cell: ({ row }) => {
      const client = row.original
      return <ClientInfo client={client} />
    },
  },
  {
    accessorKey: 'secret',
    header: function Header() {
      const { t } = useTranslation()
      return t('client.secret')
    },
    cell: ({ row }) => {
      const client = row.original
      return <ClientSecret client={client} />
    },
  },
  {
    id: 'action',
    cell: ({ row, table }) => {
      const client = row.original
      return <ClientActions client={client} table={table as Table<ClientTableSchema> & { options: { meta: TableMetaType } }} />
    },
  },
]

export const ClientID = ({ client }: { client: ClientTableSchema }) => {
  const { t } = useTranslation()
  const platformInfo = useStore($platformInfo)

  if (!platformInfo) {
    return (
      <Button variant="link" className="px-0">
        {client.id}
      </Button>
    )
  }

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="link" className="px-0 font-mono">
          {client.id}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-80">
        <div className="grid gap-4">
          <div className="space-y-2">
            <h4 className="font-medium leading-none">{t('client.install.title')}</h4>
            <p className="text-sm text-muted-foreground">{t('client.install.description')}</p>
          </div>
          <div className="grid gap-2">
            <div className="grid grid-cols-2 items-center gap-4">
              <Input
                readOnly
                value={WindowsInstallCommand('frpc', client, platformInfo)}
                className="flex-1"
              />
              <Button
                onClick={() => navigator.clipboard.writeText(WindowsInstallCommand('frpc', client, platformInfo))}
                disabled={!platformInfo}
                size="sm"
                variant="outline"
              >
                {t('client.install.windows')}
              </Button>
            </div>
            <div className="grid grid-cols-2 items-center gap-4">
              <Input
                readOnly
                value={LinuxInstallCommand('frpc', client, platformInfo)}
                className="flex-1"
              />
              <Button
                onClick={() => navigator.clipboard.writeText(LinuxInstallCommand('frpc', client, platformInfo))}
                disabled={!platformInfo}
                size="sm"
                variant="outline"
              >
                {t('client.install.linux')}
              </Button>
            </div>
          </div>
        </div>
      </PopoverContent>
    </Popover>
  )
}

export const ClientInfo = ({ client }: { client: ClientTableSchema }) => {
  const { t } = useTranslation()
  const { data: clientsStatus } = useQuery({
    queryKey: ['clientsStatus', client.id],
    queryFn: async () => {
      return await getClientsStatus({
        clientIds: [client.id],
        clientType: ClientType.FRPC,
      })
    },
  })

  const trans = (info: ClientStatus | undefined) => {
    let statusText: 'client.status_online' | 'client.status_offline' |
      'client.status_error' | 'client.status_pause' |
      'client.status_unknown' = 'client.status_unknown'
    if (info === undefined) {
      return statusText
    }
    if (info.status === ClientStatus_Status.ONLINE) {
      statusText = 'client.status_online'
      if (client.stopped) {
        statusText = 'client.status_pause'
      }
    } else if (info.status === ClientStatus_Status.OFFLINE) {
      statusText = 'client.status_offline'
    } else if (info.status === ClientStatus_Status.ERROR) {
      statusText = 'client.status_error'
    }
    return statusText
  }

  const infoColor =
    clientsStatus?.clients[client.id]?.status === ClientStatus_Status.ONLINE ? (
      client.stopped ? 'text-yellow-500' : 'text-green-500') : 'text-red-500'

  return (
    <div className="flex items-center gap-2 flex-row">
      <Badge variant={"secondary"} className={`p-2 border rounded font-mono w-fit ${infoColor} text-nowrap rounded-full h-6`}>
        {`${clientsStatus?.clients[client.id].ping}ms,${t(trans(clientsStatus?.clients[client.id]))}`}
      </Badge>
      {clientsStatus?.clients[client.id].version &&
        <ClientDetail clientStatus={clientsStatus?.clients[client.id]} />
      }
    </div>
  )
}

export const ClientSecret = ({ client }: { client: ClientTableSchema }) => {
  const { t } = useTranslation()
  const platformInfo = useStore($platformInfo)

  if (!platformInfo) {
    return (
      <Button variant="link" className="px-0">
        {client.secret}
      </Button>
    )
  }

  return (
    <Popover>
      <PopoverTrigger asChild>
        <div className="group relative cursor-pointer inline-block font-mono text-nowrap">
          <span className="opacity-0 group-hover:opacity-100 transition-opacity duration-200">
            {client.secret}
          </span>
          <span className="absolute inset-0 opacity-100 group-hover:opacity-0 transition-opacity duration-200">
            {'*'.repeat(client.secret.length)}
          </span>
        </div>
      </PopoverTrigger>
      <PopoverContent className="w-[32rem] max-w-[95vw]">
        <div className="grid gap-4">
          <div className="space-y-2">
            <h4 className="font-medium leading-none">{t('client.start.title')}</h4>
            <p className="text-sm text-muted-foreground">
              {t('client.install.description')} (<a className='text-blue-500' href='https://github.com/VaalaCat/frp-panel/releases' target="_blank" rel="noopener noreferrer">{t('common.download')}</a>)
            </p>
          </div>
          <div className="grid gap-2">
            <pre className="bg-muted p-3 rounded-md font-mono text-sm overflow-x-auto whitespace-pre-wrap break-all">
              {ExecCommandStr('client', client, platformInfo)}
            </pre>
            <Button
              size="sm"
              variant="outline"
              className="w-full"
              onClick={() => navigator.clipboard.writeText(ExecCommandStr('client', client, platformInfo))}
              disabled={!platformInfo}
            >
              {t('common.copy')}
            </Button>
          </div>
        </div>
      </PopoverContent>
    </Popover>
  )
}

export interface ClientItemProps {
  client: ClientTableSchema
  table: Table<ClientTableSchema>
}

export const ClientActions: React.FC<ClientItemProps> = ({ client, table }) => {
  const { t } = useTranslation()
  const { toast } = useToast()
  const router = useRouter()
  const platformInfo = useStore($platformInfo)

  // placeholder for refetch
  const refetchList = () => {}

  const removeClient = useMutation({
    mutationFn: deleteClient,
    onSuccess: () => {
      toast({ description: t('client.delete.success') })
      refetchList()
    },
    onError: () => {
      toast({ description: t('client.delete.failed') })
    },
  })

  const stopClient = useMutation({
    mutationFn: stopFrpc,
    onSuccess: () => {
      toast({ description: t('client.operation.stop_success') })
      refetchList()
    },
    onError: () => {
      toast({ description: t('client.operation.stop_failed') })
    },
  })

  const startClient = useMutation({
    mutationFn: startFrpc,
    onSuccess: () => {
      toast({ description: t('client.operation.start_success') })
      refetchList()
    },
    onError: () => {
      toast({ description: t('client.operation.start_failed') })
    },
  })

  const createAndDownloadFile = (fileName: string, content: string) => {
    const aTag = document.createElement('a')
    const blob = new Blob([content])
    aTag.download = fileName
    aTag.href = URL.createObjectURL(blob)
    aTag.click()
    URL.revokeObjectURL(aTag.href)
  }

  return (
    <Dialog>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" className="h-8 w-8 p-0">
            <span className="sr-only">{t('client.actions_menu.open_menu')}</span>
            <MoreHorizontal className="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuLabel>{t('client.actions_menu.title')}</DropdownMenuLabel>
          <DropdownMenuItem
            onClick={() => {
              try {
                if (platformInfo) {
                  navigator.clipboard.writeText(ExecCommandStr('client', client, platformInfo))
                  toast({ description: t('client.actions_menu.copy_success') })
                } else {
                  toast({ description: t('client.actions_menu.copy_failed') })
                }
              } catch (error) {
                toast({ description: t('client.actions_menu.copy_failed') })
              }
            }}
          >
            {t('client.actions_menu.copy_command')}
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            onClick={() => {
              router.push({ pathname: '/clientedit', query: { clientID: client.id } })
            }}
          >
            {t('client.actions_menu.edit_config')}
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={() => {
              try {
                if (platformInfo) {
                  createAndDownloadFile('.env', ClientEnvFile(client, platformInfo))
                }
              } catch (error) {
                toast({ description: t('client.actions_menu.download_failed') })
              }
            }}
          >
            {t('client.actions_menu.download_config')}
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={() => {
              router.push({ pathname: '/streamlog', query: { clientID: client.id, clientType: ClientType.FRPC.toString() } })
            }}
          >
            {t('client.actions_menu.realtime_log')}
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={() => {
              router.push({ pathname: '/console', query: { clientID: client.id, clientType: ClientType.FRPC.toString() } })
            }}
          >
            {t('client.actions_menu.remote_terminal')}
          </DropdownMenuItem>
          {!client.stopped && (
            <DropdownMenuItem className="text-destructive" onClick={() => stopClient.mutate({ clientId: client.id })}>
              {t('client.actions_menu.pause')}
            </DropdownMenuItem>
          )}
          {client.stopped && (
            <DropdownMenuItem onClick={() => startClient.mutate({ clientId: client.id })}>
              {t('client.actions_menu.resume')}
            </DropdownMenuItem>
          )}
          <DialogTrigger asChild>
            <DropdownMenuItem className="text-destructive">{t('client.actions_menu.delete')}</DropdownMenuItem>
          </DialogTrigger>
        </DropdownMenuContent>
      </DropdownMenu>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t('client.delete.title')}</DialogTitle>
          <DialogDescription>
            <p className="text-destructive">{t('client.delete.description')}</p>
            <p className="text-gray-500 border-l-4 border-gray-500 pl-4 py-2">
              {t('client.delete.warning')}
            </p>
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose asChild>
            <Button type="submit" onClick={() => removeClient.mutate({ clientId: client.id })}>
              {t('client.delete.confirm')}
            </Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
