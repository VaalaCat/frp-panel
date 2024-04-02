import { Label } from '@radix-ui/react-label'
import { Textarea } from './ui/textarea'
import { FRPCFormProps } from './frpc_form'
import { getClient } from '@/api/client'
import { useMutation, useQuery } from '@tanstack/react-query'
import { useEffect, useState } from 'react'
import { Button } from './ui/button'
import { updateFRPC } from '@/api/frp'
import { useToast } from './ui/use-toast'
import { RespCode } from '@/lib/pb/common'

export const FRPCEditor: React.FC<FRPCFormProps> = ({ clientID, serverID }) => {
  const { toast } = useToast()
  const { data: client, refetch: refetchClient } = useQuery({
    queryKey: ['getClient', clientID],
    queryFn: () => {
      return getClient({ clientId: clientID })
    },
  })

  const [configContent, setConfigContent] = useState<string>('{}')
  const [clientComment, setClientComment] = useState<string>('')
  const updateFrpc = useMutation({ mutationFn: updateFRPC })
  const [editorValue, setEditorValue] = useState<string>('')

  const handleSubmit = async () => {
    try {
      let res = await updateFrpc.mutateAsync({
        clientId: clientID,
        config: Buffer.from(editorValue),
        serverId: serverID,
        comment: clientComment,
      })
      if (res.status?.code !== RespCode.SUCCESS) {
        toast({ title: '更新失败' })
        return
      }
      toast({ title: '更新成功' })
    } catch (error) {
      toast({ title: '更新失败' })
    }
  }

  useEffect(() => {
    refetchClient()
    try {
      setConfigContent(
        JSON.stringify(
          JSON.parse(
            client?.client?.config == undefined ? '{}' || client?.client?.config == '' : client?.client?.config,
          ),
          null,
          2,
        ),
      )
      setEditorValue(
        JSON.stringify(
          JSON.parse(
            client?.client?.config == undefined || client?.client?.config == '' ? '{}' : client?.client?.config,
          ),
          null,
          2,
        ),
      )
      setClientComment(client?.client?.comment || '')
    } catch (error) {
      setConfigContent('{}')
      setEditorValue('{}')
      setClientComment('')
    }
  }, [client, refetchClient])

  return (
    <div className="grid w-full gap-1.5">
      <Label className="text-sm font-medium">节点 {clientID} 的备注</Label>
      <Textarea
        key={client?.client?.comment}
        placeholder="备注"
        id="message"
        defaultValue={client?.client?.comment}
        onChange={(e) => setClientComment(e.target.value)}
        className="h-12"
      />
      <Label className="text-sm font-medium">客户端 {clientID} 配置文件`frpc.json`内容</Label>
      <p className="text-sm text-muted-foreground">
        只需要配置proxies和visitors字段，认证信息和服务器连接信息会由系统补全
      </p>
      <Textarea
        key={configContent}
        placeholder="配置文件内容"
        id="message"
        defaultValue={configContent}
        onChange={(e) => setEditorValue(e.target.value)}
        className="h-72"
      />
      <div className="grid grid-cols-2 gap-2 mt-1">
        <Button size="sm" onClick={handleSubmit}>
          提交
        </Button>
        {/* <Button variant="outline" size="sm" onClick={async () => {
				await refetchClient()
				setConfigContent(client?.client?.config == undefined ? "{}" : client?.client?.config)
			}}>加载服务端配置</Button> */}
      </div>
    </div>
  )
}
