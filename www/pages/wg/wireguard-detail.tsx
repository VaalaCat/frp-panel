'use client'

import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import WireGuardDetail from '@/components/wg/wireguard-detail'

export default function WireGuardDetailPage() {
	return (
		<Providers>
			<RootLayout mainHeader={<Header />}>
				<div className="w-full flex flex-col gap-4">
					<WireGuardDetail />
				</div>
			</RootLayout>
		</Providers>
	)
}


