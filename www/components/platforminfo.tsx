import { useQuery } from '@tanstack/react-query'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { TbDeviceHeartMonitor, TbEngine, TbEngineOff, TbServer2, TbServerBolt, TbServerOff } from 'react-icons/tb'
import { useEffect } from 'react'
import { $platformInfo } from '@/store/user'
import { getPlatformInfo } from '@/api/platform'
export const PlatformInfo = () => {
  const platformInfo = useQuery({
    queryKey: ['platformInfo'],
    queryFn: getPlatformInfo,
  })
  useEffect(() => {
    $platformInfo.set(platformInfo.data)
  }, [platformInfo])
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
      <Card>
        <CardHeader>
          <div className="flex justify-between">
            <h3 className="tracking-tight text-sm font-medium">已配置服务端数</h3>
            <TbServerBolt className="mt-1" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{platformInfo.data?.configuredServerCount} 个</div>
          <p className="text-xs text-muted-foreground">请前往左侧🫲菜单修改</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <div className="flex justify-between">
            <h3 className="tracking-tight text-sm font-medium">已配置客户端数</h3>
            <TbEngine className="mt-1" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{platformInfo.data?.configuredClientCount} 个</div>
          <p className="text-xs text-muted-foreground">请前往左侧🫲菜单修改</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <div className="flex justify-between">
            <h3 className="tracking-tight text-sm font-medium">未配置服务端数</h3>
            <TbServerOff className="mt-1" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{platformInfo.data?.unconfiguredServerCount} 个</div>
          <p className="text-xs text-muted-foreground">请前往左侧🫲菜单修改</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <div className="flex justify-between">
            <h3 className="tracking-tight text-sm font-medium">未配置客户端数</h3>
            <TbEngineOff className="mt-1" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{platformInfo.data?.unconfiguredClientCount} 个</div>
          <p className="text-xs text-muted-foreground">请前往左侧🫲菜单修改</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <div className="flex justify-between">
            <h3 className="tracking-tight text-sm font-medium">服务端总数</h3>
            <TbServer2 className="mt-1" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{platformInfo.data?.totalServerCount} 个</div>
          <p className="text-xs text-muted-foreground">请前往左侧🫲菜单修改</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <div className="flex justify-between">
            <h3 className="tracking-tight text-sm font-medium">客户端总数</h3>
            <TbDeviceHeartMonitor className="mt-1" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{platformInfo.data?.totalClientCount} 个</div>
          <p className="text-xs text-muted-foreground">请前往左侧🫲菜单修改</p>
        </CardContent>
      </Card>
    </div>
  )
}
