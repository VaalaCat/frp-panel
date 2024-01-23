import { useEffect, useState } from 'react'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from './ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select'
import { getServer, listServer } from '@/api/server'
import { useQuery } from '@tanstack/react-query'
import { Switch } from './ui/switch'
import { FRPSEditor } from './frps_editor'
import FRPSForm from './frps_form'
import { useSearchParams } from 'next/navigation'

export interface FRPSFormCardProps {
  serverID?: string
}
export const FRPSFormCard: React.FC<FRPSFormCardProps> = ({ serverID: defaultServerID }: FRPSFormCardProps = {}) => {
  const [advanceMode, setAdvanceMode] = useState<boolean>(false)
  const [serverID, setServerID] = useState<string | undefined>()
  const searchParams = useSearchParams()
  const paramServerID = searchParams.get('serverID')
  const { data: serverList, refetch: refetchServers } = useQuery({
    queryKey: ['listServer'],
    queryFn: () => {
      return listServer({ page: 1, pageSize: 100 })
    },
  })
  const { data: server, refetch: refetchServer } = useQuery({
    queryKey: ['getServer', serverID],
    queryFn: () => {
      return getServer({ serverId: serverID })
    },
  })

  useEffect(() => {
    if (defaultServerID) {
      setServerID(defaultServerID)
    }
  }, [defaultServerID])

  useEffect(() => {
    if (paramServerID) {
      setServerID(paramServerID)
    }
  }, [paramServerID])

  const handleServerChange = (value: string) => {
    setServerID(value)
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>服务端配置</CardTitle>
        <CardDescription>
          <div>
            注意⚠️：修改服务端配置文件后，服务端<a className='text-red-600'>会退出</a>
            <br />如果您使用的是docker容器且启动命令中包含了 --restart=unless-stopped 或 --restart=always 则无需担心。
            <br />如果您使用的是systemd安装也无需担心。
          </div>
          <div>
            选择服务端以管理Frps服务
          </div></CardDescription>
      </CardHeader>
      <CardContent>
        <div className=" flex items-center space-x-4 rounded-md border p-4">
          <div className="flex-1 space-y-1">
            <p className="text-sm font-medium leading-none">高级模式</p>
            <p className="text-sm text-muted-foreground">编辑服务端原始配置文件</p>
          </div>
          <Switch onCheckedChange={setAdvanceMode} />
        </div>
        <div className="flex flex-col w-full pt-2">
          <Label className="text-sm font-medium">服务端</Label>
          <Select
            onValueChange={handleServerChange}
            value={serverID}
            onOpenChange={() => {
              refetchServers()
              refetchServer()
            }}
          >
            <SelectTrigger className="my-2">
              <SelectValue placeholder="节点名称" />
            </SelectTrigger>
            <SelectContent>
              {serverList?.servers.map(
                (serverItem) =>
                  serverItem.id && (
                    <SelectItem key={serverItem.id} value={serverItem.id}>
                      {serverItem.id}
                    </SelectItem>
                  ),
              )}
            </SelectContent>
          </Select>
        </div>
        {serverID && server && server.server && !advanceMode && <FRPSForm serverID={serverID} server={server.server} />}
        {serverID && server && server.server && advanceMode && (
          <FRPSEditor serverID={serverID} server={server.server} />
        )}
      </CardContent>
      <CardFooter></CardFooter>
    </Card>
  )
}
