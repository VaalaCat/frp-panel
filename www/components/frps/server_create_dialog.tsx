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

export const CreateServerDialog = ({refetchTrigger}: {refetchTrigger?: (randStr: string) => void}) => {
  const [serverID, setServerID] = useState<string | undefined>()
  const [serverIP, setServerIP] = useState<string | undefined>()
  const newServer = useMutation({
    mutationFn: initServer,
  })
  const { toast } = useToast()

  const handleNewServer = async () => {
    toast({ title: '已提交创建请求' })
    try {
      let resp = await newServer.mutateAsync({ serverId: serverID, serverIp: serverIP })
      if (resp.status?.code !== RespCode.SUCCESS) {
        toast({ title: '创建服务端失败' })
        return
      }
      toast({ title: '创建服务端成功' })
      refetchTrigger && refetchTrigger(JSON.stringify(Math.random()))
    } catch (error) {
      toast({ title: '创建服务端失败' })
    }
  }

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" size={'sm'}>
          新建
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>新建服务端</DialogTitle>
          <DialogDescription>创建新的服务端用于提供服务，服务端ID必须唯一</DialogDescription>
        </DialogHeader>

        <Label>服务端ID</Label>
        <Input className="mt-2" value={serverID} onChange={(e) => setServerID(e.target.value)} />
        <Label>IP地址/域名</Label>
        <Input className="mt-2" value={serverIP} onChange={(e) => setServerIP(e.target.value)} />
        <DialogFooter>
          <Button onClick={handleNewServer}>创建</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
