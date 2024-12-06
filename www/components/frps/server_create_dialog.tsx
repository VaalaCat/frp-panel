"use client"

import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { initServer } from '@/api/server'
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

export const CreateServerDialog = ({refetchTrigger}: {refetchTrigger?: (randStr: string) => void}) => {
  const { t } = useTranslation()
  const [serverID, setServerID] = useState<string | undefined>()
  const [serverIP, setServerIP] = useState<string | undefined>()
  const newServer = useMutation({
    mutationFn: initServer,
  })
  const { toast } = useToast()

  const handleNewServer = async () => {
    toast({ title: t('server.create.submitting') })
    try {
      let resp = await newServer.mutateAsync({ serverId: serverID, serverIp: serverIP })
      if (resp.status?.code !== RespCode.SUCCESS) {
        toast({ title: t('server.create.error') })
        return
      }
      toast({ title: t('server.create.success') })
      refetchTrigger && refetchTrigger(JSON.stringify(Math.random()))
    } catch (error) {
      toast({ title: t('server.create.error') })
    }
  }

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" size={'sm'}>
          {t('server.create.button')}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t('server.create.title')}</DialogTitle>
          <DialogDescription>{t('server.create.description')}</DialogDescription>
        </DialogHeader>

        <Label>{t('server.create.id')}</Label>
        <Input className="mt-2" value={serverID} onChange={(e) => setServerID(e.target.value)} />
        <Label>{t('server.create.ip')}</Label>
        <Input className="mt-2" value={serverIP} onChange={(e) => setServerIP(e.target.value)} />
        <DialogFooter>
          <Button onClick={handleNewServer}>{t('server.create.submit')}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
