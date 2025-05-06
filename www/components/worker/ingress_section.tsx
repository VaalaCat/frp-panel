'use client'

import React, { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQuery, useMutation } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { createWorkerIngress, getWorkerIngress } from '@/api/worker'
import { Client, ProxyConfig } from '@/lib/pb/common'
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { ProxyConfigMutateForm } from '../proxy/mutate_proxy_config'
import { Loader2, Settings, Trash, Network } from 'lucide-react'
import { TypedProxyConfig } from '@/types/proxy'
import { ClientSelector } from '../base/client-selector'
import { ServerSelector } from '../base/server-selector'
import { Label } from '@/components/ui/label'
import { deleteProxyConfig, getProxyConfig } from '@/api/proxy'
import { useStore } from '@nanostores/react'
import { $proxyTableRefetchTrigger } from '@/store/refetch-trigger'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { Badge } from '@/components/ui/badge'
import { ServerSideVisitPreview, VisitPreview } from '../base/visit-preview'
import { getServer } from '@/api/server'

interface WorkerIngressProps {
  workerId: string
  refetchWorker: () => void
  clients?: Client[]
}

// 创建 Worker Ingress 表单组件
const CreateWorkerIngressForm = ({
  workerId,
  onSuccess,
  clients,
}: {
  workerId: string
  onSuccess: () => void
  clients?: Client[]
}) => {
  const { t } = useTranslation()
  const [clientId, setClientId] = useState<string>('')
  const [serverId, setServerId] = useState<string>('')

  const createWorkerIngressMutation = useMutation({
    mutationFn: createWorkerIngress,
    onSuccess: () => {
      toast.success(t('worker.ingress.create_success'))
      onSuccess()
    },
    onError: (error) => {
      toast.error(t('worker.ingress.create_failed'), {
        description: error instanceof Error ? error.message : String(error),
      })
    },
  })

  const handleSubmit = () => {
    if (!clientId) {
      toast.error(t('worker.ingress.client_required'))
      return
    }

    if (!serverId) {
      toast.error(t('worker.ingress.server_required'))
      return
    }

    createWorkerIngressMutation.mutate({
      clientId,
      serverId,
      workerId,
    })
  }

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="client-select">{t('worker.ingress.select_client')}</Label>
        <ClientSelector setClientID={setClientId} clientID={clientId} clients={clients} />
      </div>

      <div className="space-y-2">
        <Label htmlFor="server-select">{t('worker.ingress.select_server')}</Label>
        <ServerSelector setServerID={setServerId} serverID={serverId} />
      </div>

      <DialogFooter>
        <Button disabled={!clientId || !serverId || createWorkerIngressMutation.isPending} onClick={handleSubmit}>
          {createWorkerIngressMutation.isPending ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              {t('worker.ingress.creating')}
            </>
          ) : (
            t('worker.ingress.create_submit')
          )}
        </Button>
      </DialogFooter>
    </div>
  )
}

function ProxyStatusBadge({
  clientId,
  serverId,
  proxyName,
}: {
  clientId?: string
  serverId?: string
  proxyName?: string
}) {
  const { t } = useTranslation()
  const refetchTrigger = useStore($proxyTableRefetchTrigger)

  const { data } = useQuery({
    queryKey: ['getProxyConfig', clientId, serverId, proxyName, refetchTrigger],
    queryFn: () => {
      return getProxyConfig({
        clientId,
        serverId,
        name: proxyName,
      })
    },
    enabled: !!clientId && !!serverId && !!proxyName,
    refetchInterval: 10000,
  })

  function getStatusInfo(status: string): {
    color: string
    text: string
    variant: 'outline' | 'default' | 'secondary' | 'destructive'
  } {
    switch (status) {
      case 'new':
        return { color: 'bg-blue-100 border-blue-400', text: t('status.new'), variant: 'secondary' }
      case 'wait start':
        return { color: 'bg-yellow-100 border-yellow-400', text: t('status.wait_start'), variant: 'secondary' }
      case 'start error':
        return { color: 'bg-red-100 border-red-400', text: t('status.start_error'), variant: 'destructive' }
      case 'running':
        return { color: 'bg-green-100 border-green-400', text: t('status.running'), variant: 'default' }
      case 'check failed':
        return { color: 'bg-orange-100 border-orange-400', text: t('status.check_failed'), variant: 'secondary' }
      case 'error':
        return { color: 'bg-red-100 border-red-400', text: t('status.error'), variant: 'destructive' }
      default:
        return { color: 'bg-gray-100 border-gray-400', text: t('status.unknown'), variant: 'outline' }
    }
  }

  const status = data?.workingStatus?.status || 'unknown'
  const { text, variant } = getStatusInfo(status)

  return (
    <Badge variant={variant} className={'font-normal whitespace-nowrap'}>
      {text}
    </Badge>
  )
}

