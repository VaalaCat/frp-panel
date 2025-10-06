"use client"

import { useTranslation } from 'react-i18next'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import WireGuardLinkForm from './wireguard-link-form'
import { WireGuardLink } from '@/lib/pb/types_wg'

export default function WireGuardLinkEditDialog({
	link,
	children,
	onSaved,
	open,
	onOpenChange,
}: {
	link?: WireGuardLink
	children?: React.ReactNode
	onSaved?: () => void
	open: boolean
	onOpenChange: (open: boolean) => void
}) {
	const { t } = useTranslation()

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			{children && <DialogTrigger asChild>{children}</DialogTrigger>}
			<DialogContent className="sm:max-w-[520px]">
				<DialogHeader>
					<DialogTitle>{link?.id ? t('wg.linkEdit.title') : t('wg.linkCreate.title')}</DialogTitle>
				</DialogHeader>
				<WireGuardLinkForm
					link={link}
					submitText={link?.id ? (t('wg.linkEdit.submit') as string) : (t('wg.linkCreate.submit') as string)}
					onSuccess={() => { onSaved?.(); onOpenChange(false) }}
				/>
			</DialogContent>
		</Dialog>
	)
}


