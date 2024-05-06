import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { ClientList } from '@/components/client_list'
import { Header } from '@/components/header'
import { SideBar } from '@/components/sidebar'
import { CreateClientDialog } from '@/components/client_create_dialog'
import { IdInput } from '@/components/id_input'
import { useState } from 'react'

export default function ClientListPage() {
  const [keyword, setKeyword] = useState('')
  const [triggerSearch, setTriggerSearch] = useState('')

  return (
    <RootLayout header={<Header />} sidebar={<SideBar />}>
      <Providers>
        <div className="w-full">
          <div className="flex flex-1 flex-col">
            <div className="flex flex-1 flex-row mb-2 gap-2">
              <CreateClientDialog />
              <IdInput setKeyword={setKeyword} keyword={keyword} refetchTrigger={setTriggerSearch} />
            </div>
            <ClientList Clients={[]} Keyword={keyword} TriggerRefetch={triggerSearch} />
          </div>
        </div>
      </Providers>
    </RootLayout>
  )
}
