"use client"

import { useTranslation } from 'react-i18next'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import NetworkForm from './network-form'

export default function NetworkEditDialog({
	network,
	children,
	onSaved,
	open,
	onOpenChange,
}: {
	network: { id: number; name: string; cidr: string; aclString?: string }
	children?: React.ReactNode
	onSaved?: () => void
	open: boolean
	onOpenChange: (open: boolean) => void
}) {
	const { t } = useTranslation()

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			{children && <DialogTrigger asChild>{children}</DialogTrigger>}
			<DialogContent>
				<DialogHeader>
					<DialogTitle>{t('wg.networkEdit.title')}</DialogTitle>
				</DialogHeader>
				<NetworkForm
					mode="edit"
					networkId={network.id}
					initial={{ name: network.name, cidr: network.cidr, aclString: network.aclString }}
					disableFields={['cidr']}
					submitText={t('wg.networkEdit.submit') as string}
					onSuccess={() => { onSaved?.() }}
				/>
			</DialogContent>
		</Dialog>
	)
}


