import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { ServerList } from '@/components/server_list'
import { Header } from '@/components/header'
import { SideBar } from '@/components/sidebar'
import { CreateServerDialog } from '@/components/server_create_dialog'
import { useState } from 'react'
import { IdInput } from '@/components/id_input'

export default function ServerListPage() {
  const [keyword, setKeyword] = useState('')
  const [triggerSearch, setTriggerSearch] = useState('')

  return (
    <RootLayout header={<Header />} sidebar={<SideBar />}>
      <Providers>
        <div className="w-full">
          <div className="flex flex-1 flex-col">
            <div className="flex flex-1 flex-row mb-2 gap-2">
              <CreateServerDialog />
              <IdInput setKeyword={setKeyword} keyword={keyword} refetchTrigger={setTriggerSearch} />
            </div>
            <ServerList Servers={[]} Keyword={keyword} TriggerRefetch={triggerSearch} />
          </div>
        </div>
      </Providers>
    </RootLayout>
  )
}
