import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ProxyInfo } from "@/lib/pb/common"
import { formatBytes } from "@/lib/utils"

export function ProxyTrafficOverview({ proxyInfo }: { proxyInfo: ProxyInfo }) {
  const todayTotal = Number(proxyInfo.todayTrafficIn) + Number(proxyInfo.todayTrafficOut)
  const historyTotal = Number(proxyInfo.historyTrafficIn) + Number(proxyInfo.historyTrafficOut)

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle>今日总流量</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatBytes(todayTotal)}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle>历史总流量</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatBytes(historyTotal)}</div>
        </CardContent>
      </Card>
    </div>
  )
}