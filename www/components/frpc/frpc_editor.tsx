import { Label } from '@radix-ui/react-label'
import { Textarea } from '@/components/ui/textarea'
import { FRPCFormProps } from './frpc_form'
import { useMutation } from '@tanstack/react-query'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import { updateFRPC } from '@/api/frp'
import { useToast } from '@/components/ui/use-toast'
import { RespCode } from '@/lib/pb/common'
import { useTranslation } from 'react-i18next'

export const FRPCEditor: React.FC<FRPCFormProps> = ({ clientID, serverID, client, refetchClient }) => {
  const { t } = useTranslation()
  const { toast } = useToast()

  const [configContent, setConfigContent] = useState<string>('{}')
  const [clientComment, setClientComment] = useState<string>('')
  const updateFrpc = useMutation({ mutationFn: updateFRPC })
  const [editorValue, setEditorValue] = useState<string>('')

  const handleSubmit = async () => {
    try {
      let res = await updateFrpc.mutateAsync({
        clientId: clientID,
        //@ts-ignore
        config: Buffer.from(editorValue),
        serverId: serverID,
        comment: clientComment,
      })
      if (res.status?.code !== RespCode.SUCCESS) {
        toast({ title: t('client.operation.update_failed') })
        return
      }
      toast({ title: t('client.operation.update_success') })
    } catch (error) {
      toast({ title: t('client.operation.update_failed') })
    }
  }

  useEffect(() => {
    refetchClient().then((cliData) => {
      setConfigContent(
        JSON.stringify(
          JSON.parse(
            //@ts-ignore
            cliData?.data?.client?.config == undefined ? '{}' || cliData?.data?.client?.config == '' : cliData?.data?.client?.config,
          ),
          null,
          2,
        ),
      )
      setEditorValue(
        JSON.stringify(
          JSON.parse(
            cliData?.data?.client?.config == undefined || cliData?.data?.client?.config == '' ? '{}' : cliData?.data?.client?.config,
          ),
          null,
          2,
        ),
      )
      setClientComment(cliData?.data?.client?.comment || '')
    }).catch(() => {
      setConfigContent('{}')
      setEditorValue('{}')
      setClientComment('')
    })
  }, [clientID, refetchClient])

  return (
    <div className="grid w-full gap-1.5">
      <Label className="text-sm font-medium">{t('client.editor.comment_title', { id: clientID })}</Label>
      <Textarea
        key={client?.comment}
        placeholder={t('client.editor.comment_placeholder')}
        id="message"
        defaultValue={client?.comment}
        onChange={(e) => setClientComment(e.target.value)}
        className="h-12"
      />
      <Label className="text-sm font-medium">{t('client.editor.config_title', { id: clientID })}</Label>
      <p className="text-sm text-muted-foreground">
        {t('client.editor.config_description')}
      </p>
      <Textarea
        key={configContent}
        placeholder={t('client.editor.config_placeholder')}
        id="message"
        defaultValue={configContent}
        onChange={(e) => setEditorValue(e.target.value)}
        className="h-72"
      />
      <div className="grid grid-cols-2 gap-2 mt-1">
        <Button size="sm" onClick={handleSubmit}>
          {t('common.submit')}
        </Button>
        {/* <Button variant="outline" size="sm" onClick={async () => {
				await refetchClient()
				setConfigContent(client?.client?.config == undefined ? "{}" : client?.client?.config)
			}}>加载服务端配置</Button> */}
      </div>
    </div>
  )
}
