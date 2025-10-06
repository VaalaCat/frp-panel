import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import PlatformInfo from '@/components/platforminfo'
import { useTranslation } from 'react-i18next'
// import Link from 'next/link'
// import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
// import { ROUTES } from '@/lib/routes'

// const quickLinks = [
//   { href: ROUTES.wg.networks, key: 'wg.nav.networks' },
//   { href: ROUTES.wg.wireguards, key: 'wg.nav.wireguards' },
//   { href: ROUTES.wg.endpoints, key: 'wg.nav.endpoints' },
//   { href: ROUTES.wg.links, key: 'wg.nav.links' },
// ]

export default function Home() {
  const { t } = useTranslation()
  return (
    <Providers>
      <RootLayout mainHeader={<Header />}>
        <div className="space-y-6">
          <PlatformInfo />
          {/* <Card>
            <CardHeader>
              <CardTitle>{t('wg.dashboard.quickEntry')}</CardTitle>
            </CardHeader>
            <CardContent className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
              {quickLinks.map((item) => (
                <Link key={item.href} href={item.href} className="group block rounded-md border p-4 hover:border-primary hover:bg-primary/5 transition-colors">
                  <div className="text-sm font-medium text-primary group-hover:text-primary/80">
                    {t(item.key)}
                  </div>
                </Link>
              ))}
            </CardContent>
          </Card> */}
        </div>
      </RootLayout>
    </Providers>
  )
}
