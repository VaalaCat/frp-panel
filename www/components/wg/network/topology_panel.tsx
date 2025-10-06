'use client'

import React, { useEffect, useState } from 'react'
import { NetworkSelector } from '@/components/base/network-selector'
import { Button } from '@/components/ui/button'
import { useNetworkTopology } from './topology_hook'
import TopologyCanvas from './topology_canvas'
import TopologySidebar from './topology_sidebar'
import type { WGEdge, TopologyNode } from './types'
import { layoutNetwork } from './layout'
import { useTranslation } from 'react-i18next'
import { RefreshCw, Maximize2, Minimize2 } from 'lucide-react'
import { ReactFlowProvider } from '@xyflow/react'

export default function TopologyPanel() {
	const [networkID, setNetworkID] = React.useState<number | undefined>()
	const { topology, isFetching, refetch } = useNetworkTopology(networkID)
	const [selectedEdgeId, setSelectedEdgeId] = React.useState<string | undefined>()
	const [fullscreen, setFullscreen] = React.useState(false)
	const { t } = useTranslation()

	const [nodes, setNodes] = useState<TopologyNode[]>(topology.nodes)
	const [edges, setEdges] = useState<WGEdge[]>(topology.edges)

	const selectedEdge = React.useMemo(() => edges.find((e) => e.id === selectedEdgeId), [selectedEdgeId, edges])

	useEffect(() => {
		let alive = true
			; (async () => {
				const { nodes: n2, edges: e2 } = await layoutNetwork(topology.nodes as any, topology.edges)
				if (!alive) return
				setNodes(n2 as any)
				setEdges(e2)
			})()
		return () => {
			alive = false
		}
	}, [topology.nodes, topology.edges])

	return (
		<ReactFlowProvider>
			<div className="flex flex-col gap-3">
				<div className="flex flex-wrap items-center justify-between gap-2">
					<div className="flex flex-wrap items-end gap-2 flex-1 min-w-[240px] max-w-2xl">
						<div className="flex-1 min-w-[240px]">
							<NetworkSelector networkID={networkID} setNetworkID={setNetworkID} />
						</div>
						<Button
							variant="secondary"
							size="icon"
							onClick={() => refetch()}
							disabled={!networkID || isFetching}
							title={t('wg.topologyActions.refresh')}
						>
							<RefreshCw className={`h-4 w-4 ${isFetching ? 'animate-spin' : ''}`} />
						</Button>
					</div>
					<Button
						variant="outline"
						size="icon"
						onClick={() => setFullscreen(!fullscreen)}
						className="hidden md:flex"
						title={fullscreen ? t('wg.topologyActions.exitFullscreen') : t('wg.topologyActions.fullscreen')}
					>
						{fullscreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
					</Button>
				</div>
				<div className={`grid gap-3 ${fullscreen ? 'grid-cols-1' : 'grid-cols-1 lg:grid-cols-[1fr_320px]'}`}>
					<div className={fullscreen ? 'fixed inset-0 z-50 bg-background p-4' : ''}>
						{fullscreen && (
							<div className="absolute top-4 right-4 z-10">
								<Button
									variant="outline"
									size="icon"
									onClick={() => setFullscreen(false)}
								>
									<Minimize2 className="h-4 w-4" />
								</Button>
							</div>
						)}
						<TopologyCanvas
							data={{ nodes, edges }}
							onEdgeClick={(id) => setSelectedEdgeId(id)}
							setNodes={setNodes}
							setEdges={setEdges}
						/>
					</div>
					{!fullscreen && (
						<div className="hidden lg:block">
							<TopologySidebar selectedEdge={selectedEdge} />
						</div>
					)}
				</div>
				{!fullscreen && (
					<div className="lg:hidden">
						<TopologySidebar selectedEdge={selectedEdge} />
					</div>
				)}
			</div>
		</ReactFlowProvider>
	)
}


