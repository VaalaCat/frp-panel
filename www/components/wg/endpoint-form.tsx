"use client"

import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { toast } from 'sonner'
import { createEndpoint, updateEndpoint } from '@/api/wg'
import { CreateEndpointRequest, UpdateEndpointRequest } from '@/lib/pb/api_wg'

type DisableField = 'host' | 'port'

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
	initial?: { host?: string; port?: number | '' }
	disableFields?: DisableField[]
	submitText?: string
	onSuccess?: () => void
	onError?: (err: Error) => void
}) {
	const { t } = useTranslation()
	const [loading, setLoading] = useState(false)
	const [host, setHost] = useState(initial?.host ?? '')
	const [port, setPort] = useState<number | ''>(initial?.port ?? '')

	useEffect(() => {
		setHost(initial?.host ?? '')
		setPort(initial?.port ?? '')
	}, [initial?.host, initial?.port])

	const disabled = (k: DisableField) => !!disableFields?.includes(k)

	const onSubmit = async () => {
		if (!clientId || !host || !port) return
		setLoading(true)
		try {
			if (mode === 'edit' && endpointId) {
				await updateEndpoint(UpdateEndpointRequest.create({ endpoint: { id: endpointId, clientId, host, port: Number(port) } }))
			} else {
				await createEndpoint(CreateEndpointRequest.create({ endpoint: { clientId, host, port: Number(port) } }))
			}
			toast.success(t('common.success'))
			onSuccess?.()
			if (mode === 'create') {
				setHost('')
				setPort('')
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
				<Label className="block text-sm mb-1">{t('wg.endpointForm.host')}</Label>
				<Input value={host} onChange={(e) => setHost(e.target.value)} placeholder="example.com" disabled={disabled('host')} />
			</div>
			<div>
				<Label className="block text-sm mb-1">{t('wg.endpointForm.port')}</Label>
				<Input value={port} onChange={(e) => setPort(Number(e.target.value) || '')} placeholder="51820" disabled={disabled('port')} />
			</div>
			<div className="flex justify-end gap-2">
				<Button onClick={onSubmit} disabled={loading || !host || !port}>
					{submitText ?? t('wg.endpointForm.submit')}
				</Button>
			</div>
		</div>
	)
}


