"use client"

import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { useEffect, useState } from 'react'
import { ClientSelector } from '@/components/base/client-selector'
import { Button } from '@/components/ui/button'
import { getClientsStatus } from '@/api/platform'
import { ClientType } from '@/lib/pb/common'
import dynamic from 'next/dynamic'
import { BaseSelector } from '@/components/base/selector'
import { ServerSelector } from '@/components/base/server-selector'
import LoadingCircle from '@/components/base/status'
import { ClientStatus } from '@/lib/pb/api_master'
import { useSearchParams } from 'next/navigation'
import { useTranslation } from 'react-i18next';
import { Card, CardContent } from '@/components/ui/card'
import { PlayCircle, StopCircle, RefreshCcw, Eraser, ExternalLink } from 'lucide-react'
import { cn } from '@/lib/utils'

const TerminalComponent = dynamic(() => import('@/components/base/read-write-xterm'), {
  ssr: false
})

export default function ConsolePage() {
  const { t } = useTranslation();
  const [clientID, setClientID] = useState<string | undefined>(undefined)
  const [clear, setClear] = useState<number>(0)
  const [enabled, setEnabled] = useState<boolean>(false)
  const [timeoutID, setTimeoutID] = useState<NodeJS.Timeout | null>(null);
  const [clientType, setClientType] = useState<ClientType>(ClientType.FRPS)
  const [status, setStatus] = useState<"loading" | "success" | "error" | undefined>()

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
    if (!clientID) {
      return;
    }
    if (!enabled) {
      return;
    }
    const abortController = new AbortController();
    setStatus("loading");

    return () => {
      abortController.abort("unmount");
      setEnabled(false);
    };
  }, [clientID, enabled]);

  const handleConnect = () => {
    if (enabled) {
      setEnabled(false)
      setStatus('error')
    } else {
      if (timeoutID) {
        clearTimeout(timeoutID)
      }
      setTimeoutID(setTimeout(() => {
        setEnabled(true)
        setStatus('success')
      }, 10))
    }
  }

  const handleRefresh = () => {
    if (!clientID) {
      return;
    }
    setClear(Math.random());
    getClientsStatus({ clientIds: [clientID!], clientType: clientType })
  }

  const handleNewWindow = () => {
    window.open(`/terminal?clientType=${clientType.toString()}&clientID=${clientID}`)
  }

  return (
    <Providers>
      <RootLayout mainHeader={<Header />}>
        <Card className="w-full h-[calc(100dvh_-_80px)] flex flex-col">
          <CardContent className="p-3 flex-1 flex flex-col gap-2">
            <div className="flex flex-wrap items-center gap-1.5 shrink-0"> 
              <div className="flex items-center gap-1.5">
                <LoadingCircle status={status} />
                <Button
                  disabled={!clientID}
                  variant={enabled ? "destructive" : "default"}
                  size="icon"
                  className="h-8 w-8"
                  onClick={handleConnect}
                >
                  {enabled ? (
                    <StopCircle className="h-3.5 w-3.5" />
                  ) : (
                    <PlayCircle className="h-3.5 w-3.5" />
                  )}
                </Button>
                
                <Button
                  disabled={!clientID}
                  variant="outline"
                  size="icon"
                  className="h-8 w-8"
                  onClick={handleRefresh}
                >
                  <RefreshCcw className="h-3.5 w-3.5" />
                </Button>

                <Button
                  variant="outline"
                  size="icon"
                  className="h-8 w-8"
                  onClick={() => setClear(Math.random())}
                >
                  <Eraser className="h-3.5 w-3.5" />
                </Button>

                <Button
                  disabled={!clientID}
                  variant="outline"
                  size="icon"
                  className="h-8 w-8"
                  onClick={handleNewWindow}
                >
                  <ExternalLink className="h-3.5 w-3.5" />
                </Button>
              </div>

              <div className="flex items-center gap-1.5">
                <BaseSelector
                  dataList={[
                    { value: ClientType.FRPC.toString(), label: "frpc" },
                    { value: ClientType.FRPS.toString(), label: "frps" }
                  ]}
                  setValue={(value) => {
                    setClientType(value === ClientType.FRPC.toString() ? ClientType.FRPC : ClientType.FRPS)
                  }}
                  value={clientType.toString()}
                  label={t('common.clientType')}
                  className="h-8"
                />
              </div>
            </div>

            <div className="flex flex-col gap-1.5 min-h-0 flex-1"> 
              {clientType === ClientType.FRPC && (
                <ClientSelector clientID={clientID} setClientID={setClientID} />
              )}
              {clientType === ClientType.FRPS && (
                <ServerSelector serverID={clientID} setServerID={setClientID} />
              )}
              
              <div className={cn(
                'flex-1 min-h-0 overflow-hidden',
                'border rounded-lg overflow-hidden'
              )}>
                <TerminalComponent
                  setStatus={setStatus}
                  isLoading={!enabled}
                  clientStatus={{
                    clientId: clientID,
                    clientType: clientType,
                    version: { platform: "linux" },
                  } as ClientStatus}
                  reset={clear} />
              </div>
            </div>
          </CardContent>
        </Card>
      </RootLayout>
    </Providers>
  )
}
