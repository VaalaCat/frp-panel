import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { WorkerList } from '@/components/worker/worker_list'
import { CreateWorkerDialog } from '@/components/worker/worker_create_dialog'
import { IdInput } from '@/components/base/id_input'
import { useState } from 'react'

export default function WorkerListPage() {
  const [keyword, setKeyword] = useState('')
  const [triggerSearch, setTriggerSearch] = useState('')

  return (
    <Providers>
      <RootLayout mainHeader={<Header />}>
        <div className="w-full flex flex-col gap-4">
          <div className="flex items-center gap-2">
            <CreateWorkerDialog refetchTrigger={setTriggerSearch} />
            <IdInput keyword={keyword} setKeyword={setKeyword} refetchTrigger={setTriggerSearch} />
          </div>
          <WorkerList initialWorkers={[]} initialTotal={0} triggerRefetch={triggerSearch} keyword={keyword} />
        </div>
      </RootLayout>
    </Providers>
  )
}
