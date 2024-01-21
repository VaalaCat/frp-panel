import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { ServerList } from '@/components/server_list'
import { Header } from '@/components/header'
import { SideBar } from '@/components/sidebar'
import { CreateServerDialog } from '@/components/server_create_dialog'

export default function ServerListPage() {
  return (
    <RootLayout header={<Header />} sidebar={<SideBar />}>
      <Providers>
        <div className="w-full">
          <div className="flex-1 flex-col">
            <div className="flex-1 flex-row mb-2">
              <CreateServerDialog />
            </div>
            <ServerList Servers={[]} />
          </div>
        </div>
      </Providers>
    </RootLayout>
  )
}
