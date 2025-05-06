import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { WorkerList } from '@/components/worker/worker_list'
import { CreateWorkerDialog } from '@/components/worker/worker_create_dialog'
import { IdInput } from '@/components/base/id_input'
import { useState } from 'react'
import WorkerEdit from '@/components/worker/edit'

export default function WorkerEditPage() {
  const [keyword, setKeyword] = useState('')
  const [triggerSearch, setTriggerSearch] = useState('')

  return (
    <Providers>
      <RootLayout mainHeader={<Header />}>
        <div className="w-full flex flex-col gap-4">
          <WorkerEdit />
        </div>
      </RootLayout>
    </Providers>
  )
}
