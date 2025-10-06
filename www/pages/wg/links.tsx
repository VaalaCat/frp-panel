'use client'

import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { useTranslation } from 'react-i18next'
import { useState } from 'react'
import { IdInput } from '@/components/base/id_input'
import { WireGuardLinkList } from '@/components/wg/wireguard-link-list'

export default function WireGuardLinksPage() {
	const { t } = useTranslation()
	const [keyword, setKeyword] = useState('')
	const [trigger, setTrigger] = useState('')

	return (
		<Providers>
			<RootLayout mainHeader={<Header />}>
				<div className="w-full flex flex-col gap-4">
					<div className="flex flex-wrap items-center gap-2">
						<h1 className="text-xl font-semibold flex-1 min-w-[160px]">
							{t('wg.linkList.title')}
						</h1>
						<IdInput keyword={keyword} setKeyword={setKeyword} refetchTrigger={setTrigger} />
					</div>
					<WireGuardLinkList keyword={keyword} />
				</div>
			</RootLayout>
		</Providers>
	)
}


