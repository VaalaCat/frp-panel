'use client'

import React from 'react'
import { Client, Worker } from '@/lib/pb/common'
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { useTranslation } from 'react-i18next'
import { WorkerStatus } from './worker_status'
import { useQuery } from '@tanstack/react-query'
import { getWorkerIngress } from '@/api/worker'
import { InfoIcon } from 'lucide-react'
// import { Textarea } from '@/components/ui/textarea'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

interface WorkerInfoCardProps {
  worker: Worker
  onChange: (worker: Worker) => void
  clients?: Client[]
}

export function WorkerInfoCard({ worker, onChange, clients = [] }: WorkerInfoCardProps) {
  const { t } = useTranslation()

  // 获取 Worker Ingress 用于状态统计
  const { data: ingressResp } = useQuery({
    queryKey: ['getWorkerIngress', worker.workerId],
    queryFn: () => getWorkerIngress({ workerId: worker.workerId || '' }),
    enabled: !!worker.workerId,
  })

  return (
    <Card className="shadow-sm">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center">
            <InfoIcon className="h-5 w-5 mr-2 text-muted-foreground" />
            {t('worker.info.basic_info')}
          </CardTitle>
          <WorkerStatus workerId={worker.workerId || ''} clients={clients} />
        </div>
        <CardDescription className="text-xs">{t('worker.info.info_description')}</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4 pt-0">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label className="text-xs font-medium">{t('worker.info.worker_id')}</Label>
            <div className="flex">
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Input value={worker.workerId || ''} readOnly className="bg-muted font-mono text-sm h-9" />
                  </TooltipTrigger>
                  <TooltipContent>{worker.workerId || ''}</TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </div>
          </div>
          <div className="space-y-2">
            <Label className="text-xs font-medium">{t('worker.info.worker_name')}</Label>
            <div className="flex">
              <Input
                value={worker.name || ''}
                onChange={(e) => onChange({ ...worker, name: e.target.value })}
                className="h-9 text-sm"
                placeholder={t('worker.info.name_placeholder')}
              />
            </div>
          </div>
        </div>

        {/* 暂不实现 */}
        {/* <div className="space-y-2">
          <Label className="text-xs font-medium">{t('worker.info.description')}</Label>
          <Textarea
            value={worker.description || ''}
            onChange={(e) => onChange({ ...worker, description: e.target.value })}
            placeholder={t('worker.info.description_placeholder')}
            className="resize-none h-20 text-sm"
          />
        </div> */}

        <div className="flex items-center justify-between rounded-md bg-muted/50 p-3 border">
          <div className="space-y-1">
            <h4 className="text-sm font-medium">{t('worker.info.resources')}</h4>
            <div className="flex space-x-4 text-xs text-muted-foreground">
              <div>
                {t('worker.info.clients')}: <span className="font-medium">{clients.length}</span>
              </div>
              <div>
                {t('worker.info.ingresses')}:{' '}
                <span className="font-medium">{ingressResp?.proxyConfigs?.length || 0}</span>
              </div>
            </div>
          </div>
          {/* 暂不实现 */}
          {/* <div className="text-xs bg-primary/10 text-primary px-2 py-1 rounded-md font-medium">
            {worker.version || t('worker.info.unknown_version')}
          </div> */}
        </div>
      </CardContent>
    </Card>
  )
}
