import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ProxyInfo } from "@/lib/pb/common"
import { formatBytes } from "@/lib/utils"
import { CloudDownload, CloudUpload } from "lucide-react"
import { useTranslation } from "react-i18next"

export function ProxyTrafficOverview({ proxyInfo }: { proxyInfo: ProxyInfo }) {
  const { t } = useTranslation()
  
  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="tracking-tight text-sm font-medium">{t('traffic.today.inbound')}</CardTitle>
          <CloudUpload className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatBytes(Number(proxyInfo.todayTrafficIn))}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="tracking-tight text-sm font-medium">{t('traffic.today.outbound')}</CardTitle>
          <CloudDownload className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatBytes(Number(proxyInfo.todayTrafficOut))}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="tracking-tight text-sm font-medium">{t('traffic.history.inbound')}</CardTitle>
          <CloudUpload className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatBytes(Number(proxyInfo.historyTrafficIn))}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="tracking-tight text-sm font-medium">{t('traffic.history.outbound')}</CardTitle>
          <CloudDownload className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatBytes(Number(proxyInfo.historyTrafficOut))}</div>
        </CardContent>
      </Card>
    </div>
  )
}