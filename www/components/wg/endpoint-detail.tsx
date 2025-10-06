'use client'

import React from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import { useQuery } from '@tanstack/react-query'
import { getEndpoint, deleteEndpoint, listWireGuards } from '@/api/wg'
import { GetEndpointRequest, DeleteEndpointRequest, ListWireGuardsRequest } from '@/lib/pb/api_wg'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { toast } from 'sonner'
import EndpointEditDialog from './endpoint-edit-dialog'
import { Skeleton } from '@/components/ui/skeleton'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

const EndpointDetail: React.FC = () => {
	const params = useSearchParams()
	const router = useRouter()
	const { t } = useTranslation()

	const idParam = params.get('id')
	const endpointId = idParam ? Number(idParam) : undefined
	const [openEdit, setOpenEdit] = React.useState(false)
	const [openDelete, setOpenDelete] = React.useState(false)

	const { data, isLoading, refetch } = useQuery({
		queryKey: ['getEndpoint', endpointId],
		queryFn: () => getEndpoint(GetEndpointRequest.create({ id: endpointId! })),
		enabled: !!endpointId,
	})

	const endpoint = data?.endpoint

	const { data: wireguardOptions } = useQuery({
		queryKey: ['listWireGuards', endpoint?.clientId],
		queryFn: () =>
			listWireGuards(
				ListWireGuardsRequest.create({ page: 1, pageSize: 50, clientId: endpoint?.clientId }),
			),
		enabled: !!endpoint?.clientId,
	})

	const handleDelete = async () => {
		if (!endpointId) return
		try {
			await deleteEndpoint(DeleteEndpointRequest.create({ id: endpointId }))
			toast.success(t('common.success'))
			router.push('/wg/endpoints')
		} catch (error: any) {
			toast.error(error?.message ?? 'Error')
		}
	}

	if (!endpointId) {
		return (
			<div className="p-6">
				<Card>
					<CardHeader>
						<CardTitle>{t('wg.endpointDetail.invalid')}</CardTitle>
						<CardDescription>{t('wg.endpointDetail.invalidDesc')}</CardDescription>
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
						{endpoint?.host ? `${endpoint.host}:${endpoint.port}` : t('wg.endpointDetail.titleFallback')}
					</h1>
					<p className="text-muted-foreground">
						{t('wg.endpointDetail.subtitle', { id: endpointId })}
					</p>
				</div>
				<div className="flex flex-wrap gap-2">
					<Button variant="outline" onClick={() => router.push('/wg/endpoints')}>
						{t('wg.endpointDetail.back')}
					</Button>
					<Button variant="outline" onClick={() => setOpenEdit(true)} disabled={!endpoint}>
						{t('wg.endpointActions.edit')}
					</Button>
					<Button variant="destructive" onClick={() => setOpenDelete(true)}>
						{t('wg.endpoint.delete')}
					</Button>
				</div>
			</div>

			<AlertDialog open={openDelete} onOpenChange={setOpenDelete}>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle>{t('wg.endpoint.delete')}</AlertDialogTitle>
						<AlertDialogDescription>{t('wg.endpointDetail.deleteConfirm')}</AlertDialogDescription>
					</AlertDialogHeader>
					<AlertDialogFooter>
						<AlertDialogCancel>{t('common.cancel')}</AlertDialogCancel>
						<AlertDialogAction
							className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
							onClick={handleDelete}
						>
							{t('wg.endpoint.delete')}
						</AlertDialogAction>
					</AlertDialogFooter>
				</AlertDialogContent>
			</AlertDialog>

			<EndpointEditDialog
				clientId={endpoint?.clientId || ''}
				endpoint={endpoint || { id: endpointId }}
				open={openEdit}
				onOpenChange={setOpenEdit}
				onSaved={() => refetch()}
			/>

			<Tabs defaultValue="overview" className="space-y-4">
				<TabsList className="w-full flex flex-wrap">
					<TabsTrigger value="overview" className="flex-1 md:flex-none md:px-6">
						{t('wg.endpointDetail.tabsOverview')}
					</TabsTrigger>
					<TabsTrigger value="relations" className="flex-1 md:flex-none md:px-6">
						{t('wg.endpointDetail.tabsRelations')}
					</TabsTrigger>
				</TabsList>

				<TabsContent value="overview" className="space-y-4">
					<Card>
						<CardHeader>
							<CardTitle>{t('wg.endpointDetail.summaryTitle')}</CardTitle>
							<CardDescription>{t('wg.endpointDetail.summaryDesc')}</CardDescription>
						</CardHeader>
						<CardContent className="grid gap-4 md:grid-cols-2">
							<div className="space-y-2">
								<p className="text-sm text-muted-foreground">{t('wg.endpointForm.host')}</p>
								<p className="text-lg font-medium">
									{isLoading ? <Skeleton className="h-5 w-32" /> : endpoint?.host || t('wg.endpointDetail.noHost')}
								</p>
							</div>
							<div className="space-y-2">
								<p className="text-sm text-muted-foreground">{t('wg.endpointForm.port')}</p>
								<p className="text-lg font-medium">{endpoint?.port ?? '-'}</p>
							</div>
							<div className="space-y-2">
								<p className="text-sm text-muted-foreground">{t('wg.endpointDetail.clientId')}</p>
								<p className="text-lg font-medium">{endpoint?.clientId || '-'}</p>
							</div>
							<div className="space-y-2">
								<p className="text-sm text-muted-foreground">{t('wg.endpointDetail.wireguardId')}</p>
								<p className="text-lg font-medium">{endpoint?.wireguardId ? `#${endpoint.wireguardId}` : t('wg.interface.network_unassigned')}</p>
							</div>
						</CardContent>
					</Card>
				</TabsContent>

				<TabsContent value="relations" className="space-y-4">
					<Card>
						<CardHeader>
							<CardTitle>{t('wg.endpointDetail.relatedWireguardsTitle')}</CardTitle>
							<CardDescription>{t('wg.endpointDetail.relatedWireguardsDesc')}</CardDescription>
						</CardHeader>
						<CardContent className="space-y-2">
							{wireguardOptions?.wireguardConfigs?.length ? (
								wireguardOptions.wireguardConfigs.map((wg) => (
									<div key={wg.id} className="flex items-center justify-between text-sm">
										<span className="font-medium">{wg.interfaceName || wg.clientId}</span>
										<span className="text-muted-foreground">{wg.localAddress}</span>
									</div>
								))
							) : (
								<p className="text-sm text-muted-foreground">{t('wg.endpointDetail.relatedWireguardsEmpty')}</p>
							)}
						</CardContent>
					</Card>
				</TabsContent>
			</Tabs>
		</div>
	)
}

export default EndpointDetail


