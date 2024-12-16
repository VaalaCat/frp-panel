import { FRPCFormCard } from '@/components/frpc/frpc_card'
import { Providers } from '@/components/providers'
import { APITest } from '@/components/apitest'
import { Separator } from '@/components/ui/separator'
import { FRPSFormCard } from '@/components/frps/frps_card'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { createProxyConfig, listProxyConfig } from '@/api/proxy'
import { Button } from '@/components/ui/button'
import { TypedProxyConfig } from '@/types/proxy'
import { ClientConfig } from '@/types/client'
import { ProxyConfigList } from '@/components/proxy/proxy_config_list'
import { Input } from '@/components/ui/input'
import { useState } from 'react'

export default function Test() {
  const [name, setName] = useState<string>('')
  const [triggerRefetch, setTriggerRefetch] = useState<number>(0)

  function create() {
    const buffer = Buffer.from(
      JSON.stringify({
        proxies: [{
          name: name,
          type: 'tcp',
          localIP: '127.0.0.1',
          localPort: 1234,
          remotePort: 4321,
        } as TypedProxyConfig]
      } as ClientConfig),
    )
    const uint8Array: Uint8Array = new Uint8Array(buffer.buffer, buffer.byteOffset, buffer.byteLength);
    createProxyConfig({
      clientId: 'admin.c.test',
      config: uint8Array,
      serverId: 'default',
    })
      .then(() => {
        setTriggerRefetch(triggerRefetch + 1)
      })
      .catch((err) => {
        console.log(err)
      })
  }
  return (
    // <>
    // </>
    <Providers>
      <RootLayout mainHeader={<Header />}>
        <div className="w-full">
          <div className="flex flex-1 flex-col">
            <div className="flex flex-1 flex-row mb-2 gap-2">
              <Button onClick={create}>新建</Button>
              <Input value={name} onChange={(e) => setName(e.target.value)} ></Input>
            </div>
            <ProxyConfigList Keyword="" ProxyConfigs={[]} TriggerRefetch={triggerRefetch.toString()} />
          </div>
        </div>
      </RootLayout>
    </Providers>
  )
}
