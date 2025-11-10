"use client"

import { useTranslation } from 'react-i18next'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import EndpointForm from './endpoint-form'

export default function EndpointEditDialog({
	clientId,
	endpoint,
	children,
	onSaved,
	open,
	onOpenChange,
}: {
	clientId: string
	endpoint?: { id?: number; host?: string; port?: number }
	children?: React.ReactNode
	onSaved?: () => void
	open: boolean
	onOpenChange: (open: boolean) => void
}) {
	const { t } = useTranslation()

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			{children && <DialogTrigger asChild>{children}</DialogTrigger>}
			<DialogContent className="sm:max-w-[480px]">
				<DialogHeader>
					<DialogTitle>{endpoint?.id ? t('wg.endpointEdit.title') : t('wg.endpointCreate.title')}</DialogTitle>
				</DialogHeader>
				<EndpointForm
					mode={endpoint?.id ? 'edit' : 'create'}
					clientId={clientId}
					endpointId={endpoint?.id}
					initial={endpoint}
					submitText={endpoint?.id ? (t('wg.endpointEdit.submit') as string) : (t('wg.endpointCreate.submit') as string)}
					onSuccess={() => { onSaved?.(); onOpenChange(false) }}
				/>
			</DialogContent>
		</Dialog>
	)
}


