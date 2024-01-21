import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { ClientList } from '@/components/client_list'
import { Header } from '@/components/header'
import { SideBar } from '@/components/sidebar'
import { CreateClientDialog } from '@/components/client_create_dialog'

export default function ClientListPage() {
  return (
    <RootLayout header={<Header />} sidebar={<SideBar />}>
      <Providers>
        <div className="w-full">
          <div className="flex-1 flex-col">
            <div className="flex-1 flex-row mb-2">
              <CreateClientDialog />
            </div>
            <ClientList Clients={[]} />
          </div>
        </div>
      </Providers>
    </RootLayout>
  )
}
