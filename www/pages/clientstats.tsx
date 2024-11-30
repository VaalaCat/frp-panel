import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { SideBar } from '@/components/sidebar'
import { ClientStatsCard } from '@/components/stats/client_stats_card'

export default function ClientStatsPage() {
  return (
    <RootLayout header={<Header />} sidebar={<SideBar />}>
      <Providers>
        <div className="w-full">
          <div className="flex-1 flex-col">
            <ClientStatsCard />
          </div>
        </div>
      </Providers>
    </RootLayout>
  )
}
