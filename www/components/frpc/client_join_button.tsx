import { Button } from '@/components/ui/button'
import React, { useState } from 'react'
import { JoinCommandStr } from '@/lib/consts'
import { useStore } from '@nanostores/react'
import { $platformInfo } from '@/store/user'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { useTranslation } from 'react-i18next'
import { signToken } from '@/api/user'
import { useQuery } from '@tanstack/react-query'
import { toast } from 'sonner'
import { RespCode } from '@/lib/pb/common'

export const ClientJoinButton = () => {
  const platformInfo = useStore($platformInfo)
  const [joinToken, setJoinToken] = useState<undefined | string>(undefined)

  const { t } = useTranslation()

  if (!platformInfo) {
    return (
      <Button variant="link" className="px-0">
      </Button>
    )
  }

  const handleNewToken = async () => {
    try {
      const resp = await signToken({
        expiresIn: BigInt(1000000000),
        permissions: [
          { method: 'POST', path: '/api/v1/client/get', },
          { method: 'POST', path: '/api/v1/client/init', },
        ],
      })
      if (!resp || !resp.status || resp.status.code !== RespCode.SUCCESS) {
        toast.error('server error')
        return
      }
      setJoinToken(resp.token)
    } catch (error) {
      toast.error(JSON.stringify(error))
    }
  }

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="outline">{t('client.join.button')}</Button>
      </PopoverTrigger>
      <PopoverContent className="w-[32rem] max-w-[95vw]">
        <div className="grid gap-4">
          <div className="space-y-2">
            <h4 className="font-medium leading-none">{t('client.join.title')}</h4>
            <p className="text-sm text-muted-foreground">
              {t('client.join.description')} (<a className='text-blue-500' href='https://github.com/VaalaCat/frp-panel/releases' target="_blank" rel="noopener noreferrer">{t('common.download')}</a>)
            </p>
          </div>
          <div className="grid gap-2">
            {joinToken != undefined && <>
              <pre className="bg-muted p-3 rounded-md font-mono text-sm overflow-x-auto whitespace-pre-wrap break-all">
                {JoinCommandStr(platformInfo, joinToken)}
              </pre>
              <Button
                size="sm"
                variant="outline"
                className="w-full"
                onClick={() => {
                  if (joinToken) {
                    navigator.clipboard.writeText(JoinCommandStr(platformInfo, joinToken))
                  }
                }}
                disabled={!platformInfo}
              >
                {t('common.copy')}
              </Button>
            </>
            }
            <Button
              size="sm"
              variant="outline"
              className="w-full"
              onClick={handleNewToken}
              disabled={!platformInfo}
            >
              {t('client.join.sign_token')}
            </Button>
          </div>
        </div>
      </PopoverContent>
    </Popover>
  )
}