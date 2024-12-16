import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { FRPSFormCard } from '@/components/frps/frps_card'

export default function ServerListPage() {
  return (
    <Providers>
      <RootLayout mainHeader={<Header />}>
        <div className="w-full flex items-center justify-center">
          <div className="flex-1 flex-col max-w-2xl">
            <FRPSFormCard />
          </div>
        </div>
      </RootLayout>
    </Providers>
  )
}
