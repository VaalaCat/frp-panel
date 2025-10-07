"use client"

import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import { toast } from 'sonner'
import { createWireGuardLink, updateWireGuardLink, getWireGuard } from '@/api/wg'
import { CreateWireGuardLinkRequest, UpdateWireGuardLinkRequest, GetWireGuardRequest } from '@/lib/pb/api_wg'
import { WireGuardLink } from '@/lib/pb/types_wg'
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

	// 当 toId 改变时，获取对应的 clientId
	useEffect(() => {
		if (toId) {
			getWireGuard(GetWireGuardRequest.create({ id: toId }))
				.then((resp) => {
					if (resp.wireguardConfig?.clientId) {
						setToClientId(resp.wireguardConfig.clientId)
					}
				})
				.catch((err) => {
					console.error('Failed to get wireguard:', err)
					setToClientId('')
				})
		} else {
			setToClientId('')
			setToEndpointId(undefined)
		}
	}, [toId])

	const onSubmit = async () => {
		if (!fromId || !toId || fromId === toId) {
			toast.error(t('wg.linkForm.invalid'))
			return
		}
		setLoading(true)
		try {
			const linkData: any = {
				fromWireguardId: fromId,
				toWireguardId: toId,
				upBandwidthMbps: upBw,
				downBandwidthMbps: downBw,
				latencyMs: latency,
				active,
			}
			if (toEndpointId) {
				linkData.toEndpoint = { id: toEndpointId }
			}
			if (link?.id) {
				linkData.id = link.id
				await updateWireGuardLink(UpdateWireGuardLinkRequest.create({
					wireguardLink: WireGuardLink.create(linkData)
				}))
			} else {
				await createWireGuardLink(CreateWireGuardLinkRequest.create({
					wireguardLink: WireGuardLink.create(linkData)
				}))
			}
			toast.success(t('common.success'))
			onSuccess?.()
		} catch (e: any) {
			toast.error(e.message)
		} finally {
			setLoading(false)
		}
	}

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
