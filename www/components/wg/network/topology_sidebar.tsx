'use client'

import React from 'react'
import { useTranslation } from 'react-i18next'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import type { WGEdge } from './types'
import { Activity, TrendingUp, TrendingDown, Clock, MapPin, Zap } from 'lucide-react'
import { cn } from '@/lib/utils'

interface TopologySidebarProps {
	selectedEdge?: WGEdge
}

export default function TopologySidebar({ selectedEdge }: TopologySidebarProps) {
	const { t } = useTranslation()

	if (!selectedEdge) {
		return null
	}

	const link = selectedEdge.data?.link
	const quality = selectedEdge.data?.quality || 'fair'
	const isBidirectional = selectedEdge.data?.isBidirectional || false

	if (!link) {
		return null
	}

	const qualityConfig = {
		excellent: { label: '优秀', color: 'bg-emerald-500', textColor: 'text-emerald-700' },
		good: { label: '良好', color: 'bg-green-500', textColor: 'text-green-700' },
		fair: { label: '一般', color: 'bg-yellow-500', textColor: 'text-yellow-700' },
		poor: { label: '较差', color: 'bg-red-500', textColor: 'text-red-700' },
	}

	const config = qualityConfig[quality]

	return (
		<Card className="h-full shadow-lg border-2">
			<CardHeader className="pb-3">
				<CardTitle className="text-lg flex items-center justify-between">
					<span>连接详情</span>
					<Badge
						variant={link.active ? 'default' : 'secondary'}
						className={cn(link.active && 'bg-emerald-500')}
					>
						{link.active ? '激活' : '未激活'}
					</Badge>
				</CardTitle>
			</CardHeader>
			<CardContent className="space-y-4">
				{/* 连接质量 */}
				<div className="space-y-2">
					<div className="text-sm font-medium text-muted-foreground">连接质量</div>
					<div className="flex items-center gap-2">
						<div className={cn('w-3 h-3 rounded-full', config.color)} />
						<span className={cn('font-semibold', config.textColor)}>{config.label}</span>
						{isBidirectional && (
							<Badge variant="outline" className="ml-auto">
								<Activity className="h-3 w-3 mr-1" />
								双向
							</Badge>
						)}
					</div>
				</div>

				<Separator />

				{/* 性能指标 */}
				<div className="grid grid-cols-2 gap-3">
					<MetricCard
						icon={<Zap className="h-4 w-4" />}
						label="延迟"
						value={`${link.latencyMs}ms`}
						variant={link.latencyMs < 50 ? 'good' : link.latencyMs < 100 ? 'normal' : 'bad'}
					/>
					<MetricCard
						icon={<TrendingUp className="h-4 w-4" />}
						label="上行带宽"
						value={`${link.upBandwidthMbps}Mbps`}
						variant="normal"
					/>
					<MetricCard
						icon={<TrendingDown className="h-4 w-4" />}
						label="下行带宽"
						value={`${link.downBandwidthMbps}Mbps`}
						variant="normal"
					/>
					<MetricCard
						icon={<Activity className="h-4 w-4" />}
						label="状态"
						value={link.active ? '在线' : '离线'}
						variant={link.active ? 'good' : 'bad'}
					/>
				</div>

				<Separator />

				{/* 端点信息 */}
				{link.toEndpoint && (
					<div className="space-y-2">
						<div className="text-sm font-medium text-muted-foreground flex items-center gap-2">
							<MapPin className="h-4 w-4" />
							目标端点
						</div>
						<div className="bg-muted/50 rounded-lg p-3 space-y-1">
							<div className="flex items-center justify-between text-sm">
								<span className="text-muted-foreground">主机</span>
								<span className="font-mono font-medium">{link.toEndpoint.host}</span>
							</div>
							<div className="flex items-center justify-between text-sm">
								<span className="text-muted-foreground">端口</span>
								<span className="font-mono font-medium">{link.toEndpoint.port}</span>
							</div>
						</div>
					</div>
				)}

				{/* 连接信息 */}
				<div className="space-y-2">
					<div className="text-sm font-medium text-muted-foreground">连接信息</div>
					<div className="space-y-2 text-sm">
						<InfoRow label="链接 ID" value={`#${link.id}`} />
						<InfoRow label="源节点" value={`WG #${link.fromWireguardId}`} />
						<InfoRow label="目标节点" value={`WG #${link.toWireguardId}`} />
					</div>
				</div>

			</CardContent>
		</Card>
	)
}

function MetricCard({
	icon,
	label,
	value,
	variant = 'normal',
}: {
	icon: React.ReactNode
	label: string
	value: string
	variant?: 'good' | 'normal' | 'bad'
}) {
	const variantClasses = {
		good: 'bg-emerald-50 dark:bg-emerald-950/20 border-emerald-200 dark:border-emerald-800',
		normal: 'bg-muted/50 border-border',
		bad: 'bg-red-50 dark:bg-red-950/20 border-red-200 dark:border-red-800',
	}

	return (
		<div className={cn('rounded-lg border p-3 space-y-1', variantClasses[variant])}>
			<div className="flex items-center gap-1.5 text-muted-foreground">
				{icon}
				<span className="text-xs">{label}</span>
			</div>
			<div className="font-semibold font-mono text-sm">{value}</div>
		</div>
	)
}

function InfoRow({ label, value }: { label: string; value: string }) {
	return (
		<div className="flex items-center justify-between">
			<span className="text-muted-foreground">{label}:</span>
			<span className="font-mono font-medium">{value}</span>
		</div>
	)
}
