'use client'

import React, { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { ClientSelector } from '../base/client-selector'
import { Cpu, Download, Loader2, Trash } from 'lucide-react'
import { useMutation, useQuery } from '@tanstack/react-query'
import { getWorkerStatus, installWorkerd } from '@/api/worker'
import { Client } from '@/lib/pb/common'
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { Badge } from '@/components/ui/badge'
import { toast } from 'sonner'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

interface ClientDeploymentProps {
  workerId: string
  deployedClientIDs: string[]
  setDeployedClientIDs: (ids: string[]) => void
  clients?: Client[]
}

export function ClientDeployment({
  workerId,
  deployedClientIDs,
  setDeployedClientIDs,
  clients = [],
}: ClientDeploymentProps) {
  const { t } = useTranslation()
  const [dialogOpen, setDialogOpen] = useState(false)
  const [selectedClientId, setSelectedClientId] = useState('')
  const [downloadUrl, setDownloadUrl] = useState('')

  const { data: statusResp } = useQuery({
    queryKey: ['workerStatus', workerId],
    queryFn: () => getWorkerStatus({ workerId }),
    enabled: !!workerId,
    refetchInterval: 10000,
  })

  const installWorkerdMutation = useMutation({
    mutationFn: installWorkerd,
    onSuccess: () => {
      toast.success(t('worker.client_install_workerd.success'))
    },
    onError: () => {
      toast.error(t('worker.client_install_workerd.error'))
    },
  })

  const statusMap = statusResp?.workerStatus || {}

  const handleAddClient = () => {
    if (selectedClientId && !deployedClientIDs.includes(selectedClientId)) {
      setDeployedClientIDs([...deployedClientIDs, selectedClientId])
      setSelectedClientId('')
      setDialogOpen(false)
    }
  }

  const handleRemoveClient = (clientId: string) => {
    setDeployedClientIDs(deployedClientIDs.filter((id) => id !== clientId))
  }

  function getStatusInfo(status?: string): {
    variant: 'outline' | 'default' | 'secondary' | 'destructive'
    text: string
  } {
    if (status === 'running') {
      return { variant: 'default', text: t('worker.status_running') }
    } else if (status === 'stopped') {
      return { variant: 'destructive', text: t('worker.status_stopped') }
    } else if (status === 'error') {
      return { variant: 'secondary', text: t('worker.status_error') }
    } else {
      return { variant: 'outline', text: t('worker.status_unknown') }
    }
  }

  return (
    <Card className="shadow-sm">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center">
            <Cpu className="h-5 w-5 mr-2 text-muted-foreground" />
            {t('worker.deploy.title')}
          </CardTitle>
          <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
            <DialogTrigger asChild>
              <Button size="sm" variant="outline" className="h-8 text-xs">
                {t('worker.deploy.add_client')}
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle>{t('worker.deploy.select_client')}</DialogTitle>
                <DialogDescription>{t('worker.deploy.client_description')}</DialogDescription>
              </DialogHeader>

              <div className="py-4">
                <ClientSelector clientID={selectedClientId} setClientID={setSelectedClientId} />
              </div>

              <DialogFooter>
                <Button onClick={handleAddClient} disabled={!selectedClientId}>
                  {t('common.add')}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <div className="space-y-2">
          {deployedClientIDs.length === 0 ? (
            <div className="text-sm text-muted-foreground flex items-center justify-center py-6 border border-dashed rounded-md">
              {t('worker.deploy.no_clients')}
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
              {deployedClientIDs.map((clientId) => {
                const clientStatus = statusMap[clientId]
                const { variant, text } = getStatusInfo(clientStatus)
                const client = clients.find((c) => c.id === clientId)

                return (
                  <div
                    key={clientId}
                    className="group flex flex-col overflow-hidden rounded-md border hover:border-primary/40 hover:shadow-sm transition-all duration-200"
                  >
                    <div className="flex items-center justify-between bg-muted/30 px-3 py-2 border-b">
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <h3 className="font-semibold text-sm truncate max-w-[200px] md:max-w-[300px]">
                              {client?.originClientId || clientId}
                            </h3>
                          </TooltipTrigger>
                          <TooltipContent>{client?.originClientId || clientId}</TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                      <Badge variant={variant} className="font-normal whitespace-nowrap">
                        {text}
                      </Badge>
                    </div>
                    <div className="px-3 py-2 flex-grow">
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <div className="flex items-center">
                              <span className="text-xs font-medium text-muted-foreground mr-1">ID:</span>
                              <span className="font-mono text-xs truncate max-w-[200px]">{clientId}</span>
                            </div>
                          </TooltipTrigger>
                          <TooltipContent>{clientId}</TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    </div>
                    <div className="bg-muted/10 px-3 py-1.5 flex items-center justify-end border-t">
                      <Dialog>
                        <DialogTrigger asChild>
                          <Button variant="ghost" size="icon" className="h-7 w-7">
                            <Download className="h-4 w-4" />
                          </Button>
                        </DialogTrigger>
                        <DialogContent className="sm:max-w-md max-h-[90vh] overflow-y-auto">
                          <DialogHeader>
                            <DialogTitle>{t('worker.client_install_workerd.title')}</DialogTitle>
                            <DialogDescription>{t('worker.client_install_workerd.description')}</DialogDescription>
                          </DialogHeader>
                          <Label>{t('worker.client_install_workerd.download_url')}</Label>
                          <Input
                            placeholder={t('worker.client_install_workerd.placeholder')}
                            defaultValue={downloadUrl}
                            onChange={(e) => setDownloadUrl(e.target.value)}
                          />
                          <DialogFooter className="pt-4">
                            <Button
                              onClick={() => {
                                installWorkerdMutation.mutate({ clientId })
                              }}
                              disabled={installWorkerdMutation.isPending}
                            >
                              {installWorkerdMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                              {t('worker.client_install_workerd.button')}
                            </Button>
                          </DialogFooter>
                        </DialogContent>
                      </Dialog>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-7 w-7 text-red-500 hover:text-red-600 hover:bg-red-50"
                        onClick={() => handleRemoveClient(clientId)}
                      >
                        <Trash className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                )
              })}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
