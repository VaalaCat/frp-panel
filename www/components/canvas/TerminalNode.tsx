'use client'

import React, { useState, useEffect } from 'react'
import type { NodeProps } from '@xyflow/react'
import { NodeResizer } from '@xyflow/react'
import type { TerminalNode, NodeOperations } from './types'
import { cn } from '@/lib/utils'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Terminal as TerminalIcon, Minimize2, Maximize2, X, RefreshCcw, Eraser, Circle } from 'lucide-react'
import dynamic from 'next/dynamic'
import { ClientStatus, ClientStatus_Status } from '@/lib/pb/api_master'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ClientSelector } from '@/components/base/client-selector'
import { ServerSelector } from '@/components/base/server-selector'
import { useQuery } from '@tanstack/react-query'
import { getClientsStatus } from '@/api/platform'
import { ClientType } from '@/lib/pb/common'
import { Badge } from '@/components/ui/badge'

const TerminalComponent = dynamic(() => import('@/components/base/read-write-xterm'), {
  ssr: false,
})

const TerminalNodeComponent: React.FC<NodeProps<TerminalNode> & NodeOperations> = ({
  id,
  data,
  selected,
  onDelete,
}) => {
  const { t } = useTranslation()
  const [clientId, setClientId] = useState(data.clientId || '')
  const [enabled, setEnabled] = useState(!!data.clientId)
  const [status, setStatus] = useState<'loading' | 'success' | 'error' | undefined>()
  const [clear, setClear] = useState(0)
  const [refetch, setRefetch] = useState(0)

  const isFrps = data.clientType === ClientType.FRPS

  // 获取客户端状态
  const { data: clientStatusData } = useQuery({
    queryKey: ['clientStatus', clientId, refetch],
    queryFn: async () => {
      if (!clientId) return undefined
      return await getClientsStatus({
        clientIds: [clientId],
        clientType: data.clientType,
      })
    },
    enabled: !!clientId,
    refetchInterval: 30000,
  })

  const clientStatus = clientStatusData?.clients[clientId]
  const isOnline = clientStatus?.status === ClientStatus_Status.ONLINE

  useEffect(() => {
    if (clientId && !enabled) {
      setEnabled(true)
    }
  }, [clientId, enabled])

  const handleDelete = () => {
    onDelete?.(id)
  }

  const handleRefresh = () => {
    setClear(Math.random())
    setRefetch(refetch + 1)
  }

  return (
    <div
      className={cn('terminalNode', selected && 'ring-2 ring-primary/60 rounded-lg')}
      style={{ userSelect: 'none', pointerEvents: 'all', width: '100%', height: '100%' }}
    >
      <NodeResizer isVisible={selected} minWidth={400} minHeight={320} maxWidth={1200} maxHeight={800} />
      <Card className={cn('shadow-lg border-2 h-full flex flex-col', 'w-full')}>
        <CardHeader
          className="p-3 pb-2 flex flex-row items-center justify-between space-y-0 drag-handle"
          style={{ cursor: 'move', userSelect: 'none' }}
        >
          <div className="flex items-center gap-2 flex-1 min-w-0">
            <TerminalIcon className="h-4 w-4 flex-shrink-0" />
            <CardTitle className="text-sm font-medium truncate">{clientId || t('canvas.terminal.noClient')}</CardTitle>
            {clientId && (
              <Badge variant={isOnline ? 'default' : 'secondary'} className="text-[10px] h-4 px-1.5">
                <Circle className={cn('h-2 w-2 mr-1', isOnline ? 'fill-green-500' : 'fill-gray-500')} />
                {isOnline ? t('client.status_online') : t('client.status_offline')}
              </Badge>
            )}
            {status && (
              <Badge variant="outline" className="text-[10px] h-4 px-1.5">
                {status === 'success' && t('canvas.terminal.connected')}
                {status === 'loading' && t('canvas.terminal.connecting')}
                {status === 'error' && t('canvas.terminal.error')}
              </Badge>
            )}
          </div>
          <div className="flex items-center gap-1 flex-shrink-0">
            {clientId && (
              <>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-6 w-6"
                  onClick={handleRefresh}
                  title={t('canvas.terminal.refresh')}
                >
                  <RefreshCcw className="h-3.5 w-3.5" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-6 w-6"
                  onClick={() => setClear(Math.random())}
                  title={t('canvas.terminal.clear')}
                >
                  <Eraser className="h-3.5 w-3.5" />
                </Button>
              </>
            )}
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6 text-destructive hover:text-destructive"
              onClick={handleDelete}
            >
              <X className="h-3.5 w-3.5" />
            </Button>
          </div>
        </CardHeader>
        <CardContent className="p-3 pt-0 flex flex-col gap-2 nodrag flex-1 min-h-0">
          <div className="nodrag">
            {isFrps ? (
              <ServerSelector serverID={clientId} setServerID={setClientId} />
            ) : (
              <ClientSelector clientID={clientId} setClientID={setClientId} />
            )}
          </div>
          {clientId ? (
            <div
              className="flex-1 min-h-[200px] border rounded overflow-hidden bg-black terminal-container nodrag"
              style={{ userSelect: 'text' }}
            >
              <TerminalComponent
                setStatus={setStatus}
                isLoading={!enabled}
                clientStatus={
                  {
                    clientId: clientId,
                    clientType: data.clientType,
                    version: { platform: 'linux' },
                  } as ClientStatus
                }
                reset={clear}
              />
            </div>
          ) : (
            <div className="flex-1 min-h-[200px] border rounded bg-muted flex items-center justify-center">
              <p className="text-muted-foreground text-sm">{t('canvas.terminal.selectClient')}</p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

export default TerminalNodeComponent
