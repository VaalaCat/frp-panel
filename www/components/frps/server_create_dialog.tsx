"use client"

import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { initServer } from '@/api/server'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
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
import { toast } from 'sonner'
import { IsIDValid } from '@/lib/consts'

export const CreateServerDialog = ({ refetchTrigger }: { refetchTrigger?: (randStr: string) => void }) => {
  const { t } = useTranslation()
  const [serverID, setServerID] = useState<string | undefined>()
  const [serverIP, setServerIP] = useState<string | undefined>()
  const newServer = useMutation({
    mutationFn: initServer,
  })

  const handleNewServer = async () => {
    toast(t('server.create.submitting'))
    try {
      let resp = await newServer.mutateAsync({ serverId: serverID, serverIp: serverIP })
      if (resp.status?.code !== RespCode.SUCCESS) {
        toast(t('server.create.error'), {
          description: resp.status?.message,
        })
        return
      }
      toast(t('server.create.success'))
      refetchTrigger && refetchTrigger(JSON.stringify(Math.random()))
    } catch (error) {
      toast(t('server.create.error'), {
        description: JSON.stringify(error),
      })
    }
  }

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline">
          {t('server.create.button')}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t('server.create.title')}</DialogTitle>
          <DialogDescription>{t('server.create.description')}</DialogDescription>
        </DialogHeader>

        <Label>{t('server.create.id')}</Label>
        <Input value={serverID} onChange={(e) => setServerID(e.target.value)} />
        <Label>{t('server.create.ip')}</Label>
        <Input value={serverIP} onChange={(e) => setServerIP(e.target.value)} />
        <DialogFooter>
          <Button onClick={handleNewServer}
            disabled={!IsIDValid(serverID) || !serverIP}
            className='w-full'>{t('server.create.submit')}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
