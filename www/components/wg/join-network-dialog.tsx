import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { listNetworks, listEndpoints, createWireGuard } from '@/api/wg'
import { CreateWireGuardRequest, ListEndpointsRequest, ListNetworksRequest } from '@/lib/pb/api_wg'
import StringListInput from '@/components/base/list-input'
import { ClientSelector } from '@/components/base/client-selector'

export default function JoinNetworkDialog({ clientId: externalClientId, networkId: externalNetworkId, children, onJoined, initIfName }: { clientId?: string; networkId?: number; children: React.ReactNode; onJoined?: () => void; initIfName?: string }) {
  const { t } = useTranslation()
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [clientId, setClientId] = useState(externalClientId ?? '')
  const [networks, setNetworks] = useState<{ id: number; name: string }[]>([])
  const [endpoints, setEndpoints] = useState<{ id: number; host: string; port: number }[]>([])
  const [networkId, setNetworkId] = useState<number>(externalNetworkId ?? 0)
  const [endpointId, setEndpointId] = useState<string>('')
  const [ifName, setIfName] = useState(initIfName ?? 'frpp0')
  const [localAddr, setLocalAddr] = useState('10.10.0.2')
  const [mtu, setMtu] = useState(1420)
  const [port, setPort] = useState(51820)
  const [tags, setTags] = useState<string[]>([])
  const selectedEndpoint = useMemo(() => endpoints.find((e) => String(e.id) === endpointId), [endpoints, endpointId])

  useEffect(() => {
    if (open) {
      listNetworks(ListNetworksRequest.create({ page: 1, pageSize: 100 })).then((res) => {
        setNetworks((res.networks ?? []).map((n) => ({ id: n.id!, name: n.name! })))
      })
      if (clientId) {
        listEndpoints(ListEndpointsRequest.create({ page: 1, pageSize: 100, clientId })).then((res) => {
          setEndpoints((res.endpoints ?? []).map((e) => ({ id: e.id!, host: e.host!, port: e.port! })))
        })
      } else {
        setEndpoints([])
      }
    }
  }, [open, clientId])

  useEffect(() => {
    if (selectedEndpoint) {
      setPort(selectedEndpoint.port)
    }
  }, [selectedEndpoint])

  const onSubmit = async () => {
    if (!clientId || !networkId || !ifName || !localAddr) return
    setLoading(true)
    try {
      await createWireGuard(CreateWireGuardRequest.create({
        wireguardConfig: {
          clientId,
          networkId: Number(networkId),
          interfaceName: ifName,
          localAddress: localAddr,
          listenPort: port,
          interfaceMtu: mtu,
          advertisedEndpoints: selectedEndpoint ? [{ id: selectedEndpoint.id, host: selectedEndpoint.host, port: selectedEndpoint.port, clientId }] : [],
          tags,
        },
      }))
      onJoined?.()
      setOpen(false)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="sm:max-w-[550px]">
        <DialogHeader>
          <DialogTitle>{t('wg.joinNetwork.label')}</DialogTitle>
        </DialogHeader>
        <div className="space-y-3">
          <div>
            <label className="block text-sm mb-1">{t('wg.selector.client')}</label>
            {externalClientId ? (
              <Input value={externalClientId} disabled className="text-sm" />
            ) : (
              <ClientSelector clientID={clientId} setClientID={(value) => { setClientId(value); setEndpointId('') }} />
            )}
          </div>
          <div>
            <label className="block text-sm mb-1">{t('wg.selector.network')}</label>
            <Select value={String(networkId)} onValueChange={(value) => setNetworkId(Number(value))}>
              <SelectTrigger><SelectValue placeholder={t('wg.selector.network') as string} /></SelectTrigger>
              <SelectContent>
                {networks.map(n => <SelectItem key={n.id} value={String(n.id)}>{n.name}</SelectItem>)}
              </SelectContent>
            </Select>
          </div>
          <div>
            <label className="block text-sm mb-1">{t('wg.selector.endpoint')}</label>
            <Select value={endpointId} onValueChange={setEndpointId}>
              <SelectTrigger><SelectValue placeholder={t('wg.selector.endpoint') as string} /></SelectTrigger>
              <SelectContent>
                {endpoints.map(e => <SelectItem key={e.id} value={String(e.id)}>{e.host}:{e.port}</SelectItem>)}
              </SelectContent>
            </Select>
          </div>
          <div className="grid grid-cols-2 gap-2">
            <div>
              <label className="block text-sm mb-1">{t('wg.interfaceField.name')}</label>
              <Input value={ifName} onChange={(e) => setIfName(e.target.value)} />
            </div>
            <div>
              <label className="block text-sm mb-1">{t('wg.interfaceField.localAddress')}</label>
              <Input value={localAddr} onChange={(e) => setLocalAddr(e.target.value)} />
            </div>
          </div>
          <div className="grid grid-cols-2 gap-2">
            <div>
              <label className="block text-sm mb-1">{t('wg.interfaceField.port')}</label>
              <Input value={port} onChange={(e) => setPort(Number(e.target.value) || 0)} />
            </div>
            <div>
              <label className="block text-sm mb-1">{t('wg.interfaceField.mtu')}</label>
              <Input value={mtu} onChange={(e) => setMtu(Number(e.target.value) || 0)} />
            </div>
          </div>
          <div>
            <label className="block text-sm mb-1">{t('wg.common.tags')}</label>
            <StringListInput value={tags} onChange={setTags} />
          </div>
          <div className="flex justify-end gap-2">
            <Button variant="outline" onClick={() => setOpen(false)}>{t('wg.common.cancel')}</Button>
            <Button onClick={onSubmit} disabled={loading || !clientId || !networkId || !ifName || !localAddr}>{t('wg.common.confirm')}</Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}


