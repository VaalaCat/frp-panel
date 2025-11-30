'use client'

import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getNetworkTopology, listWireGuards } from '@/api/wg'
import { GetNetworkTopologyRequest, ListWireGuardsRequest } from '@/lib/pb/api_wg'
import type { WGEdge, WGNode, TopologyData } from './types'
import type { WireGuardLinks, WireGuardLink, WireGuardConfig } from '@/lib/pb/types_wg'
import { calculateConnectionQuality } from './layout'

/**
 * 使用网络拓扑数据的 Hook
 */
export function useNetworkTopology(networkID?: number, spf?: boolean) {
  // 获取拓扑连接数据
  const { data: topologyData, isFetching: isTopologyFetching, refetch: refetchTopology } = useQuery({
    queryKey: ['getNetworkTopology', networkID, spf],
    queryFn: async () => {
      if (!networkID) return undefined
      return await getNetworkTopology(GetNetworkTopologyRequest.create({ id: networkID, spf: spf ?? true }))
    },
    enabled: !!networkID,
    staleTime: 10000, // 10秒内数据视为新鲜
  })

  // 获取 WireGuard 配置列表
  const { data: wgList, isFetching: isWgListFetching, refetch: refetchWgList } = useQuery({
    queryKey: ['listWireGuards', networkID],
    queryFn: async () => {
      if (!networkID) return undefined
      return await listWireGuards(
        ListWireGuardsRequest.create({
          page: 1,
          pageSize: 1000,
          networkId: networkID
        })
      )
    },
    enabled: !!networkID,
    staleTime: 10000,
  })

  // 统一的刷新函数
  const refetch = () => {
    refetchTopology()
    refetchWgList()
  }

  // 构建拓扑图数据
  const topology: TopologyData = useMemo(() => {
    if (!topologyData?.adjs || !wgList?.wireguardConfigs) {
      return { nodes: [], edges: [] }
    }

    const nodes: WGNode[] = []
    const edges: WGEdge[] = []
    const nodeIds = new Set<string>()

    // 创建配置映射
    const configMap = new Map<number, WireGuardConfig>()
    wgList.wireguardConfigs.forEach((cfg) => {
      configMap.set(cfg.id, cfg as WireGuardConfig)
    })

    // 检测双向连接
    const bidirectionalEdges = new Set<string>()
    const adjs = topologyData.adjs as Record<string, WireGuardLinks>
    Object.entries(adjs).forEach(([fromIdStr, links]) => {
      const fromId = Number(fromIdStr)
        ; (links as WireGuardLinks).links.forEach((link: WireGuardLink) => {
          const reverseLinks = adjs[String(link.toWireguardId)] as WireGuardLinks | undefined
          if (reverseLinks) {
            const hasReverse = reverseLinks.links.some(
              (l: WireGuardLink) => l.toWireguardId === fromId
            )
            if (hasReverse) {
              const edgeKey1 = `${fromId}-${link.toWireguardId}`
              const edgeKey2 = `${link.toWireguardId}-${fromId}`
              bidirectionalEdges.add(edgeKey1)
              bidirectionalEdges.add(edgeKey2)
            }
          }
        })
    })

    // 构建边
    Object.entries(topologyData.adjs).forEach(([fromIdStr, links]) => {
      const fromId = Number(fromIdStr)
      nodeIds.add(fromIdStr)

        ; (links as WireGuardLinks).links.forEach((link: WireGuardLink) => {
          const toId = link.toWireguardId
          nodeIds.add(String(toId))

          const edgeKey = `${fromId}-${toId}`
          const quality = calculateConnectionQuality(
            link.latencyMs ?? 999,
            Math.max(link.upBandwidthMbps || 0, link.downBandwidthMbps || 0)
          )

          // 使用唯一的边ID：结合link.id、源和目标ID
          const edgeId = link.id && link.id > 0
            ? `edge-${link.id}`
            : `edge-${fromId}-${toId}-${link.latencyMs || 0}-${Date.now()}`

          edges.push({
            id: edgeId,
            source: String(fromId),
            target: String(toId),
            type: 'wgEdge',
            animated: link.active,
            data: {
              link,
              quality,
              isBidirectional: bidirectionalEdges.has(edgeKey),
            },
            label: link.active
              ? `${link.latencyMs}ms | ${link.upBandwidthMbps}↑ ${link.downBandwidthMbps}↓`
              : '未激活',
          })
        })
    })

    // 构建节点
    const sortedIds = Array.from(nodeIds).sort((a, b) => Number(a) - Number(b))

    sortedIds.forEach((idStr) => {
      const id = Number(idStr)
      const config = configMap.get(id)

      // 构建节点标签
      let label = `WG #${id}`
      if (config) {
        const parts = []
        if (config.clientId) parts.push(config.clientId)
        if (config.tags && config.tags.length > 0) {
          parts.push(`(${config.tags.slice(0, 2).join(', ')})`)
        }
        if (config.localAddress) parts.push(config.localAddress)
        label = parts.join(' ') || label
      }

      nodes.push({
        id: idStr,
        type: 'wg',
        position: { x: 0, y: 0 }, // 将由布局算法计算
        data: {
          label,
          config,
        },
        // dragHandle: '.drag-handle', // 移除 dragHandle 限制，让整个节点可拖拽
      })
    })

    return { nodes, edges }
  }, [topologyData?.adjs, wgList?.wireguardConfigs])

  const isFetching = isTopologyFetching || isWgListFetching

  return {
    topology,
    isFetching,
    refetch,
    hasData: topology.nodes.length > 0,
  }
}
