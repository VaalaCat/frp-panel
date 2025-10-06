'use client'

import React from 'react'
import type { NodeProps } from '@xyflow/react'
import { Handle, Position, useConnection } from '@xyflow/react'
import type { WGNode } from './types'
import { cn } from '@/lib/utils'
import { useTranslation } from 'react-i18next'
import { useQuery } from '@tanstack/react-query'
import { getClientsStatus } from '@/api/platform'
import { getWireGuardRuntime } from '@/api/wg'
import { ClientType } from '@/lib/pb/common'
import { GetWireGuardRuntimeInfoRequest } from '@/lib/pb/api_wg'
import { ClientStatus_Status } from '@/lib/pb/api_master'
import { Badge } from '@/components/ui/badge'
import { Activity, Wifi, WifiOff, Terminal } from 'lucide-react'
import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuTrigger } from '@/components/ui/context-menu'

const Node: React.FC<NodeProps<WGNode> & { onOpenTerminal?: (clientId: string, clientType: number) => void }> = ({
  id,
  data,
  selected,
  onOpenTerminal,
}) => {
  const connection = useConnection()
  const isTarget = !!connection.inProgress && connection.fromNode?.id !== id
  const { t } = useTranslation()

  const idLabel = t('wg.topologyNode.id', { id })
  const clientId = data.original?.clientId

  const handleOpenTerminal = React.useCallback(
    (event: React.MouseEvent) => {
      event.stopPropagation()
      event.preventDefault()

      if (clientId) {
        onOpenTerminal?.(clientId, ClientType.FRPC)
      }
    },
    [clientId, onOpenTerminal],
  )

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
    refetchInterval: 30000, // 30秒刷新一次
  })

  // 获取WireGuard运行时信息
  const { data: runtimeData } = useQuery({
    queryKey: ['wgRuntime', id],
    queryFn: async () => {
      return await getWireGuardRuntime(GetWireGuardRuntimeInfoRequest.create({ id: Number(id) }))
    },
    enabled: !!id && !isNaN(Number(id)),
    refetchInterval: 30000, // 30秒刷新一次
  })

  const clientStatus = clientStatusData?.clients[clientId || '']
  const isOnline = clientStatus?.status === ClientStatus_Status.ONLINE
  const ping = clientStatus?.ping
  const runtime = runtimeData?.wgDeviceRuntimeInfo

  return (
    <ContextMenu>
      <ContextMenuTrigger asChild>
        <div
          className={cn('customNode', selected && 'ring-2 ring-primary/60 rounded-md')}
          style={{ minWidth: 200, userSelect: 'none' }}
        >
          <div className="customNodeBody bg-card rounded-md border p-3 flex flex-col gap-2 text-sm shadow-md hover:shadow-lg transition-shadow">
            <div className="flex items-start justify-between gap-2">
              <div className="flex flex-col flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  {isOnline ? (
                    <Wifi className="h-3.5 w-3.5 text-green-500 flex-shrink-0" />
                  ) : (
                    <WifiOff className="h-3.5 w-3.5 text-red-500 flex-shrink-0" />
                  )}
                  <span className="font-medium truncate" title={data.original?.clientId}>
                    {data.original?.clientId || t('wg.topologyNode.unknown')}
                  </span>
                </div>
                <span className="font-mono text-xs text-muted-foreground truncate">
                  {data.original?.localAddress || '-'}
                </span>
                <div className="flex items-center gap-1 text-xs text-muted-foreground mt-1">
                  <span title={idLabel}>WG #{id}</span>
                  {ping !== undefined && (
                    <Badge variant="secondary" className="text-[10px] h-4 px-1">
                      {ping}ms
                    </Badge>
                  )}
                </div>
                {runtime && runtime.privateKey && (
                  <div className="flex items-center gap-1 mt-1">
                    <Activity className="h-3 w-3 text-blue-500" />
                    <span className="text-[10px] text-muted-foreground">
                      {runtime.peers?.length || 0} {t('wg.topologyNode.peers')}
                    </span>
                  </div>
                )}
                {data.original?.tags && data.original.tags.length > 0 && (
                  <div className="flex flex-wrap gap-1 mt-1">
                    {data.original.tags.slice(0, 2).map((tag) => (
                      <Badge key={tag} variant="outline" className="text-[10px] h-4 px-1">
                        {tag}
                      </Badge>
                    ))}
                    {data.original.tags.length > 2 && (
                      <Badge variant="outline" className="text-[10px] h-4 px-1">
                        +{data.original.tags.length - 2}
                      </Badge>
                    )}
                  </div>
                )}
              </div>
              <div className="drag-handle cursor-grab active:cursor-grabbing p-1 hover:bg-muted rounded flex-shrink-0">
                <div className="w-4 h-4 flex items-center justify-center">
                  <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor" className="text-muted-foreground">
                    <circle cx="2" cy="2" r="1.5" />
                    <circle cx="6" cy="2" r="1.5" />
                    <circle cx="10" cy="2" r="1.5" />
                    <circle cx="2" cy="6" r="1.5" />
                    <circle cx="6" cy="6" r="1.5" />
                    <circle cx="10" cy="6" r="1.5" />
                    <circle cx="2" cy="10" r="1.5" />
                    <circle cx="6" cy="10" r="1.5" />
                    <circle cx="10" cy="10" r="1.5" />
                  </svg>
                </div>
              </div>
            </div>
            {!connection.inProgress && <Handle className="customHandle" position={Position.Right} type="source" />}
            {(!connection.inProgress || isTarget) && (
              <Handle className="customHandle" position={Position.Left} type="target" isConnectableStart={false} />
            )}
          </div>
        </div>
      </ContextMenuTrigger>
      <ContextMenuContent className="w-48">
        <ContextMenuItem onClick={handleOpenTerminal} disabled={!clientId} className="gap-2">
          <Terminal className="h-4 w-4" />
          <span>{t('wg.contextMenu.openTerminal')}</span>
        </ContextMenuItem>
      </ContextMenuContent>
    </ContextMenu>
  )
}

export default Node
