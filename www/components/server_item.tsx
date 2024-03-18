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
import { ClientEnvFile, ExecCommandStr, LinuxInstallCommand, WindowsInstallCommand } from '@/lib/consts'
import { useMutation, useQuery } from '@tanstack/react-query'
import { deleteServer, listServer } from '@/api/server'
import { useRouter } from 'next/router'
import { getUserInfo } from '@/api/user'
import { useStore } from '@nanostores/react'
import { $platformInfo } from '@/store/user'
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover'
import { getClientsStatus } from '@/api/platform'
import { ClientType } from '@/lib/pb/common'
import { ClientStatus, ClientStatus_Status } from '@/lib/pb/api_master'

export type ServerTableSchema = {
  id: string
  status: 'invalid' | 'valid'
  secret: string
  info?: string
  ip: string
  config?: string
}

export const columns: ColumnDef<ServerTableSchema>[] = [
  {
    accessorKey: 'id',
    header: 'ID(点击查看安装命令)',
    cell: ({ row }) => {
      return <ServerID server={row.original} />
    },
  },
  {
    accessorKey: 'status',
    header: '是否配置',
    cell: ({ row }) => {
      const Server = row.original
      return (
        <div className={`font-mono ${Server.status === 'valid' ? 'text-green-500' : 'text-red-500'} min-w-12`}>
          {
            {
              valid: '已配置',
              invalid: '未配置',
            }[Server.status]
          }
        </div>
      )
    },
  },
  {
    accessorKey: 'info',
    header: '运行信息',
    cell: ({ row }) => {
      const Server = row.original
      return <ServerInfo server={Server} />
    },
  },
  {
    accessorKey: 'ip',
    header: 'IP',
    cell: ({ row }) => {
      const Server = row.original
      return <div className="font-mono">{Server.ip}</div>
    },
  },
  {
    accessorKey: 'secret',
    header: '连接密钥(点击查看启动命令)',
    cell: ({ row }) => {
      const Server = row.original
      return <ServerSecret server={Server} />
    },
  },
  {
    id: 'action',
    cell: ({ row, table }) => {
      const Server = row.original
      return <ServerActions server={Server} table={table} />
    },
  },
]

export const ServerID = ({ server }: { server: ServerTableSchema }) => {
  const platformInfo = useStore($platformInfo)
  return (
    <Popover>
      <PopoverTrigger asChild>
        <div className="font-mono">{server.id}</div>
      </PopoverTrigger>
      <PopoverContent className="w-fit overflow-auto max-w-72 max-h-72">
        <div>Linux安装到systemd</div>
        <div className="p-2 border rounded font-mono w-fit">
          {platformInfo === undefined ? '获取平台信息失败' : LinuxInstallCommand('server', server, platformInfo)}
        </div>
        <div>Windows安装到系统服务</div>
        <div className="p-2 border rounded font-mono w-fit">
          {platformInfo === undefined ? '获取平台信息失败' : WindowsInstallCommand('server', server, platformInfo)}
        </div>
      </PopoverContent>
    </Popover>
  )
}

export const ServerInfo = ({ server }: { server: ServerTableSchema }) => {
  const clientsInfo = useQuery({
    queryKey: ['getClientsStatus', [server.id]],
    queryFn: async () => {
      return await getClientsStatus({
        clientIds: [server.id],
        clientType: ClientType.FRPS,
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
    clientsInfo.data?.clients[server.id]?.status === ClientStatus_Status.ONLINE ? 'text-green-500' : 'text-red-500'

  return (
    <div className={`p-2 border rounded font-mono w-fit ${infoColor}`}>
      {`${clientsInfo.data?.clients[server.id].ping}ms, ${trans(clientsInfo.data?.clients[server.id])}`}
    </div>
  )
}

export const ServerSecret = ({ server }: { server: ServerTableSchema }) => {
  const [showSecrect, setShowSecrect] = useState<boolean>(false)
  const fakeSecret = Array.from({ length: server.secret.length })
    .map(() => '*')
    .join('')
  const { toast } = useToast()
  const platformInfo = useStore($platformInfo)

  return (
    <Popover>
      <PopoverTrigger asChild>
        <div
          onMouseEnter={() => setShowSecrect(true)}
          onMouseLeave={() => setShowSecrect(false)}
          className="font-medium hover:rounded hover:bg-slate-100 p-2 font-mono whitespace-nowrap"
        >
          {showSecrect ? server.secret : fakeSecret}
        </div>
      </PopoverTrigger>
      <PopoverContent className="w-fit overflow-auto max-w-48">
        <div>运行命令(需要<a className='text-blue-500' href='https://github.com/VaalaCat/frp-panel/releases'>点击这里</a>自行下载文件)</div>
        <div className="p-2 border rounded font-mono w-fit">
          {platformInfo === undefined ? '获取平台信息失败' : ExecCommandStr('server', server, platformInfo)}
        </div>
      </PopoverContent>
    </Popover>
  )
}

export interface ServerItemProps {
  server: ServerTableSchema
  table: Table<ServerTableSchema>
}

export const ServerActions: React.FC<ServerItemProps> = ({ server, table }) => {
  const { toast } = useToast()
  const router = useRouter()
  const platformInfo = useStore($platformInfo)

  const fetchDataOptions = {
    pageIndex: table.getState().pagination.pageIndex,
    pageSize: table.getState().pagination.pageSize,
  }

  const dataQuery = useQuery({
    queryKey: ['listServer', fetchDataOptions],
    queryFn: async () => {
      return await listServer({
        page: fetchDataOptions.pageIndex + 1,
        pageSize: fetchDataOptions.pageSize,
      })
    },
  })
  const removeServer = useMutation({
    mutationFn: deleteServer,
    onSuccess: () => {
      toast({ description: '删除成功' })
      dataQuery.refetch()
    },
    onError: () => {
      toast({ description: '删除失败' })
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
                  navigator.clipboard.writeText(ExecCommandStr('server', server, platformInfo))
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
              router.push({
                pathname: '/serveredit',
                query: {
                  serverID: server.id,
                },
              })
            }}
          >
            修改服务端配置
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={() => {
              try {
                if (platformInfo) {
                  createAndDownloadFile(`.env`, ClientEnvFile(server, platformInfo))
                }
              }
              catch (error) {
                toast({ description: '获取平台信息失败' })
              }
            }}
          >
            下载配置文件
          </DropdownMenuItem>
          <DialogTrigger asChild>
            <DropdownMenuItem className="text-destructive">删除</DropdownMenuItem>
          </DialogTrigger>
        </DropdownMenuContent>
      </DropdownMenu>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>确定删除该服务端?</DialogTitle>
          <DialogDescription>
            <p className="text-destructive">此操作无法撤消。您确定要永久从我们的服务器中删除该客户端?</p>
            <p className="text-gray-500 border-l-4 border-gray-500 pl-4 py-2">
              删除后运行中的服务端将无法通过现有参数再次连接，如果您需要停止服务端的服务，可以选择清空配置
            </p>
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose asChild>
            <Button type="submit" onClick={() => removeServer.mutate({ serverId: server.id })}>
              确定
            </Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
