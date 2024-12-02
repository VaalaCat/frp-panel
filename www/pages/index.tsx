import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { PlatformInfo } from '@/components/platforminfo'

export default function Home() {
  return (
    <Providers>
      <RootLayout mainHeader={<Header />}>
        <PlatformInfo />
      </RootLayout>
    </Providers>
  )
}
