'use client'

import React from 'react'
import type { NodeProps } from '@xyflow/react'
import type { ClientNode, NodeOperations } from './types'
import { cn } from '@/lib/utils'
import { useTranslation } from 'react-i18next'
import { useQuery } from '@tanstack/react-query'
import { getClientsStatus } from '@/api/platform'
import { ClientType } from '@/lib/pb/common'
import { ClientStatus_Status } from '@/lib/pb/api_master'
import { Badge } from '@/components/ui/badge'
import { Wifi, WifiOff, Terminal, FileText } from 'lucide-react'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
  ContextMenuSeparator,
} from '@/components/ui/context-menu'

const ClientNodeComponent: React.FC<NodeProps<ClientNode> & NodeOperations> = ({
  id,
  data,
  selected,
  onOpenTerminal,
  onOpenLog,
}) => {
  const { t } = useTranslation()

  const clientId = data.original?.id
  const clientType = ClientType.FRPC

  const handleOpenTerminal = React.useCallback(
    (event: React.MouseEvent) => {
      event.stopPropagation()
      event.preventDefault()

      if (clientId) {
        onOpenTerminal?.(clientId, clientType, id)
      }
    },
    [clientId, clientType, id, onOpenTerminal],
  )

  const handleOpenLog = React.useCallback(
    (event: React.MouseEvent) => {
      event.stopPropagation()
      event.preventDefault()

      if (clientId) {
        onOpenLog?.(clientId, clientType, id)
      }
    },
    [clientId, clientType, id, onOpenLog],
  )

  // 获取客户端状态
  const { data: clientStatusData } = useQuery({
    queryKey: ['clientStatus', clientId],
    queryFn: async () => {
      if (!clientId) return undefined
      return await getClientsStatus({
        clientIds: [clientId],
        clientType: clientType,
      })
    },
    enabled: !!clientId,
    refetchInterval: 30000, // 30秒刷新一次
  })

  const clientStatus = clientStatusData?.clients[clientId || '']
  const isOnline = clientStatus?.status === ClientStatus_Status.ONLINE
  const ping = clientStatus?.ping

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
                    <Wifi className="h-3.5 w-3.5 text-green-500 shrink-0" />
                  ) : (
                    <WifiOff className="h-3.5 w-3.5 text-red-500 shrink-0" />
                  )}
                  <span className="font-medium truncate text-blue-600" title={clientId}>
                    {clientId || t('canvas.client.unknown')}
                  </span>
                </div>
                <div className="flex items-center gap-1 text-xs text-muted-foreground mt-1">
                  <span>FRPC</span>
                  {ping !== undefined && (
                    <Badge variant="secondary" className="text-[10px] h-4 px-1">
                      {ping}ms
                    </Badge>
                  )}
                </div>
                {data.original?.stopped && (
                  <Badge variant="destructive" className="text-[10px] h-4 px-1 mt-1 w-fit">
                    {t('canvas.client.stopped')}
                  </Badge>
                )}
              </div>
              <div className="drag-handle cursor-grab active:cursor-grabbing p-1 hover:bg-muted rounded shrink-0">
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
          </div>
        </div>
      </ContextMenuTrigger>
      <ContextMenuContent className="w-48">
        <ContextMenuItem onClick={handleOpenTerminal} disabled={!clientId} className="gap-2">
          <Terminal className="h-4 w-4" />
          <span>{t('canvas.contextMenu.openTerminal')}</span>
        </ContextMenuItem>
        <ContextMenuSeparator />
        <ContextMenuItem onClick={handleOpenLog} disabled={!clientId} className="gap-2">
          <FileText className="h-4 w-4" />
          <span>{t('canvas.contextMenu.openLog')}</span>
        </ContextMenuItem>
      </ContextMenuContent>
    </ContextMenu>
  )
}

export default ClientNodeComponent
