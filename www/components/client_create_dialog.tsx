import { useState } from 'react'
import { useMutation, useQuery } from '@tanstack/react-query'
import { initClient, listClient } from '@/api/client'
import { Label } from './ui/label'
import { Input } from './ui/input'
import { Button } from './ui/button'
import { useToast } from './ui/use-toast'
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

export const CreateClientDialog = () => {
  const [clientID, setClientID] = useState<string | undefined>()
  const newClient = useMutation({
    mutationFn: initClient,
  })
  const dataQuery = useQuery({
    queryKey: ['listClient', { pageIndex: 0, pageSize: 10 }],
    queryFn: async () => {
      return await listClient({ page: 1, pageSize: 10 })
    },
  })
  const { toast } = useToast()

  const handleNewClient = async () => {
    toast({ title: '已提交创建请求' })
    try {
      let resp = await newClient.mutateAsync({ clientId: clientID })
      if (resp.status?.code !== RespCode.SUCCESS) {
        toast({ title: '创建客户端失败' })
        return
      }
      toast({ title: '创建客户端成功' })
      dataQuery.refetch()
    } catch (error) {
      toast({ title: '创建客户端失败' })
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
          <DialogTitle>新建客户端</DialogTitle>
          <DialogDescription>创建新的客户端用于连接，客户端ID必须唯一</DialogDescription>
        </DialogHeader>

        <Label>客户端ID</Label>
        <Input className="mt-2" value={clientID} onChange={(e) => setClientID(e.target.value)} />
        <DialogFooter>
          <Button onClick={handleNewClient}>创建</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
