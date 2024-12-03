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

const TerminalComponentProps = dynamic(() => import('@/components/base/read-write-xterm'), {
  ssr: false
})

export default function ClientStatsPage() {
  const [clientID, setClientID] = useState<string | undefined>(undefined)
  const [clear, setClear] = useState<number>(0)
  const [enabled, setEnabled] = useState<boolean>(false)
  const [timeoutID, setTimeoutID] = useState<NodeJS.Timeout | null>(null);
  const [clientType, setClientType] = useState<ClientType>(ClientType.FRPC)
  const [status, setStatus] = useState<"loading" | "success" | "error" | undefined>()

  useEffect(() => {
    setClientID(undefined)
  }, [clientType])

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

  return (
    <Providers>
      <RootLayout mainHeader={<Header />}>
        <div className="w-full">
          <div className="flex-1 flex-col space-y-2">
            <div className="flex flex-1 flex-row gap-2 items-center">
              <div className='items-center'>
                <LoadingCircle status={status} />
              </div>
              <Button
                variant="outline"
                onClick={() => {
                  if (enabled) { setEnabled(false) }
                  if (timeoutID) { clearTimeout(timeoutID) }
                  setTimeoutID(setTimeout(() => { setEnabled(true) }, 10))
                }}>连接</Button>
              <Button onClick={() => {
                setClear(Math.random());
                getClientsStatus({ clientIds: [clientID!], clientType: clientType })
              }}>刷新</Button>
              <Button variant="destructive" onClick={() => {
                setEnabled(false)
                setClear(Math.random());
              }}>断开</Button>
              <BaseSelector
                dataList={[{ value: ClientType.FRPC.toString(), label: "frpc" }, { value: ClientType.FRPS.toString(), label: "frps" }]}
                setValue={(value) => { if (value === ClientType.FRPC.toString()) { setClientType(ClientType.FRPC) } else { setClientType(ClientType.FRPS) } }}
                value={clientType.toString()}
                label="客户端类型"
              />
            </div>
            {clientType === ClientType.FRPC && <ClientSelector clientID={clientID} setClientID={setClientID} />}
            {clientType === ClientType.FRPS && <ServerSelector serverID={clientID} setServerID={setClientID} />}
            <div className='flex-1 h-[calc(100dvh_-_180px)]'>
              <TerminalComponentProps
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
        </div>
      </RootLayout>
    </Providers>
  )
}
