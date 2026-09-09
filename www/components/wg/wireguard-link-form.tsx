"use client"

import { useState, useEffect, useRef, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import { toast } from 'sonner'
import { createWireGuardLink, updateWireGuardLink, getWireGuard } from '@/api/wg'
import { CreateWireGuardLinkRequest, UpdateWireGuardLinkRequest, GetWireGuardRequest } from '@/lib/pb/api_wg'
import { Endpoint, WireGuardLink } from '@/lib/pb/types_wg'
import { WireGuardSelector } from '../base/wireguard-selector'
import { EndpointSelector } from '../base/endpoint-selector'

export default function WireGuardLinkForm({ link, onSuccess, submitText }: { link?: WireGuardLink; onSuccess?: () => void; submitText?: string }) {
	const { t } = useTranslation()
	const [loading, setLoading] = useState(false)
	const [fromId, setFromId] = useState<number | undefined>(link?.fromWireguardId ?? undefined)
	const [toId, setToId] = useState<number | undefined>(link?.toWireguardId ?? undefined)
	const [upBw, setUpBw] = useState<number>(link?.upBandwidthMbps ?? 100)
	const [downBw, setDownBw] = useState<number>(link?.downBandwidthMbps ?? 100)
	const [latency, setLatency] = useState<number>(link?.latencyMs ?? 60)
	const [active, setActive] = useState<boolean>(link?.active ?? true)
	const [toClientId, setToClientId] = useState<string>('')
	const [toEndpointId, setToEndpointId] = useState<number | undefined>(link?.toEndpoint?.id ?? undefined)
	// 记录上一次的 toId，用于检测切换目标节点
	const prevToIdRef = useRef<number | undefined>(link?.toWireguardId ?? undefined)

	// link 变化时同步表单状态，避免使用旧值
	useEffect(() => {
		setFromId(link?.fromWireguardId ?? undefined)
		setToId(link?.toWireguardId ?? undefined)
		setUpBw(link?.upBandwidthMbps ?? 100)
		setDownBw(link?.downBandwidthMbps ?? 100)
		setLatency(link?.latencyMs ?? 60)
		setActive(link?.active ?? true)
		setToEndpointId(link?.toEndpoint?.id ?? undefined)
		prevToIdRef.current = link?.toWireguardId ?? undefined
	}, [link?.id, link?.fromWireguardId, link?.toWireguardId, link?.upBandwidthMbps, link?.downBandwidthMbps, link?.latencyMs, link?.active, link?.toEndpoint?.id])

	// 当 toId 改变时，获取对应的 clientId，并在切换目标时清理旧的 endpoint 选择
	useEffect(() => {
		if (!toId) {
			setToClientId('')
			setToEndpointId(undefined)
			prevToIdRef.current = undefined
			return
		}

		if (prevToIdRef.current && prevToIdRef.current !== toId) {
			setToEndpointId(undefined)
		}
		prevToIdRef.current = toId
		setToClientId('')

		getWireGuard(GetWireGuardRequest.create({ id: toId }))
			.then((resp) => {
				if (resp.wireguardConfig?.clientId) {
					setToClientId(resp.wireguardConfig.clientId)
				}
			})
			.catch((err) => {
				console.error('Failed to get wireguard:', err)
			})
	}, [toId])

	const onSubmit = useCallback(async () => {
		if (!fromId || !toId || fromId === toId) {
			toast.error(t('wg.linkForm.invalid'))
			return
		}
		setLoading(true)
		try {
			const toEndpoint = toEndpointId ? Endpoint.create({ id: toEndpointId }) : undefined
			const linkData = WireGuardLink.create({
				fromWireguardId: fromId,
				toWireguardId: toId,
				upBandwidthMbps: upBw,
				downBandwidthMbps: downBw,
				latencyMs: latency,
				active,
				toEndpoint,
			})

			console.log('linkData', linkData)

			if (link?.id) {
				linkData.id = link.id
				await updateWireGuardLink(
					UpdateWireGuardLinkRequest.create({
						wireguardLink: linkData,
					}),
				)
			} else {
				await createWireGuardLink(
					CreateWireGuardLinkRequest.create({
						wireguardLink: linkData,
					}),
				)
			}
			toast.success(t('common.success'))
			onSuccess?.()
		} catch (e: any) {
			toast.error(e.message)
		} finally {
			setLoading(false)
		}
	}, [fromId, toId, upBw, downBw, latency, active, toEndpointId, link?.id, t, onSuccess])

	return (
		<div className="space-y-3">
			<div className="grid grid-cols-2 gap-2">
				<div>
					<Label className="block text-sm mb-1">{t('wg.link.from')}</Label>
					<WireGuardSelector wireguardID={fromId} setWireguardID={setFromId} />
				</div>
				<div>
					<Label className="block text-sm mb-1">{t('wg.link.to')}</Label>
					<WireGuardSelector wireguardID={toId} setWireguardID={setToId} />
				</div>
			</div>
			<div>
				<Label className="block text-sm mb-1">{t('wg.link.toEndpoint')}</Label>
				{toClientId ? (
					<EndpointSelector clientID={toClientId} endpointID={toEndpointId} setEndpointID={setToEndpointId} />
				) : (
					<div className="text-sm text-muted-foreground p-2 border rounded-md">{t('wg.link.selectToFirst')}</div>
				)}
			</div>
			<div className="grid grid-cols-3 gap-2">
				<div>
					<Label className="block text-sm mb-1">{t('wg.link.up_bw')}</Label>
					<Input value={upBw} onChange={(e) => setUpBw(Number(e.target.value))} />
				</div>
				<div>
					<Label className="block text-sm mb-1">{t('wg.link.down_bw')}</Label>
					<Input value={downBw} onChange={(e) => setDownBw(Number(e.target.value))} />
				</div>
				<div>
					<Label className="block text-sm mb-1">{t('wg.link.latency')}</Label>
					<Input value={latency} onChange={(e) => setLatency(Number(e.target.value))} />
				</div>
			</div>
			<div className="flex items-center gap-2">
				<Checkbox checked={active} onCheckedChange={(v) => setActive(!!v)} />
				<Label>{t('wg.link.active')}</Label>
			</div>
			<div className="flex justify-end gap-2 pt-2">
				<Button onClick={onSubmit} disabled={loading || !fromId || !toId || fromId === toId}>{submitText ?? t('wg.linkForm.submit')}</Button>
			</div>
		</div>
	)
}
