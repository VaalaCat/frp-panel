import { useQuery } from '@tanstack/react-query'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { TbDeviceHeartMonitor, TbEngine, TbEngineOff, TbServer2, TbServerBolt, TbServerOff } from 'react-icons/tb'
import { useEffect } from 'react'
import { $platformInfo } from '@/store/user'
import { getPlatformInfo } from '@/api/platform'
import { useTranslation } from 'react-i18next';

export default function PlatformInfo() {
  const { t } = useTranslation();
  const platformInfo = useQuery({
    queryKey: ['platformInfo'],
    queryFn: getPlatformInfo,
  })
  useEffect(() => {
    $platformInfo.set(platformInfo.data)
  }, [platformInfo])
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
      <Card>
        <CardHeader>
          <div className="flex justify-between">
            <h3 className="tracking-tight text-sm font-medium">{t('platform.configuredServers')}</h3>
            <TbServerBolt className="mt-1" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{platformInfo.data?.configuredServerCount} {t('platform.unit')}</div>
          <p className="text-xs text-muted-foreground">{t('platform.menuHint')}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <div className="flex justify-between">
            <h3 className="tracking-tight text-sm font-medium">{t('platform.configuredClients')}</h3>
            <TbEngine className="mt-1" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{platformInfo.data?.configuredClientCount} {t('platform.unit')}</div>
          <p className="text-xs text-muted-foreground">{t('platform.menuHint')}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <div className="flex justify-between">
            <h3 className="tracking-tight text-sm font-medium">{t('platform.unconfiguredServers')}</h3>
            <TbServerOff className="mt-1" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{platformInfo.data?.unconfiguredServerCount} {t('platform.unit')}</div>
          <p className="text-xs text-muted-foreground">{t('platform.menuHint')}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <div className="flex justify-between">
            <h3 className="tracking-tight text-sm font-medium">{t('platform.unconfiguredClients')}</h3>
            <TbEngineOff className="mt-1" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{platformInfo.data?.unconfiguredClientCount} {t('platform.unit')}</div>
          <p className="text-xs text-muted-foreground">{t('platform.menuHint')}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <div className="flex justify-between">
            <h3 className="tracking-tight text-sm font-medium">{t('platform.totalServers')}</h3>
            <TbServer2 className="mt-1" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{platformInfo.data?.totalServerCount} {t('platform.unit')}</div>
          <p className="text-xs text-muted-foreground">{t('platform.menuHint')}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <div className="flex justify-between">
            <h3 className="tracking-tight text-sm font-medium">{t('platform.totalClients')}</h3>
            <TbDeviceHeartMonitor className="mt-1" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{platformInfo.data?.totalClientCount} {t('platform.unit')}</div>
          <p className="text-xs text-muted-foreground">{t('platform.menuHint')}</p>
        </CardContent>
      </Card>
    </div>
  )
}
