"use client"

import i18n from '@/lib/i18n'
import { useState } from 'react'
import { useMutation, useQuery } from '@tanstack/react-query'
import { initClient, listClient } from '@/api/client'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/use-toast'
import { RespCode } from '@/lib/pb/common'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { useTranslation } from 'react-i18next'

export const CreateClientDialog = ({refetchTrigger}: {refetchTrigger?: (randStr: string) => void}) => {
  const { t } = useTranslation()
  const [clientID, setClientID] = useState<string | undefined>()
  const newClient = useMutation({
    mutationFn: initClient,
  })
  const { toast } = useToast()

  const handleNewClient = async () => {
    toast({ title: t('client.create.submitting') })
    try {
      let resp = await newClient.mutateAsync({ clientId: clientID })
      if (resp.status?.code !== RespCode.SUCCESS) {
        toast({ title: t('client.create.error') })
        return
      }
      toast({ title: t('client.create.success') })
      refetchTrigger && refetchTrigger(JSON.stringify(Math.random()))
    } catch (error) {
      toast({ title: t('client.create.error') })
    }
  }

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" size={'sm'}>
          {t('client.create.button')}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t('client.create.title')}</DialogTitle>
          <DialogDescription>{t('client.create.description')}</DialogDescription>
        </DialogHeader>

        <Label>{t('client.create.id')}</Label>
        <Input className="mt-2" value={clientID} onChange={(e) => setClientID(e.target.value)} />
        <DialogFooter>
          <Button onClick={handleNewClient}>{t('client.create.submit')}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
