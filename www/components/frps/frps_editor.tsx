import { Label } from '@radix-ui/react-label'
import { Textarea } from '@/components/ui/textarea'
import { FRPSFormProps } from './frps_form'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/use-toast'
import { useMutation, useQuery } from '@tanstack/react-query'
import { getServer } from '@/api/server'
import { useEffect, useState } from 'react'
import { updateFRPS } from '@/api/frp'
import { RespCode } from '@/lib/pb/common'
import { useTranslation } from 'react-i18next'

export const FRPSEditor: React.FC<FRPSFormProps> = ({ server, serverID }) => {
  const { t } = useTranslation()
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
        //@ts-ignore
        config: Buffer.from(editorValue),
        comment: serverComment,
      })
      if (res.status?.code !== RespCode.SUCCESS) {
        toast({ title: t('server.operation.update_failed') })
        return
      }
      toast({ title: t('server.operation.update_success') })
    } catch (error) {
      toast({ title: t('server.operation.update_failed') })
    }
    refetchServer()
  }

  useEffect(() => {
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
  }, [serverResp])

  return (
    <div className="grid w-full gap-1.5">
      <Label className="text-sm font-medium">{t('server.editor.comment', { id: serverID })}</Label>
      <Textarea
        key={serverResp?.server?.comment}
        placeholder={t('server.editor.comment_placeholder')}
        id="message"
        defaultValue={serverResp?.server?.comment}
        onChange={(e) => setServerComment(e.target.value)}
        className="h-12"
      />
      <Label className="text-sm font-medium">{t('server.editor.config_title', { id: serverID })}</Label>
      <p className="text-sm text-muted-foreground">{t('server.editor.config_description')}</p>
      <Textarea
        key={configContent}
        placeholder={t('server.editor.config_placeholder')}
        id="message"
        defaultValue={configContent}
        onChange={(e) => setEditorValue(e.target.value)}
        className="h-72"
      />
      <div className="grid grid-cols-2 gap-2 mt-1">
        <Button size="sm" onClick={handleSubmit}>
          {t('common.submit')}
        </Button>
      </div>
    </div>
  )
}
