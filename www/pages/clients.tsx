import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { ClientList } from '@/components/frpc/client_list'
import { Header } from '@/components/header'
import { CreateClientDialog } from '@/components/frpc/client_create_dialog'
import { IdInput } from '@/components/base/id_input'
import { useState } from 'react'

export default function ClientListPage() {
  const [keyword, setKeyword] = useState('')
  const [triggerSearch, setTriggerSearch] = useState('')

  return (
    <Providers>
      <RootLayout mainHeader={<Header />}>
        <div className="w-full">
          <div className="flex flex-1 flex-col">
            <div className="flex flex-1 flex-row mb-2 gap-2">
              <CreateClientDialog />
              <IdInput setKeyword={setKeyword} keyword={keyword} refetchTrigger={setTriggerSearch} />
            </div>
            <ClientList Clients={[]} Keyword={keyword} TriggerRefetch={triggerSearch} />
          </div>
        </div>
      </RootLayout>
    </Providers>
  )
}
