'use client'

import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import { useTranslation } from 'react-i18next'
import { useState } from 'react'
import { makeRandomTrigger } from '@/lib/utils'
import { IdInput } from '@/components/base/id_input'
import { NetworkList } from '@/components/wg/network-list'
import { Button } from '@/components/ui/button'
import { Plus } from 'lucide-react'
import { NetworkCreateDialog } from '@/components/wg/network-create-dialog'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export default function NetworksPage() {
	const { t } = useTranslation()
	const [keyword, setKeyword] = useState('')
	const [trigger, setTrigger] = useState(makeRandomTrigger())
	const [createOpen, setCreateOpen] = useState(false)
	const [summary, setSummary] = useState<{ total: number }>({ total: 0 })
	return (
		<Providers>
			<RootLayout mainHeader={<Header />}>
				<div className="w-full flex flex-col gap-4">
					<div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
						<div className="space-y-1">
							<h1 className="text-2xl font-semibold">
								{t('wg.networkList.title')}
							</h1>
							<p className="text-sm text-muted-foreground">
								{t('wg.networkList.subtitle')}
							</p>
						</div>
						<div className="flex flex-wrap items-center gap-2">
							<IdInput keyword={keyword} setKeyword={setKeyword} refetchTrigger={() => setTrigger(makeRandomTrigger())} />
							<Button onClick={() => setCreateOpen(true)} size="sm" className="gap-2">
								<Plus className="h-4 w-4" />
								{t('wg.networkCreate.trigger')}
							</Button>
						</div>
					</div>
					<Card>
						<CardHeader>
							<CardTitle>{t('wg.networkList.metricsTitle')}</CardTitle>
						</CardHeader>
						<CardContent className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
							<div className="rounded-lg border bg-muted/30 p-4">
								<div className="text-sm text-muted-foreground">{t('wg.networkList.metricsTotal')}</div>
								<div className="text-2xl font-semibold">{summary.total}</div>
							</div>
						</CardContent>
					</Card>
					<NetworkList
						keyword={keyword}
						refreshToken={trigger}
						onChanged={() => setTrigger(makeRandomTrigger())}
						onSummary={setSummary}
					/>
				</div>
			</RootLayout>
			<NetworkCreateDialog
				open={createOpen}
				onOpenChange={setCreateOpen}
				onCreated={() => {
					setTrigger(makeRandomTrigger())
				}}
			/>
		</Providers>
	)
}


