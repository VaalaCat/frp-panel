'use client'

import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getNetworkTopology, listWireGuards } from '@/api/wg'
import { GetNetworkTopologyRequest, ListWireGuardsRequest } from '@/lib/pb/api_wg'
import type { WGEdge, WGNode, TopologyData } from './types'
import { WireGuardLinks, WireGuardLink, WireGuardConfig } from '@/lib/pb/types_wg'
import { nanoid } from 'nanoid'


export function useNetworkTopology(networkID?: number) {
	const { data, isFetching, refetch } = useQuery({
		queryKey: ['getNetworkTopology', networkID],
		queryFn: async () => {
			if (!networkID) return undefined
			return await getNetworkTopology(GetNetworkTopologyRequest.create({ id: networkID }))
		},
		enabled: !!networkID,
	})

	const { data: wgList } = useQuery({
		queryKey: ['listWireGuards', networkID],
		queryFn: async () => {
			if (!networkID) return undefined
			return await listWireGuards(
				ListWireGuardsRequest.create({ page: 1, pageSize: 500, networkId: networkID }),
			)
		},
		enabled: !!networkID,
	})

	const topology: TopologyData = useMemo(() => {
		if (!data?.adjs) return { nodes: [], edges: [] }
		const nodes: WGNode[] = []
		const edges: WGEdge[] = []

		const peerIds = new Set<string>()
		Object.entries(data.adjs).forEach(([fromIdStr, links]) => {
			peerIds.add(fromIdStr);
			(links as WireGuardLinks).links.forEach((lk: WireGuardLink) => {
				peerIds.add(lk.toWireguardId.toString())
				edges.push({
					id: `${lk.fromWireguardId}-${lk.toWireguardId}-${lk.id}-${nanoid()}`,
					source: String(fromIdStr),
					target: String(lk.toWireguardId),
					label: `${lk.latencyMs}ms / ${lk.upBandwidthMbps}↑ ${lk.downBandwidthMbps}↓`,
					animated: lk.active,
					type: 'wgEdge',
					data: { original: lk },
				})
			})
		})

		// label map from WireGuard configs in this network
		const idToLabel = new Map<number, string>()
		const configs: WireGuardConfig[] = (wgList?.wireguardConfigs ?? []) as WireGuardConfig[]
		configs?.forEach((cfg) => {
			const tagText = (cfg.tags ?? []).filter(Boolean).join(', ')
			const labelParts = [cfg.clientId, tagText ? `(${tagText})` : undefined, cfg.localAddress]
			const label = labelParts.filter(Boolean).join(' ')
			idToLabel.set(cfg.id, label || `WG ${cfg.id}`)
		})

		// compute circle layout
		const ids = Array.from(peerIds).sort((a, b) => Number(a) - Number(b))
		const radius = 220
		const centerX = 300
		const centerY = 260
		ids.forEach((pid, idx) => {
			const angle = (2 * Math.PI * idx) / Math.max(ids.length, 1)
			const x = centerX + radius * Math.cos(angle)
			const y = centerY + radius * Math.sin(angle)
			nodes.push({
				id: pid,
				type: 'wg',
				dragHandle: '.drag-handle',
				data: {
					label: idToLabel.get(Number(pid)) ?? `WG ${pid}`,
					original: wgList?.wireguardConfigs?.find((cfg) => cfg.id === Number(pid))
				},
				position: { x, y }
			})
		})

		return { nodes, edges }
	}, [data?.adjs, wgList?.wireguardConfigs])

	// 进行 ELK 自动布局：注意 React 规则，useMemo 内不能使用 async；这里使用立即执行的副作用去计算并缓存
	const laid = useMemo(() => ({ nodes: topology.nodes, edges: topology.edges }), [topology.nodes, topology.edges])

	// 返回布局后的数据（同步化：先用现有坐标快速渲染，随后 panel 会用 useEffect 同步由 layoutNetwork 产生的结果）
	// 由 TopologyPanel 触发异步布局并 setNodes/setEdges，避免在 hook 内部引入副作用

	return { topology: laid, isFetching, refetch }
}


