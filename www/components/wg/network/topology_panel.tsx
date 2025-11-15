'use client'

import React, { useEffect, useState, useCallback } from 'react'
import { NetworkSelector } from '@/components/base/network-selector'
import { Button } from '@/components/ui/button'
import { useNetworkTopology } from './topology_hook'
import TopologyCanvas from './topology_canvas'
import TopologySidebar from './topology_sidebar'
import type { WGEdge, TopologyNode } from './types'
import { layoutNetwork } from './layout'
import { useTranslation } from 'react-i18next'
import { RefreshCw, Maximize2, Minimize2, Info } from 'lucide-react'
import { ReactFlowProvider } from '@xyflow/react'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { cn } from '@/lib/utils'

export default function TopologyPanel() {
	const [networkID, setNetworkID] = useState<number | undefined>()
	const { topology, isFetching, refetch, hasData } = useNetworkTopology(networkID)
	const [selectedEdgeId, setSelectedEdgeId] = useState<string | undefined>()
	const [fullscreen, setFullscreen] = useState(false)
	const { t } = useTranslation()

	const [nodes, setNodes] = useState<TopologyNode[]>([])
	const [edges, setEdges] = useState<WGEdge[]>([])

	const selectedEdge = React.useMemo(
		() => edges.find((e) => e.id === selectedEdgeId),
		[selectedEdgeId, edges]
	)

	// 布局nodes
	useEffect(() => {
		let cancelled = false

		const doLayout = async () => {
			if (topology.nodes.length === 0) {
				setNodes([])
				return
			}

			try {
				const { nodes: layoutedNodes } = await layoutNetwork(topology.nodes, topology.edges)
				if (!cancelled) {
					setNodes(layoutedNodes)
				}
			} catch (error) {
				console.error('Layout error:', error)
				if (!cancelled) {
					setNodes(topology.nodes)
				}
			}
		}

		doLayout()

		return () => {
			cancelled = true
		}
	}, [topology.nodes, topology.edges])

	// 同步edges
	useEffect(() => {
		setEdges(topology.edges)
	}, [topology.edges])

	// 全屏切换
	const toggleFullscreen = useCallback(() => {
		setFullscreen((prev) => !prev)
	}, [])

	// 清除选中的边
	const clearSelection = useCallback(() => {
		setSelectedEdgeId(undefined)
	}, [])

	return (
		<ReactFlowProvider>
			<div className="flex flex-col gap-4 h-full">
				{/* 顶部工具栏 */}
				<div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 p-4 bg-gradient-to-r from-card to-card/80 rounded-xl border-2 border-border shadow-lg">
					<div className="flex-1 min-w-[240px]">
						<NetworkSelector networkID={networkID} setNetworkID={setNetworkID} />
					</div>

					<div className="flex gap-2">
						<Button
							variant="secondary"
							size="icon"
							onClick={() => refetch()}
							disabled={!networkID || isFetching}
							title={t('wg.topologyActions.refresh')}
							className="shadow-md hover:shadow-lg transition-shadow"
						>
							<RefreshCw className={cn('h-4 w-4', isFetching && 'animate-spin')} />
						</Button>

						<Button
							variant="outline"
							size="icon"
							onClick={toggleFullscreen}
							title={fullscreen ? t('wg.topologyActions.exitFullscreen') : t('wg.topologyActions.fullscreen')}
							className="shadow-md hover:shadow-lg transition-shadow"
						>
							{fullscreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
						</Button>
					</div>
				</div>

				{/* 提示信息 */}
				{!networkID && (
					<Alert>
						<Info className="h-4 w-4" />
						<AlertDescription>
							{t('wg.topology.selectNetwork')}
						</AlertDescription>
					</Alert>
				)}

				{networkID && !hasData && !isFetching && (
					<Alert>
						<Info className="h-4 w-4" />
						<AlertDescription>
							{t('wg.topology.noData')}
						</AlertDescription>
					</Alert>
				)}

				{/* 主内容区 */}
				{networkID && hasData && (fullscreen ? (
					/* 全屏模式 */
					<div className="fixed inset-0 z-50 bg-background p-6 flex flex-col gap-4">
						<div className="flex items-center justify-between">
							<h2 className="text-2xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
								{t('wg.topology.title')}
							</h2>
							<div className="flex gap-2">
								<Button
									variant="secondary"
									size="icon"
									onClick={() => refetch()}
									disabled={!networkID || isFetching}
								>
									<RefreshCw className={cn('h-4 w-4', isFetching && 'animate-spin')} />
								</Button>
								<Button variant="outline" size="icon" onClick={toggleFullscreen}>
									<Minimize2 className="h-4 w-4" />
								</Button>
							</div>
						</div>

						<div className="flex-1 min-h-0">
							<TopologyCanvas
								data={{ nodes, edges }}
								onEdgeClick={setSelectedEdgeId}
								onPaneClick={() => setSelectedEdgeId(undefined)}
								setNodes={setNodes}
								setEdges={setEdges}
								fullscreen={fullscreen}
								onFullscreenToggle={toggleFullscreen}
								onLinkCreated={() => refetch()}
							/>
						</div>
					</div>
				) : (
					/* 常规模式 */
					<div className={cn(
						'grid gap-4 transition-all duration-300',
						selectedEdge
							? 'grid-cols-1 lg:grid-cols-[1fr_360px] xl:grid-cols-[1fr_400px]'
							: 'grid-cols-1'
					)}>
						{/* Canvas区域 */}
						<div className="w-full h-[700px] md:h-[800px] lg:h-[900px]">
							<TopologyCanvas
								data={{ nodes, edges }}
								onEdgeClick={setSelectedEdgeId}
								onPaneClick={() => setSelectedEdgeId(undefined)}
								setNodes={setNodes}
								setEdges={setEdges}
								fullscreen={fullscreen}
								onFullscreenToggle={toggleFullscreen}
								onLinkCreated={() => refetch()}
							/>
						</div>

						{/* Sidebar区域 - 仅在大屏幕且有选中边时显示 */}
						{selectedEdge && (
							<div className="hidden lg:block h-[800px] xl:h-[900px]">
								<TopologySidebar selectedEdge={selectedEdge} />
							</div>
						)}
					</div>
				))}

				{/* 移动端Sidebar - 在底部显示 */}
				{networkID && hasData && !fullscreen && selectedEdge && (
					<div className="lg:hidden">
						<TopologySidebar selectedEdge={selectedEdge} />
					</div>
				)}
			</div>
		</ReactFlowProvider>
	)
}
