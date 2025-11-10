"use client"

import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { toast } from 'sonner'
import { createEndpoint, updateEndpoint } from '@/api/wg'
import { CreateEndpointRequest, UpdateEndpointRequest } from '@/lib/pb/api_wg'
import { ClientSelector } from '../base/client-selector'

type DisableField = 'host' | 'port' | 'type' | 'uri'

export default function EndpointForm({
	mode,
	clientId,
	endpointId,
	initial,
	disableFields,
	submitText,
	onSuccess,
	onError,
}: {
	mode: 'create' | 'edit'
	clientId: string
	endpointId?: number
	initial?: { host?: string; port?: number | ''; type?: string; uri?: string }
	disableFields?: DisableField[]
	submitText?: string
	onSuccess?: () => void
	onError?: (err: Error) => void
}) {
	const { t } = useTranslation()
	const [loading, setLoading] = useState(false)
	const [host, setHost] = useState(initial?.host ?? '')
	const [port, setPort] = useState<number | ''>(initial?.port ?? '')
	const [type, setType] = useState<string>(initial?.type ?? 'udp')
	const [uri, setUri] = useState<string>(initial?.uri ?? '')

	const [inputClientId, setInputClientId] = useState(clientId)

	useEffect(() => {
		setHost(initial?.host ?? '')
		setPort(initial?.port ?? '')
		setInputClientId(clientId)
		setType(initial?.type ?? 'udp')
		setUri(initial?.uri ?? '')
	}, [initial?.host, initial?.port, initial?.type, initial?.uri])

	const disabled = (k: DisableField) => !!disableFields?.includes(k)

	const onSubmit = async () => {
		if (!inputClientId || !host || !port || !type) {
			toast.error(t('wg.endpointForm.invalid'))
			return
		}
		setLoading(true)
		try {
			if (mode === 'edit' && endpointId) {
				await updateEndpoint(UpdateEndpointRequest.create({ endpoint: { id: endpointId, clientId: inputClientId, host, port: Number(port), type, uri } }))
			} else {
				await createEndpoint(CreateEndpointRequest.create({ endpoint: { clientId: inputClientId, host, port: Number(port), type, uri } }))
			}
			toast.success(t('common.success'))
			onSuccess?.()
			if (mode === 'create') {
				setHost('')
				setPort('')
				setType('udp')
				setUri('')
			}
		} catch (e: any) {
			onError?.(e)
			toast.error(e?.message || String(e))
		} finally {
			setLoading(false)
		}
	}

	return (
		<div className="space-y-3">
			<div>
				<Label className="block text-sm mb-1">{t('wg.endpointForm.clientId')}</Label>
				<ClientSelector clientID={inputClientId} setClientID={(value) => { setInputClientId(value); setHost('') }} />
			</div>
			<div>
				<Label className="block text-sm mb-1">{t('wg.endpointForm.host')}</Label>
				<Input value={host} onChange={(e) => setHost(e.target.value)} placeholder="example.com" disabled={disabled('host')} />
			</div>
			<div>
				<Label className="block text-sm mb-1">{t('wg.endpointForm.port')}</Label>
				<Input value={port} onChange={(e) => setPort(Number(e.target.value) || '')} placeholder="51820" disabled={disabled('port')} />
			</div>
			<div>
				<Label className="block text-sm mb-1">{t('wg.endpointForm.type')}</Label>
				<Input value={type} onChange={(e) => setType(e.target.value)} placeholder="udp" disabled={disabled('type')} />
			</div>
			<div>
				<Label className="block text-sm mb-1">{t('wg.endpointForm.uri')}</Label>
				<Input value={uri} onChange={(e) => setUri(e.target.value)} placeholder="ws://example.com" disabled={disabled('uri')} />
			</div>
			<div className="flex justify-end gap-2">
				<Button onClick={onSubmit} disabled={loading || !inputClientId || !host || !port}>
					{submitText ?? t('wg.endpointForm.submit')}
				</Button>
			</div>
		</div>
	)
}


