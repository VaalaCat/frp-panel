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
		excellent: { label: t('wg.linkQuality.excellent'), color: 'bg-emerald-500', textColor: 'text-emerald-700' },
		good: { label: t('wg.linkQuality.good'), color: 'bg-green-500', textColor: 'text-green-700' },
		fair: { label: t('wg.linkQuality.fair'), color: 'bg-yellow-500', textColor: 'text-yellow-700' },
		poor: { label: t('wg.linkQuality.poor'), color: 'bg-red-500', textColor: 'text-red-700' },
	}

	const config = qualityConfig[quality]

	return (
		<Card className="h-full shadow-lg border-2">
			<CardHeader className="pb-3">
				<CardTitle className="text-lg flex items-center justify-between">
					<span>{t('wg.topologySidebar.title')}</span>
					<Badge
						variant={link.active ? 'default' : 'secondary'}
						className={cn(link.active && 'bg-emerald-500')}
					>
						{link.active ? t('wg.topologySidebar.statusActive') : t('wg.topologySidebar.statusInactive')}
					</Badge>
				</CardTitle>
			</CardHeader>
			<CardContent className="space-y-4">
				{/* 连接质量 */}
				<div className="space-y-2">
					<div className="text-sm font-medium text-muted-foreground">{t('wg.linkQuality.label')}</div>
					<div className="flex items-center gap-2">
						<div className={cn('w-3 h-3 rounded-full', config.color)} />
						<span className={cn('font-semibold', config.textColor)}>{config.label}</span>
						{isBidirectional && (
							<Badge variant="outline" className="ml-auto">
								<Activity className="h-3 w-3 mr-1" />
								{t('wg.linkQuality.bidirectional')}
							</Badge>
						)}
					</div>
				</div>

				<Separator />

				{/* 性能指标 */}
				<div className="grid grid-cols-2 gap-3">
					<MetricCard
						icon={<Zap className="h-4 w-4" />}
						label={t('wg.topologySidebar.latency')}
						value={`${link.latencyMs}ms`}
						variant={link.latencyMs < 50 ? 'good' : link.latencyMs < 100 ? 'normal' : 'bad'}
					/>
					<MetricCard
						icon={<TrendingUp className="h-4 w-4" />}
						label={t('wg.topologySidebar.bandwidthUp')}
						value={`${link.upBandwidthMbps}Mbps`}
						variant="normal"
					/>
					<MetricCard
						icon={<TrendingDown className="h-4 w-4" />}
						label={t('wg.topologySidebar.bandwidthDown')}
						value={`${link.downBandwidthMbps}Mbps`}
						variant="normal"
					/>
					<MetricCard
						icon={<Activity className="h-4 w-4" />}
						label={t('wg.runtime.status')}
						value={link.active ? t('wg.runtime.status_ok') : t('wg.runtime.status_unknown')}
						variant={link.active ? 'good' : 'bad'}
					/>
				</div>

				<Separator />

				{/* 端点信息 */}
				{link.toEndpoint && (
					<div className="space-y-2">
						<div className="text-sm font-medium text-muted-foreground flex items-center gap-2">
							<MapPin className="h-4 w-4" />
							{t('wg.topologySidebar.endpoint')}
						</div>
						<div className="bg-muted/50 rounded-lg p-3 space-y-1">
							<div className="flex items-center justify-between text-sm">
								<span className="text-muted-foreground">{t('wg.endpointForm.host')}</span>
								<span className="font-mono font-medium">{link.toEndpoint.host}</span>
							</div>
							<div className="flex items-center justify-between text-sm">
								<span className="text-muted-foreground">{t('wg.endpointForm.port')}</span>
								<span className="font-mono font-medium">{link.toEndpoint.port}</span>
							</div>
						</div>
					</div>
				)}

				{/* 连接信息 */}
				<div className="space-y-2">
					<div className="text-sm font-medium text-muted-foreground">{t('wg.linkInfo.label')}</div>
					<div className="space-y-2 text-sm">
						<InfoRow label={t('wg.linkInfo.id')} value={`#${link.id}`} />
						<InfoRow label={t('wg.topologySidebar.from')} value={`WG #${link.fromWireguardId}`} />
						<InfoRow label={t('wg.topologySidebar.to')} value={`WG #${link.toWireguardId}`} />
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
