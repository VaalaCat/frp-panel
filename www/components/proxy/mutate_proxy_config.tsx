'use client'

import { useEffect, useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { useTranslation } from 'react-i18next'
import { ServerSelector } from '../base/server-selector'
import { ClientSelector } from '../base/client-selector'
import { TypedProxyForm } from '../frpc/proxy_form'
import { ProxyType, TypedProxyConfig } from '@/types/proxy'
import { BaseSelector } from '../base/selector'
import { createProxyConfig } from '@/api/proxy'
import { ClientConfig } from '@/types/client'
import { ObjToUint8Array } from '@/lib/utils'
import { VisitPreview } from '../base/visit-preview'
import { ProxyConfig, Server } from '@/lib/pb/common'
import { TypedProxyConfigValid } from '@/lib/consts'
import { toast } from 'sonner'
import { $proxyTableRefetchTrigger } from '@/store/refetch-trigger'

export type ProxyConfigMutateDialogProps = {
  overwrite?: boolean
  defaultProxyConfig?: TypedProxyConfig
  defaultOriginalProxyConfig?: ProxyConfig
  disableChangeProxyName?: boolean
  onSuccess?: () => void
}

export const ProxyConfigMutateDialog = ({ ...props }: ProxyConfigMutateDialogProps) => {
  const { t } = useTranslation()

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" className="w-fit">
          {t('proxy.config.create')}
        </Button>
      </DialogTrigger>
      <DialogContent className="max-h-[90vh] overflow-auto">
        <DialogHeader>
          <DialogTitle>{t('proxy.config.create_proxy')}</DialogTitle>
          <DialogDescription>{t('proxy.config.create_proxy_description')}</DialogDescription>
        </DialogHeader>
        <ProxyConfigMutateForm {...props} />
      </DialogContent>
    </Dialog>
  )
}

export const ProxyConfigMutateForm = ({
  overwrite,
  defaultProxyConfig,
  defaultOriginalProxyConfig,
  disableChangeProxyName,
  onSuccess,
}: ProxyConfigMutateDialogProps) => {
  const { t } = useTranslation()
  const [newClientID, setNewClientID] = useState<string | undefined>()
  const [newServerID, setNewServerID] = useState<string | undefined>()
  const [proxyConfigs, setProxyConfigs] = useState<TypedProxyConfig[]>([])
  const [proxyName, setProxyName] = useState<string | undefined>('')
  const [proxyType, setProxyType] = useState<ProxyType>('http')
  const [selectedServer, setSelectedServer] = useState<Server | undefined>()
  const supportedProxyTypes: ProxyType[] = ['http', 'tcp', 'udp']

  const createProxyConfigMutation = useMutation({
    mutationKey: ['createProxyConfig', newClientID, newServerID],
    mutationFn: () =>
      createProxyConfig({
        clientId: newClientID!,
        serverId: newServerID!,
        config: ObjToUint8Array({
          proxies: proxyConfigs,
        } as ClientConfig),
        overwrite,
      }),
    onSuccess: () => {
      toast(t('proxy.config.create_success'))
      $proxyTableRefetchTrigger.set(Math.random())
      onSuccess?.()
    },
    onError: (e) => {
      toast(t('proxy.config.create_failed'), {
        description: JSON.stringify(e),
      })
      $proxyTableRefetchTrigger.set(Math.random())
    },
  })

  useEffect(() => {
    if (proxyName && proxyType) {
      setProxyConfigs([{...defaultProxyConfig, name: proxyName, type: proxyType }])
    }
  }, [proxyName, proxyType])

  useEffect(() => {
    if (defaultProxyConfig && defaultOriginalProxyConfig) {
      setProxyConfigs([defaultProxyConfig])
      setProxyType(defaultProxyConfig.type)
      setProxyName(defaultProxyConfig.name)
      setNewClientID(defaultOriginalProxyConfig.originClientId)
      setNewServerID(defaultOriginalProxyConfig.serverId)
    }
  }, [defaultProxyConfig, defaultOriginalProxyConfig])

  return (
    <>
      <Label>{t('proxy.config.select_server')} </Label>
      <ServerSelector setServerID={setNewServerID} serverID={newServerID} setServer={setSelectedServer} />
      <Label>{t('proxy.config.select_client')} </Label>
      <ClientSelector setClientID={setNewClientID} clientID={newClientID} />
      <Label>{t('proxy.config.select_proxy_type')} </Label>
      <BaseSelector
        dataList={supportedProxyTypes.map((type) => ({ value: type, label: type }))}
        value={proxyType}
        setValue={(value) => {
          setProxyType(value as ProxyType)
        }}
      />
      {proxyConfigs &&
        selectedServer &&
        proxyConfigs.length > 0 &&
        proxyConfigs[0] &&
        TypedProxyConfigValid(proxyConfigs[0]) && (
          <div className="flex flex-row w-full overflow-auto">
            <div className="flex flex-col">
              <VisitPreview server={selectedServer} typedProxyConfig={proxyConfigs[0]} />
            </div>
          </div>
        )}
      <Label>{t('proxy.config.proxy_name')} </Label>
      <Input
        className="text-sm"
        defaultValue={proxyName}
        onChange={(e) => setProxyName(e.target.value)}
        disabled={disableChangeProxyName}
      />
      {proxyName && newClientID && newServerID && (
        <TypedProxyForm
          serverID={newServerID}
          clientID={newClientID}
          proxyName={proxyName}
          defaultProxyConfig={proxyConfigs && proxyConfigs.length > 0 ? proxyConfigs[0] : undefined}
          clientProxyConfigs={proxyConfigs}
          setClientProxyConfigs={setProxyConfigs}
          enablePreview={false}
        />
      )}
      <Button
        disabled={!TypedProxyConfigValid(proxyConfigs[0])}
        onClick={() => {
          if (!TypedProxyConfigValid(proxyConfigs[0])) {
            toast(t('proxy.config.invalid_config'))
            return
          }
          createProxyConfigMutation.mutate()
        }}
      >
        {t('proxy.config.submit')}
      </Button>
    </>
  )
}
