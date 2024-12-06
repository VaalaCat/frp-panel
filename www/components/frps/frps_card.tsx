import { useEffect, useState } from 'react'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { getServer } from '@/api/server'
import { useQuery } from '@tanstack/react-query'
import { Switch } from '@/components/ui/switch'
import { FRPSEditor } from './frps_editor'
import FRPSForm from './frps_form'
import { useSearchParams } from 'next/navigation'
import { ServerSelector } from '../base/server-selector'
import { useTranslation } from 'react-i18next';

export interface FRPSFormCardProps {
  serverID?: string
}
export const FRPSFormCard: React.FC<FRPSFormCardProps> = ({ serverID: defaultServerID }: FRPSFormCardProps = {}) => {
  const [advanceMode, setAdvanceMode] = useState<boolean>(false)
  const [serverID, setServerID] = useState<string | undefined>()
  const searchParams = useSearchParams()
  const paramServerID = searchParams.get('serverID')
  const { data: server, refetch: refetchServer } = useQuery({
    queryKey: ['getServer', serverID],
    queryFn: () => {
      return getServer({ serverId: serverID })
    },
  })
  const { t } = useTranslation();

  useEffect(() => {
    if (defaultServerID) {
      setServerID(defaultServerID)
    }
  }, [defaultServerID])

  useEffect(() => {
    if (paramServerID) {
      setServerID(paramServerID)
    }
  }, [paramServerID])

  const handleServerChange = (value: string) => {
    setServerID(value)
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>{t('server.configuration')}</CardTitle>
        <CardDescription>
          <div>
            {t('server.warning.title')}
            <br />{t('server.warning.dockerHint')}
            <br />{t('server.warning.systemdHint')}
          </div>
          <div>
            {t('server.selectHint')}
          </div></CardDescription>
      </CardHeader>
      <CardContent>
        <div className=" flex items-center space-x-4 rounded-md border p-4">
          <div className="flex-1 space-y-1">
            <p className="text-sm font-medium leading-none">{t('server.advancedMode.title')}</p>
            <p className="text-sm text-muted-foreground">{t('server.advancedMode.description')}</p>
          </div>
          <Switch onCheckedChange={setAdvanceMode} />
        </div>
        <div className="flex flex-col w-full pt-2">
          <Label className="text-sm font-medium">{t('server.serverLabel')}</Label>
          <ServerSelector serverID={serverID} setServerID={handleServerChange} onOpenChange={refetchServer} />
        </div>
        {serverID && server && server.server && !advanceMode && <FRPSForm key={serverID} serverID={serverID} server={server.server} />}
        {serverID && server && server.server && advanceMode && (
          <FRPSEditor serverID={serverID} server={server.server} />
        )}
      </CardContent>
      <CardFooter></CardFooter>
    </Card>
  )
}
