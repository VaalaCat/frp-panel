import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { FRPCFormCard } from '@/components/frpc/frpc_card'

export default function ClientEditPage() {
  return (
    <Providers>
      <RootLayout mainHeader={<Header />}>
        <div className="w-full">
          <div className="flex-1 flex-col">
            <FRPCFormCard />
          </div>
        </div>
      </RootLayout>
    </Providers>
  )
}
