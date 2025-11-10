"use client"

import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import StringListInput from '../base/list-input'
import { toast } from 'sonner'
import { listEndpoints, updateWireGuard } from '@/api/wg'
import { ListEndpointsRequest, UpdateWireGuardRequest } from '@/lib/pb/api_wg'
import { Endpoint, WireGuardConfig } from '@/lib/pb/types_wg'

export default function WireGuardForm({
	clientId,
	wg,
	onSuccess,
	submitText,
}: {
	clientId: string
	wg: WireGuardConfig
	onSuccess?: () => void
	submitText?: string
}) {
	const { t } = useTranslation()
	const [loading, setLoading] = useState(false)
	const [ifName, setIfName] = useState(wg.interfaceName ?? '')
	const [localAddr, setLocalAddr] = useState(wg.localAddress ?? '')
	const [mtu, setMtu] = useState<number>(wg.interfaceMtu ?? 1420)
	const [port, setPort] = useState<number>(wg.listenPort ?? 51820)
	const [privKey, setPrivKey] = useState(wg.privateKey ?? '')
	const [endpoints, setEndpoints] = useState<Endpoint[]>([])
	const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set<number>())
	const [tags, setTags] = useState(wg.tags ?? [])
	const [wsListenPort, setWsListenPort] = useState<number>(wg.wsListenPort ?? 0)
	const [useGvisorNet, setUseGvisorNet] = useState<boolean>(wg.useGvisorNet ?? false)

	useEffect(() => {
		listEndpoints(ListEndpointsRequest.create({ page: 1, pageSize: 200, clientId }))
			.then((res) => {
				const eps = (res.endpoints ?? []).map((e) => Endpoint.create({ id: e.id, host: e.host, port: e.port, clientId }))
				setEndpoints(eps)
				const preset = new Set<number>((wg.advertisedEndpoints ?? []).map((e) => e.id || 0).filter((id) => id > 0))
				setSelectedIds(preset)
			})
			.catch((err) => toast.error(err.message))
	}, [clientId, wg.id])

	const onToggle = (id: number, checked: boolean) => {
		setSelectedIds((prev) => {
			const next = new Set(prev)
			if (checked) next.add(id)
			else next.delete(id)
			return next
		})
	}

	const onSubmit = async () => {
		if (!wg.id || !clientId || !ifName || !localAddr || !privKey) return
		setLoading(true)
		try {
			const selected = endpoints.filter((e) => selectedIds.has(e.id || 0))
			await updateWireGuard(
				UpdateWireGuardRequest.create({
					wireguardConfig: WireGuardConfig.create({
						id: wg.id,
						clientId,
						networkId: wg.networkId,
						interfaceName: ifName,
						localAddress: localAddr,
						privateKey: privKey,
						listenPort: port,
						interfaceMtu: mtu,
						dnsServers: wg.dnsServers ?? [],
						advertisedEndpoints: selected,
						wsListenPort: wsListenPort,
						tags,
						useGvisorNet,
					}),
				}),
			)
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
					<Label className="block text-sm mb-1">{t('wg.wireguardForm.interfaceName')}</Label>
					<Input value={ifName} onChange={(e) => setIfName(e.target.value)} />
				</div>
				<div>
					<Label className="block text-sm mb-1">{t('wg.wireguardForm.localAddress')}</Label>
					<Input value={localAddr} onChange={(e) => setLocalAddr(e.target.value)} />
				</div>
			</div>
			<div className="grid grid-cols-2 gap-2">
				<div>
					<Label className="block text-sm mb-1">{t('wg.wireguardForm.port')}</Label>
					<Input value={port} onChange={(e) => setPort(Number(e.target.value) || 0)} />
				</div>
				<div>
					<Label className="block text-sm mb-1">{t('wg.wireguardForm.mtu')}</Label>
					<Input value={mtu} onChange={(e) => setMtu(Number(e.target.value) || 0)} />
				</div>
			</div>
			<div>
				<Label className="block text-sm mb-1">{t('wg.wireguardForm.privateKey')}</Label>
				<Input value={privKey} onChange={(e) => setPrivKey(e.target.value)} placeholder="base64 private key" />
			</div>
			<div>
				<Label className="block text-sm mb-1">{t('wg.wireguardForm.tags')}</Label>
				<StringListInput value={tags} onChange={setTags} />
			</div>
			<div>
				<Label className="block text-sm mb-1">{t('wg.wireguardForm.wsListenPort')}</Label>
				<Input value={wsListenPort} onChange={(e) => setWsListenPort(Number(e.target.value) || 0)} />
			</div>
			<div>
				<Label className="block text-sm mb-1">{t('wg.wireguardForm.useGvisorNet')}</Label>
				<Checkbox checked={useGvisorNet} onCheckedChange={(v) => setUseGvisorNet(!!v)} />
			</div>
			<div className="space-y-2">
				<div className="font-medium text-sm">{t('wg.wireguardForm.selectEndpoint')}</div>
				<div className="max-h-40 overflow-auto border rounded p-2 space-y-2">
					{endpoints.map((e) => (
						<Label key={e.id} className="flex items-center gap-2 text-sm">
							<Checkbox checked={selectedIds.has(e.id || 0)} onCheckedChange={(v) => onToggle(e.id || 0, !!v)} />
							<span>{e.host}:{e.port}</span>
						</Label>
					))}
					{endpoints.length === 0 && <div className="text-sm text-muted-foreground">{t('table.noData')}</div>}
				</div>
			</div>
			<div className="flex justify-end gap-2 pt-2">
				<Button onClick={onSubmit} disabled={loading || !ifName || !localAddr || !privKey}>{submitText ?? t('wg.wireguardForm.submit')}</Button>
			</div>
		</div>
	)
}


