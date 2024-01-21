import { Label } from '@radix-ui/react-label'
import { Textarea } from './ui/textarea'
import { FRPSFormProps } from './frps_form'
import { Button } from './ui/button'
import { useToast } from './ui/use-toast'
import { useMutation, useQuery } from '@tanstack/react-query'
import { getServer } from '@/api/server'
import { useEffect, useState } from 'react'
import { updateFRPS } from '@/api/frp'
import { RespCode } from '@/lib/pb/common'

export const FRPSEditor: React.FC<FRPSFormProps> = ({ server, serverID }) => {
  const { toast } = useToast()
  const { data: serverResp, refetch: refetchServer } = useQuery({
    queryKey: ['getServer', serverID],
    queryFn: () => {
      return getServer({ serverId: serverID })
    },
  })

  const [configContent, setConfigContent] = useState<string>('{}')
  const updateFrps = useMutation({ mutationFn: updateFRPS })
  const [editorValue, setEditorValue] = useState<string>('')
  const [serverComment, setServerComment] = useState<string>('')

  const handleSubmit = async () => {
    try {
      let res = await updateFrps.mutateAsync({
        serverId: serverID,
        config: Buffer.from(editorValue),
        comment: serverComment,
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
    refetchServer()
    try {
      setConfigContent(
        JSON.stringify(
          JSON.parse(
            serverResp?.server?.config == undefined || serverResp?.server?.config == ''
              ? '{}'
              : serverResp?.server?.config,
          ),
          null,
          2,
        ),
      )
      setEditorValue(
        JSON.stringify(
          JSON.parse(
            serverResp?.server?.config == undefined || serverResp?.server?.config == ''
              ? '{}'
              : serverResp?.server?.config,
          ),
          null,
          2,
        ),
      )
      setServerComment(serverResp?.server?.comment || '')
    } catch (error) {
      setConfigContent('{}')
      setEditorValue('{}')
      setServerComment('')
    }
  }, [serverResp, refetchServer])

  return (
    <div className="grid w-full gap-1.5">
      <Label className="text-sm font-medium">节点 {serverID} 的备注</Label>
      <Textarea
        key={serverResp?.server?.comment}
        placeholder="备注"
        id="message"
        defaultValue={serverResp?.server?.comment}
        onChange={(e) => setServerComment(e.target.value)}
        className="h-12"
      />
      <Label className="text-sm font-medium">节点 {serverID} 配置文件`frps.json`内容</Label>
      <p className="text-sm text-muted-foreground">只需要配置端口和IP等字段，认证信息会由系统补全</p>
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
      </div>
    </div>
  )
}
