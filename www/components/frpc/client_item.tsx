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
import { ClientDetail } from '../base/client_detail'
import { Input } from '../ui/input'

export type ClientTableSchema = {
  id: string
  status: 'invalid' | 'valid'
  secret: string
  stopped: boolean
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
    header: '运行信息/版本信息',
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
      <PopoverContent className="w-fit overflow-auto max-w-72 max-h-72 text-nowrap">
        <div>请点击命令框全选复制</div>
        <div>Linux安装到systemd</div>
        <Input readOnly value={platformInfo === undefined
          ? '获取平台信息失败'
          : LinuxInstallCommand('client', client, platformInfo)}></Input>
        <div>Windows安装到系统服务</div>
        <Input readOnly value={
          platformInfo === undefined
            ? "获取平台信息失败"
            : WindowsInstallCommand("client", client, platformInfo)
        }>
        </Input>
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
    let statusText: '在线' | '离线' | '错误' | '暂停' | '未知' = '未知'
    if (info === undefined) {
      return statusText
    }
    if (info.status === ClientStatus_Status.ONLINE) {
      statusText = '在线'
      if (client.stopped) {
        statusText = '暂停'
      }
    } else if (info.status === ClientStatus_Status.OFFLINE) {
      statusText = '离线'
    } else if (info.status === ClientStatus_Status.ERROR) {
      statusText = '错误'
    }
    return statusText
  }

  const infoColor =
    clientsInfo.data?.clients[client.id]?.status === ClientStatus_Status.ONLINE ? (
      client.stopped ? 'text-yellow-500' : 'text-green-500') : 'text-red-500'

  return (
    <div className="flex items-center gap-2 flex-row">
      <Badge variant={"secondary"} className={`p-2 border rounded font-mono w-fit ${infoColor} text-nowrap rounded-full h-6`}>
        {`${clientsInfo.data?.clients[client.id].ping}ms,${trans(clientsInfo.data?.clients[client.id])}`}
      </Badge>
      {clientsInfo.data?.clients[client.id].version &&
        <ClientDetail clientStatus={clientsInfo.data?.clients[client.id]} />
      }
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
        <div className="p-2 border rounded font-mono w-full break-all">
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

  const refetchList = () => {}

  const removeClient = useMutation({
    mutationFn: deleteClient,
    onSuccess: () => {
      toast({ description: '删除成功' })
      refetchList()
    },
    onError: () => {
      toast({ description: '删除失败' })
    },
  })

  const stopClient = useMutation({
    mutationFn: stopFrpc,
    onSuccess: () => {
      toast({ description: '停止成功' })
      refetchList()
    },
    onError: () => {
      toast({ description: '停止失败' })
    },
  })

  const startClient = useMutation({
    mutationFn: startFrpc,
    onSuccess: () => {
      toast({ description: '启动成功' })
      refetchList()
    },
    onError: () => {
      toast({ description: '启动失败' })
    },
  })

  const createAndDownloadFile = (fileName: string, content: string) => {
    var aTag = document.createElement('a');
    var blob = new Blob([content]);
    aTag.download = fileName;
    aTag.href = URL.createObjectURL(blob);
    aTag.click();
    URL.revokeObjectURL(aTag.href);
  }

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
          <DropdownMenuItem
            onClick={() => {
              try {
                if (platformInfo) {
                  createAndDownloadFile(`.env`, ClientEnvFile(client, platformInfo))
                }
              }
              catch (error) {
                toast({ description: '获取平台信息失败' })
              }
            }}
          >
            下载配置文件
          </DropdownMenuItem>
          {!client.stopped && <DropdownMenuItem className="text-destructive" onClick={() => stopClient.mutate({ clientId: client.id })}>暂停</DropdownMenuItem>}
          {client.stopped && <DropdownMenuItem onClick={() => startClient.mutate({ clientId: client.id })}>启动</DropdownMenuItem>}
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
              删除后运行中的客户端将无法通过现有参数再次连接，如果您需要删除客户端对外的连接，可以选择暂停客户端
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
    </Dialog >
  )
}
