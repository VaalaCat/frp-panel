import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { SideBar } from '@/components/sidebar'
import { FRPCFormCard } from '@/components/frpc_card'

export default function ClientEditPage() {
  return (
    <RootLayout header={<Header />} sidebar={<SideBar />}>
      <Providers>
        <div className="w-full">
          <div className="flex-1 flex-col">
            <FRPCFormCard />
          </div>
        </div>
      </Providers>
    </RootLayout>
  )
}
