'use client'

import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import type { WGEdge } from './types'
import { useTranslation } from 'react-i18next'
import { Activity, ArrowUp, ArrowDown, Zap, Link2 } from 'lucide-react'

export default function TopologySidebar({ selectedEdge }: { selectedEdge?: WGEdge }) {
	const { t } = useTranslation()
	return (
		<Card className="w-full">
			<CardHeader className="pb-3">
				<CardTitle className="flex items-center gap-2 text-base">
					<Link2 className="h-4 w-4" />
					{t('wg.topologySidebar.title')}
				</CardTitle>
			</CardHeader>
			<CardContent className="space-y-3 text-sm">
				{selectedEdge ? (
					<div className="space-y-3">
						<div className="flex items-center gap-2 pb-2 border-b">
							<Badge variant={selectedEdge.data?.original?.active ? "default" : "secondary"}>
								{selectedEdge.data?.original?.active ? t('wg.topologySidebar.statusActive') : t('wg.topologySidebar.statusInactive')}
							</Badge>
						</div>

						<div className="space-y-2">
							<div className="flex items-start gap-2">
								<span className="text-muted-foreground min-w-[80px]">{t('wg.topologySidebar.from')}:</span>
								<span className="font-mono text-xs break-all">{String(selectedEdge.source)}</span>
							</div>
							<div className="flex items-start gap-2">
								<span className="text-muted-foreground min-w-[80px]">{t('wg.topologySidebar.to')}:</span>
								<span className="font-mono text-xs break-all">{String(selectedEdge.target)}</span>
							</div>
						</div>

						{selectedEdge.data && (
							<div className="space-y-2 pt-2 border-t">
								<div className="flex items-center justify-between p-2 bg-muted/50 rounded">
									<div className="flex items-center gap-2">
										<Zap className="h-4 w-4 text-yellow-500" />
										<span className="text-muted-foreground">{t('wg.topologySidebar.latency')}</span>
									</div>
									<span className="font-semibold">{selectedEdge.data.original.latencyMs}ms</span>
								</div>

								<div className="flex items-center justify-between p-2 bg-muted/50 rounded">
									<div className="flex items-center gap-2">
										<ArrowUp className="h-4 w-4 text-blue-500" />
										<span className="text-muted-foreground">{t('wg.topologySidebar.bandwidthUp')}</span>
									</div>
									<span className="font-semibold">{selectedEdge.data.original.upBandwidthMbps} Mbps</span>
								</div>

								<div className="flex items-center justify-between p-2 bg-muted/50 rounded">
									<div className="flex items-center gap-2">
										<ArrowDown className="h-4 w-4 text-green-500" />
										<span className="text-muted-foreground">{t('wg.topologySidebar.bandwidthDown')}</span>
									</div>
									<span className="font-semibold">{selectedEdge.data.original.downBandwidthMbps} Mbps</span>
								</div>
							</div>
						)}
					</div>
				) : (
					<div className="flex flex-col items-center justify-center py-8 text-muted-foreground">
						<Activity className="h-8 w-8 mb-2 opacity-50" />
						<p className="text-center">{t('wg.topologySidebar.empty')}</p>
					</div>
				)}
			</CardContent>
		</Card>
	)
}


