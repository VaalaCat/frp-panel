import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { SideBar } from '@/components/sidebar'
import { FRPSFormCard } from '@/components/frps_card'

export default function ServerListPage() {
  return (
    <RootLayout header={<Header />} sidebar={<SideBar />}>
      <Providers>
        <div className="w-full">
          <div className="flex-1 flex-col">
            <FRPSFormCard />
          </div>
        </div>
      </Providers>
    </RootLayout>
  )
}
