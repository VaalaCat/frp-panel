"use client"

import { useTranslation } from 'react-i18next'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'

export function NetworkBasicFields({
	name,
	cidr,
	onNameChange,
	onCidrChange,
}: {
	name: string
	cidr: string
	onNameChange: (value: string) => void
	onCidrChange: (value: string) => void
}) {
	const { t } = useTranslation()

	return (
		<div className="space-y-4">
			<div className="grid gap-4 md:grid-cols-2">
				<div className="space-y-2">
					<Label className="text-sm text-muted-foreground" htmlFor="network-name">
						{t('wg.networkForm.name')}
					</Label>
					<Input
						id="network-name"
						value={name}
						onChange={(e) => onNameChange(e.target.value)}
						placeholder="wg-net"
						autoComplete="off"
					/>
				</div>
				<div className="space-y-2">
					<Label className="text-sm text-muted-foreground" htmlFor="network-cidr">
						{t('wg.networkForm.cidr')}
					</Label>
					<Input
						id="network-cidr"
						value={cidr}
						onChange={(e) => onCidrChange(e.target.value)}
						placeholder="10.10.0.0/24"
						autoComplete="off"
					/>
				</div>
			</div>
			<p className="text-sm text-muted-foreground">
				{t('wg.networkForm.basicHelper')}
			</p>
		</div>
	)
}
