import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ProxyInfo } from "@/lib/pb/common"
import { formatBytes } from "@/lib/utils"
import { CloudDownload, CloudUpload } from "lucide-react"

export function ProxyTrafficOverview({ proxyInfo }: { proxyInfo: ProxyInfo }) {
  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="tracking-tight text-sm font-medium">今日入站流量</CardTitle>
          <CloudUpload className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatBytes(Number(proxyInfo.todayTrafficIn))}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="tracking-tight text-sm font-medium">今日出站流量</CardTitle>
          <CloudDownload className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatBytes(Number(proxyInfo.todayTrafficOut))}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="tracking-tight text-sm font-medium">历史入站流量</CardTitle>
          <CloudUpload className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatBytes(Number(proxyInfo.historyTrafficIn))}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="tracking-tight text-sm font-medium">历史出站流量</CardTitle>
          <CloudDownload className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatBytes(Number(proxyInfo.historyTrafficOut))}</div>
        </CardContent>
      </Card>
    </div>
  )
}