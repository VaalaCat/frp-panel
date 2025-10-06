'use client'

import React, { useState, useEffect } from 'react'
import type { NodeProps } from '@xyflow/react'
import { NodeResizer } from '@xyflow/react'
import type { LogNode, NodeOperations } from './types'
import { cn } from '@/lib/utils'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { FileText, Minimize2, Maximize2, X, PlayCircle, StopCircle, Eraser, Circle } from 'lucide-react'
import dynamic from 'next/dynamic'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ClientSelector } from '@/components/base/client-selector'
import { ServerSelector } from '@/components/base/server-selector'
import { Badge } from '@/components/ui/badge'
import { BaseSelector } from '@/components/base/selector'
import { ClientType } from '@/lib/pb/common'
import { parseStreaming } from '@/lib/stream'
import { ClientStatus_Status } from '@/lib/pb/api_master'
import { useQuery } from '@tanstack/react-query'
import { getClientsStatus } from '@/api/platform'

const LogTerminalComponent = dynamic(() => import('@/components/base/readonly-xterm'), {
  ssr: false,
})

const LogNodeComponent: React.FC<NodeProps<LogNode> & NodeOperations> = ({ id, data, selected, onDelete }) => {
  const { t } = useTranslation()
  const [clientId, setClientId] = useState(data.clientId || '')
  const [log, setLog] = useState<string | undefined>(undefined)
  const [clear, setClear] = useState<number>(0)
  const [enabled, setEnabled] = useState<boolean>(false)
  const [status, setStatus] = useState<'loading' | 'success' | 'error' | undefined>()
  const [pkgs, setPkgs] = useState<string[]>(data.pkgs || ['all'])

  const isFrps = data.clientType === ClientType.FRPS

  // 获取客户端状态
  const { data: clientStatusData } = useQuery({
    queryKey: ['clientStatus', clientId],
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
    setClear(Math.random())
    setStatus(undefined)
    if (!clientId || !enabled) {
      return
    }

    const abortController = new AbortController()
    setStatus('loading')

    parseStreaming(
      abortController,
      clientId,
      pkgs[0] === 'all' ? [] : pkgs,
      setLog,
      (status: number) => {
        if (status === 200) {
          setStatus('success')
        } else {
          setStatus('error')
        }
      },
      () => {
        setStatus('success')
      },
    )

    return () => {
      abortController.abort('unmount')
      setEnabled(false)
    }
  }, [clientId, enabled, pkgs])

  const handleDelete = () => {
    onDelete?.(id)
  }

  const handleConnect = () => {
    if (!clientId) return
    setEnabled(!enabled)
  }

  return (
    <div
      className={cn('logNode', selected && 'ring-2 ring-primary/60 rounded-lg')}
      style={{ userSelect: 'none', pointerEvents: 'all', width: '100%', height: '100%' }}
    >
      <NodeResizer isVisible={selected} minWidth={400} minHeight={320} maxWidth={1200} maxHeight={800} />
      <Card className={cn('shadow-lg border-2 h-full flex flex-col', 'w-full')}>
        <CardHeader
          className="p-3 pb-2 flex flex-row items-center justify-between space-y-0 drag-handle"
          style={{ cursor: 'move', userSelect: 'none' }}
        >
          <div className="flex items-center gap-2 flex-1 min-w-0">
            <FileText className="h-4 w-4 flex-shrink-0" />
            <CardTitle className="text-sm font-medium truncate">
              {clientId ? `${clientId}` : t('canvas.log.noClient')}
            </CardTitle>
            {clientId && (
              <Badge variant={isOnline ? 'default' : 'secondary'} className="text-[10px] h-4 px-1.5">
                <Circle className={cn('h-2 w-2 mr-1', isOnline ? 'fill-green-500' : 'fill-gray-500')} />
                {isOnline ? t('client.status_online') : t('client.status_offline')}
              </Badge>
            )}
            {status && (
              <Badge variant="outline" className="text-[10px] h-4 px-1.5">
                {status === 'success' && t('canvas.log.connected')}
                {status === 'loading' && t('canvas.log.connecting')}
                {status === 'error' && t('canvas.log.error')}
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
                  onClick={handleConnect}
                  title={enabled ? t('canvas.log.stop') : t('canvas.log.start')}
                >
                  {enabled ? (
                    <StopCircle className="h-3.5 w-3.5 text-red-500" />
                  ) : (
                    <PlayCircle className="h-3.5 w-3.5 text-green-500" />
                  )}
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-6 w-6"
                  onClick={() => setClear(Math.random())}
                  title={t('canvas.log.clear')}
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
          <div className="flex gap-2 nodrag">
            <div className="flex-1">
              {isFrps ? (
                <ServerSelector serverID={clientId} setServerID={setClientId} />
              ) : (
                <ClientSelector clientID={clientId} setClientID={setClientId} />
              )}
            </div>
            <div className="w-32">
              <BaseSelector
                dataList={[
                  { value: 'all', label: 'all' },
                  { value: 'frp', label: 'frp' },
                  { value: 'workerd', label: 'workerd' },
                ]}
                setValue={(value) => {
                  setPkgs([value])
                }}
                label={t('common.stream_log_pkgs_select')}
                className="h-10"
              />
            </div>
          </div>
          {clientId ? (
            <div
              className="flex-1 min-h-[200px] border rounded overflow-hidden bg-black nodrag"
              style={{ userSelect: 'text' }}
            >
              <LogTerminalComponent logs={log || ''} reset={clear} />
            </div>
          ) : (
            <div className="flex-1 min-h-[200px] border rounded bg-muted flex items-center justify-center">
              <p className="text-muted-foreground text-sm">{t('canvas.log.selectClient')}</p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

export default LogNodeComponent
