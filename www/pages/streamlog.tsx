'use client'

import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { useEffect, useState } from 'react'
import { ClientSelector } from '@/components/base/client-selector'
import { parseStreaming } from '@/lib/stream'
import { Button } from '@/components/ui/button'
import { getClientsStatus } from '@/api/platform'
import { ClientType } from '@/lib/pb/common'
import dynamic from 'next/dynamic'
import { BaseSelector } from '@/components/base/selector'
import { ServerSelector } from '@/components/base/server-selector'
import LoadingCircle from '@/components/base/status'
import { useSearchParams } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
// import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { PlayCircle, StopCircle, RefreshCcw, Eraser } from 'lucide-react'
import { cn } from '@/lib/utils'

const LogTerminalComponent = dynamic(() => import('@/components/base/readonly-xterm'), {
  ssr: false,
})

export default function StreamLogPage() {
  const { t } = useTranslation()
  const [clientID, setClientID] = useState<string | undefined>(undefined)
  const [log, setLog] = useState<string | undefined>(undefined)
  const [clear, setClear] = useState<number>(0)
  const [enabled, setEnabled] = useState<boolean>(false)
  const [timeoutID, setTimeoutID] = useState<NodeJS.Timeout | null>(null)
  const [clientType, setClientType] = useState<ClientType>(ClientType.FRPS)
  const [status, setStatus] = useState<'loading' | 'success' | 'error' | undefined>()
  const [pkgs, setPkgs] = useState<string[]>([])

  const searchParams = useSearchParams()
  const paramClientID = searchParams.get('clientID')
  const paramClientType = searchParams.get('clientType')

  useEffect(() => {
    if (paramClientID) {
      setClientID(paramClientID)
    }
    if (paramClientType) {
      if (paramClientType == ClientType.FRPC.toString()) {
        setClientType(ClientType.FRPC)
      } else if (paramClientType == ClientType.FRPS.toString()) {
        setClientType(ClientType.FRPS)
      }
    }
    if (paramClientID && paramClientType) {
      setEnabled(true)
    }
  }, [paramClientID, paramClientType])

  useEffect(() => {
    setClear(Math.random())
    setStatus(undefined)
    if (!clientID || !enabled) {
      return
    }

    const abortController = new AbortController()
    setStatus('loading')

    parseStreaming(
      abortController,
      clientID,
      pkgs,
      setLog,
      (status: number) => {
        if (status === 200) {
          setStatus('success')
        } else {
          setStatus('error')
        }
      },
      () => {
        console.log('parseStreaming success')
        setStatus('success')
      },
    )

    return () => {
      abortController.abort('unmount')
      setEnabled(false)
    }
  }, [clientID, enabled, pkgs])

  const handleConnect = () => {
    if (enabled) {
      setEnabled(false)
    }
    if (timeoutID) {
      clearTimeout(timeoutID)
    }
    setTimeoutID(
      setTimeout(() => {
        setEnabled(true)
      }, 10),
    )
  }

  const handleRefresh = () => {
    setClear(Math.random())
    if (clientID) {
      getClientsStatus({ clientIds: [clientID], clientType: clientType })
    }
  }

  const handleDisconnect = () => {
    setEnabled(false)
    setClear(Math.random())
  }

  return (
    <Providers>
      <RootLayout mainHeader={<Header />}>
        <Card className="w-full h-[calc(100dvh_-_80px)] flex flex-col">
          <CardContent className="p-3 flex-1 flex flex-col gap-2 first-letter:">
            <div className="flex flex-wrap items-center gap-1.5 shrink-0">
              <div className="flex items-center gap-1.5">
                <LoadingCircle status={status} />
                <Button
                  disabled={!clientID}
                  variant={enabled ? 'destructive' : 'default'}
                  className="h-8 px-2 text-sm gap-1.5"
                  onClick={enabled ? handleDisconnect : handleConnect}
                >
                  {enabled ? (
                    <>
                      <StopCircle className="h-3.5 w-3.5" />
                      {t('common.disconnect')}
                    </>
                  ) : (
                    <>
                      <PlayCircle className="h-3.5 w-3.5" />
                      {t('common.connect')}
                    </>
                  )}
                </Button>

                <Button disabled={!clientID} variant="outline" className="h-8 w-8 p-0" onClick={handleRefresh}>
                  <RefreshCcw className="h-3.5 w-3.5" />
                </Button>

                <Button variant="outline" className="h-8 w-8 p-0" onClick={() => setClear(Math.random())}>
                  <Eraser className="h-3.5 w-3.5" />
                </Button>
              </div>

              <div className="flex items-center gap-1.5">
                <BaseSelector
                  dataList={[
                    { value: ClientType.FRPC.toString(), label: 'frpc' },
                    { value: ClientType.FRPS.toString(), label: 'frps' },
                  ]}
                  setValue={(value) => {
                    setClientType(value === ClientType.FRPC.toString() ? ClientType.FRPC : ClientType.FRPS)
                  }}
                  value={clientType.toString()}
                  label={t('common.clientType')}
                  className="h-8"
                />
              </div>
              <div className="flex items-center gap-1.5">
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
                  className="h-8"
                />
              </div>
            </div>

            <div className="flex flex-col gap-1.5 min-h-0 flex-1">
              {clientType === ClientType.FRPC && <ClientSelector clientID={clientID} setClientID={setClientID} />}
              {clientType === ClientType.FRPS && <ServerSelector serverID={clientID} setServerID={setClientID} />}

              <div className={cn('flex-1 min-h-0 overflow-hidden', 'border rounded-lg overflow-hidden')}>
                <LogTerminalComponent logs={log || ''} reset={clear} />
              </div>
            </div>
          </CardContent>
        </Card>
      </RootLayout>
    </Providers>
  )
}
