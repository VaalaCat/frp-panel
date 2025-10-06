import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { toast } from 'sonner'
import { createNetwork, updateNetwork } from '@/api/wg'
import { CreateNetworkRequest, UpdateNetworkRequest } from '@/lib/pb/api_wg'
import { AclConfig } from '@/lib/pb/types_wg'

type DisableField = 'name' | 'cidr' | 'acl'

export default function NetworkForm({
  mode,
  networkId,
  initial,
  disableFields,
  submitText,
  onSuccess,
  onError,
}: {
  mode: 'create' | 'edit'
  networkId?: number
  initial?: { name?: string; cidr?: string; aclString?: string }
  disableFields?: DisableField[]
  submitText?: string
  onSuccess?: () => void
  onError?: (err: Error) => void
}) {
  const { t } = useTranslation()
  const [name, setName] = useState(initial?.name ?? '')
  const [cidr, setCidr] = useState(initial?.cidr ?? '')
  const [aclString, setAclString] = useState(initial?.aclString ?? '')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    setName(initial?.name ?? '')
    setCidr(initial?.cidr ?? '')
    setAclString(initial?.aclString ?? '')
  }, [initial?.name, initial?.cidr, initial?.aclString])

  const parseAcl = (text: string): AclConfig | undefined => {
    const trimmed = (text || '').trim()
    if (!trimmed) return undefined
    try {
      const parsed = JSON.parse(trimmed)
      if (parsed && typeof parsed === 'object') {
        return parsed as AclConfig
      }
      toast.error(t('wg.networkForm.invalidAcl') as string)
      return undefined
    } catch (e) {
      toast.error(t('wg.networkForm.invalidAcl') as string)
      return undefined
    }
  }

  const disabled = (key: DisableField) => !!disableFields?.includes(key)

  const onSubmit = async () => {
    if (!name || !cidr) return
    if (mode === 'edit' && !networkId) return
    setLoading(true)
    try {
      const acl = parseAcl(aclString)
      if (mode === 'create') {
        await createNetwork(
          CreateNetworkRequest.create({ network: { name, cidr, acl } })
        )
      } else {
        await updateNetwork(
          UpdateNetworkRequest.create({ network: { id: networkId!, name, cidr, acl } })
        )
      }
      toast.success(t('common.success'))
      onSuccess?.()
      if (mode === 'create') {
        setName('')
        setCidr('')
        setAclString('')
      }
    } catch (e: any) {
      onError?.(e)
      toast.error(e?.message || String(e))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex flex-col gap-3">
      <div className="flex gap-2 items-end">
        <div className="flex-1">
          <Label className="block text-sm mb-1">{t('wg.networkForm.name')}</Label>
          <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="wg-net" disabled={disabled('name')} />
        </div>
        <div className="flex-1">
          <Label className="block text-sm mb-1">{t('wg.networkForm.cidr')}</Label>
          <Input value={cidr} onChange={(e) => setCidr(e.target.value)} placeholder="10.10.0.0/24" disabled={disabled('cidr')} />
        </div>
      </div>
      <div>
        <Label className="block text-sm mb-1">{t('wg.networkForm.acl')}</Label>
        <Textarea
          value={aclString}
          onChange={(e) => setAclString(e.target.value)}
          placeholder={t('wg.networkForm.aclPlaceholder') as string}
          className="h-28"
          disabled={disabled('acl')}
        />
      </div>
      <div className="flex justify-end">
        <Button disabled={loading || !name || !cidr} onClick={onSubmit}>
          {submitText ?? t('wg.networkForm.submit')}
        </Button>
      </div>
    </div>
  )
}

