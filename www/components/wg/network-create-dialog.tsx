"use client"

import { useState, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { ScrollArea } from '@/components/ui/scroll-area'
import { toast } from 'sonner'
import { createNetwork } from '@/api/wg'
import { CreateNetworkRequest } from '@/lib/pb/api_wg'
import { AclConfig } from '@/lib/pb/types_wg'
import { Button } from '@/components/ui/button'
import { NetworkBasicFields } from './network-form/basic-section'
import { NetworkAclFields } from './network-form/acl-section'

export function NetworkCreateDialog({
	open,
	onOpenChange,
	onCreated,
}: {
	open: boolean
	onOpenChange: (open: boolean) => void
	onCreated?: () => void
}) {
	const { t } = useTranslation()
	const [activeTab, setActiveTab] = useState<'basic' | 'advanced'>('basic')
	const [name, setName] = useState('')
	const [cidr, setCidr] = useState('')
	const [aclString, setAclString] = useState('')
	const [loading, setLoading] = useState(false)

	const parsedAcl = useMemo(() => {
		const trimmed = aclString.trim()
		if (!trimmed) return undefined
		try {
			return JSON.parse(trimmed) as AclConfig
		} catch {
			return undefined
		}
	}, [aclString])

	const onSubmit = async () => {
		if (!name || !cidr) return
		if (aclString.trim() && !parsedAcl) {
			toast.error(t('wg.networkForm.invalidAcl') as string)
			setActiveTab('advanced')
			return
		}
		setLoading(true)
		try {
			await createNetwork(
				CreateNetworkRequest.create({ network: { name, cidr, acl: parsedAcl } })
			)
			toast.success(t('common.success'))
			onCreated?.()
			onOpenChange(false)
			setName('')
			setCidr('')
			setAclString('')
			setActiveTab('basic')
		} catch (err: any) {
			toast.error(err?.message || String(err))
		} finally {
			setLoading(false)
		}
	}

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent className="sm:max-w-[640px]">
				<DialogHeader>
					<DialogTitle>{t('wg.networkCreate.title')}</DialogTitle>
				</DialogHeader>
				<Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as 'basic' | 'advanced')} className="space-y-4">
					<TabsList className="grid w-full grid-cols-2">
						<TabsTrigger value="basic">{t('wg.networkCreate.tabsBasic')}</TabsTrigger>
						<TabsTrigger value="advanced">{t('wg.networkCreate.tabsAcl')}</TabsTrigger>
					</TabsList>
					<TabsContent value="basic">
						<ScrollArea className="max-h-[55vh] pr-4">
							<NetworkBasicFields name={name} cidr={cidr} onNameChange={setName} onCidrChange={setCidr} />
						</ScrollArea>
					</TabsContent>
					<TabsContent value="advanced">
						<ScrollArea className="max-h-[55vh] pr-4">
							<NetworkAclFields value={aclString} onChange={setAclString} />
						</ScrollArea>
					</TabsContent>
				</Tabs>
				<div className="flex justify-end gap-2">
					<Button variant="outline" onClick={() => onOpenChange(false)} disabled={loading}>
						{t('common.cancel')}
					</Button>
					<Button onClick={onSubmit} disabled={loading || !name || !cidr}>
						{loading ? t('common.loading') : t('wg.networkCreate.submit')}
					</Button>
				</div>
			</DialogContent>
		</Dialog>
	)
}
