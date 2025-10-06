'use client'

import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import EndpointDetail from '@/components/wg/endpoint-detail'

export default function EndpointDetailPage() {
	return (
		<Providers>
			<RootLayout mainHeader={<Header />}>
				<div className="w-full flex flex-col gap-4">
					<EndpointDetail />
				</div>
			</RootLayout>
		</Providers>
	)
}


