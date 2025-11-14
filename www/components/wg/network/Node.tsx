'use client'

import React, { useCallback, useMemo } from 'react'
import type { NodeProps } from '@xyflow/react'
import { Handle, Position, useConnection } from '@xyflow/react'
import type { WGNode } from './types'
import { cn, formatBytes } from '@/lib/utils'
import { useTranslation } from 'react-i18next'
import { useQuery } from '@tanstack/react-query'
import { getClientsStatus } from '@/api/platform'
import { getWireGuardRuntime } from '@/api/wg'
import { ClientType } from '@/lib/pb/common'
import { GetWireGuardRuntimeInfoRequest } from '@/lib/pb/api_wg'
import { ClientStatus_Status } from '@/lib/pb/api_master'
import { Badge } from '@/components/ui/badge'
import {
  Activity,
  Wifi,
  WifiOff,
  Terminal,
  Network,
  Clock,
  Signal,
  TrendingUp,
  TrendingDown,
  Server,
  FileText
} from 'lucide-react'
import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuTrigger } from '@/components/ui/context-menu'

interface NodeComponentProps extends NodeProps<WGNode> {
  onOpenTerminal?: (clientId: string, clientType: number) => void
  onOpenLog?: (clientId: string, clientType: number) => void
}

const WGNodeComponent: React.FC<NodeComponentProps> = ({
  id,
  data,
  selected,
  onOpenTerminal,
  onOpenLog,
}) => {
  const connection = useConnection()
  const isTarget = !!connection.inProgress && connection.fromNode?.id !== id
  const { t } = useTranslation()

  const clientId = data.config?.clientId

  // 获取客户端状态
  const { data: clientStatusData } = useQuery({
    queryKey: ['clientStatus', clientId],
    queryFn: async () => {
      if (!clientId) return undefined
      return await getClientsStatus({
        clientIds: [clientId],
        clientType: ClientType.FRPC,
      })
    },
    enabled: !!clientId,
    refetchInterval: 15000,
  })

  // 获取运行时信息
  const { data: runtimeData } = useQuery({
    queryKey: ['wgRuntime', id],
    queryFn: async () => {
      return await getWireGuardRuntime(GetWireGuardRuntimeInfoRequest.create({ id: Number(id) }))
    },
    enabled: !!id && !isNaN(Number(id)),
    refetchInterval: 15000,
  })

  // 计算统计数据
  const stats = useMemo(() => {
    const runtime = runtimeData?.wgDeviceRuntimeInfo
    const clientStatus = clientStatusData?.clients[clientId || '']
    const peers = runtime?.peers || []

    const isOnline = clientStatus?.status === ClientStatus_Status.ONLINE
    const ping = clientStatus?.ping
    const peerCount = peers.length

    const totalTx = peers.reduce((sum, p) => sum + (Number(p.txBytes) || 0), 0)
    const totalRx = peers.reduce((sum, p) => sum + (Number(p.rxBytes) || 0), 0)

    const lastHandshake = peers
      .map((p) => (p.lastHandshakeTimeSec ? Number(p.lastHandshakeTimeSec) : 0))
      .filter((t) => t > 0)
      .sort((a, b) => b - a)[0]

    return {
      isOnline,
      ping,
      peerCount,
      totalTx,
      totalRx,
      lastHandshake,
      runtime,
    }
  }, [runtimeData, clientStatusData, clientId])

  const handleOpenTerminal = React.useCallback(
    (event: React.MouseEvent) => {
      event.stopPropagation()
      event.preventDefault()
      if (clientId) {
        onOpenTerminal?.(clientId, ClientType.FRPC)
      }
    },
    [clientId, onOpenTerminal]
  )

  const handleOpenLog = React.useCallback(
    (event: React.MouseEvent) => {
      event.stopPropagation()
      event.preventDefault()
      if (clientId) {
        onOpenLog?.(clientId, ClientType.FRPC)
      }
    },
    [clientId, onOpenLog]
  )

  // 获取连接状态徽章 - 显示延迟
  const getStatusBadge = () => {
    if (!stats.isOnline) {
      return <Badge variant="secondary" className="text-[10px] px-1.5 py-0">离线</Badge>
    }
    if (stats.ping !== undefined) {
      if (stats.ping <= 50) {
        return <Badge className="text-[10px] px-1.5 py-0 bg-emerald-500">{stats.ping}ms</Badge>
      }
      if (stats.ping <= 200) {
        return <Badge variant="secondary" className="text-[10px] px-1.5 py-0">{stats.ping}ms</Badge>
      }
      return <Badge variant="destructive" className="text-[10px] px-1.5 py-0">{stats.ping}ms</Badge>
    }
    return <Badge className="text-[10px] px-1.5 py-0 bg-emerald-500">在线</Badge>
  }

  const endpoint = useCallback(() => {
    if (!data.config?.advertisedEndpoints?.[0]) {
      return '-'
    }

    if (data.config?.advertisedEndpoints?.[0]?.uri) {
      return data.config?.advertisedEndpoints?.[0]?.uri
    }

    if (data.config?.advertisedEndpoints?.[0]?.host && data.config?.advertisedEndpoints?.[0]?.port) {
      return data.config?.advertisedEndpoints?.[0]?.host + ':' + data.config?.advertisedEndpoints?.[0]?.port
    }

    return '-'
  }, [data.config?.advertisedEndpoints])

  return (
    <ContextMenu>
      <ContextMenuTrigger asChild>
        <div
          className={cn(
            'group relative',
            selected && 'z-10'
          )}
          style={{ width: 240, height: 150, userSelect: 'none' }}
        >
          {/* 选中时的外层光晕 */}
          {selected && (
            <div className="absolute inset-0 rounded-lg bg-primary/10 blur-lg" />
          )}

          {/* 主卡片 */}
          <div
            className={cn(
              'relative h-full rounded-lg border-2 bg-gradient-to-br from-card to-card/95',
              'shadow-md transition-all duration-200',
              'hover:shadow-xl hover:scale-[1.01]',
              selected
                ? 'border-primary shadow-primary/20'
                : 'border-border/50 hover:border-primary/30'
            )}
          >
            {/* 顶部状态栏 - 更紧凑 */}
            <div className="flex items-center justify-between gap-1.5 border-b bg-muted/30 px-2 py-1.5 rounded-t-lg">
              <div className="flex items-center gap-1.5 flex-1 min-w-0">
                {stats.isOnline ? (
                  <Wifi className="h-3.5 w-3.5 text-emerald-500 shrink-0" />
                ) : (
                  <WifiOff className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
                )}
                <span className="font-semibold text-xs truncate" title={clientId}>
                  {clientId || `Node ${id}`}
                </span>
              </div>
              {getStatusBadge()}
              <div className="drag-handle cursor-grab active:cursor-grabbing p-0.5 hover:bg-muted rounded transition-colors shrink-0">
                <Server className="h-3 w-3 text-muted-foreground" />
              </div>
            </div>

            {/* 主要信息区 - 紧凑布局 */}
            <div className="p-2 space-y-1.5">
              {/* ID (可点击) 和 Endpoint */}
              <div className="flex items-center justify-between text-[10px]">
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    // 跳转到 WireGuard 详情页
                    window.location.href = `/wg/wireguards?id=${id}`
                  }}
                  className="flex items-center gap-1 hover:text-primary transition-colors cursor-pointer"
                >
                  <span className="text-muted-foreground">ID:</span>
                  <span className="font-mono font-semibold underline decoration-dotted">#{id}</span>
                </button>
                <div className="flex items-center gap-1 min-w-0">
                  <span className="font-mono text-[9px] truncate text-muted-foreground"
                    title={endpoint()}>
                    {endpoint()}
                  </span>
                </div>
              </div>

              {/* 虚拟IP和端口 - 紧凑显示 */}
              {stats.runtime && (
                <div className="flex items-center justify-between bg-blue-50 dark:bg-blue-950/20 rounded px-2 py-1">
                  <div className="flex items-center gap-1.5 flex-1 min-w-0">
                    <Network className="h-3 w-3 text-blue-600 dark:text-blue-400 shrink-0" />
                    <span className="font-mono text-[10px] font-medium truncate" title={stats.runtime.virtualIp}>
                      {stats.runtime.virtualIp || '-'}
                    </span>
                    {stats.runtime.listenPort && (
                      <span className="text-[9px] text-muted-foreground">:{stats.runtime.listenPort}</span>
                    )}
                  </div>
                </div>
              )}

              {/* Peers 和流量 - 横向紧凑布局 */}
              <div className="flex items-center justify-between gap-1 text-[9px] bg-muted/30 rounded px-2 py-1">
                <div className="flex items-center gap-1">
                  <Activity className="h-3 w-3 text-emerald-500 shrink-0" />
                  <span className="font-semibold">{stats.peerCount}</span>
                  <span className="text-[8px] text-muted-foreground">P</span>
                </div>
                <div className="flex items-center gap-1">
                  <TrendingUp className="h-3 w-3 text-blue-500 shrink-0" />
                  <span className="font-mono font-semibold text-[9px]">{formatBytes(stats.totalTx)}</span>
                </div>
                <div className="flex items-center gap-1">
                  <TrendingDown className="h-3 w-3 text-green-500 shrink-0" />
                  <span className="font-mono font-semibold text-[9px]">{formatBytes(stats.totalRx)}</span>
                </div>
              </div>


              <div className="flex items-center justify-between gap-1 text-[9px] bg-muted/30 rounded px-2 py-1">
                {/* 底部信息 - 最近握手 */}
                {stats.lastHandshake && (
                  <div className="flex items-center gap-0.5 text-muted-foreground text-[9px]">
                    <Clock className="h-2.5 w-2.5 shrink-0" />
                    <span className="truncate">{new Date(stats.lastHandshake * 1000).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })}</span>
                  </div>
                )}

                {/* 标签 - 完整展示 */}
                {data.config?.tags && data.config.tags.length > 0 && (
                  <div className="flex flex-wrap gap-0.5">
                    {data.config.tags.map((tag) => (
                      <Badge key={tag} variant="outline" className="text-[8px] h-3.5 px-1 py-0 leading-none">
                        #{tag}
                      </Badge>
                    ))}
                  </div>
                )}
              </div>
            </div>

            {/* Handles */}
            {!connection.inProgress && (
              <Handle
                className="!w-2.5 !h-2.5 !bg-primary !border-2 !border-background"
                position={Position.Right}
                type="source"
              />
            )}
            {(!connection.inProgress || isTarget) && (
              <Handle
                className="!w-2.5 !h-2.5 !bg-primary !border-2 !border-background"
                position={Position.Left}
                type="target"
                isConnectableStart={false}
              />
            )}
          </div>
        </div>
      </ContextMenuTrigger>

      <ContextMenuContent className="w-48">
        <ContextMenuItem onClick={handleOpenTerminal} disabled={!clientId} className="gap-2">
          <Terminal className="h-4 w-4" />
          <span>{t('wg.contextMenu.openTerminal')}</span>
        </ContextMenuItem>
        <ContextMenuItem onClick={handleOpenLog} disabled={!clientId} className="gap-2">
          <FileText className="h-4 w-4" />
          <span>{t('wg.contextMenu.viewLogs')}</span>
        </ContextMenuItem>
      </ContextMenuContent>
    </ContextMenu>
  )
}

export default WGNodeComponent
