'use client'

import React from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { getWireGuard, deleteWireGuard, getWireGuardRuntime } from '@/api/wg'
import { GetWireGuardRequest, DeleteWireGuardRequest, GetWireGuardRuntimeInfoRequest } from '@/lib/pb/api_wg'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { toast } from 'sonner'
import { format } from 'date-fns'
import WireGuardEditDialog from './wireguard-edit-dialog'
import { EndpointList } from './endpoint-list'
import { Skeleton } from '@/components/ui/skeleton'
import { Badge } from '@/components/ui/badge'
import WireGuardRuntimeCard from './wireguard-runtime-card'

const WireGuardDetail: React.FC = () => {
	const params = useSearchParams()
	const router = useRouter()
	const { t } = useTranslation()
	const queryClient = useQueryClient()

	const idParam = params.get('id')
	const wireguardId = idParam ? Number(idParam) : undefined

	const [openEdit, setOpenEdit] = React.useState(false)

	const { data, isLoading, refetch } = useQuery({
		queryKey: ['getWireGuard', wireguardId],
		queryFn: () => getWireGuard(GetWireGuardRequest.create({ id: wireguardId! })),
		enabled: !!wireguardId,
	})

	const {
		data: runtimeData,
		isLoading: runtimeLoading,
		refetch: refetchRuntime,
	} = useQuery({
		queryKey: ['getWireGuardRuntime', wireguardId],
		queryFn: () => getWireGuardRuntime(GetWireGuardRuntimeInfoRequest.create({ id: wireguardId! })),
		enabled: !!wireguardId,
		refetchInterval: 30000,
	})

	const wireguard = data?.wireguardConfig

	const handleDelete = async () => {
		if (!wireguardId) return
		try {
			await deleteWireGuard(DeleteWireGuardRequest.create({ id: wireguardId }))
			toast.success(t('common.success'))
			router.push('/wg/wireguards')
		} catch (error: any) {
			toast.error(error?.message ?? 'Error')
		}
	}

	if (!wireguardId) {
		return (
			<div className="p-6">
				<Card>
					<CardHeader>
						<CardTitle>{t('wg.wireguardDetail.invalid')}</CardTitle>
						<CardDescription>{t('wg.wireguardDetail.invalidDesc')}</CardDescription>
					</CardHeader>
				</Card>
			</div>
		)
	}

	return (
		<div className="space-y-6">
			<div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
				<div className="space-y-1">
					<h1 className="text-2xl font-semibold">
						{wireguard?.interfaceName || t('wg.wireguardDetail.titleFallback')}
					</h1>
					<p className="text-muted-foreground">
						{t('wg.wireguardDetail.subtitle', { id: wireguardId })}
					</p>
				</div>
				<Dialog>
					<div className="flex flex-wrap gap-2">
						<Button variant="outline" onClick={() => router.push('/wg/wireguards')}>
							{t('wg.wireguardDetail.back')}
						</Button>
						<Button variant="outline" onClick={() => setOpenEdit(true)} disabled={!wireguard}>
							{t('wg.wireguardDetail.edit')}
						</Button>
						<DialogTrigger asChild>
							<Button variant="destructive">
								{t('wg.wireguardDetail.delete')}
							</Button>
						</DialogTrigger>
					</div>

					<DialogContent>
						<DialogHeader>
							<DialogTitle>{t('wg.wireguardDetail.delete')}</DialogTitle>
							<DialogDescription>{t('wg.wireguardDetail.deleteConfirm')}</DialogDescription>
						</DialogHeader>
						<DialogFooter>
							<DialogClose asChild>
								<Button variant="outline">
									{t('common.cancel')}
								</Button>
							</DialogClose>
							<DialogClose asChild>
								<Button variant="destructive" onClick={handleDelete}>
									{t('wg.wireguardDetail.delete')}
								</Button>
							</DialogClose>
						</DialogFooter>
					</DialogContent>
				</Dialog>
			</div>

			<WireGuardEditDialog
				clientId={wireguard?.clientId || ''}
				wg={wireguard || ({} as any)}
				open={openEdit}
				onOpenChange={setOpenEdit}
				onUpdated={() => {
					refetch()
					queryClient.invalidateQueries({ queryKey: ['listWireGuards'] })
				}}
			/>

			<Tabs defaultValue="overview" className="space-y-4">
				<TabsList className="w-full flex flex-wrap">
					<TabsTrigger value="overview" className="flex-1 md:flex-none md:px-6">
						{t('wg.wireguardDetail.tabsOverview')}
					</TabsTrigger>
					<TabsTrigger value="endpoints" className="flex-1 md:flex-none md:px-6">
						{t('wg.wireguardDetail.tabsEndpoints')}
					</TabsTrigger>
				</TabsList>

				<TabsContent value="overview" className="space-y-4">
					<Card>
						<CardHeader>
							<CardTitle>{t('wg.wireguardDetail.summaryTitle')}</CardTitle>
							<CardDescription>{t('wg.wireguardDetail.summaryDesc')}</CardDescription>
						</CardHeader>
						<CardContent className="grid gap-4 md:grid-cols-2">
							<div className="space-y-2">
								<p className="text-sm text-muted-foreground">{t('wg.wireguardForm.interfaceName')}</p>
								<p className="text-lg font-medium">
									{isLoading ? <Skeleton className="h-5 w-32" /> : wireguard?.interfaceName || t('wg.wireguardDetail.noName')}
								</p>
							</div>
							<div className="space-y-2">
								<p className="text-sm text-muted-foreground">{t('wg.wireguardForm.localAddress')}</p>
								<p className="text-lg font-medium">
									{isLoading ? <Skeleton className="h-5 w-32" /> : wireguard?.localAddress || t('wg.wireguardDetail.noAddress')}
								</p>
							</div>
							<div className="space-y-2">
								<p className="text-sm text-muted-foreground">{t('wg.wireguardForm.port')}</p>
								<p className="text-lg font-medium">{wireguard?.listenPort ?? '-'}</p>
							</div>
							<div className="space-y-2">
								<p className="text-sm text-muted-foreground">{t('wg.wireguardDetail.networkId')}</p>
								<p className="text-lg font-medium">
									{wireguard?.networkId ? `#${wireguard.networkId}` : t('wg.interface.network_unassigned')}
								</p>
							</div>
						</CardContent>
					</Card>
					<WireGuardRuntimeCard
						runtime={runtimeData?.wgDeviceRuntimeInfo}
						loading={runtimeLoading}
						onRefresh={() => refetchRuntime()}
					/>
					<Card>
						<CardHeader>
							<CardTitle>{t('wg.wireguardDetail.tagsTitle')}</CardTitle>
						</CardHeader>
						<CardContent className="flex flex-wrap gap-2">
							{wireguard?.tags?.length ? (
								wireguard.tags.map((tag) => (
									<Badge key={tag} variant="secondary">
										#{tag}
									</Badge>
								))
							) : (
								<span className="text-sm text-muted-foreground">{t('wg.wireguardDetail.tagsEmpty')}</span>
							)}
						</CardContent>
					</Card>
				</TabsContent>

				<TabsContent value="endpoints" className="space-y-4">
					<Card>
						<CardHeader>
							<CardTitle>{t('wg.wireguardDetail.endpointsTitle')}</CardTitle>
							<CardDescription>{t('wg.wireguardDetail.endpointsDesc')}</CardDescription>
						</CardHeader>
						<CardContent>
							<EndpointList wireguardId={wireguardId} />
						</CardContent>
					</Card>
				</TabsContent>
			</Tabs>
		</div>
	)
}

export default WireGuardDetail


