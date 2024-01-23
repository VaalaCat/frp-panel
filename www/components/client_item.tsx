import { ColumnDef, Table } from '@tanstack/react-table'
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
import { useToast } from './ui/use-toast'
import React, { useState } from 'react'
import { ExecCommandStr, LinuxInstallCommand, WindowsInstallCommand } from '@/lib/consts'
import { useMutation, useQuery } from '@tanstack/react-query'
import { deleteClient, listClient } from '@/api/client'
import { useRouter } from 'next/router'
import { useStore } from '@nanostores/react'
import { $platformInfo } from '@/store/user'
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover'
import { getClientsStatus } from '@/api/platform'
import { ClientType } from '@/lib/pb/common'
import { ClientStatus, ClientStatus_Status } from '@/lib/pb/api_master'

export type ClientTableSchema = {
  id: string
  status: 'invalid' | 'valid'
  secret: string
  info?: string
  config?: string
}

export const columns: ColumnDef<ClientTableSchema>[] = [
  {
    accessorKey: 'id',
    header: 'ID(点击查看安装命令)',
    cell: ({ row }) => {
      return <ClientID client={row.original} />
    },
  },
  {
    accessorKey: 'status',
    header: '是否配置',
    cell: ({ row }) => {
      const client = row.original
      return (
        <div className={`font-medium ${client.status === 'valid' ? 'text-green-500' : 'text-red-500'} min-w-12`}>
          {
            {
              valid: '已配置',
              invalid: '未配置',
            }[client.status]
          }
        </div>
      )
    },
  },
  {
    accessorKey: 'info',
    header: '运行信息',
    cell: ({ row }) => {
      const client = row.original
      return <ClientInfo client={client} />
    },
  },
  {
    accessorKey: 'secret',
    header: '连接密钥(点击查看启动命令)',
    cell: ({ row }) => {
      const client = row.original
      return <ClientSecret client={client} />
    },
  },
  {
    id: 'action',
    cell: ({ row, table }) => {
      const client = row.original
      return <ClientActions client={client} table={table} />
    },
  },
]

export const ClientID = ({ client }: { client: ClientTableSchema }) => {
  const platformInfo = useStore($platformInfo)
  return (
    <Popover>
      <PopoverTrigger asChild>
        <div className="font-mono">{client.id}</div>
      </PopoverTrigger>
      <PopoverContent className="w-fit overflow-auto max-w-64">
        <div>Linux安装到systemd</div>
        <div className="p-2 border rounded font-mono w-fit">
          {platformInfo === undefined ? '获取平台信息失败' : LinuxInstallCommand('client', client, platformInfo)}
        </div>
        {/* <div>Windows</div>
            <div className="p-2 border rounded font-mono w-fit">
                {platformInfo === undefined ? "获取平台信息失败" : WindowsInstallCommand("client", client, platformInfo)}
            </div> */}
      </PopoverContent>
    </Popover>
  )
}

export const ClientInfo = ({ client }: { client: ClientTableSchema }) => {
  const clientsInfo = useQuery({
    queryKey: ['getClientsStatus', [client.id]],
    queryFn: async () => {
      return await getClientsStatus({
        clientIds: [client.id],
        clientType: ClientType.FRPC,
      })
    },
  })

  const trans = (info: ClientStatus | undefined) => {
    let statusText: '在线' | '离线' | '错误' | '未知' = '未知'
    if (info === undefined) {
      return statusText
    }
    if (info.status === ClientStatus_Status.ONLINE) {
      statusText = '在线'
    } else if (info.status === ClientStatus_Status.OFFLINE) {
      statusText = '离线'
    } else if (info.status === ClientStatus_Status.ERROR) {
      statusText = '错误'
    }
    return statusText
  }

  const infoColor =
    clientsInfo.data?.clients[client.id]?.status === ClientStatus_Status.ONLINE ? 'text-green-500' : 'text-red-500'

  return (
    <div className={`p-2 border rounded font-mono w-fit ${infoColor}`}>
      {`${clientsInfo.data?.clients[client.id].ping}ms, ${trans(clientsInfo.data?.clients[client.id])}`}
    </div>
  )
}

