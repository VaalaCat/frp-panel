'use client'

import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import NetworkDetail from '@/components/wg/network-detail'

export default function NetworkDetailPage() {
	return (
		<Providers>
			<RootLayout mainHeader={<Header />}>
				<div className="w-full flex flex-col gap-4">
					<NetworkDetail />
				</div>
			</RootLayout>
		</Providers>
	)
}


