'use client'

import React from 'react'
import type { NodeProps } from '@xyflow/react'
import type { ServerNode, NodeOperations } from './types'
import { cn } from '@/lib/utils'
import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'
import { Server, Terminal, FileText } from 'lucide-react'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
  ContextMenuSeparator,
} from '@/components/ui/context-menu'
import { ClientType } from '@/lib/pb/common'
import { useQuery } from '@tanstack/react-query'
import { getClientsStatus } from '@/api/platform'
import { ClientStatus_Status } from '@/lib/pb/api_master'

const ServerNodeComponent: React.FC<NodeProps<ServerNode> & NodeOperations> = ({
  id,
  data,
  selected,
  onOpenTerminal,
  onOpenLog,
}) => {
  const { t } = useTranslation()

  const serverId = data.original?.id
  const clientType = ClientType.FRPS

  const handleOpenTerminal = React.useCallback(
    (event: React.MouseEvent) => {
      event.stopPropagation()
      event.preventDefault()

      if (serverId) {
        onOpenTerminal?.(serverId, clientType, id)
      }
    },
    [serverId, clientType, id, onOpenTerminal],
  )

  const handleOpenLog = React.useCallback(
    (event: React.MouseEvent) => {
      event.stopPropagation()
      event.preventDefault()

      if (serverId) {
        onOpenLog?.(serverId, clientType, id)
      }
    },
    [serverId, clientType, id, onOpenLog],
  )

  const { data: serverStatusData } = useQuery({
    queryKey: ['serverStatus', serverId],
    queryFn: async () => {
      if (!serverId) return undefined
      return await getClientsStatus({
        clientIds: [serverId],
        clientType: clientType,
      })
    },
    enabled: !!serverId,
    refetchInterval: 30000, // 30秒刷新一次
  })

  const serverStatus = serverStatusData?.clients[serverId || '']
  const isOnline = serverStatus?.status === ClientStatus_Status.ONLINE
  const ping = serverStatus?.ping

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
                  <Server className={cn('h-3.5 w-3.5 flex-shrink-0', isOnline ? 'text-green-500' : 'text-gray-500')} />
                  <span className="font-medium truncate text-purple-600" title={serverId}>
                    {serverId || t('canvas.server.unknown')}
                  </span>
                </div>
                {data.original?.ip && (
                  <span className="font-mono text-xs text-muted-foreground truncate mt-1">{data.original.ip}</span>
                )}
                <div className="flex items-center gap-1 text-xs text-muted-foreground mt-1">
                  <span>FRPS</span>
                  <Badge variant={isOnline ? 'default' : 'secondary'} className="text-[10px] h-4 px-1">
                    {isOnline ? t('canvas.server.online') : t('canvas.server.offline')}
                  </Badge>
                </div>
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
          </div>
        </div>
      </ContextMenuTrigger>
      <ContextMenuContent className="w-48">
        <ContextMenuItem onClick={handleOpenTerminal} disabled={!serverId} className="gap-2">
          <Terminal className="h-4 w-4" />
          <span>{t('canvas.contextMenu.openTerminal')}</span>
        </ContextMenuItem>
        <ContextMenuSeparator />
        <ContextMenuItem onClick={handleOpenLog} disabled={!serverId} className="gap-2">
          <FileText className="h-4 w-4" />
          <span>{t('canvas.contextMenu.openLog')}</span>
        </ContextMenuItem>
      </ContextMenuContent>
    </ContextMenu>
  )
}

export default ServerNodeComponent
