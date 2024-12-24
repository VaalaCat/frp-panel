import { useEffect, useState } from 'react'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { useQuery } from '@tanstack/react-query'
import { useSearchParams } from 'next/navigation'
import { getProxyStatsByClientID } from '@/api/stats'
import { ProxyTrafficPieChart } from '../charts/proxy-traffic-pie-chart'
import { ProxyTrafficOverview } from '../charts/proxy-traffic-overview'
import { ClientSelector } from '../base/client-selector'
import { ProxySelector } from '../base/proxy-selector'
import { ProxyInfo } from '@/lib/pb/common'
import { Button } from '../ui/button'
import { CheckCircle2, CircleX, RefreshCcw } from "lucide-react"
import { useTranslation } from 'react-i18next';
import TrafficStatsCard from './stats-item'
import { Separator } from '@radix-ui/react-separator'
import { formatBytes } from '@/lib/utils'

export interface ClientStatsCardProps {
  clientID?: string
}
export const ClientStatsCard: React.FC<ClientStatsCardProps> = ({ clientID: defaultClientID }: ClientStatsCardProps = {}) => {
  const { t } = useTranslation();
  const [clientID, setClientID] = useState<string | undefined>()
  const [status, setStatus] = useState<"loading" | "success" | "error" | undefined>()
  const [timeoutId, setTimeoutId] = useState<NodeJS.Timeout | null>(null);

  const searchParams = useSearchParams()
  const paramClientID = searchParams.get('clientID')

  const { data: clientStatsList, refetch: refetchClientStats } = useQuery({
    queryKey: ['clientStats', clientID],
    queryFn: () => {
      return getProxyStatsByClientID({ clientId: clientID! })
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

  const mergeProxyInfos = (proxyInfos: ProxyInfo[]): ProxyInfo[] => {
    const mergedMap: Map<string, ProxyInfo> = new Map();

    for (const proxyInfo of proxyInfos) {
      const key = `${proxyInfo.clientId}:${proxyInfo.name}`;

      if (!mergedMap.has(key)) {
        mergedMap.set(key, { ...proxyInfo });
      } else {
        const existingProxyInfo = mergedMap.get(key)!;
        existingProxyInfo.todayTrafficIn = (existingProxyInfo.todayTrafficIn || BigInt(0)) + (proxyInfo.todayTrafficIn || BigInt(0));
        existingProxyInfo.todayTrafficOut = (existingProxyInfo.todayTrafficOut || BigInt(0)) + (proxyInfo.todayTrafficOut || BigInt(0));
        existingProxyInfo.historyTrafficIn = (existingProxyInfo.historyTrafficIn || BigInt(0)) + (proxyInfo.historyTrafficIn || BigInt(0));
        existingProxyInfo.historyTrafficOut = (existingProxyInfo.historyTrafficOut || BigInt(0)) + (proxyInfo.historyTrafficOut || BigInt(0));
      }
    }

    return Array.from(mergedMap.values());
  };

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>{t('client.stats.title')}</CardTitle>
        <CardDescription>
          <div>
            {t('client.stats.description')}
          </div>
        </CardDescription>
      </CardHeader>
      <CardContent className='space-y-4 flex flex-col flex-1'>
        <Label>{t('client.stats.label')}</Label>
        <ClientSelector clientID={clientID} setClientID={handleClientChange} onOpenChange={() => {
          refetchClientStats()
        }} />
        <Label>{t('proxy.stats.label')}</Label>
        <div className="w-full grid gap-2 grid-cols-1 overflow-x-auto">
          {clientStatsList && clientStatsList.proxyInfos.length > 0 &&
            clientStatsList.proxyInfos.map((proxyInfo) => {
              return (
                <ProxyStatusCard key={proxyInfo.name} proxyInfo={proxyInfo} />
              )
            })
          }
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
            if (timeoutId) { clearTimeout(timeoutId); }

            const id = setTimeout(() => {
              setStatus(undefined)
            }, 3000)
            setTimeoutId(id)
          })
        }}>
          {status === "loading" && <RefreshCcw className="w-4 h-4 animate-spin" />}
          {status === "success" && <CheckCircle2 className="w-4 h-4" />}
          {status === "error" && <CircleX className="w-4 h-4" />}
          <p>{t('refresh.data')}</p></Button>
      </CardFooter>
    </Card>
  )
}

const ProxyStatusCard: React.FC<{ proxyInfo: ProxyInfo | undefined }> = ({ proxyInfo }) => {
  const { t } = useTranslation();

  if (!proxyInfo) {
    return null;
  }

  return (
    <Card key={proxyInfo.name} className="flex flex-row gap-2 p-4 w-full min-w-[900px] shadow-none justify-between">
      <Label>{t('proxy.stats.tunnel_traffic', { name: proxyInfo.name })}</Label>
      <Separator orientation="vertical" />
      <ProxyTrafficField label={t('traffic.today.total')} value={formatBytes(Number(proxyInfo.todayTrafficOut || BigInt(0)))} />
      <ProxyTrafficField label={t('traffic.today.inbound')} value={formatBytes(Number(proxyInfo.todayTrafficIn || BigInt(0)))} />
      <ProxyTrafficField label={t('traffic.today.outbound')} value={formatBytes(Number(proxyInfo.todayTrafficOut || BigInt(0)))} />
      <Separator orientation="vertical" />
      <ProxyTrafficField label={t('traffic.history.total')} value={formatBytes(Number(proxyInfo.historyTrafficOut || BigInt(0)))} />
      <ProxyTrafficField label={t('traffic.history.inbound')} value={formatBytes(Number(proxyInfo.historyTrafficIn || BigInt(0)))} />
      <ProxyTrafficField label={t('traffic.history.outbound')} value={formatBytes(Number(proxyInfo.historyTrafficOut || BigInt(0)))} />
    </Card>
  );
}

const ProxyTrafficField = ({ label, value }: { label: string, value: string }) => {
  return <div className="flex w-fit flex-col">
    <p className="text-xs text-muted-foreground text-nowrap">{label}</p>
    <div className="flex items-center text-xs font-semibold text-nowrap">{value}</div>
  </div>
}