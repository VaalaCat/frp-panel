'use client'

import React, { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useSearchParams } from 'next/navigation'
import { useQuery, useMutation } from '@tanstack/react-query'
import { toast } from 'sonner'
import { useTranslation } from 'react-i18next'

import { getWorker, updateWorker } from '@/api/worker'
import { UpdateWorkerRequest } from '@/lib/pb/api_client'
import { Worker } from '@/lib/pb/common'

import { Button } from '@/components/ui/button'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { WorkerInfoCard } from './info_card'
import { WorkerIngress } from './ingress_section'
import { ClientDeployment } from './client_deploy'
import { WorkerCodeEditor } from './code_editor'
import { WorkerTemplateEditor } from './template_editor'

export default function WorkerEdit() {
  const router = useRouter()
  const params = useSearchParams()
  const workerId = params.get('workerId')!
  const { t } = useTranslation()

  // 本地状态
  const [worker, setWorker] = useState<Worker>({} as Worker)
  const [code, setCode] = useState('')
  const [template, setTemplate] = useState('')
  const [deployedClientIDs, setDeployedClientIDs] = useState<string[]>([])

  // 获取 Worker
  const { data: resp, refetch: refetchWorker } = useQuery({
    queryKey: ['getWorker', workerId],
    queryFn: () => getWorker({ workerId }),
    enabled: !!workerId,
  })

  useEffect(() => {
    if (resp?.worker) {
      setWorker(resp.worker)
      setCode(resp.worker.code ?? '')
      setTemplate(resp.worker.configTemplate ?? '')
      // @ts-ignore
      setDeployedClientIDs(resp.clients.map((client) => client.id).filter((id) => id !== undefined) || [])
    }
  }, [resp])

  // 更新 Worker
  const updateMut = useMutation({
    mutationFn: () => {
      const req: UpdateWorkerRequest = {
        clientIds: deployedClientIDs,
        worker: {
          ...worker,
          code,
          configTemplate: template,
        },
      }
      return updateWorker(req)
    },
    onSuccess: () => {
      toast.success(t('worker.edit.save_success'))
      refetchWorker()
    },
    onError: (e) => toast.error(`${t('worker.edit.save_error')}: ${e.message}`),
  })

  return (
    <div className="container p-4 mx-auto space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <h1 className="text-xl md:text-2xl font-semibold flex flex-row gap-1 sm:items-center">
          <span>{t('worker.edit.title')}</span>
          <span className="font-mono text-gray-500">「{worker?.name}」</span>
        </h1>
        <div className="flex gap-2 w-full sm:w-auto">
          <Button variant="outline" onClick={() => router.back()} className="flex-1 sm:flex-none">
            {t('common.back')}
          </Button>
          <Button onClick={() => updateMut.mutate()} disabled={updateMut.isPending} className="flex-1 sm:flex-none">
            {updateMut.isPending ? t('common.saving') : t('common.save')}
          </Button>
        </div>
      </div>

      <div className="space-y-4">
        <Tabs defaultValue="info" className="w-full">
          <TabsList className="w-full">
            <TabsTrigger value="info" className="flex-1">
              {t('worker.edit.info_tab')}
            </TabsTrigger>
            <TabsTrigger value="code" className="flex-1">
              {t('worker.edit.code_tab')}
            </TabsTrigger>
            <TabsTrigger value="template" className="flex-1">
              {t('worker.edit.template_tab')}
            </TabsTrigger>
          </TabsList>

          <TabsContent value="info" className="overflow-hidden space-y-4 mt-4">
            <WorkerInfoCard worker={worker} onChange={setWorker} clients={resp?.clients} />

            {/* <div className="grid grid-cols-1 lg:grid-cols-2 gap-4"> */}
            <ClientDeployment
              workerId={workerId}
              deployedClientIDs={deployedClientIDs}
              setDeployedClientIDs={setDeployedClientIDs}
              clients={resp?.clients}
            />
            <WorkerIngress workerId={workerId} refetchWorker={refetchWorker} clients={resp?.clients} />
            {/* </div> */}
          </TabsContent>

          <TabsContent value="code" className="h-[calc(100vh-240px)] border rounded-md overflow-hidden mt-4">
            <WorkerCodeEditor code={code} onChange={setCode} />
          </TabsContent>

          <TabsContent value="template" className="h-[calc(100vh-240px)] border rounded-md overflow-hidden mt-4">
            <WorkerTemplateEditor content={template} onChange={setTemplate} />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  )
}