export const ClientSecret = ({ client }: { client: ClientTableSchema }) => {
  const [showSecrect, setShowSecrect] = useState<boolean>(false)
  const fakeSecret = Array.from({ length: client.secret.length })
    .map(() => '*')
    .join('')
  const platformInfo = useStore($platformInfo)
  const { toast } = useToast()
  return (
    <Popover>
      <PopoverTrigger asChild>
        <div
          onMouseEnter={() => setShowSecrect(true)}
          onMouseLeave={() => setShowSecrect(false)}
          className="font-medium hover:rounded hover:bg-slate-100 p-2 font-mono whitespace-nowrap"
        >
          {showSecrect ? client.secret : fakeSecret}
        </div>
      </PopoverTrigger>
      <PopoverContent className="w-fit overflow-auto max-w-80">
        <div>运行命令(需要<a className='text-blue-500' href='https://github.com/VaalaCat/frp-panel/releases'>点击这里</a>自行下载文件)</div>
        <div className="p-2 border rounded font-mono w-fit">
          {platformInfo === undefined ? '获取平台信息失败' : ExecCommandStr('client', client, platformInfo)}
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
  const { toast } = useToast()
  const router = useRouter()
  const platformInfo = useStore($platformInfo)
  const fetchDataOptions = {
    pageIndex: table.getState().pagination.pageIndex,
    pageSize: table.getState().pagination.pageSize,
  }

  const dataQuery = useQuery({
    queryKey: ['listClient', fetchDataOptions],
    queryFn: async () => {
      return await listClient({
        page: fetchDataOptions.pageIndex + 1,
        pageSize: fetchDataOptions.pageSize,
      })
    },
  })

  const removeClient = useMutation({
    mutationFn: deleteClient,
    onSuccess: () => {
      toast({ description: '删除成功' })
      dataQuery.refetch()
    },
    onError: () => {
      toast({ description: '删除失败' })
    },
  })

  return (
    <Dialog>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" className="h-8 w-8 p-0">
            <span className="sr-only">打开菜单</span>
            <MoreHorizontal className="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuLabel>操作</DropdownMenuLabel>
          <DropdownMenuItem
            onClick={() => {
              try {
                if (platformInfo) {
                  navigator.clipboard.writeText(ExecCommandStr('client', client, platformInfo))
                  toast({ description: '复制成功，如果复制不成功，请点击ID字段手动复制' })
                } else {
                  toast({ description: '获取平台信息失败，如果复制不成功，请点击ID字段手动复制' })
                }
              } catch (error) {
                toast({ description: '获取平台信息失败，如果复制不成功，请点击ID字段手动复制' })
              }
            }}
          >
            复制启动命令(也可点击列表中的密钥查看)
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            onClick={() => {
              router.push({ pathname: '/clientedit', query: { clientID: client.id } })
            }}
          >
            修改客户端配置
          </DropdownMenuItem>
          <DialogTrigger asChild>
            <DropdownMenuItem className="text-destructive">删除</DropdownMenuItem>
          </DialogTrigger>
        </DropdownMenuContent>
      </DropdownMenu>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>确定删除该客户端?</DialogTitle>
          <DialogDescription>
            <p className="text-destructive">此操作无法撤消。您确定要永久从我们的服务器中删除该客户端?</p>
            <p className="text-gray-500 border-l-4 border-gray-500 pl-4 py-2">
              删除后运行中的客户端将无法通过现有参数再次连接，如果您需要删除客户端对外的连接，可以选择清空配置
            </p>
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose asChild>
            <Button type="submit" onClick={() => removeClient.mutate({ clientId: client.id })}>
              确定
            </Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