export function WorkerIngress({ workerId, refetchWorker, clients }: WorkerIngressProps) {
  const { t } = useTranslation()

  // 获取 Worker Ingress
  const { data: ingresses, refetch: refetchIngresses } = useQuery({
    queryKey: ['getWorkerIngress', workerId],
    queryFn: () => getWorkerIngress({ workerId }),
    enabled: !!workerId,
  })

  const handleIngressCreated = () => {
    refetchIngresses()
    refetchWorker()
  }

  return (
    <Card className="shadow-sm">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center">
            <Network className="h-5 w-5 mr-2 text-muted-foreground" />
            {t('worker.ingress.title')}
          </CardTitle>
          <Dialog>
            <DialogTrigger asChild>
              <Button size="sm" variant="outline" className="h-8 text-xs">
                {t('worker.ingress.create')}
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-md max-h-[90vh] overflow-y-auto">
              <DialogHeader>
                <DialogTitle>{t('worker.ingress.create_title')}</DialogTitle>
                <DialogDescription>{t('worker.ingress.create_description')}</DialogDescription>
              </DialogHeader>
              <CreateWorkerIngressForm workerId={workerId} onSuccess={handleIngressCreated} clients={clients} />
            </DialogContent>
          </Dialog>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <div className="space-y-2">
          {!ingresses || !ingresses.proxyConfigs || ingresses.proxyConfigs.length === 0 ? (
            <div className="text-sm text-muted-foreground flex items-center justify-center py-6 border border-dashed rounded-md">
              {t('worker.ingress.no_ingress')}
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-2 rounded-md">
              {ingresses.proxyConfigs.map((ingress: ProxyConfig) => (
                <div
                  key={ingress.id}
                  className="group overflow-hidden rounded-md border hover:border-primary/40 hover:shadow-sm transition-all duration-200 flex flex-col"
                >
                  <div className="flex items-center justify-between bg-muted/30 px-3 py-2 border-b">
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <h3 className="font-semibold text-sm truncate max-w-[200px] md:max-w-[300px]">
                            {ingress.name}
                          </h3>
                        </TooltipTrigger>
                        <TooltipContent>{ingress.name}</TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                    <ProxyStatusBadge
                      clientId={ingress.clientId}
                      serverId={ingress.serverId}
                      proxyName={ingress.name}
                    />
                  </div>
                  <div className="px-3 py-2 flex-grow">
                    <div className="flex flex-col space-y-1">
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <div className="flex items-center text-center">
                              <span className="text-xs font-medium text-muted-foreground mr-1 font-mono">Server:</span>
                              <span className="font-mono text-xs items-center hide-scroll-bar">
                                <WorkerIngressPreview
                                  serverId={ingress.serverId}
                                  proxyCfg={JSON.parse(ingress.config || '{}') as TypedProxyConfig} />
                              </span>
                            </div>
                          </TooltipTrigger>
                          <TooltipContent>
                            <div>
                              ServerID: {ingress.serverId}
                            </div>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>

                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <div className="flex items-center">
                              <span className="text-xs font-medium text-muted-foreground mr-1 font-mono">Client:</span>
                              <span className="font-mono text-xs truncate max-w-[150px] md:max-w-[200px]">
                                {ingress.originClientId}
                              </span>
                            </div>
                          </TooltipTrigger>
                          <TooltipContent>{ingress.originClientId}</TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    </div>
                  </div>
                  <div className="bg-muted/10 px-3 py-1.5 flex items-center justify-end gap-1 border-t">
                    <Dialog>
                      <DialogTrigger asChild>
                        <Button variant="ghost" size="icon" className="h-7 w-7">
                          <Settings className="h-4 w-4" />
                        </Button>
                      </DialogTrigger>
                      <DialogContent className="sm:max-w-md max-h-[90vh] overflow-y-auto">
                        <DialogHeader>
                          <DialogTitle>{t('worker.ingress.edit')}</DialogTitle>
                        </DialogHeader>
                        <ProxyConfigMutateForm
                          disableChangeProxyName
                          defaultProxyConfig={JSON.parse(ingress.config || '{}') as TypedProxyConfig}
                          overwrite={true}
                          defaultOriginalProxyConfig={ingress}
                          onSuccess={handleIngressCreated}
                        />
                      </DialogContent>
                    </Dialog>
                    <Dialog>
                      <DialogTrigger asChild>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-7 w-7 text-red-500 hover:text-red-600 hover:bg-red-50"
                        >
                          <Trash className="h-4 w-4" />
                        </Button>
                      </DialogTrigger>
                      <DialogContent className="sm:max-w-md">
                        <DialogHeader>
                          <DialogTitle>{t('worker.ingress.delete.title')}</DialogTitle>
                          <DialogDescription>{t('worker.ingress.delete.description')}</DialogDescription>
                        </DialogHeader>
                        <IngressDeleteForm
                          clientId={ingress.clientId}
                          serverId={ingress.serverId}
                          proxyName={ingress.name}
                          onSuccess={handleIngressCreated}
                        />
                      </DialogContent>
                    </Dialog>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}

export const IngressDeleteForm = ({
  clientId,
  serverId,
  proxyName,
  onSuccess,
}: {
  clientId?: string
  serverId?: string
  proxyName?: string
  onSuccess?: () => void
}) => {
  const { t } = useTranslation()
  // 删除 Ingress 对应的 ProxyConfig
  const deleteProxyConfigMutation = useMutation({
    mutationFn: deleteProxyConfig,
    onSuccess: () => {
      onSuccess?.()
      toast(t('worker.ingress.delete_success'), {
        description: t('worker.ingress.delete_description', { proxyName }),
      })
    },
  })
  return (
    <DialogFooter className="pt-4">
      <Button
        variant={'destructive'}
        onClick={() =>
          deleteProxyConfigMutation.mutate({
            clientId,
            serverId,
            name: proxyName,
          })
        }
        disabled={deleteProxyConfigMutation.isPending}
      >
        {deleteProxyConfigMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        {t('worker.ingress.delete.button')}
      </Button>
    </DialogFooter>
  )
}

const WorkerIngressPreview = ({ serverId, proxyCfg }: { serverId?: string; proxyCfg: TypedProxyConfig }) => {
  const { data: getServerResp } = useQuery({
    queryKey: ['getServer', serverId],
    queryFn: () => {
      return getServer({ serverId: serverId })
    },
  })

  return (<div className='font-mono text-xs'>
    <ServerSideVisitPreview server={getServerResp?.server || { frpsUrls: [] }} typedProxyConfig={proxyCfg} withIcon={false} />
  </div>)
}