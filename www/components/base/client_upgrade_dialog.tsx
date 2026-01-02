'use client'

import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { UpgradeFrppRequest } from '@/lib/pb/api_client'
import { upgradeFrpp } from '@/api/client'

import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

export interface ClientUpgradeDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  clientId: string
  defaultVersion?: string
  defaultGithubProxy?: string
  defaultUseGithubProxy?: boolean
  defaultServiceName?: string
  onDispatched?: () => void
}

export function ClientUpgradeDialog({
  open,
  onOpenChange,
  clientId,
  defaultVersion = 'latest',
  defaultGithubProxy,
  defaultUseGithubProxy = true,
  defaultServiceName = 'frpp',
  onDispatched,
}: ClientUpgradeDialogProps) {
  const { t } = useTranslation()

  const [version, setVersion] = useState(defaultVersion)
  const [downloadUrl, setDownloadUrl] = useState('')
  const [httpProxy, setHttpProxy] = useState('')

  const upgradeMutation = useMutation({
    mutationFn: upgradeFrpp,
    onSuccess: () => {
      toast.success(t('client.upgrade.dispatched'))
      onOpenChange(false)
      onDispatched?.()
    },
    onError: (e: any) => {
      toast.error(t('client.upgrade.failed'), { description: e.message })
    },
  })

  const onSubmit = () => {
    upgradeMutation.mutate(
      UpgradeFrppRequest.create({
        clientIds: [clientId],
        version: version || 'latest',
        downloadUrl: downloadUrl || undefined,
        useGithubProxy: defaultUseGithubProxy,
        githubProxy: defaultGithubProxy || undefined,
        httpProxy: httpProxy || undefined,
        backup: true,
        serviceName: defaultServiceName,
        restartService: true,
        workdir: undefined,
        serviceArgs: [],
      }),
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t('client.upgrade.title')}</DialogTitle>
          <DialogDescription>
            <p className="text-destructive">{clientId + ' ' + t('client.upgrade.warning')}</p>
          </DialogDescription>
        </DialogHeader>

        <div className="grid gap-3 py-2">
          <Label>{t('client.upgrade.version')}</Label>
          <Input value={version} onChange={(e) => setVersion(e.target.value)} placeholder="latest" />

          <Label>{t('client.upgrade.download_url')}</Label>
          <Input
            value={downloadUrl}
            onChange={(e) => setDownloadUrl(e.target.value)}
            placeholder={t('client.upgrade.download_url_placeholder')}
          />

          <Label>{t('client.upgrade.http_proxy')}</Label>
          <Input value={httpProxy} onChange={(e) => setHttpProxy(e.target.value)} placeholder="http://127.0.0.1:7890" />
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            {t('common.cancel')}
          </Button>
          <Button onClick={onSubmit} disabled={upgradeMutation.isPending}>
            {upgradeMutation.isPending ? t('common.loading') : t('client.upgrade.confirm')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}


