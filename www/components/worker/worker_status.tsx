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

  // 统计部署状态
  const clientStatuses = statusResp?.workerStatus || {}
  const deployedClients = clients || []
  const totalClients = deployedClients.length

  const runningClients = Object.entries(clientStatuses).filter(([_, status]) => status === 'running').length
  const errorClients = Object.entries(clientStatuses).filter(([_, status]) => status === 'error').length
  const stoppedClients = Object.entries(clientStatuses).filter(([_, status]) => status === 'stopped').length

  // 统计入口状态
  const ingresses = ingressResp?.proxyConfigs || []
  const totalIngresses = ingresses.length

  // 查询所有入口的状态
  const ingressStatuses = useQuery({
    queryKey: ['getIngressStatuses', workerId, ingresses.map((i: ProxyConfig) => i.id).join(','), refetchTrigger],
    queryFn: async () => {
      const statuses: Record<string, string> = {}

      await Promise.all(
        ingresses.map(async (ingress: ProxyConfig) => {
          try {
            const proxyStatus = await getProxyConfig({
              clientId: ingress.clientId,
              serverId: ingress.serverId,
              name: ingress.name,
            })
            statuses[ingress.id || ''] = proxyStatus?.workingStatus?.status || 'unknown'
          } catch (e) {
            statuses[ingress.id || ''] = 'error'
          }
        }),
      )

      return statuses
    },
    enabled: ingresses.length > 0,
    refetchInterval: 10000,
  })

  const runningIngresses = Object.values(ingressStatuses.data || {}).filter((status) => status === 'running').length
  const errorIngresses = Object.values(ingressStatuses.data || {}).filter((status) =>
    ['error', 'start error', 'check failed'].includes(status),
  ).length

  // 计算总体状态
  const getOverallStatus = () => {
    if (totalClients === 0 && totalIngresses === 0) {
      return { variant: 'outline' as const, text: t('worker.status.no_resources'), color: 'bg-gray-100 text-gray-700' }
    }

    // 资源完全不可用
    if ((totalClients > 0 && runningClients === 0) || (totalIngresses > 0 && runningIngresses === 0)) {
      return { variant: 'destructive' as const, text: t('worker.status.unusable'), color: 'bg-red-500 text-white' }
    }

    // 资源不健康但部分可用
    if (errorClients > 0 || errorIngresses > 0) {
      return { variant: 'warning' as const, text: t('worker.status.unhealthy'), color: 'bg-amber-500 text-white' }
    }

    // 所有资源健康
    if (runningClients === totalClients && runningIngresses === totalIngresses) {
      return { variant: 'default' as const, text: t('worker.status.healthy'), color: 'bg-green-500 text-white' }
    }

    // 部分资源降级但无错误
    return { variant: 'secondary' as const, text: t('worker.status.degraded'), color: 'bg-orange-500 text-white' }
  }

  const { variant, text, color } = getOverallStatus()

  // 生成客户端资源状态指示器
  const renderClientIndicators = () => {
    if (totalClients === 0) return null

    return (
      <div className="flex items-center group relative">
        <div className="flex items-center space-x-1 rounded-md bg-muted/30 px-1 py-0.5 border border-muted">
          <div className="flex space-x-0.5">
            {Array.from({ length: Math.min(totalClients, 3) }).map((_, i) => (
              <div
                key={`client-${i}`}
                className={`h-2.5 w-2.5 rounded-sm ${
                  i < runningClients ? 'bg-green-500' : i < runningClients + errorClients ? 'bg-red-500' : 'bg-gray-300'
                }`}
              />
            ))}
          </div>
          <span className="text-xs text-muted-foreground">
            <Cpu className="w-3 h-3 min-w-3 min-h-3 max-w-3 max-h-3" />
          </span>
          {totalClients > 3 && <span className="text-xs text-muted-foreground">+{totalClients - 3}</span>}
        </div>
      </div>
    )
  }

  // 生成入口资源状态指示器
  const renderIngressIndicators = () => {
    if (totalIngresses === 0) return null

    return (
      <div className="flex items-center group relative">
        <div className="flex items-center space-x-1 rounded-md bg-muted/30 px-1 py-0.5 border border-muted">
          <div className="flex space-x-0.5">
            {Array.from({ length: Math.min(totalIngresses, 3) }).map((_, i) => (
              <div
                key={`ingress-${i}`}
                className={`h-2.5 w-2.5 rounded-sm ${
                  i < runningIngresses
                    ? 'bg-green-500'
                    : i < runningIngresses + errorIngresses
                      ? 'bg-red-500'
                      : 'bg-gray-300'
                }`}
              />
            ))}
          </div>
          <span className="text-xs text-muted-foreground">
            <Network className="w-3 h-3 min-w-3 min-h-3 max-w-3 max-h-3" />
          </span>
          {totalIngresses > 3 && <span className="text-xs text-muted-foreground">+{totalIngresses - 3}</span>}
        </div>
      </div>
    )
  }

  if (compact) {
    return (
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger className="flex items-center space-x-1">
            <div
              className={`h-2 w-2 rounded-sm ${
                variant === 'default'
                  ? 'bg-green-500'
                  : variant === 'destructive'
                    ? 'bg-red-500'
                    : variant === 'warning'
                      ? 'bg-amber-500'
                      : variant === 'secondary'
                        ? 'bg-blue-500'
                        : 'bg-gray-300'
              }`}
            />
          </TooltipTrigger>
          <TooltipContent>
            <div className="space-y-1">
              <p className="text-sm font-medium">{text}</p>
              <div className="text-xs">
                <div className="flex items-center space-x-1 font-mono">
                  {t('worker.status.clients')}: {runningClients}/{totalClients} {t('worker.status.running')}
                </div>
                <div className="flex items-center space-x-1 font-mono">
                  {t('worker.status.ingresses')}: {runningIngresses}/{totalIngresses} {t('worker.status.running')}
                </div>
              </div>
            </div>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    )
  }

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <div className="flex items-center space-x-2">
            <div className="flex space-x-1">
              {renderIngressIndicators()}
              {renderClientIndicators()}
            </div>
            <Badge className={`px-2 py-0.5 ${color} whitespace-nowrap`}>{text}</Badge>
          </div>
        </TooltipTrigger>
        <TooltipContent>
          <div className="space-y-2">
            <p className="text-sm font-medium">{text}</p>
            <div className="space-y-1 text-xs">
              <div className="flex items-center space-x-1 font-mono">
                {t('worker.status.clients')}: {runningClients}/{totalClients} {t('worker.status.running')}
              </div>
              <div className="flex items-center space-x-1 font-mono">
                {t('worker.status.ingresses')}: {runningIngresses}/{totalIngresses} {t('worker.status.running')}
              </div>
            </div>
          </div>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
}
