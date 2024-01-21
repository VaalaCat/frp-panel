import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { ServerList } from '@/components/server_list'
import { Header } from '@/components/header'
import { SideBar } from '@/components/sidebar'
import { PlatformInfo } from '@/components/platforminfo'

export default function Home() {
  return (
    <RootLayout header={<Header />} sidebar={<SideBar />}>
      <Providers>
        <PlatformInfo />
      </Providers>
    </RootLayout>
  )
}
