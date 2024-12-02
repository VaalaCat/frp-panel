import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { ServerList } from '@/components/frps/server_list'
import { Header } from '@/components/header'
import { CreateServerDialog } from '@/components/frps/server_create_dialog'
import { useState } from 'react'
import { IdInput } from '@/components/base/id_input'

export default function ServerListPage() {
  const [keyword, setKeyword] = useState('')
  const [triggerSearch, setTriggerSearch] = useState('')

  return (
    <Providers>
      <RootLayout mainHeader={<Header />}>
        <div className="w-full">
          <div className="flex flex-1 flex-col">
            <div className="flex flex-1 flex-row mb-2 gap-2">
              <CreateServerDialog />
              <IdInput setKeyword={setKeyword} keyword={keyword} refetchTrigger={setTriggerSearch} />
            </div>
            <ServerList Servers={[]} Keyword={keyword} TriggerRefetch={triggerSearch} />
          </div>
        </div>
      </RootLayout>
    </Providers>
  )
}
