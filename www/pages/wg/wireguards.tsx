'use client'

import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { useTranslation } from 'react-i18next'
import { useState } from 'react'
import { IdInput } from '@/components/base/id_input'
import { WireGuardList } from '@/components/wg/wireguard-list'

export default function WireGuardsPage() {
	const { t } = useTranslation()
	const [keyword, setKeyword] = useState('')
	const [trigger, setTrigger] = useState('')

	return (
		<Providers>
			<RootLayout mainHeader={<Header />}>
				<div className="w-full flex flex-col gap-4">
					<div className="flex flex-wrap items-center gap-2">
						<h1 className="text-xl font-semibold flex-1 min-w-[160px]">
							{t('wg.wireguardList.title')}
						</h1>
						<IdInput keyword={keyword} setKeyword={setKeyword} refetchTrigger={setTrigger} />
					</div>
					<WireGuardList keyword={keyword} onChanged={() => setTrigger(String(Math.random()))} />
				</div>
			</RootLayout>
		</Providers>
	)
}


