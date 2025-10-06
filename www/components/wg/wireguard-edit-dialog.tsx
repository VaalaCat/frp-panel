"use client"

import { useTranslation } from 'react-i18next'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { WireGuardConfig } from '@/lib/pb/types_wg'
import WireGuardForm from './wireguard-form'

export default function WireGuardEditDialog({
  clientId,
  wg,
  children,
  onUpdated,
  open,
  onOpenChange,
}: {
  clientId: string
  wg: WireGuardConfig
  children?: React.ReactNode
  onUpdated?: () => void
  open: boolean
  onOpenChange: (open: boolean) => void
}) {
  const { t } = useTranslation()

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      {children && <DialogTrigger asChild>{children}</DialogTrigger>}
      <DialogContent className="sm:max-w-[620px]">
        <DialogHeader>
          <DialogTitle>{t('wg.wireguardEdit.title') || 'Edit'}</DialogTitle>
        </DialogHeader>
        <WireGuardForm clientId={clientId} wg={wg} onSuccess={() => { onUpdated?.(); onOpenChange(false) }} submitText={t('wg.wireguardEdit.submit') as string} />
      </DialogContent>
    </Dialog>
  )
}


