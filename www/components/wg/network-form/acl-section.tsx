"use client"

import { useTranslation } from 'react-i18next'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

export function NetworkAclFields({ value, onChange }: { value: string; onChange: (value: string) => void }) {
	const { t } = useTranslation()

	return (
		<div className="space-y-3">
			<div className="space-y-2">
				<Label className="text-sm text-muted-foreground" htmlFor="network-acl">
					{t('wg.networkForm.acl')}
				</Label>
				<Textarea
					id="network-acl"
					value={value}
					onChange={(e) => onChange(e.target.value)}
					placeholder={t('wg.networkForm.aclPlaceholder') as string}
					className="min-h-[220px]"
				/>
			</div>
			<div className="rounded-md border bg-muted/50 p-3 text-xs text-muted-foreground">
				{t('wg.networkForm.aclHelper')}
			</div>
		</div>
	)
}
