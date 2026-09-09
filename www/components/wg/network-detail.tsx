'use client'

import React from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import { useQuery } from '@tanstack/react-query'
import { getNetwork, deleteNetwork } from '@/api/wg'
import { GetNetworkRequest, DeleteNetworkRequest } from '@/lib/pb/api_wg'
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
import { WireGuardList } from './wireguard-list'
import { WireGuardLinkList } from './wireguard-link-list'
import { useNetworkTopology } from './network/topology_hook'
import TopologyCanvas from './network/topology_canvas'
import TopologySidebar from './network/topology_sidebar'
import type { TopologyNode, WGEdge, WGNode } from './network/types'
import { Skeleton } from '@/components/ui/skeleton'
import NetworkEditDialog from './network-edit-dialog'

const NetworkDetail: React.FC = () => {
	const params = useSearchParams()
	const router = useRouter()
	const { t } = useTranslation()
	const networkIdParam = params.get('networkId')
	const networkId = networkIdParam ? Number(networkIdParam) : undefined

	const [openEdit, setOpenEdit] = React.useState(false)
	const [nodes, setNodes] = React.useState<TopologyNode[]>([])
	const [edges, setEdges] = React.useState<WGEdge[]>([])
	const [spf, setSpf] = React.useState(true)
	const [selectedEdgeId, setSelectedEdgeId] = React.useState<string>()

	const { data, isLoading, refetch } = useQuery({
		queryKey: ['getNetwork', networkId],
		queryFn: () => getNetwork(GetNetworkRequest.create({ id: networkId! })),
		enabled: !!networkId,
		refetchOnWindowFocus: false,
	})

	const { topology, isFetching: topologyLoading, refetch: refetchTopology } = useNetworkTopology(networkId, spf)

	React.useEffect(() => {
		setNodes((prev) => {
			const customNodes = prev.filter(n => n.type !== 'wg')
			return [...topology.nodes, ...customNodes]
		})
		setEdges(topology.edges)
	}, [topology.nodes, topology.edges])

	const selectedEdge = React.useMemo(() => edges.find((edge) => edge.id === selectedEdgeId), [edges, selectedEdgeId])

	const network = data?.network

	const handleDelete = async () => {
		if (!networkId) return
		try {
			await deleteNetwork(DeleteNetworkRequest.create({ id: networkId }))
			toast.success(t('common.success'))
			router.push('/wg/networks')
		} catch (error: any) {
			toast.error(error?.message ?? 'Error')
		}
	}

	if (!networkId) {
		return (
			<div className="p-6">
				<Card>
					<CardHeader>
						<CardTitle>{t('wg.networkDetail.invalid')}</CardTitle>
						<CardDescription>{t('wg.networkDetail.invalidDesc')}</CardDescription>
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
						{network?.name || t('wg.networkDetail.titleFallback')}
					</h1>
					<p className="text-muted-foreground">
						{t('wg.networkDetail.subtitle', { id: networkId })}
					</p>
				</div>
				<Dialog>
					<div className="flex flex-wrap gap-2">
						<Button variant="outline" onClick={() => router.push('/wg/networks')}>
							{t('wg.networkDetail.back')}
						</Button>
						<Button variant="outline" onClick={() => setOpenEdit(true)} disabled={!network}>
							{t('wg.networkDetail.edit')}
						</Button>
						<DialogTrigger asChild>
							<Button variant="destructive">
								{t('wg.networkDetail.delete')}
							</Button>
						</DialogTrigger>
					</div>

					<DialogContent>
						<DialogHeader>
							<DialogTitle>{t('wg.networkDetail.delete')}</DialogTitle>
							<DialogDescription>{t('wg.networkDetail.deleteConfirm')}</DialogDescription>
						</DialogHeader>
						<DialogFooter>
							<DialogClose asChild>
								<Button variant="outline">
									{t('common.cancel')}
								</Button>
							</DialogClose>
							<DialogClose asChild>
								<Button variant="destructive" onClick={handleDelete}>
									{t('wg.networkDetail.delete')}
								</Button>
							</DialogClose>
						</DialogFooter>
					</DialogContent>
				</Dialog>
			</div>

			<NetworkEditDialog
				open={openEdit}
				onOpenChange={setOpenEdit}
				network={network ? { id: network.id!, name: network.name || '', cidr: network.cidr || '', acl: network.acl } : { id: networkId, name: '', cidr: '' }}
				onSaved={() => {
					refetch()
					refetchTopology()
				}}
			/>

			<Tabs defaultValue="overview" className="space-y-4">
				<TabsList className="w-full flex flex-wrap">
					<TabsTrigger value="overview" className="flex-1 md:flex-none md:px-6">
						{t('wg.networkDetail.tabsOverview')}
					</TabsTrigger>
					<TabsTrigger value="wireguards" className="flex-1 md:flex-none md:px-6">
						{t('wg.networkDetail.tabsWireguards')}
					</TabsTrigger>
					<TabsTrigger value="links" className="flex-1 md:flex-none md:px-6">
						{t('wg.networkDetail.tabsLinks')}
					</TabsTrigger>
					<TabsTrigger value="topology" className="flex-1 md:flex-none md:px-6">
						{t('wg.networkDetail.tabsTopology')}
					</TabsTrigger>
				</TabsList>

				<TabsContent value="overview" className="space-y-4">
					<Card>
						<CardHeader>
							<CardTitle>{t('wg.networkDetail.summaryTitle')}</CardTitle>
							<CardDescription>{t('wg.networkDetail.summaryDesc')}</CardDescription>
						</CardHeader>
						<CardContent className="grid gap-4 md:grid-cols-2">
							<div className="space-y-2">
								<p className="text-sm text-muted-foreground">{t('wg.networkForm.name')}</p>
								<p className="text-lg font-medium">
									{isLoading ? <Skeleton className="h-5 w-32" /> : network?.name || t('wg.networkDetail.unnamed')}
								</p>
							</div>
							<div className="space-y-2">
								<p className="text-sm text-muted-foreground">{t('wg.networkForm.cidr')}</p>
								<p className="text-lg font-medium">
									{isLoading ? <Skeleton className="h-5 w-28" /> : network?.cidr || t('wg.networkDetail.noCidr')}
								</p>
							</div>
							<div className="space-y-2">
								<p className="text-sm text-muted-foreground">{t('wg.networkDetail.id')}</p>
								<p className="text-lg font-medium">#{networkId}</p>
							</div>
							{/* <div className="space-y-2">
								<p className="text-sm text-muted-foreground">{t('wg.networkDetail.updatedAt')}</p>
								<p className="text-lg font-medium">
									{network?.updatedAt ? format(new Date(network.updatedAt as any), 'yyyy-MM-dd HH:mm:ss') : t('wg.networkDetail.noTime')}
								</p>
							</div> */}
						</CardContent>
					</Card>
					<Card>
						<CardHeader>
							<CardTitle>{t('wg.networkDetail.aclTitle')}</CardTitle>
						</CardHeader>
						<CardContent>
							{network?.acl ? (
								<pre className="bg-muted text-xs p-4 rounded border overflow-auto max-h-64 whitespace-pre-wrap">
									{JSON.stringify(network.acl, null, 2)}
								</pre>
							) : (
								<div className="text-sm text-muted-foreground">{t('wg.networkDetail.aclEmpty')}</div>
							)}
						</CardContent>
					</Card>
				</TabsContent>

				<TabsContent value="wireguards" className="space-y-4">
					<Card>
						<CardHeader className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
							<div>
								<CardTitle>{t('wg.networkDetail.wireguardsTitle')}</CardTitle>
								<CardDescription>{t('wg.networkDetail.wireguardsDesc')}</CardDescription>
							</div>
						</CardHeader>
						<CardContent>
							<WireGuardList clientId={undefined} networkId={networkId} onChanged={() => { }} />
						</CardContent>
					</Card>
				</TabsContent>

				<TabsContent value="links" className="space-y-4">
					<Card>
						<CardHeader>
							<CardTitle>{t('wg.networkDetail.linksTitle')}</CardTitle>
							<CardDescription>{t('wg.networkDetail.linksDesc')}</CardDescription>
						</CardHeader>
						<CardContent>
							<WireGuardLinkList networkId={networkId} />
						</CardContent>
					</Card>
				</TabsContent>

				<TabsContent value="topology" className="space-y-4">
					<Card>
						<CardHeader className="flex flex-col md:flex-row md:items-center md:justify-between gap-2">
							<div>
								<CardTitle>{t('wg.topologySidebar.title')}</CardTitle>
								<CardDescription>{t('wg.networkDetail.topologyDesc')}</CardDescription>
							</div>
							<div className="flex gap-2">
								<Button variant="outline" onClick={() => refetchTopology()} disabled={topologyLoading}>
									{topologyLoading ? t('wg.topologyActions.loading') : t('wg.topologyActions.refresh')}
								</Button>
								<Button variant="outline" onClick={() => setSpf(!spf)} disabled={topologyLoading}>
									{spf ? t('wg.topologyActions.spf') : t('wg.topologyActions.full')}
								</Button>
							</div>
						</CardHeader>
						<CardContent className={selectedEdge ? "grid gap-4 md:grid-cols-[2fr_1fr]" : ""}>
							<div className="h-[700px] md:h-[800px] lg:h-[900px]">
								<TopologyCanvas
									data={{ nodes, edges }}
									onEdgeClick={setSelectedEdgeId}
									onPaneClick={() => setSelectedEdgeId(undefined)}
									setNodes={setNodes}
									setEdges={setEdges}
									onLinkCreated={() => refetchTopology()}
								/>
							</div>
							{selectedEdge && (
								<div>
									<TopologySidebar selectedEdge={selectedEdge} />
								</div>
							)}
						</CardContent>
					</Card>
				</TabsContent>
			</Tabs>
		</div>
	)
}

export default NetworkDetail


