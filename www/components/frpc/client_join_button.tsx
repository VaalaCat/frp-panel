import { Button } from '@/components/ui/button'
import React from 'react'
import { JoinCommandStr } from '@/lib/consts'
import { useStore } from '@nanostores/react'
import { $platformInfo, $token } from '@/store/user'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { useTranslation } from 'react-i18next'

export const ClientJoinButton = () => {
  const platformInfo = useStore($platformInfo)
  const token = useStore($token)

  const { t } = useTranslation()

  if (!platformInfo) {
    return (
      <Button variant="link" className="px-0">
      </Button>
    )
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
          {token != undefined && <div className="grid gap-2">
            <pre className="bg-muted p-3 rounded-md font-mono text-sm overflow-x-auto whitespace-pre-wrap break-all">
              {JoinCommandStr(platformInfo, token)}
            </pre>
            <Button
              size="sm"
              variant="outline"
              className="w-full"
              onClick={() => {
                if (token) {
                  navigator.clipboard.writeText(JoinCommandStr(platformInfo, token))
                }
              }}
              disabled={!platformInfo}
            >
              {t('common.copy')}
            </Button>
          </div>}
        </div>
      </PopoverContent>
    </Popover>
  )
}