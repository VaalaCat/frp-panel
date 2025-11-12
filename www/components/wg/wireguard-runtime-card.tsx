"use client"

import { useTranslation } from 'react-i18next'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import { ArrowUpRight, RefreshCw } from 'lucide-react'
import type { WGDeviceRuntimeInfo, WGPeerRuntimeInfo } from '@/lib/pb/types_wg'
import { formatBytes } from '@/lib/utils'
import { useRouter } from 'next/router'

export default function WireGuardRuntimeCard({
	runtime,
	loading,
	onRefresh,
}: {
	runtime?: WGDeviceRuntimeInfo | null
	loading?: boolean
	onRefresh?: () => void
}) {
	const { t } = useTranslation()
	const peers = runtime?.peers ?? []

	return (
		<Card>
			<CardHeader className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
				<div>
					<CardTitle>{runtime?.clientId ? runtime.clientId : t('wg.runtime.title')}</CardTitle>
					<p className="text-sm text-muted-foreground">{t('wg.runtime.subtitle')}</p>
				</div>
				<Button variant="outline" size="sm" onClick={onRefresh} disabled={loading}>
					<RefreshCw className={`mr-2 h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
					{t('wg.runtime.refresh')}
				</Button>
			</CardHeader>
			<CardContent className="space-y-4">
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
					<RuntimeStat
						title={t('wg.runtime.listen_port')}
						value={runtime?.listenPort ? `:${runtime.listenPort}` : '-'}
						loading={loading}
					/>
					<RuntimeStat
						title={t('wg.runtime.virt_ip')}
						value={runtime?.virtualIp ?? '-'}
						loading={loading}
					/>
					<RuntimeStat title={t('wg.runtime.peer_count')} value={peers.length} loading={loading} />
					<RuntimeStat
						title={t('wg.runtime.status')}
						value={runtime ? t('wg.runtime.status_ok') : t('wg.runtime.status_unknown')}
						loading={loading}
					/>
				</div>
				<div className="space-y-3">
					<h4 className="text-sm font-medium text-muted-foreground">{t('wg.runtime.peer_list')}</h4>
					<div className="space-y-2">
						{loading ? (
							<Skeleton className="h-12 w-full" />
						) : peers.length ? (
							peers.sort((a, b) => a.clientId.localeCompare(b.clientId)).map((peer) => <PeerItem key={peer.publicKey} peer={peer}
								wireguardId={
									runtime?.peerConfigMap?.[peer.publicKey]?.id ?? 0
								} />)
						) : (
							<p className="text-sm text-muted-foreground">{t('wg.runtime.peer_empty')}</p>
						)}
					</div>
				</div>
			</CardContent>
		</Card>
	)
}

function RuntimeStat({ title, value, loading }: { title: string; value: React.ReactNode; loading?: boolean }) {
	return (
		<div className="rounded-lg border bg-muted/40 p-3">
			<p className="text-xs uppercase tracking-wide text-muted-foreground">{title}</p>
			<div className="mt-2 text-lg font-semibold">
				{loading ? <Skeleton className="h-5 w-24" /> : value}
			</div>
		</div>
	)
}

function PeerItem({ peer, wireguardId }: { peer: WGPeerRuntimeInfo; wireguardId: number }) {
	const { t } = useTranslation()
	const router = useRouter()
	const lastHandshake = peer.lastHandshakeTimeSec ? new Date(Number(peer.lastHandshakeTimeSec) * 1000) : undefined

	return (
		<div className="flex flex-col gap-2 rounded-md border p-3">
			<div className="flex flex-wrap items-center gap-2">
				<span className="font-mono text-sm truncate max-w-[240px]" title={peer.publicKey}>
					{peer.clientId || peer.publicKey || t('wg.runtime.peer_unknown')}
				</span>
				<Badge variant="outline">
					<p className='font-mono text-xs truncate max-w-[240px] text-nowrap w-fit'>{peer.publicKey || t('wg.runtime.peer_unknown')}</p>
				</Badge>
				<Button
					variant="ghost"
					size="icon"
					className="h-6 w-6 rounded-full"
					onClick={(e) => {
						e.preventDefault()
						e.stopPropagation()
						router.push({ pathname: '/wg/wireguard-detail', query: { id: wireguardId } })
					}}
					aria-label={t('wg.wireguardActions.view')}
				>
					<ArrowUpRight className="h-3.5 w-3.5" />
				</Button>
			</div>
			<div className="grid gap-2 text-xs text-muted-foreground sm:grid-cols-2">
				<div>
					<span className="font-medium text-foreground">{t('wg.runtime.peer_endpoint')}:</span> {peer.endpoint ?? '-'}
				</div>
				<div>
					<span className="font-medium text-foreground">{t('wg.runtime.peer_last_handshake')}:</span>{' '}
					{lastHandshake ? lastHandshake.toLocaleString() : t('wg.runtime.peer_no_handshake')}
				</div>
				<div>
					<span className="font-medium text-foreground">TX:</span> {formatBytes(peer.txBytes ?? 0)}
				</div>
				<div>
					<span className="font-medium text-foreground">RX:</span> {formatBytes(peer.rxBytes ?? 0)}
				</div>
				<div className='flex flex-row gap-2'>
					<span className="font-medium text-foreground">Route</span>
					<div className='flex flex-wrap items-center gap-2'>
						{
							peer.allowedIps.map((ip) => (
								<Badge key={ip} variant="outline">
									<span className="font-mono text-xs truncate max-w-[240px] text-nowrap w-fit">{ip}</span>
								</Badge>
							))
						}
					</div>
				</div>
			</div>
		</div>
	)
}

