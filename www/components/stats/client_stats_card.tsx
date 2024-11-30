import { useEffect, useState } from 'react'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { useQuery } from '@tanstack/react-query'
import { useSearchParams } from 'next/navigation'
import { getProxyStatsByClientID } from '@/api/stats'
import { ProxyTrafficBarChart } from '../charts/proxy-traffic-bar-chart'
import { ProxyTrafficOverview } from '../charts/proxy-traffic-overview'
import { ClientSelector } from '../base/client-selector'
import { ProxySelector } from '../base/proxy-selector'
import { ProxyInfo } from '@/lib/pb/common'
import { Button } from '../ui/button'
import { CheckCircle2, CircleX, RefreshCcw } from "lucide-react"

export interface ClientStatsCardProps {
  clientID?: string
}
export const ClientStatsCard: React.FC<ClientStatsCardProps> = ({ clientID: defaultClientID }: ClientStatsCardProps = {}) => {
  const [clientID, setClientID] = useState<string | undefined>()
  const [proxyName, setProxyName] = useState<string | undefined>()
  const [status, setStatus] = useState<"loading" | "success" | "error" | undefined>()

  const searchParams = useSearchParams()
  const paramClientID = searchParams.get('clientID')

  const { data: clientStatsList, refetch: refetchClientStats } = useQuery({
    queryKey: ['clientStats', clientID],
    queryFn: async () => {
      return await getProxyStatsByClientID({ clientId: clientID! })
    },
  })

  useEffect(() => {
    if (defaultClientID) {
      setClientID(defaultClientID)
    }
  }, [defaultClientID])

  useEffect(() => {
    if (paramClientID) {
      setClientID(paramClientID)
    }
  }, [paramClientID])

  const handleClientChange = (value: string) => {
    setClientID(value)
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>客户端统计</CardTitle>
        <CardDescription>
          <div>
            按照客户端名称对所有隧道的流量使用统计
          </div>
        </CardDescription>
      </CardHeader>
      <CardContent className='space-y-4 flex flex-col flex-1'>
        <Label>客户端</Label>
        <ClientSelector clientID={clientID} setClientID={handleClientChange} onOpenChange={() => {
          refetchClientStats()
          setProxyName(undefined)
        }} />
        <Label>隧道名称</Label>
        <ProxySelector
          // @ts-ignore
          proxyNames={clientStatsList?.proxyInfos.map((proxyInfo) => proxyInfo.name).filter((value) => value !== undefined) || []}
          proxyName={proxyName}
          setProxyname={setProxyName} />
        <div className="w-full grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
          {clientStatsList && clientStatsList.proxyInfos.length > 0 &&
            ProxyStatusCard(clientStatsList.proxyInfos.find((proxyInfo) => proxyInfo.name === proxyName))}
        </div>
      </CardContent>
      <CardFooter>
        <Button className="space-x-2" onClick={() => {
          setStatus("loading")
          refetchClientStats().then(() => {
            setStatus("success")
          }).catch(() => {
            setStatus("error")
          }).finally(() => {
            const timer = setTimeout(() => {
              setStatus(undefined)
            }, 3000)
            return () => clearTimeout(timer)
          })
        }}>
          {status === "loading" && <RefreshCcw className="w-4 h-4 animate-spin" />}
          {status === "success" && <CheckCircle2 className="w-4 h-4" />}
          {status === "error" && <CircleX className="w-4 h-4" />}
          <p>刷新数据</p></Button>
      </CardFooter>
    </Card>
  )
}

const ProxyStatusCard = (proxyInfo: ProxyInfo | undefined) => {
  return (<>{proxyInfo &&
    <div key={proxyInfo.name} className="flex flex-col w-full space-y-4">
      <Label>{`隧道 ${proxyInfo.name} 流量使用`}</Label>
      <ProxyTrafficOverview proxyInfo={proxyInfo} />
      <ProxyTrafficBarChart proxyInfo={proxyInfo} />
    </div>
  }</>)
}