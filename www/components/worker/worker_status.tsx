'use client'

import React from 'react'
import { useTranslation } from 'react-i18next'
import { useQuery } from '@tanstack/react-query'
import { getWorkerStatus, getWorkerIngress } from '@/api/worker'
import { getProxyConfig } from '@/api/proxy'
import { useStore } from '@nanostores/react'
import { $proxyTableRefetchTrigger } from '@/store/refetch-trigger'
import { Client, ProxyConfig } from '@/lib/pb/common'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { Badge } from '@/components/ui/badge'
import { Cpu, Network } from 'lucide-react'

interface WorkerStatusProps {
  workerId: string
  clients?: Client[]
  compact?: boolean
}

export function WorkerStatus({ workerId, clients = [], compact = false }: WorkerStatusProps) {
  const { t } = useTranslation()
  const refetchTrigger = useStore($proxyTableRefetchTrigger)

  // 获取 Worker 状态
  const { data: statusResp } = useQuery({
    queryKey: ['workerStatus', workerId],
    queryFn: () => getWorkerStatus({ workerId }),
    enabled: !!workerId,
    refetchInterval: 10000,
  })

  // 获取 Worker Ingress
  const { data: ingressResp } = useQuery({
    queryKey: ['getWorkerIngress', workerId],
    queryFn: () => getWorkerIngress({ workerId }),
    enabled: !!workerId,
  })

  // 状态统计
  const clientStatuses = statusResp?.workerStatus || {}
  const deployedClients = clients
  const totalClients = deployedClients.length
  const runningClients = Object.values(clientStatuses).filter((s) => s === 'running').length
  const Clients = Object.values(clientStatuses).filter((s) => s === '').length
  const stoppedClients = Object.values(clientStatuses).filter((s) => s === 'stopped').length

  const ingresses = ingressResp?.proxyConfigs || []
  const totalIngresses = ingresses.length

  // 针对每个 ingress 再次拉取状态
  const ingressStatuses = useQuery({
    queryKey: ['getIngressStatuses', workerId, ingresses.map((i) => i.id).join(','), refetchTrigger],
    queryFn: async () => {
      const statuses: Record<string, string> = {}
      await Promise.all(
        ingresses.map(async (i) => {
          try {
            const ps = await getProxyConfig({
              clientId: i.clientId,
              serverId: i.serverId,
              name: i.name,
            })
            statuses[i.id || ''] = ps?.workingStatus?.status || 'unknown'
          } catch {
            statuses[i.id || ''] = ''
          }
        }),
      )
      return statuses
    },
    enabled: ingresses.length > 0,
    refetchInterval: 10000,
  })

  const runningIngresses = Object.values(ingressStatuses.data || {}).filter((s) => s === 'running').length
  const Ingresses = Object.values(ingressStatuses.data || {}).filter((s) =>
    ['', 'start', 'check failed'].includes(s),
  ).length

  // Overall 状态
  const getOverallStatus = () => {
    if (totalClients === 0 && totalIngresses === 0) {
      return { variant: 'outline' as const, text: t('worker.status.no_resources'), color: 'bg-gray-100 text-gray-700' }
    }
    if ((totalClients > 0 && runningClients === 0) || (totalIngresses > 0 && runningIngresses === 0)) {
      return { variant: 'destructive' as const, text: t('worker.status.unusable'), color: 'bg-red-500 text-white' }
    }
    if (Clients > 0 || Ingresses > 0) {
      return { variant: 'warning' as const, text: t('worker.status.unhealthy'), color: 'bg-amber-500 text-white' }
    }
    if (runningClients === totalClients && runningIngresses === totalIngresses) {
      return { variant: 'default' as const, text: t('worker.status.healthy'), color: 'bg-green-500 text-white' }
    }
    return { variant: 'secondary' as const, text: t('worker.status.degraded'), color: 'bg-orange-500 text-white' }
  }
  const { text: overallText, color: overallColor } = getOverallStatus()

  // per-client indicators
  const renderClientIndicators = () => {
    if (totalClients === 0) return null
    const showList = deployedClients.slice(0, 3)
    return (
      <div className="flex items-center space-x-1 rounded-md bg-muted/30 px-1 py-0.5 border border-muted">
        <div className="flex space-x-0.5">
          {showList.map((client) => {
            const status = clientStatuses[client.id || ''] || 'unknown'
            const bg = status === 'running' ? 'bg-green-500' : status === '' ? 'bg-red-500' : 'bg-gray-300'
            return (
              <Tooltip key={client.id} delayDuration={200}>
                <TooltipTrigger asChild>
                  <div className={`h-2.5 w-2.5 rounded-sm ${bg} cursor-pointer transition-transform hover:scale-125`} />
                </TooltipTrigger>
                <TooltipContent>
                  <p className="text-sm font-medium">{client.id}</p>
                  <div className="text-xs font-mono">
                    {t('worker.status.clients')}: {status}
                  </div>
                </TooltipContent>
              </Tooltip>
            )
          })}
        </div>
        <span className="text-xs text-muted-foreground">
          <Cpu className="w-3 h-3" />
        </span>
        {totalClients > 3 && <span className="text-xs text-muted-foreground">+{totalClients - 3}</span>}
      </div>
    )
  }

  // per-ingress indicators
  const renderIngressIndicators = () => {
    if (totalIngresses === 0) return null
    const showList = ingresses.slice(0, 3)
    return (
      <div className="flex items-center space-x-1 rounded-md bg-muted/30 px-1 py-0.5 border border-muted">
        <div className="flex space-x-0.5">
          {showList.map((ing) => {
            const status = ingressStatuses.data?.[ing.id || ''] || 'unknown'
            const bg =
              status === 'running'
                ? 'bg-green-500'
                : ['', 'start', 'check failed'].includes(status)
                  ? 'bg-red-500'
                  : 'bg-gray-300'
            return (
              <Tooltip key={ing.id} delayDuration={200}>
                <TooltipTrigger asChild>
                  <div className={`h-2.5 w-2.5 rounded-sm ${bg} cursor-pointer transition-transform hover:scale-125`} />
                </TooltipTrigger>
                <TooltipContent>
                  <p className="text-sm font-medium">{ing.name}</p>
                  <div className="text-xs font-mono">
                    {t('worker.status.ingresses')}: {status}
                  </div>
                </TooltipContent>
              </Tooltip>
            )
          })}
        </div>
        <span className="text-xs text-muted-foreground">
          <Network className="w-3 h-3" />
        </span>
        {totalIngresses > 3 && <span className="text-xs text-muted-foreground">+{totalIngresses - 3}</span>}
      </div>
    )
  }

  // Compact 模式仍旧整体 hover
  if (compact) {
    return (
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger className="flex items-center space-x-1">
            <div
              className={`h-2 w-2 rounded-sm ${
                overallColor.includes('green')
                  ? 'bg-green-500'
                  : overallColor.includes('red')
                    ? 'bg-red-500'
                    : overallColor.includes('amber')
                      ? 'bg-amber-500'
                      : overallColor.includes('orange')
                        ? 'bg-blue-500'
                        : 'bg-gray-300'
              }`}
            />
          </TooltipTrigger>
          <TooltipContent>
            <div className="space-y-1">
              <p className="text-sm font-medium">{overallText}</p>
              <div className="text-xs space-y-0.5">
                <div className="font-mono">
                  {t('worker.status.clients')}: {runningClients}/{totalClients}
                </div>
                <div className="font-mono">
                  {t('worker.status.ingresses')}: {runningIngresses}/{totalIngresses}
                </div>
              </div>
            </div>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    )
  }

  // 默认模式：拆分整体与细节 hover
  return (
    <TooltipProvider>
      <div className="flex items-center space-x-2">
        {renderIngressIndicators()}
        {renderClientIndicators()}

        {/* 只在 Badge 上展示总体状态 */}
        <Tooltip>
          <TooltipTrigger asChild>
            <Badge className={`px-2 py-0.5 ${overallColor} whitespace-nowrap`}>{overallText}</Badge>
          </TooltipTrigger>
          <TooltipContent>
            <div className="space-y-2">
              <p className="text-sm font-medium">{overallText}</p>
              <div className="grid grid-cols-2 gap-4 text-xs font-mono">
                <div className="flex items-center space-x-1">
                  <span className="w-2 h-2 rounded-full bg-green-500" />
                  <span>
                    {t('worker.status.running')}: {runningClients}/{totalClients}
                  </span>
                </div>
                <div className="flex items-center space-x-1">
                  <span className="w-2 h-2 rounded-full bg-red-500" />
                  <span>
                    {t('worker.status_text')}: {Clients}/{totalClients}
                  </span>
                </div>
                <div className="flex items-center space-x-1">
                  <span className="w-2 h-2 rounded-full bg-green-500" />
                  <span>
                    {t('worker.status.running')}: {runningIngresses}/{totalIngresses}
                  </span>
                </div>
                <div className="flex items-center space-x-1">
                  <span className="w-2 h-2 rounded-full bg-red-500" />
                  <span>
                    {t('worker.status_text')}: {Ingresses}/{totalIngresses}
                  </span>
                </div>
              </div>
            </div>
          </TooltipContent>
        </Tooltip>
      </div>
    </TooltipProvider>
  )
}
